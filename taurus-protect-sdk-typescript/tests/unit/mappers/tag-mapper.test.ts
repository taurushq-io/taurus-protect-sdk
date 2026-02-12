/**
 * Unit tests for tag mapper functions.
 *
 * Tests tagFromDto and tagsFromDto from the user mapper module.
 */

import { tagFromDto, tagsFromDto } from '../../../src/mappers/user';

describe('tagFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'tag-1',
      name: 'Important',
      color: '#FF0000',
      createdAt: '2024-01-15T10:30:00Z',
    };

    const result = tagFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('tag-1');
    expect(result!.name).toBe('Important');
    expect(result!.color).toBe('#FF0000');
    expect(result!.createdAt).toBeDefined();
  });

  it('should handle value field as name fallback', () => {
    const dto = {
      id: 'tag-2',
      value: 'Urgent',
      color: '#FFAA00',
    };

    const result = tagFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.name).toBe('Urgent');
  });

  it('should handle snake_case field names', () => {
    const dto = {
      tag_id: 'tag-3',
      name: 'Low Priority',
      color: '#00FF00',
      created_at: '2024-03-01T00:00:00Z',
    };

    const result = tagFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('tag-3');
    expect(result!.createdAt).toBeDefined();
  });

  it('should handle tagId field', () => {
    const dto = {
      tagId: 'tag-4',
      name: 'Special',
      color: '#0000FF',
    };

    const result = tagFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('tag-4');
  });

  it('should handle creationDate fallback', () => {
    const dto = {
      id: 'tag-5',
      name: 'Archive',
      color: '#808080',
      creationDate: '2024-05-01T00:00:00Z',
    };

    const result = tagFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.createdAt).toBeDefined();
  });

  it('should return undefined for null input', () => {
    expect(tagFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(tagFromDto(undefined)).toBeUndefined();
  });

  it('should return undefined for non-object input', () => {
    expect(tagFromDto('not-an-object')).toBeUndefined();
    expect(tagFromDto(123)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = tagFromDto({});
    expect(result).toBeDefined();
    expect(result!.id).toBeUndefined();
    expect(result!.name).toBeUndefined();
    expect(result!.color).toBeUndefined();
    expect(result!.createdAt).toBeUndefined();
  });
});

describe('tagsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 'tag-1', name: 'Tag A', color: '#FF0000' },
      { id: 'tag-2', name: 'Tag B', color: '#00FF00' },
      { id: 'tag-3', name: 'Tag C', color: '#0000FF' },
    ];
    const result = tagsFromDto(dtos);
    expect(result).toHaveLength(3);
    expect(result[0].name).toBe('Tag A');
    expect(result[1].color).toBe('#00FF00');
    expect(result[2].id).toBe('tag-3');
  });

  it('should return empty array for null input', () => {
    expect(tagsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(tagsFromDto(undefined)).toEqual([]);
  });

  it('should filter out invalid entries', () => {
    const dtos = [
      { id: 'tag-1', name: 'Valid' },
      null,
      undefined,
      { id: 'tag-2', name: 'Also Valid' },
    ];
    const result = tagsFromDto(dtos as unknown[]);
    expect(result).toHaveLength(2);
  });
});
