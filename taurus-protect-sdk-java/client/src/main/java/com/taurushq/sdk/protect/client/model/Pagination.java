package com.taurushq.sdk.protect.client.model;

import java.util.regex.Pattern;

/**
 * Page-size constants, the one page-size resolver, and factories for request cursors.
 * <p>
 * Every list method follows the cross-SDK pagination contract: a page size is always
 * sent, {@code null} or {@code 0} means {@link #DEFAULT_PAGE_SIZE}, and a negative size
 * or one above {@link #MAX_PAGE_SIZE} is rejected with an {@link IllegalArgumentException}
 * before any request is made.
 * <p>
 * Offset lists return an {@link OffsetPagination}; cursor lists return a {@link CursorPage}
 * whose {@link CursorPage#getNextCursor()} is passed back until {@link CursorPage#hasMore()}
 * is false:
 * <pre>{@code
 * // Offset list
 * WalletResult wallets = client.getWalletService().getWallets(20, 0);
 * while (wallets.getPagination().hasMore()) {
 *     wallets = client.getWalletService().getWallets(20, wallets.getPagination().getNextOffset());
 * }
 *
 * // Cursor list
 * ChangeResult changes = client.getChangeService().getChanges(null, null, 20, null);
 * while (changes.getPage().hasMore()) {
 *     changes = client.getChangeService().getChanges(null, null, 20, changes.getPage().getNextCursor());
 * }
 * }</pre>
 */
public final class Pagination {

    /**
     * Page size sent when a caller leaves it unset or 0.
     */
    public static final int DEFAULT_PAGE_SIZE = 20;

    /**
     * Largest page size any list method accepts.
     */
    public static final int MAX_PAGE_SIZE = 100;

    /**
     * Largest price-history limit (one point a day for a year). That endpoint cannot page.
     */
    public static final int MAX_PRICE_HISTORY_LIMIT = 365;

    /**
     * Largest count the SDKs accept from the server: 2^53 - 1, the largest integer every
     * SDK (TypeScript included) represents exactly.
     */
    static final long MAX_COUNT = 9007199254740991L;

    /**
     * A canonical decimal count, as the server's uint64 fields arrive on the wire.
     */
    private static final Pattern COUNT = Pattern.compile("0|[1-9][0-9]{0,15}");

    /**
     * Longest server value echoed into an error message.
     */
    private static final int MAX_ECHO = 32;

    private Pagination() {
        // Static utility class
    }

    /**
     * Resolves a page size for a list method: {@code null} or 0 gives
     * {@link #DEFAULT_PAGE_SIZE}, a negative size or one above {@link #MAX_PAGE_SIZE} is
     * rejected.
     *
     * @param option the option name used in the error message
     * @param value  the caller's page size, may be null
     * @return the page size to send
     * @throws IllegalArgumentException if the size is negative or above the maximum
     */
    public static int resolvePageSize(final String option, final Integer value) {
        return resolveSize(option, value, MAX_PAGE_SIZE);
    }

    /**
     * The one page-size resolver: {@code null} or 0 gives {@link #DEFAULT_PAGE_SIZE}; a
     * negative value, or one above {@code maximum}, is rejected by name.
     *
     * @param option  the option name used in the error message
     * @param value   the caller's value, may be null
     * @param maximum the largest accepted value, or null when the endpoint has none
     * @return the value to send
     * @throws IllegalArgumentException if the value is negative or above the maximum
     */
    public static int resolveSize(final String option, final Integer value, final Integer maximum) {
        if (value == null || value == 0) {
            return DEFAULT_PAGE_SIZE;
        }
        if (maximum == null) {
            if (value < 0) {
                throw new IllegalArgumentException(option + " must not be negative, got " + value);
            }
        } else if (value < 0 || value > maximum) {
            throw new IllegalArgumentException(
                    option + " must be between 1 and " + maximum + ", got " + value);
        }
        return value;
    }

    /**
     * Resolves an offset for an offset list: it must not be negative.
     *
     * @param option the option name used in the error message
     * @param value  the caller's offset, 0 for the first page
     * @return the offset to send
     * @throws IllegalArgumentException if the offset is negative
     */
    public static long resolveOffset(final String option, final long value) {
        if (value < 0) {
            throw new IllegalArgumentException(option + " must not be negative, got " + value);
        }
        return value;
    }

    /**
     * Creates the request cursor for a page of a cursor list: the page size, plus the
     * previous page's {@link CursorPage#getNextCursor()} when there is one.
     * <p>
     * Without a cursor only the page size is sent (the first page); with one the SDK sends
     * {@code currentPage=<cursor>} and {@code pageRequest=NEXT}.
     *
     * @param pageSize the page size, null or 0 for {@link #DEFAULT_PAGE_SIZE}
     * @param cursor   a previous {@link CursorPage#getNextCursor()}, or null/empty for the
     *                 first page
     * @return the request cursor
     * @throws IllegalArgumentException if the page size is negative or above the maximum
     */
    public static ApiRequestCursor page(final Integer pageSize, final String cursor) {
        int size = resolvePageSize("pageSize", pageSize);
        if (cursor == null || cursor.isEmpty()) {
            return new ApiRequestCursor((PageRequest) null, size);
        }
        return new ApiRequestCursor(cursor, PageRequest.NEXT, size);
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
     * @throws IllegalArgumentException if pageSize is not between 1 and {@link #MAX_PAGE_SIZE}
     */
    public static ApiRequestCursor of(String pageToken, PageRequest pageRequest, int pageSize) {
        return new ApiRequestCursor(pageToken, pageRequest, pageSize);
    }

    /**
     * Checks an explicit page size, as carried by an {@link ApiRequestCursor}.
     *
     * @param pageSize the page size
     * @return the same page size
     * @throws IllegalArgumentException if it is not between 1 and {@link #MAX_PAGE_SIZE}
     */
    static long checkPageSize(final long pageSize) {
        if (pageSize < 1 || pageSize > MAX_PAGE_SIZE) {
            throw new IllegalArgumentException(
                    "pageSize must be between 1 and " + MAX_PAGE_SIZE + ", got " + pageSize);
        }
        return pageSize;
    }

    /**
     * Parses a count from a reply. Absent means 0, because the server omits zero values;
     * anything but a canonical decimal in [0, 2^53 - 1] is a malformed reply.
     *
     * @param field the reply field, for the error message
     * @param value the value as received, may be null
     * @return the count
     * @throws IntegrityException if the value is not a canonical count in range
     */
    static long parseCount(final String field, final String value) {
        if (value == null) {
            return 0L;
        }
        if (COUNT.matcher(value).matches()) {
            long parsed = Long.parseLong(value);
            if (parsed <= MAX_COUNT) {
                return parsed;
            }
        }
        String shown = value.length() > MAX_ECHO ? value.substring(0, MAX_ECHO) + "..." : value;
        throw new IntegrityException("malformed " + field + " in the server reply: \"" + shown
                + "\" is not a count between 0 and " + MAX_COUNT);
    }
}
