"""Balance service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.currency import (
    asset_balances_from_dto,
    nft_collection_balances_from_dto,
)
from taurus_protect.models.currency import AssetBalance, NFTCollectionBalance
from taurus_protect.models.pagination import Pagination
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.balances_api import BalancesApi


class BalanceService(BaseService):
    """
    Service for retrieving balance information.

    Provides methods to list balances for all assets and NFT collections
    across the tenant.

    Example:
        >>> # List all balances
        >>> balances, pagination = client.balances.list(limit=50)
        >>> for balance in balances:
        ...     print(f"{balance.currency}: {balance.balance}")
        >>>
        >>> # List balances for a specific currency
        >>> balances, pagination = client.balances.list(currency="ETH")
        >>>
        >>> # List NFT collection balances
        >>> nft_balances, pagination = client.balances.list_nft_collections(
        ...     blockchain="ETH",
        ...     network="mainnet",
        ... )
    """

    def __init__(
        self,
        api_client: Any,
        balances_api: "BalancesApi",
    ) -> None:
        """
        Initialize the balance service.

        Args:
            api_client: The OpenAPI client instance.
            balances_api: The BalancesApi instance from OpenAPI client.
        """
        super().__init__(api_client)
        self._api = balances_api

    def list(
        self,
        currency: Optional[str] = None,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[AssetBalance], Optional[Pagination]]:
        """
        Get all balances for the tenant, optionally filtered by currency.

        Each asset is identified by a full triplet of attributes (blockchain,
        contract address, and token ID).

        Args:
            currency: Filter by currency ID or symbol. If None, returns all balances.
            limit: Maximum number of balances to return (must be positive).
            offset: Number of balances to skip (must be non-negative).

        Returns:
            Tuple of (balances list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._api.wallet_service_get_balances(
                currency=currency,
                limit=str(limit),
                cursor=None,
                token_id=None,
                request_cursor_current_page=None,
                request_cursor_page_request=None,
                request_cursor_page_size=str(limit),
            )

            result = getattr(resp, "balances", None) or getattr(resp, "result", None)
            balances = asset_balances_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=offset,
                limit=limit,
            )

            return balances, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_nft_collections(
        self,
        blockchain: str,
        network: str,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[NFTCollectionBalance], Optional[Pagination]]:
        """
        Get NFT collection balances for the tenant.

        Args:
            blockchain: Blockchain to filter by (e.g., "ETH").
            network: Network to filter by (e.g., "mainnet").
            limit: Maximum number of collections to return (must be positive).
            offset: Number of collections to skip (must be non-negative).

        Returns:
            Tuple of (NFT collection balances list, pagination info).

        Raises:
            ValueError: If required arguments are missing or invalid.
            APIError: If API request fails.
        """
        self._validate_required(blockchain, "blockchain")
        self._validate_required(network, "network")
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._api.wallet_service_get_nft_collection_balances(
                blockchain=blockchain,
                network=network,
                query=None,
                cursor_current_page=None,
                cursor_page_request=None,
                cursor_page_size=str(limit),
                only_positive_balance=None,
            )

            result = getattr(resp, "collections", None) or getattr(resp, "result", None)
            balances = nft_collection_balances_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=offset,
                limit=limit,
            )

            return balances, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
