package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Action;
import com.taurushq.sdk.protect.client.model.ActionAttribute;
import com.taurushq.sdk.protect.client.model.ActionEnvelope;
import com.taurushq.sdk.protect.client.model.ActionTask;
import com.taurushq.sdk.protect.client.model.ActionTrail;
import com.taurushq.sdk.protect.client.model.ActionTrigger;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAction;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionAttribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionEnvelopeTrail;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ActionMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime creationDate = OffsetDateTime.now();
        OffsetDateTime updateDate = OffsetDateTime.now().plusHours(1);
        OffsetDateTime lastCheckedDate = OffsetDateTime.now().plusHours(2);

        TgvalidatordAction action = new TgvalidatordAction();

        TgvalidatordActionEnvelope dto = new TgvalidatordActionEnvelope();
        dto.setId("action-123");
        dto.setTenantId("tenant-456");
        dto.setLabel("Test Action");
        dto.setAction(action);
        dto.setStatus("ACTIVE");
        dto.setCreationDate(creationDate);
        dto.setUpdateDate(updateDate);
        dto.setLastcheckeddate(lastCheckedDate);
        dto.setAutoApprove(true);

        ActionEnvelope result = ActionMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("action-123", result.getId());
        assertEquals("tenant-456", result.getTenantId());
        assertEquals("Test Action", result.getLabel());
        assertNotNull(result.getAction());
        assertEquals("ACTIVE", result.getStatus());
        assertEquals(creationDate, result.getCreationDate());
        assertEquals(updateDate, result.getUpdateDate());
        assertEquals(lastCheckedDate, result.getLastCheckedDate());
        assertTrue(result.getAutoApprove());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordActionEnvelope dto = new TgvalidatordActionEnvelope();
        dto.setId("action-minimal");

        ActionEnvelope result = ActionMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("action-minimal", result.getId());
        assertNull(result.getTenantId());
        assertNull(result.getLabel());
        assertNull(result.getAction());
        assertNull(result.getStatus());
        assertNull(result.getCreationDate());
        assertNull(result.getUpdateDate());
        assertNull(result.getLastCheckedDate());
        assertNull(result.getAutoApprove());
    }

    @Test
    void fromDTO_handlesNullDto() {
        ActionEnvelope result = ActionMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordActionEnvelope dto1 = new TgvalidatordActionEnvelope();
        dto1.setId("action-1");
        dto1.setLabel("Action 1");

        TgvalidatordActionEnvelope dto2 = new TgvalidatordActionEnvelope();
        dto2.setId("action-2");
        dto2.setLabel("Action 2");

        List<ActionEnvelope> result = ActionMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("action-1", result.get(0).getId());
        assertEquals("Action 1", result.get(0).getLabel());
        assertEquals("action-2", result.get(1).getId());
        assertEquals("Action 2", result.get(1).getLabel());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<ActionEnvelope> result = ActionMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<ActionEnvelope> result = ActionMapper.INSTANCE.fromDTOList(null);
        assertNull(result);
    }

    @Test
    void fromActionDTO_mapsFields() {
        TgvalidatordAction dto = new TgvalidatordAction();

        Action result = ActionMapper.INSTANCE.fromActionDTO(dto);

        assertNotNull(result);
    }

    @Test
    void fromActionDTO_handlesNullDto() {
        Action result = ActionMapper.INSTANCE.fromActionDTO(null);
        assertNull(result);
    }

    @Test
    void fromTriggerDTO_mapsFields() {
        com.taurushq.sdk.protect.openapi.model.ActionTrigger dto =
                new com.taurushq.sdk.protect.openapi.model.ActionTrigger();
        dto.setKind("balance");

        ActionTrigger result = ActionMapper.INSTANCE.fromTriggerDTO(dto);

        assertNotNull(result);
        assertEquals("balance", result.getKind());
    }

    @Test
    void fromTaskDTO_mapsFields() {
        com.taurushq.sdk.protect.openapi.model.ActionTask dto =
                new com.taurushq.sdk.protect.openapi.model.ActionTask();
        dto.setKind("transfer");

        ActionTask result = ActionMapper.INSTANCE.fromTaskDTO(dto);

        assertNotNull(result);
        assertEquals("transfer", result.getKind());
    }

    @Test
    void fromAttributeDTO_mapsFields() {
        TgvalidatordActionAttribute dto = new TgvalidatordActionAttribute();
        dto.setId("attr-123");
        dto.setTenantId("tenant-456");
        dto.setKey("myKey");
        dto.setValue("myValue");
        dto.setContentType("text/plain");

        ActionAttribute result = ActionMapper.INSTANCE.fromAttributeDTO(dto);

        assertNotNull(result);
        assertEquals("attr-123", result.getId());
        assertEquals("tenant-456", result.getTenantId());
        assertEquals("myKey", result.getKey());
        assertEquals("myValue", result.getValue());
        assertEquals("text/plain", result.getContentType());
    }

    @Test
    void fromTrailDTO_mapsFields() {
        OffsetDateTime date = OffsetDateTime.now();

        TgvalidatordActionEnvelopeTrail dto = new TgvalidatordActionEnvelopeTrail();
        dto.setId("trail-123");
        dto.setAction("create");
        dto.setComment("Created action");
        dto.setDate(date);
        dto.setActionStatus("SUCCESS");

        ActionTrail result = ActionMapper.INSTANCE.fromTrailDTO(dto);

        assertNotNull(result);
        assertEquals("trail-123", result.getId());
        assertEquals("create", result.getAction());
        assertEquals("Created action", result.getComment());
        assertEquals(date, result.getDate());
        assertEquals("SUCCESS", result.getActionStatus());
    }

    @Test
    void fromAttributeDTOList_mapsList() {
        TgvalidatordActionAttribute dto1 = new TgvalidatordActionAttribute();
        dto1.setId("attr-1");
        dto1.setKey("key1");

        TgvalidatordActionAttribute dto2 = new TgvalidatordActionAttribute();
        dto2.setId("attr-2");
        dto2.setKey("key2");

        List<ActionAttribute> result = ActionMapper.INSTANCE.fromAttributeDTOList(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("attr-1", result.get(0).getId());
        assertEquals("key1", result.get(0).getKey());
        assertEquals("attr-2", result.get(1).getId());
        assertEquals("key2", result.get(1).getKey());
    }

    @Test
    void fromTrailDTOList_mapsList() {
        TgvalidatordActionEnvelopeTrail dto1 = new TgvalidatordActionEnvelopeTrail();
        dto1.setId("trail-1");
        dto1.setAction("create");

        TgvalidatordActionEnvelopeTrail dto2 = new TgvalidatordActionEnvelopeTrail();
        dto2.setId("trail-2");
        dto2.setAction("update");

        List<ActionTrail> result = ActionMapper.INSTANCE.fromTrailDTOList(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("trail-1", result.get(0).getId());
        assertEquals("create", result.get(0).getAction());
        assertEquals("trail-2", result.get(1).getId());
        assertEquals("update", result.get(1).getAction());
    }
}
