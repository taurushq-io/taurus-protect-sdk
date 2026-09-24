package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of webhooks.
 * <p>
 * This class contains a list of webhooks and pagination information
 * for iterating through large result sets.
 *
 * @see Webhook
 */
public class WebhookResult extends CursorPagedResult {

    private List<Webhook> webhooks;

    /**
     * Gets the webhooks of this page.
     *
     * @return the webhooks
     */
    public List<Webhook> getWebhooks() {
        return webhooks;
    }

    /**
     * Sets the webhooks of this page.
     *
     * @param webhooks the webhooks
     */
    public void setWebhooks(final List<Webhook> webhooks) {
        this.webhooks = webhooks;
    }
}
