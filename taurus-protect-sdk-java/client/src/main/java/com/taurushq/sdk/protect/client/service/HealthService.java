package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.HealthMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.HealthCheck;
import com.taurushq.sdk.protect.client.model.HealthComponent;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.HealthApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAllHealthChecksReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordHealthComponent;

import java.util.HashMap;
import java.util.Map;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving health status of the Taurus Protect system.
 * <p>
 * This service provides access to health check information for monitoring
 * system components and their operational status.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get all health checks
 * HealthCheck health = client.getHealthService().getAllHealthChecks();
 *
 * // Get health checks with fail-on-unhealthy option
 * HealthCheck health = client.getHealthService().getAllHealthChecks(null, true);
 * }</pre>
 *
 * @see HealthCheck
 */
public class HealthService {

    private final HealthApi healthApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final HealthMapper healthMapper;

    /**
     * Creates a new HealthService.
     *
     * @param apiClient          the API client for making HTTP requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public HealthService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(apiClient, "apiClient must not be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.healthApi = new HealthApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.healthMapper = HealthMapper.INSTANCE;
    }

    /**
     * Retrieves all health checks for the system.
     *
     * @return the health check results
     * @throws ApiException if the API call fails
     */
    public HealthCheck getAllHealthChecks() throws ApiException {
        return getAllHealthChecks(null, null);
    }

    /**
     * Retrieves all health checks with optional parameters.
     *
     * @param tenantId        optional tenant ID to filter results
     * @param failIfUnhealthy if true, the API will return an error if any check is unhealthy
     * @return the health check results
     * @throws ApiException if the API call fails
     */
    public HealthCheck getAllHealthChecks(final String tenantId,
                                          final Boolean failIfUnhealthy) throws ApiException {
        try {
            TgvalidatordGetAllHealthChecksReply reply = healthApi.healthServiceGetAllHealthChecks(
                    tenantId, failIfUnhealthy);

            HealthCheck healthCheck = new HealthCheck();
            if (reply.getComponents() != null) {
                Map<String, HealthComponent> components = new HashMap<>();
                for (Map.Entry<String, TgvalidatordHealthComponent> entry : reply.getComponents().entrySet()) {
                    components.put(entry.getKey(), healthMapper.fromComponentDTO(entry.getValue()));
                }
                healthCheck.setComponents(components);
            }
            return healthCheck;
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
