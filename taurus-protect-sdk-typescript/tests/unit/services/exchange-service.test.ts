/**
 * Unit tests for ExchangeService.
 */

import { ExchangeService } from '../../../src/services/exchange-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import type { ExchangeApi } from '../../../src/internal/openapi/apis/ExchangeApi';

function createMockApi(): jest.Mocked<ExchangeApi> {
  return {
    exchangeServiceGetExchanges: jest.fn(),
    exchangeServiceGetExchange: jest.fn(),
    exchangeServiceGetExchangeCounterparties: jest.fn(),
    exchangeServiceGetExchangeWithdrawalFee: jest.fn(),
    exchangeServiceExportExchanges: jest.fn(),
  } as unknown as jest.Mocked<ExchangeApi>;
}

describe('ExchangeService', () => {
  let mockApi: jest.Mocked<ExchangeApi>;
  let service: ExchangeService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new ExchangeService(mockApi);
  });

  describe('list', () => {
    it('should return exchanges', async () => {
      mockApi.exchangeServiceGetExchanges.mockResolvedValue({
        result: [{ id: 'ex-1', name: 'Exchange A' }],
      } as never);

      const result = await service.list();
      expect(result.items).toBeDefined();
    });

    it('should handle empty results', async () => {
      mockApi.exchangeServiceGetExchanges.mockResolvedValue({
        result: [],
      } as never);

      const result = await service.list();
      expect(result.items).toHaveLength(0);
    });
  });

  describe('get', () => {
    it('should throw ValidationError when id is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
    });

    it('should return exchange when found', async () => {
      mockApi.exchangeServiceGetExchange.mockResolvedValue({
        result: { id: 'ex-123', name: 'Test Exchange' },
      } as never);

      const exchange = await service.get('ex-123');
      expect(exchange).toBeDefined();
    });
  });

  describe('getCounterparties', () => {
    it('should return counterparties', async () => {
      mockApi.exchangeServiceGetExchangeCounterparties.mockResolvedValue({
        result: [{ id: 'cp-1', name: 'Counterparty A' }],
      } as never);

      const counterparties = await service.getCounterparties();
      expect(counterparties).toBeDefined();
    });
  });
});
