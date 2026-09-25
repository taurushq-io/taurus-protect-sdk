"""Unit tests for FeePayerService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import FeePayersApi
from taurus_protect.models.pagination import Pagination
from taurus_protect.services.fee_payer_service import FeePayerService
from tests.unit.transport_stub import StubTransport, api_client


class TestFeePayerServiceList:
    """FeePayerService.list over the real generated client."""

    def _service(self) -> FeePayerService:
        ac = api_client()
        return FeePayerService(ac, FeePayersApi(ac))

    def test_empty_page_still_carries_pagination(self) -> None:
        """An empty page used to return None pagination."""
        with StubTransport({}):
            fee_payers, pagination = self._service().list()

        assert fee_payers == []
        assert pagination == Pagination(limit=20, offset=0)

    def test_filters_and_page_window_reach_the_wire(self) -> None:
        with StubTransport(
            {"result": [{"id": "fp1", "blockchain": "SOL"}], "totalItems": "31"}
        ) as transport:
            fee_payers, pagination = self._service().list(
                limit=10, offset=30, blockchain="SOL", network="mainnet"
            )

        assert [fp.id for fp in fee_payers] == ["fp1"]
        assert pagination == Pagination(
            limit=10, offset=30, total_items=31, next_offset=31, has_more=False
        )
        assert transport.last.query == sorted(
            [("limit", "10"), ("offset", "30"), ("blockchain", "SOL"), ("network", "mainnet")]
        )

    @pytest.mark.parametrize("kwargs,name", [({"limit": 101}, "limit"), ({"offset": -1}, "offset")])
    def test_invalid_page_window_sends_nothing(self, kwargs: dict, name: str) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match=name):
                self._service().list(**kwargs)

        assert transport.requests == []


class TestFeePayerServiceGet:
    """Tests for FeePayerService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        fee_payers_api = MagicMock()
        service = FeePayerService(api_client=api_client, fee_payers_api=fee_payers_api)
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
