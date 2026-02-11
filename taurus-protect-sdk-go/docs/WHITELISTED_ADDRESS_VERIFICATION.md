# Whitelisted Address Integrity Verification

This document explains the cryptographic integrity verification process for whitelisted addresses in the Taurus-PROTECT Go SDK.

## Overview

The SDK implements a **6-step cryptographic verification pipeline** to ensure whitelisted addresses are authentic and properly approved according to governance rules. Verification ensures data integrity and proper authorization.

## Verification Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                    API Response (WhitelistedAddress)                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────┐ │
│  │ SignedAddress   │  │ Metadata        │  │ RulesContainer      │ │
│  │ • Payload       │  │ • Address       │  │ • Groups            │ │
│  │ • Signatures[]  │  │ • Name          │  │ • Thresholds        │ │
│  └─────────────────┘  └─────────────────┘  │ RulesSignatures     │ │
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
│  Uses ECDSA P-256 (secp256r1) signatures                            │
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
│  Parse the verified payloadAsString JSON into domain model          │
│  All fields come from the cryptographically verified source         │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
                    ✓ Verified WhitelistedAddress returned
```

## Step-by-Step Details

### Step 1: Metadata Hash Verification

**Purpose:** Ensure the payload data hasn't been tampered with.

**Process:**
1. Extract the raw payload string from the metadata
2. Compute SHA-256 hash using `crypto.CalculateHexHash()`
3. Compare with the provided hash
4. Use constant-time comparison to prevent timing attacks

**Go Implementation:**
```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"

computedHash := crypto.CalculateHexHash(metadata.PayloadAsString)
if !helper.ConstantTimeCompare(computedHash, providedHash) {
    return nil, &model.IntegrityError{Message: "hash mismatch"}
}
```

**Failure:** Returns `IntegrityError` if hashes don't match.

### Step 2: Rules Container Signature Verification

**Purpose:** Verify that SuperAdmins approved the governance rules.

**Process:**
1. Decode `rulesSignatures` from the response
2. For each signature, verify against configured SuperAdmin public keys
3. Count valid signatures
4. Require `validCount >= minValidSignatures`

**Cryptographic Algorithm:** ECDSA with P-256 curve (secp256r1)

**Go Implementation:**
```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"

validCount := 0
rulesContainerBytes := []byte(rulesContainer)

for _, sig := range signatures {
    for _, pubKey := range superAdminKeys {
        valid, err := crypto.VerifySignature(pubKey, rulesContainerBytes, sig.Signature)
        if err == nil && valid {
            validCount++
            break
        }
    }
}

if validCount < minValidSignatures {
    return nil, &IntegrityError{
        Message: fmt.Sprintf("only %d valid signatures, minimum %d required",
            validCount, minValidSignatures),
    }
}
```

**Failure:** Returns `IntegrityError` if insufficient valid signatures.

### Step 3: Rules Container Decoding

**Purpose:** Extract governance structure (groups, users, thresholds).

**Process:**
1. Base64-decode the `rulesContainer` string
2. Parse as protobuf message
3. Map to `DecodedRulesContainer` domain model

**Contains:**
- User groups with member IDs
- Parallel threshold paths (OR logic)
- Sequential group thresholds (AND logic)
- User public keys for signature verification

### Step 4: Hash Coverage Verification

**Purpose:** Confirm the metadata hash is covered by at least one signature.

**Process:**
1. Get `metadata.hash` (the computed hash from Step 1)
2. Check each `signature.Hashes[]` list
3. At least one signature must include this hash

**Why needed:** Signatures cover a list of hashes. This ensures the specific address hash was actually included in what was signed (not just any signature exists).

**Go Implementation:**
```go
// Try the provided hash first (using constant-time comparison)
if helper.VerifyHashCoverage(metadataHash, signedAddress.Signatures) {
    return metadataHash, nil
}

// Try legacy hashes for backward compatibility
legacyHashes := helper.ComputeLegacyHashes(metadata.PayloadAsString)
for _, legacyHash := range legacyHashes {
    if helper.VerifyHashCoverage(legacyHash, signedAddress.Signatures) {
        return legacyHash, nil
    }
}

