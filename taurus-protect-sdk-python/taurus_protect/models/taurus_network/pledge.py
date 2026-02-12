"""Pledge models for Taurus-PROTECT SDK Taurus Network."""

from __future__ import annotations

from datetime import datetime
from enum import Enum
from typing import List, Optional

from pydantic import BaseModel, Field


class PledgeStatus(str, Enum):
    """Pledge status enum."""

    PENDING = "PENDING"
    CONFIRMED = "CONFIRMED"
    REJECTED = "REJECTED"
    CANCELED = "CANCELED"
    UNPLEDGED = "UNPLEDGED"


class PledgeType(str, Enum):
    """Pledge type for withdrawal rights."""

    NO_WITHDRAWALS_RIGHTS = "NO_WITHDRAWALS_RIGHTS"
    PLEDGEE_WITHDRAWALS_RIGHTS = "PLEDGEE_WITHDRAWALS_RIGHTS"
    PLEDGEE_AUTO_APPROVED_WITHDRAWALS_RIGHTS = "PLEDGEE_AUTO_APPROVED_WITHDRAWALS_RIGHTS"


class PledgeActionStatus(str, Enum):
    """Pledge action status enum."""

    PENDING = "PENDING"
    APPROVED = "APPROVED"
    REJECTED = "REJECTED"
    EXECUTED = "EXECUTED"
    CANCELED = "CANCELED"


class PledgeActionType(str, Enum):
    """Pledge action type enum."""

    CREATE_PLEDGE = "CREATE_PLEDGE"
    ADD_COLLATERAL = "ADD_COLLATERAL"
    WITHDRAW = "WITHDRAW"
    INITIATE_WITHDRAW = "INITIATE_WITHDRAW"
    UNPLEDGE = "UNPLEDGE"
    REJECT = "REJECT"


class PledgeWithdrawalStatus(str, Enum):
    """Pledge withdrawal status enum."""

    PENDING = "PENDING"
    PENDING_APPROVAL = "PENDING_APPROVAL"
    APPROVED = "APPROVED"
    REJECTED = "REJECTED"
    EXECUTED = "EXECUTED"
    CANCELED = "CANCELED"


class PledgeAttribute(BaseModel):
    """Custom attribute on a pledge."""

    key: str = Field(description="Attribute key")
    value: str = Field(description="Attribute value")

    model_config = {"frozen": True}


class PledgeDurationSetup(BaseModel):
    """Duration configuration for a pledge."""

    start_date: Optional[datetime] = Field(default=None, description="Start date of the pledge")
    end_date: Optional[datetime] = Field(default=None, description="End date of the pledge")

    model_config = {"frozen": True}


class PledgeTrail(BaseModel):
    """Audit trail entry for a pledge."""

    id: str = Field(description="Trail entry ID")
    action: str = Field(default="", description="Action performed")
    actor: str = Field(default="", description="Who performed the action")
    timestamp: Optional[datetime] = Field(default=None, description="When it occurred")
    comment: Optional[str] = Field(default=None, description="Optional comment")

    model_config = {"frozen": True}


class Pledge(BaseModel):
    """
    Taurus Network pledge.

    A pledge represents funds that are reserved and pledged from one
    participant (pledgor) to another (pledgee) on a shared address.

    Attributes:
        id: Unique pledge identifier.
        shared_address_id: ID of the shared address holding pledged funds.
        owner_participant_id: ID of the pledgor (owner) participant.
        target_participant_id: ID of the pledgee (target) participant.
        currency_id: Currency identifier.
        blockchain: Blockchain name.
        network: Network name.
        amount: Pledged amount in smallest currency unit.
        status: Current pledge status.
        pledge_type: Withdrawal rights configuration.
        direction: Direction of the pledge (incoming/outgoing).
        external_reference_id: External reference for reconciliation.
        reconciliation_note: Internal reconciliation note.
        wl_address_id: Whitelisted address ID if pledgee.
        origin_creation_date: Original creation timestamp.
        unpledge_date: When pledge was unpledged (if applicable).
        duration_setup: Duration configuration.
        attributes: Custom key-value attributes.
        trails: Audit trail entries.
        created_at: When the pledge was created.
        updated_at: When the pledge was last updated.
    """

    id: str = Field(description="Unique pledge identifier")
    shared_address_id: str = Field(
        default="", description="Shared address ID holding the pledged funds"
    )
    owner_participant_id: str = Field(default="", description="Pledgor (owner) participant ID")
    target_participant_id: str = Field(default="", description="Pledgee (target) participant ID")
    currency_id: str = Field(default="", description="Currency identifier")
    blockchain: str = Field(default="", description="Blockchain name")
    network: str = Field(default="", description="Network name")
    arg1: Optional[str] = Field(default=None, description="Additional argument 1")
    arg2: Optional[str] = Field(default=None, description="Additional argument 2")
    amount: str = Field(default="0", description="Pledged amount in smallest unit")
    status: str = Field(default="", description="Pledge status")
    pledge_type: str = Field(default="", description="Pledge type for withdrawal rights")
    direction: str = Field(default="", description="Direction (incoming/outgoing)")
    external_reference_id: Optional[str] = Field(default=None, description="External reference ID")
    reconciliation_note: Optional[str] = Field(
        default=None, description="Internal reconciliation note"
    )
    wl_address_id: Optional[str] = Field(default=None, description="Whitelisted address ID")
    origin_creation_date: Optional[datetime] = Field(
        default=None, description="Original creation timestamp"
    )
    unpledge_date: Optional[datetime] = Field(default=None, description="Unpledge timestamp")
    duration_setup: Optional[PledgeDurationSetup] = Field(
        default=None, description="Duration configuration"
    )
    attributes: List[PledgeAttribute] = Field(default_factory=list, description="Custom attributes")
    trails: List[PledgeTrail] = Field(default_factory=list, description="Audit trail")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}


