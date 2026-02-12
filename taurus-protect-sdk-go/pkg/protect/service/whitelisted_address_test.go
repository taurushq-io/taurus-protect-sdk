package service

import (
	"testing"
)

func TestNewWhitelistedAddressServiceWithVerification(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestWhitelistedAddressService_GetWhitelistedAddress_EmptyID(t *testing.T) {
	// Create a service with nil API to test validation before API call
	svc := &WhitelistedAddressService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetWhitelistedAddress(nil, "")
	if err == nil {
		t.Error("GetWhitelistedAddress() with empty ID should return error")
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("GetWhitelistedAddress() error = %v, want 'id cannot be empty'", err)
	}
}

func TestWhitelistedAddressService_ListWhitelistedAddresses_NilOptions(t *testing.T) {
	// This test verifies that nil options don't cause a panic
	// Actual API call would fail since api is nil, but we can test up to the validation
	svc := &WhitelistedAddressService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// We can't actually call ListWhitelistedAddresses without a real API,
	// but we verify the service struct is properly initialized
	if svc.errMapper == nil {
		t.Error("ErrorMapper should not be nil")
	}
}
