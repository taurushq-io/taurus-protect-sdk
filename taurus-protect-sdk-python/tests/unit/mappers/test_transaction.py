"""Unit tests for transaction mapper."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.transaction import map_transaction, map_transactions


class TestMapTransaction:
    """Tests for map_transaction function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="tx-1",
            request_id="req-1",
            wallet_id="w-1",
            address_id="a-1",
            currency="BTC",
            blockchain="BTC",
            tx_hash="0xabc",
            hash=None,
            block_height=100,
            block_number=None,
            block_hash="0xblock",
            amount="1.5",
            fee="0.001",
            direction="OUTGOING",
            status="CONFIRMED",
            confirmations=6,
            created_at="2024-01-15T10:30:00Z",
            creation_date=None,
            confirmed_at="2024-01-15T11:00:00Z",
            confirmation_date=None,
        )
        result = map_transaction(dto)
        assert result.id == "tx-1"
        assert result.request_id == "req-1"
        assert result.wallet_id == "w-1"
        assert result.currency == "BTC"
        assert result.tx_hash == "0xabc"
        assert result.block_height == 100
        assert result.amount == "1.5"
        assert result.fee == "0.001"
        assert result.direction == "OUTGOING"
        assert result.status == "CONFIRMED"
        assert result.confirmations == 6

    def test_uses_fallback_hash_field(self) -> None:
        dto = SimpleNamespace(
            id="tx-2",
            request_id=None,
            wallet_id=None,
            address_id=None,
            currency="ETH",
            blockchain="ETH",
            tx_hash=None,
            hash="0xfallback",
            block_height=None,
            block_number=200,
            block_hash=None,
            amount="0",
            fee="0",
            direction="INCOMING",
            status="PENDING",
            confirmations=0,
            created_at=None,
            creation_date="2024-06-01T10:30:00Z",
            confirmed_at=None,
            confirmation_date=None,
        )
        result = map_transaction(dto)
        assert result.tx_hash == "0xfallback"
        assert result.block_height == 200

    def test_handles_all_none_fields(self) -> None:
        dto = SimpleNamespace(
            id=None,
            request_id=None,
            wallet_id=None,
            address_id=None,
            currency=None,
            blockchain=None,
            tx_hash=None,
            hash=None,
            block_height=None,
            block_number=None,
            block_hash=None,
            amount=None,
            fee=None,
            direction=None,
            status=None,
            confirmations=None,
            created_at=None,
            creation_date=None,
            confirmed_at=None,
            confirmation_date=None,
        )
        result = map_transaction(dto)
        assert result.id == ""
        assert result.confirmations == 0


class TestMapTransactions:
    """Tests for map_transactions function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", request_id=None, wallet_id=None, address_id=None,
                currency="BTC", blockchain="BTC", tx_hash="0xa", hash=None,
                block_height=1, block_number=None, block_hash=None,
                amount="1", fee="0", direction="OUT", status="OK",
                confirmations=1, created_at=None, creation_date=None,
                confirmed_at=None, confirmation_date=None,
            ),
        ]
        result = map_transactions(dtos)
        assert len(result) == 1

    def test_returns_empty_for_none(self) -> None:
        assert map_transactions(None) == []

    def test_returns_empty_for_empty(self) -> None:
        assert map_transactions([]) == []
