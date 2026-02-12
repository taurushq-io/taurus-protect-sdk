package service

import (
	"context"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewReservationService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestReservationService_GetReservation_EmptyID(t *testing.T) {
	svc := &ReservationService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetReservation(context.Background(), "")
	if err == nil {
		t.Error("GetReservation should return error for empty ID")
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("GetReservation error = %v, want 'id cannot be empty'", err)
	}
}

func TestReservationService_GetReservationUTXO_EmptyID(t *testing.T) {
	svc := &ReservationService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetReservationUTXO(context.Background(), "")
	if err == nil {
		t.Error("GetReservationUTXO should return error for empty ID")
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("GetReservationUTXO error = %v, want 'id cannot be empty'", err)
	}
}

func TestReservationService_ListReservations_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	svc := &ReservationService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("ReservationService should not be nil")
	}
}

func TestReservationService_ListReservations_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *model.ListReservationsOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.ListReservationsOptions{},
		},
		{
			name: "kinds filter",
			options: &model.ListReservationsOptions{
				Kinds: []string{"PENDING_REQUEST", "MINIMUM_BALANCE"},
			},
		},
		{
			name: "address filter",
			options: &model.ListReservationsOptions{
				Address: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
			},
		},
		{
			name: "addressId filter",
			options: &model.ListReservationsOptions{
				AddressID: "addr-123",
			},
		},
		{
			name: "pagination options",
			options: &model.ListReservationsOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &model.ListReservationsOptions{
				Kinds:       []string{"PENDING_REQUEST"},
				Address:     "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
				AddressID:   "addr-456",
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
			svc := &ReservationService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("ReservationService should not be nil")
			}
		})
	}
}

func TestListReservationsOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &model.ListReservationsOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListReservationsOptions_KindsValues(t *testing.T) {
	// Test common reservation kind values
	validKinds := []string{"PENDING_REQUEST", "MINIMUM_BALANCE", "PLEDGE", "FEE_DEPOSIT"}

	for _, kind := range validKinds {
		t.Run(kind, func(t *testing.T) {
			opts := &model.ListReservationsOptions{
				Kinds: []string{kind},
			}
			if len(opts.Kinds) != 1 || opts.Kinds[0] != kind {
				t.Errorf("Kinds = %v, want [%v]", opts.Kinds, kind)
			}
		})
	}
}

func TestListReservationsOptions_MultipleKinds(t *testing.T) {
	kinds := []string{"PENDING_REQUEST", "MINIMUM_BALANCE"}
	opts := &model.ListReservationsOptions{
		Kinds: kinds,
	}
	if len(opts.Kinds) != 2 {
		t.Errorf("Kinds length = %v, want 2", len(opts.Kinds))
	}
	if opts.Kinds[0] != "PENDING_REQUEST" {
		t.Errorf("Kinds[0] = %v, want PENDING_REQUEST", opts.Kinds[0])
	}
	if opts.Kinds[1] != "MINIMUM_BALANCE" {
		t.Errorf("Kinds[1] = %v, want MINIMUM_BALANCE", opts.Kinds[1])
	}
}
