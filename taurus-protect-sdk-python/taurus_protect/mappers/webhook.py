"""Webhook mapper for converting OpenAPI DTOs to domain models."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect.mappers._base import (
    safe_bool,
    safe_datetime,
    safe_float,
    safe_int,
    safe_string,
)
from taurus_protect.models.webhook import TenantConfig, Webhook, WebhookCall

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


def webhook_from_dto(dto: Any) -> Optional[Webhook]:
    """
    Convert OpenAPI TgvalidatordWebhook to domain Webhook.

    Args:
        dto: OpenAPI webhook DTO (TgvalidatordWebhook).

    Returns:
        Domain Webhook model or None if dto is None.
    """
    if dto is None:
        return None

    return Webhook(
        id=safe_string(getattr(dto, "id", None)),
        type=safe_string(getattr(dto, "type", None)),
        url=safe_string(getattr(dto, "url", None)),
        status=safe_string(getattr(dto, "status", None)),
        timeout_until=safe_datetime(getattr(dto, "timeout_until", None)),
        created_at=safe_datetime(getattr(dto, "created_at", None)),
        updated_at=safe_datetime(getattr(dto, "updated_at", None)),
    )


def webhooks_from_dto(dtos: Optional[List[Any]]) -> List[Webhook]:
    """
    Convert list of OpenAPI webhook DTOs to domain Webhooks.

    Args:
        dtos: List of OpenAPI webhook DTOs.

    Returns:
        List of domain Webhook models.
    """
    if dtos is None:
        return []
    return [w for dto in dtos if (w := webhook_from_dto(dto)) is not None]


def webhook_call_from_dto(dto: Any) -> Optional[WebhookCall]:
    """
    Convert OpenAPI TgvalidatordWebhookCall to domain WebhookCall.

    Args:
        dto: OpenAPI webhook call DTO (TgvalidatordWebhookCall).

    Returns:
        Domain WebhookCall model or None if dto is None.
    """
    if dto is None:
        return None

    # Attempts may be a string from the API
    attempts = safe_int(getattr(dto, "attempts", None))

    return WebhookCall(
        id=safe_string(getattr(dto, "id", None)),
        event_id=safe_string(getattr(dto, "event_id", None)),
        webhook_id=safe_string(getattr(dto, "webhook_id", None)),
        payload=safe_string(getattr(dto, "payload", None)),
        status=safe_string(getattr(dto, "status", None)),
        status_message=safe_string(getattr(dto, "status_message", None)),
        attempts=attempts,
        created_at=safe_datetime(getattr(dto, "created_at", None)),
        updated_at=safe_datetime(getattr(dto, "updated_at", None)),
    )


def webhook_calls_from_dto(dtos: Optional[List[Any]]) -> List[WebhookCall]:
    """
    Convert list of OpenAPI webhook call DTOs to domain WebhookCalls.

    Args:
        dtos: List of OpenAPI webhook call DTOs.

    Returns:
        List of domain WebhookCall models.
    """
    if dtos is None:
        return []
    return [c for dto in dtos if (c := webhook_call_from_dto(dto)) is not None]


def tenant_config_from_dto(dto: Any) -> Optional[TenantConfig]:
    """
    Convert OpenAPI TgvalidatordTenantConfig to domain TenantConfig.

    Args:
        dto: OpenAPI tenant config DTO (TgvalidatordTenantConfig).

    Returns:
        Domain TenantConfig model or None if dto is None.
    """
    if dto is None:
        return None

    # super_admin_minimum_signatures may be a string from the API
    super_admin_min_sigs = safe_int(getattr(dto, "super_admin_minimum_signatures", None))

    return TenantConfig(
        tenant_id=safe_string(getattr(dto, "tenant_id", None)),
        base_currency=safe_string(getattr(dto, "base_currency", None)),
        super_admin_minimum_signatures=super_admin_min_sigs,
        is_mfa_mandatory=safe_bool(getattr(dto, "is_mfa_mandatory", None)),
        exclude_container=safe_bool(getattr(dto, "exclude_container", None)),
        fee_limit_factor=safe_float(getattr(dto, "fee_limit_factor", None)),
        protect_engine_version=safe_string(getattr(dto, "protect_engine_version", None)),
        restrict_sources_for_whitelisted_addresses=safe_bool(
            getattr(dto, "restrict_sources_for_whitelisted_addresses", None)
        ),
        is_protect_engine_cold=safe_bool(getattr(dto, "is_protect_engine_cold", None)),
        is_cold_protect_engine_offline=safe_bool(
            getattr(dto, "is_cold_protect_engine_offline", None)
        ),
        is_physical_air_gap_enabled=safe_bool(getattr(dto, "is_physical_air_gap_enabled", None)),
    )
