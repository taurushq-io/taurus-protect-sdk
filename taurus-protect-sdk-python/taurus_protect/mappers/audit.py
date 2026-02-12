"""Audit, Change, and Job mappers for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, Dict, List, Optional

from taurus_protect.mappers._base import safe_datetime, safe_int, safe_string
from taurus_protect.models.audit import Audit, Change, CreateChangeRequest, Job

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def audit_from_dto(dto: Any) -> Optional[Audit]:
    """Convert OpenAPI audit DTO to domain Audit."""
    if dto is None:
        return None

    return Audit(
        id=safe_string(getattr(dto, "id", None)),
        type=safe_string(getattr(dto, "type", None)),
        timestamp=safe_datetime(getattr(dto, "timestamp", None)),
        description=safe_string(getattr(dto, "description", None)),
    )


def audits_from_dto(dtos: Optional[List[Any]]) -> List[Audit]:
    """Convert list of OpenAPI audit DTOs to domain Audits."""
    if dtos is None:
        return []
    return [a for dto in dtos if (a := audit_from_dto(dto)) is not None]


def _resolve(dto: Any, snake_name: str, camel_name: str) -> Any:
    """Resolve a field value trying snake_case first, then camelCase."""
    val = getattr(dto, snake_name, None)
    if val is not None:
        return val
    return getattr(dto, camel_name, None)


def change_from_dto(dto: Any) -> Optional[Change]:
    """Convert OpenAPI change DTO to domain Change.

    Maps all fields from TgvalidatordChange:
    - creationDate -> created_at (matching Java @Mapping annotation)
    - tenantId string -> int (matching Java Integer.parseInt)
    """
    if dto is None:
        return None

    return Change(
        id=safe_string(getattr(dto, "id", None)),
        tenant_id=safe_int(_resolve(dto, "tenant_id", "tenantId")),
        creator_id=safe_string(_resolve(dto, "creator_id", "creatorId")),
        creator_external_id=safe_string(_resolve(dto, "creator_external_id", "creatorExternalId")),
        action=safe_string(getattr(dto, "action", None)),
        entity=safe_string(getattr(dto, "entity", None)),
        entity_id=safe_string(_resolve(dto, "entity_id", "entityId")),
        entity_uuid=safe_string(_resolve(dto, "entity_uuid", "entityUUID")),
        changes=getattr(dto, "changes", None),
        comment=safe_string(getattr(dto, "comment", None)),
        created_at=safe_datetime(_resolve(dto, "creation_date", "creationDate")),
    )


def changes_from_dto(dtos: Optional[List[Any]]) -> List[Change]:
    """Convert list of OpenAPI change DTOs to domain Changes."""
    if dtos is None:
        return []
    return [c for dto in dtos if (c := change_from_dto(dto)) is not None]


def create_change_request_to_dto(request: CreateChangeRequest) -> Dict[str, Any]:
    """Convert SDK CreateChangeRequest to OpenAPI DTO dict.

    Maps comment -> changeComment (matching Java ChangeMapper.toDTO).
    """
    dto: Dict[str, Any] = {
        "action": request.action,
        "entity": request.entity,
    }
    if request.entity_id is not None:
        dto["entityId"] = request.entity_id
    if request.changes is not None:
        dto["changes"] = request.changes
    if request.comment is not None:
        dto["changeComment"] = request.comment
    return dto


def job_from_dto(dto: Any) -> Optional[Job]:
    """Convert OpenAPI job DTO to domain Job."""
    if dto is None:
        return None

    return Job(
        id=safe_string(getattr(dto, "id", None)),
        type=safe_string(getattr(dto, "type", None)),
        timestamp=safe_datetime(getattr(dto, "timestamp", None)),
        description=safe_string(getattr(dto, "description", None)),
    )


def jobs_from_dto(dtos: Optional[List[Any]]) -> List[Job]:
    """Convert list of OpenAPI job DTOs to domain Jobs."""
    if dtos is None:
        return []
    return [j for dto in dtos if (j := job_from_dto(dto)) is not None]
