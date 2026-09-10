"""
Strict base64 decoding for security-critical paths.

Lenient base64 decoding -- Python's default -- silently DISCARDS characters outside the
base64 alphabet and absorbs the rest. That makes the decoding non-injective: two
different strings can decode to the same bytes, and, worse, a string carrying embedded
separators decodes to a *valid container with attacker-chosen bytes appended*. Protobuf
treats concatenation as merge, so appended bytes ADD entries to repeated fields such as
``users`` -- which is how an attacker-controlled HSM or PRICEUPDATER key would reach the
governance trust root.

Every governance path that turns a base64 string into bytes must agree byte-for-byte on
what that string means: the verification memo key, the container decode, and the
signature verification that decides what a signature actually covers. Go uses strict
``base64.StdEncoding`` throughout, and that is precisely why Go was not exploitable by
this route.

This module is deliberately a leaf: stdlib only, no SDK imports, so it is safe to import
from both ``mappers`` and ``helpers`` without a cycle. It raises ``ValueError`` and lets
each caller wrap that in its own error type.
"""

from __future__ import annotations

import base64
import re

# ASCII whitespace is stripped before validation. Line-wrapped base64 is legitimate and
# introduces no ambiguity about the decoded bytes; every other out-of-alphabet byte is
# rejected rather than silently dropped.
_WHITESPACE = re.compile(rb"[ \t\n\r\v\f]+")


def strict_b64decode(data: str) -> bytes:
    """
    Decode base64, rejecting any character outside the base64 alphabet.

    Args:
        data: Base64 text. None and the empty string decode to empty bytes.

    Returns:
        The decoded bytes.

    Raises:
        ValueError: If the input is not well-formed base64. ``binascii.Error``
            subclasses ``ValueError``, so a single ``except ValueError`` covers both.
    """
    raw = (data or "").encode("utf-8")
    return base64.b64decode(_WHITESPACE.sub(b"", raw), validate=True)
