"""Unit tests for GroupService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import NotFoundError
from taurus_protect.services.group_service import GroupService


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
    """Tests for GroupService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        groups_api = MagicMock()
        service = GroupService(api_client=api_client, groups_api=groups_api)
        return service, groups_api

    def test_list_returns_groups_and_pagination(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "5"
        api.user_service_get_groups.return_value = reply

        with patch(
            "taurus_protect.services.group_service.groups_from_dto",
            return_value=[MagicMock()],
        ):
            groups, pagination = service.list()

        assert len(groups) == 1

    def test_list_raises_for_invalid_limit(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_list_raises_for_negative_offset(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_list_returns_empty_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        reply.total_items = None
        api.user_service_get_groups.return_value = reply

        groups, pagination = service.list()
        assert groups == []
