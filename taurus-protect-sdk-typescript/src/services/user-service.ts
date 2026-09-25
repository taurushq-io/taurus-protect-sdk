/**
 * User service for Taurus-PROTECT SDK.
 *
 * Provides methods for user management operations.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { UsersApi } from '../internal/openapi/apis/UsersApi';
import { userFromDto, usersFromDto } from '../mappers/user';
import {
  buildOffsetPagination,
  offsetRequest,
  type PaginatedResult,
} from '../models/pagination';
import type { ListUsersOptions, User } from '../models/user';
import { BaseService } from './base';
import { offsetQuery } from './paging';

/**
 * For the reads whose endpoint computes enforcedInRules: validatord omits a false bool, so
 * absent means false. publicKeyEnforcedInRules is a BoolValue and is never defaulted.
 */
function withComputedRulesFlags(user: User): User {
  return {
    ...user,
    enforcedInRules: user.enforcedInRules ?? false,
    groups: user.groups?.map((group) => ({
      ...group,
      enforcedInRules: group.enforcedInRules ?? false,
    })),
  };
}

/**
 * Service for user management operations.
 *
 * Provides methods to list and retrieve users.
 *
 * @example
 * ```typescript
 * // List users
 * const result = await userService.list({ limit: 50 });
 * for (const user of result.items) {
 *   console.log(`${user.email}: ${user.status}`);
 * }
 *
 * // Get single user
 * const user = await userService.get('user-123');
 * console.log(`Name: ${user.firstName} ${user.lastName}`);
 *
 * // Get current user
 * const currentUser = await userService.getCurrentUser();
 * console.log(`Logged in as: ${currentUser.email}`);
 * ```
 */
export class UserService extends BaseService {
  private readonly usersApi: UsersApi;

  /**
   * Creates a new UserService instance.
   *
   * @param usersApi - The UsersApi instance from the OpenAPI client
   */
  constructor(usersApi: UsersApi) {
    super();
    this.usersApi = usersApi;
  }

  /**
   * Gets a user by ID.
   *
   * @param userId - The user ID to retrieve
   * @returns The user
   * @throws {@link ValidationError} If userId is invalid
   * @throws {@link NotFoundError} If user not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const user = await userService.get('user-123');
   * console.log(`${user.firstName} ${user.lastName} (${user.email})`);
   * ```
   */
  async get(userId: string): Promise<User> {
    if (!userId || userId.trim() === '') {
      throw new ValidationError('userId is required');
    }

    return this.execute(async () => {
      const response = await this.usersApi.userServiceGetUser({ id: userId });

      const result = response.result;
      const user = userFromDto(result);

      if (!user) {
        throw new NotFoundError(`User ${userId} not found`);
      }

      return user;
    });
  }

  /**
   * Gets the current authenticated user.
   *
   * The reply carries the governance rules flags (`enforcedInRules`,
   * `publicKeyEnforcedInRules`, and each group's `enforcedInRules`).
   *
   * @returns The current user
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const currentUser = await userService.getCurrentUser();
   * console.log(`Logged in as: ${currentUser.email}`);
   * console.log(`Roles: ${currentUser.roles?.join(', ')}`);
   * ```
   */
  async getCurrentUser(): Promise<User> {
    return this.execute(async () => {
      // validatord computes the flags on this endpoint only when asked.
      const response = await this.usersApi.userServiceGetMe({ checkEnforcedInRules: true });

      const result = response.result;
      const user = userFromDto(result);

      if (!user) {
        throw new ValidationError('Failed to get current user: no result returned');
      }

      return withComputedRulesFlags(user);
    });
  }

  /**
   * Lists users with pagination.
   *
   * @param options - Optional filtering and pagination options
   * @returns Paginated result containing users and pagination info
   * @throws {@link ValidationError} If limit or offset are out of bounds
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List the first page of users, then the next one
   * const result = await userService.list({ limit: 50 });
   * console.log(`Found ${result.pagination.totalItems} users`);
   * if (result.pagination.hasMore) {
   *   await userService.list({ limit: 50, offset: result.pagination.nextOffset });
   * }
   *
   * // List with filters
   * const admins = await userService.list({
   *   roles: ['ADMIN'],
   *   excludeTechnicalUsers: true,
   * });
   * ```
   */
  async list(options?: ListUsersOptions): Promise<PaginatedResult<User>> {
    const page = offsetRequest(options);

    return this.execute(async () => {
      const response = await this.usersApi.userServiceGetUsers({
        ...offsetQuery(page),
        ids: options?.ids,
        externalUserIds: options?.externalUserIds,
        emails: options?.emails,
        query: options?.query,
        excludeTechnicalUsers: options?.excludeTechnicalUsers,
        roles: options?.roles,
        status: options?.status,
        totpEnabled: options?.totpEnabled,
        groupIds: options?.groupIds,
      });

      const rows = response.result ?? [];
      return {
        items: usersFromDto(rows).map(withComputedRulesFlags),
        // The server may append a synthetic daemon user beyond `limit`.
        pagination: buildOffsetPagination('plus_min_rows_limit', page, response, rows.length),
      };
    });
  }
}
