"""Audit, Change, and Job models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import Dict, List, Optional

from pydantic import BaseModel, Field

from taurus_protect.models.pagination import CursorListOptions, CursorPage


class Audit(BaseModel):
    """Audit event record."""

    id: str = Field(description="Unique audit event identifier")
    type: str = Field(default="", description="Type of the audit event")
    timestamp: Optional[datetime] = Field(default=None, description="When the event occurred")
    description: str = Field(default="", description="Human-readable description")

    model_config = {"frozen": True}


class Change(BaseModel):
    """Change record representing a configuration modification requiring approval."""

    id: str = Field(description="Unique change identifier")
    tenant_id: int = Field(default=0, description="Tenant ID")
    creator_id: str = Field(default="", description="Internal creator user ID")
    creator_external_id: str = Field(default="", description="External creator user ID")
    action: str = Field(default="", description="Action type (create, update, delete)")
    entity: str = Field(default="", description="Entity type (businessrule, user, group, etc)")
    entity_id: str = Field(default="", description="Entity ID")
    entity_uuid: str = Field(default="", description="Entity UUID")
    changes: Optional[Dict[str, str]] = Field(default=None, description="Map of field changes")
    comment: str = Field(default="", description="Change description")
    created_at: Optional[datetime] = Field(default=None, description="When the change was created")

    model_config = {"frozen": True}


class CreateChangeRequest(BaseModel):
    """Request to create a configuration change."""

    action: str = Field(description="Action type")
    entity: str = Field(description="Entity type")
    entity_id: Optional[str] = Field(default=None, description="Entity ID")
    changes: Optional[Dict[str, str]] = Field(default=None, description="Field changes")
    comment: Optional[str] = Field(default=None, description="Optional comment")

    model_config = {"frozen": True}


class ChangeResult(BaseModel):
    """
    One page of changes.

    Continue with ``ListChangesOptions(cursor=result.page.next_cursor)`` while
    ``result.page.has_more`` is true.
    """

    changes: List[Change] = Field(default_factory=list)
    page: CursorPage = Field(default_factory=CursorPage)

    model_config = {"frozen": True}


class ListChangesOptions(CursorListOptions):
    """
    Options for listing changes. Every field reaches the wire or is refused by name.

    ``page_size`` defaults to 20 and may not exceed 100; continue with
    ``cursor=result.page.next_cursor``. The approval queue accepts neither ``status``,
    ``creator_id`` nor ``entity_id``.
    """

    entity: Optional[str] = None
    entity_id: Optional[str] = None
    status: Optional[str] = None
    creator_id: Optional[str] = None
    sort_order: Optional[str] = None
    entity_ids: Optional[List[str]] = None
    entity_uuids: Optional[List[str]] = None


class Job(BaseModel):
    """Job record. A job is identified by its name, which is also its ``id``."""

    id: str = Field(description="Unique job identifier (the job name)")
    name: str = Field(default="", description="Job name")
    type: str = Field(default="", description="Type of the job")
    timestamp: Optional[datetime] = Field(default=None, description="When the job was created")
    description: str = Field(default="", description="Human-readable description")

    model_config = {"frozen": True}
