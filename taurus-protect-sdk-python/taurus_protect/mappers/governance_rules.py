"""Governance rules mapper utilities."""

from __future__ import annotations

import base64
import json
import logging
from typing import Any, Dict, List, Optional

from taurus_protect.errors import IntegrityError
from taurus_protect.models.governance_rules import (
    RULE_SOURCE_TYPE_INTERNAL_WALLET,
    AddressWhitelistingLine,
    AddressWhitelistingRules,
    ContractAddressWhitelistingRules,
    DecodedRulesContainer,
    GroupThreshold,
    RuleColumn,
    RuleGroup,
    RuleLine,
    RuleSource,
    RuleSourceInternalWallet,
    RuleUser,
    RuleUserSignature,
    SequentialThresholds,
    TransactionRuleDetails,
    TransactionRules,
)

_logger = logging.getLogger(__name__)


def _try_protobuf_decode(data: bytes) -> Optional[DecodedRulesContainer]:
    """
    Attempt to decode rules container from protobuf bytes.

    Returns None if protobuf decoding fails or is not available.
    """
    try:
        # Import protobuf modules - these may fail if imports aren't properly configured
        from taurus_protect._internal.proto import request_reply_pb2

        # Parse the protobuf message (same approach as Java SDK)
        pb_container = request_reply_pb2.RulesContainer()
        pb_container.ParseFromString(data)

        return _rules_container_from_proto(pb_container)
    except ImportError as e:
        _logger.debug("Protobuf import failed (using JSON fallback): %s", e)
        return None
    except Exception as e:
        # Log the full error for debugging
        _logger.warning("Protobuf parsing failed: %s", e)
        return None


def _rules_container_from_proto(pb: Any) -> DecodedRulesContainer:
    """Convert a protobuf RulesContainer to the model."""
    from taurus_protect._internal.proto import request_reply_pb2

    users = []
    for u in pb.users:
        # Convert enum integer values to string names using Role.Name()
        roles = [request_reply_pb2.Role.Name(role) for role in u.roles]
        users.append(
            RuleUser(
                id=u.id,
                name=getattr(u, "name", None),
                public_key_pem=u.publicKey,
                roles=roles,
            )
        )

    groups = []
    for g in pb.groups:
        groups.append(
            RuleGroup(
                id=g.id,
                name=getattr(g, "name", None),
                user_ids=(
                    list(g.userIds) if hasattr(g, "userIds") else list(getattr(g, "user_ids", []))
                ),
            )
        )

    address_whitelisting_rules = []
    for r in pb.addressWhitelistingRules:
        parallel_thresholds = [_sequential_thresholds_from_proto(pt) for pt in r.parallelThresholds]
        lines = [_address_whitelisting_line_from_proto(line) for line in r.lines]
        address_whitelisting_rules.append(
            AddressWhitelistingRules(
                currency=r.currency,
                network=r.network,
                parallel_thresholds=parallel_thresholds,
                lines=lines,
            )
        )

    contract_address_whitelisting_rules = []
    for r in pb.contractAddressWhitelistingRules:
        parallel_thresholds = [_sequential_thresholds_from_proto(pt) for pt in r.parallelThresholds]
        # Convert protobuf Blockchain enum int to string name (e.g., 5 -> "ETH")
        # Same as Java's blockchain.name() via @Named("blockchainToString")
        blockchain_name = request_reply_pb2.Blockchain.Name(r.blockchain)
        contract_address_whitelisting_rules.append(
            ContractAddressWhitelistingRules(
                blockchain=blockchain_name,
                network=r.network,
                parallel_thresholds=parallel_thresholds,
            )
        )

    # Parse transaction rules
    transaction_rules = []
    for tr in pb.transactionRules:
        transaction_rules.append(_transaction_rules_from_proto(tr))

    return DecodedRulesContainer(
        users=users,
        groups=groups,
        minimum_distinct_user_signatures=pb.minimumDistinctUserSignatures,
        minimum_distinct_group_signatures=pb.minimumDistinctGroupSignatures,
        transaction_rules=transaction_rules,
        address_whitelisting_rules=address_whitelisting_rules,
        contract_address_whitelisting_rules=contract_address_whitelisting_rules,
        enforced_rules_hash=pb.enforcedRulesHash,
        timestamp=pb.timestamp,
        hsm_slot_id=pb.hsmSlotId if hasattr(pb, "hsmSlotId") else 0,
    )


