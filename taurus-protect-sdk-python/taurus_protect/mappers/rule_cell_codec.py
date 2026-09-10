"""Codec between the typed RuleCell union and transaction-rule cell wire bytes.

Mirrors the Go/TypeScript SDKs. A cell is a wrapper message selected by the
column's type, ``{type: <cell-type enum>, payload: <bytes>}``. Message-family
payloads nest a serialized sub-message; RuleStringEqual/RuleBytesEqual carry
raw scalar bytes; the integer families carry big-endian magnitude bytes with
the sign expressed by the enum arm. The ``*Any`` cells are the protobuf zero
values, so their serialized form is the empty cell.

Decode never fails: any cell that does not round-trip byte-identically through
the typed layer (unknown column type, unknown cell type, unknown protobuf
sub-fields) is returned as a :class:`RawCell` so nothing is silently dropped.
"""

from __future__ import annotations

from typing import Optional

from google.protobuf.message import DecodeError

from taurus_protect._internal.proto import request_reply_pb2 as pb
from taurus_protect.models import rule_cell as rc
from taurus_protect.models.rule_cell import RuleCell, cell_family

_RS = pb.RuleSource
_RD = pb.RuleDestination
_FA = pb.RuleFiatAmount
_SE = pb.RuleStringEqual
_BE = pb.RuleBytesEqual
_SAE = pb.RuleStringArrayEqual
_IG = pb.RuleIntegerGreater
_UIG = pb.RuleUIntegerGreater
_WC = pb.RuleWhitelistedContract


