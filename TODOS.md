# TODOs

Deferred work from the 2026-09-03 cross-SDK alignment pass (see
[`docs/SDK_ALIGNMENT_REPORT.md`](docs/SDK_ALIGNMENT_REPORT.md)). Each entry carries enough
context to pick up cold.

---

## Java: give the verified envelopes an unforgeable witness

**What:** Have `WhitelistedAssetVerifier` / `WhitelistedAddressVerifier` (once extracted —
see the entry below) return a `VerifiedWhitelistedAsset` / `VerifiedWhitelistedAddress`
value with private fields, and have the envelopes hold that instead of exposing a public
`markVerified`.

**Why:** the *exploitable* half is already closed. `setVerifiedWhitelistedAsset(asset)` and
`setVerifiedWhitelistedAddress(address)` were public and flipped the same `isInitialized`
gate the getters check, so a caller could construct an envelope, inject a fabricated
value, and read it back with no exception — a "verified" marker returning attacker-chosen
data. Both are now `markVerified(rulesContainer)`, which **derives** the value from the
envelope's own signed payload, so forging the marker is possible but useless (the property
Go's `helper.VerifiedAsset` has). What remains is that `markVerified` must be *public*,
because the services live in another package — so the marker can still be flipped, just
not usefully.

**Context:** `SignedWhitelistedAssetEnvelope.markVerified` and its address peer carry the
reasoning inline. `VerifiedEnvelopeDerivationTest` pins the derivation, including that a
payload-less envelope cannot be marked and that the getter still refuses an unmarked one.
Go is the model to copy (`helper.VerifiedAsset`, unexported fields); TypeScript's
`Verified<T>` brand is the compile-time equivalent. Python has no equivalent and could
only get a runtime assert.

**Depends on:** the verifier extraction below — a witness type is only unforgeable if the
verifier is the sole thing that can mint it, which needs it out of the services.

---

## Seal Java's `ProtectClient.getOpenApiClient()` and TypeScript's deep imports

**What:** Deprecate and then remove Java's public `getOpenApiClient()`, and add an
`exports` map to the TypeScript `package.json`.

**Why:** these are the two ways a caller reaches the generated client and skips every
verifying service, and both are sanctioned by the current shape of the packages. Java's
`ProtectClient.getOpenApiClient()` (`ProtectClient.java:374`) is public and its javadoc
advertises *"allowing for raw API access"*, so any holder can call
`new GovernanceRulesApi(client.getOpenApiClient()).ruleServiceGetRules()` and read the
trust root unverified. TypeScript ships `dist/internal/openapi/` with no `exports` map, so
`import ... from '@taurushq/protect-sdk/dist/internal/openapi/apis/GovernanceRulesApi'`
resolves — which defeats the already-made decision to keep `governanceApi()` private
(that getter is the one deliberately-private accessor among 61).

**Context:** Go is unreachable by construction (`internal/openapi`, compiler-enforced) and
Python's `_internal` is convention-only but not re-exported at the package root. The
`exports` map is a few lines and breaks no legitimate consumer; Java's getter is a real
API decision because removing it breaks anyone relying on documented raw access. Both were
declared out of scope for the 2026-09-07 verification pass — the front door was hardened
while these stayed open.

**Depends on:** nothing for TypeScript. Java's needs a deprecation cycle.

---

## `MultiFactorSignatureService.approve` takes an opaque caller signature

**What:** Decide whether the MFA approval should re-read and sign like the request and
whitelist paths now do, or stay a pass-through — and align the four SDKs either way.

**Why:** `approveMultiFactorSignature(id, signature, comment)` forwards a caller-supplied
signature string with only a non-empty check (Go `multi_factor_signature.go:71`, Java
`MultiFactorSignatureService.java:120`, TS `multi-factor-signature-service.ts:155`). That
is the same shape the deprecated `ApproveWhitelistedContract` had before the verified
asset approval replaced it. It may well be correct — the trust model could be
"signed externally, like air-gap" — but nothing in the code says so, and Python does not
even have this method: it exposes `create_challenge` / `verify_challenge` instead
(`multi_factor_signature_service.py:136`), so there is no cross-SDK shape to compare.

