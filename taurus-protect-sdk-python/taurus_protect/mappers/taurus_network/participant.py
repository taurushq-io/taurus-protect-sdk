"""Mapper functions for Taurus Network participant DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import (
    safe_bool,
    safe_datetime,
    safe_int,
    safe_list,
    safe_string,
)
from taurus_protect.models.taurus_network.participant import (
    MyParticipant,
    Participant,
    ParticipantAttribute,
    ParticipantSettings,
)

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def participant_attribute_from_dto(dto: Any) -> Optional[ParticipantAttribute]:
    """
    Convert participant attribute DTO to domain model.

    Args:
        dto: The OpenAPI participant attribute DTO.

    Returns:
        ParticipantAttribute or None if DTO is None.
    """
    if dto is None:
        return None

    return ParticipantAttribute(
        id=safe_string(getattr(dto, "id", None)),
        key=safe_string(getattr(dto, "key", None)),
        value=safe_string(getattr(dto, "value", None)),
        content_type=safe_string(getattr(dto, "content_type", None)),
        attribute_type=safe_string(getattr(dto, "type", None)),
        subtype=safe_string(getattr(dto, "subtype", None)),
        shared=safe_bool(getattr(dto, "shared", None)),
        created_at=safe_datetime(getattr(dto, "created_at", None)),
        updated_at=safe_datetime(getattr(dto, "updated_at", None)),
    )


def participant_from_dto(dto: Any) -> Optional[Participant]:
    """
    Convert participant DTO to domain model.

    Args:
        dto: The OpenAPI TgvalidatordTnParticipant DTO.

    Returns:
        Participant or None if DTO is None.
    """
    if dto is None:
        return None

    # Parse attributes if present
    attributes: List[ParticipantAttribute] = []
    raw_attributes = getattr(dto, "attributes", None)
    if raw_attributes:
        for attr_dto in raw_attributes:
            attr = participant_attribute_from_dto(attr_dto)
            if attr:
                attributes.append(attr)

    return Participant(
        id=safe_string(getattr(dto, "id", None)),
        name=safe_string(getattr(dto, "name", None)),
        legal_address=safe_string(getattr(dto, "legal_address", None)),
        country=safe_string(getattr(dto, "country", None)),
        public_key=safe_string(getattr(dto, "public_key", None)),
        shield=safe_string(getattr(dto, "shield", None)),
        status=safe_string(getattr(dto, "status", None)),
        public_subname=safe_string(getattr(dto, "public_subname", None)),
        legal_entity_identifier=safe_string(getattr(dto, "legal_entity_identifier", None)),
        owned_shared_addresses_count=safe_int(getattr(dto, "owned_shared_addresses_count", None)),
        targeted_shared_addresses_count=safe_int(
            getattr(dto, "targeted_shared_addresses_count", None)
        ),
        outgoing_total_pledges_valuation=safe_string(
            getattr(dto, "outgoing_total_pledges_valuation_base_currency", None)
        ),
        incoming_total_pledges_valuation=safe_string(
            getattr(dto, "incoming_total_pledges_valuation_base_currency", None)
        ),
        attributes=attributes,
        origin_registration_date=safe_datetime(getattr(dto, "origin_registration_date", None)),
        origin_deletion_date=safe_datetime(getattr(dto, "origin_deletion_date", None)),
        created_at=safe_datetime(getattr(dto, "created_at", None)),
        updated_at=safe_datetime(getattr(dto, "updated_at", None)),
    )


def participants_from_dto(dto_list: Any) -> List[Participant]:
    """
    Convert list of participant DTOs to domain models.

    Args:
        dto_list: List of OpenAPI TgvalidatordTnParticipant DTOs.

    Returns:
        List of Participant domain models.
    """
    if not dto_list:
        return []

    result: List[Participant] = []
    for dto in dto_list:
        participant = participant_from_dto(dto)
        if participant:
            result.append(participant)

    return result


def participant_settings_from_dto(dto: Any) -> Optional[ParticipantSettings]:
    """
    Convert participant settings DTO to domain model.

    Args:
        dto: The OpenAPI TgvalidatordTnParticipantSettings DTO.

    Returns:
        ParticipantSettings or None if DTO is None.
    """
    if dto is None:
        return None

    return ParticipantSettings(
        status=safe_string(getattr(dto, "status", None)),
        interacting_allowed_countries=safe_list(
            getattr(dto, "interacting_allowed_countries", None)
        ),
        terms_and_conditions_accepted_at=safe_datetime(
            getattr(dto, "terms_and_conditions_accepted_at", None)
        ),
    )


def my_participant_from_dto(dto: Any) -> Optional[MyParticipant]:
    """
    Convert GetMyParticipantReply DTO to domain model.

    Args:
        dto: The OpenAPI TgvalidatordGetMyParticipantReply DTO.

    Returns:
        MyParticipant or None if DTO is None.
    """
    if dto is None:
        return None

    participant = participant_from_dto(getattr(dto, "result", None))
    settings = participant_settings_from_dto(getattr(dto, "settings", None))

    return MyParticipant(
        participant=participant,
        settings=settings,
    )
