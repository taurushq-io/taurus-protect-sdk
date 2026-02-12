package model

import "time"

// Wallet represents a cryptocurrency wallet.
type Wallet struct {
	// ID is the unique identifier for the wallet.
	ID string `json:"id"`
	// Name is the human-readable name of the wallet.
	Name string `json:"name"`
	// Currency is the currency symbol (e.g., "BTC", "ETH").
	Currency string `json:"currency"`
	// Blockchain is the blockchain network.
	Blockchain string `json:"blockchain"`
	// Network is the network type (e.g., "mainnet", "testnet").
	Network string `json:"network"`
	// Balance is the current wallet balance.
	Balance *Balance `json:"balance,omitempty"`
	// IsOmnibus indicates if this is an omnibus wallet.
	IsOmnibus bool `json:"is_omnibus"`
	// Disabled indicates if the wallet is disabled.
	Disabled bool `json:"disabled"`
	// Comment is an optional description.
	Comment string `json:"comment,omitempty"`
	// CustomerID is an optional customer identifier.
	CustomerID string `json:"customer_id,omitempty"`
	// ExternalWalletID is an optional external identifier.
	ExternalWalletID string `json:"external_wallet_id,omitempty"`
	// VisibilityGroupID is the visibility group this wallet belongs to.
	VisibilityGroupID string `json:"visibility_group_id,omitempty"`
	// AddressesCount is the number of addresses in the wallet.
	AddressesCount int64 `json:"addresses_count"`
	// CreatedAt is when the wallet was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the wallet was last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// Attributes are custom key-value attributes.
	Attributes []WalletAttribute `json:"attributes,omitempty"`
	// AccountPath is the HD wallet account derivation path.
	AccountPath string `json:"account_path,omitempty"`
	// CurrencyInfo contains detailed information about the wallet's currency.
	CurrencyInfo *Currency `json:"currency_info,omitempty"`
}

// WalletAttribute represents a custom attribute on a wallet.
type WalletAttribute struct {
	// ID is the attribute identifier.
	ID string `json:"id"`
	// Key is the attribute name.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
	// ContentType is the content type of the attribute.
	ContentType string `json:"content_type,omitempty"`
	// Owner is the owner of the attribute.
	Owner string `json:"owner,omitempty"`
	// Type is the attribute type.
	Type string `json:"type,omitempty"`
	// Subtype is the attribute subtype.
	Subtype string `json:"subtype,omitempty"`
	// IsFile indicates whether this attribute is a file.
	IsFile bool `json:"is_file,omitempty"`
}

// Balance represents a wallet or address balance.
type Balance struct {
	// TotalConfirmed is the total confirmed balance in the smallest currency unit.
	TotalConfirmed string `json:"total_confirmed"`
	// TotalUnconfirmed is the total balance including unconfirmed transactions.
	TotalUnconfirmed string `json:"total_unconfirmed"`
	// AvailableConfirmed is the available confirmed balance ready to be spent.
	AvailableConfirmed string `json:"available_confirmed"`
	// AvailableUnconfirmed is the available balance including unconfirmed transactions.
	AvailableUnconfirmed string `json:"available_unconfirmed"`
	// ReservedConfirmed is the confirmed reserved balance held for pending transactions.
	ReservedConfirmed string `json:"reserved_confirmed"`
	// ReservedUnconfirmed is the reserved unconfirmed balance.
	ReservedUnconfirmed string `json:"reserved_unconfirmed"`
}

// CreateWalletRequest contains parameters for creating a wallet.
type CreateWalletRequest struct {
	// Name is the wallet name (required).
	Name string `json:"name"`
	// Currency is the currency symbol (required).
	Currency string `json:"currency"`
	// Comment is an optional description.
	Comment string `json:"comment,omitempty"`
	// CustomerID is an optional customer identifier.
	CustomerID string `json:"customer_id,omitempty"`
	// ExternalWalletID is an optional external identifier.
	ExternalWalletID string `json:"external_wallet_id,omitempty"`
	// VisibilityGroupID is the visibility group to assign.
	VisibilityGroupID string `json:"visibility_group_id,omitempty"`
}

// ListWalletsOptions contains options for listing wallets.
type ListWalletsOptions struct {
	// Limit is the maximum number of wallets to return.
	Limit int64
	// Offset is the number of wallets to skip.
	Offset int64
	// Currency filters by currency symbol.
	Currency string
	// Query searches wallet names.
	Query string
	// ExcludeDisabled excludes disabled wallets from results.
	ExcludeDisabled bool
}

// Pagination contains pagination information for list responses.
type Pagination struct {
	// Limit is the maximum number of items per page.
	Limit int64 `json:"limit"`
	// Offset is the number of items skipped.
	Offset int64 `json:"offset"`
	// TotalItems is the total number of items available.
	TotalItems int64 `json:"total_items"`
	// HasMore indicates if there are more items to fetch.
	HasMore bool `json:"has_more"`
}

// BalanceHistoryPoint represents a balance at a specific point in time.
type BalanceHistoryPoint struct {
	// PointDate is the timestamp of this balance snapshot.
	PointDate time.Time `json:"point_date"`
	// Balance contains the balance amounts at this point in time.
	Balance *Balance `json:"balance,omitempty"`
}
