package com.taurushq.sdk.protect.client.model;

import java.util.Objects;

/**
 * Pagination of one page of a cursor list, never null on a successful call.
 * <p>
 * Continue by passing {@link #getNextCursor()} back as the {@code cursor} argument while
 * {@link #hasMore()} is true. A cursor is an opaque string (base64 text, possibly containing
 * {@code + / =}); the SDK URL-encodes it on the wire.
 */
public final class CursorPage {

    private final int pageSize;
    private final String nextCursor;
    private final boolean hasMore;
    private final Long totalItems;

    /**
     * Creates a cursor page. Services build it with {@link #fromCursor} or
     * {@link #fromToken}.
     *
     * @param pageSize   the page size sent
     * @param nextCursor the cursor of the next page, empty when there is none
     * @param hasMore    whether another page exists
     * @param totalItems the server total, or null when the endpoint reports none
     */
    public CursorPage(final int pageSize, final String nextCursor, final boolean hasMore,
                      final Long totalItems) {
        this.pageSize = pageSize;
        this.nextCursor = nextCursor == null ? "" : nextCursor;
        this.hasMore = hasMore;
        this.totalItems = totalItems;
    }

    /**
     * The cursor-family builder: {@code hasMore} is the reply cursor's {@code hasNext}, and
     * the next cursor is the reply's {@code currentPage} when there is more, else empty.
     * Never derive a cursor any other way: sending {@code NEXT} past the end is a 400 on the
     * v2 lists and silently returns page 1 on the keyset lists.
     *
     * @param pageSize    the page size sent
     * @param replyCursor the reply's cursor, null when absent
     * @param totalItems  the reply's count as received, null when absent
     * @param hasTotal    whether the endpoint reports a count at all
     * @return the cursor page
     * @throws IntegrityException if the reply says there is a next page but carries no
     *                            cursor, or its count is malformed
     */
    public static CursorPage fromCursor(final int pageSize, final ApiResponseCursor replyCursor,
                                        final String totalItems, final boolean hasTotal) {
        boolean hasMore = replyCursor != null && replyCursor.hasNext();
        String current = replyCursor == null || replyCursor.getCurrentPage() == null
                ? "" : replyCursor.getCurrentPage();
        if (hasMore && current.isEmpty()) {
            throw new IntegrityException(
                    "malformed cursor in the server reply: hasNext is set but currentPage is empty");
        }
        Long total = hasTotal ? Pagination.parseCount("totalItems", totalItems) : null;
        return new CursorPage(pageSize, hasMore ? current : "", hasMore, total);
    }

    /**
     * The token-family builder, for the operations whose reply carries only a next-page
     * token (wallet tokens, governance rules history): the token IS the next cursor.
     *
     * @param pageSize   the page size sent
     * @param replyToken the reply's token, null or empty on the last page
     * @param totalItems the reply's count as received, null when absent
     * @param hasTotal   whether the endpoint reports a count at all
     * @return the cursor page
     * @throws IntegrityException if the count is malformed
     */
    public static CursorPage fromToken(final int pageSize, final String replyToken,
                                       final String totalItems, final boolean hasTotal) {
        String token = replyToken == null ? "" : replyToken;
        Long total = hasTotal ? Pagination.parseCount("totalItems", totalItems) : null;
        return new CursorPage(pageSize, token, !token.isEmpty(), total);
    }

    /**
     * Returns the page size that was sent.
     *
     * @return the page size
     */
    public int getPageSize() {
        return pageSize;
    }

    /**
     * Returns the cursor to request the next page with.
     *
     * @return the next cursor, empty when {@link #hasMore()} is false
     */
    public String getNextCursor() {
        return nextCursor;
    }

    /**
     * Reports whether another page exists.
     *
     * @return true while the walk should continue
     */
    public boolean hasMore() {
        return hasMore;
    }

    /**
     * Returns the server total, where the endpoint reports one (balances, asset balances,
     * wallet tokens, rules history).
     *
     * @return the total, or null when the endpoint reports none
     */
    public Long getTotalItems() {
        return totalItems;
    }

    @Override
    public boolean equals(final Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof CursorPage)) {
            return false;
        }
        CursorPage that = (CursorPage) o;
        return pageSize == that.pageSize && hasMore == that.hasMore
                && nextCursor.equals(that.nextCursor) && Objects.equals(totalItems, that.totalItems);
    }

    @Override
    public int hashCode() {
        return Objects.hash(pageSize, nextCursor, hasMore, totalItems);
    }

    @Override
    public String toString() {
        return "CursorPage{pageSize=" + pageSize + ", nextCursor='" + nextCursor + "', hasMore="
                + hasMore + ", totalItems=" + totalItems + "}";
    }
}
