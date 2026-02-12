"""Unit tests for UserService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import APIError, NotFoundError
from taurus_protect.services.user_service import UserService


class TestGet:
    """Tests for UserService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        users_api = MagicMock()
        service = UserService(api_client=api_client, users_api=users_api)
        return service, users_api

    def test_get_returns_user(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.user_service_get_user.return_value = reply

        mock_user = MagicMock()
        with patch(
            "taurus_protect.services.user_service.user_from_dto",
            return_value=mock_user,
        ):
            result = service.get("user-123")

        assert result is mock_user
        api.user_service_get_user.assert_called_once_with("user-123")

    def test_get_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="user_id"):
            service.get("")

    def test_get_raises_not_found_when_result_is_none(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.user_service_get_user.return_value = reply

        with pytest.raises(NotFoundError):
            service.get("user-123")


class TestGetCurrent:
    """Tests for UserService.get_current()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        users_api = MagicMock()
        service = UserService(api_client=api_client, users_api=users_api)
        return service, users_api

    def test_get_current_returns_user(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = MagicMock()
        api.user_service_get_me.return_value = reply

        mock_user = MagicMock()
        with patch(
            "taurus_protect.services.user_service.user_from_dto",
            return_value=mock_user,
        ):
            result = service.get_current()

        assert result is mock_user

    def test_get_current_raises_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.user_service_get_me.return_value = reply

        with pytest.raises(APIError):
            service.get_current()


class TestList:
    """Tests for UserService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        users_api = MagicMock()
        service = UserService(api_client=api_client, users_api=users_api)
        return service, users_api

    def test_list_returns_users_and_pagination(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        reply.total_items = "10"
        api.user_service_get_users.return_value = reply

        with patch(
            "taurus_protect.services.user_service.users_from_dto",
            return_value=[MagicMock()],
        ):
            users, pagination = service.list()

        assert len(users) == 1

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
        api.user_service_get_users.return_value = reply

        users, pagination = service.list()
        assert users == []


class TestGetUsersByEmail:
    """Tests for UserService.get_users_by_email()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        users_api = MagicMock()
        service = UserService(api_client=api_client, users_api=users_api)
        return service, users_api

    def test_get_users_by_email_raises_for_empty_list(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="emails cannot be empty"):
            service.get_users_by_email([])

    def test_get_users_by_email_returns_users(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = [MagicMock()]
        api.user_service_get_users.return_value = reply

        with patch(
            "taurus_protect.services.user_service.users_from_dto",
            return_value=[MagicMock()],
        ):
            result = service.get_users_by_email(["user@example.com"])

        assert len(result) == 1


class TestCreateUserAttribute:
    """Tests for UserService.create_user_attribute()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        users_api = MagicMock()
        service = UserService(api_client=api_client, users_api=users_api)
        return service, users_api

    def test_create_user_attribute_raises_for_empty_user_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="user_id"):
            service.create_user_attribute("", "key", "value")

    def test_create_user_attribute_raises_for_empty_key(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="key"):
            service.create_user_attribute("user-1", "", "value")

    def test_create_user_attribute_calls_api(self) -> None:
        service, api = self._make_service()

        service.create_user_attribute("user-1", "role", "admin")

        api.user_service_create_attribute.assert_called_once()
