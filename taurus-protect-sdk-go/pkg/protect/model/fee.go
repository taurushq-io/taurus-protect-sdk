package model

import "time"

// Fee represents a fee estimate for a blockchain transaction.
// This is returned by the deprecated v1 GetFees endpoint.
type Fee struct {
	// Key identifies the fee type (e.g., blockchain:network or currency ID).
	Key string `json:"key"`
	// Value is the fee amount as a string.
	Value string `json:"value"`
}

// FeeV2 represents a native currency fee estimate with detailed information.
// This is returned by the v2 GetFees endpoint.
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

// GetFeesResult contains the result of the deprecated v1 GetFees call.
type GetFeesResult struct {
	// Fees is the list of fee estimates.
	Fees []*Fee `json:"fees"`
}

// GetFeesV2Result contains the result of the v2 GetFees call.
type GetFeesV2Result struct {
	// Fees is the list of native currency fee estimates.
	Fees []*FeeV2 `json:"fees"`
}
