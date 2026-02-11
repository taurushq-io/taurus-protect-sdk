package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.ConfigMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.TenantConfig;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.ConfigApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetConfigTenantReply;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving tenant configuration in the Taurus Protect system.
 * <p>
 * This service provides access to the tenant's configuration settings,
 * including security requirements, feature flags, and system parameters.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get tenant configuration
 * TenantConfig config = client.getConfigService().getTenantConfig();
 *
 * // Check if MFA is mandatory
 * if (Boolean.TRUE.equals(config.getMFAMandatory())) {
 *     System.out.println("MFA is required for this tenant");
 * }
 *
 * // Get the base currency
 * String baseCurrency = config.getBaseCurrency();
 * }</pre>
 *
 * @see TenantConfig
 */
public class ConfigService {

    private final ConfigApi configApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final ConfigMapper configMapper;

    /**
     * Creates a new ConfigService.
     *
     * @param apiClient          the API client for making HTTP requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public ConfigService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(apiClient, "apiClient must not be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.configApi = new ConfigApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.configMapper = ConfigMapper.INSTANCE;
    }

    /**
     * Retrieves the tenant configuration.
     * <p>
     * Returns the configuration settings for the tenant associated
     * with the authenticated user.
     *
     * @return the tenant configuration
     * @throws ApiException if the API call fails
     */
    public TenantConfig getTenantConfig() throws ApiException {
        try {
            TgvalidatordGetConfigTenantReply reply = configApi.statusServiceGetConfigTenant();
            return configMapper.fromDTO(reply.getConfig());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
