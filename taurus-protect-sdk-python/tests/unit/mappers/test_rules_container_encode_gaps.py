"""Per-node unknown-field preservation, enum passthrough, the RuleSource variant
codec, and the Xtz details arm — the encode/decode branches the baseline
round-trip test leaves at the container top level only."""

import base64

import pytest

from taurus_protect._internal.proto import request_reply_pb2 as pb
from taurus_protect.mappers.governance_rules import rules_container_from_base64
from taurus_protect.mappers.rules_container_encode import (
    _enum_number,
    rules_container_to_base64,
    rules_container_to_bytes,
)
from taurus_protect.models.governance_rules import (
    RULE_SOURCE_TYPE_ANY_EXCHANGE,
    RULE_SOURCE_TYPE_EXCHANGE,
    RULE_SOURCE_TYPE_EXTERNAL_ADDRESS,
    RULE_SOURCE_TYPE_INTERNAL_ADDRESS,
    AddressWhitelistingLine,
    AddressWhitelistingRules,
    ContractAddressWhitelistingRules,
    DecodedRulesContainer,
    GroupThreshold,
    RuleColumn,
    RuleGroup,
    RuleLine,
    RuleSource,
    RuleSourceExchange,
    RuleSourceExternalAddress,
    RuleSourceInternalAddress,
    CashSettlement,
    CosmosDetails,
    SequentialThresholds,
    TransactionRuleDetails,
    TransactionRules,
    RuleUser,
    XtzCallContract,
)

UF = b"\xc0\x0c\x2a"  # unknown protobuf field 200, varint 42


def _b64(data: bytes) -> str:
    return base64.b64encode(data).decode()


def test_unknown_fields_survive_per_node():
    threshold = SequentialThresholds(
        thresholds=[GroupThreshold(group_id="g", minimum_signatures=1, unknown_fields=UF)],
        unknown_fields=UF,
    )
    container = DecodedRulesContainer(
        users=[RuleUser(id="u", unknown_fields=UF)],
        groups=[RuleGroup(id="g", unknown_fields=UF)],
        transaction_rules=[
            TransactionRules(
                key="k",
                columns=[RuleColumn(type="RuleFiatAmount", unknown_fields=UF)],
                lines=[RuleLine(cells=[], parallel_thresholds=[threshold], unknown_fields=UF)],
                details=TransactionRuleDetails(domain="RuleDomainTransfer", unknown_fields=UF),
                unknown_fields=UF,
            )
        ],
        address_whitelisting_rules=[
            AddressWhitelistingRules(
                currency="ETH",
                lines=[AddressWhitelistingLine(cells=[], unknown_fields=UF)],
                unknown_fields=UF,
            )
        ],
        contract_address_whitelisting_rules=[
            ContractAddressWhitelistingRules(blockchain="ETH", unknown_fields=UF)
        ],
        unknown_fields=UF,
    )

    c = rules_container_from_base64(_b64(rules_container_to_bytes(container)))
    assert c.has_unknown_fields()

    assert c.unknown_fields == UF
    assert c.users[0].unknown_fields == UF
    assert c.groups[0].unknown_fields == UF
    tr = c.transaction_rules[0]
    assert tr.unknown_fields == UF
    assert tr.columns[0].unknown_fields == UF
    assert tr.lines[0].unknown_fields == UF
    assert tr.details.unknown_fields == UF
    assert tr.lines[0].parallel_thresholds[0].unknown_fields == UF
    assert tr.lines[0].parallel_thresholds[0].thresholds[0].unknown_fields == UF
    assert c.address_whitelisting_rules[0].unknown_fields == UF
    assert c.address_whitelisting_rules[0].lines[0].unknown_fields == UF
    assert c.contract_address_whitelisting_rules[0].unknown_fields == UF

    # Byte-stable re-encode.
    assert rules_container_to_bytes(c) == rules_container_to_bytes(container)


def test_all_enums_pass_through_numerically():
    container = DecodedRulesContainer(
        users=[RuleUser(id="u", roles=["991"])],
        transaction_rules=[
            TransactionRules(
                key="k",
                details=TransactionRuleDetails(domain="992", sub_domain="993"),
            )
        ],
        contract_address_whitelisting_rules=[ContractAddressWhitelistingRules(blockchain="994")],
    )
    c = rules_container_from_base64(_b64(rules_container_to_bytes(container)))
    assert c.users[0].roles == ["991"]
    assert c.transaction_rules[0].details.domain == "992"
    assert c.transaction_rules[0].details.sub_domain == "993"
    assert c.contract_address_whitelisting_rules[0].blockchain == "994"


def test_enum_number_helper():
    assert _enum_number(pb.Role, "", "role") == 0
    assert _enum_number(pb.Role, None, "role") == 0
    assert _enum_number(pb.Role, "991", "role") == 991  # decimal passthrough
    with pytest.raises(ValueError, match="unknown role"):
        _enum_number(pb.Role, "NotARole", "role")


