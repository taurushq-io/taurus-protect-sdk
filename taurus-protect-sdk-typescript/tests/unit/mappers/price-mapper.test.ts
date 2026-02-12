/**
 * Unit tests for price mapper functions.
 */

import {
  priceFromDto,
  pricesFromDto,
  priceHistoryPointFromDto,
  priceHistoryPointsFromDto,
  conversionResultFromDto,
  conversionResultsFromDto,
} from '../../../src/mappers/price';

describe('priceFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      blockchain: 'ETH',
      currencyFrom: 'ETH',
      currencyTo: 'USD',
      decimals: '18',
      rate: '3500.50',
      changePercent24Hour: '-2.5',
      source: 'coingecko',
      creationDate: new Date('2024-01-01'),
      updateDate: new Date('2024-06-01'),
      currencyFromInfo: { id: 'c1', symbol: 'ETH', name: 'Ethereum' },
      currencyToInfo: { id: 'c2', symbol: 'USD', name: 'US Dollar' },
    };

    const result = priceFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.blockchain).toBe('ETH');
    expect(result!.currencyFrom).toBe('ETH');
    expect(result!.currencyTo).toBe('USD');
    expect(result!.decimals).toBe('18');
    expect(result!.rate).toBe('3500.50');
    expect(result!.changePercent24Hour).toBe('-2.5');
    expect(result!.source).toBe('coingecko');
    expect(result!.createdAt).toBeInstanceOf(Date);
    expect(result!.updatedAt).toBeInstanceOf(Date);
    expect(result!.currencyFromInfo).toBeDefined();
    expect(result!.currencyToInfo).toBeDefined();
  });

  it('should handle snake_case field names', () => {
    const dto = {
      currency_from: 'BTC',
      currency_to: 'EUR',
      change_percent_24_hour: '1.2',
      creation_date: '2024-03-01T00:00:00Z',
      update_date: '2024-04-01T00:00:00Z',
      currency_from_info: { id: 'c3', symbol: 'BTC' },
      currency_to_info: { id: 'c4', symbol: 'EUR' },
    };

    const result = priceFromDto(dto);

    expect(result!.currencyFrom).toBe('BTC');
    expect(result!.currencyTo).toBe('EUR');
    expect(result!.changePercent24Hour).toBe('1.2');
    expect(result!.currencyFromInfo).toBeDefined();
    expect(result!.currencyToInfo).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(priceFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(priceFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = priceFromDto({});
    expect(result).toBeDefined();
    expect(result!.rate).toBeUndefined();
  });
});

describe('pricesFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { currencyFrom: 'ETH', rate: '3500' },
      { currencyFrom: 'BTC', rate: '42000' },
    ];

    const result = pricesFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].rate).toBe('3500');
    expect(result[1].rate).toBe('42000');
  });

  it('should return empty array for null input', () => {
    expect(pricesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(pricesFromDto(undefined)).toEqual([]);
  });
});

describe('priceHistoryPointFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      periodStartDate: new Date('2024-01-01'),
      blockchain: 'ETH',
      currencyFrom: 'ETH',
      currencyTo: 'USD',
      high: '3600',
      low: '3400',
      open: '3450',
      close: '3550',
      volumeFrom: '10000',
      volumeTo: '35000000',
      changePercent: '2.9',
      currencyFromInfo: { id: 'c1', symbol: 'ETH' },
      currencyToInfo: { id: 'c2', symbol: 'USD' },
    };

    const result = priceHistoryPointFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.periodStartDate).toBeInstanceOf(Date);
    expect(result!.blockchain).toBe('ETH');
    expect(result!.high).toBe('3600');
    expect(result!.low).toBe('3400');
    expect(result!.open).toBe('3450');
    expect(result!.close).toBe('3550');
    expect(result!.volumeFrom).toBe('10000');
    expect(result!.volumeTo).toBe('35000000');
    expect(result!.changePercent).toBe('2.9');
  });

  it('should handle snake_case field names', () => {
    const dto = {
      period_start_date: '2024-02-01T00:00:00Z',
      currency_from: 'BTC',
      currency_to: 'EUR',
      volume_from: '500',
      volume_to: '21000000',
      change_percent: '-1.5',
    };

    const result = priceHistoryPointFromDto(dto);

    expect(result!.periodStartDate).toBeInstanceOf(Date);
    expect(result!.currencyFrom).toBe('BTC');
    expect(result!.volumeFrom).toBe('500');
    expect(result!.changePercent).toBe('-1.5');
  });

  it('should return undefined for null input', () => {
    expect(priceHistoryPointFromDto(null)).toBeUndefined();
  });
});

describe('priceHistoryPointsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { high: '100', low: '90' },
      { high: '110', low: '95' },
    ];

    const result = priceHistoryPointsFromDto(dtos);

    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(priceHistoryPointsFromDto(null)).toEqual([]);
  });
});

describe('conversionResultFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      symbol: 'ETH',
      value: '1000000000000000000',
      mainUnitValue: '1.0',
      currencyInfo: { id: 'c1', symbol: 'ETH', name: 'Ethereum' },
    };

    const result = conversionResultFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.symbol).toBe('ETH');
    expect(result!.value).toBe('1000000000000000000');
    expect(result!.mainUnitValue).toBe('1.0');
    expect(result!.currencyInfo).toBeDefined();
  });

  it('should handle snake_case field names', () => {
    const dto = {
      symbol: 'BTC',
      main_unit_value: '0.5',
      currency_info: { id: 'c2', symbol: 'BTC' },
    };

    const result = conversionResultFromDto(dto);

    expect(result!.mainUnitValue).toBe('0.5');
    expect(result!.currencyInfo).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(conversionResultFromDto(null)).toBeUndefined();
  });
});

describe('conversionResultsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { symbol: 'ETH', value: '1.0' },
      { symbol: 'BTC', value: '0.5' },
    ];

    const result = conversionResultsFromDto(dtos);

    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(conversionResultsFromDto(null)).toEqual([]);
  });
});
