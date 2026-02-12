# Verification & Integrity Helpers — TypeScript SDK

## Critical Files

- `signature-verifier.ts` — SuperAdmin signature verification (isValidSignature, verifyGovernanceRules)
- `address-signature-verifier.ts` — HSM signature verification for addresses
- `whitelist-hash-helper.ts` — Hash computation for whitelisted address/asset payloads
- `whitelisted-address-verifier.ts` — WhitelistedAddressVerifier (6-step verification)
- `whitelisted-asset-verifier.ts` — WhitelistedAssetVerifier (5-step verification)
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

Same as address steps 1-5 using `ContractAddressWhitelistingRules`.

## Legacy Hash Computation

Address: remove `contractType`, remove `label` from `linkedInternalAddresses`, remove both.
Asset: remove `isNFT`, remove `kindType`, remove both.
**CRITICAL:** Address and Asset use DIFFERENT legacy hash functions.

## Field Sourcing (CRITICAL SECURITY)

- Security-critical fields MUST come from verified payload only
- Non-security fields can fall back to DTO
- `WhitelistedAddressService.list()` verifies each envelope (strict mode — fail on first error)

## Constant-Time Comparison

- Uses Node.js `crypto.timingSafeEqual()` internally
- Never break/return early in multi-signature loops

## AddressService Mandatory Verification

- `RulesContainerCache` REQUIRED at construction (not optional)
- No `setRulesCache()` method — cache provided at construction
- `ProtectClient.addresses` getter requires `superAdminKeysPem` configured
- Rules container must be signature-verified by SuperAdmin keys before trusting HSM public keys

## Governance Rules Mapper (`mappers/governance-rules.ts`)

- `rulesContainerFromBase64()` — decodes base64 rules container (JSON format)
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

```typescript
// WRONG
catch {} // bare catch

// CORRECT
catch (error: unknown) {
    if (error instanceof Error && isCryptoError(error)) continue;
    throw error; // re-throw unexpected
}
```

### minValidSignatures Validation

Validated > 0 whenever explicitly set (not just when superAdminKeysPem present).
