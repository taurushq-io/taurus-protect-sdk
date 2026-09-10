# SDK Alignment Report

**Date**: 2026-09-03
**Branch**: `feature/sdk_gov_rules`
**Reference**: Java for naming and API shape; the shared vector files for byte-level behaviour

## How alignment is enforced

Alignment is checked by gates, not by inspection. Four of them, all consumed by every SDK's
unit suite or `build.sh`:

| Gate | File / command | Covers | Enforced |
|---|---|---|---|
| Cell wire parity | `scripts/resources/governance-cell-vectors.json` (39 vectors) | Each typed `RuleCell` encodes to exact recorded bytes and decodes back | all 4 unit suites |
| Lossless / non-canonical parity | `scripts/resources/governance-lossless-vectors.json` (8 vectors) | Schema-newer and deliberately non-canonical wire forms survive a round trip | all 4 unit suites |
| Authorization-error parsing | `scripts/resources/authorization-error-vectors.json` (8 vectors) | The 403 role-list parser agrees on every server wording | all 4 unit suites |
| API-surface parity | `scripts/api-surface/diff.py` | Every service exists in all four SDKs; method-count deltas reported | `build.sh docs` |
| Docs match code | `build.sh docs --check` | No documented method that does not exist; generated index current | `build.sh docs --check` |

The lossless vectors and the API-surface differ are **new in this pass**. Before it, the
eight lossless base64 strings were hand-copied into four separate test suites with nothing
comparing the copies, and no gate covered the service surface at all — which is how three
capability gaps below reached only one or two SDKs unnoticed.

The differ ships a self-test (`scripts/api-surface/selftest.sh`) that runs it against a
fixture with a deliberate gap. A differ that always prints "OK" is worse than no differ,
because its report reads as evidence of parity.

## Status

Verified by running each gate, not asserted:

| SDK | Services | Public service methods | Unit tests | Other gates |
|---|---|---|---|---|
| Go | 43 | 191 | 1050 test funcs, 7 packages | `go vet`, `GOARCH=386` build, `golangci-lint` 0 issues |
| Java | 43 | 184 | 1273 | Checkstyle 0 violations |
| Python | 43 | 202 | 1419 | new-file flake8 clean (repo lint is a known-red baseline) |
| TypeScript | 43 | 207 | 1722 in 103 suites | `tsc --noEmit`, ESLint 0 |

Service counts come from `scripts/resources/api-surface.<lang>.json`, regenerated from
source by `build.sh docs`. The "43 services (38 core + 5 TaurusNetwork)" claim is now
measured rather than repeated.

`tg-protect-mcpd`, which consumes the Go SDK through a local-path `replace`, builds and
passes its own suite against this branch.

## What this pass changed

### Correctness on the governance sign-then-approve path

1. **Container JSON bridge was non-deterministic in Java, Python and TypeScript.** The
   container carries `map<string, bytes> properties` at five levels and map order is
   unspecified unless forced, so the same rules JSON encoded to different bytes per run and
   per SDK — and these are the bytes SuperAdmins sign. The typed encoder had been fixed
   earlier; the newer JSON bridge missed it. Each SDK now routes through the deterministic
   serializer it already had (Python `deterministic=True`, Java `deterministicBytes`, TS
   `sortedMap`). Pinned by a test in all four that feeds JSON with reverse-ordered keys and
   asserts one shared expected base64 — feeding back already-sorted JSON passes even with a
   broken encoder, which is why the first version of that test was worthless.
2. **TypeScript released no credential material on `close()`.** It overwrote
   `config.apiSecret`, which is `undefined` under the supported `credentials:` path. The
   secret lives in the auth-middleware closure, which JavaScript cannot zero, so `close()`
   now empties the middleware array `Configuration` shares, making it unreachable and
   GC-eligible.
3. **Go's documented error handling could never match.** Three same-named `APIError`
   families existed and only `service.APIError` was live, so every documented
   `errors.Is(err, protect.ErrNotFound)` / `protect.IsAPIError(...)` branch was dead code.
   Consolidated to one type: `protect` now type-aliases the live types and the sentinels
   match by status code. `IsRetryable` / `IsClientError` / `IsServerError` /
   `SuggestedRetryDelay` / `RetryAfter` moved onto the live type — Go had none of them,
   while the report used to claim all four SDKs implemented `isRetryable()`.
