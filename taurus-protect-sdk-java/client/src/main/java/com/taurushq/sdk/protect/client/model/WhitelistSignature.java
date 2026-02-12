package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.ArrayList;
import java.util.List;

/**
 * Represents a cryptographic signature for whitelist operations.
 * <p>
 * This class contains the signature data used to verify the authenticity
 * and authorization of whitelist modifications. It includes hashes of the
 * signed content and the actual user signature.
 *
 * @see WhitelistUserSignature
 * @see WhitelistedAddress
 */
public class WhitelistSignature {

    /**
     * List of cryptographic hashes of the signed content.
     */
    private final List<String> hashes = new ArrayList<>();

    /**
     * The user's cryptographic signature for the whitelist operation.
     */
    private WhitelistUserSignature signature;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the user signature for this whitelist operation.
     *
     * @return the user signature
     */
    public WhitelistUserSignature getSignature() {
        return signature;
    }

    /**
     * Sets the user signature for this whitelist operation.
     *
     * @param signature the user signature to set
     */
    public void setSignature(WhitelistUserSignature signature) {
        this.signature = signature;
    }

    /**
     * Returns the list of cryptographic hashes of the signed content.
     *
     * @return the list of content hashes
     */
    public List<String> getHashes() {
        return hashes;
    }
}
