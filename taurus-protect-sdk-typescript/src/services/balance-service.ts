/**
 * Balance service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving balance information across
 * wallets and addresses.
 */

import type { BalancesApi } from '../internal/openapi/apis/BalancesApi';
import { assetBalancesFromDto, nftCollectionBalancesFromDto } from '../mappers/balance';
import type {
  ListBalancesOptions,
  ListBalancesResult,
  ListNFTCollectionBalancesOptions,
  ListNFTCollectionBalancesResult,
} from '../models/balance';
import { buildCursorPage, cursorRequest } from '../models/pagination';
import { BaseService } from './base';
import { cursorQuery, requestCursorQuery } from './paging';

/**
 * Service for balance operations.
 *
 * Both lists are cursor lists: pass `pagination.nextCursor` back as `cursor` while
 * `pagination.hasMore` is true.
 *
 * @example
 * ```typescript
 * const page = await balanceService.list({ currency: 'ETH' });
 * for (const balance of page.items) {
 *   console.log(`${balance.currency}: ${balance.balance}`);
 * }
 * console.log(`${page.pagination.totalItems} balances in total`);
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
   * Lists a page of asset balances for the tenant.
   *
   * Each asset is identified by a full triplet of attributes (blockchain,
   * contract address, and token ID). Paged through `requestCursor` only; the legacy
   * `limit` / bytes `cursor` parameters are never sent.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of balances and its pagination, including the server's total
   * @throws {@link ValidationError} If the page size or cursor options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * let cursor: string | undefined;
   * do {
   *   const page = await balanceService.list({ pageSize: 100, cursor });
   *   page.items.forEach((b) => console.log(b.currency, b.balance));
   *   cursor = page.pagination.hasMore ? page.pagination.nextCursor : undefined;
   * } while (cursor);
   * ```
   */
  async list(options?: ListBalancesOptions): Promise<ListBalancesResult> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.balancesApi.walletServiceGetBalances({
        currency: options?.currency,
        tokenId: options?.tokenId,
        ...requestCursorQuery(page),
      });

      return {
        items: assetBalancesFromDto(response.balances),
        pagination: buildCursorPage(page.pageSize, response.cursor, { total: response.total }),
      };
    });
  }

  /**
   * Lists a page of NFT collection balances for the tenant.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of NFT collection balances and its pagination
   * @throws {@link ValidationError} If the page size or cursor options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const page = await balanceService.listNFTCollections({
   *   blockchain: 'ETH',
   *   network: 'mainnet',
   *   onlyPositiveBalance: true,
   * });
   * for (const collection of page.items) {
   *   console.log(`${collection.name}: ${collection.count} NFTs`);
   * }
   * ```
   */
  async listNFTCollections(
    options?: ListNFTCollectionBalancesOptions
  ): Promise<ListNFTCollectionBalancesResult> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.balancesApi.walletServiceGetNFTCollectionBalances({
        blockchain: options?.blockchain,
        network: options?.network,
        query: options?.query,
        onlyPositiveBalance: options?.onlyPositiveBalance,
        ...cursorQuery(page),
      });

      return {
        items: nftCollectionBalancesFromDto(response.balances),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }
}
