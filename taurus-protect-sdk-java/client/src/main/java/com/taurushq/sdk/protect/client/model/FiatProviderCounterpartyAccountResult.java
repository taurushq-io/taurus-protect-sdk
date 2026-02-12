package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of fiat provider counterparty accounts.
 *
 * @see FiatService
 */
public class FiatProviderCounterpartyAccountResult {

    private List<FiatProviderCounterpartyAccount> accounts;
    private ApiResponseCursor cursor;

    public List<FiatProviderCounterpartyAccount> getAccounts() {
        return accounts;
    }

    public void setAccounts(final List<FiatProviderCounterpartyAccount> accounts) {
        this.accounts = accounts;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
