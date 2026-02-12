package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;

/**
 * Represents an audit trail entry for governance rules changes.
 * <p>
 * Every modification to governance rules creates an audit trail entry that
 * records who made the change, what action was performed, and when it occurred.
 * This provides a complete history of governance changes for compliance and auditing.
 *
 * @see GovernanceRules
 */
public class GovernanceRulesTrail {

    /**
     * Unique identifier for this audit trail entry.
     */
    private String id;

    /**
     * The internal user ID of the user who performed the action.
     */
    private String userId;

    /**
     * The external user ID (from SSO/IdP) of the user who performed the action.
     */
    private String externalUserId;

    /**
     * The action that was performed (e.g., "create", "update", "approve").
     */
    private String action;

    /**
     * Optional comment provided by the user explaining the change.
     */
    private String comment;

    /**
     * Timestamp when the action was performed.
     */
    private OffsetDateTime date;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier for this audit trail entry.
     *
     * @return the trail entry ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the unique identifier for this audit trail entry.
     *
     * @param id the trail entry ID to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Returns the internal user ID of the user who performed the action.
     *
     * @return the user ID
     */
    public String getUserId() {
        return userId;
    }

    /**
     * Sets the internal user ID of the user who performed the action.
     *
     * @param userId the user ID to set
     */
    public void setUserId(String userId) {
        this.userId = userId;
    }

    /**
     * Returns the external user ID (from SSO/IdP) of the user who performed the action.
     *
     * @return the external user ID, or {@code null} if not set
     */
    public String getExternalUserId() {
        return externalUserId;
    }

    /**
     * Sets the external user ID (from SSO/IdP) of the user who performed the action.
     *
     * @param externalUserId the external user ID to set
     */
    public void setExternalUserId(String externalUserId) {
        this.externalUserId = externalUserId;
    }

    /**
     * Returns the action that was performed.
     *
     * @return the action (e.g., "create", "update", "approve")
     */
    public String getAction() {
        return action;
    }

    /**
     * Sets the action that was performed.
     *
     * @param action the action to set
     */
    public void setAction(String action) {
        this.action = action;
    }

    /**
     * Returns the optional comment provided by the user explaining the change.
     *
     * @return the comment, or {@code null} if none provided
     */
    public String getComment() {
        return comment;
    }

    /**
     * Sets the optional comment explaining the change.
     *
     * @param comment the comment to set
     */
    public void setComment(String comment) {
        this.comment = comment;
    }

    /**
     * Returns the timestamp when the action was performed.
     *
     * @return the action timestamp
     */
    public OffsetDateTime getDate() {
        return date;
    }

    /**
     * Sets the timestamp when the action was performed.
     *
     * @param date the action timestamp to set
     */
    public void setDate(OffsetDateTime date) {
        this.date = date;
    }
}
