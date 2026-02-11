"""Currency models for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import Optional

from pydantic import BaseModel, Field


class Currency(BaseModel):
    """
    Cryptocurrency or fiat currency.

    Attributes:
        id: Unique currency identifier.
        name: Human-readable currency name.
        symbol: Currency symbol (e.g., "ETH", "BTC").
        blockchain: Blockchain name.
        network: Network name (e.g., "mainnet", "testnet").
        decimals: Number of decimal places.
        logo_url: URL to the currency logo.
        enabled: Whether the currency is enabled.
        is_token: Whether this is a token on another blockchain.
        contract_address: Token contract address (if token).
        display_name: Display name.
        type: Currency type.
        coin_type_index: BIP-44 coin type index.
        token_id: Token ID.
        wlca_id: Whitelisted contract address ID.
        is_erc20: Whether this is an ERC-20 token.
        is_fa12: Whether this is an FA1.2 token.
        is_fa20: Whether this is an FA2.0 token.
        is_nft: Whether this is an NFT.
        is_utxo_based: Whether this is UTXO-based.
        is_account_based: Whether this is account-based.
        is_fiat: Whether this is a fiat currency.
        has_staking: Whether staking is supported.
    """

    id: str = Field(description="Unique currency identifier")
    name: Optional[str] = Field(default=None, description="Currency name")
    symbol: Optional[str] = Field(default=None, description="Currency symbol")
    blockchain: Optional[str] = Field(default=None, description="Blockchain name")
    network: Optional[str] = Field(default=None, description="Network name")
    decimals: int = Field(default=0, description="Decimal places")
    logo_url: Optional[str] = Field(default=None, description="Logo URL")
    enabled: bool = Field(default=True, description="Whether enabled")
    is_token: bool = Field(default=False, description="Whether this is a token")
    contract_address: Optional[str] = Field(default=None, description="Token contract address")
    display_name: Optional[str] = Field(default=None, description="Display name")
    type: Optional[str] = Field(default=None, description="Currency type")
    coin_type_index: Optional[str] = Field(default=None, description="BIP-44 coin type index")
    token_id: Optional[str] = Field(default=None, description="Token ID")
    wlca_id: Optional[int] = Field(default=None, description="Whitelisted contract address ID")
    is_erc20: bool = Field(default=False, description="Whether this is an ERC-20 token")
    is_fa12: bool = Field(default=False, description="Whether this is an FA1.2 token")
    is_fa20: bool = Field(default=False, description="Whether this is an FA2.0 token")
    is_nft: bool = Field(default=False, description="Whether this is an NFT")
    is_utxo_based: bool = Field(default=False, description="Whether this is UTXO-based")
    is_account_based: bool = Field(default=False, description="Whether this is account-based")
    is_fiat: bool = Field(default=False, description="Whether this is a fiat currency")
    has_staking: bool = Field(default=False, description="Whether staking is supported")

    model_config = {"frozen": True}


class AssetBalance(BaseModel):
    """
    Asset balance in a wallet or address.

    Attributes:
        currency_id: Currency identifier.
        currency: Currency symbol.
        blockchain: Blockchain name.
        network: Network name.
        balance: Available balance.
        pending_balance: Pending balance (unconfirmed).
        wallet_id: Associated wallet ID.
        address_id: Associated address ID.
    """

    currency_id: Optional[str] = Field(default=None, description="Currency ID")
    currency: Optional[str] = Field(default=None, description="Currency symbol")
    blockchain: Optional[str] = Field(default=None, description="Blockchain")
    network: Optional[str] = Field(default=None, description="Network")
    balance: Optional[str] = Field(default=None, description="Available balance")
    pending_balance: Optional[str] = Field(default=None, description="Pending balance")
    wallet_id: Optional[str] = Field(default=None, description="Wallet ID")
    address_id: Optional[str] = Field(default=None, description="Address ID")

    model_config = {"frozen": True}


class NFTCollectionBalance(BaseModel):
    """
    NFT collection balance.

    Attributes:
        collection_name: NFT collection name.
        contract_address: Collection contract address.
        blockchain: Blockchain name.
        network: Network name.
        count: Number of NFTs owned.
    """

    collection_name: Optional[str] = Field(default=None, description="Collection name")
    contract_address: Optional[str] = Field(default=None, description="Contract address")
    blockchain: Optional[str] = Field(default=None, description="Blockchain")
    network: Optional[str] = Field(default=None, description="Network")
    count: int = Field(default=0, description="NFT count")

    model_config = {"frozen": True}
