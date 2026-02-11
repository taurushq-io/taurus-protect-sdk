package com.taurushq.sdk.protect.client.model;

/**
 * Represents the entity type for a multi-factor signature.
 * <p>
 * Entity types identify what kind of entity is associated with the
 * multi-factor signature request.
 *
 * @see MultiFactorSignatureInfo
 */
public class MultiFactorSignatureEntityType {

    private String id;
    private String kind;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getKind() {
        return kind;
    }

    public void setKind(final String kind) {
        this.kind = kind;
    }
}
