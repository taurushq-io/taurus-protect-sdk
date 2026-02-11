"""Unit tests for FiatService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.fiat_service import FiatService


class TestFiatServiceList:
    """Tests for FiatService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fiat_api = MagicMock()
        currencies_api = MagicMock()
        service = FiatService(
            api_client=api_client, fiat_api=fiat_api, currencies_api=currencies_api
        )
        return service, fiat_api, currencies_api

    def test_raises_on_invalid_limit(self) -> None:
        service, _, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_raises_on_negative_offset(self) -> None:
        service, _, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_returns_empty_when_no_results(self) -> None:
        service, fiat_api, _ = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.accounts = None
        resp.total_items = None
        resp.totalItems = None
        resp.offset = None
        fiat_api.fiat_provider_service_get_fiat_provider_accounts.return_value = resp

        accounts, pagination = service.list()

        assert accounts == []


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
