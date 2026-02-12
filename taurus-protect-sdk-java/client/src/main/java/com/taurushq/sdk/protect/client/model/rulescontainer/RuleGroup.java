package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.List;

/**
 * Represents a group of users defined in the governance rules container.
 * <p>
 * Groups allow governance rules to require approvals from a collection of users
 * rather than specific individuals. For example, a rule might require 2 signatures
 * from the "Finance" group rather than requiring signatures from specific users.
 *
 * @see RuleUser
 * @see GroupThreshold
 * @see DecodedRulesContainer
 */
public class RuleGroup {

    /**
     * Unique identifier for the group within the rules container.
     */
    private String id;

    /**
     * List of user IDs that belong to this group.
     */
    private List<String> userIds;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the group id.
     *
     * @return the id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the group id.
     *
     * @param id the id
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the user ids in this group.
     *
     * @return the user ids
     */
    public List<String> getUserIds() {
        return userIds;
    }

    /**
     * Sets the user ids in this group.
     *
     * @param userIds the user ids
     */
    public void setUserIds(List<String> userIds) {
        this.userIds = userIds;
    }
}
