package com.taurushq.sdk.protect.client.model;

/**
 * Represents a user's membership in a group in the Taurus Protect system.
 *
 * @see User
 */
public class UserGroup {

    private String id;
    private String externalGroupId;
    private Boolean enforcedInRules;

    /**
     * Gets the group's unique identifier.
     *
     * @return the group ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the group ID.
     *
     * @param id the ID to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the external group ID (e.g., from an identity provider).
     *
     * @return the external group ID
     */
    public String getExternalGroupId() {
        return externalGroupId;
    }

    /**
     * Sets the external group ID.
     *
     * @param externalGroupId the external group ID to set
     */
    public void setExternalGroupId(String externalGroupId) {
        this.externalGroupId = externalGroupId;
    }

    /**
     * Returns whether this group is enforced in governance rules.
     *
     * @return the flag, or {@code null} when the endpoint that returned the user does not compute it
     */
    public Boolean getEnforcedInRules() {
        return enforcedInRules;
    }

    /**
     * Sets whether this group is enforced in governance rules.
     *
     * @param enforcedInRules the enforcement flag
     */
    public void setEnforcedInRules(Boolean enforcedInRules) {
        this.enforcedInRules = enforcedInRules;
    }
}
