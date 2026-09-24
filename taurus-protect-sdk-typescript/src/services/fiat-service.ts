/**
 * Fiat service for Taurus-PROTECT SDK.
 *
 * Provides methods for managing fiat provider accounts, counterparty accounts,
 * and operations in the Taurus-PROTECT system.
 */

import { ValidationError } from '../errors';
import type { FiatApi } from '../internal/openapi/apis/FiatApi';
import {
  fiatProvidersFromDto,
  fiatProviderAccountFromDto,
  fiatProviderAccountsFromDto,
  fiatProviderCounterpartyAccountFromDto,
  fiatProviderCounterpartyAccountsFromDto,
  fiatProviderEntitiesFromDto,
  fiatProviderOperationFromDto,
  fiatProviderOperationsFromDto,
} from '../mappers/fiat';
import type {
  FiatProvider,
  FiatProviderAccount,
  FiatProviderAccountResult,
  FiatProviderCounterpartyAccount,
  FiatProviderCounterpartyAccountResult,
  FiatProviderOperation,
  FiatProviderOperationResult,
  ListFiatProviderAccountsOptions,
  ListFiatProviderCounterpartyAccountsOptions,
  ListFiatProviderEntitiesOptions,
  ListFiatProviderEntitiesResult,
  ListFiatProviderOperationsOptions,
} from '../models/fiat';
import { buildCursorPage, cursorRequest } from '../models/pagination';
import { BaseService } from './base';
import { cursorQuery } from './paging';

/**
 * Service for managing fiat provider operations in the Taurus-PROTECT system.
 *
 * This service provides access to fiat provider accounts, counterparty accounts,
 * and operations for fiat currency management.
 *
 * @example
 * ```typescript
 * // List all fiat providers
 * const providers = await fiatService.getFiatProviders();
 *
 * // Get a specific account
 * const account = await fiatService.getFiatProviderAccount('account-123');
 *
 * // List accounts with filtering
 * const result = await fiatService.getFiatProviderAccounts({
 *   provider: 'circle',
 *   label: 'main',
 * });
 * ```
 */
export class FiatService extends BaseService {
  private readonly fiatApi: FiatApi;

  /**
   * Creates a new FiatService instance.
   *
   * @param fiatApi - The FiatApi instance from the OpenAPI client
   */
  constructor(fiatApi: FiatApi) {
    super();
    this.fiatApi = fiatApi;
  }

  /**
   * Retrieves all configured fiat providers.
   *
   * @returns Array of fiat providers
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const providers = await fiatService.getFiatProviders();
   * for (const provider of providers) {
   *   console.log(`${provider.provider}: ${provider.baseCurrencyValuation}`);
   * }
   * ```
   */
  async getFiatProviders(): Promise<FiatProvider[]> {
    return this.execute(async () => {
      const response = await this.fiatApi.fiatProviderServiceGetFiatProviders();
      return fiatProvidersFromDto(response.fiatProviders);
    });
  }

