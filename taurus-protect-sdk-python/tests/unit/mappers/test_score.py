"""Unit tests for score mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

from taurus_protect.mappers.statistics import score_from_dto, scores_from_dto


class TestScoreFromDto:
    """Tests for score_from_dto function."""

    def test_maps_all_fields(self) -> None:
        updated = datetime(2024, 6, 15, 12, 0, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id="score-1",
            provider="chainalysis",
            type="risk",
            score="85",
            update_date=updated,
        )
        result = score_from_dto(dto)
        assert result is not None
        assert result.id == "score-1"
        assert result.provider == "chainalysis"
        assert result.score_type == "risk"
        assert result.score == "85"
        assert result.updated_at == updated

    def test_returns_none_for_none(self) -> None:
        assert score_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id=None,
            provider=None,
            type=None,
            score=None,
            update_date=None,
        )
        result = score_from_dto(dto)
        assert result is not None
        assert result.id == ""
        assert result.provider == ""
        assert result.score_type == ""
        assert result.score == ""
        assert result.updated_at is None


class TestScoresFromDto:
    """Tests for scores_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="s1", provider="chainalysis", type="risk",
                score="90", update_date=None,
            ),
            SimpleNamespace(
                id="s2", provider="elliptic", type="trust",
                score="75", update_date=None,
            ),
        ]
        result = scores_from_dto(dtos)
        assert len(result) == 2
        assert result[0].provider == "chainalysis"
        assert result[1].provider == "elliptic"

    def test_returns_empty_for_none(self) -> None:
        assert scores_from_dto(None) == []

    def test_filters_none_dtos(self) -> None:
        dtos = [
            None,
            SimpleNamespace(
                id="s1", provider="p", type="risk", score="50", update_date=None,
            ),
        ]
        result = scores_from_dto(dtos)
        assert len(result) == 1
