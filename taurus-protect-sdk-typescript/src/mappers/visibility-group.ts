/**
 * Visibility group mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { VisibilityGroup, VisibilityGroupUser } from '../models/visibility-group';
import { safeDate, safeMap, safeString } from './base';

/**
 * Maps a visibility group user DTO to a VisibilityGroupUser domain model.
 */
export function visibilityGroupUserFromDto(dto: unknown): VisibilityGroupUser | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    externalUserId: safeString(d.externalUserId ?? d.external_user_id),
  };
}

/**
 * Maps an array of visibility group user DTOs to VisibilityGroupUser domain models.
 */
export function visibilityGroupUsersFromDto(
  dtos: unknown[] | null | undefined
): VisibilityGroupUser[] {
  return safeMap(dtos, visibilityGroupUserFromDto);
}

/**
 * Maps a visibility group DTO to a VisibilityGroup domain model.
 */
export function visibilityGroupFromDto(dto: unknown): VisibilityGroup | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const users = d.users as unknown[] | undefined;

  return {
    id: safeString(d.id),
    tenantId: safeString(d.tenantId ?? d.tenant_id),
    name: safeString(d.name),
    description: safeString(d.description),
    users: users ? visibilityGroupUsersFromDto(users) : undefined,
    creationDate: safeDate(d.creationDate ?? d.creation_date),
    updateDate: safeDate(d.updateDate ?? d.update_date),
    userCount: safeString(d.userCount ?? d.user_count),
  };
}

/**
 * Maps an array of visibility group DTOs to VisibilityGroup domain models.
 */
export function visibilityGroupsFromDto(
  dtos: unknown[] | null | undefined
): VisibilityGroup[] {
  return safeMap(dtos, visibilityGroupFromDto);
}
