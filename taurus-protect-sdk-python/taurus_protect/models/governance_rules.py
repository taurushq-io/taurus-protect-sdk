"""Governance rules domain models."""

from __future__ import annotations

import threading
from datetime import datetime
from typing import TYPE_CHECKING, Dict, List, Optional

from pydantic import BaseModel, Field, PrivateAttr, field_validator

from taurus_protect.models.rule_cell import RuleCell

if TYPE_CHECKING:
    from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey


class SuperAdminPublicKey(BaseModel):
    """A SuperAdmin public key."""

    user_id: Optional[str] = Field(default=None, description="ID of the SuperAdmin user")
    public_key: Optional[str] = Field(default=None, description="PEM-encoded public key")

    model_config = {"frozen": True}


class RuleUserSignature(BaseModel):
    """A user's signature on governance rules."""

    user_id: Optional[str] = Field(default=None, description="ID of the signing user")
    signature: Optional[str] = Field(default=None, description="Base64-encoded signature")

    model_config = {"frozen": True}


class RuleUser(BaseModel):
    """A user defined in governance rules."""

    id: Optional[str] = Field(default=None, description="User ID")
    name: Optional[str] = Field(default=None, description="User name")
    public_key_pem: Optional[str] = Field(default=None, description="PEM-encoded public key")
    roles: List[str] = Field(default_factory=list, description="User roles")
    properties: Dict[str, bytes] = Field(default_factory=dict)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class RuleGroup(BaseModel):
    """A group defined in governance rules."""

    id: Optional[str] = Field(default=None, description="Group ID")
    name: Optional[str] = Field(default=None, description="Group name")
    user_ids: List[str] = Field(default_factory=list, description="Member user IDs")
    properties: Dict[str, bytes] = Field(default_factory=dict)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class GroupThreshold(BaseModel):
    """Threshold configuration for a group."""

    group_id: Optional[str] = Field(default=None)
    minimum_signatures: int = Field(default=0, description="Minimum signatures required")
    # Alias for backward compatibility
    threshold: int = Field(default=0)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}

    def get_min_signatures(self) -> int:
        """Get minimum signatures, preferring minimum_signatures over threshold."""
        return self.minimum_signatures if self.minimum_signatures > 0 else self.threshold


class SequentialThresholds(BaseModel):
    """Sequential approval thresholds."""

    thresholds: List[GroupThreshold] = Field(default_factory=list)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class RuleSourceInternalWallet(BaseModel):
    """Internal wallet specification in a rule source."""

    path: Optional[str] = Field(default=None, description="Wallet path")

    model_config = {"frozen": True}


class RuleSourceInternalAddress(BaseModel):
    """Internal address specification in a rule source."""

    address: Optional[str] = Field(default=None)
    path: Optional[str] = Field(default=None)

    model_config = {"frozen": True}


class RuleSourceExchange(BaseModel):
    """Exchange specification in a rule source."""

    label: Optional[str] = Field(default=None)

    model_config = {"frozen": True}


class RuleSourceExternalAddress(BaseModel):
    """External address specification in a rule source."""

    address: Optional[str] = Field(default=None)
    memo: Optional[str] = Field(default=None)

    model_config = {"frozen": True}


class RuleSource(BaseModel):
    """Source specification in a whitelist rule.

    Sources whose type/content is unknown to this SDK version keep their exact
    wire bytes in ``raw`` and are re-emitted verbatim on encode.
    """

    type: int = Field(default=0, description="Source type (0=Any, 1=InternalWallet, ...)")
    internal_wallet: Optional[RuleSourceInternalWallet] = Field(
        default=None, description="Internal wallet details (when type=1)"
    )
    internal_address: Optional[RuleSourceInternalAddress] = Field(default=None)
    exchange: Optional[RuleSourceExchange] = Field(default=None)
    external_address: Optional[RuleSourceExternalAddress] = Field(default=None)
    raw: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


