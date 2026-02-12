"""Action domain models for Taurus-PROTECT SDK."""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional


@dataclass
class ActionAttribute:
    """Attribute associated with an action."""

    id: str = ""
    key: str = ""
    value: str = ""


@dataclass
class ActionTrail:
    """Trail entry for action history."""

    id: str = ""
    user_id: str = ""
    action: str = ""
    status: str = ""
    timestamp: Optional[datetime] = None


@dataclass
class ActionDetails:
    """Details of the action to be performed."""

    type: str = ""
    parameters: dict = field(default_factory=dict)


@dataclass
class Action:
    """
    Represents an action in Taurus-PROTECT.

    Actions represent pending or completed operations that may require
    approval or have been executed.
    """

    id: str = ""
    tenant_id: str = ""
    label: str = ""
    status: str = ""
    auto_approve: bool = False
    action: Optional[ActionDetails] = None
    attributes: List[ActionAttribute] = field(default_factory=list)
    trails: List[ActionTrail] = field(default_factory=list)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None
    last_checked_at: Optional[datetime] = None
