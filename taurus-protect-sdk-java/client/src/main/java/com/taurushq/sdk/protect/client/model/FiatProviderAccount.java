package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents a fiat provider account in the Taurus Protect system.
 *
 * @see FiatService
 */
public class FiatProviderAccount {

    private String id;
    private String provider;
    private String label;
    private String accountType;
    private String accountIdentifier;
    private String accountName;
    private String totalBalance;
    private String currencyID;
    private String baseCurrencyValuation;
    private OffsetDateTime creationDate;
    private OffsetDateTime updateDate;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getProvider() {
        return provider;
    }

    public void setProvider(final String provider) {
        this.provider = provider;
    }

    public String getLabel() {
        return label;
    }

    public void setLabel(final String label) {
        this.label = label;
    }

    public String getAccountType() {
        return accountType;
    }

    public void setAccountType(final String accountType) {
        this.accountType = accountType;
    }

    public String getAccountIdentifier() {
        return accountIdentifier;
    }

    public void setAccountIdentifier(final String accountIdentifier) {
        this.accountIdentifier = accountIdentifier;
    }

    public String getAccountName() {
        return accountName;
    }

    public void setAccountName(final String accountName) {
        this.accountName = accountName;
    }

    public String getTotalBalance() {
        return totalBalance;
    }

    public void setTotalBalance(final String totalBalance) {
        this.totalBalance = totalBalance;
    }

    public String getCurrencyID() {
        return currencyID;
    }

    public void setCurrencyID(final String currencyID) {
        this.currencyID = currencyID;
    }

    public String getBaseCurrencyValuation() {
        return baseCurrencyValuation;
    }

    public void setBaseCurrencyValuation(final String baseCurrencyValuation) {
        this.baseCurrencyValuation = baseCurrencyValuation;
    }

    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    public void setCreationDate(final OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    public void setUpdateDate(final OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }
}
