package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Exchange source restriction.
 * <p>
 * Populated when a {@link RuleSource} has type
 * {@link RuleSourceType#RuleSourceExchange}.
 *
 * @see RuleSource
 */
public class RuleSourceExchange {

    private String label;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the exchange label.
     *
     * @return the exchange label
     */
    public String getLabel() {
        return label;
    }

    /**
     * Sets the exchange label.
     *
     * @param label the exchange label
     */
    public void setLabel(String label) {
        this.label = label;
    }
}
