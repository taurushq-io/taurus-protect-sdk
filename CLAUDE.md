# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# General guidelines

## Workflow orchestration

### 1. Plan Mode Default

* Enter plan mode for ANY non-trivial task (3+ steps or architectural decisions)
* If something goes sideways, STOP and re-plan immediately – don't keep pushing
* Use plan mode for verification steps, not just building
* Write detailed specs upfront to reduce ambiguity
* Write down the plan (to be able to restart later if needed)
* At the end of each plan, give me a list of unresolved questions to answer, if any.

### 2. Subagent Strategy

* Use subagents or teams liberally to keep main context window clean
* Offload research, exploration, and parallel analysis to subagents or teams
* For complex problems, throw more compute at it via subagents or teams
* One task per subagent for focused execution

### 3. Self-Improvement Loop

* After ANY correction from the user: update the "Lessons Learned" section below
* Write rules for yourself that prevent the same mistake
* Ruthlessly iterate on these lessons until mistake rate drops
* Review lessons at session start for relevant project

### 4. Verification Before Done

* Never mark a task complete without proving it works
* Diff behavior between main and your changes when relevant
* Ask yourself: "Would a staff engineer approve this?"
* Run tests, check logs, demonstrate correctness

### 5. Demand Elegance (Balanced)

* For non-trivial changes: pause and ask "is there a more elegant way?"
* If a fix feels hacky: "Knowing everything I know now, implement the elegant solution"
* Skip this for simple, obvious fixes – don't over-engineer
* Challenge your own work before presenting it

### 6. Autonomous Bug Fixing

* When given a bug report: just fix it. Don't ask for hand-holding
* Point at logs, errors, failing tests – then resolve them
* Zero context switching required from the user
* Go fix failing CI tests without being told how

## Task Management

* Plan First: Write the plan to .claude/tasks/ with checkable items
* Verify Plan: Check in before starting implementation
* Track Progress: Mark items complete as you go
* Explain Changes: High-level summary at each step
* Document Results: Add review section to .claude/tasks/
* Capture Lessons: learn (command: `learn.md`) and update `CLAUDE.md` for shared or global project
  context, and `CLAUDE.local.md` for private, developer-specific notes after corrections

## Core Principles

* Simplicity First: Make every change as simple as possible. Impact minimal code.
* No Laziness: Find root causes. No temporary fixes. Senior developer standards.
* Minimal Impact: Changes should only touch what's necessary. Avoid introducing bugs.

## Repository Overview

This is a monorepo containing SDKs for the Taurus-PROTECT API, a cryptocurrency custody and transaction management platform.

### SDK Directories

| Directory | Language | Status | Services |
|-----------|----------|--------|----------|
| `taurus-protect-sdk-java/` | Java | Active development | 43 (38 + 5 TaurusNetwork) |
| `taurus-protect-sdk-go/` | Go | Active development | 43 (38 + 5 TaurusNetwork) |
| `taurus-protect-sdk-python/` | Python | Active development | 43 (38 + 5 TaurusNetwork) |
| `taurus-protect-sdk-typescript/` | TypeScript | Active development | 43 (38 + 5 TaurusNetwork) |

**Service Parity (All SDKs):** All four SDKs have 43 high-level service wrappers (38 core services + 5 TaurusNetwork services). TypeScript also provides access to all 61 OpenAPI-generated APIs. All SDKs split Taurus Network into 5 services (Participants, Pledges, Lending, Settlements, Sharing).


## Build gates — there is no CI pipeline in this repo

`.github/` contains **only `CODEOWNERS`**: no GitHub Actions workflow, no Jenkinsfile, no other
pipeline definition. The only executable definition of "green" is each SDK's `./build.sh`. Don't
claim "CI passes" without naming which command was run.

Baselines measured against committed `main`/`master` (see `CLAUDE.local.md` for the
`git archive HEAD` technique) — know these before blaming a branch for a red gate:

