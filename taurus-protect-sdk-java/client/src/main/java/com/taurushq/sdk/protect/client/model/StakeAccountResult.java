package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of stake accounts.
 * <p>
 * This class wraps a list of stake accounts along with pagination information
 * to support cursor-based navigation.
 *
 * @see StakeAccount
 * @see ApiResponseCursor
 */
public class StakeAccountResult {

    private List<StakeAccount> stakeAccounts;
    private ApiResponseCursor cursor;

    /**
     * Gets the list of stake accounts.
     *
     * @return the stake accounts
     */
    public List<StakeAccount> getStakeAccounts() {
        return stakeAccounts;
    }

    /**
     * Sets the list of stake accounts.
     *
     * @param stakeAccounts the stake accounts
     */
    public void setStakeAccounts(List<StakeAccount> stakeAccounts) {
        this.stakeAccounts = stakeAccounts;
    }

    /**
     * Gets the pagination cursor.
     *
     * @return the cursor
     */
    public ApiResponseCursor getCursor() {
        return cursor;
    }

    /**
     * Sets the pagination cursor.
     *
     * @param cursor the cursor
     */
    public void setCursor(ApiResponseCursor cursor) {
        this.cursor = cursor;
    }

    /**
     * Checks if there are more results available.
     *
     * @return true if more results exist
     */
    public boolean hasNext() {
        return cursor != null && cursor.hasNext();
    }

    /**
     * Creates a cursor for fetching the next page.
     *
     * @param pageSize the page size for the next request
     * @return a cursor for the next page, or null if no more pages
     */
    public ApiRequestCursor nextCursor(int pageSize) {
        if (!hasNext()) {
            return null;
        }
        return new ApiRequestCursor(cursor.getCurrentPage(), PageRequest.NEXT, pageSize);
    }
}
