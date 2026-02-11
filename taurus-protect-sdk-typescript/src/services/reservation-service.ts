/**
 * Reservation service for Taurus-PROTECT SDK.
 *
 * Provides methods for managing UTXO reservations used to lock specific UTXOs
 * (Unspent Transaction Outputs) for UTXO-based blockchains like Bitcoin and
 * Litecoin, preventing double-spending during transaction creation.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { ReservationsApi } from '../internal/openapi/apis/ReservationsApi';
import {
  reservationFromDto,
  reservationsFromDto,
  reservationUtxoFromDto,
} from '../mappers/reservation';
import type {
  ListReservationsOptions,
  Reservation,
  ReservationUtxo,
} from '../models/reservation';
import { BaseService } from './base';

/**
 * Service for managing UTXO reservations.
 *
 * Reservations are used to lock specific UTXOs for UTXO-based blockchains
 * like Bitcoin and Litecoin, preventing double-spending during transaction
 * creation.
 *
 * @example
 * ```typescript
 * // List all reservations
 * const reservations = await reservationService.list();
 * for (const reservation of reservations) {
 *   console.log(`${reservation.id}: ${reservation.amount}`);
 * }
 *
 * // Get a specific reservation
 * const reservation = await reservationService.get('res-123');
 * console.log(`Address: ${reservation.address}`);
 *
 * // Get UTXO details for a reservation
 * const utxo = await reservationService.getUtxo('res-123');
 * console.log(`UTXO hash: ${utxo.hash}`);
 *
 * // List reservations with filters
 * const filtered = await reservationService.list({
 *   addressId: 'addr-456',
 *   kinds: ['OUTGOING', 'CONSOLIDATION'],
 * });
 * ```
 */
export class ReservationService extends BaseService {
  private readonly reservationsApi: ReservationsApi;

  /**
   * Creates a new ReservationService instance.
   *
   * @param reservationsApi - The ReservationsApi instance from the OpenAPI client
   */
  constructor(reservationsApi: ReservationsApi) {
    super();
    this.reservationsApi = reservationsApi;
  }

  /**
   * Lists all reservations with optional filtering.
   *
   * Results are sorted by the most recent reservations first.
   * Default limit is 100 reservations.
   *
   * @param options - Optional filtering and pagination options
   * @returns Array of reservations
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all reservations
   * const all = await reservationService.list();
   *
   * // Filter by address ID
   * const byAddress = await reservationService.list({ addressId: 'addr-123' });
   *
   * // Filter by multiple kinds
   * const byKinds = await reservationService.list({
   *   kinds: ['OUTGOING', 'CONSOLIDATION'],
   * });
   * ```
   */
  async list(options?: ListReservationsOptions): Promise<Reservation[]> {
    return this.execute(async () => {
      const response = await this.reservationsApi.walletServiceGetReservations({
        kind: options?.kind,
        address: options?.address,
        addressId: options?.addressId,
        kinds: options?.kinds,
        cursorCurrentPage: options?.cursorCurrentPage,
        cursorPageRequest: options?.cursorPageRequest,
        cursorPageSize: options?.cursorPageSize,
      });

      return reservationsFromDto(response.result);
    });
  }

  /**
   * Gets a reservation by ID.
   *
   * @param id - The reservation ID
   * @returns The reservation
   * @throws {@link ValidationError} If id is empty
   * @throws {@link NotFoundError} If reservation is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const reservation = await reservationService.get('res-123');
   * console.log(`Amount: ${reservation.amount}`);
   * console.log(`Kind: ${reservation.kind}`);
   * ```
   */
  async get(id: string): Promise<Reservation> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response = await this.reservationsApi.walletServiceGetReservation({
        id,
      });

      const reservation = reservationFromDto(response.result);
      if (!reservation) {
        throw new NotFoundError(`Reservation with id '${id}' not found`);
      }

      return reservation;
    });
  }

  /**
   * Gets the UTXO details for a reservation.
   *
   * Returns the UTXO (Unspent Transaction Output) associated with a reservation,
   * if any. This is relevant for UTXO-based blockchains like Bitcoin.
   *
   * @param id - The reservation ID
   * @returns The UTXO details
   * @throws {@link ValidationError} If id is empty
   * @throws {@link NotFoundError} If reservation or UTXO is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const utxo = await reservationService.getUtxo('res-123');
   * console.log(`Hash: ${utxo.hash}`);
   * console.log(`Value: ${utxo.value}`);
   * console.log(`Block height: ${utxo.blockHeight}`);
   * ```
   */
  async getUtxo(id: string): Promise<ReservationUtxo> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response = await this.reservationsApi.walletServiceGetReservationUTXO({
        id,
      });

      const utxo = reservationUtxoFromDto(response.result);
      if (!utxo) {
        throw new NotFoundError(`UTXO for reservation with id '${id}' not found`);
      }

      return utxo;
    });
  }
}
