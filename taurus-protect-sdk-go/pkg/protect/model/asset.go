package model

import "time"

// AssetFilter represents filter criteria for asset queries.
type AssetFilter struct {
	// Currency is the currency name (e.g., "ALGO", "AVAX", "USDt").
	// Set to "Unknown" when the currency is unknown.
	Currency string `json:"currency"`
	// Kind is an optional asset kind filter ("NFT" or "Unknown").
	Kind string `json:"kind,omitempty"`
	// NFT contains NFT-specific filter criteria when Kind is "NFT".
	NFT *AssetNFTFilter `json:"nft,omitempty"`
	// Unknown contains unknown asset filter criteria when Kind is "Unknown".
	Unknown *AssetUnknownFilter `json:"unknown,omitempty"`
}

// AssetNFTFilter represents NFT-specific filter criteria.
type AssetNFTFilter struct {
	// TokenID is the NFT token identifier.
	TokenID string `json:"token_id,omitempty"`
}

// AssetUnknownFilter represents filter criteria for unknown assets.
type AssetUnknownFilter struct {
	// Blockchain is the blockchain name.
	Blockchain string `json:"blockchain,omitempty"`
	// Arg1 is the first argument (e.g., contract number).
	Arg1 string `json:"arg1,omitempty"`
	// Arg2 is the second argument (e.g., token ID).
	Arg2 string `json:"arg2,omitempty"`
	// Network is the network name.
	Network string `json:"network,omitempty"`
}

// GetAssetAddressesRequest contains parameters for retrieving address-level balances for an asset.
type GetAssetAddressesRequest struct {
	// Asset is the asset filter criteria (required).
	Asset AssetFilter `json:"asset"`
	// SortOrder is the sort order ("ASC" or "DESC").
	SortOrder string `json:"sort_order,omitempty"`
	// WalletID filters for addresses of a specific wallet.
	WalletID string `json:"wallet_id,omitempty"`
	// AddressID filters for a specific address.
	AddressID string `json:"address_id,omitempty"`
	// Addresses filters for specific address strings.
	Addresses []string `json:"addresses,omitempty"`
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64 `json:"page_size,omitempty"`
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string `json:"cursor,omitempty"`
	// CurrentPage and PageRequest (FIRST, PREVIOUS, NEXT, LAST) page by hand; neither can be
	// combined with Cursor.
	CurrentPage string `json:"current_page,omitempty"`
	PageRequest string `json:"page_request,omitempty"`
}

// GetAssetAddressesResult contains the result of retrieving address-level balances for an asset.
type GetAssetAddressesResult struct {
	// Addresses is the list of addresses with balance information.
	Addresses []*Address `json:"addresses"`
	// Page continues the list; TotalItems is the server's total.
	Page CursorPage `json:"page"`
}

// GetAssetWalletsRequest contains parameters for retrieving wallet-level balances for an asset.
type GetAssetWalletsRequest struct {
	// Asset is the asset filter criteria (required).
	Asset AssetFilter `json:"asset"`
	// WalletID filters for a specific wallet.
	WalletID string `json:"wallet_id,omitempty"`
	// WalletName filters by wallet name.
	WalletName string `json:"wallet_name,omitempty"`
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64 `json:"page_size,omitempty"`
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string `json:"cursor,omitempty"`
	// CurrentPage and PageRequest (FIRST, PREVIOUS, NEXT, LAST) page by hand; neither can be
	// combined with Cursor.
	CurrentPage string `json:"current_page,omitempty"`
	PageRequest string `json:"page_request,omitempty"`
}

// GetAssetWalletsResult contains the result of retrieving wallet-level balances for an asset.
type GetAssetWalletsResult struct {
	// Wallets is the list of wallets with balance information.
	Wallets []*Wallet `json:"wallets"`
	// Page continues the list; TotalItems is the server's total.
	Page CursorPage `json:"page"`
}

