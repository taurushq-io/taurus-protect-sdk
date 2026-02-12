/**
 * Webhook service for Taurus-PROTECT SDK.
 *
 * Provides methods for webhook management operations.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { WebhooksApi } from '../internal/openapi/apis/WebhooksApi';
import { webhookFromDto, webhooksFromDto } from '../mappers/webhook';
import type { CreateWebhookRequest, ListWebhooksOptions, Webhook } from '../models/webhook';
import { BaseService } from './base';

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
   * Note: This method lists webhooks and filters by ID since the API
   * does not provide a direct get-by-ID endpoint.
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
      // List webhooks and find the one with matching ID
      const response = await this.webhooksApi.webhookServiceGetWebhooks({
        cursorPageSize: '100',
      });

      const resp = response as Record<string, unknown>;
      const webhooksDto = resp.webhooks as unknown[];

      if (webhooksDto) {
        for (const dto of webhooksDto) {
          const d = dto as Record<string, unknown>;
          if (d.id === webhookId) {
            const webhook = webhookFromDto(dto);
            if (webhook) {
              return webhook;
            }
          }
        }
      }

      throw new NotFoundError(`Webhook ${webhookId} not found`);
    });
  }

  /**
   * Lists webhooks.
   *
   * @param options - Optional filtering options
   * @returns Array of webhooks
   * @throws {@link ValidationError} If limit is invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all webhooks
   * const webhooks = await webhookService.list();
   *
   * // List with filters
   * const webhooks = await webhookService.list({
   *   limit: 50,
   *   sortOrder: 'DESC',
   * });
   * ```
   */
  async list(options?: ListWebhooksOptions): Promise<Webhook[]> {
    const limit = options?.limit ?? 50;

    if (limit <= 0) {
      throw new ValidationError('limit must be positive');
    }

    return this.execute(async () => {
      const response = await this.webhooksApi.webhookServiceGetWebhooks({
        type: options?.type,
        url: options?.url,
        cursorPageSize: String(limit),
        sortOrder: options?.sortOrder,
      });

      const resp = response as Record<string, unknown>;
      const webhooksDto = resp.webhooks;
      return webhooksFromDto(webhooksDto as unknown[]);
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

      const resp = response as Record<string, unknown>;
      const webhookDto = resp.webhook ?? resp.result;
      const webhook = webhookFromDto(webhookDto);

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
