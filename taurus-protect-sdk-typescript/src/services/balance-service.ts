/**
 * Balance service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving balance information across
 * wallets and addresses.
 */

import { ValidationError } from '../errors';
import type { BalancesApi } from '../internal/openapi/apis/BalancesApi';
import { assetBalancesFromDto, nftCollectionBalancesFromDto } from '../mappers/balance';
import type {
  AssetBalance,
  ListBalancesOptions,
  ListNFTCollectionBalancesOptions,
  NFTCollectionBalance,
} from '../models/balance';
import { BaseService } from './base';

/**
 * Service for retrieving balance information.
 *
 * Provides methods to list balances for all assets and NFT collections
 * across the tenant.
 *
 * @example
 * ```typescript
 * // List all balances
 * const balances = await balanceService.list();
 * for (const balance of balances) {
 *   console.log(`${balance.currency}: ${balance.balance}`);
 * }
 *
 * // List balances for a specific currency
 * const ethBalances = await balanceService.list({ currency: 'ETH' });
 *
 * // List NFT collection balances
 * const nftBalances = await balanceService.listNFTCollections({
 *   blockchain: 'ETH',
 *   network: 'mainnet',
 * });
 * ```
 */
export class BalanceService extends BaseService {
  private readonly balancesApi: BalancesApi;

  /**
   * Creates a new BalanceService instance.
   *
   * @param balancesApi - The BalancesApi instance from the OpenAPI client
   */
  constructor(balancesApi: BalancesApi) {
    super();
    this.balancesApi = balancesApi;
  }

  /**
   * Lists asset balances for the tenant.
   *
   * Each asset is identified by a full triplet of attributes (blockchain,
   * contract address, and token ID).
   *
   * @param options - Optional filtering options
   * @returns Array of asset balances
   * @throws {@link ValidationError} If limit is invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all balances
   * const balances = await balanceService.list();
   *
   * // List balances for a specific currency
   * const ethBalances = await balanceService.list({ currency: 'ETH', limit: 100 });
   * ```
   */
  async list(options?: ListBalancesOptions): Promise<AssetBalance[]> {
    const limit = options?.limit ?? 50;

    if (limit <= 0) {
      throw new ValidationError('limit must be positive');
    }

    return this.execute(async () => {
      const response = await this.balancesApi.walletServiceGetBalances({
        currency: options?.currency,
        limit: String(limit),
        requestCursorPageSize: String(limit),
      });

      const result =
        (response as Record<string, unknown>).balances ??
        (response as Record<string, unknown>).result;
      return assetBalancesFromDto(result as unknown[]);
    });
  }

  /**
   * Lists NFT collection balances for the tenant.
   *
   * @param options - Filtering options (blockchain and network are required)
   * @returns Array of NFT collection balances
   * @throws {@link ValidationError} If required arguments are missing or invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const nftBalances = await balanceService.listNFTCollections({
   *   blockchain: 'ETH',
   *   network: 'mainnet',
   *   onlyPositiveBalance: true,
   * });
   * for (const collection of nftBalances) {
   *   console.log(`${collection.name}: ${collection.count} NFTs`);
   * }
   * ```
   */
  async listNFTCollections(
    options: ListNFTCollectionBalancesOptions
  ): Promise<NFTCollectionBalance[]> {
    if (!options.blockchain || options.blockchain.trim() === '') {
      throw new ValidationError('blockchain is required');
    }
    if (!options.network || options.network.trim() === '') {
      throw new ValidationError('network is required');
    }

    const limit = options?.limit ?? 50;
    if (limit <= 0) {
      throw new ValidationError('limit must be positive');
    }

    return this.execute(async () => {
      const response = await this.balancesApi.walletServiceGetNFTCollectionBalances({
        blockchain: options.blockchain,
        network: options.network,
        cursorPageSize: String(limit),
        onlyPositiveBalance: options.onlyPositiveBalance,
      });

      const result =
        (response as Record<string, unknown>).collections ??
        (response as Record<string, unknown>).result;
      return nftCollectionBalancesFromDto(result as unknown[]);
    });
  }
}
