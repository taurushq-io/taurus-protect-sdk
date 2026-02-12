package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingAgreement;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingOffer;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Participant;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Pledge;
import com.taurushq.sdk.protect.client.model.taurusnetwork.PledgeResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.PledgeWithdrawal;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Settlement;
import com.taurushq.sdk.protect.client.model.taurusnetwork.SharedAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPledgesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordLendingAgreement;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnLendingOffer;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnParticipant;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnPledge;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnPledgeWithdrawal;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnSettlement;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnSharedAddress;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

class TaurusNetworkMapperTest {

    @Test
    void fromParticipantDTO_mapsFields() {
        TgvalidatordTnParticipant dto = new TgvalidatordTnParticipant();
        dto.setId("participant-123");
        dto.setName("Test Participant");
        dto.setCountry("CH");

        Participant result = TaurusNetworkMapper.INSTANCE.fromParticipantDTO(dto);

        assertNotNull(result);
        assertEquals("participant-123", result.getId());
        assertEquals("Test Participant", result.getName());
        assertEquals("CH", result.getCountry());
    }

    @Test
    void fromParticipantDTO_handlesNull() {
        Participant result = TaurusNetworkMapper.INSTANCE.fromParticipantDTO(null);
        assertNull(result);
    }

    @Test
    void fromParticipantDTOList_mapsList() {
        TgvalidatordTnParticipant dto1 = new TgvalidatordTnParticipant();
        dto1.setId("participant-1");

        TgvalidatordTnParticipant dto2 = new TgvalidatordTnParticipant();
        dto2.setId("participant-2");

        List<Participant> result = TaurusNetworkMapper.INSTANCE.fromParticipantDTOList(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
    }

    @Test
    void fromPledgeDTO_mapsFields() {
        TgvalidatordTnPledge dto = new TgvalidatordTnPledge();
        dto.setId("pledge-123");
        dto.setOwnerParticipantID("owner-1");
        dto.setTargetParticipantID("target-1");
        dto.setAmount("1000");
        dto.setStatus("ACTIVE");

        Pledge result = TaurusNetworkMapper.INSTANCE.fromPledgeDTO(dto);

        assertNotNull(result);
        assertEquals("pledge-123", result.getId());
        assertEquals("owner-1", result.getOwnerParticipantID());
        assertEquals("target-1", result.getTargetParticipantID());
        assertEquals("1000", result.getAmount());
        assertEquals("ACTIVE", result.getStatus());
    }

    @Test
    void fromPledgeDTO_handlesNull() {
        Pledge result = TaurusNetworkMapper.INSTANCE.fromPledgeDTO(null);
        assertNull(result);
    }

    @Test
    void fromPledgesReply_mapsPledgesAndCursor() {
        TgvalidatordTnPledge pledge = new TgvalidatordTnPledge();
        pledge.setId("pledge-123");

        TgvalidatordResponseCursor cursor = new TgvalidatordResponseCursor();
        cursor.setCurrentPage("page-1");

        TgvalidatordGetPledgesReply reply = new TgvalidatordGetPledgesReply();
        reply.setPledges(Arrays.asList(pledge));
        reply.setCursor(cursor);

        PledgeResult result = TaurusNetworkMapper.INSTANCE.fromPledgesReply(reply);

        assertNotNull(result);
        assertNotNull(result.getPledges());
        assertEquals(1, result.getPledges().size());
        assertNotNull(result.getCursor());
    }

    @Test
    void fromPledgeWithdrawalDTO_mapsFields() {
        TgvalidatordTnPledgeWithdrawal dto = new TgvalidatordTnPledgeWithdrawal();
        dto.setId("withdrawal-123");
        dto.setPledgeID("pledge-1");
        dto.setAmount("500");
        dto.setStatus("COMPLETED");

        PledgeWithdrawal result = TaurusNetworkMapper.INSTANCE.fromPledgeWithdrawalDTO(dto);

        assertNotNull(result);
        assertEquals("withdrawal-123", result.getId());
        assertEquals("pledge-1", result.getPledgeID());
        assertEquals("500", result.getAmount());
        assertEquals("COMPLETED", result.getStatus());
    }

    @Test
    void fromSharedAddressDTO_mapsFields() {
        TgvalidatordTnSharedAddress dto = new TgvalidatordTnSharedAddress();
        dto.setId("shared-123");
        dto.setBlockchain("ETH");
        dto.setNetwork("mainnet");
        dto.setAddress("0x1234");

        SharedAddress result = TaurusNetworkMapper.INSTANCE.fromSharedAddressDTO(dto);

        assertNotNull(result);
        assertEquals("shared-123", result.getId());
        assertEquals("ETH", result.getBlockchain());
        assertEquals("mainnet", result.getNetwork());
        assertEquals("0x1234", result.getAddress());
    }

    @Test
    void fromSettlementDTO_mapsFields() {
        TgvalidatordTnSettlement dto = new TgvalidatordTnSettlement();
        dto.setId("settlement-123");
        dto.setCreatorParticipantID("creator-1");
        dto.setTargetParticipantID("target-1");
        dto.setStatus("COMPLETED");

        Settlement result = TaurusNetworkMapper.INSTANCE.fromSettlementDTO(dto);

        assertNotNull(result);
        assertEquals("settlement-123", result.getId());
        assertEquals("creator-1", result.getCreatorParticipantID());
        assertEquals("target-1", result.getTargetParticipantID());
        assertEquals("COMPLETED", result.getStatus());
    }

    @Test
    void fromLendingOfferDTO_mapsFields() {
        TgvalidatordTnLendingOffer dto = new TgvalidatordTnLendingOffer();
        dto.setId("offer-123");
        dto.setParticipantID("participant-1");
        dto.setBlockchain("ETH");
        dto.setAmount("1000");

        LendingOffer result = TaurusNetworkMapper.INSTANCE.fromLendingOfferDTO(dto);

        assertNotNull(result);
        assertEquals("offer-123", result.getId());
        assertEquals("participant-1", result.getParticipantID());
        assertEquals("ETH", result.getBlockchain());
        assertEquals("1000", result.getAmount());
    }

    @Test
    void fromLendingAgreementDTO_mapsFields() {
        TgvalidatordLendingAgreement dto = new TgvalidatordLendingAgreement();
        dto.setId("agreement-123");
        dto.setLenderParticipantID("lender-1");
        dto.setBorrowerParticipantID("borrower-1");
        dto.setAmount("5000");
        dto.setStatus("ACTIVE");

        LendingAgreement result = TaurusNetworkMapper.INSTANCE.fromLendingAgreementDTO(dto);

        assertNotNull(result);
        assertEquals("agreement-123", result.getId());
        assertEquals("lender-1", result.getLenderParticipantID());
        assertEquals("borrower-1", result.getBorrowerParticipantID());
        assertEquals("5000", result.getAmount());
        assertEquals("ACTIVE", result.getStatus());
    }
}
