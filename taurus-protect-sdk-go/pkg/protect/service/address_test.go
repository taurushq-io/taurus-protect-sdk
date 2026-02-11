package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewAddressService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestAddressService_GetAddress_EmptyID(t *testing.T) {
	svc := &AddressService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetAddress(nil, "")
	if err == nil {
		t.Error("GetAddress() with empty ID should return error")
	}
	if err.Error() != "addressID cannot be empty" {
		t.Errorf("GetAddress() error = %v, want 'addressID cannot be empty'", err)
	}
}

func TestAddressService_CreateAddress_NilRequest(t *testing.T) {
	svc := &AddressService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateAddress(nil, nil)
	if err == nil {
		t.Error("CreateAddress() with nil request should return error")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("CreateAddress() error = %v, want 'request cannot be nil'", err)
	}
}

func TestAddressService_CreateAddress_EmptyWalletID(t *testing.T) {
	svc := &AddressService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateAddress(nil, &model.CreateAddressRequest{})
	if err == nil {
		t.Error("CreateAddress() with empty walletID should return error")
	}
	if err.Error() != "walletID is required" {
		t.Errorf("CreateAddress() error = %v, want 'walletID is required'", err)
	}
}

func TestAddressService_CreateAddress_EmptyLabel(t *testing.T) {
	svc := &AddressService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateAddress(nil, &model.CreateAddressRequest{WalletID: "wallet-123"})
	if err == nil {
		t.Error("CreateAddress() with empty label should return error")
	}
	if err.Error() != "label is required" {
		t.Errorf("CreateAddress() error = %v, want 'label is required'", err)
	}
}
