"""Currency service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.errors import NotFoundError
from taurus_protect.mappers.currency import currencies_from_dto, currency_from_dto
from taurus_protect.models.currency import Currency
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.currencies_api import CurrenciesApi


class CurrencyService(BaseService):
    """
    Service for retrieving currency information.

    Provides methods to list all currencies, get specific currencies,
    and retrieve the base currency configured for the tenant.

    Example:
        >>> # List all enabled currencies
        >>> currencies = client.currencies.list()
        >>> for currency in currencies:
        ...     print(f"{currency.symbol}: {currency.name}")
        >>>
        >>> # Get currency by blockchain and network
        >>> eth = client.currencies.get_by_blockchain("ETH", "mainnet")
        >>> print(f"{eth.name} has {eth.decimals} decimals")
        >>>
        >>> # Get the base currency (e.g., USD, EUR, CHF)
        >>> base = client.currencies.get_base_currency()
        >>> print(f"Base currency: {base.symbol}")
    """

    def __init__(
        self,
        api_client: Any,
        currencies_api: "CurrenciesApi",
    ) -> None:
        """
        Initialize the currency service.

        Args:
            api_client: The OpenAPI client instance.
            currencies_api: The CurrenciesApi instance from OpenAPI client.
        """
        super().__init__(api_client)
        self._api = currencies_api

    def list(
        self,
        show_disabled: bool = False,
        include_logo: bool = False,
    ) -> List[Currency]:
        """
        Get all currencies.

        Args:
            show_disabled: If True, includes currencies disabled by business rules.
                Defaults to False (only enabled currencies).
            include_logo: If True, includes logo URLs in the response.
                Defaults to False (logos omitted for performance).

        Returns:
            List of currencies.

        Raises:
            APIError: If API request fails.
        """
        try:
            resp = self._api.wallet_service_get_currencies(
                show_disabled=show_disabled if show_disabled else None,
                include_logo=include_logo if include_logo else None,
            )

            result = getattr(resp, "currencies", None) or getattr(resp, "result", None)
            return currencies_from_dto(result) if result else []
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def get(self, currency_id: str) -> Currency:
        """
        Get a currency by ID.

        This method retrieves currency details using the currency ID.
        For looking up by blockchain and network, use `get_by_blockchain()`.

        Args:
            currency_id: The unique currency identifier.

        Returns:
            The currency.

        Raises:
            ValueError: If currency_id is empty.
            NotFoundError: If currency is not found.
            APIError: If API request fails.
        """
        self._validate_required(currency_id, "currency_id")

        # Get all currencies and filter by ID
        # The API doesn't have a direct get-by-id endpoint
        try:
            currencies = self.list(show_disabled=True)
            for currency in currencies:
                if currency.id == currency_id:
                    return currency

            raise NotFoundError(f"Currency with id '{currency_id}' not found")
        except NotFoundError:
            raise
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_by_blockchain(
        self,
        blockchain: str,
        network: str,
        contract_address: Optional[str] = None,
        token_id: Optional[str] = None,
    ) -> Currency:
        """
        Get a currency by blockchain and network.

        This retrieves the native currency for a blockchain/network pair,
        or a specific token if contract_address is provided.

        Args:
            blockchain: Blockchain type (e.g., "ETH", "BTC").
            network: Network name (e.g., "mainnet", "testnet").
            contract_address: Token contract address. If not provided,
                returns the native currency.
            token_id: For blockchains with multi-asset contracts (e.g., ALGO, XTZ),
                specifies which token within the contract.

        Returns:
            The currency.

        Raises:
            ValueError: If required arguments are missing.
            NotFoundError: If currency is not found.
            APIError: If API request fails.
        """
        self._validate_required(blockchain, "blockchain")
        self._validate_required(network, "network")

        try:
            resp = self._api.wallet_service_get_currency(
                unique_currency_filter_blockchain=blockchain,
                unique_currency_filter_network=network,
                unique_currency_filter_token_contract_address=contract_address,
                unique_currency_filter_token_id=token_id,
                show_disabled=True,
                currency_id=None,
                include_logo=None,
            )

            result = getattr(resp, "currency", None) or getattr(resp, "result", None)
            currency = currency_from_dto(result)

            if currency is None:
                raise NotFoundError(
                    f"Currency for blockchain '{blockchain}' network '{network}' not found"
                )

            return currency
        except NotFoundError:
            raise
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_base_currency(self) -> Currency:
        """
        Get the base currency configured for the tenant.

        The base currency is used for fiat valuations and is typically
        configured as CHF, EUR, or USD.

        Returns:
            The base currency.

        Raises:
            NotFoundError: If no base currency is configured.
            APIError: If API request fails.
        """
        try:
            resp = self._api.wallet_service_get_base_currency()

            result = getattr(resp, "currency", None) or getattr(resp, "result", None)
            currency = currency_from_dto(result)

            if currency is None:
                raise NotFoundError("Base currency not configured")

            return currency
        except NotFoundError:
            raise
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e
