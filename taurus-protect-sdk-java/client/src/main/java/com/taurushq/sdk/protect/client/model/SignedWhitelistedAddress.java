package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Contains the signed payload and cryptographic signatures for a whitelisted address.
 * <p>
 * The payload is a protobuf-encoded WhitelistedAddress that has been signed by
 * authorized approvers. This class is a data holder; verification is performed
 * by the WhitelistedAddressService.
 *
 * @see WhitelistSignature
 * @see SignedWhitelistedAddressEnvelope
 * @see com.taurushq.sdk.protect.client.service.WhitelistedAddressService
 */
public class SignedWhitelistedAddress {

    /**
     * List of cryptographic signatures from approvers who authorized this address.
     */
    private List<WhitelistSignature> signatures;

    /**
     * Protobuf-encoded WhitelistedAddress data.
     */
    private byte[] payload;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the list of cryptographic signatures from approvers.
     *
     * @return the list of signatures authorizing this address
     */
    public List<WhitelistSignature> getSignatures() {
        return signatures;
    }

    /**
     * Sets the list of cryptographic signatures from approvers.
     *
     * @param signatures the list of signatures to set
     */
    public void setSignatures(List<WhitelistSignature> signatures) {
        this.signatures = signatures;
    }

    /**
     * Returns the protobuf-encoded WhitelistedAddress data.
     *
     * @return the binary payload containing the whitelisted address details
     */
    public byte[] getPayload() {
        return payload;
    }

    /**
     * Sets the protobuf-encoded WhitelistedAddress data.
     *
     * @param payload the binary payload to set
     */
    public void setPayload(byte[] payload) {
        this.payload = payload;
    }
}
