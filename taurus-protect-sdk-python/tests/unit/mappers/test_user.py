"""Unit tests for user, group, and tag mappers."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.user import (
    group_from_dto,
    group_user_from_dto,
    groups_from_dto,
    tag_from_dto,
    tags_from_dto,
    user_attribute_from_dto,
    user_from_dto,
    user_group_from_dto,
    users_from_dto,
)


class TestUserFromDto:
    """Tests for user_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="user-1",
            external_user_id="ext-1",
            tenant_id="t-1",
            username="jdoe",
            email="jdoe@example.com",
            first_name="John",
            last_name="Doe",
            status="ACTIVE",
            roles=["ADMIN", "OPERATOR"],
            public_key="-----BEGIN PUBLIC KEY-----",
            groups=[
                SimpleNamespace(id="g-1", name="Admins"),
            ],
            totp_enabled=True,
            password_changed=True,
            enforced_in_rules=False,
            creation_date="2024-01-01T00:00:00Z",
            update_date="2024-06-01T00:00:00Z",
            last_login="2024-06-15T00:00:00Z",
            attributes=[
                SimpleNamespace(id="a1", key="dept", value="IT"),
            ],
        )
        result = user_from_dto(dto)
        assert result is not None
        assert result.id == "user-1"
        assert result.email == "jdoe@example.com"
        assert result.first_name == "John"
        assert result.roles == ["ADMIN", "OPERATOR"]
        assert result.totp_enabled is True
        assert len(result.groups) == 1
        assert len(result.attributes) == 1

    def test_returns_none_for_none(self) -> None:
        assert user_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="user-2",
            external_user_id=None,
            tenant_id=None,
            username=None,
            email=None,
            first_name=None,
            last_name=None,
            status=None,
            roles=None,
            public_key=None,
            groups=None,
            totp_enabled=None,
            password_changed=None,
            enforced_in_rules=None,
            creation_date=None,
            update_date=None,
            last_login=None,
            attributes=None,
        )
        result = user_from_dto(dto)
        assert result is not None
        assert result.groups == []
        assert result.attributes == []
        assert result.totp_enabled is False


class TestUsersFromDto:
    """Tests for users_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert users_from_dto(None) == []


class TestUserGroupFromDto:
    """Tests for user_group_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(id="g-1", name="Admins")
        result = user_group_from_dto(dto)
        assert result.id == "g-1"
        assert result.name == "Admins"


class TestUserAttributeFromDto:
    """Tests for user_attribute_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(id="a1", key="dept", value="IT")
        result = user_attribute_from_dto(dto)
        assert result.id == "a1"
        assert result.key == "dept"
        assert result.value == "IT"


class TestGroupFromDto:
    """Tests for group_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="g-1",
            external_group_id="ext-g-1",
            tenant_id="t-1",
            name="Operators",
            email="ops@example.com",
            description="Operations team",
            users=[
                SimpleNamespace(id="u-1", email="op1@example.com"),
            ],
            enforced_in_rules=True,
            creation_date="2024-01-01T00:00:00Z",
            update_date="2024-06-01T00:00:00Z",
        )
        result = group_from_dto(dto)
        assert result is not None
        assert result.id == "g-1"
        assert result.name == "Operators"
        assert result.enforced_in_rules is True
        assert len(result.users) == 1

    def test_returns_none_for_none(self) -> None:
        assert group_from_dto(None) is None


class TestGroupsFromDto:
    """Tests for groups_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert groups_from_dto(None) == []


class TestGroupUserFromDto:
    """Tests for group_user_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(id="u-1", email="user@example.com")
        result = group_user_from_dto(dto)
        assert result.id == "u-1"
        assert result.email == "user@example.com"


class TestTagFromDto:
    """Tests for tag_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="tag-1",
            value="important",
            color="#FF0000",
            creation_date="2024-01-01T00:00:00Z",
        )
        result = tag_from_dto(dto)
        assert result is not None
        assert result.id == "tag-1"
        assert result.name == "important"
        assert result.color == "#FF0000"

    def test_returns_none_for_none(self) -> None:
        assert tag_from_dto(None) is None


class TestTagsFromDto:
    """Tests for tags_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert tags_from_dto(None) == []
