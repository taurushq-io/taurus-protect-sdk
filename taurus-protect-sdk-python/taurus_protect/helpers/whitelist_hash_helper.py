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

import hmac
import json
import re
from dataclasses import dataclass
from typing import Any, Dict, List, Optional, Tuple

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.errors import IntegrityError, WhitelistError
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


@dataclass(frozen=True)
class LegacyPayloadVariant:
    """
    One backward-compatible rewrite of a signed payload: the exact byte string a
    pre-schema-change signer covered, together with its hash.

    Step 4 must carry the PAYLOAD forward, not just the hash. The strips below are not
    injective, so a response-controlling server can append a member the strip removes --
    a duplicate ``,"label":"X"`` immediately before the closing brace -- to a genuinely
    signed payload. The residue is then the signed bytes exactly, every signature check
    passes, and a step 6 that parsed the DELIVERED payload would return the appended
    value as verified (``json.loads`` keeps the LAST of two duplicate keys). Parsing the
    MATCHED VARIANT is what makes step 6's contract true: every field came from bytes a
    counted signature covered.

    The strips are global, and that does NOT bound the exposure the way it first
    appears. A row whose DELIVERED payload carries inner ``linkedInternalAddresses``
    labels is still exposed, because validatord rebuilds those labels on every read from
    live DB relations rather than from the signed envelope -- so for a row signed before
    per-object labels existed, removing every label (the inner ones the server added AND
    the one the attacker appended) lands exactly on the signed bytes. Both injectable
    members reach the caller on any legacy row: ``label`` at either level, and
    ``contractType``.

    What does bound it is the regex alphabet: ``[^"]*`` cannot contain a quote, so
    nothing beyond those two string values can be smuggled in.

    Attributes:
        hash: The hex SHA-256 of ``payload``.
        payload: The rewritten payload whose hash a signer may have covered. This, not
            the delivered payload, is what step 6 must parse when this variant is the
            match.
    """

    hash: str
    payload: str


def compute_legacy_payload_variants(payload_as_string: str) -> List[LegacyPayloadVariant]:
    """
    Compute the backward-compatible payload rewrites for an address, each paired with
    its hash.

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
        List of unique variants, in strategy order (may be empty if no
        transformation changes the payload).
    """
    if not payload_as_string:
        return []

    seen: set = set()
    variants: List[LegacyPayloadVariant] = []

    def add_variant(payload: str) -> None:
        hash_value = calculate_hex_hash(payload)
        if hash_value not in seen:
            seen.add(hash_value)
            variants.append(LegacyPayloadVariant(hash=hash_value, payload=payload))

    # Strategy 1: Remove contractType only
    # Handles addresses signed before contractType was added to schema
    without_contract_type = _CONTRACT_TYPE_PATTERN.sub("", payload_as_string)
    if without_contract_type != payload_as_string:
        add_variant(without_contract_type)

    # Strategy 2: Remove labels from linkedInternalAddresses objects only
    # (keep contractType)
    # Handles addresses signed after contractType was added but before labels
    without_labels = _LABEL_IN_OBJECT_PATTERN.sub("}", payload_as_string)
    if without_labels != payload_as_string:
        add_variant(without_labels)

    # Strategy 3: Remove BOTH contractType AND labels from linkedInternalAddresses
    # Handles addresses signed before both fields were added
    without_both = _LABEL_IN_OBJECT_PATTERN.sub("}", payload_as_string)
    without_both = _CONTRACT_TYPE_PATTERN.sub("", without_both)
    if without_both != payload_as_string:
        add_variant(without_both)

    return variants


def compute_legacy_hashes(payload_as_string: str) -> List[str]:
    """
    Compute alternative hashes for backward compatibility, discarding the payload each
    one came from.

    Verification must use :func:`compute_legacy_payload_variants` instead: step 6 needs
    the payload, not just the hash. This form remains because the cross-SDK vector
    oracle (``docs/test-vectors/crypto-test-vectors.json``) asserts hashes through this
    name, and because it is part of the public helper surface.

    Args:
        payload_as_string: The original payload JSON string.

    Returns:
        List of unique legacy hashes (may be empty if no transformations apply).
    """
    return [variant.hash for variant in compute_legacy_payload_variants(payload_as_string)]


