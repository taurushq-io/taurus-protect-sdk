"""Unit tests for TaurusNetwork SettlementService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.taurus_network.settlement_service import (
    SettlementService,
)


class TestGetSettlement:
    """Tests for SettlementService.get_settlement()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        settlement_api = MagicMock()
        service = SettlementService(
            api_client=api_client, settlement_api=settlement_api
        )
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
    """Tests for SettlementService.list_settlements()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        settlement_api = MagicMock()
        service = SettlementService(
            api_client=api_client, settlement_api=settlement_api
        )
        return service, settlement_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.taurus_network_service_get_settlements.return_value = resp

        settlements, pagination = service.list_settlements()

        assert settlements == []


class TestCreateSettlement:
    """Tests for SettlementService.create_settlement()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        settlement_api = MagicMock()
        service = SettlementService(
            api_client=api_client, settlement_api=settlement_api
        )
        return service, settlement_api

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.create_settlement(request=None)
