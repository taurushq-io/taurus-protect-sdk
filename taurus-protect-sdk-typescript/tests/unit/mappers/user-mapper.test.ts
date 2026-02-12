/**
 * Unit tests for user, group, and tag mapper functions.
 */

import {
  userFromDto,
  usersFromDto,
  groupFromDto,
  groupsFromDto,
  tagFromDto,
  tagsFromDto,
  userAttributeFromDto,
} from '../../../src/mappers/user';
import { UserStatus } from '../../../src/models/user';

describe('userFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'u-1',
      externalUserId: 'ext-1',
      email: 'jdoe@example.com',
      firstName: 'John',
      lastName: 'Doe',
      status: 'ACTIVE',
      roles: ['ADMIN', 'OPERATOR'],
      totpEnabled: true,
      publicKey: '-----BEGIN PUBLIC KEY-----',
      createdAt: new Date('2024-01-01'),
      updatedAt: new Date('2024-06-01'),
      attributes: [
        { id: 'a1', key: 'dept', value: 'IT' },
      ],
      groupIds: ['g-1', 'g-2'],
    };

    const result = userFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('u-1');
    expect(result!.email).toBe('jdoe@example.com');
    expect(result!.firstName).toBe('John');
    expect(result!.lastName).toBe('Doe');
    expect(result!.status).toBe(UserStatus.ACTIVE);
    expect(result!.roles).toEqual(['ADMIN', 'OPERATOR']);
    expect(result!.totpEnabled).toBe(true);
    expect(result!.attributes).toHaveLength(1);
    expect(result!.groupIds).toEqual(['g-1', 'g-2']);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      user_id: 'u-2',
      external_user_id: 'ext-2',
      first_name: 'Jane',
      last_name: 'Smith',
      totp_enabled: false,
      public_key: 'PEM',
      created_at: new Date('2024-01-01'),
      group_ids: ['g-1'],
    };

    const result = userFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('u-2');
    expect(result!.firstName).toBe('Jane');
    expect(result!.groupIds).toEqual(['g-1']);
  });

  it('should map status enum values', () => {
    expect(userFromDto({ id: '1', status: 'ACTIVE' })!.status).toBe(UserStatus.ACTIVE);
    expect(userFromDto({ id: '1', status: 'INACTIVE' })!.status).toBe(UserStatus.INACTIVE);
    expect(userFromDto({ id: '1', status: 'PENDING' })!.status).toBe(UserStatus.PENDING);
    expect(userFromDto({ id: '1', status: 'LOCKED' })!.status).toBe(UserStatus.LOCKED);
    expect(userFromDto({ id: '1', status: 'UNKNOWN' })!.status).toBeUndefined();
  });

  it('should return undefined for null input', () => {
    expect(userFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(userFromDto(undefined)).toBeUndefined();
  });
});

describe('usersFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: '1', email: 'a@example.com' },
      { id: '2', email: 'b@example.com' },
    ];
    const result = usersFromDto(dtos);
    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(usersFromDto(null)).toEqual([]);
  });
});

describe('groupFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'g-1',
      externalGroupId: 'ext-g-1',
      name: 'Operators',
      description: 'Operations team',
      userIds: ['u-1', 'u-2'],
      createdAt: new Date('2024-01-01'),
      updatedAt: new Date('2024-06-01'),
    };

    const result = groupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('g-1');
    expect(result!.name).toBe('Operators');
    expect(result!.description).toBe('Operations team');
    expect(result!.userIds).toEqual(['u-1', 'u-2']);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      group_id: 'g-2',
      external_group_id: 'ext-g-2',
      user_ids: ['u-3'],
    };
    const result = groupFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('g-2');
    expect(result!.userIds).toEqual(['u-3']);
  });

  it('should return undefined for null input', () => {
    expect(groupFromDto(null)).toBeUndefined();
  });
});

describe('groupsFromDto', () => {
  it('should return empty array for null input', () => {
    expect(groupsFromDto(null)).toEqual([]);
  });
});

describe('tagFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'tag-1',
      name: 'important',
      color: '#FF0000',
      createdAt: new Date('2024-01-01'),
    };

    const result = tagFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('tag-1');
    expect(result!.name).toBe('important');
    expect(result!.color).toBe('#FF0000');
  });

  it('should use value field for name fallback', () => {
    const dto = {
      id: 'tag-2',
      value: 'urgent',
    };
    const result = tagFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.name).toBe('urgent');
  });

  it('should return undefined for null input', () => {
    expect(tagFromDto(null)).toBeUndefined();
  });
});

describe('tagsFromDto', () => {
  it('should return empty array for null input', () => {
    expect(tagsFromDto(null)).toEqual([]);
  });
});

describe('userAttributeFromDto', () => {
  it('should map fields', () => {
    const dto = { id: 'a1', key: 'dept', value: 'IT' };
    const result = userAttributeFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('a1');
    expect(result!.key).toBe('dept');
    expect(result!.value).toBe('IT');
  });

  it('should return undefined for null input', () => {
    expect(userAttributeFromDto(null)).toBeUndefined();
  });
});
