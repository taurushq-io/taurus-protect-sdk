"""Unit tests for CurrencyService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import NotFoundError
from taurus_protect.services.currency_service import CurrencyService


class TestList:
    """Tests for CurrencyService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        currencies_api = MagicMock()
        service = CurrencyService(api_client=api_client, currencies_api=currencies_api)
        return service, currencies_api

    def test_list_returns_currencies(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.currencies = [MagicMock(), MagicMock()]
        reply.result = None
        api.wallet_service_get_currencies.return_value = reply

        with patch(
            "taurus_protect.services.currency_service.currencies_from_dto",
            return_value=[MagicMock(), MagicMock()],
        ):
            result = service.list()

        assert len(result) == 2

    def test_list_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.currencies = None
        reply.result = None
        api.wallet_service_get_currencies.return_value = reply

        result = service.list()
        assert result == []

    def test_list_passes_show_disabled(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.currencies = None
        reply.result = None
        api.wallet_service_get_currencies.return_value = reply

        service.list(show_disabled=True)

        api.wallet_service_get_currencies.assert_called_once_with(
            show_disabled=True,
            include_logo=None,
        )


class TestGet:
    """Tests for CurrencyService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        currencies_api = MagicMock()
        service = CurrencyService(api_client=api_client, currencies_api=currencies_api)
        return service, currencies_api

    def test_get_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="currency_id"):
            service.get("")

    def test_get_raises_not_found(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.currencies = None
        reply.result = None
        api.wallet_service_get_currencies.return_value = reply

        with patch(
            "taurus_protect.services.currency_service.currencies_from_dto",
            return_value=[],
        ):
            # list() returns empty, so get() won't find the currency
            with pytest.raises(NotFoundError, match="not found"):
                service.get("nonexistent")

    def test_get_returns_matching_currency(self) -> None:
        service, api = self._make_service()

        mock_currency = MagicMock()
        mock_currency.id = "ETH"

        reply = MagicMock()
        reply.currencies = [MagicMock()]
        reply.result = None
        api.wallet_service_get_currencies.return_value = reply

        with patch(
            "taurus_protect.services.currency_service.currencies_from_dto",
            return_value=[mock_currency],
        ):
            result = service.get("ETH")

        assert result is mock_currency


class TestGetByBlockchain:
    """Tests for CurrencyService.get_by_blockchain()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        currencies_api = MagicMock()
        service = CurrencyService(api_client=api_client, currencies_api=currencies_api)
        return service, currencies_api

    def test_get_by_blockchain_raises_for_empty_blockchain(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="blockchain"):
            service.get_by_blockchain("", "mainnet")

    def test_get_by_blockchain_raises_for_empty_network(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="network"):
            service.get_by_blockchain("ETH", "")

    def test_get_by_blockchain_returns_currency(self) -> None:
        service, api = self._make_service()

        mock_currency = MagicMock()
        reply = MagicMock()
        reply.currency = MagicMock()
        reply.result = None
        api.wallet_service_get_currency.return_value = reply

        with patch(
            "taurus_protect.services.currency_service.currency_from_dto",
            return_value=mock_currency,
        ):
            result = service.get_by_blockchain("ETH", "mainnet")

        assert result is mock_currency

    def test_get_by_blockchain_raises_not_found(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.currency = None
        reply.result = None
        api.wallet_service_get_currency.return_value = reply

        with patch(
            "taurus_protect.services.currency_service.currency_from_dto",
            return_value=None,
        ):
            with pytest.raises(NotFoundError, match="not found"):
                service.get_by_blockchain("ETH", "mainnet")


class TestGetBaseCurrency:
    """Tests for CurrencyService.get_base_currency()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        currencies_api = MagicMock()
        service = CurrencyService(api_client=api_client, currencies_api=currencies_api)
        return service, currencies_api

    def test_get_base_currency_returns_currency(self) -> None:
        service, api = self._make_service()

        mock_currency = MagicMock()
        reply = MagicMock()
        reply.currency = MagicMock()
        reply.result = None
        api.wallet_service_get_base_currency.return_value = reply

        with patch(
            "taurus_protect.services.currency_service.currency_from_dto",
            return_value=mock_currency,
        ):
            result = service.get_base_currency()

        assert result is mock_currency

    def test_get_base_currency_raises_not_found(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.currency = None
        reply.result = None
        api.wallet_service_get_base_currency.return_value = reply

        with patch(
            "taurus_protect.services.currency_service.currency_from_dto",
            return_value=None,
        ):
            with pytest.raises(NotFoundError, match="Base currency not configured"):
                service.get_base_currency()
