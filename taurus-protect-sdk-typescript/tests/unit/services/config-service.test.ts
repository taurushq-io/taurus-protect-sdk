/**
 * Unit tests for ConfigService.
 */

import { ConfigService } from '../../../src/services/config-service';
import type { ConfigApi } from '../../../src/internal/openapi/apis/ConfigApi';

function createMockApi(): jest.Mocked<ConfigApi> {
  return {
    statusServiceGetConfigTenant: jest.fn(),
  } as unknown as jest.Mocked<ConfigApi>;
}

describe('ConfigService', () => {
  let mockApi: jest.Mocked<ConfigApi>;
  let service: ConfigService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new ConfigService(mockApi);
  });

  describe('getTenantConfig', () => {
    it('should return tenant config', async () => {
      mockApi.statusServiceGetConfigTenant.mockResolvedValue({
        config: {
          tenantId: 'tenant-123',
          baseCurrency: 'USD',
        },
      } as never);

      const config = await service.getTenantConfig();
      expect(config).toBeDefined();
    });

    it('should handle response with missing fields', async () => {
      mockApi.statusServiceGetConfigTenant.mockResolvedValue({
        config: {},
      } as never);

      const config = await service.getTenantConfig();
      expect(config).toBeDefined();
    });
  });
});
