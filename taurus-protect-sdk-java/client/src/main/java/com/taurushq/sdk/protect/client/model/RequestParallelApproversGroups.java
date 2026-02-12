package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents a set of approval groups that can operate in parallel.
 * <p>
 * Within this structure, the sequential list contains groups that must
 * approve in order - each group's approval requirements must be met
 * before the next group can approve.
 * <p>
 * Multiple {@code RequestParallelApproversGroups} can exist at the
 * {@link RequestApprovers} level, and only one parallel track needs
 * to complete for the request to be approved.
 *
 * @see RequestApproversGroup
 * @see RequestApprovers
 */
public class RequestParallelApproversGroups {

    /**
     * List of approval groups that must approve sequentially within this parallel track.
     */
    private List<RequestApproversGroup> sequential;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the list of approval groups that must approve sequentially.
     * <p>
     * The groups in this list are processed in order - each group must
     * meet its minimum signature requirements before the next group
     * can begin approving.
     *
     * @return list of sequential approval groups, or {@code null} if not set
     */
    public List<RequestApproversGroup> getSequential() {
        return sequential;
    }

    /**
     * Sets the list of approval groups that must approve sequentially.
     *
     * @param sequential the list of sequential approval groups to set
     */
    public void setSequential(List<RequestApproversGroup> sequential) {
        this.sequential = sequential;
    }
}
