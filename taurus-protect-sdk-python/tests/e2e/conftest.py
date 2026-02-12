"""E2E test configuration and fixtures for Taurus-PROTECT SDK.

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
        "markers", "e2e: mark test as E2E test (requires API access)"
    )


def pytest_collection_modifyitems(config: pytest.Config, items: list) -> None:
    """Skip E2E tests if not enabled."""
    tc = get_config()
    if tc.is_enabled():
        return

    skip_e2e = pytest.mark.skip(
        reason="E2E tests disabled. Set PROTECT_INTEGRATION_TEST=true or configure test.properties."
    )
    for item in items:
        if "e2e" in item.keywords:
            item.add_marker(skip_e2e)


@pytest.fixture
def client() -> Generator[ProtectClient, None, None]:
    """Create a ProtectClient for E2E testing."""
    skip_if_not_enabled()
    protect_client = get_test_client(1)
    try:
        yield protect_client
    finally:
        protect_client.close()
