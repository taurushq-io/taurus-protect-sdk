/**
 * Unit tests for group mapper functions.
 *
 * Tests groupFromDto and groupsFromDto from the user mapper module.
 */

import { groupFromDto, groupsFromDto } from '../../../src/mappers/user';

describe('groupFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'grp-1',
      externalGroupId: 'ext-g-1',
      name: 'Admins',
      description: 'Admin group',
      userIds: ['u-1', 'u-2', 'u-3'],
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-06-15T12:00:00Z',
    };

    const result = groupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('grp-1');
    expect(result!.externalGroupId).toBe('ext-g-1');
    expect(result!.name).toBe('Admins');
    expect(result!.description).toBe('Admin group');
    expect(result!.userIds).toEqual(['u-1', 'u-2', 'u-3']);
    expect(result!.createdAt).toBeDefined();
    expect(result!.updatedAt).toBeDefined();
  });

  it('should handle snake_case field names', () => {
    const dto = {
      group_id: 'grp-2',
      external_group_id: 'ext-g-2',
      name: 'Operators',
      user_ids: ['u-10', 'u-20'],
      created_at: '2024-02-01T00:00:00Z',
      updated_at: '2024-07-01T00:00:00Z',
    };

    const result = groupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('grp-2');
    expect(result!.externalGroupId).toBe('ext-g-2');
    expect(result!.userIds).toEqual(['u-10', 'u-20']);
  });

  it('should handle alternative id fields', () => {
    const dto = {
      groupId: 'grp-3',
      name: 'Viewers',
    };

    const result = groupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('grp-3');
  });

  it('should handle users field as user IDs', () => {
    const dto = {
      id: 'grp-4',
      users: ['u-100', 'u-200'],
    };

    const result = groupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.userIds).toEqual(['u-100', 'u-200']);
  });

  it('should return undefined for null input', () => {
    expect(groupFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(groupFromDto(undefined)).toBeUndefined();
  });

  it('should return undefined for non-object input', () => {
    expect(groupFromDto('not-an-object')).toBeUndefined();
    expect(groupFromDto(42)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = groupFromDto({});
    expect(result).toBeDefined();
    expect(result!.id).toBeUndefined();
    expect(result!.name).toBeUndefined();
    expect(result!.userIds).toBeUndefined();
  });

  it('should handle creationDate fallback', () => {
    const dto = {
      id: 'grp-5',
      creationDate: '2024-03-01T00:00:00Z',
      modificationDate: '2024-08-01T00:00:00Z',
    };

    const result = groupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.createdAt).toBeDefined();
    expect(result!.updatedAt).toBeDefined();
  });
});

describe('groupsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 'grp-1', name: 'Group A' },
      { id: 'grp-2', name: 'Group B' },
      { id: 'grp-3', name: 'Group C' },
    ];
    const result = groupsFromDto(dtos);
    expect(result).toHaveLength(3);
    expect(result[0].name).toBe('Group A');
    expect(result[2].name).toBe('Group C');
  });

  it('should return empty array for null input', () => {
    expect(groupsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(groupsFromDto(undefined)).toEqual([]);
  });

  it('should filter out invalid entries', () => {
    const dtos = [
      { id: 'grp-1', name: 'Valid' },
      null,
      undefined,
      { id: 'grp-2', name: 'Also Valid' },
    ];
    const result = groupsFromDto(dtos as unknown[]);
    expect(result).toHaveLength(2);
  });
});
