# Whitelisted Address Integrity Verification

This document explains the cryptographic integrity verification process for whitelisted addresses in the Taurus-PROTECT TypeScript SDK.

## Overview

The SDK implements a **6-step cryptographic verification pipeline** to ensure whitelisted addresses are authentic and properly approved according to governance rules. Verification is available via `WhitelistedAddressService.withVerification()` - callers can choose between basic retrieval or full cryptographic verification.

## Verification Flow

```
+---------------------------------------------------------------------+
|                    API Response (Envelope)                          |
|  +------------------+  +------------------+  +-------------------+  |
|  | signedAddress    |  | metadata         |  | rulesContainer    |  |
|  |   - payload      |  |   - payloadAsStr |  |   - groups        |  |
|  |   - signatures[] |  |   - hash         |  |   - thresholds    |  |
|  +------------------+  +------------------+  | rulesSignatures   |  |
|                                             +-------------------+  |
+---------------------------------------------------------------------+
                                    |
                                    v
+---------------------------------------------------------------------+
|              STEP 1: Verify Metadata Hash                           |
|  SHA-256(payloadAsString) == metadata.hash ?                        |
|  Uses constant-time comparison to prevent timing attacks            |
+---------------------------------------------------------------------+
                                    |
                                    v
+---------------------------------------------------------------------+
|         STEP 2: Verify Rules Container Signatures                   |
|  SuperAdmin signatures on rulesContainer >= minValidSignatures      |
|  Uses ECDSA with raw r||s signature format (64 bytes)               |
+---------------------------------------------------------------------+
                                    |
                                    v
+---------------------------------------------------------------------+
|              STEP 3: Decode Rules Container                         |
|  Base64 -> JSON -> DecodedRulesContainer                            |
|  Contains groups, users, and approval thresholds                    |
+---------------------------------------------------------------------+
                                    |
                                    v
+---------------------------------------------------------------------+
|         STEP 4: Verify Hash Coverage                                |
|  metadata.hash must appear in at least one signature.hashes[]       |
|  Includes legacy hash support for backward compatibility            |
+---------------------------------------------------------------------+
                                    |
                                    v
+---------------------------------------------------------------------+
|         STEP 5: Verify Whitelist Signatures                         |
|  User signatures meet governance threshold requirements             |
|  Parallel paths (OR) -> Sequential groups (AND) -> Min signatures   |
+---------------------------------------------------------------------+
                                    |
                                    v
+---------------------------------------------------------------------+
|         STEP 6: Parse WhitelistedAddress from Payload               |
|  Parse verified payload JSON into WhitelistedAddress model          |
|  Return WhitelistedAddressVerificationResult                        |
+---------------------------------------------------------------------+
                                    |
                                    v
                    Verified WhitelistedAddress returned
```

## Step-by-Step Details

### Step 1: Metadata Hash Verification

**Purpose:** Ensure the payload data has not been tampered with.

**Location:** `src/helpers/whitelisted-address-verifier.ts:215-229`

**Process:**
1. Compute SHA-256 hash of `metadata.payloadAsString` using `calculateHexHash()`
2. Compare with provided `metadata.hash`
3. Uses **constant-time comparison** (`constantTimeCompare`) to prevent timing attacks

**Failure:** Throws `IntegrityError` if hashes do not match.

```typescript
const computedHash = calculateHexHash(envelope.metadata.payloadAsString);
if (!constantTimeCompare(computedHash, envelope.metadata.hash)) {
  throw new IntegrityError(`computed hash does not match provided hash`);
}
```

### Step 2: Rules Container Signature Verification

**Purpose:** Verify that SuperAdmins approved the governance rules.

**Location:** `src/helpers/whitelisted-address-verifier.ts:238-292`

**Process:**
1. Decode `rulesSignatures` from Base64 using `userSignaturesDecoder`
2. Decode `rulesContainer` from Base64 to raw bytes
3. For each signature, verify against SuperAdmin public keys using `isValidSignature()`
4. Count valid signatures
5. Require `validCount >= minValidSignatures`

