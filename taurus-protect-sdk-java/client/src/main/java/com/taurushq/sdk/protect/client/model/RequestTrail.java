package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;

/**
 * Represents an audit trail entry for a transaction request.
 * <p>
 * Each trail entry records a specific action taken on a request, including
 * who performed the action, when it occurred, and any comments provided.
 * This provides a complete history of a request's lifecycle for audit
 * and compliance purposes.
 * <p>
 * Common actions include:
 * <ul>
 *   <li>Request creation</li>
 *   <li>Approval or rejection</li>
 *   <li>Status changes</li>
 *   <li>Comments or notes added</li>
 * </ul>
 *
 * @see Request
 * @see RequestStatus
 */
public class RequestTrail {

    /**
     * Unique identifier for this trail entry.
     */
    private long id;

    /**
     * Internal user ID of the person who performed the action.
     */
    private long userId;

    /**
     * External user ID from an identity provider (e.g., SSO system).
     */
    private String externalUserId;

    /**
     * The action that was performed (e.g., "approve", "reject", "create").
     */
    private String action;

    /**
     * Optional comment provided by the user when performing the action.
     */
    private String comment;

    /**
     * Timestamp when the action was performed.
     */
    private OffsetDateTime date;

    /**
     * The request status after this action was performed.
     */
    private String requestStatus;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier for this trail entry.
     *
     * @return the trail entry ID
     */
    public long getId() {
        return id;
    }

    /**
     * Sets the unique identifier for this trail entry.
     *
     * @param id the trail entry ID to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Returns the internal user ID of the person who performed the action.
     *
     * @return the internal user ID
     */
    public long getUserId() {
        return userId;
    }

    /**
     * Sets the internal user ID of the person who performed the action.
     *
     * @param userId the internal user ID to set
     */
    public void setUserId(long userId) {
        this.userId = userId;
    }

    /**
     * Returns the external user ID from an identity provider.
     *
     * @return the external user ID, or {@code null} if not available
     */
    public String getExternalUserId() {
        return externalUserId;
    }

    /**
     * Sets the external user ID from an identity provider.
     *
     * @param externalUserId the external user ID to set
     */
    public void setExternalUserId(String externalUserId) {
        this.externalUserId = externalUserId;
    }

    /**
     * Returns the action that was performed.
     * <p>
     * Common actions include "create", "approve", "reject", "cancel".
     *
     * @return the action name
     */
    public String getAction() {
        return action;
    }

    /**
     * Sets the action that was performed.
     *
     * @param action the action name to set
     */
    public void setAction(String action) {
        this.action = action;
    }

    /**
     * Returns the optional comment provided with the action.
     *
     * @return the comment, or {@code null} if none was provided
     */
    public String getComment() {
        return comment;
    }

    /**
     * Sets the optional comment for the action.
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

    /**
     * Returns the request status after this action was performed.
     *
     * @return the resulting request status
     * @see RequestStatus
     */
    public String getRequestStatus() {
        return requestStatus;
    }

    /**
     * Sets the request status after this action was performed.
     *
     * @param requestStatus the resulting request status to set
     */
    public void setRequestStatus(String requestStatus) {
        this.requestStatus = requestStatus;
    }
}
