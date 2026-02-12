# Whitelisted Address Verification

This document describes the 6-step verification flow for whitelisted addresses in the Python SDK.

## Overview

Whitelisted addresses undergo a rigorous 6-step cryptographic verification pipeline to ensure:

1. **Data Integrity** - The payload hasn't been tampered with
2. **Authorization** - SuperAdmin signatures are valid
3. **Governance Compliance** - Approval thresholds are met
4. **Hash Coverage** - The specific address hash is signed
5. **User Approval** - Required users have signed

This verification is performed automatically by `WhitelistedAddressService` when SuperAdmin keys are configured.

---

## Verification Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    5-STEP VERIFICATION PIPELINE                              │
└─────────────────────────────────────────────────────────────────────────────┘

    Step 1: Verify Metadata Hash
    ┌─────────────────────────────────────────────────────────────────────────┐
    │  SHA-256(payloadAsString) == metadata.hash                              │
    │                                                                         │
    │  ┌─────────────────┐         ┌─────────────────┐                       │
    │  │ payloadAsString │ ──SHA256──► computed_hash │                       │
    │  └─────────────────┘         └─────────────────┘                       │
    │                                      │                                  │
    │                                      ▼                                  │
    │                              ┌─────────────────┐                       │
    │                              │ metadata.hash   │                       │
    │                              └─────────────────┘                       │
    │                                      │                                  │
    │                           constant_time_compare()                       │
    └─────────────────────────────────────────────────────────────────────────┘
                                          │
                                          ▼
    Step 2: Verify Rules Container Signatures
    ┌─────────────────────────────────────────────────────────────────────────┐
    │  For each signature in rulesSignatures:                                 │
    │    Verify ECDSA signature against SuperAdmin public keys                │
    │    Count valid signatures                                               │
    │  Require: validCount >= minValidSignatures                              │
    │                                                                         │
    │  ┌─────────────────┐       ┌───────────────────────┐                   │
    │  │ rulesContainer  │       │ SuperAdmin Keys [N]   │                   │
    │  │   (base64)      │       │                       │                   │
    │  └─────────────────┘       └───────────────────────┘                   │
    │          │                           │                                  │
    │          ▼                           ▼                                  │
    │  ┌─────────────────┐       ┌───────────────────────┐                   │
    │  │ rulesSignatures │ ────► │ ECDSA Verify (P-256)  │                   │
    │  └─────────────────┘       └───────────────────────┘                   │
    └─────────────────────────────────────────────────────────────────────────┘
                                          │
                                          ▼
    Step 3: Decode Rules Container
    ┌─────────────────────────────────────────────────────────────────────────┐
    │  Base64 decode ─► JSON parse ─► DecodedRulesContainer                   │
    │                                                                         │
    │  ┌─────────────────┐       ┌───────────────────────┐                   │
    │  │ rulesContainer  │ ────► │ DecodedRulesContainer │                   │
    │  │   (base64)      │       │   - users             │                   │
    │  └─────────────────┘       │   - groups            │                   │
    │                            │   - whitelisting rules│                   │
    │                            └───────────────────────┘                   │
    └─────────────────────────────────────────────────────────────────────────┘
                                          │
                                          ▼
    Step 4: Verify Hash Coverage
    ┌─────────────────────────────────────────────────────────────────────────┐
    │  metadata.hash must appear in at least one signature's hashes list      │
    │  (with legacy hash fallback for backward compatibility)                 │
    │                                                                         │
    │  ┌─────────────────┐       ┌───────────────────────┐                   │
    │  │ metadata.hash   │ ────► │ signatures[].hashes[] │                   │
    │  └─────────────────┘       └───────────────────────┘                   │
    │         │                                                               │
    │         │ (if not found)                                                │
    │         ▼                                                               │
    │  ┌─────────────────┐       ┌───────────────────────┐                   │
    │  │ legacy hashes   │ ────► │ signatures[].hashes[] │                   │
    │  └─────────────────┘       └───────────────────────┘                   │
    └─────────────────────────────────────────────────────────────────────────┘
                                          │
                                          ▼
    Step 5: Verify Whitelist Signatures
    ┌─────────────────────────────────────────────────────────────────────────┐
    │  Find matching AddressWhitelistingRules                                 │
    │  Determine applicable thresholds (rule lines or default)                │
    │  For each parallel threshold path (OR logic):                           │
    │    For each group threshold (AND logic):                                │
    │      Count valid signatures from group members                          │
    │      Require: count >= threshold.minimumSignatures                      │
    │  At least one path must succeed                                         │
    │                                                                         │
    │  ┌─────────────────────────────────────────────────────────────────┐   │
    │  │ AddressWhitelistingRules                                        │   │
    │  │   ├── lines[] (wallet-specific overrides)                       │   │
    │  │   │     └── parallelThresholds[]                                │   │
    │  │   └── parallelThresholds[] (default)                            │   │
    │  │         └── SequentialThresholds                                │   │
    │  │               └── thresholds[] (GroupThreshold)                 │   │
    │  │                     ├── groupId                                 │   │
    │  │                     └── minimumSignatures                       │   │
    │  └─────────────────────────────────────────────────────────────────┘   │
    └─────────────────────────────────────────────────────────────────────────┘
                                          │
                                          ▼
                              ┌───────────────────────┐
                              │  VERIFICATION PASSED  │
                              └───────────────────────┘
