"""Tag service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers.user import tag_from_dto, tags_from_dto
from taurus_protect.models.user import Tag
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class TagService(BaseService):
    """
    Service for tag management operations.

    Provides methods to list, get, create, and delete tags.
    Tags can be applied to wallets, addresses, and other entities
    for organization and filtering.

    Example:
        >>> # List tags
        >>> tags = client.tags.list()
        >>> for tag in tags:
        ...     print(f"{tag.name}: {tag.color}")
        >>>
        >>> # Get single tag
        >>> tag = client.tags.get("tag-123")
        >>> print(f"Tag: {tag.name}")
        >>>
        >>> # Create tag
        >>> tag = client.tags.create("Important", "#FF0000")
        >>> print(f"Created tag: {tag.id}")
        >>>
        >>> # Delete tag
        >>> client.tags.delete("tag-123")
    """

    def __init__(self, api_client: Any, tags_api: Any) -> None:
        """
        Initialize tag service.

        Args:
            api_client: The OpenAPI client instance.
            tags_api: The TagsApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._tags_api = tags_api

    def get(self, tag_id: str) -> Tag:
        """
        Get a tag by ID.

        Args:
            tag_id: The tag ID to retrieve.

        Returns:
            The tag.

        Raises:
            ValueError: If tag_id is invalid.
            NotFoundError: If tag not found.
            APIError: If API request fails.
        """
        self._validate_required(tag_id, "tag_id")

        try:
            # The API doesn't have a direct get-by-id endpoint,
            # so we use the list endpoint with id filter
            resp = self._tags_api.tag_service_get_tags(
                ids=[tag_id],
                query=None,
            )

            result = getattr(resp, "result", None)
            if result is None or len(result) == 0:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Tag {tag_id} not found")

            tag = tag_from_dto(result[0])
            if tag is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Tag {tag_id} not found")

            return tag
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, NotFoundError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        limit: int = 50,
        offset: int = 0,
    ) -> List[Tag]:
        """
        List tags.

        Note: The tags API does not support pagination parameters,
        so limit and offset are provided for API consistency but
        filtering is done client-side.

        Args:
            limit: Maximum number of tags to return (must be positive).
            offset: Number of tags to skip (must be non-negative).

        Returns:
            List of tags.

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._tags_api.tag_service_get_tags(
                ids=None,
                query=None,
            )

            result = getattr(resp, "result", None)
            tags = tags_from_dto(result) if result else []

            # Apply client-side pagination since API doesn't support it
            return tags[offset : offset + limit]
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create(self, name: str, color: str) -> Tag:
        """
        Create a new tag.

        Args:
            name: The tag name/value.
            color: The tag color (hex code or color name).

        Returns:
            The created tag.

        Raises:
            ValueError: If name or color is empty.
            APIError: If API request fails.
        """
        self._validate_required(name, "name")
        self._validate_required(color, "color")

        try:
            body = {
                "value": name,
                "color": color,
            }

            resp = self._tags_api.tag_service_create_tag(body=body)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create tag: no result returned")

            tag = tag_from_dto(result)
            if tag is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create tag: invalid response")

            return tag
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def delete(self, tag_id: str) -> None:
        """
        Delete a tag.

        This removes the tag and all its assignments from entities.

        Args:
            tag_id: The tag ID to delete.

        Raises:
            ValueError: If tag_id is empty.
            NotFoundError: If tag not found.
            APIError: If API request fails.
        """
        self._validate_required(tag_id, "tag_id")

        try:
            self._tags_api.tag_service_delete_tag(id=tag_id)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
