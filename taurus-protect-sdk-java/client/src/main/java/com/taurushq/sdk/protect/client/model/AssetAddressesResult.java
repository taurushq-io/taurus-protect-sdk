package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * One page of addresses holding an asset, each one's HSM signature verified from a cursor list. Continue with {@code getPage()}.
 *
 * @see com.taurushq.sdk.protect.client.service.AssetService
 */
public class AssetAddressesResult extends CursorPagedResult {

    private List<Address> addresses;

    /**
     * Gets the addresses holding an asset of this page.
     *
     * @return the addresses
     */
    public List<Address> getAddresses() {
        return addresses;
    }

    /**
     * Sets the addresses holding an asset of this page.
     *
     * @param addresses the addresses
     */
    public void setAddresses(final List<Address> addresses) {
        this.addresses = addresses;
    }
}
