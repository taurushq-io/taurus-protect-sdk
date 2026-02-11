package com.taurushq.sdk.protect.client.model;

import com.google.gson.annotations.SerializedName;

/**
 * Represents an internal address linked to a whitelisted address.
 * <p>
 * Internal addresses are addresses managed within Taurus Protect that are
 * linked to external whitelisted addresses. This allows tracking of which
 * internal addresses are authorized to interact with specific external addresses.
 *
 * @see WhitelistedAddress
 */
public class InternalAddress {

    /**
     * Unique identifier for the internal address.
     */
    @SerializedName("id")
    private String id;

    /**
     * The blockchain address string.
     */
    @SerializedName("address")
    private String address;

    /**
     * Human-readable label for the address.
     */
    @SerializedName("label")
    private String label;

    /**
     * Returns the unique identifier for this internal address.
     *
     * @return the address ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the unique identifier for this internal address from a long value.
     *
     * @param id the address ID as a long
     */
    public void setId(long id) {
        this.id = String.valueOf(id);
    }

    /**
     * Sets the unique identifier for this internal address.
     *
     * @param id the address ID
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Returns the blockchain address string.
     *
     * @return the blockchain address
     */
    public String getAddress() {
        return address;
    }

    /**
     * Sets the blockchain address string.
     *
     * @param address the blockchain address to set
     */
    public void setAddress(String address) {
        this.address = address;
    }

    /**
     * Returns the human-readable label for this address.
     *
     * @return the address label
     */
    public String getLabel() {
        return label;
    }

    /**
     * Sets the human-readable label for this address.
     *
     * @param label the address label to set
     */
    public void setLabel(String label) {
        this.label = label;
    }
}
