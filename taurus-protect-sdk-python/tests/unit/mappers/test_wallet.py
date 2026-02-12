"""Unit tests for wallet mapper."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.wallet import (
    asset_balance_from_dto,
    asset_from_dto,
    balance_from_dto,
    balance_history_point_from_dto,
    wallet_attribute_from_dto,
    wallet_from_create_dto,
    wallet_from_dto,
    wallets_from_dto,
)


class TestWalletFromDto:
    """Tests for wallet_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="123",
            name="My Wallet",
            currency="ETH",
            blockchain="ETH",
            network="mainnet",
            balance=SimpleNamespace(
                total_confirmed="1000",
                total_unconfirmed="500",
                available_confirmed="800",
                available_unconfirmed="400",
                reserved_confirmed="200",
                reserved_unconfirmed="100",
            ),
            is_omnibus=False,
            disabled=False,
            comment="Test wallet",
            customer_id="cust-1",
            external_wallet_id="ext-123",
            visibility_group_id="vg-1",
            account_path="m/44'/60'/0'",
            addresses_count="5",
            creation_date="2024-01-01T00:00:00Z",
            update_date="2024-06-01T00:00:00Z",
            attributes=[],
            currency_info=None,
        )
        result = wallet_from_dto(dto)
        assert result is not None
        assert result.id == "123"
        assert result.name == "My Wallet"
        assert result.blockchain == "ETH"
        assert result.is_omnibus is False
        assert result.disabled is False
        assert result.addresses_count == 5
        assert result.balance is not None
        assert result.balance.total_confirmed == "1000"

    def test_returns_none_for_none(self) -> None:
        assert wallet_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="1",
            name="W",
            currency="BTC",
            blockchain=None,
            network=None,
            balance=None,
            is_omnibus=None,
            disabled=None,
            comment=None,
            customer_id=None,
            external_wallet_id=None,
            visibility_group_id=None,
            account_path=None,
            addresses_count=None,
            creation_date=None,
            update_date=None,
            attributes=None,
            currency_info=None,
        )
        result = wallet_from_dto(dto)
        assert result is not None
        assert result.balance is None
        assert result.is_omnibus is False
        assert result.addresses_count == 0

    def test_maps_attributes(self) -> None:
        dto = SimpleNamespace(
            id="1", name="W", currency="BTC", blockchain=None, network=None,
            balance=None, is_omnibus=None, disabled=None, comment=None,
            customer_id=None, external_wallet_id=None, visibility_group_id=None,
            account_path=None, addresses_count=None, creation_date=None,
            update_date=None,
            attributes=[
                SimpleNamespace(
                    id="a1", key="dept", value="treasury",
                    content_type=None, owner=None, type=None, subtype=None,
                    isfile=None, is_file=None,
                ),
            ],
            currency_info=None,
        )
        result = wallet_from_dto(dto)
        assert result is not None
        assert len(result.attributes) == 1
        assert result.attributes[0].key == "dept"


class TestWalletFromCreateDto:
    """Tests for wallet_from_create_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(
            id="456",
            name="New Wallet",
            currency="BTC",
            blockchain="BTC",
            network="mainnet",
            balance=None,
            is_omnibus=True,
            disabled=False,
            comment=None,
            customer_id=None,
            external_wallet_id=None,
            visibility_group_id=None,
            account_path=None,
            addresses_count=None,
            creation_date=None,
            update_date=None,
            attributes=None,
            currency_info=None,
        )
        result = wallet_from_create_dto(dto)
        assert result is not None
        assert result.id == "456"
        assert result.is_omnibus is True

    def test_returns_none_for_none(self) -> None:
        assert wallet_from_create_dto(None) is None


class TestWalletsFromDto:
    """Tests for wallets_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", name="W1", currency="BTC", blockchain=None, network=None,
                balance=None, is_omnibus=None, disabled=None, comment=None,
                customer_id=None, external_wallet_id=None, visibility_group_id=None,
                account_path=None, addresses_count=None, creation_date=None,
                update_date=None, attributes=None, currency_info=None,
            ),
            SimpleNamespace(
                id="2", name="W2", currency="ETH", blockchain=None, network=None,
                balance=None, is_omnibus=None, disabled=None, comment=None,
                customer_id=None, external_wallet_id=None, visibility_group_id=None,
                account_path=None, addresses_count=None, creation_date=None,
                update_date=None, attributes=None, currency_info=None,
            ),
        ]
        result = wallets_from_dto(dtos)
        assert len(result) == 2

    def test_returns_empty_for_none(self) -> None:
        assert wallets_from_dto(None) == []


class TestWalletAttributeFromDto:
    """Tests for wallet_attribute_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="a1", key="department", value="treasury",
            content_type="text/plain", owner="admin", type="custom",
            subtype=None, isfile=None, is_file=False,
        )
        result = wallet_attribute_from_dto(dto)
        assert result.id == "a1"
        assert result.key == "department"
        assert result.value == "treasury"
        assert result.content_type == "text/plain"
        assert result.is_file is False

    def test_isfile_camelcase(self) -> None:
        dto = SimpleNamespace(
            id="a2", key="k", value="v",
            content_type=None, owner=None, type=None,
            subtype=None, isfile=True, is_file=None,
        )
        result = wallet_attribute_from_dto(dto)
        assert result.is_file is True


class TestBalanceFromDto:
    """Tests for balance_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            total_confirmed="100",
            total_unconfirmed="50",
            available_confirmed="80",
            available_unconfirmed="30",
            reserved_confirmed="20",
            reserved_unconfirmed="20",
        )
        result = balance_from_dto(dto)
        assert result is not None
        assert result.total_confirmed == "100"
        assert result.available_confirmed == "80"
        assert result.reserved_confirmed == "20"

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
        assert result.timestamp is not None

    def test_returns_none_for_none(self) -> None:
        assert balance_history_point_from_dto(None) is None


class TestAssetFromDto:
    """Tests for asset_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="asset-1",
            symbol="ETH",
            name="Ethereum",
            decimals=18,
            blockchain="ETH",
        )
        result = asset_from_dto(dto)
        assert result is not None
        assert result.id == "asset-1"
        assert result.symbol == "ETH"
        assert result.decimals == 18

    def test_returns_none_for_none(self) -> None:
        assert asset_from_dto(None) is None


class TestAssetBalanceFromDto:
    """Tests for asset_balance_from_dto function."""

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
        result = asset_balance_from_dto(dto)
        assert result is not None
        assert result.asset is not None
        assert result.balance is not None
        assert result.balance.total_confirmed == "1000"

    def test_returns_none_for_none(self) -> None:
        assert asset_balance_from_dto(None) is None
