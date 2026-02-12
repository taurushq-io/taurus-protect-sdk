package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a user's cryptographic signature on governance rules.
 * <p>
 * When governance rules are updated, super administrators must sign the changes
 * to prove their authorization. This class captures both the user who signed
 * and their cryptographic signature for verification.
 *
 * @see GovernanceRules
 * @see com.taurushq.sdk.protect.client.helper.SignatureVerifier
 */
public class RuleUserSignature {

    /**
     * The unique identifier of the user who provided this signature.
     */
    private String userId;

    /**
     * The cryptographic signature (typically base64-encoded) over the rules container.
     */
    private String signature;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier of the user who provided this signature.
     *
     * @return the user ID
     */
    public String getUserId() {
        return userId;
    }

    /**
     * Sets the unique identifier of the user who provided this signature.
     *
     * @param userId the user ID to set
     */
    public void setUserId(String userId) {
        this.userId = userId;
    }

    /**
     * Returns the cryptographic signature over the rules container.
     *
     * @return the base64-encoded signature
     */
    public String getSignature() {
        return signature;
    }

    /**
     * Sets the cryptographic signature over the rules container.
     *
     * @param signature the base64-encoded signature to set
     */
    public void setSignature(String signature) {
        this.signature = signature;
    }
}
