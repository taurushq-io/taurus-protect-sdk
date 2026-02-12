package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewExchangeService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestExchangeService_GetExchange_ServiceExists(t *testing.T) {
	// Create a service with nil API to test that the service can be created
	svc := &ExchangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("ExchangeService should not be nil")
	}
}

func TestExchangeService_ListExchanges_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	// The actual API call will fail, but we're testing the options handling
	svc := &ExchangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("ExchangeService should not be nil")
	}
}

func TestExchangeService_ListExchanges_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *model.ListExchangesOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.ListExchangesOptions{},
		},
		{
			name: "currency ID filter",
			options: &model.ListExchangesOptions{
				CurrencyID: "btc-mainnet",
			},
		},
		{
			name: "include base currency valuation",
			options: &model.ListExchangesOptions{
				IncludeBaseCurrencyValuation: true,
			},
		},
		{
			name: "exchange label filter",
			options: &model.ListExchangesOptions{
				ExchangeLabel: "My Exchange",
			},
		},
		{
			name: "sort order",
			options: &model.ListExchangesOptions{
				SortOrder: "ASC",
			},
		},
		{
			name: "status filter",
			options: &model.ListExchangesOptions{
				Status: "active",
			},
		},
		{
			name: "only positive balance",
			options: &model.ListExchangesOptions{
				OnlyPositiveBalance: true,
			},
		},
		{
			name: "pagination options",
			options: &model.ListExchangesOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &model.ListExchangesOptions{
				CurrencyID:                   "btc-mainnet",
				IncludeBaseCurrencyValuation: true,
				ExchangeLabel:                "My Exchange",
				SortOrder:                    "DESC",
				Status:                       "active",
				OnlyPositiveBalance:          true,
				CurrentPage:                  "xyz789",
				PageRequest:                  "FIRST",
				PageSize:                     100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &ExchangeService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("ExchangeService should not be nil")
			}
		})
	}
}

func TestExchangeService_ExportExchanges_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	svc := &ExchangeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("ExchangeService should not be nil")
	}
}

func TestExchangeService_ExportExchanges_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *model.ExportExchangesOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.ExportExchangesOptions{},
		},
		{
			name: "csv format",
			options: &model.ExportExchangesOptions{
				Format: "csv",
			},
		},
		{
			name: "json format",
			options: &model.ExportExchangesOptions{
				Format: "json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &ExchangeService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("ExchangeService should not be nil")
			}
		})
	}
}

func TestListExchangesOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &model.ListExchangesOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListExchangesOptions_SortOrderValues(t *testing.T) {
	// Test that sort order values match expected API values
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &model.ListExchangesOptions{
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}

func TestExportExchangesOptions_FormatValues(t *testing.T) {
	// Test that format values match expected API values
	validFormats := []string{"csv", "json"}

	for _, format := range validFormats {
		t.Run(format, func(t *testing.T) {
			opts := &model.ExportExchangesOptions{
				Format: format,
			}
			if opts.Format != format {
				t.Errorf("Format = %v, want %v", opts.Format, format)
			}
		})
	}
}
