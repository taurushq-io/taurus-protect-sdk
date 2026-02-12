"""Asset service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers._base import safe_bool, safe_int, safe_string
from taurus_protect.models.blockchain import Asset
from taurus_protect.models.pagination import Pagination
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def asset_from_dto(dto: Any) -> Optional[Asset]:
    """
    Convert an OpenAPI asset DTO to domain model.

    Args:
        dto: The OpenAPI DTO object.

    Returns:
        Asset model or None if dto is None.
    """
    if dto is None:
        return None

    return Asset(
        id=safe_string(getattr(dto, "id", None) or getattr(dto, "currency_id", None)),
        name=getattr(dto, "name", None),
        symbol=getattr(dto, "symbol", None) or getattr(dto, "currency", None),
        blockchain=getattr(dto, "blockchain", None),
        network=getattr(dto, "network", None),
        decimals=safe_int(getattr(dto, "decimals", 0)),
        logo_url=getattr(dto, "logo_url", None) or getattr(dto, "logoUrl", None),
        enabled=safe_bool(getattr(dto, "enabled", True)),
        is_token=safe_bool(getattr(dto, "is_token", False) or getattr(dto, "isToken", False)),
        contract_address=getattr(dto, "contract_address", None)
        or getattr(dto, "contractAddress", None),
    )


def assets_from_dto(dtos: Any) -> List[Asset]:
    """
    Convert a list of OpenAPI asset DTOs to domain models.

    Args:
        dtos: List of OpenAPI DTO objects.

    Returns:
        List of Asset models.
    """
    if dtos is None:
        return []
    return [a for dto in dtos if (a := asset_from_dto(dto)) is not None]


class AssetService(BaseService):
    """
    Service for asset operations.

    Provides methods to list and retrieve cryptocurrency/token assets
    and their balances across wallets and addresses.

    Example:
        >>> # List assets
        >>> assets, pagination = client.assets.list(limit=50)
        >>> for asset in assets:
        ...     print(f"{asset.symbol}: {asset.name}")
        >>>
        >>> # Get asset by ID
        >>> asset = client.assets.get("BTC")
        >>> print(f"Decimals: {asset.decimals}")
    """

    def __init__(self, api_client: Any, assets_api: Any) -> None:
        """
        Initialize asset service.

        Args:
            api_client: The OpenAPI client instance.
            assets_api: The AssetsAPI service from OpenAPI client.
        """
        super().__init__(api_client)
        self._assets_api = assets_api

    def list(
        self,
        currency: str = "ETH",
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Asset], Optional[Pagination]]:
        """
        List asset wallets for a given currency.

        Args:
            currency: The currency code to list (e.g., "ETH", "BTC", "USDC").
            limit: Maximum number of assets to return (must be positive).
            offset: Number of assets to skip (must be non-negative).

        Returns:
            Tuple of (assets list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid, or currency is empty.
            APIError: If API request fails.
        """
        self._validate_required(currency, "currency")
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_asset import TgvalidatordAsset
            from taurus_protect._internal.openapi.models.tgvalidatord_get_asset_wallets_request import (
                TgvalidatordGetAssetWalletsRequest,
            )

            asset = TgvalidatordAsset(currency=currency)
            body = TgvalidatordGetAssetWalletsRequest(
                asset=asset,
                limit=str(limit),
            )
            resp = self._assets_api.wallet_service_get_asset_wallets(body=body)

            result = getattr(resp, "result", None) or getattr(resp, "wallets", None)
            assets = assets_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None) or getattr(resp, "totalItems", None),
                offset=getattr(resp, "offset", None),
                limit=limit,
            )

            return assets, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, asset_id: str) -> Asset:
        """
        Get an asset by ID.

        Args:
            asset_id: The asset ID to retrieve (e.g., "BTC", "ETH").

        Returns:
            The asset.

        Raises:
            ValueError: If asset_id is invalid.
            NotFoundError: If asset not found.
            APIError: If API request fails.
        """
        self._validate_required(asset_id, "asset_id")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_asset import TgvalidatordAsset
            from taurus_protect._internal.openapi.models.tgvalidatord_get_asset_addresses_request import (
                TgvalidatordGetAssetAddressesRequest,
            )

            # Get asset addresses to retrieve asset details
            asset = TgvalidatordAsset(currency=asset_id)
            body = TgvalidatordGetAssetAddressesRequest(
                asset=asset,
                limit="1",
            )
            resp = self._assets_api.wallet_service_get_asset_addresses(body=body)

            # Extract asset info from response
            result = getattr(resp, "result", None) or getattr(resp, "addresses", None)
            if result and len(result) > 0:
                # Extract asset info from first address balance
                first_item = result[0]
                return Asset(
                    id=asset_id,
                    symbol=getattr(first_item, "currency", None) or asset_id,
                    blockchain=getattr(first_item, "blockchain", None),
                    network=getattr(first_item, "network", None),
                    decimals=safe_int(getattr(first_item, "decimals", 0)),
                    enabled=True,
                )

            # If no results, return basic asset info
            return Asset(
                id=asset_id,
                symbol=asset_id,
                enabled=True,
            )
        except Exception as e:
            from taurus_protect.errors import APIError, NotFoundError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_wallets(
        self,
        currency: str,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Any], Optional[Pagination]]:
        """
        Get wallet balances for a specific asset.

        Args:
            currency: The currency/asset ID.
            limit: Maximum number of results to return.
            offset: Number of results to skip.

        Returns:
            Tuple of (wallet balances list, pagination info).

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(currency, "currency")
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_asset import TgvalidatordAsset
            from taurus_protect._internal.openapi.models.tgvalidatord_get_asset_wallets_request import (
                TgvalidatordGetAssetWalletsRequest,
            )

            asset = TgvalidatordAsset(currency=currency)
            body = TgvalidatordGetAssetWalletsRequest(
                asset=asset,
                limit=str(limit),
            )
            resp = self._assets_api.wallet_service_get_asset_wallets(body=body)

            result = getattr(resp, "result", None) or getattr(resp, "wallets", [])

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None) or getattr(resp, "totalItems", None),
                offset=getattr(resp, "offset", None),
                limit=limit,
            )

            return result or [], pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_addresses(
        self,
        currency: str,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Any], Optional[Pagination]]:
        """
        Get address balances for a specific asset.

        Args:
            currency: The currency/asset ID.
            limit: Maximum number of results to return.
            offset: Number of results to skip.

        Returns:
            Tuple of (address balances list, pagination info).

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(currency, "currency")
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_asset import TgvalidatordAsset
            from taurus_protect._internal.openapi.models.tgvalidatord_get_asset_addresses_request import (
                TgvalidatordGetAssetAddressesRequest,
            )

            asset = TgvalidatordAsset(currency=currency)
            body = TgvalidatordGetAssetAddressesRequest(
                asset=asset,
                limit=str(limit),
            )
            resp = self._assets_api.wallet_service_get_asset_addresses(body=body)

            result = getattr(resp, "result", None) or getattr(resp, "addresses", [])

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None) or getattr(resp, "totalItems", None),
                offset=getattr(resp, "offset", None),
                limit=limit,
            )

            return result or [], pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