return "", &model.IntegrityError{Message: "hash not covered by any signature"}
```

**Legacy Hash Support:**

For backward compatibility with addresses signed before schema changes, the SDK computes alternative hashes if the provided hash isn't found:

| Strategy | Description | Use Case |
|----------|-------------|----------|
| 1 | Remove `contractType` field | Addresses signed before contractType was added to schema |
| 2 | Remove `label` from linkedInternalAddresses | Addresses signed before labels were added to linked addresses |
| 3 | Remove both `contractType` AND `label` | Earliest schema version (before both fields existed) |

The legacy hash computation uses regex patterns to strip these fields while preserving JSON structure. This is handled by `helper.ComputeLegacyHashes()`.

**Failure:** Returns `IntegrityError` if hash (including legacy variants) not found in any signature's hashes list.

### Step 5: Whitelist Signature Verification

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
1. **Find applicable rules** based on blockchain/network (using 3-tier priority: exact match > blockchain-only > global default)
2. **Determine thresholds** based on linked wallets/addresses (see Threshold Selection below)
3. **Try each parallel path** (OR logic - only one needs to succeed):
   - For each sequential threshold path:
     - Verify ALL group thresholds (AND logic)
     - If all pass: verification succeeds
4. **Verify group threshold:**
   - Find group by ID in rules container
   - For each signature from a user in this group:
     - Check user is in group
     - Check signature covers metadata hash (using constant-time comparison)
     - Get user's public key from rules container
     - Verify signature using ECDSA P-256
   - Count valid signatures
   - Require `validCount >= minimumSignatures`

**Threshold Selection:**

The SDK determines which approval thresholds to apply based on the address configuration:

| Condition | Threshold Used |
|-----------|----------------|
| Address has linked internal addresses | Default thresholds (from `rules.ParallelThresholds`) |
| Address has 0 or 2+ linked wallets | Default thresholds |
| Address has exactly 1 linked wallet AND no linked addresses | Line-specific thresholds if wallet path matches a rule line, otherwise default |

When checking rule lines:
1. Get the single linked wallet's path
2. Search rule lines for a matching wallet path (in the first cell of type `InternalWallet`)
3. If found, use that line's `ParallelThresholds`
4. If not found, fall back to the rule's default `ParallelThresholds`

**Signature Data Format:** The hashes array is JSON-encoded, then signed with ECDSA.

```go
// Verification pseudocode
hashesJSON, _ := json.Marshal(signature.Hashes)
valid, err := crypto.VerifySignature(userPublicKey, hashesJSON, signature.UserSignature.Signature)
```

## Data Models

### WhitelistedAddress

```go
type WhitelistedAddress struct {
    ID                 string
    TenantID           string
    Address            string
    Name               string
    Blockchain         string
    Network            string
    Status             string
    RulesContainer     string                      // Base64 protobuf
    RulesContainerHash string
    RulesSignatures    string                      // Serialized signatures
    SignedAddress      *SignedWhitelistedAddress   // Payload + signatures
    Approvers          *Approvers                  // Required approval structure
    Metadata           *WhitelistedAssetMetadata
    Trails             []Trail
}
```

### SignedWhitelistedAddress

```go
type SignedWhitelistedAddress struct {
    Payload    string                // Base64-encoded signed payload
    Signatures []WhitelistSignature  // List of user signatures
}
```

### WhitelistSignature

```go
type WhitelistSignature struct {
    UserSignature *WhitelistUserSignature  // Who signed and their signature
    Hashes        []string                 // Which hashes were covered
}

type WhitelistUserSignature struct {
    UserID    string  // Signer identity
    Signature string  // Base64-encoded ECDSA signature
    Comment   string  // Optional comment
}
```

### Approvers

```go
type Approvers struct {
    Parallel []ParallelApproversGroup
}

type ParallelApproversGroup struct {
    Sequential []SequentialApproversGroup
}

