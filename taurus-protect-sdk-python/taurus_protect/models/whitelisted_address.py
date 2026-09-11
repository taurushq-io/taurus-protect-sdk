"""Whitelisted address models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import Any, Dict, Iterable, List, Optional

from pydantic import BaseModel, Field, PrivateAttr

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

    # id -> the metadata hash the row carried in THIS read. Private because it must not
    # be settable from outside: it is the content pin the approval path checks against.
    # (A private attribute is assignable even on a frozen model -- pydantic routes
    # underscore names past the frozen check.)
    _reviewed_hashes: Dict[str, str] = PrivateAttr(default_factory=dict)

    @classmethod
    def from_verified_envelopes(
        cls,
        envelopes: Iterable["SignedWhitelistedAddressEnvelope"],
        pagination: Optional[Pagination] = None,
        excluded_unverified: Optional[List[ExcludedWhitelistedAddress]] = None,
    ) -> "WhitelistedAddressListResult":
        """Build a result from the envelopes a verified read produced.

        Envelopes rather than addresses, because in THIS SDK the metadata hash is
        reachable only on the envelope -- ``WhitelistedAddress`` has no metadata field,
        unlike Go's and Java's. Taking both from the same source is what keeps
        ``addresses`` and the approval pin from drifting.

        Args:
            envelopes: The envelopes that survived verification.
            pagination: Page window, with ``total_items`` already reduced.
            excluded_unverified: Rows withheld, with their reasons.

        Returns:
            The result, carrying the reviewed hash for every row it lists.
        """
        rows = list(envelopes)
        result = cls(
            addresses=[
                e.verified_whitelisted_address
                for e in rows
                if e.verified_whitelisted_address is not None
            ],
            pagination=pagination,
            excluded_unverified=excluded_unverified or [],
        )
        result._reviewed_hashes = {
            str(e.verified_whitelisted_address.id): e.metadata.hash
            for e in rows
            if e.verified_whitelisted_address is not None
            and e.metadata is not None
            and e.metadata.hash
        }
        return result

    def select(self, *ids: str) -> "WhitelistedAddressApproval":
        """Pin the given rows from this verified read, for approval.

        An id this read did not return is an error rather than a silent omission: it
        means the caller is trying to approve something this read did not give them --
        either it was excluded as unverifiable, or it was never on the page -- and
        approving fewer rows than asked for would tell the approver they approved more
        than they did.

        Args:
            *ids: The whitelisted address IDs to approve, as strings.

        Returns:
            The pinned selection to hand to ``WhitelistedAddressService.approve``.

        Raises:
            ValueError: If no ids were given, or an id is not in this result.
            IntegrityError: If a selected row carries no metadata hash, leaving
                nothing to pin the approval to.
        """
        from taurus_protect.errors import IntegrityError

        if not ids:
            raise ValueError("cannot select an empty set of ids")

        listed = {str(a.id) for a in self.addresses}
        pinned: Dict[str, str] = {}
        for address_id in ids:
            key = str(address_id)
            if key not in listed:
                raise ValueError(
                    f"whitelisted address {key} is not in this verified read: it was "
                    "either excluded as unverifiable or not on this page"
                )
            reviewed = self._reviewed_hashes.get(key)
            if not reviewed:
                raise IntegrityError(
                    f"whitelisted address {key} carries no metadata hash, so there is "
                    "nothing to pin the approval to"
                )
            pinned[key] = reviewed
        return WhitelistedAddressApproval(pinned)

    def select_all(self) -> "WhitelistedAddressApproval":
        """Pin every row this verified read returned.

        Note it pins what SURVIVED verification, not what the server sent: rows in
        ``excluded_unverified`` are not included, so a caller who wants to know about
        them must read that field. Approving is all-or-nothing over what is pinned here.

        Returns:
            The pinned selection.

        Raises:
            ValueError: If this read returned no verified addresses.
        """
        ids = [str(a.id) for a in self.addresses]
        if not ids:
            raise ValueError("this read returned no verified addresses to approve")
        return self.select(*ids)


class WhitelistedAddressApproval:
    """The set of rows an approver reviewed, with the metadata hash each carried AT
    REVIEW TIME. It is the content pin the approval path signs against.

    Why this type exists rather than a plain list of ids. The approval API accepts only
    ids: the SDK re-reads them and signs whatever the server returns under those ids.
    Nothing bound the approver's intent to the bytes signed, so a response-controlling
    server could answer the id-filtered re-read with a DIFFERENT row -- one whose
    existing signatures already satisfy the container it presents -- and harvest a
    genuine approver signature over content the approver never saw. This is the same
    shape ``GovernanceRuleService.approve_rules_proposal`` was hardened against with its
    mandatory ``expected_container_hash``.

    Mint it with :meth:`WhitelistedAddressListResult.select` or
    :meth:`WhitelistedAddressListResult.select_all`; ``_pinned`` is private so a
    hand-built value pins nothing and is refused by the approval path. The threat model
    here is the SERVER, not the caller: the point is that the pin cannot be FORGOTTEN.
    """

    __slots__ = ("_pinned",)

    def __init__(self, pinned: Dict[str, str]) -> None:
        self._pinned = dict(pinned)

    def ids(self) -> List[str]:
        """The pinned row ids, in no particular order.

        The approval path sorts them itself, because the endpoint requires ascending
        order and the signed array must not depend on the order the caller selected in.
        """
        return list(self._pinned)

    def pinned_hash(self, address_id: str) -> Optional[str]:
        """The reviewed metadata hash for ``address_id``, or None if it was not pinned."""
        return self._pinned.get(str(address_id))

    def is_empty(self) -> bool:
        """Whether this selection pins nothing -- true for a hand-built value."""
        return not self._pinned


class WhitelistedAssetApproval:
    """The asset peer of :class:`WhitelistedAddressApproval`.

    Minted from the ASSET OBJECTS a verified read returned, rather than from a result
    object, because this SDK merged asset and envelope: ``WhitelistedAsset`` carries
    ``metadata``, so the reviewed hash is reachable off the rows themselves. (The
    address side cannot do that -- ``WhitelistedAddress`` has no metadata field, which
    is the Python-specific difference the repo already records.) The guarantee is the
    same either way: the pin cannot be forgotten, and the adversary is the server.
    """

    __slots__ = ("_pinned",)

    def __init__(self, pinned: Dict[str, str]) -> None:
        self._pinned = dict(pinned)

    @classmethod
    def select(
        cls, assets: Iterable["WhitelistedAsset"], *ids: str
    ) -> "WhitelistedAssetApproval":
        """Pin the given ids out of the assets a verified read returned.

        Args:
            assets: The rows from ``WhitelistedAssetService.list_for_approval`` (or
                ``list``) -- every one of which has been through the verifier.
            *ids: The whitelisted asset IDs to approve, as strings.

        Returns:
            The pinned selection to hand to ``WhitelistedAssetService.approve``.

        Raises:
            ValueError: If no ids were given, or an id is not among ``assets``.
            IntegrityError: If a selected row carries no metadata hash.
        """
        from taurus_protect.errors import IntegrityError

        if not ids:
            raise ValueError("cannot select an empty set of ids")

        by_id = {str(a.id): a for a in assets if a is not None}
        pinned: Dict[str, str] = {}
        for asset_id in ids:
            key = str(asset_id)
            asset = by_id.get(key)
            if asset is None:
                raise ValueError(f"whitelisted asset {key} is not in this verified read")
            if asset.metadata is None or not asset.metadata.hash:
                raise IntegrityError(
                    f"whitelisted asset {key} carries no metadata hash, so there is "
                    "nothing to pin the approval to"
                )
            pinned[key] = asset.metadata.hash
        return cls(pinned)

    @classmethod
    def select_all(
        cls, assets: Iterable["WhitelistedAsset"]
    ) -> "WhitelistedAssetApproval":
        """Pin every asset a verified read returned.

        Raises:
            ValueError: If ``assets`` is empty.
        """
        rows = [a for a in assets if a is not None]
        if not rows:
            raise ValueError("this read returned no verified assets to approve")
        return cls.select(rows, *[str(a.id) for a in rows])

    def ids(self) -> List[str]:
        """The pinned row ids, in no particular order."""
        return list(self._pinned)

    def pinned_hash(self, asset_id: str) -> Optional[str]:
        """The reviewed metadata hash for ``asset_id``, or None if it was not pinned."""
        return self._pinned.get(str(asset_id))

    def is_empty(self) -> bool:
        """Whether this selection pins nothing -- true for a hand-built value."""
        return not self._pinned
