"""Unit tests for ExchangeService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.exchange_service import ExchangeService


class TestExchangeServiceList:
    """Tests for ExchangeService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        exchange_api = MagicMock()
        service = ExchangeService(api_client=api_client, exchange_api=exchange_api)
        return service, exchange_api

    def test_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.exchange_accounts = None
        resp.exchangeAccounts = None
        resp.cursor = None
        api.exchange_service_get_exchanges.return_value = resp

        exchanges, pagination = service.list()

        assert exchanges == []


class TestExchangeServiceGet:
    """Tests for ExchangeService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        exchange_api = MagicMock()
        service = ExchangeService(api_client=api_client, exchange_api=exchange_api)
        return service, exchange_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="exchange_id"):
            service.get("")

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.exchange_account = None
        resp.exchangeAccount = None
        api.exchange_service_get_exchange.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get("123")


class TestExchangeServiceListCounterparties:
    """Tests for ExchangeService.list_counterparties()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        exchange_api = MagicMock()
        service = ExchangeService(api_client=api_client, exchange_api=exchange_api)
        return service, exchange_api

    def test_returns_empty_when_no_data(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.exchanges = None
        api.exchange_service_get_exchange_counterparties.return_value = resp

        result = service.list_counterparties()

        assert result == []


class TestExchangeServiceGetWithdrawalFee:
    """Tests for ExchangeService.get_withdrawal_fee()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        exchange_api = MagicMock()
        service = ExchangeService(api_client=api_client, exchange_api=exchange_api)
        return service, exchange_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="exchange_id"):
            service.get_withdrawal_fee(exchange_id="")

    def test_calls_api_with_params(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = {"fee": "0.001"}
        api.exchange_service_get_exchange_withdrawal_fee.return_value = resp

        result = service.get_withdrawal_fee(
            exchange_id="ex-1",
            to_address_id="addr-1",
            amount="100",
        )

        api.exchange_service_get_exchange_withdrawal_fee.assert_called_once_with(
            id="ex-1",
            to_address_id="addr-1",
            amount="100",
        )
