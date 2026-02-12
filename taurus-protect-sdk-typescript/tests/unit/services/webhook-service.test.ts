/**
 * Unit tests for WebhookService.
 */

import { WebhookService } from '../../../src/services/webhook-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { WebhooksApi } from '../../../src/internal/openapi/apis/WebhooksApi';

function createMockWebhooksApi(): jest.Mocked<WebhooksApi> {
  return {
    webhookServiceGetWebhooks: jest.fn(),
    webhookServiceCreateWebhook: jest.fn(),
    webhookServiceDeleteWebhook: jest.fn(),
  } as unknown as jest.Mocked<WebhooksApi>;
}

describe('WebhookService', () => {
  let mockApi: jest.Mocked<WebhooksApi>;
  let service: WebhookService;

  beforeEach(() => {
    mockApi = createMockWebhooksApi();
    service = new WebhookService(mockApi);
  });

  describe('get', () => {
    it('should throw ValidationError when webhookId is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('webhookId is required');
    });

    it('should throw ValidationError when webhookId is whitespace', async () => {
      await expect(service.get('   ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when webhook not found', async () => {
      mockApi.webhookServiceGetWebhooks.mockResolvedValue({
        webhooks: [],
      } as never);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });

    it('should return webhook when found', async () => {
      mockApi.webhookServiceGetWebhooks.mockResolvedValue({
        webhooks: [{ id: 'wh-123', url: 'https://example.com', type: 'TRANSACTION' }],
      } as never);

      const webhook = await service.get('wh-123');
      expect(webhook).toBeDefined();
      expect(webhook.id).toBe('wh-123');
    });
  });

  describe('list', () => {
    it('should return webhooks', async () => {
      mockApi.webhookServiceGetWebhooks.mockResolvedValue({
        webhooks: [
          { id: 'wh-1', url: 'https://a.com' },
          { id: 'wh-2', url: 'https://b.com' },
        ],
      } as never);

      const webhooks = await service.list();
      expect(webhooks).toHaveLength(2);
    });

    it('should throw ValidationError when limit is 0', async () => {
      await expect(service.list({ limit: 0 })).rejects.toThrow(ValidationError);
      await expect(service.list({ limit: 0 })).rejects.toThrow('limit must be positive');
    });

    it('should throw ValidationError when limit is negative', async () => {
      await expect(service.list({ limit: -1 })).rejects.toThrow(ValidationError);
    });
  });

  describe('create', () => {
    it('should throw ValidationError when url is empty', async () => {
      await expect(
        service.create({ url: '', events: ['TRANSACTION'] })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.create({ url: '', events: ['TRANSACTION'] })
      ).rejects.toThrow('url is required');
    });

    it('should throw ValidationError when events is empty', async () => {
      await expect(
        service.create({ url: 'https://example.com', events: [] })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.create({ url: 'https://example.com', events: [] })
      ).rejects.toThrow('events cannot be empty');
    });

    it('should create webhook with valid request', async () => {
      mockApi.webhookServiceCreateWebhook.mockResolvedValue({
        webhook: { id: 'wh-new', url: 'https://example.com', type: 'TRANSACTION' },
      } as never);

      const webhook = await service.create({
        url: 'https://example.com',
        events: ['TRANSACTION'],
      });
      expect(webhook).toBeDefined();
    });
  });

  describe('delete', () => {
    it('should throw ValidationError when webhookId is empty', async () => {
      await expect(service.delete('')).rejects.toThrow(ValidationError);
      await expect(service.delete('')).rejects.toThrow('webhookId is required');
    });

    it('should delete webhook with valid ID', async () => {
      mockApi.webhookServiceDeleteWebhook.mockResolvedValue(undefined as never);
      await service.delete('wh-123');
      expect(mockApi.webhookServiceDeleteWebhook).toHaveBeenCalledWith({ id: 'wh-123' });
    });
  });
});
