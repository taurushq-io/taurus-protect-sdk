/**
 * Score mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Score } from '../models/score';
import { safeDate, safeInt, safeMap, safeString } from './base';

/**
 * Maps a score DTO to a Score domain model.
 *
 * @param dto - The raw DTO object from the OpenAPI client
 * @returns The Score domain model, or undefined if dto is null/undefined
 */
export function scoreFromDto(dto: unknown): Score | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeInt(d.id),
    provider: safeString(d.provider),
    type: safeString(d.type),
    score: safeString(d.score),
    updateDate: safeDate(d.updateDate ?? d.update_date),
  };
}

/**
 * Maps an array of score DTOs to Score domain models.
 *
 * @param dtos - The array of raw DTO objects
 * @returns Array of Score domain models (empty array if input is null/undefined)
 */
export function scoresFromDto(dtos: unknown[] | null | undefined): Score[] {
  return safeMap(dtos, scoreFromDto);
}
