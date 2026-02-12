package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;

import java.util.List;

/**
 * Result containing a list of lending agreements with pagination.
 */
public class LendingAgreementResult {

    private List<LendingAgreement> agreements;
    private ApiResponseCursor cursor;

    public List<LendingAgreement> getAgreements() {
        return agreements;
    }

    public void setAgreements(final List<LendingAgreement> agreements) {
        this.agreements = agreements;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
