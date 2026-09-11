package model

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"strings"
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
	// Properties holds container-level opaque properties (proto map<string, bytes>).
	Properties map[string][]byte
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte

	// hsmPublicKey is the cached HSM public key (lazily resolved).
	hsmPublicKey *ecdsa.PublicKey
	// hsmKeyOnce ensures thread-safe lazy initialization of the HSM public key.
	hsmKeyOnce sync.Once
	// priceUpdaterKeys is the cached set of PRICEUPDATER public keys.
	priceUpdaterKeys []*ecdsa.PublicKey
	// priceUpdaterOnce guards lazy initialization of priceUpdaterKeys.
	priceUpdaterOnce sync.Once
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
	// Properties holds user-level opaque properties (proto map<string, bytes>).
	Properties map[string][]byte
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte

	// fingerprint is the cached KeyFingerprint result.
	fingerprint string
	// fingerprintErr is the error from computing fingerprint, if any.
	fingerprintErr error
	// fingerprintOnce guards lazy fingerprint computation.
	fingerprintOnce sync.Once
}

// KeyFingerprint identifies this user by its public key, not by its server-supplied ID.
//
// Threshold counting must key on the key: one compromised key shared by two user IDs has
// to count once, or an N-of-M quorum is satisfiable by a single secret. Cached because a
// page of rows re-walks the same users for every row and group.
func (u *RuleUser) KeyFingerprint() (string, error) {
	u.fingerprintOnce.Do(func() {
		u.fingerprint, u.fingerprintErr = KeyFingerprint(u.PublicKey)
	})
	return u.fingerprint, u.fingerprintErr
}

// KeyFingerprint identifies a key by its encoded bytes, so the same key configured twice
// counts as one signer.
func KeyFingerprint(key *ecdsa.PublicKey) (string, error) {
	if key == nil {
		return "", fmt.Errorf("public key cannot be nil")
	}
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(der)
	return string(sum[:]), nil
}

// RuleGroup represents a group of users defined in the governance rules container.
type RuleGroup struct {
	// ID is the unique identifier for the group.
	ID string
	// UserIDs are the IDs of users in this group.
	UserIDs []string
	// Properties holds group-level opaque properties (proto map<string, bytes>).
	Properties map[string][]byte
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
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
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
}

// RuleColumn represents a column in transaction rules.
type RuleColumn struct {
	// Type is the column type enum value name (e.g., "RuleSource", "RuleFiatAmount").
	Type string
	// Name is the human-readable column name.
	Name string
	// MetadataKey is the request-metadata key this column matches against.
	MetadataKey string
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
}

// RuleLine represents a line/row in transaction rules.
type RuleLine struct {
	// Cells contain the typed cell values, aligned positionally with the rule's
	// Columns. A nil cell means "match any".
	Cells []RuleCell
	// ParallelThresholds define approval requirements for this line.
	ParallelThresholds []*SequentialThresholds
	// Priority orders lines during rule evaluation.
	Priority uint32
	// Properties holds line-level opaque properties (proto map<string, bytes>).
	Properties map[string][]byte
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
}

// TransactionRuleDetails contains additional transaction rule configuration.
type TransactionRuleDetails struct {
	// Domain is the rule domain.
	Domain string
	// SubDomain is the rule sub-domain.
	SubDomain string
	// Blockchain scopes the rule to a blockchain (empty means unscoped).
	Blockchain string
	// Network scopes the rule to a network (empty means unscoped).
	Network string
	// EvmCallContract holds EVM contract-call scoping, if any.
	EvmCallContract *EvmCallContract
	// XtzCallContract holds Tezos contract-call scoping, if any.
	XtzCallContract *XtzCallContract
	// CashSettlement holds cash-settlement scoping, if any.
	CashSettlement *CashSettlement
	// CosmosDetails holds Cosmos scoping, if any.
	CosmosDetails *CosmosDetails
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
}

// EvmCallContract scopes a transaction rule to an EVM contract call.
type EvmCallContract struct {
	ContractType    string
	MethodSignature string
	// UnknownFields holds protobuf fields from a newer schema, re-emitted on encode.
	UnknownFields []byte
}

// XtzCallContract scopes a transaction rule to a Tezos contract call.
type XtzCallContract struct {
	ContractType    string
	MethodSignature string
	// UnknownFields holds protobuf fields from a newer schema, re-emitted on encode.
	UnknownFields []byte
}

