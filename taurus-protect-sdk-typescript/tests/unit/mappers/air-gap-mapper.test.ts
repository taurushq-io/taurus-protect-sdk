/**
 * Unit tests for air gap mapper functions.
 */

import {
  toGetOutgoingAirGapRequest,
  toGetOutgoingAirGapAddressRequest,
  toSubmitIncomingAirGapRequest,
} from '../../../src/mappers/air-gap';

describe('toGetOutgoingAirGapRequest', () => {
  it('should map request IDs and signature', () => {
    const options = {
      requestIds: ['req-1', 'req-2'],
      signature: 'sig-abc',
    };

    const result = toGetOutgoingAirGapRequest(options);

    expect(result.requests).toBeDefined();
    expect(result.requests!.ids).toEqual(['req-1', 'req-2']);
    expect(result.requests!.signature).toBe('sig-abc');
  });

  it('should handle empty request IDs', () => {
    const options = {
      requestIds: [] as string[],
      signature: 'sig',
    };

    const result = toGetOutgoingAirGapRequest(options);

    expect(result.requests!.ids).toEqual([]);
  });
});

describe('toGetOutgoingAirGapAddressRequest', () => {
  it('should map address IDs', () => {
    const options = {
      addressIds: ['addr-1', 'addr-2'],
    };

    const result = toGetOutgoingAirGapAddressRequest(options);

    expect(result.addresses).toBeDefined();
    expect(result.addresses!.ids).toEqual(['addr-1', 'addr-2']);
  });

  it('should handle empty address IDs', () => {
    const options = {
      addressIds: [] as string[],
    };

    const result = toGetOutgoingAirGapAddressRequest(options);

    expect(result.addresses!.ids).toEqual([]);
  });
});

describe('toSubmitIncomingAirGapRequest', () => {
  it('should map payload and signature', () => {
    const options = {
      payload: 'base64-payload',
      signature: 'sig-xyz',
    };

    const result = toSubmitIncomingAirGapRequest(options);

    expect(result.payload).toBe('base64-payload');
    expect(result.signature).toBe('sig-xyz');
  });
});
