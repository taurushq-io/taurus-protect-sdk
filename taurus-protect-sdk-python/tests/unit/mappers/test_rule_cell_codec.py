"""Branch/edge coverage for the cell codec that the golden vectors exercise
only transitively: error paths, the RawCell fallback, the magnitude helpers,
blockchain passthrough, and cell_family."""

import pytest

from taurus_protect._internal.proto import request_reply_pb2 as pb
from taurus_protect.models import rule_cell as rc
from taurus_protect.models.rule_cell import cell_family
from taurus_protect.mappers.rule_cell_codec import (
    _blockchain_name,
    _blockchain_to_int,
    _from_magnitude,
    _magnitude,
    rule_cell_from_bytes,
    rule_cell_to_bytes,
)

# --- Error branches ---


def test_encode_wrong_column_family_raises():
    with pytest.raises(ValueError, match="not valid for column type"):
        rule_cell_to_bytes("RuleSource", rc.FiatAmountAny())


def test_encode_negative_uinteger_raises():
    with pytest.raises(ValueError, match="non-negative"):
        rule_cell_to_bytes("RuleUIntegerGreater", rc.UIntegerGreaterValue(value=-1))
    with pytest.raises(ValueError, match="non-negative"):
        rule_cell_to_bytes("RuleUIntegerGreater", rc.UIntegerGreaterIsEqual(value=-1))


# --- RawCell fallback branches ---


def test_unknown_column_type_decodes_to_raw_cell():
    cell = rule_cell_from_bytes("RuleUnknownColumn", b"\x08\x01")
    assert isinstance(cell, rc.RawCell)
    assert cell.column_type == "RuleUnknownColumn"
    assert rule_cell_to_bytes("RuleUnknownColumn", cell) == b"\x08\x01"


def test_unknown_column_empty_bytes_decodes_to_raw_cell():
    cell = rule_cell_from_bytes("RuleUnknownColumn", b"")
    assert isinstance(cell, rc.RawCell)
    assert rule_cell_to_bytes("RuleUnknownColumn", cell) == b""


def test_lossless_guard_extra_subfield_decodes_to_raw_cell():
    # A RuleFiatAmount wrapper (type Any) carrying an unknown protobuf field does
    # not re-encode byte-identically through the typed layer -> RawCell.
    data = b"\xc0\x0c\x2a"  # field 200, varint 42
    cell = rule_cell_from_bytes("RuleFiatAmount", data)
    assert isinstance(cell, rc.RawCell)
    assert rule_cell_to_bytes("RuleFiatAmount", cell) == data


def test_payload_on_payload_free_type_decodes_to_raw_cell():
    # A payload-less cell type (Any) that unexpectedly carries a payload must be
    # preserved verbatim rather than silently dropped.
    data = pb.RuleFiatAmount(type=pb.RuleFiatAmount.RuleFiatAmountAny, payload=b"x").SerializeToString()
    cell = rule_cell_from_bytes("RuleFiatAmount", data)
    assert isinstance(cell, rc.RawCell)
    assert rule_cell_to_bytes("RuleFiatAmount", cell) == data


def test_unknown_cell_enum_decodes_to_raw_cell_all_families():
    families = {
        "RuleFiatAmount": pb.RuleFiatAmount,
        "RuleSource": pb.RuleSource,
        "RuleDestination": pb.RuleDestination,
        "RuleWhitelistedContract": pb.RuleWhitelistedContract,
        "RuleStringEqual": pb.RuleStringEqual,
        "RuleBytesEqual": pb.RuleBytesEqual,
        "RuleStringArrayEqual": pb.RuleStringArrayEqual,
        "RuleIntegerGreater": pb.RuleIntegerGreater,
        "RuleUIntegerGreater": pb.RuleUIntegerGreater,
    }
    for col, msg_type in families.items():
        data = msg_type(type=902).SerializeToString()  # enum value newer than this SDK
        cell = rule_cell_from_bytes(col, data)
        assert isinstance(cell, rc.RawCell), col
        assert rule_cell_to_bytes(col, cell) == data, col


# --- Magnitude helpers (direct) ---


def test_magnitude_encoding():
    assert _magnitude(0) == b"\x00"
    assert _magnitude(128) == b"\x80"        # top bit set, no sign byte
    assert _magnitude(256) == b"\x01\x00"
    assert _magnitude(-50) == b"\x32"        # magnitude of a negative


def test_from_magnitude_decoding():
    assert _from_magnitude(b"") == 0
    assert _from_magnitude(b"\x00") == 0
    assert _from_magnitude(b"\x80") == 128
    assert _from_magnitude(b"\x01\x00") == 256


def test_uinteger_zero_magnitude_roundtrips():
    data = rule_cell_to_bytes("RuleUIntegerGreater", rc.UIntegerGreaterValue(value=0))
    assert rule_cell_from_bytes("RuleUIntegerGreater", data) == rc.UIntegerGreaterValue(value=0)


# --- Blockchain enum passthrough (vectors only use "ETH") ---


def test_blockchain_numeric_passthrough():
    assert _blockchain_to_int("ETH") == pb.Blockchain.Value("ETH")
    assert _blockchain_to_int("4242") == 4242          # value newer than this SDK
    assert _blockchain_name(4242) == "4242"
    assert _blockchain_name(pb.Blockchain.Value("ETH")) == "ETH"
    with pytest.raises(ValueError):
        _blockchain_to_int("NotAChain")                 # caller-authored typo


def test_contract_address_numeric_blockchain_roundtrips():
    cell = rc.DestinationContractAddress(address="0x1", blockchain="4242")
    data = rule_cell_to_bytes("RuleDestination", cell)
    assert rule_cell_from_bytes("RuleDestination", data) == cell


# --- cell_family (direct) ---


def test_cell_family():
    assert cell_family(rc.FiatAmountRange()) == "RuleFiatAmount"
    assert cell_family(rc.StringEqualValue()) == "RuleStringEqual"
    assert cell_family(rc.StringArrayEqualValue()) == "RuleStringArrayEqual"  # distinct from RuleStringEqual
    assert cell_family(rc.IntegerGreaterValue()) == "RuleIntegerGreater"
    assert cell_family(rc.UIntegerGreaterValue()) == "RuleUIntegerGreater"    # distinct from RuleIntegerGreater
    assert cell_family(rc.WhitelistedContractAny()) == "RuleWhitelistedContract"
    # RawCell reports its own column type.
    assert cell_family(rc.RawCell(column_type="RuleFiatAmount")) == "RuleFiatAmount"
