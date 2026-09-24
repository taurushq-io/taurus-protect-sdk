/**
 * Health service for Taurus-PROTECT SDK.
 *
 * Provides methods for checking API health status.
 */

import { APIError } from '../errors';
import type { HealthApi } from '../internal/openapi/apis/HealthApi';
import type { HealthStatus } from '../models/health';
import { BaseService } from './base';

/**
 * Service for checking API health.
 *
 * Provides health check endpoints to verify API connectivity
 * and component status.
 *
 * @example
 * ```typescript
 * // Simple health check
 * const health = await healthService.check();
 * console.log(`Status: ${health.status}`);
 *
 * // Check if API is healthy
 * if (health.status === 'healthy') {
 *   console.log('API is operational');
 * } else {
 *   console.log('API has issues:', health.message);
 * }
 * ```
 */
export class HealthService extends BaseService {
  private readonly healthApi: HealthApi;

  /**
   * Creates a new HealthService instance.
   *
   * @param healthApi - The HealthApi instance from the OpenAPI client
   */
  constructor(healthApi: HealthApi) {
    super();
    this.healthApi = healthApi;
  }

  /**
   * Checks the API health status.
   *
   * Returns a health status object indicating whether the API is healthy.
   * This method catches errors and returns an unhealthy status instead
   * of throwing, making it safe to use for health monitoring.
   *
   * @returns Health status response
   *
   * @example
   * ```typescript
   * const health = await healthService.check();
   * if (health.status === 'healthy') {
   *   console.log('API is operational');
   * } else {
   *   console.error('API is unhealthy:', health.message);
   * }
   * ```
   */
  async check(): Promise<HealthStatus> {
    try {
      const response = await this.healthApi.healthServiceGetHealthChecks({});

      // The reply groups the checks by name; each check reports "success", "failure"
      // or "deactivated". Healthy means no check failed.
      const checks = Object.values(response.groups ?? {}).flatMap(
        (group) => group.healthChecks ?? []
      );
      const failed = checks.filter((check) => check.status?.toLowerCase() === 'failure');

      return {
        status: failed.length === 0 ? 'healthy' : 'unhealthy',
        message:
          failed.length === 0
            ? undefined
            : `${failed.length} of ${checks.length} health checks failed`,
      };
    } catch (error) {
      // If health check fails, return unhealthy status
      if (error instanceof APIError) {
        return {
          status: 'unhealthy',
          message: error.message,
        };
      }
      if (error instanceof Error) {
        return {
          status: 'unhealthy',
          message: error.message,
        };
      }
      return {
        status: 'unhealthy',
        message: 'Health check failed',
      };
    }
  }

  /**
   * Checks the global component status.
   *
   * Returns the overall status of all API components.
   *
   * @returns Health status response
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const status = await healthService.getGlobalStatus();
   * console.log(`Global status: ${status.status}`);
   * ```
   */
  async getGlobalStatus(): Promise<HealthStatus> {
    return this.execute(async () => {
      const response = await this.healthApi.statusServiceGetGlobalComponentStatus({});

      // clusterStatus is "up", "degraded" or "down".
      const clusterStatus = response.clusterStatus ?? 'unknown';
      const healthy = clusterStatus.toLowerCase() === 'up';

      return {
        status: healthy ? 'healthy' : clusterStatus,
        message: healthy
          ? undefined
          : `${response.working ?? '0'} of ${response.total ?? '0'} components working`,
      };
    });
  }
}
