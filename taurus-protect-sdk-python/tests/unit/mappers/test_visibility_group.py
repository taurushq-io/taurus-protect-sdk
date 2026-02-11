"""Unit tests for visibility group mapper."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.visibility_group import (
    visibility_group_from_dto,
    visibility_group_user_from_dto,
    visibility_groups_from_dto,
)


class TestVisibilityGroupFromDto:
    """Tests for visibility_group_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="vg-1",
            tenant_id="t-1",
            name="Treasury VG",
            description="Treasury team visibility",
            user_count="5",
            users=[
                SimpleNamespace(id="u-1", email="user1@example.com", name="User 1"),
                SimpleNamespace(id="u-2", email="user2@example.com", name="User 2"),
            ],
            creation_date="2024-01-01T00:00:00Z",
            update_date="2024-06-01T00:00:00Z",
        )
        result = visibility_group_from_dto(dto)
        assert result is not None
        assert result.id == "vg-1"
        assert result.tenant_id == "t-1"
        assert result.name == "Treasury VG"
        assert result.description == "Treasury team visibility"
        assert result.user_count == 5
        assert len(result.users) == 2

    def test_returns_none_for_none(self) -> None:
        assert visibility_group_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="vg-2",
            tenant_id=None,
            name="Empty VG",
            description=None,
            user_count=None,
            users=None,
            creation_date=None,
            update_date=None,
        )
        result = visibility_group_from_dto(dto)
        assert result is not None
        assert result.user_count == 0
        assert result.users == []


class TestVisibilityGroupsFromDto:
    """Tests for visibility_groups_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", tenant_id=None, name="VG1", description=None,
                user_count=None, users=None, creation_date=None, update_date=None,
            ),
        ]
        result = visibility_groups_from_dto(dtos)
        assert len(result) == 1

    def test_returns_empty_for_none(self) -> None:
        assert visibility_groups_from_dto(None) == []


class TestVisibilityGroupUserFromDto:
    """Tests for visibility_group_user_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(id="u-1", email="user@example.com", name="Test User")
        result = visibility_group_user_from_dto(dto)
        assert result is not None
        assert result.id == "u-1"
        assert result.email == "user@example.com"
        assert result.name == "Test User"

    def test_returns_none_for_none(self) -> None:
        assert visibility_group_user_from_dto(None) is None
