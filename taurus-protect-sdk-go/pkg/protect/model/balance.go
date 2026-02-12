package model

// AssetBalance represents a balance for a specific asset.
type AssetBalance struct {
	// Asset contains information about the asset.
	Asset *Asset `json:"asset,omitempty"`
	// Balance contains the balance amounts.
	Balance *Balance `json:"balance,omitempty"`
}

// Asset represents a cryptocurrency asset.
type Asset struct {
	// Currency is the currency symbol (e.g., "ETH", "BTC", "USDt").
	Currency string `json:"currency"`
	// Kind indicates the type of asset.
	Kind string `json:"kind,omitempty"`
	// CurrencyInfo contains detailed currency information.
	CurrencyInfo *CurrencyInfo `json:"currency_info,omitempty"`
}

// CurrencyInfo represents detailed information about a currency.
type CurrencyInfo struct {
	// ID is the unique identifier of the currency.
	ID string `json:"id,omitempty"`
	// Name is the name of the currency.
	Name string `json:"name,omitempty"`
	// Symbol is the shorthand symbol for the currency.
	Symbol string `json:"symbol,omitempty"`
	// DisplayName is the display name for the currency.
	DisplayName string `json:"display_name,omitempty"`
	// Blockchain is the blockchain the currency is associated with.
	Blockchain string `json:"blockchain,omitempty"`
	// Network is the network (e.g., "mainnet", "testnet").
	Network string `json:"network,omitempty"`
	// Decimals is the number of decimal places.
	Decimals string `json:"decimals,omitempty"`
	// ContractAddress is the smart contract address if applicable.
	ContractAddress string `json:"contract_address,omitempty"`
	// TokenID is the unique token ID if applicable (e.g., for NFTs).
	TokenID string `json:"token_id,omitempty"`
	// Type is the type of currency (e.g., "token", "fiat", "native").
	Type string `json:"type,omitempty"`
	// IsToken indicates if this is a token.
	IsToken bool `json:"is_token,omitempty"`
	// IsERC20 indicates if this is an ERC-20 token.
	IsERC20 bool `json:"is_erc20,omitempty"`
	// IsNFT indicates if this is an NFT.
	IsNFT bool `json:"is_nft,omitempty"`
	// IsFiat indicates if this is a fiat currency.
	IsFiat bool `json:"is_fiat,omitempty"`
	// Enabled indicates if the currency is enabled.
	Enabled bool `json:"enabled,omitempty"`
}

// GetBalancesOptions contains options for listing balances.
type GetBalancesOptions struct {
	// Currency filters by currency ID or symbol.
	Currency string
	// TokenID filters by token ID.
	TokenID string
	// Limit is the maximum number of balances to return.
	Limit int64
	// Cursor is the pagination cursor for fetching the next page.
	Cursor string
}

// GetBalancesResult contains the result of a GetBalances call.
type GetBalancesResult struct {
	// Balances is the list of asset balances.
	Balances []*AssetBalance `json:"balances"`
	// Total is the total number of balances available.
	Total int64 `json:"total"`
	// NextCursor is the cursor for fetching the next page.
	NextCursor string `json:"next_cursor,omitempty"`
}
