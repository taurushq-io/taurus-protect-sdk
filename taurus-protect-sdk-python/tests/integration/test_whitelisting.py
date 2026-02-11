"""Integration tests for WhitelistedAddressService and WhitelistedAssetService.

These tests verify whitelisting operations against a live Taurus-PROTECT API.
"""

from __future__ import annotations

import pytest

from taurus_protect.client import ProtectClient


@pytest.mark.integration
def test_paginate_all_whitelisted_addresses(
    client: ProtectClient,
) -> None:
    """Test paginating through all whitelisted addresses."""
    page_size = 50
    all_addresses = []
    offset = 0

    # Fetch all whitelisted addresses using pagination
    while True:
        addresses, pagination = client.whitelisted_addresses.list(
            limit=page_size,
            offset=offset,
        )

        all_addresses.extend(addresses)
        print(f"Fetched {len(addresses)} whitelisted addresses (offset={offset})")

        for addr in addresses:
            print(f"  Blockchain: {addr.currency}, Network: {addr.network}, Address: {addr.address}")

        if pagination is None or not pagination.has_more:
            break
        offset += page_size

        # Safety limit for tests
        if offset > 2000:
            print("Stopping pagination test at 2000 items")
            break

    print(f"Total whitelisted addresses fetched via pagination: {len(all_addresses)}")


@pytest.mark.integration
def test_list_whitelisted_addresses(client: ProtectClient) -> None:
    """Test listing whitelisted addresses (single page)."""
    addresses, pagination = client.whitelisted_addresses.list(
        limit=10,
        offset=0,
    )

    print(f"Found {len(addresses)} whitelisted addresses")
    for addr in addresses:
        print(f"  {addr.currency}/{addr.network}: {addr.address}")

    assert addresses is not None


@pytest.mark.integration
def test_paginate_all_whitelisted_assets(
    client: ProtectClient,
) -> None:
    """Test paginating through all whitelisted assets."""
    page_size = 10
    all_assets = []
    offset = 0

    # Fetch all whitelisted assets using pagination
    while True:
        assets, pagination = client.whitelisted_assets.list(
            limit=page_size,
            offset=offset,
        )

        all_assets.extend(assets)
        print(f"Fetched {len(assets)} whitelisted assets (offset={offset})")

        for asset in assets:
            print(f"  WhitelistedAsset:")
            print(f"    ID: {asset.id}")
            print(f"    Name: {asset.name}")
            print(f"    Symbol: {asset.symbol}")
            print(f"    Blockchain: {asset.blockchain}")
            print(f"    Network: {asset.network}")
            print(f"    Contract Address: {asset.contract_address}")
            print(f"    Status: {asset.status}")
            print(f"    Action: {asset.action}")
            print(f"    Rule: {asset.rule}")
            print(f"    Created At: {asset.created_at}")
            print(f"    Tenant ID: {asset.tenant_id}")
            print(f"    Business Rule Enabled: {asset.business_rule_enabled}")
            print()

        if pagination is None or not pagination.has_more:
            break
        offset += page_size

        # Safety limit for tests
        if offset > 5000:
            print("Stopping pagination test at 5000 items")
            break

    print(f"Total whitelisted assets fetched via pagination: {len(all_assets)}")


@pytest.mark.integration
def test_list_whitelisted_assets(client: ProtectClient) -> None:
    """Test listing whitelisted assets (single page)."""
    assets, pagination = client.whitelisted_assets.list(
        limit=10,
        offset=0,
    )

    print(f"Found {len(assets)} whitelisted assets")
    for asset in assets:
        print(f"  ID={asset.id}, Blockchain={asset.blockchain}, Symbol={asset.symbol}")

    assert assets is not None


@pytest.mark.integration
def test_get_whitelisted_address(client: ProtectClient) -> None:
    """Test getting a single whitelisted address by ID."""
    # First list addresses to find a valid ID
    addresses, _ = client.whitelisted_addresses.list(limit=1, offset=0)

    if not addresses:
        pytest.skip("No whitelisted addresses available for testing")

    address_id = int(addresses[0].id)
    address = client.whitelisted_addresses.get(address_id)

    print("Whitelisted address:")
    print(f"  ID: {address.id}")
    print(f"  Blockchain: {address.currency}")
    print(f"  Network: {address.network}")
    print(f"  Address: {address.address}")
    print(f"  Label: {address.label}")

    assert address is not None
    assert address.id == str(address_id)


