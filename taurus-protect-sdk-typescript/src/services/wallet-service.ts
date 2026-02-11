/**
 * Wallet service for Taurus-PROTECT SDK.
 *
 * Provides methods for wallet management operations including listing,
 * getting, creating wallets, and managing wallet attributes.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { WalletsApi } from '../internal/openapi/apis/WalletsApi';
import { assetBalancesFromDto } from '../mappers/balance';
import {
  balanceHistoryFromDto,
  walletFromCreateDto,
  walletFromDto,
  walletsFromDto,
} from '../mappers/wallet';
import type { AssetBalance } from '../models/balance';
import type { Pagination, PaginatedResult } from '../models/pagination';
import type {
  BalanceHistoryPoint,
  CreateWalletRequest,
  ListWalletsOptions,
  Wallet,
} from '../models/wallet';
import { BaseService } from './base';

/**
 * Service for wallet management operations.
 *
 * Provides methods to list, get, and create wallets, as well as
 * manage wallet attributes and retrieve balance history.
 *
 * @example
 * ```typescript
 * // List wallets
 * const result = await walletService.list({ limit: 50, offset: 0 });
 * for (const wallet of result.items) {
 *   console.log(`${wallet.name}: ${wallet.currency}`);
 * }
 *
 * // Get single wallet
 * const wallet = await walletService.get(123);
 * console.log(`Balance: ${wallet.balance?.totalConfirmed}`);
 *
 * // Create wallet
 * const newWallet = await walletService.create({
 *   blockchain: 'ETH',
 *   network: 'mainnet',
 *   name: 'Trading Wallet',
 * });
 * ```
 */
export class WalletService extends BaseService {
  private readonly walletsApi: WalletsApi;

  /**
   * Creates a new WalletService instance.
   *
   * @param walletsApi - The WalletsApi instance from the OpenAPI client
   */
  constructor(walletsApi: WalletsApi) {
    super();
    this.walletsApi = walletsApi;
  }

  /**
   * Gets a wallet by ID.
   *
   * @param walletId - The wallet ID to retrieve
   * @returns The wallet
   * @throws {@link ValidationError} If walletId is invalid
   * @throws {@link NotFoundError} If wallet not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const wallet = await walletService.get(123);
   * console.log(`Wallet: ${wallet.name}, Blockchain: ${wallet.blockchain}`);
   * ```
   */
  async get(walletId: number): Promise<Wallet> {
    if (walletId <= 0) {
      throw new ValidationError('walletId must be positive');
    }

    return this.execute(async () => {
      const response = await this.walletsApi.walletServiceGetWalletV2({
        id: String(walletId),
      });

      const wallet = walletFromDto(response.result);
      if (!wallet) {
        throw new NotFoundError(`Wallet ${walletId} not found`);
      }

      return wallet;
    });
  }

