"""Unit tests for WalletService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import APIError, NotFoundError
from taurus_protect.services.wallet_service import WalletService


class TestGet:
    """Tests for WalletService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_get_returns_wallet(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.wallet_service_get_wallet_v2.return_value = reply

        mock_wallet = MagicMock()
        with patch(
            "taurus_protect.services.wallet_service.wallet_from_dto",
            return_value=mock_wallet,
        ):
            result = service.get(1)

        assert result is mock_wallet
        api.wallet_service_get_wallet_v2.assert_called_once_with("1")

    def test_get_raises_value_error_for_non_positive_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.get(0)

    def test_get_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()
        reply = MagicMock()
        reply.result = None
        api.wallet_service_get_wallet_v2.return_value = reply

        with pytest.raises(NotFoundError):
            service.get(1)


class TestList:
    """Tests for WalletService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_list_returns_wallets_and_pagination(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "10"
        reply.offset = "0"
        api.wallet_service_get_wallets_v2.return_value = reply

        with patch(
            "taurus_protect.services.wallet_service.wallets_from_dto",
            return_value=[MagicMock()],
        ):
            wallets, pagination = service.list(limit=50)

        assert len(wallets) == 1
        api.wallet_service_get_wallets_v2.assert_called_once()

    def test_list_raises_for_invalid_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_list_raises_for_negative_offset(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_list_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        reply.total_items = None
        reply.offset = None
        api.wallet_service_get_wallets_v2.return_value = reply

        wallets, pagination = service.list()
        assert wallets == []


class TestCreate:
    """Tests for WalletService.create()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_create_raises_for_none_request(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="request cannot be None"):
            service.create(None)

    def test_create_wallet_raises_for_empty_blockchain(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="blockchain"):
            service.create_wallet(blockchain="", network="mainnet", name="test")

    def test_create_wallet_raises_for_empty_name(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="name"):
            service.create_wallet(blockchain="ETH", network="mainnet", name="")

    def test_create_wallet_returns_wallet(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.wallet_service_create_wallet.return_value = reply

        mock_wallet = MagicMock()
        with patch(
            "taurus_protect.services.wallet_service.wallet_from_create_dto",
            return_value=mock_wallet,
        ):
            result = service.create_wallet("ETH", "mainnet", "Test Wallet")

        assert result is mock_wallet


class TestGetByName:
    """Tests for WalletService.get_by_name()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_get_by_name_raises_for_empty_name(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="name"):
            service.get_by_name("")

    def test_get_by_name_passes_name_to_api(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "1"
        reply.offset = "0"
        api.wallet_service_get_wallets_v2.return_value = reply

        with patch(
            "taurus_protect.services.wallet_service.wallets_from_dto",
            return_value=[MagicMock()],
        ):
            wallets, _ = service.get_by_name("Test")

        call_kwargs = api.wallet_service_get_wallets_v2.call_args
        assert call_kwargs[1]["name"] == "Test" or call_kwargs.kwargs.get("name") == "Test"


class TestCreateAttribute:
    """Tests for WalletService.create_attribute()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_create_attribute_raises_for_non_positive_wallet_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.create_attribute(0, "key", "value")

    def test_create_attribute_raises_for_empty_key(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="key"):
            service.create_attribute(1, "", "value")

    def test_create_attribute_calls_api(self) -> None:
        service, api = self._make_service()

        service.create_attribute(1, "key", "value")

        api.wallet_service_create_wallet_attributes.assert_called_once()


class TestGetBalanceHistory:
    """Tests for WalletService.get_balance_history()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_get_balance_history_validates_wallet_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.get_balance_history(0, 24)

    def test_get_balance_history_validates_interval(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="interval_hours must be positive"):
            service.get_balance_history(1, 0)

    def test_get_balance_history_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.wallet_service_get_wallet_balance_history.return_value = reply

        result = service.get_balance_history(1, 24)
        assert result == []


class TestGetTokens:
    """Tests for WalletService.get_tokens()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        wallets_api = MagicMock()
        service = WalletService(api_client=api_client, wallets_api=wallets_api)
        return service, wallets_api

    def test_get_tokens_validates_wallet_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="wallet_id must be positive"):
            service.get_tokens(0)

    def test_get_tokens_validates_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.get_tokens(1, limit=0)

    def test_get_tokens_returns_empty_when_no_balances(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.balances = None
        api.wallet_service_get_wallet_tokens.return_value = reply

        result = service.get_tokens(1)
        assert result == []
