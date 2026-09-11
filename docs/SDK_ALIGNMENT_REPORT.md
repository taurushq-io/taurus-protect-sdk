# SDK Alignment Report

**Date**: 2026-09-03
**Branch**: `feature/sdk_gov_rules`
**Reference**: Java for naming and API shape; the shared vector files for byte-level behaviour

## How alignment is enforced

Alignment is checked by gates, not by inspection. All of them are consumed by every SDK's unit
suite, by `build.sh`, or by a script at the repo root:

| Gate | File / command | Covers | Enforced |
|---|---|---|---|
| Cell wire parity | `scripts/resources/governance-cell-vectors.json` (39 vectors) | Each typed `RuleCell` encodes to exact recorded bytes and decodes back | all 4 unit suites |
| Lossless / non-canonical parity | `scripts/resources/governance-lossless-vectors.json` (8 vectors) | Schema-newer and deliberately non-canonical wire forms survive a round trip | all 4 unit suites |
| Authorization-error parsing | `scripts/resources/authorization-error-vectors.json` (8 vectors) | The 403 role-list parser agrees on every server wording | all 4 unit suites |
| Verification behaviour | `scripts/resources/verification-behaviour-vectors.json` (`rule_key` 12, `hash_coverage` 8, `contains_hash` 7, `memo_key` 8, **`legacy_hash` 7**, **`rule_tier_candidates` 5**) | The verification primitives: which `(blockchain, network)` selects the rules, hash coverage, memo-key injectivity, **which payload step 6 parses**, **which rule tiers apply when the network is unsigned** | all 4 unit suites |
| Signed fixtures | `scripts/resources/verification-signed-fixtures.json` (10 SuperAdmin + 8 group vectors) | Both thresholds count DISTINCT SIGNING KEYS, not signature entries | all 4 unit suites |
| Crypto + legacy hashes | `docs/test-vectors/crypto-test-vectors.json` (`legacy_hash_address` 3, `legacy_hash_asset` 3, **`canonical_string` 6**) | Hash/HMAC primitives, the legacy-hash strategies, **and the TPV1 canonical string** | all 4 unit suites |
| Signing-site inventory | `python3 scripts/signing-sites/check.py` | Every ECDSA signing site declares `verifies` or `signs-own-bytes` | repo-root script |
| API-surface parity | `scripts/api-surface/diff.py` | Every service exists in all four SDKs; method-count deltas reported | `build.sh docs` |
| Docs match code | `build.sh docs --check` | No documented method that does not exist; generated index current | `build.sh docs --check` |

The bold rows are new in the 2026-09-10 security-scan pass, and two of them exist because an
existing gate was structurally unable to catch the defect it looked like it covered.

**A note on "Enforced: all 4 unit suites", because it was briefly untrue.** The three new
sections — `legacy_hash`, `rule_tier_candidates` and `canonical_string` — shipped consumed by
the **Go suite alone**. The other three loaders asserted the `counts` block, which proves the
file is *well formed*, and never asserted the *behaviour*. That is not a gate, and it had a
direct cost: **Python shipped the exact single-tier rule lookup `rule_tier_candidates` exists
to forbid** while three SDKs carried the candidate-set fix, and nothing anywhere went red. It
was found by reading the four implementations side by side, not by a test — the drift pattern
this repo keeps producing. All four loaders now consume all three sections; verified by
truncating the file (the Python total drops 53 → 39, the documented regression signal) and by
deleting it (every loader hard-fails).

The lesson worth keeping: **a `counts` assertion is not consumption.** When adding a section,
grep for a reader in each of the four suites before claiming the row is enforced.

- **`legacy_hash`** asserts the *parsed model*, not a hash. `crypto-test-vectors.json`'s
  `legacy_hash_*` groups assert hash values and counts — and the legacy-strip injection moves no
  hash, so those groups passed against it unchanged (they still do, correctly).
- **`canonical_string`** is the first thing to pin the TPV1 canonical message at all. The
  `hmac_sha256` group HMACs a hardcoded string that *looks* like a canonical message, and no
  consumer routed through `calculateSignedHeader` — which is how Python and TypeScript came to
  upper-case the HTTP method while Java and Go signed it verbatim, leaving a lowercase-method
  caller unable to authenticate against one of the two families.

A gate that cannot fail is worse than no gate, so each of the new sections was verified by
reverting the fix and confirming the suite goes red, and by deleting the vector file and
confirming the loader hard-fails rather than silently skipping.

The lossless vectors and the API-surface differ are **new in this pass**. Before it, the
eight lossless base64 strings were hand-copied into four separate test suites with nothing
comparing the copies, and no gate covered the service surface at all — which is how three
capability gaps below reached only one or two SDKs unnoticed.

The differ ships a self-test (`scripts/api-surface/selftest.sh`) that runs it against a
fixture with a deliberate gap. A differ that always prints "OK" is worse than no differ,
because its report reads as evidence of parity.

## Status

Verified by running each gate, not asserted:

