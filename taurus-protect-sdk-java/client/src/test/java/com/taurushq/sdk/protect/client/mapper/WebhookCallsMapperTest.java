package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.WebhookCall;
import com.taurushq.sdk.protect.client.model.WebhookCallResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWebhookCallsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWebhookCall;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class WebhookCallsMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime createdAt = OffsetDateTime.now();
        OffsetDateTime updatedAt = OffsetDateTime.now().plusHours(1);

        TgvalidatordWebhookCall dto = new TgvalidatordWebhookCall();
        dto.setId("call-123");
        dto.setEventId("event-456");
        dto.setWebhookId("webhook-789");
        dto.setPayload("{\"test\": \"data\"}");
        dto.setStatus("SUCCESS");
        dto.setStatusMessage("Delivered successfully");
        dto.setAttempts("3");
        dto.setCreatedAt(createdAt);
        dto.setUpdatedAt(updatedAt);

        WebhookCall result = WebhookCallsMapper.INSTANCE.fromDTO(dto);

        assertNotNull(result);
        assertEquals("call-123", result.getId());
        assertEquals("event-456", result.getEventId());
        assertEquals("webhook-789", result.getWebhookId());
        assertEquals("{\"test\": \"data\"}", result.getPayload());
        assertEquals("SUCCESS", result.getStatus());
        assertEquals("Delivered successfully", result.getStatusMessage());
        assertEquals("3", result.getAttempts());
        assertEquals(createdAt, result.getCreatedAt());
        assertEquals(updatedAt, result.getUpdatedAt());
    }

    @Test
    void fromDTO_handlesNullDto() {
        WebhookCall result = WebhookCallsMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordWebhookCall dto1 = new TgvalidatordWebhookCall();
        dto1.setId("call-1");
        dto1.setStatus("SUCCESS");

        TgvalidatordWebhookCall dto2 = new TgvalidatordWebhookCall();
        dto2.setId("call-2");
        dto2.setStatus("FAILED");

        List<WebhookCall> result = WebhookCallsMapper.INSTANCE.fromDTOList(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals("call-1", result.get(0).getId());
        assertEquals("SUCCESS", result.get(0).getStatus());
        assertEquals("call-2", result.get(1).getId());
        assertEquals("FAILED", result.get(1).getStatus());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<WebhookCall> result = WebhookCallsMapper.INSTANCE.fromDTOList(
                Collections.emptyList());
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<WebhookCall> result = WebhookCallsMapper.INSTANCE.fromDTOList(null);
        assertNull(result);
    }

    @Test
    void fromReply_mapsCallsAndCursor() {
        TgvalidatordWebhookCall call = new TgvalidatordWebhookCall();
        call.setId("call-123");
        call.setStatus("SUCCESS");

        TgvalidatordResponseCursor cursor = new TgvalidatordResponseCursor();
        cursor.setCurrentPage("page-1");

        TgvalidatordGetWebhookCallsReply reply = new TgvalidatordGetWebhookCallsReply();
        reply.setCalls(Arrays.asList(call));
        reply.setCursor(cursor);

        WebhookCallResult result = WebhookCallsMapper.INSTANCE.fromReply(reply);

        assertNotNull(result);
        assertNotNull(result.getCalls());
        assertEquals(1, result.getCalls().size());
        assertEquals("call-123", result.getCalls().get(0).getId());
        assertNotNull(result.getCursor());
        assertEquals("page-1", result.getCursor().getCurrentPage());
    }

    @Test
    void fromReply_handlesNullReply() {
        WebhookCallResult result = WebhookCallsMapper.INSTANCE.fromReply(null);
        assertNull(result);
    }
}
