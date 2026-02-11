/**
 * Group service for Taurus-PROTECT SDK.
 *
 * Provides methods for group management operations.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { GroupsApi } from '../internal/openapi/apis/GroupsApi';
import { groupFromDto, groupsFromDto } from '../mappers/user';
import type { Pagination, PaginatedResult } from '../models/pagination';
import type { Group, ListGroupsOptions } from '../models/user';
import { BaseService } from './base';

/**
 * Service for group management operations.
 *
 * Provides methods to list and retrieve groups.
 *
 * @example
 * ```typescript
 * // List groups
 * const result = await groupService.list({ limit: 50 });
 * for (const group of result.items) {
 *   console.log(`${group.name}: ${group.userIds?.length ?? 0} users`);
 * }
 *
 * // Get single group
 * const group = await groupService.get('group-123');
 * console.log(`Group: ${group.name}`);
 * ```
 */
export class GroupService extends BaseService {
  private readonly groupsApi: GroupsApi;

  /**
   * Creates a new GroupService instance.
   *
   * @param groupsApi - The GroupsApi instance from the OpenAPI client
   */
  constructor(groupsApi: GroupsApi) {
    super();
    this.groupsApi = groupsApi;
  }

  /**
   * Gets a group by ID.
   *
   * @param groupId - The group ID to retrieve
   * @returns The group
   * @throws {@link ValidationError} If groupId is invalid
   * @throws {@link NotFoundError} If group not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const group = await groupService.get('group-123');
   * console.log(`Group: ${group.name}, Members: ${group.userIds?.length}`);
   * ```
   */
  async get(groupId: string): Promise<Group> {
    if (!groupId || groupId.trim() === '') {
      throw new ValidationError('groupId is required');
    }

    return this.execute(async () => {
      // The API doesn't have a direct get-by-id endpoint,
      // so we use the list endpoint with id filter
      const response = await this.groupsApi.userServiceGetGroups({
        limit: '1',
        offset: '0',
        ids: [groupId],
      });

      const resp = response as Record<string, unknown>;
      const result = resp.result as unknown[];

      if (!result || result.length === 0) {
        throw new NotFoundError(`Group ${groupId} not found`);
      }

      const group = groupFromDto(result[0]);
      if (!group) {
        throw new NotFoundError(`Group ${groupId} not found`);
      }

      return group;
    });
  }

  /**
   * Lists groups with pagination.
   *
   * @param options - Optional filtering and pagination options
   * @returns Paginated result containing groups and pagination info
   * @throws {@link ValidationError} If limit or offset are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List first 50 groups
   * const result = await groupService.list({ limit: 50 });
   * console.log(`Found ${result.pagination.totalItems} groups`);
   *
   * // Search groups
   * const result = await groupService.list({ query: 'admin' });
   * ```
   */
  async list(options?: ListGroupsOptions): Promise<PaginatedResult<Group>> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError('limit must be positive');
    }
    if (offset < 0) {
      throw new ValidationError('offset cannot be negative');
    }

    return this.execute(async () => {
      const response = await this.groupsApi.userServiceGetGroups({
        limit: String(limit),
        offset: String(offset),
        ids: options?.ids,
        query: options?.query,
      });

      const resp = response as Record<string, unknown>;
      const result = resp.result;
      const groups = groupsFromDto(result as unknown[]);

      const pagination: Pagination = {
        totalItems: parseInt((resp.totalItems ?? resp.total_items ?? '0') as string, 10),
        offset,
        limit,
      };

      return {
        items: groups,
        pagination,
      };
    });
  }
}
