"""Base mapper utilities for Taurus-PROTECT SDK."""

from __future__ import annotations

import re
from datetime import datetime
from typing import Any, Optional, Union


def safe_string(value: Optional[str]) -> str:
    """
    Safely convert optional string to string.

    Args:
        value: Optional string value.

    Returns:
        The string value or empty string if None.
    """
    return value if value is not None else ""


def safe_bool(value: Optional[bool]) -> bool:
    """
    Safely convert optional bool to bool.

    Args:
        value: Optional bool value.

    Returns:
        The bool value or False if None.
    """
    return value if value is not None else False


def safe_int(value: Optional[Union[int, str]]) -> int:
    """
    Safely convert optional int or string to int.

    Args:
        value: Optional int or string value.

    Returns:
        The int value or 0 if None or invalid.
    """
    if value is None:
        return 0
    if isinstance(value, int):
        return value
    try:
        return int(value)
    except (ValueError, TypeError):
        return 0


def safe_float(value: Optional[Union[float, str]]) -> float:
    """
    Safely convert optional float or string to float.

    Args:
        value: Optional float or string value.

    Returns:
        The float value or 0.0 if None or invalid.
    """
    if value is None:
        return 0.0
    if isinstance(value, float):
        return value
    try:
        return float(value)
    except (ValueError, TypeError):
        return 0.0


def safe_datetime(value: Optional[Union[datetime, str]]) -> Optional[datetime]:
    """
    Safely convert optional datetime or string to datetime.

    Args:
        value: Optional datetime or ISO string value.

    Returns:
        The datetime value or None if invalid.
    """
    if value is None:
        return None
    if isinstance(value, datetime):
        return value
    try:
        # Handle ISO format strings
        # Replace Z suffix with +00:00
        normalized = value.replace("Z", "+00:00")
        # Handle timezone without colon (e.g., +0000 -> +00:00)
        # Common cases: +0000, -0500, +0530
        normalized = re.sub(r'([+-])(\d{2})(\d{2})$', r'\1\2:\3', normalized)
        return datetime.fromisoformat(normalized)
    except (ValueError, TypeError, AttributeError):
        return None


def safe_list(value: Optional[list[Any]]) -> list[Any]:
    """
    Safely convert optional list to list.

    Args:
        value: Optional list value.

    Returns:
        The list value or empty list if None.
    """
    return value if value is not None else []


def parse_string_to_int(value: Optional[str], default: int = 0) -> int:
    """
    Parse string to int, commonly used for API pagination fields.

    Args:
        value: String value to parse.
        default: Default value if parsing fails.

    Returns:
        Parsed int or default.
    """
    if value is None:
        return default
    try:
        return int(value)
    except (ValueError, TypeError):
        return default
