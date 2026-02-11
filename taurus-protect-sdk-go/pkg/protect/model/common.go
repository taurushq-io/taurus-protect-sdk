package model

import "time"

// Cursor represents pagination cursor information.
type Cursor struct {
	// Limit is the maximum number of items per page.
	Limit int64 `json:"limit"`
	// Offset is the number of items to skip.
	Offset int64 `json:"offset"`
	// Total is the total number of items available.
	Total int64 `json:"total"`
	// HasMore indicates if there are more items after this page.
	HasMore bool `json:"has_more"`
}

// PageRequest represents a pagination request.
type PageRequest struct {
	// Limit is the maximum number of items to return.
	Limit int64
	// Offset is the number of items to skip.
	Offset int64
}

// Timestamps contains common timestamp fields.
type Timestamps struct {
	// CreatedAt is when the resource was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the resource was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}
