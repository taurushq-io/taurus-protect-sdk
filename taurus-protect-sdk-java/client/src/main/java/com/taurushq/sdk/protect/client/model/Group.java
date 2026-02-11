package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents a user group in the Taurus Protect system.
 * <p>
 * Groups are used to organize users and define access permissions
 * and approval workflows.
 *
 * @see GroupService
 */
public class Group {

    private String id;
    private String tenantId;
    private String externalGroupId;
    private String name;
    private String email;
    private List<GroupUser> users;
    private OffsetDateTime creationDate;
    private OffsetDateTime updateDate;
    private String description;
    private Boolean enforcedInRules;

    /**
     * Gets the unique identifier of the group.
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
     * Gets the tenant ID this group belongs to.
     *
     * @return the tenant ID
     */
    public String getTenantId() {
        return tenantId;
    }

    /**
     * Sets the tenant ID.
     *
     * @param tenantId the tenant ID to set
     */
    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
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
     * Gets the group name.
     *
     * @return the group name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the group name.
     *
     * @param name the name to set
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Gets the group email address.
     *
     * @return the group email
     */
    public String getEmail() {
        return email;
    }

    /**
     * Sets the group email.
     *
     * @param email the email to set
     */
    public void setEmail(String email) {
        this.email = email;
    }

    /**
     * Gets the list of users in this group.
     *
     * @return the list of group users
     */
    public List<GroupUser> getUsers() {
        return users;
    }

    /**
     * Sets the list of users.
     *
     * @param users the users to set
     */
    public void setUsers(List<GroupUser> users) {
        this.users = users;
    }

    /**
     * Gets the creation date of the group.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the creation date.
     *
     * @param creationDate the creation date to set
     */
    public void setCreationDate(OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Gets the last update date of the group.
     *
     * @return the update date
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the update date.
     *
     * @param updateDate the update date to set
     */
    public void setUpdateDate(OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }

    /**
     * Gets the group description.
     *
     * @return the description
     */
    public String getDescription() {
        return description;
    }

    /**
     * Sets the description.
     *
     * @param description the description to set
     */
    public void setDescription(String description) {
        this.description = description;
    }

    /**
     * Returns whether this group is enforced in business rules.
     *
     * @return true if enforced in rules, false otherwise
     */
    public Boolean getEnforcedInRules() {
        return enforcedInRules;
    }

    /**
     * Sets whether this group is enforced in rules.
     *
     * @param enforcedInRules the enforcement flag
     */
    public void setEnforcedInRules(Boolean enforcedInRules) {
        this.enforcedInRules = enforcedInRules;
    }
}