# Constants for RuleSource.type matching the protobuf RuleSourceType enum.
RULE_SOURCE_TYPE_ANY = 0
RULE_SOURCE_TYPE_INTERNAL_WALLET = 1
RULE_SOURCE_TYPE_INTERNAL_ADDRESS = 2
RULE_SOURCE_TYPE_ANY_EXCHANGE = 3
RULE_SOURCE_TYPE_EXCHANGE = 4
RULE_SOURCE_TYPE_EXTERNAL_ADDRESS = 5


class AddressWhitelistingLine(BaseModel):
    """Source-specific rule line for address whitelisting."""

    cells: List[RuleSource] = Field(default_factory=list, description="Rule source specifications")
    parallel_thresholds: List[SequentialThresholds] = Field(
        default_factory=list, description="Approval requirements for this line"
    )
    properties: Dict[str, bytes] = Field(default_factory=dict)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class AddressWhitelistingRules(BaseModel):
    """Rules for address whitelisting."""

    currency: Optional[str] = Field(default=None, description="Blockchain/currency")
    network: Optional[str] = Field(default=None, description="Network")
    parallel_thresholds: List[SequentialThresholds] = Field(default_factory=list)
    lines: List[AddressWhitelistingLine] = Field(default_factory=list)
    include_network_in_payload: bool = Field(default=False)
    properties: Dict[str, bytes] = Field(default_factory=dict)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class ContractAddressWhitelistingRules(BaseModel):
    """Rules for contract address whitelisting."""

    blockchain: Optional[str] = Field(default=None)
    network: Optional[str] = Field(default=None)
    parallel_thresholds: List[SequentialThresholds] = Field(default_factory=list)
    properties: Dict[str, bytes] = Field(default_factory=dict)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class RuleColumn(BaseModel):
    """Column definition in transaction rules."""

    type: Optional[str] = Field(default=None, description="Column type enum value name")
    name: Optional[str] = Field(default=None, description="Human-readable column name")
    metadata_key: Optional[str] = Field(default=None, description="Request-metadata key matched")
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class EvmCallContract(BaseModel):
    """EVM contract-call scoping."""

    contract_type: Optional[str] = Field(default=None)
    method_signature: Optional[str] = Field(default=None)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class XtzCallContract(BaseModel):
    """Tezos contract-call scoping."""

    contract_type: Optional[str] = Field(default=None)
    method_signature: Optional[str] = Field(default=None)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class CashSettlement(BaseModel):
    """Cash-settlement scoping."""

    provider: Optional[str] = Field(default=None)
    request_type: Optional[str] = Field(default=None)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class CosmosDetails(BaseModel):
    """Cosmos scoping."""

    method_signatures: List[str] = Field(default_factory=list)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class TransactionRuleDetails(BaseModel):
    """Additional transaction rule configuration."""

    domain: Optional[str] = Field(default=None, description="Rule domain")
    sub_domain: Optional[str] = Field(default=None, description="Rule sub-domain")
    blockchain: Optional[str] = Field(default=None)
    network: Optional[str] = Field(default=None)
    evm_call_contract: Optional[EvmCallContract] = Field(default=None)
    xtz_call_contract: Optional[XtzCallContract] = Field(default=None)
    cash_settlement: Optional[CashSettlement] = Field(default=None)
    cosmos_details: Optional[CosmosDetails] = Field(default=None)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class RuleLine(BaseModel):
    """Line/row in transaction rules. Cells align positionally with the rule's columns."""

    cells: List[RuleCell] = Field(default_factory=list, description="Typed cell values")
    parallel_thresholds: List[SequentialThresholds] = Field(
        default_factory=list, description="Approval requirements for this line"
    )
    priority: int = Field(default=0)
    properties: Dict[str, bytes] = Field(default_factory=dict)
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class TransactionRules(BaseModel):
    """Transaction approval rules for a specific action type."""

    key: Optional[str] = Field(default=None, description="Rule key (e.g., blockchain/action type)")
    columns: List[RuleColumn] = Field(default_factory=list, description="Column definitions")
    lines: List[RuleLine] = Field(default_factory=list, description="Rule lines")
    details: Optional[TransactionRuleDetails] = Field(
        default=None, description="Additional rule configuration"
    )
    unknown_fields: bytes = Field(default=b"", repr=False)

    model_config = {"frozen": True}


