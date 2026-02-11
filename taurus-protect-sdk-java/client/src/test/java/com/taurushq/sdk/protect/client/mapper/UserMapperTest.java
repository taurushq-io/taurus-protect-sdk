package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.User;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalUser;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class UserMapperTest {

    @Test
    void fromDTO_withCompleteData_mapsAllFields() {
        // Given
        TgvalidatordInternalUser dto = new TgvalidatordInternalUser();
        dto.setId("user-123");
        dto.setTenantId("42");
        dto.setExternalUserId("ext-456");
        dto.setFirstName("John");
        dto.setLastName("Doe");
        dto.setStatus("ACTIVE");
        dto.setEmail("john.doe@example.com");
        dto.setUsername("johndoe");
        dto.setPublicKey("pubkey-789");
        dto.setRoles(Arrays.asList("ADMIN", "USER"));
        dto.setTotpEnabled(true);
        dto.setEnforcedInRules(false);

        OffsetDateTime creationDate = OffsetDateTime.now();
        OffsetDateTime updateDate = creationDate.plusDays(1);
        OffsetDateTime lastLogin = creationDate.plusHours(2);
        dto.setCreationDate(creationDate);
        dto.setUpdateDate(updateDate);
        dto.setLastLogin(lastLogin);

        // When
        User result = UserMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals("user-123", result.getId());
        assertEquals(42, result.getTenantId());
        assertEquals("ext-456", result.getExternalUserId());
        assertEquals("John", result.getFirstName());
        assertEquals("Doe", result.getLastName());
        assertEquals("ACTIVE", result.getStatus());
        assertEquals("john.doe@example.com", result.getEmail());
        assertEquals("johndoe", result.getUsername());
        assertEquals("pubkey-789", result.getPublicKey());
        assertEquals(2, result.getRoles().size());
        assertTrue(result.getRoles().contains("ADMIN"));
        assertTrue(result.getRoles().contains("USER"));
        assertTrue(result.getTotpEnabled());
        assertFalse(result.getEnforcedInRules());
        assertEquals(creationDate, result.getCreatedAt());
        assertEquals(updateDate, result.getUpdatedAt());
        assertEquals(lastLogin, result.getLastLogin());
    }

    @Test
    void fromDTO_withNullTenantId_defaultsToZero() {
        // Given
        TgvalidatordInternalUser dto = new TgvalidatordInternalUser();
        dto.setId("user-1");
        dto.setTenantId(null);

        // When
        User result = UserMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(0, result.getTenantId());
    }

    @Test
    void fromDTO_withNullOptionalFields_handlesGracefully() {
        // Given
        TgvalidatordInternalUser dto = new TgvalidatordInternalUser();
        dto.setId("user-1");
        // All other fields left null

        // When
        User result = UserMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals("user-1", result.getId());
        assertNull(result.getEmail());
        assertNull(result.getFirstName());
        assertNull(result.getLastName());
        assertNull(result.getCreatedAt());
    }

    @Test
    void fromDTO_list_mapsMultipleUsers() {
        // Given
        TgvalidatordInternalUser dto1 = new TgvalidatordInternalUser();
        dto1.setId("user-1");
        dto1.setEmail("user1@example.com");

        TgvalidatordInternalUser dto2 = new TgvalidatordInternalUser();
        dto2.setId("user-2");
        dto2.setEmail("user2@example.com");

        List<TgvalidatordInternalUser> dtos = Arrays.asList(dto1, dto2);

        // When
        List<User> results = UserMapper.INSTANCE.fromDTO(dtos);

        // Then
        assertNotNull(results);
        assertEquals(2, results.size());
        assertEquals("user-1", results.get(0).getId());
        assertEquals("user1@example.com", results.get(0).getEmail());
        assertEquals("user-2", results.get(1).getId());
        assertEquals("user2@example.com", results.get(1).getEmail());
    }

    @Test
    void fromDTO_mapsDateFieldsCorrectly() {
        // Given - verify creationDate → createdAt and updateDate → updatedAt
        TgvalidatordInternalUser dto = new TgvalidatordInternalUser();
        dto.setId("user-1");

        OffsetDateTime creationDate = OffsetDateTime.parse("2024-01-15T10:30:00Z");
        OffsetDateTime updateDate = OffsetDateTime.parse("2024-01-16T11:45:00Z");
        dto.setCreationDate(creationDate);
        dto.setUpdateDate(updateDate);

        // When
        User result = UserMapper.INSTANCE.fromDTO(dto);

        // Then
        assertEquals(creationDate, result.getCreatedAt());
        assertEquals(updateDate, result.getUpdatedAt());
    }
}
