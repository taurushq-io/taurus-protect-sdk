"""Integration tests for miscellaneous operations.

These tests verify tags, statistics, and client lifecycle operations
against a live Taurus-PROTECT API.
"""

from __future__ import annotations

import pytest

from taurus_protect.client import ProtectClient

from tests.testutil import get_test_client, skip_if_not_enabled


@pytest.mark.integration
def test_list_tags(client: ProtectClient) -> None:
    """Test listing tags."""
    tags = client.tags.list()

    print(f"Found {len(tags)} tags")
    for tag in tags:
        print(f"Tag: ID={tag.id}, Name={tag.name}")


@pytest.mark.integration
def test_get_portfolio_statistics(client: ProtectClient) -> None:
    """Test getting portfolio statistics."""
    stats = client.statistics.get_summary()

    print("Portfolio statistics:")
    if stats:
        print(f"  TotalBalance: {stats.total_balance}")
        print(f"  TotalBalanceBaseCurrency: {stats.total_balance_base_currency}")
        print(f"  WalletsCount: {stats.wallets_count}")
        print(f"  AddressesCount: {stats.addresses_count}")
    else:
        print("  No statistics available")


@pytest.mark.integration
def test_client_lifecycle() -> None:
    """Test client service lazy initialization and close() idempotency."""
    skip_if_not_enabled()
    client = get_test_client(1)

    # Test that services are lazily initialized
    _ = client.wallets
    _ = client.addresses
    _ = client.requests
    _ = client.transactions

    print("Services lazily initialized successfully")

    # Test that close() works
    client.close()
    assert client.is_closed, "Client should be closed after close()"
    print("First close() succeeded")

    # Test that close() is idempotent
    client.close()
    assert client.is_closed, "Client should still be closed after second close()"
    print("Second close() succeeded (idempotent)")
