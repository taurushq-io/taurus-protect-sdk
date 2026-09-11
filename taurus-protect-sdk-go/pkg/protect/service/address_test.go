package service

import (
	"context"
	"errors"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"strings"
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

	_, err := svc.GetAddress(context.TODO(), "")
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

	_, err := svc.CreateAddress(context.TODO(), nil)
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

	_, err := svc.CreateAddress(context.TODO(), &model.CreateAddressRequest{})
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

// T3: every path that returns an Address verifies its HSM signature. CreateAddress was the
// remaining one that did not — and it is the highest-value moment for substitution, because
// its caller is about to publish or fund a fresh deposit address.
//
// These drive verifiedAddress directly. The seam is what the fix is: a per-path check would
// have been fixed once and missed again, which is exactly what happened when
// AssetService.GetAssetAddresses was hardened and CreateAddress was not.
func TestVerifiedAddress_RefusesAnAddressWithNoSignature(t *testing.T) {
	svc := &AddressService{errMapper: NewErrorMapper()}

	addr := "0xATTACKERCONTROLLED"
	id := "42"
	status := "created"
	_, err := svc.verifiedAddress(context.Background(), &openapi.TgvalidatordAddress{
		Id:      &id,
		Address: &addr,
		Status:  &status,
		// No Signature: the server produced an address string it did not sign.
	})
	if err == nil {
		t.Fatal("returned an unverified address string as an ordinary Address")
	}
	var integrityErr *model.IntegrityError
	if !errors.As(err, &integrityErr) {
		t.Errorf("want an IntegrityError so callers can tell this from a transport fault, got %T", err)
	}
	if !strings.Contains(err.Error(), "no HSM") {
		t.Errorf("error should name the missing signature, got %q", err)
	}
}

// The asynchronous case: `creating` with no address yet is legitimate and must not be an error,
// because there is no destination to misuse. This is why the seam checks the address string
// rather than the status.
func TestVerifiedAddress_AllowsAPendingAddressWithNoAddressString(t *testing.T) {
	svc := &AddressService{errMapper: NewErrorMapper()}

	id := "42"
	status := "creating"
	got, err := svc.verifiedAddress(context.Background(), &openapi.TgvalidatordAddress{
		Id:     &id,
		Status: &status,
	})
	if err != nil {
		t.Fatalf("a not-yet-generated address is not an integrity failure: %v", err)
	}
	if got == nil || got.Address != "" {
		t.Errorf("want an Address carrying no destination, got %+v", got)
	}
	if got.Status != "creating" {
		t.Errorf("Status = %q, want it surfaced so the caller knows to re-read", got.Status)
	}
}
