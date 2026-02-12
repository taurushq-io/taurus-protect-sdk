"""Helper utilities for Taurus-PROTECT SDK."""

from taurus_protect.helpers.address_signature_verifier import verify_address_signature
from taurus_protect.helpers.constant_time import (
    constant_time_compare,
    constant_time_compare_bytes,
)
from taurus_protect.helpers.signature_verifier import (
    is_valid_signature,
    verify_governance_rules,
    verify_raw_signature,
)
from taurus_protect.helpers.whitelist_hash_helper import (
    compute_asset_legacy_hashes,
    compute_legacy_hashes,
    compute_whitelist_hash,
    parse_whitelisted_address_from_json,
)
from taurus_protect.helpers.whitelist_integrity_helper import (
    extract_whitelisted_address_from_envelope,
    verify_envelope_field_match,
    verify_whitelist_envelope,
)
from taurus_protect.helpers.whitelisted_address_verifier import (
    AddressVerificationResult,
    WhitelistedAddressVerifier,
)
from taurus_protect.helpers.whitelisted_asset_verifier import (
    AssetVerificationResult,
    WhitelistedAssetVerifier,
    verify_hash_coverage,
)

__all__ = [
    # Constant time comparison
    "constant_time_compare",
    "constant_time_compare_bytes",
    # Governance rules verification
    "verify_governance_rules",
    "is_valid_signature",
    "verify_raw_signature",
    # Address signature verification
    "verify_address_signature",
    # Whitelisted address verification
    "WhitelistedAddressVerifier",
    "AddressVerificationResult",
    # Whitelisted asset verification
    "WhitelistedAssetVerifier",
    "AssetVerificationResult",
    "verify_hash_coverage",
    # Whitelist hash helpers
    "compute_whitelist_hash",
    "compute_legacy_hashes",
    "compute_asset_legacy_hashes",
    "parse_whitelisted_address_from_json",
    # Whitelist integrity helpers
    "verify_whitelist_envelope",
    "extract_whitelisted_address_from_envelope",
    "verify_envelope_field_match",
]
