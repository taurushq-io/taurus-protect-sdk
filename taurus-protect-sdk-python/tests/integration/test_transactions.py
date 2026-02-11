"""Integration tests for TransactionService.

These tests require a live API connection. Enable by:
    export PROTECT_INTEGRATION_TEST=true
    export PROTECT_API_HOST="https://your-api-host.com"
    export PROTECT_API_KEY="your-api-key"
    export PROTECT_API_SECRET="your-hex-encoded-secret"
    pytest tests/integration/test_transactions.py -v

Or configure defaults in conftest.py.
"""

from __future__ import annotations

import logging
from typing import TYPE_CHECKING

import pytest

if TYPE_CHECKING:
    from taurus_protect.client import ProtectClient

logger = logging.getLogger(__name__)


@pytest.mark.integration
def test_list_transactions(client: ProtectClient) -> None:
    """Test listing transactions with pagination."""
    transactions, pagination = client.transactions.list(limit=10)

    logger.info(f"Found {len(transactions)} transactions")

    if pagination is not None:
        logger.info(f"Total items: {pagination.total_items}, HasMore: {pagination.has_more}")

    for tx in transactions:
        logger.info(f"Transaction: ID={tx.id}, Hash={tx.tx_hash}, Amount={tx.amount}")


@pytest.mark.integration
def test_list_transactions_by_currency(client: ProtectClient) -> None:
    """Test listing transactions filtered by currency."""
    transactions, _ = client.transactions.list(limit=10, currency="ETH")

    logger.info(f"Found {len(transactions)} ETH transactions")

    for tx in transactions:
        logger.info(f"Transaction: ID={tx.id}, Direction={tx.direction}, Amount={tx.amount}")


@pytest.mark.integration
def test_get_transaction_by_id(client: ProtectClient) -> None:
    """Test getting a transaction by ID."""
    # First list transactions to find a valid ID
    transactions, _ = client.transactions.list(limit=1)

    if len(transactions) == 0:
        pytest.skip("No transactions available for testing")

    tx_id = int(transactions[0].id)
    tx = client.transactions.get(tx_id)

    logger.info("Transaction details:")
    logger.info(f"  ID: {tx.id}")
    logger.info(f"  Hash: {tx.tx_hash}")
    logger.info(f"  Currency: {tx.currency}")
    logger.info(f"  Amount: {tx.amount}")

    assert tx is not None
    assert tx.id == str(tx_id)


@pytest.mark.integration
def test_get_transaction_by_hash(client: ProtectClient) -> None:
    """Test getting a transaction by hash."""
    # First list transactions to find a valid hash
    transactions, _ = client.transactions.list(limit=1)

    if len(transactions) == 0:
        pytest.skip("No transactions available for testing")

    tx_hash = transactions[0].tx_hash
    if not tx_hash:
        pytest.skip("No transaction hash available for testing")

    tx = client.transactions.get_by_hash(tx_hash)

    logger.info("Transaction by hash:")
    logger.info(f"  Hash: {tx.tx_hash}")
    logger.info(f"  Currency: {tx.currency}")

    assert tx is not None
    assert tx.tx_hash == tx_hash


@pytest.mark.integration
def test_list_transactions_by_address(client: ProtectClient) -> None:
    """Test listing transactions filtered by address."""
    # Get an address from the addresses service
    from taurus_protect.models.address import ListAddressesOptions

    options = ListAddressesOptions(limit=1)
    addresses, _ = client.addresses.list_with_options(options)

    if len(addresses) == 0:
        pytest.skip("No addresses available for testing")

    address = addresses[0].address
    if not address:
        pytest.skip("No address value available for testing")

    # Get transactions for that address
    transactions, _ = client.transactions.list_by_address(address, limit=10)

    logger.info(f"Found {len(transactions)} transactions for address {address}")

    for tx in transactions:
        logger.info(f"  Transaction: ID={tx.id}, Hash={tx.tx_hash}, Amount={tx.amount}")


@pytest.mark.integration
def test_export_transactions(client: ProtectClient) -> None:
    """Test exporting transactions to CSV."""
    from datetime import datetime, timedelta, timezone

    # Use time range of last 30 days
    to_date = datetime.now(timezone.utc)
    from_date = to_date - timedelta(days=30)

    csv_content = client.transactions.export_csv(
        from_date=from_date,
        to_date=to_date,
        limit=10,
    )

    logger.info(f"Exported CSV content length: {len(csv_content)} characters")
    if csv_content:
        # Log a snippet of the CSV (first 500 chars)
        snippet = csv_content[:500] + "..." if len(csv_content) > 500 else csv_content
        logger.info(f"CSV snippet:\n{snippet}")

    assert csv_content is not None


@pytest.mark.integration
def test_paginate_transactions(client: ProtectClient) -> None:
    """Test pagination through transactions."""
    page_size = 2
    max_items = 10  # Safety limit
    all_transactions = []
    offset = 0

    while len(all_transactions) < max_items:
        transactions, pagination = client.transactions.list(limit=page_size, offset=offset)

        if len(transactions) == 0:
            break

        all_transactions.extend(transactions)
        logger.info(
            f"Page at offset {offset}: got {len(transactions)} transactions, "
            f"total collected: {len(all_transactions)}"
        )

        # Check if there are more items
        if pagination is None or not pagination.has_more:
            break

        offset += page_size

    logger.info(f"Total transactions collected through pagination: {len(all_transactions)}")

    # Verify we got some transactions (if any exist)
    if len(all_transactions) > 0:
        # Check for duplicates but don't fail - offset-based pagination on live data
        # can return overlapping items when new data is inserted during the test
        ids = [tx.id for tx in all_transactions]
        unique_ids = set(ids)
        if len(ids) != len(unique_ids):
            logger.warning(
                f"Found {len(ids) - len(unique_ids)} duplicate transaction IDs "
                f"during pagination (expected with offset-based pagination on live data)"
            )
