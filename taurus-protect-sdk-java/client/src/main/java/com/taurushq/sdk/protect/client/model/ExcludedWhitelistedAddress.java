package com.taurushq.sdk.protect.client.model;

/**
 * A row dropped from a list because it failed integrity verification, and why: a
 * whitelisted address, or a v2 asset holder that its verified counterpart did not confirm.
 *
 * @see WhitelistedAddressListResult
 * @see AssetAddressV2Result
 */
public final class ExcludedWhitelistedAddress {

    private final String id;
    private final String reason;

    /**
     * Constructs an exclusion record.
     *
     * @param id     the row's ID, or what identifies it when it carried no usable ID
     * @param reason why the row failed verification
     */
    public ExcludedWhitelistedAddress(final String id, final String reason) {
        this.id = id;
        this.reason = reason;
    }

    /**
     * Returns the row's ID. A whitelisted-address row without a usable ID gives null; an
     * asset holder without one gives its address.
     *
     * @return the ID
     */
    public String getId() {
        return id;
    }

    /**
     * Returns why the row failed verification.
     *
     * @return the failure reason
     */
    public String getReason() {
        return reason;
    }

    @Override
    public String toString() {
        return "ExcludedWhitelistedAddress{id=" + id + ", reason=" + reason + "}";
    }
}
