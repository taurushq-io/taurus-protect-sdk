package model

import "time"

// WhitelistedAsset represents a whitelisted contract address (token/NFT).
// In the API, these are referred to as "whitelisted contracts" but the SDK
// exposes them as "whitelisted assets" for clarity.
type WhitelistedAsset struct {
	// ID is the unique identifier for the whitelisted asset.
	ID string `json:"id"`
	// TenantID is the tenant identifier.
	TenantID string `json:"tenant_id,omitempty"`
	// Status is the approval status (e.g., "approved", "pending").
	Status string `json:"status,omitempty"`
	// Action is the action type.
	Action string `json:"action,omitempty"`
	// Blockchain is the blockchain symbol (e.g., "ETH", "XTZ").
	Blockchain string `json:"blockchain,omitempty"`
	// Network is the network type (e.g., "mainnet", "testnet").
	Network string `json:"network,omitempty"`
	// Rule is the governance rule applied.
	Rule string `json:"rule,omitempty"`
	// RulesContainer is the rules container ID.
	RulesContainer string `json:"rules_container,omitempty"`
	// RulesSignatures contains the super-admin signatures for verification.
	RulesSignatures string `json:"rules_signatures,omitempty"`
	// BusinessRuleEnabled indicates if the currency is enabled in business rules.
	BusinessRuleEnabled bool `json:"business_rule_enabled"`
	// Metadata contains additional metadata about the asset.
	Metadata *WhitelistedAssetMetadata `json:"metadata,omitempty"`
	// SignedContractAddress contains the signed payload and signatures.
	SignedContractAddress *SignedContractAddress `json:"signed_contract_address,omitempty"`
	// Approvers contains approval information.
	Approvers *Approvers `json:"approvers,omitempty"`
	// Attributes are custom key-value attributes.
	Attributes []WhitelistedAssetAttribute `json:"attributes,omitempty"`
	// Trails contains audit trail information.
	Trails []Trail `json:"trails,omitempty"`
}

// WhitelistedAssetMetadata contains metadata about a whitelisted asset.
type WhitelistedAssetMetadata struct {
	// Hash is the hash of the payload.
	Hash string `json:"hash,omitempty"`
	// SECURITY: Payload field intentionally omitted - use PayloadAsString only.
	// The raw payload object could be tampered with by an attacker while
	// payloadAsString remains unchanged (hash still verifies). By not having
	// this field, we enforce that all data extraction uses the verified source.
	// PayloadAsString is the payload serialized as a string.
	PayloadAsString string `json:"payload_as_string,omitempty"`
}

// SignedContractAddress contains the signed contract address data.
type SignedContractAddress struct {
	// Payload is the base64-encoded signed payload.
	Payload string `json:"payload,omitempty"`
	// Signatures are the cryptographic signatures.
	Signatures []WhitelistSignature `json:"signatures,omitempty"`
}

// WhitelistSignature represents a signature on a whitelist entry.
type WhitelistSignature struct {
	// UserSignature contains the user's signature details.
	UserSignature *WhitelistUserSignature `json:"user_signature,omitempty"`
	// Hashes are the hashes covered by this signature.
	Hashes []string `json:"hashes,omitempty"`
}

// WhitelistUserSignature represents a user's signature on a whitelist entry.
type WhitelistUserSignature struct {
	// UserID is the ID of the signing user.
	UserID string `json:"user_id,omitempty"`
	// Signature is the base64-encoded cryptographic signature.
	Signature string `json:"signature,omitempty"`
	// Comment is an optional comment from the signer.
	Comment string `json:"comment,omitempty"`
}

// Approvers contains approval information for a whitelisted asset.
type Approvers struct {
	// Parallel contains parallel approval groups.
	Parallel []ParallelApproversGroup `json:"parallel,omitempty"`
}

// ParallelApproversGroup contains sequential approval groups that can run in parallel.
type ParallelApproversGroup struct {
	// Sequential contains the sequential approval groups.
	Sequential []ApproversGroup `json:"sequential,omitempty"`
}

// ApproversGroup represents a group of approvers with requirements.
type ApproversGroup struct {
	// ExternalGroupID is the external identifier for the group.
	ExternalGroupID string `json:"external_group_id,omitempty"`
	// MinimumSignatures is the minimum number of signatures required from this group.
	MinimumSignatures int64 `json:"minimum_signatures,omitempty"`
}

