package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Approvers;
import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAsset;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAssetEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistMetadata;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistTrail;
import com.taurushq.sdk.protect.client.model.WhitelistUserSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApprovers;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddressEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTrail;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistUserSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistedContractAddressAttribute;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

class WhitelistedAssetMapperTest {

    @Test
    void fromSignedContractAddressDTO_mapsFields() {
        TgvalidatordWhitelistSignature sigDto = new TgvalidatordWhitelistSignature();
        sigDto.setHashes(Arrays.asList("hash1"));

        TgvalidatordSignedWhitelistedContractAddress dto =
                new TgvalidatordSignedWhitelistedContractAddress();
        dto.setSignatures(Arrays.asList(sigDto));
        dto.setPayload(new byte[]{1, 2, 3});

        SignedWhitelistedAsset result = WhitelistedAssetMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNotNull(result.getSignatures());
        assertEquals(1, result.getSignatures().size());
    }

    @Test
    void fromSignedContractAddressDTO_handlesNullDto() {
        SignedWhitelistedAsset result = WhitelistedAssetMapper.INSTANCE.fromDTO(
                (TgvalidatordSignedWhitelistedContractAddress) null);
        assertNull(result);
    }

    @Test
    void fromEnvelopeDTO_mapsAllFields() {
        TgvalidatordMetadata metadata = new TgvalidatordMetadata();
        metadata.setHash("hash123");
        metadata.setPayloadAsString("{\"name\":\"USDC\"}");

        TgvalidatordSignedWhitelistedContractAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedContractAddressEnvelope();
        dto.setId("789");
        dto.setTenantId("100");
        dto.setMetadata(metadata);
        dto.setAction("create");
        dto.setStatus("APPROVED");

        SignedWhitelistedAssetEnvelope result = WhitelistedAssetMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(789L, result.getId());
        assertEquals(100L, result.getTenantId());
        assertNotNull(result.getMetadata());
        assertEquals("hash123", result.getMetadata().getHash());
        assertEquals("create", result.getAction());
        assertEquals("APPROVED", result.getStatus());
    }

    @Test
    void fromEnvelopeDTO_handlesNullId() {
        TgvalidatordSignedWhitelistedContractAddressEnvelope dto =
                new TgvalidatordSignedWhitelistedContractAddressEnvelope();

        SignedWhitelistedAssetEnvelope result = WhitelistedAssetMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals(0L, result.getId());
        assertEquals(0L, result.getTenantId());
    }

    @Test
    void fromEnvelopeDTO_handlesNullDto() {
        SignedWhitelistedAssetEnvelope result = WhitelistedAssetMapper.INSTANCE.fromDTO(
                (TgvalidatordSignedWhitelistedContractAddressEnvelope) null);
        assertNull(result);
    }

    @Test
    void fromTrailDTO_mapsFields() {
        OffsetDateTime now = OffsetDateTime.now();

        TgvalidatordTrail dto = new TgvalidatordTrail();
        dto.setId("1");
        dto.setUserId("u1");
        dto.setAction("approve");
        dto.setDate(now);

        WhitelistTrail result = WhitelistedAssetMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("u1", result.getUserId());
        assertEquals("approve", result.getAction());
        assertEquals(now, result.getDate());
    }

    @Test
    void fromTrailDTO_handlesNullDto() {
        WhitelistTrail result = WhitelistedAssetMapper.INSTANCE.fromDTO(
                (TgvalidatordTrail) null);
        assertNull(result);
    }

    @Test
    void fromMetadataDTO_mapsFields() {
        TgvalidatordMetadata dto = new TgvalidatordMetadata();
        dto.setHash("myhash");
        dto.setPayloadAsString("{\"blockchain\":\"ETH\"}");

        WhitelistMetadata result = WhitelistedAssetMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("myhash", result.getHash());
        assertEquals("{\"blockchain\":\"ETH\"}", result.getPayloadAsString());
    }

    @Test
    void fromMetadataDTO_handlesNullDto() {
        WhitelistMetadata result = WhitelistedAssetMapper.INSTANCE.fromDTO(
                (TgvalidatordMetadata) null);
        assertNull(result);
    }

    @Test
    void fromWhitelistSignatureDTO_mapsFields() {
        TgvalidatordWhitelistUserSignature userSig = new TgvalidatordWhitelistUserSignature();
        userSig.setUserId("user-1");

        TgvalidatordWhitelistSignature dto = new TgvalidatordWhitelistSignature();
        dto.setSignature(userSig);
        dto.setHashes(Arrays.asList("h1", "h2"));

        WhitelistSignature result = WhitelistedAssetMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertNotNull(result.getSignature());
        assertEquals(2, result.getHashes().size());
    }

    @Test
    void fromWhitelistSignatureDTO_handlesNullDto() {
        WhitelistSignature result = WhitelistedAssetMapper.INSTANCE.fromDTO(
                (TgvalidatordWhitelistSignature) null);
        assertNull(result);
    }

    @Test
    void fromWhitelistUserSignatureDTO_mapsFields() {
        TgvalidatordWhitelistUserSignature dto = new TgvalidatordWhitelistUserSignature();
        dto.setUserId("user-asset-1");
        dto.setComment("Asset approved");

        WhitelistUserSignature result = WhitelistedAssetMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("user-asset-1", result.getUserId());
        assertEquals("Asset approved", result.getComment());
    }

    @Test
    void fromWhitelistUserSignatureDTO_handlesNullDto() {
        WhitelistUserSignature result = WhitelistedAssetMapper.INSTANCE.fromDTO(
                (TgvalidatordWhitelistUserSignature) null);
        assertNull(result);
    }

    @Test
    void fromApproversDTO_handlesNullDto() {
        Approvers result = WhitelistedAssetMapper.INSTANCE.fromDTO(
                (TgvalidatordApprovers) null);
        assertNull(result);
    }

    @Test
    void fromAttributeDTO_mapsFields() {
        TgvalidatordWhitelistedContractAddressAttribute dto =
                new TgvalidatordWhitelistedContractAddressAttribute();
        dto.setId("1");
        dto.setKey("chain");
        dto.setValue("ethereum");
        dto.setSubtype("metadata");
        dto.setIsfile(false);

        Attribute result = WhitelistedAssetMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("chain", result.getKey());
        assertEquals("ethereum", result.getValue());
        assertEquals("metadata", result.getSubType());
    }

    @Test
    void fromAttributeDTO_handlesNullDto() {
        Attribute result = WhitelistedAssetMapper.INSTANCE.fromDTO(
                (TgvalidatordWhitelistedContractAddressAttribute) null);
        assertNull(result);
    }
}
