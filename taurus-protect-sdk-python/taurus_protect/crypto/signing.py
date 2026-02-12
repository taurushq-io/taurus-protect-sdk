"""ECDSA signing utilities for Taurus-PROTECT SDK."""

from __future__ import annotations

import base64

from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.asymmetric.ec import (
    SECP256R1,
    EllipticCurvePrivateKey,
    EllipticCurvePublicKey,
)
from cryptography.hazmat.primitives.asymmetric.utils import (
    decode_dss_signature,
    encode_dss_signature,
)


def _validate_p256_curve(public_key: EllipticCurvePublicKey) -> None:
    """
    Validate that the public key uses the P-256 (secp256r1) curve.

    Args:
        public_key: The ECDSA public key to validate.

    Raises:
        ValueError: If the key does not use P-256 curve.
    """
    if not isinstance(public_key.curve, SECP256R1):
        raise ValueError(
            f"Only P-256 (secp256r1) curve is supported, got {type(public_key.curve).__name__}"
        )


def sign_data(private_key: EllipticCurvePrivateKey, data: bytes) -> str:
    """
    Sign data using ECDSA with SHA-256 (plain format).

    Returns signature in raw r||s format (base64-encoded), matching
    the Java SDK's SHA256withPLAIN-ECDSA output.

    Args:
        private_key: ECDSA private key.
        data: Data to sign.

    Returns:
        Base64-encoded raw r||s signature.

    Example:
        >>> signature = sign_data(private_key, b"request hash data")
        >>> # signature is base64-encoded, e.g., "MEUCIQDw..."
    """
    # Sign with ECDSA - this returns DER-encoded signature
    der_signature = private_key.sign(data, ec.ECDSA(hashes.SHA256()))

    # Convert DER to raw r||s format (matching Java's PLAIN-ECDSA)
    r, s = decode_dss_signature(der_signature)

    # Get key size in bytes (e.g., 32 for P-256)
    key_size = (private_key.curve.key_size + 7) // 8

    # Convert r and s to fixed-size bytes
    r_bytes = r.to_bytes(key_size, byteorder="big")
    s_bytes = s.to_bytes(key_size, byteorder="big")

    # Concatenate r||s and base64 encode
    raw_signature = r_bytes + s_bytes
    return base64.b64encode(raw_signature).decode("utf-8")


def verify_signature(
    public_key: EllipticCurvePublicKey,
    data: bytes,
    signature_b64: str,
) -> bool:
    """
    Verify an ECDSA signature.

    Expects signature in raw r||s format (base64-encoded), matching
    the Java SDK's SHA256withPLAIN-ECDSA format. Only P-256 (secp256r1)
    curve is supported.

    Args:
        public_key: ECDSA public key (must be P-256/secp256r1).
        data: The signed data.
        signature_b64: Base64-encoded raw r||s signature.

    Returns:
        True if signature is valid, False otherwise.

    Raises:
        ValueError: If the public key does not use P-256 curve.

    Example:
        >>> valid = verify_signature(public_key, b"data", "MEUCIQDw...")
        >>> if not valid:
        ...     raise IntegrityError("Signature verification failed")
    """
    # Validate curve type before verification
    _validate_p256_curve(public_key)

    try:
        # Decode base64 signature
        sig_bytes = base64.b64decode(signature_b64)

        # Get key size in bytes
        key_size = (public_key.curve.key_size + 7) // 8

        # Signature should be exactly 2 * key_size bytes (r||s)
        if len(sig_bytes) != 2 * key_size:
            return False

        # Extract r and s
        r = int.from_bytes(sig_bytes[:key_size], byteorder="big")
        s = int.from_bytes(sig_bytes[key_size:], byteorder="big")

        # Convert to DER format for verification
        der_signature = encode_dss_signature(r, s)

        # Verify
        public_key.verify(der_signature, data, ec.ECDSA(hashes.SHA256()))
        return True

    except (ValueError, TypeError):
        # ValueError: Invalid signature format, base64 decode error
        # TypeError: Invalid data types
        return False
    except InvalidSignature:
        # Explicit cryptography InvalidSignature - signature verification failed
        return False


def get_public_key_from_private(
    private_key: EllipticCurvePrivateKey,
) -> EllipticCurvePublicKey:
    """
    Extract public key from private key.

    Args:
        private_key: ECDSA private key.

    Returns:
        Corresponding ECDSA public key.
    """
    return private_key.public_key()
