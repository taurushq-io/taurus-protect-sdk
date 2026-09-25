"""Unit tests for VisibilityGroupService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import RestrictedVisibilityGroupsApi
from taurus_protect.services.visibility_group_service import VisibilityGroupService
from tests.unit.transport_stub import StubTransport, api_client


class TestVisibilityGroupServiceList:
    """VisibilityGroupService.list: the endpoint does not page, so every group comes back."""

    def _service(self) -> VisibilityGroupService:
        ac = api_client()
        return VisibilityGroupService(ac, RestrictedVisibilityGroupsApi(ac))

    def test_returns_every_group(self) -> None:
        reply = {"result": [{"id": str(i), "name": f"g{i}"} for i in range(25)]}
        with StubTransport(reply) as transport:
            groups = self._service().list()

        assert [g.id for g in groups] == [str(i) for i in range(25)]
        assert transport.last.query == []

    def test_empty_reply(self) -> None:
        with StubTransport({}):
            assert self._service().list() == []

    def test_get_finds_a_group_past_the_first_twenty(self) -> None:
        """get scans list(); a list cut at 20 reported the 21st group as not found."""
        groups = {"result": [{"id": str(i), "name": f"g{i}"} for i in range(25)]}
        with StubTransport(groups, {"result": []}):
            group = self._service().get("21")

        assert group.name == "g21"


class TestVisibilityGroupServiceGet:
    """Tests for VisibilityGroupService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        visibility_groups_api = MagicMock()
        service = VisibilityGroupService(
            api_client=api_client, visibility_groups_api=visibility_groups_api
        )
        return service, visibility_groups_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="group_id"):
            service.get(group_id="")

    def test_raises_not_found_when_missing(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.user_service_get_visibility_groups.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get(group_id="missing-id")


class TestVisibilityGroupServiceGetUsers:
    """Tests for VisibilityGroupService.get_users()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        visibility_groups_api = MagicMock()
        service = VisibilityGroupService(
            api_client=api_client, visibility_groups_api=visibility_groups_api
        )
        return service, visibility_groups_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="group_id"):
            service.get_users(group_id="")

    def test_returns_empty_when_no_users(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.user_service_get_users_by_visibility_group_id.return_value = resp

        result = service.get_users(group_id="g-1")

        assert result == []
