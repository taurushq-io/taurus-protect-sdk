"""Tests for hashing utilities."""

import pytest

from taurus_protect.crypto.hashing import (
    calculate_hex_hash,
    calculate_sha256_bytes,
    constant_time_compare,
)


class TestCalculateHexHash:
    """Tests for calculate_hex_hash function."""

    def test_hash_known_value(self) -> None:
        """Test hash of known string produces expected output."""
        # "hello" SHA-256 hash is well-known
        result = calculate_hex_hash("hello")
        expected = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
        assert result == expected

    def test_hash_empty_string(self) -> None:
        """Test hash of empty string."""
        result = calculate_hex_hash("")
        # SHA-256 of empty string
        expected = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
        assert result == expected

    def test_hash_unicode_string(self) -> None:
        """Test hash of unicode string."""
        result = calculate_hex_hash("hello \u4e16\u754c")  # hello 世界
        assert len(result) == 64  # SHA-256 hex is 64 chars
        assert all(c in "0123456789abcdef" for c in result)

    def test_hash_json_payload(self) -> None:
        """Test hash of JSON-like string (common use case)."""
        payload = '{"amount":"1000","to":"0x1234"}'
        result = calculate_hex_hash(payload)
        assert len(result) == 64
        # Verify deterministic
        assert result == calculate_hex_hash(payload)

    def test_hash_is_lowercase_hex(self) -> None:
        """Test that hash is lowercase hex string."""
        result = calculate_hex_hash("test data")
        assert result == result.lower()
        assert all(c in "0123456789abcdef" for c in result)

    def test_hash_different_inputs_differ(self) -> None:
        """Test that different inputs produce different hashes."""
        hash1 = calculate_hex_hash("input1")
        hash2 = calculate_hex_hash("input2")
        assert hash1 != hash2

    def test_hash_whitespace_sensitive(self) -> None:
        """Test that whitespace matters in hash."""
        hash1 = calculate_hex_hash("hello")
        hash2 = calculate_hex_hash("hello ")
        hash3 = calculate_hex_hash(" hello")
        assert hash1 != hash2
        assert hash1 != hash3
        assert hash2 != hash3


class TestCalculateSha256Bytes:
    """Tests for calculate_sha256_bytes function."""

    def test_hash_bytes_known_value(self) -> None:
        """Test hash of known bytes produces expected output."""
        result = calculate_sha256_bytes(b"hello")
        assert len(result) == 32  # SHA-256 is 32 bytes
        # Verify by comparing hex representation
        assert result.hex() == "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

    def test_hash_empty_bytes(self) -> None:
        """Test hash of empty bytes."""
        result = calculate_sha256_bytes(b"")
        expected_hex = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
        assert result.hex() == expected_hex

    def test_hash_returns_bytes(self) -> None:
        """Test that result is bytes type."""
        result = calculate_sha256_bytes(b"test")
        assert isinstance(result, bytes)
        assert len(result) == 32

    def test_hash_binary_data(self) -> None:
        """Test hash of binary data with non-printable bytes."""
        binary_data = bytes([0x00, 0x01, 0x02, 0xFF, 0xFE])
        result = calculate_sha256_bytes(binary_data)
        assert len(result) == 32
        # Verify deterministic
        assert result == calculate_sha256_bytes(binary_data)


class TestConstantTimeCompare:
    """Tests for constant_time_compare function."""

    def test_equal_strings_return_true(self) -> None:
        """Test that equal strings compare as equal."""
        assert constant_time_compare("abc", "abc") is True
        assert constant_time_compare("", "") is True
        assert constant_time_compare("hello world", "hello world") is True

    def test_unequal_strings_return_false(self) -> None:
        """Test that unequal strings compare as not equal."""
        assert constant_time_compare("abc", "def") is False
        assert constant_time_compare("abc", "abcd") is False
        assert constant_time_compare("abc", "ABC") is False  # Case sensitive

    def test_compare_hashes(self) -> None:
        """Test comparing hash-like strings (common use case)."""
        hash1 = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
        hash2 = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
        hash3 = "3cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

        assert constant_time_compare(hash1, hash2) is True
        assert constant_time_compare(hash1, hash3) is False

    def test_compare_different_lengths(self) -> None:
        """Test comparing strings of different lengths."""
        assert constant_time_compare("short", "longer string") is False
        assert constant_time_compare("a", "aa") is False

    def test_compare_unicode(self) -> None:
        """Test comparing unicode strings."""
        assert constant_time_compare("\u4e16\u754c", "\u4e16\u754c") is True
        assert constant_time_compare("\u4e16\u754c", "world") is False

    def test_compare_special_characters(self) -> None:
        """Test comparing strings with special characters."""
        assert constant_time_compare("a\nb\tc", "a\nb\tc") is True
        assert constant_time_compare("a\nb\tc", "a b c") is False
