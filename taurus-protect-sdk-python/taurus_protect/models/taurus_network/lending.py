"""Lending models for Taurus Network."""

from __future__ import annotations

from datetime import datetime
from enum import Enum
from typing import Any, Dict, List, Optional

from pydantic import BaseModel, Field


class LendingAgreementStatus(str, Enum):
    """Lending agreement status enum."""

    PENDING = "PENDING"
    ACTIVE = "ACTIVE"
    COMPLETED = "COMPLETED"
    CANCELED = "CANCELED"
    DEFAULTED = "DEFAULTED"


class LendingCollateralRequirement(BaseModel):
    """Collateral requirement for a lending offer."""

    currency_id: Optional[str] = Field(default=None, description="Currency ID")
    percentage: Optional[str] = Field(default=None, description="Collateral percentage")

    model_config = {"frozen": True}


class LendingOffer(BaseModel):
    """
    Lending offer from a participant.

    Represents an offer to lend assets at a specified yield.

    Attributes:
        id: Unique offer identifier.
        annual_percentage_yield: APY in basis points (525000 = 5.25%).
        duration: Loan duration (e.g., "3M", "1Y").
        collateral_requirement: Required collateral.
        participant_id: Lender participant ID.
        blockchain: Blockchain name.
        network: Network name.
        arg1: Currency argument 1.
        arg2: Currency argument 2.
        currency_info: Currency details.
        annual_percentage_yield_main_unit: APY as human-readable string.
        origin_created_at: Original creation on network.
        created_at: Local creation timestamp.
        updated_at: Local update timestamp.
        amount: Available loan amount.
        amount_main_unit: Amount in main unit.
    """

    id: Optional[str] = Field(default=None, description="Offer ID")
    annual_percentage_yield: Optional[str] = Field(default=None, description="APY in basis points")
    duration: Optional[str] = Field(default=None, description="Loan duration")
    collateral_requirement: Optional[LendingCollateralRequirement] = Field(
        default=None, description="Collateral requirements"
    )
    participant_id: Optional[str] = Field(default=None, description="Lender participant ID")
    blockchain: Optional[str] = Field(default=None, description="Blockchain name")
    network: Optional[str] = Field(default=None, description="Network name")
    arg1: Optional[str] = Field(default=None, description="Currency argument 1")
    arg2: Optional[str] = Field(default=None, description="Currency argument 2")
    currency_info: Optional[Dict[str, Any]] = Field(default=None, description="Currency details")
    annual_percentage_yield_main_unit: Optional[str] = Field(
        default=None, description="APY as percentage"
    )
    origin_created_at: Optional[datetime] = Field(
        default=None, description="Network creation timestamp"
    )
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Update timestamp")
    amount: Optional[str] = Field(default=None, description="Available amount")
    amount_main_unit: Optional[str] = Field(default=None, description="Amount in main unit")

    model_config = {"frozen": True}


class LendingAgreementCollateral(BaseModel):
    """Collateral provided for a lending agreement."""

    currency_id: Optional[str] = Field(default=None, description="Currency ID")
    amount: Optional[str] = Field(default=None, description="Collateral amount")
    pledge_id: Optional[str] = Field(default=None, description="Associated pledge ID")

    model_config = {"frozen": True}


class LendingAgreementTransaction(BaseModel):
    """Transaction related to a lending agreement."""

    id: Optional[str] = Field(default=None, description="Transaction ID")
    transaction_type: Optional[str] = Field(default=None, description="Transaction type")
    amount: Optional[str] = Field(default=None, description="Transaction amount")
    status: Optional[str] = Field(default=None, description="Transaction status")
    tx_hash: Optional[str] = Field(default=None, description="Transaction hash")

    model_config = {"frozen": True}


