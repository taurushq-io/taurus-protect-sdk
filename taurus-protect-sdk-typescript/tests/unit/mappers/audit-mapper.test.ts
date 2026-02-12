/**
 * Unit tests for audit mapper functions.
 */

import { auditTrailFromDto, auditTrailsFromDto } from '../../../src/mappers/audit';

describe('auditTrailFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'aud-1',
      entity: 'WALLET',
      entityId: 'w-1',
      action: 'CREATE',
      userId: 'u-1',
      userEmail: 'admin@example.com',
      externalUserId: 'ext-u-1',
      description: 'Wallet created',
      ipAddress: '192.168.1.1',
      createdAt: new Date('2024-01-15'),
    };

    const result = auditTrailFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('aud-1');
    expect(result!.entity).toBe('WALLET');
    expect(result!.entityId).toBe('w-1');
    expect(result!.action).toBe('CREATE');
    expect(result!.userId).toBe('u-1');
    expect(result!.userEmail).toBe('admin@example.com');
    expect(result!.description).toBe('Wallet created');
    expect(result!.ipAddress).toBe('192.168.1.1');
  });

  it('should handle snake_case field names', () => {
    const dto = {
      audit_id: 'aud-2',
      entity_type: 'ADDRESS',
      entity_id: 'a-1',
      action_type: 'UPDATE',
      user_id: 'u-2',
      user_email: 'user@example.com',
      external_user_id: 'ext-2',
      ip_address: '10.0.0.1',
      created_at: new Date('2024-06-01'),
    };

    const result = auditTrailFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('aud-2');
    expect(result!.entity).toBe('ADDRESS');
    expect(result!.userId).toBe('u-2');
    expect(result!.ipAddress).toBe('10.0.0.1');
  });

  it('should handle details as object', () => {
    const dto = {
      id: 'aud-3',
      details: { oldValue: 'A', newValue: 'B' },
    };
    const result = auditTrailFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.details).toEqual({ oldValue: 'A', newValue: 'B' });
  });

  it('should handle details as JSON string', () => {
    const dto = {
      id: 'aud-4',
      details: '{"oldValue":"A","newValue":"B"}',
    };
    const result = auditTrailFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.details).toEqual({ oldValue: 'A', newValue: 'B' });
  });

  it('should handle invalid JSON details string gracefully', () => {
    const dto = {
      id: 'aud-5',
      details: 'not-json',
    };
    const result = auditTrailFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.details).toBeUndefined();
  });

  it('should return undefined for null input', () => {
    expect(auditTrailFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(auditTrailFromDto(undefined)).toBeUndefined();
  });
});

describe('auditTrailsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: '1', entity: 'WALLET' },
      { id: '2', entity: 'ADDRESS' },
    ];
    const result = auditTrailsFromDto(dtos);
    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(auditTrailsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(auditTrailsFromDto(undefined)).toEqual([]);
  });
});
