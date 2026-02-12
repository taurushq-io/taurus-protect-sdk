package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents an audit trail entry in the Taurus Protect system.
 * <p>
 * Audit trails record all significant actions performed in the system,
 * providing a complete history for compliance and security purposes.
 *
 * @see AuditService
 */
public class AuditTrail {

    private String id;
    private String entity;
    private String action;
    private String subAction;
    private String details;
    private OffsetDateTime creationDate;

    /**
     * Gets the unique identifier.
     *
     * @return the id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the unique identifier.
     *
     * @param id the id to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the entity type that was affected.
     *
     * @return the entity type (e.g., "address", "wallet", "request")
     */
    public String getEntity() {
        return entity;
    }

    /**
     * Sets the entity type.
     *
     * @param entity the entity to set
     */
    public void setEntity(String entity) {
        this.entity = entity;
    }

    /**
     * Gets the action that was performed.
     *
     * @return the action (e.g., "create", "approve", "reject")
     */
    public String getAction() {
        return action;
    }

    /**
     * Sets the action.
     *
     * @param action the action to set
     */
    public void setAction(String action) {
        this.action = action;
    }

    /**
     * Gets the sub-action (more specific action detail).
     *
     * @return the sub-action
     */
    public String getSubAction() {
        return subAction;
    }

    /**
     * Sets the sub-action.
     *
     * @param subAction the sub-action to set
     */
    public void setSubAction(String subAction) {
        this.subAction = subAction;
    }

    /**
     * Gets the detailed description of the action.
     *
     * @return the details
     */
    public String getDetails() {
        return details;
    }

    /**
     * Sets the details.
     *
     * @param details the details to set
     */
    public void setDetails(String details) {
        this.details = details;
    }

    /**
     * Gets the creation date of the audit trail entry.
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
}
