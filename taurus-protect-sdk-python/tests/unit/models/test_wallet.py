"""Tests for wallet domain models."""

from datetime import datetime, timezone

import pytest
from pydantic import ValidationError

from taurus_protect.models.wallet import (
    CreateWalletRequest,
    ListWalletsOptions,
    Wallet,
    WalletAttribute,
)


class TestWalletAttribute:
    """Tests for WalletAttribute model."""

    def test_create_attribute(self) -> None:
        """Test creating a wallet attribute."""
        attr = WalletAttribute(id="attr1", key="env", value="production")

        assert attr.id == "attr1"
        assert attr.key == "env"
        assert attr.value == "production"

    def test_attribute_is_frozen(self) -> None:
        """Test that attribute is immutable."""
        attr = WalletAttribute(id="attr1", key="env", value="production")

        with pytest.raises(ValidationError):
            attr.key = "new_key"  # type: ignore

    def test_missing_required_fields_raises(self) -> None:
        """Test that missing required fields raise ValidationError."""
        with pytest.raises(ValidationError):
            WalletAttribute(id="attr1")  # type: ignore


class TestWallet:
    """Tests for Wallet model."""

    def test_create_minimal_wallet(self) -> None:
        """Test creating wallet with minimal required fields."""
        wallet = Wallet(id="wallet1", name="Test Wallet", currency="BTC")

        assert wallet.id == "wallet1"
        assert wallet.name == "Test Wallet"
        assert wallet.currency == "BTC"
        assert wallet.blockchain == ""
        assert wallet.network == ""
        assert wallet.is_omnibus is False
        assert wallet.disabled is False
        assert wallet.addresses_count == 0
        assert wallet.attributes == []

    def test_create_full_wallet(self) -> None:
        """Test creating wallet with all fields."""
        now = datetime.now(timezone.utc)
        wallet = Wallet(
            id="wallet1",
            name="Trading Wallet",
            currency="ETH",
            blockchain="Ethereum",
            network="mainnet",
            is_omnibus=True,
            disabled=False,
            comment="Primary trading account",
            customer_id="cust123",
            external_wallet_id="ext456",
            visibility_group_id="vg789",
            account_path="m/44'/60'/0'",
            addresses_count=5,
            created_at=now,
            updated_at=now,
            attributes=[
                WalletAttribute(id="a1", key="type", value="hot"),
            ],
        )

        assert wallet.id == "wallet1"
        assert wallet.name == "Trading Wallet"
        assert wallet.currency == "ETH"
        assert wallet.blockchain == "Ethereum"
        assert wallet.network == "mainnet"
        assert wallet.is_omnibus is True
        assert wallet.comment == "Primary trading account"
        assert wallet.customer_id == "cust123"
        assert wallet.external_wallet_id == "ext456"
        assert wallet.visibility_group_id == "vg789"
        assert wallet.account_path == "m/44'/60'/0'"
        assert wallet.addresses_count == 5
        assert wallet.created_at == now
        assert wallet.updated_at == now
        assert len(wallet.attributes) == 1
        assert wallet.attributes[0].key == "type"

    def test_wallet_is_frozen(self) -> None:
        """Test that wallet is immutable."""
        wallet = Wallet(id="wallet1", name="Test", currency="BTC")

        with pytest.raises(ValidationError):
            wallet.name = "New Name"  # type: ignore

    def test_missing_required_fields_raises(self) -> None:
        """Test that missing required fields raise ValidationError."""
        with pytest.raises(ValidationError):
            Wallet(id="wallet1")  # type: ignore

    def test_wallet_serialization(self) -> None:
        """Test wallet serialization to dict."""
        wallet = Wallet(id="wallet1", name="Test", currency="BTC")
        data = wallet.model_dump()

        assert data["id"] == "wallet1"
        assert data["name"] == "Test"
        assert data["currency"] == "BTC"


class TestCreateWalletRequest:
    """Tests for CreateWalletRequest model."""

    def test_create_minimal_request(self) -> None:
        """Test creating minimal wallet request."""
        request = CreateWalletRequest(
            blockchain="ETH",
            network="mainnet",
            name="New Wallet",
        )

        assert request.blockchain == "ETH"
        assert request.network == "mainnet"
        assert request.name == "New Wallet"
        assert request.is_omnibus is False
        assert request.comment is None
        assert request.customer_id is None

    def test_create_full_request(self) -> None:
        """Test creating wallet request with all fields."""
        request = CreateWalletRequest(
            blockchain="BTC",
            network="testnet",
            name="Test Wallet",
            is_omnibus=True,
            comment="Test comment",
            customer_id="cust123",
        )

        assert request.blockchain == "BTC"
        assert request.network == "testnet"
        assert request.name == "Test Wallet"
        assert request.is_omnibus is True
        assert request.comment == "Test comment"
        assert request.customer_id == "cust123"

    def test_missing_required_fields_raises(self) -> None:
        """Test that missing required fields raise ValidationError."""
        with pytest.raises(ValidationError):
            CreateWalletRequest(blockchain="ETH")  # type: ignore


class TestListWalletsOptions:
    """Tests for ListWalletsOptions model."""

    def test_default_options(self) -> None:
        """Test default options values."""
        options = ListWalletsOptions()

        assert options.limit == 50
        assert options.offset == 0
        assert options.currency is None
        assert options.query is None
        assert options.exclude_disabled is False

    def test_custom_options(self) -> None:
        """Test custom options values."""
        options = ListWalletsOptions(
            limit=100,
            offset=50,
            currency="BTC",
            query="trading",
            exclude_disabled=True,
        )

        assert options.limit == 100
        assert options.offset == 50
        assert options.currency == "BTC"
        assert options.query == "trading"
        assert options.exclude_disabled is True

    def test_limit_validation_min(self) -> None:
        """Test limit minimum validation."""
        with pytest.raises(ValidationError):
            ListWalletsOptions(limit=0)

    def test_limit_validation_max(self) -> None:
        """Test limit maximum validation."""
        with pytest.raises(ValidationError):
            ListWalletsOptions(limit=1001)

    def test_offset_validation_min(self) -> None:
        """Test offset minimum validation."""
        with pytest.raises(ValidationError):
            ListWalletsOptions(offset=-1)
