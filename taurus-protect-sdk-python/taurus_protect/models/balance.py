"""Balance model for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import Optional

from pydantic import BaseModel, Field


class Balance(BaseModel):
    """
    Balance information for a wallet or address.

    All balance values are returned as strings to preserve precision
    for cryptocurrency amounts.

    Attributes:
        total_confirmed: Total confirmed balance.
        total_unconfirmed: Total balance including unconfirmed.
        available_confirmed: Available confirmed balance.
        available_unconfirmed: Available balance including unconfirmed.
        reserved_confirmed: Reserved confirmed balance.
        reserved_unconfirmed: Reserved balance including unconfirmed.
    """

    total_confirmed: str = Field(default="0", description="Total confirmed balance")
    total_unconfirmed: str = Field(default="0", description="Total including unconfirmed")
    available_confirmed: str = Field(default="0", description="Available confirmed balance")
    available_unconfirmed: str = Field(default="0", description="Available including unconfirmed")
    reserved_confirmed: str = Field(default="0", description="Reserved confirmed balance")
    reserved_unconfirmed: str = Field(default="0", description="Reserved including unconfirmed")

    model_config = {"frozen": True}


class BalanceHistoryPoint(BaseModel):
    """
    A point in the balance history timeline.

    Represents a snapshot of balance at a specific point in time.
    """

    timestamp: Optional[datetime] = Field(default=None, description="Timestamp of the snapshot")
    total_confirmed: str = Field(default="0", description="Total confirmed balance at this time")
    total_unconfirmed: str = Field(
        default="0", description="Total including unconfirmed at this time"
    )
    available_confirmed: str = Field(default="0", description="Available confirmed at this time")
    available_unconfirmed: str = Field(
        default="0", description="Available unconfirmed at this time"
    )

    model_config = {"frozen": True}


class Asset(BaseModel):
    """
    Represents a cryptocurrency or token asset.

    Attributes:
        id: Unique asset identifier.
        symbol: Asset symbol (e.g., "ETH", "USDC").
        name: Human-readable asset name.
        decimals: Number of decimal places.
        blockchain: Blockchain the asset is on.
    """

    id: str = Field(default="", description="Unique asset identifier")
    symbol: str = Field(default="", description="Asset symbol")
    name: str = Field(default="", description="Human-readable name")
    decimals: int = Field(default=0, description="Number of decimal places")
    blockchain: str = Field(default="", description="Blockchain name")

    model_config = {"frozen": True}


class AssetBalance(BaseModel):
    """
    Represents the balance of a specific asset.

    Combines asset identification with current balance amounts.
    Used when retrieving balances across multiple assets for a wallet.
    """

    asset: Optional[Asset] = Field(default=None, description="The asset this balance represents")
    balance: Optional[Balance] = Field(default=None, description="Balance amounts for this asset")

    model_config = {"frozen": True}
