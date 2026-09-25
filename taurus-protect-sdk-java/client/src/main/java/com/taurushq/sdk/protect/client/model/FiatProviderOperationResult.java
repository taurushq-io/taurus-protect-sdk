package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of fiat provider operations.
 *
 * @see FiatService
 */
public class FiatProviderOperationResult extends CursorPagedResult {

    private List<FiatProviderOperation> operations;

    /**
     * Gets the fiat provider operations of this page.
     *
     * @return the operations
     */
    public List<FiatProviderOperation> getOperations() {
        return operations;
    }

    /**
     * Sets the fiat provider operations of this page.
     *
     * @param operations the operations
     */
    public void setOperations(final List<FiatProviderOperation> operations) {
        this.operations = operations;
    }
}
