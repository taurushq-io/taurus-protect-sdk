# Verification & Integrity Helpers -- Python SDK

## Critical Files

- `signature_verifier.py` -- SuperAdmin signature verification (verify_governance_rules, is_valid_signature,
  **verify_governance_rules_signatures** — the shared threshold core)
- `address_signature_verifier.py` -- HSM signature verification for addresses
- `whitelist_hash_helper.py` -- Hash computation for whitelisted address/asset payloads
- `whitelist_integrity_helper.py` -- Integrity verification orchestration
- `whitelisted_address_verifier.py` -- WhitelistedAddressVerifier (6-step verification)
- `whitelisted_asset_verifier.py` -- WhitelistedAssetVerifier (steps 1-5; **step 6 runs in
  `services/whitelisted_asset_service.py`**, which maps the asset from the verified payload)
- `constant_time.py` -- Timing-safe comparison using hmac.compare_digest()

## Whitelisted Address 6-Step Verification

1. Verify metadata hash (SHA-256 of `payload_as_string` == `metadata.hash`)
2. Verify rules container signatures (SuperAdmin keys)
3. Decode rules container (base64 -> protobuf/JSON -> model)
4. Verify hash coverage (`metadata.hash` in signature's hashes list)
5. Verify whitelist signatures per governance thresholds
6. Parse WhitelistedAddress from verified payload

## Whitelisted Asset 5-Step Verification
**Steps 1-5 are the signature checks; step 6 (parse from the VERIFIED payload) is what stops
an unsigned value reaching the caller.** See "Critical Files" above for where it runs here.


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
- `WhitelistedAddressService` list/`List()` verifies every envelope **leniently**: an
  unverifiable row is excluded and reported, not fatal. One bad row used to deny access to
  every good one — and listing is how an operator finds the bad row. Excluding stays
  fail-closed: an omitted destination cannot be selected. **Rows returned but none surviving
  is an error**, so a filtered page never reads as an empty whitelist.
- `_map_address_from_dto()` removed -- all mapping goes through verified envelope path

## Constant-Time Comparison

- Uses `hmac.compare_digest()` for timing-safe comparison
- `_verify_hash_coverage()` and `_contains_hash()` in verifiers use `hmac.compare_digest()`
- Never break/return early when comparing SECRET or hash material — that is what
  `verifyHashCoverage` / `containsHash` are for. It does NOT apply to a threshold
  loop exiting once the required count is reached, nor to trying a signature against
  a set of PUBLIC keys: all four SDKs return early on group-threshold success, and
  that is correct. Do not "fix" it.

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

``python
# WRONG
except Exception: continue

# CORRECT
from cryptography.exceptions import InvalidSignature
except (InvalidSignature, ValueError, binascii.Error): continue
``

### WhitelistedAddress Signature Mapping

The `TgvalidatordWhitelistSignature` has nested structure: `sig_dto.signature` is a nested `TgvalidatordWhitelistUserSignature` object. Access `user_id`, `signature` from the nested object. `hashes` (plural) is a list -- use `hashes[0]` for single hash.

### linkedInternalAddresses Parsing

Payload contains objects like `{'id': ..., 'address': ..., 'label': ...}`. Extract address strings: `if isinstance(item, dict) and item.get("address"): linked.append(str(item["address"]))`.

### Thread-Safe Lazy Init

Use `PrivateAttr(default_factory=threading.Lock)` on Pydantic models with cached state (DecodedRulesContainer, GovernanceRules).

### JSON Separator Compatibility

`json.dumps(hashes, separators=(",", ":"))` -- MUST match Java GSON compact output (no spaces) for signature verification.

## The threshold is evaluated in ONE function — route everything through it

`verify_governance_rules_signatures(rules_container_data, signatures, super_admin_keys, min_valid_signatures)` is the single place `minValidSignatures` is
checked. It counts **distinct signing keys** (SHA-256 of the encoded public key), skips absent
signatures, and rejects a non-positive threshold. Every caller delegates: the governance service and
both whitelisted-address/asset paths.

Do **not** hand-roll a `validCount++` loop over signatures, and do not dedupe on `userId` — it is
server-supplied, and ECDSA is randomized, so one key can emit as many valid signatures as any
threshold demands. That bug shipped once: the hardened helper existed but had zero callers while the
live paths counted entries.

The **per-group** step-5 threshold (`GroupThreshold.minimumSignatures`) is a different
threshold — but it counts distinct signers for the same reason, keyed on the
container-resolved public key and never on the server-supplied `user_id`. Counting entries
let a duplicated entry from one group member satisfy an N-of-M group; two user IDs sharing
one key likewise count once, because that is one compromised secret. Gated cross-SDK by the
`group_threshold` section of `scripts/resources/verification-signed-fixtures.json`
(`tests/unit/helpers/test_signed_fixtures.py`), the same file that gates step 2.

## Local rules for this package

The reasoning behind these lives in the repo-root `CLAUDE.md` → "Verification flows — facts
that cost time to rediscover". What matters when editing *here*:

- **Exactly TWO hash-comparison functions.** `verifyHashCoverage` (across signatures) and
  `containsHash` (within one). Both constant-time, neither breaks early. **Never add a private
  copy to a verifier** — that is how 14 implementations accumulated across the repo and drifted.
  Pinned cross-SDK by `scripts/resources/verification-behaviour-vectors.json`.
- **`resolveRuleKey` keys rule selection on the SIGNED payload.** Chain required (absent is an
  error, never a wildcard); network falls back to the DTO only when the payload has none,
  because `includeNetworkInPayload` may be off. Do not tighten it to require both.
- **`minimumSignatures = 0` on a populated group is rejected.** There is no post-loop threshold
  check, so a zero would mean "one signature suffices". An empty group with zero is fine.
