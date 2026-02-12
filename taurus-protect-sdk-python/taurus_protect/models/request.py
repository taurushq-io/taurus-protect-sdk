"""Request models for Taurus-PROTECT SDK."""

from __future__ import annotations

import json
from datetime import datetime
from enum import Enum
from typing import Any, Dict, List, Optional

from pydantic import BaseModel, ConfigDict, Field

from taurus_protect.models.currency import Currency


class RequestStatus(str, Enum):
    """
    Status of a transaction request.

    Request statuses track the progress of a request from creation through
    approval, signing, broadcasting, and confirmation on the blockchain.

    Common status transitions:
    - CREATED -> PENDING -> APPROVED -> BROADCASTED -> CONFIRMED
    - PENDING -> REJECTED (if rejected by an approver)
    - BROADCASTED -> PERMANENT_FAILURE (if transaction fails)
    """

    # Secondary approval received
    APPROVED_2 = "APPROVED_2"
    # Request has been approved and is ready for signing
    APPROVED = "APPROVED"
    # Request is in the process of being approved
    APPROVING = "APPROVING"
    # Auto prepared secondary status
    AUTO_PREPARED_2 = "AUTO_PREPARED_2"
    # Auto prepared status
    AUTO_PREPARED = "AUTO_PREPARED"
    # Broadcasting secondary status
    BROADCASTING_2 = "BROADCASTING_2"
    # Broadcasting status
    BROADCASTING = "BROADCASTING"
    # Broadcasted to the blockchain
    BROADCASTED = "BROADCASTED"
    # Bundle has been approved
    BUNDLE_APPROVED = "BUNDLE_APPROVED"
    # Bundle is being broadcast
    BUNDLE_BROADCASTING = "BUNDLE_BROADCASTING"
    # Bundle is ready
    BUNDLE_READY = "BUNDLE_READY"
    # Request was canceled
    CANCELED = "CANCELED"
    # Transaction confirmed on blockchain
    CONFIRMED = "CONFIRMED"
    # Request has been created
    CREATED = "CREATED"
    # Diem burn MBS approved
    DIEM_BURN_MBS_APPROVED = "DIEM_BURN_MBS_APPROVED"
    # Diem burn MBS pending
    DIEM_BURN_MBS_PENDING = "DIEM_BURN_MBS_PENDING"
    # Diem mint MBS approved
    DIEM_MINT_MBS_APPROVED = "DIEM_MINT_MBS_APPROVED"
    # Diem mint MBS completed
    DIEM_MINT_MBS_COMPLETED = "DIEM_MINT_MBS_COMPLETED"
    # Diem mint MBS pending
    DIEM_MINT_MBS_PENDING = "DIEM_MINT_MBS_PENDING"
    # Request has expired
    EXPIRED = "EXPIRED"
    # Fast approved secondary status
    FAST_APPROVED_2 = "FAST_APPROVED_2"
    # HSM signing failed (secondary)
    HSM_FAILED_2 = "HSM_FAILED_2"
    # HSM signing failed
    HSM_FAILED = "HSM_FAILED"
    # HSM ready for signing (secondary)
    HSM_READY_2 = "HSM_READY_2"
    # HSM ready for signing
    HSM_READY = "HSM_READY"
    # HSM signed (secondary)
    HSM_SIGNED_2 = "HSM_SIGNED_2"
    # HSM has signed the transaction
    HSM_SIGNED = "HSM_SIGNED"
    # Manual broadcast required
    MANUAL_BROADCAST = "MANUAL_BROADCAST"
    # Transaction has been mined
    MINED = "MINED"
    # Transaction partially confirmed
    PARTIALLY_CONFIRMED = "PARTIALLY_CONFIRMED"
    # Request is pending approval
    PENDING = "PENDING"
    # Transaction permanently failed
    PERMANENT_FAILURE = "PERMANENT_FAILURE"
    # Request is ready
    READY = "READY"
    # Request was rejected
    REJECTED = "REJECTED"
    # Transaction has been sent
    SENT = "SENT"
    # Signet transaction completed
    SIGNET_COMPLETED = "SIGNET_COMPLETED"
    # Signet transaction pending
    SIGNET_PENDING = "SIGNET_PENDING"
    # Unknown status
    UNKNOWN = "UNKNOWN"

    @classmethod
    def from_string(cls, value: str) -> "RequestStatus":
        """Convert string to RequestStatus, with fallback to UNKNOWN."""
        if not value:
            return cls.UNKNOWN
        try:
            return cls(value.upper())
        except ValueError:
            return cls.UNKNOWN


