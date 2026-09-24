"""Price service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, Dict, List, Optional, Tuple

from taurus_protect.helpers.price_verifier import verify_prices
from taurus_protect.mappers.statistics import (
    price_history_from_dto,
    prices_from_dto,
)
from taurus_protect.models.pagination import (
    PRICE_HISTORY_MAX_LIMIT,
    CursorPage,
    cursor_page,
    cursor_request,
    resolve_page_size,
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
        >>> # Walk every current price
        >>> cursor = None
        >>> while True:
        ...     prices, page = client.prices.get_current(page_size=100, cursor=cursor)
        ...     for price in prices:
        ...         print(f"{price.currency_from}/{price.currency_to}: {price.rate}")
        ...     if not page.has_more:
        ...         break
        ...     cursor = page.next_cursor
        >>>
        >>> # Get historical prices
        >>> history = client.prices.get_historical(
        ...     base_currency="BTC",
        ...     quote_currency="USD",
        ...     limit=100,
        ... )
    """

    def __init__(self, api_client: Any, prices_api: Any, rules_cache: Any) -> None:
        """
        Initialize price service.

        Rate and decimals feed amount conversion, so an unverified price is a wrong
        number a caller acts on. Whether prices must be signed is decided by the
        SuperAdmin-verified rules container: see helpers.price_verifier.verify_price.

        Args:
            api_client: The OpenAPI client instance.
            prices_api: The PricesApi service from OpenAPI client.
            rules_cache: Rules container cache supplying the PRICEUPDATER keys. Required.

        Raises:
            ValueError: If rules_cache is None.
        """
        super().__init__(api_client)
        if rules_cache is None:
            raise ValueError(
                "rules_cache cannot be None - price signature verification is mandatory"
            )
        self._prices_api = prices_api
        self._rules_cache = rules_cache

    def get_current(
        self,
        *,
        from_currency_id: Optional[str] = None,
        to_currency_ids: Optional[List[str]] = None,
        only_primary: Optional[bool] = None,
        sort_order: Optional[str] = None,
        page_size: Optional[int] = None,
        cursor: Optional[str] = None,
        current_page: Optional[str] = None,
        page_request: Optional[str] = None,
    ) -> Tuple[List[Price], CursorPage]:
        """
        List current prices, one page at a time, each signature verified.

        ``from_currency_id`` alone selects the prices quoted from that currency,
        ``to_currency_ids`` alone the prices quoted into those currencies, and both
        together the pairs between them.

        Args:
            from_currency_id: The quoted-from currency ID.
            to_currency_ids: The quoted-into currency IDs.
            only_primary: Only primary prices.
            sort_order: ASC or DESC.
            page_size: Page size (default 20, max 100).
            cursor: ``page.next_cursor`` from the previous page, to continue.
            current_page: Low-level page token; not with ``cursor``.
            page_request: Low-level page direction (FIRST, PREVIOUS, NEXT, LAST).

        Returns:
            Tuple of (verified prices, page).

        Raises:
            ValueError: If the page size is invalid or cursor options conflict.
            IntegrityError: If a price signature does not verify.
            APIError: If API request fails.

        Example:
            >>> prices, page = client.prices.get_current(from_currency_id="<currency-id>")
        """
        from taurus_protect._internal.openapi.models.tgvalidatord_currency_from_filter import (
            TgvalidatordCurrencyFromFilter,
        )
        from taurus_protect._internal.openapi.models.tgvalidatord_currency_from_to_filter import (
            TgvalidatordCurrencyFromToFilter,
        )
        from taurus_protect._internal.openapi.models.tgvalidatord_currency_to_filter import (
            TgvalidatordCurrencyToFilter,
        )
        from taurus_protect._internal.openapi.models.tgvalidatord_query_prices_v2_request import (
            TgvalidatordQueryPricesV2Request,
        )

        req = cursor_request(
            page_size, cursor, current_page=current_page, page_request=page_request
        )

        # The request carries one of from / fromTo / to.
        filters: Dict[str, Any] = {}
        if from_currency_id and to_currency_ids:
            filters["from_to"] = TgvalidatordCurrencyFromToFilter(
                currency_from_id=from_currency_id, currency_to_ids=list(to_currency_ids)
            )
        elif from_currency_id:
            filters["var_from"] = TgvalidatordCurrencyFromFilter(currency_from_id=from_currency_id)
        elif to_currency_ids:
            filters["to"] = TgvalidatordCurrencyToFilter(currency_to_ids=list(to_currency_ids))

        try:
            body = TgvalidatordQueryPricesV2Request(
                only_primary=only_primary,
                sort_order=sort_order,
                cursor=self._request_cursor_body(req),
                **filters,
            )
            resp = self._prices_api.price_service_query_prices_v2(body=body)

            prices = prices_from_dto(resp.result or [])
            if prices:
                verify_prices(prices, self._rules_cache.get_decoded_rules_container())

            return prices, cursor_page(req.page_size, resp.cursor)
        except Exception as e:
            from taurus_protect.errors import APIError, IntegrityError

            # IntegrityError is NOT an APIError, so without naming it a failed price
            # signature check would be remapped to a retryable ServerError.
            if isinstance(e, (APIError, IntegrityError, ValueError)):
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
            limit: Number of daily points, newest first (default 20, max 365). The
                series cannot page, so this is the whole reply.

        Returns:
            List of price history points, ordered from oldest to newest.

        Raises:
            ValueError: If base_currency or quote_currency is empty, or limit is invalid.
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
        size = resolve_page_size(limit, "limit", maximum=PRICE_HISTORY_MAX_LIMIT)

        try:
            resp = self._prices_api.price_service_get_prices_history(
                base=base_currency,
                quote=quote_currency,
                limit=str(size),
            )

            return price_history_from_dto(resp.result or [])
        except Exception as e:
            from taurus_protect.errors import APIError, IntegrityError

            # Consistent with the siblings in this file: IntegrityError is not an APIError, so omitting it here turns a
            # security failure into a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, ValueError)):
                raise
            raise self._handle_error(e) from e
