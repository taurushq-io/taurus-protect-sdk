package com.taurushq.sdk.protect.client.model;

/**
 * Represents a user assigned to a visibility group.
 *
 * @see VisibilityGroup
 */
public class VisibilityGroupUser {

    private String id;
    private String externalUserId;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getExternalUserId() {
        return externalUserId;
    }

    public void setExternalUserId(final String externalUserId) {
        this.externalUserId = externalUserId;
    }
}
