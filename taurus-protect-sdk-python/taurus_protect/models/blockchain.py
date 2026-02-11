"""Blockchain-related models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, Field


class Asset(BaseModel):
    """
    A cryptocurrency or token asset.

    Attributes:
        id: Unique asset identifier.
        name: Human-readable asset name.
        symbol: Asset symbol (e.g., "BTC", "ETH").
        blockchain: Blockchain name.
        network: Network name (e.g., "mainnet", "testnet").
        decimals: Number of decimal places.
        logo_url: URL to the asset logo.
        enabled: Whether the asset is enabled.
        is_token: Whether this is a token on another blockchain.
        contract_address: Token contract address (if token).
    """

    id: str = Field(description="Unique asset identifier")
    name: Optional[str] = Field(default=None, description="Asset name")
    symbol: Optional[str] = Field(default=None, description="Asset symbol")
    blockchain: Optional[str] = Field(default=None, description="Blockchain name")
    network: Optional[str] = Field(default=None, description="Network name")
    decimals: int = Field(default=0, description="Decimal places")
    logo_url: Optional[str] = Field(default=None, description="Logo URL")
    enabled: bool = Field(default=True, description="Whether enabled")
    is_token: bool = Field(default=False, description="Whether this is a token")
    contract_address: Optional[str] = Field(default=None, description="Token contract address")

    model_config = {"frozen": True}


class Blockchain(BaseModel):
    """
    A blockchain network.

    Attributes:
        id: Unique blockchain identifier.
        name: Blockchain name (e.g., "ETH", "BTC").
        network: Network name (e.g., "mainnet", "testnet").
        display_name: Human-readable blockchain name.
        enabled: Whether the blockchain is enabled.
        native_currency: Native currency symbol.
        block_height: Current block height (if requested).
        block_time: Average block time in seconds.
        confirmations_required: Required confirmations for transactions.
    """

    id: str = Field(description="Unique blockchain identifier")
    name: str = Field(description="Blockchain name")
    network: str = Field(default="", description="Network name")
    display_name: Optional[str] = Field(default=None, description="Display name")
    enabled: bool = Field(default=True, description="Whether enabled")
    native_currency: Optional[str] = Field(default=None, description="Native currency symbol")
    block_height: Optional[int] = Field(default=None, description="Current block height")
    block_time: Optional[int] = Field(default=None, description="Average block time (seconds)")
    confirmations_required: Optional[int] = Field(
        default=None, description="Required confirmations"
    )

    model_config = {"frozen": True}


class Exchange(BaseModel):
    """
    An exchange account.

    Attributes:
        id: Unique exchange account identifier.
        name: Exchange account name.
        exchange_label: Exchange label (e.g., "Binance", "Kraken").
        status: Exchange account status.
        currency_id: Currency identifier.
        currency: Currency symbol.
        balance: Current balance.
        pending_balance: Pending balance.
        enabled: Whether the exchange account is enabled.
        created_at: When the exchange account was created.
        updated_at: When the exchange account was last updated.
    """

    id: str = Field(description="Unique exchange account identifier")
    name: Optional[str] = Field(default=None, description="Exchange account name")
    exchange_label: Optional[str] = Field(default=None, description="Exchange label")
    status: Optional[str] = Field(default=None, description="Account status")
    currency_id: Optional[str] = Field(default=None, description="Currency ID")
    currency: Optional[str] = Field(default=None, description="Currency symbol")
    balance: Optional[str] = Field(default=None, description="Current balance")
    pending_balance: Optional[str] = Field(default=None, description="Pending balance")
    enabled: bool = Field(default=True, description="Whether enabled")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}


class FiatCurrency(BaseModel):
    """
    A fiat currency.

    Attributes:
        id: Currency identifier.
        code: Currency code (e.g., "USD", "EUR", "CHF").
        name: Human-readable currency name.
        symbol: Currency symbol (e.g., "$", "EUR").
        decimals: Number of decimal places.
        enabled: Whether the currency is enabled.
    """

    id: str = Field(description="Currency identifier")
    code: str = Field(description="Currency code")
    name: Optional[str] = Field(default=None, description="Currency name")
    symbol: Optional[str] = Field(default=None, description="Currency symbol")
    decimals: int = Field(default=2, description="Decimal places")
    enabled: bool = Field(default=True, description="Whether enabled")

    model_config = {"frozen": True}


class FiatProviderAccount(BaseModel):
    """
    A fiat provider account.

    Attributes:
        id: Unique account identifier.
        name: Account name.
        provider: Provider name.
        currency_code: Currency code.
        balance: Current balance.
        enabled: Whether the account is enabled.
    """

    id: str = Field(description="Unique account identifier")
    name: Optional[str] = Field(default=None, description="Account name")
    provider: Optional[str] = Field(default=None, description="Provider name")
    currency_code: Optional[str] = Field(default=None, description="Currency code")
    balance: Optional[str] = Field(default=None, description="Current balance")
    enabled: bool = Field(default=True, description="Whether enabled")

    model_config = {"frozen": True}


class ExchangeRate(BaseModel):
    """
    An exchange rate between two currencies.

    Attributes:
        from_currency: Source currency code.
        to_currency: Target currency code.
        rate: Exchange rate.
        timestamp: When the rate was fetched.
    """

    from_currency: str = Field(description="Source currency code")
    to_currency: str = Field(description="Target currency code")
    rate: str = Field(description="Exchange rate")
    timestamp: Optional[datetime] = Field(default=None, description="Rate timestamp")

    model_config = {"frozen": True}
