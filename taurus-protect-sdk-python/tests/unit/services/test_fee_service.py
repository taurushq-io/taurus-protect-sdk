"""Unit tests for FeeService."""

from __future__ import annotations

from decimal import Decimal
from unittest.mock import MagicMock

import pytest

from taurus_protect.services.fee_service import FeeService


class TestFeeServiceEstimate:
    """Tests for FeeService.estimate()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fee_api = MagicMock()
        service = FeeService(api_client=api_client, fee_api=fee_api)
        return service, fee_api

    def test_estimate_raises_on_empty_currency(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="currency"):
            service.estimate(currency="")

    def test_estimate_returns_empty_when_no_data(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.fees = None
        api.fee_service_get_fees_v2.return_value = resp

        result = service.estimate(currency="ETH")

        assert result.currency == "ETH"
        assert result.fee_low is None

    def test_estimate_finds_currency_in_list(self) -> None:
        service, api = self._make_service()
        fee_dto = MagicMock()
        fee_dto.currency = "ETH"
        fee_dto.symbol = None
        fee_dto.blockchain = "Ethereum"
        fee_dto.network = "mainnet"
        fee_dto.fee_low = "0.001"
        fee_dto.low = None
        fee_dto.slow = None
        fee_dto.fee_medium = "0.005"
        fee_dto.medium = None
        fee_dto.standard = None
        fee_dto.average = None
        fee_dto.fee_high = "0.01"
        fee_dto.high = None
        fee_dto.fast = None
        fee_dto.gas_limit = "21000"
        fee_dto.gas_price = "50"
        fee_dto.gasPrice = None
        resp = MagicMock()
        resp.result = [fee_dto]
        resp.fees = None
        api.fee_service_get_fees_v2.return_value = resp

        result = service.estimate(currency="ETH")

        assert result.currency == "ETH"
        assert result.fee_low == Decimal("0.001")
        assert result.fee_medium == Decimal("0.005")
        assert result.fee_high == Decimal("0.01")

    def test_estimate_returns_empty_when_currency_not_found(self) -> None:
        service, api = self._make_service()
        fee_dto = MagicMock()
        fee_dto.currency = "BTC"
        fee_dto.symbol = None
        resp = MagicMock()
        resp.result = [fee_dto]
        resp.fees = None
        api.fee_service_get_fees_v2.return_value = resp

        result = service.estimate(currency="ETH")

        assert result.currency == "ETH"
        assert result.fee_low is None


class TestFeeServiceList:
    """Tests for FeeService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fee_api = MagicMock()
        service = FeeService(api_client=api_client, fee_api=fee_api)
        return service, fee_api

    def test_list_returns_empty_when_no_data(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.fees = None
        api.fee_service_get_fees_v2.return_value = resp

        result = service.list()

        assert result == []

    def test_list_wraps_api_error(self) -> None:
        service, api = self._make_service()
        error = Exception("connection refused")
        error.status = 503
        error.body = None
        error.headers = {}
        api.fee_service_get_fees_v2.side_effect = error

        from taurus_protect.errors import APIError

        with pytest.raises(APIError):
            service.list()
