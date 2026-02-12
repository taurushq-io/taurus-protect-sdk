# Verification & Integrity Helpers -- Go SDK

## Critical Files

- `signature_verifier.go` -- SuperAdmin signature verification (IsValidSignature, VerifyGovernanceRules)
- `address_signature_verifier.go` -- HSM signature verification for addresses
- `whitelist_hash.go` -- Hash computation and legacy hash strategies
- `whitelisted_address_verifier.go` -- WhitelistedAddressVerifier (6-step verification)
- `whitelisted_asset_verifier.go` -- WhitelistedAssetVerifier (5-step verification)

## Whitelisted Address 6-Step Verification

1. Verify metadata hash (SHA-256 of `PayloadAsString` == `Metadata.Hash`)
2. Verify rules container signatures (SuperAdmin keys)
3. Decode rules container (base64 -> protobuf -> model)
4. Verify hash coverage (`Metadata.Hash` in signature's hashes list)
5. Verify whitelist signatures meet governance thresholds
6. Parse WhitelistedAddress from verified payload

## Whitelisted Asset 5-Step Verification

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
- `WhitelistedAddressService.List()` verifies each envelope (strict mode -- fail on first error)
- `Payload` is `map[string]interface{}` -- type assertions required

## Constant-Time Comparison

- `helper.ConstantTimeCompare()` for timing-safe comparison
- Dummy comparison on length mismatch
- Never break/return early in multi-signature loops

## Request Hash Verification (`service/request.go`)

- Computes: `crypto.CalculateHexHash(payload)`
- Error MUST include computed/provided values: `fmt.Sprintf("...computed=%s, provided=%s", ...)`
- When hash exists but PayloadAsString is empty -> MUST return error

## Address Verification Fail-Fast

```go
if rulesContainer == nil {
    return nil, fmt.Errorf("rules container required for address signature verification")
}
```
Never silently skip verification.

## Security-Specific Lessons

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
