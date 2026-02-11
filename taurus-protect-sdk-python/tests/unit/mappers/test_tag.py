"""Unit tests for tag mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

from taurus_protect.mappers.user import tag_from_dto, tags_from_dto


class TestTagFromDto:
    """Tests for tag_from_dto function."""

    def test_maps_all_fields(self) -> None:
        created = datetime(2024, 1, 15, 10, 30, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id="tag-1",
            value="Important",
            color="#FF0000",
            creation_date=created,
        )
        result = tag_from_dto(dto)
        assert result is not None
        assert result.id == "tag-1"
        assert result.name == "Important"  # API 'value' mapped to 'name'
        assert result.color == "#FF0000"
        assert result.created_at == created

    def test_returns_none_for_none(self) -> None:
        assert tag_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="tag-2",
            value="Simple",
            color=None,
            creation_date=None,
        )
        result = tag_from_dto(dto)
        assert result is not None
        assert result.id == "tag-2"
        assert result.name == "Simple"
        assert result.color is None
        assert result.created_at is None


class TestTagsFromDto:
    """Tests for tags_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(id="t1", value="Alpha", color="#000", creation_date=None),
            SimpleNamespace(id="t2", value="Beta", color="#FFF", creation_date=None),
        ]
        result = tags_from_dto(dtos)
        assert len(result) == 2
        assert result[0].name == "Alpha"
        assert result[1].name == "Beta"

    def test_returns_empty_for_none(self) -> None:
        assert tags_from_dto(None) == []

    def test_filters_none_dtos(self) -> None:
        dtos = [
            None,
            SimpleNamespace(id="t1", value="Tag", color=None, creation_date=None),
        ]
        result = tags_from_dto(dtos)
        assert len(result) == 1
