package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Contains a signed whitelisted asset (contract address) with cryptographic signatures.
 * <p>
 * This class holds the binary payload of a whitelisted asset along with the signatures
 * from approvers who authorized the whitelisting. The payload contains the serialized
 * protobuf representation of the whitelisted asset data.
 * <p>
 * The signatures can be verified against the payload to ensure the whitelisting
 * was properly authorized according to governance rules.
 *
 * @see SignedWhitelistedAssetEnvelope
 * @see WhitelistSignature
 * @see WhitelistedAsset
 */
public class SignedWhitelistedAsset {

    /**
     * List of cryptographic signatures from approvers who authorized this whitelisting.
     * Each signature contains the user ID, timestamp, and cryptographic proof.
     */
    private List<WhitelistSignature> signatures;

    /**
     * Binary payload containing the serialized whitelisted asset data.
     * This is the protobuf-encoded representation of the asset being whitelisted.
     */
    private byte[] payload;

    /**
     * Gets the list of signatures authorizing this whitelisting.
     *
     * @return the list of signatures, may be empty if not yet approved
     */
    public List<WhitelistSignature> getSignatures() {
        return signatures;
    }

    /**
     * Sets the list of signatures.
     *
     * @param signatures the signatures to set
     */
    public void setSignatures(List<WhitelistSignature> signatures) {
        this.signatures = signatures;
    }

    /**
     * Gets the binary payload containing the serialized asset data.
     *
     * @return the payload bytes
     */
    public byte[] getPayload() {
        return payload;
    }

    /**
     * Sets the binary payload.
     *
     * @param payload the payload bytes to set
     */
    public void setPayload(byte[] payload) {
        this.payload = payload;
    }

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }
}
