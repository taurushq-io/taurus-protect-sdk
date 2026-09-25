"""Group service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.user import group_from_dto, groups_from_dto
from taurus_protect.models.pagination import (
    PLUS_MIN_ROWS_LIMIT,
    Pagination,
    offset_pagination,
    offset_query,
    resolve_offset,
    resolve_page_size,
)
from taurus_protect.models.user import Group
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def _with_computed_flags(group: Group) -> Group:
    """validatord omits a false bool, so where it computed the flags an absent one is False."""
    return group.model_copy(
        update={
            "enforced_in_rules": bool(group.enforced_in_rules),
            "users": [
                u.model_copy(update={"enforced_in_rules": bool(u.enforced_in_rules)})
                for u in group.users
            ],
        }
    )


class GroupService(BaseService):
    """
    Service for group management operations.

    Provides methods to list and retrieve groups.

    Example:
        >>> # List groups
        >>> groups, pagination = client.groups.list(limit=100, offset=0)
        >>> for group in groups:
        ...     print(f"{group.name}: {len(group.users)} users")
        >>>
        >>> # Get single group
        >>> group = client.groups.get("group-123")
        >>> print(f"Group: {group.name}")
    """

    def __init__(self, api_client: Any, groups_api: Any) -> None:
        """
        Initialize group service.

        Args:
            api_client: The OpenAPI client instance.
            groups_api: The GroupsApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._groups_api = groups_api

    def get(self, group_id: str) -> Group:
        """
        Get a group by ID.

        Args:
            group_id: The group ID to retrieve.

        Returns:
            The group.

        Raises:
            ValueError: If group_id is invalid.
            NotFoundError: If group not found.
            APIError: If API request fails.
        """
        self._validate_required(group_id, "group_id")

        try:
            # The API doesn't have a direct get-by-id endpoint,
            # so we use the list endpoint with id filter
            resp = self._groups_api.user_service_get_groups(
                limit=str(1),
                offset=str(0),
                ids=[group_id],
                external_group_ids=None,
                query=None,
            )

            result = getattr(resp, "result", None)
            if result is None or len(result) == 0:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Group {group_id} not found")

            group = group_from_dto(result[0])
            if group is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Group {group_id} not found")

            return _with_computed_flags(group)
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
    ) -> Tuple[List[Group], Pagination]:
        """
        List groups, one page at a time.

        Args:
            limit: Page size (default 20, max 100).
            offset: Number of groups to skip; pass ``pagination.next_offset`` to continue.

        Returns:
            Tuple of (groups, pagination). A synthetic technical group can be appended
            beyond ``limit``; ``next_offset`` accounts for it.

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        page_size = resolve_page_size(limit, "limit")
        start = resolve_offset(offset)

        try:
            resp = self._groups_api.user_service_get_groups(**offset_query(page_size, start))

            rows = resp.result or []
            groups = [_with_computed_flags(g) for g in groups_from_dto(rows)]
            pagination = offset_pagination(
                PLUS_MIN_ROWS_LIMIT,
                limit=page_size,
                offset=start,
                served_rows=len(rows),
                total_items=resp.total_items,
                excluded=len(rows) - len(groups),
            )
            return groups, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