4. **TypeScript let callers bypass governance verification entirely.** `ProtectClient`
   exposed `governanceRulesApi`, whose raw `ruleServiceGetRules` skips SuperAdmin signature
   verification and whose `ruleServiceUpdateRulesProposal` accepts an arbitrary base64 blob.
   Go makes this impossible (`internal/openapi`); the getter is now private in TypeScript
   too, so the typed, verified path is the only way in.

### Capability gaps closed

| Gap | Was missing from | Now |
|---|---|---|
| `getPublicKeys` + a `SuperAdminPublicKey` model | TypeScript (absent entirely, though mcpd calls it) | present in all 4 |
| `RuleCell` codec entry points | Go (unexported — an external Go caller could not encode or decode a cell at all) | `RuleCellToBytes` / `RuleCellFromBytes` / `CellFamily` exported |
| `UpdateTransactionsEnabled` (the transactions kill switch) | Java, Python, TypeScript | present in all 4 |
| `blockchain` and `network` transaction filters (list + export) | Java, Python | present in all 4 |
| Helper-level `VerifyGovernanceRules` | Go | present in all 4 |
| `verifyHashCoverage` | Java | present in all 4 |
| Single-key signature verify | TypeScript | present in all 4 |
| `ParseRequiredRoles` exported | Go (unexported), Python/TS barrels | exported in all 4 |

### Shape and naming aligned to Java

- `RawCell.payload` — was `Bytes` in Go and `bytes` in TypeScript.
- Go `RuleCell.Kind()` — Go had no discriminator, forcing a type switch where the other
  three expose `kind`.
- Go `cellFamily(RawCell)` returns the cell's column type; it returned `""`, discarding the
  one piece of grammar information a raw cell still carries.
- TypeScript governance reads renamed `get`/`getById`/`getProposal`/`getHistory` →
  `getRules`/`getRulesById`/`getRulesProposal`/`getRulesHistory`.
- `verifyGovernanceRules(rules)` returns the verified rules in all four; Go returned only an
  error and TypeScript returned `void`. Java keeps its two-argument overload but the
  single-argument form is the cross-SDK shape, because passing a threshold that differs from
  the configured one verifies against something the rest of the service does not enforce.
- Python `comment` is required on `approve_rules_proposal` / `reject_rules_proposal`; it
  defaulted to `""`, so a Python caller could approve a governance change with no rationale
  recorded while the same call in any other SDK would not compile.

### Test-suite parity

- The eight lossless vectors moved into the shared file (above).
- **TypeScript never exercised the encoder** for the `RuleSource` lossless scenario, because
  `ruleSourceToBytes` was module-private. A regression dropping `RuleSource.raw` on encode
  passed there while breaking the other three SDKs' byte parity. Exported and asserted.
- TypeScript's nested-unknown-field assertion checked presence, not content; an empty bag
  satisfied it. Now asserts non-emptiness like its peers.
- Java packed all five `minValidSignatures` distinct-key cases into one `@Test`; JUnit fails
  fast, so one break hid four. Split into five.
- Java's round-trip fixture used `"Ethereum"` where the others use `"ETH"`, so the four
  "rich containers" were not the same fixture.

### Documentation

- **~100 documented methods did not exist** across the four `SERVICES.md` files (Go 21,
  Java 26, TypeScript 4 plus a TaurusNetwork section naming the wrong receiver types on
  every signature). All corrected against the extracted surface, and `build.sh docs --check`
  now fails on any new one. Hand-written prose was preserved; only a generated method index
  was added.
- Every "Basic Usage" example in the Python README omitted the mandatory
  `super_admin_keys_pem` and therefore raised — and implied the keys were optional.
- The Go README examples did not compile: 3-value returns from 2-value functions, `me.Name`
  on a result type with no `Name`, and `Limit:` on option structs whose field is `PageSize`.
  Verified now by compiling them.
- `taurus-protect-sdk-go/docs/SDK_OVERVIEW.md` listed the SuperAdmin key options as
  "Optional" against code that refuses to construct a client without them.
- Verification **Step 6** ("parse from the verified payload") was documented only by Java
  while Go, Python and TypeScript all claimed a 6-step flow and stopped at 5 — the step that
  actually stops an attacker-supplied label reaching the caller.
