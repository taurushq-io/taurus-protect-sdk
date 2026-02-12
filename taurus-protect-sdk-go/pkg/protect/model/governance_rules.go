package model

import "time"

// GovernanceRuleset represents a governance ruleset from the API.
type GovernanceRuleset struct {
	// RulesContainer is the base64-encoded rules container protobuf.
	RulesContainer string `json:"rules_container"`
	// Signatures contains the SuperAdmin signatures for this ruleset.
	Signatures []RuleUserSignature `json:"signatures,omitempty"`
	// Locked indicates if the rules are locked.
	Locked bool `json:"locked"`
	// CreatedAt is when the rules were created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the rules were last updated.
	UpdatedAt time.Time `json:"updated_at"`
	// Trails contains the audit trail for this ruleset.
	Trails []RulesTrail `json:"trails,omitempty"`
}

// RuleUserSignature represents a SuperAdmin signature on the rules.
type RuleUserSignature struct {
	// UserID is the ID of the user who signed.
	UserID string `json:"user_id"`
	// Signature is the base64-encoded signature.
	Signature string `json:"signature"`
}

// RulesTrail represents an audit trail entry for governance rules.
type RulesTrail struct {
	// ID is the trail entry identifier.
	ID string `json:"id"`
	// UserID is the ID of the user who performed the action.
	UserID string `json:"user_id"`
	// ExternalUserID is the external user identifier.
	ExternalUserID string `json:"external_user_id,omitempty"`
	// Action is the action performed.
	Action string `json:"action"`
	// Comment is an optional comment about the action.
	Comment string `json:"comment,omitempty"`
	// Date is when the action was performed.
	Date time.Time `json:"date"`
}

// SuperAdminPublicKey represents a SuperAdmin's public key.
type SuperAdminPublicKey struct {
	// UserID is the ID of the SuperAdmin user.
	UserID string `json:"user_id"`
	// PublicKey is the PEM-encoded public key.
	PublicKey string `json:"public_key"`
}

// GovernanceRulesHistoryResult contains the result of a rules history query.
type GovernanceRulesHistoryResult struct {
	// Rules is the list of governance rulesets in history.
	Rules []*GovernanceRuleset `json:"rules"`
	// TotalItems is the total number of items.
	TotalItems int64 `json:"total_items"`
	// Cursor is the pagination cursor for the next page.
	Cursor string `json:"cursor,omitempty"`
}

// ListRulesHistoryOptions contains options for listing rules history.
type ListRulesHistoryOptions struct {
	// Limit is the maximum number of items to return.
	Limit int64
	// Cursor is the pagination cursor from a previous request.
	Cursor string
}
