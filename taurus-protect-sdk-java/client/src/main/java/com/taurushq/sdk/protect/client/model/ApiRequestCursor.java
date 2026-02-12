package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Cursor pagination parameters for API requests.
 * <p>
 * Use this class to request specific pages of paginated results.
 * For initial requests, use {@link Pagination#first(int)}. For subsequent
 * pages, use the cursor from the previous response.
 * <p>
 * Example:
 * <pre>{@code
 * // First page request
 * ApiRequestCursor cursor = Pagination.first(50);
 * BalanceResult result = service.getBalances(cursor);
 *
 * // Next page request (if available)
 * if (result.hasNext()) {
 *     cursor = result.nextCursor(50);
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
     * The number of items to return per page.
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
     * @param pageSize    the page size
     */
    public ApiRequestCursor(PageRequest pageRequest, long pageSize) {
        this.pageRequest = pageRequest;
        this.pageSize = pageSize;
    }

    /**
     * Creates a cursor for subsequent page requests.
     *
     * @param currentPage the current page token
     * @param pageRequest the page request type (typically NEXT)
     * @param pageSize    the page size
     */
    public ApiRequestCursor(String currentPage, PageRequest pageRequest, long pageSize) {
        this.currentPage = currentPage;
        this.pageRequest = pageRequest;
        this.pageSize = pageSize;
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
     * @return the page size
     */
    public long getPageSize() {
        return pageSize;
    }

    /**
     * Sets the number of items to return per page.
     *
     * @param pageSize the page size
     */
    public void setPageSize(long pageSize) {
        this.pageSize = pageSize;
    }
}
