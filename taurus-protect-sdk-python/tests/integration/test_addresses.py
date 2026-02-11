"""Integration tests for AddressService."""

from __future__ import annotations

import logging

import pytest

from taurus_protect.client import ProtectClient
from taurus_protect.models.address import ListAddressesOptions

logger = logging.getLogger(__name__)


@pytest.mark.integration
def test_list_addresses(client: ProtectClient) -> None:
    """Test listing addresses with pagination."""
    options = ListAddressesOptions(limit=10)
    addresses, pagination = client.addresses.list_with_options(options)

    logger.info("Found %d addresses", len(addresses))
    if pagination is not None:
        logger.info("Total items: %d, HasMore: %s", pagination.total_items, pagination.has_more)

    for addr in addresses:
        logger.info("Address: ID=%s, Label=%s, Currency=%s", addr.id, addr.label, addr.currency)


@pytest.mark.integration
def test_get_address(client: ProtectClient) -> None:
    """Test getting a single address by ID."""
    # First, get a list to find a valid address ID
    options = ListAddressesOptions(limit=1)
    addresses, _ = client.addresses.list_with_options(options)

    if len(addresses) == 0:
        pytest.skip("No addresses available for testing")

    address_id = int(addresses[0].id)
    address = client.addresses.get(address_id)

    logger.info("Address details:")
    logger.info("  ID: %s", address.id)
    logger.info("  Address: %s", address.address)
    logger.info("  Label: %s", address.label)
    logger.info("  Currency: %s", address.currency)
    logger.info("  WalletID: %s", address.wallet_id)
    if address.balance is not None:
        logger.info("  Balance (confirmed): %s", address.balance.total_confirmed)
