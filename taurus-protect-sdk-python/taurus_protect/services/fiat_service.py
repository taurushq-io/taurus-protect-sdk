"""Fiat service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from taurus_protect.mappers._base import safe_bool, safe_int, safe_string
from taurus_protect.models.blockchain import ExchangeRate, FiatCurrency, FiatProviderAccount
from taurus_protect.models.pagination import Pagination
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def fiat_currency_from_dto(dto: Any) -> Optional[FiatCurrency]:
    """
    Convert an OpenAPI fiat currency DTO to domain model.

    Args:
        dto: The OpenAPI DTO object.

    Returns:
        FiatCurrency model or None if dto is None.
    """
    if dto is None:
        return None

    return FiatCurrency(
        id=safe_string(getattr(dto, "id", None)),
        code=safe_string(
            getattr(dto, "code", None) or getattr(dto, "currency", None) or getattr(dto, "id", None)
        ),
        name=getattr(dto, "name", None),
        symbol=getattr(dto, "symbol", None),
        decimals=safe_int(getattr(dto, "decimals", 2)),
        enabled=safe_bool(getattr(dto, "enabled", True)),
    )


def fiat_currencies_from_dto(dtos: Any) -> List[FiatCurrency]:
    """
    Convert a list of OpenAPI fiat currency DTOs to domain models.

    Args:
        dtos: List of OpenAPI DTO objects.

    Returns:
        List of FiatCurrency models.
    """
    if dtos is None:
        return []
    return [f for dto in dtos if (f := fiat_currency_from_dto(dto)) is not None]


def fiat_provider_account_from_dto(dto: Any) -> Optional[FiatProviderAccount]:
    """
    Convert an OpenAPI fiat provider account DTO to domain model.

    Args:
        dto: The OpenAPI DTO object.

    Returns:
        FiatProviderAccount model or None if dto is None.
    """
    if dto is None:
        return None

    return FiatProviderAccount(
        id=safe_string(getattr(dto, "id", None)),
        name=getattr(dto, "name", None),
        provider=getattr(dto, "provider", None),
        currency_code=getattr(dto, "currency_code", None) or getattr(dto, "currencyCode", None),
        balance=getattr(dto, "balance", None),
        enabled=safe_bool(getattr(dto, "enabled", True)),
    )


def fiat_provider_accounts_from_dto(dtos: Any) -> List[FiatProviderAccount]:
    """
    Convert a list of OpenAPI fiat provider account DTOs to domain models.

    Args:
        dtos: List of OpenAPI DTO objects.

    Returns:
        List of FiatProviderAccount models.
    """
    if dtos is None:
        return []
    return [f for dto in dtos if (f := fiat_provider_account_from_dto(dto)) is not None]


class FiatService(BaseService):
    """
    Service for fiat currency and fiat provider operations.

    Provides methods to list fiat currencies, retrieve exchange rates,
    and manage fiat provider accounts.

    Example:
        >>> # List fiat provider accounts
        >>> accounts = client.fiat.list()
        >>> for account in accounts:
        ...     print(f"{account.name}: {account.balance} {account.currency_code}")
        >>>
        >>> # Get exchange rate
        >>> rate = client.fiat.get_rate("USD", "EUR")
        >>> print(f"1 USD = {rate.rate} EUR")
    """

    def __init__(self, api_client: Any, fiat_api: Any, currencies_api: Any) -> None:
        """
        Initialize fiat service.

        Args:
            api_client: The OpenAPI client instance.
            fiat_api: The FiatAPI service from OpenAPI client.
            currencies_api: The CurrenciesAPI service from OpenAPI client.
        """
        super().__init__(api_client)
        self._fiat_api = fiat_api
        self._currencies_api = currencies_api

    def list(
        self,
        limit: int = 50,
        offset: int = 0,
    ) -> Tuple[List[FiatProviderAccount], Optional[Pagination]]:
        """
        List fiat provider accounts with pagination.

        Args:
            limit: Maximum number of accounts to return (must be positive).
            offset: Number of accounts to skip (must be non-negative).

        Returns:
            Tuple of (accounts list, pagination info).

        Raises:
            ValueError: If limit or offset are invalid.
            APIError: If API request fails.
        """
        if limit <= 0:
            raise ValueError("limit must be positive")
        if offset < 0:
            raise ValueError("offset cannot be negative")

        try:
            resp = self._fiat_api.fiat_provider_service_get_fiat_provider_accounts(
                limit=str(limit),
                offset=str(offset),
            )

            result = getattr(resp, "result", None) or getattr(resp, "accounts", None)
            accounts = fiat_provider_accounts_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None) or getattr(resp, "totalItems", None),
                offset=getattr(resp, "offset", None),
                limit=limit,
            )

            return accounts, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_account(self, account_id: str) -> FiatProviderAccount:
        """
        Get a fiat provider account by ID.

        Args:
            account_id: The fiat provider account ID to retrieve.

        Returns:
            The fiat provider account.

        Raises:
            ValueError: If account_id is invalid.
            NotFoundError: If account not found.
            APIError: If API request fails.
        """
        self._validate_required(account_id, "account_id")

        try:
            resp = self._fiat_api.fiat_provider_service_get_fiat_provider_account(id=account_id)

            result = getattr(resp, "result", None) or getattr(resp, "account", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Fiat provider account {account_id} not found")

            account = fiat_provider_account_from_dto(result)
            if account is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Fiat provider account {account_id} not found")

            return account
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_base_currency(self) -> FiatCurrency:
        """
        Get the configured base currency.

        Returns:
            The base fiat currency (e.g., USD, EUR, CHF).

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._currencies_api.wallet_service_get_base_currency()

            result = getattr(resp, "result", None) or getattr(resp, "currency", None)
            if result is None:
                # Default to USD if not configured
                return FiatCurrency(id="USD", code="USD", name="US Dollar", symbol="$")

            currency = fiat_currency_from_dto(result)
            if currency is None:
                return FiatCurrency(id="USD", code="USD", name="US Dollar", symbol="$")

            return currency
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def get_rate(self, from_currency: str, to_currency: str) -> ExchangeRate:
        """
        Get the exchange rate between two currencies.

        Note: This method retrieves the rate by fetching currency information.
        The actual rate calculation may depend on the API implementation.

        Args:
            from_currency: Source currency code (e.g., "USD").
            to_currency: Target currency code (e.g., "EUR").

        Returns:
            The exchange rate.

        Raises:
            ValueError: If currency codes are invalid.
            NotFoundError: If rate not available.
            APIError: If API request fails.
        """
        self._validate_required(from_currency, "from_currency")
        self._validate_required(to_currency, "to_currency")

        try:
            # Try to get the rate from the currencies API
            resp = self._currencies_api.wallet_service_get_currencies(
                blockchain=None,
                network=None,
                include_base_currency_valuation=True,
            )

            result = getattr(resp, "result", None) or getattr(resp, "currencies", None)

            # Look for matching currencies to calculate rate
            from_rate = None
            to_rate = None

            if result:
                for currency in result:
                    code = getattr(currency, "currency", None) or getattr(currency, "symbol", None)
                    valuation = getattr(currency, "base_currency_valuation", None) or getattr(
                        currency, "baseCurrencyValuation", None
                    )

                    if code == from_currency and valuation:
                        from_rate = valuation
                    if code == to_currency and valuation:
                        to_rate = valuation

            # If we found both rates, calculate the exchange rate
            if from_rate is not None and to_rate is not None:
                try:
                    rate = float(to_rate) / float(from_rate)
                    return ExchangeRate(
                        from_currency=from_currency,
                        to_currency=to_currency,
                        rate=str(rate),
                    )
                except (ValueError, ZeroDivisionError):
                    pass

            # If rate not found, return a placeholder
            from taurus_protect.errors import NotFoundError

            raise NotFoundError(
                f"Exchange rate from {from_currency} to {to_currency} not available"
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_providers(self) -> List[Any]:
        """
        List available fiat providers.

        Returns:
            List of fiat providers.

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._fiat_api.fiat_provider_service_get_fiat_providers()

            result = getattr(resp, "result", None) or getattr(resp, "providers", None)
            return result or []
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e
