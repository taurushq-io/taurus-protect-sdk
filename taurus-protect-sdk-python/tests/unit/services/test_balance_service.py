"""Unit tests for BalanceService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.services.balance_service import BalanceService


class TestList:
    """Tests for BalanceService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        balances_api = MagicMock()
        service = BalanceService(api_client=api_client, balances_api=balances_api)
        return service, balances_api

    def test_list_returns_balances_and_pagination(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.balances = [MagicMock()]
        reply.result = None
        reply.total_items = "10"
        api.wallet_service_get_balances.return_value = reply

        with patch(
            "taurus_protect.services.balance_service.asset_balances_from_dto",
            return_value=[MagicMock()],
        ):
            balances, pagination = service.list()

        assert len(balances) == 1

    def test_list_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.balances = None
        reply.result = None
        reply.total_items = None
        api.wallet_service_get_balances.return_value = reply

        balances, pagination = service.list()
        assert balances == []

    def test_list_raises_for_invalid_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_list_raises_for_negative_offset(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_list_passes_currency_filter(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.balances = None
        reply.result = None
        reply.total_items = None
        api.wallet_service_get_balances.return_value = reply

        service.list(currency="ETH")

        call_kwargs = api.wallet_service_get_balances.call_args
        assert call_kwargs[1]["currency"] == "ETH" or call_kwargs.kwargs.get("currency") == "ETH"


class TestListNftCollections:
    """Tests for BalanceService.list_nft_collections()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        balances_api = MagicMock()
        service = BalanceService(api_client=api_client, balances_api=balances_api)
        return service, balances_api

    def test_list_nft_collections_raises_for_empty_blockchain(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="blockchain"):
            service.list_nft_collections(blockchain="", network="mainnet")

    def test_list_nft_collections_raises_for_empty_network(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="network"):
            service.list_nft_collections(blockchain="ETH", network="")

    def test_list_nft_collections_raises_for_invalid_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.list_nft_collections(blockchain="ETH", network="mainnet", limit=0)

    def test_list_nft_collections_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.collections = None
        reply.result = None
        reply.total_items = None
        api.wallet_service_get_nft_collection_balances.return_value = reply

        balances, pagination = service.list_nft_collections("ETH", "mainnet")
        assert balances == []
