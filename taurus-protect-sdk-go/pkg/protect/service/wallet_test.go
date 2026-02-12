package service

import (
	"net/http"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewWalletService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestWalletService_GetWallet_EmptyID(t *testing.T) {
	// Create a service with nil API to test validation before API call
	svc := &WalletService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetWallet(nil, "")
	if err == nil {
		t.Error("GetWallet() with empty ID should return error")
	}
	if err.Error() != "walletID cannot be empty" {
		t.Errorf("GetWallet() error = %v, want 'walletID cannot be empty'", err)
	}
}

func TestWalletService_CreateWallet_NilRequest(t *testing.T) {
	svc := &WalletService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateWallet(nil, nil)
	if err == nil {
		t.Error("CreateWallet() with nil request should return error")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("CreateWallet() error = %v, want 'request cannot be nil'", err)
	}
}

func TestWalletService_CreateWallet_EmptyName(t *testing.T) {
	svc := &WalletService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateWallet(nil, &model.CreateWalletRequest{})
	if err == nil {
		t.Error("CreateWallet() with empty name should return error")
	}
	if err.Error() != "name is required" {
		t.Errorf("CreateWallet() error = %v, want 'name is required'", err)
	}
}

func TestWalletService_CreateWallet_EmptyCurrency(t *testing.T) {
	svc := &WalletService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateWallet(nil, &model.CreateWalletRequest{Name: "Test Wallet"})
	if err == nil {
		t.Error("CreateWallet() with empty currency should return error")
	}
	if err.Error() != "currency is required" {
		t.Errorf("CreateWallet() error = %v, want 'currency is required'", err)
	}
}

func TestNewErrorMapper(t *testing.T) {
	mapper := NewErrorMapper()
	if mapper == nil {
		t.Error("NewErrorMapper() returned nil")
	}
}

func TestErrorMapper_MapError_NilError(t *testing.T) {
	mapper := NewErrorMapper()
	err := mapper.MapError(nil, nil)
	if err != nil {
		t.Errorf("MapError(nil, nil) = %v, want nil", err)
	}
}

func TestErrorMapper_MapError_NonOpenAPIError(t *testing.T) {
	mapper := NewErrorMapper()
	originalErr := &testError{msg: "test error"}
	err := mapper.MapError(originalErr, nil)
	if err != originalErr {
		t.Errorf("MapError() should return original error for non-OpenAPI errors")
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name        string
		apiErr      *APIError
		wantContain string
	}{
		{
			name:        "with description",
			apiErr:      &APIError{Code: 404, Message: "wallet not found", Description: "Not Found"},
			wantContain: "Not Found",
		},
		{
			name:        "without description",
			apiErr:      &APIError{Code: 500, Message: "internal error"},
			wantContain: "internal error",
		},
		{
			name:        "contains code",
			apiErr:      &APIError{Code: 401, Message: "unauthorized", Description: "Unauthorized"},
			wantContain: "401",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.apiErr.Error()
			if !contains(got, tt.wantContain) {
				t.Errorf("APIError.Error() = %v, want to contain %v", got, tt.wantContain)
			}
		})
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	innerErr := &testError{msg: "inner error"}
	apiErr := &APIError{Code: 500, Message: "outer", Err: innerErr}

	if apiErr.Unwrap() != innerErr {
		t.Error("APIError.Unwrap() should return wrapped error")
	}
}

func TestAPIError_Unwrap_Nil(t *testing.T) {
	apiErr := &APIError{Code: 500, Message: "no wrapped error"}

	if apiErr.Unwrap() != nil {
		t.Error("APIError.Unwrap() should return nil when no wrapped error")
	}
}

func TestMapOpenAPIError_StatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantCode   int
		wantDescr  string
	}{
		{"bad request", 400, 400, "Bad Request"},
		{"unauthorized", 401, 401, "Unauthorized"},
		{"forbidden", 403, 403, "Forbidden"},
		{"not found", 404, 404, "Not Found"},
		{"rate limited", 429, 429, "Rate Limited"},
		{"server error 500", 500, 500, "Server Error"},
		{"server error 502", 502, 502, "Server Error"},
		{"server error 503", 503, 503, "Server Error"},
		{"unknown status", 418, 418, ""}, // I'm a teapot - no special description
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			// Create an actual GenericOpenAPIError
			mockErr := &openapi.GenericOpenAPIError{}

			err := mapOpenAPIError(mockErr, resp)
			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("mapOpenAPIError() returned %T, want *APIError", err)
			}
			if apiErr.Code != tt.wantCode {
				t.Errorf("APIError.Code = %v, want %v", apiErr.Code, tt.wantCode)
			}
			if apiErr.Description != tt.wantDescr {
				t.Errorf("APIError.Description = %v, want %v", apiErr.Description, tt.wantDescr)
			}
		})
	}
}

func TestMapOpenAPIError_NilResponse(t *testing.T) {
	mockErr := &openapi.GenericOpenAPIError{}
	err := mapOpenAPIError(mockErr, nil)

	// Should return the original error when response is nil
	if err != mockErr {
		t.Errorf("mapOpenAPIError(err, nil) should return original error")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
