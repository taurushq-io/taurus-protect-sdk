package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents an audit trail entry for an action.
 */
public class ActionTrail {

    private String id;
    private String action;
    private String comment;
    private OffsetDateTime date;
    private String actionStatus;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getAction() {
        return action;
    }

    public void setAction(final String action) {
        this.action = action;
    }

    public String getComment() {
        return comment;
    }

    public void setComment(final String comment) {
        this.comment = comment;
    }

    public OffsetDateTime getDate() {
        return date;
    }

    public void setDate(final OffsetDateTime date) {
        this.date = date;
    }

    public String getActionStatus() {
        return actionStatus;
    }

    public void setActionStatus(final String actionStatus) {
        this.actionStatus = actionStatus;
    }
}
