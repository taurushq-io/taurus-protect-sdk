package model

// EarnReward is a reward earned by an address (EarnService.ListRewards).
type EarnReward struct {
	// ID is the reward's identifier.
	ID string `json:"id"`
	// RecipientAddressID and RecipientAddress identify the address the reward accrues to.
	RecipientAddressID string `json:"recipient_address_id,omitempty"`
	RecipientAddress   string `json:"recipient_address,omitempty"`
	// RewardType is the kind of reward, e.g. RewardTypeMerklToken.
	RewardType string `json:"reward_type,omitempty"`
	// MerklToken is set for a Merkl token reward.
	MerklToken *MerklTokenReward `json:"merkl_token,omitempty"`
}

// MerklTokenReward is the detail of a Merkl token reward.
type MerklTokenReward struct {
	// Amount, Claimed and Pending are token amounts.
	Amount  string `json:"amount,omitempty"`
	Claimed string `json:"claimed,omitempty"`
	Pending string `json:"pending,omitempty"`
	// Token is the rewarded token.
	Token *EarnRewardToken `json:"token,omitempty"`
}

// EarnRewardToken identifies a rewarded token.
type EarnRewardToken struct {
	Address string `json:"address,omitempty"`
	Symbol  string `json:"symbol,omitempty"`
	AssetID string `json:"asset_id,omitempty"`
}

// ListEarnRewardsOptions filters and pages EarnService.ListRewards.
type ListEarnRewardsOptions struct {
	// RecipientAddressID keeps the rewards of one address.
	RecipientAddressID string
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string
}

// ListEarnRewardsResult is one page of EarnService.ListRewards.
type ListEarnRewardsResult struct {
	Rewards []*EarnReward `json:"rewards"`
	// Page continues the list: pass Page.NextCursor as the next Cursor until HasMore is false.
	Page CursorPage `json:"page"`
}