**Context:** the 2026-09-07 signing-site sweep classified every ECDSA signing site as
`verifies` or `signs-own-bytes` (`scripts/signing-sites/manifest.json`). MFA is not in
that manifest because the SDK does not sign there — the caller does. That is exactly why
it needs a decision rather than a patch: if the answer is "externally signed by design",
say so in the javadoc and keep it; if not, it needs the re-read treatment.

**Depends on:** a server-side answer on who is supposed to produce that signature.

---

## TaurusNetwork `ProofOfOwnership` is never verified anywhere

**What:** Determine whether `signedPayloadHash` / `signedPayloadAsString` on a shared
address can be verified client-side, and if so verify it on the sharing read paths.

**Why:** the field pair is the same shape as every other signed artefact in this SDK, and
it is straight-copied and never compared in any of the four. Go maps it
(`model/taurusnetwork/shared.go:83-89`) and never checks it; TypeScript the same
(`sharing-service.ts` field mapping); Java does not model it at all; Python defines
`proof_of_ownership` on `SharedAddress` but the mapper never populates it, so it is always
`None` — dead wiring on top of an unverified field.

**Context:** verifying it needs the OWNER participant's public key, and whether an SDK
consumer can obtain that is unresolved — which is why the 2026-09-07 pass recorded it
rather than guessing. Start by establishing where that key comes from (the participant
record? the governance container?); the answer decides whether this is a verification job
or a documentation one. Python's dead field is worth fixing either way.

**Depends on:** the key-provenance answer.

---
## 1. Seal TypeScript's remaining low-level `*Api` getters

**What:** Make the ~60 generated `*Api` getters on `ProtectClient` non-public, or narrow
them to the endpoints the high-level services genuinely do not wrap.

**Why:** A raw `*Api` call skips the DTO mapping and, on the security paths, the signature
verification the services perform. The governance getter was the acute case and is already
sealed: `client.governanceRulesApi.ruleServiceGetRules()` returned an unverified DTO and
`ruleServiceUpdateRulesProposal({rulesContainer})` accepted an arbitrary base64 blob, so
client-side verification was opt-out in this SDK alone. The same shape applies, less
sharply, to whitelist / request / change writes: `addressWhitelistingApi` and
`requestsApi` can both bypass a verified path that exists right next to them.

**Context:** Go makes this impossible by construction — its generated client lives under
`internal/openapi`, which the compiler forbids external packages from importing. Java's
`openapi` module is `provided` scope, so it is not transitively exposed. Python's is
`taurus_protect/_internal/openapi` (convention only). TypeScript is the outlier.

Start at `src/client.ts`: `governanceApi()` (formerly `get governanceRulesApi()`) is the
pattern to copy — private method, `@internal`, with a comment saying why. The tests that
enumerate getters are `tests/unit/client/protect-client.test.ts` (`apiGetters` array and
the count assertion), plus `governance low-level API is not publicly reachable`, which is
the assertion shape to reuse.

Note this is a breaking change for any consumer using a raw API, and
`taurus-protect-sdk-typescript/README.md` documents low-level access as a feature — that
paragraph needs to change with it.

**Depends on / blocked by:** Nothing. The governance seal (done) is the precedent.

---

## 2. Java's TaurusNetwork services are systematically thinner than the peers

**What:** Implement the missing TaurusNetwork methods in Java.

**Why:** These are not naming differences — the operations are absent, so a Java caller
cannot perform them at all. Sharing is the starkest: Java has `listSharedAddresses` and
nothing else, so a Java caller cannot share or unshare an address or an asset.

**Context:** `scripts/api-surface/diff.py` reports the spread (regenerate with
`./scripts/api-surface/generate.sh all`):

| Service | Go | Java | Python | TypeScript |
|---|---|---|---|---|
| lending | 14 | **4** | 16 | 14 |
| pledge | 14 | **3** | 14 | 14 |
| settlement | 6 | **2** | 6 | 6 |
| sharing | 6 | **1** | 6 | 6 |

