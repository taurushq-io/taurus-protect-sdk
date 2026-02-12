"""Exception classes for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import timedelta
from typing import Any, Optional


class APIError(Exception):
    """
    Base exception for all Taurus-PROTECT API errors.

    Attributes:
        message: Human-readable error message.
        code: HTTP status code.
        description: Short error description.
        error_code: Application-specific error code from the API.
        retry_after: Suggested retry delay for rate-limited requests.
        original_error: The underlying exception, if any.
    """

    def __init__(
        self,
        message: str,
        code: int = 0,
        description: Optional[str] = None,
        error_code: Optional[str] = None,
        retry_after: Optional[timedelta] = None,
        original_error: Optional[Exception] = None,
    ) -> None:
        super().__init__(message)
        self.message = message
        self.code = code
        self.description = description
        self.error_code = error_code
        self.retry_after = retry_after
        self.original_error = original_error

    def is_retryable(self) -> bool:
        """
        Check if the request might succeed on retry.

        Returns:
            True for rate limit (429) and server errors (5xx).
        """
        return self.code == 429 or self.code >= 500

    def is_client_error(self) -> bool:
        """
        Check if this is a client error.

        Returns:
            True for 4xx status codes.
        """
        return 400 <= self.code < 500

    def is_server_error(self) -> bool:
        """
        Check if this is a server error.

        Returns:
            True for 5xx status codes.
        """
        return self.code >= 500

    def suggested_retry_delay(self) -> timedelta:
        """
        Get suggested delay before retrying.

        Returns:
            Recommended wait time before retry.
        """
        if self.code == 429:
            return self.retry_after or timedelta(seconds=1)
        if self.code >= 500:
            return timedelta(seconds=5)
        return timedelta(0)

    def __str__(self) -> str:
        parts = [self.message]
        if self.code:
            parts.append(f"(HTTP {self.code})")
        if self.error_code:
            parts.append(f"[{self.error_code}]")
        return " ".join(parts)

    def __repr__(self) -> str:
        return (
            f"{self.__class__.__name__}("
            f"message={self.message!r}, "
            f"code={self.code}, "
            f"error_code={self.error_code!r})"
        )


class ValidationError(APIError):
    """400 Bad Request - Input validation failed."""

    def __init__(
        self,
        message: str,
        error_code: Optional[str] = None,
        original_error: Optional[Exception] = None,
    ) -> None:
        super().__init__(
            message=message,
            code=400,
            description="Bad Request",
            error_code=error_code,
            original_error=original_error,
        )


class AuthenticationError(APIError):
    """401 Unauthorized - Invalid or missing credentials."""

    def __init__(
        self,
        message: str,
        error_code: Optional[str] = None,
        original_error: Optional[Exception] = None,
    ) -> None:
        super().__init__(
            message=message,
            code=401,
            description="Unauthorized",
            error_code=error_code,
            original_error=original_error,
        )


class AuthorizationError(APIError):
    """403 Forbidden - Insufficient permissions."""

    def __init__(
        self,
        message: str,
        error_code: Optional[str] = None,
        original_error: Optional[Exception] = None,
    ) -> None:
        super().__init__(
            message=message,
            code=403,
            description="Forbidden",
            error_code=error_code,
            original_error=original_error,
        )


class NotFoundError(APIError):
    """404 Not Found - Resource does not exist."""

    def __init__(
        self,
        message: str,
        error_code: Optional[str] = None,
        original_error: Optional[Exception] = None,
    ) -> None:
        super().__init__(
            message=message,
            code=404,
            description="Not Found",
            error_code=error_code,
            original_error=original_error,
        )


class RateLimitError(APIError):
    """429 Too Many Requests - Rate limit exceeded."""

    def __init__(
        self,
        message: str,
        retry_after: Optional[timedelta] = None,
        error_code: Optional[str] = None,
        original_error: Optional[Exception] = None,
    ) -> None:
        super().__init__(
            message=message,
            code=429,
            description="Rate Limited",
            error_code=error_code,
            retry_after=retry_after or timedelta(seconds=1),
            original_error=original_error,
        )


class ServerError(APIError):
    """5xx Server Error - Internal server error."""

    def __init__(
        self,
        message: str,
        code: int = 500,
        error_code: Optional[str] = None,
        original_error: Optional[Exception] = None,
    ) -> None:
        super().__init__(
            message=message,
            code=code,
            description="Server Error",
            error_code=error_code,
            original_error=original_error,
        )


class IntegrityError(Exception):
    """
    Cryptographic integrity verification failure.

    This exception indicates that hash or signature verification failed.
    This is a security-critical error and should NEVER be retried.

    Examples:
        - Request hash mismatch
        - Invalid address signature
        - Invalid whitelist signature
        - SuperAdmin signature verification failure
    """

    def __init__(self, message: str) -> None:
        super().__init__(message)
        self.message = message

    def __str__(self) -> str:
        return f"IntegrityError: {self.message}"

    def __repr__(self) -> str:
        return f"IntegrityError(message={self.message!r})"


class WhitelistError(Exception):
    """
    Whitelisted address verification failure.

    This exception indicates that whitelist verification failed,
    which may be due to:
        - Invalid hash in metadata
        - Insufficient governance rule signatures
        - Missing required fields in payload
    """

    def __init__(self, message: str) -> None:
        super().__init__(message)
        self.message = message

    def __str__(self) -> str:
        return f"WhitelistError: {self.message}"

    def __repr__(self) -> str:
        return f"WhitelistError(message={self.message!r})"


class ConfigurationError(Exception):
    """
    SDK configuration error.

    This exception indicates invalid SDK configuration such as:
        - Missing required credentials
        - Invalid SuperAdmin key format
        - Invalid min_valid_signatures value
    """

    def __init__(self, message: str) -> None:
        super().__init__(message)
        self.message = message


class RequestMetadataError(Exception):
    """
    Request metadata error.

    This exception is thrown when request metadata cannot be parsed or extracted.
    Common causes include:
        - Missing required fields in the metadata payload
        - Malformed JSON structure
        - Type mismatches when parsing values
    """

    def __init__(self, message: str, cause: Optional[Exception] = None) -> None:
        super().__init__(message)
        self.message = message
        self.cause = cause

    def __str__(self) -> str:
        return f"RequestMetadataError: {self.message}"

    def __repr__(self) -> str:
        return f"RequestMetadataError(message={self.message!r})"


def map_http_error(
    status_code: int,
    message: str,
    error_code: Optional[str] = None,
    retry_after: Optional[timedelta] = None,
    original_error: Optional[Exception] = None,
) -> APIError:
    """
    Map HTTP status code to appropriate exception type.

    Args:
        status_code: HTTP status code.
        message: Error message.
        error_code: Optional application error code.
        retry_after: Optional retry delay for rate limits.
        original_error: Optional underlying exception.

    Returns:
        Appropriate APIError subclass.
    """
    if status_code == 400:
        return ValidationError(message, error_code, original_error)
    elif status_code == 401:
        return AuthenticationError(message, error_code, original_error)
    elif status_code == 403:
        return AuthorizationError(message, error_code, original_error)
    elif status_code == 404:
        return NotFoundError(message, error_code, original_error)
    elif status_code == 429:
        return RateLimitError(message, retry_after, error_code, original_error)
    elif status_code >= 500:
        return ServerError(message, status_code, error_code, original_error)
    else:
        return APIError(
            message=message,
            code=status_code,
            error_code=error_code,
            original_error=original_error,
        )
