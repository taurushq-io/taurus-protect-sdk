"""Settlement models for Taurus Network."""

from __future__ import annotations

from datetime import datetime
from enum import Enum
from typing import Any, Dict, List, Optional

from pydantic import BaseModel, Field


class SettlementStatus(str, Enum):
    """Settlement status enum."""

    PENDING = "PENDING"
    PENDING_APPROVAL = "PENDING_APPROVAL"
    APPROVED = "APPROVED"
    EXECUTING = "EXECUTING"
    COMPLETED = "COMPLETED"
    REJECTED = "REJECTED"
    CANCELED = "CANCELED"
    FAILED = "FAILED"


class SettlementAssetTransfer(BaseModel):
    """
    Asset transfer within a settlement.

    Represents one leg of an asset transfer in a settlement.

    Attributes:
        currency_id: Currency being transferred.
        amount: Amount to transfer.
        source_shared_address_id: Source shared address.
        destination_shared_address_id: Destination shared address.
    """

    currency_id: Optional[str] = Field(default=None, description="Currency ID")
    amount: Optional[str] = Field(default=None, description="Transfer amount")
    source_shared_address_id: Optional[str] = Field(
        default=None, description="Source shared address ID"
    )
    destination_shared_address_id: Optional[str] = Field(
        default=None, description="Destination shared address ID"
    )

    model_config = {"frozen": True}


class SettlementClipTransaction(BaseModel):
    """Transaction within a settlement clip."""

    id: Optional[str] = Field(default=None, description="Transaction ID")
    tx_hash: Optional[str] = Field(default=None, description="Blockchain transaction hash")
    status: Optional[str] = Field(default=None, description="Transaction status")
    amount: Optional[str] = Field(default=None, description="Transaction amount")
    currency_id: Optional[str] = Field(default=None, description="Currency ID")

    model_config = {"frozen": True}


class SettlementClip(BaseModel):
    """Settlement clip representing a portion of the settlement."""

    id: Optional[str] = Field(default=None, description="Clip ID")
    status: Optional[str] = Field(default=None, description="Clip status")
    transactions: List[SettlementClipTransaction] = Field(
        default_factory=list, description="Clip transactions"
    )

    model_config = {"frozen": True}


class SettlementTrail(BaseModel):
    """Audit trail entry for settlement changes."""

    id: Optional[str] = Field(default=None, description="Trail entry ID")
    timestamp: Optional[datetime] = Field(default=None, description="Trail timestamp")
    action: Optional[str] = Field(default=None, description="Action performed")
    actor: Optional[str] = Field(default=None, description="Actor who performed action")
    comment: Optional[str] = Field(default=None, description="Optional comment")

    model_config = {"frozen": True}


class Settlement(BaseModel):
    """
    Settlement between participants.

    Represents an atomic exchange of assets between two participants.

    Attributes:
        id: Unique settlement identifier.
        creator_participant_id: Participant who created the settlement.
        target_participant_id: Counter-party participant.
        first_leg_participant_id: Who executes the first leg.
        first_leg_assets: Assets in first leg.
        second_leg_assets: Assets in second leg.
        clips: Settlement clips.
        start_execution_date: When execution started.
        status: Settlement status.
        created_at: Creation timestamp.
        updated_at: Update timestamp.
        workflow_id: Associated workflow ID.
        trails: Audit trail.
    """

    id: Optional[str] = Field(default=None, description="Settlement ID")
    creator_participant_id: Optional[str] = Field(
        default=None, description="Creator participant ID"
    )
    target_participant_id: Optional[str] = Field(default=None, description="Target participant ID")
    first_leg_participant_id: Optional[str] = Field(
        default=None, description="First leg participant ID"
    )
    first_leg_assets: List[SettlementAssetTransfer] = Field(
        default_factory=list, description="First leg asset transfers"
    )
    second_leg_assets: List[SettlementAssetTransfer] = Field(
        default_factory=list, description="Second leg asset transfers"
    )
    clips: List[SettlementClip] = Field(default_factory=list, description="Clips")
    start_execution_date: Optional[datetime] = Field(
        default=None, description="Execution start date"
    )
    status: Optional[str] = Field(default=None, description="Settlement status")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Update timestamp")
    workflow_id: Optional[str] = Field(default=None, description="Workflow ID")
    trails: List[SettlementTrail] = Field(default_factory=list, description="Audit trail")

    model_config = {"frozen": True}


# Request models


class SettlementAssetTransferRequest(BaseModel):
    """Asset transfer specification for creating a settlement."""

    currency_id: str = Field(description="Currency ID")
    amount: str = Field(description="Amount to transfer")
    source_shared_address_id: str = Field(description="Source shared address ID")
    destination_shared_address_id: str = Field(description="Destination shared address ID")


class CreateSettlementRequest(BaseModel):
    """Request to create a new settlement."""

    target_participant_id: str = Field(description="Target participant ID")
    first_leg_participant_id: str = Field(description="Who executes the first leg")
    first_leg_assets: List[SettlementAssetTransferRequest] = Field(
        description="First leg asset transfers"
    )
    second_leg_assets: List[SettlementAssetTransferRequest] = Field(
        description="Second leg asset transfers"
    )


class AcceptSettlementRequest(BaseModel):
    """Request to accept a settlement."""

    comment: str = Field(default="", description="Acceptance comment")


class RejectSettlementRequest(BaseModel):
    """Request to reject a settlement."""

    comment: str = Field(description="Rejection reason")


# Filter options


class ListSettlementsOptions(BaseModel):
    """Options for listing settlements."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    statuses: Optional[List[str]] = Field(default=None, description="Filter by statuses")
    participant_id: Optional[str] = Field(
        default=None, description="Filter by participant ID (creator or target)"
    )
