# Integrity Verification

This document describes the cryptographic integrity verification processes used by Taurus-PROTECT SDKs. These verification flows ensure data authenticity and proper authorization.

## Overview

The SDKs implement cryptographic verification at multiple points to ensure:

- **Data Integrity** - Data hasn't been tampered with in transit or at rest
- **Authorization** - Actions were approved by the correct parties
- **Non-repudiation** - Approvers cannot deny having signed
- **Audit Trail** - All signatures are recorded for compliance

---

## What Gets Verified

| Entity | Verification |
|--------|--------------|
| **Governance Rules** | SuperAdmin signatures on rules container |
| **Whitelisted Address** | Metadata hash + SuperAdmin signatures + user signatures |
| **Request Metadata** | SHA-256 hash verification |
| **Address** | Signature proves address was created by Taurus-PROTECT |

---

## Governance Rules Verification

Governance rules define who can approve transactions and whitelist addresses. They are signed by SuperAdmins.

### Data Structure

```
GovernanceRuleset
├── rulesContainer     (Base64-encoded protobuf)
├── signatures[]       (SuperAdmin ECDSA signatures)
└── locked             (Whether rules can be modified)
```

### Verification Process

1. **Decode rules container** - Base64 decode to raw bytes
2. **Calculate hash** - SHA-256 of the raw bytes
3. **Verify signatures** - For each signature:
   - Decode Base64 signature
   - Try verification against each configured SuperAdmin public key
   - Count successful verifications
4. **Check threshold** - Ensure `validCount >= minValidSignatures`

### Failure Conditions

- Hash calculation fails
- Signature decode fails
- No signature matches any SuperAdmin key
- Fewer valid signatures than required threshold

---

## Whitelisted Address Verification

Whitelisted addresses undergo comprehensive multi-step verification.

### Data Structure

```
WhitelistedAddress
├── metadata
│   ├── payloadAsString   (Raw JSON for hashing)
│   └── hash              (Expected SHA-256 hash)
├── rulesContainer        (Governance rules)
├── rulesSignatures       (SuperAdmin signatures)
└── signedAddress
    ├── payload           (Base64-encoded address data)
    └── signatures[]      (User approval signatures)
        ├── userSignature
        │   ├── userId
        │   └── signature
        └── hashes[]      (Hashes covered by this signature)
```

### 6-Step Verification Flow

```
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
│  Uses ECDSA P-256 signature verification                            │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│              STEP 3: Decode Rules Container                         │
│  Base64 -> Protobuf -> Decoded rules structure                      │
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
│         STEP 5: Verify User Signatures                              │
│  User signatures meet governance threshold requirements             │
│  Parallel paths (OR) -> Sequential groups (AND) -> Min signatures   │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│    STEP 6: Parse WhitelistedAddress from Verified Payload           │
│  Parse the verified payloadAsString JSON into domain model          │
│  All fields come from the cryptographically verified source         │
└─────────────────────────────────────────────────────────────────────┘
```

### Step Details

#### Step 1: Metadata Hash Verification

**Purpose:** Ensure the payload data hasn't been tampered with.

**Process:**
1. Compute SHA-256 hash of `metadata.payloadAsString`
2. Compare with provided `metadata.hash`
3. Use constant-time comparison to prevent timing attacks

#### Step 2: Rules Container Signature Verification

**Purpose:** Verify that SuperAdmins approved the governance rules.

**Process:**
1. For each signature in `rulesSignatures`:
   - Verify against SuperAdmin public keys
   - Count valid signatures
2. Ensure `validCount >= minValidSignatures`

#### Step 3: Decode Rules Container

**Purpose:** Extract governance structure for threshold verification.

**Contains:**
- User groups with member IDs
- Parallel threshold paths (OR logic)
- Sequential group thresholds (AND logic)
- User public keys for signature verification

#### Step 4: Hash Coverage Verification

**Purpose:** Confirm the metadata hash is covered by at least one signature.

