"""Hashing utilities for Taurus-PROTECT SDK."""

from __future__ import annotations

import hashlib
import hmac


def calculate_hex_hash(data: str) -> str:
    """
    Calculate SHA-256 hash and return hex-encoded string.

    This is used for request hash verification.

    Args:
        data: The data to hash.

    Returns:
        Hex-encoded SHA-256 hash.

    Example:
        >>> calculate_hex_hash("hello")
        '2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824'
    """
    return hashlib.sha256(data.encode("utf-8")).hexdigest()


def constant_time_compare(a: str, b: str) -> bool:
    """
    Compare two strings in constant time to prevent timing attacks.

    This MUST be used when comparing cryptographic hashes or signatures.

    Args:
        a: First string.
        b: Second string.

    Returns:
        True if strings are equal.

    Example:
        >>> constant_time_compare("abc", "abc")
        True
        >>> constant_time_compare("abc", "def")
        False
    """
    return hmac.compare_digest(a.encode("utf-8"), b.encode("utf-8"))


def calculate_sha256_bytes(data: bytes) -> bytes:
    """
    Calculate SHA-256 hash of bytes.

    Args:
        data: The bytes to hash.

    Returns:
        Raw SHA-256 hash bytes (32 bytes).
    """
    return hashlib.sha256(data).digest()
