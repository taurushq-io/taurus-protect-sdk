package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents the approvers configuration for a transaction request.
 * <p>
 * This class defines the approval structure required for a request to be processed.
 * Approvals can be organized in parallel groups, where each group can be approved
 * independently of others.
 * <p>
 * The approval hierarchy is:
 * <ul>
 *   <li>{@code RequestApprovers} - Top-level container</li>
 *   <li>{@code RequestParallelApproversGroups} - Groups that can approve in parallel</li>
 *   <li>{@code RequestApproversGroup} - Sequential approval requirements within a parallel group</li>
 * </ul>
 *
 * @see RequestParallelApproversGroups
 * @see RequestApproversGroup
 * @see Request
 */
public class RequestApprovers {

    /**
     * List of parallel approver groups that can approve the request independently.
     */
    private List<RequestParallelApproversGroups> parallel;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the list of parallel approver groups.
     * <p>
     * Each group in this list can approve the request independently. The request
     * is considered approved when any one of the parallel groups completes its
     * approval requirements.
     *
     * @return list of parallel approver groups, or {@code null} if not set
     */
    public List<RequestParallelApproversGroups> getParallel() {
        return parallel;
    }

    /**
     * Sets the list of parallel approver groups.
     *
     * @param parallel the list of parallel approver groups to set
     */
    public void setParallel(List<RequestParallelApproversGroups> parallel) {
        this.parallel = parallel;
    }
}
