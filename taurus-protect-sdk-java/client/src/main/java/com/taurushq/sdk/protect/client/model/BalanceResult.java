package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Result of a balance query with cursor-based pagination.
 * <p>
 * Contains a page of asset balances and cursor information for fetching
 * additional pages: pass {@code getPage().getNextCursor()} back while
 * {@code getPage().hasMore()} is true.
 *
 * @see AssetBalance
 */
public class BalanceResult extends CursorPagedResult {

    private List<AssetBalance> balances;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the asset balances of this page.
     *
     * @return the balances
     */
    public List<AssetBalance> getBalances() {
        return balances;
    }

    /**
     * Sets the asset balances of this page.
     *
     * @param balances the balances
     */
    public void setBalances(final List<AssetBalance> balances) {
        this.balances = balances;
    }
}
