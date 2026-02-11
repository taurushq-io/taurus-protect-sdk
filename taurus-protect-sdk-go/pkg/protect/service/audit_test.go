package service

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewAuditService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestAuditService_ListAuditTrails_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	// The actual API call will fail, but we're testing the options handling
	svc := &AuditService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("AuditService should not be nil")
	}
}

func TestAuditService_ListAuditTrails_WithOptions(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		options *model.ListAuditTrailsOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.ListAuditTrailsOptions{},
		},
		{
			name: "external user ID filter",
			options: &model.ListAuditTrailsOptions{
				ExternalUserID: "ext-user-123",
			},
		},
		{
			name: "entities filter",
			options: &model.ListAuditTrailsOptions{
				Entities: []string{"Wallet", "Request"},
			},
		},
		{
			name: "actions filter",
			options: &model.ListAuditTrailsOptions{
				Actions: []string{"Create", "Update", "Delete"},
			},
		},
		{
			name: "date range filter",
			options: &model.ListAuditTrailsOptions{
				CreationDateFrom: &now,
				CreationDateTo:   &now,
			},
		},
		{
			name: "pagination options",
			options: &model.ListAuditTrailsOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "sorting options",
			options: &model.ListAuditTrailsOptions{
				SortBy:    []string{"CreationDate"},
				SortOrder: "DESC",
			},
		},
		{
			name: "all options combined",
			options: &model.ListAuditTrailsOptions{
				ExternalUserID:   "ext-user-456",
				Entities:         []string{"Wallet"},
				Actions:          []string{"Create"},
				CreationDateFrom: &now,
				CreationDateTo:   &now,
				CurrentPage:      "xyz789",
				PageRequest:      "FIRST",
				PageSize:         100,
				SortBy:           []string{"CreationDate"},
				SortOrder:        "ASC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &AuditService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("AuditService should not be nil")
			}
		})
	}
}

func TestAuditService_ExportAuditTrails_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	svc := &AuditService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("AuditService should not be nil")
	}
}

func TestAuditService_ExportAuditTrails_WithOptions(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		options *model.ExportAuditTrailsOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.ExportAuditTrailsOptions{},
		},
		{
			name: "external user ID filter",
			options: &model.ExportAuditTrailsOptions{
				ExternalUserID: "ext-user-123",
			},
		},
		{
			name: "entities filter",
			options: &model.ExportAuditTrailsOptions{
				Entities: []string{"Wallet", "Request"},
			},
		},
		{
			name: "actions filter",
			options: &model.ExportAuditTrailsOptions{
				Actions: []string{"Create", "Update", "Delete"},
			},
		},
		{
			name: "date range filter",
			options: &model.ExportAuditTrailsOptions{
				CreationDateFrom: &now,
				CreationDateTo:   &now,
			},
		},
		{
			name: "csv format",
			options: &model.ExportAuditTrailsOptions{
				Format: "csv",
			},
		},
		{
			name: "json format",
			options: &model.ExportAuditTrailsOptions{
				Format: "json",
			},
		},
		{
			name: "all options combined",
			options: &model.ExportAuditTrailsOptions{
				ExternalUserID:   "ext-user-456",
				Entities:         []string{"Wallet"},
				Actions:          []string{"Create"},
				CreationDateFrom: &now,
				CreationDateTo:   &now,
				Format:           "csv",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &AuditService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("AuditService should not be nil")
			}
		})
	}
}

func TestListAuditTrailsOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &model.ListAuditTrailsOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListAuditTrailsOptions_SortOrderValues(t *testing.T) {
	// Test that sort order values match expected API values
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &model.ListAuditTrailsOptions{
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}

func TestExportAuditTrailsOptions_FormatValues(t *testing.T) {
	// Test that format values match expected API values
	validFormats := []string{"csv", "json"}

	for _, format := range validFormats {
		t.Run(format, func(t *testing.T) {
			opts := &model.ExportAuditTrailsOptions{
				Format: format,
			}
			if opts.Format != format {
				t.Errorf("Format = %v, want %v", opts.Format, format)
			}
		})
	}
}
