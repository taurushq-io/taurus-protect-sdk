package com.taurushq.sdk.protect.client.model;

import java.util.Objects;

/**
 * Pagination of one page of an offset list, never null on a successful call.
 * <p>
 * Continue by passing {@link #getNextOffset()} back as the offset while
 * {@link #hasMore()} is true. An empty reply ({@code {}}) is a valid last page: zero
 * totals, {@code hasMore() == false} and {@code getNextOffset() == getOffset()}.
 */
public final class OffsetPagination {

    private final int limit;
    private final long offset;
    private final long totalItems;
    private final long nextOffset;
    private final boolean hasMore;

    /**
     * Creates a pagination value. Services build it with {@link #of}.
     *
     * @param limit      the page size sent
     * @param offset     the offset sent
     * @param totalItems the total, reduced by the rows the SDK withheld
     * @param nextOffset the offset of the next page
     * @param hasMore    whether another page exists
     */
    public OffsetPagination(final int limit, final long offset, final long totalItems,
                            final long nextOffset, final boolean hasMore) {
        this.limit = limit;
        this.offset = offset;
        this.totalItems = totalItems;
        this.nextOffset = nextOffset;
        this.hasMore = hasMore;
    }

    /**
     * The one builder for offset lists.
     * <ul>
     *   <li>{@code limit} and {@code offset} are the request's, never the reply's.</li>
     *   <li>{@code nextOffset} follows the endpoint's {@link OffsetRule}; {@code hasMore} is
     *       {@code offset < nextOffset < serverTotal}, so a page that makes no progress
     *       ends a walk (some totals are upper bounds).</li>
     *   <li>{@code totalItems} is the server total reduced by the rows the SDK withheld,
     *       clamped at 0; {@code nextOffset} and {@code hasMore} use the server's
     *       numbers.</li>
     * </ul>
     *
     * @param rule         the endpoint's rule
     * @param limit        the page size sent
     * @param offset       the offset sent
     * @param serverRows   the rows the server returned, before any SDK exclusion
     * @param excludedRows the rows the SDK withheld (failed verification)
     * @param totalItems   the reply's {@code totalItems} as received, null when absent
     * @param replyOffset  the reply's {@code offset} as received, null when absent
     * @return the pagination value
     * @throws IntegrityException if a reply count is not a canonical count in range
     */
    public static OffsetPagination of(final OffsetRule rule, final int limit, final long offset,
                                      final int serverRows, final int excludedRows,
                                      final String totalItems, final String replyOffset) {
        Objects.requireNonNull(rule, "rule cannot be null");
        long serverTotal = Pagination.parseCount("totalItems", totalItems);
        long next;
        switch (rule) {
            case REPLY_OFFSET:
                next = replyOffset != null
                        ? Pagination.parseCount("offset", replyOffset) : offset + serverRows;
                break;
            case PLUS_ROWS:
            case PLUS_SERVER_ROWS:
                next = offset + serverRows;
                break;
            case PLUS_MIN_ROWS_LIMIT:
                next = offset + Math.min(serverRows, limit);
                break;
            case PLUS_LIMIT:
                next = offset + limit;
                break;
            default:
                throw new IllegalStateException("unhandled offset rule " + rule);
        }
        return new OffsetPagination(limit, offset, Math.max(0L, serverTotal - excludedRows),
                next, offset < next && next < serverTotal);
    }

    /**
     * Returns the page size that was sent.
     *
     * @return the limit
     */
    public int getLimit() {
        return limit;
    }

    /**
     * Returns the offset that was sent.
     *
     * @return the offset of this page
     */
    public long getOffset() {
        return offset;
    }

    /**
     * Returns the total the server reported, reduced by the rows the SDK withheld.
     *
     * @return the total item count, 0 when the server reported none
     */
    public long getTotalItems() {
        return totalItems;
    }

    /**
     * Returns the offset to request the next page with.
     *
     * @return the next offset
     */
    public long getNextOffset() {
        return nextOffset;
    }

    /**
     * Reports whether another page exists.
     *
     * @return true while the walk should continue
     */
    public boolean hasMore() {
        return hasMore;
    }

    @Override
    public boolean equals(final Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof OffsetPagination)) {
            return false;
        }
        OffsetPagination that = (OffsetPagination) o;
        return limit == that.limit && offset == that.offset && totalItems == that.totalItems
                && nextOffset == that.nextOffset && hasMore == that.hasMore;
    }

    @Override
    public int hashCode() {
        return Objects.hash(limit, offset, totalItems, nextOffset, hasMore);
    }

    @Override
    public String toString() {
        return "OffsetPagination{limit=" + limit + ", offset=" + offset + ", totalItems=" + totalItems
                + ", nextOffset=" + nextOffset + ", hasMore=" + hasMore + "}";
    }
}
