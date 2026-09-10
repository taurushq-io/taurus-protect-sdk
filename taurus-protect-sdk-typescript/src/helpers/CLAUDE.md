# Verification & Integrity Helpers — TypeScript SDK

## Critical Files

- `signature-verifier.ts` — SuperAdmin signature verification (isValidSignature,
  verifySignatureWithKey, verifyGovernanceRules, **verifyGovernanceRulesSignatures** — the shared
  threshold core)

  `verifySignatureWithKey(data, sigB64, key)` is the single-key form the other three SDKs already had.
  Without it, checking one known signer meant passing a one-element array to `isValidSignature`, which
  reads as a quorum check. It returns `false` for a missing key or malformed signature rather than
  throwing: data, signature and key are all public here, so there is nothing to leak.
- `address-signature-verifier.ts` — HSM signature verification for addresses
- `whitelist-hash-helper.ts` — Hash computation for whitelisted address/asset payloads
- `whitelisted-address-verifier.ts` — WhitelistedAddressVerifier (6-step verification)
- `whitelisted-asset-verifier.ts` — WhitelistedAssetVerifier (**all six**: `verify()` returns
  the parsed asset and the branded envelope, not just a verdict)
- `constant-time.ts` — Timing-safe comparison using crypto.timingSafeEqual()
- `metadata-utils.ts` — jsonValueToString() helper for amount field parsing

## Whitelisted Address 6-Step Verification

1. Verify metadata hash (SHA-256 of `payloadAsString` == `metadata.hash`)
2. Verify rules container signatures (SuperAdmin keys)
3. Decode rules container (base64 -> JSON -> model)
4. Verify hash coverage (`metadata.hash` in signature's hashes list)
5. Verify whitelist signatures per governance thresholds
6. Parse WhitelistedAddress from verified payload

## Whitelisted Asset 5-Step Verification
**Steps 1-5 are the signature checks; step 6 (parse from the VERIFIED payload) is what stops
an unsigned value reaching the caller.** See "Critical Files" above for where it runs here.


Same as address steps 1-5 using `ContractAddressWhitelistingRules`.

## Legacy Hash Computation

Address: remove `contractType`, remove `label` from `linkedInternalAddresses`, remove both.
Asset: remove `isNFT`, remove `kindType`, remove both.
**CRITICAL:** Address and Asset use DIFFERENT legacy hash functions.

## Field Sourcing (CRITICAL SECURITY)

- Security-critical fields MUST come from verified payload only
- Non-security fields can fall back to DTO
- `WhitelistedAddressService` list/`List()` verifies every envelope **leniently**: an
  unverifiable row is excluded and reported, not fatal. One bad row used to deny access to
  every good one — and listing is how an operator finds the bad row. Excluding stays
  fail-closed: an omitted destination cannot be selected. **Rows returned but none surviving
  is an error**, so a filtered page never reads as an empty whitelist.

## Constant-Time Comparison

- Uses Node.js `crypto.timingSafeEqual()` internally
- Never break/return early when comparing SECRET or hash material — that is what
  `verifyHashCoverage` / `containsHash` are for. It does NOT apply to a threshold
  loop exiting once the required count is reached, nor to trying a signature against
  a set of PUBLIC keys: all four SDKs return early on group-threshold success, and
  that is correct. Do not "fix" it.

## AddressService Mandatory Verification

- `RulesContainerCache` REQUIRED at construction (not optional)
- No `setRulesCache()` method — cache provided at construction
- `ProtectClient.addresses` getter requires `superAdminKeysPem` configured
- Rules container must be signature-verified by SuperAdmin keys before trusting HSM public keys

## Governance Rules Mapper (`mappers/governance-rules.ts`)

- `rulesContainerFromBase64()` — decodes a base64 rules container: protobuf first, JSON as a fallback
- `userSignaturesFromBase64()` — decodes base64 user signatures
- Handles both camelCase and snake_case property names

## Security-Specific Lessons

### DER Signature Parsing Bounds Checks

- Minimum DER signature length check (< 8 bytes)
- Bounds check before each buffer access: `if (offset + length > derSignature.length)`
- Extended length encoding (> 127 bytes): reject 0x80 alone, validate declared vs remaining

**DER Extended Length Validation:**
DER signatures with extended length encoding (length > 127 bytes) require full validation:
1. Verify 0x80 alone (indefinite length) is rejected
2. Read the actual length from the extended bytes
3. Validate declared length matches remaining buffer size

### Exception Handling in Verification Loops

``typescript
// WRONG
catch {} // bare catch

// CORRECT
catch (error: unknown) {
    if (error instanceof Error && isCryptoError(error)) continue;
    throw error; // re-throw unexpected
}
``

### minValidSignatures Validation

Validated > 0 whenever explicitly set (not just when superAdminKeysPem present).

## The threshold is evaluated in ONE function — route everything through it

`verifyGovernanceRulesSignatures(rulesContainerData, signatures, superAdminKeys, minValidSignatures)` is the single place `minValidSignatures` is
checked. It counts **distinct signing keys** (SHA-256 of the encoded public key), skips absent
signatures, and rejects a non-positive threshold. Every caller delegates: the governance service and
both whitelisted-address/asset paths.

Do **not** hand-roll a `validCount++` loop over signatures, and do not dedupe on `userId` — it is
server-supplied, and ECDSA is randomized, so one key can emit as many valid signatures as any
threshold demands. That bug shipped once: the hardened helper existed but had zero callers while the
live paths counted entries.

The **per-group** step-5 threshold (`GroupThreshold.minimumSignatures`) is a different
threshold — different scope, **same counting rule**. It lives in the two verifiers'
`verifyGroupThreshold`, keyed on `keyFingerprint` of the container-resolved public key and
never on the entry's `userId`. Counting entries let a duplicated entry from one group member
satisfy an N-of-M group, promoting an under-approved whitelist entry to approved; two user IDs
sharing one key likewise count once, because that is one compromised secret. Pinned by
`tests/unit/helpers/group-threshold.test.ts` and, cross-SDK, by the `group_threshold` section
of `scripts/resources/verification-signed-fixtures.json`
(`tests/unit/helpers/signed-fixtures.test.ts`).

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
