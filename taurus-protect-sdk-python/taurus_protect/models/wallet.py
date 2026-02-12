"""Wallet models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, Field

from taurus_protect.models.balance import Balance
from taurus_protect.models.currency import Currency


class WalletAttribute(BaseModel):
    """Custom attribute on a wallet."""

    id: str = Field(description="Attribute identifier")
    key: str = Field(description="Attribute name")
    value: str = Field(description="Attribute value")
    content_type: Optional[str] = Field(default=None, description="Content type")
    owner: Optional[str] = Field(default=None, description="Attribute owner")
    type: Optional[str] = Field(default=None, description="Attribute type")
    subtype: Optional[str] = Field(default=None, description="Attribute subtype")
    is_file: bool = Field(default=False, description="Whether this is a file attribute")

    model_config = {"frozen": True}


class Wallet(BaseModel):
    """
    Cryptocurrency wallet.

    A wallet is a container for blockchain addresses of a specific currency.

    Attributes:
        id: Unique wallet identifier.
        name: Human-readable wallet name.
        currency: Currency symbol (e.g., "BTC", "ETH").
        blockchain: Blockchain name.
        network: Network name (e.g., "mainnet", "testnet").
        balance: Current wallet balance.
        is_omnibus: Whether this is an omnibus (pooled) wallet.
        disabled: Whether the wallet is disabled.
        comment: Optional description.
        customer_id: Optional customer identifier.
        external_wallet_id: Optional external identifier.
        visibility_group_id: Optional visibility group.
        account_path: HD wallet account derivation path.
        addresses_count: Number of addresses in the wallet.
        created_at: When the wallet was created.
        updated_at: When the wallet was last updated.
        attributes: Custom key-value attributes.
        currency_info: Detailed currency information.
    """

    id: str = Field(description="Unique wallet identifier")
    name: str = Field(description="Human-readable wallet name")
    currency: str = Field(description="Currency symbol")
    blockchain: str = Field(default="", description="Blockchain name")
    network: str = Field(default="", description="Network name")
    balance: Optional[Balance] = Field(default=None, description="Current wallet balance")
    is_omnibus: bool = Field(default=False, description="Whether this is an omnibus wallet")
    disabled: bool = Field(default=False, description="Whether the wallet is disabled")
    comment: Optional[str] = Field(default=None, description="Optional description")
    customer_id: Optional[str] = Field(default=None, description="Optional customer identifier")
    external_wallet_id: Optional[str] = Field(default=None, description="Optional external ID")
    visibility_group_id: Optional[str] = Field(default=None, description="Visibility group ID")
    account_path: Optional[str] = Field(default=None, description="HD account derivation path")
    addresses_count: int = Field(default=0, description="Number of addresses")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")
    attributes: List[WalletAttribute] = Field(default_factory=list, description="Custom attributes")
    currency_info: Optional[Currency] = Field(default=None, description="Detailed currency information")

    model_config = {"frozen": True}


class CreateWalletRequest(BaseModel):
    """
    Request to create a new wallet.

    Example:
        >>> request = CreateWalletRequest(
        ...     blockchain="ETH",
        ...     network="mainnet",
        ...     name="Trading Wallet",
        ...     is_omnibus=True,
        ...     comment="Primary trading account",
        ... )
    """

    blockchain: str = Field(description="Blockchain type (e.g., 'ETH', 'BTC', 'SOL')")
    network: str = Field(description="Network identifier (e.g., 'mainnet', 'testnet')")
    name: str = Field(description="Human-readable wallet name")
    is_omnibus: bool = Field(default=False, description="Whether this is an omnibus wallet")
    comment: Optional[str] = Field(default=None, description="Optional description")
    customer_id: Optional[str] = Field(default=None, description="Optional customer ID")


class ListWalletsOptions(BaseModel):
    """Options for listing wallets."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    currency: Optional[str] = Field(default=None, description="Filter by currency")
    query: Optional[str] = Field(default=None, description="Search query")
    exclude_disabled: bool = Field(default=False, description="Exclude disabled wallets")