def compute_asset_legacy_payload_variants(
    payload_as_string: str,
) -> List[LegacyPayloadVariant]:
    """
    Asset peer of :func:`compute_legacy_payload_variants`.

    This handles assets signed before schema changes by removing certain
    fields and recomputing the hash.

    Strategies (aligned with Java SDK WhitelistedAssetService.computeLegacyHashes):
    1. Remove isNFT field (assets signed before isNFT was added)
    2. Remove kindType field (assets signed before kindType was added)
    3. Remove both isNFT and kindType (assets signed before both fields were added)

    No asset identity field is currently injectable through it -- the stripped members
    (``isNFT``, ``kindType``) are not read by the asset payload parser, so for a variant
    to hash-match, the inserted text has to be exactly what the regexes remove. This
    carries the payload anyway, so the two flows stay symmetric and a future schema
    change that makes a stripped field readable does not silently reopen the address
    defect on the asset side.

    Args:
        payload_as_string: The original payload JSON string.

    Returns:
        List of unique variants, in strategy order.
    """
    if not payload_as_string:
        return []

    seen: set = set()
    variants: List[LegacyPayloadVariant] = []

    def add_variant(payload: str) -> None:
        hash_value = calculate_hex_hash(payload)
        if hash_value not in seen:
            seen.add(hash_value)
            variants.append(LegacyPayloadVariant(hash=hash_value, payload=payload))

    # Strategy 1: Remove isNFT only
    # Handles assets signed before isNFT was added to schema
    without_is_nft = _IS_NFT_PATTERN_LEADING_COMMA.sub("", payload_as_string)
    without_is_nft = _IS_NFT_PATTERN_TRAILING_COMMA.sub("", without_is_nft)
    if without_is_nft != payload_as_string:
        add_variant(without_is_nft)

    # Strategy 2: Remove kindType only
    # Handles assets signed before kindType was added to schema
    without_kind_type = _KIND_TYPE_PATTERN_LEADING_COMMA.sub("", payload_as_string)
    without_kind_type = _KIND_TYPE_PATTERN_TRAILING_COMMA.sub("", without_kind_type)
    if without_kind_type != payload_as_string:
        add_variant(without_kind_type)

    # Strategy 3: Remove BOTH isNFT AND kindType
    # Handles assets signed before both fields were added
    # Note: Order matches Java implementation - remove isNFT first, then kindType
    without_both = _IS_NFT_PATTERN_LEADING_COMMA.sub("", payload_as_string)
    without_both = _IS_NFT_PATTERN_TRAILING_COMMA.sub("", without_both)
    without_both = _KIND_TYPE_PATTERN_LEADING_COMMA.sub("", without_both)
    without_both = _KIND_TYPE_PATTERN_TRAILING_COMMA.sub("", without_both)
    if without_both != payload_as_string:
        add_variant(without_both)

    return variants


def compute_asset_legacy_hashes(payload_as_string: str) -> List[str]:
    """
    Compute alternative hashes for backward compatibility with assets, discarding the
    payload each one came from.

    See :func:`compute_asset_legacy_payload_variants`: verification uses that form,
    because step 6 needs the payload rather than only its hash. This form remains for
    the cross-SDK vector oracle.

    Args:
        payload_as_string: The original payload JSON string.

    Returns:
        List of unique legacy hashes (may be empty if no transformations apply).
    """
    return [
        variant.hash for variant in compute_asset_legacy_payload_variants(payload_as_string)
    ]


class _DuplicateJSONKey(Exception):
    """Raised by the object_pairs_hook below; never escapes this module."""

    def __init__(self, key: str) -> None:
        super().__init__(key)
        self.key = key


