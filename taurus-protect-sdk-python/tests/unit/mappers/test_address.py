"""Unit tests for address mapper."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.address import (
    address_attribute_from_dto,
    address_from_dto,
    addresses_from_dto,
)


class TestAddressFromDto:
    """Tests for address_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="456",
            wallet_id="123",
            address="0xabc123def456",
            alternate_address="0xalt",
            label="My Address",
            comment="Test address",
            currency="ETH",
            customer_id="cust-1",
            external_address_id="ext-456",
            address_path="m/44'/60'/0'/0/0",
            address_index="0",
            nonce="42",
            status="active",
            balance=None,
            signature="sig-data",
            disabled=False,
            can_use_all_funds=True,
            creation_date="2024-01-15T10:30:00Z",
            update_date="2024-06-01T10:30:00Z",
            attributes=[],
            linked_whitelisted_address_ids=["wl-1", "wl-2"],
        )
        result = address_from_dto(dto)
        assert result is not None
        assert result.id == "456"
        assert result.wallet_id == "123"
        assert result.address == "0xabc123def456"
        assert result.alternate_address == "0xalt"
        assert result.label == "My Address"
        assert result.currency == "ETH"
        assert result.customer_id == "cust-1"
        assert result.disabled is False
        assert result.can_use_all_funds is True
        assert result.linked_whitelisted_address_ids == ["wl-1", "wl-2"]

    def test_returns_none_for_none(self) -> None:
        assert address_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="1",
            wallet_id="10",
            address="0xminimal",
            alternate_address=None,
            label=None,
            comment=None,
            currency="ETH",
            customer_id=None,
            external_address_id=None,
            address_path=None,
            address_index=None,
            nonce=None,
            status=None,
            balance=None,
            signature=None,
            disabled=None,
            can_use_all_funds=None,
            creation_date=None,
            update_date=None,
            attributes=None,
            linked_whitelisted_address_ids=None,
        )
        result = address_from_dto(dto)
        assert result is not None
        assert result.id == "1"
        assert result.label is None
        assert result.disabled is False
        assert result.linked_whitelisted_address_ids == []

    def test_maps_balance(self) -> None:
        dto = SimpleNamespace(
            id="1",
            wallet_id="10",
            address="0xabc",
            alternate_address=None,
            label=None,
            comment=None,
            currency="ETH",
            customer_id=None,
            external_address_id=None,
            address_path=None,
            address_index=None,
            nonce=None,
            status=None,
            balance=SimpleNamespace(
                total_confirmed="1000",
                total_unconfirmed="500",
                available_confirmed="800",
                available_unconfirmed="400",
                reserved_confirmed="200",
                reserved_unconfirmed="100",
            ),
            signature=None,
            disabled=None,
            can_use_all_funds=None,
            creation_date=None,
            update_date=None,
            attributes=None,
            linked_whitelisted_address_ids=None,
        )
        result = address_from_dto(dto)
        assert result is not None
        assert result.balance is not None
        assert result.balance.total_confirmed == "1000"

    def test_maps_attributes(self) -> None:
        dto = SimpleNamespace(
            id="1",
            wallet_id="10",
            address="0xabc",
            alternate_address=None,
            label=None,
            comment=None,
            currency="ETH",
            customer_id=None,
            external_address_id=None,
            address_path=None,
            address_index=None,
            nonce=None,
            status=None,
            balance=None,
            signature=None,
            disabled=None,
            can_use_all_funds=None,
            creation_date=None,
            update_date=None,
            attributes=[
                SimpleNamespace(id="a1", key="tag", value="vip"),
            ],
            linked_whitelisted_address_ids=None,
        )
        result = address_from_dto(dto)
        assert result is not None
        assert len(result.attributes) == 1
        assert result.attributes[0].key == "tag"
        assert result.attributes[0].value == "vip"


class TestAddressesFromDto:
    """Tests for addresses_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", wallet_id="10", address="0xaaa", alternate_address=None,
                label=None, comment=None, currency="ETH", customer_id=None,
                external_address_id=None, address_path=None, address_index=None,
                nonce=None, status=None, balance=None, signature=None,
                disabled=None, can_use_all_funds=None, creation_date=None,
                update_date=None, attributes=None, linked_whitelisted_address_ids=None,
            ),
            SimpleNamespace(
                id="2", wallet_id="10", address="0xbbb", alternate_address=None,
                label=None, comment=None, currency="ETH", customer_id=None,
                external_address_id=None, address_path=None, address_index=None,
                nonce=None, status=None, balance=None, signature=None,
                disabled=None, can_use_all_funds=None, creation_date=None,
                update_date=None, attributes=None, linked_whitelisted_address_ids=None,
            ),
        ]
        result = addresses_from_dto(dtos)
        assert len(result) == 2

    def test_returns_empty_for_none(self) -> None:
        assert addresses_from_dto(None) == []

    def test_returns_empty_for_empty(self) -> None:
        assert addresses_from_dto([]) == []


class TestAddressAttributeFromDto:
    """Tests for address_attribute_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(id="a1", key="tag", value="important")
        result = address_attribute_from_dto(dto)
        assert result.id == "a1"
        assert result.key == "tag"
        assert result.value == "important"