// CashSettlement scopes a transaction rule to a cash-settlement flow.
type CashSettlement struct {
	Provider    string
	RequestType string
	// UnknownFields holds protobuf fields from a newer schema, re-emitted on encode.
	UnknownFields []byte
}

// CosmosDetails scopes a transaction rule to Cosmos method signatures.
type CosmosDetails struct {
	MethodSignatures []string
	// UnknownFields holds protobuf fields from a newer schema, re-emitted on encode.
	UnknownFields []byte
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
	// Properties holds rule-level opaque properties (proto map<string, bytes>).
	Properties map[string][]byte
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
}

// AddressWhitelistingLine represents a source-specific rule line.
type AddressWhitelistingLine struct {
	// Cells contain the rule source specifications.
	Cells []*RuleSource
	// ParallelThresholds define approval requirements for this line.
	ParallelThresholds []*SequentialThresholds
	// Properties holds line-level opaque properties (proto map<string, bytes>).
	Properties map[string][]byte
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
}

// RuleSource represents a source specification in a whitelist rule.
type RuleSource struct {
	// Type is the source type.
	Type RuleSourceType
	// InternalWallet is set when Type is RuleSourceTypeInternalWallet.
	InternalWallet *RuleSourceInternalWallet
	// InternalAddress is set when Type is RuleSourceTypeInternalAddress.
	InternalAddress *RuleSourceInternalAddress
	// Exchange is set when Type is RuleSourceTypeExchange.
	Exchange *RuleSourceExchange
	// ExternalAddress is set when Type is RuleSourceTypeExternalAddress.
	ExternalAddress *RuleSourceExternalAddress
	// Raw preserves a source cell whose type or content is unknown to this SDK
	// version; it is re-emitted verbatim on encode. When Raw is set the typed
	// fields above are not populated. Do not modify.
	Raw []byte
}

// RuleSourceType represents the type of rule source.
type RuleSourceType int

const (
	// RuleSourceTypeUnknown is the default/unknown type (wire value of RuleSourceAny).
	RuleSourceTypeUnknown RuleSourceType = iota
	// RuleSourceTypeInternalWallet indicates an internal wallet source.
	RuleSourceTypeInternalWallet
	// RuleSourceTypeInternalAddress indicates an internal address source.
	RuleSourceTypeInternalAddress
	// RuleSourceTypeAnyExchange indicates any exchange source.
	RuleSourceTypeAnyExchange
	// RuleSourceTypeExchange indicates a specific exchange source.
	RuleSourceTypeExchange
	// RuleSourceTypeExternalAddress indicates an external address source.
	RuleSourceTypeExternalAddress
)

// RuleSourceTypeAny matches any source (alias of the wire zero value).
const RuleSourceTypeAny = RuleSourceTypeUnknown

// RuleSourceInternalWallet represents an internal wallet rule source.
type RuleSourceInternalWallet struct {
	// Path is the wallet path.
	Path string
}

// RuleSourceInternalAddress represents an internal address rule source.
type RuleSourceInternalAddress struct {
	Address string
	Path    string
}

// RuleSourceExchange represents an exchange rule source.
type RuleSourceExchange struct {
	Label string
}

// RuleSourceExternalAddress represents an external address rule source.
type RuleSourceExternalAddress struct {
	Address string
	Memo    string
}

// ContractAddressWhitelistingRules represents contract address whitelisting rules.
type ContractAddressWhitelistingRules struct {
	// Blockchain is the blockchain identifier.
	Blockchain string
	// Network is the network identifier.
	Network string
	// ParallelThresholds define approval requirements.
	ParallelThresholds []*SequentialThresholds
	// Properties holds rule-level opaque properties (proto map<string, bytes>).
	Properties map[string][]byte
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
}

// SequentialThresholds represents a sequence of group thresholds that must be satisfied in order.
type SequentialThresholds struct {
	// Thresholds are the ordered group thresholds to satisfy.
	Thresholds []*GroupThreshold
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
}

