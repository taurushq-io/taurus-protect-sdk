/**
 * Unit tests for job mapper functions.
 */

import {
  jobFromDto,
  jobsFromDto,
  jobStatusFromDto,
  jobStatisticsFromDto,
} from '../../../src/mappers/job';

describe('jobStatusFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      id: 'js-1',
      startedAt: '2024-01-01T10:00:00Z',
      updatedAt: '2024-01-01T10:05:00Z',
      timeoutAt: '2024-01-01T11:00:00Z',
      message: 'Processing complete',
      status: 'SUCCESS',
    };

    const result = jobStatusFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('js-1');
    expect(result!.startedAt).toBeInstanceOf(Date);
    expect(result!.updatedAt).toBeInstanceOf(Date);
    expect(result!.timeoutAt).toBeInstanceOf(Date);
    expect(result!.message).toBe('Processing complete');
    expect(result!.status).toBe('SUCCESS');
  });

  it('should handle snake_case field names', () => {
    const dto = {
      id: 'js-2',
      started_at: '2024-02-01T00:00:00Z',
      updated_at: '2024-02-01T01:00:00Z',
      timeout_at: '2024-02-01T02:00:00Z',
    };

    const result = jobStatusFromDto(dto);

    expect(result!.startedAt).toBeInstanceOf(Date);
    expect(result!.updatedAt).toBeInstanceOf(Date);
    expect(result!.timeoutAt).toBeInstanceOf(Date);
  });

  it('should return undefined for null input', () => {
    expect(jobStatusFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(jobStatusFromDto(undefined)).toBeUndefined();
  });
});

describe('jobStatisticsFromDto', () => {
  it('should map all fields', () => {
    const dto = {
      pending: '5',
      successes: '100',
      failures: '2',
      lastSuccess: { id: 'ls-1', status: 'SUCCESS' },
      lastFailure: { id: 'lf-1', status: 'FAILURE' },
      avgDuration: '1500',
      maxDuration: '5000',
      minDuration: '200',
    };

    const result = jobStatisticsFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.pending).toBe('5');
    expect(result!.successes).toBe('100');
    expect(result!.failures).toBe('2');
    expect(result!.lastSuccess).toBeDefined();
    expect(result!.lastSuccess!.id).toBe('ls-1');
    expect(result!.lastFailure).toBeDefined();
    expect(result!.lastFailure!.id).toBe('lf-1');
    expect(result!.avgDuration).toBe('1500');
    expect(result!.maxDuration).toBe('5000');
    expect(result!.minDuration).toBe('200');
  });

  it('should handle snake_case field names', () => {
    const dto = {
      pending: '0',
      last_success: { id: 'ls-2', status: 'SUCCESS' },
      last_failure: null,
      avg_duration: '1000',
      max_duration: '3000',
      min_duration: '100',
    };

    const result = jobStatisticsFromDto(dto);

    expect(result!.lastSuccess).toBeDefined();
    expect(result!.avgDuration).toBe('1000');
  });

  it('should return undefined for null input', () => {
    expect(jobStatisticsFromDto(null)).toBeUndefined();
  });
});

describe('jobFromDto', () => {
  it('should map job with name and statistics', () => {
    const dto = {
      name: 'broadcast-job',
      statistics: {
        pending: '2',
        successes: '50',
        failures: '1',
        avgDuration: '800',
      },
    };

    const result = jobFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.name).toBe('broadcast-job');
    expect(result!.statistics).toBeDefined();
    expect(result!.statistics!.pending).toBe('2');
  });

  it('should handle missing statistics', () => {
    const dto = { name: 'sync-job' };

    const result = jobFromDto(dto);

    expect(result!.name).toBe('sync-job');
    expect(result!.statistics).toBeUndefined();
  });

  it('should return undefined for null input', () => {
    expect(jobFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(jobFromDto(undefined)).toBeUndefined();
  });
});

describe('jobsFromDto', () => {
  it('should map array of job DTOs', () => {
    const dtos = [
      { name: 'job-1', statistics: { pending: '0' } },
      { name: 'job-2', statistics: { pending: '3' } },
    ];

    const result = jobsFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].name).toBe('job-1');
    expect(result[1].name).toBe('job-2');
  });

  it('should return empty array for null input', () => {
    expect(jobsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(jobsFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty array', () => {
    expect(jobsFromDto([])).toEqual([]);
  });
});
