package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents a webhook call record in the Taurus Protect system.
 * <p>
 * A webhook call represents a single invocation of a configured webhook,
 * including the payload sent and the delivery status.
 *
 * @see WebhookCallsService
 */
public class WebhookCall {

    private String id;
    private String eventId;
    private String webhookId;
    private String payload;
    private String status;
    private String statusMessage;
    private String attempts;
    private OffsetDateTime updatedAt;
    private OffsetDateTime createdAt;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getEventId() {
        return eventId;
    }

    public void setEventId(final String eventId) {
        this.eventId = eventId;
    }

    public String getWebhookId() {
        return webhookId;
    }

    public void setWebhookId(final String webhookId) {
        this.webhookId = webhookId;
    }

    public String getPayload() {
        return payload;
    }

    public void setPayload(final String payload) {
        this.payload = payload;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
    }

    public String getStatusMessage() {
        return statusMessage;
    }

    public void setStatusMessage(final String statusMessage) {
        this.statusMessage = statusMessage;
    }

    public String getAttempts() {
        return attempts;
    }

    public void setAttempts(final String attempts) {
        this.attempts = attempts;
    }

    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(final OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }

    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(final OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }
}
