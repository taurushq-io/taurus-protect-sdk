/**
 * Asset service for Taurus-PROTECT SDK.
 *
 * Provides methods for querying asset balances at address and wallet levels, and the v2
 * asset registry (assets, their holders and their lifecycle operations).
 */

import type { RulesContainerCache } from '../cache';
import { ConfigurationError, IntegrityError, ValidationError } from '../errors';
import { verifyAddressSignature } from '../helpers';
import type { AssetV2Api } from '../internal/openapi/apis/AssetV2Api';
import type { AssetsApi } from '../internal/openapi/apis/AssetsApi';
import { addressesFromDto } from '../mappers/address';
import {
  assetAddressesV2FromDto,
  assetOperationsV2FromDto,
  assetsV2FromDto,
} from '../mappers/asset';
import { walletsFromDto } from '../mappers/wallet';
import type {
  AssetAddressV2,
  ExcludedAssetAddressV2,
  GetAssetAddressesOptions,
  GetAssetWalletsOptions,
  ListAssetAddressesResult,
  ListAssetAddressesV2Result,
  ListAssetOperationsOptions,
  ListAssetOperationsResult,
  ListAssetsV2Result,
  ListAssetWalletsResult,
  QueryAssetAddressesOptions,
  QueryAssetsOptions,
} from '../models/asset';
import { MAX_PAGE_SIZE, buildCursorPage, cursorRequest } from '../models/pagination';
import { MAX_ADDRESS_IDS_PER_READ, type AddressService } from './address-service';
import { BaseService } from './base';
import { cursorQuery } from './paging';
import type { VerifiedLookup } from './row-level-error';
import type { WhitelistedAddressService } from './whitelisted-address-service';

/**
 * The verified readers the v2 asset-holders list completes its unsigned rows through.
 * Resolved on first use, so a page without internal or whitelisted holders needs
 * neither.
 */
export interface AssetHolderReaders {
  /** The managed-address reader: HSM signature checked */
  readonly addresses: () => AddressService;
  /** The whitelisted-address reader: 6-step verification */
  readonly whitelistedAddresses: () => WhitelistedAddressService;
}

const INTERNAL_HOLDER = 'ADDRESS_TYPE_V2_INTERNAL';
const WHITELISTED_HOLDER = 'ADDRESS_TYPE_V2_WHITELISTED';

/** Every distinct id `pick` yields for the rows of one holder type, in first-seen order. */
function holderIds(
  rows: readonly AssetAddressV2[],
  type: string,
  pick: (row: AssetAddressV2) => string | undefined
): string[] {
  const ids = new Set<string>();
  for (const row of rows) {
    const id = pick(row);
    if (row.addressType === type && id) {
      ids.add(id);
    }
  }
  return [...ids];
}

/**
 * Keeps an internal or whitelisted holder row only when the verified reader returned an
 * address under its id whose address string equals the row's; the kept row carries the
 * reader's address.
 */
function completeHolder(
  row: AssetAddressV2,
  id: string | undefined,
  idField: string,
  kind: string,
  lookup: VerifiedLookup<{ readonly address: string }>
): AssetAddressV2 | ExcludedAssetAddressV2 {
  if (!id) {
    return { id: row.address ?? '', reason: `${row.addressType} holder carries no ${idField}` };
  }
  const failure = lookup.failed.get(id);
  if (failure !== undefined) {
    return { id, reason: failure };
  }
  const verified = lookup.verified.get(id);
  if (verified === undefined) {
    return { id, reason: `${kind} ${id} was not returned by the verified ${kind} read` };
  }
  if (verified.address !== row.address) {
    return {
      id,
      reason:
        `holder address ${JSON.stringify(row.address ?? '')} differs from the verified ` +
        `${kind} ${id} (${JSON.stringify(verified.address)})`,
    };
  }
  return { ...row, address: verified.address, verified: true };
}

/**
 * Re-reads `ids` through `read` in requests of at most `perRead` ids, merging the
 * outcomes. Only the ids a request asked for are taken from its reply.
 */