- `docs/CONCEPTS.md` documented a `RequestMetadata.payload` object that all four SDKs omit
  on purpose.
- The **container JSON bridge** was documented nowhere despite being the surface the MCP
  tooling is built on. `docs/CONCEPTS.md` now covers it, including the two things a caller
  must know: it is untyped (the `RuleCell` union is bypassed) and it is **not** lossless
  (protobuf-JSON cannot represent unknown fields).
- TypeScript's README and `SDK_OVERVIEW.md` listed twelve services as having no client
  getter. All twelve have one; the suggested workaround skipped DTO mapping and, on the
  security paths, signature verification.
- Corrected counts: Go 1.24 (root README said 1.21), 36 `RuleCell` variants (was "~37"),
  61 generated OpenAPI APIs (this report said 56), Go 47 mapper / 42 model files (said
  83/46), Python 25 pledge and 15 sharing models (said 26/14), TypeScript 38 service getters
  (said 26).

## Verification alignment pass (2026-09-04)

A pass over the verification logic for the four entity families — **addresses, requests,
whitelisted addresses, whitelisted assets** — and the documentation describing them. The gates
above covered governance-cell bytes, error parsing and API surface; none of them touched the
verification flows, which is where the drift below accumulated.

### One security defect, not drift

**TypeScript whitelisted assets were never verified.**
`WhitelistedAssetService` built a verifier and used it at exactly one site
(`getWithVerification`). `get()`, `list()` and `getEnvelope()` called the mapper directly — no
metadata hash check, no SuperAdmin signatures, no hash coverage, no thresholds — while the
mapper's own comments read *"Parse from verified payload"*. `client.ts` handed users that
service through `withVerification()`. Go, Java and Python all verified unconditionally, and
TypeScript's *address* service did too: asset-only, TypeScript-only.

Proven rather than assumed: reverting `get()` to the pre-fix form makes the new regression
test resolve to `{"contractAddress":"0xATTACKER",…}` instead of throwing. Five existing
service tests had been *pinning* the vulnerability — one fed `blockchain:'UNVERIFIED'` in the
DTO alongside a fabricated payload and asserted the payload won, with no signature anywhere.

### Fixed, per invariant

| Invariant | Was wrong in | Fix |
|---|---|---|
| Every asset read path verifies | **TypeScript** | `get`/`list`/`getEnvelope` route through the verifier; the unverified mapper is deleted |
| Hash coverage is constant-time | **Java** | `SignatureVerifier.verifyHashCoverage` had **zero production callers**; four `List.contains` sites now route through it, plus a new `containsHash` for the per-signature check |
| Asset flow has a step 6 | **Go** | new `ParseWhitelistedAssetFromJSON` + `Name`/`Symbol`/`ContractAddress`/`Decimals`/`TokenID` on the model, which did not exist |
| Threshold core rejects a non-positive threshold | **Go** | `len(signers) < 0` is false, so it passed with zero signers |
| Payload accessors are gated on verification | Java, Python, TypeScript | parse only inside verify; a distinct `UnverifiedMetadataError` on early access |
| Step 4 returns the hash it matched | **TypeScript** | a legacy-signed asset passed step 4 then failed step 5 |
| Wildcard matching is case-insensitive | **Go** | `"any"` selected a different rule tier than in the peers |
| Verification does not mutate its input | **Java** | both services called `setHash(legacyHash)` on the caller's envelope |
| `minValidSignatures <= #keys` at construction | Java, TypeScript | a 3-of-2 client verified nothing and failed only at call time |
| Envelope getters verify | **TypeScript** | address `getEnvelope` returned an unverified envelope |
| Address verifier throws rather than returning booleans | **TypeScript** | no empty-address check; batch form returned `boolean[]` a caller could ignore |
| `hashVerified` set on the single-get path | Python, TypeScript | verified, then returned with the flag still false |
| Metadata survives a hash-less response | **TypeScript** | the mapper dropped the whole object, so the payload never reached verification |
| Explicit "hash present, payload missing" branch | Go, Java | Go's guard was dead code — `CalculateHexHash` always returns 64 chars |
| Rule selection keyed on the SIGNED payload | Go, Java, TypeScript | the DTO is free-floating; an empty blockchain selected the broadest tier |
| Populated group with `minimumSignatures = 0` | **all four** | no post-loop threshold check exists, so a zero silently meant 1-of-N |
| Lenient lists with exclusions reported | Java, Python, TypeScript | one bad row denied access to every good one |
| Rows returned but none surviving | **all four** | an empty page was indistinguishable from an empty whitelist |
| One crypto-error predicate | **TypeScript** | five copies of `message.includes('key')`, which swallowed real `TypeError`s |
| Hash-in-list logic | **all four** | 14 implementations collapsed toward two shared functions per SDK |

