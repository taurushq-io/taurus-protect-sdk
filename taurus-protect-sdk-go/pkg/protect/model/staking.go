package model

import "time"

// StakeAccount represents a stake account in the system.
type StakeAccount struct {
	// ID is the unique identifier for the stake account.
	ID string `json:"id"`
	// AddressID is the identifier for the Protect address that the stake account is derived from.
	AddressID string `json:"address_id"`
	// AccountAddress is the address of the stake account.
	AccountAddress string `json:"account_address"`
	// CreatedAt is when the stake account was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the stake account was last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// UpdatedAtBlock is the block at which the stake account was last updated.
	UpdatedAtBlock string `json:"updated_at_block"`
	// AccountType is the type of stake account.
	AccountType string `json:"account_type"`
	// SolanaStakeAccount contains Solana-specific stake account information.
	SolanaStakeAccount *SolanaStakeAccount `json:"solana_stake_account,omitempty"`
}

// SolanaStakeAccount represents Solana-specific stake account information.
type SolanaStakeAccount struct {
	// DerivationIndex is the derivation index used for the stake account.
	DerivationIndex string `json:"derivation_index"`
	// State is the current state of the stake account.
	State string `json:"state"`
	// ValidatorAddress is the address of the validator.
	ValidatorAddress string `json:"validator_address"`
	// ActiveBalance is the active balance in the stake account.
	ActiveBalance string `json:"active_balance"`
	// InactiveBalance is the inactive balance in the stake account.
	InactiveBalance string `json:"inactive_balance"`
	// AllowMerge indicates whether the stake account can be merged.
	AllowMerge bool `json:"allow_merge"`
}

// ADAStakePoolInfo represents information about an ADA stake pool.
type ADAStakePoolInfo struct {
	// Pledge is the pledge amount of the stake pool.
	Pledge string `json:"pledge"`
	// Margin is the margin percentage of the stake pool.
	Margin float32 `json:"margin"`
	// FixedCost is the fixed cost of the stake pool.
	FixedCost string `json:"fixed_cost"`
	// URL is the URL of the stake pool metadata.
	URL string `json:"url"`
	// ActiveStake is the active stake in the pool.
	ActiveStake string `json:"active_stake"`
	// Epoch is the current epoch.
	Epoch string `json:"epoch"`
}

// ETHValidatorInfo represents information about an ETH validator.
type ETHValidatorInfo struct {
	// ID is the internal ID of the validator.
	ID string `json:"id"`
	// Pubkey is the public key of the validator.
	Pubkey string `json:"pubkey"`
	// Status is the current status of the validator.
	Status string `json:"status"`
	// Balance is the balance of the validator.
	Balance string `json:"balance"`
	// Network is the network the validator is on.
	Network string `json:"network"`
	// Provider is the staking provider.
	Provider string `json:"provider"`
	// AddressID is the Protect address ID associated with the validator.
	AddressID string `json:"address_id"`
}

// FTMValidatorInfo represents information about an FTM validator.
type FTMValidatorInfo struct {
	// ValidatorID is the ID of the validator.
	ValidatorID string `json:"validator_id"`
	// Address is the address of the validator.
	Address string `json:"address"`
	// IsActive indicates whether the validator is active.
	IsActive bool `json:"is_active"`
	// TotalStake is the total stake in the validator.
	TotalStake string `json:"total_stake"`
	// SelfStake is the validator's own stake.
	SelfStake string `json:"self_stake"`
	// DeactivatedAtDateUnix is the Unix timestamp when the validator was deactivated.
	DeactivatedAtDateUnix string `json:"deactivated_at_date_unix"`
	// CreatedAtDateUnix is the Unix timestamp when the validator was created.
	CreatedAtDateUnix string `json:"created_at_date_unix"`
}

// ICPNeuronInfo represents information about an ICP neuron.
type ICPNeuronInfo struct {
	// NeuronID is the ID of the neuron.
	NeuronID string `json:"neuron_id"`
	// RetrieveAtTimestampSeconds is the timestamp when the info was retrieved.
	RetrieveAtTimestampSeconds string `json:"retrieve_at_timestamp_seconds"`
	// NeuronState is the current state of the neuron.
	NeuronState string `json:"neuron_state"`
	// AgeSeconds is the age of the neuron in seconds.
	AgeSeconds string `json:"age_seconds"`
	// DissolveDelaySeconds is the dissolve delay in seconds.
	DissolveDelaySeconds string `json:"dissolve_delay_seconds"`
	// VotingPower is the voting power of the neuron.
	VotingPower string `json:"voting_power"`
	// CreatedTimestampSeconds is when the neuron was created.
	CreatedTimestampSeconds string `json:"created_timestamp_seconds"`
	// StakeE8S is the stake in e8s (smallest ICP unit).
	StakeE8S string `json:"stake_e8s"`
	// JoinedCommunityFundTimestampSeconds is when the neuron joined the community fund.
	JoinedCommunityFundTimestampSeconds string `json:"joined_community_fund_timestamp_seconds"`
	// KnownNeuronData contains data about known neurons.
	KnownNeuronData *ICPKnownNeuronData `json:"known_neuron_data,omitempty"`
}

// ICPKnownNeuronData represents data about a known ICP neuron.
type ICPKnownNeuronData struct {
	// Name is the name of the known neuron.
	Name string `json:"name"`
	// Description is the description of the known neuron.
	Description string `json:"description"`
}

// NEARValidatorInfo represents information about a NEAR validator.
type NEARValidatorInfo struct {
	// ValidatorAddress is the address of the validator.
	ValidatorAddress string `json:"validator_address"`
	// OwnerID is the owner ID of the validator.
	OwnerID string `json:"owner_id"`
	// TotalStakedBalance is the total staked balance.
	TotalStakedBalance string `json:"total_staked_balance"`
	// RewardFeeFraction is the reward fee fraction.
	RewardFeeFraction float32 `json:"reward_fee_fraction"`
	// StakingKey is the staking public key.
	StakingKey string `json:"staking_key"`
	// IsStakingPaused indicates whether staking is paused.
	IsStakingPaused bool `json:"is_staking_paused"`
}

// XTZStakingReward represents XTZ staking rewards information.
type XTZStakingReward struct {
	// ReceivedRewardsAmount is the total amount of staking rewards received.
	ReceivedRewardsAmount string `json:"received_rewards_amount"`
}

// ListStakeAccountsOptions contains options for listing stake accounts.
type ListStakeAccountsOptions struct {
	// AddressID filters by the Protect address ID.
	AddressID string
	// AccountType filters by stake account type.
	AccountType string
	// AccountAddress filters by stake account address.
	AccountAddress string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// ListStakeAccountsResult contains the result of listing stake accounts.
type ListStakeAccountsResult struct {
	// StakeAccounts is the list of stake accounts.
	StakeAccounts []*StakeAccount `json:"stake_accounts"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// GetETHValidatorsInfoOptions contains options for getting ETH validators info.
type GetETHValidatorsInfoOptions struct {
	// IDs is the list of validator IDs to retrieve (max 500).
	IDs []string
}

// GetXTZStakingRewardsOptions contains options for getting XTZ staking rewards.
type GetXTZStakingRewardsOptions struct {
	// From filters rewards from this date.
	From *time.Time
	// To filters rewards until this date.
	To *time.Time
}
