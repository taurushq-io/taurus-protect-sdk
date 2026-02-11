package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents the signature threshold for a specific group.
 * <p>
 * A group threshold specifies how many signatures are required from users
 * within a particular group to satisfy this part of the approval requirement.
 *
 * @see RuleGroup
 * @see SequentialThresholds
 */
public class GroupThreshold {

    /**
     * The ID of the group this threshold applies to.
     */
    private String groupId;

    /**
     * The minimum number of signatures required from users in this group.
     */
    private int minimumSignatures;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the group id.
     *
     * @return the group id
     */
    public String getGroupId() {
        return groupId;
    }

    /**
     * Sets the group id.
     *
     * @param groupId the group id
     */
    public void setGroupId(String groupId) {
        this.groupId = groupId;
    }

    /**
     * Gets the minimum signatures required.
     *
     * @return the minimum signatures
     */
    public int getMinimumSignatures() {
        return minimumSignatures;
    }

    /**
     * Sets the minimum signatures required.
     *
     * @param minimumSignatures the minimum signatures
     */
    public void setMinimumSignatures(int minimumSignatures) {
        this.minimumSignatures = minimumSignatures;
    }
}
