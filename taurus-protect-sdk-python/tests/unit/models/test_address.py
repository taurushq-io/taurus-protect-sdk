"""Tests for address domain models."""

from datetime import datetime, timezone

import pytest
from pydantic import ValidationError

from taurus_protect.models.address import (
    Address,
    AddressAttribute,
    CreateAddressRequest,
    ListAddressesOptions,
)


class TestAddressAttribute:
    """Tests for AddressAttribute model."""

    def test_create_attribute(self) -> None:
        """Test creating an address attribute."""
        attr = AddressAttribute(id="attr1", key="purpose", value="receiving")

        assert attr.id == "attr1"
        assert attr.key == "purpose"
        assert attr.value == "receiving"

    def test_attribute_is_frozen(self) -> None:
        """Test that attribute is immutable."""
        attr = AddressAttribute(id="attr1", key="purpose", value="receiving")

        with pytest.raises(ValidationError):
            attr.key = "new_key"  # type: ignore

    def test_missing_required_fields_raises(self) -> None:
        """Test that missing required fields raise ValidationError."""
        with pytest.raises(ValidationError):
            AddressAttribute(id="attr1")  # type: ignore


class TestAddress:
    """Tests for Address model."""

    def test_create_minimal_address(self) -> None:
        """Test creating address with minimal required fields."""
        addr = Address(
            id="addr1",
            wallet_id="wallet1",
            address="0x742d35Cc6634C0532925a3b844Bc9e7595f",
        )

        assert addr.id == "addr1"
        assert addr.wallet_id == "wallet1"
        assert addr.address == "0x742d35Cc6634C0532925a3b844Bc9e7595f"
        assert addr.currency == ""
        assert addr.disabled is False
        assert addr.can_use_all_funds is False
        assert addr.attributes == []
        assert addr.linked_whitelisted_address_ids == []

    def test_create_full_address(self) -> None:
        """Test creating address with all fields."""
        now = datetime.now(timezone.utc)
        addr = Address(
            id="addr1",
            wallet_id="wallet1",
            address="0x742d35Cc6634C0532925a3b844Bc9e7595f",
            alternate_address="bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
            label="Main Receiving",
            comment="Primary deposit address",
            currency="ETH",
            customer_id="cust123",
            external_address_id="ext456",
            address_path="m/44'/60'/0'/0/0",
            address_index="0",
            nonce="5",
            status="confirmed",
            signature="base64sig==",
            disabled=False,
            can_use_all_funds=True,
            created_at=now,
            updated_at=now,
            attributes=[
                AddressAttribute(id="a1", key="type", value="deposit"),
            ],
            linked_whitelisted_address_ids=["wla1", "wla2"],
        )

        assert addr.id == "addr1"
        assert addr.wallet_id == "wallet1"
        assert addr.address == "0x742d35Cc6634C0532925a3b844Bc9e7595f"
        assert addr.alternate_address == "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh"
        assert addr.label == "Main Receiving"
        assert addr.comment == "Primary deposit address"
        assert addr.currency == "ETH"
        assert addr.customer_id == "cust123"
        assert addr.external_address_id == "ext456"
        assert addr.address_path == "m/44'/60'/0'/0/0"
        assert addr.address_index == "0"
        assert addr.nonce == "5"
        assert addr.status == "confirmed"
        assert addr.signature == "base64sig=="
        assert addr.can_use_all_funds is True
        assert addr.created_at == now
        assert len(addr.attributes) == 1
        assert addr.linked_whitelisted_address_ids == ["wla1", "wla2"]

    def test_address_is_frozen(self) -> None:
        """Test that address is immutable."""
        addr = Address(id="addr1", wallet_id="wallet1", address="0x123")

        with pytest.raises(ValidationError):
            addr.label = "New Label"  # type: ignore

    def test_missing_required_fields_raises(self) -> None:
        """Test that missing required fields raise ValidationError."""
        with pytest.raises(ValidationError):
            Address(id="addr1")  # type: ignore

    def test_address_serialization(self) -> None:
        """Test address serialization to dict."""
        addr = Address(id="addr1", wallet_id="wallet1", address="0x123")
        data = addr.model_dump()

        assert data["id"] == "addr1"
        assert data["wallet_id"] == "wallet1"
        assert data["address"] == "0x123"


class TestCreateAddressRequest:
    """Tests for CreateAddressRequest model."""

    def test_create_minimal_request(self) -> None:
        """Test creating minimal address request."""
        request = CreateAddressRequest(
            wallet_id="wallet1",
            label="New Address",
        )

        assert request.wallet_id == "wallet1"
        assert request.label == "New Address"
        assert request.comment is None
        assert request.customer_id is None
        assert request.external_address_id is None
        assert request.address_type is None
        assert request.non_hardened_derivation is False

    def test_create_full_request(self) -> None:
        """Test creating address request with all fields."""
        request = CreateAddressRequest(
            wallet_id="wallet1",
            label="Test Address",
            comment="Test comment",
            customer_id="cust123",
            external_address_id="ext456",
            address_type="p2wpkh",
            non_hardened_derivation=True,
        )

        assert request.wallet_id == "wallet1"
        assert request.label == "Test Address"
        assert request.comment == "Test comment"
        assert request.customer_id == "cust123"
        assert request.external_address_id == "ext456"
        assert request.address_type == "p2wpkh"
        assert request.non_hardened_derivation is True

    def test_missing_required_fields_raises(self) -> None:
        """Test that missing required fields raise ValidationError."""
        with pytest.raises(ValidationError):
            CreateAddressRequest(wallet_id="wallet1")  # type: ignore


class TestListAddressesOptions:
    """Tests for ListAddressesOptions model."""

    def test_default_options(self) -> None:
        """Test default options values."""
        options = ListAddressesOptions()

        assert options.wallet_id is None
        assert options.limit == 50
        assert options.offset == 0
        assert options.query is None
        assert options.exclude_disabled is False

    def test_custom_options(self) -> None:
        """Test custom options values."""
        options = ListAddressesOptions(
            wallet_id="wallet1",
            limit=100,
            offset=50,
            query="deposit",
            exclude_disabled=True,
        )

        assert options.wallet_id == "wallet1"
        assert options.limit == 100
        assert options.offset == 50
        assert options.query == "deposit"
        assert options.exclude_disabled is True

    def test_limit_validation_min(self) -> None:
        """Test limit minimum validation."""
        with pytest.raises(ValidationError):
            ListAddressesOptions(limit=0)

    def test_limit_validation_max(self) -> None:
        """Test limit maximum validation."""
        with pytest.raises(ValidationError):
            ListAddressesOptions(limit=1001)

    def test_offset_validation_min(self) -> None:
        """Test offset minimum validation."""
        with pytest.raises(ValidationError):
            ListAddressesOptions(offset=-1)
