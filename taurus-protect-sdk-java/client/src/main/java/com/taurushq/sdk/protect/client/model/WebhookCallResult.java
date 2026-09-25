package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of webhook calls.
 *
 * @see WebhookCallsService
 */
public class WebhookCallResult extends CursorPagedResult {

    private List<WebhookCall> calls;

    /**
     * Gets the webhook calls of this page.
     *
     * @return the calls
     */
    public List<WebhookCall> getCalls() {
        return calls;
    }

    /**
     * Sets the webhook calls of this page.
     *
     * @param calls the calls
     */
    public void setCalls(final List<WebhookCall> calls) {
        this.calls = calls;
    }
}
