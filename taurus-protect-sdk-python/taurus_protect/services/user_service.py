"""User service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect._internal.openapi.models.user_service_create_attribute_body import (
    UserServiceCreateAttributeBody,
)
from taurus_protect.mappers.user import user_from_dto, users_from_dto
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.user import User
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class UserService(BaseService):
    """
    Service for user management operations.

    Provides methods to list and retrieve users.

    Example:
        >>> # List users
        >>> users, pagination = client.users.list(limit=50, offset=0)
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
            resp = self._users_api.user_service_get_me()

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to get current user: no result returned")

            user = user_from_dto(result)
            if user is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to get current user: invalid response")

            return user
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[User], Optional[Pagination]]:
        """
        List users with pagination.

        Args:
            limit: Maximum number of users to return (must be positive).
            offset: Number of users to skip (must be non-negative).

        Returns:
            Tuple of (users list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._users_api.user_service_get_users(
                limit=str(limit),
                offset=str(offset),
                ids=None,
                external_user_ids=None,
                emails=None,
                query=None,
                public_key=None,
                exclude_technical_users=None,
                roles=None,
                status=None,
                totp_enabled=None,
                group_ids=None,
            )

            result = getattr(resp, "result", None)
            users = users_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=offset,
                limit=limit,
            )

            return users, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_users_by_email(self, emails: List[str]) -> List[User]:
        """
        Get users by their email addresses.

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

        try:
            resp = self._users_api.user_service_get_users(
                limit=None,
                offset=None,
                ids=None,
                external_user_ids=None,
                emails=emails,
                query=None,
                public_key=None,
                exclude_technical_users=None,
                roles=None,
                status=None,
                totp_enabled=None,
                group_ids=None,
            )

            result = getattr(resp, "result", None)
            return users_from_dto(result) if result else []
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
