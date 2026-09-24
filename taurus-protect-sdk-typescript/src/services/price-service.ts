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
import type { TgvalidatordQueryPricesV2Request } from '../internal/openapi/models/TgvalidatordQueryPricesV2Request';
import { buildCursorPage, cursorRequest, resolvePageSize } from '../models/pagination';
import type {
  ConversionResult,
  ConvertOptions,
  GetPriceHistoryOptions,
  ListPricesOptions,
  ListPricesResult,
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
 * // Get the current prices, page by page
 * let cursor: string | undefined;
 * do {
 *   const page = await priceService.list({ onlyPrimary: true, cursor });
 *   for (const price of page.items) {
 *     console.log(`${price.currencyFrom}/${price.currencyTo}: ${price.rate}`);
 *   }
 *   cursor = page.pagination.hasMore ? page.pagination.nextCursor : undefined;
 * } while (cursor);
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
   * Lists a page of current prices, each verified against its signature.
   *
   * Served by `QueryPricesV2`, a cursor list. `fromCurrencyId` alone filters on the source
   * currency, `toCurrencyIds` alone on the target currencies, and both together on the
   * pair; neither lists every price.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of verified prices and its cursor pagination
   * @throws {@link ValidationError} If the page size is out of bounds
   * @throws {@link IntegrityError} If a price signature does not verify
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const page = await priceService.list({
   *   fromCurrencyId: 'ETH',
   *   toCurrencyIds: ['USD', 'CHF'],
   *   onlyPrimary: true,
   * });
   * for (const price of page.items) {
   *   console.log(`${price.currencyFrom}/${price.currencyTo}: ${price.rate}`);
   * }
   * ```
   */
  async list(options?: ListPricesOptions): Promise<ListPricesResult> {
    const page = cursorRequest(options);
    const body: TgvalidatordQueryPricesV2Request = {
      onlyPrimary: options?.onlyPrimary,
      sortOrder: options?.sortOrder,
      cursor: page.cursor,
    };
    const from = options?.fromCurrencyId;
    const to = options?.toCurrencyIds;
    if (from !== undefined && to !== undefined) {
      body.fromTo = { currencyFromId: from, currencyToIds: to };
    } else if (from !== undefined) {
      body.from = { currencyFromId: from };
    } else if (to !== undefined) {
      body.to = { currencyToIds: to };
    }

    return this.execute(async () => {
      const response = await this.pricesApi.priceServiceQueryPricesV2({ body });

      // Same signed CurrencyPrice rows as the v1 list, so the same verification.
      const prices = pricesFromDto(response.result);
      if (prices.length > 0) {
        verifyPrices(prices, await this.rulesCache.get());
      }

      return {
        items: prices,
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }

  /**
   * Gets price history for a currency pair.
   *
   * Returns OHLCV (Open, High, Low, Close, Volume) candlestick data
   * for historical price analysis and charting.
   *
   * The history cannot page: `limit` (1-365, default 20) is the number of daily points,
   * newest first, and is always sent.
   *
   * @param options - History options including base currency, quote currency, and limit
   * @returns Array of price history points
   * @throws {@link ValidationError} If required parameters are missing or the limit is out
   *   of bounds
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
    // Cannot page: the limit is the whole request, bounded by a year of daily points.
    const limit = resolvePageSize(options.limit, 'limit', 'price_history');

    return this.execute(async () => {
      const response = await this.pricesApi.priceServiceGetPricesHistory({
        base: options.base,
        quote: options.quote,
        limit: String(limit),
      });

      return priceHistoryPointsFromDto(response.result);
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

      return conversionResultsFromDto(response.result);
    });
  }
}
