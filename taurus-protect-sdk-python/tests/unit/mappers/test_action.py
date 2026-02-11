"""Unit tests for action mapper."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.action import (
    action_attribute_from_dto,
    action_details_from_dto,
    action_from_dto,
    action_trail_from_dto,
    actions_from_dto,
)


class TestActionFromDto:
    """Tests for action_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="act-1",
            tenant_id="t-1",
            label="Auto-rebalance",
            status="ACTIVE",
            auto_approve=True,
            action=SimpleNamespace(type="TRANSFER", parameters={"amount": "100"}),
            attributes=[
                SimpleNamespace(id="a1", key="priority", value="high"),
            ],
            trails=[
                SimpleNamespace(
                    id="tr-1", user_id="u-1", action="created",
                    status="OK", timestamp="2024-01-15T10:30:00Z",
                ),
            ],
            creation_date="2024-01-01T00:00:00Z",
            update_date="2024-06-01T00:00:00Z",
            lastcheckeddate="2024-06-02T00:00:00Z",
        )
        result = action_from_dto(dto)
        assert result is not None
        assert result.id == "act-1"
        assert result.tenant_id == "t-1"
        assert result.label == "Auto-rebalance"
        assert result.status == "ACTIVE"
        assert result.auto_approve is True
        assert result.action is not None
        assert result.action.type == "TRANSFER"
        assert len(result.attributes) == 1
        assert len(result.trails) == 1

    def test_returns_none_for_none(self) -> None:
        assert action_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="act-2",
            tenant_id=None,
            label=None,
            status=None,
            auto_approve=None,
            action=None,
            attributes=None,
            trails=None,
            creation_date=None,
            update_date=None,
            lastcheckeddate=None,
        )
        result = action_from_dto(dto)
        assert result is not None
        assert result.action is None
        assert result.attributes == []
        assert result.trails == []


class TestActionsFromDto:
    """Tests for actions_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert actions_from_dto(None) == []

    def test_returns_empty_for_empty(self) -> None:
        assert actions_from_dto([]) == []

    def test_filters_none_entries(self) -> None:
        result = actions_from_dto([None])
        assert result == []


class TestActionAttributeFromDto:
    """Tests for action_attribute_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(id="a1", key="priority", value="high")
        result = action_attribute_from_dto(dto)
        assert result is not None
        assert result.id == "a1"
        assert result.key == "priority"

    def test_returns_none_for_none(self) -> None:
        assert action_attribute_from_dto(None) is None


class TestActionTrailFromDto:
    """Tests for action_trail_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(
            id="tr-1", user_id="u-1", action="approved",
            status="OK", timestamp="2024-01-15T10:30:00Z",
        )
        result = action_trail_from_dto(dto)
        assert result is not None
        assert result.id == "tr-1"
        assert result.user_id == "u-1"
        assert result.action == "approved"

    def test_returns_none_for_none(self) -> None:
        assert action_trail_from_dto(None) is None


class TestActionDetailsFromDto:
    """Tests for action_details_from_dto function."""

    def test_maps_fields(self) -> None:
        dto = SimpleNamespace(type="TRANSFER", parameters={"key": "val"})
        result = action_details_from_dto(dto)
        assert result is not None
        assert result.type == "TRANSFER"
        assert result.parameters == {"key": "val"}

    def test_handles_none_parameters(self) -> None:
        dto = SimpleNamespace(type="NOTIFY", parameters=None)
        result = action_details_from_dto(dto)
        assert result is not None
        assert result.parameters == {}

    def test_returns_none_for_none(self) -> None:
        assert action_details_from_dto(None) is None
