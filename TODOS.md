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

### 2026-09-10 update — the security scan filed this three times, and the blocker is now precise

Scan findings **4284688** (Java, HIGH), **4284646** (TS, MED) and **4284645** (Go, MED) are all
this entry. Left as documentation by decision, because the reply cannot be bound to an entity
without a server change — but the *exact* blocker and the actionable client-side path are worth
recording so the next reader does not re-derive them:

- **`GetMultiFactorSignatureEntitiesInfoReply` carries only `{id, payloadToSign[], entityType}`.**
  All three fields are required; there is **no entity id, singular or plural**, and
  `payloadToSign` is a bare string array with no per-element id, hash or ordering marker. So the
  reply cannot be joined back to the entities the caller asked about, and an integrator following
  the documented flow signs opaque bytes chosen entirely by the server.
- **`entityType` is a bare kind enum** (`REQUEST` / `WHITELISTED_ADDRESS` /
  `WHITELISTED_CONTRACT`). Java's domain object `MultiFactorSignatureEntityType` has `id` and
  `kind` fields, both of which MapStruct's `fromEntityTypeDTO(enum)` left null — the mapper test
  only asserted non-null, so nothing caught it. **Fixed 2026-09-10**: the mapper method is
  hand-written and populates `kind`; `id` stays null because the reply carries none, and that
  absence is the blocker itself. The test now asserts the kind for every enum value and asserts
  `id` is null *with the reason*, so nobody "fixes" it by inventing one. The kind matters
  because it is what tells a caller which verifying reader to check the payload against.
- **The actionable path, for when the decision comes:** `create` DOES take `entityIDs`
  (`TgvalidatordCreateMultiFactorSignaturesRequest{entityType, entityIDs}`, both required). So
  carrying the create-time `(entityType, entityIDs)` through to approve — or taking them as
  parameters — would let the SDK re-read those entities through the already-verifying reader for
  that type and require each `payloadToSign` element to equal a locally recomputed, verified
  metadata hash. That is a client-side fix and needs no wire change; it just needs someone to
  decide it is the intended model.
- **`REQUESTMOBILEAPPSIGNER` / `WHITELISTEDADDRESSMOBILEAPPSIGNER` are governance `Role` values**
  (`request_reply.proto:3004-3005`), in the same enum as `REQUESTAPPROVER`. So the second-factor
  signature is a governance-level approval keyed to a user public key — precisely the artefact a
  compromised server cannot forge, which is what makes the substitution worth doing.
- Python's half of this entry is now **out of date in the SDK's favour**: its
  `create_challenge`/`verify_challenge` methods were not merely a different shape, they called
  four generated operations that do not exist and could never run. Rewritten 2026-09-10 onto the
  four real operations, so there IS now a cross-SDK shape to compare.

Until the decision lands, all three SDKs state in the method docs that `payloadToSign` is
unverified server data and that the caller must bind it to a verified entity before signing.

---

## TPV1's canonical string is not injective — needs a versioned scheme (TPV2)

**Status:** open, documentation only. Scan finding **4284627** (MEDIUM).

**What:** The HMAC input is built by dropping empty components and joining the rest with a single
space: `Stream.of("TPV1", apiKey, nonce, ts, method, host, path, query, contentType, body)
.filter(nonEmpty).collect(joining(" "))` (Java `CryptoTPV1.java:117-123`, and byte-for-byte the
same construction in Go `crypto/tpv1.go:95-110`, Python `crypto/tpv1.py:112-129`, TS
`crypto/tpv1.ts:226-246`, plus the Postman pre-request script).

**Why it is a defect:** the mapping (method, host, path, query, content-type, body) -> message is
not one-to-one, for two independent reasons. The delimiter legitimately occurs INSIDE components —
the normalised Content-Type is `application/json; charset=utf-8`, the body is arbitrary text, and
the path is signed percent-DECODED, so a wire `%20` becomes a delimiter. And empty components are
elided rather than kept as fixed slots, so a component can be emptied and its text moved into the
adjacent one with no change to the MAC.

Consequences for an on-path attacker who can read and rewrite one signed request (a TLS-terminating
proxy or CDN, corporate TLS inspection, a compromised session) — the exact position TPV1's
documented integrity property exists to defend against (`docs/AUTHENTICATION.md:51`):

| rewrite | effect |
|---|---|
| move a GET's query into the `Content-Type` header | the server runs the UNFILTERED, unpaginated operation under the integrator's key, and the response flows back through the attacker |
| move a POST's body into `Content-Type`, send an empty body | the path-identified operation (reject, cancel, update) executes with its payload blanked |
| shift the path/query split where a path parameter carries `%20` | re-targets the call |

