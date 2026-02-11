"""Unit tests for blockchain mapper functions."""

from types import SimpleNamespace

import pytest

from taurus_protect.services.blockchain_service import blockchain_from_dto, blockchains_from_dto


class TestBlockchainFromDto:
    """Tests for blockchain_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="eth-mainnet",
            blockchain="ETH",
            network="mainnet",
            display_name="Ethereum Mainnet",
            enabled=True,
            native_currency="ETH",
            block_height=19000000,
            block_time=12,
            confirmations_required=12,
        )
        result = blockchain_from_dto(dto)
        assert result is not None
        assert result.id == "eth-mainnet"
        assert result.name == "ETH"
        assert result.network == "mainnet"
        assert result.display_name == "Ethereum Mainnet"
        assert result.enabled is True
        assert result.native_currency == "ETH"
        assert result.block_height == 19000000
        assert result.block_time == 12
        assert result.confirmations_required == 12

    def test_returns_none_for_none(self) -> None:
        assert blockchain_from_dto(None) is None

    def test_generates_id_from_blockchain_and_network(self) -> None:
        dto = SimpleNamespace(
            blockchain="BTC",
            network="testnet",
            display_name=None,
            enabled=True,
            native_currency="BTC",
            block_height=None,
            block_time=None,
            confirmations_required=None,
        )
        result = blockchain_from_dto(dto)
        assert result is not None
        assert result.id == "BTC_testnet"
        assert result.name == "BTC"
        assert result.network == "testnet"

    def test_generates_id_without_network(self) -> None:
        dto = SimpleNamespace(
            blockchain="SOL",
            network=None,
            display_name=None,
            enabled=None,
            native_currency=None,
            block_height=None,
            block_time=None,
            confirmations_required=None,
        )
        result = blockchain_from_dto(dto)
        assert result is not None
        assert result.id == "SOL"
        assert result.name == "SOL"
        assert result.network == ""

    def test_handles_camel_case_field_names(self) -> None:
        dto = SimpleNamespace(
            blockchain="ETH",
            network="mainnet",
            displayName="Ethereum",
            enabled=True,
            nativeCurrency="ETH",
            blockHeight=18500000,
            blockTime=12,
            confirmationsRequired=6,
        )
        result = blockchain_from_dto(dto)
        assert result is not None
        assert result.display_name == "Ethereum"
        assert result.native_currency == "ETH"
        assert result.block_height == 18500000
        assert result.block_time == 12
        assert result.confirmations_required == 6

    def test_uses_name_field_as_fallback_for_blockchain(self) -> None:
        dto = SimpleNamespace(
            name="ADA",
            network="mainnet",
            display_name=None,
            enabled=True,
            native_currency="ADA",
            block_height=None,
            block_time=None,
            confirmations_required=None,
        )
        result = blockchain_from_dto(dto)
        assert result is not None
        assert result.name == "ADA"

    def test_handles_none_optional_fields(self) -> None:
        dto = SimpleNamespace(
            blockchain="DOT",
            network=None,
            display_name=None,
            enabled=None,
            native_currency=None,
            block_height=None,
            block_time=None,
            confirmations_required=None,
        )
        result = blockchain_from_dto(dto)
        assert result is not None
        assert result.enabled is False
        assert result.native_currency is None
        assert result.block_height is None
        assert result.block_time is None
        assert result.confirmations_required is None


class TestBlockchainsFromDto:
    """Tests for blockchains_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="eth-1", blockchain="ETH", network="mainnet",
                display_name=None, enabled=True, native_currency="ETH",
                block_height=None, block_time=None, confirmations_required=None,
            ),
            SimpleNamespace(
                id="btc-1", blockchain="BTC", network="mainnet",
                display_name=None, enabled=True, native_currency="BTC",
                block_height=None, block_time=None, confirmations_required=None,
            ),
        ]
        result = blockchains_from_dto(dtos)
        assert len(result) == 2

    def test_returns_empty_for_none(self) -> None:
        assert blockchains_from_dto(None) == []

    def test_returns_empty_for_empty_list(self) -> None:
        assert blockchains_from_dto([]) == []

    def test_filters_out_none_dtos(self) -> None:
        dtos = [
            SimpleNamespace(
                id="eth-1", blockchain="ETH", network="mainnet",
                display_name=None, enabled=True, native_currency="ETH",
                block_height=None, block_time=None, confirmations_required=None,
            ),
            None,
        ]
        result = blockchains_from_dto(dtos)
        assert len(result) == 1
