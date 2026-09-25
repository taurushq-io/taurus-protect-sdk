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
    it('should return balances and the server total', async () => {
      mockApi.walletServiceGetBalances.mockResolvedValue({
        balances: [
          { asset: { currency: 'ETH' }, balance: { totalConfirmed: '10.5' } },
          { asset: { currency: 'BTC' }, balance: { totalConfirmed: '0.5' } },
        ],
        total: '2',
      });

      const result = await service.list();

      expect(result.items.map((b) => [b.currency, b.balance])).toEqual([
        ['ETH', '10.5'],
        ['BTC', '0.5'],
      ]);
      expect(result.pagination).toEqual({
        pageSize: 20,
        nextCursor: '',
        hasMore: false,
        totalItems: 2,
      });
    });

    it('should pass the currency and token filters', async () => {
      mockApi.walletServiceGetBalances.mockResolvedValue({ balances: [] });

      await service.list({ currency: 'ETH', tokenId: 'tok-1' });

      expect(mockApi.walletServiceGetBalances).toHaveBeenCalledWith(
        expect.objectContaining({
          currency: 'ETH',
          tokenId: 'tok-1',
        })
      );
    });

    it('should throw ValidationError when the page size is out of bounds', async () => {
      await expect(service.list({ pageSize: 101 })).rejects.toThrow(ValidationError);
      await expect(service.list({ pageSize: 101 })).rejects.toThrow('pageSize must be at most 100, got 101');
      expect(mockApi.walletServiceGetBalances).not.toHaveBeenCalled();
    });

    it('should page through requestCursor only, never the legacy limit/cursor', async () => {
      mockApi.walletServiceGetBalances.mockResolvedValue({ balances: [] });

      await service.list({ cursor: 'next-page' });

      expect(mockApi.walletServiceGetBalances).toHaveBeenCalledWith({
        currency: undefined,
        tokenId: undefined,
        requestCursorCurrentPage: 'next-page',
        requestCursorPageRequest: 'NEXT',
        requestCursorPageSize: '20',
      });
    });
  });

  describe('listNFTCollections', () => {
    it('should read the rows from `balances`, the field the reply carries', async () => {
      mockApi.walletServiceGetNFTCollectionBalances.mockResolvedValue({
        balances: [
          { currencyInfo: { name: 'CryptoKitties' }, balance: { totalConfirmed: '10' } },
        ],
        cursor: { currentPage: 'next-page', hasNext: true },
      });

      const result = await service.listNFTCollections({
        blockchain: 'ETH',
        network: 'mainnet',
      });

      expect(result.items).toEqual([expect.objectContaining({ name: 'CryptoKitties', count: 10 })]);
      expect(result.pagination).toEqual({
        pageSize: 20,
        nextCursor: 'next-page',
        hasMore: true,
      });
    });

    it('should not require blockchain or network: the server filters are optional', async () => {
      mockApi.walletServiceGetNFTCollectionBalances.mockResolvedValue({});

      const result = await service.listNFTCollections();

      expect(result.items).toEqual([]);
      expect(mockApi.walletServiceGetNFTCollectionBalances).toHaveBeenCalledWith(
        expect.objectContaining({ cursorPageSize: '20' })
      );
    });

    it('should throw ValidationError when the page size is out of bounds', async () => {
      await expect(
        service.listNFTCollections({ blockchain: 'ETH', network: 'mainnet', pageSize: -1 })
      ).rejects.toThrow(ValidationError);
    });
  });
});
