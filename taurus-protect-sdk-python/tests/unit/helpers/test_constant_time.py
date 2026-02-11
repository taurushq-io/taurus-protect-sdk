"""Tests for constant-time comparison utilities."""

import pytest

from taurus_protect.helpers.constant_time import (
    constant_time_compare,
    constant_time_compare_bytes,
)


class TestConstantTimeCompare:
    """Tests for constant_time_compare function."""

    def test_equal_strings(self) -> None:
        """Test equal strings compare as equal."""
        assert constant_time_compare("test", "test") is True

    def test_unequal_strings(self) -> None:
        """Test unequal strings compare as not equal."""
        assert constant_time_compare("test", "other") is False

    def test_empty_strings(self) -> None:
        """Test empty strings compare as equal."""
        assert constant_time_compare("", "") is True

    def test_different_lengths(self) -> None:
        """Test strings of different lengths compare as not equal."""
        assert constant_time_compare("short", "longer string") is False

    def test_case_sensitive(self) -> None:
        """Test comparison is case sensitive."""
        assert constant_time_compare("Test", "test") is False

    def test_hash_comparison(self) -> None:
        """Test comparing hash-like strings."""
        hash1 = "abc123def456"
        hash2 = "abc123def456"
        hash3 = "abc123def457"  # One char different

        assert constant_time_compare(hash1, hash2) is True
        assert constant_time_compare(hash1, hash3) is False

    def test_unicode_strings(self) -> None:
        """Test unicode string comparison."""
        assert constant_time_compare("hello\u4e16\u754c", "hello\u4e16\u754c") is True
        assert constant_time_compare("hello\u4e16\u754c", "helloworld") is False

    def test_special_characters(self) -> None:
        """Test strings with special characters."""
        assert constant_time_compare("a\nb\tc", "a\nb\tc") is True
        assert constant_time_compare("a\nb\tc", "a b c") is False


class TestConstantTimeCompareBytes:
    """Tests for constant_time_compare_bytes function."""

    def test_equal_bytes(self) -> None:
        """Test equal bytes compare as equal."""
        assert constant_time_compare_bytes(b"test", b"test") is True

    def test_unequal_bytes(self) -> None:
        """Test unequal bytes compare as not equal."""
        assert constant_time_compare_bytes(b"test", b"other") is False

    def test_empty_bytes(self) -> None:
        """Test empty bytes compare as equal."""
        assert constant_time_compare_bytes(b"", b"") is True

    def test_different_lengths(self) -> None:
        """Test bytes of different lengths compare as not equal."""
        assert constant_time_compare_bytes(b"short", b"longer bytes") is False

    def test_binary_data(self) -> None:
        """Test binary data comparison."""
        data1 = bytes([0x00, 0x01, 0x02, 0xFF])
        data2 = bytes([0x00, 0x01, 0x02, 0xFF])
        data3 = bytes([0x00, 0x01, 0x02, 0xFE])

        assert constant_time_compare_bytes(data1, data2) is True
        assert constant_time_compare_bytes(data1, data3) is False

    def test_hash_bytes_comparison(self) -> None:
        """Test comparing hash bytes."""
        import hashlib

        hash1 = hashlib.sha256(b"data").digest()
        hash2 = hashlib.sha256(b"data").digest()
        hash3 = hashlib.sha256(b"other").digest()

        assert constant_time_compare_bytes(hash1, hash2) is True
        assert constant_time_compare_bytes(hash1, hash3) is False

    def test_signature_bytes(self) -> None:
        """Test comparing signature-like bytes."""
        sig1 = b"\x00" * 64
        sig2 = b"\x00" * 64
        sig3 = b"\x00" * 63 + b"\x01"

        assert constant_time_compare_bytes(sig1, sig2) is True
        assert constant_time_compare_bytes(sig1, sig3) is False
