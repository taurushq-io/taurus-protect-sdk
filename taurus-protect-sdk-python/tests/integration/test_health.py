"""Integration tests for HealthService.

These tests verify health check operations against a live Taurus-PROTECT API.
"""

from __future__ import annotations

import pytest

from taurus_protect.client import ProtectClient


@pytest.mark.integration
def test_health_check(client: ProtectClient) -> None:
    """Test getting all health checks."""
    result = client.health.get_all_health_checks()

    components = result.components or {}
    print(f"Health checks: {len(components)} components")

    for name, component in components.items():
        groups = component.groups or {}
        print(f"  {name}: {len(groups)} groups")

        for group_name, group in groups.items():
            health_checks = group.health_checks or []
            for hc in health_checks:
                print(f"    [{group_name}] {hc.name}: {hc.status}")
