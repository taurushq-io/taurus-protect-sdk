package com.taurushq.sdk.protect.client.mapper;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.Map;

import com.google.protobuf.ByteString;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingLine;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.ContractAddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.CosmosDetails;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.EvmCallContract;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleCell;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleColumn;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleLine;
import com.taurushq.sdk.protect.client.model.rulescontainer.CashSettlement;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSource;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceExchange;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceExternalAddress;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceInternalAddress;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceInternalWallet;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceType;
import com.taurushq.sdk.protect.client.model.rulescontainer.XtzCallContract;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.client.model.rulescontainer.SequentialThresholds;
import com.taurushq.sdk.protect.client.model.rulescontainer.TransactionRuleDetails;
import com.taurushq.sdk.protect.client.model.rulescontainer.TransactionRules;
import com.taurushq.sdk.protect.proto.v1.RequestReply;
import org.junit.jupiter.api.Test;

/**
 * Container-level round-trip and schema-evolution safety for the typed mapper.
 *
 * <p>Proves that a decoded container re-encodes byte-identically (lossless), that
 * server-controlled fields are stripped, and that protobuf fields / cell types unknown
 * to this SDK version are preserved rather than silently dropped.
 */
class RulesContainerRoundtripTest {

    private static final RulesContainerMapper M = RulesContainerMapper.INSTANCE;

    private static SequentialThresholds thresholds() {
        GroupThreshold gt = new GroupThreshold();
        gt.setGroupId("approvers");
        gt.setMinimumSignatures(2);
        SequentialThresholds st = new SequentialThresholds();
        st.setThresholds(Collections.singletonList(gt));
        return st;
    }

    private static ByteString cell(String columnType, RuleCell c) {
        return ByteString.copyFrom(RuleCellCodec.encode(columnType, c));
    }

