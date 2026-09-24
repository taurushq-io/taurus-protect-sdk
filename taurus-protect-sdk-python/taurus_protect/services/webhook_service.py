"""Webhook service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.webhook import webhook_from_dto, webhooks_from_dto
from taurus_protect.models.pagination import MAX_PAGE_SIZE, CursorPage, cursor_page, cursor_request
from taurus_protect.models.webhook import Webhook
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class WebhookService(BaseService):
    """
    Service for webhook management operations.

    Provides methods to list, get, create, and delete webhooks.
    Webhooks allow you to receive notifications when events occur in Taurus-PROTECT.

    Example:
        >>> # List webhooks
        >>> webhooks, page = client.webhooks.list(page_size=100)
        >>> for webhook in webhooks:
        ...     print(f"{webhook.id}: {webhook.url} ({webhook.status})")
        >>>
        >>> # Get single webhook
        >>> webhook = client.webhooks.get("webhook-123")
        >>> print(f"URL: {webhook.url}")
        >>>
        >>> # Create webhook
        >>> webhook = client.webhooks.create(
        ...     url="https://example.com/webhook",
        ...     events=["REQUEST_CREATED", "REQUEST_APPROVED"],
        ... )
        >>>
        >>> # Delete webhook
        >>> client.webhooks.delete("webhook-123")
    """

    def __init__(self, api_client: Any, webhooks_api: Any) -> None:
        """
        Initialize webhook service.

        Args:
            api_client: The OpenAPI client instance.
            webhooks_api: The WebhooksAPI service from OpenAPI client.
        """
        super().__init__(api_client)
        self._webhooks_api = webhooks_api

    def list(
        self,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        *,
        type: Optional[str] = None,
        url: Optional[str] = None,
        sort_order: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[Webhook], CursorPage]:
        """
        List webhooks, one page at a time.

        Args:
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            type: Filter by webhook type.
            url: Filter by URL.
            sort_order: ASC or DESC.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (webhooks, page).

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If API request fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            resp = self._webhooks_api.webhook_service_get_webhooks(
                type=type,
                url=url,
                sort_order=sort_order,
                **req.query_params(),
            )

            webhooks = webhooks_from_dto(resp.webhooks or [])
            return webhooks, cursor_page(req.page_size, resp.cursor)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, webhook_id: str) -> Webhook:
        """
        Get a webhook by ID.

        The API has no get-by-ID endpoint, so this walks the webhooks page by page
        until the ID is found.

        Args:
            webhook_id: The webhook ID to retrieve.

        Returns:
            The webhook.

        Raises:
            ValueError: If webhook_id is empty.
            NotFoundError: If webhook not found.
            APIError: If API request fails.
        """
        self._validate_required(webhook_id, "webhook_id")

        from taurus_protect.errors import NotFoundError

        cursor: Optional[str] = None
        while True:
            webhooks, page = self.list(page_size=MAX_PAGE_SIZE, cursor=cursor)
            for webhook in webhooks:
                if webhook.id == webhook_id:
                    return webhook
            if not page.has_more:
                raise NotFoundError(f"Webhook {webhook_id} not found")
            cursor = page.next_cursor

    def create(
        self,
        url: str,
        events: List[str],
    ) -> Webhook:
        """
        Create a new webhook.

        Args:
            url: The URL that will receive webhook notifications.
            events: List of event types to subscribe to.

        Returns:
            The created webhook.

        Raises:
            ValueError: If url is empty or events is empty.
            APIError: If API request fails.
        """
        self._validate_required(url, "url")
        if not events:
            raise ValueError("events cannot be empty")

        try:
            # Build request body
            body = {
                "url": url,
                "type": ",".join(events),  # API expects comma-separated event types
            }

            resp = self._webhooks_api.webhook_service_create_webhook(body=body)

            # Response contains the created webhook
            webhook_dto = getattr(resp, "webhook", None)
            if webhook_dto is None:
                # Try to get from result field
                webhook_dto = getattr(resp, "result", None)

            if webhook_dto is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create webhook: no result returned")

            webhook = webhook_from_dto(webhook_dto)
            if webhook is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create webhook: invalid response")

            return webhook
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def delete(self, webhook_id: str) -> None:
        """
        Delete a webhook.

        Args:
            webhook_id: The webhook ID to delete.

        Raises:
            ValueError: If webhook_id is empty.
            NotFoundError: If webhook not found.
            APIError: If API request fails.
        """
        self._validate_required(webhook_id, "webhook_id")

        try:
            self._webhooks_api.webhook_service_delete_webhook(id=webhook_id)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