def _magnitude(value: int) -> bytes:
    v = abs(value)
    if v == 0:
        return b"\x00"
    return v.to_bytes((v.bit_length() + 7) // 8, "big")


def _from_magnitude(data: bytes) -> int:
    return int.from_bytes(data, "big")


def _blockchain_to_int(name: str) -> int:
    try:
        return pb.Blockchain.Value(name)
    except ValueError:
        return int(name)  # numeric passthrough for values newer than this SDK


def _blockchain_name(value: int) -> str:
    try:
        return pb.Blockchain.Name(value)
    except ValueError:
        return str(value)


def _wrap(msg) -> bytes:
    return msg.SerializeToString()


def rule_cell_to_bytes(col_type: str, cell: RuleCell) -> bytes:
    """Encode a typed cell to the wrapped protobuf bytes for a line cell."""
    if isinstance(cell, rc.RawCell):
        return cell.payload

    family = cell_family(cell)
    if col_type and col_type != family:
        raise ValueError(f"cell {cell.kind} is not valid for column type {col_type!r}")

    k = cell.kind
    # RuleFiatAmount
    if k == "FiatAmountAny":
        return _wrap(_FA(type=_FA.RuleFiatAmountAny))
    if k == "FiatAmountIsZero":
        return _wrap(_FA(type=_FA.RuleFiatAmountIsZero))
    if k == "FiatAmountRange":
        inner = _wrap(pb.RuleFiatAmountRange(minAmount=cell.min_amount, maxAmount=cell.max_amount))
        return _wrap(_FA(type=_FA.RuleFiatAmountRange, payload=inner))

    # RuleSource
    if k == "SourceAny":
        return _wrap(_RS(type=_RS.RuleSourceAny))
    if k == "SourceAnyExchange":
        return _wrap(_RS(type=_RS.RuleSourceAnyExchange))
    if k == "SourceInternalWallet":
        inner = _wrap(pb.RuleSourceInternalWallet(path=cell.path))
        return _wrap(_RS(type=_RS.RuleSourceInternalWallet, payload=inner))
    if k == "SourceInternalAddress":
        inner = _wrap(pb.RuleSourceInternalAddress(address=cell.address, path=cell.path))
        return _wrap(_RS(type=_RS.RuleSourceInternalAddress, payload=inner))
    if k == "SourceExchange":
        inner = _wrap(pb.RuleSourceExchange(label=cell.label))
        return _wrap(_RS(type=_RS.RuleSourceExchange, payload=inner))
    if k == "SourceExternalAddress":
        inner = _wrap(pb.RuleSourceExternalAddress(address=cell.address, memo=cell.memo))
        return _wrap(_RS(type=_RS.RuleSourceExternalAddress, payload=inner))

    # RuleDestination
    if k == "DestinationAny":
        return _wrap(_RD(type=_RD.RuleDestinationAny))
    if k == "DestinationAnyExchange":
        return _wrap(_RD(type=_RD.RuleDestinationAnyExchange))
    if k == "DestinationAnyExternalAddress":
        return _wrap(_RD(type=_RD.RuleDestinationAnyExternalAddress))
    if k == "DestinationAnyContractAddress":
        return _wrap(_RD(type=_RD.RuleDestinationAnyContractAddress))
    if k == "DestinationInternalWallet":
        inner = _wrap(pb.RuleDestinationInternalWallet(path=cell.path))
        return _wrap(_RD(type=_RD.RuleDestinationInternalWallet, payload=inner))
    if k == "DestinationInternalAddress":
        inner = _wrap(pb.RuleDestinationInternalAddress(address=cell.address, path=cell.path))
        return _wrap(_RD(type=_RD.RuleDestinationInternalAddress, payload=inner))
    if k == "DestinationExternalAddress":
        inner = _wrap(pb.RuleDestinationExternalAddress(address=cell.address, memo=cell.memo))
        return _wrap(_RD(type=_RD.RuleDestinationExternalAddress, payload=inner))
    if k == "DestinationExchange":
        inner = _wrap(pb.RuleDestinationExchange(label=cell.label, memo=cell.memo))
        return _wrap(_RD(type=_RD.RuleDestinationExchange, payload=inner))
    if k == "DestinationContractAddress":
        inner = _wrap(pb.RuleDestinationContractAddress(
            address=cell.address, name=cell.name, symbol=cell.symbol,
            blockchain=_blockchain_to_int(cell.blockchain)))
        return _wrap(_RD(type=_RD.RuleDestinationContractAddress, payload=inner))

    # RuleWhitelistedContract
    if k == "WhitelistedContractAny":
        return _wrap(_WC(type=_WC.RuleWhitelistedContractAny))
    if k == "WhitelistedContractAddress":
        inner = _wrap(pb.RuleDestinationContractAddress(
            address=cell.address, name=cell.name, symbol=cell.symbol,
            blockchain=_blockchain_to_int(cell.blockchain)))
        return _wrap(_WC(type=_WC.RuleWhitelistedContract_RuleDestinationContractAddress, payload=inner))

    # RuleStringEqual (payload = raw UTF-8 string bytes)
    if k == "StringEqualAny":
        return _wrap(_SE(type=_SE.RuleStringEqualAny))
    if k == "StringEqualEmpty":
        return _wrap(_SE(type=_SE.RuleStringEqualEmpty))
    if k == "StringEqualValue":
        return _wrap(_SE(type=_SE.RuleStringEqualValue, payload=cell.value.encode("utf-8")))

    # RuleBytesEqual (payload = raw bytes)
    if k == "BytesEqualAny":
        return _wrap(_BE(type=_BE.RuleBytesEqualAny))
    if k == "BytesEqualEmpty":
        return _wrap(_BE(type=_BE.RuleBytesEqualEmpty))
    if k == "BytesEqualValue":
        return _wrap(_BE(type=_BE.RuleBytesEqualValue, payload=cell.value))

    # RuleStringArrayEqual
    if k == "StringArrayEqualAny":
        return _wrap(_SAE(type=_SAE.RuleStringArrayEqualAny))
    if k == "StringArrayEqualEmpty":
        return _wrap(_SAE(type=_SAE.RuleStringArrayEqualEmpty))
    if k == "StringArrayEqualValue":
        inner = _wrap(pb.RuleStringArrayEqualValue(values=cell.values))
        return _wrap(_SAE(type=_SAE.RuleStringArrayEqualValue, payload=inner))

    # RuleIntegerGreater (magnitude payload, sign selects the arm)
    if k == "IntegerGreaterAny":
        return _wrap(_IG(type=_IG.RuleIntegerGreaterAny))
    if k == "IntegerGreaterValue":
        arm = _IG.RuleIntegerGreaterNegValue if cell.value < 0 else _IG.RuleIntegerGreaterValue
        return _wrap(_IG(type=arm, payload=_magnitude(cell.value)))

    # RuleUIntegerGreater
    if k == "UIntegerGreaterAny":
        return _wrap(_UIG(type=_UIG.RuleUIntegerGreaterAny))
    if k == "UIntegerGreaterIsZero":
        return _wrap(_UIG(type=_UIG.RuleUIntegerGreaterIsZero))
    if k == "UIntegerGreaterValue":
        if cell.value < 0:
            raise ValueError("UIntegerGreaterValue requires a non-negative value")
        return _wrap(_UIG(type=_UIG.RuleUIntegerGreaterValue, payload=_magnitude(cell.value)))
    if k == "UIntegerGreaterIsEqual":
        if cell.value < 0:
            raise ValueError("UIntegerGreaterIsEqual requires a non-negative value")
        return _wrap(_UIG(type=_UIG.RuleUIntegerGreaterIsEqual, payload=_magnitude(cell.value)))

    raise ValueError(f"unsupported rule cell kind {k!r}")


def rule_cell_from_bytes(col_type: str, data: bytes) -> RuleCell:
    """Decode a line cell. Falls back to RawCell if it does not round-trip byte-identically.

    Decoding never fails: one malformed cell must degrade to a RawCell rather than
    abort the whole container, which would take every rule for that tenant down with
    it. A cell payload is protobuf ``bytes``, so a non-UTF-8 string cell and a
    truncated wrapper are both legal on the wire.
    """
    raw = rc.RawCell(column_type=col_type, payload=data)
    try:
        typed = _rule_cell_from_bytes_typed(col_type, data)
        if typed is None:
            return raw
        # Lossless guard: preserve verbatim if the typed value does not re-encode exactly.
        if rule_cell_to_bytes(col_type, typed) != data:
            return raw
    except (DecodeError, UnicodeDecodeError, ValueError):
        return raw
    return typed


def _rule_cell_from_bytes_typed(col_type: str, data: bytes) -> Optional[RuleCell]:
    if col_type == "RuleFiatAmount":
        w = _FA.FromString(data)
        if w.type == _FA.RuleFiatAmountAny:
            return rc.FiatAmountAny()
        if w.type == _FA.RuleFiatAmountIsZero:
            return rc.FiatAmountIsZero()
        if w.type == _FA.RuleFiatAmountRange:
            i = pb.RuleFiatAmountRange.FromString(w.payload)
            return rc.FiatAmountRange(min_amount=i.minAmount, max_amount=i.maxAmount)
        return None
    if col_type == "RuleSource":
        w = _RS.FromString(data)
        if w.type == _RS.RuleSourceAny:
            return rc.SourceAny()
        if w.type == _RS.RuleSourceAnyExchange:
            return rc.SourceAnyExchange()
        if w.type == _RS.RuleSourceInternalWallet:
            i = pb.RuleSourceInternalWallet.FromString(w.payload)
            return rc.SourceInternalWallet(path=i.path)
        if w.type == _RS.RuleSourceInternalAddress:
            i = pb.RuleSourceInternalAddress.FromString(w.payload)
            return rc.SourceInternalAddress(address=i.address, path=i.path)
        if w.type == _RS.RuleSourceExchange:
            i = pb.RuleSourceExchange.FromString(w.payload)
            return rc.SourceExchange(label=i.label)
        if w.type == _RS.RuleSourceExternalAddress:
            i = pb.RuleSourceExternalAddress.FromString(w.payload)
            return rc.SourceExternalAddress(address=i.address, memo=i.memo)
        return None
    if col_type == "RuleDestination":
        w = _RD.FromString(data)
        if w.type == _RD.RuleDestinationAny:
            return rc.DestinationAny()
        if w.type == _RD.RuleDestinationAnyExchange:
            return rc.DestinationAnyExchange()
        if w.type == _RD.RuleDestinationAnyExternalAddress:
            return rc.DestinationAnyExternalAddress()
        if w.type == _RD.RuleDestinationAnyContractAddress:
            return rc.DestinationAnyContractAddress()
        if w.type == _RD.RuleDestinationInternalWallet:
            i = pb.RuleDestinationInternalWallet.FromString(w.payload)
            return rc.DestinationInternalWallet(path=i.path)
        if w.type == _RD.RuleDestinationInternalAddress:
            i = pb.RuleDestinationInternalAddress.FromString(w.payload)
            return rc.DestinationInternalAddress(address=i.address, path=i.path)
        if w.type == _RD.RuleDestinationExternalAddress:
            i = pb.RuleDestinationExternalAddress.FromString(w.payload)
            return rc.DestinationExternalAddress(address=i.address, memo=i.memo)
        if w.type == _RD.RuleDestinationExchange:
            i = pb.RuleDestinationExchange.FromString(w.payload)
            return rc.DestinationExchange(label=i.label, memo=i.memo)
        if w.type == _RD.RuleDestinationContractAddress:
            i = pb.RuleDestinationContractAddress.FromString(w.payload)
            return rc.DestinationContractAddress(
                address=i.address, name=i.name, symbol=i.symbol, blockchain=_blockchain_name(i.blockchain))
        return None
    if col_type == "RuleWhitelistedContract":
        w = _WC.FromString(data)
        if w.type == _WC.RuleWhitelistedContractAny:
            return rc.WhitelistedContractAny()
        if w.type == _WC.RuleWhitelistedContract_RuleDestinationContractAddress:
            i = pb.RuleDestinationContractAddress.FromString(w.payload)
            return rc.WhitelistedContractAddress(
                address=i.address, name=i.name, symbol=i.symbol, blockchain=_blockchain_name(i.blockchain))
        return None
    if col_type == "RuleStringEqual":
        w = _SE.FromString(data)
        if w.type == _SE.RuleStringEqualAny:
            return rc.StringEqualAny()
        if w.type == _SE.RuleStringEqualEmpty:
            return rc.StringEqualEmpty()
        if w.type == _SE.RuleStringEqualValue:
            return rc.StringEqualValue(value=w.payload.decode("utf-8"))
        return None
    if col_type == "RuleBytesEqual":
        w = _BE.FromString(data)
        if w.type == _BE.RuleBytesEqualAny:
            return rc.BytesEqualAny()
        if w.type == _BE.RuleBytesEqualEmpty:
            return rc.BytesEqualEmpty()
        if w.type == _BE.RuleBytesEqualValue:
            return rc.BytesEqualValue(value=w.payload)
        return None
    if col_type == "RuleStringArrayEqual":
        w = _SAE.FromString(data)
        if w.type == _SAE.RuleStringArrayEqualAny:
            return rc.StringArrayEqualAny()
        if w.type == _SAE.RuleStringArrayEqualEmpty:
            return rc.StringArrayEqualEmpty()
        if w.type == _SAE.RuleStringArrayEqualValue:
            i = pb.RuleStringArrayEqualValue.FromString(w.payload)
            return rc.StringArrayEqualValue(values=list(i.values))
        return None
    if col_type == "RuleIntegerGreater":
        w = _IG.FromString(data)
        if w.type == _IG.RuleIntegerGreaterAny:
            return rc.IntegerGreaterAny()
        if w.type == _IG.RuleIntegerGreaterValue:
            return rc.IntegerGreaterValue(value=_from_magnitude(w.payload))
        if w.type == _IG.RuleIntegerGreaterNegValue:
            return rc.IntegerGreaterValue(value=-_from_magnitude(w.payload))
        return None
    if col_type == "RuleUIntegerGreater":
        w = _UIG.FromString(data)
        if w.type == _UIG.RuleUIntegerGreaterAny:
            return rc.UIntegerGreaterAny()
        if w.type == _UIG.RuleUIntegerGreaterIsZero:
            return rc.UIntegerGreaterIsZero()
        if w.type == _UIG.RuleUIntegerGreaterValue:
            return rc.UIntegerGreaterValue(value=_from_magnitude(w.payload))
        if w.type == _UIG.RuleUIntegerGreaterIsEqual:
            return rc.UIntegerGreaterIsEqual(value=_from_magnitude(w.payload))
        return None
    # Unknown column type (incl. RuleAny): no cell family to type it.
    return None