    private static DecodedRulesContainer richContainer() {
        DecodedRulesContainer c = new DecodedRulesContainer();

        RuleUser u1 = new RuleUser();
        u1.setId("user-1");
        u1.setRoles(Collections.singletonList("SUPERADMIN"));
        Map<String, ByteString> uprops = new LinkedHashMap<>();
        uprops.put("team", ByteString.copyFromUtf8("ops"));
        u1.setProperties(uprops);
        RuleUser u2 = new RuleUser();
        u2.setId("user-2");
        c.setUsers(Arrays.asList(u1, u2));

        RuleGroup g = new RuleGroup();
        g.setId("approvers");
        g.setUserIds(Arrays.asList("user-1", "user-2"));
        c.setGroups(Collections.singletonList(g));

        c.setMinimumDistinctUserSignatures(2);
        c.setMinimumDistinctGroupSignatures(1);

        RuleColumn amount = new RuleColumn();
        amount.setType("RuleFiatAmount");
        amount.setName("amount");
        amount.setMetadataKey("amount");
        RuleColumn dest = new RuleColumn();
        dest.setType("RuleDestination");
        dest.setName("to");
        dest.setMetadataKey("destination");
        RuleColumn contract = new RuleColumn();
        contract.setType("RuleStringEqual");
        contract.setName("contract");
        contract.setMetadataKey("contract_id");

        RuleLine line1 = new RuleLine();
        line1.setCells(Arrays.asList(
                cell("RuleFiatAmount", new RuleCell.FiatAmountRange("1000", "50000")),
                cell("RuleDestination", new RuleCell.DestinationInternalWallet("m/44'/60'/1'")),
                cell("RuleStringEqual", new RuleCell.StringEqualValue("contract-42"))));
        line1.setParallelThresholds(Collections.singletonList(thresholds()));
        line1.setPriority(1);
        Map<String, ByteString> lprops = new LinkedHashMap<>();
        lprops.put("note", ByteString.copyFromUtf8("hv"));
        line1.setProperties(lprops);

        RuleLine line2 = new RuleLine();
        line2.setCells(Arrays.asList(
                cell("RuleFiatAmount", new RuleCell.FiatAmountAny()),
                cell("RuleDestination", new RuleCell.DestinationAny()),
                cell("RuleStringEqual", new RuleCell.StringEqualAny())));
        line2.setParallelThresholds(Collections.singletonList(thresholds()));

        EvmCallContract evm = new EvmCallContract();
        evm.setContractType("ERC20");
        evm.setMethodSignature("transfer(address,uint256)");
        TransactionRuleDetails details = new TransactionRuleDetails();
        details.setDomain("RuleDomainTransfer");
        details.setSubDomain("RuleSubDomainERC20");
        details.setBlockchain("ETH");
        details.setNetwork("mainnet");
        details.setEvmCallContract(evm);

        TransactionRules tr = new TransactionRules();
        tr.setKey("ETH/ERC20_transfer");
        tr.setColumns(Arrays.asList(amount, dest, contract));
        tr.setLines(Arrays.asList(line1, line2));
        tr.setDetails(details);
        c.setTransactionRules(Collections.singletonList(tr));

        RuleSourceInternalWallet w = new RuleSourceInternalWallet();
        w.setPath("m/44'/60'/0'");
        RuleSource src = new RuleSource();
        src.setType(RuleSourceType.RuleSourceInternalWallet);
        src.setInternalWallet(w);
        RuleSource anySrc = new RuleSource();
        anySrc.setType(RuleSourceType.RuleSourceAny);
        AddressWhitelistingLine awl = new AddressWhitelistingLine();
        awl.setCells(Arrays.asList(src, anySrc));
        awl.setParallelThresholds(Collections.singletonList(thresholds()));
        AddressWhitelistingRules awr = new AddressWhitelistingRules();
        awr.setCurrency("ETH");
        awr.setNetwork("mainnet");
        awr.setParallelThresholds(Collections.singletonList(thresholds()));
        awr.setLines(Collections.singletonList(awl));
        c.setAddressWhitelistingRules(Collections.singletonList(awr));

        ContractAddressWhitelistingRules cawr = new ContractAddressWhitelistingRules();
        cawr.setBlockchain("ETH");
        cawr.setNetwork("mainnet");
        cawr.setParallelThresholds(Collections.singletonList(thresholds()));
        c.setContractAddressWhitelistingRules(Collections.singletonList(cawr));

        c.setEnforcedRulesHash("server-hash");
        c.setTimestamp(1750000000L);
        c.setHsmSlotId(7);
        c.setMinimumCommitmentSignatures(1);
        c.setEngineIdentities(Arrays.asList("hsm-1", "hsm-2"));
        Map<String, ByteString> cprops = new LinkedHashMap<>();
        cprops.put("tenant", ByteString.copyFromUtf8("42"));
        c.setProperties(cprops);
        return c;
    }

