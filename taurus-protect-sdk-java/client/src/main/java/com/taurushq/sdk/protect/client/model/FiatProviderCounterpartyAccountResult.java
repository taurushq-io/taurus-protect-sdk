package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of fiat provider counterparty accounts.
 *
 * @see FiatService
 */
public class FiatProviderCounterpartyAccountResult extends CursorPagedResult {

    private List<FiatProviderCounterpartyAccount> accounts;

    /**
     * Gets the fiat provider counterparty accounts of this page.
     *
     * @return the accounts
     */
    public List<FiatProviderCounterpartyAccount> getAccounts() {
        return accounts;
    }

    /**
     * Sets the fiat provider counterparty accounts of this page.
     *
     * @param accounts the accounts
     */
    public void setAccounts(final List<FiatProviderCounterpartyAccount> accounts) {
        this.accounts = accounts;
    }
}
