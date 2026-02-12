package taurusnetwork

import "time"

// SharedAddressStatus represents the status of a shared address.
type SharedAddressStatus string

// SharedAddressStatus constants.
const (
	SharedAddressStatusNew      SharedAddressStatus = "new"
	SharedAddressStatusPending  SharedAddressStatus = "pending"
	SharedAddressStatusRejected SharedAddressStatus = "rejected"
	SharedAddressStatusAccepted SharedAddressStatus = "accepted"
	SharedAddressStatusUnshared SharedAddressStatus = "unshared"
)

// SharedAssetStatus represents the status of a shared asset.
type SharedAssetStatus string

// SharedAssetStatus constants.
const (
	SharedAssetStatusNew      SharedAssetStatus = "new"
	SharedAssetStatusPending  SharedAssetStatus = "pending"
	SharedAssetStatusRejected SharedAssetStatus = "rejected"
	SharedAssetStatusAccepted SharedAssetStatus = "accepted"
	SharedAssetStatusUnshared SharedAssetStatus = "unshared"
)

// SharedAddress represents a shared address in Taurus-NETWORK.
type SharedAddress struct {
	// ID is the unique identifier of the shared address (uuid).
	ID string `json:"id"`
	// InternalAddressID is the ID of the deposit address if your participant is the owner.
	InternalAddressID string `json:"internal_address_id,omitempty"`
	// WhitelistedAddressID is the ID of the whitelisted address if your participant is the target.
	WhitelistedAddressID string `json:"whitelisted_address_id,omitempty"`
	// OwnerParticipantID is the ID of the participant who owns the address.
	OwnerParticipantID string `json:"owner_participant_id"`
	// TargetParticipantID is the ID of the participant the address is shared with.
	TargetParticipantID string `json:"target_participant_id"`
	// Blockchain is the blockchain of the shared address.
	Blockchain string `json:"blockchain"`
	// Network is the network of the shared address.
	Network string `json:"network"`
	// Address is the actual address value.
	Address string `json:"address"`
	// OriginLabel is the label from the origin.
	OriginLabel string `json:"origin_label,omitempty"`
	// OriginCreationDate is when the original address was created.
	OriginCreationDate time.Time `json:"origin_creation_date,omitempty"`
	// OriginDeletionDate is when the original address was deleted.
	OriginDeletionDate time.Time `json:"origin_deletion_date,omitempty"`
	// CreatedAt is when the shared address was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the shared address was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// TargetAcceptedAt is when the target participant accepted the shared address.
	TargetAcceptedAt time.Time `json:"target_accepted_at,omitempty"`
	// Status is the current status of the shared address.
	Status SharedAddressStatus `json:"status"`
	// ProofOfOwnership contains proof of ownership information.
	ProofOfOwnership *ProofOfOwnership `json:"proof_of_ownership,omitempty"`
	// PledgesCount is the count of active pledges linked to the shared address.
	PledgesCount int64 `json:"pledges_count"`
	// Trails contains the history of status changes.
	Trails []*SharedAddressTrail `json:"trails,omitempty"`
}

// SharedAddressTrail represents a status change trail for a shared address.
type SharedAddressTrail struct {
	// ID is the unique identifier of the trail entry.
	ID string `json:"id"`
	// SharedAddressID is the ID of the related shared address.
	SharedAddressID string `json:"shared_address_id"`
	// AddressStatus is the status at this point in the trail.
	AddressStatus string `json:"address_status"`
	// Comment is an optional comment for the status change.
	Comment string `json:"comment,omitempty"`
	// CreatedAt is when this trail entry was created.
	CreatedAt time.Time `json:"created_at"`
}

// ProofOfOwnership contains proof of ownership information for a shared address.
type ProofOfOwnership struct {
	// SignedPayloadHash is the hex encoded SHA256 hash of the payload.
	SignedPayloadHash string `json:"signed_payload_hash,omitempty"`
	// SignedPayloadAsString is the signed payload as a string.
	SignedPayloadAsString string `json:"signed_payload_as_string,omitempty"`
}

