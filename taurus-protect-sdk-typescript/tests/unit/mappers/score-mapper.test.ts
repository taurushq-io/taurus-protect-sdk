/**
 * Unit tests for score mapper functions.
 */

import { scoreFromDto, scoresFromDto } from '../../../src/mappers/score';

describe('scoreFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 42,
      provider: 'chainalysis',
      type: 'risk',
      score: '85',
      updateDate: new Date('2024-06-01'),
    };

    const result = scoreFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe(42);
    expect(result!.provider).toBe('chainalysis');
    expect(result!.type).toBe('risk');
    expect(result!.score).toBe('85');
    expect(result!.updateDate).toBeInstanceOf(Date);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      id: 10,
      provider: 'elliptic',
      update_date: '2024-07-01T00:00:00Z',
    };

    const result = scoreFromDto(dto);

    expect(result!.updateDate).toBeInstanceOf(Date);
  });

  it('should return undefined for null input', () => {
    expect(scoreFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(scoreFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = scoreFromDto({});
    expect(result).toBeDefined();
    expect(result!.id).toBeUndefined();
    expect(result!.provider).toBeUndefined();
  });
});

describe('scoresFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 1, provider: 'chainalysis', score: '85' },
      { id: 2, provider: 'elliptic', score: '70' },
    ];

    const result = scoresFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].provider).toBe('chainalysis');
    expect(result[1].provider).toBe('elliptic');
  });

  it('should return empty array for null input', () => {
    expect(scoresFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(scoresFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty array', () => {
    expect(scoresFromDto([])).toEqual([]);
  });
});
