package model

import (
	"crypto/ecdsa"
	"sync"
)

// DecodedRulesContainer represents the decoded rules container with all governance rules.
// It contains users, groups, signature thresholds, and rules for transactions and address whitelisting.
type DecodedRulesContainer struct {
	// Users defined in the governance rules.
	Users []*RuleUser
	// Groups defined in the governance rules.
	Groups []*RuleGroup
	// MinimumDistinctUserSignatures required for rules container updates.
	MinimumDistinctUserSignatures int
	// MinimumDistinctGroupSignatures required for rules container updates.
	MinimumDistinctGroupSignatures int
	// TransactionRules organized by key (blockchain/action type).
	TransactionRules []*TransactionRules
	// AddressWhitelistingRules organized by blockchain and network.
	AddressWhitelistingRules []*AddressWhitelistingRules
	// ContractAddressWhitelistingRules organized by blockchain and network.
	ContractAddressWhitelistingRules []*ContractAddressWhitelistingRules
	// EnforcedRulesHash is the SHA-256 hash of the enforced rules.
	EnforcedRulesHash string
	// Timestamp when the rules container was created or updated.
	Timestamp int64
	// MinimumCommitmentSignatures required from HSM engines.
	MinimumCommitmentSignatures int
	// EngineIdentities are the HSM serial numbers authorized for this tenant.
	EngineIdentities []string
	// HsmSlotId is the HSM slot ID that these rules are intended for.
	HsmSlotId uint32

	// hsmPublicKey is the cached HSM public key (lazily resolved).
	hsmPublicKey *ecdsa.PublicKey
	// hsmKeyOnce ensures thread-safe lazy initialization of the HSM public key.
	hsmKeyOnce sync.Once
}

// RuleUser represents a user defined in the governance rules container.
type RuleUser struct {
	// ID is the unique identifier for the user.
	ID string
	// PublicKeyPEM is the PEM-encoded public key for signature verification.
	PublicKeyPEM string
	// PublicKey is the decoded ECDSA public key.
	PublicKey *ecdsa.PublicKey
	// Roles assigned to this user (e.g., "SUPERADMIN", "HSMSLOT").
	Roles []string
}

// RuleGroup represents a group of users defined in the governance rules container.
type RuleGroup struct {
	// ID is the unique identifier for the group.
	ID string
	// UserIDs are the IDs of users in this group.
	UserIDs []string
}

// TransactionRules represents transaction approval rules.
type TransactionRules struct {
	// Key is the rule key (e.g., blockchain/action type).
	Key string
	// Columns define the rule structure.
	Columns []*RuleColumn
	// Lines define the rule rows.
	Lines []*RuleLine
	// Details contain additional rule configuration.
	Details *TransactionRuleDetails
}

// RuleColumn represents a column in transaction rules.
type RuleColumn struct {
	// Type is the column type (e.g., "AMOUNT", "SOURCE").
	Type string
}

// RuleLine represents a line/row in transaction rules.
type RuleLine struct {
	// Cells contain the cell values.
	Cells []string
	// ParallelThresholds define approval requirements for this line.
	ParallelThresholds []*SequentialThresholds
}

// TransactionRuleDetails contains additional transaction rule configuration.
type TransactionRuleDetails struct {
	// Domain is the rule domain.
	Domain string
	// SubDomain is the rule sub-domain.
	SubDomain string
}

// AddressWhitelistingRules represents address whitelisting rules for a blockchain/network.
type AddressWhitelistingRules struct {
	// Currency is the blockchain identifier (e.g., "ETH", "BTC"). Empty means global default.
	Currency string
	// Network is the network identifier (e.g., "mainnet", "testnet"). Empty matches any network.
	Network string
	// ParallelThresholds define default approval requirements.
	ParallelThresholds []*SequentialThresholds
	// Lines contain source-specific rule overrides.
	Lines []*AddressWhitelistingLine
	// IncludeNetworkInPayload indicates whether to include network in the payload hash.
	IncludeNetworkInPayload bool
}

// AddressWhitelistingLine represents a source-specific rule line.
type AddressWhitelistingLine struct {
	// Cells contain the rule source specifications.
	Cells []*RuleSource
	// ParallelThresholds define approval requirements for this line.
	ParallelThresholds []*SequentialThresholds
}

// RuleSource represents a source specification in a whitelist rule.
type RuleSource struct {
	// Type is the source type.
	Type RuleSourceType
	// InternalWallet is set when Type is RuleSourceInternalWallet.
	InternalWallet *RuleSourceInternalWallet
}

// RuleSourceType represents the type of rule source.
type RuleSourceType int

const (
	// RuleSourceTypeUnknown is the default/unknown type.
	RuleSourceTypeUnknown RuleSourceType = iota
	// RuleSourceTypeInternalWallet indicates an internal wallet source.
	RuleSourceTypeInternalWallet
)

// RuleSourceInternalWallet represents an internal wallet rule source.
type RuleSourceInternalWallet struct {
	// Path is the wallet path.
	Path string
}

// ContractAddressWhitelistingRules represents contract address whitelisting rules.
type ContractAddressWhitelistingRules struct {
	// Blockchain is the blockchain identifier.
	Blockchain string
	// Network is the network identifier.
	Network string
	// ParallelThresholds define approval requirements.
	ParallelThresholds []*SequentialThresholds
}