Use Go or TypeScript as the model — they agree closely. The Java service classes are in
`client/src/main/java/com/taurushq/sdk/protect/client/service/TaurusNetwork*Service.java`
and the generated APIs they need are already present in the `openapi` module (the
`taurusNetworkService*` operations), so this is wiring plus mappers, not code generation.

`docs/SDK_ALIGNMENT_REPORT.md` used to claim SharingService was "Aligned" across all four;
that row is now corrected, so the report and this TODO agree.

**Depends on / blocked by:** Nothing, but do it before promoting any TaurusNetwork
behaviour into a shared vector gate.

---

## 3. Promote the rich-container fixture to a shared vector

**What:** Move the hand-built "rich container" fixture out of the four round-trip suites
into `scripts/resources/`, the way the eight lossless vectors now are, and assert its
encoded bytes are identical across SDKs.

**Why:** It is the most complex container the suites build — users, groups, thresholds,
transaction rules with nested EVM call-contract details, address and contract whitelisting
rules, properties maps. Cross-SDK byte equality on it would cover far more of the encoder
than the current vectors do, and it is exactly the shape a real tenant's rules take.

**Context:** This was blocked until now because the fixtures were not the same container:
Java set `blockchain` to `"Ethereum"` where Go, Python and TypeScript used `"ETH"`. That is
fixed (`RulesContainerRoundtripTest.java`), so the four fixtures should now encode
identically — verify that first, because any remaining field-level difference will show up
as a byte mismatch and needs resolving before the vector is trustworthy.

The pattern to follow is `scripts/resources/governance-lossless-vectors.json` plus the four
loaders (`lossless_vectors_test.go`, `tests/unit/mappers/lossless_vectors.py`,
`tests/unit/mappers/lossless-vectors.ts`, `mapper/LosslessVectors.java`). Each loader
already asserts a vector count, so bump it in lockstep.

**Depends on / blocked by:** Confirming the four fixtures now produce identical bytes.

---

## 4. Revive or retire the Python lint and Java PMD gates

**What:** Decide, and record the decision, whether these two gates are gates.

**Why:** Both are red on committed `main` and have been ignored for long enough that every
alignment pass rediscovers them, spends time establishing they are pre-existing, and skips
them. Either state is fine; the ambiguity is what costs time.

**Context:** Measured baselines (see the per-SDK `CLAUDE.md` files):

- **Python** `./build.sh lint`: 96 files fail flake8 (mostly E501 at the 100-char limit),
  79 would be reformatted by black, mypy reports 3008 errors in 124 files. Reformatting the
  package to chase it is a repo-wide project that would bury any review in churn.
- **Java** `mvn pmd:check -pl client`: 10 violations on master. Checkstyle, by contrast, is
  a real gate at 0 violations, and SpotBugs is at 0 (note `spotbugs:check` must run
  **without** `-o`; `findsecbugs-plugin` is not in the local `~/.m2` cache).

The realistic middle path, already the working convention, is to keep *new* files free of
real findings (F401 dead imports, F841 unused locals) and leave E501 alone. If that is the
decision, encode it — a flake8 config with E501 disabled and a PMD baseline file — so the
command's exit code means something.

**Depends on / blocked by:** Nothing. Independent of the other three.

---

## GetAddresses attribute filters and NFT selector

**What:** Wire `attributeFiltersJson`, `attributeFiltersOperator` and `nfts` from the
generated `ApiWalletServiceGetAddressesRequest` into `model.ListAddressesOptions`.

**Why:** `GetAddresses` exposes 33 client parameters; the service layer now plumbs the
identity, balance and score filters, leaving these three the only reachable-but-unexposed
ones. A filter the SDK drops forces every consumer to over-fetch and narrow in memory —
which is what motivated the wider filter pass in the first place.

**Context:** Deliberately excluded from that pass. `attributeFiltersJson` takes a JSON
blob whose shape is not described anywhere in the swagger, so exposing it needs a schema
decision of its own: either a typed Go struct mirroring validatord's
`attributes_model.AttributeFilter` (which pins the SDK to an internal type) or a
pass-through `json.RawMessage` (which pushes the shape problem onto the caller and gives
an MCP tool nothing to put in a `jsonschema` description). Neither is obviously right, and
the balance and score filters — which ARE wired — cover the operator-facing questions.

