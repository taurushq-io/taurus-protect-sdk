"""Unit tests for balance-related mapper functions."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.currency import (
    asset_balance_from_dto as currency_asset_balance_from_dto,
    asset_balances_from_dto,
)
from taurus_protect.mappers.wallet import (
    asset_balance_from_dto as wallet_asset_balance_from_dto,
    balance_from_dto,
    balance_history_point_from_dto,
)


class TestBalanceFromDto:
    """Tests for balance_from_dto function (wallet module)."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            total_confirmed="1000000",
            total_unconfirmed="500000",
            available_confirmed="800000",
            available_unconfirmed="300000",
            reserved_confirmed="200000",
            reserved_unconfirmed="200000",
        )
        result = balance_from_dto(dto)
        assert result is not None
        assert result.total_confirmed == "1000000"
        assert result.total_unconfirmed == "500000"
        assert result.available_confirmed == "800000"
        assert result.available_unconfirmed == "300000"
        assert result.reserved_confirmed == "200000"
        assert result.reserved_unconfirmed == "200000"

    def test_returns_none_for_none(self) -> None:
        assert balance_from_dto(None) is None


class TestBalanceHistoryPointFromDto:
    """Tests for balance_history_point_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(
            timestamp="2024-01-01T00:00:00Z",
            total_confirmed="1000",
            total_unconfirmed="500",
            available_confirmed="800",
            available_unconfirmed="400",
        )
        result = balance_history_point_from_dto(dto)
        assert result is not None
        assert result.total_confirmed == "1000"
        assert result.total_unconfirmed == "500"
        assert result.available_confirmed == "800"
        assert result.available_unconfirmed == "400"
        assert result.timestamp is not None

    def test_returns_none_for_none(self) -> None:
        assert balance_history_point_from_dto(None) is None

    def test_handles_none_timestamp(self) -> None:
        dto = SimpleNamespace(
            timestamp=None,
            total_confirmed="0",
            total_unconfirmed="0",
            available_confirmed="0",
            available_unconfirmed="0",
        )
        result = balance_history_point_from_dto(dto)
        assert result is not None
        assert result.timestamp is None


class TestCurrencyAssetBalanceFromDto:
    """Tests for asset_balance_from_dto function (currency module)."""

    def test_maps_nested_fields(self) -> None:
        dto = SimpleNamespace(
            asset=SimpleNamespace(
                currency="ETH",
                symbol="ETH",
                currency_info=SimpleNamespace(
                    blockchain="ETH",
                    network="mainnet",
                ),
            ),
            balance=SimpleNamespace(
                total_confirmed="1000",
                total_unconfirmed="500",
            ),
            wallet_id="w-1",
            address_id="a-1",
        )
        result = currency_asset_balance_from_dto(dto)
        assert result is not None
        assert result.currency == "ETH"
        assert result.balance == "1000"
        assert result.blockchain == "ETH"
        assert result.wallet_id == "w-1"

    def test_returns_none_for_none(self) -> None:
        assert currency_asset_balance_from_dto(None) is None


class TestAssetBalancesFromDto:
    """Tests for asset_balances_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert asset_balances_from_dto(None) == []

    def test_returns_empty_for_empty_list(self) -> None:
        assert asset_balances_from_dto([]) == []


class TestWalletAssetBalanceFromDto:
    """Tests for asset_balance_from_dto function (wallet module)."""

    def test_maps_nested_fields(self) -> None:
        dto = SimpleNamespace(
            asset=SimpleNamespace(
                id="a1", symbol="ETH", name="Ethereum",
                decimals=18, blockchain="ETH",
            ),
            balance=SimpleNamespace(
                total_confirmed="1000",
                total_unconfirmed="500",
                available_confirmed="800",
                available_unconfirmed="400",
                reserved_confirmed="200",
                reserved_unconfirmed="100",
            ),
        )
        result = wallet_asset_balance_from_dto(dto)
        assert result is not None
        assert result.asset is not None
        assert result.balance is not None
        assert result.balance.total_confirmed == "1000"

    def test_returns_none_for_none(self) -> None:
        assert wallet_asset_balance_from_dto(None) is None
