/**
 * Group service for Taurus-PROTECT SDK.
 *
 * Provides methods for group management operations.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { GroupsApi } from '../internal/openapi/apis/GroupsApi';
import { groupFromDto, groupsFromDto } from '../mappers/user';
import {
  buildOffsetPagination,
  offsetRequest,
  type PaginatedResult,
} from '../models/pagination';
import type { Group, ListGroupsOptions } from '../models/user';
import { BaseService } from './base';
import { offsetQuery } from './paging';

/**
 * GetGroups computes enforcedInRules for each group and its users, and validatord omits a
 * false bool, so absent means false.
 */
function withComputedRulesFlags(group: Group): Group {
  return {
    ...group,
    enforcedInRules: group.enforcedInRules ?? false,
    users: group.users?.map((user) => ({
      ...user,
      enforcedInRules: user.enforcedInRules ?? false,
    })),
  };
}

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

      const result = response.result;

      if (!result || result.length === 0) {
        throw new NotFoundError(`Group ${groupId} not found`);
      }

      const group = groupFromDto(result[0]);
      if (!group) {
        throw new NotFoundError(`Group ${groupId} not found`);
      }

      return withComputedRulesFlags(group);
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
    const page = offsetRequest(options);

    return this.execute(async () => {
      const response = await this.groupsApi.userServiceGetGroups({
        ...offsetQuery(page),
        ids: options?.ids,
        query: options?.query,
      });

      const rows = response.result ?? [];
      return {
        items: groupsFromDto(rows).map(withComputedRulesFlags),
        // The server may append a synthetic technical group beyond `limit`.
        pagination: buildOffsetPagination('plus_min_rows_limit', page, response, rows.length),
      };
    });
  }
}
