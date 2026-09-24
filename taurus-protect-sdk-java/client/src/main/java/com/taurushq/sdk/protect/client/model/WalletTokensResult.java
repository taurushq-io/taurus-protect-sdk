package com.taurushq.sdk.protect.client.model;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Objects;

/**
 * One page of a wallet's token balances, with its {@link CursorPage}.
 * <p>
 * This list pages by an opaque token: pass {@code getPage().getNextCursor()} back while
 * {@code getPage().hasMore()} is true. The page carries the server's total.
 *
 * @see com.taurushq.sdk.protect.client.service.WalletService
 */
public final class WalletTokensResult {

    private final List<AssetBalance> balances;
    private final CursorPage page;

    /**
     * Creates a result.
     *
     * @param balances the balances of this page, null for none
     * @param page     the page, required
     */
    public WalletTokensResult(final List<AssetBalance> balances, final CursorPage page) {
        this.balances = balances == null ? Collections.<AssetBalance>emptyList()
                : Collections.unmodifiableList(new ArrayList<>(balances));
        this.page = Objects.requireNonNull(page, "page cannot be null");
    }

    /**
     * Gets the token balances of this page.
     *
     * @return an unmodifiable list, empty when the page has no rows
     */
    public List<AssetBalance> getBalances() {
        return balances;
    }

    /**
     * Gets the page: next cursor, whether more pages exist, and the server total.
     *
     * @return the page, never null
     */
    public CursorPage getPage() {
        return page;
    }
}
