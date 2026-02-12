package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;
import java.util.List;

/**
 * Represents an action envelope containing an automated action configuration
 * with its metadata and execution history.
 * <p>
 * Actions in Taurus Protect allow automated workflows to be triggered based
 * on specific conditions such as balance thresholds.
 *
 * @see Action
 * @see ActionTrigger
 * @see ActionTask
 */
public class ActionEnvelope {

    private String id;
    private String tenantId;
    private String label;
    private Action action;
    private String status;
    private OffsetDateTime creationDate;
    private OffsetDateTime updateDate;
    private OffsetDateTime lastCheckedDate;
    private Boolean autoApprove;
    private List<ActionAttribute> attributes;
    private List<ActionTrail> trails;

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

    public String getLabel() {
        return label;
    }

    public void setLabel(final String label) {
        this.label = label;
    }

    public Action getAction() {
        return action;
    }

    public void setAction(final Action action) {
        this.action = action;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
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

    public OffsetDateTime getLastCheckedDate() {
        return lastCheckedDate;
    }

    public void setLastCheckedDate(final OffsetDateTime lastCheckedDate) {
        this.lastCheckedDate = lastCheckedDate;
    }

    public Boolean getAutoApprove() {
        return autoApprove;
    }

    public void setAutoApprove(final Boolean autoApprove) {
        this.autoApprove = autoApprove;
    }

    public List<ActionAttribute> getAttributes() {
        return attributes;
    }

    public void setAttributes(final List<ActionAttribute> attributes) {
        this.attributes = attributes;
    }

    public List<ActionTrail> getTrails() {
        return trails;
    }

    public void setTrails(final List<ActionTrail> trails) {
        this.trails = trails;
    }
}
