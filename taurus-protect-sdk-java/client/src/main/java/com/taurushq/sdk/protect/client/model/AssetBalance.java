package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents the balance of a specific asset (cryptocurrency or token).
 * <p>
 * Asset balances combine asset identification with the current balance amounts,
 * providing a complete view of holdings for a particular digital asset.
 * This is commonly used when retrieving balances across multiple assets
 * for a wallet or address.
 *
 * @see Asset
 * @see Balance
 */
public class AssetBalance {

    /**
     * The asset (cryptocurrency or token) this balance represents.
     */
    private Asset asset;

    /**
     * The balance amounts for this asset (total, available, pending, etc.).
     */
    private Balance balance;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the asset this balance represents.
     *
     * @return the asset (cryptocurrency or token)
     */
    public Asset getAsset() {
        return asset;
    }

    /**
     * Sets the asset this balance represents.
     *
     * @param asset the asset to set
     */
    public void setAsset(Asset asset) {
        this.asset = asset;
    }

    /**
     * Returns the balance amounts for this asset.
     *
     * @return the balance including total, available, and pending amounts
     */
    public Balance getBalance() {
        return balance;
    }

    /**
     * Sets the balance amounts for this asset.
     *
     * @param balance the balance to set
     */
    public void setBalance(Balance balance) {
        this.balance = balance;
    }
}