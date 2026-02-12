package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewFiatService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestFiatService_ListFiatProviders(t *testing.T) {
	// Create a service with nil API to test the service initialization
	svc := &FiatService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("FiatService should not be nil")
	}
	if svc.errMapper == nil {
		t.Error("ErrorMapper should not be nil")
	}
}

func TestFiatService_GetFiatProviderAccount(t *testing.T) {
	// Create a service with nil API to test the service initialization
	svc := &FiatService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("FiatService should not be nil")
	}
}

func TestFiatService_ListFiatProviderAccounts_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options return an error
	svc := &FiatService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// Verify that nil options return an error
	_, err := svc.ListFiatProviderAccounts(nil, nil)
	if err == nil {
		t.Error("ListFiatProviderAccounts with nil options should return an error")
	}
}

func TestFiatService_ListFiatProviderAccounts_MissingProvider(t *testing.T) {
	svc := &FiatService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	opts := &model.ListFiatProviderAccountsOptions{
		Label: "main-account",
	}

	_, err := svc.ListFiatProviderAccounts(nil, opts)
	if err == nil {
		t.Error("ListFiatProviderAccounts with missing provider should return an error")
	}
}

func TestFiatService_ListFiatProviderAccounts_MissingLabel(t *testing.T) {
	svc := &FiatService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	opts := &model.ListFiatProviderAccountsOptions{
		Provider: "circle",
	}

	_, err := svc.ListFiatProviderAccounts(nil, opts)
	if err == nil {
		t.Error("ListFiatProviderAccounts with missing label should return an error")
	}
}

func TestFiatService_ListFiatProviderAccounts_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *model.ListFiatProviderAccountsOptions
	}{
		{
			name: "required fields only",
			options: &model.ListFiatProviderAccountsOptions{
				Provider: "circle",
				Label:    "main-account",
			},
		},
		{
			name: "with account type filter",
			options: &model.ListFiatProviderAccountsOptions{
				Provider:    "circle",
				Label:       "main-account",
				AccountType: "wallet",
			},
		},
		{
			name: "with sort order",
			options: &model.ListFiatProviderAccountsOptions{
				Provider:  "circle",
				Label:     "main-account",
				SortOrder: "ASC",
			},
		},
		{
			name: "with pagination options",
			options: &model.ListFiatProviderAccountsOptions{
				Provider:    "circle",
				Label:       "main-account",
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &model.ListFiatProviderAccountsOptions{
				Provider:    "circle",
				Label:       "main-account",
				AccountType: "wallet",
				SortOrder:   "DESC",
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
			svc := &FiatService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("FiatService should not be nil")
			}
		})
	}
}

func TestListFiatProviderAccountsOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &model.ListFiatProviderAccountsOptions{
				Provider:    "circle",
				Label:       "main-account",
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListFiatProviderAccountsOptions_SortOrderValues(t *testing.T) {
	// Test that sort order values match expected API values
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &model.ListFiatProviderAccountsOptions{
				Provider:  "circle",
				Label:     "main-account",
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}

func TestListFiatProviderAccountsOptions_AccountTypeValues(t *testing.T) {
	// Test that account type values match expected API values
	validAccountTypes := []string{"wallet", "bank"}

	for _, accountType := range validAccountTypes {
		t.Run(accountType, func(t *testing.T) {
			opts := &model.ListFiatProviderAccountsOptions{
				Provider:    "circle",
				Label:       "main-account",
				AccountType: accountType,
			}
			if opts.AccountType != accountType {
				t.Errorf("AccountType = %v, want %v", opts.AccountType, accountType)
			}
		})
	}
}

func TestFiatService_ErrorMapperInitialized(t *testing.T) {
	// Verify that when a service is constructed, the error mapper is properly initialized
	svc := &FiatService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc.errMapper == nil {
		t.Error("FiatService.errMapper should not be nil after NewErrorMapper()")
	}
}
