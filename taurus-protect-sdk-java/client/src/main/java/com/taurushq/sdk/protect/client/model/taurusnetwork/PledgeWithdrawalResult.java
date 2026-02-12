package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;

import java.util.List;

/**
 * Result containing a list of pledge withdrawals with pagination.
 */
public class PledgeWithdrawalResult {

    private List<PledgeWithdrawal> withdrawals;
    private ApiResponseCursor cursor;

    public List<PledgeWithdrawal> getWithdrawals() {
        return withdrawals;
    }

    public void setWithdrawals(final List<PledgeWithdrawal> withdrawals) {
        this.withdrawals = withdrawals;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
