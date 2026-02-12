package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestStatus;
import com.taurushq.sdk.protect.client.model.SignedRequest;
import com.taurushq.sdk.protect.openapi.model.RequestSignedRequest;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequest;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class RequestMapperTest {

    @Test
    void fromDTO_withCompleteData_mapsAllFields() {
        // Given
        TgvalidatordRequest dto = new TgvalidatordRequest();
        dto.setId("12345");
        dto.setTenantId("99");
        dto.setCurrency("ETH");
        dto.setEnvelope("envelope-data");
        dto.setStatus("CREATED");
        dto.setRule("rule-123");
        dto.setType("TRANSFER");
        dto.setRequestBundleId("bundle-456");
        dto.setExternalRequestId("ext-789");

        OffsetDateTime now = OffsetDateTime.now();
        dto.setCreationDate(now);
        dto.setUpdateDate(now);

        // When
        Request result = RequestMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(12345L, result.getId());
        assertEquals(99, result.getTenantId());
        assertEquals("ETH", result.getCurrency());
        assertEquals("envelope-data", result.getEnvelope());
        assertEquals(RequestStatus.CREATED, result.getStatus());
        assertEquals("rule-123", result.getRule());
        assertEquals("TRANSFER", result.getType());
        assertEquals("bundle-456", result.getRequestBundleId());
        assertEquals("ext-789", result.getExternalRequestId());
        assertEquals(now, result.getCreationDate());
        assertEquals(now, result.getUpdateDate());
    }

    @Test
    void fromDTO_withStatusApproved_mapsStatusCorrectly() {
        // Given
        TgvalidatordRequest dto = new TgvalidatordRequest();
        dto.setId("1");
        dto.setStatus("APPROVED");

        // When
        Request result = RequestMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(RequestStatus.APPROVED, result.getStatus());
    }

    @Test
    void fromDTO_withStatusRejected_mapsStatusCorrectly() {
        // Given
        TgvalidatordRequest dto = new TgvalidatordRequest();
        dto.setId("1");
        dto.setStatus("rejected");

        // When
        Request result = RequestMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(RequestStatus.REJECTED, result.getStatus());
    }

    @Test
    void fromDTO_withUnknownStatus_throws() {
        // Given
        TgvalidatordRequest dto = new TgvalidatordRequest();
        dto.setId("1");
        dto.setStatus("INVALID_STATUS");

        // When / Then
        assertThrows(IllegalArgumentException.class, () -> RequestMapper.INSTANCE.fromDTO(dto));
    }

    @Test
    void fromDTO_withNullStatus_throws() {
        // Given
        TgvalidatordRequest dto = new TgvalidatordRequest();
        dto.setId("1");
        dto.setStatus(null);

        // When / Then
        assertThrows(IllegalArgumentException.class, () -> RequestMapper.INSTANCE.fromDTO(dto));
    }

    @Test
    void fromDTO_withNullOptionalFields_handlesGracefully() {
        // Given
        TgvalidatordRequest dto = new TgvalidatordRequest();
        dto.setId("1");
        dto.setStatus("CREATED");

        // When
        Request result = RequestMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals(1L, result.getId());
        assertNull(result.getCurrency());
        assertNull(result.getEnvelope());
        assertNull(result.getRule());
        assertNull(result.getCreationDate());
        assertNull(result.getMetadata());
    }

    @Test
    void fromDTO_signedRequest_mapsCorrectly() {
        // Given
        RequestSignedRequest dto = new RequestSignedRequest();
        dto.setHash("abc123hash");
        dto.setStatus("HSM_SIGNED");

        // When
        SignedRequest result = RequestMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals("abc123hash", result.getHash());
        assertEquals(RequestStatus.HSM_SIGNED, result.getStatus());
    }

    @Test
    void fromDTO_signedRequest_withConfirmedStatus() {
        // Given
        RequestSignedRequest dto = new RequestSignedRequest();
        dto.setHash("xyz789hash");
        dto.setStatus("CONFIRMED");

        // When
        SignedRequest result = RequestMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(RequestStatus.CONFIRMED, result.getStatus());
    }

    @Test
    void fromDTO_withNeedsApprovalFrom_mapsList() {
        // Given
        TgvalidatordRequest dto = new TgvalidatordRequest();
        dto.setId("1");
        dto.setStatus("CREATED");
        dto.setNeedsApprovalFrom(Arrays.asList("user1", "user2", "user3"));

        // When
        Request result = RequestMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result.getNeedsApprovalFrom());
        assertEquals(3, result.getNeedsApprovalFrom().size());
        assertTrue(result.getNeedsApprovalFrom().contains("user1"));
        assertTrue(result.getNeedsApprovalFrom().contains("user2"));
        assertTrue(result.getNeedsApprovalFrom().contains("user3"));
    }

    @Test
    void fromDTO_withMetadata_mapsMetadata() {
        // Given
        TgvalidatordRequest dto = new TgvalidatordRequest();
        dto.setId("1");
        dto.setStatus("CREATED");

        TgvalidatordMetadata metadata = new TgvalidatordMetadata();
        metadata.setHash("metadataHash123");
        // payloadAsString needs to be valid JSON for RequestMetadata.setPayloadAsString to work
        metadata.setPayloadAsString("[]");
        dto.setMetadata(metadata);

        // When
        Request result = RequestMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result.getMetadata());
        assertEquals("metadataHash123", result.getMetadata().getHash());
    }
}
