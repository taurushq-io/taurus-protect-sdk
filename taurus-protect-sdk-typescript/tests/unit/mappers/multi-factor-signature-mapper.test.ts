/**
 * Unit tests for multi-factor signature mapper functions.
 */

import {
  multiFactorSignatureEntityTypeFromDto,
  multiFactorSignatureEntityTypeToDto,
  multiFactorSignatureInfoFromDto,
} from '../../../src/mappers/multi-factor-signature';
import { MultiFactorSignatureEntityType } from '../../../src/models/multi-factor-signature';

describe('multiFactorSignatureEntityTypeFromDto', () => {
  it('should map REQUEST type', () => {
    const result = multiFactorSignatureEntityTypeFromDto('REQUEST');
    expect(result).toBe(MultiFactorSignatureEntityType.REQUEST);
  });

  it('should map WHITELISTED_ADDRESS type', () => {
    const result = multiFactorSignatureEntityTypeFromDto('WHITELISTED_ADDRESS');
    expect(result).toBe(MultiFactorSignatureEntityType.WHITELISTED_ADDRESS);
  });

  it('should map WHITELISTED_CONTRACT type', () => {
    const result = multiFactorSignatureEntityTypeFromDto('WHITELISTED_CONTRACT');
    expect(result).toBe(MultiFactorSignatureEntityType.WHITELISTED_CONTRACT);
  });

  // An unknown kind used to come back as REQUEST, so an entity of a kind this client has never
  // seen was presented as a request. It now stays the raw wire value, as in every other SDK.
  it('should pass an unknown type through verbatim', () => {
    const result = multiFactorSignatureEntityTypeFromDto('FUTURE_KIND' as any);
    expect(result).toBe('FUTURE_KIND');
  });

  it('should keep undefined absent', () => {
    const result = multiFactorSignatureEntityTypeFromDto(undefined);
    expect(result).toBeUndefined();
  });

  it('should keep null absent', () => {
    const result = multiFactorSignatureEntityTypeFromDto(null);
    expect(result).toBeUndefined();
  });
});

describe('multiFactorSignatureEntityTypeToDto', () => {
  it('should convert REQUEST to DTO', () => {
    const result = multiFactorSignatureEntityTypeToDto(MultiFactorSignatureEntityType.REQUEST);
    expect(result).toBe('REQUEST');
  });

  it('should convert WHITELISTED_ADDRESS to DTO', () => {
    const result = multiFactorSignatureEntityTypeToDto(
      MultiFactorSignatureEntityType.WHITELISTED_ADDRESS
    );
    expect(result).toBe('WHITELISTED_ADDRESS');
  });

  it('should convert WHITELISTED_CONTRACT to DTO', () => {
    const result = multiFactorSignatureEntityTypeToDto(
      MultiFactorSignatureEntityType.WHITELISTED_CONTRACT
    );
    expect(result).toBe('WHITELISTED_CONTRACT');
  });
});

describe('multiFactorSignatureInfoFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'mfs-1',
      payloadToSign: ['payload1', 'payload2'],
      entityType: 'REQUEST' as const,
    };

    const result = multiFactorSignatureInfoFromDto(dto as any);

    expect(result).not.toBeNull();
    expect(result!.id).toBe('mfs-1');
    expect(result!.payloadToSign).toEqual(['payload1', 'payload2']);
    expect(result!.entityType).toBe(MultiFactorSignatureEntityType.REQUEST);
  });

  it('should use defaults for missing fields', () => {
    const dto = {} as any;

    const result = multiFactorSignatureInfoFromDto(dto);

    expect(result).not.toBeNull();
    expect(result!.id).toBe('');
    expect(result!.payloadToSign).toEqual([]);
  });

  it('should return null for null input', () => {
    expect(multiFactorSignatureInfoFromDto(null)).toBeNull();
  });

  it('should return null for undefined input', () => {
    expect(multiFactorSignatureInfoFromDto(undefined)).toBeNull();
  });
});
