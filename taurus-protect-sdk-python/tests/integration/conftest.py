"""Integration test configuration and fixtures for Taurus-PROTECT SDK.

Uses shared test utilities from tests.testutil for configuration.
"""

from __future__ import annotations

from typing import Generator

import pytest

from taurus_protect.client import ProtectClient

from tests.testutil import get_config, get_test_client, skip_if_not_enabled


def pytest_configure(config: pytest.Config) -> None:
    """Register custom markers."""
    config.addinivalue_line(
        "markers", "integration: mark test as integration test (requires API access)"
    )


def pytest_collection_modifyitems(config: pytest.Config, items: list) -> None:
    """Skip integration tests if not enabled."""
    tc = get_config()
    if tc.is_enabled():
        return

    skip_integration = pytest.mark.skip(
        reason="Integration tests disabled. Set PROTECT_INTEGRATION_TEST=true or configure test.properties."
    )
    for item in items:
        if "integration" in item.keywords:
            item.add_marker(skip_integration)


@pytest.fixture
def client() -> Generator[ProtectClient, None, None]:
    """Create a ProtectClient for integration testing."""
    skip_if_not_enabled()
    protect_client = get_test_client(1)
    try:
        yield protect_client
    finally:
        protect_client.close()
