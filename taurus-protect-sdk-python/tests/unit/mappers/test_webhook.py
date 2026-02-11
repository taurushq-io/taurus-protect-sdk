"""Unit tests for webhook, webhook call, and tenant config mappers."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.webhook import (
    tenant_config_from_dto,
    webhook_call_from_dto,
    webhook_calls_from_dto,
    webhook_from_dto,
    webhooks_from_dto,
)


class TestWebhookFromDto:
    """Tests for webhook_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="wh-1",
            type="TRANSACTION",
            url="https://example.com/webhook",
            status="ACTIVE",
            timeout_until=None,
            created_at="2024-01-01T00:00:00Z",
            updated_at="2024-06-01T00:00:00Z",
        )
        result = webhook_from_dto(dto)
        assert result is not None
        assert result.id == "wh-1"
        assert result.type == "TRANSACTION"
        assert result.url == "https://example.com/webhook"
        assert result.status == "ACTIVE"

    def test_returns_none_for_none(self) -> None:
        assert webhook_from_dto(None) is None


class TestWebhooksFromDto:
    """Tests for webhooks_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert webhooks_from_dto(None) == []

    def test_maps_list(self) -> None:
        dtos = [
            SimpleNamespace(
                id="1", type="TX", url="https://example.com/1",
                status="ACTIVE", timeout_until=None,
                created_at=None, updated_at=None,
            ),
        ]
        result = webhooks_from_dto(dtos)
        assert len(result) == 1


class TestWebhookCallFromDto:
    """Tests for webhook_call_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            id="wc-1",
            event_id="evt-1",
            webhook_id="wh-1",
            payload='{"event":"transaction"}',
            status="SUCCESS",
            status_message="200 OK",
            attempts=3,
            created_at="2024-01-15T10:30:00Z",
            updated_at="2024-01-15T10:31:00Z",
        )
        result = webhook_call_from_dto(dto)
        assert result is not None
        assert result.id == "wc-1"
        assert result.event_id == "evt-1"
        assert result.webhook_id == "wh-1"
        assert result.status == "SUCCESS"
        assert result.attempts == 3

    def test_returns_none_for_none(self) -> None:
        assert webhook_call_from_dto(None) is None

    def test_handles_string_attempts(self) -> None:
        dto = SimpleNamespace(
            id="wc-2",
            event_id=None,
            webhook_id=None,
            payload=None,
            status=None,
            status_message=None,
            attempts="5",
            created_at=None,
            updated_at=None,
        )
        result = webhook_call_from_dto(dto)
        assert result is not None
        assert result.attempts == 5


class TestWebhookCallsFromDto:
    """Tests for webhook_calls_from_dto function."""

    def test_returns_empty_for_none(self) -> None:
        assert webhook_calls_from_dto(None) == []


class TestTenantConfigFromDto:
    """Tests for tenant_config_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            tenant_id="t-1",
            base_currency="USD",
            super_admin_minimum_signatures=2,
            is_mfa_mandatory=True,
            exclude_container=False,
            fee_limit_factor=1.5,
            protect_engine_version="2.0",
            restrict_sources_for_whitelisted_addresses=True,
            is_protect_engine_cold=False,
            is_cold_protect_engine_offline=False,
            is_physical_air_gap_enabled=True,
        )
        result = tenant_config_from_dto(dto)
        assert result is not None
        assert result.tenant_id == "t-1"
        assert result.base_currency == "USD"
        assert result.super_admin_minimum_signatures == 2
        assert result.is_mfa_mandatory is True
        assert result.fee_limit_factor == 1.5
        assert result.is_physical_air_gap_enabled is True

    def test_returns_none_for_none(self) -> None:
        assert tenant_config_from_dto(None) is None

    def test_handles_string_min_signatures(self) -> None:
        dto = SimpleNamespace(
            tenant_id="t-2",
            base_currency=None,
            super_admin_minimum_signatures="3",
            is_mfa_mandatory=None,
            exclude_container=None,
            fee_limit_factor=None,
            protect_engine_version=None,
            restrict_sources_for_whitelisted_addresses=None,
            is_protect_engine_cold=None,
            is_cold_protect_engine_offline=None,
            is_physical_air_gap_enabled=None,
        )
        result = tenant_config_from_dto(dto)
        assert result is not None
        assert result.super_admin_minimum_signatures == 3
