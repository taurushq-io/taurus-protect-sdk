package model

import "time"

// AuditTrail represents an audit trail entry for tracking actions in the system.
type AuditTrail struct {
	// ID is the unique identifier for the audit trail entry.
	ID string `json:"id"`
	// User contains information about the user who performed the action.
	User *UserInfo `json:"user,omitempty"`
	// Entity is the type of entity that was affected (e.g., "Wallet", "Request").
	Entity string `json:"entity"`
	// Action is the action that was performed (e.g., "Create", "Update", "Delete").
	Action string `json:"action"`
	// SubAction is an optional sub-action that provides more detail.
	SubAction string `json:"sub_action,omitempty"`
	// Details contains additional details about the action in JSON format.
	Details string `json:"details,omitempty"`
	// CreationDate is when the audit trail entry was created.
	CreationDate time.Time `json:"creation_date"`
}

// UserInfo represents information about a user in an audit trail.
type UserInfo struct {
	// ID is the unique identifier for the user.
	ID string `json:"id"`
	// ExternalUserID is the external identifier for the user.
	ExternalUserID string `json:"external_user_id,omitempty"`
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// Deleted indicates if the user has been deleted.
	Deleted bool `json:"deleted"`
}

// ListAuditTrailsOptions contains options for listing audit trails.
type ListAuditTrailsOptions struct {
	// ExternalUserID filters by the external user ID who performed the action.
	ExternalUserID string
	// Entities filters by entity types (e.g., ["Wallet", "Request"]).
	Entities []string
	// Actions filters by action types (e.g., ["Create", "Update"]).
	Actions []string
	// CreationDateFrom filters for entries created on or after this date.
	CreationDateFrom *time.Time
	// CreationDateTo filters for entries created on or before this date.
	CreationDateTo *time.Time
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
	// SortBy specifies the columns to sort by (e.g., ["CreationDate"]).
	SortBy []string
	// SortOrder specifies the sort order (ASC or DESC).
	SortOrder string
}

// ListAuditTrailsResult contains the result of listing audit trails.
type ListAuditTrailsResult struct {
	// AuditTrails is the list of audit trail entries.
	AuditTrails []*AuditTrail `json:"audit_trails"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// ExportAuditTrailsOptions contains options for exporting audit trails.
type ExportAuditTrailsOptions struct {
	// ExternalUserID filters by the external user ID who performed the action.
	ExternalUserID string
	// Entities filters by entity types (e.g., ["Wallet", "Request"]).
	Entities []string
	// Actions filters by action types (e.g., ["Create", "Update"]).
	Actions []string
	// CreationDateFrom filters for entries created on or after this date.
	CreationDateFrom *time.Time
	// CreationDateTo filters for entries created on or before this date.
	CreationDateTo *time.Time
	// Format specifies the export format ("csv" or "json").
	Format string
}

// ExportAuditTrailsResult contains the result of exporting audit trails.
type ExportAuditTrailsResult struct {
	// Result contains the exported data (CSV or JSON string).
	Result string `json:"result"`
	// TotalItems is the total number of items exported.
	TotalItems int64 `json:"total_items"`
}