| Gate | Command | Baseline on master |
|---|---|---|
| Java tests | **`mvn test -o`** (whole reactor — `-pl client` skips the openapi module's own suite) | 0 failures |
| Java Checkstyle | `mvn checkstyle:check -pl client` | **0 — a real gate; do not break it** |
| Java PMD | `mvn pmd:check -pl client` | **0 — a real gate since 2026-09-10** (was 10; all of it unused imports, unnecessary FQNs and two unchained causes) |
| Java SpotBugs | `mvn spotbugs:check -pl client` (**no `-o`**) | 0 bugs |
| Go | `go build ./... && go vet ./...` (whole module — see below) `&& golangci-lint run` | 0 issues |
| Go 32-bit | `GOARCH=386 go build ./pkg/...` | worth checking; a `math.MaxUint32` compare broke it once |
| TS types | `npx tsc --noEmit` | clean — but **does not cover `tests/`** (see the TS CLAUDE.md) |
| TS ESLint | `npx eslint src --ext .ts` | was **dead** (deps + `lint` script declared, no config); `.eslintrc.json` now exists |
| Python lint | `./build.sh lint` | **abandoned: 96 files fail flake8, 3008 mypy errors.** Never expect green; don't reformat the repo to chase it |
| Docs match code | `./build.sh docs --check` (each SDK) | 0 — a real gate; exits 1 on a documented method that does not exist, or a stale generated index |
| Cross-SDK API surface | `./scripts/api-surface/generate.sh all` | 0 missing services — see "API-surface parity gate" below |

Java static analysis is worth running on any change: it is cheap, and `RuleCell.java` alone once
carried 37 PMD violations that went unnoticed because the file was untracked.

## `.gitignore` hides source whose name looks like a secret

The secrets block contains `*credentials*`, `*api*key*`, `*.secret`, `*.pem`, `*.key`. `*credentials*`
matched all 8 hand-written `Credentials` auth-API files in the four SDKs, so they were silently
untracked and a fresh clone did not compile in any language — while every local build stayed green.
Explicit `!` negations for those paths now sit directly under the pattern. **If you regenerate or
reorder `.gitignore`, keep them**, and after adding any source file whose name contains
`credential`/`apikey`/`secret`, check `git check-ignore -v <path>` and `git status --porcelain | grep '^??'`.

## SDK Alignment Reference

When aligning SDKs, **Java is the reference for naming and API shape**; for byte-level behaviour the
shared vector files are the oracle, not any one SDK (past passes found Java itself was the outlier on
unknown-enum passthrough, deterministic encoding and aggregate `hasUnknownFields`). The report lives
at [`docs/SDK_ALIGNMENT_REPORT.md`](docs/SDK_ALIGNMENT_REPORT.md).

**Alignment is gated, not inspected.** Six gates, all runnable locally:

| Gate | File / command |
|---|---|
| Cell wire parity (39 vectors) | `scripts/resources/governance-cell-vectors.json` |
| Lossless / non-canonical parity (8 vectors) | `scripts/resources/governance-lossless-vectors.json` |
| Authorization-error parsing (8 vectors) | `scripts/resources/authorization-error-vectors.json` |
| **Verification behaviour (35 vectors)** | `scripts/resources/verification-behaviour-vectors.json` |
| **Signed fixtures — both thresholds (10 SuperAdmin + 8 group vectors)** | `scripts/resources/verification-signed-fixtures.json` |
| **Signing-site inventory** | `python3 scripts/signing-sites/check.py` |
| API-surface parity | `./scripts/api-surface/generate.sh all` |
| Docs match code | `<sdk>/build.sh docs --check` |

**Signing-site inventory** enumerates every ECDSA signing call site across the four SDKs
and requires each to declare, in `scripts/signing-sites/manifest.json`, whether it
`verifies` (the bytes signed come from data this SDK verified in the same call) or
`signs-own-bytes` (a governance proposal carries 0..N signatures — signing IS the approval
step, so there is nothing to verify against yet). Currently 17 sites: 13 verify, 4 sign
own bytes.

It exists because the same defect was made four times and found by grep, not by a test:
`approveRequests` checked that a metadata hash was non-empty and then signed it. A fifth
site (Python's `approve_pledge_actions`) had no verification at all and drifted alone,
because nothing enumerated the signing surface so a new signer inherited no rule.

Rationale, keying design and how to add a site are in
[`scripts/signing-sites/CLAUDE.md`](scripts/signing-sites/CLAUDE.md), loaded when you work
there.

**A `counts` assertion is NOT consumption — check for a reader in each of the four suites.**
Every loader asserts the file's `counts` block, and that is what the repo describes as the
"a case added and consumed by nobody fails loudly" guard. It is weaker than it reads: it proves
each section has the declared NUMBER of vectors, not that anything asserts their OUTCOMES.
Three sections shipped consumed by the Go suite alone (`legacy_hash`, `rule_tier_candidates`,
and `canonical_string` in `crypto-test-vectors.json`), and the cost was concrete — **Python
shipped the exact single-tier rule lookup `rule_tier_candidates` exists to forbid**, while
three SDKs carried the fix and no gate went red. So when adding a section, grep for a reader in
all four suites; and when trusting a section, do the same before believing it covers you.

**Verification behaviour vectors** cover the primitives that had drifted apart: which
`(blockchain, network)` pair selects the governance rules (`resolveRuleKey`), whether a
hash is covered by a signature (`verifyHashCoverage` / `containsHash`), and whether the
governance verification memo key is injective (`memo_key` — see the Cross-SDK Security
Rules). All are pure input→outcome, so no key material is needed and all four SDKs assert
them from one file. The file declares its own per-section counts and every loader asserts
them, so a case added without being consumed fails loudly. Loaders: Go
`pkg/protect/helper/verification_behaviour_vectors_test.go`, Java
`.../client/helper/VerificationBehaviourVectorsTest.java`, Python
`tests/unit/helpers/test_verification_behaviour_vectors.py`, TS
`tests/unit/helpers/verification-behaviour-vectors.test.ts`.

The `memo_key` section needs a **second** reader in Go and Java — `memo_key_vectors_test.go`
and `MemoKeyVectorsTest.java`, both in the `service` package — because the key function is
unexported/package-private there, the same reason the signed fixtures need two consumers.
Python and TS consume it from their single loader. The generic count-vs-declared loop stays
in the primary loader, so the two cannot drift on it.

**Signed fixtures** close the gap this file used to record as *"No automated cross-SDK gate
covers this"* for the five `minValidSignatures` distinct-key cases. They carry **PUBLIC keys and
signatures only** — the gate verifies, it never signs — so no private key material is committed
and the `.gitignore` secrets block is irrelevant. Generated by Go with fresh keys each run:

```bash
UPDATE_VERIFICATION_SIGNED_FIXTURES=1 go test ./pkg/protect/helper -run TestUpdateSignedFixtures
```

Regeneration yields a different but equally valid set; the other three SDKs must still agree
with every recorded outcome. Each loader also asserts the file contains a case where the entry
count meets the threshold but the distinct-key count does not — without one, the whole file
would pass against an implementation that counts entries. Verified by regressing Go to entry
counting: two vectors go red.

**The file has two sections, one per threshold**, because the counting rule is shared and the
four SDKs must not be able to drift on one without drifting on the other:

| Section | Threshold | Signed bytes | Loader per SDK |
|---|---|---|---|
| `superadmin_threshold` | `minValidSignatures`, tenant-wide (step 2) | the rules container | Go `helper/signed_fixtures_test.go`, Java `helper/SignedFixturesTest.java`, Python `tests/unit/helpers/test_signed_fixtures.py`, TS `tests/unit/helpers/signed-fixtures.test.ts` |
| `group_threshold` | `GroupThreshold.minimumSignatures`, per group (step 5) | `JSON(hashes)` — the array on the entry, NOT the metadata hash alone | same three files, plus Java `service/SignedFixturesGroupThresholdTest.java` |

Java needs a second file because its group walk is a **package-private service method** —
whitelist verification lives in the services there, not `helper/` — so the shared file loader
was lifted to `testutil/SignedFixtures.java` rather than duplicated. Each section carries its
own non-vacuity guard; the group one counts only signers **inside the group**, since an
out-of-group entry is skipped without being an error. `counts` is asserted section by section,
so a section added and consumed by nobody fails loudly.

Regressing to entry counting takes **three** group vectors red in Go (duplicated entry,
re-signed entry, two user IDs sharing one key), one in Java, and the whole group section in
TypeScript and Python.

**Java's models are the outlier a new loader trips over.** The fixture carries base64
signatures and PEM keys because three of the four consume both verbatim; Java does not:

| | signature field | rules-container user key |
|---|---|---|
| Go | `string` (base64) | `*ecdsa.PublicKey` — decode the PEM |
| Java | **`byte[]` — base64-DECODE it** | `setPublicKey(PublicKey)` via `CryptoTPV1.decodePublicKey(pem)`, or `setPublicKeyPem(String)` |
| Python | `str` (base64) | `public_key_pem: str` |
| TypeScript | `string` (base64) | `publicKeyPem: string` |

Two more Java-only shapes: the entry's nested setter is `WhitelistSignature.setSignature(WhitelistUserSignature)`
(not `setUserSignature`), and `hashes` has no setter at all — populate it with `getHashes().addAll(...)`.
The Java SDK's `CLAUDE.md` repeats these under "Building a whitelist signature entry in a test",
for a reader who never opens this file; keep the two in step.

### API-surface and docs gates (`scripts/api-surface/`)

`./scripts/api-surface/generate.sh all` regenerates a per-SDK API surface from source, re-renders the
generated method index in each `docs/SERVICES.md`, and runs the cross-SDK differ. Each SDK's
`build.sh docs` / `docs --check` calls it for its own language.

Two rules to know before touching either gate:

- **Service parity is a hard gate; method-count deltas are advisory.** Method names legitimately
  differ by idiom (Go `ListWallets` / Python `list` / Java `getWallets`), so normalising them either
  hides real gaps or invents false ones. Accepted naming pairs go in
  `scripts/api-surface/aliases.json`, and **every entry must correspond to a row in the report's
  "Remaining known differences" table** — that pairing is what keeps an alias a recorded decision.
- **Never regenerate a `docs/SERVICES.md` wholesale.** Those files hold ~7k lines of hand-written
  prose no extractor can reproduce; only the marker region is rewritten. `docs --check` additionally
  fails on any documented method that does not exist — it found ~100 on its first run.

Pipeline internals (extractor tooling and why, the doc-shape parsing rules, how to add a service)
are in [`scripts/api-surface/CLAUDE.md`](scripts/api-surface/CLAUDE.md), loaded when you work there.

## Client Authentication (`Credentials`) — cross-SDK

All four SDKs construct the client with a single `Credentials` sum-type (not scattered
apiKey/apiSecret/bearer params). Three variants, aligned everywhere:
`apiKey(key, secret)` (TPV1-HMAC) · `bearerToken(token)` · `bearerTokenProvider(fn)`
(per-request token — e.g. one shared client serving many callers, each carrying its own token).

- **Clean constructor takes ONE `Credentials`:** Go `WithCredentials(Credentials)`; Java
  `create(host, Credentials, keys, minSig)` / `builder().credentials(...)`; Python
  `create(host, credentials, super_admin_keys_pem=…)`; TS `create({ host, credentials, … })`.
- **SuperAdmin keys are MANDATORY for every mechanism** (api-key AND bearer). There is NO
  keys-optional / skip-client-side-verification-under-bearer path — it was removed because it
  violated the verification-mandatory invariant. **Do not re-introduce it.**
- **Flat api-key params are kept but deprecated** (additive; they build an api-key `Credentials`
  internally) — do not remove them (external callers + many tests). Go's `WithCredentials`
  was retyped to take `Credentials` (its only callers were SDK tests + tg-protect-mcpd).
- **Bearer convenience constructors were removed** — Go `WithBearerToken*`; Java
  `createWithBearerToken*` + builder `bearerToken*`; Python `bearer_token=`/`bearer_token_provider=`
  create() params; TS `bearerToken`/`bearerTokenProvider` config fields. Bearer is reachable
  ONLY via `Credentials`. **Do not re-add them.** The transport mechanism code stays
  (Go `bearerTransport`, Java `ApiKeyTPV1Auth.setBearerTokenProvider`, Python `bearer_rest.py`, TS `bearer-middleware.ts`).
- `Credentials` lives at: Go `pkg/protect/credentials.go`; Java `client/.../client/Credentials.java`;
  Python `taurus_protect/credentials.py`; TS `src/credentials.ts`.
- **BOTH transports refuse redirects by default AND pin the host, and both guards are required.**
  The bearer half got this in 2026-09; the api-key (TPV1) half was left with neither until
  2026-09-10, which was worse rather than equivalent. Because the Authorization header is created
  *inside* `RoundTrip`, below net/http's redirect handling, the follow-up request net/http
  synthesises from a server-supplied `Location` carries nothing for net/http to strip — and the
  TPV1 transport then mints a **brand-new, fully valid signature** (fresh nonce and timestamp)
  over the attacker-chosen method, host, path, query and body. That is a signing oracle, not a
  replayed credential: a `303` turns any signed call into an authenticated GET of the
  intermediary's choosing, with the response flowing back through them, and the 10-hop limit
  allows ten per legitimate call. `apiKeyCredentials.apply` therefore keeps the host it is handed
  (it used to discard it with `_ string`) and fails closed on an unparseable one, exactly as the
  bearer half does. Tests:
  `TestTPV1SignatureIsNotSentToARedirectTarget` / `…CrossHostEvenWhenRedirectsAreFollowed`,
  mirroring the bearer pair — plus a same-host test, or "refuse everything" would pass both.
- **The TPV1 canonical string upper-cases the HTTP method, in all four SDKs and the Postman
  script.** Python and TS always did; Java and Go signed it verbatim, so a caller issuing a
  lowercase `get` produced a different canonical string in the two families and one of them could
  not authenticate. Methods are case-sensitive per RFC 9110 and every real caller sends upper
  case, so normalising removes the split without changing what a working deployment signs.
  Pinned by the `canonical_string` group in `docs/test-vectors/crypto-test-vectors.json` — which
  did not exist before: the `hmac_sha256` group HMACs a hardcoded string that merely *looks* like
  a canonical message, and no consumer routed through `calculateSignedHeader`, which is how the
  split survived. **This is not the injectivity defect** — that one changes what a correct caller
  signs and needs a versioned scheme (TPV2); see `TODOS.md`.

## Generated clients: a 2xx with no body must not yield a nil reply

Go's generated `decode` reported SUCCESS without populating the typed reply for two bodies a
server fully controls — an empty one, and the JSON literal `null` — so `Execute` returned
`(nil, resp, nil)` and the very next line in every hand-written service was a nil-pointer
dereference. Measured: ~110 unchecked `resp.<Field>` dereferences across 38 service files, and
even the apparently-careful `if resp.Result == nil` is itself a dereference of a nil `resp`.

That is a whole-process kill rather than a failed call — nothing in the SDK recovers — so in
`tg-protect-mcpd`, which runs one client per tenant in one daemon, a single empty 200 aimed at one
tenant takes down every tenant, repeatably. A server that merely errors cannot do that.

The guard (`assertReplyDecoded`) lives in the SHARED decode path, and is applied in **two places
that must stay in step**: `taurus-protect-sdk-go/internal/openapi/client.go` (the committed
generated file) and `scripts/resources/templates/go/client.mustache` (a newly vendored template
override passed to the generator with `-t`). The template is what makes it survive
`./build.sh generate`, which does `rm -rf internal/openapi` — a hand-edit to the generated file
alone is silently reverted. Only the "asked to populate a typed reply and got nil" shape is an
error: `*string`, `*os.File`, `[]byte`, value types and `google.protobuf.Empty` endpoints all
legitimately accept an empty body.

## Security Invariants (Cross-SDK)

### Implemented Security Features (All SDKs)

1. **Request approval with private key signing** — `approveRequest(request, privateKey[, comment])` signs the JSON hashes array with ECDSA. The comment is optional in all four (variadic in Go, an overload in Java).
2. **Address signature verification** — every `AddressService` **and `AssetService`** read verifies against the HSMSLOT key from a `RulesContainerCache` supplied at construction, never a setter. Both fail fast; both take the cache as a mandatory parameter.
3. **Request hash verification** — every read path (get, list, listForApproval) verifies `computedHash == providedHash` in constant time, marks `hashVerified`, and excludes-and-reports failures on list paths. Payload accessors **throw** on unverified metadata.
4. **Whitelisted address 6-step verification** — with legacy hash computation for backward compatibility.
5. **Whitelisted asset 6-step verification** — the same chain against `ContractAddressWhitelistingRules`, plus the parse from the verified payload. Often called "5-step" for the signature checks; the parse is what stops an unsigned value reaching the caller, so count it.
6. **Governance rules signature verification** — `verifyGovernanceRules()`, thresholded on DISTINCT signing keys.
7. **Per-group whitelist threshold** — step 5 counts DISTINCT SIGNERS keyed on the container-resolved public key, in both whitelist flows.
8. **Approval refuses unverified metadata** — `approveRequests` will not sign a hash that verification has not cleared.
9. **Verified whitelist approval** — `approveWhitelistedAssets` re-reads, verifies, and signs those hashes, all-or-nothing.
10. **Price signature verification** — against `PRICEUPDATER` keys, with the verified container deciding whether prices must be signed at all.

### Anti-Patterns

- Never use bare `except Exception` / `catch {}` — always catch specific exceptions
- Always add explicit null/nil checks before constant-time comparison
- Log security cleanup failures — don't silently ignore secret wiping failures
- Pre-fetch shared resources before loops to avoid N+1 patterns

### Cross-SDK Security Rules

- **Constant-time comparison**: Perform dummy comparison on length mismatch; never break/return early in multi-signature loops. Scope: this is about comparing SECRET or hash material. It does NOT apply to trying a signature against a set of PUBLIC keys — `matchingKey`/`IsValidSignature` returning on the first key that verifies is correct, since data, signature and keys are all public and there is no timing channel to leak. Don't "fix" that into a full scan.
- **Address verification mandatory**: `RulesContainerCache` must be provided at construction, never optional
- **Legacy hash**: Address and Asset verifiers use DIFFERENT functions (`ComputeLegacyHashes` vs `ComputeAssetLegacyHashes`)
- **Hash values belong in verification errors.** Request hash errors carry computed/provided;
  whitelisted address/asset step-5 skip reasons carry the metadata hash, the signer's `userId`
  and the signature's covered-hashes list. None of it is secret — a hash is SHA-256 of a
  payload the caller already holds — and without it a threshold failure is undebuggable. Note
  the covered-hashes list spans the whole approval **batch**, so it names hashes of the other
  entities approved alongside this one; that was a deliberate call, not an oversight.
  All four SDKs must agree, so **do not** strip them from one.
- **The signing path RE-VERIFIES; it never trusts a `hashVerified` flag.** The flag is a
  serialized field in Go/Python/TS (`json:"hash_verified,omitempty"`, a Pydantic field, a
  compile-time-only `readonly`), so a request decoded from a queue, webhook or cached blob
  can arrive claiming true and get an attacker-chosen hash signed by the approver's real
  key. `approveRequests` therefore calls the verify primitive again — a SHA-256 over a
  string already in hand, no network. The flag stays as advisory metadata. Python was the
  worst case: the forged flag also opened the payload accessors, having no equivalent of
  Go's unexported `entries`.
- **Every batch approval re-reads through the verifying LIST path, filtered by ids, and
  checks completeness.** One round trip and one rules-container fetch for the whole batch;
  the per-id GET it replaced cost both per id, so a 50-id approval was 50 sequential round
  trips each running the full verification chain. `includeForApproval` is required — the
  rows being approved are pending, so the default list omits them. And every requested id
  must appear in the returned page or the call aborts: a filtered page that silently omits
  an unverifiable row must not become an approval of fewer rows than the caller asked for.
  Approvals are **all-or-nothing** — one signature covers every hash in the batch.
- **One verification seam per request service; nothing else may call the mapper.**
  `verifiedRequest(dto)` maps, verifies, marks and returns, and all nine paths (get, two
  lists, six creates) go through it. The list paths skipped verification in one pass and
  the create paths in the next, both because "remember to verify" was a rule rather than
  the only available construction path. Do not call `RequestFromDTO` / `requestFromDto` /
  `request_from_dto` / `RequestMapper.INSTANCE.fromDTO` from a service again.
- **Field sourcing**: Security-critical fields MUST come from verified payload only, never from unverified DTO
- **Hash exists + no payload**: MUST fail explicitly (never silently return)
- **P-256 curve validation**: All SDKs MUST validate ECDSA keys use P-256 (secp256r1) before use
- **Verification functions must not mutate input**: Return state in result struct instead
- **RulesContainerCache requires SuperAdmin verification**: The rules container must be verified by SuperAdmin key signatures before trusting HSM public keys.
  **The cache must fetch through the governance SERVICE, never the generated API.** That
  service is where verification lives; the generated `ruleServiceGetRules` returns an
  unverified DTO and the raw mapper decodes without checking anything. TypeScript's
  `ProtectClient.getRulesCache()` did exactly that — decoded unverified and discarded
  `rulesSignatures` — so `AddressService` verified addresses against an HSM key nothing
  had authenticated, and anyone able to influence that response could have
  attacker-chosen addresses verify clean. Go, Java and Python were correct; TS was the
  lone outlier and no SDK tested the path, which is why it survived review.
  **Now gated in all four**, each feeding an unsigned container through the cache and
  asserting it is refused: Go `pkg/protect/rules_cache_verification_test.go`, TS
  `tests/unit/cache/rules-cache-verification.test.ts`, Python
  `tests/unit/cache/test_rules_container_cache.py` (`TestRulesContainerCacheVerification`),
  Java `RulesContainerCacheTest`. Two traps when editing those: the container must be
  **wire-valid**, or the test passes on a decode/parse error whether verification runs
  or not (this bit twice — Python and Java both need the `_decodesCleanly` guard test
  that sits beside them); and Go's cache no longer has `Set`/`SetFetcher`, so the
  constructor's fetcher is the only way a container gets in.
- **`minValidSignatures` counts DISTINCT signing keys, never signature entries**: ECDSA is randomized, so counting entries lets one compromised key satisfy any threshold, making 2-of-N no stronger than 1-of-N. Identify a signer by a hash of its encoded public key (so the same key configured twice counts once) — not by list index and not by the caller-supplied `userId`. Implemented in each SDK's single `signature_verifier` helper, which also serves the whitelisted-address and whitelisted-asset paths. **Gated cross-SDK** by `scripts/resources/verification-signed-fixtures.json`, which pins exactly these five cases with real signatures: N distinct keys pass at N; N entries from one key fail at N but pass at 1; a repeated identical signature counts once; a key configured twice counts once; an unconfigured key contributes nothing. Loaded by all four suites — do not delete the per-SDK copies without checking the fixture still covers what they did.
- **`minValidSignatures` and `minimumSignatures` are DIFFERENT thresholds — do not unify them**,
  but both count DISTINCT SIGNERS. `minValidSignatures` is the SuperAdmin rules-container threshold
  (above). `GroupThreshold.minimumSignatures` is the per-group whitelist threshold in step 5, and it
  counts distinct members of that group, keyed on the container-resolved public key — never on the
  server-supplied `userId`. It previously counted signature *entries*, which let a duplicated entry
  from one member satisfy an N-of-M group and promote an under-approved whitelist entry to approved;
  two user IDs sharing one key have to count once for the same reason, since that is one compromised
  secret. What differs between the two thresholds is scope (tenant SuperAdmins vs one group's
  members), not the counting rule. The Java flow lives in the services, not a `helper/` file, so grep
  for both. **Gated cross-SDK** by the `group_threshold` section of
  `scripts/resources/verification-signed-fixtures.json` (see "Signed fixtures" above).
- **Governance access is typed-only, and that must be ENFORCED, not conventional.** Writes take a
  `DecodedRulesContainer` and the SDK encodes it internally, so no caller-supplied base64 reaches the
  wire; `approveRulesProposal` re-fetches the pending proposal and signs *its* decoded bytes, so a
  caller cannot get arbitrary bytes signed. The generated OpenAPI client must not be publicly
  reachable for governance: Go hides it under `internal/openapi` (compiler-enforced), Java's module
  is `provided` scope, Python's is `_internal`. TypeScript exposed `governanceRulesApi` publicly,
  which made client-side verification opt-out in that SDK alone — the getter is private now
  (`governanceApi()`); **do not re-expose it.** Guarded by
  `tests/unit/client/protect-client.test.ts` → "governance low-level API is not publicly reachable".
- **The JSON bridge is an untyped, LOSSY escape hatch — never the submit path.** protobuf-JSON has no
  representation for unknown fields, so a schema-newer container loses data through it, and cells
  appear as opaque base64 (the `RuleCell` union is bypassed). Author in JSON if you must, then decode
  back into the typed container before submitting, so the typed encoder produces the signed bytes.
  Documented in `docs/CONCEPTS.md` → "Governance Rules: Typed Model vs JSON Bridge".
- **Deterministic serialization applies to every path whose output is hashed or signed**, including
  the JSON bridge — `properties` is `map<string, bytes>` at five levels and map order is unspecified
  unless forced. Go `deterministicMarshal`, Python `SerializeToString(deterministic=True)`, Java
  `RulesContainerMapper.INSTANCE.deterministicBytes`, TS `sortedMap`/`sortPropertiesMaps`. Add it to
  any new encode entry point.
- **Nested nodes need their own checks.** `hasUnknownFields` / `HasUnknownFields` must walk the four
  nested rule-detail sub-messages (`evmCallContract`, `xtzCallContract`, `cashSettlement`,
  `cosmosDetails`), not just the details node — a container reported "clean" while a nested node
  carried schema-newer data. Java additionally needs the aggregate override on
  `DecodedRulesContainer`, because `RulesNode.hasUnknownFields()` is per-node only.
- **An unverifiable ROW is excluded; an unusable CONTAINER aborts the call. Both must hold, and
  the error taxonomy is what enforces it.** Two SDKs got this wrong in opposite directions and
  neither was tested:
  - **TypeScript** caught everything with a blanket `catch (error: unknown)` and turned a
    `ContainerIntegrityError` into a row exclusion, so a page returned 200-shaped success with
    the affected rows silently missing. It now routes both whitelist services through one
    `rethrowIfNotRowLevel` seam (`src/services/row-level-error.ts`) — extracted, not copied,
    because the address and asset services are the two readers that would otherwise drift.
  - **Java** caught only `WhitelistException`, which is *checked*, while `IntegrityException
    extends SecurityException` is **unchecked** — so steps 1, 2 and 4 escaped the row loop and
    one tampered row aborted the whole listing, the exact failure
    `WhitelistedAddressListResult`'s javadoc claims was fixed. The catch is now
    `catch (ContainerIntegrityException e) { throw e; }` followed by
    `catch (WhitelistException | IntegrityException e)`; the order is required, since the first
    is a subclass of the second.

  Anything that is neither must still propagate: a `TypeError`/`NullPointerException` from a
  real defect reported to the caller as "this address failed verification" is indistinguishable
  from a tampered row. Gated by `tests/unit/services/whitelisted-address-container-abort.test.ts`
  and `service/WhitelistedAddressExclusionTest.java`; both go red if the catch is widened back.
  **Java's fixtures need a `signedAddress` payload AND one signature entry** to reach step 1 at
  all — `initializeEnvelope`'s two preconditions throw the checked exception first, which is
  precisely why every pre-existing bare-DTO test missed the unchecked path.
- **ONE reader per entity. A second, unverified reader of the same endpoint is the bypass.**
  A whitelisted asset and a whitelisted contract are one server entity
  (`/api/rest/v1/whitelists/contracts`) that the SDKs exposed twice: through the verified asset
  service, and again through the contract-whitelisting service with no verification at all
  (TypeScript's copy additionally parsed the signed payload and returned its contents as fact).
  The contract service is now **write-only** in all four; the reads and the DTO-projection
  mappers are **deleted**, not deprecated, because a bypass that still compiles is one someone
  will reach for. For the same reason the verified reader gained the capabilities that used to
  exist only on the unverified one — a for-approval read, the `ids`/`kindTypes`/`query` filters,
  and the payload identity fields (`contractAddress`/`name`/`symbol`/`decimals`/`tokenId`) that
  made `GetWhitelistedAsset` unable to answer the question callers actually ask. **Do not re-add
  a read to the contract service, or a mapper that takes a bare contract envelope.**
  Same rule applied to `AssetService.getAssetAddresses`: it returns `Address`, so it verifies
  the HSM signature and fails fast, exactly as `AddressService` does.
- **The approver signs the CURRENT `metadata.hash`, never the legacy variant step 4 matched.**
  Verification must clear the row first — that is what the re-read is for — but the value signed
  is the row's current hash. `ApproveWhitelistedAddressRequest` and
  `ApproveWhitelistedContractAddressReq` carry **no `hashes` field**: validatord rebuilds the
  array itself, sorted by row id ascending, from `helper.ToWLAMetadata` / `ToWLCAMetadata`, both
  hard-wired to the current schema version, and verifies the submitted signature against those
  bytes (`pkg/whitelist/service/whitelist-service.go`, `wl_address_signatures.go`). The three
  legacy variants are accepted **only** by the per-signature tolerance loop in
  `VerifyWLAddressesIntegrity`, which keeps OLD signatures valid — it does not make a legacy hash
  acceptable for a NEW one. A legacy row therefore ends up carrying a mix, and both kinds count
  toward the group threshold. Go, Java and Python were right; **TypeScript alone signed the
  verifier's `verifiedHash`** and would have been rejected on any legacy row. Its collector is
  now `hashesToSignByID` in both whitelist services — the rename is deliberate, so the old name
  cannot be reintroduced by muscle memory. Do not "align" the other three to TypeScript here.
  The caller's ascending-id sort is not cosmetic either: it reproduces the server's own ordering.
- **Step 6 parses the payload the matched signature COVERED, not the payload the server sent.**
  The legacy-hash tolerance in step 4 accepts a hash over a regex-STRIPPED rewrite of
  `payloadAsString`, and the strips are not injective — so a server can append a member the strip
  removes (`,"label":"X"` before the closing brace, or a `contractType` on a row that has none),
  have the residue land exactly on the genuinely signed bytes, pass every signature check, and
  have step 6 return the appended value as verified. Step 4 therefore carries the matched
  PAYLOAD forward alongside the hash it already returned, and step 6 parses that.
  Two consequences worth knowing before "simplifying" it:
  - The exposure is NOT limited to rows with no inner labels. validatord rebuilds
    `linkedInternalAddresses[].label` on every read from live DB relations rather than from the
    signed envelope, so for a row signed before per-object labels existed, stripping *every*
    label lands on the signed bytes too. Both injectable members (`label` at either level,
    `contractType`) are reachable on any legacy row; only the regex alphabet bounds it, since
    `[^"]*` cannot contain a quote.
  - **Accepted regression:** `linkedInternalAddresses[].label` comes back empty for
    strategy-2-era rows. Those labels were never signed. Do not "fix" it by merging the
    delivered labels back in.
  Duplicate JSON object keys are rejected in all four parsers as the second defence, because the
  appended-`contractType` shape is not a duplicate of anything and the appended-`label` shape is.
  Gated by the `legacy_hash` section of `scripts/resources/verification-behaviour-vectors.json`,
  which asserts the PARSED MODEL — `docs/test-vectors/crypto-test-vectors.json` asserts hashes
  and is structurally blind to this, because the attack moves no hash.
- **A whitelist approval signs a CONTENT PIN, not a list of ids.** `approve` takes a selection
  minted by the preceding verified read (`result.Select(ids...)` / `SelectAll()`), carrying the
  metadata hash each row had at review time, and refuses to sign when a re-read row's hash
  differs. Without it the approval re-read by id and signed whatever came back, so a
  response-controlling server could substitute a row whose existing signatures already satisfy
  the container it presents and harvest a genuine approver signature over content nobody
  reviewed. Same mitigation as `approveRulesProposal`'s mandatory `expectedContainerHash`;
  compare constant-time, and **an empty pin must be an error** — an empty one silently restores
  the unpinned behaviour. The witness type's map is private/unexported in every SDK, so a
  hand-built value is forgeable but useless (the `helper.VerifiedAsset` property).
  A mismatch is a REAL signal: `metadata.hash` is recomputed server-side on every read from the
  immutable envelope plus the row's LIVE linked-address/linked-wallet rows, so renaming a linked
  address moves it. Report both hashes and say to re-read, re-review, re-approve.
- **Every path that returns an `Address` verifies the HSM signature, through ONE seam.**
  `createAddress` was the last one that did not, and it is the highest-value moment for
  substitution — its caller is about to publish or fund a fresh deposit address. A seam
  (`verifiedAddress`, like `RequestService.verifiedRequest`) rather than a per-path check,
  because `AssetService.getAssetAddresses` had already been fixed for exactly this and the create
  path was missed anyway. **The async case is real** (`status` is one of
  `created`/`creating`/`signed`/`observed`/`confirmed`), so the rule is not "always throw" but
  *never return a non-empty address string that has not been verified*: signature present →
  verify; signature absent and address non-empty → refuse; address empty → fine, nothing to
  misuse. Java needed `status` added to its `Address` model to express this at all.
- **`includeNetworkInPayload` is NOT signed and must never be consulted as if it were.**
  `AddressWhitelistingRules` in `scripts/resources/proto/schema/v1/request_reply.proto:3434-3449`
  carries only `currency`, `parallelThresholds`, `properties`, `network`, `lines` — the flag has
  no proto backing in any of the four SDKs (model-only in Go/Java/Python, absent in TS, populated
  only from the JSON bridge in Python). So when the signed payload omits `network`, there is no
  authenticated way to learn whether that was legitimate. The rule instead: enforce EVERY rule
  tier the unsigned DTO network could have selected, rather than the one it named. Reachability
  mirrors the tier walk — all rules for the chain, plus the chain's wildcard-network rule, plus
  the global default *only* when the chain has no wildcard-network rule. With a single reachable
  tier this is a no-op, which is the common case. Gated by the `rule_tier_candidates` section of
  `verification-behaviour-vectors.json`.
- **Every signing path emits COMPACT JSON, and a space is a rejected signature.** The array of
  hashes an approver signs is rebuilt server-side and the submitted signature verified against
  *those* bytes, so the separator is part of the protocol. Go `json.Marshal`, Java/TS
  `JSON.stringify` and Python `json.dumps(..., separators=(",", ":"))` all agree — but Python's
  DEFAULT separators put a space after each comma, and `approve_pledge_actions` was signing
  `["a", "b"]` while every whitelist path in the same SDK signed `["a","b"]`. One-element
  batches masked it. Pass the separators explicitly on any new Python signing path.
- **A row-level failure and a call-level failure need DIFFERENT exception types, and the
  distinction has to survive a new check being added.** The rule the SDKs enforce — an
  unverifiable ROW is excluded, an unusable CONTAINER aborts the call — is only as good as the
  type each new check throws. Java is where this bites, because `WhitelistException` is checked
  and `IntegrityException extends SecurityException` is not: the duplicate-JSON-key rejection
  first shipped throwing the unchecked one, which escapes every `catch (WhitelistException)`,
  so one unparseable row would have aborted a whole listing through a brand-new door. Both
  public parse entry points declare `throws WhitelistException`, which is what made the break
  visible instead of silent. **Throw the row-level type for anything that is one row's
  problem.**
- **The row→container binding is authenticated, not taken on trust.** In normalized list mode a
  row picks its rules container out of `rulesContainers` by `rulesContainerHash` — and both the
  label and the container come from the same response. All four SDKs now recompute the label and
  reject a mismatch, because otherwise a server can file container A under container B's label
  and steer any row to any **other validly-signed** container: an older ruleset with a weaker
  group threshold, say. Both pass step 2, so SuperAdmin verification alone does not catch it.

  The convention is validatord's, one line in `internal/api/v1/whitelist-controller.go`:
  `base64(SHA256([]byte(e.GetRulesContainer())))`. **`GetRulesContainer()` is ALREADY a base64
  string there**, so the digest is over the base64 TEXT, not the decoded protobuf, and the output
  is base64, not hex — hash the decoded bytes instead and every container is rejected, breaking
  all list calls. Pinned by `TestContainerHashLabel_MatchesValidatordConvention` (Go). Do **not**
  confuse it with `enforcedRulesHash`, which is `base64(SHA256(raw protobuf))`, is a backlink to
  a ruleset's *predecessor* rather than its own identity, and is never verified server-side
  despite its comment. Both are 44-char base64 SHA-256 digests, so the mix-up is silent.

  This closes intra-page mixing only. It does **not** close a uniformly stale container: the
  container arrives in-band, so a hostile server can serve one older-but-validly-signed container
  for every row with no mismatch to detect. validatord exposes no ruleset identity to pin against
  (`Rules` returns no `id` and no `enforced` flag), so the SDKs cannot close it — see `TODOS.md`.
- **Approving a whitelist batch is ALL-OR-NOTHING, over hashes this SDK verified.** The API
  takes one signature covering a JSON array of every metadata hash in the batch, so there is no
  partial submission: `approveWhitelistedAssets(ids, privateKey, comment)` re-reads each row,
  signs the hashes *those* rows carry, and aborts the whole call if any row is missing or fails
  verification. Signing the survivors would tell the approver they approved less than they did.
  The contract service's raw-signature form is deprecated — an opaque blob over hashes nothing
  verified means the approver cannot know what they signed.
- **Approval must not sign an unverified hash.** `approveRequests` refuses metadata whose hash
  is not verified, in all four. The flag existed on every read path and was read by nobody. The
  check is ordered **after** the hash-present check (and, in Java, after `checkNotNull(privateKey)`)
  so absent metadata still reports as absent and argument errors stay argument errors. Java has
  no `setHashVerified`: `RequestMetadata.verifyAndMaterialise()` performs the check and sets the
  flag in one operation, so the payload accessors cannot be unlocked without it having run.
- **Prices are signature-verified, and the VERIFIED CONTAINER decides whether they must be.**
  `rate` and `decimals` feed amount conversion, so an unverified price is a wrong number a caller
  acts on. The signature covers the canonical projection of validatord's `CurrencyPrice` —
  exactly `{"blockchain","currencyFrom","currencyTo","decimals","rate"}`, in that order; do not
  reorder or extend it. A `PRICEUPDATER` (role 2) in the SuperAdmin-verified container means
  prices must be signed, so a price with no signatures is an integrity failure; **no**
  `PRICEUPDATER` means this tenant does not sign prices and it passes through. Deciding from the
  price's own `signatures` array would let the response being checked excuse itself.

## The verified HASH is not on the same object in all four SDKs

This cost real time while wiring the batch approvals, because the signing path needs the
hash verification actually cleared and each SDK exposes it somewhere different:

| SDK | Where the approver's hash comes from |
|---|---|
| Go | `asset.Metadata.Hash` / `addr.Metadata.Hash` — the models carry `Metadata` |
| Java | `envelope.getMetadata().getHash()`, and the row **id is on the ENVELOPE** (`envelope.getId()`), not on `WhitelistedAsset` |
| Python | **the envelope only** — `WhitelistedAddress` has NO metadata field (fields: id, address, label, currency, network, status, created_at, contract_type, memo, customer_id, address_type, tn_participant_id, exchange_account_id, linked_internal_addresses, linked_wallets, attributes) |
| TypeScript | `verifier.verify(...).verifiedHash` — the verification RESULT, not the model |

Consequence for Python: its shared address-verification seam returns
`SignedWhitelistedAddressEnvelope`, not `WhitelistedAddress`, and `list`/`list_for_approval`
map to addresses at the call site. Do not "simplify" it back to returning addresses — the
approval path would lose the only reachable hash.

TypeScript's is the structurally strongest of the four: the hash it signs is the value the
verifier returned, not a field re-read off a mapped model.

## Same field, different names and types across the four

Found while adding the governance history exclusions. None of these are bugs; all of them
break a copy-pasted port:

| Concept | Go | Java | Python | TypeScript |
|---|---|---|---|---|
| ruleset creation time | `CreatedAt` | `getCreationDate()` | `creation_date` | `creationDate` |
| history `total_items` | `int64` | `String` | `Optional[str]` | `number` |
| address list result rows | `Addresses` | `getEnvelopes()` | `addresses` | `items` |

**And one that differs WITHIN each SDK rather than between them, which is why it is not in the
table: the whitelisting rule's chain field.** Address rules call it `currency`, contract-address
rules call it `blockchain` — in all four SDKs (Go `Currency`/`Blockchain`, Java
`currency`/`blockchain`, Python and TS the same). It mirrors the proto, so it is not fixable.

**A shared walk over both families that reads one name is FAIL-OPEN**, which is what makes this
worth a note rather than a footnote: a rule whose chain field reads empty is treated as the
*wildcard global default* — the broadest tier — so every contract rule silently becomes a global
default and the narrow tiers vanish. Hit for real on 2026-09-10 while adding
`find_*_whitelisting_rule_candidates` to Python: one `_rule_candidates` helper served both
families and read `rule.currency`, which the asset suite caught only because an unrelated test
asserted on the error string. Python resolves it with a `_rule_chain(rule)` accessor trying
`currency` then `blockchain`; the other three have separate typed walks and cannot make the
mistake. Do not "simplify" the accessor away.

The `total_items` one bites hardest: reducing it by the exclusion count is integer
arithmetic in Go/TS and needs re-stringifying in Java/Python. Keep each SDK's existing
public type — do not unify it as a side effect of a behaviour change.

## Extract the shared seam BEFORE adding the second reader

Done three times in this repo now, and it is the structural fix for the drift that keeps
recurring. The pattern each time: verification lived inline in one method; a second method
was added that needed the same verification; the copy diverged.

- request services: verification was inline in `get`, the list paths skipped it, then the
  six `create*` paths skipped it → one `verifiedRequest(dto)` seam, all nine paths through it
- address services: ~80 lines inline in `list`, and the for-approval reader needed it →
  one `verifiedAddresses(rows, containerCache)` seam
- governance: `getRules` verified, `getRulesHistory` did not → one `verifiedRuleset` seam

So when adding a read or write path next to an existing verified one, extract first and
add second. The seams are deliberately package-private/private-with-a-test-hook (Go
unexported, Java package-private "so the exclusion behaviour can be tested directly",
Python `_`-prefixed, TS `private`) — that is what lets a test drive them without a stub API.

## Verification flows — facts that cost time to rediscover (cross-SDK)

Learned the hard way during the 2026-09-04 verification pass. Each of these was found by
reading the code or real data, not by inference, and each would mislead someone who assumed
otherwise.

### `network` is NOT in the signed address payload unless a rule says so

`AddressWhitelistingRules` carries a per-rule **`includeNetworkInPayload`** flag. When it is
off — which the captured production payload in
`taurus-protect-sdk-typescript/tests/unit/fixtures/whitelisted-address-raw-response.json`
shows is the common case — the signed payload has `currency` but **no `network` field at all**.

Consequence for anything that reads the payload: **requiring `network` rejects correctly-signed
addresses.** `resolveRuleKey` therefore requires the *chain* from the payload and falls back to
the DTO network only when the payload genuinely carries none. Do not "tighten" that into
requiring both.

(Separately, and already noted in the Java CLAUDE.md: the model field has no proto backing, so
it cannot be encoded.)

### There is NO post-loop threshold check in the group-signature walk

The only success exit is *inside* the loop, after a valid signature increments the count
(`whitelisted_address_verifier.go:499-509` and its three peers). So a group with
`minimumSignatures = 0` and users in it demands **one** signature, not zero — a 2-of-N group
silently becomes 1-of-N.

All four now reject `minimumSignatures == 0` on a **populated** group as a malformed container.
An *empty* group with a zero threshold is still fine: there is nobody to sign. This is also why
"just drop the `threshold` fallback so all four agree" is the wrong fix — it aligns Python and
TypeScript **downward**.

### Which rules judge an entity comes from the SIGNED payload, not the DTO

Nothing binds the DTO to the signatures, and `isWildcard("")` is true — so a response with an
empty blockchain would select the broadest (global-default) rule tier. `resolveRuleKey`
(all four) takes the pair from the payload, errors when the chain is absent, and errors when
the payload and DTO disagree on a field the payload carries.

### Exactly TWO hash-comparison functions per SDK

`verifyHashCoverage(hash, signatures)` and `containsHash(hashes, hash)`, both in the shared
hash helper, both used by both verifiers. There were **14** implementations across the repo
before consolidation (Go 2, Java 5, Python 4, TS 3) and they had drifted: Java's were not
constant-time, and Python carried two copies inside one SDK that disagreed on early-return.
**Do not re-add a private copy to a verifier.**

### The `Verified` brand is TypeScript + Go only, deliberately

- **TypeScript** `Verified<T>` (`src/helpers/verified.ts`) — a non-exported `unique symbol`.
  Compile-time enforced, erases at runtime. `verify()` is the only producer.
- **Go** `helper.VerifiedAsset` — unexported fields. Go permits `helper.VerifiedAsset{}` from
  another package, so the guarantee is *forgeable but useless*: a forged value returns `nil`
  from every accessor. Same property `RequestMetadata.entries` already relies on.
- **Java** would need the verifiers extracted out of the services first (see `TODOS.md`).
- **Python** could only get a runtime assert, so it has none.

This asymmetry is a decision, not an oversight. It proves verification RAN — not that it ran
with the right keys or thresholds.

### Reading unverified metadata THROWS, in all four

`ErrMetadataUnverified` (Go), `UnverifiedMetadataException` (Java),
`UnverifiedMetadataError` (Python/TS). Kept distinct from the pre-existing "key not found in
payload" error — one says the field is absent, the other says nothing can be trusted yet, and
collapsing them lets a verification failure read as an absent source address. The Java and
Python types subclass the not-found type so existing catch blocks keep working.

Go's accessors return `(value, error)` as a result; `tg-protect-mcpd` was updated in lockstep.

## Documentation Structure

Documentation is organized hierarchically to avoid duplication:

- **`docs/`** - Common documentation shared across all SDKs:
  - `README.md` - Documentation index with links to all SDK docs
  - `CONCEPTS.md` - Domain model, entities (Wallet, Address, Request, Transaction, etc.)
  - `AUTHENTICATION.md` - TPV1-HMAC-SHA256 protocol, API credentials, SuperAdmin keys
  - `INTEGRITY_VERIFICATION.md` - Cryptographic verification flows for ALL SIX entities: governance rules, whitelisted addresses, whitelisted assets (= whitelisted contracts, one entity with one verified reader), request metadata, addresses and prices; plus how the governance rules for a given entity are selected, and the all-or-nothing whitelist approval flow
  - `BUSINESS_RULES.md` - Business rules, change approval system, entity scopes, and dual-admin workflow

- **`taurus-protect-sdk-java/docs/`** - Java SDK-specific documentation
- **`taurus-protect-sdk-go/docs/`** - Go SDK-specific documentation
- **`taurus-protect-sdk-python/docs/`** - Python SDK-specific documentation
- **`taurus-protect-sdk-typescript/docs/`** - TypeScript SDK-specific documentation

Each SDK directory has:
- `README.md` - Entry point with quick start, services overview, and build commands
- `docs/SDK_OVERVIEW.md` - Architecture, package structure, design patterns
- `docs/SERVICES.md` - Complete API reference for all 43 services
- `docs/CONCEPTS.md` - SDK-specific model classes and exceptions
- `docs/AUTHENTICATION.md` - SDK-specific authentication implementation
- `docs/USAGE_EXAMPLES.md` - Code examples and patterns
- `docs/WHITELISTED_ADDRESS_VERIFICATION.md` - Address verification flow implementation
- `docs/WHITELISTED_ASSET_VERIFICATION.md` - Asset verification flow implementation

SDK-specific docs reference common docs for shared concepts. When adding features that apply to all SDKs, update the common docs. When adding SDK-specific implementation details, update the SDK-specific docs.

### Documentation Maintenance

When adding new services or features:
1. Update all four SDK `docs/SERVICES.md` files to maintain service parity documentation
2. Update the service count in `docs/SDK_OVERVIEW.md` if adding new services
3. Update the services table in each SDK's `README.md`
4. Ensure cross-references in `docs/CONCEPTS.md` files include all SDKs

## Cross-SDK Lessons Learned

### Prove a new test fails without the fix

A test that passes before the fix is worthless, and this is easy to get wrong with vectors. The
JSON-bridge determinism test originally fed the canonical JSON straight back — but protojson emits
map keys **already sorted**, so it passed against a non-deterministic encoder and proved nothing. It
only became a gate once the input carried the `properties` keys in reverse-sorted order, which is
also what an edited container actually looks like. Same rule for the file-driven gates: delete the
shared vector file and confirm the suite goes red.

**Three sharper variants of this, all found in the 2026-09-10 security-scan pass.** Each one had
a green suite sitting on top of a live vulnerability, so none of them was visible as a failure:

- **A test can assert the vulnerable behaviour.** `test_legacy_hash_without_contract_type`
  (Python) and the TS/Python create-address happy-paths *pinned* the defect: they asserted that an
  unsigned `contractType` comes back, and that a create reply with **no signature at all** yields
  an `Address`. Those had to be inverted, not adjusted. When a security fix leaves an existing
  test failing, read it before re-baselining — it may be the finding.
- **A test can exercise a private production method through an inline copy of it.** Java's
  legacy-hash tests re-implement the regexes with `String.replaceAll` because
  `computeLegacyHashes` is private (`CrossSdkCryptoVectorTest`,
  `WhitelistedAddressServiceLegacyHashTest`, `WhitelistVerificationFlowTest`). Fixing production
  and not the copies leaves the gate **green while asserting the old behaviour**. Grep for a
  duplicated implementation before trusting a test that covers a private method.
- **A test can name the property it does not test.**
  `TestCase2_LabelPatternDoesNotAffectMainLabel` (and its three peers) claims the main `label` is
  unaffected by the strip — but its fixture puts the main `label` *before* other fields, so the
  "followed by `}`" case it purports to rule out is never constructed. It passed against the
  vulnerable regex for the entire life of the defect.

And the reason the whole class survived: **the full `VerifyWhitelistedAddress` flow was only ever
called with nil inputs in the Go suite.** Every other address test drove steps 1-5 individually,
so nothing exercised step 6 with real signatures, and `VerificationResult.VerifiedAddress` had
exactly one reader in the module. If a verification flow's end-to-end test does not exist, no
amount of per-step coverage substitutes for it.

### A bulk rename can silently un-name a test, and the suite still says "passed"

Renaming `_contains_hash` → `contains_hash` with a plain string replace also rewrote
**`def test_contains_hash`** into `def testcontains_hash` — the test's own name contains the
old name as a substring. pytest stops collecting it. Nothing fails. Caught only because the
total moved 1456 → 1449.

This repo makes it likely: renames are applied across four SDKs at once, and the test names
mirror the function names.

- **Compare the test TOTAL after any bulk rename.** A drop is a regression even when everything
  reported passes. Same class as the TypeScript "Test suite failed to run" gotcha, which drops
  the count instead of showing red.
- Prefer word-boundary (`\b_contains_hash\b`) or AST/LSP renames over `str.replace`.
- Grep for the mangled shape afterwards: `grep -rnE "def test[a-z]"` catches a missing
  underscore in Python.

### `mvn … | tail` reports success for a failed build

Piping hides the exit code — `mvn compile -q 2>&1 | tail && echo OK` prints OK on a compilation
failure, because `tail` succeeded. The Java CLAUDE.md already says piping buffers output; the
worse problem is that it *masks the result*. Redirect to a file and check `$?`:

```bash
mvn compile -o -pl client -q > /tmp/build.log 2>&1; echo "exit=$?"
grep ERROR /tmp/build.log | head
```

### Deleting from a public surface: sweep all three test tiers, in every SDK

The unit gates do not see `integration/` or `e2e/`, so a removed public method leaves a
latent compile break behind a fully green run. Measured 2026-09-07, deleting
`WhitelistedContractService.ListWhitelistedContracts`: `go build ./...`,
`go vet ./pkg/...`, `golangci-lint run ./pkg/...`, `GOARCH=386` and all 7 unit packages were
green while `test/integration/extended_services_test.go` still called it — and the same stale
call sat in the TypeScript and Python integration tests.

**`go build` does not compile `_test.go` at all**, and `./pkg/...` excludes the test tiers, so
neither catches it. Per language:

| SDK | What actually compiles the test tiers |
|---|---|
| Go | **`go vet ./...`** — whole module, not `./pkg/...` |
| Java | plain `mvn test` already does; surefire excludes `*IntegrationTest`/`*E2ETest` from *running*, not from compiling |
| TypeScript | `tsc --noEmit` excludes `tests/` — run `npx jest tests/integration/<file>` and check it reports **skipped**, not "failed to run" |
| Python | nothing compiles it; `pytest --collect-only tests/integration` proves the module imports |

Then grep the docs: `<sdk>/build.sh docs --check` fails on a documented method that no longer
exists, but only in the SDK you re-ran it in.

### After changing any public API, regenerate the docs index before the cross-SDK gate

`./scripts/api-surface/generate.sh all check` fails with `generated method index is stale` if a
signature changed and `<sdk>/build.sh docs` has not been re-run. It is not a real parity
failure — run `./build.sh docs` in the affected SDK (or all four) and re-check.

### Naming shapes settled across all four (do not re-diverge)

- `RawCell.payload` — was `Bytes` in Go and `bytes` in TS.
- Every `RuleCell` variant exposes its discriminator: Go `Kind()`, Java `kind()`, Python/TS `.kind`.
- `cellFamily(RawCell)` returns the cell's own column type in all four; returning `""` discards the
  one piece of grammar information a raw cell carries.
- Governance reads are `getRules` / `getRulesById` / `getRulesProposal` / `getRulesHistory`.
- `verifyGovernanceRules(rules)` takes only the rules and returns them; the threshold comes from the
  service. Java keeps a two-argument overload, but passing a threshold that differs from the
  configured one verifies against something the rest of the service does not enforce.
- `comment` is required on `approveRulesProposal` / `rejectRulesProposal` in all four — it defaulted
  to `""` in Python, letting a governance change be approved with no rationale recorded.

### Governance reads verify in all four — RESOLVED, and the memo is why it is affordable

Was an open divergence: Go's `GetRules` returned the mapped DTO raw while Java, Python and
TypeScript verified, and `getRulesHistory` verified in TypeScript only. That asymmetry is
what let `tg-protect-mcpd` read governance rules out of an unauthenticated container.

Now, in all four: `getRules`, `getRulesById` and `getRulesHistory` verify the SuperAdmin
signatures before returning, and `getDecodedRulesContainer` verifies **unconditionally**
(Go's `if len(superAdminKeys) > 0` skip is gone, and so is the keyless
`NewGovernanceRuleService` that made the skip look necessary). TypeScript's
`GovernanceRuleServiceConfig` now REQUIRES keys and a positive threshold — the
"explicitly unverified service" it used to allow is gone.

Two things not to undo:

- **`getRulesHistory` is LENIENT, and the single-ruleset reads are STRICT.** A SuperAdmin
  key rotation makes every pre-rotation ruleset unverifiable, so a strict page would deny
  access to the whole audit trail from the rotation onwards. History excludes and NAMES
  what it dropped (`ExcludedUnverified` / `excluded_unverified` / `excludedUnverified`)
  with `totalItems` reduced, and does **not** error when nothing survives. Found by
  `tg-protect-mcpd`'s own test, which is where the key-rotation reasoning is written down.
- **The verification memo is not a second container cache.** `cache.RulesContainerCache`
  already holds the decoded+verified container for the address/asset/price paths (it
  fetches through `GetDecodedRulesContainer`). The memo covers only the three reads that
  return the RAW ruleset, which that cache does not serve, and it stores **outcomes** —
  successes only, so a failure resurfaces its error every call, and a re-signed container
  is a MISS. Do not add a third cache of the same bytes.
- **The memo key MUST be injective, in all four. This is the key that decides whether ECDSA
  runs at all.** It was `sha256(container ‖ 0x00 ‖ userId ‖ 0x00 ‖ sig ‖ …)` — a
  concatenation of variable-length, server-controlled strings with no length prefixes, so
  the boundary between the container and the signature list is not committed. A
  response-controlling attacker could therefore prime the memo with the genuine signed
  container (choosing `userId`, which **verification never reads**) and then serve a
  MODIFIED container that hashes to the same key: memo hit, no signature check, and an
  attacker-authored trust root returned as SuperAdmin-verified. Interleaving a `0x00`
  separator does **not** fix it — `(C, [(u,s)])` and `(C‖0x00‖u‖0x00‖s, [])` still collide.
  Only length prefixes do. Every SDK now writes an 8-byte big-endian length before each
  field, **drops `userId`** (a field verification ignores must not influence a
  verification-SKIP decision), and keys on the **decoded** container bytes so the key cannot
  identify something other than what was verified. An undecodable container is not
  memoised. **Gated cross-SDK** by the `memo_key` section of
  `scripts/resources/verification-behaviour-vectors.json` (8 pairs, `expect: distinct|same`),
  whose "distinct" cases each collide under one of the plausible separator-free encodings —
  so the section is red against all of them, not just the one that happened to ship. Loaders:
  Go `pkg/protect/service/memo_key_vectors_test.go` and Java
  `.../client/service/MemoKeyVectorsTest.java` (both in the **service** package, because the
  key function is unexported/package-private there), Python
  `tests/unit/helpers/test_verification_behaviour_vectors.py`, TS
  `tests/unit/helpers/verification-behaviour-vectors.test.ts` (which needed
  `rulesetVerificationKey` exported from the service module — deliberately **not** from the
  package barrel, same convention as `attestVerified`).
- **Base64 on any governance path is decoded STRICTLY.** Lenient decoders silently discard
  out-of-alphabet characters and absorb the rest, so a container carrying embedded
  separators decodes to the genuine bytes with **attacker bytes appended** — and protobuf
  treats concatenation as merge, so appended bytes ADD entries to repeated fields like
  `users`. That is how an attacker-chosen HSM or `PRICEUPDATER` key reaches the trust root,
  and it is the second half of the memo-key attack above. Go was always strict
  (`base64.StdEncoding`) and that is precisely why Go was not exploitable; the other three
  were not. Now: Python `taurus_protect/_strict_base64.py`, TS
  `src/helpers/strict-base64.ts`, Java `.../client/helper/StrictBase64.java` — each a
  dependency-free leaf so both `mappers` and `helpers` can use it, each stripping ASCII
  whitespace (line-wrapped base64 is legitimate) and rejecting everything else. Applied to
  the container decode, the signature-verification decode (so "the bytes verified" is
  unambiguous) and the memo key. One deliberate exception: `userSignaturesFromBase64`
  catches the error and returns an empty list, because its contract is "no signatures on
  parse failure" — which is fail-closed, since an empty list satisfies no threshold.
  **Validate without a quantified-group regex.** TypeScript's well-formedness test was
  `/^(?:[A-Za-z0-9+/]{4})*.../`, which recurses per four-character group in V8: a
  multi-megabyte input threw `RangeError: Maximum call stack size exceeded` out of the
  decoder instead of returning a decision (1.6 MB passed in 27 ms, 5.6 MB threw). Fails
  closed, so availability rather than a bypass — but a `RangeError` escapes every
  `catch (e) { if (e instanceof …) }` funnel in the SDK, and the decoder is reachable from any
  governance response through the signature-verification decode and the memo key, not only the
  size-capped container decode. It is a linear scan over the same accepted set now; Go, Python
  and Java were never affected. **And no SDK tests its own decoder** — the property is
  documented four times and pinned nowhere, which is the shape this repo keeps rediscovering by
  grep. TypeScript's `tests/unit/helpers/strict-base64.test.ts` is the accept/reject table to
  lift into a shared vector file; see `TODOS.md`.
- **`approveRulesProposal` PINS the content the approver reviewed.** It takes a required
  `expectedContainerHash` (`proposalContainerHash` of the reviewed proposal), re-fetches,
  and aborts **without signing** if the container differs. Without that pin a server able
  to shape responses answers the review call with the benign proposal and the re-fetch with
  a different container, obtaining a **genuine SuperAdmin signature over bytes of its
  choosing** — and `VerifyGovernanceRulesSignatures` checks ECDSA plus the distinct-key
  threshold and *nothing else*, so that container then verifies clean everywhere, including
  in independent and air-gapped verifiers that never trusted the compromised server. The
  threshold is no defence: each approver's re-fetch is independent, so the substitution
  simply repeats. This is **not** the same as verifying the proposal — a pending proposal
  legitimately carries 0..N signatures (signing IS the approval step), there is nothing to
  verify against, and `scripts/signing-sites/CLAUDE.md` rightly forbids "fixing" that.
  Verification and content-binding are different mitigations; the sites stay
  `signs-own-bytes`. Reviewing needs `decodeProposalForReview` — an explicitly-named,
  **unverified-by-design** decode — because `getDecodedRulesContainer` verifies
  unconditionally and so always fails on a pending proposal. Do not make the pin optional:
  an empty one restores the unpinned behaviour silently.
  The pin digests the **decoded** container bytes. Do not confuse it with the row-to-container
  label in the whitelist list paths, which is validatord's convention over the base64 **text**;
  both are SHA-256 digests of the same document, which is what makes the mix-up silent. The
  pin is client-side only and never reaches the wire.

### `RequestStatus` genuinely diverges — do not "align" it by guessing

Java has 40 values; Go, Python and TS have the same 38. Java-only: `FAST_APPROVED`, `INVALID`, `NEW`.
Only in the other three: `UNKNOWN` (an SDK-side sentinel). **Neither
`scripts/resources/swagger/apis.swagger.json` nor the proto schema declares the three Java-only
values**, so there is no authoritative list to align to: adding them to three SDKs copies one SDK's
unsourced list, and dropping them from Java could break a caller matching a status the server really
sends. Needs a server-side answer. Earlier reports claimed "43 values matching Java" in all four and
then listed 38 names, 23 of which were not in `RequestStatus.java` — do not trust that framing.

### Service Naming Consistency

When naming services, ensure consistency:

| Service | Java | Go | Python | TypeScript |
|---------|------|-----|--------|------------|
| Business Rules | BusinessRuleService | BusinessRuleService | BusinessRuleService | BusinessRuleService |
| Sharing | TaurusNetworkSharingService | TaurusNetworkSharingService | TaurusNetworkSharingService | SharingService |

Note: Python uses `snake_case` for file names (`business_rule_service.py`) but `PascalCase` for class names.

### Error Handling Alignment

**isRetryable() Pattern:** All SDKs should return `true` for retryable errors:
- HTTP 429 (rate limit)
- HTTP 5xx (server errors)

### RequestMetadataAmount — String Types

The API returns `valueFrom`, `valueTo`, `rate` as JSON strings for arbitrary-precision. All SDKs use string types with a `jsonValueToString()` helper for backward compatibility. Keep `decimals` as integer.

### Whitelisted Address/Asset Field Sourcing

**Key Rules (All SDKs):**
1. `WhitelistedAddressService.list()` MUST verify each envelope **leniently**: exclude and
   report an unverifiable row rather than failing the call, and error when rows came back but
   none survived (a filtered page must never read as an empty whitelist).
   **"Report" means on the RESULT, not only to the logger** — a caller cannot read the SDK's
   logger, so a logged-only exclusion still leaves a filtered page indistinguishable from a
   complete one. Every SDK returns a result carrying the exclusions: Go
   `WhitelistedAddressResult{Addresses, Pagination, ExcludedUnverified}`, TS
   `{items, pagination, excludedUnverified}`, Python `WhitelistedAddressListResult`, Java
   `getWhitelistedAddressesWithExclusions(...)` (an overload, so the `List`-returning forms
   keep compiling). Each exclusion carries `{id, reason}`.
   **`totalItems` MUST be reduced by the number excluded** — the server counts rows it
   returned, the caller receives only those that verified, so reporting the server's total
   makes `hasMore` promise a page that can never be fully read.
   Distinguish the two failure classes: a row-level integrity failure is excluded, but a
   `ContainerIntegrityError` (this SDK cannot interpret the rules container) aborts the whole
   call — it invalidates every row judged against that container, so excluding them one by one
   would empty the whitelist and report success.
2. `WhitelistedAssetService` must source `name`, `symbol`, `contract_address`, `blockchain`, `network` only from payload
3. If payload is missing a field, the result must be `None` (not DTO value)
4. Non-security fields (`status`, `action`, `rule`, `created_at`) can come from DTO
5. Remove any mapping methods that bypass the verified envelope path

**DTO field extraction:**
- `createdAt` — Extracted from trails array (find "created" action)
- `attributes` — Extracted from DTO attributes array (structured array in Go, key-value map in others)

### Test Infrastructure (Cross-SDK)

**Test hierarchy:** All SDKs use 3 tiers: `testutil/` (shared config + helpers), `integration/` (API tests), `e2e/` (multi-identity workflow tests).

**Config pattern:** All SDKs use `test.properties` file (key=value format) with env var overrides, supporting 6 identities (3 API users + 3 SuperAdmin public keys). Java's `testutil/TestConfig.java` is the reference implementation.

**testutil modules:**
| SDK | Location | Key files |
|-----|----------|-----------|
| Java | `client/src/test/java/.../testutil/` | `TestConfig.java`, `TestHelper.java` |
| Go | `test/testutil/` | `config.go`, `helpers.go`, `properties.go` |
| Python | `tests/testutil/` | `config.py`, `helpers.py`, `properties.py` |
| TypeScript | `tests/testutil/` | `config.ts`, `helpers.ts`, `properties.ts` |

**Integration config delegation:** Integration and E2E directories delegate to testutil (thin wrappers for backward compatibility). Don't duplicate config.

### Integration Test Patterns

- Use SDK-specific helpers (`getTestClientWithVerification()`) for SuperAdmin key config
- Test both positive cases (all fields match) and negative cases (one field differs)
- Fields to test for WhitelistedAddress: blockchain, network, address, customerId, label, memo, addressType
- Fields to test for WhitelistedAsset: blockchain, network, address, name, symbol, decimals, customerId, label

### ChangeService Alignment (comment→changeComment)

When creating changes via the API, the SDK `comment` field must be mapped to OpenAPI `changeComment` field. This applies across all SDKs:
- Java: MapStruct `@Mapping(source = "comment", target = "changeComment")`
- Go: Manual mapping in `CreateChange()`
- Python: `create_change_request_to_dto()` maps `comment` → `change_comment`
- TypeScript: `createChangeRequestToDto()` maps `comment` → `changeComment`

### BusinessRule Pagination (v2 API)

All SDKs must use the v2 API endpoint for listing business rules (`ruleServiceGetBusinessRulesV2`), which returns cursor-based pagination. The v1 endpoint uses offset-based pagination and is deprecated.

### Documentation Verification Quick-Checks

- `RequestMetadata.payload` — intentionally omitted from all SDKs (security)
- `RequestStatus` enum — uses CANCELED (not CANCELLED), PENDING (not PENDING_APPROVAL), BROADCASTED (not BROADCAST)
- Service naming — always singular: BusinessRuleService, WebhookCallService
- All SDKs use 6-step address verification (Step 6: Parse WhitelistedAddress from verified payload)
- Cross-references in common docs must include ALL 4 SDKs

## Shared swagger (`scripts/resources/swagger/apis.swagger.json`) is a STALE SNAPSHOT

One spec at the repo root feeds the Go, Python, Java and TypeScript generators, and it has drifted from
`tg-validatord/api/swagger/`. The generator emits only what the spec contains, so a live validatord endpoint
can be **completely absent from all four clients** with no error anywhere — it simply has no generated
method. Confirmed missing 2026-09-03: `PriceService_QueryPricesV2`, `GetPriceByID`, and the
`entities`/`fieldKey`/`fieldValue` parameters on `ChangeServiceGetChanges`.

**Grep the generated client before writing a builder call** — do not infer from validatord's proto or swagger,
and do not infer from another SDK. A plausible-looking method that does not exist costs an options field, a
mapping and a consumer before the compiler says so.

**Adopting a missing endpoint is a codegen job across four SDKs**, not a per-SDK fix:

- Patch **only** the new path plus its request/response schemas into `apis.swagger.json`, from validatord's
  per-service spec. A wholesale refresh pulls in every validatord change since the snapshot — an uncontrolled
  diff across all services, and this repo already carries a large uncommitted tree.
- Then either regenerate all four SDKs, or state plainly that the spec now leads them. Regenerating one
  leaves the other three built from a spec that has since moved, which is how they silently diverge.

Tracked in `TODOS.md`; the Go SDK's `CLAUDE.md` has the service-layer consequences.

## Shared Proto Schema (`scripts/resources/proto/schema/v1/`)

All four SDK proto generators (`taurus-protect-sdk-{go,java,python,typescript}/scripts/generate-proto.sh`) read from this directory. A missing source breaks every SDK simultaneously, because `authentication-service.proto` and `steward-service.proto` import from sibling files via plain names (no package path).

Convention for files copied in from `tg-validatord/api/proto/v1/`:
- Strip `option go_package = "github.com/taurusgroup/tg-validatord/...";`
- Add `option java_package = "com.taurushq.sdk.protect.proto.v1";`
- All other content (messages, field numbers) stays byte-for-byte identical.

The dir must contain at minimum: `apikey.proto` (defines `ApiKey`, `ApiKeyToken`) and `credentials.proto` (defines `Credentials`, `SAMLAuthRequest`, `SAMLAuthRedirect`, `SAMLSession`, `OIDCLocation`, `OIDCSession`). Both are imported by sibling protos but are easy to omit when copying schemas across — if you ever see `undefined: <Type>` errors for these names in any SDK, the proto source is missing here.

The shared swagger at `scripts/resources/swagger/apis.swagger.json` is independent of this dir and contains the openapi-side definitions for the same domain (e.g., `tgvalidatordCredentials`, `StewardServiceCreateApiKeyBody`).

## Typed Authorization Error (cross-SDK)

All four SDKs raise a typed 403 carrying the roles that would satisfy the failed check
(`AuthorizationError.requiredRoles`, Java `AuthorizationException.getRequiredRoles()`), so a caller can say
which role to ask for instead of only "forbidden". Empty when the denial was not role-based (disabled
endpoint, visibility restriction). One entry = a required role; several = any one suffices.

**Alignment gate — `scripts/resources/authorization-error-vectors.json`.** `{description, message, expected_roles}`
vectors loaded by all four unit suites, so a server wording change is caught once for everyone.

Parser contract, identical everywhere: search (never anchor) for `one of the '(...)' role is required`, split
the capture on `" - "`, trim, drop empties. **The pattern must stay unanchored** — validatord wraps the gRPC
status, so the payload really arrives as `pre-filter failed: common filters failed while constructing
authcontext: rpc error: code = PermissionDenied desc = one of the '<roles>' role is required`. An anchored
pattern matches nothing and silently yields zero roles. That exact string is the first vector.

Both of validatord's role checks (all-of and any-of) emit the same "one of" wording, so the count carries the
semantics, not the prose. Role names are lowercase alphanumeric, which is what makes `" - "` unambiguous.

**Consumers must not echo the server message.** A role list is safe to surface; the surrounding message is
server-controlled text. `tg-protect-mcpd` renders roles through a `^[a-z0-9]+$` allowlist because its tool
results reach a model API.

## Governance Rules Typed API (cross-SDK)

All four SDKs expose the same typed governance-rules surface: a lossless `DecodedRulesContainer`, a typed `RuleCell` discriminated union (36 cell types across 9 column families + a `RawCell` fallback), a bidirectional cell codec, a container encoder, and the proposal lifecycle `updateRulesProposal(container)` / `approveRulesProposal(privateKey, comment)` / `rejectRulesProposal(comment)`.

**Alignment gate — `scripts/resources/governance-cell-vectors.json`.** 39 golden vectors (`{description, column_type, cell_type, typed_value_json, wire_base64}`) produced by the Go SDK. Every SDK's unit suite loads this one file and asserts each typed cell encodes to the exact recorded bytes and decodes back to the same typed value — this is how cross-SDK byte-parity is enforced, not asserted.

Shared invariants (identical in all four SDKs — keep them so):
- `*Any` cells are first-class typed cells; their wire form is the protobuf zero value (empty cell). No special nil-casing/removal.
- Schema-evolution is lossless: unknown protobuf fields are preserved per node and re-attached on encode; unknown cell types decode to `RawCell` via a decode→re-encode→byte-compare guard; unknown enum values pass through numerically on the wire, but the encoder errors on caller-authored unknown enum *names*.
- The encoder strips server-controlled `enforcedRulesHash` + `timestamp`.
- Integer cells = big-endian magnitude payload + the sign selects the `Value`/`NegValue` enum arm; `StringEqual`/`BytesEqual` payloads are raw scalar bytes.
- `updateRulesProposal` returns void/error (the rpc returns `Empty`; the ~60s rules cache makes read-back unreliable). `approve` signs the **decoded** bytes of the pending proposal with the existing P-256 signer.

**Governance unit-test layout (each SDK):** a codec branch/error suite + a container round-trip suite (**assert per-node unknown-field preservation, not just the top-level container** — a top-level-only assertion masks a broken nested-preservation encoder) + the shared golden-vectors suite + a proposal-lifecycle suite + model find\*/HSM-helper tests. When adding a cell type or column family, regenerate the golden vectors (Go produces them) and update every SDK's codec test in lockstep or the vectors gate fails.

**Cross-SDK parity vectors beyond the golden-vectors file.** `governance-cell-vectors.json` covers
*canonical* cell encodings only — every one of its 39 vectors was produced by Go, so it is blind to
anything a decoder might silently rewrite. Five further scenarios (8 vectors) are pinned by
`scripts/resources/governance-lossless-vectors.json`, loaded by all four round-trip suites
(`rule_cell_codec_lossless_test.go`, `test_rules_container_roundtrip.py`,
`rules-container-roundtrip.test.ts`, `RulesContainerRoundtripTest.java`):

| Scenario | What it catches |
|---|---|
| `RuleSource` with an unknown field / a payload on an `Any` arm / an unknown field inside the payload | a decoder dropping schema-newer source data |
| unknown enum values (role 201, columnType 77, subDomain 202, blockchain 203) | an enum collapsing to the zero value, which rewrites a column family or widens a rule |
| ≥3 `properties` map entries, encoded 20× | non-deterministic map ordering |
| unknown fields inside all four nested rule-detail sub-messages | nested preservation, which no top-level assertion sees |
| a non-UTF-8 `RuleStringEqual` cell + a truncated `RuleFiatAmount` | one bad cell aborting the whole container |

Non-canonical cases cannot be produced by the Go encoder, which is why they are a separate file from
the cell-vector gate. To add one: build the bytes with the Python protobuf runtime (it emits
arbitrary/malformed wire forms most easily), append an entry, and bump the count constant in **all
four** loaders — `losslessVectorCount` (Go `pkg/protect/mapper/lossless_vectors_test.go`),
`VECTOR_COUNT` (Python `tests/unit/mappers/lossless_vectors.py`, path walk is `parents[4]`; TS
`tests/unit/mappers/lossless-vectors.ts`; Java `.../client/mapper/LosslessVectors.java`).

All four hard-fail on a missing file. Re-verify that after touching a loader by deleting the file and
re-running — **Go needs `-count=1`** or the test cache reports the previous pass.

**In TypeScript, load the vectors inside `it()`, never in a `describe` body.** A throw during
collection becomes "Test suite failed to run", which drops the test count instead of showing a red
test. Jest also has no per-assertion message argument (that is Vitest) — to keep a loop's failure
identifiable, assert a tuple: `expect([description, actual]).toEqual([description, expected])`.
