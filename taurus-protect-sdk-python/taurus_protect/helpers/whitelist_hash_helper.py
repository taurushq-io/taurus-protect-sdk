"""Whitelist hash computation and parsing utilities.

This module provides helpers for computing hashes and parsing
whitelisted address payloads from JSON.

SECURITY DESIGN
===============

The API returns metadata with two representations of the same data:
  - payload: Raw JSON object (UNVERIFIED)
  - payload_as_string: JSON string that is cryptographically hashed (VERIFIED)

The security model works as follows:
  1. The server computes: metadata.hash = SHA256(payload_as_string)
  2. The hash is signed by governance rules (SuperAdmin keys)
  3. Clients verify: computed_hash(payload_as_string) == metadata.hash

ATTACK VECTOR (if using raw payload):
An attacker intercepting API responses could:
  1. Modify the payload object (e.g., change destination address)
  2. Leave payload_as_string unchanged (hash still verifies)
  3. Client extracts data from modified payload -> SECURITY BYPASS

SOLUTION:
All address parsing functions in this module (e.g., parse_whitelisted_address_from_json)
expect a JSON string parameter, which should be the verified payload_as_string.
This ensures:
  - All extracted data comes from the cryptographically verified source
  - Any tampering with the raw payload object is ignored
  - The integrity chain: payload_as_string -> hash -> signature is preserved
"""

from __future__ import annotations

import json
import re
from typing import Any, Dict, List, Optional

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.errors import WhitelistError
from taurus_protect.models.whitelisted_address import (
    InternalAddress,
    InternalWallet,
    WhitelistedAddress,
)

# Regular expressions for legacy hash computation.
# These patterns are used to compute alternative hashes for backward compatibility.
_CONTRACT_TYPE_PATTERN = re.compile(r",\"contractType\":\"[^\"]*\"")
_LABEL_IN_OBJECT_PATTERN = re.compile(r",\"label\":\"[^\"]*\"}")

# Patterns for asset-specific legacy hash computation (for WhitelistedAsset).
# isNFT field patterns - matches "isNFT":(true|false) with optional comma
_IS_NFT_PATTERN_LEADING_COMMA = re.compile(r",\"isNFT\":(true|false)")
_IS_NFT_PATTERN_TRAILING_COMMA = re.compile(r"\"isNFT\":(true|false),")
# kindType field patterns - matches "kindType":"..." with optional comma
_KIND_TYPE_PATTERN_LEADING_COMMA = re.compile(r",\"kindType\":\"[^\"]*\"")
_KIND_TYPE_PATTERN_TRAILING_COMMA = re.compile(r"\"kindType\":\"[^\"]*\",")


def compute_whitelist_hash(
    address: str,
    blockchain: str,
    network: Optional[str],
    label: str,
    memo: Optional[str] = None,
) -> str:
    """
    Compute the SHA-256 hash of a whitelisted address payload.

    This computes the hash in the same format used by the Taurus-PROTECT API
    for whitelisted address verification.

    Args:
        address: The blockchain address string.
        blockchain: The blockchain/currency identifier (e.g., "ETH", "BTC").
        network: Optional network name (e.g., "mainnet", "testnet").
        label: Human-readable label for the address.
        memo: Optional memo/destination tag.

    Returns:
        Hex-encoded SHA-256 hash of the payload.

    Example:
        >>> hash_value = compute_whitelist_hash(
        ...     address="0x1234...",
        ...     blockchain="ETH",
        ...     network="mainnet",
        ...     label="My Wallet"
        ... )
    """
    # Build the payload in the same format as the API
    payload: Dict[str, Any] = {
        "address": address,
        "currency": blockchain,
        "label": label,
    }

    if network:
        payload["network"] = network
    if memo:
        payload["memo"] = memo

    # JSON encode with sorted keys and no spaces (compact format)
    payload_str = json.dumps(payload, separators=(",", ":"), sort_keys=True)

    return calculate_hex_hash(payload_str)


