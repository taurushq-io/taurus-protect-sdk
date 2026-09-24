package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of wallets from an offset list, with its {@link OffsetPagination}.
 *
 * @see com.taurushq.sdk.protect.client.service.WalletService
 */
public final class WalletResult extends OffsetPagedResult<Wallet> {

    /**
     * Creates a result.
     *
     * @param wallets the rows of this page, null for none
     * @param pagination the page's pagination, required
     */
    public WalletResult(final List<Wallet> wallets, final OffsetPagination pagination) {
        super(wallets, pagination);
    }

    /**
     * Gets the wallets of this page.
     *
     * @return an unmodifiable list, empty when the page has no rows
     */
    public List<Wallet> getWallets() {
        return items();
    }
}
