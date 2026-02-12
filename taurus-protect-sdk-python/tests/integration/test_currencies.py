"""Integration tests for CurrencyService.

These tests verify currency operations against a live Taurus-PROTECT API.
"""

from __future__ import annotations

import pytest

from taurus_protect.client import ProtectClient


@pytest.mark.integration
def test_list_currencies(client: ProtectClient) -> None:
    """Test listing currencies."""
    currencies = client.currencies.list()

    print(f"Found {len(currencies)} currencies")
    for currency in currencies:
        print(
            f"Currency: ID={currency.id}, Symbol={currency.symbol}, "
            f"Name={currency.name}, Blockchain={currency.blockchain}"
        )

    assert isinstance(currencies, list)


@pytest.mark.integration
def test_get_currency(client: ProtectClient) -> None:
    """Test getting a single currency by ID."""
    # First, list currencies to get a valid ID
    currencies = client.currencies.list()

    if len(currencies) == 0:
        pytest.skip("No currencies available for testing")

    currency_id = currencies[0].id
    currency = client.currencies.get(currency_id)

    print("Currency details:")
    print(f"  ID: {currency.id}")
    print(f"  Symbol: {currency.symbol}")
    print(f"  Name: {currency.name}")
    print(f"  Blockchain: {currency.blockchain}")

    assert currency.id == currency_id
