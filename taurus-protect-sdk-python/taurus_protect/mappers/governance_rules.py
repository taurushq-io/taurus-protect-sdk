"""Governance rules mapper utilities."""

from __future__ import annotations

import base64
import json
import logging
from typing import Any, Dict, List, Optional

from google.protobuf.message import DecodeError

from taurus_protect._strict_base64 import strict_b64decode
from taurus_protect.errors import IntegrityError
from taurus_protect.models.governance_rules import (
    RULE_SOURCE_TYPE_ANY,
    RULE_SOURCE_TYPE_ANY_EXCHANGE,
    RULE_SOURCE_TYPE_EXCHANGE,
    RULE_SOURCE_TYPE_EXTERNAL_ADDRESS,
    RULE_SOURCE_TYPE_INTERNAL_ADDRESS,
    RULE_SOURCE_TYPE_INTERNAL_WALLET,
    AddressWhitelistingLine,
    AddressWhitelistingRules,
    CashSettlement,
    ContractAddressWhitelistingRules,
    CosmosDetails,
    DecodedRulesContainer,
    EvmCallContract,
    GroupThreshold,
    RuleColumn,
    RuleGroup,
    RuleLine,
    RuleSource,
    RuleSourceExchange,
    RuleSourceExternalAddress,
    RuleSourceInternalAddress,
    RuleSourceInternalWallet,
    RuleUser,
    RuleUserSignature,
    SequentialThresholds,
    TransactionRuleDetails,
    TransactionRules,
    XtzCallContract,
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
    except DecodeError as e:
        # The payload is not protobuf at all — fall through to the JSON decoder.
        _logger.warning("Protobuf parsing failed: %s", e)
        return None


def _enum_name(wrapper: Any, value: int) -> str:
    """Enum number -> value name; decimal string for values unknown to this SDK (passthrough)."""
    try:
        return wrapper.Name(value)
    except ValueError:
        return str(value)


def _rules_container_from_proto(pb: Any) -> DecodedRulesContainer:
    """Convert a protobuf RulesContainer to the model (lossless)."""
    from taurus_protect._internal.proto import request_reply_pb2
    from taurus_protect.mappers.rules_container_encode import capture_unknown

    users = []
    for u in pb.users:
        roles = [_enum_name(request_reply_pb2.Role, role) for role in u.roles]
        users.append(
            RuleUser(
                id=u.id,
                name=getattr(u, "name", None),
                public_key_pem=u.publicKey,
                roles=roles,
                properties=dict(u.properties),
                unknown_fields=capture_unknown(u),
            )
        )

    groups = []
    for g in pb.groups:
        groups.append(
            RuleGroup(
                id=g.id,
                name=getattr(g, "name", None),
                user_ids=list(g.userIds),
                properties=dict(g.properties),
                unknown_fields=capture_unknown(g),
            )
        )

    address_whitelisting_rules = []
    for r in pb.addressWhitelistingRules:
        address_whitelisting_rules.append(
            AddressWhitelistingRules(
                currency=r.currency,
                network=r.network,
                parallel_thresholds=[_sequential_thresholds_from_proto(pt) for pt in r.parallelThresholds],
                lines=[_address_whitelisting_line_from_proto(line) for line in r.lines],
                properties=dict(r.properties),
                unknown_fields=capture_unknown(r),
            )
        )

    contract_address_whitelisting_rules = []
    for r in pb.contractAddressWhitelistingRules:
        contract_address_whitelisting_rules.append(
            ContractAddressWhitelistingRules(
                blockchain=_enum_name(request_reply_pb2.Blockchain, r.blockchain),
                network=r.network,
                parallel_thresholds=[_sequential_thresholds_from_proto(pt) for pt in r.parallelThresholds],
                properties=dict(r.properties),
                unknown_fields=capture_unknown(r),
            )
        )

    transaction_rules = [_transaction_rules_from_proto(tr) for tr in pb.transactionRules]

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
        minimum_commitment_signatures=pb.minimumCommitmentSignatures,
        engine_identities=list(pb.engineIdentities),
        hsm_slot_id=pb.hsmSlotId,
        properties=dict(pb.properties),
        unknown_fields=capture_unknown(pb),
    )


