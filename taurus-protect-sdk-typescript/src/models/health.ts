/**
 * Health models for Taurus-PROTECT SDK.
 */

/**
 * Health check status.
 */
export interface HealthStatus {
  /** Health status (e.g., "healthy", "unhealthy") */
  readonly status: string;
  /** API version */
  readonly version?: string;
  /** Status message */
  readonly message?: string;
}

/**
 * Individual health check result.
 */
export interface HealthCheck {
  /** Health check name */
  readonly name?: string;
  /** Health check group */
  readonly group?: string;
  /** Whether the check passed */
  readonly healthy?: boolean;
  /** Status message */
  readonly message?: string;
  /** Check duration in milliseconds */
  readonly durationMs?: number;
}

/**
 * Component status.
 */
export interface ComponentStatus {
  /** Component name */
  readonly name?: string;
  /** Component status */
  readonly status?: string;
  /** Whether the component is healthy */
  readonly healthy?: boolean;
  /** Last check timestamp */
  readonly lastCheck?: Date;
}
