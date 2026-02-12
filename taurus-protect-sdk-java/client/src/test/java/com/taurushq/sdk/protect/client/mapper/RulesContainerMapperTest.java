package com.taurushq.sdk.protect.client.mapper;

import com.google.protobuf.ByteString;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingLine;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.ContractAddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceType;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.client.model.rulescontainer.SequentialThresholds;
import com.taurushq.sdk.protect.proto.v1.RequestReply;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Base64;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class RulesContainerMapperTest {

    @Test
    void fromProto_mapsUsersCorrectly() {
        RequestReply.User protoUser = RequestReply.User.newBuilder()
                .setId("user-1")
                .setPublicKey("-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----")
                .addRoles(RequestReply.Role.SUPERADMIN)
                .addRoles(RequestReply.Role.HSMSLOT)
                .build();

        RequestReply.RulesContainer proto = RequestReply.RulesContainer.newBuilder()
                .addUsers(protoUser)
                .setMinimumDistinctUserSignatures(2)
                .setMinimumDistinctGroupSignatures(1)
                .setEnforcedRulesHash("hash123")
                .setTimestamp(1234567890L)
                .build();

        DecodedRulesContainer result = RulesContainerMapper.INSTANCE.fromProto(proto);

        assertNotNull(result);
        assertNotNull(result.getUsers());
        assertEquals(1, result.getUsers().size());

        RuleUser user = result.getUsers().get(0);
        assertEquals("user-1", user.getId());
        assertNotNull(user.getRoles());
        assertEquals(2, user.getRoles().size());
        assertTrue(user.getRoles().contains("SUPERADMIN"));
        assertTrue(user.getRoles().contains("HSMSLOT"));
        assertEquals("-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----",
                user.getPublicKeyPem());

        assertEquals(2, result.getMinimumDistinctUserSignatures());
        assertEquals(1, result.getMinimumDistinctGroupSignatures());
        assertEquals("hash123", result.getEnforcedRulesHash());
        assertEquals(1234567890L, result.getTimestamp());
    }

    @Test
    void fromProto_mapsGroupsCorrectly() {
        RequestReply.Group protoGroup = RequestReply.Group.newBuilder()
                .setId("group-1")
                .addUserIds("u1")
                .addUserIds("u2")
                .addUserIds("u3")
                .build();

        RequestReply.RulesContainer proto = RequestReply.RulesContainer.newBuilder()
                .addGroups(protoGroup)
                .build();

        DecodedRulesContainer result = RulesContainerMapper.INSTANCE.fromProto(proto);

        assertNotNull(result);
        assertNotNull(result.getGroups());
        assertEquals(1, result.getGroups().size());

        RuleGroup group = result.getGroups().get(0);
        assertEquals("group-1", group.getId());
        assertNotNull(group.getUserIds());
        assertEquals(3, group.getUserIds().size());
        assertEquals("u1", group.getUserIds().get(0));
        assertEquals("u2", group.getUserIds().get(1));
        assertEquals("u3", group.getUserIds().get(2));
    }

    @Test
    void fromProto_mapsAddressWhitelistingRules() {
        RequestReply.GroupThreshold protoThreshold = RequestReply.GroupThreshold.newBuilder()
                .setGroupId("g1")
                .setMinimumSignatures(2)
                .build();

        RequestReply.SequentialThresholds protoSeqThresholds =
                RequestReply.SequentialThresholds.newBuilder()
                        .addThresholds(protoThreshold)
                        .build();

        RequestReply.RulesContainer.AddressWhitelistingRules protoRule =
                RequestReply.RulesContainer.AddressWhitelistingRules.newBuilder()
                        .setCurrency("ETH")
                        .setNetwork("mainnet")
                        .addParallelThresholds(protoSeqThresholds)
                        .build();

        RequestReply.RulesContainer proto = RequestReply.RulesContainer.newBuilder()
                .addAddressWhitelistingRules(protoRule)
                .build();

        DecodedRulesContainer result = RulesContainerMapper.INSTANCE.fromProto(proto);

        assertNotNull(result);
        assertNotNull(result.getAddressWhitelistingRules());
        assertEquals(1, result.getAddressWhitelistingRules().size());

        AddressWhitelistingRules rule = result.getAddressWhitelistingRules().get(0);
        assertEquals("ETH", rule.getCurrency());
        assertEquals("mainnet", rule.getNetwork());
        assertNotNull(rule.getParallelThresholds());
        assertEquals(1, rule.getParallelThresholds().size());

        SequentialThresholds seqThresholds = rule.getParallelThresholds().get(0);
        assertNotNull(seqThresholds.getThresholds());
        assertEquals(1, seqThresholds.getThresholds().size());

        GroupThreshold threshold = seqThresholds.getThresholds().get(0);
        assertEquals("g1", threshold.getGroupId());
        assertEquals(2, threshold.getMinimumSignatures());
    }

    @Test
    void fromProto_mapsContractAddressWhitelistingRules() {
        RequestReply.GroupThreshold protoThreshold = RequestReply.GroupThreshold.newBuilder()
                .setGroupId("g2")
                .setMinimumSignatures(3)
                .build();

        RequestReply.SequentialThresholds protoSeqThresholds =
                RequestReply.SequentialThresholds.newBuilder()
                        .addThresholds(protoThreshold)
                        .build();

        RequestReply.RulesContainer.ContractAddressWhitelistingRules protoContractRule =
                RequestReply.RulesContainer.ContractAddressWhitelistingRules.newBuilder()
                        .setBlockchain(RequestReply.Blockchain.ETH)
                        .setNetwork("goerli")
                        .addParallelThresholds(protoSeqThresholds)
                        .build();

        RequestReply.RulesContainer proto = RequestReply.RulesContainer.newBuilder()
                .addContractAddressWhitelistingRules(protoContractRule)
                .build();

        DecodedRulesContainer result = RulesContainerMapper.INSTANCE.fromProto(proto);

        assertNotNull(result);
        assertNotNull(result.getContractAddressWhitelistingRules());
        assertEquals(1, result.getContractAddressWhitelistingRules().size());

        ContractAddressWhitelistingRules contractRule =
                result.getContractAddressWhitelistingRules().get(0);
        assertEquals("ETH", contractRule.getBlockchain());
        assertEquals("goerli", contractRule.getNetwork());
        assertNotNull(contractRule.getParallelThresholds());
        assertEquals(1, contractRule.getParallelThresholds().size());
    }

    @Test
    void fromProto_mapsEngineIdentities() {
        RequestReply.RulesContainer proto = RequestReply.RulesContainer.newBuilder()
                .addEngineIdentities("engine-serial-1")
                .addEngineIdentities("engine-serial-2")
                .setMinimumCommitmentSignatures(1)
                .build();

        DecodedRulesContainer result = RulesContainerMapper.INSTANCE.fromProto(proto);

        assertNotNull(result);
        assertNotNull(result.getEngineIdentities());
        assertEquals(2, result.getEngineIdentities().size());
        assertEquals("engine-serial-1", result.getEngineIdentities().get(0));
        assertEquals("engine-serial-2", result.getEngineIdentities().get(1));
        assertEquals(1, result.getMinimumCommitmentSignatures());
    }

    @Test
    void fromProto_mapsAddressWhitelistingLines() {
        RequestReply.RuleSourceInternalWallet protoWallet =
                RequestReply.RuleSourceInternalWallet.newBuilder()
                        .setPath("m/44'/60'/0'")
                        .build();

        RequestReply.RuleSource protoSource = RequestReply.RuleSource.newBuilder()
                .setType(RequestReply.RuleSource.RuleSourceType.RuleSourceInternalWallet)
                .setPayload(protoWallet.toByteString())
                .build();

        RequestReply.RulesContainer.AddressWhitelistingRules.Line protoLine =
                RequestReply.RulesContainer.AddressWhitelistingRules.Line.newBuilder()
                        .addCells(protoSource.toByteString())
                        .build();

        RequestReply.RulesContainer.AddressWhitelistingRules protoRule =
                RequestReply.RulesContainer.AddressWhitelistingRules.newBuilder()
                        .setCurrency("BTC")
                        .setNetwork("mainnet")
                        .addLines(protoLine)
                        .build();

        RequestReply.RulesContainer proto = RequestReply.RulesContainer.newBuilder()
                .addAddressWhitelistingRules(protoRule)
                .build();

        DecodedRulesContainer result = RulesContainerMapper.INSTANCE.fromProto(proto);

        assertNotNull(result);
        AddressWhitelistingRules rule = result.getAddressWhitelistingRules().get(0);
        assertNotNull(rule.getLines());
        assertEquals(1, rule.getLines().size());

        AddressWhitelistingLine line = rule.getLines().get(0);
        assertNotNull(line.getCells());
        assertEquals(1, line.getCells().size());
        assertEquals(RuleSourceType.RuleSourceInternalWallet, line.getCells().get(0).getType());
        assertNotNull(line.getCells().get(0).getInternalWallet());
        assertEquals("m/44'/60'/0'", line.getCells().get(0).getInternalWallet().getPath());
    }

    @Test
    void fromProto_handlesEmptyContainer() {
        RequestReply.RulesContainer proto = RequestReply.RulesContainer.getDefaultInstance();

        DecodedRulesContainer result = RulesContainerMapper.INSTANCE.fromProto(proto);

        assertNotNull(result);
        assertNotNull(result.getUsers());
        assertTrue(result.getUsers().isEmpty());
        assertNotNull(result.getGroups());
        assertTrue(result.getGroups().isEmpty());
        assertEquals(0, result.getMinimumDistinctUserSignatures());
        assertEquals(0, result.getMinimumDistinctGroupSignatures());
        assertEquals(0L, result.getTimestamp());
    }

    @Test
    void fromBytes_decodesProtobuf() throws InvalidProtocolBufferException {
        RequestReply.RulesContainer proto = RequestReply.RulesContainer.newBuilder()
                .setMinimumDistinctUserSignatures(3)
                .setEnforcedRulesHash("testhash")
                .build();

        byte[] bytes = proto.toByteArray();
        DecodedRulesContainer result = RulesContainerMapper.INSTANCE.fromBytes(bytes);

        assertNotNull(result);
        assertEquals(3, result.getMinimumDistinctUserSignatures());
        assertEquals("testhash", result.getEnforcedRulesHash());
    }

    @Test
    void fromBase64String_decodesBase64() throws InvalidProtocolBufferException {
        RequestReply.RulesContainer proto = RequestReply.RulesContainer.newBuilder()
                .setMinimumDistinctGroupSignatures(2)
                .setTimestamp(9876543210L)
                .build();

        String base64 = Base64.getEncoder().encodeToString(proto.toByteArray());
        DecodedRulesContainer result = RulesContainerMapper.INSTANCE.fromBase64String(base64);

        assertNotNull(result);
        assertEquals(2, result.getMinimumDistinctGroupSignatures());
        assertEquals(9876543210L, result.getTimestamp());
    }

    @Test
    void userSignaturesFromBytes_decodesSignatures() throws InvalidProtocolBufferException {
        byte[] sigBytes = "test-signature".getBytes();

        RequestReply.UserSignatures protoSigs = RequestReply.UserSignatures.newBuilder()
                .addSignatures(RequestReply.UserSignature.newBuilder()
                        .setUserId("user-1")
                        .setSignature(ByteString.copyFrom(sigBytes))
                        .build())
                .addSignatures(RequestReply.UserSignature.newBuilder()
                        .setUserId("user-2")
                        .setSignature(ByteString.copyFrom("sig-2".getBytes()))
                        .build())
                .build();

        List<RuleUserSignature> result = RulesContainerMapper.INSTANCE
                .userSignaturesFromBytes(protoSigs.toByteArray());

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("user-1", result.get(0).getUserId());
        assertNotNull(result.get(0).getSignature());
        assertEquals(Base64.getEncoder().encodeToString(sigBytes),
                result.get(0).getSignature());
        assertEquals("user-2", result.get(1).getUserId());
    }

    @Test
    void userSignaturesFromBase64String_decodesBase64Signatures()
            throws InvalidProtocolBufferException {
        RequestReply.UserSignatures protoSigs = RequestReply.UserSignatures.newBuilder()
                .addSignatures(RequestReply.UserSignature.newBuilder()
                        .setUserId("admin-1")
                        .setSignature(ByteString.copyFrom("my-sig".getBytes()))
                        .build())
                .build();

        String base64 = Base64.getEncoder().encodeToString(protoSigs.toByteArray());
        List<RuleUserSignature> result = RulesContainerMapper.INSTANCE
                .userSignaturesFromBase64String(base64);

        assertNotNull(result);
        assertEquals(1, result.size());
        assertEquals("admin-1", result.get(0).getUserId());
    }

    @Test
    void cellBytesToRuleSources_handlesNullInput() {
        List<com.taurushq.sdk.protect.client.model.rulescontainer.RuleSource> result =
                RulesContainerMapper.INSTANCE.cellBytesToRuleSources(null);

        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void rolesToStrings_handlesNullInput() {
        List<String> result = RulesContainerMapper.INSTANCE.rolesToStrings(null);

        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void rolesToStrings_convertsEnums() {
        List<String> result = RulesContainerMapper.INSTANCE.rolesToStrings(
                Arrays.asList(RequestReply.Role.SUPERADMIN, RequestReply.Role.HSMSLOT));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("SUPERADMIN", result.get(0));
        assertEquals("HSMSLOT", result.get(1));
    }

    @Test
    void blockchainToString_handlesNull() {
        assertNull(RulesContainerMapper.INSTANCE.blockchainToString(null));
    }

    @Test
    void blockchainToString_convertsEnum() {
        assertEquals("ETH", RulesContainerMapper.INSTANCE.blockchainToString(
                RequestReply.Blockchain.ETH));
        assertEquals("BTC", RulesContainerMapper.INSTANCE.blockchainToString(
                RequestReply.Blockchain.BTC));
    }

    @Test
    void columnTypeToString_handlesNull() {
        assertNull(RulesContainerMapper.INSTANCE.columnTypeToString(null));
    }

    @Test
    void columnTypeToString_convertsEnum() {
        String result = RulesContainerMapper.INSTANCE.columnTypeToString(
                RequestReply.RulesContainer.ColumnType.RuleSource);
        assertEquals("RuleSource", result);
    }
}