Note the deprecated flat score parameters (`scoreProvider`, `scoreInBelow`, `scoreOutBelow`,
`scoreExclusive`, `coinfirmScoreGreater`, `chainalysisScoreGreater`) are excluded on purpose:
`wallet-service.proto` marks all six `deprecated = true` with "Use scoreFilter instead", and
the generated client still exposes both generations, so "add the score filters" naively wires
six dead fields. `AddressScoreFilter` is a pointer precisely so nil omits the whole group.

**Depends on / blocked by:** A decision on the JSON-blob shape. Independent of the codegen work.

---

## Extract Java's whitelist verification into verifier classes

**What:** Create `WhitelistedAddressVerifier` / `WhitelistedAssetVerifier` under
`client/.../helper/`, matching Go/Python/TypeScript, and reduce `WhitelistedAddressService`
(~870 lines) and `WhitelistedAssetService` (~600 lines) to orchestration.

**Why:** Java is the only SDK that inlines verification in its services, and that is *where
the drift happened*. Four non-constant-time `List.contains` comparisons survived at
`WhitelistedAddressService:373,617` and `WhitelistedAssetService:411,553` while the correct
`SignatureVerifier.verifyHashCoverage` sat unused in another file with **zero callers** — the
same shape as the earlier threshold-core bug this repo already records.

**Context:** The 2026-09-04 verification pass relocated only the crypto helpers, to keep the
diff proportionate during an already-wide change. The structural cause remains. The alignment
report now records the layering as an accepted difference *with this TODO as the plan to stop
accepting it*.

**Depends on:** nothing further — the hash-function consolidation it needed has landed.

---

## Revisit address/asset verifier duplication

**What:** Evaluate a shared 5-step skeleton for the address and asset verifiers: roughly 400
duplicated lines per SDK (Go 520+383, Python 522+486, TypeScript 772+618), ~1600 total.

**Why:** Every divergence found in the 2026-09-04 pass had to be checked twice per SDK, once
per entity — and several had been fixed in one and not the other. Python carried two different
`_contains_hash` implementations in a single SDK for exactly this reason.

**Context:** Deliberately **not** done then: the flows genuinely differ (6 vs 5 steps,
different legacy-hash functions, different rule container types), so unifying would have been
premature abstraction. Revisit now that `verification-behaviour-vectors.json` pins the shared
primitives — the gate is what makes the refactor safe rather than speculative.

**Depends on:** the verification behaviour gate staying green (it is).

---

## Adopt the whitelisted-address export endpoint

**What:** Wire `WhitelistService_ExportWhitelistedAddresses` into a service method in all four
SDKs.

**Why:** `TgvalidatordExportWhitelistedAddressesReply` exists in the generated models of all
four, but no service wraps it — the capability is present in the client and unreachable from
the SDK.

**Context:** `scripts/resources/swagger/apis.swagger.json` is a stale snapshot shared by all
four generators, so this is a codegen job, not plumbing: patch only the new path plus its
schemas from validatord's per-service spec, then regenerate **all four** (regenerating one
leaves the others built from a spec that has since moved). Same procedure as the other
known-missing endpoints recorded in the root `CLAUDE.md`.

**Depends on:** nothing; independent.

---

## Confirm and fix Python per-signature PEM decoding

**What:** Determine whether `whitelisted_address_verifier.py` re-decodes a user's PEM public
key on every signature inside the group loop, and add a per-group key cache if so.

**Why:** TypeScript already caches within a group (`whitelisted-address-verifier.ts:560`). PEM
parsing per signature per group per OR-path multiplies with governance complexity, right next
to the ECDSA verify.

**Context:** Raised during the 2026-09-04 pass and left **unresolved rather than asserted** —
the grep was inconclusive and the loop was not read. The signal that prompted it is the
skipped-reason string `"failed to decode public key for user '…'"`, which implies decoding
happens inside the loop.

**Depends on:** nothing; a short check.

---

## Integrity failures must not be retryable

