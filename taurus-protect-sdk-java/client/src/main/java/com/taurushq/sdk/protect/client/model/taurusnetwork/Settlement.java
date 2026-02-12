package com.taurushq.sdk.protect.client.model.taurusnetwork;

import java.time.OffsetDateTime;

/**
 * Represents a settlement in the Taurus Network.
 */
public class Settlement {

    private String id;
    private String creatorParticipantID;
    private String targetParticipantID;
    private String firstLegParticipantID;
    private String status;
    private String workflowID;
    private OffsetDateTime startExecutionDate;
    private OffsetDateTime createdAt;
    private OffsetDateTime updatedAt;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public String getCreatorParticipantID() {
        return creatorParticipantID;
    }

    public void setCreatorParticipantID(final String creatorParticipantID) {
        this.creatorParticipantID = creatorParticipantID;
    }

    public String getTargetParticipantID() {
        return targetParticipantID;
    }

    public void setTargetParticipantID(final String targetParticipantID) {
        this.targetParticipantID = targetParticipantID;
    }

    public String getFirstLegParticipantID() {
        return firstLegParticipantID;
    }

    public void setFirstLegParticipantID(final String firstLegParticipantID) {
        this.firstLegParticipantID = firstLegParticipantID;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
    }

    public String getWorkflowID() {
        return workflowID;
    }

    public void setWorkflowID(final String workflowID) {
        this.workflowID = workflowID;
    }

    public OffsetDateTime getStartExecutionDate() {
        return startExecutionDate;
    }

    public void setStartExecutionDate(final OffsetDateTime startExecutionDate) {
        this.startExecutionDate = startExecutionDate;
    }

    public OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(final OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }

    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(final OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }
}
