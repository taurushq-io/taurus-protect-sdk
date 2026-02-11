"""Unit tests for ConfigService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.models.webhook import TenantConfig
from taurus_protect.services.config_service import ConfigService


class TestGet:
    """Tests for ConfigService.get()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        config_api = MagicMock()
        service = ConfigService(api_client=api_client, config_api=config_api)
        return service, config_api

    def test_get_returns_tenant_config(self) -> None:
        service, api = self._make_service()

        mock_config = TenantConfig()
        reply = MagicMock()
        reply.config = MagicMock()
        api.status_service_get_config_tenant.return_value = reply

        with patch(
            "taurus_protect.services.config_service.tenant_config_from_dto",
            return_value=mock_config,
        ):
            result = service.get()

        assert result is mock_config

    def test_get_returns_empty_config_when_none(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.config = None
        api.status_service_get_config_tenant.return_value = reply

        result = service.get()

        assert isinstance(result, TenantConfig)

    def test_get_returns_empty_config_when_mapper_returns_none(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.config = MagicMock()
        api.status_service_get_config_tenant.return_value = reply

        with patch(
            "taurus_protect.services.config_service.tenant_config_from_dto",
            return_value=None,
        ):
            result = service.get()

        assert isinstance(result, TenantConfig)


class TestGetFeatures:
    """Tests for ConfigService.get_features()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        config_api = MagicMock()
        service = ConfigService(api_client=api_client, config_api=config_api)
        return service, config_api

    def test_get_features_returns_enabled_features(self) -> None:
        service, api = self._make_service()

        config = TenantConfig(
            is_mfa_mandatory=True,
            is_protect_engine_cold=False,
            is_cold_protect_engine_offline=False,
            is_physical_air_gap_enabled=True,
            restrict_sources_for_whitelisted_addresses=False,
            exclude_container=False,
        )

        reply = MagicMock()
        reply.config = MagicMock()
        api.status_service_get_config_tenant.return_value = reply

        with patch(
            "taurus_protect.services.config_service.tenant_config_from_dto",
            return_value=config,
        ):
            features = service.get_features()

        feature_names = [f.name for f in features]
        assert "mfa_mandatory" in feature_names
        assert "physical_air_gap" in feature_names
        assert len(features) == 2

    def test_get_features_returns_empty_list_when_no_features_enabled(self) -> None:
        service, api = self._make_service()

        config = TenantConfig()

        reply = MagicMock()
        reply.config = MagicMock()
        api.status_service_get_config_tenant.return_value = reply

        with patch(
            "taurus_protect.services.config_service.tenant_config_from_dto",
            return_value=config,
        ):
            features = service.get_features()

        assert features == []
