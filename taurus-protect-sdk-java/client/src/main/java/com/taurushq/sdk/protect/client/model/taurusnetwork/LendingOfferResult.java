package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;

import java.util.List;

/**
 * Result containing a list of lending offers with pagination.
 */
public class LendingOfferResult {

    private List<LendingOffer> offers;
    private ApiResponseCursor cursor;

    public List<LendingOffer> getOffers() {
        return offers;
    }

    public void setOffers(final List<LendingOffer> offers) {
        this.offers = offers;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
