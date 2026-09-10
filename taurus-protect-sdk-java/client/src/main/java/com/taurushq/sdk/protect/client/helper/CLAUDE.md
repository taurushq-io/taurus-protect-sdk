# Verification & Integrity Helpers — Java SDK

## Critical Files
- `SignatureVerifier.java` — SuperAdmin signature verification (isValidSignature, verifyGovernanceRules,
  verifySignature, verifyHashCoverage, **verifyGovernanceRulesSignatures** — the shared threshold core)

  `verifyHashCoverage(hash, signatures)` was absent here while Go, Python and TS all had it, so Java
  callers compared hashes by hand. It uses `MessageDigest.isEqual` and deliberately does **not** break
  on a match: returning early would leak which signature covered the hash through timing.
- `AddressSignatureVerifier.java` — HSM signature verification for addresses
- `WhitelistHashHelper.java` — Hash computation for whitelisted address/asset payloads
- `WhitelistIntegrityHelper.java` — Integrity verification orchestration
- `AssetHashHelper.java` — Asset-specific hash computation
- `ValidationHelper.java` — Input validation utilities

**There is no verifier class in this package.** Unlike the other three SDKs, whitelist
verification is inlined in `service/WhitelistedAddressService.java` and
`service/WhitelistedAssetService.java` — both flows run all six steps there. That layout is
where the non-constant-time hash comparisons hid; extracting it is tracked in `TODOS.md`.

## Whitelisted Address 6-Step Verification (`WhitelistedAddressService.java`)
1. Verify metadata hash (SHA-256 of `payloadAsString` == `metadata.hash`)
2. Verify rules container signatures (SuperAdmin keys)
3. Decode rules container (base64 -> protobuf -> model)
4. Verify metadata hash is in signed hashes list
5. Verify whitelist signatures per governance thresholds
6. Parse WhitelistedAddress from verified payload

## Legacy Hash Computation (`WhitelistedAddressService.java`, `computeLegacyHashes`)
Address legacy strategies: remove `contractType`, remove `label` from `linkedInternalAddresses`, remove both.
Asset legacy strategies (`AssetHashHelper`): remove `isNFT`, remove `kindType`, remove both.
**CRITICAL:** Address and Asset use DIFFERENT legacy hash functions.

## Field Sourcing (CRITICAL SECURITY)
- Security-critical fields MUST come from verified payload, never from DTO
- `WhitelistedAddressService` list/`List()` verifies every envelope **leniently**: an
  unverifiable row is excluded and reported, not fatal. One bad row used to deny access to
  every good one — and listing is how an operator finds the bad row. Excluding stays
  fail-closed: an omitted destination cannot be selected. **Rows returned but none surviving
  is an error**, so a filtered page never reads as an empty whitelist.
- If payload is missing a field, result is null (not DTO value)
- Non-security fields (status, action, rule, createdAt) can come from DTO

## Constant-Time Comparison
- Uses BouncyCastle `constantTimeAreEqual()`
- Perform dummy comparison on length mismatch
- Never break/return early when comparing SECRET or hash material — that is what
  `verifyHashCoverage` / `containsHash` are for. It does NOT apply to a threshold
  loop exiting once the required count is reached, nor to trying a signature against
  a set of PUBLIC keys: all four SDKs return early on group-threshold success, and
  that is correct. Do not "fix" it.

## Request Hash Verification (`RequestService.java`)
- Computes: `CryptoTPV1.calculateHexHash(payload)`
- Uses `constantTimeAreEqual()` for timing-safe comparison
- Throws `IntegrityException` on mismatch — logs computed and provided values via LOGGER.warning()
- When hash exists but payloadAsString is null -> MUST throw

## Request Approval Flow (`RequestService.java`)
1. Sort requests by ID
2. Build JSON: `gson.toJson(requests.map(r -> r.getMetadata().getHash()))`
3. Sign: `CryptoTPV1.calculateBase64Signature(privateKey, hashesBytes)`
4. Submit to API

