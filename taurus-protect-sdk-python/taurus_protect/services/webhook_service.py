"""Webhook service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.webhook import webhook_from_dto, webhooks_from_dto
from taurus_protect.models.pagination import Pagination
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
        >>> webhooks, pagination = client.webhooks.list(limit=50)
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
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Webhook], Optional[Pagination]]:
        """
        List webhooks with pagination.

        Args:
            limit: Maximum number of webhooks to return (must be positive).
            offset: Number of webhooks to skip (must be non-negative).

        Returns:
            Tuple of (webhooks list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._webhooks_api.webhook_service_get_webhooks(
                type=None,
                url=None,
                cursor_current_page=None,
                cursor_page_request=None,
                cursor_page_size=str(limit),
                sort_order=None,
            )

            webhooks_dto = getattr(resp, "webhooks", None)
            webhooks = webhooks_from_dto(webhooks_dto) if webhooks_dto else []

            # Extract pagination from cursor if available
            cursor = getattr(resp, "cursor", None)
            pagination = None
            if cursor:
                total_items = getattr(cursor, "total_items", None)
                pagination = self._extract_pagination(
                    total_items=total_items,
                    offset=offset,
                    limit=limit,
                )

            return webhooks, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, webhook_id: str) -> Webhook:
        """
        Get a webhook by ID.

        Note: This method lists webhooks and filters by ID since the API
        does not provide a direct get-by-ID endpoint.

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

        try:
            # List webhooks and find the one with matching ID
            resp = self._webhooks_api.webhook_service_get_webhooks(
                type=None,
                url=None,
                cursor_current_page=None,
                cursor_page_request=None,
                cursor_page_size="100",
                sort_order=None,
            )

            webhooks_dto = getattr(resp, "webhooks", None)
            if webhooks_dto:
                for dto in webhooks_dto:
                    if getattr(dto, "id", None) == webhook_id:
                        webhook = webhook_from_dto(dto)
                        if webhook:
                            return webhook

            from taurus_protect.errors import NotFoundError

            raise NotFoundError(f"Webhook {webhook_id} not found")
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

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
