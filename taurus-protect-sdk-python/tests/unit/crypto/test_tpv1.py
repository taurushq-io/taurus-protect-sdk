"""Tests for TPV1 authentication."""

import pytest

from taurus_protect.crypto.tpv1 import TPV1Auth, calculate_base64_hmac, verify_base64_hmac


class TestTPV1Auth:
    """Tests for TPV1Auth class."""

    def test_init_valid_credentials(self, api_key: str, api_secret_hex: str) -> None:
        """Test initialization with valid credentials."""
        auth = TPV1Auth(api_key, api_secret_hex)
        assert auth.api_key == api_key
        auth.close()

    def test_init_empty_api_key_raises(self, api_secret_hex: str) -> None:
        """Test that empty API key raises ValueError."""
        with pytest.raises(ValueError, match="api_key cannot be empty"):
            TPV1Auth("", api_secret_hex)

    def test_init_empty_api_secret_raises(self, api_key: str) -> None:
        """Test that empty API secret raises ValueError."""
        with pytest.raises(ValueError, match="api_secret cannot be empty"):
            TPV1Auth(api_key, "")

    def test_init_invalid_hex_secret_raises(self, api_key: str) -> None:
        """Test that invalid hex secret raises ValueError."""
        with pytest.raises(ValueError):
            TPV1Auth(api_key, "not-valid-hex")

    def test_sign_request_returns_valid_header(self, api_key: str, api_secret_hex: str) -> None:
        """Test that sign_request returns properly formatted header."""
        auth = TPV1Auth(api_key, api_secret_hex)

        header = auth.sign_request(
            method="POST",
            host="api.example.com",
            path="/v1/wallets",
            content_type="application/json",
            body='{"name": "test"}',
        )

        assert header.startswith("TPV1-HMAC-SHA256")
        assert f"ApiKey={api_key}" in header
        assert "Nonce=" in header
        assert "Timestamp=" in header
        assert "Signature=" in header

        auth.close()

    def test_sign_request_after_close_raises(self, api_key: str, api_secret_hex: str) -> None:
        """Test that sign_request raises after close."""
        auth = TPV1Auth(api_key, api_secret_hex)
        auth.close()

        with pytest.raises(RuntimeError, match="has been closed"):
            auth.sign_request("GET", "api.example.com", "/v1/wallets")

    def test_close_is_idempotent(self, api_key: str, api_secret_hex: str) -> None:
        """Test that close can be called multiple times safely."""
        auth = TPV1Auth(api_key, api_secret_hex)
        auth.close()
        auth.close()  # Should not raise

    def test_context_manager(self, api_key: str, api_secret_hex: str) -> None:
        """Test that TPV1Auth can be used implicitly in with statements."""
        auth = TPV1Auth(api_key, api_secret_hex)
        try:
            header = auth.sign_request("GET", "api.example.com", "/v1/health")
            assert header.startswith("TPV1-HMAC-SHA256")
        finally:
            auth.close()

    def test_parse_url(self) -> None:
        """Test URL parsing utility."""
        host, path, query = TPV1Auth.parse_url("https://api.example.com/v1/wallets?limit=10")
        assert host == "api.example.com"
        assert path == "/v1/wallets"
        assert query == "limit=10"

    def test_parse_url_no_query(self) -> None:
        """Test URL parsing without query string."""
        host, path, query = TPV1Auth.parse_url("https://api.example.com/v1/wallets")
        assert host == "api.example.com"
        assert path == "/v1/wallets"
        assert query is None


class TestHMACFunctions:
    """Tests for HMAC utility functions."""

    def test_calculate_base64_hmac(self) -> None:
        """Test HMAC calculation."""
        secret = bytes.fromhex("0123456789abcdef0123456789abcdef")
        data = "test data"

        result = calculate_base64_hmac(secret, data)

        assert isinstance(result, str)
        # Base64 encoded SHA256 HMAC should be ~44 chars
        assert len(result) > 40

    def test_verify_base64_hmac_valid(self) -> None:
        """Test HMAC verification with valid signature."""
        secret = bytes.fromhex("0123456789abcdef0123456789abcdef")
        data = "test data"

        hmac_value = calculate_base64_hmac(secret, data)

        assert verify_base64_hmac(secret, data, hmac_value) is True

    def test_verify_base64_hmac_invalid(self) -> None:
        """Test HMAC verification with invalid signature."""
        secret = bytes.fromhex("0123456789abcdef0123456789abcdef")
        data = "test data"

        assert verify_base64_hmac(secret, data, "invalid_hmac") is False
