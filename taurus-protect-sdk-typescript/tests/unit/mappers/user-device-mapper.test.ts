/**
 * Unit tests for user device mapper functions.
 */

import {
  userDevicePairingFromDto,
  userDevicePairingInfoFromDto,
} from '../../../src/mappers/user-device';

describe('userDevicePairingFromDto', () => {
  it('should map fields from DTO', () => {
    const dto = {
      pairingId: 'pair-1',
    };
    const result = userDevicePairingFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.pairingId).toBe('pair-1');
  });

  it('should handle pairingID (uppercase) field', () => {
    const dto = {
      pairingID: 'pair-2',
    };
    const result = userDevicePairingFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.pairingId).toBe('pair-2');
  });

  it('should handle snake_case field', () => {
    const dto = {
      pairing_id: 'pair-3',
    };
    const result = userDevicePairingFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.pairingId).toBe('pair-3');
  });

  it('should return undefined for null input', () => {
    expect(userDevicePairingFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(userDevicePairingFromDto(undefined)).toBeUndefined();
  });

  it('should return undefined when pairingId is missing', () => {
    const dto = { someOtherField: 'value' };
    expect(userDevicePairingFromDto(dto)).toBeUndefined();
  });
});

describe('userDevicePairingInfoFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      pairingId: 'pair-1',
      status: 'ACTIVE',
      apiKey: 'key-123',
    };
    const result = userDevicePairingInfoFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.pairingId).toBe('pair-1');
    expect(result!.status).toBe('ACTIVE');
    expect(result!.apiKey).toBe('key-123');
  });

  it('should handle snake_case fields', () => {
    const dto = {
      pairing_id: 'pair-2',
      status: 'PENDING',
      api_key: 'key-456',
    };
    const result = userDevicePairingInfoFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.pairingId).toBe('pair-2');
    expect(result!.apiKey).toBe('key-456');
  });

  it('should return undefined for null input', () => {
    expect(userDevicePairingInfoFromDto(null)).toBeUndefined();
  });

  it('should return undefined when pairingId is missing', () => {
    const dto = { status: 'ACTIVE' };
    expect(userDevicePairingInfoFromDto(dto)).toBeUndefined();
  });

  it('should return undefined when status is missing', () => {
    const dto = { pairingId: 'pair-1' };
    expect(userDevicePairingInfoFromDto(dto)).toBeUndefined();
  });
});
