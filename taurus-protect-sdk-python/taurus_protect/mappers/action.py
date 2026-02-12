"""Action mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import (
    safe_bool,
    safe_datetime,
    safe_string,
)
from taurus_protect.models.action import (
    Action,
    ActionAttribute,
    ActionDetails,
    ActionTrail,
)

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def action_from_dto(dto: Any) -> Optional[Action]:
    """
    Convert OpenAPI TgvalidatordActionEnvelope to domain Action.

    Args:
        dto: OpenAPI action envelope DTO.

    Returns:
        Domain Action model or None if dto is None.
    """
    if dto is None:
        return None

    # Extract action details if present
    action_details = None
    dto_action = getattr(dto, "action", None)
    if dto_action is not None:
        action_details = action_details_from_dto(dto_action)

    # Extract attributes if present
    attributes: List[ActionAttribute] = []
    dto_attributes = getattr(dto, "attributes", None)
    if dto_attributes is not None:
        attributes = [
            attr
            for dto_attr in dto_attributes
            if (attr := action_attribute_from_dto(dto_attr)) is not None
        ]

    # Extract trails if present
    trails: List[ActionTrail] = []
    dto_trails = getattr(dto, "trails", None)
    if dto_trails is not None:
        trails = [
            trail
            for dto_trail in dto_trails
            if (trail := action_trail_from_dto(dto_trail)) is not None
        ]

    return Action(
        id=safe_string(getattr(dto, "id", None)),
        tenant_id=safe_string(getattr(dto, "tenant_id", None)),
        label=safe_string(getattr(dto, "label", None)),
        status=safe_string(getattr(dto, "status", None)),
        auto_approve=safe_bool(getattr(dto, "auto_approve", None)),
        action=action_details,
        attributes=attributes,
        trails=trails,
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
        last_checked_at=safe_datetime(getattr(dto, "lastcheckeddate", None)),
    )


def actions_from_dto(dtos: Optional[List[Any]]) -> List[Action]:
    """
    Convert list of OpenAPI action DTOs to domain Actions.

    Args:
        dtos: List of OpenAPI action DTOs.

    Returns:
        List of domain Action models.
    """
    if dtos is None:
        return []
    return [a for dto in dtos if (a := action_from_dto(dto)) is not None]


def action_attribute_from_dto(dto: Any) -> Optional[ActionAttribute]:
    """
    Convert OpenAPI TgvalidatordActionAttribute to domain ActionAttribute.

    Args:
        dto: OpenAPI action attribute DTO.

    Returns:
        Domain ActionAttribute model or None if dto is None.
    """
    if dto is None:
        return None

    return ActionAttribute(
        id=safe_string(getattr(dto, "id", None)),
        key=safe_string(getattr(dto, "key", None)),
        value=safe_string(getattr(dto, "value", None)),
    )


def action_trail_from_dto(dto: Any) -> Optional[ActionTrail]:
    """
    Convert OpenAPI TgvalidatordActionEnvelopeTrail to domain ActionTrail.

    Args:
        dto: OpenAPI action trail DTO.

    Returns:
        Domain ActionTrail model or None if dto is None.
    """
    if dto is None:
        return None

    return ActionTrail(
        id=safe_string(getattr(dto, "id", None)),
        user_id=safe_string(getattr(dto, "user_id", None)),
        action=safe_string(getattr(dto, "action", None)),
        status=safe_string(getattr(dto, "status", None)),
        timestamp=safe_datetime(getattr(dto, "timestamp", None)),
    )


def action_details_from_dto(dto: Any) -> Optional[ActionDetails]:
    """
    Convert OpenAPI TgvalidatordAction to domain ActionDetails.

    Args:
        dto: OpenAPI action DTO.

    Returns:
        Domain ActionDetails model or None if dto is None.
    """
    if dto is None:
        return None

    return ActionDetails(
        type=safe_string(getattr(dto, "type", None)),
        parameters=getattr(dto, "parameters", {}) or {},
    )