### Corrected during the pass

Two decisions were reversed by evidence found while implementing them:

1. **Dropping the `threshold` fallback would have weakened Python and TypeScript.** There is
   no post-loop `validCount >= minSigs` check, so `minimumSignatures = 0` on a populated group
   demands *one* signature, not zero. Dropping the fallback would have taken those two SDKs
   from enforcing `threshold` (say 2) down to 1-of-N. Replaced by rejecting the zero outright.
2. **Requiring `network` in the signed payload would reject correctly-signed addresses.**
   Governance rules carry a per-rule `includeNetworkInPayload` flag, and the captured
   production payload omits `network`. The chain is required from the payload; the network
   falls back to the DTO only when the payload genuinely carries none.

### New gates

Two, both loaded by all four suites and both asserting the per-section counts the file declares.

1. **`verification-behaviour-vectors.json`** — 27 crypto-free vectors over `resolveRuleKey`,
   `verifyHashCoverage` and `containsHash`.
2. **`verification-signed-fixtures.json`** — 10 vectors with real signatures over the
   SuperAdmin threshold, closing the gap the root `CLAUDE.md` recorded as *"No automated
   cross-SDK gate covers this"* for the five distinct-key cases. **Public keys only**: the gate
   verifies, it never signs, so no private key material is committed. Go regenerates it with
   fresh keys each run.

Both verified real: deleting either hard-fails every SDK rather than skipping. The signed
fixtures were additionally checked against a deliberate regression — reverting Go to count
signature *entries* instead of distinct keys turns two vectors red — and each loader asserts
the file contains a discriminating case, so a fixture set that could not catch that bug fails
on its own terms.

### Documentation

- `docs/INTEGRITY_VERIFICATION.md` gained **Whitelisted Asset**, **Address** and a real
  **Request Metadata** section (it was four bullets), a rule-selection section, and assets in
  the "What Gets Verified" table.
- **`WHITELISTED_ASSET_VERIFICATION.md` now exists in all four SDKs** — it existed in none.
- Corrected: the Go doc claiming a client works without SuperAdmin keys (it refuses to
  construct); the TypeScript doc offering "Basic Retrieval (No Verification)" (that path
  verifies); Python's "5-STEP" banner over a 6-step body; the TypeScript CLAUDE.md's 3-argument
  `withVerification`; the "strict mode" list claim in five CLAUDE.md files; and the unscoped
  "never break early" rule, which all four SDKs correctly violate on threshold success.
- Hash values in whitelist errors are now **permitted and consistent** across all four; Python
  was the outlier and now matches.

### Also removed

A 76 KB `tg-validatord` Go test file committed inside the Java SDK's test tree
(`client/src/test/java/.../client/model/whitelist_test.go`), importing
`tg-validatord/internal/...` — internal server source in a repo carrying a LICENSE, a
DISCLAIMER and a public Postman collection.

---

## Remaining known differences

Accepted, with the reason. Anything here is a decision, not an oversight — the API-surface
differ's `aliases.json` points back at this table.

