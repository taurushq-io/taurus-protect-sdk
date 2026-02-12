"""Price service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers.statistics import (
    price_history_from_dto,
    prices_from_dto,
)
from taurus_protect.models.statistics import Price, PriceHistoryPoint
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class PriceService(BaseService):
    """
    Service for cryptocurrency price operations.

    Provides methods to retrieve current and historical price data
    for various currency pairs.

    Example:
        >>> # Get all current prices
        >>> prices = client.prices.get_current()
        >>> for price in prices:
        ...     print(f"{price.currency_from}/{price.currency_to}: {price.rate}")
        >>>
        >>> # Get current price for a specific currency
        >>> btc_prices = client.prices.get_current(currency="BTC")
        >>>
        >>> # Get historical prices
        >>> history = client.prices.get_historical(
        ...     base_currency="BTC",
        ...     quote_currency="USD",
        ...     limit=100,
        ... )
    """

    def __init__(self, api_client: Any, prices_api: Any) -> None:
        """
        Initialize price service.

        Args:
            api_client: The OpenAPI client instance.
            prices_api: The PricesApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._prices_api = prices_api

    def get_current(self, currency: Optional[str] = None) -> List[Price]:
        """
        Get current prices for all currencies or a specific currency.

        Args:
            currency: Optional currency symbol to filter prices (e.g., "BTC", "ETH").
                     If not provided, returns all available prices.

        Returns:
            List of current prices.

        Raises:
            APIError: If API request fails.

        Example:
            >>> # Get all prices
            >>> all_prices = client.prices.get_current()
            >>>
            >>> # Get prices for Bitcoin
            >>> btc_prices = client.prices.get_current(currency="BTC")
        """
        try:
            resp = self._prices_api.price_service_get_prices()

            result = getattr(resp, "result", None)
            prices = prices_from_dto(result) if result else []

            # Filter by currency if specified
            if currency:
                prices = [
                    p for p in prices if p.currency_from == currency or p.currency_to == currency
                ]

            return prices
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def get_historical(
        self,
        base_currency: str,
        quote_currency: str,
        limit: Optional[int] = None,
    ) -> List[PriceHistoryPoint]:
        """
        Get historical prices for a currency pair.

        Retrieves the price history for a specific currency pair, such as
        BTC/USD or ETH/EUR.

        Args:
            base_currency: Base currency symbol (e.g., "BTC").
            quote_currency: Quote currency symbol (e.g., "USD").
            limit: Maximum number of history points to return.
                  If not provided, returns all available history.

        Returns:
            List of price history points, ordered from oldest to newest.

        Raises:
            ValueError: If base_currency or quote_currency is empty.
            APIError: If API request fails.

        Example:
            >>> # Get BTC/USD price history
            >>> history = client.prices.get_historical(
            ...     base_currency="BTC",
            ...     quote_currency="USD",
            ...     limit=100,
            ... )
            >>> for point in history:
            ...     print(f"{point.timestamp}: {point.rate}")
        """
        self._validate_required(base_currency, "base_currency")
        self._validate_required(quote_currency, "quote_currency")

        try:
            limit_str = str(limit) if limit is not None else None

            resp = self._prices_api.price_service_get_prices_history(
                base=base_currency,
                quote=quote_currency,
                limit=limit_str,
            )

            result = getattr(resp, "result", None)
            return price_history_from_dto(result) if result else []
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
