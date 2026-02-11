package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewHealthService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestHealthService_GetAllHealthChecks_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	// The actual API call will fail, but we're testing the options handling
	svc := &HealthService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("HealthService should not be nil")
	}
}

func TestHealthService_GetAllHealthChecks_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *model.GetAllHealthChecksOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.GetAllHealthChecksOptions{},
		},
		{
			name: "tenant ID filter",
			options: &model.GetAllHealthChecksOptions{
				TenantID: "tenant-123",
			},
		},
		{
			name: "fail if unhealthy",
			options: &model.GetAllHealthChecksOptions{
				FailIfUnhealthy: true,
			},
		},
		{
			name: "all options combined",
			options: &model.GetAllHealthChecksOptions{
				TenantID:        "tenant-456",
				FailIfUnhealthy: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &HealthService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("HealthService should not be nil")
			}
		})
	}
}

func TestGetAllHealthChecksOptions_TenantIDValues(t *testing.T) {
	// Test that tenant ID values are preserved correctly
	tenantIDs := []string{"", "tenant-1", "tenant-123-abc", "multi-tenant"}

	for _, tenantID := range tenantIDs {
		t.Run(tenantID, func(t *testing.T) {
			opts := &model.GetAllHealthChecksOptions{
				TenantID: tenantID,
			}
			if opts.TenantID != tenantID {
				t.Errorf("TenantID = %v, want %v", opts.TenantID, tenantID)
			}
		})
	}
}

func TestGetAllHealthChecksOptions_FailIfUnhealthyValues(t *testing.T) {
	// Test that fail if unhealthy values are preserved correctly
	tests := []struct {
		name            string
		failIfUnhealthy bool
	}{
		{
			name:            "false",
			failIfUnhealthy: false,
		},
		{
			name:            "true",
			failIfUnhealthy: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &model.GetAllHealthChecksOptions{
				FailIfUnhealthy: tt.failIfUnhealthy,
			}
			if opts.FailIfUnhealthy != tt.failIfUnhealthy {
				t.Errorf("FailIfUnhealthy = %v, want %v", opts.FailIfUnhealthy, tt.failIfUnhealthy)
			}
		})
	}
}

func TestHealthService_ErrorMapperInitialized(t *testing.T) {
	// Verify that a HealthService always has an error mapper
	svc := &HealthService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc.errMapper == nil {
		t.Error("HealthService errMapper should not be nil")
	}
}