| Difference | Detail | Why accepted |
|---|---|---|
| Java inlines whitelist verification in its services | Peers use a `helper/` verifier class | Only the crypto helpers were relocated this pass; the full extraction is tracked in `TODOS.md`, because that layout is *why* the non-constant-time comparisons hid there |
| `WhitelistedContractService` vs `ContractWhitelistingService` | Go uses the former | Go idiom; recorded in `aliases.json` |
| `WebhookCallsService` plural | Java only | Cosmetic; recorded in `aliases.json` |
| Deprecated flat api-key params absent in Go | Java/Python/TS keep them deprecated | Go's `WithCredentials` was retyped; its only callers were SDK tests and mcpd |
| `Request.id` type | Java `long`, Go `string`, Python `Optional[int]`, TS `number` | Go uses `string` for all IDs to avoid precision loss |
| Go `Request.Status` is a plain `string` | Others use a typed enum | Go idiom |
| Go-only `WithLogger` / `protect.Logger` / `protect.Field` | Go gained an injectable logger; Java uses `java.util.logging`, Python `logging`, TypeScript `console.warn` | The SDK had no report channel at all, and list-path integrity exclusions must not be silent. Each SDK uses its own idiom rather than a common abstraction none of them asked for |
| Go-only removal of generated `TgvalidatordMetadata.payload` | Java types it `Object`, Python is dynamic, TypeScript `any` — all accept the array the wire sends | Only `openapi-generator -g go` maps the untyped schema to `map[string]interface{}`, which fails to decode. Verified per SDK, not assumed |

| Whitelisted-asset envelope getter absent in Python | Go/Java/TS expose one | **Not a gap**: Python's `WhitelistedAsset` already carries the whole envelope (`metadata`, `rules_container`, `rules_signatures`, `signed_contract_address`). Go/Java/TS split asset-from-envelope; Python merged them, so a caller already has raw access |

| Helper decomposition differs | `CheckHashesSignature` (Go, Java), whitelist-integrity helper (Python, Java) | Same capability reached by different helper layout |
| `UPDATE_GOVERNANCE_CELL_VECTORS=1` regenerates the cell vectors from Go | No equivalent elsewhere | This is the sanctioned regeneration path, not a hole |

### Live divergence, not yet resolved

**`RequestStatus` values differ.** Java has 40, Go/Python/TypeScript have the same 38.
Java-only: `FAST_APPROVED`, `INVALID`, `NEW`. Present in the other three only: `UNKNOWN`
(an SDK-side sentinel). Neither `scripts/resources/swagger/apis.swagger.json` nor the proto
schema declares the three Java-only values, so there is no authoritative source to align
to: adding them to three SDKs on one SDK's unsourced list would be guesswork, and dropping
them from Java could break a caller matching a status the server really sends. This needs a
server-side answer about which values `RequestStatus` can actually take. Earlier versions of
this report claimed "43 values matching the Java source of truth" in all four, and listed 38
names, 23 of which were not in `RequestStatus.java`.

## Deferred work

Captured with full context in [`TODOS.md`](../TODOS.md):

1. Seal TypeScript's remaining low-level `*Api` getters (same bypass class as the governance
   one, for whitelist/request/change writes).
2. Java's TaurusNetwork services are systematically thinner than the peers — the differ
   reports lending 4 vs 14/16/14, pledge 3 vs 14, settlement 2 vs 6, sharing 1 vs 6. This is
   net-new Java implementation, not alignment of existing surface.
3. Promote the rich-container fixture to a shared vector now that the Java fixture matches.
4. Revive or retire the Python lint and Java PMD gates, both abandoned-red.

## Verification

```bash
# Go
cd taurus-protect-sdk-go && go build ./... && go vet ./pkg/... && go test ./pkg/... \
  && GOARCH=386 go build ./pkg/... && golangci-lint run ./pkg/...

# Java (JDK 17 on PATH; no `clean` — it fails deleting proto/target)
export JAVA_HOME=/workspace/java/jdk-17.0.19+10 && export PATH="$JAVA_HOME/bin:$PATH"
cd taurus-protect-sdk-java && mvn compile -o -pl client \
  && mvn test -o -pl client -Dspotbugs.skip=true -Dpmd.skip=true -Dcheckstyle.skip=true \
  && mvn checkstyle:check -pl client

# Python
cd taurus-protect-sdk-python && .venv/bin/python -m pytest tests/unit -o addopts= -q

# TypeScript
cd taurus-protect-sdk-typescript && ./build.sh unit && npx eslint src --ext .ts

# Cross-SDK gates
./scripts/api-surface/generate.sh all check    # surface parity + docs match code
./scripts/api-surface/selftest.sh              # proves the differ detects a gap
```

To confirm a shared vector gate is real rather than green, remove
`scripts/resources/governance-lossless-vectors.json` and re-run any SDK's suite: all four
hard-fail on a missing file, and none skips.

---

*Last updated: 2026-09-04*
