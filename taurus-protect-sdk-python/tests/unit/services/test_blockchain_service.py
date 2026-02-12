"""Unit tests for BlockchainService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import NotFoundError
from taurus_protect.services.blockchain_service import BlockchainService


class TestList:
    """Tests for BlockchainService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        blockchain_api = MagicMock()
        service = BlockchainService(api_client=api_client, blockchain_api=blockchain_api)
        return service, blockchain_api

    def test_list_returns_blockchains(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.blockchains = None
        api.blockchain_service_get_blockchains.return_value = reply

        with patch(
            "taurus_protect.services.blockchain_service.blockchains_from_dto",
            return_value=[MagicMock()],
        ):
            result = service.list()

        assert len(result) == 1

    def test_list_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        reply.blockchains = None
        api.blockchain_service_get_blockchains.return_value = reply

        result = service.list()
        assert result == []

    def test_list_passes_filters(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        reply.blockchains = None
        api.blockchain_service_get_blockchains.return_value = reply

        service.list(blockchain="ETH", network="mainnet")

        api.blockchain_service_get_blockchains.assert_called_once_with(
            blockchain="ETH",
            network="mainnet",
            include_block_height=None,
        )


class TestGet:
    """Tests for BlockchainService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        blockchain_api = MagicMock()
        service = BlockchainService(api_client=api_client, blockchain_api=blockchain_api)
        return service, blockchain_api

    def test_get_raises_for_empty_blockchain(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="blockchain"):
            service.get("")

    def test_get_returns_blockchain(self) -> None:
        service, api = self._make_service()

        mock_blockchain = MagicMock()
        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.blockchains = None
        api.blockchain_service_get_blockchains.return_value = reply

        with patch(
            "taurus_protect.services.blockchain_service.blockchains_from_dto",
            return_value=[mock_blockchain],
        ):
            result = service.get("ETH")

        assert result is mock_blockchain

    def test_get_raises_not_found(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        reply.blockchains = None
        api.blockchain_service_get_blockchains.return_value = reply

        with pytest.raises(NotFoundError, match="not found"):
            service.get("NONEXISTENT")


class TestGetById:
    """Tests for BlockchainService.get_by_id()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        blockchain_api = MagicMock()
        service = BlockchainService(api_client=api_client, blockchain_api=blockchain_api)
        return service, blockchain_api

    def test_get_by_id_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="blockchain_id"):
            service.get_by_id("")

    def test_get_by_id_parses_composite_id(self) -> None:
        service, api = self._make_service()

        mock_blockchain = MagicMock()
        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.blockchains = None
        api.blockchain_service_get_blockchains.return_value = reply

        with patch(
            "taurus_protect.services.blockchain_service.blockchains_from_dto",
            return_value=[mock_blockchain],
        ):
            result = service.get_by_id("ETH_mainnet")

        assert result is mock_blockchain
        # Verify that the API was called with the parsed blockchain and network
        call_kwargs = api.blockchain_service_get_blockchains.call_args
        assert call_kwargs[1]["blockchain"] == "ETH" or call_kwargs.kwargs.get("blockchain") == "ETH"

    def test_get_by_id_defaults_network_to_mainnet(self) -> None:
        service, api = self._make_service()

        mock_blockchain = MagicMock()
        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.blockchains = None
        api.blockchain_service_get_blockchains.return_value = reply

        with patch(
            "taurus_protect.services.blockchain_service.blockchains_from_dto",
            return_value=[mock_blockchain],
        ):
            service.get_by_id("BTC")

        call_kwargs = api.blockchain_service_get_blockchains.call_args
        assert call_kwargs[1]["network"] == "mainnet" or call_kwargs.kwargs.get("network") == "mainnet"
