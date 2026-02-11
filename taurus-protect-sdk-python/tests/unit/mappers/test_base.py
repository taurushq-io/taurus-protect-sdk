"""Tests for base mapper utilities."""

from datetime import datetime, timezone

import pytest

from taurus_protect.mappers._base import (
    parse_string_to_int,
    safe_bool,
    safe_datetime,
    safe_float,
    safe_int,
    safe_list,
    safe_string,
)


class TestSafeString:
    """Tests for safe_string function."""

    def test_string_value(self) -> None:
        """Test with actual string."""
        assert safe_string("hello") == "hello"

    def test_none_returns_empty(self) -> None:
        """Test None returns empty string."""
        assert safe_string(None) == ""

    def test_empty_string(self) -> None:
        """Test empty string returns empty."""
        assert safe_string("") == ""

    def test_whitespace_preserved(self) -> None:
        """Test whitespace is preserved."""
        assert safe_string("  hello  ") == "  hello  "

    def test_unicode_string(self) -> None:
        """Test unicode string."""
        assert safe_string("\u4e16\u754c") == "\u4e16\u754c"


class TestSafeBool:
    """Tests for safe_bool function."""

    def test_true_value(self) -> None:
        """Test with True."""
        assert safe_bool(True) is True

    def test_false_value(self) -> None:
        """Test with False."""
        assert safe_bool(False) is False

    def test_none_returns_false(self) -> None:
        """Test None returns False."""
        assert safe_bool(None) is False


class TestSafeInt:
    """Tests for safe_int function."""

    def test_int_value(self) -> None:
        """Test with actual int."""
        assert safe_int(42) == 42

    def test_none_returns_zero(self) -> None:
        """Test None returns 0."""
        assert safe_int(None) == 0

    def test_string_int(self) -> None:
        """Test with string integer."""
        assert safe_int("123") == 123

    def test_invalid_string_returns_zero(self) -> None:
        """Test invalid string returns 0."""
        assert safe_int("not a number") == 0

    def test_negative_int(self) -> None:
        """Test negative integer."""
        assert safe_int(-10) == -10
        assert safe_int("-10") == -10

    def test_float_string_returns_zero(self) -> None:
        """Test float string returns 0 (can't parse)."""
        assert safe_int("3.14") == 0

    def test_zero(self) -> None:
        """Test zero value."""
        assert safe_int(0) == 0
        assert safe_int("0") == 0


class TestSafeFloat:
    """Tests for safe_float function."""

    def test_float_value(self) -> None:
        """Test with actual float."""
        assert safe_float(3.14) == 3.14

    def test_none_returns_zero(self) -> None:
        """Test None returns 0.0."""
        assert safe_float(None) == 0.0

    def test_string_float(self) -> None:
        """Test with string float."""
        assert safe_float("3.14") == 3.14

    def test_string_int(self) -> None:
        """Test with string integer."""
        assert safe_float("42") == 42.0

    def test_invalid_string_returns_zero(self) -> None:
        """Test invalid string returns 0.0."""
        assert safe_float("not a number") == 0.0

    def test_negative_float(self) -> None:
        """Test negative float."""
        assert safe_float(-3.14) == -3.14
        assert safe_float("-3.14") == -3.14

    def test_scientific_notation(self) -> None:
        """Test scientific notation string."""
        assert safe_float("1.5e10") == 1.5e10


class TestSafeDatetime:
    """Tests for safe_datetime function."""

    def test_datetime_value(self) -> None:
        """Test with actual datetime."""
        dt = datetime(2024, 1, 15, 10, 30, 0, tzinfo=timezone.utc)
        assert safe_datetime(dt) == dt

    def test_none_returns_none(self) -> None:
        """Test None returns None."""
        assert safe_datetime(None) is None

    def test_iso_string_with_z(self) -> None:
        """Test ISO string with Z suffix."""
        result = safe_datetime("2024-01-15T10:30:00Z")
        assert result is not None
        assert result.year == 2024
        assert result.month == 1
        assert result.day == 15
        assert result.hour == 10
        assert result.minute == 30

    def test_iso_string_with_offset(self) -> None:
        """Test ISO string with timezone offset."""
        result = safe_datetime("2024-01-15T10:30:00+00:00")
        assert result is not None
        assert result.year == 2024

    def test_invalid_string_returns_none(self) -> None:
        """Test invalid string returns None."""
        assert safe_datetime("not a date") is None

    def test_empty_string_returns_none(self) -> None:
        """Test empty string returns None."""
        assert safe_datetime("") is None

    def test_partial_date_string(self) -> None:
        """Test partial date string."""
        # Behavior depends on implementation - may or may not parse
        result = safe_datetime("2024-01-15")
        # Either None or parsed date is acceptable
        if result is not None:
            assert result.year == 2024


class TestSafeList:
    """Tests for safe_list function."""

    def test_list_value(self) -> None:
        """Test with actual list."""
        assert safe_list([1, 2, 3]) == [1, 2, 3]

    def test_none_returns_empty(self) -> None:
        """Test None returns empty list."""
        assert safe_list(None) == []

    def test_empty_list(self) -> None:
        """Test empty list returns empty."""
        assert safe_list([]) == []

    def test_mixed_types(self) -> None:
        """Test list with mixed types."""
        result = safe_list(["a", 1, True, None])
        assert result == ["a", 1, True, None]

    def test_nested_list(self) -> None:
        """Test nested list."""
        result = safe_list([[1, 2], [3, 4]])
        assert result == [[1, 2], [3, 4]]


class TestParseStringToInt:
    """Tests for parse_string_to_int function."""

    def test_valid_string(self) -> None:
        """Test valid string integer."""
        assert parse_string_to_int("123") == 123

    def test_none_returns_default(self) -> None:
        """Test None returns default."""
        assert parse_string_to_int(None) == 0
        assert parse_string_to_int(None, default=10) == 10

    def test_invalid_string_returns_default(self) -> None:
        """Test invalid string returns default."""
        assert parse_string_to_int("invalid") == 0
        assert parse_string_to_int("invalid", default=99) == 99

    def test_empty_string_returns_default(self) -> None:
        """Test empty string returns default."""
        assert parse_string_to_int("") == 0

    def test_negative_string(self) -> None:
        """Test negative string."""
        assert parse_string_to_int("-42") == -42

    def test_custom_default(self) -> None:
        """Test custom default value."""
        assert parse_string_to_int(None, default=50) == 50
        assert parse_string_to_int("abc", default=50) == 50
