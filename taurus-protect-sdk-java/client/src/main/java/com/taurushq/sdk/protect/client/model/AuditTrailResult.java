package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of audit trail entries.
 *
 * @see AuditTrail
 * @see AuditService
 */
public class AuditTrailResult {

    private List<AuditTrail> auditTrails;
    private ApiResponseCursor cursor;

    /**
     * Gets the list of audit trails.
     *
     * @return the audit trails
     */
    public List<AuditTrail> getAuditTrails() {
        return auditTrails;
    }

    /**
     * Sets the list of audit trails.
     *
     * @param auditTrails the audit trails to set
     */
    public void setAuditTrails(List<AuditTrail> auditTrails) {
        this.auditTrails = auditTrails;
    }

    /**
     * Gets the pagination cursor.
     *
     * @return the cursor
     */
    public ApiResponseCursor getCursor() {
        return cursor;
    }

    /**
     * Sets the pagination cursor.
     *
     * @param cursor the cursor to set
     */
    public void setCursor(ApiResponseCursor cursor) {
        this.cursor = cursor;
    }

    /**
     * Checks if there are more results available.
     *
     * @return true if more results are available
     */
    public boolean hasNext() {
        return cursor != null && cursor.hasNext();
    }
}
