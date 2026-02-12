package com.taurushq.sdk.protect.client.model;

/**
 * Represents a network fee entry for a blockchain.
 * <p>
 * Fees are represented as key-value pairs where the key is typically
 * the blockchain/currency identifier and the value is the fee amount.
 *
 * @see FeeService
 */
public class Fee {

    private String key;
    private String value;

    /**
     * Gets the fee key (blockchain/currency identifier).
     *
     * @return the key
     */
    public String getKey() {
        return key;
    }

    /**
     * Sets the fee key.
     *
     * @param key the key to set
     */
    public void setKey(String key) {
        this.key = key;
    }

    /**
     * Gets the fee value (amount).
     *
     * @return the value
     */
    public String getValue() {
        return value;
    }

    /**
     * Sets the fee value.
     *
     * @param value the value to set
     */
    public void setValue(String value) {
        this.value = value;
    }
}
