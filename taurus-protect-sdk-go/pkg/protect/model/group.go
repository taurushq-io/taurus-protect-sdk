package model

import "time"

// Group represents a user group in the system.
type Group struct {
	// ID is the unique identifier for the group.
	ID string `json:"id"`
	// TenantID is the tenant this group belongs to.
	TenantID string `json:"tenant_id,omitempty"`
	// ExternalGroupID is the external identifier for the group.
	ExternalGroupID string `json:"external_group_id,omitempty"`
	// Name is the name of the group.
	Name string `json:"name"`
	// Email is the email address associated with the group.
	Email string `json:"email,omitempty"`
	// Description is a description of the group.
	Description string `json:"description,omitempty"`
	// EnforcedInRules indicates if the group is enforced in rules.
	EnforcedInRules bool `json:"enforced_in_rules"`
	// Users is the list of users in the group.
	Users []GroupUser `json:"users,omitempty"`
	// CreatedAt is when the group was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the group was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// GroupUser represents a user within a group.
type GroupUser struct {
	// ID is the unique identifier for the user.
	ID string `json:"id"`
	// ExternalUserID is the external identifier for the user.
	ExternalUserID string `json:"external_user_id,omitempty"`
	// EnforcedInRules indicates if the user is enforced in rules.
	EnforcedInRules bool `json:"enforced_in_rules"`
}

// ListGroupsOptions contains options for listing groups.
type ListGroupsOptions struct {
	// Limit is the maximum number of groups to return.
	Limit int64
	// Offset is the number of groups to skip.
	Offset int64
	// IDs filters by group IDs.
	IDs []string
	// ExternalGroupIDs filters by external group IDs.
	ExternalGroupIDs []string
	// Query is a search query to filter groups.
	Query string
}

// ListGroupsResult contains the result of listing groups.
type ListGroupsResult struct {
	// Groups is the list of groups.
	Groups []*Group `json:"groups"`
	// TotalItems is the total number of groups matching the filter.
	TotalItems int64 `json:"total_items"`
	// Offset is the offset used for pagination.
	Offset int64 `json:"offset"`
}
