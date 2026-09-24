/**
 * Webhook call service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving webhook call history.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { WebhookCallsApi } from '../internal/openapi/apis/WebhookCallsApi';
import { webhookCallFromDto, webhookCallsFromDto } from '../mappers/webhook-call';
import { buildCursorPage, cursorRequest } from '../models/pagination';
import type {
  ListWebhookCallsOptions,
  WebhookCall,
  WebhookCallResult,
} from '../models/webhook-call';
import { BaseService } from './base';
import { cursorQuery, scanPages } from './paging';

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
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of webhook calls and its cursor pagination
   * @throws {@link ValidationError} If the paging options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List recent webhook calls
   * const result = await webhookCallService.list({ pageSize: 100 });
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
   * let result = await webhookCallService.list({ pageSize: 50 });
   * while (result.pagination.hasMore) {
   *   result = await webhookCallService.list({
   *     pageSize: 50,
   *     cursor: result.pagination.nextCursor,
   *   });
   * }
   * ```
   */
  async list(options?: ListWebhookCallsOptions): Promise<WebhookCallResult> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.webhookCallsApi.webhookServiceGetWebhookCalls({
        eventID: options?.eventId,
        webhookID: options?.webhookId,
        status: options?.status,
        ...cursorQuery(page),
        sortOrder: options?.sortOrder,
      });

      return {
        calls: webhookCallsFromDto(response.calls),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }

  /**
   * Gets a webhook call by ID.
   *
   * The API has no get-by-ID endpoint, so this walks the call list page by page
   * (100 per page) until the ID shows up.
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
      // No single-call read exists: walk the list page by page until the id shows up.
      const dto = await scanPages(
        async (page) => {
          const response = await this.webhookCallsApi.webhookServiceGetWebhookCalls(
            cursorQuery(page)
          );
          return { rows: response.calls, cursor: response.cursor };
        },
        (row) => row.id === callId
      );

      const call = webhookCallFromDto(dto);
      if (!call) {
        throw new NotFoundError(`Webhook call with id '${callId}' not found`);
      }
      return call;
    });
  }
}
