package model

import "time"

// FiatProvider represents a fiat provider and its valuation.
type FiatProvider struct {
	// Provider is the fiat provider name (e.g., "circle").
	Provider string `json:"provider"`
	// Label is the label of the fiat provider set in the config.
	Label string `json:"label"`
	// BaseCurrencyValuation is the valuation in the base currency main unit (CHF, EUR, USD etc.).
	BaseCurrencyValuation string `json:"base_currency_valuation,omitempty"`
}

// FiatProviderAccount represents a fiat provider account with its balance.
type FiatProviderAccount struct {
	// ID is the unique identifier of the fiat provider account.
	ID string `json:"id"`
	// Provider is the fiat provider name (e.g., "circle").
	Provider string `json:"provider"`
	// Label is the label of the fiat provider set in the config.
	Label string `json:"label"`
	// AccountType is the type of account (e.g., "wallet", "bank").
	AccountType string `json:"account_type"`
	// AccountIdentifier is the identifier of the account.
	AccountIdentifier string `json:"account_identifier"`
	// AccountName is the name of the account.
	AccountName string `json:"account_name"`
	// TotalBalance is the balance in the smallest currency unit based on currency decimals.
	TotalBalance string `json:"total_balance"`
	// CreationDate is when the account was created.
	CreationDate time.Time `json:"creation_date"`
	// UpdateDate is when the account was last updated.
	UpdateDate time.Time `json:"update_date"`
	// CurrencyID is the ID of the currency.
	CurrencyID string `json:"currency_id"`
	// CurrencyInfo contains detailed currency information.
	CurrencyInfo *Currency `json:"currency_info,omitempty"`
	// BaseCurrencyValuation is the valuation in the base currency main unit (CHF, EUR, USD etc.).
	BaseCurrencyValuation string `json:"base_currency_valuation,omitempty"`
}

// ListFiatProvidersResult contains the result of listing fiat providers.
type ListFiatProvidersResult struct {
	// FiatProviders is the list of fiat providers.
	FiatProviders []*FiatProvider `json:"fiat_providers"`
	// TotalValuation is the total valuation of all fiat providers in base currency.
	TotalValuation string `json:"total_valuation,omitempty"`
}

// ListFiatProviderAccountsOptions contains options for listing fiat provider accounts.
type ListFiatProviderAccountsOptions struct {
	// Provider is the fiat provider to fetch accounts for (required, e.g., "circle").
	Provider string
	// Label is the label of the fiat provider set in the config (required).
	Label string
	// AccountType filters by account type (e.g., "wallet", "bank").
	AccountType string
	// SortOrder specifies the sort order (ASC or DESC). Default is DESC.
	SortOrder string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// ListFiatProviderAccountsResult contains the result of listing fiat provider accounts.
type ListFiatProviderAccountsResult struct {
	// Accounts is the list of fiat provider accounts.
	Accounts []*FiatProviderAccount `json:"accounts"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}
