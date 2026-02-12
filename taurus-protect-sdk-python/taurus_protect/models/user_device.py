"""User device domain models for Taurus-PROTECT SDK."""

from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime
from typing import Optional


@dataclass
class UserDevicePairing:
    """
    Represents a user device pairing request/status.

    Used for managing the pairing of user devices for
    multi-factor authentication.
    """

    pairing_id: str = ""
    status: str = ""
    created_at: Optional[datetime] = None
    expires_at: Optional[datetime] = None


@dataclass
class UserDevicePairingInfo:
    """
    Detailed information about a user device pairing.

    Contains the full pairing status and device information.
    """

    pairing_id: str = ""
    user_id: str = ""
    status: str = ""
    device_name: str = ""
    device_type: str = ""
    encryption_key: str = ""
    created_at: Optional[datetime] = None
    expires_at: Optional[datetime] = None
