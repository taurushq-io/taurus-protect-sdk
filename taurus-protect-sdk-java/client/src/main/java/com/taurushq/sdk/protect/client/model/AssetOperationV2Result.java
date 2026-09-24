package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of operations on an asset (v2 asset service) from a cursor list. Continue with {@code getPage()}.
 *
 * @see com.taurushq.sdk.protect.client.service.AssetService
 */
public class AssetOperationV2Result extends CursorPagedResult {

    private List<AssetOperationV2> operations;

    /**
     * Gets the operations on an asset of this page.
     *
     * @return the operations
     */
    public List<AssetOperationV2> getOperations() {
        return operations;
    }

    /**
     * Sets the operations on an asset of this page.
     *
     * @param operations the operations
     */
    public void setOperations(final List<AssetOperationV2> operations) {
        this.operations = operations;
    }
}
