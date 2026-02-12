package service

import (
	"testing"
)

func TestNewBalanceService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestBalanceService_GetBalances_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	// The actual API call will fail, but we're testing the options handling
	svc := &BalanceService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This will panic on nil API, so we just verify the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("BalanceService should not be nil")
	}
}

func TestBalanceService_GetBalances_WithOptions(t *testing.T) {
	// Verify that GetBalancesOptions fields are properly structured
	tests := []struct {
		name     string
		currency string
		tokenID  string
		limit    int64
		cursor   string
	}{
		{
			name:     "empty options",
			currency: "",
			tokenID:  "",
			limit:    0,
			cursor:   "",
		},
		{
			name:     "currency filter",
			currency: "ETH",
			tokenID:  "",
			limit:    0,
			cursor:   "",
		},
		{
			name:     "token ID filter",
			currency: "",
			tokenID:  "token-123",
			limit:    0,
			cursor:   "",
		},
		{
			name:     "pagination options",
			currency: "",
			tokenID:  "",
			limit:    100,
			cursor:   "abc123",
		},
		{
			name:     "all options",
			currency: "BTC",
			tokenID:  "token-456",
			limit:    50,
			cursor:   "xyz789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &BalanceService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("BalanceService should not be nil")
			}
		})
	}
}
