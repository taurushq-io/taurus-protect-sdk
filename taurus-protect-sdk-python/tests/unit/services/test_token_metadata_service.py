"""Unit tests for TokenMetadataService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.token_metadata_service import TokenMetadataService


class TestTokenMetadataServiceGet:
    """Tests for TokenMetadataService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        token_metadata_api = MagicMock()
        service = TokenMetadataService(
            api_client=api_client, token_metadata_api=token_metadata_api
        )
        return service, token_metadata_api

    def test_raises_on_empty_blockchain(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="blockchain"):
            service.get(blockchain="", contract_address="0x123")

    def test_raises_on_empty_contract_address(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="contract_address"):
            service.get(blockchain="ETH", contract_address="")

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.token_metadata_service_get_evmerc_token_metadata.return_value = resp

        service.get(blockchain="ETH", contract_address="0xabc")

        api.token_metadata_service_get_evmerc_token_metadata.assert_called_once_with(
            network="mainnet",
            contract="0xabc",
            token="0",
            with_data=None,
            blockchain="ETH",
        )


class TestTokenMetadataServiceGetERC:
    """Tests for TokenMetadataService.get_erc()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        token_metadata_api = MagicMock()
        service = TokenMetadataService(
            api_client=api_client, token_metadata_api=token_metadata_api
        )
        return service, token_metadata_api

    def test_raises_on_empty_network(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="network"):
            service.get_erc(network="", contract_address="0x123", token_id="42")

    def test_raises_on_empty_contract(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="contract_address"):
            service.get_erc(network="mainnet", contract_address="", token_id="42")

    def test_raises_on_empty_token_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="token_id"):
            service.get_erc(network="mainnet", contract_address="0x123", token_id="")


class TestTokenMetadataServiceGetFA:
    """Tests for TokenMetadataService.get_fa()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        token_metadata_api = MagicMock()
        service = TokenMetadataService(
            api_client=api_client, token_metadata_api=token_metadata_api
        )
        return service, token_metadata_api

    def test_raises_on_empty_network(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="network"):
            service.get_fa(network="", contract_address="KT1abc")

    def test_raises_on_empty_contract(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="contract_address"):
            service.get_fa(network="mainnet", contract_address="")

    def test_calls_api(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.token_metadata_service_get_fa_token_metadata.return_value = resp

        service.get_fa(network="mainnet", contract_address="KT1abc")

        api.token_metadata_service_get_fa_token_metadata.assert_called_once_with(
            network="mainnet",
            contract="KT1abc",
            token="0",
            with_data=None,
        )


class TestTokenMetadataServiceGetCryptoPunk:
    """Tests for TokenMetadataService.get_crypto_punk()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        token_metadata_api = MagicMock()
        service = TokenMetadataService(
            api_client=api_client, token_metadata_api=token_metadata_api
        )
        return service, token_metadata_api

    def test_raises_on_empty_network(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="network"):
            service.get_crypto_punk(network="", contract_address="0x123", punk_id="42")

    def test_raises_on_empty_contract(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="contract_address"):
            service.get_crypto_punk(network="mainnet", contract_address="", punk_id="42")

    def test_raises_on_empty_punk_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="punk_id"):
            service.get_crypto_punk(
                network="mainnet", contract_address="0x123", punk_id=""
            )
