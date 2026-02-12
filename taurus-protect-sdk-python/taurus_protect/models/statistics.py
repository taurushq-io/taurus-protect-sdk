"""Statistics, Price, and Score models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import Optional

from pydantic import BaseModel, Field


class Price(BaseModel):
    """
    Currency price information.

    Represents the exchange rate between two currencies with additional
    metadata about the price source and update time.

    Attributes:
        currency_from: Source currency symbol (e.g., "BTC").
        currency_to: Target currency symbol (e.g., "USD").
        rate: Exchange rate as a string to preserve precision.
        blockchain: Associated blockchain (if applicable).
        decimals: Number of decimal places for the rate.
        change_percent_24h: 24-hour price change percentage.
        source: Price data source.
        created_at: When the price was first recorded.
        updated_at: When the price was last updated.
    """

    currency_from: str = Field(default="", description="Source currency symbol")
    currency_to: str = Field(default="", description="Target currency symbol")
    rate: str = Field(default="0", description="Exchange rate")
    blockchain: Optional[str] = Field(default=None, description="Associated blockchain")
    decimals: Optional[str] = Field(default=None, description="Rate decimal places")
    change_percent_24h: Optional[str] = Field(
        default=None, description="24-hour price change percentage"
    )
    source: Optional[str] = Field(default=None, description="Price data source")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}


class PriceHistoryPoint(BaseModel):
    """
    A point in the price history timeline.

    Represents a historical price snapshot at a specific point in time.

    Attributes:
        timestamp: When this price was recorded.
        rate: Exchange rate at this time.
    """

    timestamp: Optional[datetime] = Field(default=None, description="Timestamp of the snapshot")
    rate: str = Field(default="0", description="Exchange rate at this time")

    model_config = {"frozen": True}


class PortfolioStatistics(BaseModel):
    """
    Portfolio-level statistics summary.

    Contains aggregated statistics about the entire portfolio including
    wallet counts, address counts, and total balance information.

    Attributes:
        addresses_count: Total number of addresses across all wallets.
        wallets_count: Total number of wallets.
        total_balance: Total balance in native currency units.
        total_balance_base_currency: Total balance converted to base (fiat) currency.
        avg_balance_per_address: Average balance per address.
    """

    addresses_count: str = Field(default="0", description="Total number of addresses")
    wallets_count: str = Field(default="0", description="Total number of wallets")
    total_balance: str = Field(default="0", description="Total balance in native units")
    total_balance_base_currency: str = Field(
        default="0", description="Total balance in base currency"
    )
    avg_balance_per_address: str = Field(default="0", description="Average balance per address")

    model_config = {"frozen": True}


class TransactionStatistics(BaseModel):
    """
    Transaction statistics for a time period.

    Contains aggregated statistics about transactions within a specified
    date range.

    Attributes:
        total_count: Total number of transactions.
        incoming_count: Number of incoming transactions.
        outgoing_count: Number of outgoing transactions.
        total_volume: Total transaction volume.
        incoming_volume: Total incoming transaction volume.
        outgoing_volume: Total outgoing transaction volume.
    """

    total_count: str = Field(default="0", description="Total transaction count")
    incoming_count: str = Field(default="0", description="Incoming transaction count")
    outgoing_count: str = Field(default="0", description="Outgoing transaction count")
    total_volume: str = Field(default="0", description="Total transaction volume")
    incoming_volume: str = Field(default="0", description="Incoming transaction volume")
    outgoing_volume: str = Field(default="0", description="Outgoing transaction volume")

    model_config = {"frozen": True}


class Score(BaseModel):
    """
    Risk score for an address or transaction.

    Represents a risk assessment score from a compliance/AML provider.

    Attributes:
        id: Unique score identifier.
        provider: Score provider name (e.g., "Chainalysis", "Elliptic").
        score_type: Type of score (e.g., "risk", "trust").
        score: The score value as a string.
        updated_at: When the score was last updated.
    """

    id: str = Field(default="", description="Unique score identifier")
    provider: str = Field(default="", description="Score provider name")
    score_type: str = Field(default="", description="Type of score")
    score: str = Field(default="0", description="Score value")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}