class RequestTrail(BaseModel):
    """Audit trail entry for a request."""

    timestamp: Optional[datetime] = Field(default=None, description="When the action occurred")
    action: Optional[str] = Field(default=None, description="Action taken")
    user_id: Optional[str] = Field(default=None, description="User who took action")
    comment: Optional[str] = Field(default=None, description="Optional comment")

    model_config = {"frozen": True}


class RequestMetadataAmount(BaseModel):
    """Amount information extracted from request metadata payload."""

    value_from: Optional[str] = None
    value_to: Optional[str] = None
    rate: Optional[str] = None
    decimals: Optional[int] = None  # Always a small integer
    currency_from: Optional[str] = None
    currency_to: Optional[str] = None

    model_config = ConfigDict(frozen=True)


class RequestMetadata(BaseModel):
    """
    Metadata for a transaction request.

    Contains only the hash and the verified payload string used for
    integrity verification. Use the ``get_*()`` convenience methods
    to extract fields from the cryptographically verified payload.

    SECURITY DESIGN:
        - ``hash``: SHA-256 of ``payload_as_string``, verified against governance rules
        - ``payload_as_string``: canonical JSON used for hash computation (VERIFIED)
        - ``payload`` (raw JSON object): intentionally omitted to prevent use of
          unverified data. An attacker could modify the raw payload while leaving
          ``payload_as_string`` unchanged — the hash would still verify, but the
          client would extract tampered data.
    """

    hash: Optional[str] = Field(default=None, description="SHA-256 hash of payload_as_string")
    payload_as_string: Optional[str] = Field(default=None, description="Verified JSON payload")

    model_config = {"frozen": True}

    def _parse_payload_entries(self) -> list:
        """Parse payload_as_string as JSON array of key-value entries."""
        if not self.payload_as_string:
            return []
        try:
            entries = json.loads(self.payload_as_string)
            return entries if isinstance(entries, list) else []
        except (json.JSONDecodeError, TypeError):
            return []

    def _get_payload_value(self, key: str) -> Any:
        """Get the value for a given key from payload entries."""
        for entry in self._parse_payload_entries():
            if isinstance(entry, dict) and entry.get("key") == key:
                return entry.get("value")
        return None

    def get_source_address(self) -> Optional[str]:
        """Extract source address from verified payload."""
        value = self._get_payload_value("source")
        if isinstance(value, dict):
            payload = value.get("payload")
            if isinstance(payload, dict):
                return payload.get("address")
        return None

    def get_destination_address(self) -> Optional[str]:
        """Extract destination address from verified payload."""
        value = self._get_payload_value("destination")
        if isinstance(value, dict):
            payload = value.get("payload")
            if isinstance(payload, dict):
                return payload.get("address")
        return None

    def get_amount(self) -> Optional[RequestMetadataAmount]:
        """Extract amount information from verified payload."""
        value = self._get_payload_value("amount")
        if not isinstance(value, dict):
            return None
        return RequestMetadataAmount(
            value_from=_json_value_to_string(value.get("valueFrom")),
            value_to=_json_value_to_string(value.get("valueTo")),
            rate=_json_value_to_string(value.get("rate")),
            decimals=int(value["decimals"]) if "decimals" in value else None,
            currency_from=str(value["currencyFrom"]) if "currencyFrom" in value else None,
            currency_to=str(value["currencyTo"]) if "currencyTo" in value else None,
        )


def _json_value_to_string(value: Any) -> Optional[str]:
    """Convert a JSON value to string, handling both string and numeric inputs.

    The API returns valueFrom, valueTo, and rate as JSON strings to support
    arbitrary-precision amounts exceeding 64-bit limits. This helper handles
    both string and numeric JSON values for backward compatibility.
    """
    if value is None:
        return None
    if isinstance(value, str):
        return value
    if isinstance(value, float):
        return f"{value:.20g}"
    return str(value)


class Attribute(BaseModel):
    """Custom key-value attribute that can be attached to various entities."""

    id: Optional[int] = Field(default=None, description="Attribute identifier")
    key: Optional[str] = Field(default=None, description="Attribute key")
    value: Optional[str] = Field(default=None, description="Attribute value")
    content_type: Optional[str] = Field(default=None, description="Content type for file attributes")
    owner: Optional[str] = Field(default=None, description="Owner identifier")
    type: Optional[str] = Field(default=None, description="Type classification")
    sub_type: Optional[str] = Field(default=None, description="Sub-type classification")
    is_file: bool = Field(default=False, description="Whether the value is a file reference")

    model_config = ConfigDict(frozen=True)


