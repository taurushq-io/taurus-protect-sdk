/**
 * Fee service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving network fee information.
 */

import type { FeeApi } from '../internal/openapi/apis/FeeApi';
import { feesFromDto, feesV2FromDto } from '../mappers/fee';
import type { Fee, FeeV2 } from '../models/fee';
import { BaseService } from './base';

/**
 * Service for retrieving network fee information.
 *
 * Provides access to current network fees for various blockchains,
 * which can be used to estimate transaction costs.
 *
 * @example
 * ```typescript
 * // Get all current network fees (v2 - recommended)
 * const fees = await feeService.getFeesV2();
 * for (const fee of fees) {
 *   console.log(`${fee.currencyId}: ${fee.value} ${fee.denom}`);
 * }
 *
 * // Get fees using deprecated v1 API (key-value format)
 * const feesV1 = await feeService.getFees();
 * for (const fee of feesV1) {
 *   console.log(`${fee.key}: ${fee.value}`);
 * }
 * ```
 */
export class FeeService extends BaseService {
  private readonly feeApi: FeeApi;

  /**
   * Creates a new FeeService instance.
   *
   * @param feeApi - The FeeApi instance from the OpenAPI client
   */
  constructor(feeApi: FeeApi) {
    super();
    this.feeApi = feeApi;
  }

  /**
   * Retrieves current network fees for all supported blockchains (v1 API).
   *
   * @deprecated Use getFeesV2() for richer fee information including currency details
   * @returns Array of fees as key-value pairs
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const fees = await feeService.getFees();
   * for (const fee of fees) {
   *   console.log(`${fee.key}: ${fee.value}`);
   * }
   * ```
   */
  async getFees(): Promise<Fee[]> {
    return this.execute(async () => {
      const response = await this.feeApi.feeServiceGetFees();

      const result = (response as Record<string, unknown>).result;
      return feesFromDto(result as unknown[]);
    });
  }

  /**
   * Retrieves current native currency fees for all supported blockchains (v2 API).
   *
   * This is the recommended method for retrieving fee information as it provides
   * richer data including currency details, denominations, and update timestamps.
   *
   * @returns Array of fees with currency information
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const fees = await feeService.getFeesV2();
   * for (const fee of fees) {
   *   console.log(`${fee.currencyInfo?.symbol}: ${fee.value} ${fee.denom}`);
   *   if (fee.updateDate) {
   *     console.log(`  Last updated: ${fee.updateDate.toISOString()}`);
   *   }
   * }
   * ```
   */
  async getFeesV2(): Promise<FeeV2[]> {
    return this.execute(async () => {
      const response = await this.feeApi.feeServiceGetFeesV2();

      const result = (response as Record<string, unknown>).result;
      return feesV2FromDto(result as unknown[]);
    });
  }
}
