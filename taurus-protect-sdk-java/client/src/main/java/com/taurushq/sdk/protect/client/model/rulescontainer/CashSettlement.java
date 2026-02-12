package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents cash settlement configuration for fiat-related operations.
 * <p>
 * Cash settlement rules define how fiat currency operations are handled,
 * including the settlement provider and request types.
 *
 * @see TransactionRuleDetails
 */
public class CashSettlement {

    /**
     * The cash settlement provider (e.g., banking partner identifier).
     */
    private String provider;

    /**
     * The type of cash settlement request.
     */
    private String requestType;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the provider.
     *
     * @return the provider
     */
    public String getProvider() {
        return provider;
    }

    /**
     * Sets the provider.
     *
     * @param provider the provider
     */
    public void setProvider(String provider) {
        this.provider = provider;
    }

    /**
     * Gets the request type.
     *
     * @return the request type
     */
    public String getRequestType() {
        return requestType;
    }

    /**
     * Sets the request type.
     *
     * @param requestType the request type
     */
    public void setRequestType(String requestType) {
        this.requestType = requestType;
    }
}