  /**
   * Lists wallets with pagination and optional filtering.
   *
   * @param options - Optional filtering and pagination options
   * @returns Paginated result containing wallets and pagination info
   * @throws {@link ValidationError} If limit or offset are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List first 50 wallets
   * const result = await walletService.list({ limit: 50, offset: 0 });
   * console.log(`Found ${result.pagination.totalItems} wallets`);
   *
   * // List with filters
   * const filtered = await walletService.list({
   *   blockchain: 'ETH',
   *   excludeDisabled: true,
   *   onlyPositiveBalance: true,
   * });
   * ```
   */
  async list(options?: ListWalletsOptions): Promise<PaginatedResult<Wallet>> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError('limit must be positive');
    }
    if (offset < 0) {
      throw new ValidationError('offset cannot be negative');
    }

    return this.execute(async () => {
      // Build currencies array from options
      const currencies =
        options?.currencies ??
        (options?.currency ? [options.currency] : undefined);

      const response = await this.walletsApi.walletServiceGetWalletsV2({
        currencies,
        query: options?.query,
        limit: String(limit),
        offset: offset > 0 ? String(offset) : undefined,
        name: options?.name,
        sortOrder: options?.sortOrder,
        excludeDisabled: options?.excludeDisabled,
        tagIDs: options?.tagIds,
        onlyPositiveBalance: options?.onlyPositiveBalance,
        blockchain: options?.blockchain,
        network: options?.network,
        ids: options?.ids,
      });

      const wallets = walletsFromDto(response.result);

      const pagination: Pagination = {
        totalItems: parseInt(response.totalItems ?? '0', 10),
        offset: parseInt(response.offset ?? '0', 10),
        limit,
      };

      return {
        items: wallets,
        pagination,
      };
    });
  }

  /**
   * Creates a new wallet.
   *
   * You must specify either a `currency` (by ID or symbol) or a combination
   * of `blockchain` and `network`. The wallet name must be unique per
   * blockchain/network combination.
   *
   * @param request - Wallet creation parameters
   * @returns The created wallet
   * @throws {@link ValidationError} If required fields are missing or invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Create with blockchain/network
   * const wallet = await walletService.create({
   *   blockchain: 'ETH',
   *   network: 'mainnet',
   *   name: 'My ETH Wallet',
   *   comment: 'Trading wallet',
   * });
   *
   * // Create with currency
   * const btcWallet = await walletService.create({
   *   currency: 'BTC',
   *   name: 'My BTC Wallet',
   *   isOmnibus: true,
   * });
   * ```
   */
  async create(request: CreateWalletRequest): Promise<Wallet> {
    if (!request.name || request.name.trim() === '') {
      throw new ValidationError('name is required');
    }

    // Validate that either currency or blockchain is provided
    if (!request.currency && !request.blockchain) {
      throw new ValidationError(
        'Either currency or blockchain must be provided'
      );
    }

    return this.execute(async () => {
      const response = await this.walletsApi.walletServiceCreateWallet({
        body: {
          name: request.name,
          currency: request.currency,
          blockchain: request.blockchain,
          network: request.network,
          isOmnibus: request.isOmnibus,
          comment: request.comment,
          customerId: request.customerId,
          visibilityGroupID: request.visibilityGroupId,
          externalWalletId: request.externalWalletId,
        },
      });

      const wallet = walletFromCreateDto(response.result);
      if (!wallet) {
        throw new ValidationError('Failed to create wallet: no result returned');
      }

      return wallet;
    });
  }

  /**
   * Creates an attribute for a wallet.
   *
   * Wallet attributes are key-value pairs that can be used to store
   * custom metadata on wallets.
   *
   * @param walletId - The wallet ID
   * @param key - The attribute key
   * @param value - The attribute value
   * @throws {@link ValidationError} If any argument is invalid
   * @throws {@link NotFoundError} If wallet not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await walletService.createAttribute(123, 'department', 'treasury');
   * await walletService.createAttribute(123, 'owner', 'finance-team');
   * ```
   */
  async createAttribute(
    walletId: number,
    key: string,
    value: string
  ): Promise<void> {
    if (walletId <= 0) {
      throw new ValidationError('walletId must be positive');
    }
    if (!key || key.trim() === '') {
      throw new ValidationError('key is required');
    }
    if (!value || value.trim() === '') {
      throw new ValidationError('value is required');
    }

    return this.execute(async () => {
      await this.walletsApi.walletServiceCreateWalletAttributes({
        walletId: String(walletId),
        body: {
          attributes: [{ key, value }],
        },
      });
    });
  }

  /**
   * Deletes an attribute from a wallet.
   *
   * @param walletId - The wallet ID
   * @param attributeId - The attribute ID to delete
   * @throws {@link ValidationError} If any argument is invalid
   * @throws {@link NotFoundError} If wallet or attribute not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get wallet to find attribute ID
   * const wallet = await walletService.get(123);
   * const attr = wallet.attributes.find(a => a.key === 'department');
   * if (attr?.id) {
   *   await walletService.deleteAttribute(123, attr.id);
   * }
   * ```
   */
  async deleteAttribute(walletId: number, attributeId: string): Promise<void> {
    if (walletId <= 0) {
      throw new ValidationError('walletId must be positive');
    }
    if (!attributeId || attributeId.trim() === '') {
      throw new ValidationError('attributeId is required');
    }

    return this.execute(async () => {
      await this.walletsApi.walletServiceDeleteWalletAttribute({
        walletId: String(walletId),
        id: attributeId,
      });
    });
  }

  /**
   * Gets the balance history for a wallet.
   *
   * Returns balance snapshots at regular intervals, useful for tracking
   * balance changes over time.
   *
   * @param walletId - The wallet ID
   * @param intervalHours - The interval in hours for balance snapshots
   * @returns Array of balance history points
   * @throws {@link ValidationError} If arguments are invalid
   * @throws {@link NotFoundError} If wallet not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get hourly balance history
   * const history = await walletService.getBalanceHistory(123, 1);
   * for (const point of history) {
   *   console.log(`${point.pointDate}: ${point.balance?.totalConfirmed}`);
   * }
   *
   * // Get daily balance history
   * const dailyHistory = await walletService.getBalanceHistory(123, 24);
   * ```
   */
  async getBalanceHistory(
    walletId: number,
    intervalHours: number
  ): Promise<BalanceHistoryPoint[]> {
    if (walletId <= 0) {
      throw new ValidationError('walletId must be positive');
    }
    if (intervalHours <= 0) {
      throw new ValidationError('intervalHours must be positive');
    }

    return this.execute(async () => {
      const response = await this.walletsApi.walletServiceGetWalletBalanceHistory(
        {
          id: String(walletId),
          intervalHours: String(intervalHours),
        }
      );

      return balanceHistoryFromDto(response.result);
    });
  }

  /**
   * Gets the list of tokens (asset balances) for a wallet.
   *
   * Returns the tokens with their balances held by the given wallet.
   *
   * @param walletId - The wallet ID
   * @param limit - Maximum number of tokens to return (default: 50)
   * @returns Array of asset balances
   * @throws {@link ValidationError} If arguments are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const tokens = await walletService.getWalletTokens(123);
   * for (const token of tokens) {
   *   console.log(`${token.currency}: ${token.balance}`);
   * }
   * ```
   */
  async getWalletTokens(walletId: number, limit: number = 50): Promise<AssetBalance[]> {
    if (walletId <= 0) {
      throw new ValidationError('walletId must be positive');
    }
    if (limit <= 0) {
      throw new ValidationError('limit must be positive');
    }

    return this.execute(async () => {
      const response = await this.walletsApi.walletServiceGetWalletTokens({
        id: String(walletId),
        limit: String(limit),
      });
      return assetBalancesFromDto(response.balances);
    });
  }
}
