/**
 * Exchange service for Taurus-PROTECT SDK.
 *
 * Provides methods for managing exchange accounts in the Taurus-PROTECT system.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { ExchangeApi } from '../internal/openapi/apis/ExchangeApi';
import {
  exchangeCounterpartiesFromDto,
  exchangeFromDto,
  exchangesFromDto,
  exchangeWithdrawalFeeFromDto,
} from '../mappers/exchange';
import type {
  Exchange,
  ExchangeCounterparty,
  ExchangeWithdrawalFee,
  GetWithdrawalFeeOptions,
  ListExchangesOptions,
  ListExchangesResult,
} from '../models/exchange';
import { BaseService } from './base';

/**
 * Service for managing exchange accounts in the Taurus-PROTECT system.
 *
 * This service provides access to exchange account information, including
 * balances, counterparty summaries, and withdrawal fee calculations.
 *
 * @example
 * ```typescript
 * // List exchange accounts
 * const result = await exchangeService.list();
 * for (const exchange of result.items) {
 *   console.log(`${exchange.exchange}: ${exchange.totalBalance}`);
 * }
 *
 * // Get a specific exchange account
 * const exchange = await exchangeService.get('exchange-123');
 * console.log(`Status: ${exchange.status}`);
 *
 * // Get all counterparties (summary by exchange)
 * const counterparties = await exchangeService.getCounterparties();
 * for (const cp of counterparties) {
 *   console.log(`${cp.name}: ${cp.baseCurrencyValuation}`);
 * }
 *
 * // Calculate withdrawal fee
 * const fee = await exchangeService.getWithdrawalFee('exchange-123', {
 *   toAddressId: 'address-456',
 *   amount: '1000000000',
 * });
 * console.log(`Fee: ${fee.fee}`);
 * ```
 */
export class ExchangeService extends BaseService {
  private readonly exchangeApi: ExchangeApi;

  /**
   * Creates a new ExchangeService instance.
   *
   * @param exchangeApi - The ExchangeApi instance from the OpenAPI client
   */
  constructor(exchangeApi: ExchangeApi) {
    super();
    this.exchangeApi = exchangeApi;
  }

  /**
   * Lists exchange accounts.
   *
   * @param options - Optional filtering and pagination options
   * @returns Result containing exchange accounts and pagination cursor
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all exchange accounts
   * const result = await exchangeService.list();
   *
   * // List with filtering
   * const filtered = await exchangeService.list({
   *   currencyId: 'eth',
   *   onlyPositiveBalance: true,
   *   includeBaseCurrencyValuation: true,
   * });
   *
   * // Paginate through results
   * let cursor: string | undefined;
   * do {
   *   const page = await exchangeService.list({
   *     pageSize: 50,
   *     currentPage: cursor,
   *   });
   *   for (const exchange of page.items) {
   *     console.log(exchange.id);
   *   }
   *   cursor = page.nextCursor;
   * } while (cursor);
   * ```
   */
  async list(options?: ListExchangesOptions): Promise<ListExchangesResult> {
    return this.execute(async () => {
      const response = await this.exchangeApi.exchangeServiceGetExchanges({
        currencyID: options?.currencyId,
        includeBaseCurrencyValuation: options?.includeBaseCurrencyValuation,
        exchangeLabel: options?.exchangeLabel,
        sortOrder: options?.sortOrder,
        status: options?.status,
        onlyPositiveBalance: options?.onlyPositiveBalance,
        cursorPageSize: options?.pageSize?.toString(),
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
      });

      const result =
        (response as Record<string, unknown>).result ??
        (response as Record<string, unknown>).exchanges;
      const items = exchangesFromDto(result as unknown[]);

      // Extract next cursor from response
      const cursor = (response as Record<string, unknown>).cursor as
        | Record<string, unknown>
        | undefined;
      const nextCursor = cursor?.nextPage as string | undefined;

      return {
        items,
        nextCursor,
      };
    });
  }

  /**
   * Gets an exchange account by ID.
   *
   * @param id - The exchange account ID
   * @returns The exchange account
   * @throws {@link ValidationError} If id is empty
   * @throws {@link NotFoundError} If exchange account is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const exchange = await exchangeService.get('exchange-123');
   * console.log(`Exchange: ${exchange.exchange}`);
   * console.log(`Balance: ${exchange.totalBalance}`);
   * console.log(`Status: ${exchange.status}`);
   * ```
   */
  async get(id: string): Promise<Exchange> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response = await this.exchangeApi.exchangeServiceGetExchange({ id });

      const result =
        (response as Record<string, unknown>).result ??
        (response as Record<string, unknown>).exchange;
      const exchange = exchangeFromDto(result);

      if (!exchange) {
        throw new NotFoundError(`Exchange account with id '${id}' not found`);
      }

      return exchange;
    });
  }

  /**
   * Gets all exchange counterparties.
   *
   * Returns a summary of holdings grouped by exchange name.
   *
   * @returns List of exchange counterparties
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const counterparties = await exchangeService.getCounterparties();
   * for (const cp of counterparties) {
   *   console.log(`${cp.name}: ${cp.baseCurrencyValuation}`);
   * }
   * ```
   */
  async getCounterparties(): Promise<ExchangeCounterparty[]> {
    return this.execute(async () => {
      const response = await this.exchangeApi.exchangeServiceGetExchangeCounterparties();

      const result =
        (response as Record<string, unknown>).exchanges ??
        (response as Record<string, unknown>).result;
      return exchangeCounterpartiesFromDto(result as unknown[]);
    });
  }

  /**
   * Gets the withdrawal fee for a transfer from an exchange.
   *
   * Returns the fee that will be charged for withdrawing assets from
   * an exchange to an external address. An empty/undefined fee means
   * the exchange does not provide a live estimation of its fees.
   *
   * @param exchangeId - The exchange account ID
   * @param options - Optional destination address and amount
   * @returns The withdrawal fee, or undefined if not available
   * @throws {@link ValidationError} If exchangeId is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get fee with specific amount
   * const fee = await exchangeService.getWithdrawalFee('exchange-123', {
   *   toAddressId: 'address-456',
   *   amount: '1000000000000000000', // 1 ETH in Wei
   * });
   *
   * if (fee) {
   *   console.log(`Withdrawal fee: ${fee.fee}`);
   * } else {
   *   console.log('Fee estimation not available');
   * }
   * ```
   */
  async getWithdrawalFee(
    exchangeId: string,
    options?: GetWithdrawalFeeOptions
  ): Promise<ExchangeWithdrawalFee | undefined> {
    if (!exchangeId || exchangeId.trim() === '') {
      throw new ValidationError('exchangeId is required');
    }

    return this.execute(async () => {
      const response = await this.exchangeApi.exchangeServiceGetExchangeWithdrawalFee({
        id: exchangeId,
        toAddressId: options?.toAddressId,
        amount: options?.amount,
      });

      return exchangeWithdrawalFeeFromDto(response);
    });
  }

  /**
   * Exports all exchange accounts to a specified format.
   *
   * @param format - The export format ("csv" or "json")
   * @returns The exported data as a string
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Export as CSV
   * const csvData = await exchangeService.export('csv');
   * console.log(csvData);
   *
   * // Export as JSON
   * const jsonData = await exchangeService.export('json');
   * const parsed = JSON.parse(jsonData);
   * ```
   */
  async export(format?: string): Promise<string> {
    return this.execute(async () => {
      const response = await this.exchangeApi.exchangeServiceExportExchanges({
        format,
      });

      const result =
        (response as Record<string, unknown>).result ??
        (response as Record<string, unknown>).data;
      return (result as string) ?? '';
    });
  }
}
