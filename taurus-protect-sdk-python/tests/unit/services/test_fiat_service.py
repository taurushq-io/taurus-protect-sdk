"""Unit tests for FiatService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import CurrenciesApi, FiatApi
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.fiat_service import FiatService
from tests.unit.transport_stub import StubTransport, api_client


class TestFiatServiceListAccounts:
    """list_fiat_provider_accounts: the accounts list could not run (invalid arguments)."""

    def _service(self) -> FiatService:
        ac = api_client()
        return FiatService(ac, FiatApi(ac), CurrenciesApi(ac))

    def test_provider_and_label_are_required(self) -> None:
        with pytest.raises(ValueError, match="provider"):
            self._service().list_fiat_provider_accounts("", "label")
        with pytest.raises(ValueError, match="label"):
            self._service().list_fiat_provider_accounts("provider", "")

    def test_rows_page_and_filters(self) -> None:
        reply = {
            "result": [{"id": "acc-1", "accountName": "Main", "totalBalance": "10"}],
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            accounts, page = self._service().list_fiat_provider_accounts(
                "bank", "main", page_size=5, account_type="CURRENT", sort_order="ASC"
            )

        assert [(a.id, a.name, a.balance) for a in accounts] == [("acc-1", "Main", "10")]
        assert page == CursorPage(page_size=5, next_cursor="n", has_more=True)
        assert transport.last.path == "/api/rest/v1/fiat_providers/accounts"
        assert transport.last.query == sorted(
            [
                ("provider", "bank"),
                ("label", "main"),
                ("accountType", "CURRENT"),
                ("sortOrder", "ASC"),
                ("cursor.pageSize", "5"),
            ]
        )

    def test_continuation(self) -> None:
        with StubTransport() as transport:
            self._service().list_fiat_provider_accounts("bank", "main", cursor="n")

        assert transport.last.param("cursor.currentPage") == "n"
        assert transport.last.param("cursor.pageRequest") == "NEXT"


class TestFiatServiceListEntities:
    """list_fiat_provider_entities: new cursor list."""

    def _service(self) -> FiatService:
        ac = api_client()
        return FiatService(ac, FiatApi(ac), CurrenciesApi(ac))

    def test_rows_page_and_filters(self) -> None:
        reply = {"result": [{"id": "e1", "name": "Entity"}], "cursor": {"currentPage": "n"}}
        with StubTransport(reply) as transport:
            entities, page = self._service().list_fiat_provider_entities(
                provider="bank", label="main", sort_order="DESC"
            )

        assert [(e.id, e.name) for e in entities] == [("e1", "Entity")]
        assert page == CursorPage(page_size=20)
        assert transport.last.path == "/api/rest/v1/fiat_providers/entities"
        assert transport.last.query == sorted(
            [
                ("provider", "bank"),
                ("label", "main"),
                ("sortOrder", "DESC"),
                ("cursor.pageSize", "20"),
            ]
        )


class TestFiatServiceGetAccount:
    """Tests for FiatService.get_account()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fiat_api = MagicMock()
        currencies_api = MagicMock()
        service = FiatService(
            api_client=api_client, fiat_api=fiat_api, currencies_api=currencies_api
        )
        return service, fiat_api, currencies_api

    def test_raises_on_empty_id(self) -> None:
        service, _, _ = self._make_service()
        with pytest.raises(ValueError, match="account_id"):
            service.get_account(account_id="")

    def test_raises_not_found_when_none(self) -> None:
        service, fiat_api, _ = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.account = None
        fiat_api.fiat_provider_service_get_fiat_provider_account.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get_account(account_id="acc-1")


class TestFiatServiceGetBaseCurrency:
    """Tests for FiatService.get_base_currency()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fiat_api = MagicMock()
        currencies_api = MagicMock()
        service = FiatService(
            api_client=api_client, fiat_api=fiat_api, currencies_api=currencies_api
        )
        return service, fiat_api, currencies_api

    def test_returns_default_usd_when_none(self) -> None:
        service, _, currencies_api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.currency = None
        currencies_api.wallet_service_get_base_currency.return_value = resp

        currency = service.get_base_currency()

        assert currency.code == "USD"


class TestFiatServiceGetRate:
    """Tests for FiatService.get_rate()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fiat_api = MagicMock()
        currencies_api = MagicMock()
        service = FiatService(
            api_client=api_client, fiat_api=fiat_api, currencies_api=currencies_api
        )
        return service, fiat_api, currencies_api

    def test_raises_on_empty_from(self) -> None:
        service, _, _ = self._make_service()
        with pytest.raises(ValueError, match="from_currency"):
            service.get_rate(from_currency="", to_currency="EUR")

    def test_raises_on_empty_to(self) -> None:
        service, _, _ = self._make_service()
        with pytest.raises(ValueError, match="to_currency"):
            service.get_rate(from_currency="USD", to_currency="")

    def test_raises_not_found_when_rate_missing(self) -> None:
        service, _, currencies_api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.currencies = None
        currencies_api.wallet_service_get_currencies.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get_rate(from_currency="USD", to_currency="EUR")


class TestFiatServiceListProviders:
    """Tests for FiatService.list_providers()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fiat_api = MagicMock()
        currencies_api = MagicMock()
        service = FiatService(
            api_client=api_client, fiat_api=fiat_api, currencies_api=currencies_api
        )
        return service, fiat_api, currencies_api

    def test_returns_empty_when_no_data(self) -> None:
        service, fiat_api, _ = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.providers = None
        fiat_api.fiat_provider_service_get_fiat_providers.return_value = resp

        result = service.list_providers()

        assert result == []