def _transaction_rules_from_proto(pb_tr: Any) -> TransactionRules:
    """Convert protobuf TransactionRules to model."""
    from taurus_protect._internal.proto import request_reply_pb2

    columns = []
    for col in pb_tr.columns:
        col_type = str(col.type) if hasattr(col, "type") else None
        columns.append(RuleColumn(type=col_type))

    lines = []
    for line in pb_tr.lines:
        cells = [str(c) for c in line.cells] if hasattr(line, "cells") else []
        parallel_thresholds = [
            _sequential_thresholds_from_proto(pt)
            for pt in line.parallelThresholds
        ] if hasattr(line, "parallelThresholds") else []
        lines.append(RuleLine(cells=cells, parallel_thresholds=parallel_thresholds))

    details = None
    if hasattr(pb_tr, "details") and pb_tr.details:
        # domain and sub_domain are protobuf enums (RuleDomain, RuleSubDomain) —
        # convert to string names to match Java SDK's domain.name() behavior.
        # Use try/except for unknown enum values (API may add new values before
        # proto definitions are regenerated).
        _TRD = request_reply_pb2.RulesContainer.TransactionRules.TransactionRuleDetails
        domain_val = pb_tr.details.domain if hasattr(pb_tr.details, "domain") else 0
        sub_domain_val = pb_tr.details.subDomain if hasattr(pb_tr.details, "subDomain") else 0
        domain_name = None
        if domain_val:
            try:
                domain_name = _TRD.RuleDomain.Name(domain_val)
            except ValueError:
                domain_name = str(domain_val)
        sub_domain_name = None
        if sub_domain_val:
            try:
                sub_domain_name = _TRD.RuleSubDomain.Name(sub_domain_val)
            except ValueError:
                sub_domain_name = str(sub_domain_val)
        details = TransactionRuleDetails(
            domain=domain_name,
            sub_domain=sub_domain_name,
        )

    return TransactionRules(
        key=pb_tr.key,
        columns=columns,
        lines=lines,
        details=details,
    )


def _address_whitelisting_line_from_proto(pb_line: Any) -> AddressWhitelistingLine:
    """Convert protobuf AddressWhitelistingRules.Line to model."""
    cells = []
    for cell_bytes in pb_line.cells:
        source = _rule_source_from_bytes(cell_bytes)
        if source is not None:
            cells.append(source)

    parallel_thresholds = [
        _sequential_thresholds_from_proto(pt) for pt in pb_line.parallelThresholds
    ]

    return AddressWhitelistingLine(
        cells=cells,
        parallel_thresholds=parallel_thresholds,
    )


def _rule_source_from_bytes(data: bytes) -> Optional[RuleSource]:
    """Decode a RuleSource from serialized protobuf bytes."""
    try:
        from taurus_protect._internal.proto import request_reply_pb2

        pb_source = request_reply_pb2.RuleSource()
        pb_source.ParseFromString(data)

        source_type = int(pb_source.type)
        internal_wallet = None

        if source_type == RULE_SOURCE_TYPE_INTERNAL_WALLET and pb_source.payload:
            try:
                pb_wallet = request_reply_pb2.RuleSourceInternalWallet()
                pb_wallet.ParseFromString(pb_source.payload)
                internal_wallet = RuleSourceInternalWallet(path=pb_wallet.path)
            except Exception:
                pass

        return RuleSource(type=source_type, internal_wallet=internal_wallet)
    except Exception:
        return None


def _sequential_thresholds_from_proto(pb: Any) -> SequentialThresholds:
    """Convert protobuf SequentialThresholds to model."""
    thresholds = []
    for t in pb.thresholds:
        thresholds.append(
            GroupThreshold(
                group_id=t.groupId,
                minimum_signatures=t.minimumSignatures,
            )
        )
    return SequentialThresholds(thresholds=thresholds)


def rules_container_from_base64(base64_data: str) -> DecodedRulesContainer:
    """
    Decode a base64-encoded rules container.

    This function attempts to decode the rules container using the following
    priority:
    1. Protobuf decoding (if available and data is valid protobuf)
    2. JSON decoding (fallback for JSON-encoded containers)
    3. Empty container (if all parsing fails)

    Args:
        base64_data: Base64-encoded rules container.

    Returns:
        Decoded rules container.
    """
    if not base64_data:
        return DecodedRulesContainer()

    try:
        decoded = base64.b64decode(base64_data)

        # Try protobuf first
        result = _try_protobuf_decode(decoded)
        if result is not None:
            return result

        # Fall back to JSON parsing
        try:
            data = json.loads(decoded.decode("utf-8"))
            return _parse_rules_container_from_dict(data)
        except (json.JSONDecodeError, UnicodeDecodeError):
            # Neither protobuf nor JSON worked - this is a security-critical failure
            raise IntegrityError("Failed to decode rules container: not valid protobuf or JSON")
    except (ValueError, TypeError) as e:
        raise IntegrityError(f"Failed to decode rules container from base64: {e}")


