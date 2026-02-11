"""Health service for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING, Any, Dict, List, Optional

from pydantic import BaseModel, Field

from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class HealthStatus(BaseModel):
    """Health check status response."""

    status: str = Field(default="unknown", description="Health status (healthy, unhealthy)")
    version: Optional[str] = Field(default=None, description="API version")
    message: Optional[str] = Field(default=None, description="Status message")

    model_config = {"frozen": True}


class ClientStatus(BaseModel):
    """Status of a client connection."""

    host_port: Optional[str] = Field(default=None, description="Host:port of the client")
    connected: bool = Field(default=False, description="Whether the client is connected")
    error: Optional[str] = Field(default=None, description="Error message if connection failed")
    last_ping: Optional[datetime] = Field(
        default=None, description="When the client was last pinged"
    )

    model_config = {"frozen": True}


class HealthReport(BaseModel):
    """Detailed information about a health check."""

    name: Optional[str] = Field(default=None, description="Name of the health check report")
    status: Optional[str] = Field(default=None, description="Status of the report")
    message: Optional[str] = Field(default=None, description="Additional status message")
    duration: Optional[str] = Field(default=None, description="How long the check took")
    error: Optional[str] = Field(default=None, description="Error message if check failed")
    results: Optional[Dict[str, str]] = Field(default=None, description="Additional result data")
    vaultd_clients: Optional[List[ClientStatus]] = Field(
        default=None, description="Status of vault clients"
    )

    model_config = {"frozen": True}


class HealthCheck(BaseModel):
    """A single health check entry."""

    tenant_id: Optional[str] = Field(default=None, description="Tenant ID for the health check")
    component_name: Optional[str] = Field(default=None, description="Name of the component")
    component_id: Optional[str] = Field(
        default=None, description="Unique identifier for the component"
    )
    group: Optional[str] = Field(default=None, description="Group this health check belongs to")
    name: Optional[str] = Field(default=None, description="Name of the health check")
    status: Optional[str] = Field(default=None, description="Health status (HEALTHY, UNHEALTHY)")
    report: Optional[HealthReport] = Field(default=None, description="Detailed health check report")
    last_update_date: Optional[datetime] = Field(
        default=None, description="When the check was last updated"
    )
    valid_until_date: Optional[datetime] = Field(
        default=None, description="When the check result expires"
    )

    model_config = {"frozen": True}


class HealthGroup(BaseModel):
    """A group of health checks."""

    health_checks: Optional[List[HealthCheck]] = Field(
        default=None, description="Health checks in this group"
    )

    model_config = {"frozen": True}


class HealthComponent(BaseModel):
    """A component with health check groups."""

    groups: Optional[Dict[str, HealthGroup]] = Field(
        default=None, description="Map of group name to health group"
    )

    model_config = {"frozen": True}


class GetAllHealthChecksResult(BaseModel):
    """Result of getting all health checks."""

    components: Optional[Dict[str, HealthComponent]] = Field(
        default=None, description="Map of component name to health component"
    )

    model_config = {"frozen": True}


class HealthService(BaseService):
    """
    Service for checking API health.

    Provides a simple health check endpoint to verify API connectivity.

    Example:
        >>> health = client.health.check()
        >>> print(f"Status: {health.status}")
    """

    def __init__(self, api_client: Any, health_api: Any) -> None:
        """
        Initialize health service.

        Args:
            api_client: The OpenAPI client instance.
            health_api: The HealthAPI service from OpenAPI client.
        """
        super().__init__(api_client)
        self._health_api = health_api

    def check(self) -> HealthStatus:
        """
        Check the API health status.

        Returns:
            Health status response.

        Raises:
            APIError: If health check fails.
        """
        try:
            resp = self._health_api.health_service_health_check()

            status = getattr(resp, "status", "healthy")
            version = getattr(resp, "version", None)
            message = getattr(resp, "message", None)

            return HealthStatus(
                status=status or "healthy",
                version=version,
                message=message,
            )
        except Exception as e:
            # If health check fails, return unhealthy status
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                return HealthStatus(status="unhealthy", message=str(e))
            raise self._handle_error(e) from e

    def get_all_health_checks(
        self,
        tenant_id: Optional[str] = None,
        fail_if_unhealthy: bool = False,
    ) -> GetAllHealthChecksResult:
        """
        Get all health checks with optional filtering.

        Args:
            tenant_id: Optional filter by tenant ID.
            fail_if_unhealthy: If True, raise an error if any checks are unhealthy.

        Returns:
            Result containing all health check components.

        Raises:
            APIError: If API request fails or if fail_if_unhealthy=True and checks are unhealthy.
        """
        try:
            kwargs: Dict[str, Any] = {}
            if tenant_id:
                kwargs["tenant_id"] = tenant_id
            if fail_if_unhealthy:
                kwargs["fail_if_unhealthy"] = fail_if_unhealthy

            resp = self._health_api.health_service_get_all_health_checks(**kwargs)

            components_dto = getattr(resp, "components", None)
            components = self._map_components(components_dto)

            return GetAllHealthChecksResult(components=components)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def _map_components(
        self, components_dto: Optional[Dict[str, Any]]
    ) -> Optional[Dict[str, HealthComponent]]:
        """Map OpenAPI components DTO to domain model."""
        if components_dto is None:
            return None

        result: Dict[str, HealthComponent] = {}
        for name, component_dto in components_dto.items():
            result[name] = self._map_component(component_dto)
        return result

    def _map_component(self, component_dto: Any) -> HealthComponent:
        """Map OpenAPI component DTO to domain model."""
        groups_dto = getattr(component_dto, "groups", None)
        groups = self._map_groups(groups_dto)
        return HealthComponent(groups=groups)

    def _map_groups(self, groups_dto: Optional[Dict[str, Any]]) -> Optional[Dict[str, HealthGroup]]:
        """Map OpenAPI groups DTO to domain model."""
        if groups_dto is None:
            return None

        result: Dict[str, HealthGroup] = {}
        for name, group_dto in groups_dto.items():
            result[name] = self._map_group(group_dto)
        return result

    def _map_group(self, group_dto: Any) -> HealthGroup:
        """Map OpenAPI group DTO to domain model."""
        health_checks_dto = getattr(group_dto, "health_checks", None)
        health_checks = self._map_health_checks(health_checks_dto)
        return HealthGroup(health_checks=health_checks)

    def _map_health_checks(
        self, health_checks_dto: Optional[List[Any]]
    ) -> Optional[List[HealthCheck]]:
        """Map OpenAPI health checks DTO to domain model."""
        if health_checks_dto is None:
            return None

        return [self._map_health_check(hc) for hc in health_checks_dto]

    def _map_health_check(self, hc_dto: Any) -> HealthCheck:
        """Map OpenAPI health check DTO to domain model."""
        report_dto = getattr(hc_dto, "report", None)
        report = self._map_report(report_dto) if report_dto else None

        return HealthCheck(
            tenant_id=getattr(hc_dto, "tenant_id", None),
            component_name=getattr(hc_dto, "component_name", None),
            component_id=getattr(hc_dto, "component_id", None),
            group=getattr(hc_dto, "group", None),
            name=getattr(hc_dto, "health_check", None),
            status=getattr(hc_dto, "status", None),
            report=report,
            last_update_date=getattr(hc_dto, "last_update_date", None),
            valid_until_date=getattr(hc_dto, "valid_until_date", None),
        )

    def _map_report(self, report_dto: Any) -> HealthReport:
        """Map OpenAPI health report DTO to domain model."""
        vaultd_clients_dto = getattr(report_dto, "vaultd_clients", None)
        vaultd_clients = None
        if vaultd_clients_dto:
            vaultd_clients = [
                ClientStatus(
                    host_port=getattr(c, "host_port", None),
                    connected=getattr(c, "connected", False),
                    error=getattr(c, "error", None),
                    last_ping=getattr(c, "last_ping", None),
                )
                for c in vaultd_clients_dto
            ]

        return HealthReport(
            name=getattr(report_dto, "name", None),
            status=getattr(report_dto, "status", None),
            message=getattr(report_dto, "message", None),
            duration=getattr(report_dto, "duration", None),
            error=getattr(report_dto, "error", None),
            results=getattr(report_dto, "results", None),
            vaultd_clients=vaultd_clients,
        )
