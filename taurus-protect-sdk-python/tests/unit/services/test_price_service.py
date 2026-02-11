"""Unit tests for PriceService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.price_service import PriceService


class TestGetCurrent:
    """Tests for PriceService.get_current()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        prices_api = MagicMock()
        service = PriceService(api_client=api_client, prices_api=prices_api)
        return service, prices_api

    def test_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.price_service_get_prices.return_value = resp

        result = service.get_current()

        assert result == []

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = []
        api.price_service_get_prices.return_value = resp

        service.get_current()

        api.price_service_get_prices.assert_called_once()

    def test_wraps_api_error(self) -> None:
        service, api = self._make_service()
        error = Exception("server error")
        error.status = 500
        error.body = None
        error.headers = {}
        api.price_service_get_prices.side_effect = error

        from taurus_protect.errors import APIError

        with pytest.raises(APIError):
            service.get_current()


class TestGetHistorical:
    """Tests for PriceService.get_historical()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        prices_api = MagicMock()
        service = PriceService(api_client=api_client, prices_api=prices_api)
        return service, prices_api

    def test_raises_on_empty_base_currency(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="base_currency"):
            service.get_historical(base_currency="", quote_currency="USD")

    def test_raises_on_empty_quote_currency(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="quote_currency"):
            service.get_historical(base_currency="BTC", quote_currency="")

    def test_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.price_service_get_prices_history.return_value = resp

        result = service.get_historical(base_currency="BTC", quote_currency="USD")

        assert result == []

    def test_passes_limit_as_string(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.price_service_get_prices_history.return_value = resp

        service.get_historical(base_currency="BTC", quote_currency="USD", limit=100)

        api.price_service_get_prices_history.assert_called_once_with(
            base="BTC",
            quote="USD",
            limit="100",
        )