class PledgeActionMetadata(BaseModel):
    """Metadata for a pledge action."""

    hash: str = Field(default="", description="Hash of the action metadata")
    payload: Optional[str] = Field(default=None, description="Payload data as string")

    model_config = {"frozen": True}


class PledgeActionTrail(BaseModel):
    """Audit trail entry for a pledge action."""

    id: str = Field(description="Trail entry ID")
    action: str = Field(default="", description="Action performed")
    actor: str = Field(default="", description="Who performed the action")
    timestamp: Optional[datetime] = Field(default=None, description="When it occurred")
    comment: Optional[str] = Field(default=None, description="Optional comment")

    model_config = {"frozen": True}


class PledgeAction(BaseModel):
    """
    Pledge action requiring approval.

    Pledge actions represent pending operations on pledges that require
    approval before execution.

    Attributes:
        id: Unique action identifier.
        pledge_id: Associated pledge ID.
        action_type: Type of action.
        status: Current action status.
        metadata: Action metadata including hash for signing.
        rule: Rule applied to this action.
        needs_approval_from: List of approvers needed.
        pledge_withdrawal_id: Associated withdrawal ID if applicable.
        envelope: Envelope data.
        trails: Audit trail.
        created_at: Creation timestamp.
        last_approval_date: Last approval timestamp.
    """

    id: str = Field(description="Unique action identifier")
    pledge_id: str = Field(default="", description="Associated pledge ID")
    action_type: str = Field(default="", description="Type of action")
    status: str = Field(default="", description="Action status")
    metadata: Optional[PledgeActionMetadata] = Field(default=None, description="Action metadata")
    rule: Optional[str] = Field(default=None, description="Applied rule")
    needs_approval_from: List[str] = Field(default_factory=list, description="Required approvers")
    pledge_withdrawal_id: Optional[str] = Field(
        default=None, description="Associated withdrawal ID"
    )
    envelope: Optional[str] = Field(default=None, description="Envelope data")
    trails: List[PledgeActionTrail] = Field(default_factory=list, description="Audit trail")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    last_approval_date: Optional[datetime] = Field(
        default=None, description="Last approval timestamp"
    )

    model_config = {"frozen": True}


class PledgeWithdrawalTrail(BaseModel):
    """Audit trail entry for a pledge withdrawal."""

    id: str = Field(description="Trail entry ID")
    action: str = Field(default="", description="Action performed")
    actor: str = Field(default="", description="Who performed the action")
    timestamp: Optional[datetime] = Field(default=None, description="When it occurred")
    comment: Optional[str] = Field(default=None, description="Optional comment")

    model_config = {"frozen": True}


class PledgeWithdrawal(BaseModel):
    """
    Pledge withdrawal record.

    Represents a withdrawal request or execution from a pledge.

    Attributes:
        id: Unique withdrawal identifier.
        pledge_id: Associated pledge ID.
        destination_shared_address_id: Destination address.
        amount: Withdrawal amount.
        status: Withdrawal status.
        tx_hash: Blockchain transaction hash.
        tx_id: Internal transaction ID.
        request_id: Associated request ID.
        tx_block_number: Block number of transaction.
        initiator_participant_id: Who initiated the withdrawal.
        external_reference_id: External reference.
        trails: Audit trail.
        created_at: Creation timestamp.
    """

    id: str = Field(description="Unique withdrawal identifier")
    pledge_id: str = Field(default="", description="Associated pledge ID")
    destination_shared_address_id: str = Field(
        default="", description="Destination shared address ID"
    )
    amount: str = Field(default="0", description="Withdrawal amount")
    status: str = Field(default="", description="Withdrawal status")
    tx_hash: Optional[str] = Field(default=None, description="Transaction hash")
    tx_id: Optional[str] = Field(default=None, description="Internal transaction ID")
    request_id: Optional[str] = Field(default=None, description="Associated request ID")
    tx_block_number: Optional[str] = Field(default=None, description="Block number")
    initiator_participant_id: Optional[str] = Field(
        default=None, description="Initiator participant ID"
    )
    external_reference_id: Optional[str] = Field(default=None, description="External reference ID")
    trails: List[PledgeWithdrawalTrail] = Field(default_factory=list, description="Audit trail")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")

    model_config = {"frozen": True}


