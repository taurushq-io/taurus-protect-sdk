"""Key handling utilities for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import List, Optional

from cryptography.exceptions import UnsupportedAlgorithm
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric.ec import (
    EllipticCurvePrivateKey,
    EllipticCurvePublicKey,
    SECP256R1,
)


def decode_public_key_pem(pem_data: str) -> EllipticCurvePublicKey:
    """
    Decode a PEM-encoded ECDSA public key.

    Only P-256 (secp256r1) curve is supported for security reasons.

    Args:
        pem_data: PEM-encoded public key string.

    Returns:
        ECDSA public key object (P-256 curve).

    Raises:
        ValueError: If the key cannot be decoded, is not an EC key,
                   or uses an unsupported curve.

    Example:
        >>> key = decode_public_key_pem('''-----BEGIN PUBLIC KEY-----
        ... MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
        ... -----END PUBLIC KEY-----''')
    """
    try:
        key = serialization.load_pem_public_key(pem_data.encode("utf-8"))
        if not isinstance(key, EllipticCurvePublicKey):
            raise ValueError("Key is not an elliptic curve public key")
        # Validate curve is P-256 (secp256r1) - reject weaker curves
        if not isinstance(key.curve, SECP256R1):
            raise ValueError(
                f"Only P-256 (secp256r1) curve is supported, got {key.curve.name}"
            )
        return key
    except ValueError as e:
        # Wrap cryptography library errors to provide consistent messaging
        raise ValueError(f"Failed to decode public key: {e}") from e
    except (TypeError, UnsupportedAlgorithm) as e:
        raise ValueError(f"Failed to decode public key: {e}") from e


def decode_public_keys_pem(pem_keys: List[str]) -> List[EllipticCurvePublicKey]:
    """
    Decode multiple PEM-encoded public keys.

    Args:
        pem_keys: List of PEM-encoded public key strings.

    Returns:
        List of ECDSA public key objects.
    """
    return [decode_public_key_pem(pem) for pem in pem_keys]


def decode_private_key_pem(
    pem_data: str,
    password: Optional[bytes] = None,
) -> EllipticCurvePrivateKey:
    """
    Decode a PEM-encoded ECDSA private key.

    Supports both EC PRIVATE KEY and PKCS#8 formats.
    Only P-256 (secp256r1) curve is supported for security reasons.

    Args:
        pem_data: PEM-encoded private key string.
        password: Optional password for encrypted keys.

    Returns:
        ECDSA private key object (P-256 curve).

    Raises:
        ValueError: If the key cannot be decoded, is not an EC key,
                   or uses an unsupported curve.

    Example:
        >>> key = decode_private_key_pem('''-----BEGIN EC PRIVATE KEY-----
        ... MHQCAQEEIDfN...
        ... -----END EC PRIVATE KEY-----''')
    """
    try:
        key = serialization.load_pem_private_key(
            pem_data.encode("utf-8"),
            password=password,
        )
        if not isinstance(key, EllipticCurvePrivateKey):
            raise ValueError("Key is not an elliptic curve private key")
        # Validate curve is P-256 (secp256r1) - reject weaker curves
        if not isinstance(key.curve, SECP256R1):
            raise ValueError(
                f"Only P-256 (secp256r1) curve is supported, got {key.curve.name}"
            )
        return key
    except ValueError as e:
        # Wrap cryptography library errors to provide consistent messaging
        raise ValueError(f"Failed to decode private key: {e}") from e
    except (TypeError, UnsupportedAlgorithm) as e:
        raise ValueError(f"Failed to decode private key: {e}") from e


def encode_public_key_pem(key: EllipticCurvePublicKey) -> str:
    """
    Encode an ECDSA public key to PEM format.

    Args:
        key: ECDSA public key object.

    Returns:
        PEM-encoded public key string.
    """
    pem_bytes = key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    return pem_bytes.decode("utf-8")
