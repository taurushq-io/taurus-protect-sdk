package model

// BusinessRule represents a business rule in the system.
type BusinessRule struct {
	// ID is the unique identifier of the business rule.
	ID string `json:"id"`
	// TenantID is the tenant identifier.
	TenantID string `json:"tenant_id,omitempty"`
	// Currency is the currency this rule applies to (if applicable).
	Currency string `json:"currency,omitempty"`
	// WalletID is the wallet ID this rule applies to (deprecated, use EntityType/EntityID).
	WalletID string `json:"wallet_id,omitempty"`
	// RuleKey is the key identifying the rule type.
	RuleKey string `json:"rule_key"`
	// RuleValue is the value of the rule.
	RuleValue string `json:"rule_value"`
	// RuleGroup is the group this rule belongs to.
	RuleGroup string `json:"rule_group,omitempty"`
	// RuleDescription is a description of the rule.
	RuleDescription string `json:"rule_description,omitempty"`
	// RuleValidation is the validation pattern for the rule value.
	RuleValidation string `json:"rule_validation,omitempty"`
	// AddressID is the address ID this rule applies to (deprecated, use EntityType/EntityID).
	AddressID string `json:"address_id,omitempty"`
	// CurrencyInfo contains detailed currency information if available.
	CurrencyInfo *Currency `json:"currency_info,omitempty"`
	// EntityType indicates what this rule applies to (global, currency, wallet, address, exchange, exchange_account, tn_participant).
	EntityType string `json:"entity_type,omitempty"`
	// EntityID is the identifier of the affected entity.
	EntityID string `json:"entity_id,omitempty"`
}

// ListBusinessRulesOptions contains options for listing business rules.
type ListBusinessRulesOptions struct {
	// IDs filters by business rule IDs.
	IDs []string
	// RuleKeys filters by rule keys.
	RuleKeys []string
	// RuleGroups filters by rule groups.
	RuleGroups []string
	// WalletIDs filters by wallet IDs (deprecated, use EntityType/EntityIDs).
	WalletIDs []string
	// CurrencyIDs filters by currency IDs.
	CurrencyIDs []string
	// AddressIDs filters by address IDs (deprecated, use EntityType/EntityIDs).
	AddressIDs []string
	// Level filters by rule level (deprecated, use EntityType).
	// One of: '', 'global', 'currency', 'address', 'wallet'.
	Level string
	// EntityType filters rules by what they apply to.
	// One of: 'global', 'currency', 'wallet', 'address', 'exchange', 'exchange_account', 'tn_participant'.
	EntityType string
	// EntityIDs filters by entity identifiers.
	EntityIDs []string
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string
	// PageRequest indicates which page to request (FIRST, PREVIOUS, NEXT, LAST).
	PageRequest string
	// PageSize is the number of items per page.
	PageSize int64
}

// ListBusinessRulesResult contains the result of listing business rules.
type ListBusinessRulesResult struct {
	// BusinessRules is the list of business rules.
	BusinessRules []*BusinessRule `json:"business_rules"`
	// CurrentPage is the base64-encoded cursor for the current page.
	CurrentPage string `json:"current_page,omitempty"`
	// HasPrevious indicates if there is a previous page.
	HasPrevious bool `json:"has_previous"`
	// HasNext indicates if there is a next page.
	HasNext bool `json:"has_next"`
}

// UpdateTransactionsEnabledRequest contains the request to update the transactions enabled rule.
type UpdateTransactionsEnabledRequest struct {
	// Enabled indicates whether transactions should be enabled.
	Enabled bool `json:"enabled"`
}
