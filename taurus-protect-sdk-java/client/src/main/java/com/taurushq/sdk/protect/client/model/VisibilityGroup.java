package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents a visibility group in the Taurus Protect system.
 * <p>
 * Visibility groups are used to control which users can see specific
 * wallets, addresses, and other entities. Users can only see entities
 * that belong to their visibility groups.
 *
 * @see VisibilityGroupUser
 * @see VisibilityGroupService
 */
public class VisibilityGroup {

    private String id;
    private String tenantId;
    private String name;
    private String description;
    private List<VisibilityGroupUser> users;
    private OffsetDateTime creationDate;
    private OffsetDateTime updateDate;
    private String userCount;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(final String tenantId) {
        this.tenantId = tenantId;
    }

    public String getName() {
        return name;
    }

    public void setName(final String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(final String description) {
        this.description = description;
    }

    public List<VisibilityGroupUser> getUsers() {
        return users;
    }

    public void setUsers(final List<VisibilityGroupUser> users) {
        this.users = users;
    }

    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    public void setCreationDate(final OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    public void setUpdateDate(final OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }

    public String getUserCount() {
        return userCount;
    }

    public void setUserCount(final String userCount) {
        this.userCount = userCount;
    }
}
