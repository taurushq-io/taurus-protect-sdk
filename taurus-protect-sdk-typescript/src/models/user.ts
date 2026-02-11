/**
 * User and Group models for Taurus-PROTECT SDK.
 */

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
  /** Group creation date */
  readonly createdAt?: Date;
  /** Last modification date */
  readonly updatedAt?: Date;
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
 * Options for listing users.
 */
export interface ListUsersOptions {
  /** Maximum number of users to return */
  limit?: number;
  /** Number of users to skip */
  offset?: number;
  /** Filter by user IDs */
  ids?: string[];
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
 * Options for listing groups.
 */
export interface ListGroupsOptions {
  /** Maximum number of groups to return */
  limit?: number;
  /** Number of groups to skip */
  offset?: number;
  /** Filter by group IDs */
  ids?: string[];
  /** Search query */
  query?: string;
}

/**
 * Options for listing tags.
 */
export interface ListTagsOptions {
  /** Maximum number of tags to return */
  limit?: number;
  /** Number of tags to skip */
  offset?: number;
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
