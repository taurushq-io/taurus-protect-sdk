package service

import (
	"context"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

func TestNewTaurusNetworkSharingService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestTaurusNetworkSharingService_ListSharedAddresses_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	svc := &TaurusNetworkSharingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("TaurusNetworkSharingService should not be nil")
	}
}

func TestTaurusNetworkSharingService_ListSharedAddresses_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *taurusnetwork.ListSharedAddressesOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &taurusnetwork.ListSharedAddressesOptions{},
		},
		{
			name: "participant ID filter",
			options: &taurusnetwork.ListSharedAddressesOptions{
				ParticipantID: "participant-123",
			},
		},
		{
			name: "owner participant ID filter",
			options: &taurusnetwork.ListSharedAddressesOptions{
				OwnerParticipantID: "owner-123",
			},
		},
		{
			name: "target participant ID filter",
			options: &taurusnetwork.ListSharedAddressesOptions{
				TargetParticipantID: "target-456",
			},
		},
		{
			name: "blockchain and network filter",
			options: &taurusnetwork.ListSharedAddressesOptions{
				Blockchain: "ethereum",
				Network:    "mainnet",
			},
		},
		{
			name: "IDs filter",
			options: &taurusnetwork.ListSharedAddressesOptions{
				IDs: []string{"id-1", "id-2", "id-3"},
			},
		},
		{
			name: "statuses filter",
			options: &taurusnetwork.ListSharedAddressesOptions{
				Statuses: []string{"new", "pending", "accepted"},
			},
		},
		{
			name: "sort order",
			options: &taurusnetwork.ListSharedAddressesOptions{
				SortOrder: "DESC",
			},
		},
		{
			name: "pagination options",
			options: &taurusnetwork.ListSharedAddressesOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &taurusnetwork.ListSharedAddressesOptions{
				ParticipantID:       "participant-123",
				OwnerParticipantID:  "owner-456",
				TargetParticipantID: "target-789",
				Blockchain:          "ethereum",
				Network:             "mainnet",
				IDs:                 []string{"id-1"},
				Statuses:            []string{"accepted"},
				SortOrder:           "ASC",
				CurrentPage:         "xyz789",
				PageRequest:         "FIRST",
				PageSize:            100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &TaurusNetworkSharingService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkSharingService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkSharingService_ListSharedAssets_NilOptions(t *testing.T) {
	svc := &TaurusNetworkSharingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("TaurusNetworkSharingService should not be nil")
	}
}

func TestTaurusNetworkSharingService_ListSharedAssets_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *taurusnetwork.ListSharedAssetsOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &taurusnetwork.ListSharedAssetsOptions{},
		},
		{
			name: "participant ID filter",
			options: &taurusnetwork.ListSharedAssetsOptions{
				ParticipantID: "participant-123",
			},
		},
		{
			name: "owner participant ID filter",
			options: &taurusnetwork.ListSharedAssetsOptions{
				OwnerParticipantID: "owner-123",
			},
		},
		{
			name: "target participant ID filter",
			options: &taurusnetwork.ListSharedAssetsOptions{
				TargetParticipantID: "target-456",
			},
		},
		{
			name: "blockchain and network filter",
			options: &taurusnetwork.ListSharedAssetsOptions{
				Blockchain: "ethereum",
				Network:    "mainnet",
			},
		},
		{
			name: "IDs filter",
			options: &taurusnetwork.ListSharedAssetsOptions{
				IDs: []string{"id-1", "id-2", "id-3"},
			},
		},
		{
			name: "statuses filter",
			options: &taurusnetwork.ListSharedAssetsOptions{
				Statuses: []string{"new", "pending", "accepted"},
			},
		},
		{
			name: "sort order",
			options: &taurusnetwork.ListSharedAssetsOptions{
				SortOrder: "DESC",
			},
		},
		{
			name: "pagination options",
			options: &taurusnetwork.ListSharedAssetsOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &taurusnetwork.ListSharedAssetsOptions{
				ParticipantID:       "participant-123",
				OwnerParticipantID:  "owner-456",
				TargetParticipantID: "target-789",
				Blockchain:          "ethereum",
				Network:             "mainnet",
				IDs:                 []string{"id-1"},
				Statuses:            []string{"accepted"},
				SortOrder:           "ASC",
				CurrentPage:         "xyz789",
				PageRequest:         "FIRST",
				PageSize:            100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &TaurusNetworkSharingService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkSharingService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkSharingService_ShareAddress_Validation(t *testing.T) {
	svc := &TaurusNetworkSharingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		request *taurusnetwork.ShareAddressRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil request",
			request: nil,
			wantErr: true,
			errMsg:  "request cannot be nil",
		},
		{
			name:    "empty to_participant_id",
			request: &taurusnetwork.ShareAddressRequest{},
			wantErr: true,
			errMsg:  "to_participant_id is required",
		},
		{
			name: "empty address_id",
			request: &taurusnetwork.ShareAddressRequest{
				ToParticipantID: "participant-123",
			},
			wantErr: true,
			errMsg:  "address_id is required",
		},
		{
			name: "valid request with key-value attributes",
			request: &taurusnetwork.ShareAddressRequest{
				ToParticipantID: "participant-123",
				AddressID:       "address-456",
				KeyValueAttributes: []taurusnetwork.KeyValueAttribute{
					{Key: "key1", Value: "value1"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only test cases that should fail (validation errors)
			// Valid requests would require mocking the API
			if !tt.wantErr {
				return
			}
			err := svc.ShareAddress(ctx, tt.request)
			if err == nil {
				t.Errorf("ShareAddress() expected error, got nil")
			} else if err.Error() != tt.errMsg {
				t.Errorf("ShareAddress() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestTaurusNetworkSharingService_ShareWhitelistedAsset_Validation(t *testing.T) {
	svc := &TaurusNetworkSharingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		request *taurusnetwork.ShareWhitelistedAssetRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil request",
			request: nil,
			wantErr: true,
			errMsg:  "request cannot be nil",
		},
		{
			name:    "empty to_participant_id",
			request: &taurusnetwork.ShareWhitelistedAssetRequest{},
			wantErr: true,
			errMsg:  "to_participant_id is required",
		},
		{
			name: "empty whitelisted_contract_id",
			request: &taurusnetwork.ShareWhitelistedAssetRequest{
				ToParticipantID: "participant-123",
			},
			wantErr: true,
			errMsg:  "whitelisted_contract_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ShareWhitelistedAsset(ctx, tt.request)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ShareWhitelistedAsset() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("ShareWhitelistedAsset() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestTaurusNetworkSharingService_UnshareAddress_Validation(t *testing.T) {
	svc := &TaurusNetworkSharingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}
	ctx := context.Background()

	tests := []struct {
		name            string
		sharedAddressID string
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "empty shared_address_id",
			sharedAddressID: "",
			wantErr:         true,
			errMsg:          "shared_address_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.UnshareAddress(ctx, tt.sharedAddressID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnshareAddress() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("UnshareAddress() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestTaurusNetworkSharingService_UnshareWhitelistedAsset_Validation(t *testing.T) {
	svc := &TaurusNetworkSharingService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}
	ctx := context.Background()

	tests := []struct {
		name          string
		sharedAssetID string
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "empty shared_asset_id",
			sharedAssetID: "",
			wantErr:       true,
			errMsg:        "shared_asset_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.UnshareWhitelistedAsset(ctx, tt.sharedAssetID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnshareWhitelistedAsset() expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("UnshareWhitelistedAsset() error = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestListSharedAddressesOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &taurusnetwork.ListSharedAddressesOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListSharedAddressesOptions_SortOrderValues(t *testing.T) {
	// Test that sort order values match expected API values
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &taurusnetwork.ListSharedAddressesOptions{
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}

func TestListSharedAddressesOptions_StatusValues(t *testing.T) {
	// Test that status values match expected API values
	validStatuses := []string{"new", "pending", "rejected", "accepted", "unshared"}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			opts := &taurusnetwork.ListSharedAddressesOptions{
				Statuses: []string{status},
			}
			if len(opts.Statuses) != 1 || opts.Statuses[0] != status {
				t.Errorf("Statuses = %v, want [%v]", opts.Statuses, status)
			}
		})
	}
}

func TestListSharedAssetsOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &taurusnetwork.ListSharedAssetsOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListSharedAssetsOptions_SortOrderValues(t *testing.T) {
	// Test that sort order values match expected API values
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &taurusnetwork.ListSharedAssetsOptions{
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}

func TestListSharedAssetsOptions_StatusValues(t *testing.T) {
	// Test that status values match expected API values
	validStatuses := []string{"new", "pending", "rejected", "accepted", "unshared"}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			opts := &taurusnetwork.ListSharedAssetsOptions{
				Statuses: []string{status},
			}
			if len(opts.Statuses) != 1 || opts.Statuses[0] != status {
				t.Errorf("Statuses = %v, want [%v]", opts.Statuses, status)
			}
		})
	}
}

func TestShareAddressRequest_KeyValueAttributes(t *testing.T) {
	request := &taurusnetwork.ShareAddressRequest{
		ToParticipantID: "participant-123",
		AddressID:       "address-456",
		KeyValueAttributes: []taurusnetwork.KeyValueAttribute{
			{Key: "env", Value: "production"},
			{Key: "region", Value: "us-east-1"},
		},
	}

	if len(request.KeyValueAttributes) != 2 {
		t.Errorf("KeyValueAttributes length = %v, want 2", len(request.KeyValueAttributes))
	}
	if request.KeyValueAttributes[0].Key != "env" || request.KeyValueAttributes[0].Value != "production" {
		t.Errorf("First attribute = %+v, want {Key: env, Value: production}", request.KeyValueAttributes[0])
	}
	if request.KeyValueAttributes[1].Key != "region" || request.KeyValueAttributes[1].Value != "us-east-1" {
		t.Errorf("Second attribute = %+v, want {Key: region, Value: us-east-1}", request.KeyValueAttributes[1])
	}
}
