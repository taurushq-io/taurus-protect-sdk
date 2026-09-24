"""Base service class for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import timedelta
from typing import TYPE_CHECKING, Any

from taurus_protect.errors import APIError, map_http_error
from taurus_protect.models.pagination import CursorRequest

if TYPE_CHECKING:
    pass  # Import types for type checking only


class BaseService:
    """
    Base class for all service implementations.

    Provides common error handling. Pagination is built by the helpers in
    ``taurus_protect.models.pagination``, never here.
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
            # "code" is the numeric HTTP status; the application code is "error_code".
            error_code = body.get("error_code")

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

    @staticmethod
    def _request_cursor_body(request: CursorRequest) -> Any:
        """The generated ``RequestCursor`` for operations that take the cursor in the body."""
        from taurus_protect._internal.openapi.models.tgvalidatord_request_cursor import (
            TgvalidatordRequestCursor,
        )

        return TgvalidatordRequestCursor(
            current_page=request.current_page,
            page_request=request.page_request,
            page_size=str(request.page_size),
        )

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
