/**
 * Unit tests for the earn reward mappers.
 */

import { earnRewardFromDto, earnRewardsFromDto } from '../../../src/mappers/earn';

describe('earnRewardFromDto', () => {
  it('should map a Merkl token reward', () => {
    expect(
      earnRewardFromDto({
        id: 'r1',
        recipientAddressId: '42',
        recipientAddress: '0xabc',
        rewardType: 'RewardTypeMerklToken',
        merklTokenReward: {
          amount: '100',
          claimed: '40',
          pending: '60',
          token: { address: '0xtok', symbol: 'MRK', assetId: 'a1' },
        },
      })
    ).toEqual({
      id: 'r1',
      recipientAddressId: '42',
      recipientAddress: '0xabc',
      rewardType: 'RewardTypeMerklToken',
      merklTokenReward: {
        amount: '100',
        claimed: '40',
        pending: '60',
        token: { address: '0xtok', symbol: 'MRK', assetId: 'a1' },
      },
    });
  });

  it('should leave the reward details undefined when absent', () => {
    expect(earnRewardFromDto({ id: 'r2' })?.merklTokenReward).toBeUndefined();
  });

  it('should return undefined for null input and [] for an absent array', () => {
    expect(earnRewardFromDto(null)).toBeUndefined();
    expect(earnRewardsFromDto(undefined)).toEqual([]);
  });
});
