package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
import com.taurushq.sdk.protect.client.model.Webhook;
import com.taurushq.sdk.protect.client.model.WebhookResult;
import com.taurushq.sdk.protect.client.model.WebhookStatus;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class WebhookServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private WebhookService webhookService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        webhookService = new WebhookService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new WebhookService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new WebhookService(apiClient, null));
    }

    @Test
    void createWebhook_throwsOnNullUrl() {
        assertThrows(IllegalArgumentException.class, () ->
                webhookService.createWebhook(null, "TRANSACTION", "secret"));
    }

    @Test
    void createWebhook_throwsOnEmptyUrl() {
        assertThrows(IllegalArgumentException.class, () ->
                webhookService.createWebhook("", "TRANSACTION", "secret"));
    }

    @Test
    void createWebhook_throwsOnNullType() {
        assertThrows(IllegalArgumentException.class, () ->
                webhookService.createWebhook("https://example.com", null, "secret"));
    }

    @Test
    void createWebhook_throwsOnEmptyType() {
        assertThrows(IllegalArgumentException.class, () ->
                webhookService.createWebhook("https://example.com", "", "secret"));
    }

    @Test
    void deleteWebhook_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                webhookService.deleteWebhook(null));
    }

    @Test
    void deleteWebhook_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                webhookService.deleteWebhook(""));
    }

    @Test
    void updateWebhookStatus_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                webhookService.updateWebhookStatus(null, WebhookStatus.ENABLED));
    }

    @Test
    void updateWebhookStatus_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                webhookService.updateWebhookStatus("", WebhookStatus.ENABLED));
    }

    @Test
    void updateWebhookStatus_throwsOnNullStatus() {
        assertThrows(NullPointerException.class, () ->
                webhookService.updateWebhookStatus("webhook-123", null));
    }

    @Test
    void webhookStatusEnum_fromValue() {
        assertEquals(WebhookStatus.ENABLED, WebhookStatus.fromValue("ENABLED"));
        assertEquals(WebhookStatus.DISABLED, WebhookStatus.fromValue("DISABLED"));
        assertEquals(WebhookStatus.TIMEOUT, WebhookStatus.fromValue("TIMEOUT"));
        assertEquals(WebhookStatus.ENABLED, WebhookStatus.fromValue("enabled"));
    }

    @Test
    void webhookStatusEnum_fromValueReturnsNullForUnknown() {
        assertEquals(null, WebhookStatus.fromValue("UNKNOWN"));
        assertEquals(null, WebhookStatus.fromValue(null));
    }

    @Test
    void webhookStatusEnum_getValue() {
        assertEquals("ENABLED", WebhookStatus.ENABLED.getValue());
        assertEquals("DISABLED", WebhookStatus.DISABLED.getValue());
        assertEquals("TIMEOUT", WebhookStatus.TIMEOUT.getValue());
    }

    @Test
    void apiRequestCursor_construction() {
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 50);
        assertEquals(PageRequest.FIRST, cursor.getPageRequest());
        assertEquals(50, cursor.getPageSize());
    }

    @Test
    void apiRequestCursor_constructionWithCurrentPage() {
        ApiRequestCursor cursor = new ApiRequestCursor("page-token", PageRequest.NEXT, 25);
        assertEquals("page-token", cursor.getCurrentPage());
        assertEquals(PageRequest.NEXT, cursor.getPageRequest());
        assertEquals(25, cursor.getPageSize());
    }
}
