# Verification & Integrity Helpers — Java SDK

## Critical Files
- `SignatureVerifier.java` — SuperAdmin signature verification (isValidSignature, verifyGovernanceRules)
- `AddressSignatureVerifier.java` — HSM signature verification for addresses
- `WhitelistHashHelper.java` — Hash computation for whitelisted address/asset payloads
- `WhitelistIntegrityHelper.java` — Integrity verification orchestration
- `AssetHashHelper.java` — Asset-specific hash computation
- `ValidationHelper.java` — Input validation utilities

## Whitelisted Address 6-Step Verification (`WhitelistedAddressService.java`)
1. Verify metadata hash (SHA-256 of `payloadAsString` == `metadata.hash`)
2. Verify rules container signatures (SuperAdmin keys)
3. Decode rules container (base64 -> protobuf -> model)
4. Verify metadata hash is in signed hashes list
5. Verify whitelist signatures per governance thresholds
6. Parse WhitelistedAddress from verified payload

## Legacy Hash Computation (`WhitelistedAddressService.java:367-399`)
Address legacy strategies: remove `contractType`, remove `label` from `linkedInternalAddresses`, remove both.
Asset legacy strategies (`AssetHashHelper`): remove `isNFT`, remove `kindType`, remove both.
**CRITICAL:** Address and Asset use DIFFERENT legacy hash functions.

## Field Sourcing (CRITICAL SECURITY)
- Security-critical fields MUST come from verified payload, never from DTO
- `WhitelistedAddressService.list()` verifies each envelope (strict mode — fail on first error)
- If payload is missing a field, result is null (not DTO value)
- Non-security fields (status, action, rule, createdAt) can come from DTO

## Constant-Time Comparison
- Uses BouncyCastle `constantTimeAreEqual()`
- Perform dummy comparison on length mismatch
- Never break/return early in multi-signature loops

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
