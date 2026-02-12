/**
 * Unit tests for WalletService.
 */

import { WalletService } from '../../../src/services/wallet-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { WalletsApi } from '../../../src/internal/openapi/apis/WalletsApi';

function createMockWalletsApi(): jest.Mocked<WalletsApi> {
  return {
    walletServiceGetWalletV2: jest.fn(),
    walletServiceGetWalletsV2: jest.fn(),
    walletServiceCreateWallet: jest.fn(),
    walletServiceCreateWalletAttributes: jest.fn(),
    walletServiceDeleteWalletAttribute: jest.fn(),
    walletServiceGetWalletBalanceHistory: jest.fn(),
  } as unknown as jest.Mocked<WalletsApi>;
}

describe('WalletService', () => {
  let mockApi: jest.Mocked<WalletsApi>;
  let service: WalletService;

  beforeEach(() => {
    mockApi = createMockWalletsApi();
    service = new WalletService(mockApi);
  });

  describe('get', () => {
    it('should return a wallet for a valid ID', async () => {
      mockApi.walletServiceGetWalletV2.mockResolvedValue({
        result: {
          id: '123',
          name: 'Test Wallet',
          blockchain: 'ETH',
          network: 'mainnet',
          currency: 'ETH',
          creationDate: new Date('2024-01-01'),
        },
      });

      const wallet = await service.get(123);

      expect(wallet).toBeDefined();
      expect(wallet.id).toBe('123');
      expect(wallet.name).toBe('Test Wallet');
      expect(wallet.blockchain).toBe('ETH');
      expect(mockApi.walletServiceGetWalletV2).toHaveBeenCalledWith({ id: '123' });
    });

    it('should throw ValidationError when walletId is 0', async () => {
      await expect(service.get(0)).rejects.toThrow(ValidationError);
      await expect(service.get(0)).rejects.toThrow('walletId must be positive');
    });

    it('should throw ValidationError when walletId is negative', async () => {
      await expect(service.get(-1)).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when wallet is not found', async () => {
      mockApi.walletServiceGetWalletV2.mockResolvedValue({
        result: undefined,
      });

      await expect(service.get(999)).rejects.toThrow(NotFoundError);
    });
  });

  describe('list', () => {
    it('should return paginated wallets', async () => {
      mockApi.walletServiceGetWalletsV2.mockResolvedValue({
        result: [
          { id: '1', name: 'Wallet 1', currency: 'ETH' },
          { id: '2', name: 'Wallet 2', currency: 'BTC' },
        ],
        totalItems: '10',
        offset: '0',
      });

      const result = await service.list({ limit: 50 });

      expect(result.items).toHaveLength(2);
      expect(result.items[0].name).toBe('Wallet 1');
      expect(result.items[1].name).toBe('Wallet 2');
      expect(result.pagination.totalItems).toBe(10);
      expect(result.pagination.limit).toBe(50);
    });

    it('should throw ValidationError when limit is 0', async () => {
      await expect(service.list({ limit: 0 })).rejects.toThrow(ValidationError);
      await expect(service.list({ limit: 0 })).rejects.toThrow('limit must be positive');
    });

    it('should throw ValidationError when limit is negative', async () => {
      await expect(service.list({ limit: -1 })).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when offset is negative', async () => {
      await expect(service.list({ offset: -1 })).rejects.toThrow(ValidationError);
      await expect(service.list({ offset: -1 })).rejects.toThrow('offset cannot be negative');
    });

    it('should use default limit and offset when not provided', async () => {
      mockApi.walletServiceGetWalletsV2.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      await service.list();

      expect(mockApi.walletServiceGetWalletsV2).toHaveBeenCalledWith(
        expect.objectContaining({
          limit: '50',
        })
      );
    });
  });

  describe('create', () => {
    it('should create a wallet with valid request', async () => {
      mockApi.walletServiceCreateWallet.mockResolvedValue({
        result: {
          id: '456',
          name: 'New Wallet',
          blockchain: 'ETH',
          currency: 'ETH',
        },
      });

      const wallet = await service.create({
        name: 'New Wallet',
        blockchain: 'ETH',
        network: 'mainnet',
      });

      expect(wallet).toBeDefined();
      expect(wallet.name).toBe('New Wallet');
      expect(mockApi.walletServiceCreateWallet).toHaveBeenCalledWith({
        body: expect.objectContaining({
          name: 'New Wallet',
          blockchain: 'ETH',
          network: 'mainnet',
        }),
      });
    });

    it('should throw ValidationError when name is missing', async () => {
      await expect(
        service.create({ name: '', blockchain: 'ETH' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.create({ name: '', blockchain: 'ETH' })
      ).rejects.toThrow('name is required');
    });

    it('should throw ValidationError when neither currency nor blockchain is provided', async () => {
      await expect(
        service.create({ name: 'Test' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.create({ name: 'Test' })
      ).rejects.toThrow('Either currency or blockchain must be provided');
    });
  });

  describe('createAttribute', () => {
    it('should create an attribute for a wallet', async () => {
      mockApi.walletServiceCreateWalletAttributes.mockResolvedValue({});

      await service.createAttribute(123, 'department', 'treasury');

      expect(mockApi.walletServiceCreateWalletAttributes).toHaveBeenCalledWith({
        walletId: '123',
        body: {
          attributes: [{ key: 'department', value: 'treasury' }],
        },
      });
    });

    it('should throw ValidationError when walletId is invalid', async () => {
      await expect(service.createAttribute(0, 'key', 'value')).rejects.toThrow(
        ValidationError
      );
    });

    it('should throw ValidationError when key is empty', async () => {
      await expect(service.createAttribute(1, '', 'value')).rejects.toThrow(
        ValidationError
      );
      await expect(service.createAttribute(1, '', 'value')).rejects.toThrow(
        'key is required'
      );
    });

    it('should throw ValidationError when value is empty', async () => {
      await expect(service.createAttribute(1, 'key', '')).rejects.toThrow(
        ValidationError
      );
      await expect(service.createAttribute(1, 'key', '')).rejects.toThrow(
        'value is required'
      );
    });
  });

  describe('deleteAttribute', () => {
    it('should delete an attribute from a wallet', async () => {
      mockApi.walletServiceDeleteWalletAttribute.mockResolvedValue({});

      await service.deleteAttribute(123, 'attr-456');

      expect(mockApi.walletServiceDeleteWalletAttribute).toHaveBeenCalledWith({
        walletId: '123',
        id: 'attr-456',
      });
    });

    it('should throw ValidationError when walletId is invalid', async () => {
      await expect(service.deleteAttribute(0, 'attr-1')).rejects.toThrow(
        ValidationError
      );
    });

    it('should throw ValidationError when attributeId is empty', async () => {
      await expect(service.deleteAttribute(1, '')).rejects.toThrow(
        ValidationError
      );
      await expect(service.deleteAttribute(1, '')).rejects.toThrow(
        'attributeId is required'
      );
    });
  });

  describe('getBalanceHistory', () => {
    it('should return balance history for a wallet', async () => {
      mockApi.walletServiceGetWalletBalanceHistory.mockResolvedValue({
        result: [
          {
            pointDate: new Date('2024-01-01'),
            balance: { totalConfirmed: '1000' },
          },
          {
            pointDate: new Date('2024-01-02'),
            balance: { totalConfirmed: '1500' },
          },
        ],
      });

      const history = await service.getBalanceHistory(123, 24);

      expect(history).toHaveLength(2);
      expect(mockApi.walletServiceGetWalletBalanceHistory).toHaveBeenCalledWith({
        id: '123',
        intervalHours: '24',
      });
    });

    it('should throw ValidationError when walletId is invalid', async () => {
      await expect(service.getBalanceHistory(0, 24)).rejects.toThrow(
        ValidationError
      );
    });

    it('should throw ValidationError when intervalHours is invalid', async () => {
      await expect(service.getBalanceHistory(1, 0)).rejects.toThrow(
        ValidationError
      );
      await expect(service.getBalanceHistory(1, 0)).rejects.toThrow(
        'intervalHours must be positive'
      );
    });
  });
});
