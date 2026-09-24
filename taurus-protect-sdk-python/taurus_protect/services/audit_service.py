"""Audit service for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.audit import audits_from_dto
from taurus_protect.models.audit import Audit
from taurus_protect.models.pagination import MAX_PAGE_SIZE, CursorPage, cursor_page, cursor_request
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class AuditService(BaseService):
    """
    Service for audit event operations.

    Provides methods to list and retrieve audit events.

    Example:
        >>> # List audit events
        >>> audits, page = client.audits.list(page_size=100)
        >>> for audit in audits:
        ...     print(f"{audit.id}: {audit.description}")
        >>>
        >>> # Get single audit event
        >>> audit = client.audits.get("123")
        >>> print(f"Type: {audit.type}")
    """

    def __init__(self, api_client: Any, audit_api: Any) -> None:
        """
        Initialize audit service.

        Args:
            api_client: The OpenAPI client instance.
            audit_api: The AuditApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._audit_api = audit_api

    def list(
        self,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        *,
        external_user_id: Optional[str] = None,
        entities: Optional[List[str]] = None,
        actions: Optional[List[str]] = None,
        from_date: Optional[datetime] = None,
        to_date: Optional[datetime] = None,
        sort_by: Optional[List[str]] = None,
        sort_order: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[Audit], CursorPage]:
        """
        List audit trails, one page at a time.

        Args:
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            external_user_id: Filter by external user ID.
            entities: Filter by entity types.
            actions: Filter by action types.
            from_date: Only trails created at or after this date.
            to_date: Only trails created before this date.
            sort_by: Sort fields.
            sort_order: ASC or DESC.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (audits, page).

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If API request fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            resp = self._audit_api.audit_service_get_audit_trails(
                external_user_id=external_user_id,
                entities=entities,
                actions=actions,
                creation_date_from=from_date,
                creation_date_to=to_date,
                sorting_sort_by=sort_by,
                sorting_sort_order=sort_order,
                **req.query_params(),
            )

            return audits_from_dto(resp.result or []), cursor_page(req.page_size, resp.cursor)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, audit_id: str) -> Audit:
        """
        Get an audit event by ID.

        The API has no get-by-ID endpoint, so this walks the audit trail page by page
        until the ID is found.

        Args:
            audit_id: The audit event ID to retrieve.

        Returns:
            The audit event.

        Raises:
            ValueError: If audit_id is invalid.
            NotFoundError: If audit event not found.
            APIError: If API request fails.
        """
        self._validate_required(audit_id, "audit_id")

        from taurus_protect.errors import NotFoundError

        cursor: Optional[str] = None
        while True:
            audits, page = self.list(page_size=MAX_PAGE_SIZE, cursor=cursor)
            for audit in audits:
                if audit.id == audit_id:
                    return audit
            if not page.has_more:
                raise NotFoundError(f"Audit {audit_id} not found")
            cursor = page.next_cursor

    def export_audit_trails(
        self,
        external_user_id: Optional[str] = None,
        entities: Optional[List[str]] = None,
        actions: Optional[List[str]] = None,
        from_date: Optional[datetime] = None,
        to_date: Optional[datetime] = None,
        format: Optional[str] = None,
    ) -> str:
        """
        Export audit trails in the specified format.

        Args:
            external_user_id: Filter by external user ID.
            entities: Filter by entity types.
            actions: Filter by action types.
            from_date: Filter from date.
            to_date: Filter to date.
            format: Export format, e.g. "csv" or "json".

        Returns:
            The exported data as a string.

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._audit_api.audit_service_export_audit_trails(
                external_user_id=external_user_id,
                entities=entities,
                actions=actions,
                creation_date_from=from_date,
                creation_date_to=to_date,
                format=format,
            )
            return getattr(resp, "result", "") or ""
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
