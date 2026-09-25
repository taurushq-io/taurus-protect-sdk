package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of stake accounts.
 * <p>
 * This class wraps a list of stake accounts along with pagination information
 * to support cursor-based navigation.
 *
 * @see StakeAccount
 */
public class StakeAccountResult extends CursorPagedResult {

    private List<StakeAccount> stakeAccounts;

    /**
     * Gets the stake accounts of this page.
     *
     * @return the stakeAccounts
     */
    public List<StakeAccount> getStakeAccounts() {
        return stakeAccounts;
    }

    /**
     * Sets the stake accounts of this page.
     *
     * @param stakeAccounts the stakeAccounts
     */
    public void setStakeAccounts(final List<StakeAccount> stakeAccounts) {
        this.stakeAccounts = stakeAccounts;
    }
}
