package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of fiat provider accounts.
 *
 * @see FiatService
 */
public class FiatProviderAccountResult {

    private List<FiatProviderAccount> accounts;
    private ApiResponseCursor cursor;

    public List<FiatProviderAccount> getAccounts() {
        return accounts;
    }

    public void setAccounts(final List<FiatProviderAccount> accounts) {
        this.accounts = accounts;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
