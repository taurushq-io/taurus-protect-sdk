package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;
import java.util.Map;

/**
 * Represents a change record in the Taurus Protect audit system.
 * <p>
 * Changes track modifications made to system entities (wallets, addresses, users, etc.)
 * providing a complete audit trail of who changed what, when, and why. Each change
 * captures the before and after state of modified fields.
 *
 * @see ChangeResult
 */
public class Change {

    /**
     * Unique identifier for this change record.
     */
    private String id;

    /**
     * ID of the tenant (organization) where this change occurred.
     */
    private int tenantId;

    /**
     * Internal user ID of the user who made the change.
     */
    private String creatorId;

    /**
     * External user ID (from SSO/IdP) of the user who made the change.
     */
    private String creatorExternalId;

    /**
     * The type of action performed (e.g., "create", "update", "delete").
     */
    private String action;

    /**
     * The type of entity that was changed (e.g., "wallet", "address", "user").
     */
    private String entity;

    /**
     * The numeric ID of the entity that was changed.
     */
    private String entityId;

    /**
     * The UUID of the entity that was changed (if applicable).
     */
    private String entityUUID;

    /**
     * Map of field names to their new values after the change.
     */
    private Map<String, String> changes;

    /**
     * Optional comment explaining the reason for the change.
     */
    private String comment;

    /**
     * Timestamp when the change was made.
     */
    private OffsetDateTime createdAt;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier for this change record.
     *
     * @return the change record ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the unique identifier for this change record.
     *
     * @param id the change record ID to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Returns the tenant ID where this change occurred.
     *
     * @return the tenant ID
     */
    public int getTenantId() {
        return tenantId;
    }

    /**
     * Sets the tenant ID where this change occurred.
     *
     * @param tenantId the tenant ID to set
     */
    public void setTenantId(int tenantId) {
        this.tenantId = tenantId;
    }

    /**
     * Returns the internal user ID of the user who made the change.
     *
     * @return the creator's user ID
     */
    public String getCreatorId() {
        return creatorId;
    }

    /**
     * Sets the internal user ID of the user who made the change.
     *
     * @param creatorId the creator's user ID to set
     */
    public void setCreatorId(String creatorId) {
        this.creatorId = creatorId;
    }

    /**
     * Returns the external user ID (from SSO/IdP) of the user who made the change.
     *
     * @return the creator's external user ID, or {@code null} if not set
     */
    public String getCreatorExternalId() {
        return creatorExternalId;
    }

    /**
     * Sets the external user ID (from SSO/IdP) of the user who made the change.
     *
     * @param creatorExternalId the creator's external user ID to set
     */
    public void setCreatorExternalId(String creatorExternalId) {
        this.creatorExternalId = creatorExternalId;
    }

    /**
     * Returns the type of action performed.
     *
     * @return the action (e.g., "create", "update", "delete")
     */
    public String getAction() {
        return action;
    }

    /**
     * Sets the type of action performed.
     *
     * @param action the action to set
     */
    public void setAction(String action) {
        this.action = action;
    }

    /**
     * Returns the type of entity that was changed.
     *
     * @return the entity type (e.g., "wallet", "address", "user")
     */
    public String getEntity() {
        return entity;
    }

    /**
     * Sets the type of entity that was changed.
     *
     * @param entity the entity type to set
     */
    public void setEntity(String entity) {
        this.entity = entity;
    }

    /**
     * Returns the numeric ID of the entity that was changed.
     *
     * @return the entity ID
     */
    public String getEntityId() {
        return entityId;
    }

    /**
     * Sets the numeric ID of the entity that was changed.
     *
     * @param entityId the entity ID to set
     */
    public void setEntityId(String entityId) {
        this.entityId = entityId;
    }

    /**
     * Returns the UUID of the entity that was changed.
     *
     * @return the entity UUID, or {@code null} if not applicable
     */
    public String getEntityUUID() {
        return entityUUID;
    }

    /**
     * Sets the UUID of the entity that was changed.
     *
     * @param entityUUID the entity UUID to set
     */
    public void setEntityUUID(String entityUUID) {
        this.entityUUID = entityUUID;
    }

    /**
     * Returns the map of field names to their new values after the change.
     *
     * @return the changes map
     */
    public Map<String, String> getChanges() {
        return changes;
    }

    /**
     * Sets the map of field names to their new values after the change.
     *
     * @param changes the changes map to set
     */
    public void setChanges(Map<String, String> changes) {
        this.changes = changes;
    }

    /**
     * Returns the optional comment explaining the reason for the change.
     *
     * @return the comment, or {@code null} if none provided
     */
    public String getComment() {
        return comment;
    }

    /**
     * Sets the optional comment explaining the reason for the change.
     *
     * @param comment the comment to set
     */
    public void setComment(String comment) {
        this.comment = comment;
    }

    /**
     * Returns the timestamp when the change was made.
     *
     * @return the creation timestamp
     */
    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    /**
     * Sets the timestamp when the change was made.
     *
     * @param createdAt the creation timestamp to set
     */
    public void setCreatedAt(OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }
}