// Trail represents an audit trail entry.
type Trail struct {
	// ID is the trail entry identifier.
	ID string `json:"id,omitempty"`
	// UserID is the ID of the user who performed the action.
	UserID string `json:"user_id,omitempty"`
	// ExternalUserID is the external user identifier.
	ExternalUserID string `json:"external_user_id,omitempty"`
	// Action is the action performed.
	Action string `json:"action,omitempty"`
	// Comment is an optional comment.
	Comment string `json:"comment,omitempty"`
	// Date is when the action occurred.
	Date time.Time `json:"date,omitempty"`
}

// WhitelistedAssetAttribute represents a custom attribute on a whitelisted asset.
type WhitelistedAssetAttribute struct {
	// ID is the attribute identifier.
	ID string `json:"id"`
	// Key is the attribute name.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
	// ContentType is the content type of the value.
	ContentType string `json:"content_type,omitempty"`
	// Owner is the owner of the attribute.
	Owner string `json:"owner,omitempty"`
	// Type is the attribute type.
	Type string `json:"type,omitempty"`
	// Subtype is the attribute subtype.
	Subtype string `json:"subtype,omitempty"`
	// IsFile indicates if the attribute value is a file.
	IsFile bool `json:"is_file,omitempty"`
}

// WhitelistedAssetEnvelope wraps a whitelisted asset with its verification data.
// This is returned by GetWhitelistedAssetEnvelope after the 5-step verification flow.
type WhitelistedAssetEnvelope struct {
	// ID is the unique identifier for the whitelisted asset.
	ID string `json:"id"`
	// TenantID is the tenant identifier.
	TenantID string `json:"tenant_id,omitempty"`
	// Blockchain is the blockchain network (e.g., "ETH", "XTZ").
	Blockchain string `json:"blockchain"`
	// Network is the network type (e.g., "mainnet", "testnet").
	Network string `json:"network"`
	// Status is the approval status (e.g., "approved", "pending").
	Status string `json:"status,omitempty"`
	// Action is the action type.
	Action string `json:"action,omitempty"`
	// Rule is the governance rule applied.
	Rule string `json:"rule,omitempty"`
	// BusinessRuleEnabled indicates if the currency is enabled in business rules.
	BusinessRuleEnabled bool `json:"business_rule_enabled"`
	// Metadata contains the asset metadata.
	Metadata *WhitelistedAssetMetadata `json:"metadata,omitempty"`
	// SignedContractAddress contains the signed payload and signatures.
	SignedContractAddress *SignedContractAddress `json:"signed_contract_address,omitempty"`
	// RulesContainer is the base64-encoded rules container.
	RulesContainer string `json:"rules_container,omitempty"`
	// RulesSignatures is the base64-encoded rules signatures.
	RulesSignatures string `json:"rules_signatures,omitempty"`
	// Approvers contains approval information.
	Approvers *Approvers `json:"approvers,omitempty"`
	// Attributes are custom key-value attributes.
	Attributes []WhitelistedAssetAttribute `json:"attributes,omitempty"`
	// Trails contains audit trail information.
	Trails []Trail `json:"trails,omitempty"`

	// verifiedWhitelistedAsset is the verified whitelisted asset parsed from the payload.
	verifiedWhitelistedAsset *WhitelistedAsset
	// verifiedRulesContainer is the verified and decoded rules container.
	verifiedRulesContainer *DecodedRulesContainer
}

// WhitelistedAsset returns the verified whitelisted asset.
// This is populated after successful verification.
func (e *WhitelistedAssetEnvelope) WhitelistedAsset() *WhitelistedAsset {
	return e.verifiedWhitelistedAsset
}

// DecodedRulesContainer returns the verified and decoded rules container.
// This is populated after successful verification.
func (e *WhitelistedAssetEnvelope) DecodedRulesContainer() *DecodedRulesContainer {
	return e.verifiedRulesContainer
}

// SetVerified sets the verified asset and rules container.
// This is called internally after successful verification.
func (e *WhitelistedAssetEnvelope) SetVerified(asset *WhitelistedAsset, rules *DecodedRulesContainer) {
	e.verifiedWhitelistedAsset = asset
	e.verifiedRulesContainer = rules
}

// ListWhitelistedAssetsOptions contains options for listing whitelisted assets.
type ListWhitelistedAssetsOptions struct {
	// Limit is the maximum number of assets to return (max 100).
	Limit int64
	// Offset is the number of assets to skip.
	Offset int64
	// Query searches across multiple fields (address, blockchain, label, etc.).
	Query string
	// Blockchain filters by blockchain symbol.
	Blockchain string
	// Network filters by network (e.g., "mainnet", "testnet").
	Network string
	// IncludeForApproval includes assets pending approval in results.
	IncludeForApproval bool
	// KindTypes filters by contract kind type ("nft" or "token").
	KindTypes []string
	// IDs filters by specific whitelisted asset IDs.
	IDs []string
}
