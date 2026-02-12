/**
 * Currency service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving currency information.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { CurrenciesApi } from '../internal/openapi/apis/CurrenciesApi';
import { currenciesFromDto, currencyFromDto } from '../mappers/currency';
import type {
  Currency,
  GetCurrencyByBlockchainOptions,
  ListCurrenciesOptions,
} from '../models/currency';
import { BaseService } from './base';

/**
 * Service for retrieving currency information.
 *
 * Provides methods to list all currencies, get specific currencies,
 * and retrieve the base currency configured for the tenant.
 *
 * @example
 * ```typescript
 * // List all enabled currencies
 * const currencies = await currencyService.list();
 * for (const currency of currencies) {
 *   console.log(`${currency.symbol}: ${currency.name}`);
 * }
 *
 * // Get currency by blockchain and network
 * const eth = await currencyService.getByBlockchain({
 *   blockchain: 'ETH',
 *   network: 'mainnet',
 * });
 * console.log(`${eth.name} has ${eth.decimals} decimals`);
 *
 * // Get the base currency (e.g., USD, EUR, CHF)
 * const base = await currencyService.getBaseCurrency();
 * console.log(`Base currency: ${base.symbol}`);
 * ```
 */
export class CurrencyService extends BaseService {
  private readonly currenciesApi: CurrenciesApi;

  /**
   * Creates a new CurrencyService instance.
   *
   * @param currenciesApi - The CurrenciesApi instance from the OpenAPI client
   */
  constructor(currenciesApi: CurrenciesApi) {
    super();
    this.currenciesApi = currenciesApi;
  }

  /**
   * Lists all currencies.
   *
   * @param options - Optional filtering options
   * @returns Array of currencies
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List only enabled currencies
   * const currencies = await currencyService.list();
   *
   * // List all currencies including disabled ones
   * const allCurrencies = await currencyService.list({ showDisabled: true });
   *
   * // Include logo URLs
   * const withLogos = await currencyService.list({ includeLogo: true });
   * ```
   */
  async list(options?: ListCurrenciesOptions): Promise<Currency[]> {
    return this.execute(async () => {
      const response = await this.currenciesApi.walletServiceGetCurrencies({
        showDisabled: options?.showDisabled,
        includeLogo: options?.includeLogo,
      });

      const result =
        (response as Record<string, unknown>).currencies ??
        (response as Record<string, unknown>).result;
      return currenciesFromDto(result as unknown[]);
    });
  }

  /**
   * Gets a currency by ID.
   *
   * This method retrieves currency details using the currency ID.
   * For looking up by blockchain and network, use `getByBlockchain()`.
   *
   * @param currencyId - The unique currency identifier
   * @returns The currency
   * @throws {@link ValidationError} If currencyId is empty
   * @throws {@link NotFoundError} If currency is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const currency = await currencyService.get('eth-mainnet');
   * console.log(`Symbol: ${currency.symbol}, Decimals: ${currency.decimals}`);
   * ```
   */
  async get(currencyId: string): Promise<Currency> {
    if (!currencyId || currencyId.trim() === '') {
      throw new ValidationError('currencyId is required');
    }

    return this.execute(async () => {
      // Get all currencies and filter by ID since API doesn't have a direct get-by-id endpoint
      const currencies = await this.list({ showDisabled: true });
      const currency = currencies.find((c) => c.id === currencyId);

      if (!currency) {
        throw new NotFoundError(`Currency with id '${currencyId}' not found`);
      }

      return currency;
    });
  }

  /**
   * Gets a currency by blockchain and network.
   *
   * This retrieves the native currency for a blockchain/network pair,
   * or a specific token if contractAddress is provided.
   *
   * @param options - Blockchain lookup options
   * @returns The currency
   * @throws {@link ValidationError} If required arguments are missing
   * @throws {@link NotFoundError} If currency is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get native currency
   * const eth = await currencyService.getByBlockchain({
   *   blockchain: 'ETH',
   *   network: 'mainnet',
   * });
   *
   * // Get a specific token
   * const usdc = await currencyService.getByBlockchain({
   *   blockchain: 'ETH',
   *   network: 'mainnet',
   *   contractAddress: '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48',
   * });
   * ```
   */
  async getByBlockchain(options: GetCurrencyByBlockchainOptions): Promise<Currency> {
    if (!options.blockchain || options.blockchain.trim() === '') {
      throw new ValidationError('blockchain is required');
    }
    if (!options.network || options.network.trim() === '') {
      throw new ValidationError('network is required');
    }

    return this.execute(async () => {
      const response = await this.currenciesApi.walletServiceGetCurrency({
        uniqueCurrencyFilterBlockchain: options.blockchain,
        uniqueCurrencyFilterNetwork: options.network,
        uniqueCurrencyFilterTokenContractAddress: options.contractAddress,
        uniqueCurrencyFilterTokenID: options.tokenId,
        showDisabled: true,
      });

      const result =
        (response as Record<string, unknown>).currency ??
        (response as Record<string, unknown>).result;
      const currency = currencyFromDto(result);

      if (!currency) {
        throw new NotFoundError(
          `Currency for blockchain '${options.blockchain}' network '${options.network}' not found`
        );
      }

      return currency;
    });
  }

  /**
   * Gets the base currency configured for the tenant.
   *
   * The base currency is used for fiat valuations and is typically
   * configured as CHF, EUR, or USD.
   *
   * @returns The base currency
   * @throws {@link NotFoundError} If no base currency is configured
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const baseCurrency = await currencyService.getBaseCurrency();
   * console.log(`Base currency: ${baseCurrency.symbol}`);
   * ```
   */
  async getBaseCurrency(): Promise<Currency> {
    return this.execute(async () => {
      const response = await this.currenciesApi.walletServiceGetBaseCurrency();

      const result =
        (response as Record<string, unknown>).currency ??
        (response as Record<string, unknown>).result;
      const currency = currencyFromDto(result);

      if (!currency) {
        throw new NotFoundError('Base currency not configured');
      }

      return currency;
    });
  }
}
