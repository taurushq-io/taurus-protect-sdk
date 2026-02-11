package model

import "time"

// WhitelistedAddress represents a whitelisted address in the system.
type WhitelistedAddress struct {
	// ID is the unique identifier for the whitelisted address.
	ID string `json:"id"`
	// TenantID is the tenant identifier.
	TenantID string `json:"tenant_id"`
	// Address is the blockchain address string (extracted from metadata).
	Address string `json:"address"`
	// Label is the human-readable label for this address (extracted from metadata).
	Label string `json:"label"`
	// Memo is the optional memo field for Stellar or destination tag for Ripple.
	Memo string `json:"memo,omitempty"`
	// CustomerId is the customer ID for external reconciliation.
	CustomerId string `json:"customer_id,omitempty"`
	// Blockchain is the blockchain network (e.g., "ETH", "BTC").
	Blockchain string `json:"blockchain"`
	// Network is the network type (e.g., "mainnet", "testnet").
	Network string `json:"network"`
	// Status is the current status of the whitelisted address.
	Status string `json:"status"`
	// Action is the pending action type if any.
	Action string `json:"action,omitempty"`
	// AddressType is the type of address (individual, exchange, contract, etc.).
	AddressType string `json:"address_type,omitempty"`
	// ContractType is the smart contract type (e.g., "CMTA20", "ERC20") for contract addresses.
	ContractType string `json:"contract_type,omitempty"`
	// ExchangeAccountId is the exchange account ID when the address belongs to an exchange.
	ExchangeAccountId int64 `json:"exchange_account_id,omitempty"`
	// Rule is the governance rule applied to this address.
	Rule string `json:"rule,omitempty"`
	// RulesContainer is the serialized rules container.
	RulesContainer string `json:"rules_container,omitempty"`
	// RulesContainerHash is the hash of the rules container.
	RulesContainerHash string `json:"rules_container_hash,omitempty"`
	// RulesSignatures contains the super-admin signatures for the rules.
	RulesSignatures string `json:"rules_signatures,omitempty"`
	// VisibilityGroupID is the visibility group this address belongs to.
	VisibilityGroupID string `json:"visibility_group_id,omitempty"`
	// TnParticipantID is the Taurus Network participant ID.
	TnParticipantID string `json:"tn_participant_id,omitempty"`
	// Metadata contains additional metadata about the address.
	// Uses WhitelistedAssetMetadata which is shared between assets and addresses.
	Metadata *WhitelistedAssetMetadata `json:"metadata,omitempty"`
	// Scores contains the risk scores from various providers.
	Scores []Score `json:"scores,omitempty"`
	// Trails contains the audit trail of actions on this address.
	// Uses shared Trail type.
	Trails []Trail `json:"trails,omitempty"`
	// CreatedAt is the creation timestamp extracted from the "created" trail action.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Approvers contains the approval requirements and status.
	// Uses shared Approvers type.
	Approvers *Approvers `json:"approvers,omitempty"`
	// Attributes are custom key-value attributes.
	Attributes []WhitelistedAddressAttribute `json:"attributes,omitempty"`
	// SignedAddress contains the cryptographic signature data.
	SignedAddress *SignedWhitelistedAddress `json:"signed_address,omitempty"`
	// LinkedInternalAddresses is the list of internal addresses linked to this whitelisted address.
	LinkedInternalAddresses []InternalAddress `json:"linked_internal_addresses,omitempty"`
	// LinkedWallets is the list of internal wallets that can send to this whitelisted address.
	LinkedWallets []InternalWallet `json:"linked_wallets,omitempty"`
}

// InternalAddress represents an internal address linked to a whitelisted address.
type InternalAddress struct {
	// ID is the internal address identifier.
	ID int64 `json:"id"`
	// Label is the human-readable label for the address.
	Label string `json:"label,omitempty"`
}

// InternalWallet represents an internal wallet linked to a whitelisted address.
type InternalWallet struct {
	// ID is the internal wallet identifier.
	ID int64 `json:"id"`
	// Path is the wallet path.
	Path string `json:"path,omitempty"`
	// Label is the human-readable label for the wallet.
	Label string `json:"label,omitempty"`
}

