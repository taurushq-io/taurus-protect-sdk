"""Constant-time comparison utilities to prevent timing attacks."""

from __future__ import annotations

import hmac


def constant_time_compare(a: str, b: str) -> bool:
    """
    Compare two strings in constant time.

    Uses hmac.compare_digest to prevent timing attacks by ensuring
    the comparison time does not leak information about the content.

    Args:
        a: First string to compare.
        b: Second string to compare.

    Returns:
        True if strings are equal, False otherwise.

    Example:
        >>> constant_time_compare(computed_hash, provided_hash)
        True
    """
    return hmac.compare_digest(a.encode("utf-8"), b.encode("utf-8"))


def constant_time_compare_bytes(a: bytes, b: bytes) -> bool:
    """
    Compare two byte sequences in constant time.

    Uses hmac.compare_digest to prevent timing attacks by ensuring
    the comparison time does not leak information about the content.

    Args:
        a: First bytes to compare.
        b: Second bytes to compare.

    Returns:
        True if bytes are equal, False otherwise.
    """
    return hmac.compare_digest(a, b)
