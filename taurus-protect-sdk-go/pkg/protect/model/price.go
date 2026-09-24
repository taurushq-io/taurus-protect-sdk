package model

import "time"

// Price represents a currency price with exchange rate information.
type Price struct {
	// Blockchain is the blockchain network.
	Blockchain string `json:"blockchain,omitempty"`
	// CurrencyFrom is the base currency symbol (e.g., "BTC").
	CurrencyFrom string `json:"currency_from,omitempty"`
	// CurrencyTo is the quote currency symbol (e.g., "USD").
	CurrencyTo string `json:"currency_to,omitempty"`
	// Decimals is the number of decimal places for the rate.
	Decimals string `json:"decimals,omitempty"`
	// Rate is the exchange rate between the currencies.
	Rate string `json:"rate,omitempty"`
	// Signatures contains price signatures for verification.
	Signatures []PriceSignature `json:"signatures,omitempty"`
	// ChangePercent24Hour is the 24-hour price change percentage.
	ChangePercent24Hour string `json:"change_percent_24hour,omitempty"`
	// Source is the price data source.
	Source string `json:"source,omitempty"`
	// CreationDate is when the price was created.
	CreationDate time.Time `json:"creation_date,omitempty"`
	// UpdateDate is when the price was last updated.
	UpdateDate time.Time `json:"update_date,omitempty"`
	// CurrencyFromInfo contains detailed information about the base currency.
	CurrencyFromInfo *CurrencyInfo `json:"currency_from_info,omitempty"`
	// CurrencyToInfo contains detailed information about the quote currency.
	CurrencyToInfo *CurrencyInfo `json:"currency_to_info,omitempty"`
	// ID is the price's identifier.
	ID string `json:"id,omitempty"`
	// IsPrimary marks the primary price of its currency pair.
	IsPrimary bool `json:"is_primary,omitempty"`
	// Status is the deviation status of a primary price.
	Status string `json:"status,omitempty"`
}

// PriceSignature represents a signature for price verification.
type PriceSignature struct {
	// UserID is the ID of the user who signed the price.
	UserID string `json:"user_id,omitempty"`
	// Signature is the cryptographic signature.
	Signature string `json:"signature,omitempty"`
}

// PriceHistoryPoint represents a single point in price history (OHLCV data).
type PriceHistoryPoint struct {
	// PeriodStartDate is the start time of the period.
	PeriodStartDate time.Time `json:"period_start_date,omitempty"`
	// Blockchain is the blockchain network.
	Blockchain string `json:"blockchain,omitempty"`
	// CurrencyFrom is the base currency symbol.
	CurrencyFrom string `json:"currency_from,omitempty"`
	// CurrencyTo is the quote currency symbol.
	CurrencyTo string `json:"currency_to,omitempty"`
	// High is the highest price in the period.
	High string `json:"high,omitempty"`
	// Low is the lowest price in the period.
	Low string `json:"low,omitempty"`
	// Open is the opening price of the period.
	Open string `json:"open,omitempty"`
	// Close is the closing price of the period.
	Close string `json:"close,omitempty"`
	// VolumeFrom is the trading volume in the base currency.
	VolumeFrom string `json:"volume_from,omitempty"`
	// VolumeTo is the trading volume in the quote currency.
	VolumeTo string `json:"volume_to,omitempty"`
	// ChangePercent is the price change percentage for the period.
	ChangePercent string `json:"change_percent,omitempty"`
	// CurrencyFromInfo contains detailed information about the base currency.
	CurrencyFromInfo *CurrencyInfo `json:"currency_from_info,omitempty"`
	// CurrencyToInfo contains detailed information about the quote currency.
	CurrencyToInfo *CurrencyInfo `json:"currency_to_info,omitempty"`
}

// ConversionValue represents a converted amount in a target currency.
type ConversionValue struct {
	// Symbol is the currency symbol.
	Symbol string `json:"symbol,omitempty"`
	// Value is the value in the smallest currency unit.
	Value string `json:"value,omitempty"`
	// MainUnitValue is the value in the main currency unit.
	MainUnitValue string `json:"main_unit_value,omitempty"`
	// CurrencyInfo contains detailed information about the currency.
	CurrencyInfo *CurrencyInfo `json:"currency_info,omitempty"`
}