    @Test
    void roundtripTypedCellsAndStrip() throws Exception {
        byte[] encoded = M.toBytes(richContainer());
        DecodedRulesContainer c = M.fromBytes(encoded);

        assertFalse(c.hasUnknownFields());
        // server-controlled fields stripped by the encoder
        assertEquals("", c.getEnforcedRulesHash());
        assertEquals(0L, c.getTimestamp());

        TransactionRules tr = c.getTransactionRules().get(0);
        assertEquals("ETH/ERC20_transfer", tr.getKey());
        assertEquals("RuleFiatAmount", tr.getColumns().get(0).getType());
        assertEquals("amount", tr.getColumns().get(0).getMetadataKey());

        RuleLine l0 = tr.getLines().get(0);
        assertEquals(new RuleCell.FiatAmountRange("1000", "50000"), RuleCellCodec.decode("RuleFiatAmount", l0.getCells().get(0).toByteArray()));
        assertEquals(new RuleCell.DestinationInternalWallet("m/44'/60'/1'"), RuleCellCodec.decode("RuleDestination", l0.getCells().get(1).toByteArray()));
        assertEquals(new RuleCell.StringEqualValue("contract-42"), RuleCellCodec.decode("RuleStringEqual", l0.getCells().get(2).toByteArray()));
        assertEquals(1, l0.getPriority());
        assertEquals(ByteString.copyFromUtf8("hv"), l0.getProperties().get("note"));

        // empty cells decode to the columns' typed *Any values
        RuleLine l1 = tr.getLines().get(1);
        assertEquals(new RuleCell.FiatAmountAny(), RuleCellCodec.decode("RuleFiatAmount", l1.getCells().get(0).toByteArray()));
        assertEquals(new RuleCell.DestinationAny(), RuleCellCodec.decode("RuleDestination", l1.getCells().get(1).toByteArray()));

        assertEquals("ETH", tr.getDetails().getBlockchain());
        assertEquals("transfer(address,uint256)", tr.getDetails().getEvmCallContract().getMethodSignature());

        RuleSource src = c.getAddressWhitelistingRules().get(0).getLines().get(0).getCells().get(0);
        assertEquals(RuleSourceType.RuleSourceInternalWallet, src.getType());
        assertEquals("m/44'/60'/0'", src.getInternalWallet().getPath());

        assertEquals(ByteString.copyFromUtf8("ops"), c.getUsers().get(0).getProperties().get("team"));
        assertEquals(7, c.getHsmSlotId());

        // re-encoding the decoded container is byte-stable
        assertArrayEquals(encoded, M.toBytes(c));
    }

    @Test
    void unknownFieldsSurviveRoundtrip() throws Exception {
        byte[] encoded = M.toBytes(richContainer());
        // append an unknown top-level field (field 500, varint 42) — valid proto concatenation
        byte[] withUnknown = Arrays.copyOf(encoded, encoded.length + 3);
        withUnknown[encoded.length] = (byte) 0xA0;
        withUnknown[encoded.length + 1] = (byte) 0x1F;
        withUnknown[encoded.length + 2] = (byte) 0x2A;

        DecodedRulesContainer c = M.fromBytes(withUnknown);
        assertTrue(c.hasUnknownFields());

        byte[] reencoded = M.toBytes(c);
        // the unknown field is preserved through decode -> model -> encode, and survives another cycle
        assertTrue(M.fromBytes(reencoded).hasUnknownFields());
        assertFalse(Arrays.equals(reencoded, M.toBytes(richContainer())));
    }

    @Test
    void unknownCellTypePreservedAsRaw() {
        // A cell wrapper with a cell-type enum value newer than this SDK -> RawCell.
        byte[] future = RequestReply.RuleFiatAmount.newBuilder()
                .setTypeValue(902).setPayload(ByteString.copyFromUtf8("future")).build().toByteArray();

        RuleCell cell = RuleCellCodec.decode("RuleFiatAmount", future);
        assertTrue(cell instanceof RuleCell.RawCell);
        assertArrayEquals(future, RuleCellCodec.encode("RuleFiatAmount", cell));
    }

    @Test
    void integerCellSignRoundtrips() {
        // negative and large unsigned values survive the magnitude+sign encoding
        assertEquals(new RuleCell.IntegerGreaterValue(BigInteger.valueOf(-50)),
                RuleCellCodec.decode("RuleIntegerGreater", RuleCellCodec.encode("RuleIntegerGreater", new RuleCell.IntegerGreaterValue(BigInteger.valueOf(-50)))));
        RuleCell big = new RuleCell.UIntegerGreaterValue(new BigInteger("18446744073709551615"));
        assertEquals(big, RuleCellCodec.decode("RuleUIntegerGreater", RuleCellCodec.encode("RuleUIntegerGreater", big)));
    }

    private static final ByteString UF = ByteString.copyFrom(new byte[] {(byte) 0xA0, 0x1F, 0x2A}); // field 500, varint 42

