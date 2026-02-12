"""Wallet mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import (
    parse_string_to_int,
    safe_bool,
    safe_datetime,
    safe_int,
    safe_string,
)
from taurus_protect.mappers.currency import currency_from_dto
from taurus_protect.models.balance import Asset, AssetBalance, Balance, BalanceHistoryPoint
from taurus_protect.models.wallet import Wallet, WalletAttribute

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def wallet_from_dto(dto: Any) -> Optional[Wallet]:
    """
    Convert OpenAPI TgvalidatordWalletInfo to domain Wallet.

    Args:
        dto: OpenAPI wallet DTO (TgvalidatordWalletInfo).

    Returns:
        Domain Wallet model or None if dto is None.
    """
    if dto is None:
        return None

    # Extract balance if present
    balance = None
    dto_balance = getattr(dto, "balance", None)
    if dto_balance is not None:
        balance = balance_from_dto(dto_balance)

    # Extract attributes if present
    attributes: List[WalletAttribute] = []
    dto_attributes = getattr(dto, "attributes", None)
    if dto_attributes is not None:
        attributes = [wallet_attribute_from_dto(attr) for attr in dto_attributes]

    # Parse addresses count (API returns string)
    addresses_count = parse_string_to_int(getattr(dto, "addresses_count", None))

    # Convert currency info if present
    currency_info = None
    dto_currency_info = getattr(dto, "currency_info", None)
    if dto_currency_info is not None:
        currency_info = currency_from_dto(dto_currency_info)

    return Wallet(
        id=safe_string(getattr(dto, "id", None)),
        name=safe_string(getattr(dto, "name", None)),
        currency=safe_string(getattr(dto, "currency", None)),
        blockchain=safe_string(getattr(dto, "blockchain", None)),
        network=safe_string(getattr(dto, "network", None)),
        balance=balance,
        is_omnibus=safe_bool(getattr(dto, "is_omnibus", None)),
        disabled=safe_bool(getattr(dto, "disabled", None)),
        comment=getattr(dto, "comment", None),
        customer_id=getattr(dto, "customer_id", None),
        external_wallet_id=getattr(dto, "external_wallet_id", None),
        visibility_group_id=getattr(dto, "visibility_group_id", None),
        account_path=getattr(dto, "account_path", None),
        addresses_count=addresses_count,
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
        attributes=attributes,
        currency_info=currency_info,
    )


def wallet_from_create_dto(dto: Any) -> Optional[Wallet]:
    """
    Convert OpenAPI TgvalidatordWallet (from create response) to domain Wallet.

    Args:
        dto: OpenAPI wallet DTO (TgvalidatordWallet).

    Returns:
        Domain Wallet model or None if dto is None.
    """
    if dto is None:
        return None

    # Extract balance if present
    balance = None
    dto_balance = getattr(dto, "balance", None)
    if dto_balance is not None:
        balance = balance_from_dto(dto_balance)

    # Extract attributes if present
    attributes: List[WalletAttribute] = []
    dto_attributes = getattr(dto, "attributes", None)
    if dto_attributes is not None:
        attributes = [wallet_attribute_from_dto(attr) for attr in dto_attributes]

    # Parse addresses count (API returns string)
    addresses_count = parse_string_to_int(getattr(dto, "addresses_count", None))

    # Convert currency info if present
    currency_info = None
    dto_currency_info = getattr(dto, "currency_info", None)
    if dto_currency_info is not None:
        currency_info = currency_from_dto(dto_currency_info)

    return Wallet(
        id=safe_string(getattr(dto, "id", None)),
        name=safe_string(getattr(dto, "name", None)),
        currency=safe_string(getattr(dto, "currency", None)),
        blockchain=safe_string(getattr(dto, "blockchain", None)),
        network=safe_string(getattr(dto, "network", None)),
        balance=balance,
        is_omnibus=safe_bool(getattr(dto, "is_omnibus", None)),
        disabled=safe_bool(getattr(dto, "disabled", None)),
        comment=getattr(dto, "comment", None),
        customer_id=getattr(dto, "customer_id", None),
        external_wallet_id=getattr(dto, "external_wallet_id", None),
        visibility_group_id=getattr(dto, "visibility_group_id", None),
        account_path=getattr(dto, "account_path", None),
        addresses_count=addresses_count,
        created_at=safe_datetime(getattr(dto, "creation_date", None)),
        updated_at=safe_datetime(getattr(dto, "update_date", None)),
        attributes=attributes,
        currency_info=currency_info,
    )


def wallets_from_dto(dtos: Optional[List[Any]]) -> List[Wallet]:
    """
    Convert list of OpenAPI wallet DTOs to domain Wallets.

    Args:
        dtos: List of OpenAPI wallet DTOs.

    Returns:
        List of domain Wallet models.
    """
    if dtos is None:
        return []
    return [w for dto in dtos if (w := wallet_from_dto(dto)) is not None]


def wallet_attribute_from_dto(dto: Any) -> WalletAttribute:
    """
    Convert OpenAPI TgvalidatordWalletAttribute to domain WalletAttribute.

    Args:
        dto: OpenAPI wallet attribute DTO.

    Returns:
        Domain WalletAttribute model.
    """
    return WalletAttribute(
        id=safe_string(getattr(dto, "id", None)),
        key=safe_string(getattr(dto, "key", None)),
        value=safe_string(getattr(dto, "value", None)),
        content_type=getattr(dto, "content_type", None),
        owner=getattr(dto, "owner", None),
        type=getattr(dto, "type", None),
        subtype=getattr(dto, "subtype", None),
        is_file=safe_bool(
            getattr(dto, "isfile", None)
            if getattr(dto, "isfile", None) is not None
            else getattr(dto, "is_file", None)
        ),
    )


def balance_from_dto(dto: Any) -> Optional[Balance]:
    """
    Convert OpenAPI TgvalidatordBalance to domain Balance.

    Args:
        dto: OpenAPI balance DTO.

    Returns:
        Domain Balance model or None if dto is None.
    """
    if dto is None:
        return None

    return Balance(
        total_confirmed=safe_string(getattr(dto, "total_confirmed", None)),
        total_unconfirmed=safe_string(getattr(dto, "total_unconfirmed", None)),
        available_confirmed=safe_string(getattr(dto, "available_confirmed", None)),
        available_unconfirmed=safe_string(getattr(dto, "available_unconfirmed", None)),
        reserved_confirmed=safe_string(getattr(dto, "reserved_confirmed", None)),
        reserved_unconfirmed=safe_string(getattr(dto, "reserved_unconfirmed", None)),
    )


def balance_history_point_from_dto(dto: Any) -> Optional[BalanceHistoryPoint]:
    """
    Convert OpenAPI TgvalidatordBalanceHistoryPoint to domain BalanceHistoryPoint.

    Args:
        dto: OpenAPI balance history point DTO.

    Returns:
        Domain BalanceHistoryPoint model or None if dto is None.
    """
    if dto is None:
        return None

    return BalanceHistoryPoint(
        timestamp=safe_datetime(getattr(dto, "timestamp", None)),
        total_confirmed=safe_string(getattr(dto, "total_confirmed", None)),
        total_unconfirmed=safe_string(getattr(dto, "total_unconfirmed", None)),
        available_confirmed=safe_string(getattr(dto, "available_confirmed", None)),
        available_unconfirmed=safe_string(getattr(dto, "available_unconfirmed", None)),
    )


def asset_from_dto(dto: Any) -> Optional[Asset]:
    """
    Convert OpenAPI asset DTO to domain Asset.

    Args:
        dto: OpenAPI asset DTO.

    Returns:
        Domain Asset model or None if dto is None.
    """
    if dto is None:
        return None

    return Asset(
        id=safe_string(getattr(dto, "id", None)),
        symbol=safe_string(getattr(dto, "symbol", None)),
        name=safe_string(getattr(dto, "name", None)),
        decimals=safe_int(getattr(dto, "decimals", None)),
        blockchain=safe_string(getattr(dto, "blockchain", None)),
    )


def asset_balance_from_dto(dto: Any) -> Optional[AssetBalance]:
    """
    Convert OpenAPI TgvalidatordAssetBalance to domain AssetBalance.

    Args:
        dto: OpenAPI asset balance DTO.

    Returns:
        Domain AssetBalance model or None if dto is None.
    """
    if dto is None:
        return None

    asset = None
    dto_asset = getattr(dto, "asset", None)
    if dto_asset is not None:
        asset = asset_from_dto(dto_asset)

    balance = None
    dto_balance = getattr(dto, "balance", None)
    if dto_balance is not None:
        balance = balance_from_dto(dto_balance)

    return AssetBalance(
        asset=asset,
        balance=balance,
    )
