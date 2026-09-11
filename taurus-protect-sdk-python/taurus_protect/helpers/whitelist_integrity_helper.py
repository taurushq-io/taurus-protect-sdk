"""Whitelist integrity verification utilities.

This module provides helpers for verifying the integrity of signed
whitelist envelopes and extracting verified addresses.
"""

from __future__ import annotations

from typing import Optional

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.errors import IntegrityError, WhitelistError
from taurus_protect.helpers.constant_time import constant_time_compare
from taurus_protect.helpers.whitelist_hash_helper import (
    parse_whitelisted_address_from_json,
)
from taurus_protect.models.whitelisted_address import (
    SignedWhitelistedAddressEnvelope,
    WhitelistedAddress,
)


def verify_whitelist_envelope(envelope: SignedWhitelistedAddressEnvelope) -> bool:
    """
    Verify the integrity of a signed whitelist envelope.

    This verifies that the hash in the metadata matches the computed hash
    of the payload. This is a basic integrity check and does NOT verify
    cryptographic signatures.

    For full verification including signature validation, use the
    WhitelistedAddressVerifier class.

    Args:
        envelope: The signed whitelist envelope to verify.

    Returns:
        True if the envelope hash is valid.

    Raises:
        IntegrityError: If the envelope is invalid or hash verification fails.

    Example:
        >>> if verify_whitelist_envelope(envelope):
        ...     address = extract_whitelisted_address_from_envelope(envelope)
    """
    if envelope is None:
        raise IntegrityError("envelope cannot be None")

    if envelope.metadata is None:
        raise IntegrityError("envelope metadata cannot be None")

    payload_as_string = envelope.metadata.payload_as_string
    if not payload_as_string:
        raise IntegrityError("payloadAsString is empty")

    provided_hash = envelope.metadata.hash
    if not provided_hash:
        raise IntegrityError("metadata hash is empty")

    # Compute hash of payload.
    #
    # There is deliberately NO legacy-variant tolerance here. This used to accept
    # `metadata.hash` when it equalled the hash of a STRIPPED variant -- i.e. it did not
    # require sha256(payload_as_string) == metadata.hash at all -- and then
    # extract_whitelisted_address_from_envelope parsed the UNSTRIPPED payload below.
    # That is the exact inverse of the invariant: it admitted a payload whose own hash
    # nothing had committed to, and returned every field of it as fact.
    #
    # The legacy tolerance belongs in step 4 (is metadata.hash covered by a signature),
    # where it is paired with the payload the matching variant carries -- see
    # WhitelistedAddressVerifier._verify_hash_in_signed_hashes and
    # helpers.whitelist_hash_helper.LegacyPayloadVariant. Step 1 is only ever
    # "does the delivered payload hash to the hash the response claims for it".
    computed_hash = calculate_hex_hash(payload_as_string)
    if constant_time_compare(computed_hash, provided_hash):
        return True

    raise IntegrityError("metadata hash verification failed")


def extract_whitelisted_address_from_envelope(
    envelope: SignedWhitelistedAddressEnvelope,
) -> WhitelistedAddress:
    """
    Extract and parse the whitelisted address from a signed envelope.

    Checks the metadata hash only — step 1 of six. It has no SuperAdmin keys and no
    rules container, so it CANNOT check the container signatures, hash coverage or the
    governance thresholds. A hash match proves the payload was not altered without also
    updating its hash; it does not prove the entry was approved.

    ``WhitelistedAddressService`` is the only route that runs all six steps, and the
    supported way to obtain a whitelisted address.

    The former ``verify`` parameter is gone: passing False returned a parsed address with
    nothing checked at all.

    Args:
        envelope: The signed whitelist envelope containing the address.

    Returns:
        The WhitelistedAddress parsed from the envelope payload.

    Raises:
        IntegrityError: If the metadata hash does not match the payload.
        WhitelistError: If parsing the address fails.
    """
    if envelope is None:
        raise IntegrityError("envelope cannot be None")

    verify_whitelist_envelope(envelope)

    if envelope.metadata is None:
        raise IntegrityError("envelope metadata cannot be None")

    payload_as_string = envelope.metadata.payload_as_string
    if not payload_as_string:
        raise IntegrityError("payloadAsString is empty")

    return parse_whitelisted_address_from_json(payload_as_string)


def verify_envelope_field_match(
    db_address: WhitelistedAddress,
    envelope: SignedWhitelistedAddressEnvelope,
) -> None:
    """
    Verify that the whitelisted address envelope fields match the database fields.

    This validates that critical fields in the envelope match what's stored
    in the database, ensuring the data hasn't been tampered with.

    Args:
        db_address: The address from the database.
        envelope: The signed envelope to validate.

    Raises:
        IntegrityError: If validation fails or fields don't match.

    Example:
        >>> verify_envelope_field_match(db_address, envelope)
    """
    if db_address is None:
        raise IntegrityError("database address cannot be None")
    if envelope is None:
        raise IntegrityError("envelope cannot be None")

    # Extract address from envelope (with verification)
    envelope_address = extract_whitelisted_address_from_envelope(envelope)

    # Validate fields match
    _validate_field_match("Address", db_address.address, envelope_address.address)
    _validate_field_match("Label", db_address.label, envelope_address.label)
    _validate_field_match("Currency", db_address.currency, envelope_address.currency)
    _validate_field_match("Network", db_address.network, envelope_address.network)
    _validate_field_match("ContractType", db_address.contract_type, envelope_address.contract_type)


def _validate_field_match(
    field_name: str,
    db_value: Optional[str],
    envelope_value: Optional[str],
) -> None:
    """
    Validate that two field values match.

    Args:
        field_name: The name of the field being validated.
        db_value: The value from the database.
        envelope_value: The value from the envelope.

    Raises:
        IntegrityError: If the values don't match.
    """
    # Treat None and empty string as equivalent
    db_normalized = db_value if db_value else None
    envelope_normalized = envelope_value if envelope_value else None

    if db_normalized != envelope_normalized:
        raise IntegrityError(
            f"invalid whitelist signature: field {field_name} mismatch "
            f"(db='{db_value}', envelope='{envelope_value}')"
        )