```

---

## Step-by-Step Details

### Step 1: Verify Metadata Hash

Ensures the payload data hasn't been modified.

```python
from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.helpers.constant_time import constant_time_compare

# Compute hash of the raw payload string
computed_hash = calculate_hex_hash(envelope.metadata.payload_as_string)

# Compare with provided hash using constant-time comparison
if not constant_time_compare(computed_hash, envelope.metadata.hash):
    raise IntegrityError("metadata hash verification failed")
```

**Why constant-time?** Prevents timing attacks that could reveal information about the expected hash.

### Step 2: Verify Rules Container Signatures

Ensures the governance rules were signed by trusted SuperAdmins.

```python
import base64
from taurus_protect.helpers.signature_verifier import is_valid_signature

# Decode rules container data
rules_data = base64.b64decode(envelope.rules_container)

# Decode signatures (protobuf UserSignatures)
signatures = user_signatures_decoder(envelope.rules_signatures)

# Count valid signatures
valid_count = 0
for sig in signatures:
    if sig.signature and is_valid_signature(rules_data, sig.signature, super_admin_keys):
        valid_count += 1

if valid_count < min_valid_signatures:
    raise IntegrityError(
        f"rules container signature verification failed: only {valid_count} valid signatures, "
        f"minimum {min_valid_signatures} required"
    )
```

### Step 3: Decode Rules Container

Parses the governance rules for use in subsequent verification.

```python
from taurus_protect.mappers.governance_rules import rules_container_from_base64

# Decode and parse the rules container
rules_container = rules_container_from_base64(envelope.rules_container)
# Returns: DecodedRulesContainer with users, groups, whitelisting rules
```

### Step 4: Verify Hash Coverage

Ensures the specific address's hash was actually signed (not just any hash). Includes legacy hash fallback for backward compatibility.

```python
from taurus_protect.helpers.whitelist_hash_helper import compute_legacy_hashes

metadata_hash = envelope.metadata.hash
signatures = envelope.signed_address.signatures

# Try the provided hash first using constant-time comparison
if _verify_hash_coverage(metadata_hash, signatures):
    return metadata_hash

# Try legacy hashes for backward compatibility
legacy_hashes = compute_legacy_hashes(envelope.metadata.payload_as_string)
for legacy_hash in legacy_hashes:
    if _verify_hash_coverage(legacy_hash, signatures):
        return legacy_hash

raise IntegrityError("metadata hash is not covered by any signature")
```

**Legacy hash strategies** (for backward compatibility):
- Remove `contractType` field from payload
- Remove `label` from objects in payload
- Remove both `contractType` and `label`

### Step 5: Verify Whitelist Signatures

Ensures the required users approved the address according to governance rules. This step uses `AddressWhitelistingRules` from the decoded rules container.

**Threshold selection logic:**

The verifier first checks whether wallet-specific rule lines should be used. Rule lines are checked only when:
1. The envelope has **no linked internal addresses**, AND
2. The envelope has **exactly 1 linked wallet**

If a matching rule line is found for the linked wallet's path, its thresholds are used. Otherwise, the default `parallel_thresholds` from the `AddressWhitelistingRules` are used.

```python
# Find matching address whitelisting rules for this blockchain/network
whitelist_rules = rules_container.find_address_whitelisting_rules(
    envelope.blockchain, envelope.network
)

# Determine applicable thresholds (wallet-specific line or default)
parallel_thresholds = self._get_applicable_thresholds(whitelist_rules, envelope)