The attacker cannot INJECT content, only redistribute the signed bytes across adjacent fields, so
the gain is bounded — but the signature no longer uniquely binds the request it protects.

**Why this is documentation and not a patch:** the server verifies the same string, so any fix is a
protocol change that has to land in lock-step across four SDKs, the Postman collection and
tg-validatord. The fix itself is well understood: a fixed number of fields with empty fields
represented explicitly, each field length-prefixed or replaced by its SHA-256 digest (at minimum
body and Content-Type), and the path signed in the exact raw form sent on the wire. Ship it as
**TPV2** rather than mutating TPV1, or every deployed client breaks at once.

**Assumption stated explicitly:** that tg-validatord reconstructs the documented string from the
received request, taking the Content-Type header verbatim and dropping empty components. That is
what interoperability with these SDKs requires — the vendor's own Postman script signs
`application/json` while the SDKs sign `application/json; charset=utf-8`, so the server must accept
whatever it is sent. A server that re-normalised Content-Type, or rejected a Content-Type header on
a bodiless request, would defeat the header-based rewrites.

**What WAS fixed on 2026-09-10:** the separate, non-protocol half — Python and TS upper-cased the
HTTP method while Java and Go signed it verbatim, so a caller issuing lowercase `get` produced a
different canonical string in the two families and one of them could not authenticate. All four now
normalise, and a new canonical-string vector pins it. That is X2 in the security-fix plan; it is
NOT this entry.

**Depends on:** a validatord decision on adopting a versioned canonicalisation.

## The single-Address seam is not literally single, and the three SDKs disagree

**What:** `AddressService.verifiedAddress` is documented as the ONE construction seam for an
`Address`, and in Java it is: `getAddress`, `getAddresses`, `createAddress` and
`AssetService.getAssetAddresses` all route through `verifiedAddress` / `verifiedAddresses`. In
**Go and TypeScript the two PAGE paths do not** — `ListAddresses` and
`AssetService.GetAssetAddresses` call the batch `helper.VerifyAddressSignatures` instead.

**Why it is not a hole, and why it still matters:** the batch verifier is *stricter* than the
seam — it errors on an empty address string, where the seam returns it — so the security
invariant ("never return a non-empty address that has not been verified") holds on all four
paths in all three SDKs. What differs is the **asynchronous-creation case**: address creation is
async (`status` is one of `created`/`creating`/`signed`/`observed`/`confirmed`), so a page can
legitimately contain an address that has no address string yet. In Java that row comes back with
its `status`; in Go and TypeScript **the whole page fails**.

So a caller listing a wallet's addresses immediately after creating one gets a working list in
Java and an `IntegrityError` in Go and TypeScript. That is a usability divergence on a read path,
not a verification difference.

**Decide which behaviour is wanted before aligning.** The seam's behaviour is the better one — an
in-flight address should not deny access to the whole page, and the batch verifier's strictness
buys nothing, since an empty address carries no destination to misuse. But changing it alters
list semantics in two SDKs on a path no scan finding named, so it was left recorded rather than
changed. **Do not "align" it by loosening the batch verifier alone** — route the page paths
through the seam, so there is one rule rather than two that happen to agree.

**Context:** surfaced 2026-09-10 while porting the T3 `createAddress` fix. Go's seam doc comment
claimed all four paths went through it; that claim is corrected in place and now points here.

**Depends on:** a decision on the async-list behaviour. No blocker otherwise.

---

## No SDK tests its own strict base64 decoder

**What:** the repo-root `CLAUDE.md` records "Base64 on any governance path is decoded STRICTLY"
as a security property with a named implementation per SDK (Go `base64.StdEncoding`, Python
`_strict_base64.py`, TypeScript `helpers/strict-base64.ts`, Java `helper/StrictBase64.java`).
**Nothing pins it.** A grep of all four test trees for `strict_b64decode` / `StrictBase64` /
`strictBase64Decode` returns nothing outside the TypeScript suite added 2026-09-10.

**Why:** this is the primitive that decides what bytes a signature covers. A lenient decoder
silently discards out-of-alphabet characters and absorbs the rest, so a container carrying
embedded separators decodes to the genuine bytes **with attacker bytes appended** — and protobuf
treats concatenation as merge, so appended bytes ADD entries to repeated fields like `users`.
That is how an attacker-chosen HSM or `PRICEUPDATER` key reaches the trust root. Go was never
exploitable by this route *only* because it happened to use the strict encoder. A property with
no test, in four independent implementations, is the exact shape this repo keeps rediscovering by
grep.

