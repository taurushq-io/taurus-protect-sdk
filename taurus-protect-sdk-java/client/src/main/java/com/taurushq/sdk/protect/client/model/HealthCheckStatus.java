package com.taurushq.sdk.protect.client.model;

/**
 * Represents the status of an individual health check.
 *
 * @see HealthGroup
 */
public class HealthCheckStatus {

    private String tenantId;
    private String componentName;
    private String componentId;
    private String group;
    private String healthCheck;
    private String status;

    /**
     * Gets the tenant ID.
     *
     * @return the tenant ID
     */
    public String getTenantId() {
        return tenantId;
    }

    /**
     * Sets the tenant ID.
     *
     * @param tenantId the tenant ID to set
     */
    public void setTenantId(String tenantId) {
        this.tenantId = tenantId;
    }

    /**
     * Gets the component name.
     *
     * @return the component name
     */
    public String getComponentName() {
        return componentName;
    }

    /**
     * Sets the component name.
     *
     * @param componentName the component name to set
     */
    public void setComponentName(String componentName) {
        this.componentName = componentName;
    }

    /**
     * Gets the component ID.
     *
     * @return the component ID
     */
    public String getComponentId() {
        return componentId;
    }

    /**
     * Sets the component ID.
     *
     * @param componentId the component ID to set
     */
    public void setComponentId(String componentId) {
        this.componentId = componentId;
    }

    /**
     * Gets the group name.
     *
     * @return the group name
     */
    public String getGroup() {
        return group;
    }

    /**
     * Sets the group name.
     *
     * @param group the group name to set
     */
    public void setGroup(String group) {
        this.group = group;
    }

    /**
     * Gets the health check name.
     *
     * @return the health check name
     */
    public String getHealthCheck() {
        return healthCheck;
    }

    /**
     * Sets the health check name.
     *
     * @param healthCheck the health check name to set
     */
    public void setHealthCheck(String healthCheck) {
        this.healthCheck = healthCheck;
    }

    /**
     * Gets the health check status (e.g., "healthy", "unhealthy").
     *
     * @return the status
     */
    public String getStatus() {
        return status;
    }

    /**
     * Sets the status.
     *
     * @param status the status to set
     */
    public void setStatus(String status) {
        this.status = status;
    }
}
