/**
 * User, Group, and Tag mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Group, Tag, User, UserAttribute } from '../models/user';
import { UserStatus } from '../models/user';
import { safeBool, safeDate, safeMap, safeString } from './base';

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
    groupIds: Array.isArray(groupIds) ? groupIds.map((id) => String(id)) : undefined,
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
    userIds: Array.isArray(userIds) ? userIds.map((id) => String(id)) : undefined,
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
