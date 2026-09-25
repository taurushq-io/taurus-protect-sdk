"""Wallet service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers.wallet import (
    asset_balance_from_dto,
    balance_history_point_from_dto,
    wallet_from_create_dto,
    wallet_from_dto,
    wallets_from_dto,
)
from taurus_protect.models.balance import AssetBalance, BalanceHistoryPoint
from taurus_protect.models.pagination import (
    REPLY_OFFSET,
    CursorPage,
    Pagination,
    cursor_page,
    offset_pagination,
    offset_query,
    resolve_offset,
    resolve_page_size,
)
from taurus_protect.models.wallet import CreateWalletRequest, ListWalletsOptions, Wallet
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class WalletService(BaseService):
    """
    Service for wallet management operations.

    Provides methods to list, get, and create wallets, as well as
    manage wallet attributes and retrieve balance history.

    Example:
        >>> # Walk every wallet, one page at a time
        >>> offset = 0
        >>> while True:
        ...     wallets, pagination = client.wallets.list(limit=100, offset=offset)
        ...     for wallet in wallets:
        ...         print(f"{wallet.name}: {wallet.currency}")
        ...     if not pagination.has_more:
        ...         break
        ...     offset = pagination.next_offset
        >>>
        >>> # Get single wallet
        >>> wallet = client.wallets.get(123)
        >>> print(f"Balance: {wallet.balance.total_confirmed}")
        >>>
        >>> # Create wallet
        >>> request = CreateWalletRequest(
        ...     blockchain="ETH",
        ...     network="mainnet",
        ...     name="Trading Wallet",
        ... )
        >>> wallet = client.wallets.create(request)
    """

    def __init__(self, api_client: Any, wallets_api: Any) -> None:
        """
        Initialize wallet service.

        Args:
            api_client: The OpenAPI client instance.
            wallets_api: The WalletsAPI service from OpenAPI client.
        """
        super().__init__(api_client)
        self._wallets_api = wallets_api

    def get(self, wallet_id: int) -> Wallet:
        """
        Get a wallet by ID.

        Args:
            wallet_id: The wallet ID to retrieve.

        Returns:
            The wallet.

        Raises:
            ValueError: If wallet_id is invalid.
            NotFoundError: If wallet not found.
            APIError: If API request fails.
        """
        if wallet_id <= 0:
            raise ValueError("wallet_id must be positive")

        try:
            resp = self._wallets_api.wallet_service_get_wallet_v2(str(wallet_id))

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Wallet {wallet_id} not found")

            wallet = wallet_from_dto(result)
            if wallet is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Wallet {wallet_id} not found")

            return wallet
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(
        self,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        exclude_disabled: Optional[bool] = None,
    ) -> Tuple[List[Wallet], Pagination]:
        """
        List wallets, one page at a time.

        Args:
            limit: Page size (default 20, max 100).
            offset: Number of wallets to skip; pass ``pagination.next_offset`` to continue.
            exclude_disabled: True hides every disabled wallet; unset hides only
                wallets whose currency is disabled.

        Returns:
            Tuple of (wallets, pagination).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        return self.list_with_options(
            ListWalletsOptions(limit=limit, offset=offset, exclude_disabled=exclude_disabled)
        )

    def list_with_options(
        self,
        options: Optional[ListWalletsOptions] = None,
    ) -> Tuple[List[Wallet], Pagination]:
        """
        List wallets with every filter the endpoint supports.

        Args:
            options: Filters and page window; every field reaches the wire.

        Returns:
            Tuple of (wallets, pagination).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        opts = options or ListWalletsOptions()
        limit = resolve_page_size(opts.limit, "limit")
        offset = resolve_offset(opts.offset)

        try:
            resp = self._wallets_api.wallet_service_get_wallets_v2(
                currencies=[opts.currency] if opts.currency else None,
                query=opts.query,
                name=opts.name,
                sort_order=opts.sort_order,
                exclude_disabled=opts.exclude_disabled,
                tag_ids=opts.tag_ids,
                only_positive_balance=opts.only_positive_balance,
                blockchain=opts.blockchain,
                network=opts.network,
                ids=opts.ids,
                **offset_query(limit, offset),
            )

            rows = resp.result or []
            wallets = wallets_from_dto(rows)
            pagination = offset_pagination(
                REPLY_OFFSET,
                limit=limit,
                offset=offset,
                served_rows=len(rows),
                total_items=resp.total_items,
                reply_offset=resp.offset,
                excluded=len(rows) - len(wallets),
            )
            return wallets, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_by_name(
        self,
        name: str,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        exclude_disabled: Optional[bool] = None,
    ) -> Tuple[List[Wallet], Pagination]:
        """
        List wallets whose name matches (case-insensitive, partial), one page at a time.

        Args:
            name: The wallet name to search for.
            limit: Page size (default 20, max 100).
            offset: Number of wallets to skip; pass ``pagination.next_offset`` to continue.
            exclude_disabled: True hides every disabled wallet.

        Returns:
            Tuple of (wallets, pagination).

        Raises:
            ValueError: If name is empty or limit/offset invalid.
            APIError: If API request fails.
        """
        self._validate_required(name, "name")
        return self.list_with_options(
            ListWalletsOptions(
                name=name, limit=limit, offset=offset, exclude_disabled=exclude_disabled
            )
        )

    def create(self, request: CreateWalletRequest) -> Wallet:
        """
        Create a new wallet.

        Args:
            request: Wallet creation parameters.

        Returns:
            The created wallet.

        Raises:
            ValueError: If required fields are missing.
            ValidationError: If request is invalid.
            APIError: If API request fails.
        """
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.blockchain, "blockchain")
        self._validate_required(request.network, "network")
        self._validate_required(request.name, "name")

        return self._create_wallet(
            blockchain=request.blockchain,
            network=request.network,
            name=request.name,
            is_omnibus=request.is_omnibus,
            comment=request.comment or "",
            customer_id=request.customer_id or "",
        )

    def create_wallet(
        self,
        blockchain: str,
        network: str,
        name: str,
        is_omnibus: bool = False,
        comment: str = "",
        customer_id: str = "",
    ) -> Wallet:
        """
        Create a new wallet with explicit parameters.

        Args:
            blockchain: Blockchain type (e.g., "ETH", "BTC", "SOL").
            network: Network identifier (e.g., "mainnet", "testnet").
            name: Human-readable wallet name.
            is_omnibus: Whether this is an omnibus wallet.
            comment: Optional description.
            customer_id: Optional customer identifier.

        Returns:
            The created wallet.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        self._validate_required(blockchain, "blockchain")
        self._validate_required(network, "network")
        self._validate_required(name, "name")

        return self._create_wallet(
            blockchain=blockchain,
            network=network,
            name=name,
            is_omnibus=is_omnibus,
            comment=comment,
            customer_id=customer_id,
        )

    def _create_wallet(
        self,
        blockchain: str,
        network: str,
        name: str,
        is_omnibus: bool,
        comment: str,
        customer_id: str,
    ) -> Wallet:
        """Internal wallet creation implementation."""
        try:
            # Build request - field names depend on generated OpenAPI client
            body = {
                "blockchain": blockchain,
                "network": network,
                "name": name,
                "is_omnibus": is_omnibus,
            }
            if comment:
                body["comment"] = comment
            if customer_id:
                body["customer_id"] = customer_id

            resp = self._wallets_api.wallet_service_create_wallet(body=body)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create wallet: no result returned")

            wallet = wallet_from_create_dto(result)
            if wallet is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create wallet: invalid response")

            return wallet
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def create_attribute(
        self,
        wallet_id: int,
        key: str,
        value: str,
    ) -> None:
        """
        Create an attribute for a wallet.

        Args:
            wallet_id: The wallet ID.
            key: The attribute key.
            value: The attribute value.

        Raises:
            ValueError: If any argument is invalid.
            APIError: If API request fails.
        """
        if wallet_id <= 0:
            raise ValueError("wallet_id must be positive")
        self._validate_required(key, "key")
        self._validate_required(value, "value")

        try:
            body = {
                "attributes": [{"key": key, "value": value}],
            }
            self._wallets_api.wallet_service_create_wallet_attributes(str(wallet_id), body=body)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_balance_history(
        self,
        wallet_id: int,
        interval_hours: int,
    ) -> List[BalanceHistoryPoint]:
        """
        Get wallet balance history.

        Args:
            wallet_id: The wallet ID.
            interval_hours: The interval in hours for balance snapshots.

        Returns:
            List of balance history points.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        if wallet_id <= 0:
            raise ValueError("wallet_id must be positive")
        if interval_hours <= 0:
            raise ValueError("interval_hours must be positive")

        try:
            resp = self._wallets_api.wallet_service_get_wallet_balance_history(
                str(wallet_id),
                str(interval_hours),
            )

            result = getattr(resp, "result", None)
            if result is None:
                return []

            return [p for dto in result if (p := balance_history_point_from_dto(dto)) is not None]
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_tokens(
        self,
        wallet_id: int,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
    ) -> Tuple[List[AssetBalance], CursorPage]:
        """
        List a wallet's token balances, one page at a time.

        Args:
            wallet_id: The wallet ID.
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.

        Returns:
            Tuple of (asset balances, page); the page carries the server's total.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        if wallet_id <= 0:
            raise ValueError("wallet_id must be positive")
        size = resolve_page_size(page_size)

        try:
            resp = self._wallets_api.wallet_service_get_wallet_tokens(
                str(wallet_id),
                limit=str(size),
                cursor=cursor or None,
            )

            balances = [
                b for dto in resp.balances or [] if (b := asset_balance_from_dto(dto)) is not None
            ]
            return balances, cursor_page(size, token=resp.next, total=resp.total, has_total=True)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
