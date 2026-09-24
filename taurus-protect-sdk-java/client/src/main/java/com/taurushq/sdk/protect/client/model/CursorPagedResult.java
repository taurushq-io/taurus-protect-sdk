package com.taurushq.sdk.protect.client.model;

/**
 * Base of every cursor-list result: the {@link CursorPage} a service builds, plus the raw
 * server cursor.
 * <p>
 * Walk a list with {@link #getPage()}: pass {@link CursorPage#getNextCursor()} back while
 * {@link CursorPage#hasMore()} is true. {@link #getCursor()} is the cursor as the server sent
 * it, for low-level navigation such as {@link ApiResponseCursor#previousPage(int)}.
 */
public abstract class CursorPagedResult {

    /**
     * The raw server cursor.
     */
    private ApiResponseCursor cursor;

    /**
     * The contract pagination value.
     */
    private CursorPage page;

    /**
     * Gets the raw server cursor.
     *
     * @return the cursor as received, may be null when the server sent none
     */
    public ApiResponseCursor getCursor() {
        return cursor;
    }

    /**
     * Sets the raw server cursor.
     *
     * @param cursor the cursor
     */
    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }

    /**
     * Gets the page: page size, next cursor, whether more pages exist, and the total where
     * the endpoint reports one.
     *
     * @return the page, never null on a result returned by a service
     */
    public CursorPage getPage() {
        return page;
    }

    /**
     * Sets the page.
     *
     * @param page the page
     */
    public void setPage(final CursorPage page) {
        this.page = page;
    }

    /**
     * Checks whether another page exists.
     *
     * @return true if a next page is available
     */
    public boolean hasNext() {
        if (page != null) {
            return page.hasMore();
        }
        return cursor != null && cursor.hasNext();
    }

    /**
     * Creates a low-level request cursor for the next page.
     *
     * @param pageSize the page size for the next request
     * @return the cursor for the next page, or null when there is none
     * @throws IllegalArgumentException if the page size is out of range
     */
    public ApiRequestCursor nextCursor(final long pageSize) {
        if (!hasNext()) {
            return null;
        }
        String token = page != null ? page.getNextCursor() : cursor.getCurrentPage();
        return new ApiRequestCursor(token, PageRequest.NEXT, pageSize);
    }
}
