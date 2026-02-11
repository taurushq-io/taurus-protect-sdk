/**
 * Fee models for Taurus-PROTECT SDK.
 *
 * Provides models for network fee information retrieved from the API.
 */

import type { Currency } from './currency';

/**
 * Represents a network fee entry (v1 API format).
 *
 * Fees are represented as key-value pairs where the key is typically
 * the blockchain/currency identifier and the value is the fee amount.
 *
 * @deprecated Use FeeV2 for richer fee information
 */
export interface Fee {
  /** Fee key (blockchain/currency identifier) */
  readonly key?: string;
  /** Fee value (amount) */
  readonly value?: string;
}

/**
 * Represents a network fee entry with currency information (v2 API format).
 *
 * Provides detailed fee information including currency metadata.
 */
export interface FeeV2 {
  /** Currency identifier */
  readonly currencyId?: string;
  /** Fee value (amount) */
  readonly value?: string;
  /** Denomination (e.g., "wei", "satoshi") */
  readonly denom?: string;
  /** Currency information */
  readonly currencyInfo?: Currency;
  /** Last update timestamp */
  readonly updateDate?: Date;
}
