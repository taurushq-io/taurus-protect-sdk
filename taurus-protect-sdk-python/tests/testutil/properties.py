"""Simple key=value properties file parser."""

from __future__ import annotations

from typing import Dict


def load_properties(path: str) -> Dict[str, str]:
    """Load a Java-style .properties file.

    Supports:
    - key=value pairs
    - # comments
    - Blank lines (ignored)
    - \\n escape sequences in values (for PEM keys)

    Args:
        path: Path to the properties file.

    Returns:
        Dictionary of key-value pairs.
    """
    props: Dict[str, str] = {}
    try:
        with open(path, "r") as f:
            for line in f:
                line = line.strip()
                if not line or line.startswith("#"):
                    continue
                eq_idx = line.find("=")
                if eq_idx < 0:
                    continue
                key = line[:eq_idx].strip()
                value = line[eq_idx + 1:].strip()
                # Unescape \n in PEM keys
                value = value.replace("\\n", "\n")
                props[key] = value
    except FileNotFoundError:
        pass
    return props
