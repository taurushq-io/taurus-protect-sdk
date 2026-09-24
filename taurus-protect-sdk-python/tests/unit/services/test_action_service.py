"""Unit tests for ActionService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect._internal.openapi import ActionsApi
from taurus_protect.errors import NotFoundError
from taurus_protect.models.pagination import Pagination
from taurus_protect.services.action_service import ActionService
from tests.unit.transport_stub import StubTransport, api_client


class TestGet:
    """Tests for ActionService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        actions_api = MagicMock()
        service = ActionService(api_client=api_client, actions_api=actions_api)
        return service, actions_api

    def test_get_returns_action(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.action = MagicMock()
        api.action_service_get_action.return_value = reply

        mock_action = MagicMock()
        with patch(
            "taurus_protect.services.action_service.action_from_dto",
            return_value=mock_action,
        ):
            result = service.get("action-1")

        assert result is mock_action
        api.action_service_get_action.assert_called_once_with("action-1")

    def test_get_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="action_id"):
            service.get("")

    def test_get_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.action = None
        api.action_service_get_action.return_value = reply

        with pytest.raises(NotFoundError, match="not found"):
            service.get("action-999")


class TestList:
    """ActionService.list over the real generated client (next offset = offset + rows)."""

    def _service(self) -> ActionService:
        ac = api_client()
        return ActionService(ac, ActionsApi(ac))

    def test_rows_and_pagination(self) -> None:
        with StubTransport(
            {"result": [{"id": "a1"}, {"id": "a2"}], "totalItems": "3"}
        ) as transport:
            actions, pagination = self._service().list(limit=2)

        assert [a.id for a in actions] == ["a1", "a2"]
        assert pagination == Pagination(
            limit=2, offset=0, total_items=3, next_offset=2, has_more=True
        )
        assert transport.last.query == [("limit", "2")]

    def test_empty_reply(self) -> None:
        with StubTransport({}):
            actions, pagination = self._service().list()

        assert actions == []
        assert pagination == Pagination(limit=20, offset=0)

    @pytest.mark.parametrize("kwargs,name", [({"limit": 101}, "limit"), ({"offset": -1}, "offset")])
    def test_invalid_page_window_sends_nothing(self, kwargs: dict, name: str) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match=name):
                self._service().list(**kwargs)

        assert transport.requests == []
