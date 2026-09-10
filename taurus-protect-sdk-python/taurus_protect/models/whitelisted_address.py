"""Whitelisted address models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import Any, Dict, List, Optional

from pydantic import BaseModel, Field

from taurus_protect.models.pagination import Pagination


class InternalAddress(BaseModel):
    """An internal address linked to a whitelisted address."""

    id: Optional[str] = Field(default=None, description="Address identifier")
    address: Optional[str] = Field(default=None, description="Blockchain address")
    label: Optional[str] = Field(default=None, description="Human-readable label")

    model_config = {"frozen": True}


class InternalWallet(BaseModel):
    """An internal wallet linked to a whitelisted address."""

    id: int = Field(description="Wallet identifier")
    path: Optional[str] = Field(default=None, description="Wallet path")
    label: Optional[str] = Field(default=None, description="Human-readable label")

    model_config = {"frozen": True}


class WhitelistedAddress(BaseModel):
    """
    A whitelisted external address.

    Whitelisted addresses are pre-approved destinations for withdrawals.
    They must be verified with cryptographic signatures before use.

    Attributes:
        id: Unique identifier.
        address: Blockchain address string.
        label: Human-readable label.
        currency: Currency/blockchain.
        network: Network name.
        status: Whitelisting status.
        created_at: Creation timestamp.
        contract_type: Optional contract type for smart contracts.
        memo: Optional memo/destination tag.
        customer_id: Optional customer identifier.
        address_type: Optional address type.
        tn_participant_id: Optional Taurus Network participant ID.
        exchange_account_id: Optional exchange account ID.
    """

    id: str = Field(description="Unique identifier")
    address: Optional[str] = Field(default=None, description="Blockchain address")
    label: Optional[str] = Field(default=None, description="Human-readable label")
    currency: Optional[str] = Field(default=None, description="Currency/blockchain")
    network: Optional[str] = Field(default=None, description="Network")
    status: Optional[str] = Field(default=None, description="Whitelisting status")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    contract_type: Optional[str] = Field(default=None, description="Contract type")
    memo: Optional[str] = Field(default=None, description="Memo/destination tag")
    customer_id: Optional[str] = Field(default=None, description="Customer identifier")
    address_type: Optional[str] = Field(default=None, description="Address type")
    tn_participant_id: Optional[str] = Field(
        default=None, description="Taurus Network participant ID"
    )
    exchange_account_id: Optional[str] = Field(
        default=None, description="Exchange account ID"
    )
    linked_internal_addresses: List[InternalAddress] = Field(
        default_factory=list, description="Linked internal addresses"
    )
    linked_wallets: List[InternalWallet] = Field(
        default_factory=list, description="Linked internal wallets"
    )
    attributes: Dict[str, Any] = Field(default_factory=dict, description="Custom attributes")

    model_config = {"frozen": True}


class WhitelistSignature(BaseModel):
    """A signature on a whitelisted address."""

    user_id: Optional[str] = Field(default=None, description="Signing user ID")
    signature: Optional[str] = Field(default=None, description="Base64 signature")
    hash: Optional[str] = Field(default=None, description="Hash that was signed")
    hashes: List[str] = Field(default_factory=list, description="All hashes covered by signature")

    model_config = {"frozen": True}


class WhitelistMetadata(BaseModel):
    """Metadata for a whitelisted address envelope."""

    hash: Optional[str] = Field(default=None, description="Hash of payload")
    payload_as_string: Optional[str] = Field(default=None, description="Signed payload JSON")

    model_config = {"frozen": True}


class SignedWhitelistedAddress(BaseModel):
    """Signed whitelisted address data with signatures."""

    payload: Optional[str] = Field(default=None, description="Base64 signed payload")
    signatures: List[WhitelistSignatureEntry] = Field(
        default_factory=list, description="Cryptographic signatures"
    )

    model_config = {"frozen": True}


class SignedWhitelistedAddressEnvelope(BaseModel):
    """
    Envelope containing a whitelisted address with signatures.

    This envelope contains all the cryptographic signatures needed
    to verify the whitelisted address was properly approved.
    """

    metadata: Optional[WhitelistMetadata] = Field(default=None)
    blockchain: Optional[str] = Field(default=None, description="Blockchain identifier")
    network: Optional[str] = Field(default=None, description="Network identifier")
    rules_container: Optional[str] = Field(default=None, description="Base64 rules container")
    rules_signatures: Optional[str] = Field(default=None, description="Base64 rules signatures")
    signatures: List[WhitelistSignature] = Field(default_factory=list)
    signed_address: Optional[SignedWhitelistedAddress] = Field(
        default=None, description="Signed address data with signatures"
    )
    linked_wallets: List[InternalWallet] = Field(
        default_factory=list, description="Linked internal wallets"
    )
    rules_container_hash: Optional[str] = Field(
        default=None, description="Hash of the rules container (for normalized caching)"
    )
    verified_whitelisted_address: Optional[WhitelistedAddress] = Field(
        default=None, description="Verified whitelisted address (set after 6-step verification)"
    )
    verified_rules_container: Optional[Any] = Field(
        default=None, description="Verified rules container (set after 6-step verification)"
    )

    model_config = {"frozen": False}


class CreateWhitelistedAddressRequest(BaseModel):
    """Request to create a whitelisted address."""

    address: str = Field(description="Blockchain address to whitelist")
    label: str = Field(description="Human-readable label")
    currency: str = Field(description="Currency/blockchain")
    network: Optional[str] = Field(default=None, description="Network")
    contract_type: Optional[str] = Field(default=None, description="Contract type")
    attributes: Dict[str, Any] = Field(default_factory=dict, description="Custom attributes")

    model_config = {"frozen": True}


class WhitelistedAssetMetadata(BaseModel):
    """Metadata for a whitelisted asset envelope."""

    hash: Optional[str] = Field(default=None, description="Hash of payload")
    # SECURITY: payload field intentionally omitted - use payload_as_string only.
    # The raw payload object could be tampered with by an attacker while
    # payload_as_string remains unchanged (hash still verifies). By not having
    # this field, we enforce that all data extraction uses the verified source.
    payload_as_string: Optional[str] = Field(default=None, description="Signed payload JSON")

    model_config = {"frozen": True}


class WhitelistUserSignature(BaseModel):
    """A user's signature on a whitelist entry."""

    user_id: Optional[str] = Field(default=None, description="Signing user ID")
    signature: Optional[str] = Field(default=None, description="Base64 signature")
    comment: Optional[str] = Field(default=None, description="Optional comment")

    model_config = {"frozen": True}


