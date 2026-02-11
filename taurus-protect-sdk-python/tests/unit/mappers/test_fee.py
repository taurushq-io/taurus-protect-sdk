"""Unit tests for fee mapper functions."""

from decimal import Decimal
from types import SimpleNamespace
from unittest.mock import MagicMock

import pytest

from taurus_protect.services.fee_service import FeeService


class TestMapFeeEstimate:
    """Tests for FeeService._map_fee_estimate method."""

    def _make_service(self) -> FeeService:
        """Create a FeeService with mocked dependencies."""
        return FeeService(api_client=MagicMock(), fee_api=MagicMock())

    def test_maps_all_fields(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            blockchain="ETH",
            network="mainnet",
            fee_low="0.001",
            fee_medium="0.005",
            fee_high="0.01",
            gas_limit=21000,
            gas_price="50",
        )
        result = service._map_fee_estimate(dto, "ETH", "1.0")
        assert result.currency == "ETH"
        assert result.blockchain == "ETH"
        assert result.network == "mainnet"
        assert result.fee_low == Decimal("0.001")
        assert result.fee_medium == Decimal("0.005")
        assert result.fee_high == Decimal("0.01")
        assert result.gas_limit == 21000
        assert result.gas_price == Decimal("50")
        assert result.amount == Decimal("1.0")

    def test_handles_none_dto(self) -> None:
        service = self._make_service()
        result = service._map_fee_estimate(None, "BTC", None)
        assert result.currency == "BTC"
        assert result.fee_low is None
        assert result.fee_medium is None
        assert result.fee_high is None

    def test_handles_alternative_fee_names(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            blockchain="BTC",
            network="mainnet",
            low="0.0001",
            standard="0.0005",
            fast="0.001",
            gas_limit=None,
            gas_price=None,
        )
        result = service._map_fee_estimate(dto, "BTC", None)
        assert result.fee_low == Decimal("0.0001")
        assert result.fee_medium == Decimal("0.0005")
        assert result.fee_high == Decimal("0.001")

    def test_handles_slow_average_fast_names(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            blockchain="SOL",
            network="mainnet-beta",
            slow="0.00001",
            average="0.00005",
            fast="0.0001",
            gas_limit=None,
            gas_price=None,
        )
        result = service._map_fee_estimate(dto, "SOL", None)
        assert result.fee_low == Decimal("0.00001")
        assert result.fee_medium == Decimal("0.00005")
        assert result.fee_high == Decimal("0.0001")

    def test_handles_none_fee_values(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            blockchain="ADA",
            network="mainnet",
            fee_low=None,
            fee_medium=None,
            fee_high=None,
            gas_limit=None,
            gas_price=None,
        )
        result = service._map_fee_estimate(dto, "ADA", None)
        assert result.fee_low is None
        assert result.fee_medium is None
        assert result.fee_high is None
        assert result.gas_limit is None
        assert result.gas_price is None

    def test_handles_evm_gas_fields(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            blockchain="ETH",
            network="mainnet",
            fee_low="0.001",
            fee_medium="0.002",
            fee_high="0.003",
            gas_limit="21000",
            gasPrice="100",
        )
        result = service._map_fee_estimate(dto, "ETH", None)
        assert result.gas_limit == 21000
        assert result.gas_price == Decimal("100")

    def test_amount_not_set_when_none(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            blockchain="ETH",
            network="mainnet",
            fee_low=None,
            fee_medium=None,
            fee_high=None,
            gas_limit=None,
            gas_price=None,
        )
        result = service._map_fee_estimate(dto, "ETH", None)
        assert result.amount is None

    def test_handles_missing_blockchain_network(self) -> None:
        service = self._make_service()
        dto = SimpleNamespace(
            fee_low="0.001",
            fee_medium="0.002",
            fee_high="0.003",
        )
        result = service._map_fee_estimate(dto, "XRP", None)
        assert result.blockchain == ""
        assert result.network == ""
