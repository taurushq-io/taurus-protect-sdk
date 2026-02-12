"""Unit tests for group mapper functions."""

from types import SimpleNamespace

from taurus_protect.mappers.user import (
    group_from_dto,
    group_user_from_dto,
    groups_from_dto,
)


class TestGroupFromDto:
    """Tests for group_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="g-1",
            external_group_id="ext-g-1",
            tenant_id="t-1",
            name="Admins",
            email="admins@example.com",
            description="Admin group",
            users=[
                SimpleNamespace(id="u-1", email="alice@example.com"),
                SimpleNamespace(id="u-2", email="bob@example.com"),
            ],
            enforced_in_rules=True,
            creation_date="2024-01-01T00:00:00Z",
            update_date="2024-06-01T00:00:00Z",
        )
        result = group_from_dto(dto)
        assert result is not None
        assert result.id == "g-1"
        assert result.external_group_id == "ext-g-1"
        assert result.name == "Admins"
        assert result.email == "admins@example.com"
        assert result.description == "Admin group"
        assert len(result.users) == 2
        assert result.users[0].id == "u-1"
        assert result.enforced_in_rules is True

    def test_returns_none_for_none(self) -> None:
        assert group_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="g-2",
            external_group_id=None,
            tenant_id=None,
            name="Empty Group",
            email=None,
            description=None,
            users=None,
            enforced_in_rules=None,
            creation_date=None,
            update_date=None,
        )
        result = group_from_dto(dto)
        assert result is not None
        assert result.id == "g-2"
        assert result.users == []
        assert result.enforced_in_rules is False


class TestGroupsFromDto:
    """Tests for groups_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="g-1", external_group_id=None, tenant_id=None,
                name="G1", email=None, description=None, users=None,
                enforced_in_rules=None, creation_date=None, update_date=None,
            ),
            SimpleNamespace(
                id="g-2", external_group_id=None, tenant_id=None,
                name="G2", email=None, description=None, users=None,
                enforced_in_rules=None, creation_date=None, update_date=None,
            ),
        ]
        result = groups_from_dto(dtos)
        assert len(result) == 2

    def test_returns_empty_for_none(self) -> None:
        assert groups_from_dto(None) == []


class TestGroupUserFromDto:
    """Tests for group_user_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(id="u-1", email="test@example.com")
        result = group_user_from_dto(dto)
        assert result.id == "u-1"
        assert result.email == "test@example.com"

    def test_handles_none_email(self) -> None:
        dto = SimpleNamespace(id="u-2", email=None)
        result = group_user_from_dto(dto)
        assert result.id == "u-2"
        assert result.email is None