// ConversionResult represents the result of a currency conversion.
type ConversionResult struct {
	// CurrencyFrom is the source currency symbol.
	CurrencyFrom string `json:"currency_from,omitempty"`
	// BaseCurrency is the base currency symbol.
	BaseCurrency string `json:"base_currency,omitempty"`
	// Values contains the converted amounts for each target currency.
	Values []*ConversionValue `json:"values,omitempty"`
	// FullCurrencyFrom contains detailed information about the source currency.
	FullCurrencyFrom *CurrencyInfo `json:"full_currency_from,omitempty"`
	// FullBaseCurrency contains detailed information about the base currency.
	FullBaseCurrency *CurrencyInfo `json:"full_base_currency,omitempty"`
}

// ListPricesOptions filters and pages PriceService.ListPrices.
//
// The currency filter is one of three shapes, chosen by which fields are set: FromCurrencyID
// alone (prices from that currency), FromCurrencyID with ToCurrencyIDs (from it to those), or
// ToCurrencyIDs alone (prices to those currencies). Neither sets no currency filter.
type ListPricesOptions struct {
	// OnlyPrimary keeps only the primary price of each currency pair.
	OnlyPrimary bool
	// SortOrder is ASC or DESC (the server's default).
	SortOrder string
	// FromCurrencyID is the currency the prices convert from.
	FromCurrencyID string
	// ToCurrencyIDs are the currencies the prices convert to.
	ToCurrencyIDs []string
	// PageSize is the page size: 0 selects DefaultPageSize, above MaxPageSize is an error.
	PageSize int64
	// Cursor is a previous page's Page.NextCursor; the SDK then requests the NEXT page.
	Cursor string
}

// ListPricesResult is one page of verified prices.
type ListPricesResult struct {
	// BaseCurrency is the tenant's base currency.
	BaseCurrency string `json:"base_currency,omitempty"`
	// Prices is the list of currency prices, each signature-verified.
	Prices []*Price `json:"prices"`
	// Page continues the list: pass Page.NextCursor as the next Cursor until HasMore is false.
	Page CursorPage `json:"page"`
}

// ConvertOptions contains options for the Convert method.
type ConvertOptions struct {
	// Currency is the source currency to convert from (required).
	Currency string
	// Amount is the amount to convert (required).
	Amount string
	// Symbols filters the target currencies by symbol.
	Symbols []string
	// TargetCurrencyIds filters the target currencies by ID.
	TargetCurrencyIds []string
}

// GetPriceHistoryOptions contains options for the GetPriceHistory method.
type GetPriceHistoryOptions struct {
	// Base is the base currency symbol (required).
	Base string
	// Quote is the quote currency symbol (required).
	Quote string
	// Limit is the number of newest daily points to return: 0 selects DefaultPageSize, above
	// MaxPriceHistoryLimit is an error. Price history cannot page.
	Limit int64
}

// GetPriceHistoryResult contains the result of a GetPriceHistory call.
type GetPriceHistoryResult struct {
	// History is the list of price history points.
	History []*PriceHistoryPoint `json:"history"`
	// Period is the time period between history points.
	Period string `json:"period,omitempty"`
}

// ExportPriceHistoryOptions contains options for the ExportPriceHistory method.
type ExportPriceHistoryOptions struct {
	// CurrencyPairs is the list of currency pairs to export (e.g., ["BTC/USD", "ETH/USD"]).
	CurrencyPairs []string
	// Limit is the number of newest points to export: 0 selects DefaultPageSize; there is no SDK
	// maximum (the server bounds it). The export cannot page.
	Limit int64
	// Format is the export format ("csv" or "json").
	Format string
}

// ExportPriceHistoryResult contains the result of an ExportPriceHistory call.
type ExportPriceHistoryResult struct {
	// Data contains the exported data in the requested format.
	Data string `json:"data,omitempty"`
	// Period is the time period between history points.
	Period string `json:"period,omitempty"`
}
