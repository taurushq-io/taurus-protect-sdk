"""Webhook call service for Taurus-PROTECT SDK."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import TYPE_CHECKING, Any, List, Optional, Tuple, Union

from taurus_protect.mappers.webhook import webhook_calls_from_dto
from taurus_protect.models.pagination import (
    MAX_PAGE_SIZE,
    CursorPage,
    CursorRequest,
    cursor_page,
    cursor_request,
)
from taurus_protect.models.webhook import WebhookCall
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


@dataclass
class ApiRequestCursor:
    """
    Low-level cursor for :meth:`WebhookCallService.get_webhook_calls`.

    Prefer passing ``page.next_cursor`` as a string; this form is for callers that need
    an explicit page direction.

    Attributes:
        current_page: Page token.
        page_request: Page direction (FIRST, PREVIOUS, NEXT, LAST).
        page_size: Page size (default 20, max 100).
    """

    current_page: Optional[str] = None
    page_request: Optional[str] = None
    page_size: Optional[int] = None


@dataclass
class WebhookCallResult:
    """
    One page of webhook calls.

    Attributes:
        calls: The webhook calls in this page.
        page: The page window; continue with ``page.next_cursor`` while ``has_more``.
    """

    calls: List[WebhookCall]
    page: CursorPage = field(default_factory=CursorPage)


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
        >>> calls, page = client.webhook_calls.list(
        ...     webhook_id="webhook-123",
        ...     page_size=100,
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
        cursor: Optional[Union[str, ApiRequestCursor]] = None,
        *,
        page_size: Optional[int] = None,
    ) -> WebhookCallResult:
        """
        Retrieve webhook call history, one page at a time.

        Args:
            event_id: Filter by event ID (optional).
            webhook_id: Filter by webhook ID (optional).
            status: Filter by call status (optional, e.g., "SUCCESS", "FAILED").
            sort_order: Sort order for results (optional, "ASC" or "DESC", default "DESC").
            cursor: ``result.page.next_cursor`` from the previous page, or an
                :class:`ApiRequestCursor` for an explicit page direction.
            page_size: Page size (default 20, max 100); overrides the page size of an
                :class:`ApiRequestCursor`.

        Returns:
            WebhookCallResult with the calls and the page.

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If the API call fails.

        Example:
            >>> result = client.webhook_calls.get_webhook_calls(webhook_id="webhook-123")
            >>> while result.page.has_more:
            ...     result = client.webhook_calls.get_webhook_calls(
            ...         webhook_id="webhook-123", cursor=result.page.next_cursor
            ...     )
        """
        if isinstance(cursor, ApiRequestCursor):
            req = cursor_request(
                page_size if page_size is not None else cursor.page_size,
                current_page=cursor.current_page,
                page_request=cursor.page_request,
            )
        else:
            req = cursor_request(page_size, cursor)

        calls, page = self._page(
            req, event_id=event_id, webhook_id=webhook_id, status=status, sort_order=sort_order
        )
        return WebhookCallResult(calls=calls, page=page)

    def list(
        self,
        webhook_id: Optional[str] = None,
        event_id: Optional[str] = None,
        status: Optional[str] = None,
        sort_order: Optional[str] = None,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
    ) -> Tuple[List[WebhookCall], CursorPage]:
        """
        List webhook calls, one page at a time.

        Args:
            webhook_id: Filter by webhook ID (optional).
            event_id: Filter by event ID (optional).
            status: Filter by call status (optional, e.g., "SUCCESS", "FAILED").
            sort_order: Sort order for results (optional, "ASC" or "DESC").
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.

        Returns:
            Tuple of (webhook calls, page).

        Raises:
            ValueError: If the page size is invalid.
            APIError: If API request fails.

        Example:
            >>> calls, page = client.webhook_calls.list(webhook_id="webhook-123")
            >>> if page.has_more:
            ...     more, page = client.webhook_calls.list(cursor=page.next_cursor)
        """
        return self._page(
            cursor_request(page_size, cursor),
            event_id=event_id,
            webhook_id=webhook_id,
            status=status,
            sort_order=sort_order,
        )

    def get(self, call_id: str) -> WebhookCall:
        """
        Get a webhook call by ID.

        The API has no get-by-ID endpoint, so this walks the calls page by page until
        the ID is found.

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

        from taurus_protect.errors import NotFoundError

        cursor: Optional[str] = None
        while True:
            calls, page = self.list(page_size=MAX_PAGE_SIZE, cursor=cursor)
            for call in calls:
                if call.id == call_id:
                    return call
            if not page.has_more:
                raise NotFoundError(f"Webhook call {call_id} not found")
            cursor = page.next_cursor

    def _page(self, req: CursorRequest, **filters: Any) -> Tuple[List[WebhookCall], CursorPage]:
        try:
            resp = self._webhook_calls_api.webhook_service_get_webhook_calls(
                **filters, **req.query_params()
            )

            calls = webhook_calls_from_dto(resp.calls or [])
            return calls, cursor_page(req.page_size, resp.cursor)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
