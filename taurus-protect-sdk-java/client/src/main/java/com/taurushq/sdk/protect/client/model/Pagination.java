package com.taurushq.sdk.protect.client.model;

/**
 * Factory for creating pagination cursors with a fluent API.
 * <p>
 * This class provides static factory methods for creating pagination cursors
 * that are more intuitive than constructing {@link ApiRequestCursor} directly.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get first page
 * List<Wallet> wallets = client.getWalletService().getWallets(
 *     "ETH", Pagination.first(50));
 *
 * // Get next page using response cursor
 * if (cursor.hasNext()) {
 *     List<Wallet> nextPage = client.getWalletService().getWallets(
 *         "ETH", Pagination.next(cursor, 50));
 * }
 * }</pre>
 */
public final class Pagination {

    /**
     * Default page size used when not specified.
     */
    public static final int DEFAULT_PAGE_SIZE = 50;

    /**
     * Maximum allowed page size.
     */
    public static final int MAX_PAGE_SIZE = 1000;

    private Pagination() {
        // Static utility class
    }

    /**
     * Creates a cursor for the first page with default page size.
     *
     * @return a cursor for the first page
     */
    public static ApiRequestCursor first() {
        return first(DEFAULT_PAGE_SIZE);
    }

    /**
     * Creates a cursor for the first page with the specified page size.
     *
     * @param pageSize the number of items per page
     * @return a cursor for the first page
     * @throws IllegalArgumentException if pageSize is not between 1 and {@link #MAX_PAGE_SIZE}
     */
    public static ApiRequestCursor first(int pageSize) {
        validatePageSize(pageSize);
        return new ApiRequestCursor(PageRequest.FIRST, pageSize);
    }

    /**
     * Creates a cursor for the next page based on a response cursor.
     *
     * @param responseCursor the response cursor from the previous request
     * @param pageSize       the number of items per page
     * @return a cursor for the next page
     * @throws IllegalArgumentException if pageSize is invalid or responseCursor is null
     * @throws IllegalStateException    if there is no next page
     */
    public static ApiRequestCursor next(ApiResponseCursor responseCursor, int pageSize) {
        if (responseCursor == null) {
            throw new IllegalArgumentException("responseCursor cannot be null");
        }
        if (!responseCursor.hasNext()) {
            throw new IllegalStateException("No next page available");
        }
        validatePageSize(pageSize);
        return new ApiRequestCursor(responseCursor.getCurrentPage(), PageRequest.NEXT, pageSize);
    }

    /**
     * Creates a cursor for the next page with default page size.
     *
     * @param responseCursor the response cursor from the previous request
     * @return a cursor for the next page
     */
    public static ApiRequestCursor next(ApiResponseCursor responseCursor) {
        return next(responseCursor, DEFAULT_PAGE_SIZE);
    }

    /**
     * Creates a cursor for the previous page based on a response cursor.
     *
     * @param responseCursor the response cursor from the previous request
     * @param pageSize       the number of items per page
     * @return a cursor for the previous page
     * @throws IllegalArgumentException if pageSize is invalid or responseCursor is null
     * @throws IllegalStateException    if there is no previous page
     */
    public static ApiRequestCursor previous(ApiResponseCursor responseCursor, int pageSize) {
        if (responseCursor == null) {
            throw new IllegalArgumentException("responseCursor cannot be null");
        }
        if (!responseCursor.hasPrevious()) {
            throw new IllegalStateException("No previous page available");
        }
        validatePageSize(pageSize);
        return new ApiRequestCursor(responseCursor.getCurrentPage(), PageRequest.PREVIOUS, pageSize);
    }

    /**
     * Creates a cursor for the previous page with default page size.
     *
     * @param responseCursor the response cursor from the previous request
     * @return a cursor for the previous page
     */
    public static ApiRequestCursor previous(ApiResponseCursor responseCursor) {
        return previous(responseCursor, DEFAULT_PAGE_SIZE);
    }

    /**
     * Creates a cursor for the last page with the specified page size.
     *
     * @param pageSize the number of items per page
     * @return a cursor for the last page
     * @throws IllegalArgumentException if pageSize is not between 1 and {@link #MAX_PAGE_SIZE}
     */
    public static ApiRequestCursor last(int pageSize) {
        validatePageSize(pageSize);
        return new ApiRequestCursor(PageRequest.LAST, pageSize);
    }

    /**
     * Creates a cursor for the last page with default page size.
     *
     * @return a cursor for the last page
     */
    public static ApiRequestCursor last() {
        return last(DEFAULT_PAGE_SIZE);
    }

    /**
     * Creates a cursor with a specific page token for direct navigation.
     *
     * @param pageToken   the page token to navigate to
     * @param pageRequest the type of page request
     * @param pageSize    the number of items per page
     * @return a cursor for the specified page
     */
    public static ApiRequestCursor of(String pageToken, PageRequest pageRequest, int pageSize) {
        validatePageSize(pageSize);
        return new ApiRequestCursor(pageToken, pageRequest, pageSize);
    }

    private static void validatePageSize(int pageSize) {
        if (pageSize < 1 || pageSize > MAX_PAGE_SIZE) {
            throw new IllegalArgumentException(
                    "pageSize must be between 1 and " + MAX_PAGE_SIZE + ", got: " + pageSize);
        }
    }
}