**Cryptographic Algorithm:** ECDSA P-256 with raw r||s signature format (64 bytes)

**Failure:** Throws `IntegrityError` if insufficient valid signatures.

### Step 3: Rules Container Decoding

**Purpose:** Extract governance structure (groups, users, thresholds).

**Location:** `src/helpers/whitelisted-address-verifier.ts:302-313`

**Process:**
1. Decode Base64 `rulesContainer` using `rulesContainerDecoder`
2. Map to `DecodedRulesContainer` domain model

**Contains:**
- User groups with member IDs
- Parallel threshold paths (OR logic)
- Sequential group thresholds (AND logic)
- User public keys for signature verification

### Step 4: Hash Coverage Verification

**Purpose:** Confirm the metadata hash is covered by at least one signature.

**Location:** `src/helpers/whitelisted-address-verifier.ts:325-354`

**Process:**
1. Get `metadata.hash`
2. Check each `signature.hashes[]` list using `verifyHashCoverage()`
3. At least one signature must include this hash
4. If not found, try legacy hashes for backward compatibility

**Why needed:** Signatures cover a list of hashes. This ensures the specific address hash was actually included in what was signed (not just any signature exists).

**Failure:** Throws `IntegrityError` if hash not found in any signature's hashes list.

### Step 5: Whitelist Signature Verification

**Purpose:** Verify user approvals meet governance threshold requirements.

**Location:** `src/helpers/whitelisted-address-verifier.ts:364-405`

**Governance Structure:**
```
ParallelThresholds (OR paths)
  +-- SequentialThresholds (AND groups within a path)
        +-- GroupThreshold
              +-- groupId
              +-- minimumSignatures
```

**Algorithm:**
1. **Find applicable rules** based on blockchain/network using `findAddressWhitelistingRules()`
2. **Determine thresholds** based on linked wallets/addresses:
   - If no linked addresses AND exactly 1 linked wallet: check rule lines for wallet-specific thresholds
   - Otherwise: use default parallel thresholds
3. **Try each parallel path** (OR logic - only one needs to succeed):
   - For each sequential threshold path:
     - Verify ALL group thresholds (AND logic)
     - If all pass: verification succeeds
4. **Verify group threshold:**
   - Find group by ID in rules container
   - For each signature from a user in this group:
     - Check user is in group
     - Check signature covers metadata hash
     - Get user's public key from rules container
     - Verify signature: `ECDSA(JSON(hashes[]), userPublicKey)`
   - Count valid signatures
   - Require `validCount >= minimumSignatures`

**Signature Data Format:** JSON array of hashes, then ECDSA signed.

```typescript
const hashesJson = JSON.stringify(sig.hashes);
const hashesData = Buffer.from(hashesJson, 'utf-8');
if (verifySignature(publicKey, hashesData, sig.userSignature.signature)) {
  validCount++;
}
```

## Legacy Hash Support

For backward compatibility with addresses signed before schema changes, the SDK supports legacy hash computation.

**Location:** `src/helpers/whitelist-hash-helper.ts:46-88`

**Legacy Hash Strategies:**

1. **Remove contractType only** - Handles addresses signed before `contractType` was added
2. **Remove labels from linkedInternalAddresses** - Handles addresses signed after `contractType` but before labels
3. **Remove both** - Handles addresses signed before both fields were added

```typescript
// Strategy 1: Remove contractType
const withoutContractType = payloadAsString.replace(/,"contractType":"[^"]*"/g, '');

// Strategy 2: Remove labels from linkedInternalAddresses
const withoutLabels = payloadAsString.replace(/,"label":"[^"]*"}/g, '}');

// Strategy 3: Remove both
let withoutBoth = payloadAsString.replace(/,"label":"[^"]*"}/g, '}');
withoutBoth = withoutBoth.replace(/,"contractType":"[^"]*"/g, '');
```

## Cryptographic Primitives

