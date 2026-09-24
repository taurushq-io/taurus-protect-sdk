"""Unit tests for GroupService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect._internal.openapi import GroupsApi
from taurus_protect.errors import NotFoundError
from taurus_protect.models.pagination import Pagination
from taurus_protect.services.group_service import GroupService
from tests.unit.transport_stub import StubTransport, api_client


class TestGet:
    """Tests for GroupService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        groups_api = MagicMock()
        service = GroupService(api_client=api_client, groups_api=groups_api)
        return service, groups_api

    def test_get_returns_group(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        api.user_service_get_groups.return_value = reply

        mock_group = MagicMock()
        with patch(
            "taurus_protect.services.group_service.group_from_dto",
            return_value=mock_group,
        ):
            result = service.get("group-1")

        assert result is mock_group

    def test_get_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="group_id"):
            service.get("")

    def test_get_raises_not_found_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = []
        api.user_service_get_groups.return_value = reply

        with pytest.raises(NotFoundError, match="not found"):
            service.get("group-999")

    def test_get_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.user_service_get_groups.return_value = reply

        with pytest.raises(NotFoundError, match="not found"):
            service.get("group-999")


class TestList:
    """GroupService.list: next offset = offset + min(rows, limit)."""

    def _service(self) -> GroupService:
        ac = api_client()
        return GroupService(ac, GroupsApi(ac))

    def test_a_synthetic_group_beyond_the_limit_does_not_skip_a_row(self) -> None:
        rows = [{"id": str(i)} for i in range(3)]
        with StubTransport({"result": rows, "totalItems": "5"}) as transport:
            groups, pagination = self._service().list(limit=2, offset=2)

        assert len(groups) == 3
        assert pagination == Pagination(
            limit=2, offset=2, total_items=5, next_offset=4, has_more=True
        )
        assert transport.last.query == [("limit", "2"), ("offset", "2")]

    def test_empty_reply(self) -> None:
        with StubTransport({}):
            groups, pagination = self._service().list()

        assert groups == []
        assert pagination == Pagination(limit=20, offset=0)

    @pytest.mark.parametrize("kwargs,name", [({"limit": 101}, "limit"), ({"offset": -1}, "offset")])
    def test_invalid_page_window_sends_nothing(self, kwargs: dict, name: str) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match=name):
                self._service().list(**kwargs)

        assert transport.requests == []
