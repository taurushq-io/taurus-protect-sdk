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

// GetPricesResult contains the result of a GetPrices call.
type GetPricesResult struct {
	// BaseCurrency is the base currency used for prices.
	BaseCurrency string `json:"base_currency,omitempty"`
	// Prices is the list of currency prices.
	Prices []*Price `json:"prices"`
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
	// Limit is the maximum number of history points to return.
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
	// Limit is the maximum number of history points to return.
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
