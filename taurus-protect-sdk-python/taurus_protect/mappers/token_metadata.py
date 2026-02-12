"""Token metadata mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, Optional

from taurus_protect.mappers._base import safe_string
from taurus_protect.models.token_metadata import (
    CryptoPunkMetadata,
    FATokenMetadata,
    TokenMetadata,
)

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def token_metadata_from_dto(dto: Any) -> Optional[TokenMetadata]:
    """
    Convert OpenAPI TgvalidatordERCTokenMetadata to domain TokenMetadata.

    Args:
        dto: OpenAPI token metadata DTO.

    Returns:
        Domain TokenMetadata model or None if dto is None.
    """
    if dto is None:
        return None

    return TokenMetadata(
        name=safe_string(getattr(dto, "name", None)),
        description=safe_string(getattr(dto, "description", None)),
        decimals=safe_string(getattr(dto, "decimals", None)),
        uri=safe_string(getattr(dto, "uri", None)),
        data_type=safe_string(getattr(dto, "data_type", None)),
        base64_data=safe_string(getattr(dto, "base64_data", None)),
    )


def fa_token_metadata_from_dto(dto: Any) -> Optional[FATokenMetadata]:
    """
    Convert OpenAPI TgvalidatordFATokenMetadata to domain FATokenMetadata.

    Args:
        dto: OpenAPI FA token metadata DTO.

    Returns:
        Domain FATokenMetadata model or None if dto is None.
    """
    if dto is None:
        return None

    return FATokenMetadata(
        name=safe_string(getattr(dto, "name", None)),
        symbol=safe_string(getattr(dto, "symbol", None)),
        decimals=safe_string(getattr(dto, "decimals", None)),
        description=safe_string(getattr(dto, "description", None)),
        uri=safe_string(getattr(dto, "uri", None)),
        data_type=safe_string(getattr(dto, "data_type", None)),
        base64_data=safe_string(getattr(dto, "base64_data", None)),
    )


def crypto_punk_metadata_from_dto(dto: Any) -> Optional[CryptoPunkMetadata]:
    """
    Convert OpenAPI TgvalidatordCryptoPunkMetadata to domain CryptoPunkMetadata.

    Args:
        dto: OpenAPI CryptoPunk metadata DTO.

    Returns:
        Domain CryptoPunkMetadata model or None if dto is None.
    """
    if dto is None:
        return None

    return CryptoPunkMetadata(
        punk_index=safe_string(getattr(dto, "punk_index", None)),
        image_url=safe_string(getattr(dto, "image_url", None)),
        attributes=getattr(dto, "attributes", None),
    )
