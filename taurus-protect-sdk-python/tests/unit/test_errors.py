"""Tests for error classes and error mapping."""

from datetime import timedelta

import pytest

from taurus_protect.errors import (
    APIError,
    AuthenticationError,
    AuthorizationError,
    ConfigurationError,
    IntegrityError,
    NotFoundError,
    RateLimitError,
    ServerError,
    ValidationError,
    WhitelistError,
    map_http_error,
)


class TestAPIError:
    """Tests for base APIError class."""

    def test_init_with_message_only(self) -> None:
        """Test creating error with just message."""
        error = APIError("Something went wrong")
        assert error.message == "Something went wrong"
        assert error.code == 0
        assert error.description is None
        assert error.error_code is None
        assert error.retry_after is None
        assert error.original_error is None

    def test_init_with_all_params(self) -> None:
        """Test creating error with all parameters."""
        original = ValueError("original")
        error = APIError(
            message="Test error",
            code=500,
            description="Server Error",
            error_code="ERR_500",
            retry_after=timedelta(seconds=5),
            original_error=original,
        )
        assert error.message == "Test error"
        assert error.code == 500
        assert error.description == "Server Error"
        assert error.error_code == "ERR_500"
        assert error.retry_after == timedelta(seconds=5)
        assert error.original_error is original

    def test_is_retryable_for_rate_limit(self) -> None:
        """Test is_retryable returns True for 429."""
        error = APIError("Rate limited", code=429)
        assert error.is_retryable() is True

    def test_is_retryable_for_server_errors(self) -> None:
        """Test is_retryable returns True for 5xx."""
        for code in [500, 502, 503, 504]:
            error = APIError("Server error", code=code)
            assert error.is_retryable() is True, f"Expected 429 to be retryable"

    def test_is_retryable_for_client_errors(self) -> None:
        """Test is_retryable returns False for 4xx."""
        for code in [400, 401, 403, 404]:
            error = APIError("Client error", code=code)
            assert error.is_retryable() is False

    def test_is_client_error(self) -> None:
        """Test is_client_error for various codes."""
        assert APIError("", code=400).is_client_error() is True
        assert APIError("", code=404).is_client_error() is True
        assert APIError("", code=499).is_client_error() is True
        assert APIError("", code=500).is_client_error() is False
        assert APIError("", code=200).is_client_error() is False

    def test_is_server_error(self) -> None:
        """Test is_server_error for various codes."""
        assert APIError("", code=500).is_server_error() is True
        assert APIError("", code=503).is_server_error() is True
        assert APIError("", code=400).is_server_error() is False

    def test_suggested_retry_delay_rate_limit(self) -> None:
        """Test suggested_retry_delay for rate limit."""
        error = APIError("", code=429, retry_after=timedelta(seconds=30))
        assert error.suggested_retry_delay() == timedelta(seconds=30)

    def test_suggested_retry_delay_rate_limit_no_header(self) -> None:
        """Test suggested_retry_delay for rate limit without retry_after."""
        error = APIError("", code=429)
        assert error.suggested_retry_delay() == timedelta(seconds=1)

    def test_suggested_retry_delay_server_error(self) -> None:
        """Test suggested_retry_delay for server error."""
        error = APIError("", code=500)
        assert error.suggested_retry_delay() == timedelta(seconds=5)

    def test_suggested_retry_delay_client_error(self) -> None:
        """Test suggested_retry_delay for non-retryable error."""
        error = APIError("", code=400)
        assert error.suggested_retry_delay() == timedelta(0)

    def test_str_representation(self) -> None:
        """Test string representation."""
        error = APIError("Test error", code=500, error_code="ERR_500")
        s = str(error)
        assert "Test error" in s
        assert "500" in s
        assert "ERR_500" in s

    def test_repr_representation(self) -> None:
        """Test repr representation."""
        error = APIError("Test", code=500)
        r = repr(error)
        assert "APIError" in r
        assert "Test" in r


class TestValidationError:
    """Tests for ValidationError class."""

    def test_has_code_400(self) -> None:
        """Test ValidationError has code 400."""
        error = ValidationError("Invalid input")
        assert error.code == 400
        assert error.description == "Bad Request"

    def test_is_not_retryable(self) -> None:
        """Test ValidationError is not retryable."""
        error = ValidationError("Invalid")
        assert error.is_retryable() is False


class TestAuthenticationError:
    """Tests for AuthenticationError class."""

    def test_has_code_401(self) -> None:
        """Test AuthenticationError has code 401."""
        error = AuthenticationError("Invalid credentials")
        assert error.code == 401
        assert error.description == "Unauthorized"


class TestAuthorizationError:
    """Tests for AuthorizationError class."""

    def test_has_code_403(self) -> None:
        """Test AuthorizationError has code 403."""
        error = AuthorizationError("Insufficient permissions")
        assert error.code == 403
        assert error.description == "Forbidden"


class TestNotFoundError:
    """Tests for NotFoundError class."""

    def test_has_code_404(self) -> None:
        """Test NotFoundError has code 404."""
        error = NotFoundError("Resource not found")
        assert error.code == 404
        assert error.description == "Not Found"