type SequentialApproversGroup struct {
    GroupID           string
    MinimumSignatures int
    Users             []ApproverUser
}
```

## Cryptographic Primitives

| Operation | Algorithm | Implementation |
|-----------|-----------|----------------|
| Hash | SHA-256 | `crypto/sha256` |
| Signature | ECDSA P-256 | `crypto/ecdsa` |
| Comparison | Constant-time | `crypto/subtle` |
| Encoding | Base64 | `encoding/base64` |

## Error Types

| Error | Type | When Returned |
|-------|------|---------------|
| `IntegrityError` | Unchecked | Hash mismatch, insufficient signatures, cryptographic failure |
| `WhitelistError` | Checked | Decode errors, missing data, no matching rules |

## Key Functions

### Core Cryptographic Functions

| Function | Location | Purpose |
|----------|----------|---------|
| `CalculateHexHash` | `crypto/tpv1.go` | SHA-256 hash computation (returns hex string) |
| `VerifySignature` | `crypto/tpv1.go` | ECDSA P-256 signature verification |
| `DecodePublicKeyPEM` | `crypto/tpv1.go` | Parse PEM-encoded ECDSA public key |

### Verification Helper Functions

| Function | Location | Purpose |
|----------|----------|---------|
| `VerifyGovernanceRulesSignatures` | `helper/signature_verifier.go` | Verify SuperAdmin signatures on rules container |
| `VerifyHashCoverage` | `helper/signature_verifier.go` | Check if hash appears in any signature's hashes list |
| `ConstantTimeCompare` | `helper/signature_verifier.go` | Timing-safe string comparison |
| `IsValidSignature` | `helper/signature_verifier.go` | Check signature against multiple public keys |
| `DecodeBase64` | `helper/signature_verifier.go` | Decode base64-encoded data |

### Hash and Payload Functions

| Function | Location | Purpose |
|----------|----------|---------|
| `ComputeLegacyHashes` | `helper/whitelist_hash.go` | Compute backward-compatible hash variants |
| `ParseWhitelistedAddressFromJSON` | `helper/whitelist_hash.go` | Parse verified JSON payload to domain model |

### Rules Container Functions

| Function | Location | Purpose |
|----------|----------|---------|
| `RulesContainerFromBase64` | `mapper/rules_container.go` | Decode base64 protobuf to DecodedRulesContainer |
| `UserSignaturesFromBase64` | `mapper/rules_container.go` | Decode base64 protobuf to RuleUserSignature slice |
| `FindAddressWhitelistingRules` | `model/rules_container.go` | Find matching rules by blockchain/network |
| `FindUserByID` | `model/rules_container.go` | Look up user by ID in rules container |
| `FindGroupByID` | `model/rules_container.go` | Look up group by ID in rules container |

## Security Properties

1. **Tamper Detection:** SHA-256 hash of payload ensures data integrity
2. **Timing Attack Resistance:** Constant-time hash comparison prevents side-channel attacks
3. **Multi-Party Approval:** Governance rules require threshold of user signatures
4. **SuperAdmin Authorization:** Rules container must be signed by trusted administrators
5. **Explicit Hash Coverage:** Each signature explicitly lists which hashes it covers

## Usage Example

```go
import (
    "context"
    "fmt"

    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

func verifyWhitelistedAddress(ctx context.Context, client *protect.Client, id string) error {
    // Verification happens automatically when retrieving addresses
    // The SDK will return an error if verification fails
    addr, err := client.WhitelistedAddresses().GetWhitelistedAddress(ctx, id)
    if err != nil {
        if protect.IsIntegrityError(err) {
            // Cryptographic verification failed - investigate!
            return fmt.Errorf("integrity check failed: %w", err)
        }
        if protect.IsWhitelistError(err) {
            // Whitelist-specific error
            return fmt.Errorf("whitelist error: %w", err)
        }
        return err
    }

    // Address passed all verification checks
    fmt.Printf("Verified address: %s\n", addr.Address)
    fmt.Printf("  Blockchain: %s/%s\n", addr.Blockchain, addr.Network)
    fmt.Printf("  Name: %s\n", addr.Name)
    fmt.Printf("  Status: %s\n", addr.Status)

    // Access approval information
    if addr.Approvers != nil {
        fmt.Printf("  Approval groups: %d\n", len(addr.Approvers.Parallel))
    }

    return nil
}
```

## Verification Configuration

When creating the client, configure verification parameters:

```go
client, err := protect.NewClient(
    host,
    protect.WithCredentials(apiKey, apiSecret),
    protect.WithSuperAdminKeysPEM(superAdminKeys),  // Required for verification
    protect.WithMinValidSignatures(2),               // Minimum SuperAdmin sigs
    protect.WithRulesCacheTTL(5 * time.Minute),     // Cache validated rules
)
```

If SuperAdmin keys are not configured, the client can still retrieve whitelisted addresses but won't perform cryptographic verification of the governance rules.

## Manual Verification

For advanced use cases, you can manually verify signatures:

```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"

// Verify a signature manually
data := []byte("data to verify")
signature := "base64-encoded-signature"
publicKey, _ := crypto.DecodePublicKeyPEM(pemKeyString)

valid, err := crypto.VerifySignature(publicKey, data, signature)
if err != nil {
    return fmt.Errorf("verification error: %w", err)
}
if !valid {
    return fmt.Errorf("invalid signature")
}
```

## Troubleshooting

### Common Verification Failures

| Error | Likely Cause | Solution |
|-------|--------------|----------|
| "hash mismatch" | Data modified in transit | Check network security, retry request |
| "insufficient signatures" | SuperAdmin keys don't match | Verify configured keys match production keys |
| "hash not covered" | Signature doesn't cover this address | Address may need re-approval |
| "group threshold not met" | Not enough user approvals | Wait for required approvers |

### Debugging Tips

1. Enable verbose logging to see verification steps
2. Check that SuperAdmin keys match the production environment
3. Verify the `minValidSignatures` setting matches governance requirements
4. Ensure rules cache is not stale (reduce TTL for debugging)

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Authentication](AUTHENTICATION.md) - Security and signing
- [Services Reference](SERVICES.md) - Complete API documentation
- [Key Concepts](CONCEPTS.md) - Domain model overview