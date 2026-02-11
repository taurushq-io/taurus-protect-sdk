/**
 * Unit tests for ContractWhitelistingService.
 */

import { ContractWhitelistingService } from '../../../src/services/contract-whitelisting-service';
import { NotFoundError, ValidationError } from '../../../src/errors';
import type { ContractWhitelistingApi } from '../../../src/internal/openapi/apis/ContractWhitelistingApi';

function createMockApi(): jest.Mocked<ContractWhitelistingApi> {
  return {
    whitelistServiceGetWhitelistedContract: jest.fn(),
    whitelistServiceGetWhitelistedContracts: jest.fn(),
    whitelistServiceGetWhitelistedContractsForApproval: jest.fn(),
    whitelistServiceCreateWhitelistedContract: jest.fn(),
    whitelistServiceApproveWhitelistedContract: jest.fn(),
    whitelistServiceRejectWhitelistedContract: jest.fn(),
    whitelistServiceUpdateWhitelistedContract: jest.fn(),
    whitelistServiceCreateWhitelistedContractAttributes: jest.fn(),
    whitelistServiceGetWhitelistedContractAttribute: jest.fn(),
    whitelistServiceDeleteWhitelistedContractAttribute: jest.fn(),
  } as unknown as jest.Mocked<ContractWhitelistingApi>;
}

