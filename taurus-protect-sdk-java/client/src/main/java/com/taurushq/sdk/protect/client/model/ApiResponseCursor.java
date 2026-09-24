package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Cursor pagination information from API responses, as the server sent it.
 * <p>
 * This is the raw cursor: the current page token and whether a previous or next page
 * exists. The contract value built from it is {@link CursorPage}, exposed by every cursor
 * result as {@code getPage()}; use that to walk a list:
 * <pre>{@code
 * ChangeResult page = client.getChangeService().getChanges(null, null, 20, null);
 * while (page.getPage().hasMore()) {
 *     page = client.getChangeService().getChanges(null, null, 20, page.getPage().getNextCursor());
 * }
 * }</pre>
 * The methods below build a low-level {@link ApiRequestCursor}, which can also go
 * backwards ({@link #previousPage(int)}).
 */
public class ApiResponseCursor {

    /**
     * Token identifying the current page position.
     */
    private String currentPage;

    /**
     * Whether a previous page of results is available.
     */
    private Boolean hasPrevious;

    /**
     * Whether a next page of results is available.
     */
    private Boolean hasNext;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets current page identifier.
     *
     * @return the current page
     */
    public String getCurrentPage() {
        return currentPage;
    }

    /**
     * Sets current page identifier.
     *
     * @param currentPage the current page
     */
    public void setCurrentPage(String currentPage) {
        this.currentPage = currentPage;
    }

    /**
     * Gets whether a previous page is available.
     *
     * @return true if previous page available, false otherwise
     */
    public boolean hasPrevious() {
        return Boolean.TRUE.equals(hasPrevious);
    }

    /**
     * Sets whether a previous page is available.
     *
     * @param hasPrevious true if previous page available
     */
    public void setHasPrevious(Boolean hasPrevious) {
        this.hasPrevious = hasPrevious;
    }

    /**
     * Gets whether a next page is available.
     *
     * @return true if next page available, false otherwise
     */
    public boolean hasNext() {
        return Boolean.TRUE.equals(hasNext);
    }

    /**
     * Sets whether a next page is available.
     *
     * @param hasNext true if next page available
     */
    public void setHasNext(Boolean hasNext) {
        this.hasNext = hasNext;
    }

    /**
     * Creates a cursor for the next page with the specified page size.
     * <p>
     * This is a convenience method equivalent to calling
     * {@code Pagination.next(this, pageSize)}.
     *
     * @param pageSize the number of items per page
     * @return a cursor for the next page
     * @throws IllegalStateException if there is no next page
     */
    public ApiRequestCursor nextPage(int pageSize) {
        return Pagination.next(this, pageSize);
    }

    /**
     * Creates a cursor for the next page with default page size.
     *
     * @return a cursor for the next page
     * @throws IllegalStateException if there is no next page
     */
    public ApiRequestCursor nextPage() {
        return Pagination.next(this);
    }

    /**
     * Creates a cursor for the previous page with the specified page size.
     * <p>
     * This is a convenience method equivalent to calling
     * {@code Pagination.previous(this, pageSize)}.
     *
     * @param pageSize the number of items per page
     * @return a cursor for the previous page
     * @throws IllegalStateException if there is no previous page
     */
    public ApiRequestCursor previousPage(int pageSize) {
        return Pagination.previous(this, pageSize);
    }

    /**
     * Creates a cursor for the previous page with default page size.
     *
     * @return a cursor for the previous page
     * @throws IllegalStateException if there is no previous page
     */
    public ApiRequestCursor previousPage() {
        return Pagination.previous(this);
    }
}
