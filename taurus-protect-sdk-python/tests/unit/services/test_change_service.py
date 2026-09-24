"""Unit tests for ChangeService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect._internal.openapi import ChangesApi
from taurus_protect.errors import NotFoundError
from taurus_protect.models.audit import (
    ChangeResult,
    CreateChangeRequest,
    ListChangesOptions,
)
from taurus_protect.models.pagination import CursorPage
from taurus_protect.services.change_service import ChangeService
from tests.unit.transport_stub import StubTransport, api_client


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
    """ChangeService.list over the real generated client."""

    def _service(self) -> ChangeService:
        ac = api_client()
        return ChangeService(ac, ChangesApi(ac))

    def test_rows_page_and_filters(self) -> None:
        reply = {"result": [{"id": "c1"}], "cursor": {"currentPage": "p2", "hasNext": True}}
        options = ListChangesOptions(
            entity="businessrule",
            entity_id="7",
            status="Created",
            creator_id="u1",
            sort_order="ASC",
            entity_ids=["1"],
            entity_uuids=["u"],
            page_size=25,
        )
        with StubTransport(reply) as transport:
            result = self._service().list(options)

        assert isinstance(result, ChangeResult)
        assert [c.id for c in result.changes] == ["c1"]
        assert result.page == CursorPage(page_size=25, next_cursor="p2", has_more=True)
        assert dict(transport.last.query) == {
            "entity": "businessrule",
            "entityId": "7",
            "status": "Created",
            "creatorId": "u1",
            "sortOrder": "ASC",
            "entityIDs": "1",
            "entityUUIDs": "u",
            "cursor.pageSize": "25",
        }

    def test_continuation_sends_next(self) -> None:
        with StubTransport() as transport:
            self._service().list(ListChangesOptions(cursor="p2"))

        assert transport.last.param("cursor.currentPage") == "p2"
        assert transport.last.param("cursor.pageRequest") == "NEXT"
        assert transport.last.param("cursor.pageSize") == "20"

    def test_empty_reply(self) -> None:
        with StubTransport({}):
            result = self._service().list()

        assert result.changes == []
        assert result.page == CursorPage(page_size=20)


class TestListForApproval:
    """ChangeService.list_for_approval: the queue's own filters, the rest refused by name."""

    def _service(self) -> ChangeService:
        ac = api_client()
        return ChangeService(ac, ChangesApi(ac))

    def test_entity_filters_the_queue(self) -> None:
        with StubTransport() as transport:
            self._service().list_for_approval(ListChangesOptions(entity="user", sort_order="DESC"))

        assert transport.last.path == "/api/rest/v1/changes/for-approval"
        assert transport.last.query == sorted(
            [("entities", "user"), ("sortOrder", "DESC"), ("cursor.pageSize", "20")]
        )

    @pytest.mark.parametrize("field", ["status", "creator_id", "entity_id"])
    def test_options_the_queue_cannot_apply_are_refused(self, field: str) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match=field):
                self._service().list_for_approval(ListChangesOptions(**{field: "x"}))

        assert transport.requests == []


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