# Request models


class CreatePledgeRequest(BaseModel):
    """
    Request to create a new pledge.

    Example:
        >>> request = CreatePledgeRequest(
        ...     shared_address_id="123",
        ...     currency_id="BTC",
        ...     amount="1000000000",
        ...     pledge_type=PledgeType.PLEDGEE_WITHDRAWALS_RIGHTS,
        ... )
    """

    shared_address_id: str = Field(description="Shared address ID for pledge")
    currency_id: str = Field(description="Currency identifier")
    amount: str = Field(description="Amount in smallest currency unit")
    pledge_type: str = Field(
        default=PledgeType.NO_WITHDRAWALS_RIGHTS.value,
        description="Withdrawal rights type",
    )
    pledge_duration_setup: Optional[PledgeDurationSetup] = Field(
        default=None, description="Duration configuration"
    )
    key_value_attributes: Optional[List[PledgeAttribute]] = Field(
        default=None, description="Custom attributes"
    )
    external_reference_id: Optional[str] = Field(default=None, description="External reference")
    reconciliation_note: Optional[str] = Field(default=None, description="Reconciliation note")


class UpdatePledgeRequest(BaseModel):
    """Request to update a pledge's default destination."""

    default_destination_shared_address_id: Optional[str] = Field(
        default=None, description="Default destination shared address ID"
    )
    default_destination_internal_address_id: Optional[str] = Field(
        default=None, description="Default destination internal address ID"
    )


class AddPledgeCollateralRequest(BaseModel):
    """Request to add collateral to an existing pledge."""

    amount: str = Field(description="Amount to add in smallest currency unit")


class WithdrawPledgeRequest(BaseModel):
    """
    Request for pledgee withdrawal from pledge.

    Use either destination_shared_address_id or destination_internal_address_id.
    """

    amount: str = Field(description="Amount to withdraw")
    destination_shared_address_id: Optional[str] = Field(
        default=None, description="Destination shared address ID"
    )
    destination_internal_address_id: Optional[str] = Field(
        default=None, description="Destination internal address ID"
    )
    external_reference_id: Optional[str] = Field(default=None, description="External reference ID")


class InitiateWithdrawPledgeRequest(BaseModel):
    """Request for pledgor-initiated withdrawal from pledge."""

    amount: str = Field(description="Amount to withdraw")
    destination_shared_address_id: Optional[str] = Field(
        default=None, description="Destination shared address ID"
    )


class RejectPledgeRequest(BaseModel):
    """Request to reject a pledge."""

    comment: str = Field(description="Rejection comment")


class ApprovePledgeActionsRequest(BaseModel):
    """Request to approve pledge actions."""

    ids: List[str] = Field(description="IDs of pledge actions to approve")
    comment: str = Field(default="", description="Approval comment")
    signature: str = Field(description="ECDSA signature over action hashes")


class RejectPledgeActionsRequest(BaseModel):
    """Request to reject pledge actions."""

    ids: List[str] = Field(description="IDs of pledge actions to reject")
    comment: str = Field(description="Rejection comment")


# Filter options


class ListPledgesOptions(BaseModel):
    """Options for listing pledges."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    statuses: Optional[List[str]] = Field(default=None, description="Filter by statuses")
    currency_id: Optional[str] = Field(default=None, description="Filter by currency")
    direction: Optional[str] = Field(
        default=None, description="Filter by direction (incoming/outgoing)"
    )
    participant_id: Optional[str] = Field(default=None, description="Filter by participant ID")


class ListPledgeActionsOptions(BaseModel):
    """Options for listing pledge actions."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    statuses: Optional[List[str]] = Field(default=None, description="Filter by statuses")
    action_types: Optional[List[str]] = Field(default=None, description="Filter by action types")
    pledge_id: Optional[str] = Field(default=None, description="Filter by pledge ID")


class ListPledgeWithdrawalsOptions(BaseModel):
    """Options for listing pledge withdrawals."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    statuses: Optional[List[str]] = Field(default=None, description="Filter by statuses")
    pledge_id: Optional[str] = Field(default=None, description="Filter by pledge ID")
