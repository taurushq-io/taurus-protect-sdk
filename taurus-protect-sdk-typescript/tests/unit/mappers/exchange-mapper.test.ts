/**
 * Unit tests for exchange mapper functions.
 */

import {
  exchangeFromDto,
  exchangesFromDto,
  exchangeCounterpartyFromDto,
  exchangeCounterpartiesFromDto,
  exchangeWithdrawalFeeFromDto,
} from '../../../src/mappers/exchange';

describe('exchangeFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'ex-1',
      exchange: 'Binance',
      account: 'main-account',
      currency: 'BTC',
      type: 'spot',
      totalBalance: '10.5',
      status: 'active',
      container: 'default',
      label: 'Primary',
      displayLabel: 'Primary Binance',
      baseCurrencyValuation: '420000',
      hasWLA: true,
      currencyInfo: { id: 'c1', symbol: 'BTC', name: 'Bitcoin' },
      creationDate: new Date('2024-01-01'),
      updateDate: new Date('2024-06-01'),
    };

    const result = exchangeFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('ex-1');
    expect(result!.exchange).toBe('Binance');
    expect(result!.account).toBe('main-account');
    expect(result!.currency).toBe('BTC');
    expect(result!.type).toBe('spot');
    expect(result!.totalBalance).toBe('10.5');
    expect(result!.status).toBe('active');
    expect(result!.container).toBe('default');
    expect(result!.label).toBe('Primary');
    expect(result!.displayLabel).toBe('Primary Binance');
    expect(result!.baseCurrencyValuation).toBe('420000');
    expect(result!.hasWLA).toBe(true);
    expect(result!.currencyInfo).toBeDefined();
    expect(result!.creationDate).toBeInstanceOf(Date);
    expect(result!.updateDate).toBeInstanceOf(Date);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      id: 'ex-2',
      total_balance: '20.0',
      display_label: 'Label',
      base_currency_valuation: '100',
      has_wla: false,
      currency_info: { id: 'c2', symbol: 'ETH' },
      creation_date: '2024-03-01T00:00:00Z',
      update_date: '2024-04-01T00:00:00Z',
    };

    const result = exchangeFromDto(dto);

    expect(result!.totalBalance).toBe('20.0');
    expect(result!.displayLabel).toBe('Label');
    expect(result!.baseCurrencyValuation).toBe('100');
    expect(result!.hasWLA).toBe(false);
    expect(result!.currencyInfo).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(exchangeFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(exchangeFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = exchangeFromDto({});
    expect(result).toBeDefined();
    expect(result!.id).toBeUndefined();
  });
});

describe('exchangesFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 'ex-1', exchange: 'Binance' },
      { id: 'ex-2', exchange: 'Kraken' },
    ];

    const result = exchangesFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].exchange).toBe('Binance');
    expect(result[1].exchange).toBe('Kraken');
  });

  it('should return empty array for null input', () => {
    expect(exchangesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(exchangesFromDto(undefined)).toEqual([]);
  });
});

describe('exchangeCounterpartyFromDto', () => {
  it('should map counterparty fields', () => {
    const dto = {
      name: 'Counterparty A',
      baseCurrencyValuation: '50000',
    };

    const result = exchangeCounterpartyFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.name).toBe('Counterparty A');
    expect(result!.baseCurrencyValuation).toBe('50000');
  });

  it('should return undefined for null input', () => {
    expect(exchangeCounterpartyFromDto(null)).toBeUndefined();
  });
});

describe('exchangeCounterpartiesFromDto', () => {
  it('should map array of counterparties', () => {
    const dtos = [
      { name: 'A', baseCurrencyValuation: '100' },
      { name: 'B', baseCurrencyValuation: '200' },
    ];

    const result = exchangeCounterpartiesFromDto(dtos);

    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(exchangeCounterpartiesFromDto(null)).toEqual([]);
  });
});

describe('exchangeWithdrawalFeeFromDto', () => {
  it('should map fee from result field', () => {
    const dto = { result: '0.001' };

    const result = exchangeWithdrawalFeeFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.fee).toBe('0.001');
  });

  it('should map fee from fee field', () => {
    const dto = { fee: '0.002' };

    const result = exchangeWithdrawalFeeFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.fee).toBe('0.002');
  });

  it('should return undefined when no fee present', () => {
    const dto = {};

    const result = exchangeWithdrawalFeeFromDto(dto);

    expect(result).toBeUndefined();
  });

  it('should return undefined for null input', () => {
    expect(exchangeWithdrawalFeeFromDto(null)).toBeUndefined();
  });
});
