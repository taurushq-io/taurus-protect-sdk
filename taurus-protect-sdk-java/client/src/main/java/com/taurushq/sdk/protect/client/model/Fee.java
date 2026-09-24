package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * The current network fee of a currency.
 *
 * @see com.taurushq.sdk.protect.client.service.FeeService
 */
public class Fee {

    private String currencyId;
    private String value;
    private String denom;
    private Currency currencyInfo;
    private OffsetDateTime updateDate;

    /**
     * Gets the id of the currency the fee applies to.
     *
     * @return the currency id
     */
    public String getCurrencyId() {
        return currencyId;
    }

    /**
     * Sets the currency id.
     *
     * @param currencyId the currency id
     */
    public void setCurrencyId(final String currencyId) {
        this.currencyId = currencyId;
    }

    /**
     * Gets the fee amount.
     *
     * @return the value
     */
    public String getValue() {
        return value;
    }

    /**
     * Sets the fee amount.
     *
     * @param value the value
     */
    public void setValue(final String value) {
        this.value = value;
    }

    /**
     * Gets the denomination the value is expressed in.
     *
     * @return the denomination
     */
    public String getDenom() {
        return denom;
    }

    /**
     * Sets the denomination.
     *
     * @param denom the denomination
     */
    public void setDenom(final String denom) {
        this.denom = denom;
    }

    /**
     * Gets the currency the fee applies to.
     *
     * @return the currency, may be null
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets the currency.
     *
     * @param currencyInfo the currency
     */
    public void setCurrencyInfo(final Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }

    /**
     * Gets when the fee was last updated.
     *
     * @return the update date
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the update date.
     *
     * @param updateDate the update date
     */
    public void setUpdateDate(final OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }
}
