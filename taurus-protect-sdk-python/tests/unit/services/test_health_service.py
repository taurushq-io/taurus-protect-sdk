"""Unit tests for HealthService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.errors import APIError
from taurus_protect.services.health_service import HealthService, HealthStatus


class TestCheck:
    """Tests for HealthService.check()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        health_api = MagicMock()
        service = HealthService(api_client=api_client, health_api=health_api)
        return service, health_api

    def test_check_returns_healthy_status(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.status = "healthy"
        reply.version = "1.0.0"
        reply.message = None
        api.health_service_health_check.return_value = reply

        result = service.check()

        assert isinstance(result, HealthStatus)
        assert result.status == "healthy"
        assert result.version == "1.0.0"

    def test_check_returns_unhealthy_on_api_error(self) -> None:
        service, api = self._make_service()

        error = APIError("connection refused", code=503)
        api.health_service_health_check.side_effect = error

        result = service.check()

        assert result.status == "unhealthy"

    def test_check_defaults_to_healthy_when_status_is_none(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.status = None
        reply.version = None
        reply.message = None
        api.health_service_health_check.return_value = reply

        result = service.check()
        assert result.status == "healthy"


class TestGetAllHealthChecks:
    """Tests for HealthService.get_all_health_checks()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        health_api = MagicMock()
        service = HealthService(api_client=api_client, health_api=health_api)
        return service, health_api

    def test_get_all_health_checks_returns_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.components = None
        api.health_service_get_all_health_checks.return_value = reply

        result = service.get_all_health_checks()

        assert result is not None
        assert result.components is None

    def test_get_all_health_checks_with_tenant_id(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.components = None
        api.health_service_get_all_health_checks.return_value = reply

        service.get_all_health_checks(tenant_id="tenant-1")

        call_kwargs = api.health_service_get_all_health_checks.call_args
        assert call_kwargs[1].get("tenant_id") == "tenant-1" or \
               call_kwargs.kwargs.get("tenant_id") == "tenant-1"

    def test_get_all_health_checks_maps_components(self) -> None:
        service, api = self._make_service()

        # Build nested mock structure
        hc_dto = MagicMock()
        hc_dto.tenant_id = "t1"
        hc_dto.component_name = "api"
        hc_dto.component_id = "api-1"
        hc_dto.group = "core"
        hc_dto.health_check = "db"
        hc_dto.status = "HEALTHY"
        hc_dto.report = None
        hc_dto.last_update_date = None
        hc_dto.valid_until_date = None

        group_dto = MagicMock()
        group_dto.health_checks = [hc_dto]

        component_dto = MagicMock()
        component_dto.groups = {"core": group_dto}

        reply = MagicMock()
        reply.components = {"api": component_dto}
        api.health_service_get_all_health_checks.return_value = reply

        result = service.get_all_health_checks()

        assert result.components is not None
        assert "api" in result.components