    @Test
    void unknownFieldsSurvivePerNode() throws Exception {
        GroupThreshold gt = new GroupThreshold();
        gt.setGroupId("g");
        gt.setMinimumSignatures(1);
        gt.setUnknownFields(UF);
        SequentialThresholds st = new SequentialThresholds();
        st.setThresholds(Collections.singletonList(gt));
        st.setUnknownFields(UF);

        RuleColumn col = new RuleColumn();
        col.setType("RuleFiatAmount");
        col.setUnknownFields(UF);
        RuleLine line = new RuleLine();
        line.setParallelThresholds(Collections.singletonList(st));
        line.setUnknownFields(UF);
        TransactionRuleDetails details = new TransactionRuleDetails();
        details.setDomain("RuleDomainTransfer");
        details.setUnknownFields(UF);
        TransactionRules tr = new TransactionRules();
        tr.setKey("k");
        tr.setColumns(Collections.singletonList(col));
        tr.setLines(Collections.singletonList(line));
        tr.setDetails(details);
        tr.setUnknownFields(UF);

        RuleUser u = new RuleUser();
        u.setId("u");
        u.setUnknownFields(UF);
        RuleGroup g = new RuleGroup();
        g.setId("g");
        g.setUnknownFields(UF);

        AddressWhitelistingLine awl = new AddressWhitelistingLine();
        awl.setUnknownFields(UF);
        AddressWhitelistingRules awr = new AddressWhitelistingRules();
        awr.setCurrency("ETH");
        awr.setLines(Collections.singletonList(awl));
        awr.setUnknownFields(UF);
        ContractAddressWhitelistingRules cawr = new ContractAddressWhitelistingRules();
        cawr.setBlockchain("ETH");
        cawr.setUnknownFields(UF);

        DecodedRulesContainer c = new DecodedRulesContainer();
        c.setUsers(Collections.singletonList(u));
        c.setGroups(Collections.singletonList(g));
        c.setTransactionRules(Collections.singletonList(tr));
        c.setAddressWhitelistingRules(Collections.singletonList(awr));
        c.setContractAddressWhitelistingRules(Collections.singletonList(cawr));
        c.setUnknownFields(UF);

        DecodedRulesContainer out = M.fromBytes(M.toBytes(c));
        assertTrue(out.hasUnknownFields());

        assertEquals(UF, out.getUnknownFields());
        assertEquals(UF, out.getUsers().get(0).getUnknownFields());
        assertEquals(UF, out.getGroups().get(0).getUnknownFields());
        TransactionRules otr = out.getTransactionRules().get(0);
        assertEquals(UF, otr.getUnknownFields());
        assertEquals(UF, otr.getColumns().get(0).getUnknownFields());
        assertEquals(UF, otr.getLines().get(0).getUnknownFields());
        assertEquals(UF, otr.getDetails().getUnknownFields());
        assertEquals(UF, otr.getLines().get(0).getParallelThresholds().get(0).getUnknownFields());
        assertEquals(UF, otr.getLines().get(0).getParallelThresholds().get(0).getThresholds().get(0).getUnknownFields());
        assertEquals(UF, out.getAddressWhitelistingRules().get(0).getUnknownFields());
        assertEquals(UF, out.getAddressWhitelistingRules().get(0).getLines().get(0).getUnknownFields());
        assertEquals(UF, out.getContractAddressWhitelistingRules().get(0).getUnknownFields());

        // byte-stable re-encode
        assertArrayEquals(M.toBytes(c), M.toBytes(out));
    }

