/**
 * Unit tests for FeeService.
 */

import { FeeService } from '../../../src/services/fee-service';
import type { FeeApi } from '../../../src/internal/openapi/apis/FeeApi';

function createMockApi(): jest.Mocked<FeeApi> {
  return {
    feeServiceGetFees: jest.fn(),
    feeServiceGetFeesV2: jest.fn(),
  } as unknown as jest.Mocked<FeeApi>;
}

describe('FeeService', () => {
  let mockApi: jest.Mocked<FeeApi>;
  let service: FeeService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new FeeService(mockApi);
  });

  describe('getFees', () => {
    it('should return fees', async () => {
      mockApi.feeServiceGetFees.mockResolvedValue({
        result: [
          { key: 'ETH', value: '0.001' },
          { key: 'BTC', value: '0.0001' },
        ],
      } as never);

      const fees = await service.getFees();
      expect(fees).toHaveLength(2);
    });

    it('should handle empty results', async () => {
      mockApi.feeServiceGetFees.mockResolvedValue({
        result: [],
      } as never);

      const fees = await service.getFees();
      expect(fees).toHaveLength(0);
    });
  });

  describe('getFeesV2', () => {
    it('should return v2 fees', async () => {
      mockApi.feeServiceGetFeesV2.mockResolvedValue({
        result: [
          { currencyId: 'ETH', value: '0.001', denom: 'gwei' },
        ],
      } as never);

      const fees = await service.getFeesV2();
      expect(fees).toHaveLength(1);
    });

    it('should handle empty results', async () => {
      mockApi.feeServiceGetFeesV2.mockResolvedValue({
        result: [],
      } as never);

      const fees = await service.getFeesV2();
      expect(fees).toHaveLength(0);
    });
  });
});
