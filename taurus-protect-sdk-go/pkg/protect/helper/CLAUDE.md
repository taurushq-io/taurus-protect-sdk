# Verification & Integrity Helpers -- Go SDK

## Critical Files

- `signature_verifier.go` -- SuperAdmin signature verification (IsValidSignature, VerifyGovernanceRules,
  VerifySignatureWithKey, VerifyHashCoverage, **VerifyGovernanceRulesSignatures** — the shared
  threshold core)

  `VerifyGovernanceRules(rules, minValidSignatures, superAdminKeys)` takes a `*model.GovernanceRuleset`
  and decodes its base64 container itself, so a caller outside the service layer does not have to.
  This file was listed here for a long time **before it existed** — the doc was ahead of the code, and
  Java/Python/TS all had the convenience while Go had only the `*Signatures` core. If you add a name to
  this list, add the function in the same change.
- `address_signature_verifier.go` -- HSM signature verification for addresses
- `whitelist_hash.go` -- Hash computation and legacy hash strategies
- `whitelisted_address_verifier.go` -- WhitelistedAddressVerifier (6-step verification)
- `whitelisted_asset_verifier.go` -- WhitelistedAssetVerifier (steps 1-5; **step 6 runs in
  `service/whitelisted_asset.go`**, which calls `ParseWhitelistedAssetFromJSON`)

## Whitelisted Address 6-Step Verification