| SDK | Services | Unit tests | Other gates |
|---|---|---|---|
| Go | 43 | 7 packages | `go vet ./...` (whole module), `GOARCH=386` build, `golangci-lint` **0 issues** |
| Java | 43 | **1386** | Checkstyle **0**, PMD **0**, SpotBugs **0** |
| Python | 43 | **1583** | touched-file flake8 no worse than the committed baseline (repo lint is known-red) |
| TypeScript | 43 | **1851** in 118 suites | `tsc --noEmit`, ESLint **0** |

Measured 2026-09-10, after the security-scan pass. Previous run (2026-09-07): Java 1273,
Python 1419, TypeScript 1722 in 103 suites. **A drop in any of these is a regression even
when nothing reports failure** — a broken TypeScript test file reports "Test suite failed to
run" and drops the count instead of showing red, and a bulk rename can un-name a Python test
so pytest stops collecting it. Compare the totals, not just the pass/fail line.

**Java's PMD moved 10 → 0 in this pass**, so it is a real gate now rather than known-red debt.

Service counts come from `scripts/resources/api-surface.<lang>.json`, regenerated from
source by `build.sh docs`. The "43 services (38 core + 5 TaurusNetwork)" claim is measured
rather than repeated. Per-service method counts are deliberately **advisory**, not a gate:
names differ by language idiom (Go `ListWallets` / Python `list` / Java `getWallets`), so
normalising them either hides real gaps or invents false ones.

`tg-protect-mcpd` consumes the Go SDK and **no longer uses a local-path `replace`** (removed
2026-09-10) — it pins a released commit through the goproxy, so edits here are invisible to it
until they are pushed and it is bumped. Checked against this branch with a temporary
`replace`: it breaks in **exactly 5 places**, all the `ListWhitelistedAssets` tuple →
`WhitelistedAssetResult` change, and needs a deliberate bump. That is a separate repo.

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

## Security-scan pass (2026-09-10)

A scan against commit `20df2ca` produced **30 findings — 19 HIGH, 11 MEDIUM — across all four
SDKs**. They were not 30 independent bugs: most were **one design defect implemented four
times**, which is the drift pattern this repo keeps producing and the reason every fix below
landed with a gate rather than a review. Threat model throughout: **the API server is the
adversary** — the adversary client-side verification exists to defeat, so each of these was a
gap in the product's core promise rather than a hardening nicety.

Three defects were found during the work and are not in the 30. Two of them are arguably worse
than findings that were.

### The legacy-hash strip was injective in nobody's implementation

**All four SDKs.** Step 4 accepts three backward-compatible rewrites of the delivered payload
so that OLD signatures stay valid. Those strips are not injective, and every SDK then parsed
the **delivered** text rather than the variant that matched — so a server could append a
duplicate `,"label":"X"` (or a `contractType` the row never had) immediately before the closing
brace: the strip recovers the genuinely signed bytes, steps 1 through 5 all pass, and
last-duplicate-wins parsing hands the appended value back to the caller as *verified*.

Four exploitable shapes, confirmed by simulating the real regexes. The fix has two halves,
because neither alone is enough:

1. **Step 4 returns the payload it matched, and step 6 parses THAT.** By preimage resistance a
   matching variant *is* the byte string governance signed, so this is closure, not narrowing.
2. **Duplicate object keys are refused**, structurally, inside the parse functions — closing the
   shapes the strip does not reach, including the third parse in `resolveRuleKey` where a
   duplicated `currency`/`network` would re-point rule selection at a weaker quorum.

One accepted regression, and it must not be "fixed": `linkedInternalAddresses[].label` comes
back empty for legacy-era rows. Those labels are live DB values validatord rebuilds on every
read and were never signed — returning them was the bug.

Gated by the new `legacy_hash` section of `verification-behaviour-vectors.json`, which asserts
the **parsed model** rather than a hash. The pre-existing `crypto-test-vectors.json` legacy
groups pass unchanged, which is correct: the fix moves no hash, and that is exactly why they
could never have caught it.

### Whitelist approval signed server-selected rows

**All four SDKs.** `approve` took a list of ids, re-read them, and signed whatever came back
under those ids. Nothing bound the approver's intent to the bytes signed, so a
response-controlling server could answer the id-filtered re-read with a **different** row — one
whose existing signatures already satisfy the container it presents — and harvest a genuine
approver signature over content the approver never saw. Verification is no defence: the
substituted row is a real, validly-signed whitelist entry, just not the reviewed one.

`approve` now takes the **rows a verified read returned**, carrying the metadata hash each had
at review time, and aborts on a mismatch. This is a **breaking change**, taken deliberately: the
pin has to be impossible to forget, so the pin types have no public constructor and are minted
only by `select(ids)` / `selectAll()` on a read result. An empty selection raises rather than
meaning "approve nothing" — the same rule `approveRulesProposal`'s mandatory
`expectedContainerHash` follows.

The value **signed** is still the row's current `metadata.hash`, never the legacy variant step 4
matched: validatord rebuilds the hash array from the current schema and verifies the submitted
signature against those bytes.

### Six more themes, same shape

