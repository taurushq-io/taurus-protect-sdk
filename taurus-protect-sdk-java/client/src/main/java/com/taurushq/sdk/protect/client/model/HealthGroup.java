package com.taurushq.sdk.protect.client.model;

import java.util.List;

/**
 * Represents a group of health checks in the Taurus Protect system.
 * <p>
 * Contains a list of individual health check results.
 *
 * @see HealthComponent
 */
public class HealthGroup {

    private List<HealthCheckStatus> healthChecks;

    /**
     * Gets the list of health check statuses.
     *
     * @return the health checks list
     */
    public List<HealthCheckStatus> getHealthChecks() {
        return healthChecks;
    }

    /**
     * Sets the health checks list.
     *
     * @param healthChecks the health checks to set
     */
    public void setHealthChecks(List<HealthCheckStatus> healthChecks) {
        this.healthChecks = healthChecks;
    }
}
