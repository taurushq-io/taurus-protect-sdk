/**
 * Unit tests for request mapper functions.
 */

import { requestFromDto, requestsFromDto } from '../../../src/mappers/request';

describe('requestFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: '123',
      type: 'TRANSFER',
      status: 'APPROVED',
      tenantId: 't-1',
      currency: 'BTC',
      rule: 'default',
      externalRequestId: 'ext-1',
      requestBundleId: 'bundle-1',
      needsApprovalFrom: ['user-1'],
      creationDate: new Date('2024-01-01'),
      updateDate: new Date('2024-06-01'),
      metadata: {
        hash: 'abc123',
        payloadAsString: '{"key":"value"}',
      },
      approvers: null,
      signedRequests: [],
      currencyInfo: {
        id: 'BTC',
        symbol: 'BTC',
        name: 'Bitcoin',
        decimals: 8,
        blockchain: 'BTC',
        network: 'mainnet',
      },
    };

    const result = requestFromDto(dto as any);

    expect(result).toBeDefined();
    expect(result!.id).toBe(123);
    expect(result!.type).toBe('TRANSFER');
    expect(result!.status).toBe('APPROVED');
    expect(result!.tenantId).toBeDefined();
    expect(result!.currency).toBeDefined();
    expect(result!.metadata).toBeDefined();
    expect(result!.metadata!.hash).toBe('abc123');
    expect(result!.metadata!.payloadAsString).toBe('{"key":"value"}');
    expect(result!.currencyInfo).toBeDefined();
    expect(result!.currencyInfo!.symbol).toBe('BTC');
  });

  it('should return undefined for null input', () => {
    expect(requestFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(requestFromDto(undefined)).toBeUndefined();
  });

  it('should handle missing optional fields', () => {
    const dto = {
      id: '1',
      status: 'PENDING',
    };

    const result = requestFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.id).toBe(1);
  });

  it('should not include metadata when hash is missing', () => {
    const dto = {
      id: '1',
      metadata: {
        hash: undefined,
        payloadAsString: 'something',
      },
    };

    const result = requestFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.metadata).toBeUndefined();
  });

  it('should map signed requests', () => {
    const dto = {
      id: '1',
      signedRequests: [
        {
          id: 'sr-1',
          signedRequest: '0xabc',
          status: 'BROADCASTED',
          hash: 'txhash',
          block: 100,
          details: 'ok',
          creationDate: new Date('2024-01-01'),
          updateDate: new Date('2024-06-01'),
          broadcastDate: new Date('2024-01-02'),
          confirmationDate: new Date('2024-01-03'),
        },
      ],
    };

    const result = requestFromDto(dto as any);
    expect(result).toBeDefined();
    expect(result!.signedRequests).toHaveLength(1);
    expect(result!.signedRequests[0].id).toBe('sr-1');
    expect(result!.signedRequests[0].status).toBe('BROADCASTED');
    expect(result!.signedRequests[0].block).toBe(100);
  });
});

describe('requestsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: '1', status: 'PENDING', metadata: { hash: 'h1', payloadAsString: 'p1' } },
      { id: '2', status: 'APPROVED', metadata: { hash: 'h2', payloadAsString: 'p2' } },
    ];

    const result = requestsFromDto(dtos as any);
    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(requestsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(requestsFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty input', () => {
    expect(requestsFromDto([])).toEqual([]);
  });
});
