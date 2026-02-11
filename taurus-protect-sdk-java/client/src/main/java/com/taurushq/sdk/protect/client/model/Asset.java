package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a digital asset (cryptocurrency or token) in the Taurus Protect system.
 * <p>
 * An asset identifies a specific cryptocurrency or token by its currency code
 * and kind (native coin, ERC-20 token, etc.). The asset also includes
 * detailed currency information such as decimals and display name.
 *
 * @see Currency
 * @see AssetBalance
 */
public class Asset {

    /**
     * The currency code or identifier (e.g., "ETH", "BTC", "USDC").
     */
    private String currency;

    /**
     * The type/kind of asset (e.g., "native", "erc20", "fa2").
     */
    private String kind;

    /**
     * Detailed information about the currency including decimals and display name.
     */
    private Currency currencyInfo;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the currency code or identifier.
     *
     * @return the currency code (e.g., "ETH", "BTC", "USDC")
     */
    public String getCurrency() {
        return currency;
    }

    /**
     * Sets the currency code or identifier.
     *
     * @param currency the currency code to set
     */
    public void setCurrency(String currency) {
        this.currency = currency;
    }

    /**
     * Returns the type/kind of this asset.
     *
     * @return the asset kind (e.g., "native", "erc20", "fa2")
     */
    public String getKind() {
        return kind;
    }

    /**
     * Sets the type/kind of this asset.
     *
     * @param kind the asset kind to set
     */
    public void setKind(String kind) {
        this.kind = kind;
    }

    /**
     * Returns detailed currency information.
     *
     * @return the currency information including decimals and display name
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets detailed currency information.
     *
     * @param currencyInfo the currency information to set
     */
    public void setCurrencyInfo(Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }
}