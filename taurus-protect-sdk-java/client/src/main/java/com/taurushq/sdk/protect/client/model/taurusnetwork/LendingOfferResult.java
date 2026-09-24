package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.CursorPagedResult;

import java.util.List;

/**
 * Result containing a list of lending offers with pagination.
 */
public class LendingOfferResult extends CursorPagedResult {

    private List<LendingOffer> offers;

    /**
     * Gets the lending offers of this page.
     *
     * @return the offers
     */
    public List<LendingOffer> getOffers() {
        return offers;
    }

    /**
     * Sets the lending offers of this page.
     *
     * @param offers the offers
     */
    public void setOffers(final List<LendingOffer> offers) {
        this.offers = offers;
    }
}