// WhitelistedAddressAttribute represents a custom attribute on a whitelisted address.
type WhitelistedAddressAttribute struct {
	// ID is the attribute identifier.
	ID string `json:"id"`
	// Key is the attribute name.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
	// ContentType is the MIME type of the value.
	ContentType string `json:"content_type,omitempty"`
	// Owner is the user who created the attribute.
	Owner string `json:"owner,omitempty"`
	// Type is the attribute type.
	Type string `json:"type,omitempty"`
	// Subtype is the attribute subtype.
	Subtype string `json:"subtype,omitempty"`
	// IsFile indicates if the value is a file reference.
	IsFile bool `json:"is_file"`
}

// SignedWhitelistedAddress contains the cryptographic signature data for a whitelisted address.
type SignedWhitelistedAddress struct {
	// Payload is the base64-encoded signed payload.
	Payload string `json:"payload,omitempty"`
	// Signatures contains the list of signatures.
	// Uses shared WhitelistSignature type.
	Signatures []WhitelistSignature `json:"signatures,omitempty"`
}

// Score represents a risk score from a scoring provider.
type Score struct {
	// ID is the score identifier.
	ID string `json:"id"`
	// Provider is the scoring provider name.
	Provider string `json:"provider"`
	// Type is the type of score.
	Type string `json:"type"`
	// Score is the score value.
	Score string `json:"score"`
	// UpdateDate is when the score was last updated as a string.
	UpdateDate string `json:"update_date"`
}

// ListWhitelistedAddressesOptions contains options for listing whitelisted addresses.
type ListWhitelistedAddressesOptions struct {
	// Limit is the maximum number of addresses to return.
	Limit int64
	// Offset is the number of addresses to skip.
	Offset int64
	// Blockchain filters by blockchain.
	Blockchain string
	// Network filters by network.
	Network string
	// Currency filters by currency.
	Currency string
	// Query searches address names or values.
	Query string
	// AddressType filters by address type.
	AddressType string
	// IDs filters by specific whitelisted address IDs.
	IDs []string
	// Addresses filters by specific address values.
	Addresses []string
	// IncludeForApproval includes addresses pending approval.
	IncludeForApproval bool
}

// WhitelistedAddressEnvelope wraps a whitelisted address with its verification data.
// This is returned by GetWhitelistedAddressEnvelope after the 6-step verification flow.
type WhitelistedAddressEnvelope struct {
	// ID is the unique identifier for the whitelisted address.
	ID string `json:"id"`
	// Blockchain is the blockchain network (e.g., "ETH", "BTC").
	Blockchain string `json:"blockchain"`
	// Network is the network type (e.g., "mainnet", "testnet").
	Network string `json:"network"`
	// Metadata contains the address metadata.
	Metadata *WhitelistedAssetMetadata `json:"metadata,omitempty"`
	// SignedAddress contains the cryptographic signature data.
	SignedAddress *SignedWhitelistedAddress `json:"signed_address,omitempty"`
	// RulesContainer is the base64-encoded rules container.
	RulesContainer string `json:"rules_container,omitempty"`
	// RulesSignatures is the base64-encoded rules signatures.
	RulesSignatures string `json:"rules_signatures,omitempty"`
	// LinkedInternalAddresses is the list of internal addresses linked to this whitelisted address.
	LinkedInternalAddresses []InternalAddress `json:"linked_internal_addresses,omitempty"`
	// LinkedWallets is the list of internal wallets that can send to this whitelisted address.
	LinkedWallets []InternalWallet `json:"linked_wallets,omitempty"`

	// verifiedWhitelistedAddress is the verified whitelisted address parsed from the payload.
	verifiedWhitelistedAddress *WhitelistedAddress
	// verifiedRulesContainer is the verified and decoded rules container.
	verifiedRulesContainer *DecodedRulesContainer
}

// WhitelistedAddress returns the verified whitelisted address.
// This is populated after successful verification.
func (e *WhitelistedAddressEnvelope) WhitelistedAddress() *WhitelistedAddress {
	return e.verifiedWhitelistedAddress
}

// DecodedRulesContainer returns the verified and decoded rules container.
// This is populated after successful verification.
func (e *WhitelistedAddressEnvelope) DecodedRulesContainer() *DecodedRulesContainer {
	return e.verifiedRulesContainer
}

// SetVerified sets the verified address and rules container.
// This is called internally after successful verification.
func (e *WhitelistedAddressEnvelope) SetVerified(addr *WhitelistedAddress, rules *DecodedRulesContainer) {
	e.verifiedWhitelistedAddress = addr
	e.verifiedRulesContainer = rules
}
