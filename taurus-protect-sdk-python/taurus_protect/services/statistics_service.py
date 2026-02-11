"""Statistics service for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING, Any, Optional

from taurus_protect.mappers.statistics import portfolio_statistics_from_dto
from taurus_protect.models.statistics import PortfolioStatistics, TransactionStatistics
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class StatisticsService(BaseService):
    """
    Service for portfolio and transaction statistics.

    Provides methods to retrieve summary statistics about the portfolio
    and transaction activity over time.

    Example:
        >>> # Get portfolio summary
        >>> summary = client.statistics.get_summary()
        >>> print(f"Total wallets: {summary.wallets_count}")
        >>> print(f"Total addresses: {summary.addresses_count}")
        >>> print(f"Total balance: {summary.total_balance_base_currency}")
        >>>
        >>> # Get transaction statistics for a date range
        >>> from datetime import datetime
        >>> stats = client.statistics.get_transaction_stats(
        ...     from_date=datetime(2024, 1, 1),
        ...     to_date=datetime(2024, 12, 31),
        ... )
    """

    def __init__(self, api_client: Any, statistics_api: Any) -> None:
        """
        Initialize statistics service.

        Args:
            api_client: The OpenAPI client instance.
            statistics_api: The StatisticsApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._statistics_api = statistics_api

    def get_summary(self) -> Optional[PortfolioStatistics]:
        """
        Get summary statistics for the portfolio.

        Returns aggregated statistics about the entire portfolio including
        wallet counts, address counts, and total balance information.

        Returns:
            Portfolio statistics or None if not available.

        Raises:
            APIError: If API request fails.

        Example:
            >>> summary = client.statistics.get_summary()
            >>> print(f"Total wallets: {summary.wallets_count}")
            >>> print(f"Total addresses: {summary.addresses_count}")
            >>> print(f"Total balance: {summary.total_balance_base_currency}")
        """
        try:
            resp = self._statistics_api.statistics_service_get_portfolio_statistics()

            result = getattr(resp, "result", None)
            return portfolio_statistics_from_dto(result)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def get_transaction_stats(
        self,
        from_date: Optional[datetime] = None,
        to_date: Optional[datetime] = None,
    ) -> TransactionStatistics:
        """
        Get transaction statistics for a date range.

        Returns aggregated transaction statistics for the specified time period.
        If no dates are provided, returns statistics for all time.

        Note: This is a placeholder implementation. The actual API endpoint
        for transaction-specific statistics may vary. This method demonstrates
        the intended interface pattern.

        Args:
            from_date: Start date for the statistics period.
            to_date: End date for the statistics period.

        Returns:
            Transaction statistics for the specified period.

        Raises:
            ValueError: If from_date is after to_date.
            APIError: If API request fails.

        Example:
            >>> from datetime import datetime
            >>> stats = client.statistics.get_transaction_stats(
            ...     from_date=datetime(2024, 1, 1),
            ...     to_date=datetime(2024, 12, 31),
            ... )
            >>> print(f"Total transactions: {stats.total_count}")
            >>> print(f"Total volume: {stats.total_volume}")
        """
        if from_date and to_date and from_date > to_date:
            raise ValueError("from_date cannot be after to_date")

        try:
            # The API may have different endpoints for transaction statistics.
            # For now, we return a default TransactionStatistics.
            # This can be expanded when the specific API endpoint is identified.
            #
            # Potential API endpoints that could be used:
            # - statistics_service_get_currency_statistics
            # - statistics_service_get_currency_statistics_history
            #
            # For a full implementation, these would need to be mapped to
            # transaction-specific statistics.

            return TransactionStatistics(
                total_count="0",
                incoming_count="0",
                outgoing_count="0",
                total_volume="0",
                incoming_volume="0",
                outgoing_volume="0",
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e
