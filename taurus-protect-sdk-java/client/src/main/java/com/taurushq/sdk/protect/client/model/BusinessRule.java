package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a business rule in the Taurus Protect system.
 * <p>
 * Business rules define operational constraints and configurations at various
 * scopes (tenant, wallet, address, currency). Rules can control behaviors such
 * as transaction limits, approval requirements, and other policy settings.
 * <p>
 * Rules are identified by a key-value pair and can be associated with specific
 * entities (wallets, addresses) or apply globally to a tenant.
 *
 * @see GovernanceRules
 * @see Currency
 */
public class BusinessRule {

    /**
     * Unique identifier for this business rule.
     */
    private String id;

    /**
     * ID of the tenant (organization) this rule belongs to.
     */
    private int tenantId;

    /**
     * The currency code this rule applies to, if currency-specific.
     */
    private String currency;

    /**
     * The wallet ID this rule applies to, if wallet-specific.
     */
    private String walletId;

    /**
     * The address ID this rule applies to, if address-specific.
     */
    private String addressId;

    /**
     * The rule key/name that identifies the rule type (e.g., "max_transaction_amount").
     */
    private String ruleKey;

    /**
     * The rule value/setting (e.g., "1000000" for a limit).
     */
    private String ruleValue;

    /**
     * The group/category this rule belongs to for organizational purposes.
     */
    private String ruleGroup;

    /**
     * Human-readable description of what this rule does.
     */
    private String ruleDescription;

    /**
     * Validation pattern or constraint for the rule value.
     */
    private String ruleValidation;

    /**
     * The type of entity this rule is associated with (e.g., "wallet", "address").
     */
    private String entityType;

    /**
     * The ID of the entity this rule is associated with.
     */
    private String entityID;

    /**
     * Detailed information about the currency this rule applies to.
     */
    private Currency currencyInfo;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier for this business rule.
     *
     * @return the rule ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the unique identifier for this business rule.
     *
     * @param id the rule ID to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Returns the tenant ID this rule belongs to.
     *
     * @return the tenant ID
     */
    public int getTenantId() {
        return tenantId;
    }

    /**
     * Sets the tenant ID this rule belongs to.
     *
     * @param tenantId the tenant ID to set
     */
    public void setTenantId(int tenantId) {
        this.tenantId = tenantId;
    }

    /**
     * Returns the currency code this rule applies to.
     *
     * @return the currency code, or {@code null} if not currency-specific
     */
    public String getCurrency() {
        return currency;
    }

    /**
     * Sets the currency code this rule applies to.
     *
     * @param currency the currency code to set
     */
    public void setCurrency(String currency) {
        this.currency = currency;
    }

    /**
     * Returns the wallet ID this rule applies to.
     *
     * @return the wallet ID, or {@code null} if not wallet-specific
     */
    public String getWalletId() {
        return walletId;
    }

    /**
     * Sets the wallet ID this rule applies to.
     *
     * @param walletId the wallet ID to set
     */
    public void setWalletId(String walletId) {
        this.walletId = walletId;
    }

    /**
     * Returns the address ID this rule applies to.
     *
     * @return the address ID, or {@code null} if not address-specific
     */
    public String getAddressId() {
        return addressId;
    }

    /**
     * Sets the address ID this rule applies to.
     *
     * @param addressId the address ID to set
     */
    public void setAddressId(String addressId) {
        this.addressId = addressId;
    }

    /**
     * Returns the rule key/name that identifies the rule type.
     *
     * @return the rule key (e.g., "max_transaction_amount")
     */
    public String getRuleKey() {
        return ruleKey;
    }

    /**
     * Sets the rule key/name that identifies the rule type.
     *
     * @param ruleKey the rule key to set
     */
    public void setRuleKey(String ruleKey) {
        this.ruleKey = ruleKey;
    }

    /**
     * Returns the rule value/setting.
     *
     * @return the rule value (e.g., "1000000" for a limit)
     */
    public String getRuleValue() {
        return ruleValue;
    }

    /**
     * Sets the rule value/setting.
     *
     * @param ruleValue the rule value to set
     */
    public void setRuleValue(String ruleValue) {
        this.ruleValue = ruleValue;
    }

    /**
     * Returns the group/category this rule belongs to.
     *
     * @return the rule group name
     */
    public String getRuleGroup() {
        return ruleGroup;
    }

    /**
     * Sets the group/category this rule belongs to.
     *
     * @param ruleGroup the rule group name to set
     */
    public void setRuleGroup(String ruleGroup) {
        this.ruleGroup = ruleGroup;
    }

    /**
     * Returns the human-readable description of this rule.
     *
     * @return the rule description
     */
    public String getRuleDescription() {
        return ruleDescription;
    }

    /**
     * Sets the human-readable description of this rule.
     *
     * @param ruleDescription the rule description to set
     */
    public void setRuleDescription(String ruleDescription) {
        this.ruleDescription = ruleDescription;
    }

    /**
     * Returns the validation pattern or constraint for the rule value.
     *
     * @return the rule validation pattern
     */
    public String getRuleValidation() {
        return ruleValidation;
    }

    /**
     * Sets the validation pattern or constraint for the rule value.
     *
     * @param ruleValidation the rule validation pattern to set
     */
    public void setRuleValidation(String ruleValidation) {
        this.ruleValidation = ruleValidation;
    }

    /**
     * Returns the type of entity this rule is associated with.
     *
     * @return the entity type (e.g., "wallet", "address")
     */
    public String getEntityType() {
        return entityType;
    }

    /**
     * Sets the type of entity this rule is associated with.
     *
     * @param entityType the entity type to set
     */
    public void setEntityType(String entityType) {
        this.entityType = entityType;
    }

    /**
     * Returns the ID of the entity this rule is associated with.
     *
     * @return the entity ID
     */
    public String getEntityID() {
        return entityID;
    }

    /**
     * Sets the ID of the entity this rule is associated with.
     *
     * @param entityID the entity ID to set
     */
    public void setEntityID(String entityID) {
        this.entityID = entityID;
    }

    /**
     * Returns detailed information about the currency this rule applies to.
     *
     * @return the currency metadata, or {@code null} if not currency-specific
     */
    public Currency getCurrencyInfo() {
        return currencyInfo;
    }

    /**
     * Sets detailed information about the currency this rule applies to.
     *
     * @param currencyInfo the currency metadata to set
     */
    public void setCurrencyInfo(Currency currencyInfo) {
        this.currencyInfo = currencyInfo;
    }
}
