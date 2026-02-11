/**
 * Unit tests for AddressService.
 */

import { AddressService } from '../../../src/services/address-service';
import { ConfigurationError, ValidationError, NotFoundError } from '../../../src/errors';
import type { AddressesApi } from '../../../src/internal/openapi/apis/AddressesApi';
import type { RulesContainerCache } from '../../../src/cache';
import type { DecodedRulesContainer } from '../../../src/models/governance-rules';

// Mock the address signature verification helper so unit tests don't need real keys
jest.mock('../../../src/helpers', () => ({
  ...jest.requireActual('../../../src/helpers'),
  verifyAddressSignature: jest.fn().mockReturnValue(true),
}));

function createMockAddressesApi(): jest.Mocked<AddressesApi> {
  return {
    walletServiceGetAddress: jest.fn(),
    walletServiceGetAddresses: jest.fn(),
    walletServiceCreateAddress: jest.fn(),
    walletServiceCreateAddressAttributes: jest.fn(),
    walletServiceDeleteAddressAttribute: jest.fn(),
    walletServiceGetAddressProofOfReserve: jest.fn(),
  } as unknown as jest.Mocked<AddressesApi>;
}

function createMockRulesCache(): jest.Mocked<RulesContainerCache> {
  return {
    get: jest.fn().mockResolvedValue({
      users: [],
      groups: [],
      getHsmPublicKey: jest.fn().mockReturnValue(undefined),
    } as unknown as DecodedRulesContainer),
    clear: jest.fn(),
  } as unknown as jest.Mocked<RulesContainerCache>;
}

describe('AddressService', () => {
  let mockApi: jest.Mocked<AddressesApi>;
  let mockRulesCache: jest.Mocked<RulesContainerCache>;
  let service: AddressService;

  beforeEach(() => {
    mockApi = createMockAddressesApi();
    mockRulesCache = createMockRulesCache();
    service = new AddressService(mockApi, mockRulesCache);
  });

  describe('constructor', () => {
    it('should throw ConfigurationError when rulesCache is not provided', () => {
      expect(() => new AddressService(mockApi, undefined as unknown as RulesContainerCache)).toThrow(
        ConfigurationError
      );
      expect(() => new AddressService(mockApi, null as unknown as RulesContainerCache)).toThrow(
        ConfigurationError
      );
    });
  });

  describe('get', () => {
    it('should return an address for a valid ID', async () => {
      mockApi.walletServiceGetAddress.mockResolvedValue({
        result: {
          id: '456',
          walletId: '123',
          address: '0xabc123',
          currency: 'ETH',
          label: 'Test Address',
          signature: 'dGVzdHNpZw==',
          creationDate: new Date('2024-01-01'),
        },
      });

      const address = await service.get(456);

      expect(address).toBeDefined();
      expect(address.id).toBe('456');
      expect(address.address).toBe('0xabc123');
      expect(address.walletId).toBe('123');
      expect(mockApi.walletServiceGetAddress).toHaveBeenCalledWith({ id: '456' });
    });

    it('should throw ValidationError when addressId is invalid', async () => {
      await expect(service.get(0)).rejects.toThrow(ValidationError);
      await expect(service.get(0)).rejects.toThrow('addressId must be positive');
      await expect(service.get(-1)).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when address is not found', async () => {
      mockApi.walletServiceGetAddress.mockResolvedValue({
        result: undefined,
      });

      await expect(service.get(999)).rejects.toThrow(NotFoundError);
    });
  });

  describe('list', () => {
    it('should return addresses for a wallet', async () => {
      mockApi.walletServiceGetAddresses.mockResolvedValue({
        result: [
          { id: '1', walletId: '100', address: '0xaaa', currency: 'ETH', signature: 'c2ln' },
          { id: '2', walletId: '100', address: '0xbbb', currency: 'ETH', signature: 'c2ln' },
        ],
        totalItems: '2',
      });

      const result = await service.list(100);

      expect(result.items).toHaveLength(2);
      expect(result.items[0].address).toBe('0xaaa');
      expect(result.items[1].address).toBe('0xbbb');
    });

    it('should throw ValidationError when walletId is invalid', async () => {
      await expect(service.list(0)).rejects.toThrow(ValidationError);
      await expect(service.list(0)).rejects.toThrow('walletId must be positive');
    });

    it('should throw ValidationError when limit is invalid', async () => {
      await expect(service.list(1, { limit: 0 })).rejects.toThrow(ValidationError);
      await expect(service.list(1, { limit: 0 })).rejects.toThrow('limit must be positive');
    });

    it('should throw ValidationError when offset is negative', async () => {
      await expect(service.list(1, { offset: -1 })).rejects.toThrow(ValidationError);
      await expect(service.list(1, { offset: -1 })).rejects.toThrow('offset cannot be negative');
    });
  });

  describe('create', () => {
    it('should create an address with valid request', async () => {
      mockApi.walletServiceCreateAddress.mockResolvedValue({
        result: {
          id: '789',
          walletId: '123',
          address: '0xnew',
          currency: 'ETH',
          label: 'New Address',
        },
      });

      const address = await service.create({
        walletId: '123',
        label: 'New Address',
        comment: 'test',
      });

      expect(address).toBeDefined();
      expect(address.id).toBe('789');
      expect(address.label).toBe('New Address');
    });

    it('should throw ValidationError when walletId is missing', async () => {
      await expect(
        service.create({ walletId: '', label: 'test' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.create({ walletId: '', label: 'test' })
      ).rejects.toThrow('walletId is required');
    });

    it('should throw ValidationError when label is missing', async () => {
      await expect(
        service.create({ walletId: '123', label: '' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.create({ walletId: '123', label: '' })
      ).rejects.toThrow('label is required');
    });
  });

  describe('createAttribute', () => {
    it('should create an attribute for an address', async () => {
      mockApi.walletServiceCreateAddressAttributes.mockResolvedValue({});

      await service.createAttribute(456, 'tag', 'vip');

      expect(mockApi.walletServiceCreateAddressAttributes).toHaveBeenCalledWith({
        addressId: '456',
        body: {
          attributes: [{ key: 'tag', value: 'vip' }],
        },
      });
    });

    it('should throw ValidationError when addressId is invalid', async () => {
      await expect(service.createAttribute(0, 'key', 'value')).rejects.toThrow(
        ValidationError
      );
    });

    it('should throw ValidationError when key is empty', async () => {
      await expect(service.createAttribute(1, '', 'value')).rejects.toThrow(
        ValidationError
      );
    });
  });

  describe('deleteAttribute', () => {
    it('should delete an attribute from an address', async () => {
      mockApi.walletServiceDeleteAddressAttribute.mockResolvedValue({});

      await service.deleteAttribute(456, 789);

      expect(mockApi.walletServiceDeleteAddressAttribute).toHaveBeenCalledWith({
        addressId: '456',
        id: '789',
      });
    });

    it('should throw ValidationError when addressId is invalid', async () => {
      await expect(service.deleteAttribute(0, 1)).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when attributeId is invalid', async () => {
      await expect(service.deleteAttribute(1, 0)).rejects.toThrow(ValidationError);
      await expect(service.deleteAttribute(1, 0)).rejects.toThrow('attributeId must be positive');
    });
  });
});
