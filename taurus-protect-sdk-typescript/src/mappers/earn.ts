/**
 * Earn mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { TgvalidatordEarnReward } from '../internal/openapi/models/TgvalidatordEarnReward';
import type { EarnReward } from '../models/earn';
import { safeMap } from './base';

/**
 * Maps a reward row to an EarnReward domain model.
 */
export function earnRewardFromDto(
  dto: TgvalidatordEarnReward | null | undefined
): EarnReward | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }
  const merkl = dto.merklTokenReward;
  return {
    id: dto.id,
    recipientAddressId: dto.recipientAddressId,
    recipientAddress: dto.recipientAddress,
    rewardType: dto.rewardType,
    merklTokenReward: merkl
      ? {
          amount: merkl.amount,
          claimed: merkl.claimed,
          pending: merkl.pending,
          token: merkl.token
            ? {
                address: merkl.token.address,
                symbol: merkl.token.symbol,
                assetId: merkl.token.assetId,
              }
            : undefined,
        }
      : undefined,
  };
}

/**
 * Maps reward rows to EarnReward domain models.
 */
export function earnRewardsFromDto(
  dtos: TgvalidatordEarnReward[] | null | undefined
): EarnReward[] {
  return safeMap(dtos, earnRewardFromDto);
}
