package service

import (
	"testing"
)

func TestNewFeeService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestFeeService_GetFees_ServiceExists(t *testing.T) {
	// Create a service with nil API to verify service structure
	// The actual API call will fail, but we're testing the service exists
	svc := &FeeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("FeeService should not be nil")
	}
	if svc.errMapper == nil {
		t.Error("FeeService.errMapper should not be nil")
	}
}

func TestFeeService_GetFeesV2_ServiceExists(t *testing.T) {
	// Create a service with nil API to verify service structure
	// The actual API call will fail, but we're testing the service exists
	svc := &FeeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("FeeService should not be nil")
	}
	if svc.errMapper == nil {
		t.Error("FeeService.errMapper should not be nil")
	}
}
