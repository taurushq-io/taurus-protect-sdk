/**
 * Unit tests for FiatService.
 */

import { FiatService } from '../../../src/services/fiat-service';
import { ValidationError } from '../../../src/errors';
import type { FiatApi } from '../../../src/internal/openapi/apis/FiatApi';

function createMockApi(): jest.Mocked<FiatApi> {
  return {
    fiatProviderServiceGetFiatProviders: jest.fn(),
    fiatProviderServiceGetFiatProviderAccount: jest.fn(),
    fiatProviderServiceGetFiatProviderAccounts: jest.fn(),
    fiatProviderServiceGetFiatProviderCounterpartyAccount: jest.fn(),
    fiatProviderServiceGetFiatProviderCounterpartyAccounts: jest.fn(),
    fiatProviderServiceGetFiatProviderOperation: jest.fn(),
    fiatProviderServiceGetFiatProviderOperations: jest.fn(),
  } as unknown as jest.Mocked<FiatApi>;
}

describe('FiatService', () => {
  let mockApi: jest.Mocked<FiatApi>;
  let service: FiatService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new FiatService(mockApi);
  });

  describe('getFiatProviders', () => {
    it('should return fiat providers', async () => {
      mockApi.fiatProviderServiceGetFiatProviders.mockResolvedValue({
        fiatProviders: [
          { provider: 'circle', baseCurrencyValuation: 'USD' },
        ],
      } as never);

      const providers = await service.getFiatProviders();
      expect(providers).toBeDefined();
      expect(providers.length).toBeGreaterThanOrEqual(0);
    });

    it('should handle empty results', async () => {
      mockApi.fiatProviderServiceGetFiatProviders.mockResolvedValue({
        fiatProviders: [],
      } as never);

      const providers = await service.getFiatProviders();
      expect(providers).toHaveLength(0);
    });
  });

  describe('getFiatProviderAccount', () => {
    it('should throw ValidationError when id is empty', async () => {
      await expect(service.getFiatProviderAccount('')).rejects.toThrow(ValidationError);
      await expect(service.getFiatProviderAccount('')).rejects.toThrow('id is required');
    });

    it('should throw ValidationError when id is whitespace', async () => {
      await expect(service.getFiatProviderAccount('  ')).rejects.toThrow(ValidationError);
    });

    it('should return account for valid id', async () => {
      mockApi.fiatProviderServiceGetFiatProviderAccount.mockResolvedValue({
        result: {
          id: 'account-123',
          accountName: 'Main Account',
          totalBalance: '1000000',
          currencyId: 'USD',
        },
      } as never);

      const account = await service.getFiatProviderAccount('account-123');
      expect(account).toBeDefined();
      expect(mockApi.fiatProviderServiceGetFiatProviderAccount).toHaveBeenCalledWith({
        id: 'account-123',
      });
    });
  });

  describe('getFiatProviderAccounts', () => {
    it('should throw ValidationError when provider is empty', async () => {
      await expect(
        service.getFiatProviderAccounts({ provider: '', label: 'main' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getFiatProviderAccounts({ provider: '', label: 'main' })
      ).rejects.toThrow('provider is required');
    });

    it('should throw ValidationError when label is empty', async () => {
      await expect(
        service.getFiatProviderAccounts({ provider: 'circle', label: '' })
      ).rejects.toThrow(ValidationError);
      await expect(
        service.getFiatProviderAccounts({ provider: 'circle', label: '' })
      ).rejects.toThrow('label is required');
    });

    it('should return accounts for valid options', async () => {
      mockApi.fiatProviderServiceGetFiatProviderAccounts.mockResolvedValue({
        accounts: [
          { id: '1', accountName: 'Account 1' },
        ],
      } as never);

      const result = await service.getFiatProviderAccounts({
        provider: 'circle',
        label: 'main',
      });

      expect(result).toBeDefined();
      expect(mockApi.fiatProviderServiceGetFiatProviderAccounts).toHaveBeenCalledWith(
        expect.objectContaining({
          provider: 'circle',
          label: 'main',
        })
      );
    });
  });

  describe('getFiatProviderCounterpartyAccount', () => {
    it('should throw ValidationError when id is empty', async () => {
      await expect(service.getFiatProviderCounterpartyAccount('')).rejects.toThrow(
        ValidationError
      );
      await expect(service.getFiatProviderCounterpartyAccount('')).rejects.toThrow(
        'id is required'
      );
    });

    it('should return counterparty account for valid id', async () => {
      mockApi.fiatProviderServiceGetFiatProviderCounterpartyAccount.mockResolvedValue({
        result: {
          id: 'cp-123',
          counterpartyName: 'Partner Corp',
        },
      } as never);

      const account = await service.getFiatProviderCounterpartyAccount('cp-123');
      expect(account).toBeDefined();
    });
  });

  describe('getFiatProviderCounterpartyAccounts', () => {
    it('should throw ValidationError when provider is empty', async () => {
      await expect(
        service.getFiatProviderCounterpartyAccounts({ provider: '', label: 'main' })
      ).rejects.toThrow(ValidationError);
    });

    it('should throw ValidationError when label is empty', async () => {
      await expect(
        service.getFiatProviderCounterpartyAccounts({ provider: 'cubnet', label: '' })
      ).rejects.toThrow(ValidationError);
    });

    it('should return counterparty accounts for valid options', async () => {
      mockApi.fiatProviderServiceGetFiatProviderCounterpartyAccounts.mockResolvedValue({
        accounts: [],
      } as never);

      const result = await service.getFiatProviderCounterpartyAccounts({
        provider: 'cubnet',
        label: 'main',
      });

      expect(result).toBeDefined();
    });
  });

  describe('getFiatProviderOperation', () => {
    it('should throw ValidationError when id is empty', async () => {
      await expect(service.getFiatProviderOperation('')).rejects.toThrow(ValidationError);
      await expect(service.getFiatProviderOperation('')).rejects.toThrow('id is required');
    });

    it('should return operation for valid id', async () => {
      mockApi.fiatProviderServiceGetFiatProviderOperation.mockResolvedValue({
        result: {
          id: 'op-123',
          status: 'COMPLETED',
          amount: '5000',
        },
      } as never);

      const operation = await service.getFiatProviderOperation('op-123');
      expect(operation).toBeDefined();
    });
  });

  describe('getFiatProviderOperations', () => {
    it('should return operations without options', async () => {
      mockApi.fiatProviderServiceGetFiatProviderOperations.mockResolvedValue({
        operations: [],
      } as never);

      const result = await service.getFiatProviderOperations();
      expect(result).toBeDefined();
    });

    it('should pass filter options to API', async () => {
      mockApi.fiatProviderServiceGetFiatProviderOperations.mockResolvedValue({
        operations: [],
      } as never);

      await service.getFiatProviderOperations({
        provider: 'circle',
        label: 'main',
        sortOrder: 'DESC',
      });

      expect(mockApi.fiatProviderServiceGetFiatProviderOperations).toHaveBeenCalledWith(
        expect.objectContaining({
          provider: 'circle',
          label: 'main',
          sortOrder: 'DESC',
        })
      );
    });
  });
});
