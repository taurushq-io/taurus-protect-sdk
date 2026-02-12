package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a single approval group requirement for whitelist operations.
 * <p>
 * An approvers group defines the external group (e.g., a team or role) that must
 * provide approvals and the minimum number of signatures required from that group.
 *
 * @see ParallelApproversGroups
 * @see Approvers
 */
public class ApproversGroup {

    /**
     * External identifier for the approval group (e.g., team or role ID).
     */
    private String externalGroupID;

    /**
     * Minimum number of signatures required from this group.
     */
    private long minimumSignatures;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the external identifier for this approval group.
     *
     * @return the external group identifier
     */
    public String getExternalGroupID() {
        return externalGroupID;
    }

    /**
     * Sets the external identifier for this approval group.
     *
     * @param externalGroupID the external group identifier to set
     */
    public void setExternalGroupID(String externalGroupID) {
        this.externalGroupID = externalGroupID;
    }

    /**
     * Returns the minimum number of signatures required from this group.
     *
     * @return the minimum number of required signatures
     */
    public long getMinimumSignatures() {
        return minimumSignatures;
    }

    /**
     * Sets the minimum number of signatures required from this group.
     *
     * @param minimumSignatures the minimum number of required signatures
     */
    public void setMinimumSignatures(long minimumSignatures) {
        this.minimumSignatures = minimumSignatures;
    }
}
