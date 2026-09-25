"""Balance service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.currency import (
    asset_balances_from_dto,
    nft_collection_balances_from_dto,
)
from taurus_protect.models.currency import AssetBalance, NFTCollectionBalance
from taurus_protect.models.pagination import CursorPage, cursor_page, cursor_request
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
        >>> balances, page = client.balances.list(page_size=100)
        >>> for balance in balances:
        ...     print(f"{balance.currency}: {balance.balance}")
        >>>
        >>> # List balances for a specific currency
        >>> balances, page = client.balances.list(currency="ETH")
        >>>
        >>> # List NFT collection balances
        >>> nft_balances, page = client.balances.list_nft_collections(
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
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        *,
        token_id: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[AssetBalance], CursorPage]:
        """
        List the tenant's balances, one page at a time, optionally by currency.

        Each asset is identified by a full triplet of attributes (blockchain,
        contract address, and token ID).

        Args:
            currency: Filter by currency ID or symbol. If None, returns all balances.
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            token_id: Filter by token ID.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (balances, page); the page carries the server's total.

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            APIError: If API request fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            resp = self._api.wallet_service_get_balances(
                currency=currency,
                token_id=token_id,
                **req.query_params("request_cursor"),
            )

            balances = asset_balances_from_dto(resp.balances or [])
            page = cursor_page(req.page_size, resp.cursor, total=resp.total, has_total=True)
            return balances, page
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_nft_collections(
        self,
        blockchain: Optional[str] = None,
        network: Optional[str] = None,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        *,
        query: Optional[str] = None,
        only_positive_balance: Optional[bool] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[NFTCollectionBalance], CursorPage]:
        """
        List the tenant's NFT collection balances, one page at a time.

        Args:
            blockchain: Filter by blockchain (e.g., "ETH").
            network: Filter by network (e.g., "mainnet").
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            query: Search query.
            only_positive_balance: Only collections with a non-zero balance.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (NFT collection balances, page).

        Raises:
            ValueError: If paging options are invalid.
            APIError: If API request fails.
        """
        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        try:
            resp = self._api.wallet_service_get_nft_collection_balances(
                blockchain=blockchain,
                network=network,
                query=query,
                only_positive_balance=only_positive_balance,
                **req.query_params(),
            )

            balances = nft_collection_balances_from_dto(resp.balances or [])
            return balances, cursor_page(req.page_size, resp.cursor)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
