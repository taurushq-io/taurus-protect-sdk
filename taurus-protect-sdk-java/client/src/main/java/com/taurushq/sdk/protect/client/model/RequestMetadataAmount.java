package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents the amount details extracted from request metadata.
 * <p>
 * This class contains the transaction amount in both the source and target
 * currencies, along with the conversion rate and decimal precision used.
 * It supports currency conversion tracking for cross-currency transactions.
 *
 * @see RequestMetadata#getAmount()
 * @see Request
 */
public class RequestMetadataAmount {

    /**
     * The amount value in the source currency's smallest unit (e.g., satoshis, wei).
     * Stored as a string representation to preserve arbitrary precision.
     */
    private String valueFrom;

    /**
     * The converted amount value in the target currency.
     * Stored as a string representation to preserve arbitrary precision.
     */
    private String valueTo;

    /**
     * The exchange rate used for currency conversion.
     * Stored as a string representation to preserve arbitrary precision.
     */
    private String rate;

    /**
     * The number of decimal places for the source currency.
     */
    private int decimals;

    /**
     * The source currency code (e.g., "BTC", "ETH").
     */
    private String currencyFrom;

    /**
     * The target currency code (e.g., "USD", "EUR").
     */
    private String currencyTo;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the amount value in the source currency's smallest unit.
     * <p>
     * For example, for Bitcoin this would be in satoshis, for Ethereum in wei.
     * Returned as a string representation to preserve arbitrary precision.
     *
     * @return the amount in the source currency's smallest unit
     */
    public String getValueFrom() {
        return valueFrom;
    }

    /**
     * Sets the amount value in the source currency's smallest unit.
     *
     * @param valueFrom the amount in the source currency's smallest unit
     */
    public void setValueFrom(String valueFrom) {
        this.valueFrom = valueFrom;
    }

    /**
     * Returns the converted amount value in the target currency.
     * Returned as a string representation to preserve arbitrary precision.
     *
     * @return the converted amount in the target currency
     */
    public String getValueTo() {
        return valueTo;
    }

    /**
     * Sets the converted amount value in the target currency.
     *
     * @param valueTo the converted amount in the target currency
     */
    public void setValueTo(String valueTo) {
        this.valueTo = valueTo;
    }

    /**
     * Returns the exchange rate used for currency conversion.
     * Returned as a string representation to preserve arbitrary precision.
     *
     * @return the exchange rate from source to target currency
     */
    public String getRate() {
        return rate;
    }

    /**
     * Sets the exchange rate used for currency conversion.
     *
     * @param rate the exchange rate from source to target currency
     */
    public void setRate(String rate) {
        this.rate = rate;
    }

    /**
     * Returns the number of decimal places for the source currency.
     * <p>
     * For example, Bitcoin has 8 decimals, Ethereum has 18.
     *
     * @return the number of decimal places
     */
    public int getDecimals() {
        return decimals;
    }

    /**
     * Sets the number of decimal places for the source currency.
     *
     * @param decimals the number of decimal places
     */
    public void setDecimals(int decimals) {
        this.decimals = decimals;
    }

    /**
     * Returns the source currency code.
     *
     * @return the source currency code (e.g., "BTC", "ETH")
     */
    public String getCurrencyFrom() {
        return currencyFrom;
    }

    /**
     * Sets the source currency code.
     *
     * @param currencyFrom the source currency code (e.g., "BTC", "ETH")
     */
    public void setCurrencyFrom(String currencyFrom) {
        this.currencyFrom = currencyFrom;
    }

    /**
     * Returns the target currency code.
     *
     * @return the target currency code (e.g., "USD", "EUR")
     */
    public String getCurrencyTo() {
        return currencyTo;
    }

    /**
     * Sets the target currency code.
     *
     * @param currencyTo the target currency code (e.g., "USD", "EUR")
     */
    public void setCurrencyTo(String currencyTo) {
        this.currencyTo = currencyTo;
    }
}