// SequentialThresholds represents a sequence of group thresholds that must be satisfied in order.
type SequentialThresholds struct {
	// Thresholds are the ordered group thresholds to satisfy.
	Thresholds []*GroupThreshold
}

// GroupThreshold represents the signature threshold for a specific group.
type GroupThreshold struct {
	// GroupID is the ID of the group this threshold applies to.
	GroupID string
	// MinimumSignatures is the minimum number of signatures required from this group.
	MinimumSignatures int
}

// anyWildcard is the string that represents "any" in rules matching.
const anyWildcard = "Any"

// FindAddressWhitelistingRules finds AddressWhitelistingRules matching the given blockchain and network.
// Uses a three-tier priority system:
// 1. Exact match - both blockchain and network match exactly
// 2. Blockchain-only match - blockchain matches, rule has wildcard network
// 3. Global default - rule has wildcard blockchain
func (r *DecodedRulesContainer) FindAddressWhitelistingRules(blockchain, network string) *AddressWhitelistingRules {
	if r.AddressWhitelistingRules == nil {
		return nil
	}

	var blockchainOnlyMatch *AddressWhitelistingRules
	var globalDefault *AddressWhitelistingRules

	for _, rule := range r.AddressWhitelistingRules {
		ruleIsGlobalDefault := isWildcard(rule.Currency)
		blockchainMatches := !ruleIsGlobalDefault && rule.Currency == blockchain
		networkMatches := rule.Network == network
		ruleHasWildcardNetwork := isWildcard(rule.Network)

		// Priority 1: Exact match (blockchain + network)
		if blockchainMatches && networkMatches {
			return rule
		}

		// Priority 2: Blockchain match with wildcard network
		if blockchainMatches && ruleHasWildcardNetwork && blockchainOnlyMatch == nil {
			blockchainOnlyMatch = rule
		}

		// Priority 3: Global default (wildcard blockchain)
		if ruleIsGlobalDefault && globalDefault == nil {
			globalDefault = rule
		}
	}

	// Return best match by priority
	if blockchainOnlyMatch != nil {
		return blockchainOnlyMatch
	}
	return globalDefault
}

// FindContractAddressWhitelistingRules finds ContractAddressWhitelistingRules matching the given blockchain and network.
func (r *DecodedRulesContainer) FindContractAddressWhitelistingRules(blockchain, network string) *ContractAddressWhitelistingRules {
	if r.ContractAddressWhitelistingRules == nil {
		return nil
	}

	var blockchainOnlyMatch *ContractAddressWhitelistingRules
	var globalDefault *ContractAddressWhitelistingRules

	for _, rule := range r.ContractAddressWhitelistingRules {
		ruleIsGlobalDefault := isWildcard(rule.Blockchain)
		blockchainMatches := !ruleIsGlobalDefault && rule.Blockchain == blockchain
		networkMatches := rule.Network == network
		ruleHasWildcardNetwork := isWildcard(rule.Network)

		// Priority 1: Exact match (blockchain + network)
		if blockchainMatches && networkMatches {
			return rule
		}

		// Priority 2: Blockchain match with wildcard network
		if blockchainMatches && ruleHasWildcardNetwork && blockchainOnlyMatch == nil {
			blockchainOnlyMatch = rule
		}

		// Priority 3: Global default (wildcard blockchain)
		if ruleIsGlobalDefault && globalDefault == nil {
			globalDefault = rule
		}
	}

	if blockchainOnlyMatch != nil {
		return blockchainOnlyMatch
	}
	return globalDefault
}

// FindUserByID finds a RuleUser by ID.
func (r *DecodedRulesContainer) FindUserByID(userID string) *RuleUser {
	if r.Users == nil || userID == "" {
		return nil
	}
	for _, user := range r.Users {
		if user.ID == userID {
			return user
		}
	}
	return nil
}

// FindGroupByID finds a RuleGroup by ID.
func (r *DecodedRulesContainer) FindGroupByID(groupID string) *RuleGroup {
	if r.Groups == nil || groupID == "" {
		return nil
	}
	for _, group := range r.Groups {
		if group.ID == groupID {
			return group
		}
	}
	return nil
}

// GetHsmPublicKey returns the HSM slot public key (cached).
// Finds the first user with the HSMSLOT role.
// Thread-safe: uses sync.Once for lazy initialization.
func (r *DecodedRulesContainer) GetHsmPublicKey() *ecdsa.PublicKey {
	r.hsmKeyOnce.Do(func() {
		r.hsmPublicKey = r.findHsmPublicKey()
	})
	return r.hsmPublicKey
}

func (r *DecodedRulesContainer) findHsmPublicKey() *ecdsa.PublicKey {
	if r.Users == nil || len(r.Users) == 0 {
		return nil
	}
	for _, user := range r.Users {
		if user.Roles != nil && containsRole(user.Roles, "HSMSLOT") {
			if user.PublicKey != nil {
				return user.PublicKey
			}
		}
	}
	return nil
}

// isWildcard checks if a value represents a wildcard (empty or "Any").
func isWildcard(value string) bool {
	return value == "" || value == anyWildcard
}

// containsRole checks if a role list contains a specific role.
func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasRole checks if the user has the specified role.
func (u *RuleUser) HasRole(role string) bool {
	return containsRole(u.Roles, role)
}

// ContainsUser checks if the group contains the specified user.
func (g *RuleGroup) ContainsUser(userID string) bool {
	for _, id := range g.UserIDs {
		if id == userID {
			return true
		}
	}
	return false
}
