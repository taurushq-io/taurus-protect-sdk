/**
 * Unit tests for reservation mapper functions.
 */

import {
  reservationFromDto,
  reservationsFromDto,
  reservationUtxoFromDto,
} from '../../../src/mappers/reservation';

describe('reservationFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'res-1',
      amount: '1000',
      creationDate: new Date('2024-01-01'),
      kind: 'UTXO',
      comment: 'Reserved for transaction',
      addressid: 'addr-1',
      address: '0xABC123',
      currencyInfo: { id: 'c1', symbol: 'BTC', name: 'Bitcoin' },
      resourceId: 'req-1',
      resourceType: 'REQUEST',
    };

    const result = reservationFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('res-1');
    expect(result!.amount).toBe('1000');
    expect(result!.creationDate).toBeInstanceOf(Date);
    expect(result!.kind).toBe('UTXO');
    expect(result!.comment).toBe('Reserved for transaction');
    expect(result!.addressId).toBe('addr-1');
    expect(result!.address).toBe('0xABC123');
    expect(result!.currencyInfo).toBeDefined();
    expect(result!.currencyInfo!.symbol).toBe('BTC');
    expect(result!.resourceId).toBe('req-1');
    expect(result!.resourceType).toBe('REQUEST');
  });

  it('should handle addressId (camelCase) instead of addressid', () => {
    const dto = { id: 'res-2', addressId: 'addr-2' };

    const result = reservationFromDto(dto);

    expect(result!.addressId).toBe('addr-2');
  });

  it('should handle snake_case field names', () => {
    const dto = {
      id: 'res-3',
      creation_date: '2024-03-01T00:00:00Z',
      address_id: 'addr-3',
      currency_info: { id: 'c2', symbol: 'ETH' },
      resource_id: 'req-3',
      resource_type: 'REQUEST',
    };

    const result = reservationFromDto(dto);

    expect(result!.addressId).toBe('addr-3');
    expect(result!.resourceId).toBe('req-3');
    expect(result!.resourceType).toBe('REQUEST');
  });

  it('should return undefined for null input', () => {
    expect(reservationFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(reservationFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = reservationFromDto({});
    expect(result).toBeDefined();
    expect(result!.id).toBeUndefined();
  });
});

describe('reservationsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 'res-1', amount: '100' },
      { id: 'res-2', amount: '200' },
    ];

    const result = reservationsFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].amount).toBe('100');
    expect(result[1].amount).toBe('200');
  });

  it('should return empty array for null input', () => {
    expect(reservationsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(reservationsFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty array', () => {
    expect(reservationsFromDto([])).toEqual([]);
  });
});

describe('reservationUtxoFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      id: 'utxo-1',
      hash: 'txhash123',
      outputIndex: 0,
      script: '76a914...',
      value: '50000',
      blockHeight: '800000',
      reservedByRequestId: 'req-1',
      reservationId: 'res-1',
      valueString: '0.0005',
    };

    const result = reservationUtxoFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('utxo-1');
    expect(result!.hash).toBe('txhash123');
    expect(result!.outputIndex).toBe(0);
    expect(result!.script).toBe('76a914...');
    expect(result!.value).toBe('50000');
    expect(result!.blockHeight).toBe('800000');
    expect(result!.reservedByRequestId).toBe('req-1');
    expect(result!.reservationId).toBe('res-1');
    expect(result!.valueString).toBe('0.0005');
  });

  it('should handle snake_case field names', () => {
    const dto = {
      id: 'utxo-2',
      output_index: 1,
      block_height: '900000',
      reserved_by_request_id: 'req-2',
      reservation_id: 'res-2',
      value_string: '0.001',
    };

    const result = reservationUtxoFromDto(dto);

    expect(result!.outputIndex).toBe(1);
    expect(result!.blockHeight).toBe('900000');
    expect(result!.reservedByRequestId).toBe('req-2');
  });

  it('should return undefined for null input', () => {
    expect(reservationUtxoFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(reservationUtxoFromDto(undefined)).toBeUndefined();
  });
});
