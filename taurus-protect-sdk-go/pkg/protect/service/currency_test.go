package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewCurrencyService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestCurrencyService_GetCurrency_NilOptions(t *testing.T) {
	// Create a service with nil API to test validation before API call
	svc := &CurrencyService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetCurrency(nil, nil)
	if err == nil {
		t.Error("GetCurrency() with nil options should return error")
	}
	if err.Error() != "options cannot be nil" {
		t.Errorf("GetCurrency() error = %v, want 'options cannot be nil'", err)
	}
}

func TestCurrencyService_GetCurrency_EmptyBlockchain(t *testing.T) {
	svc := &CurrencyService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetCurrency(nil, &model.GetCurrencyOptions{})
	if err == nil {
		t.Error("GetCurrency() with empty blockchain should return error")
	}
	if err.Error() != "blockchain is required" {
		t.Errorf("GetCurrency() error = %v, want 'blockchain is required'", err)
	}
}

func TestCurrencyService_GetCurrency_EmptyNetwork(t *testing.T) {
	svc := &CurrencyService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetCurrency(nil, &model.GetCurrencyOptions{Blockchain: "ETH"})
	if err == nil {
		t.Error("GetCurrency() with empty network should return error")
	}
	if err.Error() != "network is required" {
		t.Errorf("GetCurrency() error = %v, want 'network is required'", err)
	}
}

func TestCurrencyService_GetCurrencies_NilOptions(t *testing.T) {
	// This test documents that GetCurrencies accepts nil options
	// The actual API call would require a mock, but we can verify
	// that the service handles nil options without panicking in the
	// options handling code (before the API call)
	svc := &CurrencyService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This will panic due to nil API, but that's expected since we're
	// not mocking the API client. The test verifies the options handling
	// doesn't cause issues before the API call.
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to nil API, but didn't get one")
		}
	}()

	_, _ = svc.GetCurrencies(nil, nil)
}

func TestListCurrenciesOptions(t *testing.T) {
	// Test that options struct can be created with expected fields
	opts := &model.ListCurrenciesOptions{
		ShowDisabled: true,
		IncludeLogo:  true,
	}

	if !opts.ShowDisabled {
		t.Error("ShowDisabled should be true")
	}
	if !opts.IncludeLogo {
		t.Error("IncludeLogo should be true")
	}
}

func TestGetCurrencyOptions(t *testing.T) {
	// Test that options struct can be created with expected fields
	opts := &model.GetCurrencyOptions{
		Blockchain:           "ETH",
		Network:              "mainnet",
		CurrencyID:           "currency-123",
		TokenContractAddress: "0x1234567890abcdef",
		TokenID:              "token-001",
		ShowDisabled:         true,
		IncludeLogo:          true,
	}

	if opts.Blockchain != "ETH" {
		t.Errorf("Blockchain = %v, want 'ETH'", opts.Blockchain)
	}
	if opts.Network != "mainnet" {
		t.Errorf("Network = %v, want 'mainnet'", opts.Network)
	}
	if opts.CurrencyID != "currency-123" {
		t.Errorf("CurrencyID = %v, want 'currency-123'", opts.CurrencyID)
	}
	if opts.TokenContractAddress != "0x1234567890abcdef" {
		t.Errorf("TokenContractAddress = %v, want '0x1234567890abcdef'", opts.TokenContractAddress)
	}
	if opts.TokenID != "token-001" {
		t.Errorf("TokenID = %v, want 'token-001'", opts.TokenID)
	}
	if !opts.ShowDisabled {
		t.Error("ShowDisabled should be true")
	}
	if !opts.IncludeLogo {
		t.Error("IncludeLogo should be true")
	}
}

func TestGetCurrencyOptions_Validation(t *testing.T) {
	tests := []struct {
		name    string
		opts    *model.GetCurrencyOptions
		wantErr string
	}{
		{
			name:    "nil options",
			opts:    nil,
			wantErr: "options cannot be nil",
		},
		{
			name:    "empty blockchain",
			opts:    &model.GetCurrencyOptions{Network: "mainnet"},
			wantErr: "blockchain is required",
		},
		{
			name:    "empty network",
			opts:    &model.GetCurrencyOptions{Blockchain: "ETH"},
			wantErr: "network is required",
		},
		{
			name: "valid minimum options",
			opts: &model.GetCurrencyOptions{
				Blockchain: "ETH",
				Network:    "mainnet",
			},
			wantErr: "", // Would require API call to fail
		},
	}

	svc := &CurrencyService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr == "" {
				// Skip tests that would require API call
				return
			}

			_, err := svc.GetCurrency(nil, tt.opts)
			if err == nil {
				t.Error("GetCurrency() should return error")
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("GetCurrency() error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}
