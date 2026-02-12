package com.taurushq.sdk.protect.client.model;

/**
 * Represents a user within a group in the Taurus Protect system.
 *
 * @see Group
 */
public class GroupUser {

    private String id;
    private String externalUserId;
    private Boolean enforcedInRules;

    /**
     * Gets the user's unique identifier.
     *
     * @return the user ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the user ID.
     *
     * @param id the ID to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the external user ID (e.g., from an identity provider).
     *
     * @return the external user ID
     */
    public String getExternalUserId() {
        return externalUserId;
    }

    /**
     * Sets the external user ID.
     *
     * @param externalUserId the external user ID to set
     */
    public void setExternalUserId(String externalUserId) {
        this.externalUserId = externalUserId;
    }

    /**
     * Returns whether this user is enforced in business rules.
     *
     * @return true if enforced in rules, false otherwise
     */
    public Boolean getEnforcedInRules() {
        return enforcedInRules;
    }

    /**
     * Sets whether this user is enforced in rules.
     *
     * @param enforcedInRules the enforcement flag
     */
    public void setEnforcedInRules(Boolean enforcedInRules) {
        this.enforcedInRules = enforcedInRules;
    }
}
