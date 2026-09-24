package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.CursorPagedResult;

import java.util.List;

/**
 * Result containing a list of lending agreements with pagination.
 */
public class LendingAgreementResult extends CursorPagedResult {

    private List<LendingAgreement> agreements;

    /**
     * Gets the lending agreements of this page.
     *
     * @return the agreements
     */
    public List<LendingAgreement> getAgreements() {
        return agreements;
    }

    /**
     * Sets the lending agreements of this page.
     *
     * @param agreements the agreements
     */
    public void setAgreements(final List<LendingAgreement> agreements) {
        this.agreements = agreements;
    }
}
