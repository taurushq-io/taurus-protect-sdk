/**
 * Webhook call mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { WebhookCall } from '../models/webhook-call';
import type { TgvalidatordWebhookCall } from '../internal/openapi/models';
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
