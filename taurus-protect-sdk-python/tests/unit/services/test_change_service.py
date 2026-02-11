"""Unit tests for ChangeService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import NotFoundError
from taurus_protect.models.audit import (
    ChangeResult,
    CreateChangeRequest,
    ListChangesOptions,
)
from taurus_protect.services.change_service import ChangeService


class TestCreateChange:
    """Tests for ChangeService.create_change()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        changes_api = MagicMock()
        service = ChangeService(api_client=api_client, changes_api=changes_api)
        return service, changes_api

    def test_create_change_returns_id(self) -> None:
        service, api = self._make_service()

        result_obj = MagicMock()
        result_obj.id = "change-42"
        reply = MagicMock()
        reply.result = result_obj
        api.change_service_create_change.return_value = reply

        request = CreateChangeRequest(
            action="update",
            entity="businessrule",
            entity_id="10",
            changes={"rulevalue": "100"},
            comment="Test",
        )
        change_id = service.create_change(request)
        assert change_id == "change-42"
        api.change_service_create_change.assert_called_once()

    def test_create_change_raises_for_empty_action(self) -> None:
        service, _ = self._make_service()
        request = CreateChangeRequest(action="", entity="user")
        with pytest.raises(ValueError, match="action cannot be empty"):
            service.create_change(request)

    def test_create_change_raises_for_empty_entity(self) -> None:
        service, _ = self._make_service()
        request = CreateChangeRequest(action="update", entity="")
        with pytest.raises(ValueError, match="entity cannot be empty"):
            service.create_change(request)

    def test_create_change_raises_for_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.create_change(None)


class TestList:
    """Tests for ChangeService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        changes_api = MagicMock()
        service = ChangeService(api_client=api_client, changes_api=changes_api)
        return service, changes_api

    def test_list_returns_change_result(self) -> None:
        from taurus_protect.models.audit import Change

        service, api = self._make_service()

        cursor_obj = MagicMock()
        cursor_obj.current_page = "page1"
        cursor_obj.has_next = True
        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.cursor = cursor_obj
        api.change_service_get_changes.return_value = reply

        mock_change = Change(id="1")
        with patch(
            "taurus_protect.services.change_service.changes_from_dto",
            return_value=[mock_change],
        ):
            result = service.list()

        assert isinstance(result, ChangeResult)
        assert len(result.changes) == 1
        assert result.has_next is True
        assert result.current_page == "page1"

    def test_list_passes_options(self) -> None:
        service, api = self._make_service()

        cursor_obj = MagicMock(spec=["current_page", "has_next"])
        cursor_obj.current_page = None
        cursor_obj.has_next = False
        reply = MagicMock()
        reply.result = None
        reply.cursor = cursor_obj
        api.change_service_get_changes.return_value = reply

        opts = ListChangesOptions(entity="businessrule", status="pending", page_size=25)
        result = service.list(options=opts)

        api.change_service_get_changes.assert_called_once_with(
            entity="businessrule",
            entity_id=None,
            status="pending",
            creator_id=None,
            sort_order=None,
            cursor_current_page=None,
            cursor_page_request="FIRST",
            cursor_page_size="25",
            entity_ids=None,
            entity_uuids=None,
        )
        assert result.changes == []

    def test_list_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        reply.cursor = None
        api.change_service_get_changes.return_value = reply

        result = service.list()
        assert result.changes == []
        assert result.has_next is False


class TestListForApproval:
    """Tests for ChangeService.list_for_approval()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        changes_api = MagicMock()
        service = ChangeService(api_client=api_client, changes_api=changes_api)
        return service, changes_api

    def test_list_for_approval_calls_api(self) -> None:
        service, api = self._make_service()

        cursor_obj = MagicMock()
        cursor_obj.current_page = "abc"
        cursor_obj.has_next = False
        reply = MagicMock()
        reply.result = []
        reply.cursor = cursor_obj
        api.change_service_get_changes_for_approval.return_value = reply

        result = service.list_for_approval()

        api.change_service_get_changes_for_approval.assert_called_once()
        assert isinstance(result, ChangeResult)


class TestGet:
    """Tests for ChangeService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        changes_api = MagicMock()
        service = ChangeService(api_client=api_client, changes_api=changes_api)
        return service, changes_api

    def test_get_returns_change(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.change_service_get_change.return_value = reply

        mock_change = MagicMock()
        with patch(
            "taurus_protect.services.change_service.change_from_dto",
            return_value=mock_change,
        ):
            result = service.get("123")

        assert result is mock_change
        api.change_service_get_change.assert_called_once_with("123")

    def test_get_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="change_id"):
            service.get("")

    def test_get_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.change_service_get_change.return_value = reply

        with pytest.raises(NotFoundError, match="not found"):
            service.get("123")


class TestApproveChange:
    """Tests for ChangeService.approve_change()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        changes_api = MagicMock()
        service = ChangeService(api_client=api_client, changes_api=changes_api)
        return service, changes_api

    def test_approve_change_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="change_id"):
            service.approve_change("")

    def test_approve_change_calls_api(self) -> None:
        service, api = self._make_service()

        service.approve_change("123")

        api.change_service_approve_change.assert_called_once_with("123", body={})


class TestApproveChanges:
    """Tests for ChangeService.approve_changes()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        changes_api = MagicMock()
        service = ChangeService(api_client=api_client, changes_api=changes_api)
        return service, changes_api

    def test_approve_changes_raises_for_empty_list(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="change_ids cannot be empty"):
            service.approve_changes([])

    def test_approve_changes_calls_api(self) -> None:
        service, api = self._make_service()

        service.approve_changes(["1", "2", "3"])

        api.change_service_approve_changes.assert_called_once()


class TestRejectChange:
    """Tests for ChangeService.reject_change()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        changes_api = MagicMock()
        service = ChangeService(api_client=api_client, changes_api=changes_api)
        return service, changes_api

    def test_reject_change_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="change_id"):
            service.reject_change("")

    def test_reject_change_calls_api(self) -> None:
        service, api = self._make_service()

        service.reject_change("123")

        api.change_service_reject_change.assert_called_once_with("123", body={})


class TestRejectChanges:
    """Tests for ChangeService.reject_changes()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        changes_api = MagicMock()
        service = ChangeService(api_client=api_client, changes_api=changes_api)
        return service, changes_api

    def test_reject_changes_raises_for_empty_list(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="change_ids cannot be empty"):
            service.reject_changes([])

    def test_reject_changes_calls_api(self) -> None:
        service, api = self._make_service()

        service.reject_changes(["1", "2"])

        api.change_service_reject_changes.assert_called_once()
