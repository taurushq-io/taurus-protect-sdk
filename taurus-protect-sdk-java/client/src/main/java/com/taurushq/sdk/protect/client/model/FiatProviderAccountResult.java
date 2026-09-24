package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of fiat provider accounts.
 *
 * @see FiatService
 */
public class FiatProviderAccountResult extends CursorPagedResult {

    private List<FiatProviderAccount> accounts;

    /**
     * Gets the fiat provider accounts of this page.
     *
     * @return the accounts
     */
    public List<FiatProviderAccount> getAccounts() {
        return accounts;
    }

    /**
     * Sets the fiat provider accounts of this page.
     *
     * @param accounts the accounts
     */
    public void setAccounts(final List<FiatProviderAccount> accounts) {
        this.accounts = accounts;
    }
}
