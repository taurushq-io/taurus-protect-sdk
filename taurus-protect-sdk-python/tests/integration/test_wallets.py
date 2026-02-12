"""Integration tests for WalletService.

These tests verify wallet operations against a live Taurus-PROTECT API.
"""

from __future__ import annotations

import pytest

from taurus_protect.client import ProtectClient


@pytest.mark.integration
def test_list_wallets(client: ProtectClient) -> None:
    """Test listing wallets with pagination."""
    wallets, pagination = client.wallets.list(limit=10)

    print(f"Found {len(wallets)} wallets")
    if pagination:
        print(f"Total items: {pagination.total_items}, HasMore: {pagination.has_more}")

    for wallet in wallets:
        print(f"Wallet: ID={wallet.id}, Name={wallet.name}, Currency={wallet.currency}")


@pytest.mark.integration
def test_get_wallet(client: ProtectClient) -> None:
    """Test getting a single wallet by ID."""
    # First, get a list to find a valid wallet ID
    wallets, _ = client.wallets.list(limit=1)

    if len(wallets) == 0:
        pytest.skip("No wallets available for testing")

    wallet_id = int(wallets[0].id)
    wallet = client.wallets.get(wallet_id)

    print("Wallet details:")
    print(f"  ID: {wallet.id}")
    print(f"  Name: {wallet.name}")
    print(f"  Currency: {wallet.currency}")
    print(f"  Blockchain: {wallet.blockchain}")
    print(f"  AddressesCount: {wallet.addresses_count}")


@pytest.mark.integration
def test_pagination(client: ProtectClient) -> None:
    """Test pagination through all wallets with safety limit."""
    page_size = 2
    all_wallets = []
    offset = 0

    # Fetch all wallets using pagination
    while True:
        wallets, pagination = client.wallets.list(limit=page_size, offset=offset)

        all_wallets.extend(wallets)
        print(f"Fetched {len(wallets)} wallets (offset={offset})")
        for wallet in wallets:
            print(f"  {wallet}")

        if pagination is None or not pagination.has_more:
            break
        offset += page_size

        # Safety limit for tests
        if offset > 100:
            print("Stopping pagination test at 100 items")
            break

    print(f"Total wallets fetched via pagination: {len(all_wallets)}")
