/**
 * Job mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Job, JobStatistics, JobStatus } from '../models/job';
import { safeDate, safeMap, safeString } from './base';

/**
 * Maps a JobStatus DTO to a JobStatus domain model.
 */
export function jobStatusFromDto(dto: unknown): JobStatus | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    startedAt: safeDate(d.startedAt ?? d.started_at),
    updatedAt: safeDate(d.updatedAt ?? d.updated_at),
    timeoutAt: safeDate(d.timeoutAt ?? d.timeout_at),
    message: safeString(d.message),
    status: safeString(d.status),
  };
}

/**
 * Maps a JobStatistics DTO to a JobStatistics domain model.
 */
export function jobStatisticsFromDto(dto: unknown): JobStatistics | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    pending: safeString(d.pending),
    successes: safeString(d.successes),
    failures: safeString(d.failures),
    lastSuccess: jobStatusFromDto(d.lastSuccess ?? d.last_success),
    lastFailure: jobStatusFromDto(d.lastFailure ?? d.last_failure),
    avgDuration: safeString(d.avgDuration ?? d.avg_duration),
    maxDuration: safeString(d.maxDuration ?? d.max_duration),
    minDuration: safeString(d.minDuration ?? d.min_duration),
  };
}

/**
 * Maps a Job DTO to a Job domain model.
 */
export function jobFromDto(dto: unknown): Job | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    name: safeString(d.name),
    statistics: jobStatisticsFromDto(d.statistics),
  };
}

/**
 * Maps an array of Job DTOs to Job domain models.
 */
export function jobsFromDto(dtos: unknown[] | null | undefined): Job[] {
  return safeMap(dtos, jobFromDto);
}
