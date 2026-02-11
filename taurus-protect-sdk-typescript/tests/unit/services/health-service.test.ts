/**
 * Unit tests for HealthService.
 */

import { HealthService } from '../../../src/services/health-service';
import type { HealthApi } from '../../../src/internal/openapi/apis/HealthApi';

function createMockApi(): jest.Mocked<HealthApi> {
  return {
    healthServiceGetHealthChecks: jest.fn(),
    statusServiceGetGlobalComponentStatus: jest.fn(),
  } as unknown as jest.Mocked<HealthApi>;
}

describe('HealthService', () => {
  let mockApi: jest.Mocked<HealthApi>;
  let service: HealthService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new HealthService(mockApi);
  });

  describe('check', () => {
    it('should return healthy status', async () => {
      mockApi.healthServiceGetHealthChecks.mockResolvedValue({
        checks: [{ healthy: true }],
      } as never);

      const status = await service.check();
      expect(status).toBeDefined();
      expect(status.status).toBe('healthy');
    });

    it('should return unhealthy status on error', async () => {
      mockApi.healthServiceGetHealthChecks.mockRejectedValue(new Error('Connection refused'));

      const status = await service.check();
      expect(status.status).toBe('unhealthy');
    });
  });

  describe('getGlobalStatus', () => {
    it('should return health status', async () => {
      mockApi.statusServiceGetGlobalComponentStatus.mockResolvedValue({
        status: 'HEALTHY',
      } as never);

      const status = await service.getGlobalStatus();
      expect(status).toBeDefined();
    });

    it('should handle unhealthy status', async () => {
      mockApi.statusServiceGetGlobalComponentStatus.mockResolvedValue({
        status: 'UNHEALTHY',
      } as never);

      const status = await service.getGlobalStatus();
      expect(status).toBeDefined();
      expect(status.status).not.toBe('healthy');
    });
  });
});
