/**
 * User service for Taurus-PROTECT SDK.
 *
 * Provides methods for user management operations.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { UsersApi } from '../internal/openapi/apis/UsersApi';
import { userFromDto, usersFromDto } from '../mappers/user';
import type { Pagination, PaginatedResult } from '../models/pagination';
import type { ListUsersOptions, User } from '../models/user';
import { BaseService } from './base';

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

      const result = (response as Record<string, unknown>).result;
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
      const response = await this.usersApi.userServiceGetMe({});

      const result = (response as Record<string, unknown>).result;
      const user = userFromDto(result);

      if (!user) {
        throw new ValidationError('Failed to get current user: no result returned');
      }

      return user;
    });
  }

  /**
   * Lists users with pagination.
   *
   * @param options - Optional filtering and pagination options
   * @returns Paginated result containing users and pagination info
   * @throws {@link ValidationError} If limit or offset are invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List first 50 users
   * const result = await userService.list({ limit: 50, offset: 0 });
   * console.log(`Found ${result.pagination.totalItems} users`);
   *
   * // List with filters
   * const admins = await userService.list({
   *   roles: ['ADMIN'],
   *   excludeTechnicalUsers: true,
   * });
   * ```
   */
  async list(options?: ListUsersOptions): Promise<PaginatedResult<User>> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError('limit must be positive');
    }
    if (offset < 0) {
      throw new ValidationError('offset cannot be negative');
    }

    return this.execute(async () => {
      const response = await this.usersApi.userServiceGetUsers({
        limit: String(limit),
        offset: String(offset),
        ids: options?.ids,
        emails: options?.emails,
        query: options?.query,
        excludeTechnicalUsers: options?.excludeTechnicalUsers,
        roles: options?.roles,
        status: options?.status,
        totpEnabled: options?.totpEnabled,
        groupIds: options?.groupIds,
      });

      const resp = response as Record<string, unknown>;
      const result = resp.result;
      const users = usersFromDto(result as unknown[]);

      const pagination: Pagination = {
        totalItems: parseInt((resp.totalItems ?? resp.total_items ?? '0') as string, 10),
        offset,
        limit,
      };

      return {
        items: users,
        pagination,
      };
    });
  }
}
