package model

// MultiFactorSignatureEntityType represents the entity type for a multi-factor signature.
type MultiFactorSignatureEntityType string

const (
	// MFSEntityTypeRequest is for transaction request entities.
	MFSEntityTypeRequest MultiFactorSignatureEntityType = "REQUEST"
	// MFSEntityTypeWhitelistedAddress is for whitelisted address entities.
	MFSEntityTypeWhitelistedAddress MultiFactorSignatureEntityType = "WHITELISTED_ADDRESS"
	// MFSEntityTypeWhitelistedContract is for whitelisted contract entities.
	MFSEntityTypeWhitelistedContract MultiFactorSignatureEntityType = "WHITELISTED_CONTRACT"
)

// MultiFactorSignatureInfo represents information about a multi-factor signature request.
type MultiFactorSignatureInfo struct {
	// ID is the multi-factor signature ID.
	ID string `json:"id"`
	// PayloadToSign is the list of payloads that need to be signed.
	PayloadToSign []string `json:"payload_to_sign,omitempty"`
	// EntityType is the type of entity associated with this signature request.
	EntityType MultiFactorSignatureEntityType `json:"entity_type"`
}

// MultiFactorSignatureResult represents the result of creating multi-factor signatures.
type MultiFactorSignatureResult struct {
	// ID is the ID of the created multi-factor signature batch.
	ID string `json:"id"`
}

// MultiFactorSignatureApprovalResult represents the result of approving a multi-factor signature.
type MultiFactorSignatureApprovalResult struct {
	// SignatureCount is the number of signatures applied.
	SignatureCount string `json:"signature_count"`
}
