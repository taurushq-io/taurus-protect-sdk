package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * External address source restriction.
 * <p>
 * Populated when a {@link RuleSource} has type
 * {@link RuleSourceType#RuleSourceExternalAddress}.
 *
 * @see RuleSource
 */
public class RuleSourceExternalAddress {

    private String address;
    private String memo;

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
     * Gets the memo.
     *
     * @return the memo
     */
    public String getMemo() {
        return memo;
    }

    /**
     * Sets the memo.
     *
     * @param memo the memo
     */
    public void setMemo(String memo) {
        this.memo = memo;
    }
}
