package com.taurushq.sdk.protect.client.model;

/**
 * A whitelisted-address row dropped from a list because it failed integrity
 * verification, and why.
 *
 * @see WhitelistedAddressListResult
 */
public final class ExcludedWhitelistedAddress {

    private final String id;
    private final String reason;

    /**
     * Constructs an exclusion record.
     *
     * @param id     the address ID, or null when the row carried no usable ID
     * @param reason why the row failed verification
     */
    public ExcludedWhitelistedAddress(final String id, final String reason) {
        this.id = id;
        this.reason = reason;
    }

    /**
     * Returns the address ID, or null when the row carried no usable ID.
     *
     * @return the address ID
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
