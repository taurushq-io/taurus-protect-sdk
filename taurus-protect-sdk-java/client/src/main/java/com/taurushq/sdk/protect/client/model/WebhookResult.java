package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of webhooks.
 * <p>
 * This class contains a list of webhooks and pagination information
 * for iterating through large result sets.
 *
 * @see Webhook
 * @see ApiResponseCursor
 */
public class WebhookResult {

    private List<Webhook> webhooks;
    private ApiResponseCursor cursor;

    /**
     * Gets the list of webhooks in this result page.
     *
     * @return the list of webhooks
     */
    public List<Webhook> getWebhooks() {
        return webhooks;
    }

    /**
     * Sets the list of webhooks in this result page.
     *
     * @param webhooks the list of webhooks
     */
    public void setWebhooks(List<Webhook> webhooks) {
        this.webhooks = webhooks;
    }

    /**
     * Gets the pagination cursor for navigating results.
     *
     * @return the response cursor
     */
    public ApiResponseCursor getCursor() {
        return cursor;
    }

    /**
     * Sets the pagination cursor.
     *
     * @param cursor the response cursor
     */
    public void setCursor(ApiResponseCursor cursor) {
        this.cursor = cursor;
    }

    /**
     * Checks if there are more results available.
     *
     * @return true if more results exist, false otherwise
     */
    public boolean hasNext() {
        return cursor != null && cursor.hasNext();
    }

    /**
     * Creates a cursor for fetching the next page of results.
     *
     * @param pageSize the page size for the next request
     * @return an ApiRequestCursor for the next page
     */
    public ApiRequestCursor nextCursor(int pageSize) {
        if (cursor == null) {
            return null;
        }
        return new ApiRequestCursor(cursor.getCurrentPage(), PageRequest.NEXT, pageSize);
    }
}
