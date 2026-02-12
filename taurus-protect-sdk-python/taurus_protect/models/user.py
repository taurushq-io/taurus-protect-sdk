"""User, Group, and Tag models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, Field


class UserGroup(BaseModel):
    """Group membership information for a user."""

    id: str = Field(description="Group identifier")
    name: str = Field(default="", description="Group name")

    model_config = {"frozen": True}


class UserAttribute(BaseModel):
    """Custom attribute on a user."""

    id: str = Field(description="Attribute identifier")
    key: str = Field(description="Attribute name")
    value: str = Field(description="Attribute value")

    model_config = {"frozen": True}


class User(BaseModel):
    """
    System user.

    Represents a user account in the Taurus-PROTECT system.

    Attributes:
        id: Unique user identifier.
        external_user_id: External user ID for integration.
        tenant_id: Tenant identifier.
        username: Username for login.
        email: User email address.
        first_name: User's first name.
        last_name: User's last name.
        status: Account status (e.g., "ACTIVE", "DISABLED").
        roles: List of role names assigned to the user.
        public_key: User's public key for signing.
        groups: Groups the user belongs to.
        totp_enabled: Whether TOTP 2FA is enabled.
        password_changed: Whether password has been changed.
        enforced_in_rules: Whether user is enforced in governance rules.
        created_at: When the user was created.
        updated_at: When the user was last updated.
        last_login: When the user last logged in.
        attributes: Custom key-value attributes.
    """

    id: str = Field(description="Unique user identifier")
    external_user_id: Optional[str] = Field(default=None, description="External user ID")
    tenant_id: Optional[str] = Field(default=None, description="Tenant identifier")
    username: Optional[str] = Field(default=None, description="Username for login")
    email: Optional[str] = Field(default=None, description="User email address")
    first_name: Optional[str] = Field(default=None, description="User's first name")
    last_name: Optional[str] = Field(default=None, description="User's last name")
    status: Optional[str] = Field(default=None, description="Account status")
    roles: List[str] = Field(default_factory=list, description="Assigned roles")
    public_key: Optional[str] = Field(default=None, description="User's public key")
    groups: List[UserGroup] = Field(default_factory=list, description="Group memberships")
    totp_enabled: bool = Field(default=False, description="Whether TOTP 2FA is enabled")
    password_changed: bool = Field(default=False, description="Whether password has been changed")
    enforced_in_rules: bool = Field(default=False, description="Whether user is enforced in rules")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")
    last_login: Optional[datetime] = Field(default=None, description="Last login timestamp")
    attributes: List[UserAttribute] = Field(default_factory=list, description="Custom attributes")

    model_config = {"frozen": True}


class GroupUser(BaseModel):
    """User membership information within a group."""

    id: str = Field(description="User identifier")
    email: Optional[str] = Field(default=None, description="User email")

    model_config = {"frozen": True}


class Group(BaseModel):
    """
    User group.

    Represents a group of users for access control and governance rules.

    Attributes:
        id: Unique group identifier.
        external_group_id: External group ID for integration.
        tenant_id: Tenant identifier.
        name: Group name.
        email: Group email address.
        description: Group description.
        users: Users in the group.
        enforced_in_rules: Whether group is enforced in governance rules.
        created_at: When the group was created.
        updated_at: When the group was last updated.
    """

    id: str = Field(description="Unique group identifier")
    external_group_id: Optional[str] = Field(default=None, description="External group ID")
    tenant_id: Optional[str] = Field(default=None, description="Tenant identifier")
    name: str = Field(default="", description="Group name")
    email: Optional[str] = Field(default=None, description="Group email address")
    description: Optional[str] = Field(default=None, description="Group description")
    users: List[GroupUser] = Field(default_factory=list, description="Users in the group")
    enforced_in_rules: bool = Field(default=False, description="Whether group is enforced in rules")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}


class Tag(BaseModel):
    """
    Tag for categorizing entities.

    Tags can be applied to wallets, addresses, and other entities
    for organization and filtering.

    Attributes:
        id: Unique tag identifier.
        name: Tag name/value.
        color: Tag color (hex code or color name).
        created_at: When the tag was created.
    """

    id: str = Field(description="Unique tag identifier")
    name: str = Field(description="Tag name/value")
    color: Optional[str] = Field(default=None, description="Tag color")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")

    model_config = {"frozen": True}
