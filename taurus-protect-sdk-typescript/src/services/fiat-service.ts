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
  fiatProviderAccountResultFromDto,
  fiatProviderCounterpartyAccountFromDto,
  fiatProviderCounterpartyAccountResultFromDto,
  fiatProviderOperationFromDto,
  fiatProviderOperationResultFromDto,
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
  ListFiatProviderOperationsOptions,
} from '../models/fiat';
import { BaseService } from './base';

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
      const result =
        (response as Record<string, unknown>).fiatProviders ??
        (response as Record<string, unknown>).result;
      return fiatProvidersFromDto(result as unknown[]);
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
      const result =
        (response as Record<string, unknown>).result ?? response;
      const account = fiatProviderAccountFromDto(result);
      if (!account) {
        throw new ValidationError(`Fiat provider account '${id}' not found`);
      }
      return account;
    });
  }

  /**
   * Retrieves fiat provider accounts with optional filtering.
   *
   * @param options - Filter and pagination options
   * @returns Paginated result containing fiat provider accounts
   * @throws {@link ValidationError} If required options are missing
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

    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderAccounts({
          provider: options.provider,
          label: options.label,
          accountType: options.accountType,
          sortOrder: options.sortOrder,
          cursorCurrentPage: options.cursor?.currentPage,
          cursorPageRequest: options.cursor?.pageRequest,
          cursorPageSize: options.cursor?.pageSize?.toString(),
        });
      return fiatProviderAccountResultFromDto(response);
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
      const result =
        (response as Record<string, unknown>).result ?? response;
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
   * Retrieves fiat provider counterparty accounts with optional filtering.
   *
   * @param options - Filter and pagination options
   * @returns Paginated result containing fiat provider counterparty accounts
   * @throws {@link ValidationError} If required options are missing
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

    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderCounterpartyAccounts(
          {
            provider: options.provider,
            label: options.label,
            counterpartyID: options.counterpartyId,
            sortOrder: options.sortOrder,
            cursorCurrentPage: options.cursor?.currentPage,
            cursorPageRequest: options.cursor?.pageRequest,
            cursorPageSize: options.cursor?.pageSize?.toString(),
          }
        );
      return fiatProviderCounterpartyAccountResultFromDto(response);
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
      const result =
        (response as Record<string, unknown>).result ?? response;
      const operation = fiatProviderOperationFromDto(result);
      if (!operation) {
        throw new ValidationError(`Fiat provider operation '${id}' not found`);
      }
      return operation;
    });
  }

  /**
   * Retrieves fiat provider operations with optional filtering.
   *
   * @param options - Filter and pagination options (optional)
   * @returns Paginated result containing fiat provider operations
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
    return this.execute(async () => {
      const response =
        await this.fiatApi.fiatProviderServiceGetFiatProviderOperations({
          provider: options?.provider,
          label: options?.label,
          sortOrder: options?.sortOrder,
          cursorCurrentPage: options?.cursor?.currentPage,
          cursorPageRequest: options?.cursor?.pageRequest,
          cursorPageSize: options?.cursor?.pageSize?.toString(),
        });
      return fiatProviderOperationResultFromDto(response);
    });
  }
}
