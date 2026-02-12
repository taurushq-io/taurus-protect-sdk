package com.taurushq.sdk.protect.client.model;

/**
 * Represents the result of creating a webhook.
 * <p>
 * Contains both the webhook configuration and the secret used for
 * verifying webhook signatures.
 *
 * @see WebhookService
 */
public class WebhookCreateResult {

    private Webhook webhook;
    private String secret;

    public Webhook getWebhook() {
        return webhook;
    }

    public void setWebhook(final Webhook webhook) {
        this.webhook = webhook;
    }

    public String getSecret() {
        return secret;
    }

    public void setSecret(final String secret) {
        this.secret = secret;
    }
}