**What it should be:** a shared accept/reject vector file — the accepted set (padding variants,
line wrapping) and the rejected set (embedded separators, misplaced padding, non-multiple-of-four,
non-ASCII) — loaded by all four suites, same pattern as `authorization-error-vectors.json`. It is
pure input→outcome, so no key material is involved. The TypeScript test added in the 2026-09-10
pass (`tests/unit/helpers/strict-base64.test.ts`) has the table to lift; it deliberately asserts
the accept/reject SET rather than the implementation, which is what let a linear-scan rewrite be
proven behaviour-preserving.

**One live defect found while writing it, fixed in TypeScript only:** the well-formedness regex
`/^(?:[A-Za-z0-9+/]{4})*.../` recurses per four-character group in V8, so a multi-megabyte input
threw `RangeError: Maximum call stack size exceeded` out of the decoder instead of returning a
decision. Measured: 1.6 MB passed in 27 ms, 5.6 MB threw. It failed **closed**, so it was
availability rather than a verification bypass — but a `RangeError` escapes every
`catch (e) { if (e instanceof ...) }` funnel in the SDK, and it was reachable from any governance
response through the signature-verification decode and the memo key, only one of which sits
behind a size cap. Replaced with a linear scan over the same accepted set. Go, Python and Java
validate without a quantified-group regex and were never affected — TypeScript was the lone
outlier, again.

**Depends on:** nothing.

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

## 4. Revive or retire the Python lint gate (Java PMD: done)

**What:** Decide, and record the decision, whether the Python lint gate is a gate.
Java PMD was the other half and is resolved — see below.

**Why:** Both are red on committed `main` and have been ignored for long enough that every
alignment pass rediscovers them, spends time establishing they are pre-existing, and skips
them. Either state is fine; the ambiguity is what costs time.

**Context:** Measured baselines (see the per-SDK `CLAUDE.md` files):

- **Python** `./build.sh lint`: 96 files fail flake8 (mostly E501 at the 100-char limit),
  79 would be reformatted by black, mypy reports 3008 errors in 124 files. Reformatting the
  package to chase it is a repo-wide project that would bury any review in churn.
- **Java** — **RESOLVED 2026-09-10: PMD is a real gate now, at 0.** It was 10 on master, and
  every one was mechanical: unused imports and unnecessary fully-qualified names left behind
  when logic moved, plus two `PreserveStackTrace` in `RequestService` (chain the cause). Fixed
  in passing during the security-scan pass, so `mvn pmd:check -pl client` exits 0 and its exit
  code now means something. Checkstyle was already a real gate at 0, and SpotBugs is at 0 (note
  `spotbugs:check` must run **without** `-o`; `findsecbugs-plugin` is not in the local `~/.m2`
  cache). **Only the Python half of this item is still open.**

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

### 2026-09-10 update — the verifying paths are done in Python; the mechanical remainder is not

The 2026-09-10 security scan filed this twice against Python (`address_service.py`,
`request_service.py`) and once against TypeScript (`base.ts`). The decision taken was to fix
**every funnel on a path that can raise `IntegrityError` today** and record the rest here.

Done in Python (all 7 funnels widened per file unless noted, and the module-level import
hoisted — see the landmine note below): `address_service.py` (7), `request_service.py` (11),
`whitelisted_address_service.py` (4), `whitelisted_asset_service.py` (4),
`governance_rule_service.py` (8), `asset_service.py` (4), `price_service.py` (2), and
`taurus_network/pledge_service.py` (14, after the read paths gained hash verification).
The whitelist and pledge funnels also carry `WhitelistError`.

Two shapes were fixed while in there, both worth knowing before touching the rest:

- **A function-local `from taurus_protect.errors import X` is a landmine, not a style choice.**
  Python makes the name local to the WHOLE function, so an `except X` or `raise X` *earlier* in
  the same function raises `UnboundLocalError`. This was a **live defect**:
  `request_service.py`'s `approve_requests` caught `IntegrityError` before an except funnel
  that re-imported it, so a genuine verification failure crashed with `UnboundLocalError`
  instead of reporting. Hoisted to module scope in the three files touched (35 local imports
  removed from `pledge_service.py` alone); the pattern is still present elsewhere.
