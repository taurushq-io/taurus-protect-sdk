/**
 * Fee mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Fee, FeeV2 } from '../models/fee';
import { currencyFromDto } from './currency';
import { safeDate, safeMap, safeString } from './base';

/**
 * Maps a fee DTO (v1 format) to a Fee domain model.
 *
 * @param dto - The OpenAPI KeyValue DTO
 * @returns The Fee domain model, or undefined if dto is null/undefined
 */
export function feeFromDto(dto: unknown): Fee | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    key: safeString(d.key),
    value: safeString(d.value),
  };
}

/**
 * Maps an array of fee DTOs (v1 format) to Fee domain models.
 *
 * @param dtos - Array of OpenAPI KeyValue DTOs
 * @returns Array of Fee domain models
 */
export function feesFromDto(dtos: unknown[] | null | undefined): Fee[] {
  return safeMap(dtos, feeFromDto);
}

/**
 * Maps a fee DTO (v2 format) to a FeeV2 domain model.
 *
 * @param dto - The OpenAPI TgvalidatordFee DTO
 * @returns The FeeV2 domain model, or undefined if dto is null/undefined
 */
export function feeV2FromDto(dto: unknown): FeeV2 | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    currencyId: safeString(d.currencyId ?? d.currency_id),
    value: safeString(d.value),
    denom: safeString(d.denom),
    currencyInfo: currencyFromDto(d.currencyInfo ?? d.currency_info),
    updateDate: safeDate(d.updateDate ?? d.update_date),
  };
}

/**
 * Maps an array of fee DTOs (v2 format) to FeeV2 domain models.
 *
 * @param dtos - Array of OpenAPI TgvalidatordFee DTOs
 * @returns Array of FeeV2 domain models
 */
export function feesV2FromDto(dtos: unknown[] | null | undefined): FeeV2[] {
  return safeMap(dtos, feeV2FromDto);
}
