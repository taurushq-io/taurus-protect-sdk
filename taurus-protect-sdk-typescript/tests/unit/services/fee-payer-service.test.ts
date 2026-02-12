/**
 * Unit tests for FeePayerService.
 */

import { FeePayerService } from '../../../src/services/fee-payer-service';
import { NotFoundError, ValidationError } from '../../../src/errors';
import type { FeePayersApi } from '../../../src/internal/openapi/apis/FeePayersApi';

function createMockApi(): jest.Mocked<FeePayersApi> {
  return {
    feePayerServiceGetFeePayers: jest.fn(),
    feePayerServiceGetFeePayer: jest.fn(),
  } as unknown as jest.Mocked<FeePayersApi>;
}

describe('FeePayerService', () => {
  let mockApi: jest.Mocked<FeePayersApi>;
  let service: FeePayerService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new FeePayerService(mockApi);
  });

  describe('list', () => {
    it('should return fee payers', async () => {
      mockApi.feePayerServiceGetFeePayers.mockResolvedValue({
        result: [
          { id: 'fp-1', name: 'Fee Payer 1', blockchain: 'ETH', network: 'mainnet' },
          { id: 'fp-2', name: 'Fee Payer 2', blockchain: 'ETH', network: 'goerli' },
        ],
      } as never);

      const feePayers = await service.list();
      expect(feePayers).toBeDefined();
      expect(feePayers.length).toBeGreaterThanOrEqual(0);
    });

    it('should pass filter options to API', async () => {
      mockApi.feePayerServiceGetFeePayers.mockResolvedValue({
        result: [],
      } as never);

      await service.list({
        blockchain: 'ETH',
        network: 'mainnet',
        limit: 10,
        offset: 5,
        ids: ['fp-1'],
      });

      expect(mockApi.feePayerServiceGetFeePayers).toHaveBeenCalledWith({
        blockchain: 'ETH',
        network: 'mainnet',
        limit: '10',
        offset: '5',
        ids: ['fp-1'],
      });
    });

    it('should handle empty results', async () => {
      mockApi.feePayerServiceGetFeePayers.mockResolvedValue({
        result: [],
      } as never);

      const feePayers = await service.list();
      expect(feePayers).toHaveLength(0);
    });

    it('should work without options', async () => {
      mockApi.feePayerServiceGetFeePayers.mockResolvedValue({
        result: [],
      } as never);

      const feePayers = await service.list();
      expect(feePayers).toBeDefined();
    });
  });

  describe('get', () => {
    it('should throw ValidationError when id is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('id is required');
    });

    it('should throw ValidationError when id is whitespace', async () => {
      await expect(service.get('  ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when fee payer is not found', async () => {
      mockApi.feePayerServiceGetFeePayer.mockResolvedValue({
        feepayer: undefined,
        feePayer: undefined,
        result: undefined,
      } as never);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });

    it('should return fee payer for valid id', async () => {
      mockApi.feePayerServiceGetFeePayer.mockResolvedValue({
        feePayer: {
          id: 'fp-123',
          name: 'Main Fee Payer',
          blockchain: 'ETH',
          network: 'mainnet',
        },
      } as never);

      const feePayer = await service.get('fp-123');
      expect(feePayer).toBeDefined();
      expect(mockApi.feePayerServiceGetFeePayer).toHaveBeenCalledWith({ id: 'fp-123' });
    });
  });
});