    @Test
    void ruleSourceVariantsRoundtrip() {
        RuleSourceInternalAddress ia = new RuleSourceInternalAddress();
        ia.setAddress("0xabc");
        ia.setPath("m/44'/60'/0'/0/0");
        RuleSource internalAddress = new RuleSource();
        internalAddress.setType(RuleSourceType.RuleSourceInternalAddress);
        internalAddress.setInternalAddress(ia);

        RuleSourceExchange ex = new RuleSourceExchange();
        ex.setLabel("kraken");
        RuleSource exchange = new RuleSource();
        exchange.setType(RuleSourceType.RuleSourceExchange);
        exchange.setExchange(ex);

        RuleSourceExternalAddress ea = new RuleSourceExternalAddress();
        ea.setAddress("0xdef");
        ea.setMemo("m");
        RuleSource external = new RuleSource();
        external.setType(RuleSourceType.RuleSourceExternalAddress);
        external.setExternalAddress(ea);

        RuleSource anyExchange = new RuleSource();
        anyExchange.setType(RuleSourceType.RuleSourceAnyExchange);

        for (RuleSource s : Arrays.asList(internalAddress, exchange, external, anyExchange)) {
            RuleSource decoded = M.ruleSourceFromBytes(M.ruleSourceToBytes(s));
            assertEquals(s.getType(), decoded.getType());
        }
        RuleSource decodedIa = M.ruleSourceFromBytes(M.ruleSourceToBytes(internalAddress));
        assertEquals("0xabc", decodedIa.getInternalAddress().getAddress());
        assertEquals("kraken", M.ruleSourceFromBytes(M.ruleSourceToBytes(exchange)).getExchange().getLabel());
        assertEquals("m", M.ruleSourceFromBytes(M.ruleSourceToBytes(external)).getExternalAddress().getMemo());
    }

    @Test
    void ruleSourceMalformedPayloadSurvivesAsRaw() throws Exception {
        // A source claiming InternalWallet but carrying an undecodable payload is
        // preserved verbatim rather than dropped.
        byte[] bad = RequestReply.RuleSource.newBuilder()
                .setType(RequestReply.RuleSource.RuleSourceType.RuleSourceInternalWallet)
                .setPayload(ByteString.copyFrom(new byte[] {(byte) 0xff, (byte) 0xff, (byte) 0xff, (byte) 0xff}))
                .build().toByteArray();
        RuleSource s = M.ruleSourceFromBytes(ByteString.copyFrom(bad));
        assertNotNull(s.getRaw());
        assertEquals(ByteString.copyFrom(bad), M.ruleSourceToBytes(s));
    }

    @Test
    void detailsXtzCashCosmosRoundtrip() throws Exception {
        XtzCallContract xtz = new XtzCallContract();
        xtz.setContractType("FA2");
        xtz.setMethodSignature("transfer");
        CashSettlement cash = new CashSettlement();
        cash.setProvider("prov");
        cash.setRequestType("settle");
        TransactionRuleDetails d = new TransactionRuleDetails();
        d.setDomain("RuleDomainCallContract");
        d.setXtzCallContract(xtz);
        d.setCashSettlement(cash);
        CosmosDetails cosmos = new CosmosDetails();
        cosmos.setMethodSignatures(Collections.singletonList("/cosmos.bank.v1beta1.MsgSend"));
        d.setCosmosDetails(cosmos);
        TransactionRules tr = new TransactionRules();
        tr.setKey("k");
        tr.setDetails(d);
        DecodedRulesContainer c = new DecodedRulesContainer();
        c.setTransactionRules(Collections.singletonList(tr));

        TransactionRuleDetails od = M.fromBytes(M.toBytes(c)).getTransactionRules().get(0).getDetails();
        assertEquals("FA2", od.getXtzCallContract().getContractType());
        assertEquals("settle", od.getCashSettlement().getRequestType());
        assertEquals(Collections.singletonList("/cosmos.bank.v1beta1.MsgSend"),
                od.getCosmosDetails().getMethodSignatures());
    }

    @Test
    void reattachMalformedUnknownFieldsIsIgnored() {
        // Malformed preserved bytes must not crash the encoder (they are logged and skipped).
        DecodedRulesContainer c = new DecodedRulesContainer();
        c.setUnknownFields(ByteString.copyFrom(new byte[] {(byte) 0xff, (byte) 0xff}));
        assertDoesNotThrow(() -> M.toBytes(c));
    }

