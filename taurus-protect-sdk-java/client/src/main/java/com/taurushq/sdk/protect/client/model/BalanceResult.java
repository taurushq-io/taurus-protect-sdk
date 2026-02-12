package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Result of a balance query with cursor-based pagination.
 * <p>
 * Contains a page of asset balances and cursor information for fetching
 * additional pages. Use {@link #hasNext()} to check for more pages and
 * {@link #nextCursor(long)} to create a cursor for the next page.
 *
 * @see AssetBalance
 * @see ApiResponseCursor
 */
public class BalanceResult {

    /**
     * The list of asset balances in this page of results.
     */
    private List<AssetBalance> balances;

    /**
     * Pagination cursor containing page information.
     */
    private ApiResponseCursor cursor;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the list of balances.
     *
     * @return the balances
     */
    public List<AssetBalance> getBalances() {
        return balances;
    }

    /**
     * Sets the list of balances.
     *
     * @param balances the balances
     */
    public void setBalances(List<AssetBalance> balances) {
        this.balances = balances;
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
     * Checks if there is a next result
     *
     * @return a boolean indicating if there is a next result
     */
    public boolean hasNext() {
        return cursor != null && cursor.hasNext();
    }
}
