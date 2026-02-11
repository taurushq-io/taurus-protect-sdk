"""Signature verification utilities for governance rules."""

from __future__ import annotations

import base64
import binascii
from typing import List

from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect.crypto.signing import verify_signature
from taurus_protect.errors import IntegrityError
from taurus_protect.models.governance_rules import GovernanceRules


def verify_governance_rules(
    rules: GovernanceRules,
    min_valid_signatures: int,
    super_admin_keys: List[EllipticCurvePublicKey],
) -> None:
    """
    Verify that governance rules have enough valid SuperAdmin signatures.

    This function verifies the cryptographic signatures on the rules container
    to ensure they were signed by the required number of SuperAdmin keys.

    Args:
        rules: The governance rules to verify.
        min_valid_signatures: Minimum number of valid signatures required.
        super_admin_keys: List of SuperAdmin public keys for verification.

    Raises:
        IntegrityError: If verification fails or not enough valid signatures.
        ValueError: If arguments are invalid.

    Example:
        >>> verify_governance_rules(rules, min_valid_signatures=2, super_admin_keys=keys)
    """
    if rules is None:
        raise ValueError("rules cannot be None")
    if min_valid_signatures <= 0:
        raise ValueError("min_valid_signatures must be positive")
    if not super_admin_keys:
        raise ValueError("super_admin_keys cannot be empty")

    if rules.rules_container is None:
        raise IntegrityError("Governance rules verification failed: rulesContainer is null")

    signatures = rules.rules_signatures
    if not signatures:
        raise IntegrityError("Governance rules verification failed: no signatures present")

    # Decode the rules container
    try:
        rules_data = base64.b64decode(rules.rules_container)
    except (binascii.Error, ValueError) as e:
        raise IntegrityError(f"Governance rules verification failed: invalid base64 encoding: {e}") from e

    valid_count = 0
    for sig in signatures:
        if sig.signature and is_valid_signature(rules_data, sig.signature, super_admin_keys):
            valid_count += 1

    if valid_count < min_valid_signatures:
        raise IntegrityError(
            f"Governance rules verification failed: only {valid_count} valid signatures found, "
            f"minimum {min_valid_signatures} required"
        )


def is_valid_signature(
    data: bytes,
    signature_b64: str,
    super_admin_keys: List[EllipticCurvePublicKey],
) -> bool:
    """
    Verify a signature against the provided SuperAdmin public keys.

    Tries each key in sequence and returns True if any key validates
    the signature.

    Args:
        data: The data that was signed.
        signature_b64: Base64-encoded signature.
        super_admin_keys: List of SuperAdmin public keys to try.

    Returns:
        True if the signature is valid for any of the keys.
    """
    for public_key in super_admin_keys:
        try:
            if verify_signature(public_key, data, signature_b64):
                return True
        except (InvalidSignature, ValueError):
            # Signature verification failed for this key, try next
            # InvalidSignature: cryptographic verification failed
            # ValueError: malformed signature or key format
            continue
    return False


def verify_raw_signature(
    data: bytes,
    signature: bytes,
    public_key: EllipticCurvePublicKey,
) -> bool:
    """
    Verify a raw signature against a single public key.

    Args:
        data: The data that was signed.
        signature: The raw signature bytes (not base64).
        public_key: The public key to verify against.

    Returns:
        True if the signature is valid.
    """
    if data is None or signature is None or public_key is None:
        return False

    # Convert raw bytes to base64 for verify_signature
    signature_b64 = base64.b64encode(signature).decode("utf-8")
    return verify_signature(public_key, data, signature_b64)