## Governance Rules
- `DecodedRulesContainer.getHsmPublicKey()` — finds user with HSMSLOT role (thread-safe with synchronized)
- `findAddressWhitelistingRules(blockchain, network)` — three-tier priority (exact -> blockchain-only -> global)
- `RulesContainerCache` (`cache/RulesContainerCache.java`) — thread-safe with configurable TTL

## The threshold is evaluated in ONE function — route everything through it

`SignatureVerifier.verifyGovernanceRulesSignatures(rulesContainerData, signatures, superAdminPublicKeys, minValidSignatures)` is the single place `minValidSignatures` is
checked. It counts **distinct signing keys** (SHA-256 of the encoded public key), skips absent
signatures, and rejects a non-positive threshold. Every caller delegates: the governance service and
both whitelisted-address/asset paths.

Do **not** hand-roll a `validCount++` loop over signatures, and do not dedupe on `userId` — it is
server-supplied, and ECDSA is randomized, so one key can emit as many valid signatures as any
threshold demands. That bug shipped once: the hardened helper existed but had zero callers while the
live paths counted entries.

The **per-group** step-5 threshold (`GroupThreshold.minimumSignatures`) is a different
threshold — different scope, **same counting rule**. In this SDK it lives in the two whitelist
**services**, not here: `WhitelistedAddressService.verifyGroupThreshold` (package-private, for
test access) and its asset peer, keyed on `SignatureVerifier.keyFingerprint` of the
container-resolved public key and never on the entry's `userId`. Counting entries let a
duplicated entry from one group member satisfy an N-of-M group, promoting an under-approved
whitelist entry to approved; two user IDs sharing one key likewise count once, because that is
one compromised secret. Pinned by `service/GroupThresholdTest` and, cross-SDK, by the
`group_threshold` section of `scripts/resources/verification-signed-fixtures.json`
(`service/SignedFixturesGroupThresholdTest` — a *service*-package test for that reason).

An earlier revision of this file said the per-group threshold "legitimately counts entries…
leave it alone". That was wrong and is what let the bug survive review — do not restore it.

## Local rules for this package

The reasoning behind these lives in the repo-root `CLAUDE.md` → "Verification flows — facts
that cost time to rediscover". What matters when editing *here*:

- **Exactly TWO hash-comparison functions.** `verifyHashCoverage` (across signatures) and
  `containsHash` (within one). Both constant-time, neither breaks early. **Never add a private
  copy to a verifier** — that is how 14 implementations accumulated across the repo and drifted.
  Pinned cross-SDK by `scripts/resources/verification-behaviour-vectors.json`.
- **`resolveRuleKey` keys rule selection on the SIGNED payload.** Chain required (absent is an
  error, never a wildcard); network falls back to the DTO only when the payload has none,
  because `includeNetworkInPayload` may be off. Do not tighten it to require both.
- **`minimumSignatures = 0` on a populated group is rejected.** There is no post-loop threshold
  check, so a zero would mean "one signature suffices". An empty group with zero is fine.

## The MODEL package now calls in — mind the signatures

`SignedWhitelistedAssetEnvelope.markVerified` and its address peer call
`AssetHashHelper.parseWhitelistedAssetFromJson` / `WhitelistHashHelper.parseWhitelistedAddressFromJson`
directly, so `client.model` depends on `client.helper` (a package cycle, since helper
returns model types — legal in Java, and precedented by `RequestMetadata` calling
`CryptoTPV1`).

Why: those setters used to be public and took the asset/address, flipping the same
`isInitialized` gate the getter checks — so a caller could inject a fabricated value into
an envelope and read it back as "verified". `markVerified` **derives** the value from the
envelope's own signed payload instead, which needs these two parse functions.

Consequence: changing either signature now breaks the model package, not just the services.
`VerifiedEnvelopeDerivationTest` pins the behaviour.
