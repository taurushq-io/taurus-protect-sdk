package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.Map;

/**
 * Represents a request to create a configuration change in the Taurus Protect system.
 * <p>
 * Changes follow an approval workflow: one admin proposes the change, then another
 * admin approves it before it takes effect.
 *
 * @see Change
 */
public class CreateChangeRequest {

    /**
     * The type of action to perform (e.g., "update", "create", "delete").
     */
    private String action;

    /**
     * The type of entity to change (e.g., "businessrule", "user", "group").
     */
    private String entity;

    /**
     * The ID of the entity to change.
     */
    private String entityId;

    /**
     * Map of field names to their new values.
     */
    private Map<String, String> changes;

    /**
     * Optional description of the change request.
     */
    private String comment;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the type of action to perform.
     *
     * @return the action (e.g., "update", "create", "delete")
     */
    public String getAction() {
        return action;
    }

    /**
     * Sets the type of action to perform.
     *
     * @param action the action to set
     */
    public void setAction(String action) {
        this.action = action;
    }

    /**
     * Returns the type of entity to change.
     *
     * @return the entity type (e.g., "businessrule", "user", "group")
     */
    public String getEntity() {
        return entity;
    }

    /**
     * Sets the type of entity to change.
     *
     * @param entity the entity type to set
     */
    public void setEntity(String entity) {
        this.entity = entity;
    }

    /**
     * Returns the ID of the entity to change.
     *
     * @return the entity ID
     */
    public String getEntityId() {
        return entityId;
    }

    /**
     * Sets the ID of the entity to change.
     *
     * @param entityId the entity ID to set
     */
    public void setEntityId(String entityId) {
        this.entityId = entityId;
    }

    /**
     * Returns the map of field names to their new values.
     *
     * @return the changes map
     */
    public Map<String, String> getChanges() {
        return changes;
    }

    /**
     * Sets the map of field names to their new values.
     *
     * @param changes the changes map to set
     */
    public void setChanges(Map<String, String> changes) {
        this.changes = changes;
    }

    /**
     * Returns the optional description of the change request.
     *
     * @return the comment, or {@code null} if none provided
     */
    public String getComment() {
        return comment;
    }

    /**
     * Sets the optional description of the change request.
     *
     * @param comment the comment to set
     */
    public void setComment(String comment) {
        this.comment = comment;
    }
}
