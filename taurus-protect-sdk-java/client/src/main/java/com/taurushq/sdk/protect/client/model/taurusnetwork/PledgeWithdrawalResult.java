package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.CursorPagedResult;

import java.util.List;

/**
 * Result containing a list of pledge withdrawals with pagination.
 */
public class PledgeWithdrawalResult extends CursorPagedResult {

    private List<PledgeWithdrawal> withdrawals;

    /**
     * Gets the pledge withdrawals of this page.
     *
     * @return the withdrawals
     */
    public List<PledgeWithdrawal> getWithdrawals() {
        return withdrawals;
    }

    /**
     * Sets the pledge withdrawals of this page.
     *
     * @param withdrawals the withdrawals
     */
    public void setWithdrawals(final List<PledgeWithdrawal> withdrawals) {
        this.withdrawals = withdrawals;
    }
}
