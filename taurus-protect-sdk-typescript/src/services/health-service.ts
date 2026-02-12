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

      const resp = response as Record<string, unknown>;
      const checks = resp.checks ?? resp.healthChecks ?? resp.result;

      // Determine overall health status
      let healthy = true;
      if (Array.isArray(checks)) {
        healthy = checks.every((check) => {
          const c = check as Record<string, unknown>;
          return c.healthy === true || c.status === 'healthy' || c.status === 'HEALTHY';
        });
      }

      return {
        status: healthy ? 'healthy' : 'unhealthy',
        version: resp.version as string | undefined,
        message: resp.message as string | undefined,
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

      const resp = response as Record<string, unknown>;
      const status = resp.status ?? resp.componentStatus ?? resp.result;

      let statusStr = 'unknown';
      let healthy = false;

      if (typeof status === 'string') {
        statusStr = status;
        healthy = status.toLowerCase() === 'healthy' || status.toLowerCase() === 'ok';
      } else if (typeof status === 'object' && status) {
        const s = status as Record<string, unknown>;
        statusStr = (s.status as string) ?? (s.state as string) ?? 'unknown';
        healthy = s.healthy === true || statusStr.toLowerCase() === 'healthy';
      }

      return {
        status: healthy ? 'healthy' : statusStr,
        version: resp.version as string | undefined,
        message: resp.message as string | undefined,
      };
    });
  }
}
