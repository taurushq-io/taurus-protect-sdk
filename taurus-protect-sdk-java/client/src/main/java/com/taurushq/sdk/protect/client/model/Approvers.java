package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.ArrayList;
import java.util.List;

/**
 * Represents the approvers configuration for whitelist operations.
 * <p>
 * This class defines the approval structure required for whitelist entries to be
 * approved. Approvals can be organized in parallel groups, where each group can
 * approve independently of others.
 * <p>
 * The approval hierarchy is:
 * <ul>
 *   <li>{@code Approvers} - Top-level container</li>
 *   <li>{@code ParallelApproversGroups} - Groups that can approve in parallel</li>
 *   <li>{@code ApproversGroup} - Sequential approval requirements within a parallel group</li>
 * </ul>
 *
 * @see ParallelApproversGroups
 * @see ApproversGroup
 * @see SignedWhitelistedAddress
 */
public class Approvers {

    /**
     * List of parallel approver groups that can approve the whitelist entry independently.
     */
    private List<ParallelApproversGroups> parallel = new ArrayList<>();

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the list of parallel approver groups.
     * <p>
     * Each group in this list can approve the whitelist entry independently.
     *
     * @return list of parallel approver groups
     */
    public List<ParallelApproversGroups> getParallel() {
        return parallel;
    }

    /**
     * Sets the list of parallel approver groups.
     *
     * @param parallel the list of parallel approver groups to set
     */
    public void setParallel(List<ParallelApproversGroups> parallel) {
        this.parallel = parallel;
    }
}
