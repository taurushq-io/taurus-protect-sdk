package com.taurushq.sdk.protect.client.model;

/**
 * Represents the status of a webhook configuration.
 * <p>
 * Webhooks can be enabled or disabled. When disabled, no notifications
 * will be sent to the webhook URL. A webhook may also be in a timeout
 * state if it has failed repeatedly.
 */
public enum WebhookStatus {

    /**
     * The webhook is active and will receive notifications.
     */
    ENABLED("ENABLED"),

    /**
     * The webhook is disabled and will not receive notifications.
     */
    DISABLED("DISABLED"),

    /**
     * The webhook is temporarily disabled due to repeated failures.
     */
    TIMEOUT("TIMEOUT");

    private final String value;

    WebhookStatus(String value) {
        this.value = value;
    }

    /**
     * Gets the string value of the status.
     *
     * @return the status value
     */
    public String getValue() {
        return value;
    }

    /**
     * Converts a string value to a WebhookStatus enum.
     *
     * @param value the string value
     * @return the corresponding WebhookStatus, or null if not found
     */
    public static WebhookStatus fromValue(String value) {
        if (value == null) {
            return null;
        }
        for (WebhookStatus status : values()) {
            if (status.value.equalsIgnoreCase(value)) {
                return status;
            }
        }
        return null;
    }

    @Override
    public String toString() {
        return value;
    }
}