**What:** Stop mapping integrity/configuration failures onto retryable `ServerError(500)`.
TypeScript: `WhitelistError` and `ConfigurationError` in `errors.ts:452` /
`services/base.ts:99-101`. Python: `IntegrityError` in `errors.py:247`, reached through the
catch-all funnels in `address_service.py`.

**Why:** `isRetryable()` is the documented pattern for callers, and every SDK's docs tell them
to branch on it. A suspected forgery presented as a transient server fault means a caller
following the documentation retries it — and each retry is another chance for a different
response to slip through. It also hides the finding: an integrity failure logged as a 500 does
not look like an attack.

**Context:** Found during the 2026-09-07 verification review and deliberately deferred: it is
an error-taxonomy change touching every service in two SDKs, so it wants its own pass rather
than riding along with security fixes. Go and Java already keep the types distinct
(`model.IntegrityError` / `IntegrityException` never funnel into the API-error mapper). One
funnel WAS fixed in passing: `asset_service.py`'s three `except` blocks, because the new
per-address HSM verification made them reachable from a path that previously could not raise
`IntegrityError` at all. Use that as the shape — widen the `isinstance` tuple so the typed
error propagates instead of being re-wrapped.

**Depends on:** nothing. Best done with the two items below, which are the same taxonomy.

---

## Container failures must abort; row failures must exclude

**What:** Make the two failure classes behave consistently. A row-level integrity failure is
excluded and reported; a `ContainerIntegrityError` — this SDK cannot interpret the rules
container — aborts the whole call. Two SDKs get it wrong in **opposite** directions:
TypeScript swallows `ContainerIntegrityError` in the list catch-all
(`whitelisted-address-service.ts:339-344`), and Java's unchecked `IntegrityException` escapes
`catch (WhitelistException)` (`WhitelistedAddressService.java:882`) and aborts a whole page
over one bad row.

**Why:** A container this SDK cannot interpret invalidates every row judged against it, so
excluding them one by one empties the whitelist and reports success — the worst outcome
available, because a caller plans around a confident empty answer. Conversely, aborting a page
over one unverifiable row denies access to every good row, and listing is how an operator
finds the bad one.

