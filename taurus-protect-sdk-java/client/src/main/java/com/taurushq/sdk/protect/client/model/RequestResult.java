package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Result of a request query with cursor-based pagination.
 * <p>
 * Contains a page of transaction requests and cursor information for fetching
 * additional pages. Use {@link #hasNext()} to check for more pages and
 * {@link #nextCursor(long)} to create a cursor for the next page.
 *
 * @see Request
 * @see ApiResponseCursor
 */
public class RequestResult {

    /**
     * The list of requests in this page of results.
     */
    private List<Request> requests;

    /**
     * Pagination cursor containing page information.
     */
    private ApiResponseCursor cursor;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the list of requests.
     *
     * @return the requests
     */
    public List<Request> getRequests() {
        return requests;
    }

    /**
     * Sets the list of requests.
     *
     * @param requests the requests
     */
    public void setRequests(List<Request> requests) {
        this.requests = requests;
    }

    /**
     * Gets the response cursor for pagination.
     *
     * @return the cursor
     */
    public ApiResponseCursor getCursor() {
        return cursor;
    }

    /**
     * Sets the response cursor for pagination.
     *
     * @param cursor the cursor
     */
    public void setCursor(ApiResponseCursor cursor) {
        this.cursor = cursor;
    }

    /**
     * Creates a cursor for the next page, or null if no more pages.
     *
     * @param pageSize the page size
     * @return the next cursor, or null if no more pages
     */
    public ApiRequestCursor nextCursor(long pageSize) {
        if (hasNext()) {
            return new ApiRequestCursor(cursor.getCurrentPage(), PageRequest.NEXT, pageSize);
        }
        return null;
    }

    /**
     * Checks if there is a next page of results.
     *
     * @return true if there is a next page
     */
    public boolean hasNext() {
        return cursor != null && cursor.hasNext();
    }
}
