/**
 * Price service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving cryptocurrency prices, price history,
 * and performing currency conversions.
 */

import type { RulesContainerCache } from '../cache';
import { ConfigurationError, ValidationError } from '../errors';
import { verifyPrices } from '../helpers/price-verifier';
import type { PricesApi } from '../internal/openapi/apis/PricesApi';
import {
  conversionResultsFromDto,
  priceHistoryPointsFromDto,
  pricesFromDto,
} from '../mappers/price';
import type {
  ConversionResult,
  ConvertOptions,
  GetPriceHistoryOptions,
  Price,
  PriceHistoryPoint,
} from '../models/price';
import { BaseService } from './base';

/**
 * Service for retrieving cryptocurrency prices and performing conversions.
 *
 * This service provides operations for querying current and historical prices
 * of supported cryptocurrencies, as well as converting amounts between currencies.
 * Prices are provided against the configured base currency (typically USD).
 *
 * @example
 * ```typescript
 * // Get all current prices
 * const prices = await priceService.list();
 * for (const price of prices) {
 *   console.log(`${price.currencyFrom}/${price.currencyTo}: ${price.rate}`);
 * }
 *
 * // Get price history for a currency pair
 * const history = await priceService.getHistory({
 *   base: 'ETH',
 *   quote: 'USD',
 *   limit: 100,
 * });
 *
 * // Convert an amount to target currencies
 * const converted = await priceService.convert({
 *   currency: 'ETH',
 *   amount: '1000000000000000000',
 *   targetCurrencyIds: ['USD', 'BTC'],
 * });
 * ```
 */
export class PriceService extends BaseService {
  private readonly pricesApi: PricesApi;
  private readonly rulesCache: RulesContainerCache;

  /**
   * Creates a new PriceService instance.
   *
   * Price signature verification is MANDATORY: rate and decimals feed amount
   * conversion, so an unverified price is a wrong number a caller acts on. Whether
   * prices must be signed is decided by the SuperAdmin-verified rules container.
   *
   * @param pricesApi - The PricesApi instance from the OpenAPI client
   * @param rulesCache - Rules container cache supplying the PRICEUPDATER keys
   * @throws ConfigurationError if rulesCache is not provided
   */
  constructor(pricesApi: PricesApi, rulesCache: RulesContainerCache) {
    super();
    if (!rulesCache) {
      throw new ConfigurationError(
        "RulesContainerCache is required for PriceService — price signature verification is mandatory"
      );
    }
    this.pricesApi = pricesApi;
    this.rulesCache = rulesCache;
  }

  /**
   * Lists all current prices.
   *
   * Returns the current exchange rates for all supported currency pairs.
   *
   * @returns Array of current prices
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const prices = await priceService.list();
   * for (const price of prices) {
   *   console.log(`${price.currencyFrom}/${price.currencyTo}: ${price.rate}`);
   * }
   * ```
   */
  async list(): Promise<Price[]> {
    return this.execute(async () => {
      const response = await this.pricesApi.priceServiceGetPrices();

      const result =
        (response as Record<string, unknown>).result ??
        (response as Record<string, unknown>).prices;
      const prices = pricesFromDto(result as unknown[]);

      if (prices.length > 0) {
        verifyPrices(prices, await this.rulesCache.get());
      }

      return prices;
    });
  }

  /**
   * Gets price history for a currency pair.
   *
   * Returns OHLCV (Open, High, Low, Close, Volume) candlestick data
   * for historical price analysis and charting.
   *
   * @param options - History options including base currency, quote currency, and limit
   * @returns Array of price history points
   * @throws {@link ValidationError} If required parameters are missing
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const history = await priceService.getHistory({
   *   base: 'ETH',
   *   quote: 'USD',
   *   limit: 100,
   * });
   * for (const point of history) {
   *   console.log(`${point.periodStartDate}: Open=${point.open} Close=${point.close}`);
   * }
   * ```
   */
  async getHistory(options: GetPriceHistoryOptions): Promise<PriceHistoryPoint[]> {
    if (!options.base || options.base.trim() === '') {
      throw new ValidationError('base is required');
    }
    if (!options.quote || options.quote.trim() === '') {
      throw new ValidationError('quote is required');
    }
    if (options.limit !== undefined && options.limit <= 0) {
      throw new ValidationError('limit must be positive');
    }

    return this.execute(async () => {
      const response = await this.pricesApi.priceServiceGetPricesHistory({
        base: options.base,
        quote: options.quote,
        limit: options.limit?.toString(),
      });

      const result =
        (response as Record<string, unknown>).result ??
        (response as Record<string, unknown>).points;
      return priceHistoryPointsFromDto(result as unknown[]);
    });
  }

  /**
   * Converts an amount from one currency to target currencies.
   *
   * Takes an amount in the source currency and returns the converted
   * values in the specified target currencies.
   *
   * @param options - Conversion options including currency, amount, and target currencies
   * @returns Array of conversion results
   * @throws {@link ValidationError} If required parameters are missing
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Convert 1 ETH (in wei) to USD and BTC
   * const results = await priceService.convert({
   *   currency: 'ETH',
   *   amount: '1000000000000000000',
   *   targetCurrencyIds: ['USD', 'BTC'],
   * });
   * for (const result of results) {
   *   console.log(`${result.symbol}: ${result.mainUnitValue}`);
   * }
   * ```
   */
  async convert(options: ConvertOptions): Promise<ConversionResult[]> {
    if (!options.currency || options.currency.trim() === '') {
      throw new ValidationError('currency is required');
    }
    if (!options.amount || options.amount.trim() === '') {
      throw new ValidationError('amount is required');
    }

    return this.execute(async () => {
      const response = await this.pricesApi.priceServiceConvert({
        currency: options.currency,
        amount: options.amount,
        symbols: options.symbols,
        targetCurrencyIds: options.targetCurrencyIds,
      });

      const result =
        (response as Record<string, unknown>).result ??
        (response as Record<string, unknown>).values;
      return conversionResultsFromDto(result as unknown[]);
    });
  }
}
