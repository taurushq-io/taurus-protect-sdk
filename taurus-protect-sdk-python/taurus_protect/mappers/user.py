"""User, Group, and Tag mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import (
    safe_bool,
    safe_datetime,
    safe_list,
    safe_string,
)
from taurus_protect.models.user import Group, GroupUser, Tag, User, UserAttribute, UserGroup

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def user_from_dto(dto: Any) -> Optional[User]:
    """
    Convert OpenAPI TgvalidatordInternalUser to domain User.

    Args:
        dto: OpenAPI user DTO (TgvalidatordInternalUser).

    Returns:
        Domain User model or None if dto is None.
    """
    if dto is None:
        return None

    # Extract groups if present
    groups: List[UserGroup] = []
    dto_groups = getattr(dto, "groups", None)
    if dto_groups is not None:
        groups = [user_group_from_dto(g) for g in dto_groups if g is not None]

    # Extract attributes if present
    attributes: List[UserAttribute] = []
    dto_attributes = getattr(dto, "attributes", None)
    if dto_attributes is not None:
        attributes = [user_attribute_from_dto(attr) for attr in dto_attributes if attr is not None]

    # Extract roles
    roles = safe_list(getattr(dto, "roles", None))

    return User(
        id=safe_string(getattr(dto, "id", None)),
        external_user_id=getattr(dto, "external_user_id", None),
        tenant_id=getattr(dto, "tenant_id", None),
        username=getattr(dto, "username", None),
        email=getattr(dto, "email", None),
        first_name=getattr(dto, "first_name", None),
        last_name=getattr(dto, "last_name", None),
        status=getattr(dto, "status", None),
        roles=roles,
        public_key=getattr(dto, "public_key", None),
        groups=groups,
        totp_enabled=safe_bool(getattr(dto, "totp_enabled", None)),
        password_changed=safe_bool(getattr(dto, "password_changed", None)),
        enforced_in_rules=safe_bool(getattr(dto, "enforced_in_rules", None)),
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
        last_login=safe_datetime(getattr(dto, "last_login", None)),
        attributes=attributes,
    )


def users_from_dto(dtos: Optional[List[Any]]) -> List[User]:
    """
    Convert list of OpenAPI user DTOs to domain Users.

    Args:
        dtos: List of OpenAPI user DTOs.

    Returns:
        List of domain User models.
    """
    if dtos is None:
        return []
    return [u for dto in dtos if (u := user_from_dto(dto)) is not None]


def user_group_from_dto(dto: Any) -> UserGroup:
    """
    Convert OpenAPI TgvalidatordInternalUserGroup to domain UserGroup.

    Args:
        dto: OpenAPI user group DTO.

    Returns:
        Domain UserGroup model.
    """
    return UserGroup(
        id=safe_string(getattr(dto, "id", None)),
        name=safe_string(getattr(dto, "name", None)),
    )


def user_attribute_from_dto(dto: Any) -> UserAttribute:
    """
    Convert OpenAPI InternalUserAttribute to domain UserAttribute.

    Args:
        dto: OpenAPI user attribute DTO.

    Returns:
        Domain UserAttribute model.
    """
    return UserAttribute(
        id=safe_string(getattr(dto, "id", None)),
        key=safe_string(getattr(dto, "key", None)),
        value=safe_string(getattr(dto, "value", None)),
    )


def group_from_dto(dto: Any) -> Optional[Group]:
    """
    Convert OpenAPI TgvalidatordInternalGroup to domain Group.

    Args:
        dto: OpenAPI group DTO (TgvalidatordInternalGroup).

    Returns:
        Domain Group model or None if dto is None.
    """
    if dto is None:
        return None

    # Extract users if present
    users: List[GroupUser] = []
    dto_users = getattr(dto, "users", None)
    if dto_users is not None:
        users = [group_user_from_dto(u) for u in dto_users if u is not None]

    return Group(
        id=safe_string(getattr(dto, "id", None)),
        external_group_id=getattr(dto, "external_group_id", None),
        tenant_id=getattr(dto, "tenant_id", None),
        name=safe_string(getattr(dto, "name", None)),
        email=getattr(dto, "email", None),
        description=getattr(dto, "description", None),
        users=users,
        enforced_in_rules=safe_bool(getattr(dto, "enforced_in_rules", None)),
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
    )


def groups_from_dto(dtos: Optional[List[Any]]) -> List[Group]:
    """
    Convert list of OpenAPI group DTOs to domain Groups.

    Args:
        dtos: List of OpenAPI group DTOs.

    Returns:
        List of domain Group models.
    """
    if dtos is None:
        return []
    return [g for dto in dtos if (g := group_from_dto(dto)) is not None]


def group_user_from_dto(dto: Any) -> GroupUser:
    """
    Convert OpenAPI TgvalidatordInternalGroupUser to domain GroupUser.

    Args:
        dto: OpenAPI group user DTO.

    Returns:
        Domain GroupUser model.
    """
    return GroupUser(
        id=safe_string(getattr(dto, "id", None)),
        email=getattr(dto, "email", None),
    )


def tag_from_dto(dto: Any) -> Optional[Tag]:
    """
    Convert OpenAPI TgvalidatordTag to domain Tag.

    Args:
        dto: OpenAPI tag DTO (TgvalidatordTag).

    Returns:
        Domain Tag model or None if dto is None.
    """
    if dto is None:
        return None

    return Tag(
        id=safe_string(getattr(dto, "id", None)),
        name=safe_string(getattr(dto, "value", None)),  # API uses 'value' for tag name
        color=getattr(dto, "color", None),
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
    )


def tags_from_dto(dtos: Optional[List[Any]]) -> List[Tag]:
    """
    Convert list of OpenAPI tag DTOs to domain Tags.

    Args:
        dtos: List of OpenAPI tag DTOs.

    Returns:
        List of domain Tag models.
    """
    if dtos is None:
        return []
    return [t for dto in dtos if (t := tag_from_dto(dto)) is not None]
