package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents an exchange account in the Taurus Protect system.
 * <p>
 * An exchange represents a connection to a third-party exchange (e.g., Binance,
 * Coinbase) that can be used for trading and transfers.
 *
 * @see ExchangeService
 */
public class Exchange {

    private String id;
    private String exchange;
    private String account;
    private String currency;
    private String type;
    private String totalBalance;
    private String status;
    private String container;
    private String label;
    private String displayLabel;
    private String baseCurrencyValuation;
    private Boolean hasWLA;
    private Currency currencyInfo;
    private OffsetDateTime creationDate;
    private OffsetDateTime updateDate;

    /**
     * Gets the unique identifier.
     *
     * @return the id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the unique identifier.
     *
     * @param id the id to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the exchange name (e.g., "binance", "coinbase").
     *
     * @return the exchange name
     */
    public String getExchange() {
        return exchange;
    }

    /**
     * Sets the exchange name.
     *
     * @param exchange the exchange name to set
     */
    public void setExchange(String exchange) {
        this.exchange = exchange;
    }

    /**
     * Gets the account identifier on the exchange.
     *
     * @return the account
     */
    public String getAccount() {
        return account;
    }

    /**
     * Sets the account identifier.
     *
     * @param account the account to set
     */
    public void setAccount(String account) {
        this.account = account;
    }

    /**
     * Gets the currency code.
     *
     * @return the currency
     */
    public String getCurrency() {
        return currency;
    }

    /**
     * Sets the currency code.
     *
     * @param currency the currency to set
     */
    public void setCurrency(String currency) {
        this.currency = currency;
    }

    /**
     * Gets the account type.
     *
     * @return the type
     */
    public String getType() {
        return type;
    }

    /**
     * Sets the account type.
     *
     * @param type the type to set
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Gets the total balance in the smallest currency unit.
     *
     * @return the total balance
     */
    public String getTotalBalance() {
        return totalBalance;
    }

    /**
     * Sets the total balance.
     *
     * @param totalBalance the total balance to set
     */
    public void setTotalBalance(String totalBalance) {
        this.totalBalance = totalBalance;
    }

    /**
     * Gets the account status.
     *
     * @return the status
     */
    public String getStatus() {
        return status;
    }

    /**
     * Sets the account status.
     *
     * @param status the status to set
     */
    public void setStatus(String status) {
        this.status = status;
    }

    /**
     * Gets the container.
     *
     * @return the container
     */
    public String getContainer() {
        return container;
    }

    /**
     * Sets the container.
     *
     * @param container the container to set
     */
    public void setContainer(String container) {
        this.container = container;
    }

    /**
     * Gets the label.
     *
     * @return the label
     */
    public String getLabel() {
        return label;
    }

    /**
     * Sets the label.
     *
     * @param label the label to set
     */
    public void setLabel(String label) {
        this.label = label;
    }

    /**
     * Gets the display label.
     *
     * @return the display label
     */
    public String getDisplayLabel() {
        return displayLabel;
    }

    /**
     * Sets the display label.
     *
     * @param displayLabel the display label to set
     */
    public void setDisplayLabel(String displayLabel) {
        this.displayLabel = displayLabel;
    }

    /**
     * Gets the valuation in the base currency (CHF, EUR, USD, etc.).
     *
     * @return the base currency valuation
     */
    public String getBaseCurrencyValuation() {
        return baseCurrencyValuation;
    }

    /**
     * Sets the base currency valuation.
     *
     * @param baseCurrencyValuation the base currency valuation to set
     */
    public void setBaseCurrencyValuation(String baseCurrencyValuation) {
        this.baseCurrencyValuation = baseCurrencyValuation;
    }

    /**
     * Checks if the exchange has a whitelisted address.
     *
     * @return true if has WLA
     */
    public Boolean getHasWLA() {
        return hasWLA;
    }

    /**
     * Sets whether the exchange has a whitelisted address.
     *
     * @param hasWLA the hasWLA flag to set
     */
    public void setHasWLA(Boolean hasWLA) {
        this.hasWLA = hasWLA;
    }

    /**
     * Gets the currency information.
     *
     * @return the currency info
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets the currency information.
     *
     * @param currencyInfo the currency info to set
     */
    public void setCurrencyInfo(Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }

    /**
     * Gets the creation date.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the creation date.
     *
     * @param creationDate the creation date to set
     */
    public void setCreationDate(OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Gets the last update date.
     *
     * @return the update date
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the last update date.
     *
     * @param updateDate the update date to set
     */
    public void setUpdateDate(OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }
}
