/**
 * Unit tests for visibility group mapper functions.
 */

import {
  visibilityGroupFromDto,
  visibilityGroupsFromDto,
  visibilityGroupUserFromDto,
  visibilityGroupUsersFromDto,
} from '../../../src/mappers/visibility-group';

describe('visibilityGroupFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'vg-1',
      tenantId: 't-1',
      name: 'Treasury VG',
      description: 'Treasury team visibility',
      userCount: '5',
      users: [
        { id: 'u-1', externalUserId: 'ext-1' },
        { id: 'u-2', externalUserId: 'ext-2' },
      ],
      creationDate: new Date('2024-01-01'),
      updateDate: new Date('2024-06-01'),
    };

    const result = visibilityGroupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('vg-1');
    expect(result!.tenantId).toBe('t-1');
    expect(result!.name).toBe('Treasury VG');
    expect(result!.description).toBe('Treasury team visibility');
    expect(result!.userCount).toBe('5');
    expect(result!.users).toHaveLength(2);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      id: 'vg-2',
      tenant_id: 't-2',
      user_count: '3',
      creation_date: new Date('2024-01-01'),
      update_date: new Date('2024-06-01'),
    };

    const result = visibilityGroupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.tenantId).toBe('t-2');
  });

  it('should return undefined for null input', () => {
    expect(visibilityGroupFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(visibilityGroupFromDto(undefined)).toBeUndefined();
  });

  it('should handle missing users', () => {
    const dto = { id: 'vg-3', name: 'Empty' };
    const result = visibilityGroupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.users).toBeUndefined();
  });
});

describe('visibilityGroupsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: '1', name: 'VG1' },
      { id: '2', name: 'VG2' },
    ];
    const result = visibilityGroupsFromDto(dtos);
    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(visibilityGroupsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(visibilityGroupsFromDto(undefined)).toEqual([]);
  });
});

describe('visibilityGroupUserFromDto', () => {
  it('should map fields', () => {
    const dto = { id: 'u-1', externalUserId: 'ext-1' };
    const result = visibilityGroupUserFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('u-1');
    expect(result!.externalUserId).toBe('ext-1');
  });

  it('should handle snake_case', () => {
    const dto = { id: 'u-2', external_user_id: 'ext-2' };
    const result = visibilityGroupUserFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.externalUserId).toBe('ext-2');
  });

  it('should return undefined for null input', () => {
    expect(visibilityGroupUserFromDto(null)).toBeUndefined();
  });
});

describe('visibilityGroupUsersFromDto', () => {
  it('should return empty array for null input', () => {
    expect(visibilityGroupUsersFromDto(null)).toEqual([]);
  });
});