| Operation | Algorithm | Implementation |
|-----------|-----------|----------------|
| Hash | SHA-256 | Node.js `crypto` |
| Signature | ECDSA P-256 (raw r||s format) | Node.js `crypto` |
| Comparison | Constant-time | `crypto.timingSafeEqual()` |

## Exception Types

| Exception | When Thrown |
|-----------|-------------|
| `IntegrityError` | Hash mismatch, insufficient signatures, cryptographic failure |
| `WhitelistError` | Governance thresholds not met, no matching rules |
| `ValidationError` | Invalid input parameters |
| `NotFoundError` | Address not found |

## Key Files

| File | Purpose |
|------|---------|
| `src/helpers/whitelisted-address-verifier.ts` | Main verification orchestration |
| `src/helpers/signature-verifier.ts` | SuperAdmin signature verification |
| `src/helpers/whitelist-hash-helper.ts` | Hash operations, legacy hashes, JSON parsing |
| `src/helpers/constant-time.ts` | Timing-safe comparison |
| `src/crypto/signing.ts` | ECDSA signing and verification |
| `src/crypto/hashing.ts` | SHA-256 hash computation |
| `src/services/whitelisted-address-service.ts` | Service layer with verification integration |

## Security Properties

1. **Tamper Detection:** SHA-256 hash of payload ensures data integrity
2. **Timing Attack Resistance:** Constant-time hash comparison prevents side-channel attacks
3. **Multi-Party Approval:** Governance rules require threshold of user signatures
4. **SuperAdmin Authorization:** Rules container must be signed by trusted administrators
5. **Explicit Hash Coverage:** Each signature explicitly lists which hashes it covers
6. **Legacy Compatibility:** Backward-compatible hash computation for schema evolution

## Usage Examples

### Basic Retrieval (No Verification)

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

const client = ProtectClient.create({
  host: 'https://protect.example.com',
  apiKey: 'your-api-key',
  apiSecret: 'your-hex-secret',
});

// Get address without verification
const address = await client.whitelistedAddresses.get('123');
console.log(`Address: ${address.address}`);
```

### Full Verification

```typescript
import { WhitelistedAddressService } from '@taurushq/protect-sdk';
import { rulesContainerFromBase64, userSignaturesFromBase64 } from '@taurushq/protect-sdk';

// Create service with verification enabled
const verifiedService = WhitelistedAddressService.withVerification(api, {
  superAdminKeysPem: [superAdminKey1, superAdminKey2],
  minValidSignatures: 2,
  rulesContainerDecoder: rulesContainerFromBase64,
  userSignaturesDecoder: userSignaturesFromBase64,
});

try {
  // Get address with full 6-step verification
  const result = await verifiedService.getWithVerification('123');
  console.log('Verified address:', result.verifiedWhitelistedAddress);
  console.log('Verified hash:', result.verifiedHash);
} catch (error) {
  if (error instanceof IntegrityError) {
    console.error('Integrity verification failed:', error.message);
    // DO NOT trust the address - verification failed
  }
}
```

### Using the Verifier Directly

```typescript
import { WhitelistedAddressVerifier } from '@taurushq/protect-sdk';
import { rulesContainerFromBase64, userSignaturesFromBase64 } from '@taurushq/protect-sdk';

const verifier = new WhitelistedAddressVerifier({
  superAdminKeysPem: [superAdminKey1, superAdminKey2],
  minValidSignatures: 2,
});

// Verify an envelope manually
try {
  const result = verifier.verify(
    envelope,
    rulesContainerFromBase64,
    userSignaturesFromBase64
  );
  console.log('Verified address:', result.verifiedWhitelistedAddress);
} catch (error) {
  if (error instanceof IntegrityError) {
    console.error('Cryptographic verification failed');
  } else if (error instanceof WhitelistError) {
    console.error('Governance thresholds not met');
  }
}
```

If verification fails at any step, an exception is thrown. Callers are guaranteed to only receive cryptographically verified data when using verification methods.
