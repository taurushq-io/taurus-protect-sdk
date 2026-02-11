package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a user's cryptographic signature for a whitelist operation.
 * <p>
 * This class contains the signature data provided by a user when approving
 * or creating a whitelist entry. The signature is used to verify the user's
 * authorization and to maintain an audit trail.
 *
 * @see WhitelistSignature
 * @see WhitelistedAddress
 */
public class WhitelistUserSignature {

    /**
     * The ID of the user who created this signature.
     */
    private String userId;

    /**
     * The raw cryptographic signature bytes.
     */
    private byte[] signature;

    /**
     * Optional comment provided by the user when signing.
     */
    private String comment;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the ID of the user who created this signature.
     *
     * @return the user ID
     */
    public String getUserId() {
        return userId;
    }

    /**
     * Sets the ID of the user who created this signature.
     *
     * @param userId the user ID to set
     */
    public void setUserId(String userId) {
        this.userId = userId;
    }

    /**
     * Returns the raw cryptographic signature bytes.
     *
     * @return the signature bytes
     */
    public byte[] getSignature() {
        return signature;
    }

    /**
     * Sets the raw cryptographic signature bytes.
     *
     * @param signature the signature bytes to set
     */
    public void setSignature(byte[] signature) {
        this.signature = signature;
    }

    /**
     * Returns the optional comment provided when signing.
     *
     * @return the comment, or {@code null} if none was provided
     */
    public String getComment() {
        return comment;
    }

    /**
     * Sets the optional comment for the signature.
     *
     * @param comment the comment to set
     */
    public void setComment(String comment) {
        this.comment = comment;
    }
}
