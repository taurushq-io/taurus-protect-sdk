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
from taurus_protect.models.pagination import Pagination
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
        >>> # List wallets
        >>> wallets, pagination = client.wallets.list(limit=50, offset=0)
        >>> for wallet in wallets:
        ...     print(f"{wallet.name}: {wallet.currency}")
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
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Wallet], Optional[Pagination]]:
        """
        List wallets with pagination.

        Args:
            limit: Maximum number of wallets to return (must be positive).
            offset: Number of wallets to skip (must be non-negative).

        Returns:
            Tuple of (wallets list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._wallets_api.wallet_service_get_wallets_v2(
                currencies=None,
                query=None,
                limit=str(limit),
                offset=str(offset),
                name=None,
                sort_order=None,
                exclude_disabled=None,
                tag_ids=None,
                only_positive_balance=None,
                blockchain=None,
                network=None,
                ids=None,
            )

            result = getattr(resp, "result", None)
            wallets = wallets_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=getattr(resp, "offset", None),
                limit=limit,
            )

            return wallets, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_with_options(
        self,
        options: Optional[ListWalletsOptions] = None,
    ) -> Tuple[List[Wallet], Optional[Pagination]]:
        """
        List wallets with full filtering options.

        Args:
            options: Optional filtering and pagination options.

        Returns:
            Tuple of (wallets list, pagination info).

        Raises:
            APIError: If API request fails.
        """
        opts = options or ListWalletsOptions()

        try:
            currencies = [opts.currency] if opts.currency else None

            resp = self._wallets_api.wallet_service_get_wallets_v2(
                currencies=currencies,
                query=opts.query,
                limit=str(opts.limit) if opts.limit > 0 else None,
                offset=str(opts.offset) if opts.offset > 0 else None,
                name=None,
                sort_order=None,
                exclude_disabled=opts.exclude_disabled if opts.exclude_disabled else None,
                tag_ids=None,
                only_positive_balance=None,
                blockchain=None,
                network=None,
                ids=None,
            )

            result = getattr(resp, "result", None)
            wallets = wallets_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=getattr(resp, "offset", None),
                limit=opts.limit,
            )

            return wallets, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def get_by_name(
        self,
        name: str,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[Wallet], Optional[Pagination]]:
        """
        Get wallets by name with pagination.

        Args:
            name: The wallet name to search for.
            limit: Maximum number of wallets to return.
            offset: Number of wallets to skip.

        Returns:
            Tuple of (wallets list, pagination info).

        Raises:
            ValueError: If name is empty or limit/offset invalid.
            APIError: If API request fails.
        """
        self._validate_required(name, "name")
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._wallets_api.wallet_service_get_wallets_v2(
                currencies=None,
                query=None,
                limit=str(limit),
                offset=str(offset),
                name=name,
                sort_order=None,
                exclude_disabled=None,
                tag_ids=None,
                only_positive_balance=None,
                blockchain=None,
                network=None,
                ids=None,
            )

            result = getattr(resp, "result", None)
            wallets = wallets_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=getattr(resp, "offset", None),
                limit=limit,
            )

            return wallets, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

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
        limit: int = 50,
    ) -> List[AssetBalance]:
        """
        Get wallet tokens (asset balances).

        Args:
            wallet_id: The wallet ID.
            limit: Maximum number of tokens to return.

        Returns:
            List of asset balances.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        if wallet_id <= 0:
            raise ValueError("wallet_id must be positive")
        if limit <= 0:
            raise ValueError("limit must be positive")

        try:
            resp = self._wallets_api.wallet_service_get_wallet_tokens(
                str(wallet_id),
                str(limit),
                None,  # cursor
            )

            balances = getattr(resp, "balances", None)
            if balances is None:
                return []

            return [b for dto in balances if (b := asset_balance_from_dto(dto)) is not None]
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
