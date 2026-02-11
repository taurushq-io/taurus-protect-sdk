/**
 * Job service for Taurus-PROTECT SDK.
 *
 * Provides methods for monitoring jobs in the Taurus-PROTECT system.
 * Jobs are background tasks that process various operations such as
 * transaction monitoring, balance updates, and other async operations.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { JobsApi } from '../internal/openapi/apis/JobsApi';
import { jobFromDto, jobsFromDto, jobStatusFromDto } from '../mappers/job';
import type { Job, JobStatus } from '../models/job';
import { BaseService } from './base';

/**
 * Service for monitoring jobs in the Taurus-PROTECT system.
 *
 * Jobs are background tasks that process various operations such as
 * transaction monitoring, balance updates, and other async operations.
 *
 * @example
 * ```typescript
 * // Get all jobs
 * const jobs = await jobService.list();
 * for (const job of jobs) {
 *   console.log(`${job.name}: ${job.statistics?.successes} successes`);
 * }
 *
 * // Get a specific job
 * const job = await jobService.get('balance-sync');
 * console.log(`Job ${job.name} has ${job.statistics?.pending} pending`);
 *
 * // Get job execution status
 * const status = await jobService.getStatus('balance-sync', 'exec-123');
 * console.log(`Status: ${status.status}, Message: ${status.message}`);
 * ```
 */
export class JobService extends BaseService {
  private readonly jobsApi: JobsApi;

  /**
   * Creates a new JobService instance.
   *
   * @param jobsApi - The JobsApi instance from the OpenAPI client
   */
  constructor(jobsApi: JobsApi) {
    super();
    this.jobsApi = jobsApi;
  }

  /**
   * Lists all jobs.
   *
   * @returns Array of jobs
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const jobs = await jobService.list();
   * for (const job of jobs) {
   *   console.log(`${job.name}: ${job.statistics?.successes} successes`);
   * }
   * ```
   */
  async list(): Promise<Job[]> {
    return this.execute(async () => {
      const response = await this.jobsApi.jobServiceGetJobs();

      const result =
        (response as Record<string, unknown>).jobs ??
        (response as Record<string, unknown>).result;
      return jobsFromDto(result as unknown[]);
    });
  }

  /**
   * Gets a job by name.
   *
   * @param name - The job name
   * @returns The job
   * @throws {@link ValidationError} If name is empty
   * @throws {@link NotFoundError} If job is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const job = await jobService.get('balance-sync');
   * console.log(`Job ${job.name} has ${job.statistics?.pending} pending`);
   * ```
   */
  async get(name: string): Promise<Job> {
    if (!name || name.trim() === '') {
      throw new ValidationError('name is required');
    }

    return this.execute(async () => {
      const response = await this.jobsApi.jobServiceGetJob({ name });

      const result =
        (response as Record<string, unknown>).job ??
        (response as Record<string, unknown>).result;
      const job = jobFromDto(result);

      if (!job) {
        throw new NotFoundError(`Job with name '${name}' not found`);
      }

      return job;
    });
  }

  /**
   * Gets the status of a specific job execution.
   *
   * @param name - The job name
   * @param id - The job execution ID
   * @returns The job status
   * @throws {@link ValidationError} If name or id is empty
   * @throws {@link NotFoundError} If job status is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const status = await jobService.getStatus('balance-sync', 'exec-123');
   * console.log(`Status: ${status.status}`);
   * console.log(`Started: ${status.startedAt}`);
   * console.log(`Message: ${status.message}`);
   * ```
   */
  async getStatus(name: string, id: string): Promise<JobStatus> {
    if (!name || name.trim() === '') {
      throw new ValidationError('name is required');
    }
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response = await this.jobsApi.jobServiceGetJobStatus({ name, id });

      const result =
        (response as Record<string, unknown>).status ??
        (response as Record<string, unknown>).result;
      const status = jobStatusFromDto(result);

      if (!status) {
        throw new NotFoundError(`Job status for '${name}' execution '${id}' not found`);
      }

      return status;
    });
  }
}
