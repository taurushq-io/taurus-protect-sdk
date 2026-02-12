"""Shared Address and Shared Asset models for Taurus Network."""

from __future__ import annotations

from datetime import datetime
from enum import Enum
from typing import List, Optional

from pydantic import BaseModel, Field


class SharedAddressStatus(str, Enum):
    """Shared address status enum."""

    PENDING = "PENDING"
    ACTIVE = "ACTIVE"
    REJECTED = "REJECTED"
    DELETED = "DELETED"


class SharedAssetStatus(str, Enum):
    """Shared asset status enum."""

    PENDING = "PENDING"
    ACTIVE = "ACTIVE"
    REJECTED = "REJECTED"
    DELETED = "DELETED"


class SharedAddressTrail(BaseModel):
    """Audit trail entry for shared address changes."""

    id: Optional[str] = Field(default=None, description="Trail entry ID")
    timestamp: Optional[datetime] = Field(default=None, description="Trail timestamp")
    action: Optional[str] = Field(default=None, description="Action performed")
    actor: Optional[str] = Field(default=None, description="Actor who performed action")
    comment: Optional[str] = Field(default=None, description="Optional comment")

    model_config = {"frozen": True}


class SharedAddressProofOfOwnership(BaseModel):
    """Proof of ownership for a shared address."""

    signature: Optional[str] = Field(default=None, description="Ownership signature")
    message: Optional[str] = Field(default=None, description="Signed message")

    model_config = {"frozen": True}


class SharedAddress(BaseModel):
    """
    Shared address between participants.

    A shared address allows one participant (owner) to share an address
    with another participant (target) for pledging and transfers.

    Attributes:
        id: Unique shared address identifier (UUID).
        internal_address_id: Internal address ID (for owner).
        wladdress_id: Whitelisted address ID (for target).
        owner_participant_id: Owner participant ID.
        target_participant_id: Target participant ID.
        blockchain: Blockchain name.
        network: Network name.
        address: Blockchain address string.
        origin_label: Original label.
        origin_creation_date: Network creation date.
        origin_deletion_date: Network deletion date.
        created_at: Local creation timestamp.
        updated_at: Local update timestamp.
        target_accepted_at: When target accepted.
        status: Shared address status.
        proof_of_ownership: Proof of ownership.
        pledges_count: Number of active pledges.
        trails: Audit trail.
    """

    id: Optional[str] = Field(default=None, description="Shared address ID")
    internal_address_id: Optional[str] = Field(
        default=None, description="Internal address ID (owner)"
    )
    wladdress_id: Optional[str] = Field(default=None, description="Whitelisted address ID (target)")
    owner_participant_id: Optional[str] = Field(default=None, description="Owner participant ID")
    target_participant_id: Optional[str] = Field(default=None, description="Target participant ID")
    blockchain: Optional[str] = Field(default=None, description="Blockchain name")
    network: Optional[str] = Field(default=None, description="Network name")
    address: Optional[str] = Field(default=None, description="Blockchain address")
    origin_label: Optional[str] = Field(default=None, description="Original label")
    origin_creation_date: Optional[datetime] = Field(
        default=None, description="Network creation date"
    )
    origin_deletion_date: Optional[datetime] = Field(
        default=None, description="Network deletion date"
    )
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Update timestamp")
    target_accepted_at: Optional[datetime] = Field(
        default=None, description="Target acceptance timestamp"
    )
    status: Optional[str] = Field(default=None, description="Shared address status")
    proof_of_ownership: Optional[SharedAddressProofOfOwnership] = Field(
        default=None, description="Proof of ownership"
    )
    pledges_count: Optional[str] = Field(default=None, description="Number of active pledges")
    trails: List[SharedAddressTrail] = Field(default_factory=list, description="Audit trail")

    model_config = {"frozen": True}


class SharedAssetTrail(BaseModel):
    """Audit trail entry for shared asset changes."""

    id: Optional[str] = Field(default=None, description="Trail entry ID")
    timestamp: Optional[datetime] = Field(default=None, description="Trail timestamp")
    action: Optional[str] = Field(default=None, description="Action performed")
    actor: Optional[str] = Field(default=None, description="Actor who performed action")
    comment: Optional[str] = Field(default=None, description="Optional comment")

    model_config = {"frozen": True}


