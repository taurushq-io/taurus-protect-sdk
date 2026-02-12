"""Webhook and Config models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, Field


class Webhook(BaseModel):
    """
    Webhook configuration.

    Webhooks allow you to receive notifications when events occur in Taurus-PROTECT.

    Attributes:
        id: Unique webhook identifier.
        type: Webhook type/event type.
        url: The URL that receives webhook notifications.
        status: Webhook status (e.g., "ACTIVE", "PAUSED").
        timeout_until: If in timeout, when the timeout expires.
        created_at: When the webhook was created.
        updated_at: When the webhook was last updated.
    """

    id: str = Field(description="Unique webhook identifier")
    type: str = Field(default="", description="Webhook type/event type")
    url: str = Field(default="", description="Webhook notification URL")
    status: str = Field(default="", description="Webhook status")
    timeout_until: Optional[datetime] = Field(default=None, description="Timeout expiration")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}


class WebhookCall(BaseModel):
    """
    Record of a webhook call.

    Represents a single invocation of a webhook, including the payload and status.

    Attributes:
        id: Unique call identifier.
        event_id: ID of the event that triggered the call.
        webhook_id: ID of the webhook that was called.
        payload: The JSON payload sent to the webhook.
        status: Call status (e.g., "SUCCESS", "FAILED").
        status_message: Additional status information.
        attempts: Number of delivery attempts.
        created_at: When the call was made.
        updated_at: When the call record was last updated.
    """

    id: str = Field(description="Unique call identifier")
    event_id: str = Field(default="", description="Event ID that triggered the call")
    webhook_id: str = Field(default="", description="Webhook ID")
    payload: str = Field(default="", description="JSON payload sent")
    status: str = Field(default="", description="Call status")
    status_message: str = Field(default="", description="Status message")
    attempts: int = Field(default=0, description="Number of delivery attempts")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}


class TenantConfig(BaseModel):
    """
    Tenant configuration.

    Contains configuration settings for the tenant.

    Attributes:
        tenant_id: The tenant identifier.
        base_currency: Base currency for the tenant.
        super_admin_minimum_signatures: Minimum signatures required from SuperAdmins.
        is_mfa_mandatory: Whether MFA is mandatory for users.
        exclude_container: Whether to exclude container.
        fee_limit_factor: Fee limit multiplier.
        protect_engine_version: Version of the protect engine.
        restrict_sources_for_whitelisted_addresses: Whether to restrict sources.
        is_protect_engine_cold: Whether protect engine is in cold mode.
        is_cold_protect_engine_offline: Whether cold protect engine is offline.
        is_physical_air_gap_enabled: Whether physical air gap is enabled.
    """

    tenant_id: str = Field(default="", description="Tenant identifier")
    base_currency: str = Field(default="", description="Base currency")
    super_admin_minimum_signatures: int = Field(
        default=0, description="Minimum SuperAdmin signatures"
    )
    is_mfa_mandatory: bool = Field(default=False, description="MFA mandatory flag")
    exclude_container: bool = Field(default=False, description="Exclude container flag")
    fee_limit_factor: float = Field(default=0.0, description="Fee limit factor")
    protect_engine_version: str = Field(default="", description="Protect engine version")
    restrict_sources_for_whitelisted_addresses: bool = Field(
        default=False, description="Restrict sources for whitelisted addresses"
    )
    is_protect_engine_cold: bool = Field(default=False, description="Protect engine cold mode")
    is_cold_protect_engine_offline: bool = Field(
        default=False, description="Cold protect engine offline"
    )
    is_physical_air_gap_enabled: bool = Field(default=False, description="Physical air gap enabled")

    model_config = {"frozen": True}


class Feature(BaseModel):
    """
    Feature flag.

    Represents an enabled feature for the tenant.

    Attributes:
        name: Feature name.
        enabled: Whether the feature is enabled.
    """

    name: str = Field(description="Feature name")
    enabled: bool = Field(default=True, description="Whether feature is enabled")

    model_config = {"frozen": True}
