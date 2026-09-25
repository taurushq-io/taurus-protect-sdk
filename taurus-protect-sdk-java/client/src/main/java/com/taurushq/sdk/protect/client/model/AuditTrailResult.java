package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a paginated result of audit trail entries.
 *
 * @see AuditTrail
 * @see AuditService
 */
public class AuditTrailResult extends CursorPagedResult {

    private List<AuditTrail> auditTrails;

    /**
     * Gets the audit trail entries of this page.
     *
     * @return the auditTrails
     */
    public List<AuditTrail> getAuditTrails() {
        return auditTrails;
    }

    /**
     * Sets the audit trail entries of this page.
     *
     * @param auditTrails the auditTrails
     */
    public void setAuditTrails(final List<AuditTrail> auditTrails) {
        this.auditTrails = auditTrails;
    }
}
