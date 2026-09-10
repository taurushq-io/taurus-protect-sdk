"""Model -> protobuf encoding for RulesContainer (mirrors the Go SDK).

Inverse of the decode in ``governance_rules.py``. Server-controlled
``enforcedRulesHash`` and ``timestamp`` are stripped (the server recomputes
them and rejects submissions asserting stale values). Per-node unknown
protobuf fields captured at decode are re-attached, and unknown enum values
pass through numerically, so a container from a newer schema re-encodes
without dropping data. Enum names the caller authored that this SDK does not
recognize are a hard error.
"""

from __future__ import annotations

import base64
from typing import Any, List

from taurus_protect._internal.proto import request_reply_pb2 as pb
from taurus_protect.mappers.rule_cell_codec import rule_cell_to_bytes
from taurus_protect.models.governance_rules import (
    AddressWhitelistingRules,
    ContractAddressWhitelistingRules,
    DecodedRulesContainer,
    RuleColumn,
    RuleSource,
    SequentialThresholds,
    TransactionRuleDetails,
    TransactionRules,
    RuleUser,
    RuleGroup,
)

_TRD = pb.RulesContainer.TransactionRules.TransactionRuleDetails


def capture_unknown(msg: Any) -> bytes:
    """Return the message's unknown-field bytes only (empty if none).

    upb does not expose UnknownFieldSet iteration, but retains unknowns across
    parse/serialize: reparse a clone, clear all known fields, and serialize what
    remains.
    """
    clone = type(msg).FromString(msg.SerializeToString())
    for field in list(clone.DESCRIPTOR.fields):
        clone.ClearField(field.name)
    return clone.SerializeToString()


def _reattach(msg: Any, unknown: bytes) -> None:
    """Re-attach previously captured unknown-field bytes onto a built message."""
    if unknown:
        msg.MergeFromString(unknown)


def _enum_number(enum_wrapper: Any, name, what: str) -> int:
    """Name -> enum number. Empty/None -> 0; decimal strings pass through; unknown names raise."""
    if name is None or name == "":
        return 0
    if isinstance(name, str) and (name.isdigit() or (name.startswith("-") and name[1:].isdigit())):
        return int(name)
    try:
        return enum_wrapper.Value(name)
    except ValueError:
        raise ValueError(f"unknown {what} {name!r}")


def _blockchain_number(name) -> int:
    return _enum_number(pb.Blockchain, name, "blockchain")


def _thresholds_to_proto(items: List[SequentialThresholds]) -> List[Any]:
    out = []
    for st in items:
        pb_st = pb.SequentialThresholds()
        for t in st.thresholds:
            pb_t = pb.GroupThreshold(groupId=t.group_id or "", minimumSignatures=t.minimum_signatures)
            _reattach(pb_t, t.unknown_fields)
            pb_st.thresholds.append(pb_t)
        _reattach(pb_st, st.unknown_fields)
        out.append(pb_st)
    return out


def _user_to_proto(u: RuleUser) -> Any:
    m = pb.User(
        id=u.id or "",
        publicKey=u.public_key_pem or "",
        roles=[_enum_number(pb.Role, r, "user role") for r in u.roles],
        properties=dict(u.properties),
    )
    _reattach(m, u.unknown_fields)
    return m


def _group_to_proto(g: RuleGroup) -> Any:
    m = pb.Group(id=g.id or "", userIds=list(g.user_ids), properties=dict(g.properties))
    _reattach(m, g.unknown_fields)
    return m


def _column_to_proto(c: RuleColumn) -> Any:
    m = pb.RulesContainer.Column(
        type=_enum_number(pb.RulesContainer.ColumnType, c.type, "column type"),
        name=c.name or "",
        metadataKey=c.metadata_key or "",
    )
    _reattach(m, c.unknown_fields)
    return m


def _details_to_proto(d: TransactionRuleDetails) -> Any:
    m = _TRD(
        domain=_enum_number(_TRD.RuleDomain, d.domain, "rule domain"),
        subDomain=_enum_number(_TRD.RuleSubDomain, d.sub_domain, "rule sub-domain"),
        blockchain=d.blockchain or "",
        network=d.network or "",
    )
    if d.evm_call_contract is not None:
        m.evmCallContract.contractType = d.evm_call_contract.contract_type or ""
        m.evmCallContract.methodSignature = d.evm_call_contract.method_signature or ""
        _reattach(m.evmCallContract, d.evm_call_contract.unknown_fields)
    if d.xtz_call_contract is not None:
        m.xtzCallContract.contractType = d.xtz_call_contract.contract_type or ""
        m.xtzCallContract.methodSignature = d.xtz_call_contract.method_signature or ""
        _reattach(m.xtzCallContract, d.xtz_call_contract.unknown_fields)
    if d.cash_settlement is not None:
        m.cashSettlement.provider = d.cash_settlement.provider or ""
        m.cashSettlement.requestType = d.cash_settlement.request_type or ""
        _reattach(m.cashSettlement, d.cash_settlement.unknown_fields)
    if d.cosmos_details is not None:
        m.cosmosDetails.methodSignatures.extend(d.cosmos_details.method_signatures)
        _reattach(m.cosmosDetails, d.cosmos_details.unknown_fields)
    _reattach(m, d.unknown_fields)
    return m


