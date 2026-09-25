/**
 * Change mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Change, CreateChangeRequest } from '../models/change';
import { safeDate, safeInt, safeMap, safeString } from './base';

/**
 * Maps a change DTO to a Change domain model.
 */
export function changeFromDto(dto: unknown): Change | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;

  return {
    id: safeString(d.id),
    tenantId: safeInt(d.tenantId ?? d.tenant_id),
    creatorId: safeString(d.creatorId ?? d.creator_id),
    creatorExternalId: safeString(d.creatorExternalId ?? d.creator_external_id),
    action: safeString(d.action),
    entity: safeString(d.entity),
    entityId: safeString(d.entityId ?? d.entity_id),
    entityUUID: safeString(d.entityUUID ?? d.entity_uuid),
    changes: d.changes as Record<string, string> | undefined,
    comment: safeString(d.comment),
    // Map creationDate to createdAt (as in Java mapper)
    createdAt: safeDate(d.creationDate ?? d.creation_date ?? d.createdAt ?? d.created_at),
  };
}

/**
 * Maps an array of change DTOs to Change domain models.
 */
export function changesFromDto(dtos: unknown[] | null | undefined): Change[] {
  return safeMap(dtos, changeFromDto);
}

/**
 * Converts a CreateChangeRequest to an OpenAPI DTO.
 * Maps SDK 'comment' to OpenAPI 'changeComment' (matching Java ChangeMapper).
 */
export function createChangeRequestToDto(request: CreateChangeRequest): Record<string, unknown> {
  return {
    action: request.action,
    entity: request.entity,
    entityId: request.entityId,
    changes: request.changes,
    changeComment: request.comment,
  };
}
