package model

// WhitelistedContract represents a whitelisted contract in the system.
// A whitelisted contract is a smart contract (token or NFT) that has been
// approved for use within the Taurus-PROTECT platform.
type WhitelistedContract struct {
	// ID is the unique identifier for the whitelisted contract.
	ID string `json:"id"`
	// TenantID is the tenant identifier.
	TenantID string `json:"tenant_id,omitempty"`
	// Status is the approval status (e.g., "approved", "pending").
	Status string `json:"status,omitempty"`
	// Action is the pending action type if any.
	Action string `json:"action,omitempty"`
	// Blockchain is the blockchain symbol (e.g., "ETH", "XTZ").
	Blockchain string `json:"blockchain,omitempty"`
	// Network is the network type (e.g., "mainnet", "testnet").
	Network string `json:"network,omitempty"`
	// Rule is the governance rule applied to this contract.
	Rule string `json:"rule,omitempty"`
	// RulesContainer is the serialized rules container.
	RulesContainer string `json:"rules_container,omitempty"`
	// RulesSignatures contains the super-admin signatures for the rules.
	RulesSignatures string `json:"rules_signatures,omitempty"`
	// BusinessRuleEnabled indicates if the currency is enabled in business rules.
	BusinessRuleEnabled bool `json:"business_rule_enabled"`
	// Metadata contains additional metadata about the contract.
	Metadata *WhitelistedAssetMetadata `json:"metadata,omitempty"`
	// SignedContractAddress contains the cryptographic signature data.
	SignedContractAddress *SignedContractAddress `json:"signed_contract_address,omitempty"`
	// Approvers contains the approval requirements and status.
	Approvers *Approvers `json:"approvers,omitempty"`
	// Attributes are custom key-value attributes.
	Attributes []WhitelistedContractAttribute `json:"attributes,omitempty"`
	// Trails contains the audit trail of actions on this contract.
	Trails []Trail `json:"trails,omitempty"`
}

// WhitelistedContractAttribute represents a custom attribute on a whitelisted contract.
type WhitelistedContractAttribute struct {
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

// CreateWhitelistedContractRequest contains parameters for creating a whitelisted contract.
type CreateWhitelistedContractRequest struct {
	// Blockchain is the symbol of the blockchain on which to whitelist a contract.
	// Required. Examples: "ETH", "XTZ", "ALGO", "SOL".
	Blockchain string `json:"blockchain"`
	// Symbol is the symbol of the token.
	// Required.
	Symbol string `json:"symbol"`
	// ContractAddress is the address of the smart contract.
	// For ALGO standard assets, this is the address of the creator of the asset.
	// For XLM, this is the Stellar address.
	// For ICP, this is the ICP principal.
	ContractAddress string `json:"contract_address,omitempty"`
	// Name is the name for the contract.
	Name string `json:"name,omitempty"`
	// Decimals is the number of decimal places.
	Decimals string `json:"decimals,omitempty"`
	// TokenID is the token identifier.
	// For XTZ, this is used to create a FA2 token.
	// For ALGO standard assets, this is the globally unique asset id.
	TokenID string `json:"token_id,omitempty"`
	// Kind is the contract kind.
	// Examples: "", "NFT_AUTO", "NFT_XTZ_FA2", "NFT_EVM_ERC721", "NFT_EVM_ERC1155",
	// "SOL_TOKEN", "SOL_TOKEN2022", "NFT_EVM_CRYPTOPUNKS".
	Kind string `json:"kind,omitempty"`
	// Network is the name of the blockchain network.
	// Examples: "mainnet", "testnet", "sepolia".
	Network string `json:"network,omitempty"`
}

// UpdateWhitelistedContractRequest contains parameters for updating a whitelisted contract.
type UpdateWhitelistedContractRequest struct {
	// Symbol is the symbol of the token.
	Symbol string `json:"symbol,omitempty"`
	// Name is the name for the contract.
	Name string `json:"name,omitempty"`
	// Decimals is the number of decimal places.
	Decimals string `json:"decimals,omitempty"`
}

// ApproveWhitelistedContractRequest contains parameters for approving a whitelisted contract.
type ApproveWhitelistedContractRequest struct {
	// ID is the whitelisted contract ID to approve.
	ID string `json:"id"`
	// Comment is an optional comment for the approval.
	Comment string `json:"comment,omitempty"`
}

// RejectWhitelistedContractRequest contains parameters for rejecting a whitelisted contract.
type RejectWhitelistedContractRequest struct {
	// ID is the whitelisted contract ID to reject.
	ID string `json:"id"`
	// Comment is an optional comment for the rejection.
	Comment string `json:"comment,omitempty"`
}

// DeleteWhitelistedContractRequest contains parameters for deleting a whitelisted contract.
// Note: This operation is deprecated due to complex dependencies.
type DeleteWhitelistedContractRequest struct {
	// ID is the whitelisted contract ID to delete.
	ID string `json:"id"`
}

// CreateWhitelistedContractAttributeRequest contains parameters for creating an attribute.
type CreateWhitelistedContractAttributeRequest struct {
	// Key is the attribute name.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
	// ContentType is the MIME type of the value.
	ContentType string `json:"content_type,omitempty"`
	// Type is the attribute type.
	Type string `json:"type,omitempty"`
	// Subtype is the attribute subtype.
	Subtype string `json:"subtype,omitempty"`
	// IsFile indicates if the value is a file reference.
	IsFile bool `json:"is_file,omitempty"`
}

// ListWhitelistedContractsOptions contains options for listing whitelisted contracts.
type ListWhitelistedContractsOptions struct {
	// Limit is the maximum number of contracts to return (max 100).
	Limit int64
	// Offset is the number of contracts to skip.
	Offset int64
	// Query searches across multiple fields (customerid, address, blockchain, label, memo, addresstype).
	Query string
	// Blockchain filters by blockchain symbol.
	Blockchain string
	// Network filters by network (e.g., "mainnet", "testnet", "sepolia").
	Network string
	// IncludeForApproval includes contracts pending approval in results.
	IncludeForApproval bool
	// KindTypes filters by contract kind type ("nft" or "token").
	KindTypes []string
	// IDs filters by specific whitelisted contract IDs.
	IDs []string
}

// ListWhitelistedContractsForApprovalOptions contains options for listing contracts pending approval.
type ListWhitelistedContractsForApprovalOptions struct {
	// Limit is the maximum number of contracts to return (max 100).
	Limit int64
	// Offset is the number of contracts to skip.
	Offset int64
	// IDs filters by specific whitelisted contract IDs.
	IDs []string
}

// ListWhitelistedContractsResult contains the results of listing whitelisted contracts.
type ListWhitelistedContractsResult struct {
	// Contracts is the list of whitelisted contracts.
	Contracts []*WhitelistedContract
	// Pagination contains pagination information.
	Pagination *Pagination
}
