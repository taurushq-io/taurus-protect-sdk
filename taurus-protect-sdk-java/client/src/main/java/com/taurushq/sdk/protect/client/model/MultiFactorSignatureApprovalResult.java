package com.taurushq.sdk.protect.client.model;

/**
 * Represents the result of approving a multi-factor signature.
 *
 * @see MultiFactorSignatureService
 */
public class MultiFactorSignatureApprovalResult {

    private boolean approved;
    private String message;

    public boolean isApproved() {
        return approved;
    }

    public void setApproved(final boolean approved) {
        this.approved = approved;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(final String message) {
        this.message = message;
    }
}
