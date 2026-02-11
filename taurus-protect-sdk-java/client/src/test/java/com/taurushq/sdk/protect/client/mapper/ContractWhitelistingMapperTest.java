package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedContractAddress;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedContractAddressEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistedContractAddressResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddressEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTrail;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistedContractAddressAttribute;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ContractWhitelistingMapperTest {

    @Test
    void fromEnvelopeDTO_mapsAllFields() {
        TgvalidatordSignedWhitelistedContractAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedContractAddressEnvelope();
        dto.setId("123");
        dto.setTenantId("456");
        dto.setAction("create");
        dto.setStatus("pending");
        dto.setBlockchain("ETH");
        dto.setNetwork("mainnet");
        dto.setRule("rule-1");
        dto.setRulesContainer("container-1");
        dto.setBusinessRuleEnabled(true);
        dto.setRulesSignatures("sig123");

        SignedWhitelistedContractAddressEnvelope envelope =
                ContractWhitelistingMapper.INSTANCE.fromEnvelopeDTO(dto);

        assertEquals("123", envelope.getId());
        assertEquals("456", envelope.getTenantId());
        assertEquals("create", envelope.getAction());
        assertEquals("pending", envelope.getStatus());
        assertEquals("ETH", envelope.getBlockchain());
        assertEquals("mainnet", envelope.getNetwork());
        assertEquals("rule-1", envelope.getRule());
        assertEquals("container-1", envelope.getRulesContainer());
        assertTrue(envelope.getBusinessRuleEnabled());
        assertEquals("sig123", envelope.getRulesSignatures());
    }

    @Test
    void fromEnvelopeDTO_mapsSignedContractAddress() {
        TgvalidatordSignedWhitelistedContractAddress signedAddr =
                new TgvalidatordSignedWhitelistedContractAddress();
        signedAddr.setPayload(new byte[]{1, 2, 3});

        TgvalidatordSignedWhitelistedContractAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedContractAddressEnvelope();
        dto.setId("123");
        dto.setSignedContractAddress(signedAddr);

        SignedWhitelistedContractAddressEnvelope envelope =
                ContractWhitelistingMapper.INSTANCE.fromEnvelopeDTO(dto);

        assertNotNull(envelope.getSignedContractAddress());
        assertEquals(3, envelope.getSignedContractAddress().getPayload().length);
    }

    @Test
    void fromEnvelopeDTO_mapsTrails() {
        TgvalidatordTrail trail = new TgvalidatordTrail();
        trail.setId("456");
        trail.setUserId("user-123");
        trail.setAction("approve");
        trail.setComment("Approved");
        trail.setDate(OffsetDateTime.now());

        TgvalidatordSignedWhitelistedContractAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedContractAddressEnvelope();
        dto.setId("123");
        dto.setTrails(Arrays.asList(trail));

        SignedWhitelistedContractAddressEnvelope envelope =
                ContractWhitelistingMapper.INSTANCE.fromEnvelopeDTO(dto);

        assertNotNull(envelope.getTrails());
        assertEquals(1, envelope.getTrails().size());
        assertEquals("user-123", envelope.getTrails().get(0).getUserId());
        assertEquals("approve", envelope.getTrails().get(0).getAction());
    }

    @Test
    void fromEnvelopesReply_mapsList() {
        TgvalidatordSignedWhitelistedContractAddressEnvelope dto1 =
                new TgvalidatordSignedWhitelistedContractAddressEnvelope();
        dto1.setId("1");
        dto1.setBlockchain("ETH");

        TgvalidatordSignedWhitelistedContractAddressEnvelope dto2 =
                new TgvalidatordSignedWhitelistedContractAddressEnvelope();
        dto2.setId("2");
        dto2.setBlockchain("MATIC");

        TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply reply =
                new TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply();
        reply.setResult(Arrays.asList(dto1, dto2));
        reply.setTotalItems("100");

        WhitelistedContractAddressResult result =
                ContractWhitelistingMapper.INSTANCE.fromEnvelopesReply(reply);

        assertNotNull(result);
        assertEquals(2, result.getContracts().size());
        assertEquals("1", result.getContracts().get(0).getId());
        assertEquals("ETH", result.getContracts().get(0).getBlockchain());
        assertEquals("2", result.getContracts().get(1).getId());
        assertEquals("MATIC", result.getContracts().get(1).getBlockchain());
        assertEquals(100L, result.getTotalItems());
    }

    @Test
    void fromEnvelopesReply_handlesPagination() {
        TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply reply =
                new TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply();
        reply.setResult(Arrays.asList());
        reply.setTotalItems("150");

        WhitelistedContractAddressResult result =
                ContractWhitelistingMapper.INSTANCE.fromEnvelopesReply(reply);

        assertEquals(150L, result.getTotalItems());
        assertTrue(result.hasMore(0, 50));
        assertTrue(result.hasMore(50, 50));
        assertFalse(result.hasMore(100, 50));
    }

    @Test
    void fromSignedContractDTO_mapsPayload() {
        TgvalidatordSignedWhitelistedContractAddress dto =
                new TgvalidatordSignedWhitelistedContractAddress();
        byte[] payload = {0x01, 0x02, 0x03, 0x04};
        dto.setPayload(payload);

        SignedWhitelistedContractAddress result =
                ContractWhitelistingMapper.INSTANCE.fromSignedContractDTO(dto);

        assertNotNull(result);
        assertEquals(4, result.getPayload().length);
    }

    @Test
    void fromAttributeDTO_mapsAllFields() {
        TgvalidatordWhitelistedContractAddressAttribute dto =
                new TgvalidatordWhitelistedContractAddressAttribute();
        dto.setId("123");
        dto.setKey("description");
        dto.setValue("A test token");
        dto.setContentType("text/plain");
        dto.setOwner("user-123");
        dto.setType("metadata");
        dto.setSubtype("basic");
        dto.setIsfile(false);

        Attribute attr = ContractWhitelistingMapper.INSTANCE.fromAttributeDTO(dto);

        assertEquals("description", attr.getKey());
        assertEquals("A test token", attr.getValue());
        assertEquals("text/plain", attr.getContentType());
        assertEquals("user-123", attr.getOwner());
        assertEquals("metadata", attr.getType());
        assertEquals("basic", attr.getSubType());
        assertFalse(attr.isFile());
    }

    @Test
    void fromAttributeDTOList_mapsList() {
        TgvalidatordWhitelistedContractAddressAttribute dto1 =
                new TgvalidatordWhitelistedContractAddressAttribute();
        dto1.setId("1");
        dto1.setKey("key1");

        TgvalidatordWhitelistedContractAddressAttribute dto2 =
                new TgvalidatordWhitelistedContractAddressAttribute();
        dto2.setId("2");
        dto2.setKey("key2");

        List<Attribute> attrs = ContractWhitelistingMapper.INSTANCE
                .fromAttributeDTOList(Arrays.asList(dto1, dto2));

        assertEquals(2, attrs.size());
        assertEquals("key1", attrs.get(0).getKey());
        assertEquals("key2", attrs.get(1).getKey());
    }

    @Test
    void stringToLong_convertsValidString() {
        long result = ContractWhitelistingMapper.INSTANCE.stringToLong("12345");
        assertEquals(12345L, result);
    }

    @Test
    void stringToLong_handlesNull() {
        long result = ContractWhitelistingMapper.INSTANCE.stringToLong(null);
        assertEquals(0L, result);
    }

    @Test
    void stringToLong_handlesEmptyString() {
        long result = ContractWhitelistingMapper.INSTANCE.stringToLong("");
        assertEquals(0L, result);
    }

    @Test
    void fromEnvelopeDTO_handlesNullFields() {
        TgvalidatordSignedWhitelistedContractAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedContractAddressEnvelope();
        dto.setId("123");
        // Leave all other fields null

        SignedWhitelistedContractAddressEnvelope envelope =
                ContractWhitelistingMapper.INSTANCE.fromEnvelopeDTO(dto);

        assertEquals("123", envelope.getId());
        assertNull(envelope.getTenantId());
        assertNull(envelope.getAction());
        assertNull(envelope.getSignedContractAddress());
        assertNull(envelope.getMetadata());
        assertNull(envelope.getApprovers());
    }
}
