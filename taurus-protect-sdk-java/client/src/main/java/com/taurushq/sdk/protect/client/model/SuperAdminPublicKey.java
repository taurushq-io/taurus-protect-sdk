package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a super administrator's public key in the Taurus Protect system.
 * <p>
 * Super administrators have elevated privileges and are required to sign
 * governance rule changes. Their public keys are used to verify these signatures
 * and ensure only authorized administrators can modify system policies.
 *
 * @see GovernanceRules
 * @see com.taurushq.sdk.protect.client.helper.SignatureVerifier
 */
public class SuperAdminPublicKey {

    /**
     * The unique identifier of the super administrator user.
     */
    private String userId;

    /**
     * The super administrator's public key (base64 or PEM encoded).
     */
    private String publicKey;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier of the super administrator user.
     *
     * @return the user ID
     */
    public String getUserId() {
        return userId;
    }

    /**
     * Sets the unique identifier of the super administrator user.
     *
     * @param userId the user ID to set
     */
    public void setUserId(String userId) {
        this.userId = userId;
    }

    /**
     * Returns the super administrator's public key.
     *
     * @return the public key (base64 or PEM encoded)
     */
    public String getPublicKey() {
        return publicKey;
    }

    /**
     * Sets the super administrator's public key.
     *
     * @param publicKey the public key to set (base64 or PEM encoded)
     */
    public void setPublicKey(String publicKey) {
        this.publicKey = publicKey;
    }
}