// GroupThreshold represents the signature threshold for a specific group.
type GroupThreshold struct {
	// GroupID is the ID of the group this threshold applies to.
	GroupID string
	// MinimumSignatures is the minimum number of signatures required from this group.
	MinimumSignatures int
	// UnknownFields holds protobuf fields unknown to this SDK version, captured
	// at decode and re-emitted verbatim on encode. Do not modify.
	UnknownFields []byte
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

// FindAddressWhitelistingRuleCandidates returns every rule the tier walk could select for this
// blockchain across all possible network values.
//
// This exists because the network half of the rule key is NOT always signed. Governance carries a
// per-rule includeNetworkInPayload flag, and when it is off the signed payload has no `network`
// member at all — the common case in captured production data. The key then falls back to the
// network on the unsigned response DTO, which hands a response-controlling server the choice of
// WHICH rule judges the row, and therefore which group quorum it must meet. That is not closeable
// by reading the flag: includeNetworkInPayload has no proto backing in any of the four SDKs
// (request_reply.proto's AddressWhitelistingRules carries only currency, parallelThresholds,
// properties, network, lines), so it is never part of the SuperAdmin-signed container.
//
// So when the network is unsigned the caller must satisfy EVERY rule this returns, not the one the
// DTO named. Where a chain has a single reachable tier — again the common case — the set has one
// element and behaviour is unchanged.
//
// Reachability, mirroring the tier walk above: every rule for this chain is reachable by naming
// its network; the chain's wildcard-network rule is reachable by naming a network no exact rule
// covers; and the global default is reachable ONLY when the chain has no wildcard-network rule,
// because priority 2 would otherwise win.
func (r *DecodedRulesContainer) FindAddressWhitelistingRuleCandidates(
	blockchain string,
) []*AddressWhitelistingRules {
	if r.AddressWhitelistingRules == nil {
		return nil
	}

	var candidates []*AddressWhitelistingRules
	var globalDefault *AddressWhitelistingRules
	chainHasWildcardNetwork := false

	for _, rule := range r.AddressWhitelistingRules {
		if isWildcard(rule.Currency) {
			if globalDefault == nil {
				globalDefault = rule
			}
			continue
		}
		if rule.Currency != blockchain {
			continue
		}
		candidates = append(candidates, rule)
		if isWildcard(rule.Network) {
			chainHasWildcardNetwork = true
		}
	}

	if !chainHasWildcardNetwork && globalDefault != nil {
		candidates = append(candidates, globalDefault)
	}
	return candidates
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

// FindContractAddressWhitelistingRuleCandidates is the asset peer of
// FindAddressWhitelistingRuleCandidates. Same reasoning: when the signed payload omits
// `network`, the unsigned response DTO would otherwise choose which quorum judges the asset.
func (r *DecodedRulesContainer) FindContractAddressWhitelistingRuleCandidates(
	blockchain string,
) []*ContractAddressWhitelistingRules {
	if r.ContractAddressWhitelistingRules == nil {
		return nil
	}

	var candidates []*ContractAddressWhitelistingRules
	var globalDefault *ContractAddressWhitelistingRules
	chainHasWildcardNetwork := false

	for _, rule := range r.ContractAddressWhitelistingRules {
		if isWildcard(rule.Blockchain) {
			if globalDefault == nil {
				globalDefault = rule
			}
			continue
		}
		if rule.Blockchain != blockchain {
			continue
		}
		candidates = append(candidates, rule)
		if isWildcard(rule.Network) {
			chainHasWildcardNetwork = true
		}
	}

	if !chainHasWildcardNetwork && globalDefault != nil {
		candidates = append(candidates, globalDefault)
	}
	return candidates
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
	if len(r.Users) == 0 {
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

// PriceUpdaterKeys returns the public keys of every user carrying the PRICEUPDATER
// role (cached).
//
// Prices are signed by these keys. The container itself is SuperAdmin-verified, so
// "this tenant has no price signer" is a trustworthy statement — whereas "this price
// carries no signatures" is not, which is what makes the empty result meaningful.
func (r *DecodedRulesContainer) PriceUpdaterKeys() []*ecdsa.PublicKey {
	r.priceUpdaterOnce.Do(func() {
		for _, user := range r.Users {
			if user == nil || user.PublicKey == nil {
				continue
			}
			if containsRole(user.Roles, "PRICEUPDATER") {
				r.priceUpdaterKeys = append(r.priceUpdaterKeys, user.PublicKey)
			}
		}
	})
	return r.priceUpdaterKeys
}

// HasUnknownFields reports whether any part of the container carries data
// unknown to this SDK version (unknown protobuf fields, raw cells, or raw
// whitelisting sources). Such data is preserved verbatim on re-encode; a true
// result is a hint to upgrade the SDK before editing the container.
func (r *DecodedRulesContainer) HasUnknownFields() bool {
	if len(r.UnknownFields) > 0 {
		return true
	}
	for _, u := range r.Users {
		if u != nil && len(u.UnknownFields) > 0 {
			return true
		}
	}
	for _, g := range r.Groups {
		if g != nil && len(g.UnknownFields) > 0 {
			return true
		}
	}
	for _, tr := range r.TransactionRules {
		if transactionRulesHaveUnknown(tr) {
			return true
		}
	}
	for _, awr := range r.AddressWhitelistingRules {
		if addressWhitelistingRulesHaveUnknown(awr) {
			return true
		}
	}
	for _, cawr := range r.ContractAddressWhitelistingRules {
		if cawr == nil {
			continue
		}
		if len(cawr.UnknownFields) > 0 || thresholdsHaveUnknown(cawr.ParallelThresholds) {
			return true
		}
	}
	return false
}

func transactionRulesHaveUnknown(tr *TransactionRules) bool {
	if tr == nil {
		return false
	}
	if len(tr.UnknownFields) > 0 {
		return true
	}
	for _, c := range tr.Columns {
		if c != nil && len(c.UnknownFields) > 0 {
			return true
		}
	}
	for _, l := range tr.Lines {
		if l == nil {
			continue
		}
		if len(l.UnknownFields) > 0 || thresholdsHaveUnknown(l.ParallelThresholds) {
			return true
		}
		for _, cell := range l.Cells {
			if _, ok := cell.(RawCell); ok {
				return true
			}
		}
	}
	return ruleDetailsHaveUnknown(tr.Details)
}

// ruleDetailsHaveUnknown reports unknown fields on the details node and on each of its
// nested contract-call scoping sub-messages. Checking only the details node would report
// a container clean while a nested node carried data from a newer schema.
func ruleDetailsHaveUnknown(d *TransactionRuleDetails) bool {
	if d == nil {
		return false
	}
	if len(d.UnknownFields) > 0 {
		return true
	}
	if d.EvmCallContract != nil && len(d.EvmCallContract.UnknownFields) > 0 {
		return true
	}
	if d.XtzCallContract != nil && len(d.XtzCallContract.UnknownFields) > 0 {
		return true
	}
	if d.CashSettlement != nil && len(d.CashSettlement.UnknownFields) > 0 {
		return true
	}
	if d.CosmosDetails != nil && len(d.CosmosDetails.UnknownFields) > 0 {
		return true
	}
	return false
}

func addressWhitelistingRulesHaveUnknown(awr *AddressWhitelistingRules) bool {
	if awr == nil {
		return false
	}
	if len(awr.UnknownFields) > 0 || thresholdsHaveUnknown(awr.ParallelThresholds) {
		return true
	}
	for _, l := range awr.Lines {
		if l == nil {
			continue
		}
		if len(l.UnknownFields) > 0 || thresholdsHaveUnknown(l.ParallelThresholds) {
			return true
		}
		for _, s := range l.Cells {
			if s != nil && len(s.Raw) > 0 {
				return true
			}
		}
	}
	return false
}

func thresholdsHaveUnknown(thresholds []*SequentialThresholds) bool {
	for _, st := range thresholds {
		if st == nil {
			continue
		}
		if len(st.UnknownFields) > 0 {
			return true
		}
		for _, t := range st.Thresholds {
			if t != nil && len(t.UnknownFields) > 0 {
				return true
			}
		}
	}
	return false
}

// isWildcard checks if a value represents a wildcard (empty or "Any").
//
// Case-insensitive, matching Java (equalsIgnoreCase), Python (.lower()) and
// TypeScript (.toLowerCase()). Comparing exactly against "Any" made a container
// carrying "any" or "ANY" select a different rule tier in this SDK than in the
// other three — a literal blockchain match rather than the global default.
func isWildcard(value string) bool {
	return value == "" || strings.EqualFold(value, anyWildcard)
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
