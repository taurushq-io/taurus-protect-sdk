package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.HealthCheckStatus;
import com.taurushq.sdk.protect.client.model.HealthComponent;
import com.taurushq.sdk.protect.client.model.HealthGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordHealth;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordHealthComponent;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordHealthGroup;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class HealthMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordHealth dto = new TgvalidatordHealth();
        dto.setTenantId("tenant-123");
        dto.setComponentName("database");
        dto.setComponentId("db-001");
        dto.setGroup("storage");
        dto.setHealthCheck("connection");
        dto.setStatus("healthy");

        HealthCheckStatus status = HealthMapper.INSTANCE.fromDTO(dto);

        assertEquals("tenant-123", status.getTenantId());
        assertEquals("database", status.getComponentName());
        assertEquals("db-001", status.getComponentId());
        assertEquals("storage", status.getGroup());
        assertEquals("connection", status.getHealthCheck());
        assertEquals("healthy", status.getStatus());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordHealth dto = new TgvalidatordHealth();
        dto.setStatus("unknown");

        HealthCheckStatus status = HealthMapper.INSTANCE.fromDTO(dto);

        assertNull(status.getTenantId());
        assertNull(status.getComponentName());
        assertNull(status.getComponentId());
        assertNull(status.getGroup());
        assertNull(status.getHealthCheck());
        assertEquals("unknown", status.getStatus());
    }

    @Test
    void fromDTO_handlesNullDto() {
        HealthCheckStatus status = HealthMapper.INSTANCE.fromDTO(null);
        assertNull(status);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordHealth dto1 = new TgvalidatordHealth();
        dto1.setComponentName("api");
        dto1.setStatus("healthy");

        TgvalidatordHealth dto2 = new TgvalidatordHealth();
        dto2.setComponentName("cache");
        dto2.setStatus("unhealthy");

        List<HealthCheckStatus> statuses = HealthMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(statuses);
        assertEquals(2, statuses.size());
        assertEquals("api", statuses.get(0).getComponentName());
        assertEquals("healthy", statuses.get(0).getStatus());
        assertEquals("cache", statuses.get(1).getComponentName());
        assertEquals("unhealthy", statuses.get(1).getStatus());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<HealthCheckStatus> statuses = HealthMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(statuses);
        assertTrue(statuses.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<HealthCheckStatus> statuses = HealthMapper.INSTANCE.fromDTOList(null);
        assertNull(statuses);
    }

    @Test
    void fromGroupDTO_mapsHealthChecks() {
        TgvalidatordHealth healthDto = new TgvalidatordHealth();
        healthDto.setComponentName("test-component");
        healthDto.setStatus("healthy");

        TgvalidatordHealthGroup dto = new TgvalidatordHealthGroup();
        dto.setHealthChecks(Collections.singletonList(healthDto));

        HealthGroup group = HealthMapper.INSTANCE.fromGroupDTO(dto);

        assertNotNull(group);
        assertNotNull(group.getHealthChecks());
        assertEquals(1, group.getHealthChecks().size());
        assertEquals("test-component", group.getHealthChecks().get(0).getComponentName());
    }

    @Test
    void fromGroupDTO_handlesNullDto() {
        HealthGroup group = HealthMapper.INSTANCE.fromGroupDTO(null);
        assertNull(group);
    }

    @Test
    void fromComponentDTO_mapsGroups() {
        TgvalidatordHealthGroup groupDto = new TgvalidatordHealthGroup();
        groupDto.setHealthChecks(Collections.emptyList());

        Map<String, TgvalidatordHealthGroup> groupsMap = new HashMap<>();
        groupsMap.put("storage", groupDto);

        TgvalidatordHealthComponent dto = new TgvalidatordHealthComponent();
        dto.setGroups(groupsMap);

        HealthComponent component = HealthMapper.INSTANCE.fromComponentDTO(dto);

        assertNotNull(component);
        assertNotNull(component.getGroups());
        assertTrue(component.getGroups().containsKey("storage"));
    }

    @Test
    void fromComponentDTO_handlesNullDto() {
        HealthComponent component = HealthMapper.INSTANCE.fromComponentDTO(null);
        assertNull(component);
    }
}
