package com.taurushq.sdk.protect.client.model;

import java.util.Map;

/**
 * Represents the overall health status of the Taurus Protect system.
 * <p>
 * Contains a map of component names to their health component details,
 * which include groups of individual health checks.
 *
 * @see HealthService
 */
public class HealthCheck {

    private Map<String, HealthComponent> components;

    /**
     * Gets the map of health components by name.
     *
     * @return the components map
     */
    public Map<String, HealthComponent> getComponents() {
        return components;
    }

    /**
     * Sets the components map.
     *
     * @param components the components to set
     */
    public void setComponents(Map<String, HealthComponent> components) {
        this.components = components;
    }
}