class DecodedRulesContainer(BaseModel):
    """
    Decoded governance rules container.

    Contains users, groups, and rules for transaction and address approval.
    The HSM public key is lazily resolved from users with the HSMSLOT role.
    """

    users: List[RuleUser] = Field(default_factory=list)
    groups: List[RuleGroup] = Field(default_factory=list)
    minimum_distinct_user_signatures: int = Field(default=0)
    minimum_distinct_group_signatures: int = Field(default=0)
    transaction_rules: List[TransactionRules] = Field(default_factory=list)
    address_whitelisting_rules: List[AddressWhitelistingRules] = Field(default_factory=list)
    contract_address_whitelisting_rules: List[ContractAddressWhitelistingRules] = Field(
        default_factory=list
    )
    enforced_rules_hash: Optional[str] = Field(default=None)
    timestamp: int = Field(default=0)
    minimum_commitment_signatures: int = Field(default=0)
    engine_identities: List[str] = Field(default_factory=list)
    hsm_slot_id: int = Field(default=0)
    properties: Dict[str, bytes] = Field(default_factory=dict)
    unknown_fields: bytes = Field(default=b"", repr=False)

    # Private fields for cached HSM key (using PrivateAttr for Pydantic compatibility)
    _hsm_public_key: Optional["EllipticCurvePublicKey"] = PrivateAttr(default=None)
    _hsm_public_key_resolved: bool = PrivateAttr(default=False)
    _hsm_lock: threading.Lock = PrivateAttr(default_factory=threading.Lock)

    # Cache for decoded public keys (keyed by PEM string) to avoid repeated PEM parsing
    _key_cache: Dict[str, "EllipticCurvePublicKey"] = PrivateAttr(default_factory=dict)
    _key_cache_lock: threading.Lock = PrivateAttr(default_factory=threading.Lock)

    model_config = {"frozen": False, "arbitrary_types_allowed": True}

    def get_hsm_public_key(self) -> Optional["EllipticCurvePublicKey"]:
        """
        Get the HSM slot public key.

        Finds the first user with the HSMSLOT role and returns their public key.
        The result is cached for subsequent calls. This method is thread-safe.

        Returns:
            The HSM public key, or None if no user with HSMSLOT role exists.
        """
        with self._hsm_lock:
            if not self._hsm_public_key_resolved:
                self._hsm_public_key = self._find_hsm_public_key()
                self._hsm_public_key_resolved = True
            return self._hsm_public_key

    def _find_hsm_public_key(self) -> Optional["EllipticCurvePublicKey"]:
        """Find the HSM public key from users with HSMSLOT role."""
        for user in self.users:
            if "HSMSLOT" in user.roles and user.public_key_pem:
                try:
                    return self.get_user_public_key(user.public_key_pem)
                except (ValueError, TypeError) as e:
                    # Log and continue to next user if this key is malformed
                    # ValueError: Invalid PEM format or unsupported key type
                    # TypeError: Invalid key data type
                    import logging

                    logging.getLogger(__name__).warning(
                        "Failed to decode HSM public key for user %s: %s",
                        user.id,
                        str(e),
                    )
                    continue
        return None

    def get_user_public_key(self, pem: str) -> "EllipticCurvePublicKey":
        """
        Get a decoded public key from PEM string, using a cache to avoid repeated parsing.

        This is thread-safe and caches decoded keys for the lifetime of the rules container.

        Args:
            pem: PEM-encoded public key string.

        Returns:
            Decoded EllipticCurvePublicKey.

        Raises:
            ValueError: If PEM decoding fails or key uses unsupported curve.
        """
        with self._key_cache_lock:
            cached = self._key_cache.get(pem)
            if cached is not None:
                return cached

        # Decode outside lock to avoid holding lock during crypto operations
        from taurus_protect.crypto.keys import decode_public_key_pem

        key = decode_public_key_pem(pem)

        with self._key_cache_lock:
            self._key_cache[pem] = key

        return key

    def find_user_by_id(self, user_id: str) -> Optional[RuleUser]:
        """Find a user by ID."""
        for user in self.users:
            if user.id == user_id:
                return user
        return None

    def find_group_by_id(self, group_id: str) -> Optional[RuleGroup]:
        """Find a group by ID."""
        for group in self.groups:
            if group.id == group_id:
                return group
        return None

    def find_address_whitelisting_rules(
        self, blockchain: str, network: Optional[str] = None
    ) -> Optional[AddressWhitelistingRules]:
        """
        Find address whitelisting rules for a blockchain/network.

        Uses a three-tier priority system:
        1. Exact match - both blockchain and network match
        2. Blockchain-only match - blockchain matches, rule has wildcard network
        3. Global default - rule has wildcard blockchain ("Any" or empty)

        Args:
            blockchain: The blockchain identifier.
            network: Optional network identifier.

        Returns:
            Matching rules or None.
        """
        blockchain_only_match: Optional[AddressWhitelistingRules] = None
        global_default: Optional[AddressWhitelistingRules] = None

        for rule in self.address_whitelisting_rules:
            is_global_default = self._is_wildcard(rule.currency)
            blockchain_matches = not is_global_default and rule.currency == blockchain
            network_matches = rule.network == network
            has_wildcard_network = self._is_wildcard(rule.network)

            # Priority 1: Exact match
            if blockchain_matches and network_matches:
                return rule

            # Priority 2: Blockchain match with wildcard network
            if blockchain_matches and has_wildcard_network and blockchain_only_match is None:
                blockchain_only_match = rule

            # Priority 3: Global default
            if is_global_default and global_default is None:
                global_default = rule

        return blockchain_only_match or global_default

    def find_contract_address_whitelisting_rules(
        self, blockchain: str, network: Optional[str] = None
    ) -> Optional[ContractAddressWhitelistingRules]:
        """
        Find contract address whitelisting rules for a blockchain/network.

        Uses a three-tier priority system:
        1. Exact match - both blockchain and network match
        2. Blockchain-only match - blockchain matches, rule has wildcard network
        3. Global default - rule has wildcard blockchain ("Any" or empty)

        Args:
            blockchain: The blockchain identifier.
            network: Optional network identifier.

        Returns:
            Matching rules or None.
        """
        blockchain_only_match: Optional[ContractAddressWhitelistingRules] = None
        global_default: Optional[ContractAddressWhitelistingRules] = None

        for rule in self.contract_address_whitelisting_rules:
            is_global_default = self._is_wildcard(rule.blockchain)
            blockchain_matches = not is_global_default and rule.blockchain == blockchain
            network_matches = rule.network == network
            has_wildcard_network = self._is_wildcard(rule.network)

            # Priority 1: Exact match
            if blockchain_matches and network_matches:
                return rule

            # Priority 2: Blockchain match with wildcard network
            if blockchain_matches and has_wildcard_network and blockchain_only_match is None:
                blockchain_only_match = rule

            # Priority 3: Global default
            if is_global_default and global_default is None:
                global_default = rule

        return blockchain_only_match or global_default

    @staticmethod
    def _is_wildcard(value: Optional[str]) -> bool:
        """Check if a value is a wildcard (None, empty, or 'Any')."""
        return value is None or value == "" or value.lower() == "any"

    def has_unknown_fields(self) -> bool:
        """Report whether any node carries data unknown to this SDK version.

        Unknown protobuf fields, raw cells, and raw whitelisting sources are
        preserved verbatim on re-encode; a true result is a hint to upgrade the
        SDK before editing the container.
        """
        if self.unknown_fields:
            return True
        if any(u.unknown_fields for u in self.users):
            return True
        if any(g.unknown_fields for g in self.groups):
            return True
        for tr in self.transaction_rules:
            if tr.unknown_fields:
                return True
            if any(c.unknown_fields for c in tr.columns):
                return True
            if _rule_details_unknown(tr.details):
                return True
            for line in tr.lines:
                if line.unknown_fields or _thresholds_unknown(line.parallel_thresholds):
                    return True
                if any(getattr(cell, "kind", None) == "RawCell" for cell in line.cells):
                    return True
        for awr in self.address_whitelisting_rules:
            if awr.unknown_fields or _thresholds_unknown(awr.parallel_thresholds):
                return True
            for line in awr.lines:
                if line.unknown_fields or _thresholds_unknown(line.parallel_thresholds):
                    return True
                if any(s.raw for s in line.cells):
                    return True
        for cawr in self.contract_address_whitelisting_rules:
            if cawr.unknown_fields or _thresholds_unknown(cawr.parallel_thresholds):
                return True
        return False


