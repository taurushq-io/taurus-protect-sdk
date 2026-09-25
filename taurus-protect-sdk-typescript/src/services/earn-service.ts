/**
 * Earn service for Taurus-PROTECT SDK.
 *
 * Provides access to the rewards earned by addresses.
 */

import type { EarnApi } from '../internal/openapi/apis/EarnApi';
import { earnRewardsFromDto } from '../mappers/earn';
import type { ListEarnRewardsOptions, ListEarnRewardsResult } from '../models/earn';
import { buildCursorPage, cursorRequest } from '../models/pagination';
import { BaseService } from './base';
import { cursorQuery } from './paging';

/**
 * Service for the rewards earned by addresses.
 *
 * @example
 * ```typescript
 * let cursor: string | undefined;
 * do {
 *   const page = await client.earn.listRewards({ recipientAddressId: '42', cursor });
 *   for (const reward of page.items) {
 *     console.log(reward.rewardType, reward.merklTokenReward?.amount);
 *   }
 *   cursor = page.pagination.hasMore ? page.pagination.nextCursor : undefined;
 * } while (cursor);
 * ```
 */
export class EarnService extends BaseService {
  private readonly earnApi: EarnApi;

  /**
   * Creates a new EarnService instance.
   *
   * @param earnApi - The EarnApi instance from the OpenAPI client
   */
  constructor(earnApi: EarnApi) {
    super();
    this.earnApi = earnApi;
  }

  /**
   * Lists a page of rewards.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of rewards and its cursor pagination
   * @throws {@link ValidationError} If the paging options are invalid
   * @throws {@link APIError} If API request fails
   */
  async listRewards(options?: ListEarnRewardsOptions): Promise<ListEarnRewardsResult> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.earnApi.earnServiceGetRewards({
        recipientAddressId: options?.recipientAddressId,
        ...cursorQuery(page),
      });

      return {
        items: earnRewardsFromDto(response.rewards),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }
}
