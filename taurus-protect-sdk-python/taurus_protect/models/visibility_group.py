"""Visibility group domain models for Taurus-PROTECT SDK."""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional


@dataclass
class VisibilityGroupUser:
    """User within a visibility group."""

    id: str = ""
    email: str = ""
    name: str = ""


@dataclass
class VisibilityGroup:
    """
    Represents a restricted visibility group in Taurus-PROTECT.

    Visibility groups control which users can see and access
    specific wallets and resources.
    """

    id: str = ""
    tenant_id: str = ""
    name: str = ""
    description: str = ""
    user_count: int = 0
    users: List[VisibilityGroupUser] = field(default_factory=list)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None
