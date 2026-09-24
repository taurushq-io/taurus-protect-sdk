"""Unit tests for TaurusNetwork SettlementService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import TaurusNetworkSettlementApi
from taurus_protect.models.pagination import CursorPage
from taurus_protect.models.taurus_network.settlement import (
    ListSettlementsForApprovalOptions,
    ListSettlementsOptions,
)
from taurus_protect.services.taurus_network.settlement_service import (
    SettlementService,
)
from tests.unit.transport_stub import StubTransport, api_client


class TestGetSettlement:
    """Tests for SettlementService.get_settlement()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        settlement_api = MagicMock()
        service = SettlementService(api_client=api_client, settlement_api=settlement_api)
        return service, settlement_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="settlement_id"):
            service.get_settlement(settlement_id="")

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.taurus_network_service_get_settlement.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get_settlement(settlement_id="s-missing")


class TestListSettlements:
    """Settlement lists over the real generated client."""

    def _service(self) -> SettlementService:
        ac = api_client()
        return SettlementService(ac, TaurusNetworkSettlementApi(ac))

    def test_rows_page_and_filters(self) -> None:
        reply = {"result": [{"id": "s1"}], "cursor": {"currentPage": "n", "hasNext": True}}
        with StubTransport(reply) as transport:
            settlements, page = self._service().list_settlements(
                ListSettlementsOptions(
                    counter_participant_id="p", statuses=["A"], sort_order="DESC"
                )
            )

        assert [s.id for s in settlements] == ["s1"]
        assert page == CursorPage(page_size=20, next_cursor="n", has_more=True)
        assert transport.last.query == sorted(
            [
                ("counterParticipantID", "p"),
                ("statuses", "A"),
                ("sortOrder", "DESC"),
                ("cursor.pageSize", "20"),
            ]
        )

    def test_for_approval_continuation(self) -> None:
        with StubTransport() as transport:
            self._service().list_settlements_for_approval(
                ListSettlementsForApprovalOptions(ids=["1"], cursor="n", page_size=3)
            )

        assert transport.last.path == "/api/rest/v1/tn/settlements/for-approval"
        assert transport.last.query == sorted(
            [
                ("ids", "1"),
                ("cursor.currentPage", "n"),
                ("cursor.pageRequest", "NEXT"),
                ("cursor.pageSize", "3"),
            ]
        )


class TestCreateSettlement:
    """Tests for SettlementService.create_settlement()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        settlement_api = MagicMock()
        service = SettlementService(api_client=api_client, settlement_api=settlement_api)
        return service, settlement_api

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.create_settlement(request=None)
