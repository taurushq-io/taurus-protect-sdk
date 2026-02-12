package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents the status of a job execution in the Taurus Protect system.
 * <p>
 * Contains timing information and the current status of a job run.
 *
 * @see Job
 * @see JobStatistics
 */
public class JobStatus {

    private String id;
    private OffsetDateTime startedAt;
    private OffsetDateTime updatedAt;
    private OffsetDateTime timeoutAt;
    private String message;
    private String status;

    public String getId() {
        return id;
    }

    public void setId(final String id) {
        this.id = id;
    }

    public OffsetDateTime getStartedAt() {
        return startedAt;
    }

    public void setStartedAt(final OffsetDateTime startedAt) {
        this.startedAt = startedAt;
    }

    public OffsetDateTime getUpdatedAt() {
        return updatedAt;
    }

    public void setUpdatedAt(final OffsetDateTime updatedAt) {
        this.updatedAt = updatedAt;
    }

    public OffsetDateTime getTimeoutAt() {
        return timeoutAt;
    }

    public void setTimeoutAt(final OffsetDateTime timeoutAt) {
        this.timeoutAt = timeoutAt;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(final String message) {
        this.message = message;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(final String status) {
        this.status = status;
    }
}
