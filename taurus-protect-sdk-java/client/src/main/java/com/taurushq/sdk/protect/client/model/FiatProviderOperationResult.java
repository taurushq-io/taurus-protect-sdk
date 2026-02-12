package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of fiat provider operations.
 *
 * @see FiatService
 */
public class FiatProviderOperationResult {

    private List<FiatProviderOperation> operations;
    private ApiResponseCursor cursor;

    public List<FiatProviderOperation> getOperations() {
        return operations;
    }

    public void setOperations(final List<FiatProviderOperation> operations) {
        this.operations = operations;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
