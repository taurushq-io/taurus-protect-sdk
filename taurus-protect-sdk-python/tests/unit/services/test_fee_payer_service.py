"""Unit tests for FeePayerService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.fee_payer_service import FeePayerService


class TestFeePayerServiceList:
    """Tests for FeePayerService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fee_payers_api = MagicMock()
        service = FeePayerService(
            api_client=api_client, fee_payers_api=fee_payers_api
        )
        return service, fee_payers_api

    def test_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.fee_payers = None
        api.fee_payer_service_get_fee_payers.return_value = resp

        fee_payers, pagination = service.list()

        assert fee_payers == []
        assert pagination is None

    def test_passes_blockchain_filter(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.fee_payers = None
        api.fee_payer_service_get_fee_payers.return_value = resp

        service.list(blockchain="SOL", network="mainnet")

        api.fee_payer_service_get_fee_payers.assert_called_once_with(
            limit="50",
            offset="0",
            ids=None,
            blockchain="SOL",
            network="mainnet",
        )


class TestFeePayerServiceGet:
    """Tests for FeePayerService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fee_payers_api = MagicMock()
        service = FeePayerService(
            api_client=api_client, fee_payers_api=fee_payers_api
        )
        return service, fee_payers_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="fee_payer_id"):
            service.get(fee_payer_id="")

    def test_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.fee_payer = None
        api.fee_payer_service_get_fee_payer.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get(fee_payer_id="fp-1")