def _reject_duplicate_object_keys(pairs: List[Tuple[str, Any]]) -> Dict[str, Any]:
    """
    ``object_pairs_hook`` that refuses any object carrying the same key twice.

    ``json.loads`` keeps the LAST of two duplicate keys and reports no error, which is a
    verification bypass on this path rather than a curiosity. The legacy-hash tolerance
    strips a member the parser would still read, so a server can append
    ``,"label":"X"`` before the closing brace of a genuinely signed payload: the strip
    recovers the signed bytes, every signature check passes, and the parse then returns
    the attacker's value. Parsing the matched variant (see
    :class:`LegacyPayloadVariant`) closes that, and this closes the shapes the strip
    does not reach.

    The check is PER OBJECT -- the hook is invoked once for each object as it is
    decoded -- so the same key appearing in sibling objects is perfectly legal, which
    matters because ``linkedInternalAddresses`` is an array of objects that all carry
    ``id``/``address``/``label``.
    """
    seen: set = set()
    for key, _value in pairs:
        if key in seen:
            raise _DuplicateJSONKey(key)
        seen.add(key)
    return dict(pairs)


def _loads_strict(json_str: str, what: str) -> Any:
    """
    Decode a signed payload, refusing duplicate object keys.

    Args:
        json_str: The payload to decode.
        what: What is being parsed, for the error message.

    Returns:
        The decoded JSON value.

    Raises:
        IntegrityError: If an object carries a duplicate key, or nesting is deep
            enough to exhaust the interpreter's recursion budget -- a server-chosen
            depth must not become a crash.
        WhitelistError: If the JSON is malformed.
    """
    try:
        return json.loads(json_str, object_pairs_hook=_reject_duplicate_object_keys)
    except _DuplicateJSONKey as exc:
        raise IntegrityError(
            f"cannot parse {what}: signed payload carries duplicate key "
            f"{exc.key!r}; refusing to choose between two values for one field"
        ) from exc
    except RecursionError as exc:
        raise IntegrityError(
            f"cannot parse {what}: signed payload nests too deeply"
        ) from exc
    except json.JSONDecodeError as exc:
        raise WhitelistError(f"Failed to parse JSON: {exc}")


def parse_whitelisted_address_from_json(json_str: str) -> WhitelistedAddress:
    """
    Parse a JSON string into a WhitelistedAddress model.

    This parses the verified JSON payload from a signed whitelist envelope
    and extracts the whitelisted address fields.

    SECURITY NOTE:
        This function expects the payload a counted signature actually COVERED --
        ``AddressVerificationResult.verified_payload``, which is the matched legacy
        variant when step 4 fell back to one, and otherwise the delivered
        ``payload_as_string``. It is never the raw ``payload`` object: that can be
        tampered with while ``payload_as_string`` still hashes to ``metadata.hash``.

    Args:
        json_str: The verified payload JSON string.

    Returns:
        WhitelistedAddress model populated from the JSON fields.

    Raises:
        WhitelistError: If parsing fails or JSON is invalid.
        IntegrityError: If the payload exceeds ``MAX_PAYLOAD_BYTES`` or carries a
            duplicate object key.

    Example:
        >>> # CORRECT: parse what verification cleared
        >>> result = verifier.verify_whitelisted_address(envelope, ...)
        >>> addr = parse_whitelisted_address_from_json(result.verified_payload)
        >>> print(addr.address)
        0x123
    """
    if not json_str:
        raise WhitelistError("JSON payload cannot be null or empty")

    # Bounded here, not left to the caller's ordering: this is exported, and the
    # payload is hash-checked in step 1 but not AUTHENTICATED until step 5.
    if len(json_str) > MAX_PAYLOAD_BYTES:
        raise IntegrityError(
            f"cannot parse whitelisted address: payload exceeds {MAX_PAYLOAD_BYTES} bytes"
        )

    obj = _loads_strict(json_str, "whitelisted address")

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


