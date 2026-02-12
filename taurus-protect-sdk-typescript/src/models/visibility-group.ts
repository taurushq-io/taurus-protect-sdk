/**
 * Visibility group models for Taurus-PROTECT SDK.
 *
 * Visibility groups are used to control which users can see specific
 * wallets, addresses, and other entities. Users can only see entities
 * that belong to their visibility groups.
 */

/**
 * A user assigned to a visibility group.
 */
export interface VisibilityGroupUser {
  /** Unique user identifier */
  readonly id?: string;
  /** External user identifier */
  readonly externalUserId?: string;
}

/**
 * Represents a visibility group in the Taurus-PROTECT system.
 *
 * Visibility groups are used to control data access. Users can only see
 * wallets, addresses, and other entities that belong to their assigned
 * visibility groups. This provides fine-grained access control within
 * a tenant.
 */
export interface VisibilityGroup {
  /** Unique visibility group identifier */
  readonly id?: string;
  /** Tenant identifier */
  readonly tenantId?: string;
  /** Visibility group name */
  readonly name?: string;
  /** Visibility group description */
  readonly description?: string;
  /** Users assigned to this visibility group */
  readonly users?: VisibilityGroupUser[];
  /** Creation date */
  readonly creationDate?: Date;
  /** Last update date */
  readonly updateDate?: Date;
  /** Number of users in the group (as string) */
  readonly userCount?: string;
}
