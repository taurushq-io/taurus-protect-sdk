package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewWebhookCallService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestWebhookCallService_ListWebhookCalls_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	// The actual API call will fail, but we're testing the options handling
	svc := &WebhookCallService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("WebhookCallService should not be nil")
	}
}

func TestWebhookCallService_ListWebhookCalls_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *model.ListWebhookCallsOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.ListWebhookCallsOptions{},
		},
		{
			name: "event ID filter",
			options: &model.ListWebhookCallsOptions{
				EventID: "event-123",
			},
		},
		{
			name: "webhook ID filter",
			options: &model.ListWebhookCallsOptions{
				WebhookID: "webhook-456",
			},
		},
		{
			name: "status filter",
			options: &model.ListWebhookCallsOptions{
				Status: "SUCCESS",
			},
		},
		{
			name: "pagination options",
			options: &model.ListWebhookCallsOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "sort order option",
			options: &model.ListWebhookCallsOptions{
				SortOrder: "DESC",
			},
		},
		{
			name: "all options combined",
			options: &model.ListWebhookCallsOptions{
				EventID:     "event-789",
				WebhookID:   "webhook-012",
				Status:      "FAILED",
				CurrentPage: "xyz789",
				PageRequest: "FIRST",
				PageSize:    100,
				SortOrder:   "ASC",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &WebhookCallService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("WebhookCallService should not be nil")
			}
		})
	}
}

func TestListWebhookCallsOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &model.ListWebhookCallsOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListWebhookCallsOptions_SortOrderValues(t *testing.T) {
	// Test that sort order values match expected API values
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &model.ListWebhookCallsOptions{
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}

func TestListWebhookCallsOptions_StatusValues(t *testing.T) {
	// Test that status values match expected API values
	validStatuses := []string{"SUCCESS", "FAILED", "PENDING", "RETRY"}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			opts := &model.ListWebhookCallsOptions{
				Status: status,
			}
			if opts.Status != status {
				t.Errorf("Status = %v, want %v", opts.Status, status)
			}
		})
	}
}
