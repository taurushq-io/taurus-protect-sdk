package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents information about a multi-factor signature request.
 * <p>
 * Multi-factor signatures are used for operations that require multiple
 * approvals, such as critical governance changes or high-value transactions.
 *
 * @see MultiFactorSignatureService
 */
public class MultiFactorSignatureInfo {

    private String id;
    private List<String> payloadToSign;
    private MultiFactorSignatureEntityType entityType;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public List<String> getPayloadToSign() {
        return payloadToSign;
    }

    public void setPayloadToSign(final List<String> payloadToSign) {
        this.payloadToSign = payloadToSign;
    }

    public MultiFactorSignatureEntityType getEntityType() {
        return entityType;
    }

    public void setEntityType(final MultiFactorSignatureEntityType entityType) {
        this.entityType = entityType;
    }
}
