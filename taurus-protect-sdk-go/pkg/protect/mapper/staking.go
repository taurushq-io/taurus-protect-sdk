package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// StakeAccountFromDTO converts an OpenAPI StakeAccount to a domain StakeAccount.
func StakeAccountFromDTO(dto *openapi.TgvalidatordStakeAccount) *model.StakeAccount {
	if dto == nil {
		return nil
	}

	account := &model.StakeAccount{
		ID:             safeString(dto.Id),
		AddressID:      safeString(dto.AddressId),
		AccountAddress: safeString(dto.AccountAddress),
		UpdatedAtBlock: safeString(dto.UpdatedAtBlock),
	}

	// Convert account type
	if dto.AccountType != nil {
		account.AccountType = string(*dto.AccountType)
	}

	// Convert timestamps
	if dto.CreatedAt != nil {
		account.CreatedAt = *dto.CreatedAt
	}
	if dto.UpdatedAt != nil {
		account.UpdatedAt = *dto.UpdatedAt
	}

	// Convert Solana stake account
	if dto.SolanaStakeAccount != nil {
		account.SolanaStakeAccount = SolanaStakeAccountFromDTO(dto.SolanaStakeAccount)
	}

	return account
}

// StakeAccountsFromDTO converts a slice of OpenAPI StakeAccount to domain StakeAccounts.
func StakeAccountsFromDTO(dtos []openapi.TgvalidatordStakeAccount) []*model.StakeAccount {
	if dtos == nil {
		return nil
	}
	accounts := make([]*model.StakeAccount, len(dtos))
	for i := range dtos {
		accounts[i] = StakeAccountFromDTO(&dtos[i])
	}
	return accounts
}

// SolanaStakeAccountFromDTO converts an OpenAPI SolanaStakeAccount to a domain SolanaStakeAccount.
func SolanaStakeAccountFromDTO(dto *openapi.TgvalidatordSolanaStakeAccount) *model.SolanaStakeAccount {
	if dto == nil {
		return nil
	}

	account := &model.SolanaStakeAccount{
		DerivationIndex:  safeString(dto.DerivationIndex),
		ValidatorAddress: safeString(dto.ValidatorAddress),
		ActiveBalance:    safeString(dto.ActiveBalance),
		InactiveBalance:  safeString(dto.InactiveBalance),
		AllowMerge:       safeBool(dto.AllowMerge),
	}

	// Convert state
	if dto.State != nil {
		account.State = string(*dto.State)
	}

	return account
}

// ADAStakePoolInfoFromDTO converts an OpenAPI GetADAStakePoolInfoReply to a domain ADAStakePoolInfo.
func ADAStakePoolInfoFromDTO(dto *openapi.TgvalidatordGetADAStakePoolInfoReply) *model.ADAStakePoolInfo {
	if dto == nil {
		return nil
	}

	info := &model.ADAStakePoolInfo{
		Pledge:      safeString(dto.Pledge),
		FixedCost:   safeString(dto.FixedCost),
		URL:         safeString(dto.Url),
		ActiveStake: safeString(dto.ActiveStake),
		Epoch:       safeString(dto.Epoch),
	}

	if dto.Margin != nil {
		info.Margin = *dto.Margin
	}

	return info
}

// ETHValidatorInfoFromDTO converts an OpenAPI ETHValidatorInfo to a domain ETHValidatorInfo.
func ETHValidatorInfoFromDTO(dto *openapi.TgvalidatordETHValidatorInfo) *model.ETHValidatorInfo {
	if dto == nil {
		return nil
	}

	return &model.ETHValidatorInfo{
		ID:        safeString(dto.Id),
		Pubkey:    safeString(dto.Pubkey),
		Status:    safeString(dto.Status),
		Balance:   safeString(dto.Balance),
		Network:   safeString(dto.Network),
		Provider:  safeString(dto.Provider),
		AddressID: safeString(dto.AddressID),
	}
}

// ETHValidatorsInfoFromDTO converts a slice of OpenAPI ETHValidatorInfo to domain ETHValidatorInfo.
func ETHValidatorsInfoFromDTO(dtos []openapi.TgvalidatordETHValidatorInfo) []*model.ETHValidatorInfo {
	if dtos == nil {
		return nil
	}
	validators := make([]*model.ETHValidatorInfo, len(dtos))
	for i := range dtos {
		validators[i] = ETHValidatorInfoFromDTO(&dtos[i])
	}
	return validators
}

