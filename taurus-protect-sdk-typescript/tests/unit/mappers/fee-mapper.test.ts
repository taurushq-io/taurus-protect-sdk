/**
 * Unit tests for fee mapper functions.
 */

import {
  feeFromDto,
  feesFromDto,
  feeV2FromDto,
  feesV2FromDto,
} from '../../../src/mappers/fee';

describe('feeFromDto', () => {
  it('should map key and value', () => {
    const dto = { key: 'gasPrice', value: '20000000000' };

    const result = feeFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.key).toBe('gasPrice');
    expect(result!.value).toBe('20000000000');
  });

  it('should return undefined for null input', () => {
    expect(feeFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(feeFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = feeFromDto({});
    expect(result).toBeDefined();
    expect(result!.key).toBeUndefined();
    expect(result!.value).toBeUndefined();
  });
});

describe('feesFromDto', () => {
  it('should map array of fee DTOs', () => {
    const dtos = [
      { key: 'gasPrice', value: '20' },
      { key: 'gasLimit', value: '21000' },
    ];

    const result = feesFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].key).toBe('gasPrice');
    expect(result[1].key).toBe('gasLimit');
  });

  it('should return empty array for null input', () => {
    expect(feesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(feesFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty array', () => {
    expect(feesFromDto([])).toEqual([]);
  });
});

describe('feeV2FromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      currencyId: 'c-1',
      value: '0.001',
      denom: 'ETH',
      currencyInfo: { id: 'c-1', symbol: 'ETH', name: 'Ethereum' },
      updateDate: new Date('2024-06-01'),
    };

    const result = feeV2FromDto(dto);

    expect(result).toBeDefined();
    expect(result!.currencyId).toBe('c-1');
    expect(result!.value).toBe('0.001');
    expect(result!.denom).toBe('ETH');
    expect(result!.currencyInfo).toBeDefined();
    expect(result!.currencyInfo!.symbol).toBe('ETH');
    expect(result!.updateDate).toBeInstanceOf(Date);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      currency_id: 'c-2',
      value: '0.002',
      currency_info: { id: 'c-2', symbol: 'BTC' },
      update_date: '2024-07-01T00:00:00Z',
    };

    const result = feeV2FromDto(dto);

    expect(result!.currencyId).toBe('c-2');
    expect(result!.currencyInfo).toBeDefined();
    expect(result!.updateDate).toBeInstanceOf(Date);
  });

  it('should return undefined for null input', () => {
    expect(feeV2FromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(feeV2FromDto(undefined)).toBeUndefined();
  });
});

describe('feesV2FromDto', () => {
  it('should map array of v2 fee DTOs', () => {
    const dtos = [
      { currencyId: 'c-1', value: '0.001' },
      { currencyId: 'c-2', value: '0.002' },
    ];

    const result = feesV2FromDto(dtos);

    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(feesV2FromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(feesV2FromDto(undefined)).toEqual([]);
  });
});
