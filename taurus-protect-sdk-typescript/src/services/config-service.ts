/**
 * Config service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving tenant configuration settings.
 */

import type { ConfigApi } from '../internal/openapi/apis/ConfigApi';
import { tenantConfigFromDto } from '../mappers/config';
import type { TenantConfig } from '../models/config';
import { BaseService } from './base';

/**
 * Service for retrieving tenant configuration in the Taurus-PROTECT system.
 *
 * This service provides access to the tenant's configuration settings,
 * including security requirements, feature flags, and system parameters.
 *
 * @example
 * ```typescript
 * // Get tenant configuration
 * const config = await client.configService.getTenantConfig();
 *
 * // Check if MFA is mandatory
 * if (config.mfaMandatory) {
 *   console.log('MFA is required for this tenant');
 * }
 *
 * // Get the base currency
 * const baseCurrency = config.baseCurrency;
 * console.log(`Base currency: ${baseCurrency}`);
 *
 * // Check Protect Engine status
 * if (config.protectEngineCold) {
 *   console.log(`Protect Engine version ${config.protectEngineVersion} is in cold mode`);
 * }
 * ```
 */
export class ConfigService extends BaseService {
  private readonly configApi: ConfigApi;

  /**
   * Creates a new ConfigService instance.
   *
   * @param configApi - The ConfigApi instance from the OpenAPI client
   */
  constructor(configApi: ConfigApi) {
    super();
    this.configApi = configApi;
  }

  /**
   * Retrieves the tenant configuration.
   *
   * Returns the configuration settings for the tenant associated
   * with the authenticated user.
   *
   * @returns The tenant configuration
   * @throws {@link APIError} If the API request fails
   *
   * @example
   * ```typescript
   * const config = await configService.getTenantConfig();
   * console.log(`Tenant ID: ${config.tenantId}`);
   * console.log(`Base currency: ${config.baseCurrency}`);
   * console.log(`MFA mandatory: ${config.mfaMandatory}`);
   * console.log(`Protect Engine version: ${config.protectEngineVersion}`);
   * ```
   */
  async getTenantConfig(): Promise<TenantConfig> {
    return this.execute(async () => {
      const response = await this.configApi.statusServiceGetConfigTenant();

      // The response has a 'config' property (not 'result')
      const configDto = response.config;

      const config = tenantConfigFromDto(configDto);
      return config ?? {};
    });
  }
}
