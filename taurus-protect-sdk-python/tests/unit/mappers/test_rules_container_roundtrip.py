"""Container-level round-trip + schema-evolution safety for the typed mapper."""

from tests.unit.mappers.lossless_vectors import vector_for, vectors_for

from taurus_protect._internal.proto import request_reply_pb2 as pb
from taurus_protect.mappers.governance_rules import rules_container_from_base64
from taurus_protect.mappers.rules_container_encode import rules_container_to_bytes
from taurus_protect.models import rule_cell as rc
from taurus_protect.models.governance_rules import (
    RULE_SOURCE_TYPE_INTERNAL_WALLET,
    AddressWhitelistingLine,
    AddressWhitelistingRules,
    ContractAddressWhitelistingRules,
    DecodedRulesContainer,
    EvmCallContract,
    GroupThreshold,
    RuleColumn,
    RuleGroup,
    RuleLine,
    RuleSource,
    RuleSourceInternalWallet,
    RuleUser,
    SequentialThresholds,
    TransactionRuleDetails,
    TransactionRules,
)


def _rich_container() -> DecodedRulesContainer:
    thresholds = [SequentialThresholds(thresholds=[GroupThreshold(group_id="approvers", minimum_signatures=2)])]
    return DecodedRulesContainer(
        users=[
            RuleUser(id="user-1", roles=["SUPERADMIN"], properties={"team": b"ops"}),
            RuleUser(id="user-2"),
        ],
        groups=[RuleGroup(id="approvers", user_ids=["user-1", "user-2"])],
        minimum_distinct_user_signatures=2,
        minimum_distinct_group_signatures=1,
        transaction_rules=[
            TransactionRules(
                key="ETH/ERC20_transfer",
                columns=[
                    RuleColumn(type="RuleFiatAmount", name="amount", metadata_key="amount"),
                    RuleColumn(type="RuleDestination", name="to", metadata_key="destination"),
                    RuleColumn(type="RuleStringEqual", name="contract", metadata_key="contract_id"),
                ],
                lines=[
                    RuleLine(
                        cells=[
                            rc.FiatAmountRange(min_amount="1000", max_amount="50000"),
                            rc.DestinationInternalWallet(path="m/44'/60'/1'"),
                            rc.StringEqualValue(value="contract-42"),
                        ],
                        parallel_thresholds=thresholds,
                        priority=1,
                        properties={"note": b"hv"},
                    ),
                    RuleLine(
                        cells=[rc.FiatAmountAny(), rc.DestinationAny(), rc.StringEqualAny()],
                        parallel_thresholds=thresholds,
                    ),
                ],
                details=TransactionRuleDetails(
                    domain="RuleDomainTransfer",
                    sub_domain="RuleSubDomainERC20",
                    blockchain="ETH",
                    network="mainnet",
                    evm_call_contract=EvmCallContract(contract_type="ERC20", method_signature="transfer(address,uint256)"),
                ),
            )
        ],
        address_whitelisting_rules=[
            AddressWhitelistingRules(
                currency="ETH",
                network="mainnet",
                parallel_thresholds=thresholds,
                lines=[
                    AddressWhitelistingLine(
                        cells=[
                            RuleSource(type=RULE_SOURCE_TYPE_INTERNAL_WALLET, internal_wallet=RuleSourceInternalWallet(path="m/44'/60'/0'")),
                            RuleSource(type=0),
                        ],
                        parallel_thresholds=thresholds,
                    )
                ],
            )
        ],
        contract_address_whitelisting_rules=[
            ContractAddressWhitelistingRules(blockchain="ETH", network="mainnet", parallel_thresholds=thresholds)
        ],
        enforced_rules_hash="server-hash",
        timestamp=1750000000,
        hsm_slot_id=7,
        minimum_commitment_signatures=1,
        engine_identities=["hsm-1", "hsm-2"],
        properties={"tenant": b"42"},
    )


