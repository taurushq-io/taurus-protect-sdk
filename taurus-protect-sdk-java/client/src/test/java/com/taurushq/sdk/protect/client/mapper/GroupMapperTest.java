package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Group;
import com.taurushq.sdk.protect.client.model.GroupUser;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalGroupUser;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class GroupMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime creationDate = OffsetDateTime.now();
        OffsetDateTime updateDate = OffsetDateTime.now().plusHours(1);

        TgvalidatordInternalGroupUser userDto = new TgvalidatordInternalGroupUser();
        userDto.setId("user-123");
        userDto.setExternalUserId("ext-user-123");
        userDto.setEnforcedInRules(true);

        TgvalidatordInternalGroup dto = new TgvalidatordInternalGroup();
        dto.setId("group-123");
        dto.setTenantId("tenant-456");
        dto.setExternalGroupId("ext-group-789");
        dto.setName("Administrators");
        dto.setEmail("admin@example.com");
        dto.setDescription("System administrators group");
        dto.setEnforcedInRules(true);
        dto.setCreationDate(creationDate);
        dto.setUpdateDate(updateDate);
        dto.setUsers(Collections.singletonList(userDto));

        Group group = GroupMapper.INSTANCE.fromDTO(dto);

        assertEquals("group-123", group.getId());
        assertEquals("tenant-456", group.getTenantId());
        assertEquals("ext-group-789", group.getExternalGroupId());
        assertEquals("Administrators", group.getName());
        assertEquals("admin@example.com", group.getEmail());
        assertEquals("System administrators group", group.getDescription());
        assertEquals(true, group.getEnforcedInRules());
        assertEquals(creationDate, group.getCreationDate());
        assertEquals(updateDate, group.getUpdateDate());
        assertNotNull(group.getUsers());
        assertEquals(1, group.getUsers().size());
        assertEquals("user-123", group.getUsers().get(0).getId());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordInternalGroup dto = new TgvalidatordInternalGroup();
        dto.setId("group-minimal");

        Group group = GroupMapper.INSTANCE.fromDTO(dto);

        assertEquals("group-minimal", group.getId());
        assertNull(group.getTenantId());
        assertNull(group.getExternalGroupId());
        assertNull(group.getName());
        assertNull(group.getEmail());
        assertNull(group.getDescription());
        assertNull(group.getEnforcedInRules());
        assertNull(group.getCreationDate());
        assertNull(group.getUpdateDate());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Group group = GroupMapper.INSTANCE.fromDTO(null);
        assertNull(group);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordInternalGroup dto1 = new TgvalidatordInternalGroup();
        dto1.setId("group-1");
        dto1.setName("Group A");

        TgvalidatordInternalGroup dto2 = new TgvalidatordInternalGroup();
        dto2.setId("group-2");
        dto2.setName("Group B");

        List<Group> groups = GroupMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(groups);
        assertEquals(2, groups.size());
        assertEquals("group-1", groups.get(0).getId());
        assertEquals("Group A", groups.get(0).getName());
        assertEquals("group-2", groups.get(1).getId());
        assertEquals("Group B", groups.get(1).getName());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<Group> groups = GroupMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(groups);
        assertTrue(groups.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<Group> groups = GroupMapper.INSTANCE.fromDTOList(null);
        assertNull(groups);
    }

    @Test
    void fromUserDTO_mapsAllFields() {
        TgvalidatordInternalGroupUser dto = new TgvalidatordInternalGroupUser();
        dto.setId("user-abc");
        dto.setExternalUserId("ext-user-abc");
        dto.setEnforcedInRules(false);

        GroupUser user = GroupMapper.INSTANCE.fromUserDTO(dto);

        assertEquals("user-abc", user.getId());
        assertEquals("ext-user-abc", user.getExternalUserId());
        assertEquals(false, user.getEnforcedInRules());
    }

    @Test
    void fromUserDTO_handlesNullFields() {
        TgvalidatordInternalGroupUser dto = new TgvalidatordInternalGroupUser();
        dto.setId("user-minimal");

        GroupUser user = GroupMapper.INSTANCE.fromUserDTO(dto);

        assertEquals("user-minimal", user.getId());
        assertNull(user.getExternalUserId());
        assertNull(user.getEnforcedInRules());
    }

    @Test
    void fromUserDTO_handlesNullDto() {
        GroupUser user = GroupMapper.INSTANCE.fromUserDTO(null);
        assertNull(user);
    }
}
