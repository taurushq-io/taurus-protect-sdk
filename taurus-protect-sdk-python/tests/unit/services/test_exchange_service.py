"""Unit tests for ExchangeService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import ExchangeApi
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.exchange_service import ExchangeService
from tests.unit.transport_stub import StubTransport, api_client


class TestExchangeServiceList:
    """ExchangeService.list over the real generated client (cursor)."""

    def _service(self) -> ExchangeService:
        ac = api_client()
        return ExchangeService(ac, ExchangeApi(ac))

    def test_rows_page_and_filters(self) -> None:
        reply = {
            "result": [{"id": "x1", "exchange": "kraken"}],
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            exchanges, page = self._service().list(
                page_size=10,
                currency_id="c",
                exchange_label="kraken",
                status="ok",
                only_positive_balance=False,
                sort_order="ASC",
                include_base_currency_valuation=True,
            )

        assert [e.id for e in exchanges] == ["x1"]
        assert page == CursorPage(page_size=10, next_cursor="n", has_more=True)
        assert dict(transport.last.query) == {
            "currencyID": "c",
            "exchangeLabel": "kraken",
            "status": "ok",
            "onlyPositiveBalance": "false",
            "sortOrder": "ASC",
            "includeBaseCurrencyValuation": "true",
            "cursor.pageSize": "10",
        }

    def test_page_two_is_reachable(self) -> None:
        """The list could only ever request the FIRST page."""
        with StubTransport() as transport:
            self._service().list(cursor="n")

        assert transport.last.param("cursor.currentPage") == "n"
        assert transport.last.param("cursor.pageRequest") == "NEXT"

    def test_invalid_page_size_sends_nothing(self) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match="page_size"):
                self._service().list(page_size=101)

        assert transport.requests == []


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
