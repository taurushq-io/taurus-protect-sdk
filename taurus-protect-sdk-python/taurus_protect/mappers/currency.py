"""Currency mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import Any, List, Optional

from taurus_protect.mappers._base import safe_bool, safe_int, safe_string
from taurus_protect.models.currency import AssetBalance, Currency, NFTCollectionBalance


def currency_from_dto(dto: Any) -> Optional[Currency]:
    """
    Convert OpenAPI TgvalidatordCurrency to domain Currency.

    Args:
        dto: OpenAPI currency DTO (TgvalidatordCurrency).

    Returns:
        Domain Currency model or None if dto is None.
    """
    if dto is None:
        return None

    return Currency(
        id=safe_string(getattr(dto, "id", None)),
        name=getattr(dto, "name", None),
        symbol=getattr(dto, "symbol", None),
        blockchain=getattr(dto, "blockchain", None),
        network=getattr(dto, "network", None),
        decimals=safe_int(getattr(dto, "decimals", None)),
        logo_url=getattr(dto, "logo_url", None) or getattr(dto, "logo", None),
        enabled=(
            safe_bool(getattr(dto, "enabled", None))
            if getattr(dto, "enabled", None) is not None
            else True
        ),
        is_token=safe_bool(getattr(dto, "is_token", None)),
        contract_address=getattr(dto, "contract_address", None)
        or getattr(dto, "token_contract_address", None),
        display_name=getattr(dto, "display_name", None),
        type=getattr(dto, "type", None),
        coin_type_index=safe_string(getattr(dto, "coin_type_index", None)),
        token_id=safe_string(getattr(dto, "token_id", None)),
        wlca_id=safe_int(getattr(dto, "wlca_id", None)) if getattr(dto, "wlca_id", None) is not None else None,
        is_erc20=safe_bool(getattr(dto, "is_erc20", None)),
        is_fa12=safe_bool(getattr(dto, "is_fa12", None)),
        is_fa20=safe_bool(getattr(dto, "is_fa20", None)),
        is_nft=safe_bool(getattr(dto, "is_nft", None)),
        is_utxo_based=safe_bool(getattr(dto, "is_utxo_based", None)),
        is_account_based=safe_bool(getattr(dto, "is_account_based", None)),
        is_fiat=safe_bool(getattr(dto, "is_fiat", None)),
        has_staking=safe_bool(getattr(dto, "has_staking", None)),
    )


def currencies_from_dto(dtos: Optional[List[Any]]) -> List[Currency]:
    """
    Convert list of OpenAPI currency DTOs to domain Currencies.

    Args:
        dtos: List of OpenAPI currency DTOs.

    Returns:
        List of domain Currency models.
    """
    if dtos is None:
        return []
    return [c for dto in dtos if (c := currency_from_dto(dto)) is not None]


def asset_balance_from_dto(dto: Any) -> Optional[AssetBalance]:
    """
    Convert OpenAPI balance DTO to domain AssetBalance.

    This is used for the balances endpoint response, which returns
    aggregated balances per currency/asset. The DTO has nested structure:
    - dto.asset: TgvalidatordAsset with currency info
    - dto.balance: TgvalidatordBalance with balance fields

    Args:
        dto: OpenAPI balance DTO (TgvalidatordAssetBalance).

    Returns:
        Domain AssetBalance model or None if dto is None.
    """
    if dto is None:
        return None

    # Extract nested asset info
    asset = getattr(dto, "asset", None)
    currency_id = None
    currency = None
    blockchain = None
    network = None
    if asset is not None:
        currency_id = getattr(asset, "currency", None)
        currency = getattr(asset, "currency", None) or getattr(asset, "symbol", None)
        # Try to get blockchain/network from asset's currency_info if available
        currency_info = getattr(asset, "currency_info", None)
        if currency_info is not None:
            blockchain = getattr(currency_info, "blockchain", None)
            network = getattr(currency_info, "network", None)

    # Extract nested balance info
    balance_obj = getattr(dto, "balance", None)
    balance_value = None
    pending_balance = None
    if balance_obj is not None:
        balance_value = safe_string(getattr(balance_obj, "total_confirmed", None))
        pending_balance = getattr(balance_obj, "total_unconfirmed", None)

    return AssetBalance(
        currency_id=currency_id,
        currency=currency,
        blockchain=blockchain,
        network=network,
        balance=balance_value,
        pending_balance=pending_balance,
        wallet_id=getattr(dto, "wallet_id", None),
        address_id=getattr(dto, "address_id", None),
    )


def asset_balances_from_dto(dtos: Optional[List[Any]]) -> List[AssetBalance]:
    """
    Convert list of OpenAPI balance DTOs to domain AssetBalances.

    Args:
        dtos: List of OpenAPI balance DTOs.

    Returns:
        List of domain AssetBalance models.
    """
    if dtos is None:
        return []
    return [b for dto in dtos if (b := asset_balance_from_dto(dto)) is not None]


def nft_collection_balance_from_dto(dto: Any) -> Optional[NFTCollectionBalance]:
    """
    Convert OpenAPI NFT collection balance DTO to domain NFTCollectionBalance.

    Args:
        dto: OpenAPI NFT collection balance DTO.

    Returns:
        Domain NFTCollectionBalance model or None if dto is None.
    """
    if dto is None:
        return None

    return NFTCollectionBalance(
        collection_name=getattr(dto, "collection_name", None) or getattr(dto, "name", None),
        contract_address=getattr(dto, "contract_address", None),
        blockchain=getattr(dto, "blockchain", None),
        network=getattr(dto, "network", None),
        count=safe_int(getattr(dto, "count", None)) or safe_int(getattr(dto, "balance", None)),
    )


def nft_collection_balances_from_dto(dtos: Optional[List[Any]]) -> List[NFTCollectionBalance]:
    """
    Convert list of OpenAPI NFT collection balance DTOs to domain NFTCollectionBalances.

    Args:
        dtos: List of OpenAPI NFT collection balance DTOs.

    Returns:
        List of domain NFTCollectionBalance models.
    """
    if dtos is None:
        return []
    return [b for dto in dtos if (b := nft_collection_balance_from_dto(dto)) is not None]
