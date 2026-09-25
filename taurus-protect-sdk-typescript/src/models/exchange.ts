/**
 * Exchange models for Taurus-PROTECT SDK.
 *
 * Represents exchange account information in the Taurus-PROTECT system.
 */

import type { Currency } from './currency';
import type { CursorNavigationOptions, CursorPage } from './pagination';

/**
 * Represents an exchange account in the Taurus-PROTECT system.
 *
 * An exchange represents a connection to a third-party exchange (e.g., Binance,
 * Coinbase) that can be used for trading and transfers.
 */
export interface Exchange {
  /** Unique identifier */
  readonly id?: string;
  /** Exchange name (e.g., "binance", "coinbase") */
  readonly exchange?: string;
  /** Account identifier on the exchange */
  readonly account?: string;
  /** Currency code */
  readonly currency?: string;
  /** Account type */
  readonly type?: string;
  /** Total balance in the smallest currency unit */
  readonly totalBalance?: string;
  /** Account status */
  readonly status?: string;
  /** Container */
  readonly container?: string;
  /** Label */
  readonly label?: string;
  /** Display label */
  readonly displayLabel?: string;
  /** Valuation in the base currency (CHF, EUR, USD, etc.) */
  readonly baseCurrencyValuation?: string;
  /** Whether the exchange has a whitelisted address */
  readonly hasWLA?: boolean;
  /** Currency information */
  readonly currencyInfo?: Currency;
  /** Creation date */
  readonly creationDate?: Date;
  /** Last update date */
  readonly updateDate?: Date;
}

/**
 * Represents an exchange counterparty summary.
 *
 * A counterparty represents a grouped view of exchange holdings by exchange name,
 * with the total valuation across all currencies.
 */
export interface ExchangeCounterparty {
  /** Exchange name */
  readonly name?: string;
  /** Total valuation in the base currency (CHF, EUR, USD, etc.) */
  readonly baseCurrencyValuation?: string;
}

/**
 * Represents the withdrawal fee for an exchange transfer.
 *
 * Contains the fee amount that will be charged for withdrawing
 * assets from an exchange to an external address.
 */
export interface ExchangeWithdrawalFee {
  /** Withdrawal fee amount in the smallest currency unit */
  readonly fee?: string;
}

/**
 * Options for listing exchange accounts. A cursor list: `pageSize` 1-100 (default 20) and
 * the `cursor` of a previous page.
 */
export interface ListExchangesOptions extends CursorNavigationOptions {
  /** Filter on currency ID */
  currencyId?: string;
  /** Include base currency valuation in response */
  includeBaseCurrencyValuation?: boolean;
  /** Filter by exchange label */
  exchangeLabel?: string;
  /** Sort order: "ASC" or "DESC" (default: "DESC") */
  sortOrder?: 'ASC' | 'DESC';
  /** Filter by status (leave empty for all statuses) */
  status?: string;
  /** Exclude exchange accounts with zero balance */
  onlyPositiveBalance?: boolean;
}

/**
 * A page of exchange accounts.
 */
export interface ListExchangesResult {
  /** The exchange accounts of this page */
  items: Exchange[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  pagination: CursorPage;
}

/**
 * Options for getting withdrawal fee.
 */
export interface GetWithdrawalFeeOptions {
  /** Destination address ID (optional) */
  toAddressId?: string;
  /** Amount to withdraw in the smallest currency unit (optional) */
  amount?: string;
}
