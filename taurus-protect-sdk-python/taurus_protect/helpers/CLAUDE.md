# Verification & Integrity Helpers -- Python SDK

## Critical Files

- `signature_verifier.py` -- SuperAdmin signature verification (verify_governance_rules, is_valid_signature)
- `address_signature_verifier.py` -- HSM signature verification for addresses
- `whitelist_hash_helper.py` -- Hash computation for whitelisted address/asset payloads
- `whitelist_integrity_helper.py` -- Integrity verification orchestration
- `whitelisted_address_verifier.py` -- WhitelistedAddressVerifier (6-step verification)
- `whitelisted_asset_verifier.py` -- WhitelistedAssetVerifier (5-step verification)
- `constant_time.py` -- Timing-safe comparison using hmac.compare_digest()

## Whitelisted Address 6-Step Verification

1. Verify metadata hash (SHA-256 of `payload_as_string` == `metadata.hash`)
2. Verify rules container signatures (SuperAdmin keys)
3. Decode rules container (base64 -> protobuf/JSON -> model)
4. Verify hash coverage (`metadata.hash` in signature's hashes list)
5. Verify whitelist signatures per governance thresholds
6. Parse WhitelistedAddress from verified payload

## Whitelisted Asset 5-Step Verification

Same as address steps 1-5 but uses `ContractAddressWhitelistingRules` with `find_contract_address_whitelisting_rules(blockchain, network)`.
- `parallel_thresholds` is `List[SequentialThresholds]` NOT `List[GroupThreshold]`
- `GroupThreshold.get_min_signatures()` prefers `minimum_signatures` over `threshold` field

## Legacy Hash Computation

Address: remove `contractType`, remove `label` from `linkedInternalAddresses`, remove both.
Asset: remove `isNFT`, remove `kindType`, remove both.
**CRITICAL:** Address and Asset use DIFFERENT legacy hash functions.

## Field Sourcing (CRITICAL SECURITY)

- Security-critical fields (address, label, currency, contract_type, linked_internal_addresses) MUST come from verified payload ONLY
- Non-security fields (id, status, network, tenant_id, created_at, action, rule) can fall back to DTO
- If payload is missing a field, result is None (not DTO value)
- `WhitelistedAddressService.list()` verifies each envelope (strict mode -- fail on first error)
- `_map_address_from_dto()` removed -- all mapping goes through verified envelope path

## Constant-Time Comparison

- Uses `hmac.compare_digest()` for timing-safe comparison
- `_verify_hash_coverage()` and `_contains_hash()` in verifiers use `hmac.compare_digest()`
- Never break/return early in multi-signature loops

## Governance Rules Model (`models/governance_rules.py`)

- `get_hsm_public_key()` -- finds HSMSLOT role user, cached with thread-safe lock
- `find_address_whitelisting_rules(blockchain, network)` -- three-tier priority matching
- `find_contract_address_whitelisting_rules(blockchain, network)` -- same priority for assets

## Rules Container Mapper (`mappers/governance_rules.py`)

- `rules_container_from_base64()` -- protobuf preferred, JSON fallback
- Protobuf Role enum: `Role.Name(role_int)` converts int -> string (e.g., 5 -> "HSMSLOT")
- Protobuf uses camelCase, Python models use snake_case after conversion

## Security-Specific Lessons

### Exception Handling in Crypto Loops

```python
# WRONG
except Exception: continue

# CORRECT
from cryptography.exceptions import InvalidSignature
except (InvalidSignature, ValueError, binascii.Error): continue
```

### WhitelistedAddress Signature Mapping

The `TgvalidatordWhitelistSignature` has nested structure: `sig_dto.signature` is a nested `TgvalidatordWhitelistUserSignature` object. Access `user_id`, `signature` from the nested object. `hashes` (plural) is a list -- use `hashes[0]` for single hash.

### linkedInternalAddresses Parsing

Payload contains objects like `{'id': ..., 'address': ..., 'label': ...}`. Extract address strings: `if isinstance(item, dict) and item.get("address"): linked.append(str(item["address"]))`.

### Thread-Safe Lazy Init

Use `PrivateAttr(default_factory=threading.Lock)` on Pydantic models with cached state (DecodedRulesContainer, GovernanceRules).

### JSON Separator Compatibility

`json.dumps(hashes, separators=(",", ":"))` -- MUST match Java GSON compact output (no spaces) for signature verification.