**Why needed:** Signatures cover a list of hashes. This ensures the specific address hash was actually included in what was signed.

**Process:**
1. Get `metadata.hash`
2. Check each `signature.hashes[]` list
3. At least one signature must include this hash

#### Step 5: User Signature Verification

**Purpose:** Verify user approvals meet governance threshold requirements.

**Governance Structure:**
```
ParallelThresholds (OR paths)
  └── SequentialThresholds (AND groups within a path)
        └── GroupThreshold
              ├── groupId
              └── minimumSignatures
```

**Algorithm:**
1. Find applicable rules based on blockchain/network
2. Determine thresholds based on linked wallets/addresses
3. Try each parallel path (OR logic - only one needs to succeed):
   - For each sequential group in the path (AND logic):
     - Verify group threshold is met
   - If all groups pass: verification succeeds
4. Group threshold verification:
   - Find group by ID in rules container
   - For each signature from a user in this group:
     - Verify signature covers metadata hash
     - Verify signature using user's public key
   - Count valid signatures
   - Require `validCount >= minimumSignatures`

---

## Request Metadata Verification

When retrieving or approving requests, verify the metadata hasn't been tampered.

### Verification Process

1. Extract `metadata.payloadAsString`
2. Compute SHA-256 hash
3. Compare with `metadata.hash`
4. Verify the hash matches before signing approval

---

## Cryptographic Primitives

| Operation | Algorithm | Notes |
|-----------|-----------|-------|
| Hash | SHA-256 | Payload integrity |
| Signature | ECDSA P-256 | SuperAdmin and user signatures |
| Comparison | Constant-time | Timing attack prevention |
| Encoding | Base64 | Signature transport |

---

## Error Types

### IntegrityError

Thrown when cryptographic verification fails:

- Hash mismatch (payload tampered)
- Insufficient valid SuperAdmin signatures
- Signature decode failure
- Cryptographic verification failure

### WhitelistError

Thrown for whitelist-specific failures:

- No matching rules for blockchain/network
- User signature verification failed
- Group threshold not met
- Hash not covered by any signature

---

## Security Properties

1. **Tamper Detection** - SHA-256 hash ensures data integrity
2. **Timing Attack Resistance** - Constant-time comparison
3. **Multi-Party Approval** - Governance requires threshold signatures
4. **SuperAdmin Authorization** - Rules must be signed by trusted admins
5. **Explicit Hash Coverage** - Signatures explicitly list covered hashes

---

## When Verification Happens

### Automatic Verification

SDKs can perform automatic verification when:
- SuperAdmin keys are configured
- Retrieving governance rules
- Fetching whitelisted addresses (when verification is enabled)

### Manual Verification

Some SDKs allow manual verification for advanced use cases:
- Custom verification workflows
- Offline verification
- Audit and compliance checks

---

## Caching Verified Rules

For performance, verified governance rules are cached:

- **TTL** - Configurable time-to-live (default: 5 minutes)
- **Invalidation** - Manual cache clear available
- **Thread-safe** - Safe for concurrent access

---

## Related Documentation

### Common Documentation
- [Key Concepts](CONCEPTS.md) - Domain model and entities
- [Authentication](AUTHENTICATION.md) - TPV1 authentication protocol

### SDK-Specific Documentation

**Java SDK**
- [Whitelisted Address Verification](../taurus-protect-sdk-java/docs/WHITELISTED_ADDRESS_VERIFICATION.md) - Java implementation details

**Go SDK**
- [Whitelisted Address Verification](../taurus-protect-sdk-go/docs/WHITELISTED_ADDRESS_VERIFICATION.md) - Go implementation details

**Python SDK**
- [Whitelisted Address Verification](../taurus-protect-sdk-python/docs/WHITELISTED_ADDRESS_VERIFICATION.md) - Python implementation details

**TypeScript SDK**
- [Whitelisted Address Verification](../taurus-protect-sdk-typescript/docs/WHITELISTED_ADDRESS_VERIFICATION.md) - TypeScript implementation details