package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;

/**
 * Represents a currency price/exchange rate in the Taurus Protect system.
 * <p>
 * Prices track exchange rates between currencies (e.g., ETH to USD) along with
 * metadata about the rate source, precision, and price changes. Prices are
 * used for portfolio valuation, transaction display, and reporting.
 *
 * @see PriceHistoryPoint
 * @see Currency
 */
public class Price {

    /**
     * The blockchain this price applies to (e.g., "ethereum", "bitcoin").
     */
    private String blockchain;

    /**
     * The source currency code being converted from (e.g., "ETH", "BTC").
     */
    private String currencyFrom;

    /**
     * The target currency code being converted to (e.g., "USD", "EUR").
     */
    private String currencyTo;

    /**
     * The number of decimal places for the rate precision.
     */
    private String decimals;

    /**
     * The exchange rate value (amount of currencyTo per unit of currencyFrom).
     */
    private String rate;

    /**
     * The percentage price change over the last 24 hours.
     */
    private String changePercent24Hour;

    /**
     * The data source for this price (e.g., "coingecko", "cryptocompare").
     */
    private String source;

    /**
     * Timestamp when this price record was created.
     */
    private OffsetDateTime createdAt;

    /**
     * Timestamp when this price record was last updated.
     */
    private OffsetDateTime updatedAt;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the blockchain this price applies to.
     *
     * @return the blockchain identifier (e.g., "ethereum", "bitcoin")
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain this price applies to.
     *
     * @param blockchain the blockchain identifier to set
     */
    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Returns the source currency code being converted from.
     *
     * @return the source currency code (e.g., "ETH", "BTC")
     */
    public String getCurrencyFrom() {
        return currencyFrom;
    }

    /**
     * Sets the source currency code being converted from.
     *
     * @param currencyFrom the source currency code to set
     */
    public void setCurrencyFrom(String currencyFrom) {
        this.currencyFrom = currencyFrom;
    }

    /**
     * Returns the target currency code being converted to.
     *
     * @return the target currency code (e.g., "USD", "EUR")
     */
    public String getCurrencyTo() {
        return currencyTo;
    }

    /**
     * Sets the target currency code being converted to.
     *
     * @param currencyTo the target currency code to set
     */
    public void setCurrencyTo(String currencyTo) {
        this.currencyTo = currencyTo;
    }

    /**
     * Returns the number of decimal places for rate precision.
     *
     * @return the decimal precision as a string
     */
    public String getDecimals() {
        return decimals;
    }

    /**
     * Sets the number of decimal places for rate precision.
     *
     * @param decimals the decimal precision to set
     */
    public void setDecimals(String decimals) {
        this.decimals = decimals;
    }

    /**
     * Returns the exchange rate value.
     *
     * @return the rate (amount of currencyTo per unit of currencyFrom)
     */
    public String getRate() {
        return rate;
    }

    /**
     * Sets the exchange rate value.
     *
     * @param rate the rate to set
     */
    public void setRate(String rate) {
        this.rate = rate;
    }

    /**
     * Returns the percentage price change over the last 24 hours.
     *
     * @return the 24-hour price change percentage
     */
    public String getChangePercent24Hour() {
        return changePercent24Hour;
    }

    /**
     * Sets the percentage price change over the last 24 hours.
     *
     * @param changePercent24Hour the 24-hour change percentage to set
     */
    public void setChangePercent24Hour(String changePercent24Hour) {
        this.changePercent24Hour = changePercent24Hour;
    }

    /**
     * Returns the data source for this price.
     *
     * @return the price source (e.g., "coingecko", "cryptocompare")
     */
    public String getSource() {
        return source;
    }

    /**
     * Sets the data source for this price.
     *
     * @param source the price source to set
     */
    public void setSource(String source) {
        this.source = source;
    }

    /**
     * Returns the timestamp when this price record was created.
     *
     * @return the creation timestamp
     */
    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    /**
     * Sets the timestamp when this price record was created.
     *
     * @param createdAt the creation timestamp to set
     */
    public void setCreatedAt(OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }

    /**
     * Returns the timestamp when this price record was last updated.
     *
     * @return the last update timestamp
     */
    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    /**
     * Sets the timestamp when this price record was last updated.
     *
     * @param updatedAt the last update timestamp to set
     */
    public void setUpdatedAt(OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }
}
