package model

import "time"

// Exchange represents an exchange account in the Taurus-PROTECT system.
type Exchange struct {
	// ID is the unique identifier for the exchange account.
	ID string `json:"id"`
	// Exchange is the name of the exchange (e.g., "Binance", "Coinbase").
	Exchange string `json:"exchange"`
	// Account is the account identifier on the exchange.
	Account string `json:"account"`
	// Currency is the currency identifier for this exchange account.
	Currency string `json:"currency"`
	// Type is the type of exchange account.
	Type string `json:"type"`
	// TotalBalance is the total balance in the smallest currency unit.
	TotalBalance string `json:"total_balance"`
	// Status is the status of the exchange account.
	Status string `json:"status"`
	// Container is the container identifier.
	Container string `json:"container"`
	// Label is a user-defined label for the exchange account.
	Label string `json:"label"`
	// DisplayLabel is the display label for the exchange account.
	DisplayLabel string `json:"display_label"`
	// HasWLA indicates if the exchange account has a whitelisted address.
	HasWLA bool `json:"has_wla"`
	// BaseCurrencyValuation is the valuation in the base currency (e.g., CHF, EUR, USD).
	BaseCurrencyValuation string `json:"base_currency_valuation,omitempty"`
	// CurrencyInfo contains detailed information about the currency.
	CurrencyInfo *Currency `json:"currency_info,omitempty"`
	// CreationDate is when the exchange account was created.
	CreationDate time.Time `json:"creation_date"`
	// UpdateDate is when the exchange account was last updated.
	UpdateDate time.Time `json:"update_date"`
}

// ExchangeCounterparty represents an exchange counterparty with its total valuation.
type ExchangeCounterparty struct {
	// Name is the name of the exchange.
	Name string `json:"name"`
	// BaseCurrencyValuation is the total valuation in the base currency.
	BaseCurrencyValuation string `json:"base_currency_valuation"`
}

// ListExchangesOptions contains options for listing exchange accounts.
type ListExchangesOptions struct {
	// CurrencyID filters by currency ID.
	CurrencyID string
	// IncludeBaseCurrencyValuation includes the base currency valuation in results.
	IncludeBaseCurrencyValuation bool
	// ExchangeLabel filters by exchange label.
	ExchangeLabel string
	// SortOrder is the sort order (ASC or DESC). Default is DESC.
	SortOrder string
	// Status filters by exchange account status.
	Status string
	// OnlyPositiveBalance excludes accounts with zero balance when true.
	OnlyPositiveBalance bool
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// ListExchangesResult contains the result of listing exchange accounts.
type ListExchangesResult struct {
	// Exchanges is the list of exchange accounts.
	Exchanges []*Exchange `json:"exchanges"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// ListExchangeCounterpartiesResult contains the result of listing exchange counterparties.
type ListExchangeCounterpartiesResult struct {
	// ExchangesTotalValuation is the total valuation of all exchanges in the base currency.
	ExchangesTotalValuation string `json:"exchanges_total_valuation"`
	// Exchanges is the list of exchange counterparties.
	Exchanges []*ExchangeCounterparty `json:"exchanges"`
}

// ExportExchangesOptions contains options for exporting exchange accounts.
type ExportExchangesOptions struct {
	// Format specifies the export format ("csv" or "json").
	Format string
}

// ExportExchangesResult contains the result of exporting exchange accounts.
type ExportExchangesResult struct {
	// Result contains the exported data (CSV or JSON string).
	Result string `json:"result"`
	// TotalItems is the total number of items exported.
	TotalItems int64 `json:"total_items"`
}
