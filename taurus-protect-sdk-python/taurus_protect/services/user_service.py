"""User service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect._internal.openapi.models.user_service_create_attribute_body import (
    UserServiceCreateAttributeBody,
)
from taurus_protect.mappers.user import user_from_dto, users_from_dto
from taurus_protect.models.pagination import (
    MAX_PAGE_SIZE,
    PLUS_MIN_ROWS_LIMIT,
    Pagination,
    offset_pagination,
    offset_query,
    resolve_offset,
    resolve_page_size,
)
from taurus_protect.models.user import User
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def _with_computed_flags(user: User) -> User:
    """validatord omits a false bool, so where it computed the flags an absent one is False."""
    return user.model_copy(
        update={
            "enforced_in_rules": bool(user.enforced_in_rules),
            "groups": [
                g.model_copy(update={"enforced_in_rules": bool(g.enforced_in_rules)})
                for g in user.groups
            ],
        }
    )


class UserService(BaseService):
    """
    Service for user management operations.

    Provides methods to list and retrieve users.

    Example:
        >>> # List users
        >>> users, pagination = client.users.list(limit=100, offset=0)
        >>> for user in users:
        ...     print(f"{user.email}: {user.status}")
        >>>
        >>> # Get single user
        >>> user = client.users.get("user-123")
        >>> print(f"Name: {user.first_name} {user.last_name}")
        >>>
        >>> # Get current user
        >>> current_user = client.users.get_current()
        >>> print(f"Logged in as: {current_user.email}")
    """

    def __init__(self, api_client: Any, users_api: Any) -> None:
        """
        Initialize user service.

        Args:
            api_client: The OpenAPI client instance.
            users_api: The UsersApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._users_api = users_api

    def get(self, user_id: str) -> User:
        """
        Get a user by ID.

        Args:
            user_id: The user ID to retrieve.

        Returns:
            The user.

        Raises:
            ValueError: If user_id is invalid.
            NotFoundError: If user not found.
            APIError: If API request fails.
        """
        self._validate_required(user_id, "user_id")

        try:
            resp = self._users_api.user_service_get_user(user_id)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"User {user_id} not found")

            user = user_from_dto(result)
            if user is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"User {user_id} not found")

            return user
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_current(self) -> User:
        """
        Get the current authenticated user.

        Returns:
            The current user.

        Raises:
            APIError: If API request fails.
        """
        try:
            # GetMe computes the enforced-in-rules flags only when asked.
            resp = self._users_api.user_service_get_me(check_enforced_in_rules=True)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to get current user: no result returned")

            user = user_from_dto(result)
            if user is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to get current user: invalid response")

            return _with_computed_flags(user)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
    ) -> Tuple[List[User], Pagination]:
        """
        List users, one page at a time.

        Args:
            limit: Page size (default 20, max 100).
            offset: Number of users to skip; pass ``pagination.next_offset`` to continue.

        Returns:
            Tuple of (users, pagination). A synthetic daemon user can be appended
            beyond ``limit``; ``next_offset`` accounts for it.

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        return self._list_page(resolve_page_size(limit, "limit"), resolve_offset(offset))

    def get_users_by_email(self, emails: List[str]) -> List[User]:
        """
        Get users by their email addresses, reading every page.

        Args:
            emails: List of email addresses to search for.

        Returns:
            List of users matching the provided email addresses.

        Raises:
            ValueError: If emails list is empty.
            APIError: If API request fails.
        """
        if not emails:
            raise ValueError("emails cannot be empty")

        users: List[User] = []
        seen = set()
        offset = 0
        while True:
            page, pagination = self._list_page(MAX_PAGE_SIZE, offset, emails=emails)
            for user in page:
                if user.id not in seen:
                    seen.add(user.id)
                    users.append(user)
            if not pagination.has_more:
                return users
            offset = pagination.next_offset

    def _list_page(self, limit: int, offset: int, **filters: Any) -> Tuple[List[User], Pagination]:
        try:
            resp = self._users_api.user_service_get_users(**offset_query(limit, offset), **filters)

            rows = resp.result or []
            users = [_with_computed_flags(u) for u in users_from_dto(rows)]
            pagination = offset_pagination(
                PLUS_MIN_ROWS_LIMIT,
                limit=limit,
                offset=offset,
                served_rows=len(rows),
                total_items=resp.total_items,
                excluded=len(rows) - len(users),
            )
            return users, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create_user_attribute(self, user_id: str, key: str, value: str) -> None:
        """
        Create an attribute for a user.

        Args:
            user_id: The user ID to create the attribute for.
            key: The attribute key.
            value: The attribute value.

        Raises:
            ValueError: If user_id, key, or value is invalid.
            APIError: If API request fails.
        """
        self._validate_required(user_id, "user_id")
        self._validate_required(key, "key")

        try:
            body = UserServiceCreateAttributeBody(key=key, value=value)
            self._users_api.user_service_create_attribute(
                user_id=user_id,
                body=body,
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