def _parse_rules_container_from_dict(data: Dict[str, Any]) -> DecodedRulesContainer:
    """Parse rules container from dictionary."""
    users = []
    for user_data in data.get("users", []):
        users.append(
            RuleUser(
                id=user_data.get("id"),
                name=user_data.get("name"),
                public_key_pem=(
                    user_data.get("publicKeyPem")
                    or user_data.get("publicKey")
                    or user_data.get("public_key_pem")
                ),
                roles=user_data.get("roles", []),
            )
        )

    groups = []
    for group_data in data.get("groups", []):
        groups.append(
            RuleGroup(
                id=group_data.get("id"),
                name=group_data.get("name"),
                user_ids=group_data.get("userIds") or group_data.get("user_ids", []),
            )
        )

    address_whitelisting_rules = []
    for rule_data in data.get(
        "addressWhitelistingRules", data.get("address_whitelisting_rules", [])
    ):
        parallel_thresholds = _parse_sequential_thresholds(
            rule_data.get("parallelThresholds", rule_data.get("parallel_thresholds", []))
        )
        lines = _parse_address_whitelisting_lines(
            rule_data.get("lines", [])
        )
        address_whitelisting_rules.append(
            AddressWhitelistingRules(
                currency=rule_data.get("currency"),
                network=rule_data.get("network"),
                parallel_thresholds=parallel_thresholds,
                lines=lines,
                include_network_in_payload=rule_data.get(
                    "includeNetworkInPayload",
                    rule_data.get("include_network_in_payload", False),
                ),
            )
        )

    contract_address_whitelisting_rules = []
    for rule_data in data.get(
        "contractAddressWhitelistingRules",
        data.get("contract_address_whitelisting_rules", []),
    ):
        parallel_thresholds = _parse_sequential_thresholds(
            rule_data.get("parallelThresholds", rule_data.get("parallel_thresholds", []))
        )
        contract_address_whitelisting_rules.append(
            ContractAddressWhitelistingRules(
                blockchain=rule_data.get("blockchain"),
                network=rule_data.get("network"),
                parallel_thresholds=parallel_thresholds,
            )
        )

    # Parse transaction rules
    transaction_rules = _parse_transaction_rules(
        data.get("transactionRules", data.get("transaction_rules", []))
    )

    return DecodedRulesContainer(
        users=users,
        groups=groups,
        minimum_distinct_user_signatures=data.get(
            "minimumDistinctUserSignatures",
            data.get("minimum_distinct_user_signatures", 0),
        ),
        minimum_distinct_group_signatures=data.get(
            "minimumDistinctGroupSignatures",
            data.get("minimum_distinct_group_signatures", 0),
        ),
        transaction_rules=transaction_rules,
        address_whitelisting_rules=address_whitelisting_rules,
        contract_address_whitelisting_rules=contract_address_whitelisting_rules,
        enforced_rules_hash=data.get("enforcedRulesHash", data.get("enforced_rules_hash")),
        timestamp=data.get("timestamp", 0),
        hsm_slot_id=data.get("hsmSlotId", data.get("hsm_slot_id", 0)),
    )


def _parse_transaction_rules(
    data: List[Dict[str, Any]],
) -> List[TransactionRules]:
    """Parse transaction rules from list of dicts."""
    rules = []
    for item in data:
        columns = []
        for col_data in item.get("columns", []):
            columns.append(RuleColumn(type=col_data.get("type")))

        lines = []
        for line_data in item.get("lines", []):
            cells = line_data.get("cells", [])
            parallel_thresholds = _parse_sequential_thresholds(
                line_data.get("parallelThresholds", line_data.get("parallel_thresholds", []))
            )
            lines.append(RuleLine(cells=cells, parallel_thresholds=parallel_thresholds))

        details = None
        details_data = item.get("details")
        if details_data:
            details = TransactionRuleDetails(
                domain=details_data.get("domain"),
                sub_domain=details_data.get("subDomain") or details_data.get("sub_domain"),
            )

        rules.append(
            TransactionRules(
                key=item.get("key"),
                columns=columns,
                lines=lines,
                details=details,
            )
        )
    return rules


def _parse_group_thresholds(data: List[Dict[str, Any]]) -> List[GroupThreshold]:
    """Parse group thresholds from list of dicts."""
    thresholds = []
    for item in data:
        thresholds.append(
            GroupThreshold(
                group_id=item.get("groupId") or item.get("group_id"),
                minimum_signatures=item.get("minimumSignatures", item.get("minimum_signatures", 0)),
                threshold=item.get("threshold", 0),
            )
        )
    return thresholds


