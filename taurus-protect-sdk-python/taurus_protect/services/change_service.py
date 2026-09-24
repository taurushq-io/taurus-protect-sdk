"""Change service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers.audit import (
    change_from_dto,
    changes_from_dto,
    create_change_request_to_dto,
)
from taurus_protect.models.audit import (
    Change,
    ChangeResult,
    CreateChangeRequest,
    ListChangesOptions,
)
from taurus_protect.models.pagination import cursor_page
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class ChangeService(BaseService):
    """Service for managing configuration changes.

    Changes represent modifications to system configuration that require
    approval before taking effect (dual-admin workflow).
    """

    def __init__(self, api_client: Any, changes_api: Any) -> None:
        super().__init__(api_client)
        self._changes_api = changes_api

    def create_change(self, request: CreateChangeRequest) -> str:
        """Create a change request.

        Args:
            request: The change request with action, entity, changes, etc.

        Returns:
            The created change ID.

        Raises:
            ValueError: If action or entity are empty.
            APIError: If API request fails.
        """
        if request is None:
            raise ValueError("request cannot be None")
        if not request.action or not request.action.strip():
            raise ValueError("action cannot be empty")
        if not request.entity or not request.entity.strip():
            raise ValueError("entity cannot be empty")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_create_change_request import (
                TgvalidatordCreateChangeRequest,
            )

            dto_dict = create_change_request_to_dto(request)
            body = TgvalidatordCreateChangeRequest(**dto_dict)
            reply = self._changes_api.change_service_create_change(body)
            result = getattr(reply, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError("createChange returned no result")
            return getattr(result, "id", "")
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, change_id: str) -> Change:
        """Get a change by ID.

        Args:
            change_id: The change ID to retrieve.

        Returns:
            The change.

        Raises:
            ValueError: If change_id is invalid.
            NotFoundError: If change not found.
            APIError: If API request fails.
        """
        self._validate_required(change_id, "change_id")

        try:
            resp = self._changes_api.change_service_get_change(change_id)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Change {change_id} not found")

            change = change_from_dto(result)
            if change is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Change {change_id} not found")

            return change
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        options: Optional[ListChangesOptions] = None,
    ) -> ChangeResult:
        """List changes, one page at a time.

        Args:
            options: Filters and page window; ``options.cursor`` continues a walk.

        Returns:
            ChangeResult with the changes and the page.

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If API request fails.
        """
        opts = options or ListChangesOptions()
        req = opts.to_cursor_request()

        try:
            resp = self._changes_api.change_service_get_changes(
                entity=opts.entity,
                entity_id=opts.entity_id,
                status=opts.status,
                creator_id=opts.creator_id,
                sort_order=opts.sort_order,
                entity_ids=opts.entity_ids,
                entity_uuids=opts.entity_uuids,
                **req.query_params(),
            )

            return ChangeResult(
                changes=changes_from_dto(resp.result),
                page=cursor_page(req.page_size, resp.cursor),
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_for_approval(
        self,
        options: Optional[ListChangesOptions] = None,
    ) -> ChangeResult:
        """List changes pending approval, one page at a time.

        Args:
            options: Filters and page window. ``entity`` filters the queue by entity;
                ``status``, ``creator_id`` and ``entity_id`` do not apply to it and are
                refused.

        Returns:
            ChangeResult with the changes and the page.

        Raises:
            ValueError: If an option does not apply to the approval queue, the page size
                is invalid, or cursor options conflict.
            APIError: If API request fails.
        """
        opts = options or ListChangesOptions()
        for name in ("status", "creator_id", "entity_id"):
            if getattr(opts, name) is not None:
                raise ValueError(f"{name} cannot be applied to the change approval queue")
        req = opts.to_cursor_request()

        try:
            resp = self._changes_api.change_service_get_changes_for_approval(
                entities=[opts.entity] if opts.entity else None,
                sort_order=opts.sort_order,
                entity_ids=opts.entity_ids,
                entity_uuids=opts.entity_uuids,
                **req.query_params(),
            )

            return ChangeResult(
                changes=changes_from_dto(resp.result),
                page=cursor_page(req.page_size, resp.cursor),
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def approve_change(self, change_id: str) -> None:
        """Approve a change.

        Args:
            change_id: The change ID to approve.

        Raises:
            ValueError: If change_id is invalid.
            APIError: If API request fails.
        """
        self._validate_required(change_id, "change_id")

        try:
            self._changes_api.change_service_approve_change(change_id, body={})
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def approve_changes(self, change_ids: List[str]) -> None:
        """Approve multiple changes.

        Args:
            change_ids: List of change IDs to approve.

        Raises:
            ValueError: If change_ids is empty or invalid.
            APIError: If API request fails.
        """
        if not change_ids:
            raise ValueError("change_ids cannot be empty")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_approve_changes_request import (
                TgvalidatordApproveChangesRequest,
            )

            body = TgvalidatordApproveChangesRequest(ids=change_ids)
            self._changes_api.change_service_approve_changes(body)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def reject_change(self, change_id: str) -> None:
        """Reject a change.

        Args:
            change_id: The change ID to reject.

        Raises:
            ValueError: If change_id is invalid.
            APIError: If API request fails.
        """
        self._validate_required(change_id, "change_id")

        try:
            self._changes_api.change_service_reject_change(change_id, body={})
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def reject_changes(self, change_ids: List[str]) -> None:
        """Reject multiple changes.

        Args:
            change_ids: List of change IDs to reject.

        Raises:
            ValueError: If change_ids is empty or invalid.
            APIError: If API request fails.
        """
        if not change_ids:
            raise ValueError("change_ids cannot be empty")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_reject_changes_request import (
                TgvalidatordRejectChangesRequest,
            )

            body = TgvalidatordRejectChangesRequest(ids=change_ids)
            self._changes_api.change_service_reject_changes(body)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
