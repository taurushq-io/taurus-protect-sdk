/**
 * Unit tests for CurrencyService.
 */

import { CurrencyService } from '../../../src/services/currency-service';
import { NotFoundError, ValidationError } from '../../../src/errors';
import type { CurrenciesApi } from '../../../src/internal/openapi/apis/CurrenciesApi';

function createMockApi(): jest.Mocked<CurrenciesApi> {
  return {
    walletServiceGetCurrencies: jest.fn(),
    walletServiceGetCurrency: jest.fn(),
    walletServiceGetBaseCurrency: jest.fn(),
  } as unknown as jest.Mocked<CurrenciesApi>;
}

describe('CurrencyService', () => {
  let mockApi: jest.Mocked<CurrenciesApi>;
  let service: CurrencyService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new CurrencyService(mockApi);
  });

  describe('list', () => {
    it('should return currencies', async () => {
      mockApi.walletServiceGetCurrencies.mockResolvedValue({
        currencies: [
          { id: 'eth', symbol: 'ETH', name: 'Ethereum', decimals: '18', blockchain: 'ETH', network: 'mainnet' },
          { id: 'btc', symbol: 'BTC', name: 'Bitcoin', decimals: '8', blockchain: 'BTC', network: 'mainnet' },
        ],
      } as never);

      const currencies = await service.list();
      expect(currencies).toBeDefined();
      expect(currencies.length).toBeGreaterThanOrEqual(0);
    });

    it('should pass filter options to API', async () => {
      mockApi.walletServiceGetCurrencies.mockResolvedValue({
        currencies: [],
      } as never);

      await service.list({ showDisabled: true, includeLogo: true });

      expect(mockApi.walletServiceGetCurrencies).toHaveBeenCalledWith({
        showDisabled: true,
        includeLogo: true,
      });
    });

    it('should handle empty response', async () => {
      mockApi.walletServiceGetCurrencies.mockResolvedValue({
        currencies: [],
      } as never);

      const currencies = await service.list();
      expect(currencies).toHaveLength(0);
    });
  });

  describe('get', () => {
    it('should throw ValidationError when currencyId is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('currencyId is required');
    });

    it('should throw ValidationError when currencyId is whitespace', async () => {
      await expect(service.get('  ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when currency is not found', async () => {
      mockApi.walletServiceGetCurrencies.mockResolvedValue({
        currencies: [
          { id: 'eth', symbol: 'ETH', name: 'Ethereum' },
        ],
      } as never);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });

    it('should return currency when found', async () => {
      mockApi.walletServiceGetCurrencies.mockResolvedValue({
        currencies: [
          { id: 'eth', symbol: 'ETH', name: 'Ethereum', decimals: '18' },
          { id: 'btc', symbol: 'BTC', name: 'Bitcoin', decimals: '8' },
        ],
      } as never);

      const currency = await service.get('eth');
      expect(currency).toBeDefined();
      expect(currency.id).toBe('eth');
      expect(currency.symbol).toBe('ETH');
    });
  });

  describe('getByBlockchain', () => {
    it('should throw ValidationError when blockchain is empty', async () => {
      await expect(
        service.getByBlockchain({ blockchain: '', network: 'mainnet' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getByBlockchain({ blockchain: '', network: 'mainnet' })
      ).rejects.toThrow('blockchain is required');
    });

    it('should throw ValidationError when network is empty', async () => {
      await expect(
        service.getByBlockchain({ blockchain: 'ETH', network: '' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getByBlockchain({ blockchain: 'ETH', network: '' })
      ).rejects.toThrow('network is required');
    });

    it('should throw NotFoundError when currency is not found', async () => {
      mockApi.walletServiceGetCurrency.mockResolvedValue({
        currency: undefined,
      } as never);

      await expect(
        service.getByBlockchain({ blockchain: 'ETH', network: 'mainnet' })
      ).rejects.toThrow(NotFoundError);
    });

    it('should return currency for valid blockchain and network', async () => {
      mockApi.walletServiceGetCurrency.mockResolvedValue({
        currency: {
          id: 'eth-mainnet',
          symbol: 'ETH',
          name: 'Ethereum',
          decimals: '18',
          blockchain: 'ETH',
          network: 'mainnet',
        },
      } as never);

      const currency = await service.getByBlockchain({
        blockchain: 'ETH',
        network: 'mainnet',
      });

      expect(currency).toBeDefined();
      expect(currency.symbol).toBe('ETH');
      expect(mockApi.walletServiceGetCurrency).toHaveBeenCalledWith(
        expect.objectContaining({
          uniqueCurrencyFilterBlockchain: 'ETH',
          uniqueCurrencyFilterNetwork: 'mainnet',
        })
      );
    });

    it('should pass contractAddress and tokenId options', async () => {
      mockApi.walletServiceGetCurrency.mockResolvedValue({
        currency: {
          id: 'usdc',
          symbol: 'USDC',
          name: 'USD Coin',
        },
      } as never);

      await service.getByBlockchain({
        blockchain: 'ETH',
        network: 'mainnet',
        contractAddress: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
        tokenId: '0',
      });

      expect(mockApi.walletServiceGetCurrency).toHaveBeenCalledWith(
        expect.objectContaining({
          uniqueCurrencyFilterTokenContractAddress: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
          uniqueCurrencyFilterTokenID: '0',
        })
      );
    });
  });

  describe('getBaseCurrency', () => {
    it('should return base currency', async () => {
      mockApi.walletServiceGetBaseCurrency.mockResolvedValue({
        currency: {
          id: 'usd',
          symbol: 'USD',
          name: 'US Dollar',
        },
      } as never);

      const currency = await service.getBaseCurrency();
      expect(currency).toBeDefined();
      expect(currency.symbol).toBe('USD');
    });

    it('should throw NotFoundError when no base currency is configured', async () => {
      mockApi.walletServiceGetBaseCurrency.mockResolvedValue({
        currency: undefined,
      } as never);

      await expect(service.getBaseCurrency()).rejects.toThrow(NotFoundError);
      await expect(service.getBaseCurrency()).rejects.toThrow('Base currency not configured');
    });
  });
});
