package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Cursor pagination parameters for API requests.
 * <p>
 * This is the low-level form of a cursor-list page request: it can express every
 * {@link PageRequest}, including {@code PREVIOUS} and {@code LAST}. The contract form is a
 * page size plus the previous page's {@link CursorPage#getNextCursor()}, which the list
 * methods accept directly or through {@link Pagination#page(Integer, String)}.
 * <p>
 * An explicit page size must be between 1 and {@link Pagination#MAX_PAGE_SIZE}; a cursor
 * built with the no-argument constructor leaves it unset, and the SDK then sends
 * {@link Pagination#DEFAULT_PAGE_SIZE}.
 * <p>
 * Example:
 * <pre>{@code
 * // First page request
 * ApiRequestCursor cursor = Pagination.first(20);
 * BalanceResult result = service.getBalances(cursor);
 *
 * // Next page request (if available)
 * if (result.hasNext()) {
 *     cursor = result.nextCursor(20);
 *     result = service.getBalances(cursor);
 * }
 * }</pre>
 *
 * @see ApiResponseCursor
 * @see Pagination
 * @see PageRequest
 */
public class ApiRequestCursor {

    /**
     * Token identifying the current page position (from previous response).
     */
    private String currentPage;

    /**
     * The type of page navigation (FIRST, NEXT, PREVIOUS, etc.).
     */
    private PageRequest pageRequest;

    /**
     * The number of items to return per page; 0 when unset.
     */
    private long pageSize;

    /**
     * Default constructor for framework use (e.g., deserialization).
     */
    public ApiRequestCursor() {
        // Required for deserialization frameworks
    }

    /**
     * Creates a cursor for initial page requests.
     *
     * @param pageRequest the page request type (typically FIRST)
     * @param pageSize    the page size, between 1 and {@link Pagination#MAX_PAGE_SIZE}
     * @throws IllegalArgumentException if the page size is out of range
     */
    public ApiRequestCursor(PageRequest pageRequest, long pageSize) {
        this.pageRequest = pageRequest;
        this.pageSize = Pagination.checkPageSize(pageSize);
    }

    /**
     * Creates a cursor for subsequent page requests.
     *
     * @param currentPage the current page token
     * @param pageRequest the page request type (typically NEXT)
     * @param pageSize    the page size, between 1 and {@link Pagination#MAX_PAGE_SIZE}
     * @throws IllegalArgumentException if the page size is out of range
     */
    public ApiRequestCursor(String currentPage, PageRequest pageRequest, long pageSize) {
        this.currentPage = currentPage;
        this.pageRequest = pageRequest;
        this.pageSize = Pagination.checkPageSize(pageSize);
    }

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the token identifying the current page position.
     *
     * @return the current page token, or {@code null} for initial requests
     */
    public String getCurrentPage() {
        return currentPage;
    }

    /**
     * Sets the token identifying the current page position.
     *
     * @param currentPage the current page token
     */
    public void setCurrentPage(String currentPage) {
        this.currentPage = currentPage;
    }

    /**
     * Returns the type of page navigation.
     *
     * @return the page request type (FIRST, NEXT, PREVIOUS, etc.)
     */
    public PageRequest getPageRequest() {
        return pageRequest;
    }

    /**
     * Sets the type of page navigation.
     *
     * @param pageRequest the page request type
     */
    public void setPageRequest(PageRequest pageRequest) {
        this.pageRequest = pageRequest;
    }

    /**
     * Returns the number of items to return per page.
     *
     * @return the page size, or 0 when unset (the SDK then sends the default)
     */
    public long getPageSize() {
        return pageSize;
    }

    /**
     * Sets the number of items to return per page.
     *
     * @param pageSize the page size, between 1 and {@link Pagination#MAX_PAGE_SIZE}
     * @throws IllegalArgumentException if the page size is out of range
     */
    public void setPageSize(long pageSize) {
        this.pageSize = Pagination.checkPageSize(pageSize);
    }
}