# Try each parallel path (OR logic - any path can succeed)
for path in parallel_thresholds:
    # Each path has sequential group thresholds (AND logic)
    all_groups_satisfied = True

    for group_threshold in path.thresholds:
        group = rules_container.find_group_by_id(group_threshold.group_id)
        min_sigs = group_threshold.get_min_signatures()

        # Count valid signatures from this group's members
        valid_count = count_valid_group_signatures(
            group, signatures, metadata_hash, rules_container
        )

        if valid_count < min_sigs:
            all_groups_satisfied = False
            break

    if all_groups_satisfied:
        return  # Verification passed!

raise WhitelistError("No approval path satisfied")
```

---

## Cryptographic Primitives

| Operation | Algorithm | Python Implementation |
|-----------|-----------|----------------------|
| Hash | SHA-256 | `hashlib.sha256()` |
| Signature | ECDSA P-256 | `cryptography.hazmat.primitives.asymmetric.ec` |
| Encoding | Base64 | `base64.b64decode()` / `base64.b64encode()` |
| Comparison | Constant-time | `hmac.compare_digest()` |

### ECDSA Signature Format

Signatures use raw r||s format (not DER-encoded):

```python
from taurus_protect.crypto.signing import verify_signature

# Signature is base64-encoded raw r||s (64 bytes for P-256)
is_valid = verify_signature(public_key, data, signature_b64)
```

---

## Exception Types

| Exception | When Raised | Action |
|-----------|-------------|--------|
| `IntegrityError` | Hash mismatch, insufficient SuperAdmin signatures | Do not retry. Investigate. |
| `WhitelistError` | User signature thresholds not met | Check approvals. May need more signatures. |
| `ValueError` | Missing required fields (metadata, signatures) | Fix input data. |

---

## Key Implementation Files

| File | Purpose |
|------|---------|
| `helpers/whitelisted_address_verifier.py` | Main verification logic |
| `helpers/whitelist_hash_helper.py` | Legacy hash computation |
| `helpers/signature_verifier.py` | SuperAdmin signature verification |
| `helpers/constant_time.py` | Timing-attack-safe comparison |
| `crypto/signing.py` | ECDSA sign/verify |
| `crypto/hashing.py` | SHA-256 hash computation |
| `mappers/governance_rules.py` | Rules container decoding |
| `models/governance_rules.py` | DecodedRulesContainer, AddressWhitelistingRules |

---

## Usage Example

### Automatic Verification

When SuperAdmin keys are configured, verification happens automatically:

```python
from taurus_protect import ProtectClient

super_admin_keys = [
    "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
    "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
]

with ProtectClient.create(
    host="https://api.protect.taurushq.com",
    api_key=api_key,
    api_secret=api_secret,
    super_admin_keys_pem=super_admin_keys,
    min_valid_signatures=2,
) as client:
    # Verification happens automatically on get() and list()
    try:
        addresses, _ = client.whitelisted_addresses.list()
        for addr in addresses:
            print(f"Verified address: {addr.address}")
    except IntegrityError as e:
        print(f"Verification failed: {e.message}")
```

### Manual Verification

For custom verification scenarios:

```python
from taurus_protect.helpers.whitelisted_address_verifier import WhitelistedAddressVerifier
from taurus_protect.mappers.governance_rules import (
    rules_container_from_base64,
    user_signatures_from_base64,
)

# Create verifier with SuperAdmin keys
verifier = WhitelistedAddressVerifier(
    super_admin_keys=super_admin_keys,
    min_valid_signatures=2,
)

# Verify an address envelope
result = verifier.verify_whitelisted_address(
    envelope=envelope,
    rules_container_decoder=rules_container_from_base64,
    user_signatures_decoder=user_signatures_from_base64,
)

# Access decoded rules container from result
rules_container = result.rules_container
verified_hash = result.verified_hash
```

---

## Security Considerations

1. **Never skip verification** - Always verify whitelisted addresses before using them
2. **Use constant-time comparison** - Prevents timing attacks on hash verification
3. **Verify SuperAdmin signatures first** - Ensures rules container is trusted before using
4. **Check hash coverage** - Ensures the specific address hash (not just any hash) was signed
5. **Don't retry integrity errors** - These indicate potential tampering

---

## Related Documentation

- [Authentication](AUTHENTICATION.md) - SuperAdmin key configuration
- [Services Reference](SERVICES.md) - WhitelistedAddressService API
- [Key Concepts](CONCEPTS.md) - Exception handling
- [Common Integrity Verification](../../docs/INTEGRITY_VERIFICATION.md) - Shared verification concepts
