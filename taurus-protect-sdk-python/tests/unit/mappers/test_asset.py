"""Unit tests for asset mapper functions."""

from types import SimpleNamespace

import pytest

from taurus_protect.services.asset_service import asset_from_dto, assets_from_dto


class TestAssetFromDto:
    """Tests for asset_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="BTC",
            name="Bitcoin",
            symbol="BTC",
            blockchain="BTC",
            network="mainnet",
            decimals=8,
            logo_url="https://example.com/btc.png",
            enabled=True,
            is_token=False,
            contract_address=None,
        )
        result = asset_from_dto(dto)
        assert result is not None
        assert result.id == "BTC"
        assert result.name == "Bitcoin"
        assert result.symbol == "BTC"
        assert result.blockchain == "BTC"
        assert result.network == "mainnet"
        assert result.decimals == 8
        assert result.logo_url == "https://example.com/btc.png"
        assert result.enabled is True
        assert result.is_token is False
        assert result.contract_address is None

    def test_returns_none_for_none(self) -> None:
        assert asset_from_dto(None) is None

    def test_handles_alternative_field_names(self) -> None:
        dto = SimpleNamespace(
            currency_id="ETH",
            name="Ethereum",
            currency="ETH",
            blockchain="ETH",
            network="mainnet",
            decimals="18",
            logoUrl="https://example.com/eth.png",
            enabled=None,
            isToken=True,
            contractAddress="0xabc123",
        )
        # Remove standard names so fallbacks are used
        result = asset_from_dto(dto)
        assert result is not None
        assert result.id == "ETH"
        assert result.symbol == "ETH"
        assert result.decimals == 18
        assert result.logo_url == "https://example.com/eth.png"
        assert result.is_token is True
        assert result.contract_address == "0xabc123"

    def test_handles_none_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="SOL",
            name=None,
            symbol=None,
            blockchain=None,
            network=None,
            decimals=None,
            logo_url=None,
            enabled=None,
            is_token=None,
            contract_address=None,
        )
        result = asset_from_dto(dto)
        assert result is not None
        assert result.id == "SOL"
        assert result.name is None
        assert result.decimals == 0
        assert result.enabled is False
        assert result.is_token is False

    def test_enabled_defaults_to_true_when_field_missing(self) -> None:
        dto = SimpleNamespace(
            id="ADA",
            name="Cardano",
            symbol="ADA",
            blockchain="ADA",
            network=None,
            decimals=6,
            logo_url=None,
        )
        result = asset_from_dto(dto)
        assert result is not None
        assert result.enabled is True


class TestAssetsFromDto:
    """Tests for assets_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="BTC", name="Bitcoin", symbol="BTC", blockchain="BTC",
                network=None, decimals=8, logo_url=None, enabled=True,
                is_token=False, contract_address=None,
            ),
            SimpleNamespace(
                id="ETH", name="Ethereum", symbol="ETH", blockchain="ETH",
                network=None, decimals=18, logo_url=None, enabled=True,
                is_token=False, contract_address=None,
            ),
        ]
        result = assets_from_dto(dtos)
        assert len(result) == 2
        assert result[0].id == "BTC"
        assert result[1].id == "ETH"

    def test_returns_empty_for_none(self) -> None:
        assert assets_from_dto(None) == []

    def test_returns_empty_for_empty_list(self) -> None:
        assert assets_from_dto([]) == []

    def test_filters_out_none_dtos(self) -> None:
        dtos = [
            SimpleNamespace(
                id="BTC", name="Bitcoin", symbol="BTC", blockchain="BTC",
                network=None, decimals=8, logo_url=None, enabled=True,
                is_token=False, contract_address=None,
            ),
            None,
        ]
        result = assets_from_dto(dtos)
        assert len(result) == 1
        assert result[0].id == "BTC"
