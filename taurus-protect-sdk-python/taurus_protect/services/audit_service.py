"""Audit service for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.audit import audit_from_dto, audits_from_dto
from taurus_protect.models.audit import Audit
from taurus_protect.models.pagination import Pagination
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class AuditService(BaseService):
    """
    Service for audit event operations.

    Provides methods to list and retrieve audit events.

    Example:
        >>> # List audit events
        >>> audits, pagination = client.audits.list(limit=50, offset=0)
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
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Audit], Optional[Pagination]]:
        """
        List audit events with pagination.

        Args:
            limit: Maximum number of audit events to return (must be positive).
            offset: Number of audit events to skip (must be non-negative).

        Returns:
            Tuple of (audits list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._audit_api.audit_service_get_audit_trails(
                cursor_page_size=str(limit),
            )

            result = getattr(resp, "result", None)
            audits = audits_from_dto(result) if result else []

            # Extract pagination from cursor
            cursor = getattr(resp, "cursor", None)
            total = getattr(cursor, "total_items", None) if cursor else None
            pagination = self._extract_pagination(
                total_items=total,
                offset=offset,
                limit=limit,
            )

            return audits, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, audit_id: str) -> Audit:
        """
        Get an audit event by ID.

        Note: The underlying API does not have a direct get-by-ID endpoint.
        This method fetches the audit list and filters by ID.

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

        try:
            # The API doesn't have a direct get-by-ID endpoint
            # Fetch from the list and find the matching audit
            audits, _ = self.list(limit=1000)
            for audit in audits:
                if audit.id == audit_id:
                    return audit

            from taurus_protect.errors import NotFoundError

            raise NotFoundError(f"Audit {audit_id} not found")
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

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
