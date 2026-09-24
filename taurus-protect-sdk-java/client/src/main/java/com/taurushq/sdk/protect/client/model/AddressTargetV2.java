package com.taurushq.sdk.protect.client.model;

/**
 * The address an asset operation acts on: an internal address or a whitelisted one.
 */
public class AddressTargetV2 {

    private String addressId;
    private String whitelistedAddressId;

    /**
     * Gets the internal address id.
     *
     * @return the internal address id
     */
    public String getAddressId() {
        return addressId;
    }

    /**
     * Sets the internal address id.
     *
     * @param addressId the internal address id
     */
    public void setAddressId(final String addressId) {
        this.addressId = addressId;
    }

    /**
     * Gets the whitelisted address id.
     *
     * @return the whitelisted address id
     */
    public String getWhitelistedAddressId() {
        return whitelistedAddressId;
    }

    /**
     * Sets the whitelisted address id.
     *
     * @param whitelistedAddressId the whitelisted address id
     */
    public void setWhitelistedAddressId(final String whitelistedAddressId) {
        this.whitelistedAddressId = whitelistedAddressId;
    }
}
