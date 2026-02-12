package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.ArrayList;
import java.util.List;

/**
 * Represents a set of approval groups that can operate in parallel for whitelist operations.
 * <p>
 * Within this structure, the sequential list contains groups that must
 * approve in order - each group's approval requirements must be met
 * before the next group can approve.
 *
 * @see ApproversGroup
 * @see Approvers
 */
public class ParallelApproversGroups {

    /**
     * List of approval groups that must approve sequentially within this parallel track.
     */
    private List<ApproversGroup> sequential = new ArrayList<>();

    /**
     * Returns the list of approval groups that must approve sequentially.
     *
     * @return list of sequential approval groups
     */
    public List<ApproversGroup> getSequential() {
        return sequential;
    }

    /**
     * Sets the list of approval groups that must approve sequentially.
     *
     * @param sequential the list of sequential approval groups to set
     */
    public void setSequential(List<ApproversGroup> sequential) {
        this.sequential = sequential;
    }

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }
}
