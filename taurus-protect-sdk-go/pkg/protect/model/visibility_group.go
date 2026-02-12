package model

import "time"

// VisibilityGroup represents a restricted visibility group in the system.
// Visibility groups control which users can access specific resources like wallets.
type VisibilityGroup struct {
	// ID is the unique identifier for the visibility group.
	ID string `json:"id"`
	// TenantID is the tenant this visibility group belongs to.
	TenantID string `json:"tenant_id,omitempty"`
	// Name is the display name of the visibility group.
	Name string `json:"name"`
	// Description is an optional description of the visibility group.
	Description string `json:"description,omitempty"`
	// Users is the list of users that belong to this visibility group.
	Users []*VisibilityGroupUser `json:"users,omitempty"`
	// UserCount is the number of users in this visibility group.
	UserCount int64 `json:"user_count"`
	// CreatedAt is when the visibility group was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the visibility group was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// VisibilityGroupUser represents a user reference within a visibility group.
type VisibilityGroupUser struct {
	// ID is the unique identifier for the user.
	ID string `json:"id"`
	// ExternalUserID is the external identifier for the user.
	ExternalUserID string `json:"external_user_id,omitempty"`
}

// Note: User type is defined in user.go to avoid duplication
