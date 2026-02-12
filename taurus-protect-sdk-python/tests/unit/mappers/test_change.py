"""Unit tests for change mapper functions."""

from datetime import datetime, timezone
from types import SimpleNamespace

import pytest

from taurus_protect.mappers.audit import (
    change_from_dto,
    changes_from_dto,
    create_change_request_to_dto,
)
from taurus_protect.models.audit import CreateChangeRequest


class TestChangeFromDto:
    """Tests for change_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="change-1",
            tenant_id="42",
            creator_id="user-1",
            creator_external_id="ext-1",
            action="update",
            entity="businessrule",
            entity_id="100",
            entity_uuid="uuid-abc",
            changes={"rulevalue": "10"},
            comment="Test change",
            creation_date="2024-06-15T10:30:00Z",
        )
        result = change_from_dto(dto)
        assert result is not None
        assert result.id == "change-1"
        assert result.tenant_id == 42
        assert result.creator_id == "user-1"
        assert result.creator_external_id == "ext-1"
        assert result.action == "update"
        assert result.entity == "businessrule"
        assert result.entity_id == "100"
        assert result.entity_uuid == "uuid-abc"
        assert result.changes == {"rulevalue": "10"}
        assert result.comment == "Test change"
        assert result.created_at is not None

    def test_maps_camel_case_fields(self) -> None:
        dto = SimpleNamespace(
            id="change-2",
            tenantId="10",
            creatorId="user-2",
            creatorExternalId="ext-2",
            action="create",
            entity="user",
            entityId="200",
            entityUUID="uuid-def",
            changes=None,
            comment="",
            creationDate="2024-01-01T00:00:00Z",
        )
        result = change_from_dto(dto)
        assert result is not None
        assert result.tenant_id == 10
        assert result.creator_id == "user-2"
        assert result.creator_external_id == "ext-2"
        assert result.entity_id == "200"
        assert result.entity_uuid == "uuid-def"
        assert result.created_at is not None

    def test_returns_none_for_none(self) -> None:
        assert change_from_dto(None) is None

    def test_handles_none_fields(self) -> None:
        dto = SimpleNamespace(
            id="change-3",
            tenant_id=None,
            creator_id=None,
            creator_external_id=None,
            action=None,
            entity=None,
            entity_id=None,
            entity_uuid=None,
            changes=None,
            comment=None,
            creation_date=None,
        )
        result = change_from_dto(dto)
        assert result is not None
        assert result.id == "change-3"
        assert result.tenant_id == 0
        assert result.creator_id == ""
        assert result.action == ""
        assert result.entity == ""
        assert result.changes is None
        assert result.created_at is None

    def test_handles_datetime_object(self) -> None:
        ts = datetime(2024, 1, 15, 12, 0, 0, tzinfo=timezone.utc)
        dto = SimpleNamespace(
            id="change-4",
            tenant_id="1",
            creator_id="u1",
            creator_external_id=None,
            action="update",
            entity="wallet",
            entity_id="5",
            entity_uuid=None,
            changes=None,
            comment=None,
            creation_date=ts,
        )
        result = change_from_dto(dto)
        assert result is not None
        assert result.created_at == ts


class TestChangesFromDto:
    """Tests for changes_from_dto function."""

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", tenant_id=None, creator_id=None, creator_external_id=None,
                action="create", entity="user", entity_id=None, entity_uuid=None,
                changes=None, comment=None, creation_date=None,
            ),
            SimpleNamespace(
                id="2", tenant_id=None, creator_id=None, creator_external_id=None,
                action="update", entity="wallet", entity_id=None, entity_uuid=None,
                changes=None, comment=None, creation_date=None,
            ),
        ]
        result = changes_from_dto(dtos)
        assert len(result) == 2
        assert result[0].id == "1"
        assert result[1].id == "2"

    def test_returns_empty_for_none(self) -> None:
        assert changes_from_dto(None) == []

    def test_returns_empty_for_empty_list(self) -> None:
        assert changes_from_dto([]) == []

    def test_filters_out_none_dtos(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", tenant_id=None, creator_id=None, creator_external_id=None,
                action="create", entity="user", entity_id=None, entity_uuid=None,
                changes=None, comment=None, creation_date=None,
            ),
            None,
        ]
        result = changes_from_dto(dtos)
        assert len(result) == 1
        assert result[0].id == "1"


class TestCreateChangeRequestToDto:
    """Tests for create_change_request_to_dto function."""

    def test_maps_all_fields(self) -> None:
        request = CreateChangeRequest(
            action="update",
            entity="businessrule",
            entity_id="42",
            changes={"rulevalue": "100"},
            comment="Test comment",
        )
        dto = create_change_request_to_dto(request)
        assert dto["action"] == "update"
        assert dto["entity"] == "businessrule"
        assert dto["entityId"] == "42"
        assert dto["changes"] == {"rulevalue": "100"}
        assert dto["changeComment"] == "Test comment"

    def test_omits_none_fields(self) -> None:
        request = CreateChangeRequest(
            action="create",
            entity="user",
        )
        dto = create_change_request_to_dto(request)
        assert dto["action"] == "create"
        assert dto["entity"] == "user"
        assert "entityId" not in dto
        assert "changes" not in dto
        assert "changeComment" not in dto