// AssetResource is an asset in the v2 asset registry (AssetService.QueryAssets).
type AssetResource struct {
	// ID is the asset's identifier, the assetID the other v2 asset methods take.
	ID string `json:"id"`
	// TenantID is the owning tenant.
	TenantID string `json:"tenant_id,omitempty"`
	// Version is the asset's version.
	Version string `json:"version,omitempty"`
	// CreatedAt and UpdatedAt are the asset's timestamps.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Label is the asset's label.
	Label string `json:"label,omitempty"`
	// AssetType is the asset type.
	AssetType string `json:"asset_type,omitempty"`
	// Status is ASSET_VIEW_STATUS_V2_PENDING, _ACTIVE or _FAILED.
	Status string `json:"status,omitempty"`
	// Blockchain and Network locate the asset.
	Blockchain string `json:"blockchain,omitempty"`
	Network    string `json:"network,omitempty"`
	// CurrencyID is the internal currency the asset is registered under.
	CurrencyID string `json:"currency_id,omitempty"`
	// Name, Symbol and Decimals are the public summary fields.
	Name     string `json:"name,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	Decimals string `json:"decimals,omitempty"`
	// ContractAddress is the on-chain contract or token identifier, when the chain has one.
	ContractAddress string `json:"contract_address,omitempty"`
	// Attributes are the asset's key/value attributes.
	Attributes []AssetAttribute `json:"attributes,omitempty"`
	// CantonInstrument is set for a Canton native token.
	CantonInstrument *CantonInstrument `json:"canton_instrument,omitempty"`
}

// AssetAttribute is one key/value attribute of an AssetResource.
type AssetAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CantonInstrument is the Canton-specific part of a Canton native token asset.
type CantonInstrument struct {
	// InstrumentID is the Canton instrument identifier.
	InstrumentID string `json:"instrument_id,omitempty"`
	// ContractID is the instrument configuration contract (cid).
	ContractID string `json:"contract_id,omitempty"`
	// RequireCredentials, Paused and Operator are the instrument configuration.
	RequireCredentials bool   `json:"require_credentials"`
	Paused             bool   `json:"paused"`
	Operator           string `json:"operator,omitempty"`
}

// QueryAssetsOptions filters and pages AssetService.QueryAssets.
type QueryAssetsOptions struct {
	Blockchain      string
	Network         string
	Symbol          string
	ContractAddress string
	Label           string
	CurrencyName    string
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string
}

// QueryAssetsResult is one page of AssetService.QueryAssets.
type QueryAssetsResult struct {
	Assets []*AssetResource `json:"assets"`
	// Page continues the list: pass Page.NextCursor as the next Cursor until HasMore is false.
	Page CursorPage `json:"page"`
}

// AssetAddress is an address holding a v2 asset (AssetService.QueryAssetAddresses).
//
// The server signs none of these rows, so the SDK confirms each one against its verified
// readers before returning it: see Verified.
type AssetAddress struct {
	// Address is the on-chain address. On a Verified row it is the verified address.
	Address string `json:"address"`
	// KYCStatus is KYC_STATUS_V2_APPROVED or KYC_STATUS_V2_REVOKED.
	KYCStatus string `json:"kyc_status,omitempty"`
	// Balance is the address's balance of the asset, as a decimal string.
	Balance string `json:"balance,omitempty"`
	// AddressType is ADDRESS_TYPE_V2_INTERNAL, _WHITELISTED or _EXTERNAL.
	AddressType string `json:"address_type,omitempty"`
	// AddressID is the managed address id; empty when the address is not a managed address.
	AddressID string `json:"address_id,omitempty"`
	// WhitelistedAddressID is the whitelisted address id; empty when it is not whitelisted.
	WhitelistedAddressID string `json:"whitelisted_address_id,omitempty"`
	// Verified reports that the SDK confirmed the row: an INTERNAL row against the HSM-signed
	// managed address with its AddressID, a WHITELISTED row against the verified whitelisted
	// address with its WhitelistedAddressID, the address strings being equal. Every other row
	// (EXTERNAL, or no type) is false: it is on-chain data, never a Taurus-PROTECT address.
	Verified bool `json:"verified"`
}

// QueryAssetAddressesOptions filters and pages AssetService.QueryAssetAddresses.
type QueryAssetAddressesOptions struct {
	// AddressType keeps one kind of address: ADDRESS_TYPE_V2_INTERNAL, _WHITELISTED or _EXTERNAL.
	AddressType string
	// KYCStatus keeps one KYC status: KYC_STATUS_V2_APPROVED or KYC_STATUS_V2_REVOKED.
	KYCStatus string
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string
}

// QueryAssetAddressesResult is one page of AssetService.QueryAssetAddresses.
type QueryAssetAddressesResult struct {
	Addresses []*AssetAddress `json:"addresses"`
	// ExcludedUnverified names the INTERNAL and WHITELISTED rows withheld because their verified
	// counterpart was missing, did not verify or carried another address. ID is the row's
	// AddressID or WhitelistedAddressID, or its address when it has neither. Exclusions do not
	// change Page.
	ExcludedUnverified []ExcludedWhitelistedAddress `json:"excluded_unverified,omitempty"`
	// Page continues the list: pass Page.NextCursor as the next Cursor until HasMore is false.
	Page CursorPage `json:"page"`
}

// AssetOperation is an operation on a v2 asset (AssetService.ListAssetOperations). Which of the
// parameter fields are set depends on Type.
type AssetOperation struct {
	ID      string `json:"id"`
	AssetID string `json:"asset_id,omitempty"`
	// Type is ASSET_OPERATION_TYPE_V2_CREATE, _UPDATE, _IMPORT, _MINT, _BURN, _PAUSE, _UNPAUSE,
	// _PAUSE_ACCOUNT, _UNPAUSE_ACCOUNT or _SET_KYC.
	Type string `json:"type,omitempty"`
	// Status is ASSET_OPERATION_STATUS_V2_PENDING, _PAUSED, _COMPLETED, _FAILED or _CANCELED.
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// InitiatedByAddressID is the internal address that initiated the operation.
	InitiatedByAddressID string `json:"initiated_by_address_id,omitempty"`
	// FailureReason and BlockingReason explain a failed or blocked operation.
	FailureReason  string `json:"failure_reason,omitempty"`
	BlockingReason string `json:"blocking_reason,omitempty"`
	// Label, Price, Decimals, Blockchain, Network, AssetType and Address are the parameters of
	// a create, update or import.
	Label      string `json:"label,omitempty"`
	Price      string `json:"price,omitempty"`
	Decimals   string `json:"decimals,omitempty"`
	Blockchain string `json:"blockchain,omitempty"`
	Network    string `json:"network,omitempty"`
	AssetType  string `json:"asset_type,omitempty"`
	Address    string `json:"address,omitempty"`
	// Amount is the amount of a mint or burn.
	Amount string `json:"amount,omitempty"`
	// NFTMetadata are a mint's NFT metadata; NFTTokenIDs a burn's NFT token ids.
	NFTMetadata []string `json:"nft_metadata,omitempty"`
	NFTTokenIDs []string `json:"nft_token_ids,omitempty"`
	// Target is a mint or burn's destination, or the address a pause-account, unpause-account
	// or set-KYC operation acts on.
	Target *AssetOperationTarget `json:"target,omitempty"`
	// KYCStatus is the status a set-KYC operation applies.
	KYCStatus string `json:"kyc_status,omitempty"`
}

// AssetOperationTarget identifies the address an asset operation acts on.
type AssetOperationTarget struct {
	AddressID            string `json:"address_id,omitempty"`
	WhitelistedAddressID string `json:"whitelisted_address_id,omitempty"`
}

// ListAssetOperationsOptions filters and pages AssetService.ListAssetOperations.
type ListAssetOperationsOptions struct {
	// Type keeps one operation type (an ASSET_OPERATION_TYPE_V2_* value).
	Type string
	// Status keeps one operation status (an ASSET_OPERATION_STATUS_V2_* value).
	Status string
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string
}

// ListAssetOperationsResult is one page of AssetService.ListAssetOperations.
type ListAssetOperationsResult struct {
	Operations []*AssetOperation `json:"operations"`
	// Page continues the list: pass Page.NextCursor as the next Cursor until HasMore is false.
	Page CursorPage `json:"page"`
}
