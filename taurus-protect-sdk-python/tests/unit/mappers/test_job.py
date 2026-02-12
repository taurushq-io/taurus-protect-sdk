"""Unit tests for job mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

from taurus_protect.mappers.audit import job_from_dto, jobs_from_dto


class TestJobFromDto:
    """Tests for job_from_dto function."""

    def test_maps_all_fields(self) -> None:
        ts = datetime(2024, 6, 15, 12, 0, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id="job-1",
            type="WALLET_SYNC",
            timestamp=ts,
            description="Sync wallet balances",
        )
        result = job_from_dto(dto)
        assert result is not None
        assert result.id == "job-1"
        assert result.type == "WALLET_SYNC"
        assert result.timestamp == ts
        assert result.description == "Sync wallet balances"

    def test_returns_none_for_none(self) -> None:
        assert job_from_dto(None) is None

    def test_handles_missing_optional_fields(self) -> None:
        dto = SimpleNamespace(
            id="job-2",
            type=None,
            timestamp=None,
            description=None,
        )
        result = job_from_dto(dto)
        assert result is not None
        assert result.id == "job-2"
        assert result.type == ""
        assert result.timestamp is None
        assert result.description == ""


class TestJobsFromDto:
    """Tests for jobs_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="j1", type="SYNC", timestamp=None, description="Sync",
            ),
            SimpleNamespace(
                id="j2", type="EXPORT", timestamp=None, description="Export",
            ),
        ]
        result = jobs_from_dto(dtos)
        assert len(result) == 2
        assert result[0].id == "j1"
        assert result[1].id == "j2"

    def test_returns_empty_for_none(self) -> None:
        assert jobs_from_dto(None) == []

    def test_filters_none_dtos(self) -> None:
        dtos = [
            None,
            SimpleNamespace(id="j1", type="SYNC", timestamp=None, description=""),
        ]
        result = jobs_from_dto(dtos)
        assert len(result) == 1
