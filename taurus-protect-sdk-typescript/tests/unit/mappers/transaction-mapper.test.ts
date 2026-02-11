/**
 * Unit tests for transaction mapper functions.
 */

import {
  transactionFromDto,
  transactionsFromDto,
  transactionAddressInfoFromDto,
  transactionCurrencyInfoFromDto,
  transactionAttributeFromDto,
} from '../../../src/mappers/transaction';
import { TransactionStatus } from '../../../src/models/transaction';

describe('transactionFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'tx-1',
      direction: 'OUTGOING',
      currency: 'BTC',
      blockchain: 'BTC',
      network: 'mainnet',
      hash: '0xabc',
      block: '100',
      confirmationBlock: '106',
      amount: '1.5',
      amountMainUnit: '1.5',
      fee: '0.001',
      feeMainUnit: '0.001',
      type: 'TRANSFER',
      status: 'SUCCESS',
      isConfirmed: true,
      sources: [{ address: '0xsrc', amount: '1.5' }],
      destinations: [{ address: '0xdst', amount: '1.5' }],
      transactionId: 'tid-1',
      uniqueId: 'uid-1',
      requestId: 'req-1',
      requestVisible: true,
      receptionDate: new Date('2024-01-01'),
      confirmationDate: new Date('2024-01-02'),
      attributes: [{ key: 'tag', value: 'test' }],
      forkNumber: '0',
      currencyInfo: {
        id: 'BTC',
        symbol: 'BTC',
        name: 'Bitcoin',
        decimals: 8,
        blockchain: 'BTC',
        network: 'mainnet',
      },
    };

    const result = transactionFromDto(dto as any);

    expect(result).toBeDefined();
    expect(result!.id).toBe('tx-1');
    expect(result!.direction).toBeDefined();
    expect(result!.currency).toBeDefined();
    expect(result!.hash).toBeDefined();
    expect(result!.amount).toBeDefined();
    expect(result!.status).toBe(TransactionStatus.SUCCESS);
    expect(result!.isConfirmed).toBe(true);
  });

  it('should return undefined for null input', () => {
    expect(transactionFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(transactionFromDto(undefined)).toBeUndefined();
  });

  it('should handle minimal DTO', () => {
    const dto = { id: '1' };
    const result = transactionFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.id).toBe('1');
  });

  it('should parse known status values', () => {
    const dto = { id: '1', status: 'PENDING' };
    const result = transactionFromDto(dto as any);
    expect(result!.status).toBe(TransactionStatus.PENDING);
  });

  it('should preserve unknown status values as strings', () => {
    const dto = { id: '1', status: 'CUSTOM_STATUS' };
    const result = transactionFromDto(dto as any);
    expect(result!.status).toBe('CUSTOM_STATUS');
  });
});

describe('transactionsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: '1', currency: 'BTC' },
      { id: '2', currency: 'ETH' },
    ];
    const result = transactionsFromDto(dtos as any);
    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(transactionsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(transactionsFromDto(undefined)).toEqual([]);
  });
});

describe('transactionAddressInfoFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      address: '0xabc',
      label: 'My Address',
      container: 'wallet-1',
      customerId: 'cust-1',
      amount: '1.0',
      amountMainUnit: '1.0',
      type: 'INTERNAL',
      idx: '0',
      internalAddressId: 'addr-1',
      whitelistedAddressId: 'wl-1',
    };
    const result = transactionAddressInfoFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.address).toBeDefined();
    expect(result!.label).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(transactionAddressInfoFromDto(null)).toBeUndefined();
  });
});

describe('transactionCurrencyInfoFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      id: 'BTC',
      symbol: 'BTC',
      name: 'Bitcoin',
      decimals: 8,
      blockchain: 'BTC',
      network: 'mainnet',
    };
    const result = transactionCurrencyInfoFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.symbol).toBe('BTC');
    expect(result!.decimals).toBe(8);
  });

  it('should return undefined for null input', () => {
    expect(transactionCurrencyInfoFromDto(null)).toBeUndefined();
  });
});

describe('transactionAttributeFromDto', () => {
  it('should map fields', () => {
    const dto = { key: 'tag', value: 'test' };
    const result = transactionAttributeFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.key).toBe('tag');
    expect(result!.value).toBe('test');
  });

  it('should return undefined for null input', () => {
    expect(transactionAttributeFromDto(null)).toBeUndefined();
  });
});
