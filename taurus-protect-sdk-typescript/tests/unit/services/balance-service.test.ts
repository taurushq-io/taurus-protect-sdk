/**
 * Unit tests for BalanceService.
 */

import { BalanceService } from '../../../src/services/balance-service';
import { ValidationError } from '../../../src/errors';
import type { BalancesApi } from '../../../src/internal/openapi/apis/BalancesApi';

function createMockBalancesApi(): jest.Mocked<BalancesApi> {
  return {
    walletServiceGetBalances: jest.fn(),
    walletServiceGetNFTCollectionBalances: jest.fn(),
  } as unknown as jest.Mocked<BalancesApi>;
}

describe('BalanceService', () => {
  let mockApi: jest.Mocked<BalancesApi>;
  let service: BalanceService;

  beforeEach(() => {
    mockApi = createMockBalancesApi();
    service = new BalanceService(mockApi);
  });

  describe('list', () => {
    it('should return balances', async () => {
      mockApi.walletServiceGetBalances.mockResolvedValue({
        balances: [
          { currency: 'ETH', totalConfirmed: '10.5' },
          { currency: 'BTC', totalConfirmed: '0.5' },
        ],
      } as any);

      const result = await service.list();

      expect(result).toHaveLength(2);
    });

    it('should pass currency filter', async () => {
      mockApi.walletServiceGetBalances.mockResolvedValue({
        balances: [],
      } as any);

      await service.list({ currency: 'ETH' });

      expect(mockApi.walletServiceGetBalances).toHaveBeenCalledWith(
        expect.objectContaining({
          currency: 'ETH',
        })
      );
    });

    it('should throw ValidationError when limit is invalid', async () => {
      await expect(service.list({ limit: 0 })).rejects.toThrow(ValidationError);
      await expect(service.list({ limit: 0 })).rejects.toThrow('limit must be positive');
    });

    it('should use default limit when not provided', async () => {
      mockApi.walletServiceGetBalances.mockResolvedValue({
        balances: [],
      } as any);

      await service.list();

      expect(mockApi.walletServiceGetBalances).toHaveBeenCalledWith(
        expect.objectContaining({
          limit: '50',
        })
      );
    });
  });

  describe('listNFTCollections', () => {
    it('should return NFT collection balances', async () => {
      mockApi.walletServiceGetNFTCollectionBalances.mockResolvedValue({
        collections: [
          { name: 'CryptoKitties', count: '10' },
        ],
      } as any);

      const result = await service.listNFTCollections({
        blockchain: 'ETH',
        network: 'mainnet',
      });

      expect(result).toHaveLength(1);
    });

    it('should throw ValidationError when blockchain is missing', async () => {
      await expect(
        service.listNFTCollections({ blockchain: '', network: 'mainnet' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.listNFTCollections({ blockchain: '', network: 'mainnet' })
      ).rejects.toThrow('blockchain is required');
    });

    it('should throw ValidationError when network is missing', async () => {
      await expect(
        service.listNFTCollections({ blockchain: 'ETH', network: '' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.listNFTCollections({ blockchain: 'ETH', network: '' })
      ).rejects.toThrow('network is required');
    });

    it('should throw ValidationError when limit is invalid', async () => {
      await expect(
        service.listNFTCollections({ blockchain: 'ETH', network: 'mainnet', limit: 0 })
      ).rejects.toThrow(ValidationError);
    });
  });
});
