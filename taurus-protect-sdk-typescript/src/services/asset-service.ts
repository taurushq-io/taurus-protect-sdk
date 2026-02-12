/**
 * Asset service for Taurus-PROTECT SDK.
 *
 * Provides methods for querying asset balances at address and wallet levels.
 */

import { ValidationError } from '../errors';
import type { AssetsApi } from '../internal/openapi/apis/AssetsApi';
import { addressesFromDto } from '../mappers/address';
import { walletsFromDto } from '../mappers/wallet';
import type { Address } from '../models/address';
import type { Wallet } from '../models/wallet';
import { BaseService } from './base';

/**
 * Options for filtering asset addresses.
 */
export interface GetAssetAddressesOptions {
  /** The currency code (e.g., "ETH", "BTC", "USDC") - required */
  readonly currency: string;
  /** Optional wallet ID to filter addresses */
  readonly walletId?: string;
  /** Optional address ID to filter */
  readonly addressId?: string;
  /** Maximum number of addresses to return */
  readonly limit?: string;
}

/**
 * Options for filtering asset wallets.
 */
export interface GetAssetWalletsOptions {
  /** The currency code (e.g., "ETH", "BTC", "USDC") - required */
  readonly currency: string;
  /** Optional wallet ID to filter */
  readonly walletId?: string;
  /** Optional wallet name to filter */
  readonly walletName?: string;
  /** Maximum number of wallets to return */
  readonly limit?: string;
}

/**
 * Service for querying asset balances at address and wallet levels.
 *
 * The AssetService provides methods to retrieve addresses and wallets that hold
 * a specific asset (cryptocurrency or token). This is useful for portfolio management,
 * compliance reporting, and understanding asset distribution across the organization.
 *
 * @example
 * ```typescript
 * // Get all addresses holding ETH
 * const ethAddresses = await assetService.getAssetAddresses({ currency: 'ETH' });
 *
 * // Get all wallets holding a specific token
 * const usdcWallets = await assetService.getAssetWallets({ currency: 'USDC' });
 *
 * // Get addresses for a specific asset filtered by wallet
 * const addresses = await assetService.getAssetAddresses({
 *   currency: 'BTC',
 *   walletId: 'wallet-123',
 * });
 * ```
 */
export class AssetService extends BaseService {
  private readonly assetsApi: AssetsApi;

  /**
   * Creates a new AssetService instance.
   *
   * @param assetsApi - The AssetsApi instance from the OpenAPI client
   */
  constructor(assetsApi: AssetsApi) {
    super();
    this.assetsApi = assetsApi;
  }

  /**
   * Retrieves addresses that hold a specific asset.
   *
   * @param options - Filter options including currency (required)
   * @returns A list of addresses holding the asset
   * @throws {@link ValidationError} If currency is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get all addresses holding ETH
   * const addresses = await assetService.getAssetAddresses({ currency: 'ETH' });
   *
   * // Get addresses for a specific wallet
   * const walletAddresses = await assetService.getAssetAddresses({
   *   currency: 'ETH',
   *   walletId: 'wallet-123',
   * });
   *
   * // Limit results
   * const limitedAddresses = await assetService.getAssetAddresses({
   *   currency: 'USDC',
   *   limit: '50',
   * });
   * ```
   */
  async getAssetAddresses(options: GetAssetAddressesOptions): Promise<Address[]> {
    if (!options.currency || options.currency.trim() === '') {
      throw new ValidationError('currency is required');
    }

    return this.execute(async () => {
      const response = await this.assetsApi.walletServiceGetAssetAddresses({
        body: {
          asset: {
            currency: options.currency,
          },
          walletId: options.walletId,
          addressId: options.addressId,
          limit: options.limit,
        },
      });

      return addressesFromDto(response.addresses);
    });
  }

  /**
   * Retrieves wallets that hold a specific asset.
   *
   * @param options - Filter options including currency (required)
   * @returns A list of wallets holding the asset
   * @throws {@link ValidationError} If currency is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get all wallets holding ETH
   * const wallets = await assetService.getAssetWallets({ currency: 'ETH' });
   *
   * // Limit results
   * const limitedWallets = await assetService.getAssetWallets({
   *   currency: 'USDC',
   *   limit: '50',
   * });
   * ```
   */
  async getAssetWallets(options: GetAssetWalletsOptions): Promise<Wallet[]> {
    if (!options.currency || options.currency.trim() === '') {
      throw new ValidationError('currency is required');
    }

    return this.execute(async () => {
      const response = await this.assetsApi.walletServiceGetAssetWallets({
        body: {
          asset: {
            currency: options.currency,
          },
          walletId: options.walletId,
          walletName: options.walletName,
          limit: options.limit,
        },
      });

      return walletsFromDto(response.wallets);
    });
  }
}
