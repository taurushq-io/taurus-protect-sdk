/**
 * Fee payer service for Taurus-PROTECT SDK.
 *
 * Provides methods for managing fee payers in the system.
 * Fee payers are accounts used to pay transaction fees on behalf of other
 * addresses, commonly used for sponsored transactions on EVM-compatible
 * blockchains like Ethereum.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { FeePayersApi } from '../internal/openapi/apis/FeePayersApi';
import { feePayerFromDto, feePayersFromDto } from '../mappers/fee-payer';
import type { FeePayer, ListFeePayersOptions } from '../models/fee-payer';
import {
  buildOffsetPagination,
  offsetRequest,
  type PaginatedResult,
} from '../models/pagination';
import { BaseService } from './base';
import { offsetQuery } from './paging';

/**
 * Service for managing fee payers.
 *
 * Provides methods to list all fee payers and get specific fee payers by ID.
 *
 * @example
 * ```typescript
 * // List all fee payers
 * const feePayers = await feePayerService.list();
 * for (const fp of feePayers) {
 *   console.log(`${fp.name}: ${fp.blockchain}/${fp.network}`);
 * }
 *
 * // List fee payers for a specific blockchain
 * const ethFeePayers = await feePayerService.list({
 *   blockchain: 'ETH',
 *   network: 'mainnet',
 * });
 *
 * // Get a specific fee payer
 * const feePayer = await feePayerService.get('fp-123');
 * console.log(`Fee payer: ${feePayer.name}`);
 * ```
 */
export class FeePayerService extends BaseService {
  private readonly feePayersApi: FeePayersApi;

  /**
   * Creates a new FeePayerService instance.
   *
   * @param feePayersApi - The FeePayersApi instance from the OpenAPI client
   */
  constructor(feePayersApi: FeePayersApi) {
    super();
    this.feePayersApi = feePayersApi;
  }

  /**
   * Lists a page of fee payers with optional filtering.
   *
   * @param options - Filters, `limit` (1-100, default 20) and `offset`
   * @returns The page of fee payers and its pagination, including the server's total
   * @throws {@link ValidationError} If limit or offset are out of bounds
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // First page
   * const { items, pagination } = await feePayerService.list();
   *
   * // Next page
   * if (pagination.hasMore) {
   *   await feePayerService.list({ offset: pagination.nextOffset });
   * }
   *
   * // Filter by blockchain and network
   * const ethFeePayers = await feePayerService.list({
   *   blockchain: 'ETH',
   *   network: 'mainnet',
   * });
   *
   * // Filter by specific IDs
   * const specific = await feePayerService.list({
   *   ids: ['fp-123', 'fp-456'],
   * });
   * ```
   */
  async list(options?: ListFeePayersOptions): Promise<PaginatedResult<FeePayer>> {
    const page = offsetRequest(options);

    return this.execute(async () => {
      const response = await this.feePayersApi.feePayerServiceGetFeePayers({
        ...offsetQuery(page),
        ids: options?.ids,
        blockchain: options?.blockchain,
        network: options?.network,
      });

      const rows = response.result ?? [];
      return {
        items: feePayersFromDto(rows),
        pagination: buildOffsetPagination('plus_rows', page, response, rows.length),
      };
    });
  }

  /**
   * Gets a fee payer by ID.
   *
   * @param id - The fee payer ID
   * @returns The fee payer
   * @throws {@link ValidationError} If id is empty
   * @throws {@link NotFoundError} If fee payer is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const feePayer = await feePayerService.get('fp-123');
   * console.log(`Name: ${feePayer.name}`);
   * console.log(`Blockchain: ${feePayer.blockchain}/${feePayer.network}`);
   * if (feePayer.feePayerInfo?.eth?.local) {
   *   console.log(`Address ID: ${feePayer.feePayerInfo.eth.local.addressId}`);
   * }
   * ```
   */
  async get(id: string): Promise<FeePayer> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response = await this.feePayersApi.feePayerServiceGetFeePayer({ id });

      const feePayer = feePayerFromDto(response.feepayer);

      if (!feePayer) {
        throw new NotFoundError(`Fee payer with id '${id}' not found`);
      }

      return feePayer;
    });
  }
}