class WhitelistSignatureEntry(BaseModel):
    """A signature entry with hashes covered."""

    user_signature: Optional[WhitelistUserSignature] = Field(
        default=None, description="User signature details"
    )
    hashes: List[str] = Field(default_factory=list, description="Hashes covered by signature")

    model_config = {"frozen": True}


class SignedContractAddress(BaseModel):
    """Signed contract address data with signatures."""

    payload: Optional[str] = Field(default=None, description="Base64 signed payload")
    signatures: List[WhitelistSignatureEntry] = Field(
        default_factory=list, description="Cryptographic signatures"
    )

    model_config = {"frozen": True}


class WhitelistedAsset(BaseModel):
    """
    A whitelisted asset/token.

    Whitelisted assets are pre-approved tokens that can be transferred.
    When verification is enabled, all retrieved assets are cryptographically
    verified using the 5-step verification flow.
    """

    id: str = Field(description="Unique identifier")
    tenant_id: Optional[str] = Field(default=None, description="Tenant ID")
    name: Optional[str] = Field(default=None, description="Asset name")
    symbol: Optional[str] = Field(default=None, description="Asset symbol")
    blockchain: Optional[str] = Field(default=None, description="Blockchain")
    network: Optional[str] = Field(default=None, description="Network")
    contract_address: Optional[str] = Field(default=None, description="Token contract")
    decimals: Optional[int] = Field(
        default=None,
        description=(
            "Token decimal precision, from the verified payload. Security-critical: "
            "it scales every amount denominated in this asset, so it is never taken "
            "from the DTO. Java and TypeScript already exposed it; this SDK did not."
        ),
    )
    token_id: Optional[str] = Field(
        default=None, description="Token identifier within a contract (NFTs)"
    )
    status: Optional[str] = Field(default=None, description="Whitelisting status")
    action: Optional[str] = Field(default=None, description="Action type")
    rule: Optional[str] = Field(default=None, description="Governance rule")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    # Verification fields
    metadata: Optional[WhitelistedAssetMetadata] = Field(
        default=None, description="Asset metadata with hash"
    )
    rules_container: Optional[str] = Field(default=None, description="Base64 rules container")
    rules_signatures: Optional[str] = Field(default=None, description="Base64 rules signatures")
    signed_contract_address: Optional[SignedContractAddress] = Field(
        default=None, description="Signed payload and signatures"
    )
    business_rule_enabled: bool = Field(
        default=False, description="Whether business rule is enabled"
    )

    model_config = {"frozen": True}


class ExcludedWhitelistedAddress(BaseModel):
    """A row dropped from a list because it failed integrity verification, and why."""

    id: Optional[str] = Field(
        default=None, description="The address ID, or None when the row carried no usable ID"
    )
    reason: str = Field(description="Why the row failed verification")

    model_config = {"frozen": True}


class WhitelistedAddressListResult(BaseModel):
    """Result of a paginated whitelisted-address list query.

    Verification is lenient by design: a row that cannot be verified is excluded
    rather than failing the whole call, because one bad row used to deny access to
    every good one -- and listing is how an operator finds the bad row. Excluding
    stays fail-closed, since an omitted destination cannot be selected.

    The omission is REPORTED, not just logged: a caller cannot read the SDK's logger,
    and a shortened list must never be mistaken for a complete one.
    """

    addresses: List[WhitelistedAddress] = Field(
        default_factory=list, description="Verified addresses in the current page"
    )
    pagination: Optional[Pagination] = Field(
        default=None,
        description=(
            "Page window, with total_items already reduced by the number of excluded "
            "rows so has_more stays honest"
        ),
    )
    excluded_unverified: List[ExcludedWhitelistedAddress] = Field(
        default_factory=list, description="Rows dropped from addresses, with the reason"
    )

    model_config = {"frozen": True}
