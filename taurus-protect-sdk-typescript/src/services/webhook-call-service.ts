/**
 * Webhook call service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving webhook call history.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { WebhookCallsApi } from '../internal/openapi/apis/WebhookCallsApi';
import { webhookCallFromDto, webhookCallResultFromDto } from '../mappers/webhook-call';
import type {
  ListWebhookCallsOptions,
  WebhookCall,
  WebhookCallResponseCursor,
  WebhookCallResult,
  WebhookCallStatus,
} from '../models/webhook-call';
import { BaseService } from './base';

// Re-export types from models for convenience
export type {
  ListWebhookCallsOptions,
  WebhookCall,
  WebhookCallResponseCursor,
  WebhookCallResult,
  WebhookCallStatus,
};

/**
 * Service for retrieving webhook call history.
 *
 * Provides access to the history of webhook invocations,
 * including their delivery status and payload information.
 *
 * @example
 * ```typescript
 * // List all webhook calls
 * const result = await webhookCallService.list();
 * for (const call of result.calls) {
 *   console.log(`${call.id}: ${call.status} (${call.attempts} attempts)`);
 * }
 *
 * // Get a specific webhook call by ID
 * const call = await webhookCallService.get('call-123');
 * console.log(`Status: ${call.status}`);
 *
 * // Filter by webhook ID
 * const webhookCalls = await webhookCallService.list({
 *   webhookId: 'webhook-123',
 * });
 *
 * // Filter by status
 * const failedCalls = await webhookCallService.list({
 *   status: 'FAILED',
 *   sortOrder: 'DESC',
 * });
 * ```
 */
export class WebhookCallService extends BaseService {
  private readonly webhookCallsApi: WebhookCallsApi;

  /**
   * Creates a new WebhookCallService instance.
   *
   * @param webhookCallsApi - The WebhookCallsApi instance from the OpenAPI client
   */
  constructor(webhookCallsApi: WebhookCallsApi) {
    super();
    this.webhookCallsApi = webhookCallsApi;
  }

  /**
   * Lists webhook calls with optional filtering.
   *
   * Returns a paginated list of webhook calls that can be filtered by
   * event ID, webhook ID, or status.
   *
   * @param options - Optional filtering options
   * @returns Paginated result containing webhook calls and cursor
   * @throws {@link ValidationError} If limit is invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List recent webhook calls
   * const result = await webhookCallService.list({ limit: 100 });
   * console.log(`Found ${result.calls.length} calls`);
   *
   * // Filter by webhook
   * const webhookCalls = await webhookCallService.list({
   *   webhookId: 'webhook-123',
   * });
   *
   * // Filter by event
   * const eventCalls = await webhookCallService.list({
   *   eventId: 'event-456',
   * });
   *
   * // Filter by status
   * const failedCalls = await webhookCallService.list({
   *   status: 'FAILED',
   * });
   *
   * // Paginate through results
   * let result = await webhookCallService.list({ limit: 50 });
   * while (result.cursor?.hasNextPage) {
   *   result = await webhookCallService.list({
   *     limit: 50,
   *     cursorCurrentPage: result.cursor.nextPage,
   *     cursorPageRequest: 'NEXT',
   *   });
   * }
   * ```
   */
  async list(options?: ListWebhookCallsOptions): Promise<WebhookCallResult> {
    const limit = options?.limit ?? 50;

    if (limit <= 0) {
      throw new ValidationError('limit must be positive');
    }

    return this.execute(async () => {
      const response = await this.webhookCallsApi.webhookServiceGetWebhookCalls({
        eventID: options?.eventId,
        webhookID: options?.webhookId,
        status: options?.status,
        cursorPageSize: String(limit),
        cursorCurrentPage: options?.cursorCurrentPage,
        cursorPageRequest: options?.cursorPageRequest,
        sortOrder: options?.sortOrder,
      });

      return webhookCallResultFromDto(response);
    });
  }

  /**
   * Gets a webhook call by ID.
   *
   * Note: This method fetches webhook calls and filters by ID since the API
   * does not provide a direct get-by-ID endpoint.
   *
   * @param callId - The unique webhook call identifier
   * @returns The webhook call
   * @throws {@link ValidationError} If callId is empty
   * @throws {@link NotFoundError} If webhook call is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const call = await webhookCallService.get('call-123');
   * console.log(`Status: ${call.status}, Attempts: ${call.attempts}`);
   * console.log(`Payload: ${call.payload}`);
   * ```
   */
  async get(callId: string): Promise<WebhookCall> {
    if (!callId || callId.trim() === '') {
      throw new ValidationError('callId is required');
    }

    return this.execute(async () => {
      // Fetch webhook calls - there's no direct get-by-id endpoint,
      // so we fetch a page and search for the matching ID
      const response = await this.webhookCallsApi.webhookServiceGetWebhookCalls({
        cursorPageSize: '100',
      });

      const calls = response.calls ?? [];
      for (const dto of calls) {
        if (dto.id === callId) {
          const call = webhookCallFromDto(dto);
          if (call) {
            return call;
          }
        }
      }

      throw new NotFoundError(`Webhook call with id '${callId}' not found`);
    });
  }
}