// FTMValidatorInfoFromDTO converts an OpenAPI GetFTMValidatorInfoReply to a domain FTMValidatorInfo.
func FTMValidatorInfoFromDTO(dto *openapi.TgvalidatordGetFTMValidatorInfoReply) *model.FTMValidatorInfo {
	if dto == nil {
		return nil
	}

	return &model.FTMValidatorInfo{
		ValidatorID:          safeString(dto.ValidatorID),
		Address:              safeString(dto.Address),
		IsActive:             safeBool(dto.IsActive),
		TotalStake:           safeString(dto.TotalStake),
		SelfStake:            safeString(dto.SelfStake),
		DeactivatedAtDateUnix: safeString(dto.DeactivatedAtDateUnix),
		CreatedAtDateUnix:    safeString(dto.CreatedAtDateUnix),
	}
}

// ICPNeuronInfoFromDTO converts an OpenAPI GetICPNeuronInfoReply to a domain ICPNeuronInfo.
func ICPNeuronInfoFromDTO(dto *openapi.TgvalidatordGetICPNeuronInfoReply) *model.ICPNeuronInfo {
	if dto == nil {
		return nil
	}

	info := &model.ICPNeuronInfo{
		NeuronID:                            safeString(dto.NeuronId),
		RetrieveAtTimestampSeconds:          safeString(dto.RetrieveAtTimestampSeconds),
		AgeSeconds:                          safeString(dto.AgeSeconds),
		DissolveDelaySeconds:                safeString(dto.DissolveDelaySeconds),
		VotingPower:                         safeString(dto.VotingPower),
		CreatedTimestampSeconds:             safeString(dto.CreatedTimestampSeconds),
		StakeE8S:                            safeString(dto.StakeE8S),
		JoinedCommunityFundTimestampSeconds: safeString(dto.JoinedCommunityFundTimestampSeconds),
	}

	// Convert neuron state
	if dto.NeuronState != nil {
		info.NeuronState = string(*dto.NeuronState)
	}

	// Convert known neuron data
	if dto.KnownNeuronData != nil {
		info.KnownNeuronData = ICPKnownNeuronDataFromDTO(dto.KnownNeuronData)
	}

	return info
}

// ICPKnownNeuronDataFromDTO converts an OpenAPI ICPKnownNeuronData to a domain ICPKnownNeuronData.
func ICPKnownNeuronDataFromDTO(dto *openapi.TgvalidatordICPKnownNeuronData) *model.ICPKnownNeuronData {
	if dto == nil {
		return nil
	}

	return &model.ICPKnownNeuronData{
		Name:        safeString(dto.Name),
		Description: safeString(dto.Description),
	}
}

// NEARValidatorInfoFromDTO converts an OpenAPI GetNEARValidatorInfoReply to a domain NEARValidatorInfo.
func NEARValidatorInfoFromDTO(dto *openapi.TgvalidatordGetNEARValidatorInfoReply) *model.NEARValidatorInfo {
	if dto == nil {
		return nil
	}

	info := &model.NEARValidatorInfo{
		ValidatorAddress:   safeString(dto.ValidatorAddress),
		OwnerID:            safeString(dto.OwnerId),
		TotalStakedBalance: safeString(dto.TotalStakedBalance),
		StakingKey:         safeString(dto.StakingKey),
		IsStakingPaused:    safeBool(dto.IsStakingPaused),
	}

	if dto.RewardFeeFraction != nil {
		info.RewardFeeFraction = *dto.RewardFeeFraction
	}

	return info
}

// XTZStakingRewardFromDTO converts an OpenAPI GetXTZAddressStakingRewardsReply to a domain XTZStakingReward.
func XTZStakingRewardFromDTO(dto *openapi.TgvalidatordGetXTZAddressStakingRewardsReply) *model.XTZStakingReward {
	if dto == nil {
		return nil
	}

	return &model.XTZStakingReward{
		ReceivedRewardsAmount: safeString(dto.ReceivedRewardsAmount),
	}
}
