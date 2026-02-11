"""Address models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, Field

from taurus_protect.models.balance import Balance


class AddressAttribute(BaseModel):
    """Custom attribute on an address."""

    id: str = Field(description="Attribute identifier")
    key: str = Field(description="Attribute name")
    value: str = Field(description="Attribute value")

    model_config = {"frozen": True}


class Address(BaseModel):
    """
    Blockchain address within a wallet.

    Attributes:
        id: Unique address identifier.
        wallet_id: ID of the containing wallet.
        address: The blockchain address string.
        alternate_address: Alternate address format if available.
        label: Human-readable label.
        comment: Optional description.
        currency: Currency symbol.
        customer_id: Optional customer identifier.
        external_address_id: Optional external identifier.
        address_path: HD derivation path.
        address_index: Index used for address generation.
        nonce: Current address nonce.
        status: Creation status (created, creating, signed, observed, confirmed).
        balance: Current address balance.
        signature: Cryptographic signature for integrity verification.
        disabled: Whether the address is disabled.
        can_use_all_funds: Whether all funds can be used.
        created_at: When the address was created.
        updated_at: When the address was last updated.
        attributes: Custom key-value attributes.
        linked_whitelisted_address_ids: IDs of linked whitelisted addresses.
    """

    id: str = Field(description="Unique address identifier")
    wallet_id: str = Field(description="Containing wallet ID")
    address: str = Field(description="Blockchain address string")
    alternate_address: Optional[str] = Field(default=None, description="Alternate format")
    label: Optional[str] = Field(default=None, description="Human-readable label")
    comment: Optional[str] = Field(default=None, description="Optional description")
    currency: str = Field(default="", description="Currency symbol")
    customer_id: Optional[str] = Field(default=None, description="Customer identifier")
    external_address_id: Optional[str] = Field(default=None, description="External ID")
    address_path: Optional[str] = Field(default=None, description="HD derivation path")
    address_index: Optional[str] = Field(default=None, description="Address index")
    nonce: Optional[str] = Field(default=None, description="Current nonce")
    status: Optional[str] = Field(default=None, description="Creation status")
    balance: Optional[Balance] = Field(default=None, description="Current balance")
    signature: Optional[str] = Field(default=None, description="Integrity signature")
    disabled: bool = Field(default=False, description="Whether disabled")
    can_use_all_funds: bool = Field(default=False, description="Can use all funds")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Update timestamp")
    attributes: List[AddressAttribute] = Field(
        default_factory=list, description="Custom attributes"
    )
    linked_whitelisted_address_ids: List[str] = Field(
        default_factory=list, description="Linked whitelisted address IDs"
    )

    model_config = {"frozen": True}


class CreateAddressRequest(BaseModel):
    """Request to create a new address."""

    wallet_id: str = Field(description="Wallet ID (required)")
    label: str = Field(description="Address label (required)")
    comment: Optional[str] = Field(default=None, description="Optional description")
    customer_id: Optional[str] = Field(default=None, description="Customer identifier")
    external_address_id: Optional[str] = Field(default=None, description="External ID")
    address_type: Optional[str] = Field(
        default=None, description="Address type (e.g., p2pkh, p2wpkh)"
    )
    non_hardened_derivation: bool = Field(default=False, description="Use non-hardened derivation")


class ListAddressesOptions(BaseModel):
    """Options for listing addresses."""

    wallet_id: Optional[str] = Field(default=None, description="Filter by wallet ID")
    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    query: Optional[str] = Field(default=None, description="Search query")
    exclude_disabled: bool = Field(default=False, description="Exclude disabled addresses")
