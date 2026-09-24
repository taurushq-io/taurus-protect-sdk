package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.CursorPagedResult;

import java.util.List;

/**
 * Result containing a list of settlements with pagination.
 */
public class SettlementResult extends CursorPagedResult {

    private List<Settlement> settlements;

    /**
     * Gets the settlements of this page.
     *
     * @return the settlements
     */
    public List<Settlement> getSettlements() {
        return settlements;
    }

    /**
     * Sets the settlements of this page.
     *
     * @param settlements the settlements
     */
    public void setSettlements(final List<Settlement> settlements) {
        this.settlements = settlements;
    }
}