def _transaction_rules_from_proto(pb_tr: Any) -> TransactionRules:
    """Convert protobuf TransactionRules to model (typed cells)."""
    from taurus_protect._internal.proto import request_reply_pb2
    from taurus_protect.mappers.rule_cell_codec import rule_cell_from_bytes
    from taurus_protect.mappers.rules_container_encode import capture_unknown

    _CT = request_reply_pb2.RulesContainer.ColumnType
    columns = []
    for col in pb_tr.columns:
        columns.append(
            RuleColumn(
                type=_enum_name(_CT, col.type),
                name=col.name,
                metadata_key=col.metadataKey,
                unknown_fields=capture_unknown(col),
            )
        )

    lines = []
    for line in pb_tr.lines:
        cells = []
        for i, cell_bytes in enumerate(line.cells):
            col_type = columns[i].type if i < len(columns) else ""
            cells.append(rule_cell_from_bytes(col_type or "", cell_bytes))
        lines.append(
            RuleLine(
                cells=cells,
                parallel_thresholds=[_sequential_thresholds_from_proto(pt) for pt in line.parallelThresholds],
                priority=line.priority,
                properties=dict(line.properties),
                unknown_fields=capture_unknown(line),
            )
        )

    details = None
    if pb_tr.HasField("details"):
        d = pb_tr.details
        _TRD = request_reply_pb2.RulesContainer.TransactionRules.TransactionRuleDetails
        evm = None
        if d.HasField("evmCallContract"):
            evm = EvmCallContract(contract_type=d.evmCallContract.contractType,
                                  method_signature=d.evmCallContract.methodSignature,
                                  unknown_fields=capture_unknown(d.evmCallContract))
        xtz = None
        if d.HasField("xtzCallContract"):
            xtz = XtzCallContract(contract_type=d.xtzCallContract.contractType,
                                  method_signature=d.xtzCallContract.methodSignature,
                                  unknown_fields=capture_unknown(d.xtzCallContract))
        cash = None
        if d.HasField("cashSettlement"):
            cash = CashSettlement(provider=d.cashSettlement.provider,
                                  request_type=d.cashSettlement.requestType,
                                  unknown_fields=capture_unknown(d.cashSettlement))
        cosmos = None
        if d.HasField("cosmosDetails"):
            cosmos = CosmosDetails(method_signatures=list(d.cosmosDetails.methodSignatures),
                                   unknown_fields=capture_unknown(d.cosmosDetails))
        details = TransactionRuleDetails(
            domain=_enum_name(_TRD.RuleDomain, d.domain),
            sub_domain=_enum_name(_TRD.RuleSubDomain, d.subDomain),
            blockchain=d.blockchain,
            network=d.network,
            evm_call_contract=evm,
            xtz_call_contract=xtz,
            cash_settlement=cash,
            cosmos_details=cosmos,
            unknown_fields=capture_unknown(d),
        )

    return TransactionRules(
        key=pb_tr.key,
        columns=columns,
        lines=lines,
        details=details,
        unknown_fields=capture_unknown(pb_tr),
    )


def _address_whitelisting_line_from_proto(pb_line: Any) -> AddressWhitelistingLine:
    """Convert protobuf AddressWhitelistingRules.Line to model."""
    from taurus_protect.mappers.rules_container_encode import capture_unknown

    cells = [_rule_source_from_bytes(cell_bytes) for cell_bytes in pb_line.cells]
    return AddressWhitelistingLine(
        cells=cells,
        parallel_thresholds=[_sequential_thresholds_from_proto(pt) for pt in pb_line.parallelThresholds],
        properties=dict(pb_line.properties),
        unknown_fields=capture_unknown(pb_line),
    )


def _rule_source_from_bytes(data: bytes) -> RuleSource:
    """Decode a RuleSource, keeping anything not reproducible byte-for-byte verbatim.

    Mirrors the cell codec's lossless guard. An unknown-field check alone is too
    weak: a non-canonical encoding carries no unknown field yet still re-encodes
    differently. An explicitly-present zero-length payload (``08 01 12 00``) decodes
    to the same typed value as an absent one (``08 01``), so re-emitting the typed
    form would drop two bytes from a container the SuperAdmins signed.
    """
    typed = _rule_source_from_bytes_typed(data)
    if typed is None:
        return RuleSource(type=0, raw=data)

    from taurus_protect.mappers.rules_container_encode import _rule_source_to_bytes

    try:
        if _rule_source_to_bytes(typed) != data:
            return RuleSource(type=0, raw=data)
    except (DecodeError, ValueError, TypeError):
        return RuleSource(type=0, raw=data)
    return typed


