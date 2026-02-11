"""Group service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.user import group_from_dto, groups_from_dto
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.user import Group
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class GroupService(BaseService):
    """
    Service for group management operations.

    Provides methods to list and retrieve groups.

    Example:
        >>> # List groups
        >>> groups, pagination = client.groups.list(limit=50, offset=0)
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

            return group
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Group], Optional[Pagination]]:
        """
        List groups with pagination.

        Args:
            limit: Maximum number of groups to return (must be positive).
            offset: Number of groups to skip (must be non-negative).

        Returns:
            Tuple of (groups list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._groups_api.user_service_get_groups(
                limit=str(limit),
                offset=str(offset),
                ids=None,
                external_group_ids=None,
                query=None,
            )

            result = getattr(resp, "result", None)
            groups = groups_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=offset,
                limit=limit,
            )

            return groups, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
