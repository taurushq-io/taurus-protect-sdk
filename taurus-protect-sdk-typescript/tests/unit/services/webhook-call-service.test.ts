/**
 * Unit tests for WebhookCallService.
 */

import { WebhookCallService } from '../../../src/services/webhook-call-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { WebhookCallsApi } from '../../../src/internal/openapi/apis/WebhookCallsApi';

function createMockApi(): jest.Mocked<WebhookCallsApi> {
  return {
    webhookServiceGetWebhookCalls: jest.fn(),
  } as unknown as jest.Mocked<WebhookCallsApi>;
}

describe('WebhookCallService', () => {
  let mockApi: jest.Mocked<WebhookCallsApi>;
  let service: WebhookCallService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new WebhookCallService(mockApi);
  });

  describe('list', () => {
    it('should return webhook calls', async () => {
      mockApi.webhookServiceGetWebhookCalls.mockResolvedValue({
        result: [
          { id: 'call-1', status: 'SUCCESS' },
          { id: 'call-2', status: 'FAILED' },
        ],
      } as never);

      const result = await service.list();
      expect(result.calls).toBeDefined();
    });

    it('should handle empty results', async () => {
      mockApi.webhookServiceGetWebhookCalls.mockResolvedValue({
        result: [],
      } as never);

      const result = await service.list();
      expect(result.calls).toHaveLength(0);
    });
  });

  describe('get', () => {
    it('should throw ValidationError when callId is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('callId is required');
    });

    it('should throw NotFoundError when call not found', async () => {
      mockApi.webhookServiceGetWebhookCalls.mockResolvedValue({
        result: [],
      } as never);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });
  });
});