def parse_whitelisted_asset_identity_from_json(json_str: str) -> Dict[str, Any]:
    """
    Parse the identity fields of a whitelisted asset out of a verified JSON payload.

    This is step 6 of the asset flow: steps 1-5 prove the envelope is authentic, and
    this is what stops an unsigned value reaching the caller. It returns a MAPPING
    rather than a ``WhitelistedAsset`` because this SDK merged asset and envelope into
    one model, so the service applies these fields onto the envelope it already built
    (``asset.model_copy(update=...)``) instead of constructing a second object.

    The payload it must be given is the payload a counted signature COVERED --
    ``AssetVerificationResult.verified_payload`` -- not the delivered
    ``payload_as_string``. See :class:`LegacyPayloadVariant` for why the two can differ.

    Args:
        json_str: The verified payload JSON string.

    Returns:
        Mapping of ``WhitelistedAsset`` field names to their verified values. Every
        key is always present, so applying it clears any field the signed payload
        does not carry rather than leaving a stale value behind.

    Raises:
        WhitelistError: If the JSON is malformed or is not an object.
        IntegrityError: If the payload exceeds ``MAX_PAYLOAD_BYTES`` or carries a
            duplicate object key.
    """
    if not json_str:
        raise WhitelistError("JSON payload cannot be null or empty")

    # See parse_whitelisted_address_from_json: bounded here, not at the call sites.
    if len(json_str) > MAX_PAYLOAD_BYTES:
        raise IntegrityError(
            f"cannot parse whitelisted asset: payload exceeds {MAX_PAYLOAD_BYTES} bytes"
        )

    payload = _loads_strict(json_str, "whitelisted asset")
    if not isinstance(payload, dict):
        raise WhitelistError("Failed to parse WhitelistedAsset from JSON: not an object")

    # Field names match what the DTO-era mapper read, including the capitalised
    # tolerances -- the same signed bytes are parsed by the Java (AssetHashHelper) and
    # TypeScript (WhitelistedAssetPayload) readers, so the three must agree.
    return {
        "name": payload.get("name"),
        "symbol": payload.get("symbol"),
        "blockchain": payload.get("blockchain") or payload.get("Blockchain"),
        "network": payload.get("network") or payload.get("Network"),
        "contract_address": payload.get("contract_address") or payload.get("contractAddress"),
        "decimals": payload.get("decimals"),
        "token_id": payload.get("token_id") or payload.get("tokenId"),
    }


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


# Maximum size of a signed payload before it is parsed.
#
# The payload is hash-checked in step 1 but not AUTHENTICATED until step 5's
# signatures verify, so anything parsed in between is still attacker-influenced.
# A generous ceiling that only stops a hostile response consuming memory.
MAX_PAYLOAD_BYTES = 1 << 20