def _tx_rules_to_proto(r: TransactionRules) -> Any:
    m = pb.RulesContainer.TransactionRules(key=r.key or "")
    for c in r.columns:
        m.columns.append(_column_to_proto(c))
    for line in r.lines:
        pb_line = pb.RulesContainer.Line(priority=line.priority, properties=dict(line.properties))
        for i, cell in enumerate(line.cells):
            col_type = r.columns[i].type if i < len(r.columns) else ""
            pb_line.cells.append(rule_cell_to_bytes(col_type or "", cell))
        pb_line.parallelThresholds.extend(_thresholds_to_proto(line.parallel_thresholds))
        _reattach(pb_line, line.unknown_fields)
        m.lines.append(pb_line)
    if r.details is not None:
        m.details.CopyFrom(_details_to_proto(r.details))
    _reattach(m, r.unknown_fields)
    return m


def _rule_source_to_bytes(s: RuleSource) -> bytes:
    if s.raw:
        return s.raw
    m = pb.RuleSource(type=s.type)
    if s.type == pb.RuleSource.RuleSourceInternalWallet and s.internal_wallet is not None:
        m.payload = pb.RuleSourceInternalWallet(path=s.internal_wallet.path or "").SerializeToString()
    elif s.type == pb.RuleSource.RuleSourceInternalAddress and s.internal_address is not None:
        m.payload = pb.RuleSourceInternalAddress(
            address=s.internal_address.address or "", path=s.internal_address.path or "").SerializeToString()
    elif s.type == pb.RuleSource.RuleSourceExchange and s.exchange is not None:
        m.payload = pb.RuleSourceExchange(label=s.exchange.label or "").SerializeToString()
    elif s.type == pb.RuleSource.RuleSourceExternalAddress and s.external_address is not None:
        m.payload = pb.RuleSourceExternalAddress(
            address=s.external_address.address or "", memo=s.external_address.memo or "").SerializeToString()
    return m.SerializeToString(deterministic=True)


def _awr_to_proto(r: AddressWhitelistingRules) -> Any:
    m = pb.RulesContainer.AddressWhitelistingRules(
        currency=r.currency or "", network=r.network or "", properties=dict(r.properties))
    m.parallelThresholds.extend(_thresholds_to_proto(r.parallel_thresholds))
    for line in r.lines:
        pb_line = pb.RulesContainer.AddressWhitelistingRules.Line(properties=dict(line.properties))
        for s in line.cells:
            pb_line.cells.append(_rule_source_to_bytes(s))
        pb_line.parallelThresholds.extend(_thresholds_to_proto(line.parallel_thresholds))
        _reattach(pb_line, line.unknown_fields)
        m.lines.append(pb_line)
    _reattach(m, r.unknown_fields)
    return m


def _cawr_to_proto(r: ContractAddressWhitelistingRules) -> Any:
    m = pb.RulesContainer.ContractAddressWhitelistingRules(
        blockchain=_blockchain_number(r.blockchain),
        network=r.network or "",
        properties=dict(r.properties),
    )
    m.parallelThresholds.extend(_thresholds_to_proto(r.parallel_thresholds))
    _reattach(m, r.unknown_fields)
    return m


def rules_container_to_proto(c: DecodedRulesContainer) -> Any:
    out = pb.RulesContainer(
        minimumDistinctUserSignatures=c.minimum_distinct_user_signatures,
        minimumDistinctGroupSignatures=c.minimum_distinct_group_signatures,
        minimumCommitmentSignatures=c.minimum_commitment_signatures,
        engineIdentities=list(c.engine_identities),
        hsmSlotId=c.hsm_slot_id,
        properties=dict(c.properties),
        # enforcedRulesHash and timestamp intentionally omitted (server-controlled).
    )
    for u in c.users:
        out.users.append(_user_to_proto(u))
    for g in c.groups:
        out.groups.append(_group_to_proto(g))
    for r in c.transaction_rules:
        out.transactionRules.append(_tx_rules_to_proto(r))
    for r in c.address_whitelisting_rules:
        out.addressWhitelistingRules.append(_awr_to_proto(r))
    for r in c.contract_address_whitelisting_rules:
        out.contractAddressWhitelistingRules.append(_cawr_to_proto(r))
    _reattach(out, c.unknown_fields)
    return out


def rules_container_to_bytes(container: DecodedRulesContainer) -> bytes:
    """Encode a DecodedRulesContainer to raw protobuf bytes."""
    if container is None:
        raise ValueError("rules container cannot be None")
    return rules_container_to_proto(container).SerializeToString(deterministic=True)


def rules_container_to_base64(container: DecodedRulesContainer) -> str:
    """Encode a DecodedRulesContainer to the base64 protobuf wire format."""
    return base64.b64encode(rules_container_to_bytes(container)).decode("ascii")
