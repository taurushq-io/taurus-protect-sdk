"""Address mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import (
    safe_bool,
    safe_datetime,
    safe_list,
    safe_string,
)
from taurus_protect.mappers.wallet import balance_from_dto
from taurus_protect.models.address import Address, AddressAttribute

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def address_from_dto(dto: Any) -> Optional[Address]:
    """
    Convert OpenAPI TgvalidatordAddress to domain Address.

    Args:
        dto: OpenAPI address DTO (TgvalidatordAddress).

    Returns:
        Domain Address model or None if dto is None.
    """
    if dto is None:
        return None

    # Extract balance if present
    balance = None
    dto_balance = getattr(dto, "balance", None)
    if dto_balance is not None:
        balance = balance_from_dto(dto_balance)

    # Extract attributes if present
    attributes: List[AddressAttribute] = []
    dto_attributes = getattr(dto, "attributes", None)
    if dto_attributes is not None:
        attributes = [address_attribute_from_dto(attr) for attr in dto_attributes]

    # Extract linked whitelisted address IDs
    linked_ids = safe_list(getattr(dto, "linked_whitelisted_address_ids", None))

    return Address(
        id=safe_string(getattr(dto, "id", None)),
        wallet_id=safe_string(getattr(dto, "wallet_id", None)),
        address=safe_string(getattr(dto, "address", None)),
        alternate_address=getattr(dto, "alternate_address", None),
        label=getattr(dto, "label", None),
        comment=getattr(dto, "comment", None),
        currency=safe_string(getattr(dto, "currency", None)),
        customer_id=getattr(dto, "customer_id", None),
        external_address_id=getattr(dto, "external_address_id", None),
        address_path=getattr(dto, "address_path", None),
        address_index=getattr(dto, "address_index", None),
        nonce=getattr(dto, "nonce", None),
        status=getattr(dto, "status", None),
        balance=balance,
        signature=getattr(dto, "signature", None),
        disabled=safe_bool(getattr(dto, "disabled", None)),
        can_use_all_funds=safe_bool(getattr(dto, "can_use_all_funds", None)),
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
        attributes=attributes,
        linked_whitelisted_address_ids=linked_ids,
    )


def addresses_from_dto(dtos: Optional[List[Any]]) -> List[Address]:
    """
    Convert list of OpenAPI address DTOs to domain Addresses.

    Args:
        dtos: List of OpenAPI address DTOs.

    Returns:
        List of domain Address models.
    """
    if dtos is None:
        return []
    return [a for dto in dtos if (a := address_from_dto(dto)) is not None]


def address_attribute_from_dto(dto: Any) -> AddressAttribute:
    """
    Convert OpenAPI TgvalidatordAddressAttribute to domain AddressAttribute.

    Args:
        dto: OpenAPI address attribute DTO.

    Returns:
        Domain AddressAttribute model.
    """
    return AddressAttribute(
        id=safe_string(getattr(dto, "id", None)),
        key=safe_string(getattr(dto, "key", None)),
        value=safe_string(getattr(dto, "value", None)),
    )
