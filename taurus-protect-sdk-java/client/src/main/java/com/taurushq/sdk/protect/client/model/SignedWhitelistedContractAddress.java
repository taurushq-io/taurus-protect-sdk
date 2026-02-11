package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents a signed whitelisted contract address with its cryptographic signatures.
 * <p>
 * This class contains the raw payload (protobuf-encoded contract address data) along
 * with the signatures from users who have signed the whitelist entry. The payload
 * can be decoded using the {@link WhitelistedContractAddressMapper} to extract
 * the underlying contract address details.
 *
 * @see SignedWhitelistedContractAddressEnvelope
 * @see WhitelistedContractAddress
 */
public class SignedWhitelistedContractAddress {

    /**
     * List of signatures from users who have signed this whitelist entry.
     */
    private List<WhitelistSignature> signatures;

    /**
     * Raw protobuf-encoded payload containing the contract address data.
     */
    private byte[] payload;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the list of signatures.
     *
     * @return the signatures
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
     * Gets the raw protobuf-encoded payload.
     *
     * @return the payload bytes
     */
    public byte[] getPayload() {
        return payload;
    }

    /**
     * Sets the raw protobuf-encoded payload.
     *
     * @param payload the payload bytes to set
     */
    public void setPayload(byte[] payload) {
        this.payload = payload;
    }
}
