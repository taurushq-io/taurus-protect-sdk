"""Cross-SDK cell wire-format alignment.

Every entry in the shared golden-vector file (produced by the Go SDK) must
encode to exactly the recorded bytes and decode back to the same typed cell.
The same file is consumed by the Go/Java/TypeScript suites, pinning byte-for-byte
parity across all four SDKs.
"""

import base64
import json
from pathlib import Path

from taurus_protect.models import rule_cell as rc
from taurus_protect.mappers.rule_cell_codec import rule_cell_from_bytes, rule_cell_to_bytes

VECTORS_PATH = Path(__file__).resolve().parents[4] / "scripts" / "resources" / "governance-cell-vectors.json"

RAW_DESCRIPTION = "raw cell (unknown cell type preserved verbatim)"

# Hardcoded cell cases keyed by the shared vectors' `description`, mirroring the
# Go SDK's allCellCases table. The vectors file is the byte oracle.
CASES = {
    "fiat amount any": ("RuleFiatAmount", rc.FiatAmountAny()),
    "fiat amount is zero": ("RuleFiatAmount", rc.FiatAmountIsZero()),
    "fiat amount range": ("RuleFiatAmount", rc.FiatAmountRange(min_amount="1000", max_amount="50000")),
    "source any": ("RuleSource", rc.SourceAny()),
    "source internal wallet": ("RuleSource", rc.SourceInternalWallet(path="m/44'/60'/0'")),
    "source internal address": ("RuleSource", rc.SourceInternalAddress(address="0xabc", path="m/44'/60'/0'/0/0")),
    "source any exchange": ("RuleSource", rc.SourceAnyExchange()),
    "source exchange": ("RuleSource", rc.SourceExchange(label="kraken-main")),
    "source external address": ("RuleSource", rc.SourceExternalAddress(address="0xdef", memo="memo-1")),
    "destination any": ("RuleDestination", rc.DestinationAny()),
    "destination internal wallet": ("RuleDestination", rc.DestinationInternalWallet(path="m/44'/60'/1'")),
    "destination internal address": ("RuleDestination", rc.DestinationInternalAddress(address="0x111", path="m/44'/60'/1'/0/0")),
    "destination external address": ("RuleDestination", rc.DestinationExternalAddress(address="0x222", memo="dest-memo")),
    "destination any exchange": ("RuleDestination", rc.DestinationAnyExchange()),
    "destination exchange": ("RuleDestination", rc.DestinationExchange(label="binance-desk", memo="x")),
    "destination contract address": ("RuleDestination", rc.DestinationContractAddress(address="0x333", name="USDC", symbol="USDC", blockchain="ETH")),
    "destination contract address (unknown blockchain)": ("RuleDestination", rc.DestinationContractAddress(address="0x1", name="", symbol="", blockchain="4242")),
    "destination any external address": ("RuleDestination", rc.DestinationAnyExternalAddress()),
    "destination any contract address": ("RuleDestination", rc.DestinationAnyContractAddress()),
    "string equal any": ("RuleStringEqual", rc.StringEqualAny()),
    "string equal empty": ("RuleStringEqual", rc.StringEqualEmpty()),
    "string equal value": ("RuleStringEqual", rc.StringEqualValue(value="contract-id-42")),
    "bytes equal any": ("RuleBytesEqual", rc.BytesEqualAny()),
    "bytes equal empty": ("RuleBytesEqual", rc.BytesEqualEmpty()),
    "bytes equal value": ("RuleBytesEqual", rc.BytesEqualValue(value=bytes([0xDE, 0xAD, 0xBE, 0xEF]))),
    "string array equal any": ("RuleStringArrayEqual", rc.StringArrayEqualAny()),
    "string array equal empty": ("RuleStringArrayEqual", rc.StringArrayEqualEmpty()),
    "string array equal value": ("RuleStringArrayEqual", rc.StringArrayEqualValue(values=["a", "b", "c"])),
    "integer greater any": ("RuleIntegerGreater", rc.IntegerGreaterAny()),
    "integer greater positive value": ("RuleIntegerGreater", rc.IntegerGreaterValue(value=50)),
    "integer greater negative value": ("RuleIntegerGreater", rc.IntegerGreaterValue(value=-50)),
    "integer greater zero": ("RuleIntegerGreater", rc.IntegerGreaterValue(value=0)),
    "uinteger greater any": ("RuleUIntegerGreater", rc.UIntegerGreaterAny()),
    "uinteger greater is zero": ("RuleUIntegerGreater", rc.UIntegerGreaterIsZero()),
    "uinteger greater value": ("RuleUIntegerGreater", rc.UIntegerGreaterValue(value=18446744073709551615)),
    "uinteger greater is equal": ("RuleUIntegerGreater", rc.UIntegerGreaterIsEqual(value=10)),
    "whitelisted contract any": ("RuleWhitelistedContract", rc.WhitelistedContractAny()),
    "whitelisted contract address": ("RuleWhitelistedContract", rc.WhitelistedContractAddress(address="0x444", name="DAI", symbol="DAI", blockchain="ETH")),
}


def _load_vectors():
    return json.loads(VECTORS_PATH.read_text())


def test_vectors_cover_all_cases():
    vectors = _load_vectors()
    described = {v["description"] for v in vectors}
    for desc in CASES:
        assert desc in described, f"vector missing for {desc!r}"
    assert len(vectors) == len(CASES) + 1  # typed cases + one raw case


def test_cell_vectors_encode_and_decode():
    for v in _load_vectors():
        desc = v["description"]
        if desc == RAW_DESCRIPTION:
            wire = base64.b64decode(v["wire_base64"])
            decoded = rule_cell_from_bytes(v["column_type"], wire)
            assert isinstance(decoded, rc.RawCell), f"{desc}: expected RawCell, got {decoded!r}"
            assert rule_cell_to_bytes(v["column_type"], decoded) == wire
            continue
        col_type, cell = CASES[desc]
        assert col_type == v["column_type"], desc
        # encode(typed) == recorded wire bytes (cross-SDK byte parity)
        assert base64.b64encode(rule_cell_to_bytes(col_type, cell)).decode() == v["wire_base64"], desc
        # decode(wire) == typed
        assert rule_cell_from_bytes(col_type, base64.b64decode(v["wire_base64"])) == cell, desc