class RequestApproversGroup(BaseModel):
    """Single approval group requirement for a request."""

    external_group_id: Optional[str] = Field(default=None, description="External group identifier")
    minimum_signatures: int = Field(default=0, description="Minimum signatures required")

    model_config = ConfigDict(frozen=True)


class RequestParallelApproversGroups(BaseModel):
    """Set of approval groups that operate in parallel."""

    sequential: List[RequestApproversGroup] = Field(
        default_factory=list, description="Sequential approval groups"
    )

    model_config = ConfigDict(frozen=True)


class RequestApprovers(BaseModel):
    """Approvers configuration for a transaction request."""

    parallel: List[RequestParallelApproversGroups] = Field(
        default_factory=list, description="Parallel approver groups"
    )

    model_config = ConfigDict(frozen=True)


class SignedRequest(BaseModel):
    """A signed blockchain transaction associated with a request."""

    id: Optional[str] = None
    signed_request: Optional[str] = None
    status: Optional[RequestStatus] = None
    hash: Optional[str] = None
    block: Optional[int] = None
    details: Optional[str] = None
    creation_date: Optional[datetime] = None
    update_date: Optional[datetime] = None
    broadcast_date: Optional[datetime] = None
    confirmation_date: Optional[datetime] = None

    model_config = ConfigDict(frozen=True)


class Request(BaseModel):
    """
    Transaction request.

    Represents a pending or completed transaction request with full
    audit trail and metadata for integrity verification.
    """

    id: str = Field(description="Unique request identifier")
    tenant_id: Optional[int] = Field(default=None, description="Tenant identifier")
    currency: Optional[str] = Field(default=None, description="Currency symbol")
    envelope: Optional[str] = Field(default=None, description="Unsigned transaction")
    status: RequestStatus = Field(default=RequestStatus.PENDING, description="Request status")
    trails: List[RequestTrail] = Field(default_factory=list, description="Audit trail")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Update timestamp")
    metadata: Optional[RequestMetadata] = Field(default=None, description="Request metadata")
    rule: Optional[str] = Field(default=None, description="Applied governance rule")
    signed_requests: List[SignedRequest] = Field(default_factory=list, description="Signed transactions")
    type: Optional[str] = Field(default=None, description="Request type")
    approvers: Optional[RequestApprovers] = Field(default=None, description="Approval configuration")
    currency_info: Optional[Currency] = Field(default=None, description="Detailed currency information")
    needs_approval_from: List[str] = Field(default_factory=list, description="Pending approvers")
    request_bundle_id: Optional[str] = Field(default=None, description="Request bundle ID")
    external_request_id: Optional[str] = Field(default=None, description="External request ID")
    attributes: List[Attribute] = Field(default_factory=list, description="Custom attributes")

    model_config = {"frozen": True}


class CreateInternalTransferRequest(BaseModel):
    """Request to create an internal transfer between addresses."""

    from_address_id: str = Field(description="Source address ID")
    to_address_id: str = Field(description="Destination address ID")
    amount: str = Field(description="Amount to transfer")
    comment: Optional[str] = Field(default=None, description="Optional comment")
    external_request_id: Optional[str] = Field(default=None, description="External request ID")
    fee_level: Optional[str] = Field(default=None, description="Fee level (low, medium, high)")


class CreateExternalTransferRequest(BaseModel):
    """Request to create an external transfer to a whitelisted address."""

    from_address_id: str = Field(description="Source address ID")
    to_whitelisted_address_id: str = Field(description="Destination whitelisted address ID")
    amount: str = Field(description="Amount to transfer")
    comment: Optional[str] = Field(default=None, description="Optional comment")
    external_request_id: Optional[str] = Field(default=None, description="External request ID")
    fee_level: Optional[str] = Field(default=None, description="Fee level (low, medium, high)")
    memo: Optional[str] = Field(default=None, description="Optional memo/destination tag")


class ListRequestsOptions(BaseModel):
    """Options for listing requests."""

    limit: int = Field(default=50, ge=1, le=1000, description="Maximum items to return")
    offset: int = Field(default=0, ge=0, description="Number of items to skip")
    status: Optional[RequestStatus] = Field(default=None, description="Filter by status")
    currency: Optional[str] = Field(default=None, description="Filter by currency")
    wallet_id: Optional[str] = Field(default=None, description="Filter by wallet ID")
