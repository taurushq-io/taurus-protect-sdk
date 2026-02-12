package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewRequestService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestRequestService_GetRequest_EmptyID(t *testing.T) {
	svc := &RequestService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetRequest(nil, "")
	if err == nil {
		t.Error("GetRequest() with empty ID should return error")
	}
	if err.Error() != "requestID cannot be empty" {
		t.Errorf("GetRequest() error = %v, want 'requestID cannot be empty'", err)
	}
}

func TestRequestService_CreateOutgoingRequest_NilRequest(t *testing.T) {
	svc := &RequestService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateOutgoingRequest(nil, nil)
	if err == nil {
		t.Error("CreateOutgoingRequest() with nil request should return error")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("CreateOutgoingRequest() error = %v, want 'request cannot be nil'", err)
	}
}

func TestRequestService_CreateOutgoingRequest_EmptyAmount(t *testing.T) {
	svc := &RequestService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateOutgoingRequest(nil, &model.CreateOutgoingRequest{})
	if err == nil {
		t.Error("CreateOutgoingRequest() with empty amount should return error")
	}
	if err.Error() != "amount is required" {
		t.Errorf("CreateOutgoingRequest() error = %v, want 'amount is required'", err)
	}
}

func TestRequestService_CreateOutgoingRequest_NoFromAddress(t *testing.T) {
	svc := &RequestService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateOutgoingRequest(nil, &model.CreateOutgoingRequest{Amount: "100"})
	if err == nil {
		t.Error("CreateOutgoingRequest() without fromAddressID or fromWalletID should return error")
	}
	if err.Error() != "either fromAddressID or fromWalletID is required" {
		t.Errorf("CreateOutgoingRequest() error = %v, want 'either fromAddressID or fromWalletID is required'", err)
	}
}

func TestRequestService_CreateOutgoingRequest_NoToAddress(t *testing.T) {
	svc := &RequestService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateOutgoingRequest(nil, &model.CreateOutgoingRequest{
		Amount:        "100",
		FromAddressID: "addr-123",
	})
	if err == nil {
		t.Error("CreateOutgoingRequest() without toAddressID or toWhitelistedAddressID should return error")
	}
	if err.Error() != "either toAddressID or toWhitelistedAddressID is required" {
		t.Errorf("CreateOutgoingRequest() error = %v, want 'either toAddressID or toWhitelistedAddressID is required'", err)
	}
}

func TestRequestService_CreateOutgoingRequest_ValidWithFromWalletID(t *testing.T) {
	// Test that FromWalletID is an acceptable alternative to FromAddressID
	svc := &RequestService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateOutgoingRequest(nil, &model.CreateOutgoingRequest{
		Amount:       "100",
		FromWalletID: "wallet-123",
	})
	// Should fail on ToAddress validation, not FromAddress
	if err == nil {
		t.Error("CreateOutgoingRequest() should still require ToAddress")
	}
	if err.Error() != "either toAddressID or toWhitelistedAddressID is required" {
		t.Errorf("CreateOutgoingRequest() error = %v, want 'either toAddressID or toWhitelistedAddressID is required'", err)
	}
}

func TestRequestService_CreateOutgoingRequest_ValidWithToWhitelistedAddressID(t *testing.T) {
	// This test would require a real API client, so we can only test validation
	// The validation should pass with valid inputs, but the API call would fail with nil api
	svc := &RequestService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This will panic because api is nil - we can't test full valid flow without mock
	// Just documenting that the validation passes for valid inputs
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to nil API, but didn't get one")
		}
	}()

	_, _ = svc.CreateOutgoingRequest(nil, &model.CreateOutgoingRequest{
		Amount:                 "100",
		FromAddressID:          "addr-123",
		ToWhitelistedAddressID: "whitelist-456",
	})
}

func TestVerifyRequestHash_NilMetadata(t *testing.T) {
	r := &model.Request{ID: "1"}
	err := verifyRequestHash(r)
	if err != nil {
		t.Errorf("verifyRequestHash() with nil metadata should return nil, got %v", err)
	}
}

func TestVerifyRequestHash_EmptyHashAndPayload(t *testing.T) {
	r := &model.Request{
		ID:       "1",
		Metadata: &model.RequestMetadata{Hash: "", PayloadAsString: ""},
	}
	err := verifyRequestHash(r)
	if err != nil {
		t.Errorf("verifyRequestHash() with empty hash and payload should return nil, got %v", err)
	}
}

func TestVerifyRequestHash_EmptyProvidedHash(t *testing.T) {
	r := &model.Request{
		ID: "1",
		Metadata: &model.RequestMetadata{
			Hash:            "",
			PayloadAsString: "test-payload",
		},
	}
	err := verifyRequestHash(r)
	if err == nil {
		t.Fatal("verifyRequestHash() with empty provided hash should return error")
	}
	var intErr *model.IntegrityError
	if !errors.As(err, &intErr) {
		t.Fatalf("verifyRequestHash() error should be IntegrityError, got %T", err)
	}
	if !strings.Contains(intErr.Message, "non-empty") {
		t.Errorf("error message should mention non-empty, got %q", intErr.Message)
	}
}

func TestVerifyRequestHash_ValidHash(t *testing.T) {
	// SHA-256("test-payload") = 6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273
	r := &model.Request{
		ID: "1",
		Metadata: &model.RequestMetadata{
			Hash:            "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273",
			PayloadAsString: "test-payload",
		},
	}
	err := verifyRequestHash(r)
	if err != nil {
		t.Errorf("verifyRequestHash() with valid hash should return nil, got %v", err)
	}
}

func TestVerifyRequestHash_MismatchedHash(t *testing.T) {
	r := &model.Request{
		ID: "1",
		Metadata: &model.RequestMetadata{
			Hash:            "0000000000000000000000000000000000000000000000000000000000000000",
			PayloadAsString: "test-payload",
		},
	}
	err := verifyRequestHash(r)
	if err == nil {
		t.Fatal("verifyRequestHash() with mismatched hash should return error")
	}
	var intErr *model.IntegrityError
	if !errors.As(err, &intErr) {
		t.Fatalf("verifyRequestHash() error should be IntegrityError, got %T", err)
	}
	if !strings.Contains(intErr.Message, "request hash verification failed") {
		t.Errorf("error message should mention verification failed, got %q", intErr.Message)
	}
}
