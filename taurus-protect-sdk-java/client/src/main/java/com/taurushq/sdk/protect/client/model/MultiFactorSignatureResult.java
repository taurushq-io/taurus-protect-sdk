package com.taurushq.sdk.protect.client.model;

/**
 * Represents the result of creating multi-factor signatures.
 *
 * @see MultiFactorSignatureService
 */
public class MultiFactorSignatureResult {

    private String id;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }
}