describe('ContractWhitelistingService', () => {
  let mockApi: jest.Mocked<ContractWhitelistingApi>;
  let service: ContractWhitelistingService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new ContractWhitelistingService(mockApi);
  });

  describe('get', () => {
    it('should throw ValidationError when id is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('id is required');
    });

    it('should throw ValidationError when id is whitespace', async () => {
      await expect(service.get('  ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when contract is not found', async () => {
      mockApi.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: undefined,
      });

      await expect(service.get('999')).rejects.toThrow(NotFoundError);
    });

    it('should return contract for valid id', async () => {
      mockApi.whitelistServiceGetWhitelistedContract.mockResolvedValue({
        result: {
          id: '123',
          metadata: {
            payloadAsString: JSON.stringify({
              blockchain: 'ETH',
              network: 'mainnet',
              address: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
              name: 'USD Coin',
              symbol: 'USDC',
              decimals: 6,
            }),
          },
        },
      });

      const contract = await service.get('123');
      expect(contract).toBeDefined();
      expect(mockApi.whitelistServiceGetWhitelistedContract).toHaveBeenCalledWith({ id: '123' });
    });
  });

  describe('list', () => {
    it('should return paginated results', async () => {
      mockApi.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [
          { id: '1', metadata: { payloadAsString: JSON.stringify({ name: 'A', symbol: 'A' }) } },
          { id: '2', metadata: { payloadAsString: JSON.stringify({ name: 'B', symbol: 'B' }) } },
        ],
        totalItems: '5',
      });

      const result = await service.list({ limit: 50 });
      expect(result.items).toHaveLength(2);
      expect(result.pagination.totalItems).toBe(5);
    });

    it('should throw ValidationError when limit is 0', async () => {
      await expect(service.list({ limit: 0 })).rejects.toThrow(ValidationError);
      await expect(service.list({ limit: 0 })).rejects.toThrow('limit must be positive');
    });

    it('should throw ValidationError when offset is negative', async () => {
      await expect(service.list({ offset: -1 })).rejects.toThrow(ValidationError);
      await expect(service.list({ offset: -1 })).rejects.toThrow('offset cannot be negative');
    });

    it('should handle empty results', async () => {
      mockApi.whitelistServiceGetWhitelistedContracts.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      const result = await service.list();
      expect(result.items).toHaveLength(0);
    });
  });

  describe('listForApproval', () => {
    it('should return pending contracts', async () => {
      mockApi.whitelistServiceGetWhitelistedContractsForApproval.mockResolvedValue({
        result: [
          { id: '1', metadata: { payloadAsString: JSON.stringify({ name: 'Pending', symbol: 'PND' }) } },
        ],
        totalItems: '1',
      });

      const result = await service.listForApproval({ limit: 50 });
      expect(result.items).toHaveLength(1);
    });

    it('should throw ValidationError when limit is 0', async () => {
      await expect(service.listForApproval({ limit: 0 })).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when offset is negative', async () => {
      await expect(service.listForApproval({ offset: -1 })).rejects.toThrow(ValidationError);
    });
  });

  describe('create', () => {
    it('should create a contract and return its ID', async () => {
      mockApi.whitelistServiceCreateWhitelistedContract.mockResolvedValue({
        result: { id: '456' },
      });

      const id = await service.create({
        blockchain: 'ETH',
        network: 'mainnet',
        contractAddress: '0x1234',
        symbol: 'USDC',
        name: 'USD Coin',
        decimals: 6,
        kind: 'erc20',
      });

      expect(id).toBe('456');
      expect(mockApi.whitelistServiceCreateWhitelistedContract).toHaveBeenCalledWith({
        body: expect.objectContaining({
          blockchain: 'ETH',
          network: 'mainnet',
          symbol: 'USDC',
          name: 'USD Coin',
          kind: 'erc20',
        }),
      });
    });

    it('should throw ValidationError when blockchain is empty', async () => {
      await expect(
        service.create({ blockchain: '', network: 'mainnet', symbol: 'X', name: 'X', kind: 'erc20', decimals: 0 })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.create({ blockchain: '', network: 'mainnet', symbol: 'X', name: 'X', kind: 'erc20', decimals: 0 })
      ).rejects.toThrow('blockchain is required');
    });

    it('should throw ValidationError when network is empty', async () => {
      await expect(
        service.create({ blockchain: 'ETH', network: '', symbol: 'X', name: 'X', kind: 'erc20', decimals: 0 })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when symbol is empty', async () => {
      await expect(
        service.create({ blockchain: 'ETH', network: 'mainnet', symbol: '', name: 'X', kind: 'erc20', decimals: 0 })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when name is empty', async () => {
      await expect(
        service.create({ blockchain: 'ETH', network: 'mainnet', symbol: 'X', name: '', kind: 'erc20', decimals: 0 })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when kind is empty', async () => {
      await expect(
        service.create({ blockchain: 'ETH', network: 'mainnet', symbol: 'X', name: 'X', kind: '', decimals: 0 })
      ).rejects.toThrow(ValidationError);
    });
  });

  describe('approve', () => {
    it('should approve contracts', async () => {
      mockApi.whitelistServiceApproveWhitelistedContract.mockResolvedValue({});

      await service.approve(['1', '2'], 'base64sig', 'Approved');

      expect(mockApi.whitelistServiceApproveWhitelistedContract).toHaveBeenCalledWith({
        body: { ids: ['1', '2'], signature: 'base64sig', comment: 'Approved' },
      });
    });

    it('should throw ValidationError when ids is empty', async () => {
      await expect(service.approve([], 'sig', 'comment')).rejects.toThrow(ValidationError);
      await expect(service.approve([], 'sig', 'comment')).rejects.toThrow('ids cannot be empty');
    });

    it('should throw ValidationError when signature is empty', async () => {
      await expect(service.approve(['1'], '', 'comment')).rejects.toThrow(ValidationError);
      await expect(service.approve(['1'], '', 'comment')).rejects.toThrow('signature is required');
    });

    it('should throw ValidationError when comment is empty', async () => {
      await expect(service.approve(['1'], 'sig', '')).rejects.toThrow(ValidationError);
      await expect(service.approve(['1'], 'sig', '')).rejects.toThrow('comment is required');
    });
  });

  describe('reject', () => {
    it('should reject contracts', async () => {
      mockApi.whitelistServiceRejectWhitelistedContract.mockResolvedValue({});

      await service.reject(['1'], 'Not needed');

      expect(mockApi.whitelistServiceRejectWhitelistedContract).toHaveBeenCalledWith({
        body: { ids: ['1'], comment: 'Not needed' },
      });
    });

    it('should throw ValidationError when ids is empty', async () => {
      await expect(service.reject([], 'comment')).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when comment is empty', async () => {
      await expect(service.reject(['1'], '')).rejects.toThrow(ValidationError);
    });
  });

  describe('update', () => {
    it('should update a contract', async () => {
      mockApi.whitelistServiceUpdateWhitelistedContract.mockResolvedValue({});

      await service.update('123', { symbol: 'USDC', name: 'USD Coin V2', decimals: 6 });

      expect(mockApi.whitelistServiceUpdateWhitelistedContract).toHaveBeenCalledWith({
        id: '123',
        body: { symbol: 'USDC', name: 'USD Coin V2', decimals: '6' },
      });
    });

    it('should throw ValidationError when id is empty', async () => {
      await expect(service.update('', { symbol: 'X', name: 'X', decimals: 0 })).rejects.toThrow(
        ValidationError
      );
    });

    it('should throw ValidationError when symbol is empty', async () => {
      await expect(service.update('1', { symbol: '', name: 'X', decimals: 0 })).rejects.toThrow(
        ValidationError
      );
    });

    it('should throw ValidationError when name is empty', async () => {
      await expect(service.update('1', { symbol: 'X', name: '', decimals: 0 })).rejects.toThrow(
        ValidationError
      );
    });
  });

  describe('createAttribute', () => {
    it('should create an attribute', async () => {
      mockApi.whitelistServiceCreateWhitelistedContractAttributes.mockResolvedValue({});

      await service.createAttribute('123', 'category', 'stablecoin');

      expect(mockApi.whitelistServiceCreateWhitelistedContractAttributes).toHaveBeenCalledWith({
        whitelistedContractAddressId: '123',
        body: {
          attributes: [{ key: 'category', value: 'stablecoin', contentType: undefined, type: undefined, subtype: undefined }],
        },
      });
    });

    it('should throw ValidationError when contractId is empty', async () => {
      await expect(service.createAttribute('', 'key', 'value')).rejects.toThrow(ValidationError);
      await expect(service.createAttribute('', 'key', 'value')).rejects.toThrow('contractId is required');
    });

    it('should throw ValidationError when key is empty', async () => {
      await expect(service.createAttribute('1', '', 'value')).rejects.toThrow(ValidationError);
      await expect(service.createAttribute('1', '', 'value')).rejects.toThrow('key is required');
    });
  });

  describe('getAttribute', () => {
    it('should throw ValidationError when contractId is empty', async () => {
      await expect(service.getAttribute('', 'attr1')).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when attributeId is empty', async () => {
      await expect(service.getAttribute('1', '')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when attribute is not found', async () => {
      mockApi.whitelistServiceGetWhitelistedContractAttribute.mockResolvedValue({
        result: undefined,
      });

      await expect(service.getAttribute('1', 'attr1')).rejects.toThrow(NotFoundError);
    });
  });

  describe('deleteAttribute', () => {
    it('should delete an attribute', async () => {
      mockApi.whitelistServiceDeleteWhitelistedContractAttribute.mockResolvedValue({});

      await service.deleteAttribute('123', 'attr-456');

      expect(mockApi.whitelistServiceDeleteWhitelistedContractAttribute).toHaveBeenCalledWith({
        whitelistedContractAddressId: '123',
        id: 'attr-456',
      });
    });

    it('should throw ValidationError when contractId is empty', async () => {
      await expect(service.deleteAttribute('', 'attr1')).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when attributeId is empty', async () => {
      await expect(service.deleteAttribute('1', '')).rejects.toThrow(ValidationError);
    });
  });
});