def compute_legacy_hashes(payload_as_string: str) -> List[str]:
    """
    Compute alternative hashes for backward compatibility.

    This handles addresses signed before schema changes by removing certain
    fields and recomputing the hash.

    Strategies:
    1. Remove contractType field (addresses signed before contractType was added)
    2. Remove labels from linkedInternalAddresses (after contractType but before
       labels were added)
    3. Remove both contractType and labels (before both fields were added)

    Args:
        payload_as_string: The original payload JSON string.

    Returns:
        List of unique legacy hashes (may be empty if no transformations apply).
    """
    if not payload_as_string:
        return []

    seen: set[str] = set()
    hashes: List[str] = []

    def add_hash(payload: str) -> None:
        hash_value = calculate_hex_hash(payload)
        if hash_value not in seen:
            seen.add(hash_value)
            hashes.append(hash_value)

    # Strategy 1: Remove contractType only
    # Handles addresses signed before contractType was added to schema
    without_contract_type = _CONTRACT_TYPE_PATTERN.sub("", payload_as_string)
    if without_contract_type != payload_as_string:
        add_hash(without_contract_type)

    # Strategy 2: Remove labels from linkedInternalAddresses objects only
    # (keep contractType)
    # Handles addresses signed after contractType was added but before labels
    without_labels = _LABEL_IN_OBJECT_PATTERN.sub("}", payload_as_string)
    if without_labels != payload_as_string:
        add_hash(without_labels)

    # Strategy 3: Remove BOTH contractType AND labels from linkedInternalAddresses
    # Handles addresses signed before both fields were added
    without_both = _LABEL_IN_OBJECT_PATTERN.sub("}", payload_as_string)
    without_both = _CONTRACT_TYPE_PATTERN.sub("", without_both)
    if without_both != payload_as_string:
        add_hash(without_both)

    return hashes


def compute_asset_legacy_hashes(payload_as_string: str) -> List[str]:
    """
    Compute alternative hashes for backward compatibility with assets.

    This handles assets signed before schema changes by removing certain
    fields and recomputing the hash.

    Strategies (aligned with Java SDK WhitelistedAssetService.computeLegacyHashes):
    1. Remove isNFT field (assets signed before isNFT was added)
    2. Remove kindType field (assets signed before kindType was added)
    3. Remove both isNFT and kindType (assets signed before both fields were added)

    Args:
        payload_as_string: The original payload JSON string.

    Returns:
        List of unique legacy hashes (may be empty if no transformations apply).
    """
    if not payload_as_string:
        return []

    seen: set[str] = set()
    hashes: List[str] = []

    def add_hash(payload: str) -> None:
        hash_value = calculate_hex_hash(payload)
        if hash_value not in seen:
            seen.add(hash_value)
            hashes.append(hash_value)

    # Strategy 1: Remove isNFT only
    # Handles assets signed before isNFT was added to schema
    without_is_nft = _IS_NFT_PATTERN_LEADING_COMMA.sub("", payload_as_string)
    without_is_nft = _IS_NFT_PATTERN_TRAILING_COMMA.sub("", without_is_nft)
    if without_is_nft != payload_as_string:
        add_hash(without_is_nft)

    # Strategy 2: Remove kindType only
    # Handles assets signed before kindType was added to schema
    without_kind_type = _KIND_TYPE_PATTERN_LEADING_COMMA.sub("", payload_as_string)
    without_kind_type = _KIND_TYPE_PATTERN_TRAILING_COMMA.sub("", without_kind_type)
    if without_kind_type != payload_as_string:
        add_hash(without_kind_type)

    # Strategy 3: Remove BOTH isNFT AND kindType
    # Handles assets signed before both fields were added
    # Note: Order matches Java implementation - remove isNFT first, then kindType
    without_both = _IS_NFT_PATTERN_LEADING_COMMA.sub("", payload_as_string)
    without_both = _IS_NFT_PATTERN_TRAILING_COMMA.sub("", without_both)
    without_both = _KIND_TYPE_PATTERN_LEADING_COMMA.sub("", without_both)
    without_both = _KIND_TYPE_PATTERN_TRAILING_COMMA.sub("", without_both)
    if without_both != payload_as_string:
        add_hash(without_both)

    return hashes


