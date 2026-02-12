package model

import "time"

// User represents a user in the Taurus-PROTECT system.
type User struct {
	// ID is the unique identifier for the user.
	ID string `json:"id"`
	// TenantID is the tenant this user belongs to.
	TenantID string `json:"tenant_id"`
	// ExternalUserID is an optional external identifier.
	ExternalUserID string `json:"external_user_id,omitempty"`
	// Username is the user's login name.
	Username string `json:"username,omitempty"`
	// Email is the user's email address.
	Email string `json:"email"`
	// FirstName is the user's first name.
	FirstName string `json:"first_name,omitempty"`
	// LastName is the user's last name.
	LastName string `json:"last_name,omitempty"`
	// Status is the user's status (e.g., "ACTIVE", "INACTIVE").
	Status string `json:"status,omitempty"`
	// Roles is the list of role IDs assigned to the user.
	Roles []string `json:"roles,omitempty"`
	// Groups is the list of groups the user belongs to.
	Groups []UserGroup `json:"groups,omitempty"`
	// Attributes are custom key-value attributes.
	Attributes []UserAttribute `json:"attributes,omitempty"`
	// PublicKey is the user's public key if set.
	PublicKey string `json:"public_key,omitempty"`
	// PasswordChanged indicates if the user has changed their initial password.
	PasswordChanged bool `json:"password_changed"`
	// TotpEnabled indicates if TOTP is enabled for the user.
	TotpEnabled bool `json:"totp_enabled"`
	// EnforcedInRules indicates if the user is enforced in rules.
	EnforcedInRules bool `json:"enforced_in_rules"`
	// PublicKeyEnforcedInRules indicates if the user's public key is enforced in rules.
	PublicKeyEnforcedInRules bool `json:"public_key_enforced_in_rules"`
	// CreatedAt is when the user was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the user was last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// LastLogin is when the user last logged in.
	LastLogin time.Time `json:"last_login,omitempty"`
}

// UserGroup represents a group that a user belongs to.
type UserGroup struct {
	// ID is the unique identifier for the group.
	ID string `json:"id"`
	// ExternalGroupID is an optional external identifier.
	ExternalGroupID string `json:"external_group_id,omitempty"`
	// EnforcedInRules indicates if this group membership is enforced in rules.
	EnforcedInRules bool `json:"enforced_in_rules"`
}

// UserAttribute represents a custom attribute on a user.
type UserAttribute struct {
	// ID is the attribute identifier.
	ID string `json:"id"`
	// Key is the attribute name.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
	// ContentType is the MIME type of the attribute value.
	ContentType string `json:"content_type,omitempty"`
	// Owner is the owner of the attribute.
	Owner string `json:"owner,omitempty"`
	// Type is the attribute type.
	Type string `json:"type,omitempty"`
	// SubType is the attribute subtype.
	SubType string `json:"sub_type,omitempty"`
	// IsFile indicates if the attribute value is a file reference.
	IsFile bool `json:"is_file"`
	// CreatedAt is when the attribute was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the attribute was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// ListUsersOptions contains options for listing users.
type ListUsersOptions struct {
	// Limit is the maximum number of users to return.
	Limit int64
	// Offset is the number of users to skip.
	Offset int64
	// IDs filters by specific user IDs.
	IDs []string
	// ExternalUserIDs filters by external user IDs.
	ExternalUserIDs []string
	// Emails filters by email addresses.
	Emails []string
	// Roles filters by role IDs (users must have all specified roles).
	Roles []string
	// GroupIDs filters by group IDs.
	GroupIDs []string
	// Query searches user fields.
	Query string
	// Status filters by user status.
	Status string
	// TotpEnabled filters by TOTP enabled status.
	TotpEnabled *bool
	// ExcludeTechnicalUsers excludes technical/system users from results.
	ExcludeTechnicalUsers bool
}

// ListUsersResult contains the result of listing users.
type ListUsersResult struct {
	// Users is the list of users.
	Users []*User `json:"users"`
	// TotalItems is the total number of users matching the filter.
	TotalItems int64 `json:"total_items"`
	// Offset is the number of users skipped.
	Offset int64 `json:"offset"`
}
