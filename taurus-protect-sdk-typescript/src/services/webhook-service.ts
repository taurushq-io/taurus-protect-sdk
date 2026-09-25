/**
 * Webhook service for Taurus-PROTECT SDK.
 *
 * Provides methods for webhook management operations.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { WebhooksApi } from '../internal/openapi/apis/WebhooksApi';
import { webhookFromDto, webhooksFromDto } from '../mappers/webhook';
import { buildCursorPage, cursorRequest } from '../models/pagination';
import type {
  CreateWebhookRequest,
  ListWebhooksOptions,
  ListWebhooksResult,
  Webhook,
} from '../models/webhook';
import { BaseService } from './base';
import { cursorQuery, scanPages } from './paging';

/**
 * Service for webhook management operations.
 *
 * Provides methods to list, get, create, and delete webhooks.
 * Webhooks allow you to receive notifications when events occur in Taurus-PROTECT.
 *
 * @example
 * ```typescript
 * // List webhooks
 * const webhooks = await webhookService.list();
 * for (const webhook of webhooks) {
 *   console.log(`${webhook.id}: ${webhook.url} (${webhook.status})`);
 * }
 *
 * // Get single webhook
 * const webhook = await webhookService.get('webhook-123');
 * console.log(`URL: ${webhook.url}`);
 *
 * // Create webhook
 * const newWebhook = await webhookService.create({
 *   url: 'https://example.com/webhook',
 *   events: ['REQUEST_CREATED', 'REQUEST_APPROVED'],
 * });
 *
 * // Delete webhook
 * await webhookService.delete('webhook-123');
 * ```
 */
export class WebhookService extends BaseService {
  private readonly webhooksApi: WebhooksApi;

  /**
   * Creates a new WebhookService instance.
   *
   * @param webhooksApi - The WebhooksApi instance from the OpenAPI client
   */
  constructor(webhooksApi: WebhooksApi) {
    super();
    this.webhooksApi = webhooksApi;
  }

  /**
   * Gets a webhook by ID.
   *
   * The API has no get-by-ID endpoint, so this walks the webhook list page by page
   * (100 per page) until the ID shows up.
   *
   * @param webhookId - The webhook ID to retrieve
   * @returns The webhook
   * @throws {@link ValidationError} If webhookId is empty
   * @throws {@link NotFoundError} If webhook not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const webhook = await webhookService.get('webhook-123');
   * console.log(`URL: ${webhook.url}, Events: ${webhook.events?.join(', ')}`);
   * ```
   */
  async get(webhookId: string): Promise<Webhook> {
    if (!webhookId || webhookId.trim() === '') {
      throw new ValidationError('webhookId is required');
    }

    return this.execute(async () => {
      // No single-webhook read exists: walk the list page by page until the id shows up.
      const dto = await scanPages(
        async (page) => {
          const response = await this.webhooksApi.webhookServiceGetWebhooks(cursorQuery(page));
          return { rows: response.webhooks, cursor: response.cursor };
        },
        (row) => row.id === webhookId
      );

      const webhook = webhookFromDto(dto);
      if (!webhook) {
        throw new NotFoundError(`Webhook ${webhookId} not found`);
      }
      return webhook;
    });
  }

  /**
   * Lists a page of webhooks.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of webhooks and its cursor pagination
   * @throws {@link ValidationError} If the paging options are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // First page
   * const { items, pagination } = await webhookService.list({ sortOrder: 'DESC' });
   *
   * // Next page
   * if (pagination.hasMore) {
   *   await webhookService.list({ sortOrder: 'DESC', cursor: pagination.nextCursor });
   * }
   * ```
   */
  async list(options?: ListWebhooksOptions): Promise<ListWebhooksResult> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.webhooksApi.webhookServiceGetWebhooks({
        type: options?.type,
        url: options?.url,
        ...cursorQuery(page),
        sortOrder: options?.sortOrder,
      });

      return {
        items: webhooksFromDto(response.webhooks),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }

  /**
   * Creates a new webhook.
   *
   * @param request - Webhook creation request
   * @returns The created webhook
   * @throws {@link ValidationError} If url is empty or events is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const webhook = await webhookService.create({
   *   url: 'https://example.com/webhook',
   *   events: ['REQUEST_CREATED', 'REQUEST_APPROVED'],
   * });
   * console.log(`Created webhook: ${webhook.id}`);
   * console.log(`Secret: ${webhook.secret}`);
   * ```
   */
  async create(request: CreateWebhookRequest): Promise<Webhook> {
    if (!request.url || request.url.trim() === '') {
      throw new ValidationError('url is required');
    }
    if (!request.events || request.events.length === 0) {
      throw new ValidationError('events cannot be empty');
    }

    return this.execute(async () => {
      const response = await this.webhooksApi.webhookServiceCreateWebhook({
        body: {
          url: request.url,
          type: request.events.join(','), // API expects comma-separated event types
        },
      });

      const webhook = webhookFromDto(response.webhook);

      if (!webhook) {
        throw new ValidationError('Failed to create webhook: no result returned');
      }

      return webhook;
    });
  }

  /**
   * Deletes a webhook.
   *
   * @param webhookId - The webhook ID to delete
   * @throws {@link ValidationError} If webhookId is empty
   * @throws {@link NotFoundError} If webhook not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await webhookService.delete('webhook-123');
   * console.log('Webhook deleted');
   * ```
   */
  async delete(webhookId: string): Promise<void> {
    if (!webhookId || webhookId.trim() === '') {
      throw new ValidationError('webhookId is required');
    }

    return this.execute(async () => {
      await this.webhooksApi.webhookServiceDeleteWebhook({ id: webhookId });
    });
  }
}
