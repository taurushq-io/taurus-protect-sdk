"""Webhook call service for Taurus-PROTECT SDK."""

from __future__ import annotations

from dataclasses import dataclass
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.webhook import webhook_call_from_dto, webhook_calls_from_dto
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.webhook import WebhookCall
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


@dataclass
class ApiRequestCursor:
    """
    Cursor for paginated API requests.

    Attributes:
        current_page: Current page token for pagination.
        page_request: Page request direction (NEXT or PREVIOUS).
        page_size: Number of items per page.
    """

    current_page: Optional[str] = None
    page_request: Optional[str] = None
    page_size: int = 50


@dataclass
class WebhookCallResult:
    """
    Result from webhook calls list operation.

    Attributes:
        calls: List of webhook calls.
        cursor: Pagination cursor for next page.
        has_more: Whether more results are available.
    """

    calls: List[WebhookCall]
    cursor: Optional[str] = None
    has_more: bool = False


class WebhookCallService(BaseService):
    """
    Service for retrieving webhook call history in the Taurus-PROTECT system.

    This service provides access to the history of webhook invocations,
    including their delivery status and payload information.

    Example:
        >>> # Get all webhook calls
        >>> result = client.webhook_calls.get_webhook_calls()
        >>> for call in result.calls:
        ...     print(f"{call.id}: {call.status}")
        >>>
        >>> # Get calls for a specific webhook
        >>> result = client.webhook_calls.get_webhook_calls(webhook_id="webhook-123")
        >>>
        >>> # Get failed calls only
        >>> result = client.webhook_calls.get_webhook_calls(status="FAILED")
        >>>
        >>> # List webhook calls (alternative method)
        >>> calls, pagination = client.webhook_calls.list(
        ...     webhook_id="webhook-123",
        ...     limit=50,
        ... )
    """

    def __init__(self, api_client: Any, webhook_calls_api: Any) -> None:
        """
        Initialize webhook call service.

        Args:
            api_client: The OpenAPI client instance.
            webhook_calls_api: The WebhookCallsApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._webhook_calls_api = webhook_calls_api

    def get_webhook_calls(
        self,
        event_id: Optional[str] = None,
        webhook_id: Optional[str] = None,
        status: Optional[str] = None,
        sort_order: Optional[str] = None,
        cursor: Optional[ApiRequestCursor] = None,
    ) -> WebhookCallResult:
        """
        Retrieve webhook call history with optional filtering.

        Returns a paginated list of webhook calls that can be filtered by
        event ID, webhook ID, or status.

        Args:
            event_id: Filter by event ID (optional).
            webhook_id: Filter by webhook ID (optional).
            status: Filter by call status (optional, e.g., "SUCCESS", "FAILED").
            sort_order: Sort order for results (optional, "ASC" or "DESC", default "DESC").
            cursor: Pagination cursor (optional, None for first page).

        Returns:
            WebhookCallResult containing the calls and pagination info.

        Raises:
            APIError: If the API call fails.

        Example:
            >>> # Get first page of calls
            >>> result = client.webhook_calls.get_webhook_calls(
            ...     webhook_id="webhook-123",
            ...     status="SUCCESS",
            ... )
            >>> print(f"Found {len(result.calls)} calls")
            >>>
            >>> # Get next page using cursor
            >>> if result.has_more:
            ...     next_cursor = ApiRequestCursor(current_page=result.cursor)
            ...     result = client.webhook_calls.get_webhook_calls(cursor=next_cursor)
        """
        cursor_current_page = None
        cursor_page_request = None
        cursor_page_size = None

        if cursor is not None:
            cursor_current_page = cursor.current_page
            cursor_page_request = cursor.page_request
            cursor_page_size = str(cursor.page_size) if cursor.page_size else None

        try:
            resp = self._webhook_calls_api.webhook_service_get_webhook_calls(
                event_id=event_id,
                webhook_id=webhook_id,
                status=status,
                cursor_current_page=cursor_current_page,
                cursor_page_request=cursor_page_request,
                cursor_page_size=cursor_page_size,
                sort_order=sort_order,
            )

            # Extract calls from response
            calls_dto = getattr(resp, "calls", None) or getattr(resp, "result", None)
            calls = webhook_calls_from_dto(calls_dto) if calls_dto else []

            # Extract cursor for pagination
            cursor_resp = getattr(resp, "cursor", None)
            next_cursor = None
            has_more = False
            if cursor_resp:
                next_cursor = getattr(cursor_resp, "current_page", None)
                has_more = bool(next_cursor)

            return WebhookCallResult(
                calls=calls,
                cursor=next_cursor,
                has_more=has_more,
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        webhook_id: Optional[str] = None,
        event_id: Optional[str] = None,
        status: Optional[str] = None,
        sort_order: Optional[str] = None,
        limit: int = 50,
        cursor: Optional[str] = None,
    ) -> Tuple[List[WebhookCall], Optional[Pagination]]:
        """
        List webhook calls with optional filtering.

        This is an alternative interface to get_webhook_calls() that returns
        a tuple of (calls, pagination) for consistency with other services.

        Args:
            webhook_id: Filter by webhook ID (optional).
            event_id: Filter by event ID (optional).
            status: Filter by call status (optional, e.g., "SUCCESS", "FAILED").
            sort_order: Sort order for results (optional, "ASC" or "DESC").
            limit: Maximum number of calls to return (must be positive).
            cursor: Pagination cursor for next page (optional).

        Returns:
            Tuple of (webhook calls list, pagination info).

        Raises:
            ValueError: If limit is invalid.
            APIError: If API request fails.

        Example:
            >>> # List all webhook calls
            >>> calls, pagination = client.webhook_calls.list(limit=50)
            >>>
            >>> # Filter by webhook ID
            >>> calls, pagination = client.webhook_calls.list(webhook_id="webhook-123")
            >>>
            >>> # Filter by status
            >>> calls, pagination = client.webhook_calls.list(status="FAILED")
        """
        if limit <= 0:
            raise ValueError("limit must be positive")

        try:
            resp = self._webhook_calls_api.webhook_service_get_webhook_calls(
                event_id=event_id,
                webhook_id=webhook_id,
                status=status,
                cursor_current_page=cursor,
                cursor_page_request=None,
                cursor_page_size=str(limit),
                sort_order=sort_order,
            )

            # Extract calls from response
            calls_dto = getattr(resp, "calls", None) or getattr(resp, "result", None)
            calls = webhook_calls_from_dto(calls_dto) if calls_dto else []

            # Extract cursor for pagination
            cursor_resp = getattr(resp, "cursor", None)
            next_cursor = None
            has_more = False
            if cursor_resp:
                next_cursor = getattr(cursor_resp, "current_page", None)
                has_more = bool(next_cursor)

            pagination = Pagination(
                total_items=len(calls),
                offset=0,
                limit=limit,
                has_more=has_more,
            )

            return calls, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, call_id: str) -> WebhookCall:
        """
        Get a webhook call by ID.

        Note: This method lists webhook calls and filters by ID since the API
        does not provide a direct get-by-ID endpoint.

        Args:
            call_id: The webhook call ID to retrieve.

        Returns:
            The webhook call.

        Raises:
            ValueError: If call_id is empty.
            NotFoundError: If webhook call not found.
            APIError: If API request fails.
        """
        self._validate_required(call_id, "call_id")

        try:
            # List all calls and find the one with matching ID
            resp = self._webhook_calls_api.webhook_service_get_webhook_calls(
                event_id=None,
                webhook_id=None,
                status=None,
                cursor_current_page=None,
                cursor_page_request=None,
                cursor_page_size="100",
                sort_order=None,
            )

            calls_dto = getattr(resp, "calls", None)
            if calls_dto:
                for dto in calls_dto:
                    if getattr(dto, "id", None) == call_id:
                        call = webhook_call_from_dto(dto)
                        if call:
                            return call

            from taurus_protect.errors import NotFoundError

            raise NotFoundError(f"Webhook call {call_id} not found")
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e
