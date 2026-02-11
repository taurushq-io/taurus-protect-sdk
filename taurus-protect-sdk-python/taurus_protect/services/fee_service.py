"""Fee service for Taurus-PROTECT SDK."""

from __future__ import annotations

from decimal import Decimal
from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.models.staking import FeeEstimate
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.fee_api import FeeApi


class FeeService(BaseService):
    """
    Service for fee estimation operations.

    Provides methods to estimate transaction fees across different
    blockchains and currencies.

    Example:
        >>> # Estimate fee for an ETH transfer
        >>> estimate = client.fees.estimate(
        ...     currency="ETH",
        ...     amount="1.5",
        ...     destination="0x1234..."
        ... )
        >>> print(f"Low: {estimate.fee_low} ETH")
        >>> print(f"Medium: {estimate.fee_medium} ETH")
        >>> print(f"High: {estimate.fee_high} ETH")
    """

    def __init__(self, api_client: Any, fee_api: "FeeApi") -> None:
        """
        Initialize fee service.

        Args:
            api_client: The OpenAPI client instance.
            fee_api: The FeeApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._api = fee_api

    def estimate(
        self,
        currency: str,
        amount: Optional[str] = None,
        destination: Optional[str] = None,
    ) -> FeeEstimate:
        """
        Estimate transaction fee.

        Provides fee estimates for a potential transaction. The estimates
        are based on current network conditions.

        Args:
            currency: Currency symbol (e.g., "ETH", "BTC", "SOL").
            amount: Optional transaction amount for more accurate estimates.
            destination: Optional destination address for more accurate estimates.

        Returns:
            Fee estimate with low, medium, and high priority options.

        Raises:
            ValueError: If currency is empty.
            APIError: If the API request fails.

        Example:
            >>> estimate = client.fees.estimate(
            ...     currency="ETH",
            ...     amount="1.0",
            ...     destination="0xabc..."
            ... )
            >>> if estimate.gas_price:
            ...     print(f"Gas price: {estimate.gas_price} gwei")
        """
        self._validate_required(currency, "currency")

        try:
            # Use the v2 endpoint for better fee data
            resp = self._api.fee_service_get_fees_v2()

            # Find the fee info for the requested currency
            result = getattr(resp, "result", None) or getattr(resp, "fees", None)

            if not result:
                # Return empty estimate if no data
                return FeeEstimate(currency=currency)

            # Search for the currency in the results
            if isinstance(result, list):
                for fee_data in result:
                    fee_currency = (
                        getattr(fee_data, "currency", None)
                        or getattr(fee_data, "symbol", None)
                        or ""
                    )
                    if fee_currency.upper() == currency.upper():
                        return self._map_fee_estimate(fee_data, currency, amount)

            # If not found in list or single result, try to use it directly
            elif result:
                return self._map_fee_estimate(result, currency, amount)

            # Return empty estimate if currency not found
            return FeeEstimate(currency=currency)

        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list(self) -> List[FeeEstimate]:
        """
        List fee estimates for all supported currencies.

        Returns:
            List of fee estimates for all currencies.

        Raises:
            APIError: If the API request fails.

        Example:
            >>> fees = client.fees.list()
            >>> for fee in fees:
            ...     print(f"{fee.currency}: {fee.fee_medium}")
        """
        try:
            resp = self._api.fee_service_get_fees_v2()

            result = getattr(resp, "result", None) or getattr(resp, "fees", None)
            if not result:
                return []

            estimates: List[FeeEstimate] = []

            if isinstance(result, list):
                for fee_data in result:
                    currency = (
                        getattr(fee_data, "currency", None)
                        or getattr(fee_data, "symbol", None)
                        or ""
                    )
                    if currency:
                        estimate = self._map_fee_estimate(fee_data, currency, None)
                        estimates.append(estimate)
            else:
                # Single result
                currency = (
                    getattr(result, "currency", None) or getattr(result, "symbol", None) or ""
                )
                if currency:
                    estimates.append(self._map_fee_estimate(result, currency, None))

            return estimates

        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def _map_fee_estimate(
        self,
        dto: Any,
        currency: str,
        amount: Optional[str],
    ) -> FeeEstimate:
        """Map fee DTO to FeeEstimate model."""
        if dto is None:
            return FeeEstimate(currency=currency)

        blockchain = getattr(dto, "blockchain", None) or ""
        network = getattr(dto, "network", None) or ""

        # Extract fee levels
        fee_low_raw = (
            getattr(dto, "fee_low", None) or getattr(dto, "low", None) or getattr(dto, "slow", None)
        )
        fee_medium_raw = (
            getattr(dto, "fee_medium", None)
            or getattr(dto, "medium", None)
            or getattr(dto, "standard", None)
            or getattr(dto, "average", None)
        )
        fee_high_raw = (
            getattr(dto, "fee_high", None)
            or getattr(dto, "high", None)
            or getattr(dto, "fast", None)
        )

        fee_low = Decimal(str(fee_low_raw)) if fee_low_raw is not None else None
        fee_medium = Decimal(str(fee_medium_raw)) if fee_medium_raw is not None else None
        fee_high = Decimal(str(fee_high_raw)) if fee_high_raw is not None else None

        # EVM-specific fields
        gas_limit_raw = getattr(dto, "gas_limit", None)
        gas_limit = int(gas_limit_raw) if gas_limit_raw is not None else None

        gas_price_raw = getattr(dto, "gas_price", None) or getattr(dto, "gasPrice", None)
        gas_price = Decimal(str(gas_price_raw)) if gas_price_raw is not None else None

        # Parse amount if provided
        amount_decimal = Decimal(amount) if amount else None

        return FeeEstimate(
            currency=currency,
            blockchain=str(blockchain),
            network=str(network),
            amount=amount_decimal,
            fee_low=fee_low,
            fee_medium=fee_medium,
            fee_high=fee_high,
            gas_limit=gas_limit,
            gas_price=gas_price,
        )