def _rule_source_from_bytes_typed(data: bytes) -> Optional[RuleSource]:
    """Decode a RuleSource into its typed form, or None when it cannot be typed."""
    from taurus_protect._internal.proto import request_reply_pb2

    try:
        pb_source = request_reply_pb2.RuleSource()
        pb_source.ParseFromString(data)
    except DecodeError:
        return None

    from taurus_protect.mappers.rules_container_encode import capture_unknown

    # A schema-newer field on the source itself cannot be represented by the typed
    # model, so keep the bytes verbatim rather than dropping it on the next encode.
    if capture_unknown(pb_source):
        return None

    t = int(pb_source.type)
    try:
        if t == RULE_SOURCE_TYPE_INTERNAL_WALLET and pb_source.payload:
            m = request_reply_pb2.RuleSourceInternalWallet.FromString(pb_source.payload)
            if capture_unknown(m):
                return None
            return RuleSource(type=t, internal_wallet=RuleSourceInternalWallet(path=m.path))
        if t == RULE_SOURCE_TYPE_INTERNAL_ADDRESS and pb_source.payload:
            m = request_reply_pb2.RuleSourceInternalAddress.FromString(pb_source.payload)
            if capture_unknown(m):
                return None
            return RuleSource(type=t, internal_address=RuleSourceInternalAddress(address=m.address, path=m.path))
        if t == RULE_SOURCE_TYPE_EXCHANGE and pb_source.payload:
            m = request_reply_pb2.RuleSourceExchange.FromString(pb_source.payload)
            if capture_unknown(m):
                return None
            return RuleSource(type=t, exchange=RuleSourceExchange(label=m.label))
        if t == RULE_SOURCE_TYPE_EXTERNAL_ADDRESS and pb_source.payload:
            m = request_reply_pb2.RuleSourceExternalAddress.FromString(pb_source.payload)
            if capture_unknown(m):
                return None
            return RuleSource(type=t, external_address=RuleSourceExternalAddress(address=m.address, memo=m.memo))
    except DecodeError:
        return None

    # Payload-less arms: a present payload is data this SDK would discard.
    if t in (RULE_SOURCE_TYPE_ANY, RULE_SOURCE_TYPE_ANY_EXCHANGE) and pb_source.payload:
        return None

    if t in (RULE_SOURCE_TYPE_INTERNAL_WALLET, RULE_SOURCE_TYPE_INTERNAL_ADDRESS,
             RULE_SOURCE_TYPE_EXCHANGE, RULE_SOURCE_TYPE_EXTERNAL_ADDRESS,
             RULE_SOURCE_TYPE_ANY, RULE_SOURCE_TYPE_ANY_EXCHANGE):
        return RuleSource(type=t)
    # source type newer than this SDK: preserve verbatim
    return None


def _sequential_thresholds_from_proto(pb: Any) -> SequentialThresholds:
    """Convert protobuf SequentialThresholds to model."""
    from taurus_protect.mappers.rules_container_encode import capture_unknown

    thresholds = [
        GroupThreshold(
            group_id=t.groupId,
            minimum_signatures=t.minimumSignatures,
            unknown_fields=capture_unknown(t),
        )
        for t in pb.thresholds
    ]
    return SequentialThresholds(thresholds=thresholds, unknown_fields=capture_unknown(pb))


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
        decoded = strict_b64decode(base64_data)

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
            # JSON fallback (used only when protobuf decode fails): transaction-rule
            # cells are typed and column-scoped, which the JSON shape does not carry,
            # so leave them empty here. The protobuf path is the lossless one.
            parallel_thresholds = _parse_sequential_thresholds(
                line_data.get("parallelThresholds", line_data.get("parallel_thresholds", []))
            )
            lines.append(RuleLine(cells=[], parallel_thresholds=parallel_thresholds))

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
        decoded = strict_b64decode(base64_data)

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