1. Verify metadata hash (SHA-256 of `PayloadAsString` == `Metadata.Hash`)
2. Verify rules container signatures (SuperAdmin keys)
3. Decode rules container (base64 -> protobuf -> model)
4. Verify hash coverage (`Metadata.Hash` in signature's hashes list)
5. Verify whitelist signatures meet governance thresholds
6. Parse WhitelistedAddress from verified payload

## Whitelisted Asset 5-Step Verification
**Steps 1-5 are the signature checks; step 6 (parse from the VERIFIED payload) is what stops
an unsigned value reaching the caller.** See "Critical Files" above for where it runs here.


Same as address steps 1-5 using `ContractAddressWhitelistingRules`.
- `parallel_thresholds` is `[]*SequentialThresholds`, not `[]*GroupThreshold`
- Simpler structure than AddressWhitelistingRules (no rule lines)

## Legacy Hash Computation

Address: remove `contractType`, remove `label` from `linkedInternalAddresses`, remove both.
Asset: remove `isNFT`, remove `kindType`, remove both.
**CRITICAL:** Address and Asset use DIFFERENT legacy hash functions (`ComputeLegacyHashes` vs `ComputeAssetLegacyHashes`).

## Field Sourcing (CRITICAL SECURITY)

- Security-critical fields MUST come from verified payload only
- Non-security fields can fall back to DTO
- `WhitelistedAddressService` list/`List()` verifies every envelope **leniently**: an
  unverifiable row is excluded and reported, not fatal. One bad row used to deny access to
  every good one — and listing is how an operator finds the bad row. Excluding stays
  fail-closed: an omitted destination cannot be selected. **Rows returned but none surviving
  is an error**, so a filtered page never reads as an empty whitelist.
- `Payload` is `map[string]interface{}` -- type assertions required

## Constant-Time Comparison

- `helper.ConstantTimeCompare()` for timing-safe comparison
- Dummy comparison on length mismatch
- Never break/return early when comparing SECRET or hash material — that is what
  `verifyHashCoverage` / `containsHash` are for. It does NOT apply to a threshold
  loop exiting once the required count is reached, nor to trying a signature against
  a set of PUBLIC keys: all four SDKs return early on group-threshold success, and
  that is correct. Do not "fix" it.

## Request Hash Verification (`service/request.go`)

- Computes: `crypto.CalculateHexHash(payload)`
- Error MUST include computed/provided values: `fmt.Sprintf("...computed=%s, provided=%s", ...)`
- When hash exists but PayloadAsString is empty -> MUST return error

## Address Verification Fail-Fast

``go
if rulesContainer == nil {
    return nil, fmt.Errorf("rules container required for address signature verification")
}
``
Never silently skip verification.

## Security-Specific Lessons

### RulesContainerCache has no setter

The constructor's fetcher is the ONLY way a container enters the cache. `Set(container)`
and `SetFetcher(fn)` are **deleted** — either seated a container nothing had verified,
and this cache supplies the HSM public key address verification trusts, so an unverified
one makes that check pass against whatever key the container carries. Neither had a
caller outside tests, and no sibling SDK offers an equivalent (Java and Python inject the
governance service, TypeScript a provider closure). `protect.NewClient` wires the fetcher
to `GovernanceRuleService.GetDecodedRulesContainer`, which verifies. Guarded by
`pkg/protect/rules_cache_verification_test.go`; do not reintroduce the setters.

### RulesContainerCache Thread Safety

Uses "fetching" flag + channel pattern: network I/O must NOT occur while holding the lock. Other goroutines wait on completion channel. See `cache/rules_container.go`.

### Cache Fetch Error Propagation

Store `fetchErr` and propagate to waiters via the completion channel. Waiters check for error after channel signals.

### Legacy Hash Support

Both `WhitelistedAddressVerifier` and `WhitelistedAssetVerifier` must call `ComputeLegacyHashes()` / `ComputeAssetLegacyHashes()` to support items signed before schema changes. The address verifier had this; the asset verifier was initially missing it.

### Thread-Safe Cache with Network I/O

**Problem:** Calling network fetcher while holding a mutex lock causes deadlocks under concurrent load.

**Solution:** Use a "fetching" flag with channel coordination:
1. Check if another goroutine is already fetching
2. If so, wait on the completion channel
3. If not, mark as fetching, release lock, perform I/O
4. Reacquire lock to update cache and signal completion

See `cache/rules_container.go` for the implementation pattern.

## The threshold is evaluated in ONE function — route everything through it

`VerifyGovernanceRulesSignatures(rulesContainerData, signatures, superAdminKeys, minValidSignatures)` is the single place `minValidSignatures` is
checked. It counts **distinct signing keys** (SHA-256 of the encoded public key), skips absent
signatures, and rejects a non-positive threshold. Every caller delegates: the governance service and
both whitelisted-address/asset paths.

Do **not** hand-roll a `validCount++` loop over signatures, and do not dedupe on `userId` — it is
server-supplied, and ECDSA is randomized, so one key can emit as many valid signatures as any
threshold demands. That bug shipped once: the hardened helper existed but had zero callers while the
live paths counted entries.

The **per-group** step-5 threshold (`GroupThreshold.minimumSignatures`) is a different
threshold — but it counts distinct signers for the same reason. It lives in
`group_threshold.go`, keyed on `RuleUser.KeyFingerprint()`, i.e. the container-resolved
public key and never the server-supplied `userId`. Counting entries let a duplicated entry
from one group member satisfy an N-of-M group, promoting an under-approved whitelist entry
to approved; two user IDs sharing one key likewise count once, because that is one
compromised secret. `group_threshold_test.go` pins the behaviour and
`signed_fixtures_test.go`'s `group_threshold` section pins it cross-SDK from the same file
that gates the SuperAdmin threshold — three vectors go red against entry counting.

## Local rules for this package

The reasoning behind these lives in the repo-root `CLAUDE.md` → "Verification flows — facts
that cost time to rediscover". What matters when editing *here*:

- **Exactly TWO hash-comparison functions, and both live in `signature_verifier.go`.**
  `VerifyHashCoverage` (across signatures) and `containsHash` (within one). Both constant-time,
  neither breaks early, and `VerifyHashCoverage` **delegates** to `containsHash` rather than
  inlining its own inner loop. **Never add a private copy to a verifier** — that is how 14
  implementations accumulated across the repo and drifted.

  Two things this SDK got wrong until 2026-09-04, both worth not repeating: `containsHash`
  lived in `whitelisted_address_verifier.go` (the siblings keep the pair together — Java in
  `SignatureVerifier`, Python and TS in the hash helper), and it **returned on the first
  match** while Java, Python and TS all scanned the full list with an explicit "No break"
  comment. Python's docstring records the same drift already fixed there once.

  **The behaviour vectors cannot catch an early-break regression.**
  `scripts/resources/verification-behaviour-vectors.json` pins input→outcome, and early
  return and full scan agree on every output — they differ only in timing. Neither can any
  unit test, in any of the four SDKs; none of them guards it. Code review is the only gate,
  which is exactly why the comment on each function says so explicitly.
- **`resolveRuleKey` keys rule selection on the SIGNED payload.** Chain required (absent is an
  error, never a wildcard); network falls back to the DTO only when the payload has none,
  because `includeNetworkInPayload` may be off. Do not tighten it to require both.
- **`minimumSignatures = 0` on a populated group is rejected.** There is no post-loop threshold
  check, so a zero would mean "one signature suffices". An empty group with zero is fine.
