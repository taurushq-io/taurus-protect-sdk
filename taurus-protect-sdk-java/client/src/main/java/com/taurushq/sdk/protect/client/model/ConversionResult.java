package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents the result of a currency conversion operation.
 * <p>
 * Conversion results contain both the raw value and a human-readable main unit value
 * (accounting for decimals), along with currency metadata. This is typically used
 * when displaying converted amounts in the UI or for reporting purposes.
 *
 * @see Currency
 * @see Price
 */
public class ConversionResult {

    /**
     * The currency symbol for display (e.g., "USD", "EUR", "ETH").
     */
    private String symbol;

    /**
     * The converted value in the smallest unit (e.g., wei for ETH, satoshi for BTC).
     */
    private String value;

    /**
     * The converted value in main units (e.g., ETH instead of wei), formatted for display.
     */
    private String mainUnitValue;

    /**
     * Detailed information about the target currency.
     */
    private Currency currencyInfo;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the currency symbol for display.
     *
     * @return the currency symbol (e.g., "USD", "EUR", "ETH")
     */
    public String getSymbol() {
        return symbol;
    }

    /**
     * Sets the currency symbol for display.
     *
     * @param symbol the currency symbol to set
     */
    public void setSymbol(String symbol) {
        this.symbol = symbol;
    }

    /**
     * Returns the converted value in the smallest unit.
     *
     * @return the value in smallest units (e.g., wei, satoshi)
     */
    public String getValue() {
        return value;
    }

    /**
     * Sets the converted value in the smallest unit.
     *
     * @param value the value in smallest units to set
     */
    public void setValue(String value) {
        this.value = value;
    }

    /**
     * Returns the converted value in main units, formatted for display.
     *
     * @return the value in main units (e.g., ETH instead of wei)
     */
    public String getMainUnitValue() {
        return mainUnitValue;
    }

    /**
     * Sets the converted value in main units.
     *
     * @param mainUnitValue the value in main units to set
     */
    public void setMainUnitValue(String mainUnitValue) {
        this.mainUnitValue = mainUnitValue;
    }

    /**
     * Returns detailed information about the target currency.
     *
     * @return the currency metadata
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets detailed information about the target currency.
     *
     * @param currencyInfo the currency metadata to set
     */
    public void setCurrencyInfo(Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }
}
