package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of fee payers from an offset list, with its {@link OffsetPagination}.
 *
 * @see com.taurushq.sdk.protect.client.service.FeePayerService
 */
public final class FeePayerResult extends OffsetPagedResult<FeePayer> {

    /**
     * Creates a result.
     *
     * @param feePayers the rows of this page, null for none
     * @param pagination the page's pagination, required
     */
    public FeePayerResult(final List<FeePayer> feePayers, final OffsetPagination pagination) {
        super(feePayers, pagination);
    }

    /**
     * Gets the fee payers of this page.
     *
     * @return an unmodifiable list, empty when the page has no rows
     */
    public List<FeePayer> getFeePayers() {
        return items();
    }
}
