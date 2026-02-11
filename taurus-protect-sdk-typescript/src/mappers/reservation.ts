/**
 * Reservation mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Reservation, ReservationUtxo } from '../models/reservation';
import { safeDate, safeInt, safeMap, safeString } from './base';
import { currencyFromDto } from './currency';

/**
 * Maps a reservation DTO to a Reservation domain model.
 *
 * @param dto - The DTO from the OpenAPI client
 * @returns The Reservation domain model, or undefined if dto is null/undefined
 */
export function reservationFromDto(dto: unknown): Reservation | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    amount: safeString(d.amount),
    creationDate: safeDate(d.creationDate ?? d.creation_date),
    kind: safeString(d.kind),
    comment: safeString(d.comment),
    // Note: OpenAPI DTO uses 'addressid' (lowercase 'i'), map to 'addressId'
    addressId: safeString(d.addressid ?? d.addressId ?? d.address_id),
    address: safeString(d.address),
    currencyInfo: currencyFromDto(d.currencyInfo ?? d.currency_info),
    resourceId: safeString(d.resourceId ?? d.resource_id),
    resourceType: safeString(d.resourceType ?? d.resource_type),
  };
}

/**
 * Maps an array of reservation DTOs to Reservation domain models.
 *
 * @param dtos - Array of DTOs from the OpenAPI client
 * @returns Array of Reservation domain models
 */
export function reservationsFromDto(dtos: unknown[] | null | undefined): Reservation[] {
  return safeMap(dtos, reservationFromDto);
}

/**
 * Maps a UTXO DTO to a ReservationUtxo domain model.
 *
 * @param dto - The DTO from the OpenAPI client
 * @returns The ReservationUtxo domain model, or undefined if dto is null/undefined
 */
export function reservationUtxoFromDto(dto: unknown): ReservationUtxo | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    hash: safeString(d.hash),
    outputIndex: safeInt(d.outputIndex ?? d.output_index),
    script: safeString(d.script),
    value: safeString(d.value),
    blockHeight: safeString(d.blockHeight ?? d.block_height),
    reservedByRequestId: safeString(d.reservedByRequestId ?? d.reserved_by_request_id),
    reservationId: safeString(d.reservationId ?? d.reservation_id),
    valueString: safeString(d.valueString ?? d.value_string),
  };
}
