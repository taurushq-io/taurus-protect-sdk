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
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string
	// CurrentPage and PageRequest (FIRST, PREVIOUS, NEXT, LAST) page by hand; neither can be
	// combined with Cursor.
	CurrentPage string
	PageRequest string
}

// ListFiatProviderAccountsResult contains the result of listing fiat provider accounts.
type ListFiatProviderAccountsResult struct {
	// Accounts is the list of fiat provider accounts.
	Accounts []*FiatProviderAccount `json:"accounts"`
	// Page continues the list: pass Page.NextCursor as the next Cursor until HasMore is false.
	Page CursorPage `json:"page"`
}

// FiatProviderEntity is an entity registered with a fiat provider.
type FiatProviderEntity struct {
	// ID is the unique identifier of the entity.
	ID string `json:"id"`
	// Provider is the fiat provider name (e.g., "circle").
	Provider string `json:"provider"`
	// Label is the label of the fiat provider set in the config.
	Label string `json:"label"`
	// AccountIdentifier is the identifier of the entity's account at the provider.
	AccountIdentifier string `json:"account_identifier"`
	// Name is the entity's name.
	Name string `json:"name"`
	// Details carries provider-specific details.
	Details string `json:"details,omitempty"`
	// CreationDate is when the entity was created.
	CreationDate time.Time `json:"creation_date"`
	// UpdateDate is when the entity was last updated.
	UpdateDate time.Time `json:"update_date"`
}

// ListFiatProviderEntitiesOptions contains options for listing fiat provider entities.
type ListFiatProviderEntitiesOptions struct {
	// Provider filters by fiat provider (e.g., "circle").
	Provider string
	// Label filters by the provider label set in the config.
	Label string
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string
}

// ListFiatProviderEntitiesResult is one page of fiat provider entities.
type ListFiatProviderEntitiesResult struct {
	// Entities is the list of entities on this page.
	Entities []*FiatProviderEntity `json:"entities"`
	// Page continues the list: pass Page.NextCursor as the next Cursor until HasMore is false.
	Page CursorPage `json:"page"`
}
