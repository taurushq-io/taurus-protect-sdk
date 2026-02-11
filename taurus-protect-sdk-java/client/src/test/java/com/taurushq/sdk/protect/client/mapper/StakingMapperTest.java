package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ADAStakePoolInfo;
import com.taurushq.sdk.protect.client.model.ETHValidatorInfo;
import com.taurushq.sdk.protect.client.model.FTMValidatorInfo;
import com.taurushq.sdk.protect.client.model.ICPNeuronInfo;
import com.taurushq.sdk.protect.client.model.NEARValidatorInfo;
import com.taurushq.sdk.protect.client.model.SolanaStakeAccount;
import com.taurushq.sdk.protect.client.model.StakeAccount;
import com.taurushq.sdk.protect.client.model.StakeAccountResult;
import com.taurushq.sdk.protect.client.model.XTZStakingRewards;
import com.taurushq.sdk.protect.openapi.model.GetICPNeuronInfoReplyNeuronState;
import com.taurushq.sdk.protect.openapi.model.SolanaStakeAccountState;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordETHValidatorInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetADAStakePoolInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFTMValidatorInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetICPNeuronInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetNEARValidatorInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetStakeAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetXTZAddressStakingRewardsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSolanaStakeAccount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordStakeAccount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordStakeAccountType;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class StakingMapperTest {

    @Test
    void fromADAStakePoolInfoReply_mapsAllFields() {
        TgvalidatordGetADAStakePoolInfoReply reply = new TgvalidatordGetADAStakePoolInfoReply();
        reply.setPledge("1000000000");
        reply.setMargin(0.05f);
        reply.setFixedCost("340000000");
        reply.setUrl("https://pool.example.com");
        reply.setActiveStake("50000000000000");
        reply.setEpoch("450");

        ADAStakePoolInfo info = StakingMapper.INSTANCE.fromADAStakePoolInfoReply(reply);

        assertEquals("1000000000", info.getPledge());
        assertEquals(0.05f, info.getMargin());
        assertEquals("340000000", info.getFixedCost());
        assertEquals("https://pool.example.com", info.getUrl());
        assertEquals("50000000000000", info.getActiveStake());
        assertEquals("450", info.getEpoch());
    }

    @Test
    void fromETHValidatorInfo_mapsAllFields() {
        TgvalidatordETHValidatorInfo dto = new TgvalidatordETHValidatorInfo();
        dto.setId("12345");
        dto.setPubkey("0x1234567890abcdef");
        dto.setStatus("active_ongoing");
        dto.setBalance("32000000000");
        dto.setNetwork("mainnet");
        dto.setProvider("lido");
        dto.setAddressID("address-123");

        ETHValidatorInfo info = StakingMapper.INSTANCE.fromETHValidatorInfo(dto);

        assertEquals("12345", info.getId());
        assertEquals("0x1234567890abcdef", info.getPubkey());
        assertEquals("active_ongoing", info.getStatus());
        assertEquals("32000000000", info.getBalance());
        assertEquals("mainnet", info.getNetwork());
        assertEquals("lido", info.getProvider());
        assertEquals("address-123", info.getAddressId());
    }

    @Test
    void fromETHValidatorInfoList_mapsList() {
        TgvalidatordETHValidatorInfo dto1 = new TgvalidatordETHValidatorInfo();
        dto1.setId("1");
        dto1.setStatus("active");

        TgvalidatordETHValidatorInfo dto2 = new TgvalidatordETHValidatorInfo();
        dto2.setId("2");
        dto2.setStatus("pending");

        List<ETHValidatorInfo> infos = StakingMapper.INSTANCE.fromETHValidatorInfoList(
                Arrays.asList(dto1, dto2));

        assertEquals(2, infos.size());
        assertEquals("1", infos.get(0).getId());
        assertEquals("2", infos.get(1).getId());
    }

    @Test
    void fromFTMValidatorInfoReply_mapsAllFields() {
        TgvalidatordGetFTMValidatorInfoReply reply = new TgvalidatordGetFTMValidatorInfoReply();
        reply.setValidatorID("1");
        reply.setAddress("0xabc123");
        reply.setIsActive(true);
        reply.setTotalStake("1000000000000000000000");
        reply.setSelfStake("500000000000000000000");
        reply.setCreatedAtDateUnix("1609459200");

        FTMValidatorInfo info = StakingMapper.INSTANCE.fromFTMValidatorInfoReply(reply);

        assertEquals("1", info.getValidatorId());
        assertEquals("0xabc123", info.getAddress());
        assertTrue(info.getActive());
        assertEquals("1000000000000000000000", info.getTotalStake());
        assertEquals("500000000000000000000", info.getSelfStake());
        assertEquals("1609459200", info.getCreatedAtDateUnix());
    }

    @Test
    void fromICPNeuronInfoReply_mapsAllFields() {
        TgvalidatordGetICPNeuronInfoReply reply = new TgvalidatordGetICPNeuronInfoReply();
        reply.setNeuronId("12345678901234567890");
        reply.setNeuronState(GetICPNeuronInfoReplyNeuronState.NEURON_STATE_NOT_DISSOLVING);
        reply.setAgeSeconds("31536000");
        reply.setDissolveDelaySeconds("15552000");
        reply.setVotingPower("100000000");
        reply.setStakeE8S("1000000000");
        reply.setCreatedTimestampSeconds("1609459200");

        ICPNeuronInfo info = StakingMapper.INSTANCE.fromICPNeuronInfoReply(reply);

        assertEquals("12345678901234567890", info.getNeuronId());
        assertEquals("NeuronStateNotDissolving", info.getNeuronState());
        assertEquals("31536000", info.getAgeSeconds());
        assertEquals("15552000", info.getDissolveDelaySeconds());
        assertEquals("100000000", info.getVotingPower());
        assertEquals("1000000000", info.getStakeE8s());
    }

    @Test
    void fromICPNeuronInfoReply_handlesNullState() {
        TgvalidatordGetICPNeuronInfoReply reply = new TgvalidatordGetICPNeuronInfoReply();
        reply.setNeuronId("123");
        reply.setNeuronState(null);

        ICPNeuronInfo info = StakingMapper.INSTANCE.fromICPNeuronInfoReply(reply);

        assertEquals("123", info.getNeuronId());
        assertNull(info.getNeuronState());
    }

    @Test
    void fromNEARValidatorInfoReply_mapsAllFields() {
        TgvalidatordGetNEARValidatorInfoReply reply = new TgvalidatordGetNEARValidatorInfoReply();
        reply.setValidatorAddress("pool.near");
        reply.setOwnerId("owner.near");
        reply.setTotalStakedBalance("1000000000000000000000000000");
        reply.setRewardFeeFraction(0.10f);
        reply.setStakingKey("ed25519:abc123");
        reply.setIsStakingPaused(false);

        NEARValidatorInfo info = StakingMapper.INSTANCE.fromNEARValidatorInfoReply(reply);

        assertEquals("pool.near", info.getValidatorAddress());
        assertEquals("owner.near", info.getOwnerId());
        assertEquals("1000000000000000000000000000", info.getTotalStakedBalance());
        assertEquals(0.10f, info.getRewardFeeFraction());
        assertEquals("ed25519:abc123", info.getStakingKey());
        assertFalse(info.getStakingPaused());
    }

    @Test
    void fromStakeAccountsReply_mapsAccountsAndCursor() {
        TgvalidatordSolanaStakeAccount solanaAccount = new TgvalidatordSolanaStakeAccount();
        solanaAccount.setDerivationIndex("0");
        solanaAccount.setState(SolanaStakeAccountState.ACTIVE);
        solanaAccount.setValidatorAddress("validator1");
        solanaAccount.setActiveBalance("1000000000");
        solanaAccount.setInactiveBalance("0");
        solanaAccount.setAllowMerge(true);

        TgvalidatordStakeAccount dto = new TgvalidatordStakeAccount();
        dto.setId("stake-123");
        dto.setAddressId("address-456");
        dto.setAccountAddress("StakeAccount1...");
        dto.setAccountType(TgvalidatordStakeAccountType.STAKE_ACCOUNT_TYPE_SOLANA);
        dto.setSolanaStakeAccount(solanaAccount);
        dto.setCreatedAt(OffsetDateTime.now());
        dto.setUpdatedAt(OffsetDateTime.now());

        TgvalidatordResponseCursor cursor = new TgvalidatordResponseCursor();
        cursor.setCurrentPage("page-token");
        cursor.setHasNext(true);

        TgvalidatordGetStakeAccountsReply reply = new TgvalidatordGetStakeAccountsReply();
        reply.setStakeAccounts(Arrays.asList(dto));
        reply.setCursor(cursor);

        StakeAccountResult result = StakingMapper.INSTANCE.fromStakeAccountsReply(reply);

        assertNotNull(result);
        assertEquals(1, result.getStakeAccounts().size());

        StakeAccount account = result.getStakeAccounts().get(0);
        assertEquals("stake-123", account.getId());
        assertEquals("address-456", account.getAddressId());
        assertEquals("StakeAccount1...", account.getAccountAddress());
        assertEquals("StakeAccountTypeSolana", account.getAccountType());

        SolanaStakeAccount solana = account.getSolanaStakeAccount();
        assertNotNull(solana);
        assertEquals("0", solana.getDerivationIndex());
        assertEquals("active", solana.getState());
        assertEquals("validator1", solana.getValidatorAddress());
        assertEquals("1000000000", solana.getActiveBalance());

        assertTrue(result.hasNext());
    }

    @Test
    void fromStakeAccountsReply_handlesEmptyList() {
        TgvalidatordGetStakeAccountsReply reply = new TgvalidatordGetStakeAccountsReply();
        reply.setStakeAccounts(Arrays.asList());
        reply.setCursor(null);

        StakeAccountResult result = StakingMapper.INSTANCE.fromStakeAccountsReply(reply);

        assertNotNull(result);
        assertTrue(result.getStakeAccounts().isEmpty());
        assertFalse(result.hasNext());
    }

    @Test
    void fromXTZStakingRewardsReply_mapsField() {
        TgvalidatordGetXTZAddressStakingRewardsReply reply = new TgvalidatordGetXTZAddressStakingRewardsReply();
        reply.setReceivedRewardsAmount("1500000");

        XTZStakingRewards rewards = StakingMapper.INSTANCE.fromXTZStakingRewardsReply(reply);

        assertEquals("1500000", rewards.getReceivedRewardsAmount());
    }

    @Test
    void toAccountTypeString_convertsEnumToString() {
        String result = StakingMapper.INSTANCE.toAccountTypeString(
                TgvalidatordStakeAccountType.STAKE_ACCOUNT_TYPE_SOLANA);
        assertEquals("StakeAccountTypeSolana", result);
    }

    @Test
    void toAccountTypeString_handlesNull() {
        String result = StakingMapper.INSTANCE.toAccountTypeString(null);
        assertNull(result);
    }

    @Test
    void toNeuronStateString_convertsEnumToString() {
        String result = StakingMapper.INSTANCE.toNeuronStateString(
                GetICPNeuronInfoReplyNeuronState.NEURON_STATE_DISSOLVING);
        assertEquals("NeuronStateDissolving", result);
    }

    @Test
    void toNeuronStateString_handlesNull() {
        String result = StakingMapper.INSTANCE.toNeuronStateString(null);
        assertNull(result);
    }

    @Test
    void toSolanaStateString_convertsEnumToString() {
        String result = StakingMapper.INSTANCE.toSolanaStateString(SolanaStakeAccountState.DEACTIVATING);
        assertEquals("deactivating", result);
    }

    @Test
    void toSolanaStateString_handlesNull() {
        String result = StakingMapper.INSTANCE.toSolanaStateString(null);
        assertNull(result);
    }
}
