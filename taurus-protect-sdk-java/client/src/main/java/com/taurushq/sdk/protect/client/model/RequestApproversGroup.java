package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a single approval group requirement for a transaction request.
 * <p>
 * An approvers group defines the external group (e.g., a team or role) that must
 * provide approvals and the minimum number of signatures required from that group.
 * <p>
 * For example, a group might require at least 2 signatures from a "Treasury" team
 * before a large transfer can proceed.
 *
 * @see RequestParallelApproversGroups
 * @see RequestApprovers
 */
public class RequestApproversGroup {

    /**
     * External identifier for the approval group (e.g., team or role ID).
     */
    private String externalGroupID;

    /**
     * Minimum number of signatures required from this group.
     */
    private int minimumSignatures;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the external identifier for this approval group.
     * <p>
     * This identifier typically corresponds to a team, role, or organizational
     * unit in an external identity management system.
     *
     * @return the external group identifier, or {@code null} if not set
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
     * <p>
     * The request cannot proceed until at least this many members of
     * the group have approved it.
     *
     * @return the minimum number of required signatures
     */
    public int getMinimumSignatures() {
        return minimumSignatures;
    }

    /**
     * Sets the minimum number of signatures required from this group.
     *
     * @param minimumSignatures the minimum number of required signatures
     */
    public void setMinimumSignatures(int minimumSignatures) {
        this.minimumSignatures = minimumSignatures;
    }
}
