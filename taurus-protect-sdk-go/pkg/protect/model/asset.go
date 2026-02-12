package model

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
	// Limit is the maximum number of items to return.
	Limit int64 `json:"limit,omitempty"`
	// Cursor is the pagination cursor for the next page.
	Cursor string `json:"cursor,omitempty"`
	// SortOrder is the sort order ("ASC" or "DESC").
	SortOrder string `json:"sort_order,omitempty"`
	// WalletID filters for addresses of a specific wallet.
	WalletID string `json:"wallet_id,omitempty"`
	// AddressID filters for a specific address.
	AddressID string `json:"address_id,omitempty"`
	// Addresses filters for specific address strings.
	Addresses []string `json:"addresses,omitempty"`
	// CurrentPage is the current page cursor (for cursor-based pagination).
	CurrentPage string `json:"current_page,omitempty"`
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string `json:"page_request,omitempty"`
	// PageSize is the number of items per page.
	PageSize int64 `json:"page_size,omitempty"`
}

// GetAssetAddressesResult contains the result of retrieving address-level balances for an asset.
type GetAssetAddressesResult struct {
	// Addresses is the list of addresses with balance information.
	Addresses []*Address `json:"addresses"`
	// TotalItems is the total number of items available.
	TotalItems int64 `json:"total_items"`
	// NextCursor is the cursor for the next page.
	NextCursor string `json:"next_cursor,omitempty"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// GetAssetWalletsRequest contains parameters for retrieving wallet-level balances for an asset.
type GetAssetWalletsRequest struct {
	// Asset is the asset filter criteria (required).
	Asset AssetFilter `json:"asset"`
	// Limit is the maximum number of items to return.
	Limit int64 `json:"limit,omitempty"`
	// Cursor is the pagination cursor for the next page.
	Cursor string `json:"cursor,omitempty"`
	// WalletID filters for a specific wallet.
	WalletID string `json:"wallet_id,omitempty"`
	// WalletName filters by wallet name.
	WalletName string `json:"wallet_name,omitempty"`
	// CurrentPage is the current page cursor (for cursor-based pagination).
	CurrentPage string `json:"current_page,omitempty"`
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string `json:"page_request,omitempty"`
	// PageSize is the number of items per page.
	PageSize int64 `json:"page_size,omitempty"`
}

// GetAssetWalletsResult contains the result of retrieving wallet-level balances for an asset.
type GetAssetWalletsResult struct {
	// Wallets is the list of wallets with balance information.
	Wallets []*Wallet `json:"wallets"`
	// TotalItems is the total number of items available.
	TotalItems int64 `json:"total_items"`
	// NextCursor is the cursor for the next page.
	NextCursor string `json:"next_cursor,omitempty"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}
