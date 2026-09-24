/**
 * Webhook models for Taurus-PROTECT SDK.
 */

import type { CursorPage, CursorPageOptions } from './pagination';

/**
 * Webhook status enum.
 */
export enum WebhookStatus {
  ACTIVE = 'ACTIVE',
  INACTIVE = 'INACTIVE',
  FAILED = 'FAILED',
}

/**
 * Webhook information.
 */
export interface Webhook {
  /** Unique webhook identifier */
  readonly id?: string;
  /** Webhook URL */
  readonly url?: string;
  /** Event types subscribed to */
  readonly events?: string[];
  /** Webhook status */
  readonly status?: WebhookStatus;
  /** Secret for HMAC signature verification */
  readonly secret?: string;
  /** Creation date */
  readonly createdAt?: Date;
  /** Last modification date */
  readonly updatedAt?: Date;
  /** Number of consecutive failures */
  readonly failureCount?: number;
  /** Last failure message */
  readonly lastFailureMessage?: string;
}

/**
 * Options for listing webhooks. A cursor list: `pageSize` 1-100 (default 20) and the
 * `cursor` of a previous page.
 */
export interface ListWebhooksOptions extends CursorPageOptions {
  /** Filter by webhook type */
  type?: string;
  /** Filter by URL */
  url?: string;
  /** Sort order (ASC or DESC) */
  sortOrder?: 'ASC' | 'DESC';
}

/**
 * A page of webhooks.
 */
export interface ListWebhooksResult {
  /** The webhooks of this page */
  readonly items: Webhook[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}


/**
 * Request for creating a webhook.
 */
export interface CreateWebhookRequest {
  /** Webhook URL */
  url: string;
  /** Event types to subscribe to */
  events: string[];
}