def parse_whitelisted_address_from_json(json_str: str) -> WhitelistedAddress:
    """
    Parse a JSON string into a WhitelistedAddress model.

    This parses the verified JSON payload from a signed whitelist envelope
    and extracts the whitelisted address fields.

    SECURITY NOTE:
        This function expects the verified ``payload_as_string`` from the
        metadata, NOT the raw ``payload`` object. The ``payload_as_string``
        is the cryptographically verified source (its hash is signed by
        governance rules). Using the raw payload object would bypass
        integrity verification and could allow an attacker to inject
        tampered data.

    Args:
        json_str: JSON string from metadata.payload_as_string (verified source).

    Returns:
        WhitelistedAddress model populated from the JSON fields.

    Raises:
        WhitelistError: If parsing fails or JSON is invalid.

    Example:
        >>> # CORRECT: Use payload_as_string
        >>> addr = parse_whitelisted_address_from_json(metadata.payload_as_string)
        >>> print(addr.address)
        0x123
    """
    if not json_str:
        raise WhitelistError("JSON payload cannot be null or empty")

    try:
        obj = json.loads(json_str)
    except json.JSONDecodeError as e:
        raise WhitelistError(f"Failed to parse JSON: {e}")

    try:
        # Extract basic fields
        # Note: JSON uses "currency" but model uses "currency" (maps to blockchain)
        linked_internal = _parse_linked_internal_addresses(obj.get("linkedInternalAddresses", []))
        linked_wallets = _parse_linked_wallets(obj.get("linkedWallets", []))

        return WhitelistedAddress(
            id=str(obj.get("id", "")),
            address=_get_string_or_none(obj, "address"),
            label=_get_string_or_none(obj, "label"),
            currency=_get_string_or_none(obj, "currency"),
            network=_get_string_or_none(obj, "network"),
            contract_type=_get_string_or_none(obj, "contractType"),
            memo=_get_string_or_none(obj, "memo"),
            customer_id=_get_string_or_none(obj, "customerId"),
            address_type=_get_string_or_none(obj, "addressType"),
            tn_participant_id=_get_string_or_none(obj, "tnParticipantId"),
            exchange_account_id=_get_string_or_none(obj, "exchangeAccountId"),
            linked_internal_addresses=linked_internal,
            linked_wallets=linked_wallets,
        )
    except (KeyError, TypeError, ValueError) as e:
        raise WhitelistError(f"Failed to parse WhitelistedAddress from JSON: {e}") from e


def _get_string_or_none(obj: Dict[str, Any], key: str) -> Optional[str]:
    """Get a string value from dict, returning None for empty strings."""
    value = obj.get(key)
    if value is None or value == "":
        return None
    return str(value)


def _parse_linked_internal_addresses(arr: List[Dict[str, Any]]) -> List[InternalAddress]:
    """Parse linkedInternalAddresses array into list of InternalAddress objects."""
    result: List[InternalAddress] = []
    for item in arr:
        if isinstance(item, dict):
            result.append(
                InternalAddress(
                    id=str(item.get("id", "")) if item.get("id") is not None else None,
                    address=item.get("address"),
                    label=item.get("label"),
                )
            )
    return result


def _parse_linked_wallets(arr: List[Dict[str, Any]]) -> List[InternalWallet]:
    """Parse linkedWallets array into list of InternalWallet objects."""
    result: List[InternalWallet] = []
    for item in arr:
        if isinstance(item, dict):
            result.append(
                InternalWallet(
                    id=int(item.get("id", 0)),
                    path=item.get("path"),
                    label=item.get("label"),
                )
            )
    return result
