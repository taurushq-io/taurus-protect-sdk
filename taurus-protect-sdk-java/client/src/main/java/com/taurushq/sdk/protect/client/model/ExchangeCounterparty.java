package com.taurushq.sdk.protect.client.model;

/**
 * Represents an exchange counterparty summary.
 * <p>
 * A counterparty represents a grouped view of exchange holdings by exchange name,
 * with the total valuation across all currencies.
 *
 * @see ExchangeService
 */
public class ExchangeCounterparty {

    private String name;
    private String baseCurrencyValuation;

    /**
     * Gets the exchange name.
     *
     * @return the name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the exchange name.
     *
     * @param name the name to set
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Gets the total valuation in the base currency (CHF, EUR, USD, etc.).
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
}
