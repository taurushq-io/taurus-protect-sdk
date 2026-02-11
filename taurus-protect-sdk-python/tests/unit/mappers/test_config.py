"""Unit tests for config (tenant_config) mapper functions."""

from types import SimpleNamespace

import pytest

from taurus_protect.mappers.webhook import tenant_config_from_dto


class TestTenantConfigFromDto:
    """Tests for tenant_config_from_dto function."""

    def test_maps_all_fields(self) -> None:
        dto = SimpleNamespace(
            tenant_id="tenant-1",
            base_currency="USD",
            super_admin_minimum_signatures=2,
            is_mfa_mandatory=True,
            exclude_container=False,
            fee_limit_factor=1.5,
            protect_engine_version="3.2.1",
            restrict_sources_for_whitelisted_addresses=True,
            is_protect_engine_cold=False,
            is_cold_protect_engine_offline=False,
            is_physical_air_gap_enabled=True,
        )
        result = tenant_config_from_dto(dto)
        assert result is not None
        assert result.tenant_id == "tenant-1"
        assert result.base_currency == "USD"
        assert result.super_admin_minimum_signatures == 2
        assert result.is_mfa_mandatory is True
        assert result.exclude_container is False
        assert result.fee_limit_factor == 1.5
        assert result.protect_engine_version == "3.2.1"
        assert result.restrict_sources_for_whitelisted_addresses is True
        assert result.is_protect_engine_cold is False
        assert result.is_cold_protect_engine_offline is False
        assert result.is_physical_air_gap_enabled is True

    def test_returns_none_for_none(self) -> None:
        assert tenant_config_from_dto(None) is None

    def test_handles_none_fields(self) -> None:
        dto = SimpleNamespace(
            tenant_id=None,
            base_currency=None,
            super_admin_minimum_signatures=None,
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
        assert result.tenant_id == ""
        assert result.base_currency == ""
        assert result.super_admin_minimum_signatures == 0
        assert result.is_mfa_mandatory is False
        assert result.exclude_container is False
        assert result.fee_limit_factor == 0.0
        assert result.protect_engine_version == ""
        assert result.restrict_sources_for_whitelisted_addresses is False
        assert result.is_protect_engine_cold is False
        assert result.is_cold_protect_engine_offline is False
        assert result.is_physical_air_gap_enabled is False

    def test_handles_string_minimum_signatures(self) -> None:
        dto = SimpleNamespace(
            tenant_id="t-1",
            base_currency="CHF",
            super_admin_minimum_signatures="3",
            is_mfa_mandatory=False,
            exclude_container=False,
            fee_limit_factor="2.0",
            protect_engine_version="4.0.0",
            restrict_sources_for_whitelisted_addresses=False,
            is_protect_engine_cold=False,
            is_cold_protect_engine_offline=False,
            is_physical_air_gap_enabled=False,
        )
        result = tenant_config_from_dto(dto)
        assert result is not None
        assert result.super_admin_minimum_signatures == 3
        assert result.fee_limit_factor == 2.0

    def test_handles_missing_attributes(self) -> None:
        dto = SimpleNamespace()
        result = tenant_config_from_dto(dto)
        assert result is not None
        assert result.tenant_id == ""
        assert result.base_currency == ""
        assert result.super_admin_minimum_signatures == 0
        assert result.is_mfa_mandatory is False
