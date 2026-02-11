package model

import "time"

// Address represents a blockchain address within a wallet.
type Address struct {
	// ID is the unique identifier for the address.
	ID string `json:"id"`
	// WalletID is the ID of the wallet containing this address.
	WalletID string `json:"wallet_id"`
	// Address is the actual blockchain address string.
	Address string `json:"address"`
	// AlternateAddress is an alternate address format if available.
	AlternateAddress string `json:"alternate_address,omitempty"`
	// Label is the human-readable name for the address.
	Label string `json:"label"`
	// Comment is an optional description.
	Comment string `json:"comment,omitempty"`
	// Currency is the currency symbol.
	Currency string `json:"currency"`
	// CustomerID is an optional customer identifier.
	CustomerID string `json:"customer_id,omitempty"`
	// ExternalAddressID is an optional external identifier.
	ExternalAddressID string `json:"external_address_id,omitempty"`
	// AddressPath is the derivation path for the address.
	AddressPath string `json:"address_path,omitempty"`
	// AddressIndex is the index used for address generation.
	AddressIndex string `json:"address_index,omitempty"`
	// Nonce is the current nonce of the address.
	Nonce string `json:"nonce,omitempty"`
	// Status is the creation status (created, creating, signed, observed, confirmed).
	Status string `json:"status,omitempty"`
	// Balance is the current address balance.
	Balance *Balance `json:"balance,omitempty"`
	// Signature is the cryptographic signature for integrity verification.
	Signature string `json:"signature,omitempty"`
	// Disabled indicates if the address is disabled.
	Disabled bool `json:"disabled"`
	// CanUseAllFunds indicates if all funds can be used.
	CanUseAllFunds bool `json:"can_use_all_funds"`
	// CreatedAt is when the address was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the address was last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// Attributes are custom key-value attributes.
	Attributes []AddressAttribute `json:"attributes,omitempty"`
	// LinkedWhitelistedAddressIDs are linked whitelisted address IDs.
	LinkedWhitelistedAddressIDs []string `json:"linked_whitelisted_address_ids,omitempty"`
}

// AddressAttribute represents a custom attribute on an address.
type AddressAttribute struct {
	// ID is the attribute identifier.
	ID string `json:"id"`
	// Key is the attribute name.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
}

// CreateAddressRequest contains parameters for creating an address.
type CreateAddressRequest struct {
	// WalletID is the ID of the wallet to create the address in (required).
	WalletID string `json:"wallet_id"`
	// Label is the human-readable name for the address (required).
	Label string `json:"label"`
	// Comment is an optional description.
	Comment string `json:"comment,omitempty"`
	// CustomerID is an optional customer identifier.
	CustomerID string `json:"customer_id,omitempty"`
	// ExternalAddressID is an optional external identifier.
	ExternalAddressID string `json:"external_address_id,omitempty"`
	// Type is the address type (e.g., p2pkh, p2sh_p2wpkh, p2wpkh for BTC).
	Type string `json:"type,omitempty"`
	// NonHardenedDerivation indicates whether to use non-hardened derivation.
	NonHardenedDerivation bool `json:"non_hardened_derivation,omitempty"`
}

// ListAddressesOptions contains options for listing addresses.
type ListAddressesOptions struct {
	// WalletID filters by wallet ID.
	WalletID string
	// Limit is the maximum number of addresses to return.
	Limit int64
	// Offset is the number of addresses to skip.
	Offset int64
	// Query searches address labels or addresses.
	Query string
	// ExcludeDisabled excludes disabled addresses from results.
	ExcludeDisabled bool
}

// ProofOfReserve represents the proof of reserve for an address.
type ProofOfReserve struct {
	// Curve is the cryptographic curve used.
	Curve string `json:"curve,omitempty"`
	// Cipher is the cipher algorithm used.
	Cipher string `json:"cipher,omitempty"`
	// Path is the derivation path.
	Path string `json:"path,omitempty"`
	// Address is the blockchain address.
	Address string `json:"address,omitempty"`
	// PublicKey is the public key (base64 encoded).
	PublicKey string `json:"public_key,omitempty"`
	// Challenge is the challenge string used for the proof.
	Challenge string `json:"challenge,omitempty"`
	// ChallengeResponse is the signed response to the challenge (base64 encoded).
	ChallengeResponse string `json:"challenge_response,omitempty"`
	// Type is the reserve type.
	Type string `json:"type,omitempty"`
	// StakePublicKey is the staking public key (base64 encoded).
	StakePublicKey string `json:"stake_public_key,omitempty"`
	// StakeChallengeResponse is the staking challenge response (base64 encoded).
	StakeChallengeResponse string `json:"stake_challenge_response,omitempty"`
}