def _parse_sequential_thresholds(data: List[Dict[str, Any]]) -> List[SequentialThresholds]:
    """Parse sequential thresholds from list of dicts.

    Handles two formats:
    1. Proper nested: [{"thresholds": [{"groupId": ..., "minimumSignatures": ...}]}]
    2. Flat (legacy JSON): [{"groupId": ..., "minimumSignatures": ...}]
       In this case, each item is treated as a single GroupThreshold wrapped in SequentialThresholds.
    """
    result = []
    for item in data:
        if "thresholds" in item:
            # Proper nested format
            thresholds_data = item["thresholds"]
            thresholds = _parse_group_thresholds(thresholds_data)
            result.append(SequentialThresholds(thresholds=thresholds))
        elif "groupId" in item or "group_id" in item:
            # Flat format - wrap single group threshold in SequentialThresholds
            gt = GroupThreshold(
                group_id=item.get("groupId") or item.get("group_id"),
                minimum_signatures=item.get("minimumSignatures", item.get("minimum_signatures", 0)),
                threshold=item.get("threshold", 0),
            )
            result.append(SequentialThresholds(thresholds=[gt]))
        else:
            # Unknown format - create empty SequentialThresholds
            result.append(SequentialThresholds(thresholds=[]))
    return result


def _parse_address_whitelisting_lines(
    data: List[Dict[str, Any]],
) -> List[AddressWhitelistingLine]:
    """Parse address whitelisting lines from list of dicts."""
    lines = []
    for item in data:
        cells_data = item.get("cells", [])
        cells = []
        for cell_data in cells_data:
            internal_wallet = None
            iw_data = cell_data.get("internalWallet") or cell_data.get("internal_wallet")
            if iw_data:
                internal_wallet = RuleSourceInternalWallet(path=iw_data.get("path"))
            cells.append(
                RuleSource(
                    type=cell_data.get("type", 0),
                    internal_wallet=internal_wallet,
                )
            )
        parallel_thresholds = _parse_sequential_thresholds(
            item.get("parallelThresholds", item.get("parallel_thresholds", []))
        )
        lines.append(
            AddressWhitelistingLine(
                cells=cells,
                parallel_thresholds=parallel_thresholds,
            )
        )
    return lines


def _try_protobuf_decode_signatures(data: bytes) -> Optional[List[RuleUserSignature]]:
    """
    Attempt to decode user signatures from protobuf bytes.

    Returns None if protobuf decoding fails or is not available.
    """
    try:
        from taurus_protect._internal.proto import request_reply_pb2

        pb_sigs = request_reply_pb2.UserSignatures()
        pb_sigs.ParseFromString(data)

        signatures = []
        for sig in pb_sigs.signatures:
            signatures.append(
                RuleUserSignature(
                    user_id=sig.userId,
                    signature=base64.b64encode(sig.signature).decode("ascii"),
                )
            )
        return signatures
    except ImportError as e:
        _logger.debug("Protobuf import failed for signatures (using JSON fallback): %s", e)
        return None
    except Exception as e:
        _logger.debug("Protobuf parsing failed for signatures (using JSON fallback): %s", e)
        return None


def user_signatures_from_base64(base64_data: str) -> List[RuleUserSignature]:
    """
    Decode base64-encoded user signatures.

    This function attempts to decode signatures using the following priority:
    1. Protobuf decoding (if available and data is valid protobuf)
    2. JSON decoding (fallback for JSON-encoded signatures)
    3. Empty list (if all parsing fails)

    Args:
        base64_data: Base64-encoded user signatures.

    Returns:
        List of user signatures.
    """
    if not base64_data:
        return []

    try:
        decoded = base64.b64decode(base64_data)

        # Try protobuf first
        result = _try_protobuf_decode_signatures(decoded)
        if result is not None:
            return result

        # Fall back to JSON parsing
        try:
            data = json.loads(decoded.decode("utf-8"))
            signatures = []
            sig_list = data if isinstance(data, list) else data.get("signatures", [])
            for sig_data in sig_list:
                signatures.append(
                    RuleUserSignature(
                        user_id=sig_data.get("userId") or sig_data.get("user_id"),
                        signature=sig_data.get("signature"),
                    )
                )
            return signatures
        except (json.JSONDecodeError, UnicodeDecodeError):
            # Neither protobuf nor JSON worked
            _logger.warning("Failed to decode signatures: not valid protobuf or JSON")
            return []
    except Exception as e:
        _logger.warning("Failed to decode signatures from base64: %s", e)
        return []
