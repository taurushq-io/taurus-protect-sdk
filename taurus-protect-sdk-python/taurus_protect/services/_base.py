"""Base service class for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import timedelta
from typing import TYPE_CHECKING, Any, Optional, Tuple, Union

from taurus_protect.errors import APIError, map_http_error
from taurus_protect.models.pagination import Pagination

if TYPE_CHECKING:
    pass  # Import types for type checking only


class BaseService:
    """
    Base class for all service implementations.

    Provides common error handling and pagination utilities.
    """

    def __init__(self, api_client: Any) -> None:
        """
        Initialize base service.

        Args:
            api_client: The OpenAPI client instance.
        """
        self._api_client = api_client

    def _handle_error(self, error: Exception) -> APIError:
        """
        Convert API exceptions to domain errors.

        Args:
            error: The original exception.

        Returns:
            Appropriate APIError subclass.
        """
        # Handle OpenAPI generated exceptions
        status_code = getattr(error, "status", 500)
        message = str(error)

        # Try to extract error details from response body
        body = getattr(error, "body", None)
        error_code = None

        if body and isinstance(body, dict):
            message = body.get("message", message)
            error_code = body.get("code")

        # Check for retry-after header (rate limiting)
        retry_after = None
        headers = getattr(error, "headers", {})
        if headers and "Retry-After" in headers:
            try:
                retry_after = timedelta(seconds=int(headers["Retry-After"]))
            except (ValueError, TypeError):
                pass

        return map_http_error(
            status_code=status_code,
            message=message,
            error_code=error_code,
            retry_after=retry_after,
            original_error=error,
        )

    def _extract_pagination(
        self,
        total_items: Optional[Union[str, int]],
        offset: Optional[Union[str, int]],
        limit: int,
    ) -> Optional[Pagination]:
        """
        Extract pagination info from API response.

        Args:
            total_items: Total items from response (string or int).
            offset: Current offset from response (string or int).
            limit: Requested limit.

        Returns:
            Pagination info or None if not available.
        """
        if total_items is None and offset is None:
            return None

        return Pagination.from_response(total_items, offset, limit)

    @staticmethod
    def _validate_required(value: Any, name: str) -> None:
        """
        Validate that a required value is not empty.

        Args:
            value: The value to validate.
            name: Parameter name for error message.

        Raises:
            ValueError: If value is None or empty string.
        """
        if value is None:
            raise ValueError(f"{name} cannot be None")
        if isinstance(value, str) and not value.strip():
            raise ValueError(f"{name} cannot be empty")
