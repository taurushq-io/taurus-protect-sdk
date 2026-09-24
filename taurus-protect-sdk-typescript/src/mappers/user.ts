/**
 * User, Group, and Tag mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Group, GroupUser, Tag, User, UserAttribute, UserGroup } from '../models/user';
import { UserStatus } from '../models/user';
import { safeBool, safeDate, safeMap, safeString } from './base';

/**
 * Maps a list of membership objects. A list of bare IDs (or anything else) maps to
 * undefined rather than to an empty membership list.
 */
function membershipsFromDto<T>(value: unknown, mapper: (dto: unknown) => T | undefined): T[] | undefined {
  if (!Array.isArray(value) || !value.every((item) => item !== null && typeof item === 'object')) {
    return undefined;
  }
  return safeMap(value, mapper);
}

/**
 * Reads IDs from a list of bare IDs or, as validatord sends memberships, of objects with an `id`.
 */
function idsFromDto(value: unknown): string[] | undefined {
  if (!Array.isArray(value)) {
    return undefined;
  }
  const ids: string[] = [];
  for (const item of value) {
    const id =
      item !== null && typeof item === 'object' ? safeString((item as Record<string, unknown>).id) : safeString(item);
    if (id !== undefined) {
      ids.push(id);
    }
  }
  return ids;
}

function userGroupFromDto(dto: unknown): UserGroup | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    externalGroupId: safeString(d.externalGroupId ?? d.external_group_id),
    enforcedInRules: safeBool(d.enforcedInRules ?? d.enforced_in_rules),
  };
}

function groupUserFromDto(dto: unknown): GroupUser | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    externalUserId: safeString(d.externalUserId ?? d.external_user_id),
    enforcedInRules: safeBool(d.enforcedInRules ?? d.enforced_in_rules),
  };
}

/**
 * Maps a user attribute DTO to a UserAttribute domain model.
 */
export function userAttributeFromDto(dto: unknown): UserAttribute | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    key: safeString(d.key),
    value: safeString(d.value),
  };
}

/**
 * Maps a user status string to UserStatus enum.
 */
function mapUserStatus(status: string | undefined): UserStatus | undefined {
  if (!status) {
    return undefined;
  }
  const upper = status.toUpperCase();
  switch (upper) {
    case 'ACTIVE':
      return UserStatus.ACTIVE;
    case 'INACTIVE':
      return UserStatus.INACTIVE;
    case 'PENDING':
      return UserStatus.PENDING;
    case 'LOCKED':
      return UserStatus.LOCKED;
    default:
      return undefined;
  }
}

/**
 * Maps a user DTO to a User domain model.
 */
export function userFromDto(dto: unknown): User | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const attributes = d.attributes as unknown[] | undefined;
  const roles = d.roles as string[] | undefined;
  const groupIds = d.groupIds ?? d.group_ids ?? d.groups;

  return {
    id: safeString(d.id ?? d.userId ?? d.user_id),
    externalUserId: safeString(d.externalUserId ?? d.external_user_id),
    email: safeString(d.email),
    firstName: safeString(d.firstName ?? d.first_name ?? d.firstname),
    lastName: safeString(d.lastName ?? d.last_name ?? d.lastname),
    status: mapUserStatus(safeString(d.status)),
    roles: roles,
    totpEnabled: safeBool(d.totpEnabled ?? d.totp_enabled),
    publicKey: safeString(d.publicKey ?? d.public_key),
    createdAt: safeDate(d.createdAt ?? d.created_at ?? d.creationDate),
    updatedAt: safeDate(d.updatedAt ?? d.updated_at ?? d.modificationDate),
    attributes: attributes ? safeMap(attributes, userAttributeFromDto) : undefined,
    groupIds: idsFromDto(groupIds),
    groups: membershipsFromDto(d.groups, userGroupFromDto),
    enforcedInRules: safeBool(d.enforcedInRules ?? d.enforced_in_rules),
    publicKeyEnforcedInRules: safeBool(d.publicKeyEnforcedInRules ?? d.public_key_enforced_in_rules),
  };
}

/**
 * Maps an array of user DTOs to User domain models.
 */
export function usersFromDto(dtos: unknown[] | null | undefined): User[] {
  return safeMap(dtos, userFromDto);
}

/**
 * Maps a group DTO to a Group domain model.
 */
export function groupFromDto(dto: unknown): Group | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const userIds = d.userIds ?? d.user_ids ?? d.users;

  return {
    id: safeString(d.id ?? d.groupId ?? d.group_id),
    externalGroupId: safeString(d.externalGroupId ?? d.external_group_id),
    name: safeString(d.name),
    description: safeString(d.description),
    userIds: idsFromDto(userIds),
    users: membershipsFromDto(d.users, groupUserFromDto),
    enforcedInRules: safeBool(d.enforcedInRules ?? d.enforced_in_rules),
    createdAt: safeDate(d.createdAt ?? d.created_at ?? d.creationDate),
    updatedAt: safeDate(d.updatedAt ?? d.updated_at ?? d.modificationDate),
  };
}

/**
 * Maps an array of group DTOs to Group domain models.
 */
export function groupsFromDto(dtos: unknown[] | null | undefined): Group[] {
  return safeMap(dtos, groupFromDto);
}

/**
 * Maps a tag DTO to a Tag domain model.
 */
export function tagFromDto(dto: unknown): Tag | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id ?? d.tagId ?? d.tag_id),
    name: safeString(d.name ?? d.value),
    color: safeString(d.color),
    createdAt: safeDate(d.createdAt ?? d.created_at ?? d.creationDate),
  };
}

/**
 * Maps an array of tag DTOs to Tag domain models.
 */
export function tagsFromDto(dtos: unknown[] | null | undefined): Tag[] {
  return safeMap(dtos, tagFromDto);
}