async function lookupInChunks<T>(
  ids: readonly string[],
  perRead: number,
  read: (chunk: string[]) => Promise<VerifiedLookup<T>>
): Promise<VerifiedLookup<T>> {
  const verified = new Map<string, T>();
  const failed = new Map<string, string>();
  for (let start = 0; start < ids.length; start += perRead) {
    const chunk = ids.slice(start, start + perRead);
    const lookup = await read(chunk);
    for (const id of chunk) {
      const reason = lookup.failed.get(id);
      const value = lookup.verified.get(id);
      if (reason !== undefined) {
        failed.set(id, reason);
      } else if (value !== undefined) {
        verified.set(id, value);
      }
    }
  }
  return { verified, failed };
}

/**
 * Service for querying assets.
 *
 * `getAssetAddresses` / `getAssetWallets` list the addresses and wallets that hold an
 * asset (cryptocurrency or token); `queryAssets`, `queryAssetAddresses` and
 * `listAssetOperations` read the v2 asset registry. Every list is a cursor list: pass
 * `pagination.nextCursor` back as `cursor` while `pagination.hasMore` is true.
 *
 * @example
 * ```typescript
 * // Get the addresses holding ETH, page by page
 * let cursor: string | undefined;
 * do {
 *   const page = await assetService.getAssetAddresses({ currency: 'ETH', cursor });
 *   page.items.forEach((a) => console.log(a.address));
 *   cursor = page.pagination.hasMore ? page.pagination.nextCursor : undefined;
 * } while (cursor);
 *
 * // Get the wallets holding a token
 * const usdcWallets = await assetService.getAssetWallets({ currency: 'USDC' });
 *
 * // Query the v2 asset registry
 * const assets = await assetService.queryAssets({ blockchain: 'CANTON' });
 * ```
 */
export class AssetService extends BaseService {
  private readonly assetsApi: AssetsApi;
  private readonly rulesCache: RulesContainerCache;
  private readonly assetV2Api: AssetV2Api;
  private readonly holders: AssetHolderReaders;

  /**
   * Creates a new AssetService instance.
   *
   * Address signature verification is MANDATORY: getAssetAddresses returns the same
   * Address entity AddressService does, signature and all, and used to hand it over
   * unverified — so the mandatory verification there could be walked around by asking
   * for the same rows here.
   *
   * @param assetsApi - The AssetsApi instance from the OpenAPI client
   * @param rulesCache - Rules container cache for signature verification (required)
   * @param assetV2Api - The AssetV2Api instance from the OpenAPI client
   * @param holders - The verified readers queryAssetAddresses re-reads its rows through
   *   (required)
   * @throws ConfigurationError if rulesCache or holders is not provided
   */
  constructor(
    assetsApi: AssetsApi,
    rulesCache: RulesContainerCache,
    assetV2Api: AssetV2Api,
    holders: AssetHolderReaders
  ) {
    super();
    if (!rulesCache) {
      throw new ConfigurationError(
        "RulesContainerCache is required for AssetService — address signature verification is mandatory"
      );
    }
    if (!holders) {
      throw new ConfigurationError(
        "AssetHolderReaders are required for AssetService — asset holder verification is mandatory"
      );
    }
    this.assetsApi = assetsApi;
    this.rulesCache = rulesCache;
    this.assetV2Api = assetV2Api;
    this.holders = holders;
  }

