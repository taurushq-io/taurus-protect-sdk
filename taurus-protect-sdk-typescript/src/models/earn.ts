/**
 * Earn models for Taurus-PROTECT SDK.
 *
 * Rewards earned by addresses, such as Merkl token incentives.
 */

import type { CursorPage, CursorPageOptions } from './pagination';

/** The token a reward is paid in. */
export interface EarnRewardToken {
  /** Token contract address */
  readonly address?: string;
  /** Token symbol */
  readonly symbol?: string;
  /** Asset ID of the token */
  readonly assetId?: string;
}

/** A Merkl token reward. */
export interface MerklTokenReward {
  /** Total reward amount */
  readonly amount?: string;
  /** Amount already claimed */
  readonly claimed?: string;
  /** Amount pending */
  readonly pending?: string;
  /** The token the reward is paid in */
  readonly token?: EarnRewardToken;
}

/**
 * A reward earned by an address.
 */
export interface EarnReward {
  /** Reward ID */
  readonly id?: string;
  /** Internal ID of the address receiving the reward */
  readonly recipientAddressId?: string;
  /** Blockchain address receiving the reward */
  readonly recipientAddress?: string;
  /** Reward type (e.g., "RewardTypeMerklToken") */
  readonly rewardType?: string;
  /** Details of a Merkl token reward */
  readonly merklTokenReward?: MerklTokenReward;
}

/**
 * Options for listing rewards. A cursor list: `pageSize` 1-100 (default 20) and the
 * `cursor` of a previous page.
 */
export interface ListEarnRewardsOptions extends CursorPageOptions {
  /** Keep the rewards of this recipient address */
  readonly recipientAddressId?: string;
}

/**
 * A page of rewards.
 */
export interface ListEarnRewardsResult {
  /** The rewards of this page */
  readonly items: EarnReward[];
  /** Cursor pagination: continue with `cursor: pagination.nextCursor` while `hasMore` */
  readonly pagination: CursorPage;
}
