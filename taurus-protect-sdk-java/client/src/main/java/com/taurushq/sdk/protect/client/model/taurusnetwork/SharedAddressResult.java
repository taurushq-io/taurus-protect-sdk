package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.CursorPagedResult;

import java.util.List;

/**
 * Result containing a list of shared addresses with pagination.
 */
public class SharedAddressResult extends CursorPagedResult {

    private List<SharedAddress> sharedAddresses;

    /**
     * Gets the shared addresses of this page.
     *
     * @return the sharedAddresses
     */
    public List<SharedAddress> getSharedAddresses() {
        return sharedAddresses;
    }

    /**
     * Sets the shared addresses of this page.
     *
     * @param sharedAddresses the sharedAddresses
     */
    public void setSharedAddresses(final List<SharedAddress> sharedAddresses) {
        this.sharedAddresses = sharedAddresses;
    }
}
