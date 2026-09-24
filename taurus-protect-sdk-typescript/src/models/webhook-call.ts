/**
 * Webhook call models for Taurus-PROTECT SDK.
 */

import type { CursorNavigationOptions, CursorPage } from './pagination';

/**
 * Webhook call status enum.
 */
export type WebhookCallStatus = 'SUCCESS' | 'FAILED' | 'PENDING';

/**
 * A webhook call record representing a single invocation of a configured webhook.
 *
 * A webhook call represents a single HTTP request made to a webhook URL
 * when an event occurs in Taurus-PROTECT. It includes information about
 * the payload sent, the delivery status, and retry attempts.
 */
export interface WebhookCall {
  /** Unique call identifier */
  readonly id?: string;
  /** The event ID that triggered this call */
  readonly eventId?: string;
  /** The webhook ID that was called */
  readonly webhookId?: string;
  /** The payload sent in the webhook call */
  readonly payload?: string;
  /** The status of the call (SUCCESS, FAILED, PENDING) */
  readonly status?: string;
  /** Status message with details about the call result */
  readonly statusMessage?: string;
  /** Number of delivery attempts */
  readonly attempts?: string;
  /** Last modification date */
  readonly updatedAt?: Date;
  /** Creation date */
  readonly createdAt?: Date;
}

/**
 * Options for listing webhook calls. A cursor list: `pageSize` 1-100 (default 20) and the
 * `cursor` of a previous page.
 */
export interface ListWebhookCallsOptions extends CursorNavigationOptions {
  /** Filter by event ID */
  eventId?: string;
  /** Filter by webhook ID */
  webhookId?: string;
  /** Filter by call status (SUCCESS, FAILED, PENDING) */
  status?: WebhookCallStatus;
  /** Sort order (ASC or DESC, default DESC) */
  sortOrder?: 'ASC' | 'DESC';
}

/**
 * A page of webhook calls.
 */
export interface WebhookCallResult {
  /** The webhook calls in this page */
  readonly calls: WebhookCall[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}

