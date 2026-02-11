/**
 * Reservation models for Taurus-PROTECT SDK.
 *
 * Reservations are used to lock specific UTXOs (Unspent Transaction Outputs)
 * for UTXO-based blockchains like Bitcoin and Litecoin, preventing
 * double-spending during transaction creation.
 */

import type { Currency } from './currency';

/**
 * Represents a UTXO reservation in the Taurus-PROTECT system.
 *
 * Reservations lock specific UTXOs to prevent double-spending during
 * transaction creation on UTXO-based blockchains.
 */
export interface Reservation {
  /** Unique reservation identifier */
  readonly id?: string;
  /** Reserved amount in the currency's smallest unit */
  readonly amount?: string;
  /** Date when the reservation was created */
  readonly creationDate?: Date;
  /** The kind/type of reservation */
  readonly kind?: string;
  /** Optional comment about the reservation */
  readonly comment?: string;
  /** Internal address ID in Taurus-PROTECT */
  readonly addressId?: string;
  /** Blockchain address string/hash */
  readonly address?: string;
  /** Currency information for the reserved funds */
  readonly currencyInfo?: Currency;
  /** Resource ID associated with the reservation */
  readonly resourceId?: string;
  /** Resource type associated with the reservation */
  readonly resourceType?: string;
}

/**
 * Represents UTXO details for a reservation.
 *
 * Contains the UTXO (Unspent Transaction Output) data associated with a reservation.
 */
export interface ReservationUtxo {
  /** Unique UTXO identifier */
  readonly id?: string;
  /** Transaction hash */
  readonly hash?: string;
  /** Output index in the transaction */
  readonly outputIndex?: number;
  /** Locking script */
  readonly script?: string;
  /** Value in the smallest unit */
  readonly value?: string;
  /** Block height where the UTXO was created */
  readonly blockHeight?: string;
  /**
   * Request ID that reserved this UTXO.
   * @deprecated Replaced by reservationId since version 3.22
   */
  readonly reservedByRequestId?: string;
  /** Reservation ID */
  readonly reservationId?: string;
  /** Human-readable value string */
  readonly valueString?: string;
}

/**
 * Options for listing reservations.
 */
export interface ListReservationsOptions {
  /**
   * Filter by kind of reservation.
   * @deprecated Since 3.20: use 'kinds' instead
   */
  kind?: string;
  /** Filter by blockchain address */
  address?: string;
  /** Filter by internal address ID */
  addressId?: string;
  /** Filter by multiple kinds of reservation (takes precedence over 'kind') */
  kinds?: string[];
  /** Base64-encoded cursor for pagination */
  cursorCurrentPage?: string;
  /** Page request direction: FIRST, PREVIOUS, NEXT, or LAST */
  cursorPageRequest?: string;
  /** Page size */
  cursorPageSize?: string;
}
