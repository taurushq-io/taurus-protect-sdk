/**
 * Audit models for Taurus-PROTECT SDK.
 */

/**
 * Audit trail entry.
 */
export interface AuditTrail {
  /** Unique audit trail identifier */
  readonly id?: string;
  /** Entity type (e.g., "WALLET", "REQUEST", "ADDRESS") */
  readonly entity?: string;
  /** Entity ID */
  readonly entityId?: string;
  /** Action performed (e.g., "CREATE", "UPDATE", "DELETE") */
  readonly action?: string;
  /** User ID who performed the action */
  readonly userId?: string;
  /** User email who performed the action */
  readonly userEmail?: string;
  /** External user ID */
  readonly externalUserId?: string;
  /** Description of the action */
  readonly description?: string;
  /** IP address of the requester */
  readonly ipAddress?: string;
  /** Creation date of the audit trail */
  readonly createdAt?: Date;
  /** Additional details as JSON */
  readonly details?: Record<string, unknown>;
}

/**
 * Options for listing audit trails.
 */
export interface ListAuditTrailsOptions {
  /** Maximum number of audit trails to return */
  limit?: number;
  /** Filter by external user ID */
  externalUserId?: string;
  /** Filter by entity types */
  entities?: string[];
  /** Filter by actions */
  actions?: string[];
  /** Filter by creation date from */
  creationDateFrom?: Date;
  /** Filter by creation date to */
  creationDateTo?: Date;
  /** Sort order (ASC or DESC) */
  sortOrder?: 'ASC' | 'DESC';
}