    @Test
    void toBase64StringMatchesEncodedBytes() {
        DecodedRulesContainer c = richContainer();
        assertEquals(java.util.Base64.getEncoder().encodeToString(M.toBytes(c)), M.toBase64String(c));
    }

    // A whitelisting RuleSource this SDK cannot fully represent keeps its exact wire
    // bytes. The same three base64 vectors are asserted in all four SDKs.
    @Test
    void ruleSourcePreservesUnrepresentableVerbatim() {
        Map<String, String> cases = new LinkedHashMap<String, String>();
        java.util.List<String> descriptions =
                LosslessVectors.descriptionsForScenario("rule_source_lossless");
        java.util.List<String> wire = LosslessVectors.forScenario("rule_source_lossless");
        for (int i = 0; i < wire.size(); i++) {
            cases.put(descriptions.get(i), wire.get(i));
        }

        for (Map.Entry<String, String> e : cases.entrySet()) {
            byte[] data = java.util.Base64.getDecoder().decode(e.getValue());
            com.taurushq.sdk.protect.client.model.rulescontainer.RuleSource src =
                    M.ruleSourceFromBytes(ByteString.copyFrom(data));
            assertNotNull(src.getRaw(), e.getKey() + ": expected raw preservation");
            assertFalse(src.getRaw().isEmpty(), e.getKey() + ": expected raw preservation");
            assertArrayEquals(data, M.ruleSourceToBytes(src).toByteArray(),
                    e.getKey() + ": not re-emitted verbatim");
        }
    }

    // Enum values newer than this SDK keep their numbers across a decode/encode round
    // trip: collapsing them to the zero value would rewrite a column's family or widen a
    // rule's sub-domain. The same base64 vector is asserted in all four SDKs.
    @Test
    void unknownEnumsPassThroughNumerically() throws Exception {
        final String vector = LosslessVectors.one("unknown_enum_passthrough");

        DecodedRulesContainer c = M.fromProto(RequestReply.RulesContainer.parseFrom(
                java.util.Base64.getDecoder().decode(vector)));

        assertEquals("201", c.getUsers().get(0).getRoles().get(0));
        assertEquals("77", c.getTransactionRules().get(0).getColumns().get(0).getType());
        assertEquals("202", c.getTransactionRules().get(0).getDetails().getSubDomain());
        assertEquals("203", c.getContractAddressWhitelistingRules().get(0).getBlockchain());
        assertEquals(vector, M.toBase64String(c));
    }

    // The container carries map<string, bytes> properties at five levels. Map iteration
    // order is unspecified, so without deterministic serialization the same reviewed
    // container encodes to different bytes across runs and across SDKs. The expected
    // value is asserted byte-for-byte in all four SDKs.
    @Test
    void encodingIsDeterministicAndCrossSdkStable() {
        final String expected = LosslessVectors.one("deterministic_encoding");

        // Keys are inserted in reverse-sorted order on purpose.
        Map<String, ByteString> props = new LinkedHashMap<String, ByteString>();
        for (String k : new String[] {"kEcho", "kDelta", "kCharlie", "kBravo", "kAlpha"}) {
            props.put(k, ByteString.copyFromUtf8(k));
        }

        for (int i = 0; i < 20; i++) {
            DecodedRulesContainer c = new DecodedRulesContainer();
            RuleUser u = new RuleUser();
            u.setId("u1");
            u.setPublicKeyPem("PEM");
            u.setRoles(Arrays.asList("SUPERADMIN"));
            u.setProperties(new LinkedHashMap<String, ByteString>(props));
            c.setUsers(Arrays.asList(u));
            c.setProperties(new LinkedHashMap<String, ByteString>(props));
            assertEquals(expected, M.toBase64String(c), "run " + i + " is not stable");
        }
    }