def test_roundtrip_typed_cells_and_strip():
    encoded = rules_container_to_bytes(_rich_container())
    c = rules_container_from_base64(_b64(encoded))
    assert not c.has_unknown_fields()

    # server-controlled fields stripped by the encoder
    assert c.enforced_rules_hash == ""
    assert c.timestamp == 0

    tr = c.transaction_rules[0]
    assert tr.key == "ETH/ERC20_transfer"
    assert tr.columns[0].type == "RuleFiatAmount" and tr.columns[0].metadata_key == "amount"
    assert tr.lines[0].cells[0] == rc.FiatAmountRange(min_amount="1000", max_amount="50000")
    assert tr.lines[0].cells[1] == rc.DestinationInternalWallet(path="m/44'/60'/1'")
    assert tr.lines[0].cells[2] == rc.StringEqualValue(value="contract-42")
    assert tr.lines[0].priority == 1
    # empty cells decode to the columns' typed *Any values
    assert tr.lines[1].cells == [rc.FiatAmountAny(), rc.DestinationAny(), rc.StringEqualAny()]
    assert tr.details.blockchain == "ETH"
    assert tr.details.evm_call_contract.method_signature == "transfer(address,uint256)"

    src = c.address_whitelisting_rules[0].lines[0].cells[0]
    assert src.type == RULE_SOURCE_TYPE_INTERNAL_WALLET
    assert src.internal_wallet.path == "m/44'/60'/0'"

    assert dict(c.users[0].properties) == {"team": b"ops"}
    assert c.hsm_slot_id == 7

    # re-encoding the decoded container is byte-stable
    assert rules_container_to_bytes(c) == encoded


def test_unknown_fields_survive_roundtrip():
    encoded = rules_container_to_bytes(_rich_container())
    # append an unknown top-level field (field 500, varint 42) — valid proto concat
    with_unknown = encoded + bytes([0xA0, 0x1F, 0x2A])

    c = rules_container_from_base64(_b64(with_unknown))
    assert c.has_unknown_fields()

    reencoded = rules_container_to_bytes(c)
    # the unknown field is preserved through decode -> model -> encode
    raw = pb.RulesContainer.FromString(reencoded)
    assert raw.SerializeToString() != rules_container_to_bytes(_rich_container())  # differs by the unknown field
    # and it survives another full cycle
    assert rules_container_from_base64(_b64(reencoded)).has_unknown_fields()


def test_unknown_cell_type_preserved_as_raw():
    # A cell wrapper with a cell-type enum value newer than this SDK -> RawCell.
    future = pb.RuleFiatAmount(type=902, payload=b"future").SerializeToString()
    from taurus_protect.mappers.rule_cell_codec import rule_cell_from_bytes, rule_cell_to_bytes

    cell = rule_cell_from_bytes("RuleFiatAmount", future)
    assert isinstance(cell, rc.RawCell)
    assert rule_cell_to_bytes("RuleFiatAmount", cell) == future


def _b64(data: bytes) -> str:
    import base64

    return base64.b64encode(data).decode()


class TestRuleSourceLosslessGuards:
    """A whitelisting RuleSource this SDK cannot fully represent keeps its exact bytes.

    The same three base64 vectors are asserted in all four SDKs.
    """

    CASES = {
        v["description"]: v["wire_base64"]
        for v in vectors_for("rule_source_lossless")
    }

    def test_preserves_unrepresentable_sources(self) -> None:
        import base64 as _b64

        from taurus_protect.mappers.governance_rules import _rule_source_from_bytes
        from taurus_protect.mappers.rules_container_encode import _rule_source_to_bytes

        for name, b64 in self.CASES.items():
            data = _b64.b64decode(b64)
            src = _rule_source_from_bytes(data)
            assert src.raw, f"{name}: expected raw preservation, got {src!r}"
            assert _rule_source_to_bytes(src) == data, f"{name}: not re-emitted verbatim"


def test_unknown_enums_pass_through_numerically() -> None:
    """Enum values newer than this SDK keep their numbers across a round trip.

    Collapsing them to the zero value would rewrite a column's family or widen a
    rule's sub-domain. The same base64 vector is asserted in all four SDKs.
    """
    import base64 as _b64

    from taurus_protect.mappers.governance_rules import rules_container_from_base64
    from taurus_protect.mappers.rules_container_encode import rules_container_to_bytes

    vector = vector_for("unknown_enum_passthrough")
    c = rules_container_from_base64(vector)

    assert c.users[0].roles[0] == "201"
    assert c.transaction_rules[0].columns[0].type == "77"
    assert c.transaction_rules[0].details.sub_domain == "202"
    assert c.contract_address_whitelisting_rules[0].blockchain == "203"
    assert _b64.b64encode(rules_container_to_bytes(c)).decode() == vector


