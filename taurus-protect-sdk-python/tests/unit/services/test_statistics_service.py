"""Unit tests for StatisticsService."""

from __future__ import annotations

from datetime import datetime
from unittest.mock import MagicMock

import pytest

from taurus_protect.services.statistics_service import StatisticsService


class TestGetSummary:
    """Tests for StatisticsService.get_summary()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        statistics_api = MagicMock()
        service = StatisticsService(
            api_client=api_client, statistics_api=statistics_api
        )
        return service, statistics_api

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.statistics_service_get_portfolio_statistics.return_value = resp

        service.get_summary()

        api.statistics_service_get_portfolio_statistics.assert_called_once()

    def test_wraps_api_error(self) -> None:
        service, api = self._make_service()
        error = Exception("server error")
        error.status = 500
        error.body = None
        error.headers = {}
        api.statistics_service_get_portfolio_statistics.side_effect = error

        from taurus_protect.errors import APIError

        with pytest.raises(APIError):
            service.get_summary()


class TestGetTransactionStats:
    """Tests for StatisticsService.get_transaction_stats()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        statistics_api = MagicMock()
        service = StatisticsService(
            api_client=api_client, statistics_api=statistics_api
        )
        return service, statistics_api

    def test_raises_when_from_after_to(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="from_date cannot be after to_date"):
            service.get_transaction_stats(
                from_date=datetime(2025, 12, 31),
                to_date=datetime(2025, 1, 1),
            )

    def test_returns_default_stats(self) -> None:
        service, _ = self._make_service()

        stats = service.get_transaction_stats()

        assert stats.total_count == "0"
        assert stats.incoming_count == "0"
        assert stats.outgoing_count == "0"

    def test_accepts_valid_date_range(self) -> None:
        service, _ = self._make_service()

        stats = service.get_transaction_stats(
            from_date=datetime(2025, 1, 1),
            to_date=datetime(2025, 12, 31),
        )

        assert stats is not None