    // A malformed cell must degrade on its own and never abort the container: every rule
    // for that tenant would go down with it, taking whitelisted-address verification with
    // them. A cell payload is protobuf bytes, so a non-UTF-8 string cell and a truncated
    // wrapper are both legal on the wire. Same vector is asserted in all four SDKs.
    @Test
    void malformedCellDoesNotAbortDecode() throws Exception {
        final String vector = LosslessVectors.one("malformed_cell_degrades_alone");

        DecodedRulesContainer c = M.fromProto(RequestReply.RulesContainer.parseFrom(
                java.util.Base64.getDecoder().decode(vector)));

        assertEquals(1, c.getTransactionRules().size());
        assertEquals(2, c.getTransactionRules().get(0).getLines().get(0).getCells().size());
        assertEquals(vector, M.toBase64String(c));
    }

    // Unknown protobuf fields inside the nested contract-call scoping sub-messages must
    // survive a round trip: dropping them silently narrows which contract calls a rule
    // covers. Java captured them at decode and threw them away at encode. The same
    // base64 vector is asserted in all four SDKs.
    @Test
    void nestedRuleDetailUnknownFieldsSurviveRoundtrip() throws Exception {
        final String vector = LosslessVectors.one("nested_detail_unknown_fields");

        DecodedRulesContainer c = M.fromProto(RequestReply.RulesContainer.parseFrom(
                java.util.Base64.getDecoder().decode(vector)));
        TransactionRuleDetails d = c.getTransactionRules().get(0).getDetails();

        assertTrue(d.getEvmCallContract().hasUnknownFields(), "evmCallContract");
        assertTrue(d.getXtzCallContract().hasUnknownFields(), "xtzCallContract");
        assertTrue(d.getCashSettlement().hasUnknownFields(), "cashSettlement");
        assertTrue(d.getCosmosDetails().hasUnknownFields(), "cosmosDetails");
        // The container-level report must see the nested nodes, not just its own bytes.
        assertTrue(c.hasUnknownFields(), "container must report nested unknown fields");

        assertEquals(vector, M.toBase64String(c));
    }

    // An empty payload on a payload-carrying arm leaves the typed variant unset: parsing
    // it would materialize a default sub-message, so a caller checking the variant for
    // null would see a wallet with an empty path instead of nothing. All four SDKs.
    @Test
    void emptyPayloadLeavesVariantUnset() {
        com.taurushq.sdk.protect.client.model.rulescontainer.RuleSource src =
                M.ruleSourceFromBytes(ByteString.copyFrom(java.util.Base64.getDecoder().decode(LosslessVectors.one("empty_payload_arm"))));
        assertEquals(com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceType
                .RuleSourceInternalWallet, src.getType());
        assertNull(src.getInternalWallet());
    }

    // `08 01 12 00` and `08 01` decode to the SAME typed value: an explicitly-present
    // zero-length payload is legal on the wire and leaves the variant unset just like an
    // absent one, so re-encoding the typed form emits `08 01` and drops two bytes from a
    // container the SuperAdmins signed. No unknown field is present, so only the
    // decode -> re-encode -> byte-compare guard catches it. Java has had that guard;
    // this pins it, and pins the shared vector all four SDKs now carry.
    @Test
    void nonCanonicalEmptyPayloadStaysRaw() {
        ByteString data = ByteString.copyFrom(java.util.Base64.getDecoder()
                .decode(LosslessVectors.one("explicit_empty_payload_noncanonical")));

        com.taurushq.sdk.protect.client.model.rulescontainer.RuleSource src =
                M.ruleSourceFromBytes(data);

        assertEquals(data, src.getRaw(), "non-canonical source must be preserved verbatim");
        assertEquals(data, M.ruleSourceToBytes(src), "re-encode must not rewrite signed bytes");
    }
}
