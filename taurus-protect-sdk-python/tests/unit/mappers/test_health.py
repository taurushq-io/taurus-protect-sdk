"""Unit tests for health service mapper functions."""

from types import SimpleNamespace

from taurus_protect.services.health_service import (
    ClientStatus,
    HealthCheck,
    HealthComponent,
    HealthGroup,
    HealthReport,
    HealthService,
    HealthStatus,
)


class TestHealthStatusModel:
    """Tests for HealthStatus model."""

    def test_creates_with_defaults(self) -> None:
        status = HealthStatus()
        assert status.status == "unknown"
        assert status.version is None
        assert status.message is None

    def test_creates_with_all_fields(self) -> None:
        status = HealthStatus(
            status="healthy",
            version="1.2.3",
            message="All systems operational",
        )
        assert status.status == "healthy"
        assert status.version == "1.2.3"
        assert status.message == "All systems operational"


class TestHealthCheckModel:
    """Tests for HealthCheck model."""

    def test_creates_with_all_fields(self) -> None:
        report = HealthReport(
            name="db-check",
            status="HEALTHY",
            message="OK",
            duration="12ms",
        )
        hc = HealthCheck(
            tenant_id="t-1",
            component_name="database",
            component_id="db-1",
            group="core",
            name="postgres-check",
            status="HEALTHY",
            report=report,
        )
        assert hc.tenant_id == "t-1"
        assert hc.component_name == "database"
        assert hc.name == "postgres-check"
        assert hc.status == "HEALTHY"
        assert hc.report is not None
        assert hc.report.name == "db-check"

    def test_creates_with_defaults(self) -> None:
        hc = HealthCheck()
        assert hc.tenant_id is None
        assert hc.status is None
        assert hc.report is None


class TestHealthReportModel:
    """Tests for HealthReport model."""

    def test_maps_vaultd_clients(self) -> None:
        report = HealthReport(
            name="vault-check",
            status="HEALTHY",
            vaultd_clients=[
                ClientStatus(
                    host_port="localhost:8200",
                    connected=True,
                ),
                ClientStatus(
                    host_port="localhost:8201",
                    connected=False,
                    error="connection refused",
                ),
            ],
        )
        assert report.name == "vault-check"
        assert report.vaultd_clients is not None
        assert len(report.vaultd_clients) == 2
        assert report.vaultd_clients[0].connected is True
        assert report.vaultd_clients[1].connected is False
        assert report.vaultd_clients[1].error == "connection refused"


class TestHealthServiceMappers:
    """Tests for HealthService internal mapping methods."""

    def _make_service(self) -> HealthService:
        """Create a HealthService with a mock API."""
        return HealthService.__new__(HealthService)

    def test_map_health_check(self) -> None:
        svc = self._make_service()
        dto = SimpleNamespace(
            tenant_id="t-1",
            component_name="api",
            component_id="api-1",
            group="web",
            health_check="endpoint-check",
            status="HEALTHY",
            report=None,
            last_update_date=None,
            valid_until_date=None,
        )
        result = svc._map_health_check(dto)
        assert result.tenant_id == "t-1"
        assert result.name == "endpoint-check"
        assert result.status == "HEALTHY"

    def test_map_health_check_with_report(self) -> None:
        svc = self._make_service()
        report_dto = SimpleNamespace(
            name="db-check",
            status="OK",
            message="Connected",
            duration="5ms",
            error=None,
            results={"latency": "2ms"},
            vaultd_clients=None,
        )
        dto = SimpleNamespace(
            tenant_id="t-1",
            component_name="db",
            component_id="db-1",
            group="storage",
            health_check="pg-check",
            status="HEALTHY",
            report=report_dto,
            last_update_date=None,
            valid_until_date=None,
        )
        result = svc._map_health_check(dto)
        assert result.report is not None
        assert result.report.name == "db-check"
        assert result.report.results == {"latency": "2ms"}

    def test_map_components_none(self) -> None:
        svc = self._make_service()
        assert svc._map_components(None) is None

    def test_map_groups_none(self) -> None:
        svc = self._make_service()
        assert svc._map_groups(None) is None

    def test_map_health_checks_none(self) -> None:
        svc = self._make_service()
        assert svc._map_health_checks(None) is None
