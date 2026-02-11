"""Mapper functions for Taurus Network pledge DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import (
    safe_datetime,
    safe_list,
    safe_string,
)
from taurus_protect.models.taurus_network.pledge import (
    Pledge,
    PledgeAction,
    PledgeActionMetadata,
    PledgeActionTrail,
    PledgeAttribute,
    PledgeDurationSetup,
    PledgeTrail,
    PledgeWithdrawal,
    PledgeWithdrawalTrail,
)

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def pledge_attribute_from_dto(dto: Any) -> Optional[PledgeAttribute]:
    """
    Convert pledge attribute DTO to domain model.

    Args:
        dto: The OpenAPI pledge attribute DTO.

    Returns:
        PledgeAttribute or None if DTO is None.
    """
    if dto is None:
        return None

    return PledgeAttribute(
        key=safe_string(getattr(dto, "key", None)),
        value=safe_string(getattr(dto, "value", None)),
    )


def pledge_duration_setup_from_dto(dto: Any) -> Optional[PledgeDurationSetup]:
    """
    Convert pledge duration setup DTO to domain model.

    Args:
        dto: The OpenAPI duration setup DTO.

    Returns:
        PledgeDurationSetup or None if DTO is None.
    """
    if dto is None:
        return None

    return PledgeDurationSetup(
        start_date=safe_datetime(getattr(dto, "start_date", None)),
        end_date=safe_datetime(getattr(dto, "end_date", None)),
    )


def pledge_trail_from_dto(dto: Any) -> Optional[PledgeTrail]:
    """
    Convert pledge trail DTO to domain model.

    Args:
        dto: The OpenAPI pledge trail DTO.

    Returns:
        PledgeTrail or None if DTO is None.
    """
    if dto is None:
        return None

    return PledgeTrail(
        id=safe_string(getattr(dto, "id", None)),
        action=safe_string(getattr(dto, "action", None)),
        actor=safe_string(getattr(dto, "actor", None)),
        timestamp=safe_datetime(getattr(dto, "timestamp", None)),
        comment=getattr(dto, "comment", None),
    )


def pledge_from_dto(dto: Any) -> Optional[Pledge]:
    """
    Convert pledge DTO to domain model.

    Args:
        dto: The OpenAPI TgvalidatordTnPledge DTO.

    Returns:
        Pledge or None if DTO is None.
    """
    if dto is None:
        return None

    # Parse attributes if present
    attributes: List[PledgeAttribute] = []
    raw_attributes = getattr(dto, "attributes", None)
    if raw_attributes:
        for attr_dto in raw_attributes:
            attr = pledge_attribute_from_dto(attr_dto)
            if attr:
                attributes.append(attr)

    # Parse trails if present
    trails: List[PledgeTrail] = []
    raw_trails = getattr(dto, "trails", None)
    if raw_trails:
        for trail_dto in raw_trails:
            trail = pledge_trail_from_dto(trail_dto)
            if trail:
                trails.append(trail)

    # Parse duration setup
    duration_setup = pledge_duration_setup_from_dto(getattr(dto, "duration_setup", None))

    return Pledge(
        id=safe_string(getattr(dto, "id", None)),
        shared_address_id=safe_string(getattr(dto, "shared_address_id", None)),
        owner_participant_id=safe_string(getattr(dto, "owner_participant_id", None)),
        target_participant_id=safe_string(getattr(dto, "target_participant_id", None)),
        currency_id=safe_string(getattr(dto, "currency_id", None)),
        blockchain=safe_string(getattr(dto, "blockchain", None)),
        network=safe_string(getattr(dto, "network", None)),
        arg1=getattr(dto, "arg1", None),
        arg2=getattr(dto, "arg2", None),
        amount=safe_string(getattr(dto, "amount", None)) or "0",
        status=safe_string(getattr(dto, "status", None)),
        pledge_type=safe_string(getattr(dto, "pledge_type", None)),
        direction=safe_string(getattr(dto, "direction", None)),
        external_reference_id=getattr(dto, "external_reference_id", None),
        reconciliation_note=getattr(dto, "reconciliation_note", None),
        wl_address_id=getattr(dto, "wladdress_id", None),
        origin_creation_date=safe_datetime(getattr(dto, "origin_creation_date", None)),
        unpledge_date=safe_datetime(getattr(dto, "unpledge_date", None)),
        duration_setup=duration_setup,
        attributes=attributes,
        trails=trails,
        created_at=safe_datetime(getattr(dto, "created_at", None)),
        updated_at=safe_datetime(getattr(dto, "updated_at", None)),
    )


def pledges_from_dto(dto_list: Any) -> List[Pledge]:
    """
    Convert list of pledge DTOs to domain models.

    Args:
        dto_list: List of OpenAPI TgvalidatordTnPledge DTOs.

    Returns:
        List of Pledge domain models.
    """
    if not dto_list:
        return []

    result: List[Pledge] = []
    for dto in dto_list:
        pledge = pledge_from_dto(dto)
        if pledge:
            result.append(pledge)

    return result


def pledge_action_metadata_from_dto(dto: Any) -> Optional[PledgeActionMetadata]:
    """
    Convert pledge action metadata DTO to domain model.

    Args:
        dto: The OpenAPI metadata DTO.

    Returns:
        PledgeActionMetadata or None if DTO is None.
    """
    if dto is None:
        return None

    return PledgeActionMetadata(
        hash=safe_string(getattr(dto, "hash", None)),
        payload=getattr(dto, "payload", None),
    )


def pledge_action_trail_from_dto(dto: Any) -> Optional[PledgeActionTrail]:
    """
    Convert pledge action trail DTO to domain model.

    Args:
        dto: The OpenAPI pledge action trail DTO.

    Returns:
        PledgeActionTrail or None if DTO is None.
    """
    if dto is None:
        return None

    return PledgeActionTrail(
        id=safe_string(getattr(dto, "id", None)),
        action=safe_string(getattr(dto, "action", None)),
        actor=safe_string(getattr(dto, "actor", None)),
        timestamp=safe_datetime(getattr(dto, "timestamp", None)),
        comment=getattr(dto, "comment", None),
    )


def pledge_action_from_dto(dto: Any) -> Optional[PledgeAction]:
    """
    Convert pledge action DTO to domain model.

    Args:
        dto: The OpenAPI TgvalidatordTnPledgeAction DTO.

    Returns:
        PledgeAction or None if DTO is None.
    """
    if dto is None:
        return None

    # Parse metadata
    metadata = pledge_action_metadata_from_dto(getattr(dto, "metadata", None))

    # Parse trails if present
    trails: List[PledgeActionTrail] = []
    raw_trails = getattr(dto, "trails", None)
    if raw_trails:
        for trail_dto in raw_trails:
            trail = pledge_action_trail_from_dto(trail_dto)
            if trail:
                trails.append(trail)

    # Parse needs_approval_from
    needs_approval_from = safe_list(getattr(dto, "needs_approval_from", None))

    return PledgeAction(
        id=safe_string(getattr(dto, "id", None)),
        pledge_id=safe_string(getattr(dto, "pledge_id", None)),
        action_type=safe_string(getattr(dto, "action_type", None)),
        status=safe_string(getattr(dto, "status", None)),
        metadata=metadata,
        rule=getattr(dto, "rule", None),
        needs_approval_from=needs_approval_from,
        pledge_withdrawal_id=getattr(dto, "pledge_withdrawal_id", None),
        envelope=getattr(dto, "envelope", None),
        trails=trails,
        created_at=safe_datetime(getattr(dto, "created_at", None)),
        last_approval_date=safe_datetime(getattr(dto, "last_approval_date", None)),
    )


def pledge_actions_from_dto(dto_list: Any) -> List[PledgeAction]:
    """
    Convert list of pledge action DTOs to domain models.

    Args:
        dto_list: List of OpenAPI TgvalidatordTnPledgeAction DTOs.

    Returns:
        List of PledgeAction domain models.
    """
    if not dto_list:
        return []

    result: List[PledgeAction] = []
    for dto in dto_list:
        action = pledge_action_from_dto(dto)
        if action:
            result.append(action)

    return result


def pledge_withdrawal_trail_from_dto(dto: Any) -> Optional[PledgeWithdrawalTrail]:
    """
    Convert pledge withdrawal trail DTO to domain model.

    Args:
        dto: The OpenAPI pledge withdrawal trail DTO.

    Returns:
        PledgeWithdrawalTrail or None if DTO is None.
    """
    if dto is None:
        return None

    return PledgeWithdrawalTrail(
        id=safe_string(getattr(dto, "id", None)),
        action=safe_string(getattr(dto, "action", None)),
        actor=safe_string(getattr(dto, "actor", None)),
        timestamp=safe_datetime(getattr(dto, "timestamp", None)),
        comment=getattr(dto, "comment", None),
    )


def pledge_withdrawal_from_dto(dto: Any) -> Optional[PledgeWithdrawal]:
    """
    Convert pledge withdrawal DTO to domain model.

    Args:
        dto: The OpenAPI TgvalidatordTnPledgeWithdrawal DTO.

    Returns:
        PledgeWithdrawal or None if DTO is None.
    """
    if dto is None:
        return None

    # Parse trails if present
    trails: List[PledgeWithdrawalTrail] = []
    raw_trails = getattr(dto, "trails", None)
    if raw_trails:
        for trail_dto in raw_trails:
            trail = pledge_withdrawal_trail_from_dto(trail_dto)
            if trail:
                trails.append(trail)

    return PledgeWithdrawal(
        id=safe_string(getattr(dto, "id", None)),
        pledge_id=safe_string(getattr(dto, "pledge_id", None)),
        destination_shared_address_id=safe_string(
            getattr(dto, "destination_shared_address_id", None)
        ),
        amount=safe_string(getattr(dto, "amount", None)) or "0",
        status=safe_string(getattr(dto, "status", None)),
        tx_hash=getattr(dto, "tx_hash", None),
        tx_id=getattr(dto, "tx_id", None),
        request_id=getattr(dto, "request_id", None),
        tx_block_number=getattr(dto, "tx_block_number", None),
        initiator_participant_id=getattr(dto, "initiator_participant_id", None),
        external_reference_id=getattr(dto, "external_reference_id", None),
        trails=trails,
        created_at=safe_datetime(getattr(dto, "created_at", None)),
    )


def pledge_withdrawals_from_dto(dto_list: Any) -> List[PledgeWithdrawal]:
    """
    Convert list of pledge withdrawal DTOs to domain models.

    Args:
        dto_list: List of OpenAPI TgvalidatordTnPledgeWithdrawal DTOs.

    Returns:
        List of PledgeWithdrawal domain models.
    """
    if not dto_list:
        return []

    result: List[PledgeWithdrawal] = []
    for dto in dto_list:
        withdrawal = pledge_withdrawal_from_dto(dto)
        if withdrawal:
            result.append(withdrawal)

    return result
