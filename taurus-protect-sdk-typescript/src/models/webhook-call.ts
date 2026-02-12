/**
 * Webhook call models for Taurus-PROTECT SDK.
 */

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
 * Options for listing webhook calls.
 */
export interface ListWebhookCallsOptions {
  /** Filter by event ID */
  eventId?: string;
  /** Filter by webhook ID */
  webhookId?: string;
  /** Filter by call status (SUCCESS, FAILED, PENDING) */
  status?: WebhookCallStatus;
  /** Sort order (ASC or DESC, default DESC) */
  sortOrder?: 'ASC' | 'DESC';
  /** Maximum number of calls to return */
  limit?: number;
  /** Pagination cursor for the current page */
  cursorCurrentPage?: string;
  /** Page request direction (FIRST, PREVIOUS, NEXT, LAST) */
  cursorPageRequest?: 'FIRST' | 'PREVIOUS' | 'NEXT' | 'LAST';
}

/**
 * Response cursor for paginated webhook call results.
 */
export interface WebhookCallResponseCursor {
  /** Base64-encoded cursor for the next page */
  readonly nextPage?: string;
  /** Base64-encoded cursor for the previous page */
  readonly previousPage?: string;
  /** Whether there is a next page available */
  readonly hasNextPage?: boolean;
  /** Whether there is a previous page available */
  readonly hasPreviousPage?: boolean;
}

/**
 * Paginated result of webhook calls.
 */
export interface WebhookCallResult {
  /** The webhook calls in this page */
  readonly calls: WebhookCall[];
  /** Cursor for pagination */
  readonly cursor?: WebhookCallResponseCursor;
}