| Theme | SDKs | What was wrong |
|---|---|---|
| `createAddress` returned an unverified HSM address | all 4 | `getAssetAddresses` had been fixed and `createAddress` missed. Now ONE seam per SDK, with the rule "never return a non-empty address string that has not been verified" — branching on the address STRING, since `status` is server-controlled. Async creation is real, so an unsigned address is withheld rather than returned |
| Pledge-action approval signed an unverified hash | Go, TS | Read paths now verify `sha256(payload) == hash`, and approval verifies again and signs internally. Python was the model; Java has no pledge surface |
| An unsigned DTO `network` picked the rule tier | Go, Java, TS **+ Python** | `includeNetworkInPayload` has **no proto backing**, so it is unsigned and the finding's suggested fix was unavailable. Instead: when the payload omits `network`, EVERY reachable tier must be satisfied. **The scan filed this against three SDKs; Python had the same defect and was not filed.** Found by verifying the fix site-by-site across all four — regressing Python's walk makes the lookup return the *global default*, the broadest tier, which is fail-open. Its two families also name the chain field differently (`currency` vs `blockchain`), so a shared walk reading one name silently treats every contract rule as a wildcard |
| Integrity failures remapped to a retryable `ServerError(500)` | Python, TS | Told a caller to retry a response an attacker controls. Fixed on every path that can raise it; the mechanical remainder is in `TODOS.md` |
| Gson-onto-`Throwable` recursion · rules-cache wedge | Java | A nested error body could raise `StackOverflowError` out of every SDK call; separately, an unchecked throw from the rules fetch left the single-flight flag set, so **one crafted `/rules` response wedged address, asset and price verification process-wide** on an untimed wait |
| Nil-reply deref · TPV1 followed redirects | Go | ~110 unchecked `resp.<Field>` sites behind one generated `decode`; and TPV1 minted a **fresh valid signature** for whatever host a `Location` header named — a signing oracle, not a replayed credential |
| Quadratic bigint decode · unbounded container · unverified `getEnvelope` | TS | The address `getEnvelope` returned the unverified input envelope while the asset one was already correct |

### Not fixed, by decision — and recorded rather than shipped silently

Three themes are blocked on a server-side answer and are documentation-only, each with its
finding ids in `TODOS.md`:

- **`payloadToSign` on a multi-factor signature is unverified**, and cannot be bound: the reply
  carries no entity id. This is the ONE place in the SDK where bytes intended for a signing key
  are returned unverified, and it is now stated as such on the model field and both methods in
  all four SDKs, with the actionable client-side path written down.
- **TPV1's canonical string is not injective** — needs a versioned scheme across four SDKs, the
  Postman collection and validatord, because the server verifies the same string.
- **A uniformly stale rules container** is undetectable client-side: the container arrives
  in-band and validatord exposes no ruleset identity to pin against. The approval pin closes the
  *signing* consequence; the read consequence stays open.

### Three defects found during the work, not in the 30

1. **Python's entire MFA service was dead code that could not run** — it called four generated
   operations that do not exist and imported two request models that do not exist, so every call
   raised `AttributeError` or `ModuleNotFoundError`. **Its tests passed** because the API object
   was a bare `MagicMock()`, which answers any attribute. Rewritten onto the four real
   operations; the stub is now `MagicMock(spec=MultiFactorSignatureApi)`, which is the actual fix
   for the class of bug, and the four phantom operations are pinned as absent.
2. **TPV1 diverged on HTTP-method case** — an interop break, not just an inconsistency. Python
   and TS upper-cased; Java and Go signed verbatim. A caller issuing lowercase `get` signs a
   different canonical string in the two families, so one family cannot authenticate. Aligned on
   upper-casing (HTTP methods are case-sensitive per RFC 9110, so normalising cannot break a
   working deployment) across all four SDKs and all five Postman signing blocks, and pinned by
   the new `canonical_string` vectors.
3. **Java's `ApiExceptionMapper` let the response body override the HTTP status** — a body
   `{"code":200}` on a 503 yielded `getCode() == 200` and `isRetryable() == false`. The transport
   status always wins now. Same taxonomy inversion as the retryable-integrity finding, reached
   from the other side.

### What made each fix a gate

Every new test was run against the UNFIXED code first, and the verbatim failure captured. Three
worth recording because they show what the old suites could not see:

- Removing both legacy defences: `attacker-appended label reached the caller as verified:
  "Coinbase Prime custody"`. The end-to-end address flow had only ever been exercised with nil
  inputs, so nothing could tell "closed" from "narrowed".
- Restoring the single rule-tier lookup: *"the row met goerli's 1-of-1 but not mainnet's
  ops+compliance, and it verified anyway."*
- Removing the duplicate-key pass: **7 of 9** Java cases go red, including `resolveRuleKey`.
- Reverting the TPV1 method case: the `lowercase method` and `mixed-case method` vectors go red.
- Removing the nil-reply guard: `panic: runtime error: invalid memory address or nil pointer
  dereference`.
- Removing the rules-cache `finally`: *"a second caller never returned"* — the two cache tests
  were the last thing still red after everything else had landed, which is the gate working.

---

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
