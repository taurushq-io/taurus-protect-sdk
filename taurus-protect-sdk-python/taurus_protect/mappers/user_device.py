"""User device mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, Optional

from taurus_protect.mappers._base import safe_datetime, safe_string
from taurus_protect.models.user_device import UserDevicePairing, UserDevicePairingInfo

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def user_device_pairing_from_dto(dto: Any) -> Optional[UserDevicePairing]:
    """
    Convert OpenAPI TgvalidatordCreateUserDevicePairingReply to domain UserDevicePairing.

    Args:
        dto: OpenAPI user device pairing reply DTO.

    Returns:
        Domain UserDevicePairing model or None if dto is None.
    """
    if dto is None:
        return None

    return UserDevicePairing(
        pairing_id=safe_string(getattr(dto, "pairing_id", None)),
        status=safe_string(getattr(dto, "status", None)),
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        expires_at=safe_datetime(getattr(dto, "expiration_date", None)),
    )


def user_device_pairing_info_from_dto(dto: Any) -> Optional[UserDevicePairingInfo]:
    """
    Convert OpenAPI TgvalidatordUserDevicePairingInfo to domain UserDevicePairingInfo.

    Args:
        dto: OpenAPI user device pairing info DTO.

    Returns:
        Domain UserDevicePairingInfo model or None if dto is None.
    """
    if dto is None:
        return None

    return UserDevicePairingInfo(
        pairing_id=safe_string(getattr(dto, "pairing_id", None)),
        user_id=safe_string(getattr(dto, "user_id", None)),
        status=safe_string(getattr(dto, "status", None)),
        device_name=safe_string(getattr(dto, "device_name", None)),
        device_type=safe_string(getattr(dto, "device_type", None)),
        encryption_key=safe_string(getattr(dto, "encryption_key", None)),
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        expires_at=safe_datetime(getattr(dto, "expiration_date", None)),
    )
