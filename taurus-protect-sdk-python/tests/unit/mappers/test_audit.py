"""Unit tests for audit, change, and job mappers."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.audit import (
    audit_from_dto,
    audits_from_dto,
    change_from_dto,
    changes_from_dto,
    job_from_dto,
    jobs_from_dto,
)


class TestAuditFromDto:
    """Tests for audit_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="aud-1",
            type="LOGIN",
            timestamp="2024-01-15T10:30:00Z",
            description="User logged in",
        )
        result = audit_from_dto(dto)
        assert result is not None
        assert result.id == "aud-1"
        assert result.type == "LOGIN"
        assert result.description == "User logged in"
        assert result.timestamp is not None

    def test_returns_none_for_none(self) -> None:
        assert audit_from_dto(None) is None


class TestAuditsFromDto:
    """Tests for audits_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert audits_from_dto(None) == []

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(id="1", type="A", timestamp=None, description="d1"),
            SimpleNamespace(id="2", type="B", timestamp=None, description="d2"),
        ]
        result = audits_from_dto(dtos)
        assert len(result) == 2


class TestChangeFromDto:
    """Tests for change_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="ch-1",
            tenant_id="42",
            creator_id="user-1",
            creator_external_id="ext-1",
            action="update",
            entity="businessrule",
            entity_id="100",
            entity_uuid="uuid-abc",
            changes={"rulevalue": "10"},
            comment="Config changed",
            creation_date="2024-01-15T10:30:00Z",
        )
        result = change_from_dto(dto)
        assert result is not None
        assert result.id == "ch-1"
        assert result.tenant_id == 42
        assert result.action == "update"
        assert result.entity == "businessrule"
        assert result.comment == "Config changed"
        assert result.created_at is not None

    def test_returns_none_for_none(self) -> None:
        assert change_from_dto(None) is None


class TestChangesFromDto:
    """Tests for changes_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert changes_from_dto(None) == []


class TestJobFromDto:
    """Tests for job_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="job-1",
            type="SYNC",
            timestamp="2024-01-15T10:30:00Z",
            description="Blockchain sync",
        )
        result = job_from_dto(dto)
        assert result is not None
        assert result.id == "job-1"
        assert result.type == "SYNC"

    def test_returns_none_for_none(self) -> None:
        assert job_from_dto(None) is None


class TestJobsFromDto:
    """Tests for jobs_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert jobs_from_dto(None) == []
