/**
 * Webhook mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Webhook } from '../models/webhook';
import { WebhookStatus } from '../models/webhook';
import { safeDate, safeInt, safeMap, safeString } from './base';

/**
 * Maps a webhook status string to WebhookStatus enum.
 */
function mapWebhookStatus(status: string | undefined): WebhookStatus | undefined {
  if (!status) {
    return undefined;
  }
  const upper = status.toUpperCase();
  switch (upper) {
    case 'ACTIVE':
    case 'ENABLED':
      return WebhookStatus.ACTIVE;
    case 'INACTIVE':
    case 'DISABLED':
      return WebhookStatus.INACTIVE;
    case 'FAILED':
      return WebhookStatus.FAILED;
    default:
      return undefined;
  }
}

/**
 * Maps a webhook DTO to a Webhook domain model.
 */
export function webhookFromDto(dto: unknown): Webhook | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const eventsRaw = d.events ?? d.type ?? d.types;
  let events: string[] | undefined;

  if (Array.isArray(eventsRaw)) {
    events = eventsRaw.map((e) => String(e));
  } else if (typeof eventsRaw === 'string') {
    // API may return comma-separated event types
    events = eventsRaw.split(',').map((e) => e.trim());
  }

  return {
    id: safeString(d.id ?? d.webhookId ?? d.webhook_id),
    url: safeString(d.url),
    events,
    status: mapWebhookStatus(safeString(d.status)),
    secret: safeString(d.secret),
    createdAt: safeDate(d.createdAt ?? d.created_at ?? d.creationDate),
    updatedAt: safeDate(d.updatedAt ?? d.updated_at ?? d.modificationDate),
    failureCount: safeInt(d.failureCount ?? d.failure_count ?? d.consecutiveFailures),
    lastFailureMessage: safeString(d.lastFailureMessage ?? d.last_failure_message ?? d.lastError),
  };
}

/**
 * Maps an array of webhook DTOs to Webhook domain models.
 */
export function webhooksFromDto(dtos: unknown[] | null | undefined): Webhook[] {
  return safeMap(dtos, webhookFromDto);
}