def test_authored_unknown_enum_name_raises():
    container = DecodedRulesContainer(users=[RuleUser(id="u", roles=["NotARole"])])
    with pytest.raises(ValueError, match="unknown user role"):
        rules_container_to_bytes(container)


def test_rule_source_variants_roundtrip():
    sources = [
        RuleSource(
            type=RULE_SOURCE_TYPE_INTERNAL_ADDRESS,
            internal_address=RuleSourceInternalAddress(address="0xabc", path="m/44'/60'/0'/0/0"),
        ),
        RuleSource(type=RULE_SOURCE_TYPE_EXCHANGE, exchange=RuleSourceExchange(label="kraken")),
        RuleSource(
            type=RULE_SOURCE_TYPE_EXTERNAL_ADDRESS,
            external_address=RuleSourceExternalAddress(address="0xdef", memo="m"),
        ),
        RuleSource(type=RULE_SOURCE_TYPE_ANY_EXCHANGE),
    ]
    container = DecodedRulesContainer(
        address_whitelisting_rules=[
            AddressWhitelistingRules(currency="ETH", lines=[AddressWhitelistingLine(cells=sources)])
        ]
    )
    c = rules_container_from_base64(_b64(rules_container_to_bytes(container)))
    got = c.address_whitelisting_rules[0].lines[0].cells
    assert [s.type for s in got] == [
        RULE_SOURCE_TYPE_INTERNAL_ADDRESS,
        RULE_SOURCE_TYPE_EXCHANGE,
        RULE_SOURCE_TYPE_EXTERNAL_ADDRESS,
        RULE_SOURCE_TYPE_ANY_EXCHANGE,
    ]
    assert got[0].internal_address.address == "0xabc"
    assert got[1].exchange.label == "kraken"
    assert got[2].external_address.memo == "m"


def test_rule_source_malformed_payload_survives_as_raw():
    bad = pb.RuleSource(
        type=pb.RuleSource.RuleSourceInternalWallet, payload=b"\xff\xff\xff\xff"
    ).SerializeToString()
    pb_c = pb.RulesContainer()
    awr = pb_c.addressWhitelistingRules.add()
    awr.currency = "ETH"
    awr.lines.add().cells.append(bad)

    c = rules_container_from_base64(_b64(pb_c.SerializeToString()))
    src = c.address_whitelisting_rules[0].lines[0].cells[0]
    assert src.raw == bad  # preserved verbatim
    # re-encode keeps the opaque source intact
    reencoded = rules_container_from_base64(_b64(rules_container_to_bytes(c)))
    assert reencoded.address_whitelisting_rules[0].lines[0].cells[0].raw == bad


def test_rule_source_unknown_type_survives_as_raw():
    future = pb.RuleSource(type=99, payload=b"opaque").SerializeToString()
    pb_c = pb.RulesContainer()
    awr = pb_c.addressWhitelistingRules.add()
    awr.currency = "ETH"
    awr.lines.add().cells.append(future)

    c = rules_container_from_base64(_b64(pb_c.SerializeToString()))
    assert c.has_unknown_fields()
    assert c.address_whitelisting_rules[0].lines[0].cells[0].raw == future


def test_xtz_details_roundtrip():
    container = DecodedRulesContainer(
        transaction_rules=[
            TransactionRules(
                key="XTZ/call",
                details=TransactionRuleDetails(
                    domain="RuleDomainCallContract",
                    xtz_call_contract=XtzCallContract(contract_type="FA2", method_signature="transfer"),
                ),
            )
        ]
    )
    c = rules_container_from_base64(_b64(rules_container_to_bytes(container)))
    xtz = c.transaction_rules[0].details.xtz_call_contract
    assert xtz is not None
    assert xtz.contract_type == "FA2"
    assert xtz.method_signature == "transfer"


def test_cash_and_cosmos_details_roundtrip():
    container = DecodedRulesContainer(
        transaction_rules=[
            TransactionRules(
                key="k",
                details=TransactionRuleDetails(
                    domain="RuleDomainCashSettlement",
                    cash_settlement=CashSettlement(provider="prov", request_type="settle"),
                    cosmos_details=CosmosDetails(method_signatures=["/cosmos.bank.v1beta1.MsgSend"]),
                ),
            )
        ]
    )
    c = rules_container_from_base64(_b64(rules_container_to_bytes(container)))
    details = c.transaction_rules[0].details
    assert details.cash_settlement.provider == "prov"
    assert details.cash_settlement.request_type == "settle"
    assert details.cosmos_details.method_signatures == ["/cosmos.bank.v1beta1.MsgSend"]


def test_encode_none_container_raises():
    with pytest.raises(ValueError, match="cannot be None"):
        rules_container_to_bytes(None)
    with pytest.raises(ValueError, match="cannot be None"):
        rules_container_to_base64(None)
