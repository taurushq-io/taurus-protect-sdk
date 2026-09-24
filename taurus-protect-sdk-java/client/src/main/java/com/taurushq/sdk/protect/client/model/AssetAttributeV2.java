package com.taurushq.sdk.protect.client.model;

/**
 * A key/value attribute of a v2 asset.
 */
public class AssetAttributeV2 {

    private String key;
    private String value;

    /**
     * Gets the key.
     *
     * @return the key
     */
    public String getKey() {
        return key;
    }

    /**
     * Sets the key.
     *
     * @param key the key
     */
    public void setKey(final String key) {
        this.key = key;
    }

    /**
     * Gets the value.
     *
     * @return the value
     */
    public String getValue() {
        return value;
    }

    /**
     * Sets the value.
     *
     * @param value the value
     */
    public void setValue(final String value) {
        this.value = value;
    }
}