class LendingAgreement(BaseModel):
    """
    Lending agreement between participants.

    Represents an active loan between a lender and borrower.

    Attributes:
        id: Unique agreement identifier.
        lender_participant_id: Lender participant ID.
        borrower_participant_id: Borrower participant ID.
        lending_offer_id: Associated lending offer ID.
        amount: Loan amount.
        currency_id: Currency ID.
        annual_yield: Annual yield in basis points.
        status: Agreement status.
        duration: Loan duration.
        start_loan_date: When loan started.
        workflow_id: Associated workflow ID.
        borrower_shared_address_id: Borrower's shared address.
        lender_shared_address_id: Lender's shared address.
        collaterals: Provided collateral.
        transactions: Related transactions.
        created_at: Creation timestamp.
        updated_at: Update timestamp.
        annual_yield_main_unit: Yield as percentage.
        currency_info: Currency details.
        amount_main_unit: Amount in main unit.
        repayment_due_date: When repayment is due.
    """

    id: Optional[str] = Field(default=None, description="Agreement ID")
    lender_participant_id: Optional[str] = Field(default=None, description="Lender participant ID")
    borrower_participant_id: Optional[str] = Field(
        default=None, description="Borrower participant ID"
    )
    lending_offer_id: Optional[str] = Field(default=None, description="Associated offer ID")
    amount: Optional[str] = Field(default=None, description="Loan amount")
    currency_id: Optional[str] = Field(default=None, description="Currency ID")
    annual_yield: Optional[str] = Field(default=None, description="Annual yield in basis points")
    status: Optional[str] = Field(default=None, description="Agreement status")
    duration: Optional[str] = Field(default=None, description="Loan duration")
    start_loan_date: Optional[datetime] = Field(default=None, description="Loan start date")
    workflow_id: Optional[str] = Field(default=None, description="Workflow ID")
    borrower_shared_address_id: Optional[str] = Field(
        default=None, description="Borrower's shared address ID"
    )
    lender_shared_address_id: Optional[str] = Field(
        default=None, description="Lender's shared address ID"
    )
    collaterals: List[LendingAgreementCollateral] = Field(
        default_factory=list, description="Collateral"
    )
    transactions: List[LendingAgreementTransaction] = Field(
        default_factory=list, description="Related transactions"
    )
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Update timestamp")
    annual_yield_main_unit: Optional[str] = Field(default=None, description="Yield as percentage")
    currency_info: Optional[Dict[str, Any]] = Field(default=None, description="Currency details")
    amount_main_unit: Optional[str] = Field(default=None, description="Amount in main unit")
    repayment_due_date: Optional[datetime] = Field(default=None, description="Repayment due date")

    model_config = {"frozen": True}


class LendingAgreementAttachment(BaseModel):
    """
    Attachment on a lending agreement.

    Can be embedded content (base64) or external link.

    Attributes:
        id: Attachment identifier.
        lending_agreement_id: Associated agreement ID.
        uploader_participant_id: Participant who uploaded.
        name: Attachment name/filename.
        type: Attachment type (EMBEDDED, EXTERNAL_LINK).
        content_type: MIME type.
        value: Content (base64) or URL.
        file_size: Size in bytes (for embedded).
        created_at: Creation timestamp.
        updated_at: Update timestamp.
    """

    id: Optional[str] = Field(default=None, description="Attachment ID")
    lending_agreement_id: Optional[str] = Field(default=None, description="Associated agreement ID")
    uploader_participant_id: Optional[str] = Field(
        default=None, description="Uploader participant ID"
    )
    name: Optional[str] = Field(default=None, description="Attachment name")
    type: Optional[str] = Field(default=None, description="Type (EMBEDDED, EXTERNAL_LINK)")
    content_type: Optional[str] = Field(default=None, description="MIME type")
    value: Optional[str] = Field(default=None, description="Content or URL")
    file_size: Optional[str] = Field(default=None, description="File size in bytes")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Update timestamp")

    model_config = {"frozen": True}


# Request models


class CreateLendingOfferRequest(BaseModel):
    """Request to create a lending offer."""

    currency_id: str = Field(description="Currency identifier")
    amount: str = Field(description="Loan amount in smallest unit")
    annual_percentage_yield: str = Field(description="APY in basis points")
    duration: str = Field(description="Loan duration (e.g., '3M', '1Y')")
    collateral_requirement: Optional[LendingCollateralRequirement] = Field(
        default=None, description="Collateral requirements"
    )


class UpdateLendingOfferRequest(BaseModel):
    """Request to update a lending offer."""

    amount: Optional[str] = Field(default=None, description="New loan amount")
    annual_percentage_yield: Optional[str] = Field(default=None, description="New APY")


class CreateLendingAgreementRequest(BaseModel):
    """Request to create a lending agreement."""

    lending_offer_id: str = Field(description="Lending offer ID")
    amount: str = Field(description="Loan amount")
    borrower_shared_address_id: str = Field(description="Borrower's shared address ID")
    lender_shared_address_id: str = Field(description="Lender's shared address ID")
    collateral_pledge_ids: Optional[List[str]] = Field(
        default=None, description="Collateral pledge IDs"
    )


class CreateLendingAgreementAttachmentRequest(BaseModel):
    """Request to create a lending agreement attachment."""

    name: str = Field(description="Attachment name")
    type: str = Field(description="Type (EMBEDDED, EXTERNAL_LINK)")
    value: str = Field(description="Content (base64) or URL")
    content_type: Optional[str] = Field(default=None, description="MIME type")


# Filter options


class ListLendingOffersOptions(BaseModel):
    """Options for listing lending offers."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    participant_id: Optional[str] = Field(default=None, description="Filter by participant ID")
    currency_id: Optional[str] = Field(default=None, description="Filter by currency")


class ListLendingAgreementsOptions(BaseModel):
    """Options for listing lending agreements."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    statuses: Optional[List[str]] = Field(default=None, description="Filter by statuses")
    lender_participant_id: Optional[str] = Field(default=None, description="Filter by lender")
    borrower_participant_id: Optional[str] = Field(default=None, description="Filter by borrower")
