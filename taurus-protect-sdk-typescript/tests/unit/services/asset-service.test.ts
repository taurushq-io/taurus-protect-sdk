/**
 * Unit tests for AssetService.
 */

import { ConfigurationError, ValidationError } from '../../../src/errors';
import type { RulesContainerCache } from '../../../src/cache';
import type { DecodedRulesContainer } from '../../../src/models/governance-rules';
import { AssetService } from '../../../src/services/asset-service';

// The verification itself is covered by the address-signature-verifier tests; what
// matters here is that getAssetAddresses runs it at all, which the constructor and the
// call below pin.
jest.mock('../../../src/helpers', () => ({
  ...jest.requireActual('../../../src/helpers'),
  verifyAddressSignature: jest.fn().mockReturnValue(true),
}));

function createMockRulesCache(): jest.Mocked<RulesContainerCache> {
  return {
    get: jest.fn().mockResolvedValue({
      users: [],
      groups: [],
    } as unknown as DecodedRulesContainer),
    clear: jest.fn(),
  } as unknown as jest.Mocked<RulesContainerCache>;
}

describe('AssetService', () => {
  describe('constructor', () => {
    it('should throw ConfigurationError when rulesCache is not provided', () => {
      const mockApi = {} as never;
      expect(
        () => new AssetService(mockApi, undefined as unknown as RulesContainerCache)
      ).toThrow(ConfigurationError);
      expect(
        () => new AssetService(mockApi, null as unknown as RulesContainerCache)
      ).toThrow(ConfigurationError);
    });
  });

  describe('getAssetAddresses', () => {
    // getAssetAddresses returns the same Address entity AddressService does and used to
    // hand it over unverified, so the mandatory verification there was reachable around.
    it('should verify every address signature it returns', async () => {
      const { verifyAddressSignature } = jest.requireMock('../../../src/helpers');
      (verifyAddressSignature as jest.Mock).mockClear();

      const mockApi = {
        walletServiceGetAssetAddresses: jest.fn().mockResolvedValue({
          addresses: [
            { id: '1', address: '0x123', walletId: '1', currency: 'ETH', signature: 'sig1' },
            { id: '2', address: '0x456', walletId: '2', currency: 'ETH', signature: 'sig2' },
          ],
          totalItems: '2',
        }),
      };
      const rulesCache = createMockRulesCache();
      const service = new AssetService(mockApi as any, rulesCache);

      await service.getAssetAddresses({ currency: 'ETH' });

      expect(rulesCache.get).toHaveBeenCalled();
      expect(verifyAddressSignature).toHaveBeenCalledTimes(2);
      expect(verifyAddressSignature).toHaveBeenCalledWith(
        '0x123',
        'sig1',
        expect.anything(),
        '1'
      );
    });

    it('should propagate a verification failure rather than returning the address', async () => {
      const { verifyAddressSignature } = jest.requireMock('../../../src/helpers');
      (verifyAddressSignature as jest.Mock).mockImplementationOnce(() => {
        throw new Error('address signature verification failed');
      });

      const mockApi = {
        walletServiceGetAssetAddresses: jest.fn().mockResolvedValue({
          addresses: [{ id: '1', address: '0xEVIL', walletId: '1', currency: 'ETH' }],
          totalItems: '1',
        }),
      };
      const service = new AssetService(mockApi as any, createMockRulesCache());

      await expect(service.getAssetAddresses({ currency: 'ETH' })).rejects.toThrow();
    });

    it('should return addresses for a given currency', async () => {
      const mockAddresses = [
        { id: '1', address: '0x123', walletId: '1', currency: 'ETH' },
        { id: '2', address: '0x456', walletId: '2', currency: 'ETH' },
      ];

      const mockApi = {
        walletServiceGetAssetAddresses: jest.fn().mockResolvedValue({
          addresses: mockAddresses,
          totalItems: '2',
        }),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache());
      const result = await service.getAssetAddresses({ currency: 'ETH' });

      expect(mockApi.walletServiceGetAssetAddresses).toHaveBeenCalledWith({
        body: {
          asset: { currency: 'ETH' },
          walletId: undefined,
          addressId: undefined,
          limit: undefined,
        },
      });
      expect(result).toHaveLength(2);
      expect(result[0].id).toBe('1');
      expect(result[0].address).toBe('0x123');
    });

    it('should pass optional filters', async () => {
      const mockApi = {
        walletServiceGetAssetAddresses: jest.fn().mockResolvedValue({
          addresses: [],
        }),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache());
      await service.getAssetAddresses({
        currency: 'BTC',
        walletId: 'wallet-123',
        addressId: 'addr-456',
        limit: '50',
      });

      expect(mockApi.walletServiceGetAssetAddresses).toHaveBeenCalledWith({
        body: {
          asset: { currency: 'BTC' },
          walletId: 'wallet-123',
          addressId: 'addr-456',
          limit: '50',
        },
      });
    });

    it('should throw ValidationError when currency is empty', async () => {
      const mockApi = {
        walletServiceGetAssetAddresses: jest.fn(),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache());

      await expect(service.getAssetAddresses({ currency: '' }))
        .rejects.toThrow(ValidationError);
      await expect(service.getAssetAddresses({ currency: '   ' }))
        .rejects.toThrow(ValidationError);
    });

    it('should return empty array when no addresses found', async () => {
      const mockApi = {
        walletServiceGetAssetAddresses: jest.fn().mockResolvedValue({
          addresses: undefined,
        }),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache());
      const result = await service.getAssetAddresses({ currency: 'USDC' });

      expect(result).toEqual([]);
    });
  });

  describe('getAssetWallets', () => {
    it('should return wallets for a given currency', async () => {
      const mockWallets = [
        { id: '1', name: 'Wallet 1', currency: 'ETH' },
        { id: '2', name: 'Wallet 2', currency: 'ETH' },
      ];

      const mockApi = {
        walletServiceGetAssetWallets: jest.fn().mockResolvedValue({
          wallets: mockWallets,
          totalItems: '2',
        }),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache());
      const result = await service.getAssetWallets({ currency: 'ETH' });

      expect(mockApi.walletServiceGetAssetWallets).toHaveBeenCalledWith({
        body: {
          asset: { currency: 'ETH' },
          walletId: undefined,
          walletName: undefined,
          limit: undefined,
        },
      });
      expect(result).toHaveLength(2);
      expect(result[0].id).toBe('1');
      expect(result[0].name).toBe('Wallet 1');
    });

    it('should pass optional filters', async () => {
      const mockApi = {
        walletServiceGetAssetWallets: jest.fn().mockResolvedValue({
          wallets: [],
        }),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache());
      await service.getAssetWallets({
        currency: 'USDC',
        walletId: 'wallet-123',
        walletName: 'My Wallet',
        limit: '100',
      });

      expect(mockApi.walletServiceGetAssetWallets).toHaveBeenCalledWith({
        body: {
          asset: { currency: 'USDC' },
          walletId: 'wallet-123',
          walletName: 'My Wallet',
          limit: '100',
        },
      });
    });

    it('should throw ValidationError when currency is empty', async () => {
      const mockApi = {
        walletServiceGetAssetWallets: jest.fn(),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache());

      await expect(service.getAssetWallets({ currency: '' }))
        .rejects.toThrow(ValidationError);
      await expect(service.getAssetWallets({ currency: '   ' }))
        .rejects.toThrow(ValidationError);
    });

    it('should return empty array when no wallets found', async () => {
      const mockApi = {
        walletServiceGetAssetWallets: jest.fn().mockResolvedValue({
          wallets: undefined,
        }),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache());
      const result = await service.getAssetWallets({ currency: 'USDC' });

      expect(result).toEqual([]);
    });
  });
});
