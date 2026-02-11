package com.taurushq.sdk.protect.client.model;

/**
 * Represents a fiat currency provider in the Taurus Protect system.
 *
 * @see FiatService
 */
public class FiatProvider {

    private String provider;
    private String label;
    private String baseCurrencyValuation;

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

    public String getBaseCurrencyValuation() {
        return baseCurrencyValuation;
    }

    public void setBaseCurrencyValuation(final String baseCurrencyValuation) {
        this.baseCurrencyValuation = baseCurrencyValuation;
    }
}
