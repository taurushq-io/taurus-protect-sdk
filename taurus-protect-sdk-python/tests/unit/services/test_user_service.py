"""Unit tests for UserService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect._internal.openapi import UsersApi
from taurus_protect.errors import APIError, NotFoundError
from taurus_protect.models.pagination import Pagination
from taurus_protect.services.user_service import UserService
from tests.unit.transport_stub import StubTransport, api_client


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
    """UserService.list: next offset = offset + min(rows, limit)."""

    def _service(self) -> UserService:
        ac = api_client()
        return UserService(ac, UsersApi(ac))

    def test_a_synthetic_daemon_user_beyond_the_limit_does_not_skip_a_row(self) -> None:
        rows = [{"id": str(i)} for i in range(3)]
        with StubTransport({"result": rows, "totalItems": "30"}):
            users, pagination = self._service().list(limit=2)

        assert len(users) == 3
        assert pagination == Pagination(
            limit=2, offset=0, total_items=30, next_offset=2, has_more=True
        )

    def test_empty_reply(self) -> None:
        with StubTransport({}) as transport:
            users, pagination = self._service().list()

        assert users == []
        assert pagination == Pagination(limit=20, offset=0)
        assert transport.last.query == [("limit", "20")]

    @pytest.mark.parametrize("kwargs,name", [({"limit": 101}, "limit"), ({"offset": -1}, "offset")])
    def test_invalid_page_window_sends_nothing(self, kwargs: dict, name: str) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match=name):
                self._service().list(**kwargs)

        assert transport.requests == []


class TestGetUsersByEmail:
    """get_users_by_email reads every page, 100 at a time, instead of one unbounded page."""

    def _service(self) -> UserService:
        ac = api_client()
        return UserService(ac, UsersApi(ac))

    def test_get_users_by_email_raises_for_empty_list(self) -> None:
        with pytest.raises(ValueError, match="emails cannot be empty"):
            self._service().get_users_by_email([])

    def test_walks_every_page_and_dedupes(self) -> None:
        page1 = {"result": [{"id": str(i)} for i in range(100)], "totalItems": "101"}
        page2 = {"result": [{"id": "99"}, {"id": "100"}], "totalItems": "101"}
        with StubTransport(page1, page2) as transport:
            users = self._service().get_users_by_email(["a@example.test", "b@example.test"])

        assert [u.id for u in users] == [str(i) for i in range(101)]
        assert transport.requests[0].query == sorted(
            [("emails", "a@example.test"), ("emails", "b@example.test"), ("limit", "100")]
        )
        assert transport.requests[1].param("offset") == "100"


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
