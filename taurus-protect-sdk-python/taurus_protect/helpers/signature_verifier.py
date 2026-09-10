"""Signature verification utilities for governance rules."""

from __future__ import annotations

import base64
import binascii
import hashlib
from typing import List, Optional, Sequence

from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect._strict_base64 import strict_b64decode
from taurus_protect.crypto.signing import verify_signature
from taurus_protect.errors import IntegrityError
from taurus_protect.models.governance_rules import GovernanceRules, RuleUserSignature


def verify_governance_rules(
    rules: GovernanceRules,
    min_valid_signatures: int,
    super_admin_keys: List[EllipticCurvePublicKey],
) -> None:
    """
    Verify that governance rules are signed by enough DISTINCT SuperAdmin keys.

    min_valid_signatures counts distinct signing keys, not signature entries: ECDSA
    is randomized, so counting entries would let a single key produce as many valid
    signatures as any threshold demands, making 2-of-N no stronger than 1-of-N.

    Args:
        rules: The governance rules to verify.
        min_valid_signatures: Minimum number of distinct valid signers required.
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
        rules_data = strict_b64decode(rules.rules_container)
    except (binascii.Error, ValueError) as e:
        raise IntegrityError(f"Governance rules verification failed: invalid base64 encoding: {e}") from e

    try:
        verify_governance_rules_signatures(
            rules_data, signatures, super_admin_keys, min_valid_signatures
        )
    except IntegrityError as e:
        raise IntegrityError(f"Governance rules verification failed: {e}") from e


def verify_governance_rules_signatures(
    rules_container_data: bytes,
    signatures: Sequence[RuleUserSignature],
    super_admin_keys: List[EllipticCurvePublicKey],
    min_valid_signatures: int,
) -> None:
    """
    Verify that the rules container is signed by enough DISTINCT SuperAdmin keys.

    This is the single place the threshold is evaluated -- every caller routes here.
    min_valid_signatures counts distinct signing keys, not signature entries: ECDSA
    is randomized, so counting entries would let a single key produce as many valid
    signatures as any threshold demands, making 2-of-N no stronger than 1-of-N.

    Args:
        rules_container_data: The raw signed rules container bytes.
        signatures: The signature entries to check.
        super_admin_keys: List of SuperAdmin public keys for verification.
        min_valid_signatures: Minimum number of distinct valid signers required.

    Raises:
        IntegrityError: If too few distinct SuperAdmin keys signed.
        ValueError: If min_valid_signatures is not positive.
    """
    if min_valid_signatures <= 0:
        raise ValueError("min_valid_signatures must be positive")
    if not rules_container_data:
        raise IntegrityError("rules container data cannot be empty")
    if not super_admin_keys:
        raise IntegrityError("no SuperAdmin keys configured for verification")
    if not signatures:
        raise IntegrityError("no signatures provided")

    signers = set()
    for sig in signatures:
        if not sig.signature:
            continue
        key = _matching_key(rules_container_data, sig.signature, super_admin_keys)
        if key is not None:
            signers.add(key_fingerprint(key))

    if len(signers) < min_valid_signatures:
        raise IntegrityError(
            f"only {len(signers)} distinct valid signers found, "
            f"minimum {min_valid_signatures} required"
        )


def _matching_key(
    data: bytes,
    signature_b64: str,
    super_admin_keys: List[EllipticCurvePublicKey],
) -> Optional[EllipticCurvePublicKey]:
    """Return the first configured key that verifies the signature, else None."""
    for public_key in super_admin_keys:
        try:
            if verify_signature(public_key, data, signature_b64):
                return public_key
        except (InvalidSignature, ValueError):
            # Signature verification failed for this key, try next
            # InvalidSignature: cryptographic verification failed
            # ValueError: malformed signature or key format
            continue
    return None


def key_fingerprint(public_key: EllipticCurvePublicKey) -> bytes:
    """Identify a key by its encoded bytes, so the same key configured twice counts once."""
    der = public_key.public_bytes(
        encoding=serialization.Encoding.DER,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    return hashlib.sha256(der).digest()


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
    return _matching_key(data, signature_b64, super_admin_keys) is not None


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
