/**
 * Fee service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving network fee information.
 */

import type { FeeApi } from '../internal/openapi/apis/FeeApi';
import { feesV2FromDto } from '../mappers/fee';
import type { FeeV2 } from '../models/fee';
import { BaseService } from './base';

/**
 * Service for retrieving network fee information.
 *
 * Provides access to current network fees for various blockchains,
 * which can be used to estimate transaction costs. Served by the v2 fees endpoint only;
 * the deprecated v1 key-value endpoint is not wrapped.
 *
 * @example
 * ```typescript
 * const fees = await feeService.getFeesV2();
 * for (const fee of fees) {
 *   console.log(`${fee.currencyId}: ${fee.value} ${fee.denom}`);
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
   * Retrieves current native currency fees for all supported blockchains (v2 API).
   *
   * Provides currency details, denominations, and update timestamps.
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
      return feesV2FromDto(response.result);
    });
  }
}