- **Wrapping an `APIError` into an `IntegrityError` is the same inversion, backwards.** The
  whitelist approve paths' `except Exception -> IntegrityError("the verified read failed")`
  turned a transport failure into an integrity failure and lost `is_retryable()`. Both now
  propagate `(APIError, IntegrityError, WhitelistError)` unchanged and wrap only the rest.
  Go keeps the type reachable via `%w`; `raise ... from` does not, so the tuple is required.

**Remaining: roughly 100 funnels** across the other ~30 `services/*.py` files, none of which
can raise `IntegrityError` today. Two shapes are mixed in and should be unified with them:
funnels keyed on the **raw generated `ApiException`** rather than the SDK `APIError`, and
narrow two-element tuples. Do it as one mechanical pass; it is not urgent, but it is what stops
a new verification site inheriting the bug.


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

### 2026-09-10 update — both halves are CLOSED, and the classification rule is now written down

TypeScript routes both whitelist services through one `rethrowIfNotRowLevel` seam
(`src/services/row-level-error.ts`); Java's catch is
`catch (ContainerIntegrityException e) { throw e; }` then
`catch (WhitelistException | IntegrityException e)`, in that order (the first is a subclass of
the second). Gated by `tests/unit/services/whitelisted-address-container-abort.test.ts` and
`service/WhitelistedAddressExclusionTest.java`.

The decision the item was waiting on: **`IntegrityException` stays unchecked and the catch
sites were widened** — but the rule that came out of it is the useful part, and it is a
classification rule, not a hierarchy one:

> Throw the CHECKED `WhitelistException` for anything that is one ROW's problem; reserve the
> unchecked `IntegrityException` (and `ContainerIntegrityException`) for what invalidates the
> whole call.

Applied when the 2026-09-10 duplicate-JSON-key rejection landed:
`WhitelistHashHelper.rejectDuplicateObjectKeys` initially threw `IntegrityException`, which
silently escaped every existing `catch (WhitelistException)` — so one unparseable row would
have aborted a whole listing, re-opening this exact item through a new door. It throws
`WhitelistException` now. Two pre-existing tests caught it
(`WhitelistHashHelperTest.testParseWhitelistedAddressFromJson_InvalidJson` and
`VerificationBehaviourVectorsTest`), which is the gate working: both public parse entry points
declare `throws WhitelistException`, so changing what they actually throw was a visible API
break rather than a silent one. Keep that alignment when adding a check to either parser.

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

### 2026-09-10 update — the SIGNING consequence is now closed; the READ consequence is not

The security scan filed the signing half of this four times: findings **4284629** (Go),
**4284644** (Java), **4284643** (Python), **4284632** (TS), all HIGH. Their common shape was that
the whitelist approval took bare row ids, re-read them, and signed whatever came back under those
ids — so a stale-but-validly-signed container was not just a misleading READ, it was the thing
that let a substituted row survive verification and get an approver's signature.

**What changed:** all four approvals now take a content pin minted by the preceding verified read
(a witness type whose map of `id -> reviewed metadata hash` cannot be built by hand), and refuse to
sign when a re-read row's hash differs. That closes the *harvest*: a substituted row can no longer
be signed regardless of which container cleared it. Comparison is constant-time, matching
`approveRulesProposal`'s `expectedContainerHash`.

**What is still open, and is what this entry remains about:** a uniformly stale container still
misleads every READ. Two further notes for whoever picks this up:

- **The pin does not make the read-path fix unnecessary**, it just removes the worst consequence.
  A caller who reads a whitelist and acts on it without approving anything is still judged against
  whatever ruleset the server chose to serve.
- **Fix option 2 from the findings ("check the in-band container against an independently obtained
  current ruleset") is a bigger change than it sounds:** none of the four whitelist services holds
  a `RulesContainerCache` today — only `AddressService`, `AssetService` and `PriceService` do — so
  it means threading the cache or a `GovernanceRuleService` into a whitelist-service constructor in
  all four SDKs. Worth knowing before scoping it. And it is still racy without the ruleset identity
  this entry asks for, because a ruleset can be promoted between the two calls.

**A pin mismatch is a real signal, not a false positive.** `metadata.hash` is recomputed by the
server on every read (`enrichWLAs` -> `ToWLAMetadata`) from the immutable envelope PLUS the row's
live linked-address and linked-wallet rows, so it moves when a linked address is renamed. Usually
that also breaks signature coverage and the row is excluded anyway; for a legacy-signed row the
strip removes the inner labels, so the hash can move while the row still verifies. Either way the
content changed since review — re-read, re-review, re-approve.
