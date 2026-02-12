"""Audit, Change, and Job models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import Dict, List, Optional

from pydantic import BaseModel, Field


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
    """Result of listing changes with cursor pagination."""

    changes: List[Change] = Field(default_factory=list)
    current_page: Optional[str] = Field(default=None)
    has_next: bool = Field(default=False)

    model_config = {"frozen": True}


class ListChangesOptions(BaseModel):
    """Options for listing changes."""

    entity: Optional[str] = None
    status: Optional[str] = None
    creator_id: Optional[str] = None
    sort_order: Optional[str] = None
    page_size: Optional[int] = None
    current_page: Optional[str] = None
    page_request: Optional[str] = None
    entity_ids: Optional[List[str]] = None
    entity_uuids: Optional[List[str]] = None


class Job(BaseModel):
    """Job record."""

    id: str = Field(description="Unique job identifier")
    type: str = Field(default="", description="Type of the job")
    timestamp: Optional[datetime] = Field(default=None, description="When the job was created")
    description: str = Field(default="", description="Human-readable description")

    model_config = {"frozen": True}
