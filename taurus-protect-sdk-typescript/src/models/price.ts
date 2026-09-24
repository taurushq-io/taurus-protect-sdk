/**
 * Price models for Taurus-PROTECT SDK.
 *
 * Provides domain models for cryptocurrency prices, price history,
 * and currency conversion operations.
 */

import type { Currency } from './currency';
import type { CursorPage, CursorPageOptions } from './pagination';

/**
 * Represents a currency price/exchange rate.
 *
 * Prices track exchange rates between currencies (e.g., ETH to USD) along with
 * metadata about the rate source, precision, and price changes. Prices are
 * used for portfolio valuation, transaction display, and reporting.
 */
export interface PriceSignature {
  /** The signing user's ID */
  readonly userId?: string;
  /** Base64-encoded ECDSA signature over the canonical price form */
  readonly signature?: string;
}

export interface Price {
  /** The blockchain this price applies to (e.g., "ethereum", "bitcoin") */
  readonly blockchain?: string;
  /** The source currency code being converted from (e.g., "ETH", "BTC") */
  readonly currencyFrom?: string;
  /** The target currency code being converted to (e.g., "USD", "EUR") */
  readonly currencyTo?: string;
  /** The number of decimal places for the rate precision */
  readonly decimals?: string;
  /** The exchange rate value (amount of currencyTo per unit of currencyFrom) */
  readonly rate?: string;
  /** PRICEUPDATER signatures over the canonical price form */
  readonly signatures?: PriceSignature[];
  /** The percentage price change over the last 24 hours */
  readonly changePercent24Hour?: string;
  /** The data source for this price (e.g., "coingecko", "cryptocompare") */
  readonly source?: string;
  /** Timestamp when this price record was created */
  readonly createdAt?: Date;
  /** Timestamp when this price record was last updated */
  readonly updatedAt?: Date;
  /** Detailed information about the source currency */
  readonly currencyFromInfo?: Currency;
  /** Detailed information about the target currency */
  readonly currencyToInfo?: Currency;
}

/**
 * Represents a single data point in a price history timeline (OHLCV candlestick data).
 *
 * Price history points capture OHLCV (Open, High, Low, Close, Volume) data for a specific
 * time period, enabling historical price analysis and charting. This follows standard
 * financial candlestick data format used in trading and portfolio analysis.
 */
export interface PriceHistoryPoint {
  /** The start timestamp of the time period this data point covers */
  readonly periodStartDate?: Date;
  /** The blockchain this price data applies to (e.g., "ethereum", "bitcoin") */
  readonly blockchain?: string;
  /** The source currency code being converted from (e.g., "ETH", "BTC") */
  readonly currencyFrom?: string;
  /** The target currency code being converted to (e.g., "USD", "EUR") */
  readonly currencyTo?: string;
  /** The highest price during this time period */
  readonly high?: string;
  /** The lowest price during this time period */
  readonly low?: string;
  /** The opening price at the start of this time period */
  readonly open?: string;
  /** The closing price at the end of this time period */
  readonly close?: string;
  /** The trading volume in the source currency during this period */
  readonly volumeFrom?: string;
  /** The trading volume in the target currency during this period */
  readonly volumeTo?: string;
  /** The percentage price change during this time period */
  readonly changePercent?: string;
  /** Detailed information about the source currency */
  readonly currencyFromInfo?: Currency;
  /** Detailed information about the target currency */
  readonly currencyToInfo?: Currency;
}

/**
 * Represents the result of a currency conversion operation.
 *
 * Conversion results contain both the raw value and a human-readable main unit value
 * (accounting for decimals), along with currency metadata. This is typically used
 * when displaying converted amounts in the UI or for reporting purposes.
 */
export interface ConversionResult {
  /** The currency symbol for display (e.g., "USD", "EUR", "ETH") */
  readonly symbol?: string;
  /** The converted value in the smallest unit (e.g., wei for ETH, satoshi for BTC) */
  readonly value?: string;
  /** The converted value in main units (e.g., ETH instead of wei), formatted for display */
  readonly mainUnitValue?: string;
  /** Detailed information about the target currency */
  readonly currencyInfo?: Currency;
}

/**
 * Options for getting price history.
 */
export interface GetPriceHistoryOptions {
  /** The base currency (e.g., "ETH", "BTC") */
  base: string;
  /** The quote currency (e.g., "USD", "EUR") */
  quote: string;
  /** Number of daily points to return, newest first: 1-365 (default 20) */
  limit?: number;
}

/**
 * Options for listing prices. A cursor list: `pageSize` 1-100 (default 20) and the
 * `cursor` of a previous page.
 *
 * `fromCurrencyId` alone keeps the prices from that currency, `toCurrencyIds` alone the
 * prices into those currencies, and both together the prices from the one into the
 * others.
 */
export interface ListPricesOptions extends CursorPageOptions {
  /** Keep only the primary price of each currency pair */
  readonly onlyPrimary?: boolean;
  /** Sort order ("ASC" or "DESC") */
  readonly sortOrder?: string;
  /** Keep prices whose source currency is this currency ID */
  readonly fromCurrencyId?: string;
  /** Keep prices whose target currency is one of these currency IDs */
  readonly toCurrencyIds?: string[];
}

/**
 * A page of verified prices.
 */
export interface ListPricesResult {
  /** The prices of this page, each verified against its signature */
  readonly items: Price[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}

/**
 * Options for currency conversion.
 */
export interface ConvertOptions {
  /** The source currency (e.g., "ETH", "BTC") */
  currency: string;
  /** The amount to convert (as a string to preserve precision) */
  amount: string;
  /** Optional list of currency symbols to convert to */
  symbols?: string[];
  /** Optional list of target currency IDs to convert to */
  targetCurrencyIds?: string[];
}