def test_encoding_is_deterministic_and_cross_sdk_stable() -> None:
    """Map-valued properties must emit in a stable, cross-SDK order.

    Map iteration order is unspecified, so without deterministic marshaling the same
    reviewed container encodes to different bytes across runs and across SDKs, making
    any client-side hash of the encoded container meaningless. The expected value is
    asserted byte-for-byte in all four SDKs.
    """
    import base64 as _b64

    from taurus_protect.models.governance_rules import DecodedRulesContainer, RuleUser
    from taurus_protect.mappers.rules_container_encode import rules_container_to_bytes

    expected = vector_for("deterministic_encoding")

    # Keys are inserted in reverse-sorted order on purpose.
    props = {}
    for k in ["kEcho", "kDelta", "kCharlie", "kBravo", "kAlpha"]:
        props[k] = k.encode()

    for _ in range(20):
        c = DecodedRulesContainer(
            users=[RuleUser(id="u1", public_key_pem="PEM", roles=["SUPERADMIN"],
                            properties=dict(props))],
            properties=dict(props),
        )
        assert _b64.b64encode(rules_container_to_bytes(c)).decode() == expected


def test_malformed_cell_does_not_abort_decode() -> None:
    """A malformed cell degrades on its own; it must not abort the container.

    Every rule for that tenant would go down with it, taking whitelisted-address
    verification with them. A cell payload is protobuf bytes, so a non-UTF-8
    string cell and a truncated wrapper are both legal on the wire. Same vector is
    asserted in all four SDKs.
    """
    import base64 as _b64

    from taurus_protect.mappers.governance_rules import rules_container_from_base64
    from taurus_protect.mappers.rules_container_encode import rules_container_to_bytes

    vector = vector_for("malformed_cell_degrades_alone")
    c = rules_container_from_base64(vector)

    assert len(c.transaction_rules) == 1
    cells = c.transaction_rules[0].lines[0].cells
    assert len(cells) == 2
    # Python str cannot hold invalid UTF-8, so both cells are preserved verbatim.
    assert all(type(cell).__name__ == "RawCell" for cell in cells)
    assert _b64.b64encode(rules_container_to_bytes(c)).decode() == vector


def test_nested_rule_detail_unknown_fields_survive() -> None:
    """Unknown fields inside the nested contract-call scoping sub-messages survive.

    Dropping them silently narrows which contract calls a rule covers. The same
    base64 vector is asserted in all four SDKs.
    """
    import base64 as _b64

    from taurus_protect.mappers.governance_rules import rules_container_from_base64
    from taurus_protect.mappers.rules_container_encode import rules_container_to_bytes

    vector = vector_for("nested_detail_unknown_fields")
    c = rules_container_from_base64(vector)
    d = c.transaction_rules[0].details

    assert d.evm_call_contract.unknown_fields, "evm_call_contract"
    assert d.xtz_call_contract.unknown_fields, "xtz_call_contract"
    assert d.cash_settlement.unknown_fields, "cash_settlement"
    assert d.cosmos_details.unknown_fields, "cosmos_details"
    # The container-level report must see the nested nodes, not just the details node.
    assert c.has_unknown_fields()

    assert _b64.b64encode(rules_container_to_bytes(c)).decode() == vector


def test_empty_payload_leaves_variant_unset() -> None:
    """An empty payload on a payload-carrying arm leaves the typed variant unset.

    Parsing it would materialize a default sub-message, so a caller checking the
    variant for None would see an empty wallet instead of nothing. All four SDKs.
    """
    import base64 as _b64

    from taurus_protect.mappers.governance_rules import _rule_source_from_bytes
    from taurus_protect.models.governance_rules import RULE_SOURCE_TYPE_INTERNAL_WALLET

    src = _rule_source_from_bytes(_b64.b64decode(vector_for("empty_payload_arm")))
    assert src.type == RULE_SOURCE_TYPE_INTERNAL_WALLET
    assert src.internal_wallet is None


def test_non_canonical_empty_payload_stays_raw() -> None:
    """``08 01 12 00`` and ``08 01`` decode to the same typed value.

    An explicitly-present zero-length payload is legal on the wire and leaves the
    variant unset just like an absent one, so re-encoding the typed form emits
    ``08 01`` and drops two bytes from a signed container. There is no unknown field
    to catch it — only the decode -> re-encode -> byte-compare guard sees it.
    """
    import base64 as _b64

    from taurus_protect.mappers.governance_rules import _rule_source_from_bytes
    from taurus_protect.mappers.rules_container_encode import _rule_source_to_bytes

    data = _b64.b64decode(vector_for("explicit_empty_payload_noncanonical"))
    src = _rule_source_from_bytes(data)
    assert src.raw == data, "non-canonical source must be preserved verbatim"
    assert _rule_source_to_bytes(src) == data, "re-encode must not rewrite signed bytes"
