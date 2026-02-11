package model

import "time"

// TagStatistics represents aggregated statistics for a tag.
type TagStatistics struct {
	// TagID is the unique identifier for the tag.
	TagID string `json:"tag_id"`
	// TotalValuation is the total valuation in the base currency main unit (CHF, EUR, USD, etc.).
	TotalValuation string `json:"total_valuation"`
	// Tag contains the full tag information.
	Tag *Tag `json:"tag,omitempty"`
}

// PortfolioStatistics represents aggregated statistics for the global portfolio.
type PortfolioStatistics struct {
	// AvgBalancePerAddress is the average balance per address in the smallest currency unit.
	AvgBalancePerAddress string `json:"avg_balance_per_address"`
	// AddressesCount is the total number of addresses.
	AddressesCount string `json:"addresses_count"`
	// WalletsCount is the total number of wallets.
	WalletsCount string `json:"wallets_count"`
	// TotalBalance is the total balance in the smallest currency unit.
	TotalBalance string `json:"total_balance"`
	// TotalBalanceBaseCurrency is the balance converted in the base currency (fiat) in main unit.
	TotalBalanceBaseCurrency string `json:"total_balance_base_currency"`
}

// PortfolioStatisticsHistoryPoint represents a point in portfolio statistics history.
type PortfolioStatisticsHistoryPoint struct {
	// Timestamp is when this statistics snapshot was taken.
	Timestamp time.Time `json:"timestamp"`
	// Statistics contains the aggregated stats data.
	Statistics *PortfolioStatistics `json:"statistics,omitempty"`
}

// ListTagStatisticsOptions contains options for listing tag statistics.
type ListTagStatisticsOptions struct {
	// Query filters tags where the value contains this substring.
	Query string
	// SortBy specifies the field to sort by (id, tagname, createdat).
	SortBy string
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// ListTagStatisticsResult contains the result of listing tag statistics.
type ListTagStatisticsResult struct {
	// TagStatistics is the list of tag statistics entries.
	TagStatistics []*TagStatistics `json:"tag_statistics"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// GetPortfolioStatisticsHistoryOptions contains options for getting portfolio statistics history.
type GetPortfolioStatisticsHistoryOptions struct {
	// IntervalHours is the interval in hours between data points.
	IntervalHours int64
	// From is the start time for the history range.
	From *time.Time
	// To is the end time for the history range.
	To *time.Time
	// Limit is the maximum number of history points to return.
	Limit int64
	// SortOrder specifies the sort order based on timestamp (ASC or DESC).
	SortOrder string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// GetPortfolioStatisticsHistoryResult contains the result of getting portfolio statistics history.
type GetPortfolioStatisticsHistoryResult struct {
	// HistoryPoints is the list of historical statistics points.
	HistoryPoints []*PortfolioStatisticsHistoryPoint `json:"history_points"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}
