package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestStakeAccountFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordStakeAccount
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns stake account with zero values",
			dto:  &openapi.TgvalidatordStakeAccount{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordStakeAccount {
				id := "stake-123"
				addressId := "addr-456"
				accountAddress := "Sol1234567890"
				updatedAtBlock := "12345"
				createdAt := time.Now().Add(-24 * time.Hour)
				updatedAt := time.Now()
				accountType := openapi.TGVALIDATORDSTAKEACCOUNTTYPE_STAKE_ACCOUNT_TYPE_SOLANA
				return &openapi.TgvalidatordStakeAccount{
					Id:             &id,
					AddressId:      &addressId,
					AccountAddress: &accountAddress,
					UpdatedAtBlock: &updatedAtBlock,
					CreatedAt:      &createdAt,
					UpdatedAt:      &updatedAt,
					AccountType:    &accountType,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StakeAccountFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("StakeAccountFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("StakeAccountFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.AddressId != nil && got.AddressID != *tt.dto.AddressId {
				t.Errorf("AddressID = %v, want %v", got.AddressID, *tt.dto.AddressId)
			}
			if tt.dto.AccountAddress != nil && got.AccountAddress != *tt.dto.AccountAddress {
				t.Errorf("AccountAddress = %v, want %v", got.AccountAddress, *tt.dto.AccountAddress)
			}
			if tt.dto.UpdatedAtBlock != nil && got.UpdatedAtBlock != *tt.dto.UpdatedAtBlock {
				t.Errorf("UpdatedAtBlock = %v, want %v", got.UpdatedAtBlock, *tt.dto.UpdatedAtBlock)
			}
			if tt.dto.AccountType != nil && got.AccountType != string(*tt.dto.AccountType) {
				t.Errorf("AccountType = %v, want %v", got.AccountType, string(*tt.dto.AccountType))
			}
		})
	}
}

func TestStakeAccountsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordStakeAccount
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordStakeAccount{},
			want: 0,
		},
		{
			name: "converts multiple stake accounts",
			dtos: func() []openapi.TgvalidatordStakeAccount {
				id1 := "stake-1"
				id2 := "stake-2"
				return []openapi.TgvalidatordStakeAccount{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StakeAccountsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("StakeAccountsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("StakeAccountsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestSolanaStakeAccountFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordSolanaStakeAccount
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns account with zero values",
			dto:  &openapi.TgvalidatordSolanaStakeAccount{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordSolanaStakeAccount {
				derivationIndex := "0"
				validatorAddress := "Val123"
				activeBalance := "1000000000"
				inactiveBalance := "500000000"
				allowMerge := true
				return &openapi.TgvalidatordSolanaStakeAccount{
					DerivationIndex:  &derivationIndex,
					ValidatorAddress: &validatorAddress,
					ActiveBalance:    &activeBalance,
					InactiveBalance:  &inactiveBalance,
					AllowMerge:       &allowMerge,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SolanaStakeAccountFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("SolanaStakeAccountFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("SolanaStakeAccountFromDTO() returned nil for non-nil input")
			}
			if tt.dto.DerivationIndex != nil && got.DerivationIndex != *tt.dto.DerivationIndex {
				t.Errorf("DerivationIndex = %v, want %v", got.DerivationIndex, *tt.dto.DerivationIndex)
			}
			if tt.dto.ValidatorAddress != nil && got.ValidatorAddress != *tt.dto.ValidatorAddress {
				t.Errorf("ValidatorAddress = %v, want %v", got.ValidatorAddress, *tt.dto.ValidatorAddress)
			}
			if tt.dto.ActiveBalance != nil && got.ActiveBalance != *tt.dto.ActiveBalance {
				t.Errorf("ActiveBalance = %v, want %v", got.ActiveBalance, *tt.dto.ActiveBalance)
			}
			if tt.dto.InactiveBalance != nil && got.InactiveBalance != *tt.dto.InactiveBalance {
				t.Errorf("InactiveBalance = %v, want %v", got.InactiveBalance, *tt.dto.InactiveBalance)
			}
			if tt.dto.AllowMerge != nil && got.AllowMerge != *tt.dto.AllowMerge {
				t.Errorf("AllowMerge = %v, want %v", got.AllowMerge, *tt.dto.AllowMerge)
			}
		})
	}
}

func TestADAStakePoolInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordGetADAStakePoolInfoReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns info with zero values",
			dto:  &openapi.TgvalidatordGetADAStakePoolInfoReply{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordGetADAStakePoolInfoReply {
				pledge := "100000000"
				margin := float32(0.05)
				fixedCost := "340000000"
				url := "https://pool.example.com"
				activeStake := "5000000000"
				epoch := "400"
				return &openapi.TgvalidatordGetADAStakePoolInfoReply{
					Pledge:      &pledge,
					Margin:      &margin,
					FixedCost:   &fixedCost,
					Url:         &url,
					ActiveStake: &activeStake,
					Epoch:       &epoch,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ADAStakePoolInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ADAStakePoolInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ADAStakePoolInfoFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Pledge != nil && got.Pledge != *tt.dto.Pledge {
				t.Errorf("Pledge = %v, want %v", got.Pledge, *tt.dto.Pledge)
			}
			if tt.dto.Margin != nil && got.Margin != *tt.dto.Margin {
				t.Errorf("Margin = %v, want %v", got.Margin, *tt.dto.Margin)
			}
			if tt.dto.FixedCost != nil && got.FixedCost != *tt.dto.FixedCost {
				t.Errorf("FixedCost = %v, want %v", got.FixedCost, *tt.dto.FixedCost)
			}
			if tt.dto.Url != nil && got.URL != *tt.dto.Url {
				t.Errorf("URL = %v, want %v", got.URL, *tt.dto.Url)
			}
			if tt.dto.ActiveStake != nil && got.ActiveStake != *tt.dto.ActiveStake {
				t.Errorf("ActiveStake = %v, want %v", got.ActiveStake, *tt.dto.ActiveStake)
			}
			if tt.dto.Epoch != nil && got.Epoch != *tt.dto.Epoch {
				t.Errorf("Epoch = %v, want %v", got.Epoch, *tt.dto.Epoch)
			}
		})
	}
}

func TestETHValidatorInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordETHValidatorInfo
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns info with zero values",
			dto:  &openapi.TgvalidatordETHValidatorInfo{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordETHValidatorInfo {
				id := "val-123"
				pubkey := "0x1234567890abcdef"
				status := "active"
				balance := "32000000000"
				network := "mainnet"
				provider := "example-provider"
				addressID := "addr-456"
				return &openapi.TgvalidatordETHValidatorInfo{
					Id:        &id,
					Pubkey:    &pubkey,
					Status:    &status,
					Balance:   &balance,
					Network:   &network,
					Provider:  &provider,
					AddressID: &addressID,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ETHValidatorInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ETHValidatorInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ETHValidatorInfoFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Pubkey != nil && got.Pubkey != *tt.dto.Pubkey {
				t.Errorf("Pubkey = %v, want %v", got.Pubkey, *tt.dto.Pubkey)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.Balance != nil && got.Balance != *tt.dto.Balance {
				t.Errorf("Balance = %v, want %v", got.Balance, *tt.dto.Balance)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
			if tt.dto.Provider != nil && got.Provider != *tt.dto.Provider {
				t.Errorf("Provider = %v, want %v", got.Provider, *tt.dto.Provider)
			}
			if tt.dto.AddressID != nil && got.AddressID != *tt.dto.AddressID {
				t.Errorf("AddressID = %v, want %v", got.AddressID, *tt.dto.AddressID)
			}
		})
	}
}

func TestETHValidatorsInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordETHValidatorInfo
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1,
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordETHValidatorInfo{},
			want: 0,
		},
		{
			name: "converts multiple validators",
			dtos: func() []openapi.TgvalidatordETHValidatorInfo {
				id1 := "val-1"
				id2 := "val-2"
				return []openapi.TgvalidatordETHValidatorInfo{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ETHValidatorsInfoFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("ETHValidatorsInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("ETHValidatorsInfoFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestFTMValidatorInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordGetFTMValidatorInfoReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns info with zero values",
			dto:  &openapi.TgvalidatordGetFTMValidatorInfoReply{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordGetFTMValidatorInfoReply {
				validatorID := "1"
				address := "0xval123"
				isActive := true
				totalStake := "1000000000000000000"
				selfStake := "500000000000000000"
				deactivatedAt := "0"
				createdAt := "1640000000"
				return &openapi.TgvalidatordGetFTMValidatorInfoReply{
					ValidatorID:           &validatorID,
					Address:               &address,
					IsActive:              &isActive,
					TotalStake:            &totalStake,
					SelfStake:             &selfStake,
					DeactivatedAtDateUnix: &deactivatedAt,
					CreatedAtDateUnix:     &createdAt,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FTMValidatorInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("FTMValidatorInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("FTMValidatorInfoFromDTO() returned nil for non-nil input")
			}
			if tt.dto.ValidatorID != nil && got.ValidatorID != *tt.dto.ValidatorID {
				t.Errorf("ValidatorID = %v, want %v", got.ValidatorID, *tt.dto.ValidatorID)
			}
			if tt.dto.Address != nil && got.Address != *tt.dto.Address {
				t.Errorf("Address = %v, want %v", got.Address, *tt.dto.Address)
			}
			if tt.dto.IsActive != nil && got.IsActive != *tt.dto.IsActive {
				t.Errorf("IsActive = %v, want %v", got.IsActive, *tt.dto.IsActive)
			}
			if tt.dto.TotalStake != nil && got.TotalStake != *tt.dto.TotalStake {
				t.Errorf("TotalStake = %v, want %v", got.TotalStake, *tt.dto.TotalStake)
			}
			if tt.dto.SelfStake != nil && got.SelfStake != *tt.dto.SelfStake {
				t.Errorf("SelfStake = %v, want %v", got.SelfStake, *tt.dto.SelfStake)
			}
		})
	}
}

func TestICPNeuronInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordGetICPNeuronInfoReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns info with zero values",
			dto:  &openapi.TgvalidatordGetICPNeuronInfoReply{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordGetICPNeuronInfoReply {
				neuronId := "123456789"
				retrieveAt := "1700000000"
				ageSeconds := "86400"
				dissolveDelay := "15552000"
				votingPower := "100000000"
				createdAt := "1640000000"
				stakeE8S := "100000000"
				joinedFund := "1650000000"
				name := "My Neuron"
				desc := "A known neuron"
				return &openapi.TgvalidatordGetICPNeuronInfoReply{
					NeuronId:                            &neuronId,
					RetrieveAtTimestampSeconds:          &retrieveAt,
					AgeSeconds:                          &ageSeconds,
					DissolveDelaySeconds:                &dissolveDelay,
					VotingPower:                         &votingPower,
					CreatedTimestampSeconds:             &createdAt,
					StakeE8S:                            &stakeE8S,
					JoinedCommunityFundTimestampSeconds: &joinedFund,
					KnownNeuronData: &openapi.TgvalidatordICPKnownNeuronData{
						Name:        &name,
						Description: &desc,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ICPNeuronInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ICPNeuronInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ICPNeuronInfoFromDTO() returned nil for non-nil input")
			}
			if tt.dto.NeuronId != nil && got.NeuronID != *tt.dto.NeuronId {
				t.Errorf("NeuronID = %v, want %v", got.NeuronID, *tt.dto.NeuronId)
			}
			if tt.dto.RetrieveAtTimestampSeconds != nil && got.RetrieveAtTimestampSeconds != *tt.dto.RetrieveAtTimestampSeconds {
				t.Errorf("RetrieveAtTimestampSeconds = %v, want %v", got.RetrieveAtTimestampSeconds, *tt.dto.RetrieveAtTimestampSeconds)
			}
			if tt.dto.AgeSeconds != nil && got.AgeSeconds != *tt.dto.AgeSeconds {
				t.Errorf("AgeSeconds = %v, want %v", got.AgeSeconds, *tt.dto.AgeSeconds)
			}
			if tt.dto.KnownNeuronData != nil {
				if got.KnownNeuronData == nil {
					t.Error("KnownNeuronData should not be nil when DTO has it")
				}
			}
		})
	}
}

func TestICPKnownNeuronDataFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordICPKnownNeuronData
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns data with zero values",
			dto:  &openapi.TgvalidatordICPKnownNeuronData{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordICPKnownNeuronData {
				name := "My Neuron"
				desc := "A known neuron for voting"
				return &openapi.TgvalidatordICPKnownNeuronData{
					Name:        &name,
					Description: &desc,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ICPKnownNeuronDataFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ICPKnownNeuronDataFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ICPKnownNeuronDataFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Description != nil && got.Description != *tt.dto.Description {
				t.Errorf("Description = %v, want %v", got.Description, *tt.dto.Description)
			}
		})
	}
}

func TestNEARValidatorInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordGetNEARValidatorInfoReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns info with zero values",
			dto:  &openapi.TgvalidatordGetNEARValidatorInfoReply{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordGetNEARValidatorInfoReply {
				validatorAddress := "validator.poolv1.near"
				ownerId := "owner.near"
				totalStakedBalance := "1000000000000000000000000"
				rewardFeeFraction := float32(0.1)
				stakingKey := "ed25519:abc123"
				isStakingPaused := false
				return &openapi.TgvalidatordGetNEARValidatorInfoReply{
					ValidatorAddress:   &validatorAddress,
					OwnerId:            &ownerId,
					TotalStakedBalance: &totalStakedBalance,
					RewardFeeFraction:  &rewardFeeFraction,
					StakingKey:         &stakingKey,
					IsStakingPaused:    &isStakingPaused,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NEARValidatorInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("NEARValidatorInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("NEARValidatorInfoFromDTO() returned nil for non-nil input")
			}
			if tt.dto.ValidatorAddress != nil && got.ValidatorAddress != *tt.dto.ValidatorAddress {
				t.Errorf("ValidatorAddress = %v, want %v", got.ValidatorAddress, *tt.dto.ValidatorAddress)
			}
			if tt.dto.OwnerId != nil && got.OwnerID != *tt.dto.OwnerId {
				t.Errorf("OwnerID = %v, want %v", got.OwnerID, *tt.dto.OwnerId)
			}
			if tt.dto.TotalStakedBalance != nil && got.TotalStakedBalance != *tt.dto.TotalStakedBalance {
				t.Errorf("TotalStakedBalance = %v, want %v", got.TotalStakedBalance, *tt.dto.TotalStakedBalance)
			}
			if tt.dto.RewardFeeFraction != nil && got.RewardFeeFraction != *tt.dto.RewardFeeFraction {
				t.Errorf("RewardFeeFraction = %v, want %v", got.RewardFeeFraction, *tt.dto.RewardFeeFraction)
			}
			if tt.dto.StakingKey != nil && got.StakingKey != *tt.dto.StakingKey {
				t.Errorf("StakingKey = %v, want %v", got.StakingKey, *tt.dto.StakingKey)
			}
			if tt.dto.IsStakingPaused != nil && got.IsStakingPaused != *tt.dto.IsStakingPaused {
				t.Errorf("IsStakingPaused = %v, want %v", got.IsStakingPaused, *tt.dto.IsStakingPaused)
			}
		})
	}
}

func TestXTZStakingRewardFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordGetXTZAddressStakingRewardsReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns reward with zero values",
			dto:  &openapi.TgvalidatordGetXTZAddressStakingRewardsReply{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordGetXTZAddressStakingRewardsReply {
				amount := "1234567890"
				return &openapi.TgvalidatordGetXTZAddressStakingRewardsReply{
					ReceivedRewardsAmount: &amount,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := XTZStakingRewardFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("XTZStakingRewardFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("XTZStakingRewardFromDTO() returned nil for non-nil input")
			}
			if tt.dto.ReceivedRewardsAmount != nil && got.ReceivedRewardsAmount != *tt.dto.ReceivedRewardsAmount {
				t.Errorf("ReceivedRewardsAmount = %v, want %v", got.ReceivedRewardsAmount, *tt.dto.ReceivedRewardsAmount)
			}
		})
	}
}

func TestStakeAccountFromDTO_WithSolanaStakeAccount(t *testing.T) {
	id := "stake-123"
	derivationIndex := "0"
	validatorAddress := "Val123"
	activeBalance := "1000000000"

	dto := &openapi.TgvalidatordStakeAccount{
		Id: &id,
		SolanaStakeAccount: &openapi.TgvalidatordSolanaStakeAccount{
			DerivationIndex:  &derivationIndex,
			ValidatorAddress: &validatorAddress,
			ActiveBalance:    &activeBalance,
		},
	}

	got := StakeAccountFromDTO(dto)
	if got == nil {
		t.Fatal("StakeAccountFromDTO() returned nil for non-nil input")
	}
	if got.SolanaStakeAccount == nil {
		t.Fatal("SolanaStakeAccount should not be nil when DTO has it")
	}
	if got.SolanaStakeAccount.DerivationIndex != derivationIndex {
		t.Errorf("SolanaStakeAccount.DerivationIndex = %v, want %v", got.SolanaStakeAccount.DerivationIndex, derivationIndex)
	}
	if got.SolanaStakeAccount.ValidatorAddress != validatorAddress {
		t.Errorf("SolanaStakeAccount.ValidatorAddress = %v, want %v", got.SolanaStakeAccount.ValidatorAddress, validatorAddress)
	}
}

func TestStakeAccountFromDTO_NilSolanaStakeAccount(t *testing.T) {
	id := "stake-123"
	dto := &openapi.TgvalidatordStakeAccount{
		Id:                 &id,
		SolanaStakeAccount: nil,
	}

	got := StakeAccountFromDTO(dto)
	if got == nil {
		t.Fatal("StakeAccountFromDTO() returned nil for non-nil input")
	}
	if got.SolanaStakeAccount != nil {
		t.Errorf("SolanaStakeAccount should be nil when DTO has nil, got %v", got.SolanaStakeAccount)
	}
}
