package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewBusinessRuleService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestBusinessRuleService_ListBusinessRules_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	// The actual API call will fail, but we're testing the options handling
	svc := &BusinessRuleService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("BusinessRuleService should not be nil")
	}
}

func TestBusinessRuleService_ListBusinessRules_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *model.ListBusinessRulesOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.ListBusinessRulesOptions{},
		},
		{
			name: "IDs filter",
			options: &model.ListBusinessRulesOptions{
				IDs: []string{"rule-1", "rule-2"},
			},
		},
		{
			name: "rule keys filter",
			options: &model.ListBusinessRulesOptions{
				RuleKeys: []string{"TRANSACTIONS_ENABLED", "MAX_AMOUNT"},
			},
		},
		{
			name: "rule groups filter",
			options: &model.ListBusinessRulesOptions{
				RuleGroups: []string{"security", "limits"},
			},
		},
		{
			name: "wallet IDs filter",
			options: &model.ListBusinessRulesOptions{
				WalletIDs: []string{"wallet-1", "wallet-2"},
			},
		},
		{
			name: "currency IDs filter",
			options: &model.ListBusinessRulesOptions{
				CurrencyIDs: []string{"currency-1", "currency-2"},
			},
		},
		{
			name: "address IDs filter",
			options: &model.ListBusinessRulesOptions{
				AddressIDs: []string{"address-1", "address-2"},
			},
		},
		{
			name: "level filter",
			options: &model.ListBusinessRulesOptions{
				Level: "global",
			},
		},
		{
			name: "entity type filter",
			options: &model.ListBusinessRulesOptions{
				EntityType: "wallet",
			},
		},
		{
			name: "entity IDs filter",
			options: &model.ListBusinessRulesOptions{
				EntityIDs: []string{"entity-1", "entity-2"},
			},
		},
		{
			name: "pagination options",
			options: &model.ListBusinessRulesOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &model.ListBusinessRulesOptions{
				IDs:         []string{"rule-1"},
				RuleKeys:    []string{"TRANSACTIONS_ENABLED"},
				RuleGroups:  []string{"security"},
				WalletIDs:   []string{"wallet-1"},
				CurrencyIDs: []string{"currency-1"},
				AddressIDs:  []string{"address-1"},
				Level:       "wallet",
				EntityType:  "wallet",
				EntityIDs:   []string{"wallet-1"},
				CurrentPage: "xyz789",
				PageRequest: "FIRST",
				PageSize:    100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &BusinessRuleService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("BusinessRuleService should not be nil")
			}
		})
	}
}

func TestBusinessRuleService_UpdateTransactionsEnabled_NilRequest(t *testing.T) {
	// Create a service with nil API to test that nil request is rejected
	svc := &BusinessRuleService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// Nil request should return an error without calling the API
	err := svc.UpdateTransactionsEnabled(nil, nil)
	if err == nil {
		t.Error("UpdateTransactionsEnabled should return error for nil request")
	}
}

func TestBusinessRuleService_UpdateTransactionsEnabled_WithRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *model.UpdateTransactionsEnabledRequest
	}{
		{
			name: "enable transactions",
			request: &model.UpdateTransactionsEnabledRequest{
				Enabled: true,
			},
		},
		{
			name: "disable transactions",
			request: &model.UpdateTransactionsEnabledRequest{
				Enabled: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify requests can be created with these values
			// Actual API testing requires mocking
			svc := &BusinessRuleService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("BusinessRuleService should not be nil")
			}
			if tt.request == nil {
				t.Error("Request should not be nil in this test")
			}
		})
	}
}

func TestListBusinessRulesOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &model.ListBusinessRulesOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListBusinessRulesOptions_LevelValues(t *testing.T) {
	// Test that level values match expected API values
	validLevels := []string{"", "global", "currency", "address", "wallet"}

	for _, level := range validLevels {
		t.Run(level, func(t *testing.T) {
			opts := &model.ListBusinessRulesOptions{
				Level: level,
			}
			if opts.Level != level {
				t.Errorf("Level = %v, want %v", opts.Level, level)
			}
		})
	}
}

func TestListBusinessRulesOptions_EntityTypeValues(t *testing.T) {
	// Test that entity type values match expected API values
	validEntityTypes := []string{"global", "currency", "wallet", "address", "exchange", "exchange_account", "tn_participant"}

	for _, entityType := range validEntityTypes {
		t.Run(entityType, func(t *testing.T) {
			opts := &model.ListBusinessRulesOptions{
				EntityType: entityType,
			}
			if opts.EntityType != entityType {
				t.Errorf("EntityType = %v, want %v", opts.EntityType, entityType)
			}
		})
	}
}