// SharedAsset represents a shared asset in Taurus-NETWORK.
type SharedAsset struct {
	// ID is the unique identifier of the shared asset.
	ID string `json:"id"`
	// WhitelistedContractAddressID is the ID of the whitelisted contract address.
	WhitelistedContractAddressID string `json:"whitelisted_contract_address_id,omitempty"`
	// OwnerParticipantID is the ID of the participant who owns the asset.
	OwnerParticipantID string `json:"owner_participant_id"`
	// TargetParticipantID is the ID of the participant the asset is shared with.
	TargetParticipantID string `json:"target_participant_id"`
	// Blockchain is the blockchain of the shared asset.
	Blockchain string `json:"blockchain"`
	// Network is the network of the shared asset.
	Network string `json:"network"`
	// Name is the name of the asset.
	Name string `json:"name"`
	// Symbol is the symbol of the asset.
	Symbol string `json:"symbol"`
	// Decimals is the number of decimal places for the asset.
	Decimals string `json:"decimals"`
	// ContractAddress is the contract address of the asset.
	ContractAddress string `json:"contract_address"`
	// TokenID is the token ID (for NFTs).
	TokenID string `json:"token_id,omitempty"`
	// Kind is the type of asset.
	Kind string `json:"kind,omitempty"`
	// OriginCreationDate is when the original asset was created.
	OriginCreationDate time.Time `json:"origin_creation_date,omitempty"`
	// OriginDeletionDate is when the original asset was deleted.
	OriginDeletionDate time.Time `json:"origin_deletion_date,omitempty"`
	// CreatedAt is when the shared asset was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the shared asset was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// TargetAcceptedAt is when the target participant accepted the shared asset.
	TargetAcceptedAt time.Time `json:"target_accepted_at,omitempty"`
	// TargetRejectedAt is when the target participant rejected the shared asset.
	TargetRejectedAt time.Time `json:"target_rejected_at,omitempty"`
	// Status is the current status of the shared asset.
	Status SharedAssetStatus `json:"status"`
	// Trails contains the history of status changes.
	Trails []*SharedAssetTrail `json:"trails,omitempty"`
}

// SharedAssetTrail represents a status change trail for a shared asset.
type SharedAssetTrail struct {
	// ID is the unique identifier of the trail entry.
	ID string `json:"id"`
	// SharedAssetID is the ID of the related shared asset.
	SharedAssetID string `json:"shared_asset_id"`
	// AssetStatus is the status at this point in the trail.
	AssetStatus string `json:"asset_status"`
	// Comment is an optional comment for the status change.
	Comment string `json:"comment,omitempty"`
	// CreatedAt is when this trail entry was created.
	CreatedAt time.Time `json:"created_at"`
}

// KeyValueAttribute represents a key-value attribute pair.
type KeyValueAttribute struct {
	// Key is the attribute key.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
}

// ListSharedAddressesOptions contains options for listing shared addresses.
type ListSharedAddressesOptions struct {
	// ParticipantID filters addresses where the participant is either owner or target.
	ParticipantID string
	// OwnerParticipantID filters by the owner participant ID.
	OwnerParticipantID string
	// TargetParticipantID filters by the target participant ID.
	TargetParticipantID string
	// Blockchain filters by blockchain.
	Blockchain string
	// Network filters by network.
	Network string
	// IDs filters by specific shared address IDs.
	IDs []string
	// Statuses filters by status values.
	Statuses []string
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// ListSharedAddressesResult contains the result of listing shared addresses.
type ListSharedAddressesResult struct {
	// SharedAddresses is the list of shared addresses.
	SharedAddresses []*SharedAddress `json:"shared_addresses"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// ListSharedAssetsOptions contains options for listing shared assets.
type ListSharedAssetsOptions struct {
	// ParticipantID filters assets where the participant is either owner or target.
	ParticipantID string
	// OwnerParticipantID filters by the owner participant ID.
	OwnerParticipantID string
	// TargetParticipantID filters by the target participant ID.
	TargetParticipantID string
	// Blockchain filters by blockchain.
	Blockchain string
	// Network filters by network.
	Network string
	// IDs filters by specific shared asset IDs.
	IDs []string
	// Statuses filters by status values.
	Statuses []string
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// ListSharedAssetsResult contains the result of listing shared assets.
type ListSharedAssetsResult struct {
	// SharedAssets is the list of shared assets.
	SharedAssets []*SharedAsset `json:"shared_assets"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// ShareAddressRequest contains the request parameters for sharing an address.
type ShareAddressRequest struct {
	// ToParticipantID is the ID of the participant to share the address with.
	ToParticipantID string `json:"to_participant_id"`
	// AddressID is the ID of the internal address to share.
	AddressID string `json:"address_id"`
	// KeyValueAttributes are optional key-value attributes to attach to the shared address.
	KeyValueAttributes []KeyValueAttribute `json:"key_value_attributes,omitempty"`
}

// ShareWhitelistedAssetRequest contains the request parameters for sharing an asset.
type ShareWhitelistedAssetRequest struct {
	// ToParticipantID is the ID of the participant to share the asset with.
	ToParticipantID string `json:"to_participant_id"`
	// WhitelistedContractID is the ID of the whitelisted contract to share.
	WhitelistedContractID string `json:"whitelisted_contract_id"`
}