class TestRateLimitError:
    """Tests for RateLimitError class."""

    def test_has_code_429(self) -> None:
        """Test RateLimitError has code 429."""
        error = RateLimitError("Too many requests")
        assert error.code == 429
        assert error.description == "Rate Limited"

    def test_has_default_retry_after(self) -> None:
        """Test RateLimitError has default retry_after."""
        error = RateLimitError("Too many requests")
        assert error.retry_after == timedelta(seconds=1)

    def test_custom_retry_after(self) -> None:
        """Test RateLimitError with custom retry_after."""
        error = RateLimitError("Too many requests", retry_after=timedelta(seconds=60))
        assert error.retry_after == timedelta(seconds=60)

    def test_is_retryable(self) -> None:
        """Test RateLimitError is retryable."""
        error = RateLimitError("Too many requests")
        assert error.is_retryable() is True


class TestServerError:
    """Tests for ServerError class."""

    def test_default_code_500(self) -> None:
        """Test ServerError defaults to 500."""
        error = ServerError("Internal error")
        assert error.code == 500

    def test_custom_code(self) -> None:
        """Test ServerError with custom code."""
        error = ServerError("Bad gateway", code=502)
        assert error.code == 502

    def test_is_retryable(self) -> None:
        """Test ServerError is retryable."""
        error = ServerError("Error")
        assert error.is_retryable() is True


class TestIntegrityError:
    """Tests for IntegrityError class."""

    def test_init(self) -> None:
        """Test IntegrityError initialization."""
        error = IntegrityError("Hash mismatch")
        assert error.message == "Hash mismatch"

    def test_str_representation(self) -> None:
        """Test string representation."""
        error = IntegrityError("Hash mismatch")
        assert str(error) == "IntegrityError: Hash mismatch"

    def test_repr_representation(self) -> None:
        """Test repr representation."""
        error = IntegrityError("Hash mismatch")
        assert repr(error) == "IntegrityError(message='Hash mismatch')"

    def test_is_exception(self) -> None:
        """Test IntegrityError is an Exception."""
        error = IntegrityError("Test")
        assert isinstance(error, Exception)
        with pytest.raises(IntegrityError):
            raise error


class TestWhitelistError:
    """Tests for WhitelistError class."""

    def test_init(self) -> None:
        """Test WhitelistError initialization."""
        error = WhitelistError("Invalid whitelist entry")
        assert error.message == "Invalid whitelist entry"

    def test_str_representation(self) -> None:
        """Test string representation."""
        error = WhitelistError("Invalid entry")
        assert str(error) == "WhitelistError: Invalid entry"


class TestConfigurationError:
    """Tests for ConfigurationError class."""

    def test_init(self) -> None:
        """Test ConfigurationError initialization."""
        error = ConfigurationError("Missing API key")
        assert error.message == "Missing API key"


class TestMapHttpError:
    """Tests for map_http_error function."""

    def test_map_400_to_validation_error(self) -> None:
        """Test 400 maps to ValidationError."""
        error = map_http_error(400, "Invalid input")
        assert isinstance(error, ValidationError)
        assert error.code == 400

    def test_map_401_to_authentication_error(self) -> None:
        """Test 401 maps to AuthenticationError."""
        error = map_http_error(401, "Unauthorized")
        assert isinstance(error, AuthenticationError)
        assert error.code == 401

    def test_map_403_to_authorization_error(self) -> None:
        """Test 403 maps to AuthorizationError."""
        error = map_http_error(403, "Forbidden")
        assert isinstance(error, AuthorizationError)
        assert error.code == 403

    def test_map_404_to_not_found_error(self) -> None:
        """Test 404 maps to NotFoundError."""
        error = map_http_error(404, "Not found")
        assert isinstance(error, NotFoundError)
        assert error.code == 404

    def test_map_429_to_rate_limit_error(self) -> None:
        """Test 429 maps to RateLimitError."""
        error = map_http_error(429, "Too many requests", retry_after=timedelta(seconds=30))
        assert isinstance(error, RateLimitError)
        assert error.retry_after == timedelta(seconds=30)

    def test_map_500_to_server_error(self) -> None:
        """Test 500 maps to ServerError."""
        error = map_http_error(500, "Internal error")
        assert isinstance(error, ServerError)
        assert error.code == 500

    def test_map_502_to_server_error(self) -> None:
        """Test 502 maps to ServerError."""
        error = map_http_error(502, "Bad gateway")
        assert isinstance(error, ServerError)
        assert error.code == 502

    def test_map_unknown_to_api_error(self) -> None:
        """Test unknown codes map to base APIError."""
        error = map_http_error(418, "I'm a teapot")
        assert type(error) is APIError
        assert error.code == 418

    def test_map_preserves_error_code(self) -> None:
        """Test that error_code is preserved."""
        error = map_http_error(404, "Not found", error_code="RESOURCE_NOT_FOUND")
        assert error.error_code == "RESOURCE_NOT_FOUND"

    def test_map_preserves_original_error(self) -> None:
        """Test that original_error is preserved."""
        original = ValueError("original")
        error = map_http_error(500, "Error", original_error=original)
        assert error.original_error is original
