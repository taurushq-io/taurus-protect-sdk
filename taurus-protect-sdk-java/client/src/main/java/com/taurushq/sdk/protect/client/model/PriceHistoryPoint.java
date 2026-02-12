package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;

/**
 * Represents a single data point in a price history timeline (OHLCV candlestick data).
 * <p>
 * Price history points capture OHLCV (Open, High, Low, Close, Volume) data for a specific
 * time period, enabling historical price analysis and charting. This follows standard
 * financial candlestick data format used in trading and portfolio analysis.
 *
 * @see Price
 */
public class PriceHistoryPoint {

    /**
     * The start timestamp of the time period this data point covers.
     */
    private OffsetDateTime periodStartDate;

    /**
     * The blockchain this price data applies to (e.g., "ethereum", "bitcoin").
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
     * The highest price during this time period.
     */
    private String high;

    /**
     * The lowest price during this time period.
     */
    private String low;

    /**
     * The opening price at the start of this time period.
     */
    private String open;

    /**
     * The closing price at the end of this time period.
     */
    private String close;

    /**
     * The trading volume in the source currency during this period.
     */
    private String volumeFrom;

    /**
     * The trading volume in the target currency during this period.
     */
    private String volumeTo;

    /**
     * The percentage price change during this time period.
     */
    private String changePercent;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the start timestamp of the time period this data point covers.
     *
     * @return the period start timestamp
     */
    public OffsetDateTime getPeriodStartDate() {
        return periodStartDate;
    }

    /**
     * Sets the start timestamp of the time period this data point covers.
     *
     * @param periodStartDate the period start timestamp to set
     */
    public void setPeriodStartDate(OffsetDateTime periodStartDate) {
        this.periodStartDate = periodStartDate;
    }

    /**
     * Returns the blockchain this price data applies to.
     *
     * @return the blockchain identifier (e.g., "ethereum", "bitcoin")
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain this price data applies to.
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
     * Returns the highest price during this time period.
     *
     * @return the high price
     */
    public String getHigh() {
        return high;
    }

    /**
     * Sets the highest price during this time period.
     *
     * @param high the high price to set
     */
    public void setHigh(String high) {
        this.high = high;
    }

    /**
     * Returns the lowest price during this time period.
     *
     * @return the low price
     */
    public String getLow() {
        return low;
    }

    /**
     * Sets the lowest price during this time period.
     *
     * @param low the low price to set
     */
    public void setLow(String low) {
        this.low = low;
    }

    /**
     * Returns the opening price at the start of this time period.
     *
     * @return the open price
     */
    public String getOpen() {
        return open;
    }

    /**
     * Sets the opening price at the start of this time period.
     *
     * @param open the open price to set
     */
    public void setOpen(String open) {
        this.open = open;
    }

    /**
     * Returns the closing price at the end of this time period.
     *
     * @return the close price
     */
    public String getClose() {
        return close;
    }

    /**
     * Sets the closing price at the end of this time period.
     *
     * @param close the close price to set
     */
    public void setClose(String close) {
        this.close = close;
    }

    /**
     * Returns the trading volume in the source currency during this period.
     *
     * @return the volume in source currency
     */
    public String getVolumeFrom() {
        return volumeFrom;
    }

    /**
     * Sets the trading volume in the source currency during this period.
     *
     * @param volumeFrom the volume in source currency to set
     */
    public void setVolumeFrom(String volumeFrom) {
        this.volumeFrom = volumeFrom;
    }

    /**
     * Returns the trading volume in the target currency during this period.
     *
     * @return the volume in target currency
     */
    public String getVolumeTo() {
        return volumeTo;
    }

    /**
     * Sets the trading volume in the target currency during this period.
     *
     * @param volumeTo the volume in target currency to set
     */
    public void setVolumeTo(String volumeTo) {
        this.volumeTo = volumeTo;
    }

    /**
     * Returns the percentage price change during this time period.
     *
     * @return the price change percentage
     */
    public String getChangePercent() {
        return changePercent;
    }

    /**
     * Sets the percentage price change during this time period.
     *
     * @param changePercent the price change percentage to set
     */
    public void setChangePercent(String changePercent) {
        this.changePercent = changePercent;
    }
}
