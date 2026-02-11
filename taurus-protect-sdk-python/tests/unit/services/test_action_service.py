"""Unit tests for ActionService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import NotFoundError
from taurus_protect.services.action_service import ActionService


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
    """Tests for ActionService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        actions_api = MagicMock()
        service = ActionService(api_client=api_client, actions_api=actions_api)
        return service, actions_api

    def test_list_returns_actions_and_pagination(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "10"
        api.action_service_get_actions.return_value = reply

        with patch(
            "taurus_protect.services.action_service.actions_from_dto",
            return_value=[MagicMock()],
        ):
            actions, pagination = service.list()

        assert len(actions) == 1

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
        api.action_service_get_actions.return_value = reply

        actions, pagination = service.list()
        assert actions == []
