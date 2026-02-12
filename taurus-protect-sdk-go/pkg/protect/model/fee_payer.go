package model

import "time"

// FeePayer represents a fee payer entity for managing transaction fees.
type FeePayer struct {
	// ID is the unique identifier for the fee payer.
	ID string `json:"id"`
	// TenantID is the tenant identifier.
	TenantID string `json:"tenant_id,omitempty"`
	// Blockchain is the blockchain type (e.g., "ETH").
	Blockchain string `json:"blockchain"`
	// Network is the network identifier (e.g., "mainnet", "goerli").
	Network string `json:"network,omitempty"`
	// Name is the human-readable name for the fee payer.
	Name string `json:"name,omitempty"`
	// CreationDate is when the fee payer was created.
	CreationDate time.Time `json:"creation_date"`
	// ETH contains Ethereum-specific fee payer configuration.
	ETH *FeePayerETH `json:"eth,omitempty"`
}

// FeePayerETH contains Ethereum-specific fee payer configuration.
type FeePayerETH struct {
	// Blockchain is the blockchain type within ETH config.
	Blockchain string `json:"blockchain,omitempty"`
	// Kind indicates the fee payer type (e.g., "local", "remote").
	Kind string `json:"kind,omitempty"`
	// Local contains local fee payer configuration.
	Local *ETHLocalConfig `json:"local,omitempty"`
	// Remote contains remote fee payer configuration.
	Remote *ETHRemoteConfig `json:"remote,omitempty"`
	// RemoteEncrypted contains encrypted remote configuration (base64).
	RemoteEncrypted string `json:"remote_encrypted,omitempty"`
}

// ETHLocalConfig contains configuration for local Ethereum fee payers.
type ETHLocalConfig struct {
	// AddressID is the address ID used for fee payments.
	AddressID string `json:"address_id,omitempty"`
	// ForwarderAddressID is the forwarder contract address ID.
	ForwarderAddressID string `json:"forwarder_address_id,omitempty"`
	// AutoApprove indicates if transactions are auto-approved.
	AutoApprove bool `json:"auto_approve"`
	// CreatorAddressID is the creator address ID.
	CreatorAddressID string `json:"creator_address_id,omitempty"`
	// ForwarderKind is the type of forwarder (e.g., "OpenZeppelinForwarder").
	ForwarderKind string `json:"forwarder_kind,omitempty"`
	// DomainSeparator is the EIP-712 domain separator (base64).
	DomainSeparator string `json:"domain_separator,omitempty"`
}

// ETHRemoteConfig contains configuration for remote Ethereum fee payers.
type ETHRemoteConfig struct {
	// URL is the remote fee payer service URL.
	URL string `json:"url,omitempty"`
	// Username is the authentication username.
	Username string `json:"username,omitempty"`
	// Password is the authentication password.
	Password string `json:"password,omitempty"`
	// PrivateKey is the private key for signing.
	PrivateKey string `json:"private_key,omitempty"`
	// FromAddressID is the source address ID.
	FromAddressID string `json:"from_address_id,omitempty"`
	// ForwarderAddress is the forwarder contract address.
	ForwarderAddress string `json:"forwarder_address,omitempty"`
	// ForwarderAddressID is the forwarder contract address ID.
	ForwarderAddressID string `json:"forwarder_address_id,omitempty"`
	// CreatorAddress is the creator address.
	CreatorAddress string `json:"creator_address,omitempty"`
	// CreatorAddressID is the creator address ID.
	CreatorAddressID string `json:"creator_address_id,omitempty"`
	// ForwarderKind is the type of forwarder (e.g., "OpenZeppelinForwarder").
	ForwarderKind string `json:"forwarder_kind,omitempty"`
	// DomainSeparator is the EIP-712 domain separator (base64).
	DomainSeparator string `json:"domain_separator,omitempty"`
}

// ListFeePayersOptions contains options for listing fee payers.
type ListFeePayersOptions struct {
	// Limit is the maximum number of results to return.
	Limit int64
	// Offset is the number of results to skip.
	Offset int64
	// IDs filters by specific fee payer IDs.
	IDs []string
	// Blockchain filters by blockchain type.
	Blockchain string
	// Network filters by network.
	Network string
}

// ListFeePayersResult contains the result of listing fee payers.
type ListFeePayersResult struct {
	// FeePayers is the list of fee payers.
	FeePayers []*FeePayer `json:"fee_payers"`
	// TotalItems is the total number of items available.
	TotalItems int64 `json:"total_items"`
}

// ChecksumRequest contains the request data for computing a checksum.
type ChecksumRequest struct {
	// Data is the base64-encoded data to compute checksum for.
	Data string `json:"data"`
}

// ChecksumResult contains the result of a checksum computation.
type ChecksumResult struct {
	// Checksum is the computed checksum.
	Checksum string `json:"checksum"`
}
