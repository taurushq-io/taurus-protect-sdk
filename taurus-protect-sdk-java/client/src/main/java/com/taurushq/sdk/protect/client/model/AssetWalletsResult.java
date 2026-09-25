package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of wallets holding an asset from a cursor list. Continue with {@code getPage()}.
 *
 * @see com.taurushq.sdk.protect.client.service.AssetService
 */
public class AssetWalletsResult extends CursorPagedResult {

    private List<Wallet> wallets;

    /**
     * Gets the wallets holding an asset of this page.
     *
     * @return the wallets
     */
    public List<Wallet> getWallets() {
        return wallets;
    }

    /**
     * Sets the wallets holding an asset of this page.
     *
     * @param wallets the wallets
     */
    public void setWallets(final List<Wallet> wallets) {
        this.wallets = wallets;
    }
}
