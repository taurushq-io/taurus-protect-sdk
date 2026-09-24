package model

import "time"

// Timestamps contains common timestamp fields.
type Timestamps struct {
	// CreatedAt is when the resource was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the resource was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}
