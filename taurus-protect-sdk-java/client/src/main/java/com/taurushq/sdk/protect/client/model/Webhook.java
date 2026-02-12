package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents a webhook configuration in the Taurus Protect system.
 * <p>
 * A webhook is an HTTP callback that is triggered when specific events occur,
 * allowing external systems to receive real-time notifications about transactions,
 * request status changes, and other events.
 * <p>
 * Webhooks are created with a URL endpoint and a type that determines which
 * events will trigger the callback. The webhook can be enabled or disabled
 * using the status field.
 *
 * @see WebhookStatus
 */
public class Webhook {

    private String id;
    private String type;
    private String url;
    private WebhookStatus status;
    private OffsetDateTime timeoutUntil;
    private OffsetDateTime updatedAt;
    private OffsetDateTime createdAt;

    /**
     * Gets the unique identifier of the webhook.
     *
     * @return the webhook ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the unique identifier of the webhook.
     *
     * @param id the webhook ID
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the type of events this webhook subscribes to.
     * <p>
     * Common types include: TRANSACTION, REQUEST, ADDRESS, WALLET.
     *
     * @return the webhook type
     */
    public String getType() {
        return type;
    }

    /**
     * Sets the type of events this webhook subscribes to.
     *
     * @param type the webhook type
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Gets the URL endpoint that will receive webhook notifications.
     * <p>
     * This URL must be HTTPS and publicly accessible.
     *
     * @return the webhook URL
     */
    public String getUrl() {
        return url;
    }

    /**
     * Sets the URL endpoint that will receive webhook notifications.
     *
     * @param url the webhook URL
     */
    public void setUrl(String url) {
        this.url = url;
    }

    /**
     * Gets the current status of the webhook.
     *
     * @return the webhook status
     */
    public WebhookStatus getStatus() {
        return status;
    }

    /**
     * Sets the current status of the webhook.
     *
     * @param status the webhook status
     */
    public void setStatus(WebhookStatus status) {
        this.status = status;
    }

    /**
     * Gets the timeout timestamp if the webhook is temporarily disabled.
     * <p>
     * When a webhook fails repeatedly, it may be automatically disabled
     * until this timestamp.
     *
     * @return the timeout until timestamp, or null if not in timeout
     */
    public OffsetDateTime getTimeoutUntil() {
        return timeoutUntil;
    }

    /**
     * Sets the timeout timestamp.
     *
     * @param timeoutUntil the timeout until timestamp
     */
    public void setTimeoutUntil(OffsetDateTime timeoutUntil) {
        this.timeoutUntil = timeoutUntil;
    }

    /**
     * Gets the timestamp when the webhook was last updated.
     *
     * @return the last update timestamp
     */
    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    /**
     * Sets the timestamp when the webhook was last updated.
     *
     * @param updatedAt the last update timestamp
     */
    public void setUpdatedAt(OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }

    /**
     * Gets the timestamp when the webhook was created.
     *
     * @return the creation timestamp
     */
    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    /**
     * Sets the timestamp when the webhook was created.
     *
     * @param createdAt the creation timestamp
     */
    public void setCreatedAt(OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }

    @Override
    public String toString() {
        return "Webhook{"
                + "id='" + id + '\''
                + ", type='" + type + '\''
                + ", url='" + url + '\''
                + ", status=" + status
                + ", createdAt=" + createdAt
                + '}';
    }
}
