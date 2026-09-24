/**
 * Fee models for Taurus-PROTECT SDK.
 *
 * Provides models for network fee information retrieved from the API.
 */

import type { Currency } from './currency';

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
