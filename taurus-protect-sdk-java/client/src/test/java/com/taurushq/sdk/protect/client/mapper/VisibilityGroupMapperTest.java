package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.VisibilityGroup;
import com.taurushq.sdk.protect.client.model.VisibilityGroupUser;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalVisibilityGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalVisibilityGroupUser;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class VisibilityGroupMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime creationDate = OffsetDateTime.now();
        OffsetDateTime updateDate = OffsetDateTime.now().plusHours(1);

        TgvalidatordInternalVisibilityGroup dto = new TgvalidatordInternalVisibilityGroup();
        dto.setId("vg-123");
        dto.setTenantId("tenant-456");
        dto.setName("Test Group");
        dto.setDescription("A test visibility group");
        dto.setCreationDate(creationDate);
        dto.setUpdateDate(updateDate);
        dto.setUserCount("5");

        VisibilityGroup result = VisibilityGroupMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("vg-123", result.getId());
        assertEquals("tenant-456", result.getTenantId());
        assertEquals("Test Group", result.getName());
        assertEquals("A test visibility group", result.getDescription());
        assertEquals(creationDate, result.getCreationDate());
        assertEquals(updateDate, result.getUpdateDate());
        assertEquals("5", result.getUserCount());
    }

    @Test
    void fromDTO_handlesNullDto() {
        VisibilityGroup result = VisibilityGroupMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordInternalVisibilityGroup dto1 = new TgvalidatordInternalVisibilityGroup();
        dto1.setId("vg-1");
        dto1.setName("Group 1");

        TgvalidatordInternalVisibilityGroup dto2 = new TgvalidatordInternalVisibilityGroup();
        dto2.setId("vg-2");
        dto2.setName("Group 2");

        List<VisibilityGroup> result = VisibilityGroupMapper.INSTANCE.fromDTOList(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("vg-1", result.get(0).getId());
        assertEquals("Group 1", result.get(0).getName());
        assertEquals("vg-2", result.get(1).getId());
        assertEquals("Group 2", result.get(1).getName());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<VisibilityGroup> result = VisibilityGroupMapper.INSTANCE.fromDTOList(
                Collections.emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<VisibilityGroup> result = VisibilityGroupMapper.INSTANCE.fromDTOList(null);
        assertNull(result);
    }

    @Test
    void fromUserDTO_mapsFields() {
        TgvalidatordInternalVisibilityGroupUser dto = new TgvalidatordInternalVisibilityGroupUser();
        dto.setId("user-123");
        dto.setExternalUserId("ext-456");

        VisibilityGroupUser result = VisibilityGroupMapper.INSTANCE.fromUserDTO(dto);

        assertNotNull(result);
        assertEquals("user-123", result.getId());
        assertEquals("ext-456", result.getExternalUserId());
    }

    @Test
    void fromUserDTO_handlesNullDto() {
        VisibilityGroupUser result = VisibilityGroupMapper.INSTANCE.fromUserDTO(null);
        assertNull(result);
    }

    @Test
    void fromUserDTOList_mapsList() {
        TgvalidatordInternalVisibilityGroupUser dto1 = new TgvalidatordInternalVisibilityGroupUser();
        dto1.setId("user-1");

        TgvalidatordInternalVisibilityGroupUser dto2 = new TgvalidatordInternalVisibilityGroupUser();
        dto2.setId("user-2");

        List<VisibilityGroupUser> result = VisibilityGroupMapper.INSTANCE.fromUserDTOList(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("user-1", result.get(0).getId());
        assertEquals("user-2", result.get(1).getId());
    }
}
