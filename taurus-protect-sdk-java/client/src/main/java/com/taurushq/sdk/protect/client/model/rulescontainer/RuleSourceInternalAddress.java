package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Internal address source restriction.
 * <p>
 * Populated when a {@link RuleSource} has type
 * {@link RuleSourceType#RuleSourceInternalAddress}.
 *
 * @see RuleSource
 */
public class RuleSourceInternalAddress {

    private String address;
    private String path;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the address.
     *
     * @return the address
     */
    public String getAddress() {
        return address;
    }

    /**
     * Sets the address.
     *
     * @param address the address
     */
    public void setAddress(String address) {
        this.address = address;
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
