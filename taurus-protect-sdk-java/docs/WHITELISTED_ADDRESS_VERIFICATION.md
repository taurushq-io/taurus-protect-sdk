# Whitelisted Address 6-Step Integrity Verification

This document explains the cryptographic integrity verification process for whitelisted addresses in the Taurus PROTECT Java SDK.

## Overview

The SDK implements a **6-step cryptographic verification pipeline** to ensure whitelisted addresses are authentic and properly approved according to governance rules. Verification happens automatically when retrieving addresses - callers never receive unverified data.

## Verification Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                    API Response (Envelope)                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────┐ │
│  │ signedAddress   │  │ metadata        │  │ rulesContainer      │ │
│  │ • payload       │  │ • payloadAsStr  │  │ • groups            │ │
│  │ • signatures[]  │  │ • hash          │  │ • thresholds        │ │
│  └─────────────────┘  └─────────────────┘  │ rulesSignatures     │ │
│                                            └─────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│              STEP 1: Verify Metadata Hash                           │
│  SHA-256(payloadAsString) == metadata.hash ?                        │
│  Uses constant-time comparison to prevent timing attacks            │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│         STEP 2: Verify Rules Container Signatures                   │
│  SuperAdmin signatures on rulesContainer >= minValidSignatures      │
│  Uses SHA256withPLAIN-ECDSA algorithm                               │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│              STEP 3: Decode Rules Container                         │
│  Base64 -> Protobuf -> DecodedRulesContainer                        │
│  Contains groups, users, and approval thresholds                    │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│         STEP 4: Verify Hash Coverage                                │
│  metadata.hash must appear in at least one signature.hashes[]       │
│  Ensures the payload hash was actually signed                       │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│         STEP 5: Verify Whitelist Signatures                         │
│  User signatures meet governance threshold requirements             │
│  Parallel paths (OR) -> Sequential groups (AND) -> Min signatures   │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│         STEP 6: Parse WhitelistedAddress from Verified Payload      │
│  Extract WhitelistedAddress fields from the verified payload        │
│  Only uses data from the cryptographically verified source          │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
                    Verified WhitelistedAddress returned
```

## Step-by-Step Details

### Step 1: Metadata Hash Verification

**Purpose:** Ensure the payload data hasn't been tampered with.

**Location:** `WhitelistedAddressService.java:173-192`

**Process:**
1. Compute SHA-256 hash of `metadata.payloadAsString` using `CryptoTPV1.calculateHexHash()`
2. Compare with provided `metadata.hash`
3. Uses **constant-time comparison** (`constantTimeAreEqual`) to prevent timing attacks

**Failure:** Throws `IntegrityException` if hashes don't match.

```java
String computedHash = CryptoTPV1.calculateHexHash(envelope.getMetadata().getPayloadAsString());
if (!constantTimeAreEqual(computedHash, providedHash)) {
    throw new IntegrityException(...);
}
```

### Step 2: Rules Container Signature Verification

**Purpose:** Verify that SuperAdmins approved the governance rules.

**Location:** `WhitelistedAddressService.java:197-237`

**Process:**
1. Decode `rulesSignatures` from Base64 to Protobuf `UserSignatures`
2. For each signature, verify against SuperAdmin public keys
3. Count valid signatures
4. Require `validCount >= minValidSignatures`

**Cryptographic Algorithm:** SHA256withPLAIN-ECDSA

**Failure:** Throws `IntegrityException` if insufficient valid signatures.

### Step 3: Rules Container Decoding

**Purpose:** Extract governance structure (groups, users, thresholds).

**Location:** `WhitelistedAddressService.java:242-249`

**Process:**
1. Decode Base64 `rulesContainer` to protobuf bytes
2. Map to `DecodedRulesContainer` domain model

**Contains:**
- User groups with member IDs
- Parallel threshold paths (OR logic)
- Sequential group thresholds (AND logic)
- User public keys for signature verification

### Step 4: Hash Coverage Verification

**Purpose:** Confirm the metadata hash is covered by at least one signature.

**Location:** `WhitelistedAddressService.java:254-271`

**Process:**
1. Get `metadata.hash`
2. Check each `signature.hashes[]` list
3. At least one signature must include this hash

**Why needed:** Signatures cover a list of hashes. This ensures the specific address hash was actually included in what was signed (not just any signature exists).

**Failure:** Throws `IntegrityException` if hash not found in any signature's hashes list.

### Step 5: Whitelist Signature Verification

**Purpose:** Verify user approvals meet governance threshold requirements.

**Location:** `WhitelistedAddressService.java:276-491`

**Governance Structure:**
```
ParallelThresholds (OR paths)
  └── SequentialThresholds (AND groups within a path)
        └── GroupThreshold
              ├── groupId
              └── minimumSignatures