  /**
   * Retrieves a page of the addresses that hold a specific asset, each signature-verified.
   *
   * Paged through `requestCursor` only; the legacy `limit` / bytes `cursor` fields are
   * never sent.
   *
   * @param options - Currency (required), filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of verified addresses and its pagination, including the server's total
   * @throws {@link ValidationError} If currency is empty or the paging options are invalid
   * @throws {@link IntegrityError} If an address signature does not verify
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const page = await assetService.getAssetAddresses({
   *   currency: 'ETH',
   *   walletId: 'wallet-123',
   *   pageSize: 50,
   * });
   * console.log(`${page.pagination.totalItems} addresses hold ETH`);
   * ```
   */
  async getAssetAddresses(options: GetAssetAddressesOptions): Promise<ListAssetAddressesResult> {
    if (!options.currency || options.currency.trim() === '') {
      throw new ValidationError('currency is required');
    }
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.assetsApi.walletServiceGetAssetAddresses({
        body: {
          asset: {
            currency: options.currency,
          },
          walletId: options.walletId,
          addressId: options.addressId,
          addresses: options.addresses,
          requestCursor: page.cursor,
        },
      });

      const addresses = addressesFromDto(response.addresses);

      // Fail-fast, as AddressService does: one unverifiable address is not a row to
      // skip past when the caller is choosing where funds go.
      if (addresses.length > 0) {
        const rules = await this.rulesCache.get();
        for (const address of addresses) {
          verifyAddressSignature(
            address.address,
            address.signature ?? "",
            rules,
            address.id
          );
        }
      }

      return {
        items: addresses,
        pagination: buildCursorPage(page.pageSize, response.cursor, {
          total: response.totalItems,
        }),
      };
    });
  }

  /**
   * Retrieves a page of the wallets that hold a specific asset.
   *
   * Paged through `requestCursor` only; the legacy `limit` / bytes `cursor` fields are
   * never sent.
   *
   * @param options - Currency (required), filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of wallets and its pagination, including the server's total
   * @throws {@link ValidationError} If currency is empty or the paging options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const page = await assetService.getAssetWallets({ currency: 'ETH' });
   * for (const wallet of page.items) {
   *   console.log(`${wallet.name}: ${wallet.balance?.totalConfirmed}`);
   * }
   * ```
   */
  async getAssetWallets(options: GetAssetWalletsOptions): Promise<ListAssetWalletsResult> {
    if (!options.currency || options.currency.trim() === '') {
      throw new ValidationError('currency is required');
    }
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.assetsApi.walletServiceGetAssetWallets({
        body: {
          asset: {
            currency: options.currency,
          },
          walletId: options.walletId,
          walletName: options.walletName,
          requestCursor: page.cursor,
        },
      });

      return {
        items: walletsFromDto(response.wallets),
        pagination: buildCursorPage(page.pageSize, response.cursor, {
          total: response.totalItems,
        }),
      };
    });
  }

  /**
   * Queries a page of the v2 asset registry.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of assets and its cursor pagination
   * @throws {@link ValidationError} If the paging options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const page = await assetService.queryAssets({ blockchain: 'CANTON', network: 'mainnet' });
   * page.items.forEach((a) => console.log(a.id, a.symbol, a.status));
   * ```
   */
  async queryAssets(options?: QueryAssetsOptions): Promise<ListAssetsV2Result> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.assetV2Api.assetServiceV2QueryAssetsV2({
        body: {
          cursor: page.cursor,
          blockchain: options?.blockchain,
          network: options?.network,
          symbol: options?.symbol,
          contractAddress: options?.contractAddress,
          label: options?.label,
          currencyName: options?.currencyName,
        },
      });

      return {
        items: assetsV2FromDto(response.result),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }

  /**
   * Queries a page of the holders of a v2 asset, verifying internal and whitelisted ones.
   *
   * The holder rows carry no signature, so each page is completed through the verified
   * readers:
   *
   *   INTERNAL    --addressID------------> managed addresses by id (HSM signature, 50 per read)
   *   WHITELISTED --whitelistedAddressID-> whitelisted addresses by id (6-step, 100 per read)
   *                   |
   *                   +- verified, same address -> kept, address from the reader, verified: true
   *                   +- no id / not returned / failed / different address -> excludedUnverified
   *   anything else ----------------------> kept as on-chain data, verified: false, no request
   *
   * An API error or an unusable rules container aborts the call; so does a page whose
   * rows all failed. Exclusions never move the cursor.
   *
   * @param assetId - The v2 asset ID
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of holders, the rows excluded as unverifiable, and the cursor
   *   pagination
   * @throws {@link ValidationError} If assetId is empty or the paging options are invalid
   * @throws {@link IntegrityError} If rows came back and none could be verified, or the
   *   rules container cannot verify anything
   * @throws {@link APIError} If an API request fails
   *
   * @example
   * ```typescript
   * const page = await assetService.queryAssetAddresses('asset-1', {
   *   kycStatus: 'KYC_STATUS_V2_APPROVED',
   * });
   * page.items
   *   .filter((h) => h.verified)
   *   .forEach((h) => console.log(h.address, h.balance));
   * ```
   */
  async queryAssetAddresses(
    assetId: string,
    options?: QueryAssetAddressesOptions
  ): Promise<ListAssetAddressesV2Result> {
    if (!assetId || assetId.trim() === '') {
      throw new ValidationError('assetId is required');
    }
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.assetV2Api.assetServiceV2QueryAssetAddressesV2({
        assetID: assetId,
        body: {
          cursor: page.cursor,
          addressType: options?.addressType,
          kycStatus: options?.kycStatus,
        },
      });

      const rows = assetAddressesV2FromDto(response.result);
      const { items, excludedUnverified } = await this.verifiedHolders(rows);

      // Rows came back but none survived: a systemic failure, not an empty page.
      if (rows.length > 0 && items.length === 0) {
        throw new IntegrityError(
          `all ${rows.length} asset holder(s) failed verification; ` +
            `first failure: ${excludedUnverified[0]?.reason ?? "unknown"}`
        );
      }

      return {
        items,
        pagination: buildCursorPage(page.pageSize, response.cursor),
        excludedUnverified,
      };
    });
  }

  /**
   * Completes a page of holder rows through the verified readers; see
   * {@link queryAssetAddresses}.
   */
  private async verifiedHolders(rows: readonly AssetAddressV2[]): Promise<{
    items: AssetAddressV2[];
    excludedUnverified: ExcludedAssetAddressV2[];
  }> {
    const internalIds = holderIds(rows, INTERNAL_HOLDER, (row) => row.addressId);
    const whitelistedIds = holderIds(rows, WHITELISTED_HOLDER, (row) => row.whitelistedAddressId);
    const managed = await lookupInChunks(internalIds, MAX_ADDRESS_IDS_PER_READ, (chunk) =>
      this.holders.addresses()._verifiedByIds(chunk)
    );
    const whitelisted = await lookupInChunks(whitelistedIds, MAX_PAGE_SIZE, (chunk) =>
      this.holders.whitelistedAddresses()._verifiedByIds(chunk)
    );

    const items: AssetAddressV2[] = [];
    const excludedUnverified: ExcludedAssetAddressV2[] = [];
    for (const row of rows) {
      let outcome: AssetAddressV2 | ExcludedAssetAddressV2;
      if (row.addressType === INTERNAL_HOLDER) {
        outcome = completeHolder(row, row.addressId, 'addressID', 'managed address', managed);
      } else if (row.addressType === WHITELISTED_HOLDER) {
        outcome = completeHolder(
          row,
          row.whitelistedAddressId,
          'whitelistedAddressID',
          'whitelisted address',
          whitelisted
        );
      } else {
        outcome = { ...row, verified: false };
      }
      if ('reason' in outcome) {
        excludedUnverified.push(outcome);
      } else {
        items.push(outcome);
      }
    }
    return { items, excludedUnverified };
  }

  /**
   * Lists a page of the lifecycle operations of a v2 asset.
   *
   * @param assetId - The v2 asset ID
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of operations and its cursor pagination
   * @throws {@link ValidationError} If assetId is empty or the paging options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const page = await assetService.listAssetOperations('asset-1', {
   *   status: 'ASSET_OPERATION_STATUS_V2_PENDING',
   * });
   * page.items.forEach((op) => console.log(op.type, op.status));
   * ```
   */
  async listAssetOperations(
    assetId: string,
    options?: ListAssetOperationsOptions
  ): Promise<ListAssetOperationsResult> {
    if (!assetId || assetId.trim() === '') {
      throw new ValidationError('assetId is required');
    }
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.assetV2Api.assetServiceV2ListAssetOperationsV2({
        assetID: assetId,
        ...cursorQuery(page),
        type: options?.type,
        status: options?.status,
      });

      return {
        items: assetOperationsV2FromDto(response.result),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }
}
