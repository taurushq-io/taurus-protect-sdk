"""Unit tests for TaurusNetwork PledgeService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.taurus_network.pledge_service import PledgeService


class TestGetPledge:
    """Tests for PledgeService.get_pledge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="pledge_id"):
            service.get_pledge(pledge_id="")

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.taurus_network_service_get_pledge.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get_pledge(pledge_id="p-missing")


class TestListPledges:
    """Tests for PledgeService.list_pledges()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.taurus_network_service_get_pledges.return_value = resp

        pledges, pagination = service.list_pledges()

        assert pledges == []


class TestCreatePledge:
    """Tests for PledgeService.create_pledge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.create_pledge(req=None)

    def test_raises_on_missing_shared_address_id(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import CreatePledgeRequest

        req = CreatePledgeRequest(
            shared_address_id="",
            currency_id="ETH",
            amount="1000",
        )
        with pytest.raises(ValueError, match="shared_address_id"):
            service.create_pledge(req=req)


class TestApprovePledgeActions:
    """Tests for PledgeService.approve_pledge_actions()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_empty_actions(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="actions list cannot be empty"):
            service.approve_pledge_actions(actions=[], private_key=MagicMock())

    def test_raises_on_none_private_key(self) -> None:
        service, _ = self._make_service()
        action = MagicMock()
        action.id = "a-1"
        action.metadata = MagicMock()
        action.metadata.hash = "abc123"
        with pytest.raises(ValueError, match="private_key cannot be None"):
            service.approve_pledge_actions(actions=[action], private_key=None)

    def test_raises_on_missing_metadata(self) -> None:
        service, _ = self._make_service()
        action = MagicMock()
        action.id = "a-1"
        action.metadata = None
        with pytest.raises(ValueError, match="action metadata cannot be None"):
            service.approve_pledge_actions(
                actions=[action], private_key=MagicMock()
            )

    def test_raises_on_empty_hash(self) -> None:
        service, _ = self._make_service()
        action = MagicMock()
        action.id = "a-1"
        action.metadata = MagicMock()
        action.metadata.hash = ""
        with pytest.raises(ValueError, match="action metadata hash cannot be empty"):
            service.approve_pledge_actions(
                actions=[action], private_key=MagicMock()
            )


class TestRejectPledgeActions:
    """Tests for PledgeService.reject_pledge_actions()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.reject_pledge_actions(req=None)

    def test_raises_on_empty_ids(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import (
            RejectPledgeActionsRequest,
        )

        req = RejectPledgeActionsRequest(ids=[], comment="rejected")
        with pytest.raises(ValueError, match="ids list cannot be empty"):
            service.reject_pledge_actions(req=req)

    def test_raises_on_empty_comment(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import (
            RejectPledgeActionsRequest,
        )

        req = RejectPledgeActionsRequest(ids=["a-1"], comment="")
        with pytest.raises(ValueError, match="comment"):
            service.reject_pledge_actions(req=req)


class TestUnpledge:
    """Tests for PledgeService.unpledge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="pledge_id"):
            service.unpledge(pledge_id="")


class TestRejectPledge:
    """Tests for PledgeService.reject_pledge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import RejectPledgeRequest

        req = RejectPledgeRequest(comment="not approved")
        with pytest.raises(ValueError, match="pledge_id"):
            service.reject_pledge(pledge_id="", req=req)

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.reject_pledge(pledge_id="p-1", req=None)

    def test_raises_on_empty_comment(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import RejectPledgeRequest

        req = RejectPledgeRequest(comment="")
        with pytest.raises(ValueError, match="comment"):
            service.reject_pledge(pledge_id="p-1", req=req)


class TestListPledgeWithdrawals:
    """Tests for PledgeService.list_pledge_withdrawals()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.taurus_network_service_get_pledges_withdrawals.return_value = resp

        withdrawals, pagination = service.list_pledge_withdrawals()

        assert withdrawals == []
