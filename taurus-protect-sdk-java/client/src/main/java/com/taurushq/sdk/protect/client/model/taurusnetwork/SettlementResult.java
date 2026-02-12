package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;

import java.util.List;

/**
 * Result containing a list of settlements with pagination.
 */
public class SettlementResult {

    private List<Settlement> settlements;
    private ApiResponseCursor cursor;

    public List<Settlement> getSettlements() {
        return settlements;
    }

    public void setSettlements(final List<Settlement> settlements) {
        this.settlements = settlements;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