def _thresholds_unknown(thresholds: "List[SequentialThresholds]") -> bool:
    for st in thresholds:
        if st.unknown_fields or any(t.unknown_fields for t in st.thresholds):
            return True
    return False


def _rule_details_unknown(details: "Optional[TransactionRuleDetails]") -> bool:
    """Report unknown fields on the details node and each nested scoping sub-message.

    Checking only the details node would report a container clean while a nested
    contract-call scoping node carried data from a newer schema.
    """
    if details is None:
        return False
    if details.unknown_fields:
        return True
    nested = (
        details.evm_call_contract,
        details.xtz_call_contract,
        details.cash_settlement,
        details.cosmos_details,
    )
    return any(n is not None and n.unknown_fields for n in nested)


class GovernanceRulesTrail(BaseModel):
    """Audit trail entry for governance rules changes."""

    user_id: Optional[str] = Field(default=None)
    action: Optional[str] = Field(default=None)
    timestamp: Optional[datetime] = Field(default=None)

    model_config = {"frozen": True}


class GovernanceRules(BaseModel):
    """
    Governance rules for the Taurus-PROTECT tenant.

    Contains the signed rules container and SuperAdmin signatures.
    The rules container must be verified before use.
    """

    rules_container: Optional[str] = Field(
        default=None, description="Base64-encoded protobuf rules container"
    )
    rules_signatures: List[RuleUserSignature] = Field(
        default_factory=list, description="SuperAdmin signatures"
    )
    locked: bool = Field(default=False, description="Whether rules are locked")
    creation_date: Optional[datetime] = Field(default=None)
    update_date: Optional[datetime] = Field(default=None)
    trails: List[GovernanceRulesTrail] = Field(default_factory=list)

    @field_validator("locked", mode="before")
    @classmethod
    def coerce_locked(cls, v):
        """Coerce None to False for the locked field."""
        if v is None:
            return False
        return v

    # Cached decoded container
    _decoded_container: Optional[DecodedRulesContainer] = None

    model_config = {"frozen": False, "arbitrary_types_allowed": True}


