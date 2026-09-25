/**
 * User and Group models for Taurus-PROTECT SDK.
 */

import type { OffsetPageOptions } from './pagination';

/**
 * User status enum.
 */
export enum UserStatus {
  ACTIVE = 'ACTIVE',
  INACTIVE = 'INACTIVE',
  PENDING = 'PENDING',
  LOCKED = 'LOCKED',
}

/**
 * User role enum.
 */
export enum UserRole {
  ADMIN = 'ADMIN',
  ADMIN_READ_ONLY = 'ADMIN_READ_ONLY',
  USER = 'USER',
  TECHNICAL = 'TECHNICAL',
  SUPER_ADMIN = 'SUPER_ADMIN',
  HSM_SLOT = 'HSMSLOT',
}

/**
 * User information.
 */
export interface User {
  /** Unique user identifier */
  readonly id?: string;
  /** External user identifier */
  readonly externalUserId?: string;
  /** User email address */
  readonly email?: string;
  /** First name */
  readonly firstName?: string;
  /** Last name */
  readonly lastName?: string;
  /** User status */
  readonly status?: UserStatus;
  /** User roles */
  readonly roles?: string[];
  /** Whether TOTP is enabled */
  readonly totpEnabled?: boolean;
  /** User's public key (base64 encoded) */
  readonly publicKey?: string;
  /** User creation date */
  readonly createdAt?: Date;
  /** Last modification date */
  readonly updatedAt?: Date;
  /** User attributes */
  readonly attributes?: UserAttribute[];
  /** Group IDs the user belongs to */
  readonly groupIds?: string[];
  /** Groups the user belongs to */
  readonly groups?: UserGroup[];
  /**
   * Whether the user is in the enforced governance rules. `undefined` where the endpoint
   * does not compute it (`users.get`, `visibilityGroups.getUsersByVisibilityGroup`).
   */
  readonly enforcedInRules?: boolean;
  /**
   * Whether the user's public key is the one in the enforced governance rules. Only
   * `users.getCurrentUser` computes it; `undefined` elsewhere.
   */
  readonly publicKeyEnforcedInRules?: boolean;
}

/**
 * A group the user belongs to.
 */
export interface UserGroup {
  /** Group identifier */
  readonly id?: string;
  /** External group identifier */
  readonly externalGroupId?: string;
  /** Whether the group is in the enforced governance rules; `undefined` where not computed */
  readonly enforcedInRules?: boolean;
}

/**
 * User attribute (key-value pair).
 */
export interface UserAttribute {
  /** Attribute ID */
  readonly id?: string;
  /** Attribute key */
  readonly key?: string;
  /** Attribute value */
  readonly value?: string;
}

/**
 * Group information.
 */
export interface Group {
  /** Unique group identifier */
  readonly id?: string;
  /** External group identifier */
  readonly externalGroupId?: string;
  /** Group name */
  readonly name?: string;
  /** Group description */
  readonly description?: string;
  /** User IDs in the group */
  readonly userIds?: string[];
  /** Users in the group */
  readonly users?: GroupUser[];
  /** Whether the group is in the enforced governance rules */
  readonly enforcedInRules?: boolean;
  /** Group creation date */
  readonly createdAt?: Date;
  /** Last modification date */
  readonly updatedAt?: Date;
}

/**
 * A user within a group.
 */
export interface GroupUser {
  /** User identifier */
  readonly id?: string;
  /** External user identifier */
  readonly externalUserId?: string;
  /** Whether the user is in the enforced governance rules */
  readonly enforcedInRules?: boolean;
}

/**
 * Tag information.
 */
export interface Tag {
  /** Unique tag identifier */
  readonly id?: string;
  /** Tag name/value */
  readonly name?: string;
  /** Tag color (hex code) */
  readonly color?: string;
  /** Tag creation date */
  readonly createdAt?: Date;
}

/**
 * Options for listing users. An offset list: `limit` 1-100 (default 20) and `offset`
 * (a previous page's `pagination.nextOffset`).
 */
export interface ListUsersOptions extends OffsetPageOptions {
  /** Filter by user IDs */
  ids?: string[];
  /** Filter by external user IDs */
  externalUserIds?: string[];
  /** Filter by emails */
  emails?: string[];
  /** Search query */
  query?: string;
  /** Exclude technical users */
  excludeTechnicalUsers?: boolean;
  /** Filter by roles */
  roles?: string[];
  /** Filter by status */
  status?: string;
  /** Filter by TOTP enabled */
  totpEnabled?: boolean;
  /** Filter by group IDs */
  groupIds?: string[];
}

/**
 * Options for listing groups. An offset list: `limit` 1-100 (default 20) and `offset`
 * (a previous page's `pagination.nextOffset`).
 */
export interface ListGroupsOptions extends OffsetPageOptions {
  /** Filter by group IDs */
  ids?: string[];
  /** Search query */
  query?: string;
}

/**
 * Options for listing tags. The tags endpoint is not paged, so `list` returns every tag
 * that matches.
 */
export interface ListTagsOptions {
  /** Filter by tag IDs */
  ids?: string[];
  /** Search query */
  query?: string;
}

/**
 * Request for creating a tag.
 */
export interface CreateTagRequest {
  /** Tag name/value */
  name: string;
  /** Tag color (hex code) */
  color: string;
}
