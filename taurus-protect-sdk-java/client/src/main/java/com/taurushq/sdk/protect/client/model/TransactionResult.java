package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of transactions from an offset list, with its {@link OffsetPagination}.
 *
 * @see com.taurushq.sdk.protect.client.service.TransactionService
 */
public final class TransactionResult extends OffsetPagedResult<Transaction> {

    /**
     * Creates a result.
     *
     * @param transactions the rows of this page, null for none
     * @param pagination the page's pagination, required
     */
    public TransactionResult(final List<Transaction> transactions, final OffsetPagination pagination) {
        super(transactions, pagination);
    }

    /**
     * Gets the transactions of this page.
     *
     * @return an unmodifiable list, empty when the page has no rows
     */
    public List<Transaction> getTransactions() {
        return items();
    }
}
