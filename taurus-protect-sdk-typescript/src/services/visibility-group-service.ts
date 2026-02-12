/**
 * Visibility group service for Taurus-PROTECT SDK.
 *
 * Provides methods for managing visibility groups in the Taurus-PROTECT system.
 * Visibility groups are used to control data access. Users can only see
 * wallets, addresses, and other entities that belong to their assigned
 * visibility groups.
 */

import { ValidationError } from '../errors';
import type { RestrictedVisibilityGroupsApi } from '../internal/openapi/apis/RestrictedVisibilityGroupsApi';
import { visibilityGroupsFromDto } from '../mappers/visibility-group';
import { usersFromDto } from '../mappers/user';
import type { VisibilityGroup } from '../models/visibility-group';
import type { User } from '../models/user';
import { BaseService } from './base';

/**
 * Service for managing visibility groups in the Taurus-PROTECT system.
 *
 * Visibility groups are used to control data access. Users can only see
 * wallets, addresses, and other entities that belong to their assigned
 * visibility groups. This provides fine-grained access control within
 * a tenant.
 *
 * @example
 * ```typescript
 * // Get all visibility groups
 * const groups = await client.visibilityGroups.list();
 * for (const group of groups) {
 *   console.log(`${group.name}: ${group.description}`);
 * }
 *
 * // Get users in a specific visibility group
 * const users = await client.visibilityGroups.getUsersByVisibilityGroup('vg-123');
 * for (const user of users) {
 *   console.log(`${user.firstName} ${user.lastName}`);
 * }
 * ```
 */
export class VisibilityGroupService extends BaseService {
  private readonly visibilityGroupsApi: RestrictedVisibilityGroupsApi;

  /**
   * Creates a new VisibilityGroupService instance.
   *
   * @param visibilityGroupsApi - The RestrictedVisibilityGroupsApi instance from the OpenAPI client
   */
  constructor(visibilityGroupsApi: RestrictedVisibilityGroupsApi) {
    super();
    this.visibilityGroupsApi = visibilityGroupsApi;
  }

  /**
   * Lists all visibility groups.
   *
   * @returns Array of visibility groups
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const groups = await visibilityGroupService.list();
   * for (const group of groups) {
   *   console.log(`${group.name}: ${group.userCount} users`);
   * }
   * ```
   */
  async list(): Promise<VisibilityGroup[]> {
    return this.execute(async () => {
      const response = await this.visibilityGroupsApi.userServiceGetVisibilityGroups();

      const result =
        (response as Record<string, unknown>).result ??
        (response as Record<string, unknown>).visibilityGroups;
      return visibilityGroupsFromDto(result as unknown[]);
    });
  }

  /**
   * Gets users assigned to a specific visibility group.
   *
   * @param visibilityGroupId - The visibility group ID
   * @returns Array of users in the visibility group
   * @throws {@link ValidationError} If visibilityGroupId is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const users = await visibilityGroupService.getUsersByVisibilityGroup('vg-123');
   * for (const user of users) {
   *   console.log(`${user.firstName} ${user.lastName} (${user.email})`);
   * }
   * ```
   */
  async getUsersByVisibilityGroup(visibilityGroupId: string): Promise<User[]> {
    if (!visibilityGroupId || visibilityGroupId.trim() === '') {
      throw new ValidationError('visibilityGroupId is required');
    }

    return this.execute(async () => {
      const response =
        await this.visibilityGroupsApi.userServiceGetUsersByVisibilityGroupID({
          visibilityGroupID: visibilityGroupId,
        });

      const result =
        (response as Record<string, unknown>).result ??
        (response as Record<string, unknown>).users;
      return usersFromDto(result as unknown[]);
    });
  }
}
