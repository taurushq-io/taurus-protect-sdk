"""Action service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.action import action_from_dto, actions_from_dto
from taurus_protect.models.action import Action
from taurus_protect.models.pagination import (
    PLUS_ROWS,
    Pagination,
    offset_pagination,
    offset_query,
    resolve_offset,
    resolve_page_size,
)
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class ActionService(BaseService):
    """
    Service for action management operations.

    Provides methods to list and get actions. Actions represent
    pending or completed operations that may require approval.

    Example:
        >>> # List actions
        >>> actions, pagination = client.actions.list(limit=100, offset=0)
        >>> for action in actions:
        ...     print(f"{action.label}: {action.status}")
        >>>
        >>> # Get single action
        >>> action = client.actions.get("action-123")
        >>> print(f"Status: {action.status}")
    """

    def __init__(self, api_client: Any, actions_api: Any) -> None:
        """
        Initialize action service.

        Args:
            api_client: The OpenAPI client instance.
            actions_api: The ActionsApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._actions_api = actions_api

    def get(self, action_id: str) -> Action:
        """
        Get an action by ID.

        Args:
            action_id: The action ID to retrieve.

        Returns:
            The action.

        Raises:
            ValueError: If action_id is invalid.
            NotFoundError: If action not found.
            APIError: If API request fails.
        """
        self._validate_required(action_id, "action_id")

        try:
            resp = self._actions_api.action_service_get_action(action_id)

            result = getattr(resp, "action", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Action {action_id} not found")

            action = action_from_dto(result)
            if action is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Action {action_id} not found")

            return action
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
    ) -> Tuple[List[Action], Pagination]:
        """
        List actions, one page at a time.

        Args:
            limit: Page size (default 20, max 100).
            offset: Number of actions to skip; pass ``pagination.next_offset`` to continue.

        Returns:
            Tuple of (actions, pagination).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        page_size = resolve_page_size(limit, "limit")
        start = resolve_offset(offset)

        try:
            resp = self._actions_api.action_service_get_actions(**offset_query(page_size, start))

            rows = resp.result or []
            actions = actions_from_dto(rows)
            pagination = offset_pagination(
                PLUS_ROWS,
                limit=page_size,
                offset=start,
                served_rows=len(rows),
                total_items=resp.total_items,
                excluded=len(rows) - len(actions),
            )
            return actions, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
