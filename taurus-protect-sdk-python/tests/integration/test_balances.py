"""Integration tests for BalanceService.

These tests verify balance operations against a live Taurus-PROTECT API.
"""

from __future__ import annotations

import pytest

from taurus_protect.client import ProtectClient


@pytest.mark.integration
def test_list_balances(client: ProtectClient) -> None:
    """Test listing balances."""
    balances, pagination = client.balances.list()

    print(f"Found {len(balances)} balances")
    if pagination:
        print(f"Total items: {pagination.total_items}, HasMore: {pagination.has_more}")

    for balance in balances:
        currency = balance.currency or ""
        total_confirmed = balance.balance or ""
        print(f"Balance: Currency={currency}, TotalConfirmed={total_confirmed}")

    assert isinstance(balances, list)
