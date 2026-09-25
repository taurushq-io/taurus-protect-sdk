/**
 * Unit tests for PriceService.
 */

import type { RulesContainerCache } from '../../../src/cache';
import type { DecodedRulesContainer } from '../../../src/models/governance-rules';
import { PriceService } from '../../../src/services/price-service';

// A container with no PRICEUPDATER: this tenant does not sign prices, so verification
// passes through. Signing itself is covered by the price-verifier tests.
function createMockRulesCache(): jest.Mocked<RulesContainerCache> {
  return {
    get: jest.fn().mockResolvedValue({
      users: [],
      groups: [],
    } as unknown as DecodedRulesContainer),
    clear: jest.fn(),
  } as unknown as jest.Mocked<RulesContainerCache>;
}
import { ValidationError } from '../../../src/errors';
import type { PricesApi } from '../../../src/internal/openapi/apis/PricesApi';

function createMockApi(): jest.Mocked<PricesApi> {
  return {
    priceServiceGetPrices: jest.fn(),
    priceServiceQueryPricesV2: jest.fn(),
    priceServiceGetPricesHistory: jest.fn(),
    priceServiceConvert: jest.fn(),
  } as unknown as jest.Mocked<PricesApi>;
}

describe('PriceService', () => {
  let mockApi: jest.Mocked<PricesApi>;
  let service: PriceService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new PriceService(mockApi, createMockRulesCache());
  });

  describe('list', () => {
    it('should list through QueryPricesV2, never the unpaged v1 list', async () => {
      mockApi.priceServiceQueryPricesV2.mockResolvedValue({
        result: [
          { currencyFrom: 'ETH', currencyTo: 'USD', rate: '3000.50' },
          { currencyFrom: 'BTC', currencyTo: 'USD', rate: '65000.00' },
        ],
        cursor: { currentPage: 'next-page', hasNext: true },
      } as never);

      const prices = await service.list();
      expect(prices.items).toHaveLength(2);
      expect(prices.pagination).toEqual({ pageSize: 20, nextCursor: 'next-page', hasMore: true });
      expect(mockApi.priceServiceGetPrices).not.toHaveBeenCalled();
    });

    it('should handle empty results', async () => {
      mockApi.priceServiceQueryPricesV2.mockResolvedValue({} as never);

      const prices = await service.list();
      expect(prices.items).toHaveLength(0);
      expect(prices.pagination.hasMore).toBe(false);
    });

    it.each([
      ['source only', { fromCurrencyId: 'c1' }, { from: { currencyFromId: 'c1' } }],
      [
        'source and targets',
        { fromCurrencyId: 'c1', toCurrencyIds: ['c2', 'c3'] },
        { fromTo: { currencyFromId: 'c1', currencyToIds: ['c2', 'c3'] } },
      ],
      ['targets only', { toCurrencyIds: ['c2'] }, { to: { currencyToIds: ['c2'] } }],
      ['no filter', {}, {}],
    ])('should map the currency filter (%s) onto the request oneof', async (_n, options, expected) => {
      mockApi.priceServiceQueryPricesV2.mockResolvedValue({} as never);

      await service.list({ ...options, onlyPrimary: true, sortOrder: 'ASC' });

      expect(mockApi.priceServiceQueryPricesV2).toHaveBeenCalledWith({
        body: {
          onlyPrimary: true,
          sortOrder: 'ASC',
          cursor: { currentPage: undefined, pageRequest: undefined, pageSize: '20' },
          ...expected,
        },
      });
    });
  });

  describe('getHistory', () => {
    it('should throw ValidationError when base is empty', async () => {
      await expect(
        service.getHistory({ base: '', quote: 'USD', limit: 100 })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when quote is empty', async () => {
      await expect(
        service.getHistory({ base: 'ETH', quote: '', limit: 100 })
      ).rejects.toThrow(ValidationError);
    });

    it('should return price history', async () => {
      mockApi.priceServiceGetPricesHistory.mockResolvedValue({
        result: [
          { price: '3000', date: '2024-01-01' },
          { price: '3100', date: '2024-01-02' },
        ],
      } as never);

      const history = await service.getHistory({ base: 'ETH', quote: 'USD', limit: 100 });
      expect(history).toHaveLength(2);
    });

    it('should always send a limit, 20 by default, at most 365', async () => {
      mockApi.priceServiceGetPricesHistory.mockResolvedValue({ result: [] } as never);

      await service.getHistory({ base: 'ETH', quote: 'USD' });
      expect(mockApi.priceServiceGetPricesHistory).toHaveBeenLastCalledWith({
        base: 'ETH',
        quote: 'USD',
        limit: '20',
      });

      await service.getHistory({ base: 'ETH', quote: 'USD', limit: 365 });
      expect(mockApi.priceServiceGetPricesHistory).toHaveBeenLastCalledWith({
        base: 'ETH',
        quote: 'USD',
        limit: '365',
      });

      await expect(service.getHistory({ base: 'ETH', quote: 'USD', limit: 366 })).rejects.toThrow(
        'limit must be at most 365, got 366'
      );
    });
  });

  describe('convert', () => {
    it('should throw ValidationError when currency is empty', async () => {
      await expect(
        service.convert({ currency: '', amount: '1000', targetCurrencyIds: ['USD'] })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when amount is empty', async () => {
      await expect(
        service.convert({ currency: 'ETH', amount: '', targetCurrencyIds: ['USD'] })
      ).rejects.toThrow(ValidationError);
    });

    it('should return conversion results', async () => {
      mockApi.priceServiceConvert.mockResolvedValue({
        result: [
          { currencyId: 'USD', value: '3000.50' },
        ],
      } as never);

      const results = await service.convert({
        currency: 'ETH',
        amount: '1000000000000000000',
        targetCurrencyIds: ['USD'],
      });
      expect(results).toHaveLength(1);
    });
  });
});