class SharedAsset(BaseModel):
    """
    Shared asset (token) between participants.

    A shared asset allows participants to share information about
    tokens/contracts for use in settlements and transfers.

    Attributes:
        id: Unique shared asset identifier.
        wl_contract_address_id: Whitelisted contract address ID.
        owner_participant_id: Owner participant ID.
        target_participant_id: Target participant ID.
        blockchain: Blockchain name.
        network: Network name.
        name: Token/asset name.
        symbol: Token symbol.
        decimals: Token decimals.
        contract_address: Smart contract address.
        token_id: Token ID (for NFTs).
        kind: Asset kind (ERC20, ERC721, etc.).
        origin_creation_date: Network creation date.
        origin_deletion_date: Network deletion date.
        created_at: Local creation timestamp.
        updated_at: Local update timestamp.
        target_accepted_at: Target acceptance timestamp.
        target_rejected_at: Target rejection timestamp.
        status: Shared asset status.
        trails: Audit trail.
    """

    id: Optional[str] = Field(default=None, description="Shared asset ID")
    wl_contract_address_id: Optional[str] = Field(
        default=None, description="Whitelisted contract address ID"
    )
    owner_participant_id: Optional[str] = Field(default=None, description="Owner participant ID")
    target_participant_id: Optional[str] = Field(default=None, description="Target participant ID")
    blockchain: Optional[str] = Field(default=None, description="Blockchain name")
    network: Optional[str] = Field(default=None, description="Network name")
    name: Optional[str] = Field(default=None, description="Token/asset name")
    symbol: Optional[str] = Field(default=None, description="Token symbol")
    decimals: Optional[str] = Field(default=None, description="Token decimals")
    contract_address: Optional[str] = Field(default=None, description="Smart contract address")
    token_id: Optional[str] = Field(default=None, description="Token ID (for NFTs)")
    kind: Optional[str] = Field(default=None, description="Asset kind")
    origin_creation_date: Optional[datetime] = Field(
        default=None, description="Network creation date"
    )
    origin_deletion_date: Optional[datetime] = Field(
        default=None, description="Network deletion date"
    )
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Update timestamp")
    target_accepted_at: Optional[datetime] = Field(
        default=None, description="Target acceptance timestamp"
    )
    target_rejected_at: Optional[datetime] = Field(
        default=None, description="Target rejection timestamp"
    )
    status: Optional[str] = Field(default=None, description="Shared asset status")
    trails: List[SharedAssetTrail] = Field(default_factory=list, description="Audit trail")

    model_config = {"frozen": True}


# Request models


class CreateSharedAddressRequest(BaseModel):
    """Request to create a new shared address."""

    internal_address_id: str = Field(description="Internal address ID to share")
    target_participant_id: str = Field(description="Target participant ID")
    label: Optional[str] = Field(default=None, description="Label for the shared address")
    proof_of_ownership_signature: Optional[str] = Field(
        default=None, description="Proof of ownership signature"
    )
    proof_of_ownership_message: Optional[str] = Field(
        default=None, description="Proof of ownership message"
    )


class AcceptSharedAddressRequest(BaseModel):
    """Request to accept a shared address."""

    comment: str = Field(default="", description="Acceptance comment")


class RejectSharedAddressRequest(BaseModel):
    """Request to reject a shared address."""

    comment: str = Field(description="Rejection reason")


class CreateSharedAssetRequest(BaseModel):
    """Request to create a new shared asset."""

    wl_contract_address_id: str = Field(description="Whitelisted contract address ID")
    target_participant_id: str = Field(description="Target participant ID")


class AcceptSharedAssetRequest(BaseModel):
    """Request to accept a shared asset."""

    comment: str = Field(default="", description="Acceptance comment")


class RejectSharedAssetRequest(BaseModel):
    """Request to reject a shared asset."""

    comment: str = Field(description="Rejection reason")


# Filter options


class ListSharedAddressesOptions(BaseModel):
    """Options for listing shared addresses."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    statuses: Optional[List[str]] = Field(default=None, description="Filter by statuses")
    participant_id: Optional[str] = Field(
        default=None, description="Filter by participant ID (owner or target)"
    )
    blockchain: Optional[str] = Field(default=None, description="Filter by blockchain")
    network: Optional[str] = Field(default=None, description="Filter by network")


class ListSharedAssetsOptions(BaseModel):
    """Options for listing shared assets."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    statuses: Optional[List[str]] = Field(default=None, description="Filter by statuses")
    participant_id: Optional[str] = Field(
        default=None, description="Filter by participant ID (owner or target)"
    )
    blockchain: Optional[str] = Field(default=None, description="Filter by blockchain")
    network: Optional[str] = Field(default=None, description="Filter by network")
