/**
 * Webhook call mapper functions for converting OpenAPI DTOs to domain models.
 */

import type {
  WebhookCall,
  WebhookCallResponseCursor,
  WebhookCallResult,
} from '../models/webhook-call';
import type {
  TgvalidatordGetWebhookCallsReply,
  TgvalidatordResponseCursor,
  TgvalidatordWebhookCall,
} from '../internal/openapi/models';
import { safeDate, safeMap, safeString } from './base';

/**
 * Maps a webhook call DTO to a WebhookCall domain model.
 *
 * @param dto - The OpenAPI webhook call DTO
 * @returns The domain model webhook call, or undefined if dto is null/invalid
 */
export function webhookCallFromDto(dto: TgvalidatordWebhookCall | unknown): WebhookCall | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    eventId: safeString(d.eventId ?? d.event_id),
    webhookId: safeString(d.webhookId ?? d.webhook_id),
    payload: safeString(d.payload),
    status: safeString(d.status),
    statusMessage: safeString(d.statusMessage ?? d.status_message),
    attempts: safeString(d.attempts),
    updatedAt: safeDate(d.updatedAt ?? d.updated_at),
    createdAt: safeDate(d.createdAt ?? d.created_at),
  };
}

/**
 * Maps an array of webhook call DTOs to WebhookCall domain models.
 *
 * @param dtos - The list of OpenAPI webhook call DTOs
 * @returns The list of domain model webhook calls
 */
export function webhookCallsFromDto(
  dtos: TgvalidatordWebhookCall[] | unknown[] | null | undefined
): WebhookCall[] {
  return safeMap(dtos, webhookCallFromDto);
}

/**
 * Maps a response cursor DTO to a WebhookCallResponseCursor domain model.
 *
 * @param cursor - The OpenAPI response cursor
 * @returns The domain model cursor, or undefined if cursor is null/invalid
 */
export function webhookCallCursorFromDto(
  cursor: TgvalidatordResponseCursor | unknown
): WebhookCallResponseCursor | undefined {
  if (!cursor || typeof cursor !== 'object') {
    return undefined;
  }

  const c = cursor as Record<string, unknown>;
  return {
    nextPage: safeString(c.nextPage ?? c.next_page),
    previousPage: safeString(c.previousPage ?? c.previous_page),
    hasNextPage: c.nextPage != null || c.next_page != null,
    hasPreviousPage: c.previousPage != null || c.previous_page != null,
  };
}

/**
 * Maps a get webhook calls reply to a WebhookCallResult domain model.
 *
 * @param reply - The OpenAPI reply
 * @returns The webhook call result with pagination
 */
export function webhookCallResultFromDto(
  reply: TgvalidatordGetWebhookCallsReply | unknown
): WebhookCallResult {
  if (!reply || typeof reply !== 'object') {
    return { calls: [] };
  }

  const r = reply as Record<string, unknown>;
  return {
    calls: webhookCallsFromDto(r.calls as unknown[]),
    cursor: webhookCallCursorFromDto(r.cursor),
  };
}
