package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Approvers;
import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddress;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistMetadata;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistTrail;
import com.taurushq.sdk.protect.client.model.WhitelistUserSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApprovers;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTrail;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistUserSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistedAddressAttribute;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class WhitelistedAddressMapperTest {

    @Test
    void fromSignedWhitelistedAddressDTO_mapsFields() {
        TgvalidatordWhitelistSignature sigDto = new TgvalidatordWhitelistSignature();
        sigDto.setHashes(Arrays.asList("hash1", "hash2"));

        TgvalidatordSignedWhitelistedAddress dto = new TgvalidatordSignedWhitelistedAddress();
        dto.setSignatures(Arrays.asList(sigDto));
        dto.setPayload(new byte[]{1, 2, 3});

        SignedWhitelistedAddress result = WhitelistedAddressMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNotNull(result.getSignatures());
        assertEquals(1, result.getSignatures().size());
    }

    @Test
    void fromSignedWhitelistedAddressDTO_handlesNullDto() {
        SignedWhitelistedAddress result = WhitelistedAddressMapper.INSTANCE.fromDTO(
                (TgvalidatordSignedWhitelistedAddress) null);
        assertNull(result);
    }

    @Test
    void fromEnvelopeDTO_mapsAllFields() {
        TgvalidatordSignedWhitelistedAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedAddressEnvelope();
        dto.setId("123");
        dto.setTenantId("456");
        dto.setAction("create");
        dto.setRulesContainer("base64rules==");
        dto.setRule("default");
        dto.setStatus("PENDING");
        dto.setBlockchain("ETH");
        dto.setNetwork("mainnet");

        TgvalidatordMetadata metadata = new TgvalidatordMetadata();
        metadata.setHash("abc123");
        metadata.setPayloadAsString("{\"address\":\"0x123\"}");
        dto.setMetadata(metadata);

        SignedWhitelistedAddressEnvelope result = WhitelistedAddressMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(123L, result.getId());
        assertEquals(456L, result.getTenantId());
        assertEquals("create", result.getAction());
        assertEquals("base64rules==", result.getRulesContainer());
        assertEquals("default", result.getRule());
        assertEquals("PENDING", result.getStatus());
        assertEquals("ETH", result.getBlockchain());
        assertEquals("mainnet", result.getNetwork());
        assertNotNull(result.getMetadata());
        assertEquals("abc123", result.getMetadata().getHash());
    }

    @Test
    void fromEnvelopeDTO_handlesNullId() {
        TgvalidatordSignedWhitelistedAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedAddressEnvelope();

        SignedWhitelistedAddressEnvelope result = WhitelistedAddressMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(0L, result.getId());
        assertEquals(0L, result.getTenantId());
    }

    @Test
    void fromEnvelopeDTO_handlesNullDto() {
        SignedWhitelistedAddressEnvelope result = WhitelistedAddressMapper.INSTANCE.fromDTO(
                (TgvalidatordSignedWhitelistedAddressEnvelope) null);
        assertNull(result);
    }

    @Test
    void fromWhitelistSignatureDTO_mapsFields() {
        TgvalidatordWhitelistUserSignature userSig = new TgvalidatordWhitelistUserSignature();
        userSig.setUserId("user-1");
        userSig.setSignature(new byte[]{1, 2, 3});
        userSig.setComment("Approved");

        TgvalidatordWhitelistSignature dto = new TgvalidatordWhitelistSignature();
        dto.setSignature(userSig);
        dto.setHashes(Arrays.asList("hash1", "hash2"));

        WhitelistSignature result = WhitelistedAddressMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNotNull(result.getSignature());
        assertNotNull(result.getHashes());
        assertEquals(2, result.getHashes().size());
    }

    @Test
    void fromWhitelistSignatureDTO_handlesNullDto() {
        WhitelistSignature result = WhitelistedAddressMapper.INSTANCE.fromDTO(
                (TgvalidatordWhitelistSignature) null);
        assertNull(result);
    }

    @Test
    void fromWhitelistUserSignatureDTO_mapsFields() {
        TgvalidatordWhitelistUserSignature dto = new TgvalidatordWhitelistUserSignature();
        dto.setUserId("user-123");
        dto.setSignature(new byte[]{10, 20, 30});
        dto.setComment("Looks good");

        WhitelistUserSignature result = WhitelistedAddressMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("user-123", result.getUserId());
        assertEquals("Looks good", result.getComment());
    }

    @Test
    void fromWhitelistUserSignatureDTO_handlesNullDto() {
        WhitelistUserSignature result = WhitelistedAddressMapper.INSTANCE.fromDTO(
                (TgvalidatordWhitelistUserSignature) null);
        assertNull(result);
    }

    @Test
    void fromTrailDTO_mapsFields() {
        OffsetDateTime now = OffsetDateTime.now();

        TgvalidatordTrail dto = new TgvalidatordTrail();
        dto.setId("1");
        dto.setUserId("user-1");
        dto.setExternalUserId("ext-1");
        dto.setAction("approve");
        dto.setComment("Approved whitelist");
        dto.setDate(now);

        WhitelistTrail result = WhitelistedAddressMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("user-1", result.getUserId());
        assertEquals("ext-1", result.getExternalUserId());
        assertEquals("approve", result.getAction());
        assertEquals("Approved whitelist", result.getComment());
        assertEquals(now, result.getDate());
    }

    @Test
    void fromTrailDTO_handlesNullDto() {
        WhitelistTrail result = WhitelistedAddressMapper.INSTANCE.fromDTO(
                (TgvalidatordTrail) null);
        assertNull(result);
    }

    @Test
    void fromMetadataDTO_mapsFields() {
        TgvalidatordMetadata dto = new TgvalidatordMetadata();
        dto.setHash("sha256hash");
        dto.setPayloadAsString("{\"address\":\"0xabc\"}");

        WhitelistMetadata result = WhitelistedAddressMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("sha256hash", result.getHash());
        assertEquals("{\"address\":\"0xabc\"}", result.getPayloadAsString());
    }

    @Test
    void fromMetadataDTO_handlesNullDto() {
        WhitelistMetadata result = WhitelistedAddressMapper.INSTANCE.fromDTO(
                (TgvalidatordMetadata) null);
        assertNull(result);
    }

    @Test
    void fromApproversDTO_handlesNullDto() {
        Approvers result = WhitelistedAddressMapper.INSTANCE.fromDTO(
                (TgvalidatordApprovers) null);
        assertNull(result);
    }

    @Test
    void fromAttributeDTO_mapsFields() {
        TgvalidatordWhitelistedAddressAttribute dto = new TgvalidatordWhitelistedAddressAttribute();
        dto.setId("1");
        dto.setKey("note");
        dto.setValue("important");
        dto.setSubtype("label");
        dto.setIsfile(false);

        Attribute result = WhitelistedAddressMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("note", result.getKey());
        assertEquals("important", result.getValue());
        assertEquals("label", result.getSubType());
    }

    @Test
    void fromAttributeDTO_handlesNullDto() {
        Attribute result = WhitelistedAddressMapper.INSTANCE.fromDTO(
                (TgvalidatordWhitelistedAddressAttribute) null);
        assertNull(result);
    }
}
