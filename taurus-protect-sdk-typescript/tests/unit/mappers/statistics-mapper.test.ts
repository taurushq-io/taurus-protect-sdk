/**
 * Unit tests for statistics mapper functions.
 */

import { portfolioStatisticsFromDto } from '../../../src/mappers/statistics';

describe('portfolioStatisticsFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      avgBalancePerAddress: '5000',
      addressesCount: '100',
      walletsCount: '10',
      totalBalance: '500000',
      totalBalanceBaseCurrency: '50000000',
    };

    const result = portfolioStatisticsFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.avgBalancePerAddress).toBe('5000');
    expect(result!.addressesCount).toBe('100');
    expect(result!.walletsCount).toBe('10');
    expect(result!.totalBalance).toBe('500000');
    expect(result!.totalBalanceBaseCurrency).toBe('50000000');
  });

  it('should handle snake_case field names', () => {
    const dto = {
      avg_balance_per_address: '2500',
      addresses_count: '50',
      wallets_count: '5',
      total_balance: '250000',
      total_balance_base_currency: '25000000',
    };

    const result = portfolioStatisticsFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.avgBalancePerAddress).toBe('2500');
    expect(result!.addressesCount).toBe('50');
  });

  it('should return undefined for null input', () => {
    expect(portfolioStatisticsFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(portfolioStatisticsFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = portfolioStatisticsFromDto({});
    expect(result).toBeDefined();
  });
});
