/**
 * Change models for Taurus-PROTECT SDK.
 *
 * Changes represent modifications to system configuration that require approval
 * before taking effect. This includes user management, role assignments, and
 * other administrative operations that follow an approval workflow.
 */

/**
 * The status of a change request.
 */
export type ChangeStatus = 'Created' | 'Approved' | 'Rejected' | 'Canceled';

/**
 * Represents a change record in the Taurus-PROTECT audit system.
 *
 * Changes track modifications made to system entities (wallets, addresses, users, etc.)
 * providing a complete audit trail of who changed what, when, and why. Each change
 * captures the before and after state of modified fields.
 */
export interface Change {
  /** Unique identifier for this change record */
  readonly id?: string;
  /** ID of the tenant (organization) where this change occurred */
  readonly tenantId?: number;
  /** Internal user ID of the user who made the change */
  readonly creatorId?: string;
  /** External user ID (from SSO/IdP) of the user who made the change */
  readonly creatorExternalId?: string;
  /** The type of action performed (e.g., "create", "update", "delete") */
  readonly action?: string;
  /** The type of entity that was changed (e.g., "wallet", "address", "user") */
  readonly entity?: string;
  /** The numeric ID of the entity that was changed */
  readonly entityId?: string;
  /** The UUID of the entity that was changed (if applicable) */
  readonly entityUUID?: string;
  /** Map of field names to their new values after the change */
  readonly changes?: Record<string, string>;
  /** Optional comment explaining the reason for the change */
  readonly comment?: string;
  /** Timestamp when the change was made */
  readonly createdAt?: Date;
}

/**
 * Result of a change list query with cursor-based pagination.
 */
export interface ListChangesResult {
  /** The list of changes in this page */
  readonly changes: Change[];
  /** Current page cursor (base64 encoded) */
  readonly currentPage?: string;
  /** Whether there are more pages available */
  readonly hasNext: boolean;
}

/**
 * Options for listing changes.
 */
export interface ListChangesOptions {
  /**
   * The entity type to filter by.
   * Can be one of: user, group, usergroup, businessrule, exchange, price, action,
   * feepayer, userapikey, securitydomain, taurusnetworkparticipant, visibilitygroup,
   * uservisibilitygroup, manualaccountbalancefreeze, wallet, whitelistedaddress,
   * autotransfereventhandler
   */
  entity?: string;
  /** The status to filter by */
  status?: ChangeStatus;
  /** Filter by creator ID */
  creatorId?: string;
  /** Sort order (ASC or DESC, default DESC) */
  sortOrder?: 'ASC' | 'DESC';
  /** Page size (default 50) */
  pageSize?: number;
  /** Current page cursor for pagination */
  currentPage?: string;
  /** Page request type (FIRST, PREVIOUS, NEXT, LAST) */
  pageRequest?: 'FIRST' | 'PREVIOUS' | 'NEXT' | 'LAST';
  /** Filter by entity IDs (valid when entity type is given) */
  entityIDs?: string[];
  /** Filter by entity UUIDs (valid when entity type is given) */
  entityUUIDs?: string[];
}

/**
 * Options for listing changes for approval.
 */
/**
 * Request to create a configuration change.
 */
export interface CreateChangeRequest {
  /** Action type (e.g., "update", "create", "delete") */
  action: string;
  /** Entity type (e.g., "businessrule", "user", "group") */
  entity: string;
  /** ID of the entity to change */
  entityId?: string;
  /** Map of field names to their new values */
  changes?: Record<string, string>;
  /** Optional description of the change request */
  comment?: string;
}

/**
 * Options for listing changes for approval.
 */
export interface ListChangesForApprovalOptions {
  /** Entity types to filter by */
  entities?: string[];
  /** Sort order (ASC or DESC, default DESC) */
  sortOrder?: 'ASC' | 'DESC';
  /** Page size (default 50) */
  pageSize?: number;
  /** Current page cursor for pagination */
  currentPage?: string;
  /** Page request type (FIRST, PREVIOUS, NEXT, LAST) */
  pageRequest?: 'FIRST' | 'PREVIOUS' | 'NEXT' | 'LAST';
  /** Filter by entity IDs (valid when one entity type is given) */
  entityIDs?: string[];
  /** Filter by entity UUIDs (valid when one entity type is given) */
  entityUUIDs?: string[];
}
