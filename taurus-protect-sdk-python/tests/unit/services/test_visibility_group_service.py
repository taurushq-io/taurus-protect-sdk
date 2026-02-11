"""Unit tests for VisibilityGroupService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.visibility_group_service import VisibilityGroupService


class TestVisibilityGroupServiceList:
    """Tests for VisibilityGroupService.list()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        visibility_groups_api = MagicMock()
        service = VisibilityGroupService(
            api_client=api_client, visibility_groups_api=visibility_groups_api
        )
        return service, visibility_groups_api

    def test_raises_on_invalid_limit(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="limit must be positive"):
            service.list(limit=0)

    def test_raises_on_negative_offset(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offset cannot be negative"):
            service.list(offset=-1)

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.user_service_get_visibility_groups.return_value = resp

        groups, pagination = service.list()

        assert groups == []
        assert pagination is not None
        assert pagination.total_items == 0

    def test_applies_client_side_pagination(self) -> None:
        service, api = self._make_service()
        # Create 3 mock group DTOs
        dto1 = MagicMock()
        dto1.id = "1"
        dto1.name = "Group A"
        dto1.description = None
        dto1.users = None
        dto2 = MagicMock()
        dto2.id = "2"
        dto2.name = "Group B"
        dto2.description = None
        dto2.users = None
        dto3 = MagicMock()
        dto3.id = "3"
        dto3.name = "Group C"
        dto3.description = None
        dto3.users = None
        resp = MagicMock()
        resp.result = [dto1, dto2, dto3]
        api.user_service_get_visibility_groups.return_value = resp

        groups, pagination = service.list(limit=2, offset=0)

        assert len(groups) <= 2
        assert pagination.total_items == 3
        assert pagination.has_more is True


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
