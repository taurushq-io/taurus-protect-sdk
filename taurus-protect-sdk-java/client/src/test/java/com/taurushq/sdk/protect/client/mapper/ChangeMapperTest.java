package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Change;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordChange;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;

class ChangeMapperTest {

    @Test
    void fromDTO_withCompleteData_mapsAllFields() {
        // Given
        TgvalidatordChange dto = new TgvalidatordChange();
        dto.setId("change-123");
        dto.setTenantId("42");
        dto.setCreatorId("creator-456");
        dto.setCreatorExternalId("ext-creator-789");
        dto.setAction("UPDATE");
        dto.setEntity("WALLET");
        dto.setEntityId("entity-111");
        dto.setEntityUUID(UUID.fromString("550e8400-e29b-41d4-a716-446655440000"));
        dto.setComment("Test change comment");

        Map<String, String> changes = new HashMap<>();
        changes.put("name", "New Wallet Name");
        changes.put("status", "ACTIVE");
        dto.setChanges(changes);

        OffsetDateTime creationDate = OffsetDateTime.now();
        dto.setCreationDate(creationDate);

        // When
        Change result = ChangeMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals("change-123", result.getId());
        assertEquals(42, result.getTenantId());
        assertEquals("creator-456", result.getCreatorId());
        assertEquals("ext-creator-789", result.getCreatorExternalId());
        assertEquals("UPDATE", result.getAction());
        assertEquals("WALLET", result.getEntity());
        assertEquals("entity-111", result.getEntityId());
        assertEquals("550e8400-e29b-41d4-a716-446655440000", result.getEntityUUID());
        assertEquals("Test change comment", result.getComment());
        assertEquals(creationDate, result.getCreatedAt());
        assertNotNull(result.getChanges());
        assertEquals("New Wallet Name", result.getChanges().get("name"));
        assertEquals("ACTIVE", result.getChanges().get("status"));
    }

    @Test
    void fromDTO_withNullTenantId_defaultsToZero() {
        // Given
        TgvalidatordChange dto = new TgvalidatordChange();
        dto.setId("change-1");
        dto.setTenantId(null);

        // When
        Change result = ChangeMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(0, result.getTenantId());
    }

    @Test
    void fromDTO_withNullOptionalFields_handlesGracefully() {
        // Given
        TgvalidatordChange dto = new TgvalidatordChange();
        dto.setId("change-1");
        // All other fields left null

        // When
        Change result = ChangeMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals("change-1", result.getId());
        assertNull(result.getAction());
        assertNull(result.getEntity());
        assertNull(result.getComment());
        assertNull(result.getCreatedAt());
    }

    @Test
    void fromDTO_list_mapsMultipleChanges() {
        // Given
        TgvalidatordChange dto1 = new TgvalidatordChange();
        dto1.setId("change-1");
        dto1.setAction("CREATE");

        TgvalidatordChange dto2 = new TgvalidatordChange();
        dto2.setId("change-2");
        dto2.setAction("DELETE");

        List<TgvalidatordChange> dtos = Arrays.asList(dto1, dto2);

        // When
        List<Change> results = ChangeMapper.INSTANCE.fromDTO(dtos);

        // Then
        assertNotNull(results);
        assertEquals(2, results.size());
        assertEquals("change-1", results.get(0).getId());
        assertEquals("CREATE", results.get(0).getAction());
        assertEquals("change-2", results.get(1).getId());
        assertEquals("DELETE", results.get(1).getAction());
    }

    @Test
    void fromDTO_mapsCreationDateToCreatedAt() {
        // Given
        TgvalidatordChange dto = new TgvalidatordChange();
        dto.setId("change-1");

        OffsetDateTime creationDate = OffsetDateTime.parse("2024-06-15T14:30:00Z");
        dto.setCreationDate(creationDate);

        // When
        Change result = ChangeMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(creationDate, result.getCreatedAt());
    }
}
