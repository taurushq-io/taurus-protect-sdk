"""Pagination model for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import Optional, Union

from pydantic import BaseModel, Field


class Pagination(BaseModel):
    """
    Pagination information returned from list operations.

    Attributes:
        total_items: Total number of items available.
        offset: Current offset in the result set.
        limit: Maximum items per page.
        has_more: Whether more items are available.
    """

    total_items: int = Field(default=0, description="Total number of items available")
    offset: int = Field(default=0, description="Current offset in the result set")
    limit: int = Field(default=0, description="Maximum items per page")
    has_more: bool = Field(default=False, description="Whether more items are available")

    model_config = {"frozen": True}

    @classmethod
    def from_response(
        cls,
        total_items: Optional[Union[str, int]],
        offset: Optional[Union[str, int]],
        limit: int,
    ) -> "Pagination":
        """
        Create Pagination from API response fields.

        Args:
            total_items: Total items (string or int from API).
            offset: Current offset (string or int from API).
            limit: Requested limit.

        Returns:
            Pagination instance.
        """
        total = int(total_items) if total_items else 0
        off = int(offset) if offset else 0

        return cls(
            total_items=total,
            offset=off,
            limit=limit,
            has_more=off + limit < total,
        )
