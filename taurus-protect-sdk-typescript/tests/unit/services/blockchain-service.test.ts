/**
 * Unit tests for BlockchainService.
 */

import { BlockchainService } from '../../../src/services/blockchain-service';
import { ValidationError } from '../../../src/errors';
import type { BlockchainApi } from '../../../src/internal/openapi/apis/BlockchainApi';

function createMockApi(): jest.Mocked<BlockchainApi> {
  return {
    blockchainServiceGetBlockchains: jest.fn(),
    blockchainServiceGetBlockchainsInfo: jest.fn(),
  } as unknown as jest.Mocked<BlockchainApi>;
}

describe('BlockchainService', () => {
  let mockApi: jest.Mocked<BlockchainApi>;
  let service: BlockchainService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new BlockchainService(mockApi);
  });

  describe('list', () => {
    it('should return blockchains', async () => {
      mockApi.blockchainServiceGetBlockchains.mockResolvedValue({
        blockchains: [
          { blockchain: 'ETH', network: 'mainnet', symbol: 'ETH' },
          { blockchain: 'BTC', network: 'mainnet', symbol: 'BTC' },
        ],
      } as never);

      const blockchains = await service.list();
      expect(blockchains).toHaveLength(2);
    });

    it('should handle empty results', async () => {
      mockApi.blockchainServiceGetBlockchains.mockResolvedValue({
        blockchains: [],
      } as never);

      const blockchains = await service.list();
      expect(blockchains).toHaveLength(0);
    });
  });

  describe('get', () => {
    it('should throw ValidationError when blockchain is empty', async () => {
      await expect(service.get('', 'mainnet')).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when network is empty', async () => {
      await expect(service.get('ETH', '')).rejects.toThrow(ValidationError);
    });

    it('should return blockchain when found', async () => {
      mockApi.blockchainServiceGetBlockchains.mockResolvedValue({
        blockchains: [{ blockchain: 'ETH', network: 'mainnet', symbol: 'ETH' }],
      } as never);

      const blockchain = await service.get('ETH', 'mainnet');
      expect(blockchain).toBeDefined();
    });
  });
});
