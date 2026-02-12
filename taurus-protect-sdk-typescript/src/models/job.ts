/**
 * Job models for Taurus-PROTECT SDK.
 *
 * Jobs are background tasks that process various operations such as
 * transaction monitoring, balance updates, and other async operations.
 */

/**
 * Represents the status of a job execution.
 *
 * Contains timing information and the current status of a job run.
 */
export interface JobStatus {
  /** Unique job execution identifier */
  readonly id?: string;
  /** When the job started */
  readonly startedAt?: Date;
  /** When the job was last updated */
  readonly updatedAt?: Date;
  /** When the job will timeout */
  readonly timeoutAt?: Date;
  /** Status message or error details */
  readonly message?: string;
  /** Current status (e.g., "running", "completed", "failed") */
  readonly status?: string;
}

/**
 * Statistics for a job in the Taurus-PROTECT system.
 *
 * Contains metrics about job execution including success/failure counts
 * and duration statistics.
 */
export interface JobStatistics {
  /** Number of pending job executions */
  readonly pending?: string;
  /** Number of successful job executions */
  readonly successes?: string;
  /** Number of failed job executions */
  readonly failures?: string;
  /** Status of the last successful execution */
  readonly lastSuccess?: JobStatus;
  /** Status of the last failed execution */
  readonly lastFailure?: JobStatus;
  /** Average duration of job executions */
  readonly avgDuration?: string;
  /** Maximum duration of job executions */
  readonly maxDuration?: string;
  /** Minimum duration of job executions */
  readonly minDuration?: string;
}

/**
 * Represents a job in the Taurus-PROTECT system.
 *
 * Jobs are background tasks that process various operations such as
 * transaction monitoring, balance updates, and other async operations.
 */
export interface Job {
  /** Unique job name identifier */
  readonly name?: string;
  /** Job execution statistics */
  readonly statistics?: JobStatistics;
}
