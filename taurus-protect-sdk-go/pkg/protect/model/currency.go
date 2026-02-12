package model

// Currency represents a cryptocurrency or fiat currency.
type Currency struct {
	// ID is the unique identifier of the currency.
	ID string `json:"id"`
	// Name is the name of the currency (e.g., "ethereum").
	Name string `json:"name"`
	// Symbol is the shorthand symbol for the currency (e.g., "ETH").
	Symbol string `json:"symbol"`
	// DisplayName is the display name for the currency (e.g., "Ethereum").
	DisplayName string `json:"display_name"`
	// Type is the type of currency: "token", "fiat", "native", or "signet".
	Type string `json:"type"`
	// Blockchain is the blockchain the currency is associated with (e.g., "ETH", "BTC").
	Blockchain string `json:"blockchain"`
	// Network is the network or environment (e.g., "mainnet", "testnet").
	Network string `json:"network"`
	// Decimals is the number of decimal places the currency uses (e.g., 18 for ETH).
	Decimals int64 `json:"decimals"`
	// CoinTypeIndex is the BIP44 coin type index (e.g., "0" for Bitcoin, "60" for Ethereum).
	CoinTypeIndex string `json:"coin_type_index"`
	// ContractAddress is the smart contract address if currency is a token.
	ContractAddress string `json:"contract_address,omitempty"`
	// TokenID is the unique identifier for the token (used for NFTs and some blockchains).
	TokenID string `json:"token_id,omitempty"`
	// WlcaID is the whitelisted contract address ID associated with the currency.
	WlcaID int64 `json:"wlca_id,omitempty"`
	// Logo is the currency logo in Data URI scheme (base64 encoded).
	Logo string `json:"logo,omitempty"`
	// IsToken indicates if the currency is a token.
	IsToken bool `json:"is_token"`
	// IsERC20 indicates if the token is an ERC-20 token.
	IsERC20 bool `json:"is_erc20"`
	// IsFA12 indicates if the currency is based on FA12 standard (Tezos).
	IsFA12 bool `json:"is_fa12"`
	// IsFA20 indicates if the currency is based on FA20 standard (Tezos).
	IsFA20 bool `json:"is_fa20"`
	// IsNFT indicates if the currency represents a Non-Fungible Token.
	IsNFT bool `json:"is_nft"`
	// IsFiat indicates if the currency is a fiat currency (e.g., CHF, EUR, USD).
	IsFiat bool `json:"is_fiat"`
	// IsUTXOBased indicates if the currency is UTXO-based (e.g., Bitcoin).
	IsUTXOBased bool `json:"is_utxo_based"`
	// IsAccountBased indicates if the currency is account-based (e.g., Ethereum).
	IsAccountBased bool `json:"is_account_based"`
	// HasStaking indicates if the currency supports staking.
	HasStaking bool `json:"has_staking"`
	// Enabled indicates if the currency is enabled in the current tenant.
	Enabled bool `json:"enabled"`
}

// ListCurrenciesOptions contains options for listing currencies.
type ListCurrenciesOptions struct {
	// ShowDisabled includes currencies disabled by business rules when true.
	ShowDisabled bool
	// IncludeLogo includes currency logos in the response when true.
	IncludeLogo bool
}

// GetCurrencyOptions contains options for retrieving a single currency.
type GetCurrencyOptions struct {
	// Blockchain is the blockchain to filter by (required, e.g., "ETH", "BTC").
	Blockchain string
	// Network is the network or environment (required, e.g., "mainnet", "testnet").
	Network string
	// CurrencyID is the unique identifier of the currency to filter by (optional).
	CurrencyID string
	// TokenContractAddress filters by token contract address (optional).
	// If set, returns a token with the specified contract address.
	// If not set, returns the native currency of the blockchain.
	TokenContractAddress string
	// TokenID differentiates assets within a contract address (optional).
	// Used for blockchains like ALGO and XTZ where one contract may contain multiple assets.
	TokenID string
	// ShowDisabled includes the currency even if disabled by business rules when true.
	ShowDisabled bool
	// IncludeLogo includes the currency logo in the response when true.
	IncludeLogo bool
}
