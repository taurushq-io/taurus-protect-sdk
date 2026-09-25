/**
 * Unit tests for AssetService.
 */

import { ConfigurationError, ValidationError } from '../../../src/errors';
import type { RulesContainerCache } from '../../../src/cache';
import type { DecodedRulesContainer } from '../../../src/models/governance-rules';
import { AssetService, type AssetHolderReaders } from '../../../src/services/asset-service';

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

// The holder readers are exercised through the transport in asset-holders.test.ts; the
// replies here carry no internal or whitelisted holders, so nothing resolves them.
const NO_HOLDERS: AssetHolderReaders = {
  addresses: () => {
    throw new Error('the address reader must not be used');
  },
  whitelistedAddresses: () => {
    throw new Error('the whitelisted-address reader must not be used');
  },
};

describe('AssetService', () => {
  describe('constructor', () => {
    it('should throw ConfigurationError when the holder readers are not provided', () => {
      expect(
        () =>
          new AssetService(
            {} as never,
            createMockRulesCache(),
            {} as never,
            undefined as unknown as AssetHolderReaders
          )
      ).toThrow(ConfigurationError);
    });

    it('should throw ConfigurationError when rulesCache is not provided', () => {
      const mockApi = {} as never;
      expect(
        () => new AssetService(mockApi, undefined as unknown as RulesContainerCache, {} as never, NO_HOLDERS)
      ).toThrow(ConfigurationError);
      expect(
        () => new AssetService(mockApi, null as unknown as RulesContainerCache, {} as never, NO_HOLDERS)
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
      const service = new AssetService(mockApi as any, rulesCache, {} as never, NO_HOLDERS);

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
      const service = new AssetService(mockApi as any, createMockRulesCache(), {} as never, NO_HOLDERS);

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

      const service = new AssetService(mockApi as any, createMockRulesCache(), {} as never, NO_HOLDERS);
      const result = await service.getAssetAddresses({ currency: 'ETH' });

      // Paged through requestCursor only: the legacy limit/cursor fields are never sent.
      expect(mockApi.walletServiceGetAssetAddresses).toHaveBeenCalledWith({
        body: {
          asset: { currency: 'ETH' },
          walletId: undefined,
          addressId: undefined,
          addresses: undefined,
          requestCursor: { currentPage: undefined, pageRequest: undefined, pageSize: '20' },
        },
      });
      expect(result.items).toHaveLength(2);
      expect(result.items[0].id).toBe('1');
      expect(result.items[0].address).toBe('0x123');
      expect(result.pagination).toEqual({
        pageSize: 20,
        nextCursor: '',
        hasMore: false,
        totalItems: 2,
      });
    });

    it('should pass optional filters', async () => {
      const mockApi = {
        walletServiceGetAssetAddresses: jest.fn().mockResolvedValue({
          addresses: [],
        }),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache(), {} as never, NO_HOLDERS);
      await service.getAssetAddresses({
        currency: 'BTC',
        walletId: 'wallet-123',
        addressId: 'addr-456',
        addresses: ['bc1q'],
        pageSize: 50,
        cursor: 'next-page',
      });

      expect(mockApi.walletServiceGetAssetAddresses).toHaveBeenCalledWith({
        body: {
          asset: { currency: 'BTC' },
          walletId: 'wallet-123',
          addressId: 'addr-456',
          addresses: ['bc1q'],
          requestCursor: { currentPage: 'next-page', pageRequest: 'NEXT', pageSize: '50' },
        },
      });
    });

    it('should throw ValidationError when currency is empty', async () => {
      const mockApi = {
        walletServiceGetAssetAddresses: jest.fn(),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache(), {} as never, NO_HOLDERS);

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

      const service = new AssetService(mockApi as any, createMockRulesCache(), {} as never, NO_HOLDERS);
      const result = await service.getAssetAddresses({ currency: 'USDC' });

      expect(result.items).toEqual([]);
      expect(result.pagination.totalItems).toBe(0);
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

      const service = new AssetService(mockApi as any, createMockRulesCache(), {} as never, NO_HOLDERS);
      const result = await service.getAssetWallets({ currency: 'ETH' });

      expect(mockApi.walletServiceGetAssetWallets).toHaveBeenCalledWith({
        body: {
          asset: { currency: 'ETH' },
          walletId: undefined,
          walletName: undefined,
          requestCursor: { currentPage: undefined, pageRequest: undefined, pageSize: '20' },
        },
      });
      expect(result.items).toHaveLength(2);
      expect(result.items[0].id).toBe('1');
      expect(result.items[0].name).toBe('Wallet 1');
      expect(result.pagination.totalItems).toBe(2);
    });

    it('should pass optional filters', async () => {
      const mockApi = {
        walletServiceGetAssetWallets: jest.fn().mockResolvedValue({
          wallets: [],
        }),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache(), {} as never, NO_HOLDERS);
      await service.getAssetWallets({
        currency: 'USDC',
        walletId: 'wallet-123',
        walletName: 'My Wallet',
        pageSize: 100,
      });

      expect(mockApi.walletServiceGetAssetWallets).toHaveBeenCalledWith({
        body: {
          asset: { currency: 'USDC' },
          walletId: 'wallet-123',
          walletName: 'My Wallet',
          requestCursor: { currentPage: undefined, pageRequest: undefined, pageSize: '100' },
        },
      });
    });

    it('should throw ValidationError when currency is empty', async () => {
      const mockApi = {
        walletServiceGetAssetWallets: jest.fn(),
      };

      const service = new AssetService(mockApi as any, createMockRulesCache(), {} as never, NO_HOLDERS);

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

      const service = new AssetService(mockApi as any, createMockRulesCache(), {} as never, NO_HOLDERS);
      const result = await service.getAssetWallets({ currency: 'USDC' });

      expect(result.items).toEqual([]);
    });
  });

  describe('v2 registry reads', () => {
    it('should reject an empty assetId before sending', async () => {
      const assetV2Api = {
        assetServiceV2QueryAssetAddressesV2: jest.fn(),
        assetServiceV2ListAssetOperationsV2: jest.fn(),
      };
      const service = new AssetService({} as never, createMockRulesCache(), assetV2Api as never, NO_HOLDERS);
      await expect(service.queryAssetAddresses('')).rejects.toThrow('assetId is required');
      await expect(service.listAssetOperations(' ')).rejects.toThrow('assetId is required');
      expect(assetV2Api.assetServiceV2QueryAssetAddressesV2).not.toHaveBeenCalled();
      expect(assetV2Api.assetServiceV2ListAssetOperationsV2).not.toHaveBeenCalled();
    });

    it('should send every queryAssets filter with the page size in the body', async () => {
      const assetV2Api = {
        assetServiceV2QueryAssetsV2: jest.fn().mockResolvedValue({ result: [{ id: 'a1' }] }),
      };
      const service = new AssetService({} as never, createMockRulesCache(), assetV2Api as never, NO_HOLDERS);
      const page = await service.queryAssets({
        blockchain: 'CANTON',
        network: 'mainnet',
        symbol: 'TKN',
        contractAddress: 'c',
        label: 'l',
        currencyName: 'n',
      });
      expect(assetV2Api.assetServiceV2QueryAssetsV2).toHaveBeenCalledWith({
        body: {
          cursor: { currentPage: undefined, pageRequest: undefined, pageSize: '20' },
          blockchain: 'CANTON',
          network: 'mainnet',
          symbol: 'TKN',
          contractAddress: 'c',
          label: 'l',
          currencyName: 'n',
        },
      });
      expect(page.items.map((a) => a.id)).toEqual(['a1']);
    });
  });
});
