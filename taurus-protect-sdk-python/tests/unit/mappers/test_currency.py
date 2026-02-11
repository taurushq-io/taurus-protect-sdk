"""Unit tests for currency mapper."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.currency import (
    asset_balance_from_dto,
    asset_balances_from_dto,
    currencies_from_dto,
    currency_from_dto,
    nft_collection_balance_from_dto,
    nft_collection_balances_from_dto,
)


class TestCurrencyFromDto:
    """Tests for currency_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="BTC",
            name="Bitcoin",
            symbol="BTC",
            blockchain="BTC",
            network="mainnet",
            decimals=8,
            logo_url="https://example.com/btc.png",
            logo=None,
            enabled=True,
            is_token=False,
            contract_address=None,
            token_contract_address=None,
            display_name="Bitcoin (BTC)",
            type="NATIVE",
            coin_type_index="0",
            token_id="",
            wlca_id=None,
            is_erc20=False,
            is_fa12=False,
            is_fa20=False,
            is_nft=False,
            is_utxo_based=True,
            is_account_based=False,
            is_fiat=False,
            has_staking=False,
        )
        result = currency_from_dto(dto)
        assert result is not None
        assert result.id == "BTC"
        assert result.name == "Bitcoin"
        assert result.symbol == "BTC"
        assert result.decimals == 8
        assert result.is_utxo_based is True
        assert result.is_token is False
        assert result.enabled is True

    def test_returns_none_for_none(self) -> None:
        assert currency_from_dto(None) is None

    def test_handles_logo_fallback(self) -> None:
        dto = SimpleNamespace(
            id="ETH",
            name="Ethereum",
            symbol="ETH",
            blockchain="ETH",
            network="mainnet",
            decimals=18,
            logo_url=None,
            logo="https://example.com/eth.png",
            enabled=None,
            is_token=None,
            contract_address=None,
            token_contract_address=None,
            display_name=None,
            type=None,
            coin_type_index=None,
            token_id=None,
            wlca_id=None,
            is_erc20=None,
            is_fa12=None,
            is_fa20=None,
            is_nft=None,
            is_utxo_based=None,
            is_account_based=None,
            is_fiat=None,
            has_staking=None,
        )
        result = currency_from_dto(dto)
        assert result is not None
        assert result.logo_url == "https://example.com/eth.png"
        # enabled defaults to True when None
        assert result.enabled is True

    def test_contract_address_fallback(self) -> None:
        dto = SimpleNamespace(
            id="USDT",
            name="Tether",
            symbol="USDT",
            blockchain="ETH",
            network="mainnet",
            decimals=6,
            logo_url=None,
            logo=None,
            enabled=True,
            is_token=True,
            contract_address=None,
            token_contract_address="0xdAC17F958D2ee523a2206206994597C13D831ec7",
            display_name=None,
            type=None,
            coin_type_index=None,
            token_id=None,
            wlca_id=None,
            is_erc20=True,
            is_fa12=None,
            is_fa20=None,
            is_nft=None,
            is_utxo_based=None,
            is_account_based=None,
            is_fiat=None,
            has_staking=None,
        )
        result = currency_from_dto(dto)
        assert result is not None
        assert result.contract_address == "0xdAC17F958D2ee523a2206206994597C13D831ec7"


class TestCurrenciesFromDto:
    """Tests for currencies_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="BTC", name="Bitcoin", symbol="BTC", blockchain="BTC",
                network=None, decimals=8, logo_url=None, logo=None,
                enabled=True, is_token=False, contract_address=None,
                token_contract_address=None, display_name=None, type=None,
                coin_type_index=None, token_id=None, wlca_id=None,
                is_erc20=None, is_fa12=None, is_fa20=None, is_nft=None,
                is_utxo_based=None, is_account_based=None, is_fiat=None,
                has_staking=None,
            ),
        ]
        result = currencies_from_dto(dtos)
        assert len(result) == 1

    def test_returns_empty_for_none(self) -> None:
        assert currencies_from_dto(None) == []


class TestAssetBalanceFromDto:
    """Tests for asset_balance_from_dto (currency module) function."""

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
        result = asset_balance_from_dto(dto)
        assert result is not None
        assert result.currency == "ETH"
        assert result.balance == "1000"
        assert result.blockchain == "ETH"
        assert result.wallet_id == "w-1"

    def test_returns_none_for_none(self) -> None:
        assert asset_balance_from_dto(None) is None


class TestAssetBalancesFromDto:
    """Tests for asset_balances_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert asset_balances_from_dto(None) == []


class TestNftCollectionBalanceFromDto:
    """Tests for nft_collection_balance_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(
            collection_name="Bored Apes",
            name=None,
            contract_address="0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D",
            blockchain="ETH",
            network="mainnet",
            count=10,
            balance=None,
        )
        result = nft_collection_balance_from_dto(dto)
        assert result is not None
        assert result.collection_name == "Bored Apes"
        assert result.count == 10

    def test_name_fallback(self) -> None:
        dto = SimpleNamespace(
            collection_name=None,
            name="CryptoPunks",
            contract_address="0x1234",
            blockchain="ETH",
            network=None,
            count=None,
            balance=5,
        )
        result = nft_collection_balance_from_dto(dto)
        assert result is not None
        assert result.collection_name == "CryptoPunks"
        assert result.count == 5

    def test_returns_none_for_none(self) -> None:
        assert nft_collection_balance_from_dto(None) is None


class TestNftCollectionBalancesFromDto:
    """Tests for nft_collection_balances_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert nft_collection_balances_from_dto(None) == []
