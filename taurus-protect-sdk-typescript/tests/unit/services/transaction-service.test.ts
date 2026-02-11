/**
 * Unit tests for TransactionService.
 */

import { TransactionService } from '../../../src/services/transaction-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { TransactionsApi } from '../../../src/internal/openapi/apis/TransactionsApi';

function createMockTransactionsApi(): jest.Mocked<TransactionsApi> {
  return {
    transactionServiceGetTransactions: jest.fn(),
    transactionServiceExportTransactions: jest.fn(),
  } as unknown as jest.Mocked<TransactionsApi>;
}

describe('TransactionService', () => {
  let mockApi: jest.Mocked<TransactionsApi>;
  let service: TransactionService;

  beforeEach(() => {
    mockApi = createMockTransactionsApi();
    service = new TransactionService(mockApi);
  });

  describe('get', () => {
    it('should return a transaction for a valid ID', async () => {
      mockApi.transactionServiceGetTransactions.mockResolvedValue({
        result: [
          {
            id: '12345',
            hash: '0xabcdef',
            type: 'incoming',
            status: 'CONFIRMED',
            currency: 'ETH',
            amount: '1.5',
          },
        ],
        totalItems: '1',
      });

      const tx = await service.get('12345');

      expect(tx).toBeDefined();
      expect(tx.id).toBe('12345');
      expect(mockApi.transactionServiceGetTransactions).toHaveBeenCalledWith({
        ids: ['12345'],
        limit: '1',
      });
    });

    it('should throw ValidationError when transactionId is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('transactionId is required');
    });

    it('should throw ValidationError when transactionId is whitespace', async () => {
      await expect(service.get('   ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when transaction is not found', async () => {
      mockApi.transactionServiceGetTransactions.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });
  });

  describe('list', () => {
    it('should return paginated transactions', async () => {
      mockApi.transactionServiceGetTransactions.mockResolvedValue({
        result: [
          { id: '1', hash: '0x111', currency: 'ETH', amount: '1.0' },
          { id: '2', hash: '0x222', currency: 'ETH', amount: '2.0' },
        ],
        totalItems: '100',
      });

      const result = await service.list({ limit: 50 });

      expect(result.items).toHaveLength(2);
      expect(result.pagination.totalItems).toBe(100);
      expect(result.pagination.limit).toBe(50);
    });

    it('should pass currency filter to API', async () => {
      mockApi.transactionServiceGetTransactions.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      await service.list({ currency: 'BTC', limit: 10 });

      expect(mockApi.transactionServiceGetTransactions).toHaveBeenCalledWith(
        expect.objectContaining({
          currency: 'BTC',
          limit: '10',
        })
      );
    });

    it('should throw ValidationError when limit is invalid', async () => {
      await expect(service.list({ limit: 0 })).rejects.toThrow(ValidationError);
      await expect(service.list({ limit: 0 })).rejects.toThrow('limit must be positive');
    });

    it('should throw ValidationError when offset is negative', async () => {
      await expect(service.list({ offset: -1 })).rejects.toThrow(ValidationError);
      await expect(service.list({ offset: -1 })).rejects.toThrow('offset cannot be negative');
    });

    it('should use defaults when options are not provided', async () => {
      mockApi.transactionServiceGetTransactions.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      await service.list();

      expect(mockApi.transactionServiceGetTransactions).toHaveBeenCalledWith(
        expect.objectContaining({
          limit: '50',
        })
      );
    });
  });

  describe('getByHash', () => {
    it('should return a transaction by hash', async () => {
      mockApi.transactionServiceGetTransactions.mockResolvedValue({
        result: [
          { id: '1', hash: '0xabc', currency: 'ETH' },
        ],
        totalItems: '1',
      });

      const tx = await service.getByHash('0xabc');

      expect(tx).toBeDefined();
      expect(mockApi.transactionServiceGetTransactions).toHaveBeenCalledWith({
        hashes: ['0xabc'],
        limit: '1',
      });
    });

    it('should throw ValidationError when hash is empty', async () => {
      await expect(service.getByHash('')).rejects.toThrow(ValidationError);
      await expect(service.getByHash('')).rejects.toThrow('txHash is required');
    });

    it('should throw NotFoundError when hash is not found', async () => {
      mockApi.transactionServiceGetTransactions.mockResolvedValue({
        result: [],
        totalItems: '0',
      });

      await expect(service.getByHash('0xnonexistent')).rejects.toThrow(NotFoundError);
    });
  });

  describe('listByAddress', () => {
    it('should return transactions for an address', async () => {
      mockApi.transactionServiceGetTransactions.mockResolvedValue({
        result: [
          { id: '1', hash: '0xaaa', currency: 'ETH' },
        ],
        totalItems: '1',
      });

      const result = await service.listByAddress('0xmyaddress');

      expect(result.items).toHaveLength(1);
      expect(mockApi.transactionServiceGetTransactions).toHaveBeenCalledWith(
        expect.objectContaining({
          address: '0xmyaddress',
        })
      );
    });

    it('should throw ValidationError when address is empty', async () => {
      await expect(service.listByAddress('')).rejects.toThrow(ValidationError);
      await expect(service.listByAddress('')).rejects.toThrow('address is required');
    });
  });
});
