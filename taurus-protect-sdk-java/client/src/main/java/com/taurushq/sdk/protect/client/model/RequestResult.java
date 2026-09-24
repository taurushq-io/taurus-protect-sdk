package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Result of a request query with cursor-based pagination.
 * <p>
 * Contains a page of transaction requests and cursor information for fetching
 * additional pages: pass {@code getPage().getNextCursor()} back while
 * {@code getPage().hasMore()} is true.
 *
 * @see Request
 */
public class RequestResult extends CursorPagedResult {

    private List<Request> requests;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the requests of this page.
     *
     * @return the requests
     */
    public List<Request> getRequests() {
        return requests;
    }

    /**
     * Sets the requests of this page.
     *
     * @param requests the requests
     */
    public void setRequests(final List<Request> requests) {
        this.requests = requests;
    }
}
