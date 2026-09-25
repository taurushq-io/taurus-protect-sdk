package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// EarnRewardFromDTO converts an OpenAPI EarnReward to a domain EarnReward.
func EarnRewardFromDTO(dto *openapi.TgvalidatordEarnReward) *model.EarnReward {
	if dto == nil {
		return nil
	}
	reward := &model.EarnReward{
		ID:                 safeString(dto.Id),
		RecipientAddressID: safeString(dto.RecipientAddressId),
		RecipientAddress:   safeString(dto.RecipientAddress),
	}
	if dto.RewardType != nil {
		reward.RewardType = string(*dto.RewardType)
	}
	if m := dto.MerklTokenReward; m != nil {
		merkl := &model.MerklTokenReward{
			Amount:  safeString(m.Amount),
			Claimed: safeString(m.Claimed),
			Pending: safeString(m.Pending),
		}
		if t := m.Token; t != nil {
			merkl.Token = &model.EarnRewardToken{
				Address: safeString(t.Address),
				Symbol:  safeString(t.Symbol),
				AssetID: safeString(t.AssetId),
			}
		}
		reward.MerklToken = merkl
	}
	return reward
}

// EarnRewardsFromDTO converts a page of OpenAPI EarnRewards.
func EarnRewardsFromDTO(dtos []openapi.TgvalidatordEarnReward) []*model.EarnReward {
	if dtos == nil {
		return nil
	}
	rewards := make([]*model.EarnReward, len(dtos))
	for i := range dtos {
		rewards[i] = EarnRewardFromDTO(&dtos[i])
	}
	return rewards
}
