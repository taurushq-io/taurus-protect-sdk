package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of prices from the prices query, each one's signature verified. Continue with
 * {@code getPage()}.
 *
 * @see com.taurushq.sdk.protect.client.service.PriceService
 */
public class PriceResult extends CursorPagedResult {

    private List<Price> prices;
    private String baseCurrency;

    /**
     * Gets the prices of this page.
     *
     * @return the prices
     */
    public List<Price> getPrices() {
        return prices;
    }

    /**
     * Sets the prices of this page.
     *
     * @param prices the prices
     */
    public void setPrices(final List<Price> prices) {
        this.prices = prices;
    }

    /**
     * Gets the tenant's base currency, as reported with the prices.
     *
     * @return the base currency, may be null
     */
    public String getBaseCurrency() {
        return baseCurrency;
    }

    /**
     * Sets the base currency.
     *
     * @param baseCurrency the base currency
     */
    public void setBaseCurrency(final String baseCurrency) {
        this.baseCurrency = baseCurrency;
    }
}
