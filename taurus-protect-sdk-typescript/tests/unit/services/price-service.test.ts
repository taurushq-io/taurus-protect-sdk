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
    it('should return current prices', async () => {
      mockApi.priceServiceGetPrices.mockResolvedValue({
        result: [
          { currencyId: 'ETH', price: '3000.50' },
          { currencyId: 'BTC', price: '65000.00' },
        ],
      } as never);

      const prices = await service.list();
      expect(prices).toHaveLength(2);
    });

    it('should handle empty results', async () => {
      mockApi.priceServiceGetPrices.mockResolvedValue({
        result: [],
      } as never);

      const prices = await service.list();
      expect(prices).toHaveLength(0);
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
