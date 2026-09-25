"""Visibility group service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List

from taurus_protect.mappers.visibility_group import (
    visibility_group_from_dto,
    visibility_groups_from_dto,
)
from taurus_protect.models.visibility_group import VisibilityGroup
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class VisibilityGroupService(BaseService):
    """
    Service for visibility group management operations.

    Provides methods to list and get restricted visibility groups.
    Visibility groups control which users can see and access specific
    wallets and resources.

    Example:
        >>> # List visibility groups
        >>> groups = client.visibility_groups.list()
        >>> for group in groups:
        ...     print(f"{group.name}: {group.user_count} users")
        >>>
        >>> # Get single visibility group
        >>> group = client.visibility_groups.get("group-123")
        >>> print(f"Description: {group.description}")
    """

    def __init__(self, api_client: Any, visibility_groups_api: Any) -> None:
        """
        Initialize visibility group service.

        Args:
            api_client: The OpenAPI client instance.
            visibility_groups_api: The RestrictedVisibilityGroupsApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._visibility_groups_api = visibility_groups_api

    def get(self, group_id: str) -> VisibilityGroup:
        """
        Get a visibility group by ID.

        Note: The API returns users within the group via a separate endpoint.
        This method fetches the group from the list and enriches it with
        user information if available.

        Args:
            group_id: The visibility group ID to retrieve.

        Returns:
            The visibility group.

        Raises:
            ValueError: If group_id is invalid.
            NotFoundError: If visibility group not found.
            APIError: If API request fails.
        """
        self._validate_required(group_id, "group_id")

        try:
            # The API doesn't have a direct "get by ID" endpoint for visibility groups
            # We need to list all groups and find the matching one
            groups = self.list()

            for group in groups:
                if group.id == group_id:
                    # Try to fetch users for this group
                    try:
                        users_resp = self._visibility_groups_api.user_service_get_users_by_visibility_group_id(
                            visibility_group_id=group_id,
                        )
                        users_result = getattr(users_resp, "result", None)
                        if users_result:
                            from taurus_protect.mappers.visibility_group import (
                                visibility_group_user_from_dto,
                            )

                            group.users = [
                                u
                                for dto in users_result
                                if (u := visibility_group_user_from_dto(dto)) is not None
                            ]
                    except Exception:
                        # If we can't fetch users, return the group without them
                        pass
                    return group

            from taurus_protect.errors import NotFoundError

            raise NotFoundError(f"Visibility group {group_id} not found")
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(self) -> List[VisibilityGroup]:
        """
        List visibility groups: every group the endpoint returns, which does not page.

        Returns:
            Every visibility group.

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._visibility_groups_api.user_service_get_visibility_groups()

            return visibility_groups_from_dto(resp.result or [])
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_users(self, group_id: str) -> List[Any]:
        """
        Get users in a visibility group.

        Args:
            group_id: The visibility group ID.

        Returns:
            List of users in the group.

        Raises:
            ValueError: If group_id is invalid.
            APIError: If API request fails.
        """
        self._validate_required(group_id, "group_id")

        try:
            resp = self._visibility_groups_api.user_service_get_users_by_visibility_group_id(
                visibility_group_id=group_id,
            )

            result = getattr(resp, "result", None)
            if result is None:
                return []

            from taurus_protect.mappers.visibility_group import (
                visibility_group_user_from_dto,
            )

            return [u for dto in result if (u := visibility_group_user_from_dto(dto)) is not None]
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
