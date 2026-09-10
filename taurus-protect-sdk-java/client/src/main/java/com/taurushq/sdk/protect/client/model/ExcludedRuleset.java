package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * A governance rules history entry withheld because its SuperAdmin signatures did not
 * verify.
 *
 * <p>History is LENIENT where the single-ruleset reads are strict, and deliberately does
 * not fail when nothing survives: a SuperAdmin key rotation makes every pre-rotation
 * ruleset unverifiable, so aborting the page would deny access to the whole audit trail
 * from the rotation onwards. Naming what was dropped is what stops a shortened page
 * reading as a complete one.
 */
public class ExcludedRuleset {

    /**
     * When the excluded ruleset was created — the only stable identifier a history entry
     * carries.
     */
    private OffsetDateTime creationDate;

    /**
     * Why it was excluded.
     */
    private String reason;

    /**
     * Gets the creation date of the excluded ruleset.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the creation date of the excluded ruleset.
     *
     * @param creationDate the creation date
     */
    public void setCreationDate(final OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Gets the exclusion reason.
     *
     * @return the reason
     */
    public String getReason() {
        return reason;
    }

    /**
     * Sets the exclusion reason.
     *
     * @param reason the reason
     */
    public void setReason(final String reason) {
        this.reason = reason;
    }
}
