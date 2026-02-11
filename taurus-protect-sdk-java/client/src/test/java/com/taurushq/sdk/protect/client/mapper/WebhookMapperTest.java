package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Webhook;
import com.taurushq.sdk.protect.client.model.WebhookResult;
import com.taurushq.sdk.protect.client.model.WebhookStatus;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWebhooksReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWebhook;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class WebhookMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime now = OffsetDateTime.now();

        TgvalidatordWebhook dto = new TgvalidatordWebhook();
        dto.setId("webhook-123");
        dto.setType("TRANSACTION");
        dto.setUrl("https://example.com/webhook");
        dto.setStatus("ENABLED");
        dto.setCreatedAt(now);
        dto.setUpdatedAt(now);

        Webhook webhook = WebhookMapper.INSTANCE.fromDTO(dto);

        assertEquals("webhook-123", webhook.getId());
        assertEquals("TRANSACTION", webhook.getType());
        assertEquals("https://example.com/webhook", webhook.getUrl());
        assertEquals(WebhookStatus.ENABLED, webhook.getStatus());
        assertEquals(now, webhook.getCreatedAt());
        assertEquals(now, webhook.getUpdatedAt());
    }

    @Test
    void fromDTO_mapsDisabledStatus() {
        TgvalidatordWebhook dto = new TgvalidatordWebhook();
        dto.setId("webhook-456");
        dto.setStatus("DISABLED");

        Webhook webhook = WebhookMapper.INSTANCE.fromDTO(dto);

        assertEquals(WebhookStatus.DISABLED, webhook.getStatus());
    }

    @Test
    void fromDTO_mapsTimeoutStatus() {
        TgvalidatordWebhook dto = new TgvalidatordWebhook();
        dto.setId("webhook-789");
        dto.setStatus("TIMEOUT");

        Webhook webhook = WebhookMapper.INSTANCE.fromDTO(dto);

        assertEquals(WebhookStatus.TIMEOUT, webhook.getStatus());
    }

    @Test
    void fromDTO_handlesNullStatus() {
        TgvalidatordWebhook dto = new TgvalidatordWebhook();
        dto.setId("webhook-000");
        dto.setStatus(null);

        Webhook webhook = WebhookMapper.INSTANCE.fromDTO(dto);

        assertNull(webhook.getStatus());
    }

    @Test
    void fromDTO_handlesUnknownStatus() {
        TgvalidatordWebhook dto = new TgvalidatordWebhook();
        dto.setId("webhook-unknown");
        dto.setStatus("UNKNOWN_STATUS");

        Webhook webhook = WebhookMapper.INSTANCE.fromDTO(dto);

        assertNull(webhook.getStatus());
    }

    @Test
    void fromReply_mapsWebhooksAndCursor() {
        TgvalidatordWebhook webhook1 = new TgvalidatordWebhook();
        webhook1.setId("wh-1");
        webhook1.setType("TRANSACTION");
        webhook1.setStatus("ENABLED");

        TgvalidatordWebhook webhook2 = new TgvalidatordWebhook();
        webhook2.setId("wh-2");
        webhook2.setType("REQUEST");
        webhook2.setStatus("DISABLED");

        TgvalidatordResponseCursor cursor = new TgvalidatordResponseCursor();
        cursor.setCurrentPage("page-token-abc");
        cursor.setHasNext(true);

        TgvalidatordGetWebhooksReply reply = new TgvalidatordGetWebhooksReply();
        reply.setWebhooks(Arrays.asList(webhook1, webhook2));
        reply.setCursor(cursor);

        WebhookResult result = WebhookMapper.INSTANCE.fromReply(reply);

        assertNotNull(result);
        assertNotNull(result.getWebhooks());
        assertEquals(2, result.getWebhooks().size());

        assertEquals("wh-1", result.getWebhooks().get(0).getId());
        assertEquals("TRANSACTION", result.getWebhooks().get(0).getType());
        assertEquals(WebhookStatus.ENABLED, result.getWebhooks().get(0).getStatus());

        assertEquals("wh-2", result.getWebhooks().get(1).getId());
        assertEquals("REQUEST", result.getWebhooks().get(1).getType());
        assertEquals(WebhookStatus.DISABLED, result.getWebhooks().get(1).getStatus());

        assertNotNull(result.getCursor());
        assertEquals("page-token-abc", result.getCursor().getCurrentPage());
        assertTrue(result.hasNext());
    }

    @Test
    void fromReply_handlesEmptyWebhooks() {
        TgvalidatordGetWebhooksReply reply = new TgvalidatordGetWebhooksReply();
        reply.setWebhooks(Arrays.asList());
        reply.setCursor(null);

        WebhookResult result = WebhookMapper.INSTANCE.fromReply(reply);

        assertNotNull(result);
        assertNotNull(result.getWebhooks());
        assertTrue(result.getWebhooks().isEmpty());
        assertFalse(result.hasNext());
    }
}
