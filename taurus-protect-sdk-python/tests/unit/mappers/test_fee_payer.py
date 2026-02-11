"""Unit tests for fee payer mapper functions."""

from decimal import Decimal
from types import SimpleNamespace
from unittest.mock import MagicMock

import pytest

from taurus_protect.services.fee_payer_service import FeePayerService


class TestMapFeePayer:
    """Tests for FeePayerService._map_fee_payer method."""

    def _make_service(self) -> FeePayerService:
        """Create a FeePayerService with mocked dependencies."""
        return FeePayerService(api_client=MagicMock(), fee_payers_api=MagicMock())

    def test_maps_all_fields(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            id="fp-1",
            blockchain="SOL",
            network="mainnet-beta",
            address="8ZeUv7aXJp4YGKhVrPPjnPXNzRxh8BAiTx7qYqFBN3Vz",
            balance="100.5",
            status="active",
            created_at="2024-01-15T10:30:00Z",
            updated_at="2024-06-01T12:00:00Z",
        )
        result = service._map_fee_payer(dto)
        assert result is not None
        assert result.id == "fp-1"
        assert result.blockchain == "SOL"
        assert result.network == "mainnet-beta"
        assert result.address == "8ZeUv7aXJp4YGKhVrPPjnPXNzRxh8BAiTx7qYqFBN3Vz"
        assert result.balance == Decimal("100.5")
        assert result.status == "active"

    def test_returns_none_for_none(self) -> None:
        service = self._make_service()
        assert service._map_fee_payer(None) is None

    def test_returns_none_when_no_id(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            blockchain="ETH",
            network="mainnet",
            address="0xabc",
            balance="1.0",
            status="active",
            created_at=None,
            updated_at=None,
        )
        result = service._map_fee_payer(dto)
        assert result is None

    def test_returns_none_when_id_is_empty(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            id="",
            blockchain="ETH",
            network="mainnet",
            address="0xabc",
            balance="1.0",
            status="active",
            created_at=None,
            updated_at=None,
        )
        result = service._map_fee_payer(dto)
        assert result is None

    def test_handles_none_optional_fields(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            id="fp-2",
            blockchain=None,
            network=None,
            address=None,
            balance=None,
            status=None,
            created_at=None,
            updated_at=None,
        )
        result = service._map_fee_payer(dto)
        assert result is not None
        assert result.id == "fp-2"
        assert result.blockchain == ""
        assert result.network == ""
        assert result.address == ""
        assert result.balance is None
        assert result.status == "active"

    def test_uses_public_key_as_address_fallback(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            id="fp-3",
            blockchain="SOL",
            network="devnet",
            public_key="PublicKeyABC",
            balance="50.0",
            status="active",
            created_at=None,
            updated_at=None,
        )
        result = service._map_fee_payer(dto)
        assert result is not None
        assert result.address == "PublicKeyABC"

    def test_uses_pubkey_as_address_fallback(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            id="fp-4",
            blockchain="SOL",
            network="devnet",
            pubkey="PubkeyDEF",
            balance=None,
            status="active",
            created_at=None,
            updated_at=None,
        )
        result = service._map_fee_payer(dto)
        assert result is not None
        assert result.address == "PubkeyDEF"

    def test_handles_boolean_true_status(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            id="fp-5",
            blockchain="ETH",
            network="mainnet",
            address="0x123",
            balance="10.0",
            status=True,
            created_at=None,
            updated_at=None,
        )
        result = service._map_fee_payer(dto)
        assert result is not None
        assert result.status == "active"

    def test_handles_boolean_false_status(self) -> None:
        # Note: `False or "active"` short-circuits to "active" in the mapper,
        # so boolean False status also resolves to "active"
        service = self._make_service()
        dto = SimpleNamespace(
            id="fp-6",
            blockchain="ETH",
            network="mainnet",
            address="0x456",
            balance="0",
            status=False,
            created_at=None,
            updated_at=None,
        )
        result = service._map_fee_payer(dto)
        assert result is not None
        assert result.status == "active"

    def test_handles_integer_balance(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            id="fp-7",
            blockchain="ETH",
            network="mainnet",
            address="0xdef",
            balance=42,
            status="active",
            created_at=None,
            updated_at=None,
        )
        result = service._map_fee_payer(dto)
        assert result is not None
        assert result.balance == Decimal("42")

    def test_handles_zero_balance(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            id="fp-8",
            blockchain="ETH",
            network="mainnet",
            address="0xghi",
            balance="0",
            status="active",
            created_at=None,
            updated_at=None,
        )
        result = service._map_fee_payer(dto)
        assert result is not None
        assert result.balance == Decimal("0")
