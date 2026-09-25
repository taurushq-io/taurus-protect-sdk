package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Result of an NFT collection balance query with cursor-based pagination.
 * <p>
 * Contains a page of NFT collection balances and cursor information for fetching
 * additional pages: pass {@code getPage().getNextCursor()} back while
 * {@code getPage().hasMore()} is true.
 *
 * @see NFTCollectionBalance
 */
public class NFTCollectionBalanceResult extends CursorPagedResult {

    private List<NFTCollectionBalance> balances;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the NFT collection balances of this page.
     *
     * @return the balances
     */
    public List<NFTCollectionBalance> getBalances() {
        return balances;
    }

    /**
     * Sets the NFT collection balances of this page.
     *
     * @param balances the balances
     */
    public void setBalances(final List<NFTCollectionBalance> balances) {
        this.balances = balances;
    }
}