@pytest.mark.integration
def test_list_whitelisted_addresses_by_blockchain(
    client: ProtectClient,
) -> None:
    """Test listing whitelisted addresses filtered by blockchain (ETH)."""
    eth_addresses, pagination = client.whitelisted_addresses.list(
        currency="ETH",
        limit=10,
        offset=0,
    )

    print(f"Found {len(eth_addresses)} ETH whitelisted addresses")
    for addr in eth_addresses:
        print(f"  Address: {addr.address}")
        print(f"  Blockchain: {addr.currency}")
        print(f"  Network: {addr.network}")

    assert eth_addresses is not None


@pytest.mark.integration
def test_list_whitelisted_addresses_by_blockchain_and_network(
    client: ProtectClient,
) -> None:
    """Test listing whitelisted addresses filtered by blockchain and network."""
    # Note: The Python SDK uses 'currency' parameter for blockchain filtering
    # and the list API may not support network filtering directly.
    # First get ETH addresses, then filter by network if the API supports it
    eth_addresses, pagination = client.whitelisted_addresses.list(
        currency="ETH",
        limit=50,
        offset=0,
    )

    # Filter by network client-side if API doesn't support network param
    mainnet_addresses = [addr for addr in eth_addresses if addr.network == "mainnet"]

    print(f"Found {len(mainnet_addresses)} ETH mainnet whitelisted addresses")
    for addr in mainnet_addresses:
        print(f"  Address: {addr.address}")
        print(f"  Blockchain: {addr.currency}")
        print(f"  Network: {addr.network}")

    assert mainnet_addresses is not None or eth_addresses is not None


@pytest.mark.integration
def test_get_whitelisted_asset(client: ProtectClient) -> None:
    """Test getting a single whitelisted asset by ID."""
    # First list assets to find a valid ID
    assets, _ = client.whitelisted_assets.list(limit=1, offset=0)

    if not assets:
        pytest.skip("No whitelisted assets available for testing")

    asset_id = int(assets[0].id)
    asset = client.whitelisted_assets.get(asset_id)

    print("Whitelisted asset:")
    print(f"  ID: {asset.id}")
    print(f"  Blockchain: {asset.blockchain}")
    print(f"  Network: {asset.network}")
    print(f"  Symbol: {asset.symbol}")
    print(f"  Name: {asset.name}")
    print(f"  Contract Address: {asset.contract_address}")
    print(f"  Status: {asset.status}")

    assert asset is not None
    assert asset.id == str(asset_id)


@pytest.mark.integration
def test_list_whitelisted_assets_by_blockchain(
    client: ProtectClient,
) -> None:
    """Test listing whitelisted assets filtered by blockchain (ETH)."""
    eth_assets, pagination = client.whitelisted_assets.list(
        blockchain="ETH",
        limit=10,
        offset=0,
    )

    print(f"Found {len(eth_assets)} ETH whitelisted assets")
    for asset in eth_assets:
        print(f"  Symbol: {asset.symbol}")
        print(f"  Blockchain: {asset.blockchain}")
        print(f"  Contract: {asset.contract_address}")

    assert eth_assets is not None


@pytest.mark.integration
def test_list_whitelisted_assets_by_blockchain_and_network(
    client: ProtectClient,
) -> None:
    """Test listing whitelisted assets filtered by blockchain and network."""
    mainnet_assets, pagination = client.whitelisted_assets.list(
        blockchain="ETH",
        network="mainnet",
        limit=10,
        offset=0,
    )

    print(f"Found {len(mainnet_assets)} ETH mainnet whitelisted assets")
    for asset in mainnet_assets:
        print(f"  Symbol: {asset.symbol}")
        print(f"  Network: {asset.network}")
        print(f"  Contract: {asset.contract_address}")

    assert mainnet_assets is not None