def resolve_rule_key_with_source(
    payload_as_string: Optional[str],
    dto_blockchain: Optional[str],
    dto_network: Optional[str],
) -> Tuple[str, str, bool]:
    """
    Return the (blockchain, network) pair that selects the governance rules.

    Taken from the SIGNED payload rather than the surrounding DTO. The DTO is
    free-floating: nothing binds it to the signatures, so a response that set
    blockchain to empty would steer verification to the global-default rule tier,
    which is broader than the rule the entity belongs to. The payload pair is at
    least hash-bound to the material step 5 checks.

    Addresses name the chain ``currency``; assets name it ``blockchain``. Both are
    accepted so one helper serves both flows.

    The chain is always in the payload. The NETWORK is not: governance rules carry a
    per-rule ``includeNetworkInPayload`` flag, and real signed payloads omit
    ``network`` when it is off. Requiring it would reject correctly-signed
    addresses, so the DTO network is used only when the payload has none.

    Args:
        payload_as_string: The signed payload.
        dto_blockchain: The blockchain the response claims, may be None or empty.
        dto_network: The network the response claims, may be None or empty.

    Returns:
        The (blockchain, network) pair to look rules up with.

    Raises:
        IntegrityError: If the payload omits the chain -- never treated as a
            wildcard, because an empty value matches every tier -- or if the
            payload and the DTO disagree on a field the payload does carry.
    """
    if not payload_as_string:
        raise IntegrityError("cannot resolve governance rule key: payload is empty")
    if len(payload_as_string) > MAX_PAYLOAD_BYTES:
        raise IntegrityError(
            f"cannot resolve governance rule key: payload exceeds {MAX_PAYLOAD_BYTES} bytes"
        )

    # Through the STRICT loader, like both parse functions. This is the third place a
    # signed payload is parsed and the one with the widest consequence: `currency` and
    # `network` decide WHICH governance rule judges the row, so last-duplicate-wins here
    # could steer an address to a tier with a weaker quorum. Go, Java and TypeScript all
    # guard this site; this SDK used a bare json.loads.
    try:
        payload = _loads_strict(payload_as_string, "governance rule key")
    except WhitelistError as exc:
        raise IntegrityError(f"cannot resolve governance rule key: {exc}") from exc

    if not isinstance(payload, dict):
        raise IntegrityError("cannot resolve governance rule key: payload is not an object")

    blockchain = payload.get("blockchain") or payload.get("currency") or ""
    network = payload.get("network") or ""

    if not blockchain:
        raise IntegrityError(
            "signed payload does not carry a blockchain; "
            "refusing to fall back to a wildcard rule"
        )

    # Only reachable when include_network_in_payload is off for this rule.
    network_from_payload = bool(network)
    if not network_from_payload:
        network = dto_network or ""

    if dto_blockchain and dto_blockchain.lower() != blockchain.lower():
        raise IntegrityError(
            f"blockchain disagrees between the signed payload ({blockchain}) "
            f"and the response ({dto_blockchain})"
        )
    if network_from_payload and dto_network and dto_network.lower() != network.lower():
        raise IntegrityError(
            f"network disagrees between the signed payload ({network}) "
            f"and the response ({dto_network})"
        )

    return blockchain, network, network_from_payload


def resolve_rule_key(
    payload_as_string: Optional[str],
    dto_blockchain: Optional[str],
    dto_network: Optional[str],
) -> Tuple[str, str]:
    """
    The two-value projection of :func:`resolve_rule_key_with_source`.

    Kept because the shared ``rule_key`` vectors in
    ``scripts/resources/verification-behaviour-vectors.json`` assert exactly this
    shape across all four SDKs -- the same reason Go keeps its own two-value form.
    A caller that needs to know whether the NETWORK was signed (and therefore
    whether one rule tier or every reachable tier must be satisfied) must use the
    source-aware form.

    Args:
        payload_as_string: The signed payload.
        dto_blockchain: The unverified blockchain from the response DTO.
        dto_network: The unverified network from the response DTO.

    Returns:
        The ``(blockchain, network)`` pair that selects the governance rules.
    """
    blockchain, network, _ = resolve_rule_key_with_source(
        payload_as_string, dto_blockchain, dto_network
    )
    return blockchain, network


def contains_hash(hashes: Optional[List[str]], target: str) -> bool:
    """
    Check whether one signature's hashes list covers ``target``.

    The per-signature half of the pair; :func:`verify_hash_coverage` asks the same
    question across every signature. These two are the ONLY places this SDK compares
    hash strings. The verifiers used to carry their own copies, and they had drifted:
    the address copy returned on the first match while the asset copy did not.

    Constant-time, and the loop does not break early.

    Args:
        hashes: The hashes a signature covers, may be None.
        target: The hash to look for.

    Returns:
        True when the list covers the hash.
    """
    if not hashes or not target:
        return False

    found = False
    for h in hashes:
        if h is not None and hmac.compare_digest(target, h):
            found = True
            # No break: returning on a match would leak its position through timing.
    return found


def verify_hash_coverage(metadata_hash: str, signatures) -> bool:
    """
    Check whether the metadata hash is covered by at least one signature.

    Args:
        metadata_hash: The hash to find.
        signatures: Signature entries, each exposing ``hashes``.

    Returns:
        True when at least one signature covers the hash.
    """
    if not metadata_hash or not signatures:
        return False

    found = False
    for sig in signatures:
        if contains_hash(getattr(sig, "hashes", None), metadata_hash):
            found = True
            # No break, as above.
    return found
