"""Visibility group mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import (
    parse_string_to_int,
    safe_datetime,
    safe_string,
)
from taurus_protect.models.visibility_group import VisibilityGroup, VisibilityGroupUser

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def visibility_group_from_dto(dto: Any) -> Optional[VisibilityGroup]:
    """
    Convert OpenAPI TgvalidatordInternalVisibilityGroup to domain VisibilityGroup.

    Args:
        dto: OpenAPI visibility group DTO.

    Returns:
        Domain VisibilityGroup model or None if dto is None.
    """
    if dto is None:
        return None

    # Extract users if present
    users: List[VisibilityGroupUser] = []
    dto_users = getattr(dto, "users", None)
    if dto_users is not None:
        users = [
            u
            for dto_user in dto_users
            if (u := visibility_group_user_from_dto(dto_user)) is not None
        ]

    return VisibilityGroup(
        id=safe_string(getattr(dto, "id", None)),
        tenant_id=safe_string(getattr(dto, "tenant_id", None)),
        name=safe_string(getattr(dto, "name", None)),
        description=safe_string(getattr(dto, "description", None)),
        user_count=parse_string_to_int(getattr(dto, "user_count", None)),
        users=users,
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
    )


def visibility_groups_from_dto(dtos: Optional[List[Any]]) -> List[VisibilityGroup]:
    """
    Convert list of OpenAPI visibility group DTOs to domain VisibilityGroups.

    Args:
        dtos: List of OpenAPI visibility group DTOs.

    Returns:
        List of domain VisibilityGroup models.
    """
    if dtos is None:
        return []
    return [vg for dto in dtos if (vg := visibility_group_from_dto(dto)) is not None]


def visibility_group_user_from_dto(dto: Any) -> Optional[VisibilityGroupUser]:
    """
    Convert OpenAPI TgvalidatordInternalVisibilityGroupUser to domain VisibilityGroupUser.

    Args:
        dto: OpenAPI visibility group user DTO.

    Returns:
        Domain VisibilityGroupUser model or None if dto is None.
    """
    if dto is None:
        return None

    return VisibilityGroupUser(
        id=safe_string(getattr(dto, "id", None)),
        email=safe_string(getattr(dto, "email", None)),
        name=safe_string(getattr(dto, "name", None)),
    )
