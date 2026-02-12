package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of webhook calls.
 *
 * @see WebhookCallsService
 */
public class WebhookCallResult {

    private List<WebhookCall> calls;
    private ApiResponseCursor cursor;

    public List<WebhookCall> getCalls() {
        return calls;
    }

    public void setCalls(final List<WebhookCall> calls) {
        this.calls = calls;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
