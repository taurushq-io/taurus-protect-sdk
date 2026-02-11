package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents an internal wallet source restriction for rule matching.
 * <p>
 * When a rule line specifies {@link RuleSourceType#RuleSourceInternalWallet},
 * this class provides the wallet path pattern to match against.
 *
 * @see RuleSource
 * @see RuleSourceType
 */
public class RuleSourceInternalWallet {

    /**
     * The HD wallet path pattern to match (e.g., "m/44'/60'/0'").
     */
    private String path;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the wallet path.
     *
     * @return the wallet path
     */
    public String getPath() {
        return path;
    }

    /**
     * Sets the wallet path.
     *
     * @param path the wallet path
     */
    public void setPath(String path) {
        this.path = path;
    }
}
