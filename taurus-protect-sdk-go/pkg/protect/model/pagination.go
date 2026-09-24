package model

import "errors"

// ErrMalformedPagination is wrapped by every error about a reply's paging fields that cannot be
// read: a count that is not a canonical decimal in [0, 2^53-1], or a cursor that claims a next
// page without one. A malformed count is an error rather than a zero, which would end a walk
// early. Match it with errors.Is.
var ErrMalformedPagination = errors.New("malformed pagination")

// Page sizes shared by every list method, whatever the backend's own defaults or maxima.
const (
	// DefaultPageSize is sent when a caller leaves the page size (Limit or PageSize) unset or 0.
	DefaultPageSize = 20
	// MaxPageSize is the largest page size a list method accepts; above it, or below 0, the
	// call fails before any request is sent.
	MaxPageSize = 100
	// MaxPriceHistoryLimit bounds PriceService.GetPriceHistory, which cannot page: it returns
	// the newest Limit daily points (default DefaultPageSize).
	MaxPriceHistoryLimit = 365
)

// Pagination describes one page of an offset-paginated list. It is never nil on success.
//
// Walk a list by passing NextOffset back as the next request's Offset until HasMore is false.
// NextOffset follows the endpoint's own rule (the reply offset, offset plus the rows the server
// returned, or offset plus the limit), so it is not always Offset + len(rows).
type Pagination struct {
	// Limit is the page size that was sent.
	Limit int64 `json:"limit"`
	// Offset is the offset that was sent.
	Offset int64 `json:"offset"`
	// TotalItems is the server's total, reduced by rows the SDK withheld (never negative).
	TotalItems int64 `json:"total_items"`
	// NextOffset is the offset of the next page.
	NextOffset int64 `json:"next_offset"`
	// HasMore reports whether a request at NextOffset can return further rows.
	HasMore bool `json:"has_more"`
}

// CursorPage describes one page of a cursor-paginated list.
//
// Walk a list by passing NextCursor back as the next request's Cursor until HasMore is false.
// Cursors are opaque base64 text; never build or modify one.
type CursorPage struct {
	// PageSize is the page size that was sent.
	PageSize int64 `json:"page_size"`
	// NextCursor continues the list; empty when HasMore is false.
	NextCursor string `json:"next_cursor"`
	// HasMore reports whether another page exists.
	HasMore bool `json:"has_more"`
	// TotalItems is the server's total, only for the lists whose reply carries one.
	TotalItems *int64 `json:"total_items,omitempty"`
}
