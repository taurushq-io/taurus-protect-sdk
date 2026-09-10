package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * A PRICEUPDATER signature over the canonical price form.
 *
 * @see Price
 */
public class PriceSignature {

    /**
     * The signing user's ID.
     */
    private String userId;

    /**
     * Base64-encoded ECDSA signature.
     */
    private String signature;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the signing user's ID.
     *
     * @return the user ID
     */
    public String getUserId() {
        return userId;
    }

    /**
     * Sets the signing user's ID.
     *
     * @param userId the user ID to set
     */
    public void setUserId(String userId) {
        this.userId = userId;
    }

    /**
     * Returns the base64-encoded signature.
     *
     * @return the signature
     */
    public String getSignature() {
        return signature;
    }

    /**
     * Sets the base64-encoded signature.
     *
     * @param signature the signature to set
     */
    public void setSignature(String signature) {
        this.signature = signature;
    }
}
