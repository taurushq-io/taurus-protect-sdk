"""Exchange service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers._base import safe_bool, safe_datetime, safe_string
from taurus_protect.models.blockchain import Exchange
from taurus_protect.models.pagination import Pagination
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def exchange_from_dto(dto: Any) -> Optional[Exchange]:
    """
    Convert an OpenAPI exchange DTO to domain model.

    Args:
        dto: The OpenAPI DTO object.

    Returns:
        Exchange model or None if dto is None.
    """
    if dto is None:
        return None

    return Exchange(
        id=safe_string(getattr(dto, "id", None)),
        name=getattr(dto, "name", None),
        exchange_label=getattr(dto, "exchange_label", None) or getattr(dto, "exchangeLabel", None),
        status=getattr(dto, "status", None),
        currency_id=getattr(dto, "currency_id", None) or getattr(dto, "currencyId", None),
        currency=getattr(dto, "currency", None),
        balance=getattr(dto, "balance", None),
        pending_balance=getattr(dto, "pending_balance", None)
        or getattr(dto, "pendingBalance", None),
        enabled=safe_bool(getattr(dto, "enabled", True)),
        created_at=safe_datetime(
            getattr(dto, "created_at", None) or getattr(dto, "createdAt", None)
        ),
        updated_at=safe_datetime(
            getattr(dto, "updated_at", None) or getattr(dto, "updatedAt", None)
        ),
    )


def exchanges_from_dto(dtos: Any) -> List[Exchange]:
    """
    Convert a list of OpenAPI exchange DTOs to domain models.

    Args:
        dtos: List of OpenAPI DTO objects.

    Returns:
        List of Exchange models.
    """
    if dtos is None:
        return []
    return [e for dto in dtos if (e := exchange_from_dto(dto)) is not None]


class ExchangeService(BaseService):
    """
    Service for exchange account operations.

    Provides methods to list and retrieve exchange accounts
    connected to Taurus-PROTECT.

    Example:
        >>> # List exchange accounts
        >>> exchanges, pagination = client.exchanges.list(limit=50)
        >>> for exchange in exchanges:
        ...     print(f"{exchange.name} ({exchange.exchange_label}): {exchange.balance}")
        >>>
        >>> # Get specific exchange account
        >>> exchange = client.exchanges.get("123")
        >>> print(f"Status: {exchange.status}")
    """

    def __init__(self, api_client: Any, exchange_api: Any) -> None:
        """
        Initialize exchange service.

        Args:
            api_client: The OpenAPI client instance.
            exchange_api: The ExchangeAPI service from OpenAPI client.
        """
        super().__init__(api_client)
        self._exchange_api = exchange_api

    def list(
        self,
        limit: int = 50,
        offset: int = 0,
        currency_id: Optional[str] = None,
        exchange_label: Optional[str] = None,
        status: Optional[str] = None,
        only_positive_balance: bool = False,
    ) -> Tuple[List[Exchange], Optional[Pagination]]:
        """
        List exchange accounts with pagination.

        Args:
            limit: Maximum number of exchanges to return (must be positive).
            offset: Number of exchanges to skip (must be non-negative).
            currency_id: Optional filter by currency ID.
            exchange_label: Optional filter by exchange label.
            status: Optional filter by status.
            only_positive_balance: Whether to exclude zero-balance accounts.

        Returns:
            Tuple of (exchanges list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._exchange_api.exchange_service_get_exchanges(
                currency_id=currency_id,
                exchange_label=exchange_label,
                status=status,
                only_positive_balance=only_positive_balance if only_positive_balance else None,
                cursor_page_size=str(limit),
                cursor_page_request="FIRST" if offset == 0 else None,
            )

            result = (
                getattr(resp, "result", None)
                or getattr(resp, "exchange_accounts", None)
                or getattr(resp, "exchangeAccounts", None)
            )
            exchanges = exchanges_from_dto(result) if result else []

            # Extract pagination from cursor-based response
            cursor = getattr(resp, "cursor", None)
            total_items = None
            if cursor:
                total_items = getattr(cursor, "total_items", None) or getattr(
                    cursor, "totalItems", None
                )

            pagination = self._extract_pagination(
                total_items=total_items,
                offset=offset,
                limit=limit,
            )

            return exchanges, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get(self, exchange_id: str) -> Exchange:
        """
        Get an exchange account by ID.

        Args:
            exchange_id: The exchange account ID to retrieve.

        Returns:
            The exchange account.

        Raises:
            ValueError: If exchange_id is invalid.
            NotFoundError: If exchange account not found.
            APIError: If API request fails.
        """
        self._validate_required(exchange_id, "exchange_id")

        try:
            resp = self._exchange_api.exchange_service_get_exchange(id=exchange_id)

            result = (
                getattr(resp, "result", None)
                or getattr(resp, "exchange_account", None)
                or getattr(resp, "exchangeAccount", None)
            )
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Exchange account {exchange_id} not found")

            exchange = exchange_from_dto(result)
            if exchange is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Exchange account {exchange_id} not found")

            return exchange
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_counterparties(self) -> List[Any]:
        """
        List exchange counterparties with their exposure limits.

        Returns:
            List of exchange counterparties.

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._exchange_api.exchange_service_get_exchange_counterparties()

            result = getattr(resp, "result", None) or getattr(resp, "exchanges", None)
            return result or []
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def get_withdrawal_fee(
        self,
        exchange_id: str,
        to_address_id: Optional[str] = None,
        amount: Optional[str] = None,
    ) -> Any:
        """
        Get withdrawal fees for an exchange account.

        Args:
            exchange_id: The exchange account ID.
            to_address_id: Optional destination address ID.
            amount: Optional amount for fee estimation.

        Returns:
            Withdrawal fee information.

        Raises:
            ValueError: If exchange_id is invalid.
            APIError: If API request fails.
        """
        self._validate_required(exchange_id, "exchange_id")

        try:
            resp = self._exchange_api.exchange_service_get_exchange_withdrawal_fee(
                id=exchange_id,
                to_address_id=to_address_id,
                amount=amount,
            )

            return getattr(resp, "result", None) or getattr(resp, "fee", None) or resp
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
