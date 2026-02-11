package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;

import java.util.List;

/**
 * Result containing a list of shared addresses with pagination.
 */
public class SharedAddressResult {

    private List<SharedAddress> sharedAddresses;
    private ApiResponseCursor cursor;

    public List<SharedAddress> getSharedAddresses() {
        return sharedAddresses;
    }

    public void setSharedAddresses(final List<SharedAddress> sharedAddresses) {
        this.sharedAddresses = sharedAddresses;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
