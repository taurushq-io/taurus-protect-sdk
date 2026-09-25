package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.CursorPagedResult;

import java.util.List;

/**
 * Result containing a list of pledges with pagination.
 */
public class PledgeResult extends CursorPagedResult {

    private List<Pledge> pledges;

    /**
     * Gets the pledges of this page.
     *
     * @return the pledges
     */
    public List<Pledge> getPledges() {
        return pledges;
    }

    /**
     * Sets the pledges of this page.
     *
     * @param pledges the pledges
     */
    public void setPledges(final List<Pledge> pledges) {
        this.pledges = pledges;
    }
}
