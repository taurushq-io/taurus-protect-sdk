package model

import "time"

// FeeV2 represents a native currency fee estimate with detailed information.
// This is returned by the GetFeesV2 endpoint.
type FeeV2 struct {
	// CurrencyID is the unique identifier of the native currency.
	CurrencyID string `json:"currency_id,omitempty"`
	// Value is the fee amount in the native currency's smallest unit.
	Value string `json:"value,omitempty"`
	// Denom is the denomination/unit of the fee value.
	Denom string `json:"denom,omitempty"`
	// CurrencyInfo contains detailed information about the native currency.
	CurrencyInfo *CurrencyInfo `json:"currency_info,omitempty"`
	// UpdateDate is when this fee estimate was last updated.
	UpdateDate time.Time `json:"update_date,omitempty"`
}

// GetFeesV2Result contains the result of the GetFeesV2 call.
type GetFeesV2Result struct {
	// Fees is the list of native currency fee estimates.
	Fees []*FeeV2 `json:"fees"`
}