  /**
   * Retrieves a fiat provider account by ID.
   *
   * @param id - The account ID
   * @returns The fiat provider account
   * @throws {@link ValidationError} If id is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const account = await fiatService.getFiatProviderAccount('account-123');
   * console.log(`Balance: ${account.totalBalance} ${account.currencyId}`);
   * ```
   */
  async getFiatProviderAccount(id: string): Promise<FiatProviderAccount> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderAccount({ id });
      const result = response.result;
      const account = fiatProviderAccountFromDto(result);
      if (!account) {
        throw new ValidationError(`Fiat provider account '${id}' not found`);
      }
      return account;
    });
  }

  /**
   * Retrieves a page of fiat provider accounts.
   *
   * @param options - Provider and label (required), filters, `pageSize` (1-100,
   *   default 20) and `cursor`
   * @returns The page of accounts and its cursor pagination
   * @throws {@link ValidationError} If required options are missing or paging options
   *   are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const result = await fiatService.getFiatProviderAccounts({
   *   provider: 'circle',
   *   label: 'main',
   *   accountType: 'wallet',
   *   sortOrder: 'ASC',
   * });
   * for (const account of result.accounts) {
   *   console.log(`${account.accountName}: ${account.totalBalance}`);
   * }
   * ```
   */
  async getFiatProviderAccounts(
    options: ListFiatProviderAccountsOptions
  ): Promise<FiatProviderAccountResult> {
    if (!options.provider || options.provider.trim() === '') {
      throw new ValidationError('provider is required');
    }
    if (!options.label || options.label.trim() === '') {
      throw new ValidationError('label is required');
    }

    const page = cursorRequest(options);

    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderAccounts({
          provider: options.provider,
          label: options.label,
          accountType: options.accountType,
          sortOrder: options.sortOrder,
          ...cursorQuery(page),
        });
      return {
        accounts: fiatProviderAccountsFromDto(response.result),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }

  /**
   * Retrieves a fiat provider counterparty account by ID.
   *
   * @param id - The counterparty account ID
   * @returns The fiat provider counterparty account
   * @throws {@link ValidationError} If id is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const account = await fiatService.getFiatProviderCounterpartyAccount('cp-123');
   * console.log(`Counterparty: ${account.counterpartyName}`);
   * ```
   */
  async getFiatProviderCounterpartyAccount(
    id: string
  ): Promise<FiatProviderCounterpartyAccount> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderCounterpartyAccount(
          { id }
        );
      const result = response.result;
      const account = fiatProviderCounterpartyAccountFromDto(result);
      if (!account) {
        throw new ValidationError(
          `Fiat provider counterparty account '${id}' not found`
        );
      }
      return account;
    });
  }

  /**
   * Retrieves a page of fiat provider counterparty accounts.
   *
   * @param options - Provider and label (required), filters, `pageSize` (1-100,
   *   default 20) and `cursor`
   * @returns The page of counterparty accounts and its cursor pagination
   * @throws {@link ValidationError} If required options are missing or paging options
   *   are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const result = await fiatService.getFiatProviderCounterpartyAccounts({
   *   provider: 'cubnet',
   *   label: 'main',
   *   counterpartyId: 'counterparty-123',
   * });
   * for (const account of result.accounts) {
   *   console.log(`${account.counterpartyName}: ${account.accountType}`);
   * }
   * ```
   */
  async getFiatProviderCounterpartyAccounts(
    options: ListFiatProviderCounterpartyAccountsOptions
  ): Promise<FiatProviderCounterpartyAccountResult> {
    if (!options.provider || options.provider.trim() === '') {
      throw new ValidationError('provider is required');
    }
    if (!options.label || options.label.trim() === '') {
      throw new ValidationError('label is required');
    }

    const page = cursorRequest(options);

    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderCounterpartyAccounts({
          provider: options.provider,
          label: options.label,
          counterpartyID: options.counterpartyId,
          sortOrder: options.sortOrder,
          ...cursorQuery(page),
        });
      return {
        accounts: fiatProviderCounterpartyAccountsFromDto(response.result),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }

  /**
   * Retrieves a fiat provider operation by ID.
   *
   * @param id - The operation ID
   * @returns The fiat provider operation
   * @throws {@link ValidationError} If id is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const operation = await fiatService.getFiatProviderOperation('op-123');
   * console.log(`Status: ${operation.status}, Amount: ${operation.amount}`);
   * ```
   */
  async getFiatProviderOperation(id: string): Promise<FiatProviderOperation> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderOperation({ id });
      const result = response.result;
      const operation = fiatProviderOperationFromDto(result);
      if (!operation) {
        throw new ValidationError(`Fiat provider operation '${id}' not found`);
      }
      return operation;
    });
  }

  /**
   * Retrieves a page of fiat provider operations.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of operations and its cursor pagination
   * @throws {@link ValidationError} If the paging options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all operations
   * const result = await fiatService.getFiatProviderOperations();
   *
   * // List operations with filtering
   * const filtered = await fiatService.getFiatProviderOperations({
   *   provider: 'circle',
   *   label: 'main',
   *   sortOrder: 'DESC',
   * });
   * for (const op of filtered.operations) {
   *   console.log(`${op.operationType}: ${op.amount} (${op.status})`);
   * }
   * ```
   */
  async getFiatProviderOperations(
    options?: ListFiatProviderOperationsOptions
  ): Promise<FiatProviderOperationResult> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderOperations({
          provider: options?.provider,
          label: options?.label,
          sortOrder: options?.sortOrder,
          ...cursorQuery(page),
        });
      return {
        operations: fiatProviderOperationsFromDto(response.result),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }

  /**
   * Lists a page of fiat provider entities.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of entities and its cursor pagination
   * @throws {@link ValidationError} If the page size is out of bounds
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * let cursor: string | undefined;
   * do {
   *   const page = await fiatService.listFiatProviderEntities({ provider: 'bank', cursor });
   *   page.items.forEach((e) => console.log(e.name));
   *   cursor = page.pagination.hasMore ? page.pagination.nextCursor : undefined;
   * } while (cursor);
   * ```
   */
  async listFiatProviderEntities(
    options?: ListFiatProviderEntitiesOptions
  ): Promise<ListFiatProviderEntitiesResult> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderEntities({
          provider: options?.provider,
          label: options?.label,
          sortOrder: options?.sortOrder,
          ...cursorQuery(page),
        });
      return {
        items: fiatProviderEntitiesFromDto(response.result),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }
}
