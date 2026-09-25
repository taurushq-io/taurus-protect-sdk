/**
 * Audit mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { AuditTrail } from '../models/audit';
import { safeDate, safeMap, safeString } from './base';

/**
 * Maps an audit trail DTO to an AuditTrail domain model.
 */
export function auditTrailFromDto(dto: unknown): AuditTrail | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const user =
    d.user && typeof d.user === 'object'
      ? (d.user as { id?: unknown; email?: unknown; externalUserId?: unknown })
      : undefined;

  // Handle details field - could be object or JSON string
  let details: Record<string, unknown> | undefined;
  if (d.details) {
    if (typeof d.details === 'object') {
      details = d.details as Record<string, unknown>;
    } else if (typeof d.details === 'string') {
      try {
        details = JSON.parse(d.details);
      } catch {
        // Ignore JSON parse errors
      }
    }
  }

  return {
    id: safeString(d.id ?? d.auditId ?? d.audit_id),
    entity: safeString(d.entity ?? d.entityType ?? d.entity_type),
    entityId: safeString(d.entityId ?? d.entity_id),
    action: safeString(d.action ?? d.actionType ?? d.action_type),
    // The wire row carries the actor as a nested `user` object (TgvalidatordUserInfo).
    userId: safeString(user?.id ?? d.userId ?? d.user_id),
    userEmail: safeString(user?.email ?? d.userEmail ?? d.user_email ?? d.email),
    externalUserId: safeString(
      user?.externalUserId ?? d.externalUserId ?? d.external_user_id
    ),
    description: safeString(d.description ?? d.message),
    ipAddress: safeString(d.ipAddress ?? d.ip_address ?? d.ip),
    createdAt: safeDate(d.createdAt ?? d.created_at ?? d.creationDate ?? d.timestamp),
    details,
  };
}

/**
 * Maps an array of audit trail DTOs to AuditTrail domain models.
 */
export function auditTrailsFromDto(dtos: unknown[] | null | undefined): AuditTrail[] {
  return safeMap(dtos, auditTrailFromDto);
}
