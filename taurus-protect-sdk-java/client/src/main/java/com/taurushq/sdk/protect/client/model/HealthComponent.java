package com.taurushq.sdk.protect.client.model;

import java.util.Map;

/**
 * Represents a health component in the Taurus Protect system.
 * <p>
 * Each component contains groups of health checks that monitor
 * specific aspects of the system.
 *
 * @see HealthCheck
 */
public class HealthComponent {

    private Map<String, HealthGroup> groups;

    /**
     * Gets the map of health groups by name.
     *
     * @return the groups map
     */
    public Map<String, HealthGroup> getGroups() {
        return groups;
    }

    /**
     * Sets the groups map.
     *
     * @param groups the groups to set
     */
    public void setGroups(Map<String, HealthGroup> groups) {
        this.groups = groups;
    }
}