**Context:** The rule is already written in the root `CLAUDE.md` ("Whitelisted Address/Asset
Field Sourcing") and implemented correctly in Go and Python. The Java half needs the exception
hierarchy decided first: `IntegrityException` is unchecked and `WhitelistException` is checked,
which is why one escapes the other. Deferred from the 2026-09-07 review as a Medium.

**Depends on:** a decision on whether Java's `IntegrityException` becomes a subtype of
`WhitelistException` or the catch sites are widened.

---

## List-result parity: exclusions, totals, and the all-dropped case

**What:** Three related gaps in the request/asset list paths. (1) Only Go reports request-list
exclusions **on the result**; Java, Python and TypeScript log them where no caller can read
them. (2) Python reports the server's unadjusted `total_items` after excluding rows. (3) No SDK
errors when rows came back from the request endpoint and none survived — that guard exists only
on the whitelist path (`whitelisted_address.go:237`).

**Why:** A caller cannot read the SDK's logger, so a logged-only exclusion leaves a filtered
page indistinguishable from a complete one. An unadjusted total makes `hasMore` promise a page
that can never be fully read. And an entirely unverifiable page returning zero rows with a nil
error is the same confident-empty-answer failure as the item above. `tg-protect-mcpd` already
depends on the Go behaviour and surfaces `excluded_unverified_ids` from it, so the shape to
copy exists.

**Context:** Deferred from the 2026-09-07 review because it is a result-shape change in three
SDKs, and the whitelist paths show the target shape: a result object carrying
`{id, reason}` exclusions with `totalItems` reduced by the count excluded. Note the reduction
is right for `totalItems` and **wrong** for `hasMore`, which is a server-side pagination fact.

**Depends on:** nothing, but it is a breaking return-type change in Java/Python/TypeScript, so
it wants the same release as the two items above.

---

## Trust anchors for ProofOfReserve and TaurusNetwork SharedAddress

**What:** Decide, with validatord, what these two flows can actually prove — then either
implement the check or document the scope honestly. Today neither is verified and the SDKs
present both as ordinary reads.

**Why:** They look verifiable and are not, which is worse than an obviously unverified read.
`ProofOfReserve`'s challenge response proves possession of the key **the response supplied**,
not that the key controls the address — binding them needs per-chain pubkey-to-address
derivation the SDKs do not have, plus a challenge the *caller* chooses rather than the server.
TaurusNetwork `SharedAddress` carries the owner participant's key in the same untrusted
response, so any check is a consistency check, not an authenticity one; an anchored participant
key would have to come from somewhere the response cannot influence.

**Context:** Both were examined in the 2026-09-07 verification review and skipped rather than
fixed, because no amount of client-side work closes them without a server-side answer. Without
this entry the next reader has to re-derive that. The `SharedAddress` model also diverges four
ways across the SDKs, so aligning it is a prerequisite for any shared gate.

**Depends on:** a validatord decision on (a) caller-chosen challenges for proof-of-reserve and
(b) an out-of-band source for participant keys.

## Rules-container cache is shared across callers under a bearer token provider

**What:** Key the rules-container cache per credential identity, or otherwise stop one
client's cache from serving containers fetched with another caller's token.

**Why:** `BearerTokenProviderCredentials` exists so one client can serve many callers, each
carrying its own token (`BearerTokenProvider` takes a `ctx` and `bearerTransport.RoundTrip`
resolves per request). The rules-container cache is a single unkeyed slot. On a cache MISS the
fetch uses the calling context's token; on a HIT no request is issued at all, so the next
caller receives the container the first caller fetched and their own authorization for
`GetRules` is never exercised.

**Context:** Not a cross-tenant leak today. `tg-protect-mcpd` —
`internal/cmd/provider/tenant_registry.go` — already builds one client per tenant, each with
its own SuperAdmin keys and its own cache, LRU-bounded at 256, and its comment calls that the
multi-tenant verification boundary. Governance containers are tenant-wide and
SuperAdmin-signed, so what users of one tenant share is the same document. The residual gap is
intra-tenant: a user without the governance-read role can still obtain the container from a
warm cache.

Per-user keying was **considered and rejected on 2026-09-04**: the cache is consulted on every
address verification (Go `service/address.go`), so keying per user would require every user to
hold `GetRules` authorization merely to verify an address.
`tg-validatord/api/proto/v1/rule-service.proto` carries no role annotation on `GetRules` and
validatord enforces roles in the authcontext filter layer, so whether that would 403 ordinary
users is unconfirmed — which is why address verification was not bet on it. The caveat is
documented in all four `docs/AUTHENTICATION.md` and on Go's `Client.RulesCache()`.

**Depends on:** confirming validatord's role model for `GetRules`. If it is not role-gated,
per-credential keying becomes safe and this can proceed; if it is, any scheme must keep
per-user authorization off the address-verification path.

---

## Two encoder outliers in the governance RuleSource path

**What:** (a) Java `RulesContainerMapper.ruleSourceToBytes` treats any non-null `raw` as a
verbatim passthrough; (b) Python `_rule_source_to_bytes` is an if/elif chain with no final
`else`.

**Why:** (a) An empty-but-non-null `raw` makes Java emit an empty source cell, which decodes
as `RuleSourceAny` — "match any source" — silently widening a whitelisting rule in a container
that is then signed. (b) A source whose `type` this SDK does not model serialises as a bare
`{type: N}` with its payload dropped, so an address restriction vanishes from a signed
container.

**Context:** Go (`len(s.Raw) > 0`), TypeScript (`s.raw !== undefined && s.raw.length > 0`) and
Python (`if s.raw:`) all require non-empty before the passthrough; Java alone checks only
non-null. On the encoder side Go returns `unknown rule source type %d` and TypeScript throws;
Python does neither. Neither defect is reachable by `governance-cell-vectors.json`, which is
canonical-by-construction (all 39 vectors were produced by the Go encoder), so a new
non-canonical vector in `governance-lossless-vectors.json` is the place to pin the fix —
`explicit_empty_payload_noncanonical` is the worked example, and the count constant must be
bumped in all four loaders in lockstep.

**Depends on:** coordinate with the lossless guard added to `ruleSourceFromBytes` in Go, Python
and TypeScript — the same functions are involved, so do both in one pass rather than touching
them twice.

---

## Stale shared swagger — the service-layer consequences

**What:** Decide, per missing endpoint, whether to patch
`scripts/resources/swagger/apis.swagger.json` and regenerate all four SDKs, or to state
plainly that the spec now leads them.

**Why:** The generator emits only what the spec contains, so a live validatord endpoint can be
**absent from all four clients with no error anywhere** — it simply has no generated method. A
`List*Options` gap is then un-fillable by hand, and the absence looks like an SDK design
choice rather than a codegen lag.

**Context:** Confirmed missing 2026-09-03: `PriceService_QueryPricesV2`, `GetPriceByID`, and
the `entities` / `fieldKey` / `fieldValue` parameters on `ChangeServiceGetChanges`. The
consequences reach the service layer:

| Missing | Consequence |
|---|---|
| `PriceService_QueryPricesV2`, `GetPriceByID` | `ListPrices` is stuck on `GetPrices`, which takes `google.protobuf.Empty` — no filters, no pagination at all |
| `Entities`, `FieldKey`, `FieldValue` on `ChangeServiceGetChanges` | `ListChangesOptions` can only reach 8 of the endpoint's 11 fields |

Grep the generated client before writing any builder call — a plausible-looking method that
does not exist costs an options field, a mapping and a consumer before the compiler says so.
Do **not** do a wholesale spec refresh: that pulls in every validatord change since the
snapshot, an uncontrolled diff across all services, on top of an already-large uncommitted
tree. Patch only the new path plus its request/response schemas.

**Depends on:** a decision on regenerating all four SDKs together vs recording the spec as
ahead of them. The "Adopt the whitelisted-address export endpoint" entry above is the same
codegen job and should be batched with it.

---

## A uniformly stale rules container is undetectable by the SDKs — needs a server-side answer

**Status:** open. Intra-page mixing is closed; this is the part that cannot be.

All four SDKs now recompute the normalized-list container label
(`base64(SHA256(<the base64 container text>))`, validatord's convention in
`internal/api/v1/whitelist-controller.go`) and reject a row whose `rulesContainerHash` names
bytes that do not hash to it. That stops a server filing container A under container B's label
and steering individual rows to a different, weaker ruleset within one page.

**It does not stop the uniform case.** The rules container arrives *in-band* with the whitelist
response, so a hostile or compromised validatord can serve one older-but-still-validly-signed
container for **every** row. There is no mismatch to detect: every label is self-consistent,
every container passes the SuperAdmin threshold, and every row verifies cleanly against a
ruleset whose group thresholds may be lower than the ones currently enforced. All four SDKs
report success.

**Why the SDKs cannot close it today.** The only exact test is to compare against the currently
enforced ruleset from `GetRules`, and that is racy — a ruleset can be promoted between the two
calls — because validatord exposes **no ruleset identity to pin against**: the `Rules` message
returned by `GetRules` carries `rulesContainer`, `rulesSignatures`, `locked`, `creationDate`,
`updateDate` and `trails`, but **no `id` and no `enforced` flag**
(`internal/api/v1/rule-controller.go`), while the server tracks both internally (`rules` table,
`enforced BOOL`, index `rules_tenantid_enforced_idx`). A naive cross-check would reject
legitimate mid-rotation pages.

**`enforcedRulesHash` does not solve it and must not be used for this.** It is a *backlink* set
on a proposal to point at the ruleset it was diffed against, so after promotion a container
carries its predecessor's hash, not its own. It also uses a **different preimage** —
`base64(SHA256(raw protobuf))` rather than the base64 text
(`pkg/rule/store/cockroach/rule-storage.go`) — and both are 44-character base64 SHA-256
digests, so confusing the two is silent. Despite the comment claiming the backend checks it,
its only two non-test references in validatord are the writes.

**Depends on:** validatord exposing a ruleset identity on `GetRules` and on the whitelist
envelope — an `id`, a generation counter, or an `enforced` flag — so a client can assert that
the container it was handed is the one currently in force. This is a server-side change; raise
it with the validatord team rather than working around it in the SDKs.
