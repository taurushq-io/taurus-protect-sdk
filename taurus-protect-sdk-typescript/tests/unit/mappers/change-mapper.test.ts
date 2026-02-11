/**
 * Unit tests for change mapper functions.
 */

import {
  changeFromDto,
  changesFromDto,
  createChangeRequestToDto,
  listChangesResultFromDto,
} from '../../../src/mappers/change';

describe('changeFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'ch-1',
      tenantId: 42,
      creatorId: 'user-1',
      creatorExternalId: 'ext-user-1',
      action: 'CREATE',
      entity: 'wallet',
      entityId: 'w-1',
      entityUUID: 'uuid-123',
      changes: { name: 'New Name' },
      comment: 'Created wallet',
      creationDate: '2024-01-15T10:00:00Z',
    };

    const result = changeFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('ch-1');
    expect(result!.tenantId).toBe(42);
    expect(result!.creatorId).toBe('user-1');
    expect(result!.creatorExternalId).toBe('ext-user-1');
    expect(result!.action).toBe('CREATE');
    expect(result!.entity).toBe('wallet');
    expect(result!.entityId).toBe('w-1');
    expect(result!.entityUUID).toBe('uuid-123');
    expect(result!.changes).toEqual({ name: 'New Name' });
    expect(result!.comment).toBe('Created wallet');
    expect(result!.createdAt).toBeInstanceOf(Date);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      id: 'ch-2',
      tenant_id: 10,
      creator_id: 'u-2',
      creator_external_id: 'ext-2',
      entity_id: 'e-2',
      entity_uuid: 'uuid-456',
      creation_date: '2024-02-01T00:00:00Z',
    };

    const result = changeFromDto(dto);

    expect(result!.tenantId).toBe(10);
    expect(result!.creatorId).toBe('u-2');
    expect(result!.creatorExternalId).toBe('ext-2');
    expect(result!.entityId).toBe('e-2');
    expect(result!.entityUUID).toBe('uuid-456');
    expect(result!.createdAt).toBeInstanceOf(Date);
  });

  it('should return undefined for null input', () => {
    expect(changeFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(changeFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = changeFromDto({});
    expect(result).toBeDefined();
    expect(result!.id).toBeUndefined();
  });
});

describe('changesFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 'ch-1', action: 'CREATE' },
      { id: 'ch-2', action: 'UPDATE' },
    ];

    const result = changesFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].action).toBe('CREATE');
    expect(result[1].action).toBe('UPDATE');
  });

  it('should return empty array for null input', () => {
    expect(changesFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(changesFromDto(undefined)).toEqual([]);
  });
});

describe('listChangesResultFromDto', () => {
  it('should map response with result and cursor', () => {
    const response = {
      result: [
        { id: 'ch-1', action: 'CREATE' },
        { id: 'ch-2', action: 'UPDATE' },
      ],
      cursor: {
        currentPage: 'page-1',
        hasNext: true,
      },
    };

    const result = listChangesResultFromDto(response);

    expect(result.changes).toHaveLength(2);
    expect(result.currentPage).toBe('page-1');
    expect(result.hasNext).toBe(true);
  });

  it('should handle response with changes key instead of result', () => {
    const response = {
      changes: [{ id: 'ch-1' }],
    };

    const result = listChangesResultFromDto(response);

    expect(result.changes).toHaveLength(1);
  });

  it('should return default for null input', () => {
    const result = listChangesResultFromDto(null);

    expect(result.changes).toEqual([]);
    expect(result.hasNext).toBe(false);
  });

  it('should return default for undefined input', () => {
    const result = listChangesResultFromDto(undefined);

    expect(result.changes).toEqual([]);
    expect(result.hasNext).toBe(false);
  });

  it('should handle snake_case cursor fields', () => {
    const response = {
      result: [],
      cursor: {
        current_page: 'p2',
        has_next: true,
      },
    };

    const result = listChangesResultFromDto(response);

    expect(result.currentPage).toBe('p2');
    expect(result.hasNext).toBe(true);
  });
});

describe('createChangeRequestToDto', () => {
  it('should map all fields to DTO', () => {
    const result = createChangeRequestToDto({
      action: 'update',
      entity: 'businessrule',
      entityId: '42',
      changes: { rulevalue: '100' },
      comment: 'test comment',
    });

    expect(result.action).toBe('update');
    expect(result.entity).toBe('businessrule');
    expect(result.entityId).toBe('42');
    expect(result.changes).toEqual({ rulevalue: '100' });
    expect(result.changeComment).toBe('test comment');
  });

  it('should map comment to changeComment', () => {
    const result = createChangeRequestToDto({
      action: 'create',
      entity: 'user',
      comment: 'my comment',
    });

    expect(result.changeComment).toBe('my comment');
    expect((result as any).comment).toBeUndefined();
  });

  it('should handle missing optional fields', () => {
    const result = createChangeRequestToDto({
      action: 'delete',
      entity: 'group',
    });

    expect(result.action).toBe('delete');
    expect(result.entity).toBe('group');
    expect(result.entityId).toBeUndefined();
    expect(result.changes).toBeUndefined();
    expect(result.changeComment).toBeUndefined();
  });
});