```

**Algorithm:**
1. **Find applicable rules** based on blockchain/network
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
     - Verify signature: `SHA256withPLAIN-ECDSA(JSON(hashes[]), userPublicKey)`
   - Count valid signatures
   - Require `validCount >= minimumSignatures`

**Signature Data Format:** JSON array of hashes converted to UTF-8 bytes, then ECDSA signed.

```java
String hashesJson = GSON.toJson(hashes);  // e.g., ["abc123","def456"]
byte[] hashesBytes = hashesJson.getBytes(StandardCharsets.UTF_8);
SignatureVerifier.verifySignature(hashesBytes, userSig.getSignature(), user.getPublicKey());
```

### Step 6: Parse WhitelistedAddress from Verified Payload

**Purpose:** Extract the WhitelistedAddress domain object from the cryptographically verified payload.

**Process:**
1. Parse `metadata.payloadAsString` (the verified source) into a `WhitelistedAddress` object
2. All fields (blockchain, network, address, label, memo, etc.) are extracted only from the verified payload
3. Non-security fields (status, action, trails) may come from the DTO

**Security:** This step ensures that the returned `WhitelistedAddress` contains only data that has passed all prior verification steps. Fields are never sourced from unverified DTO attributes.

## Cryptographic Primitives

| Operation | Algorithm | Library |
|-----------|-----------|---------|
| Hash | SHA-256 | Apache Commons Codec |
| Signature | SHA256withPLAIN-ECDSA | BouncyCastle |
| Comparison | Constant-time | BouncyCastle |

## Exception Types

| Exception | Type | When Thrown |
|-----------|------|-------------|
| `IntegrityException` | Unchecked (`SecurityException`) | Hash mismatch, insufficient signatures, cryptographic failure |
| `WhitelistException` | Checked | Decode errors, missing data, no matching rules |

## Key Files

| File | Purpose |
|------|---------|
| `WhitelistedAddressService.java` | Main verification orchestration |
| `SignatureVerifier.java` | Low-level signature verification |
| `WhitelistHashHelper.java` | Hash operations and JSON parsing |
| `WhitelistIntegrityHelper.java` | Field-level integrity checks |
| `CryptoTPV1.java` | Cryptographic primitives |

## Security Properties

1. **Tamper Detection:** SHA-256 hash of payload ensures data integrity
2. **Timing Attack Resistance:** Constant-time hash comparison prevents side-channel attacks
3. **Multi-Party Approval:** Governance rules require threshold of user signatures
4. **SuperAdmin Authorization:** Rules container must be signed by trusted administrators
5. **Explicit Hash Coverage:** Each signature explicitly lists which hashes it covers

## Usage Example

```java
// Verification is automatic - addresses are verified before return
WhitelistedAddress addr = whitelistedAddressService.getWhitelistedAddress(123L);

// Or get full envelope with verification details
SignedWhitelistedAddressEnvelope envelope =
    whitelistedAddressService.getWhitelistedAddressEnvelope(123L);
WhitelistedAddress verifiedAddr = envelope.getWhitelistedAddress();
DecodedRulesContainer rules = envelope.getVerifiedRulesContainer();
```

If verification fails at any step, an exception is thrown. Callers are guaranteed to only receive cryptographically verified data.
