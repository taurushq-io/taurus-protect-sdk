/**
 * Unit tests for AssetService.
 */

import { ValidationError } from '../../../src/errors';
import { AssetService } from '../../../src/services/asset-service';

describe('AssetService', () => {
  describe('getAssetAddresses', () => {
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

      const service = new AssetService(mockApi as any);
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

      const service = new AssetService(mockApi as any);
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

      const service = new AssetService(mockApi as any);

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

      const service = new AssetService(mockApi as any);
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

      const service = new AssetService(mockApi as any);
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

      const service = new AssetService(mockApi as any);
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

      const service = new AssetService(mockApi as any);

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

      const service = new AssetService(mockApi as any);
      const result = await service.getAssetWallets({ currency: 'USDC' });

      expect(result).toEqual([]);
    });
  });
});