class ExcludedRuleset(BaseModel):
    """A history entry withheld because its SuperAdmin signatures did not verify."""

    creation_date: Optional[datetime] = Field(
        default=None,
        description="When the excluded ruleset was created -- the only stable identifier "
        "a history entry carries. Named to match GovernanceRules.creation_date; Go calls "
        "the same field CreatedAt.",
    )
    reason: str = Field(default="", description="Why it was excluded")

    model_config = {"frozen": True}


class GovernanceRulesHistoryResult(BaseModel):
    """Result of a governance rules history query with cursor-based pagination."""

    rules: List[GovernanceRules] = Field(
        default_factory=list, description="Rules in this page whose signatures verified"
    )
    excluded_unverified: List[ExcludedRuleset] = Field(
        default_factory=list,
        description="Entries withheld because their signatures did not verify, so a "
        "shortened page cannot read as a complete one. History is LENIENT and does not "
        "fail when nothing survives: a SuperAdmin key rotation makes every pre-rotation "
        "ruleset unverifiable, and aborting would deny the whole audit trail.",
    )
    cursor: Optional[str] = Field(default=None, description="Cursor for fetching the next page")
    total_items: Optional[str] = Field(
        default=None, description="Total number of items available, reduced by the exclusions"
    )

    model_config = {"frozen": True}
