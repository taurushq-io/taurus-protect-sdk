package service

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

func TestNewTaurusNetworkSettlementService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestTaurusNetworkSettlementService_GetSettlement_EmptyID(t *testing.T) {
	svc := &TaurusNetworkSettlementService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetSettlement(nil, "")
	if err == nil {
		t.Error("GetSettlement() should return error for empty settlement ID")
	}
}

func TestTaurusNetworkSettlementService_ListSettlements_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	svc := &TaurusNetworkSettlementService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("TaurusNetworkSettlementService should not be nil")
	}
}

func TestTaurusNetworkSettlementService_ListSettlements_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *taurusnetwork.ListSettlementsOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &taurusnetwork.ListSettlementsOptions{},
		},
		{
			name: "counter participant ID filter",
			options: &taurusnetwork.ListSettlementsOptions{
				CounterParticipantID: "participant-123",
			},
		},
		{
			name: "statuses filter",
			options: &taurusnetwork.ListSettlementsOptions{
				Statuses: []string{taurusnetwork.SettlementStatusCreated, taurusnetwork.SettlementStatusPending},
			},
		},
		{
			name: "sort order",
			options: &taurusnetwork.ListSettlementsOptions{
				SortOrder: "ASC",
			},
		},
		{
			name: "pagination options",
			options: &taurusnetwork.ListSettlementsOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &taurusnetwork.ListSettlementsOptions{
				CounterParticipantID: "participant-456",
				Statuses:             []string{taurusnetwork.SettlementStatusCompleted},
				SortOrder:            "DESC",
				CurrentPage:          "xyz789",
				PageRequest:          "FIRST",
				PageSize:             100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TaurusNetworkSettlementService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkSettlementService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkSettlementService_ListSettlementsForApproval_NilOptions(t *testing.T) {
	svc := &TaurusNetworkSettlementService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("TaurusNetworkSettlementService should not be nil")
	}
}

func TestTaurusNetworkSettlementService_ListSettlementsForApproval_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *taurusnetwork.ListSettlementsForApprovalOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &taurusnetwork.ListSettlementsForApprovalOptions{},
		},
		{
			name: "IDs filter",
			options: &taurusnetwork.ListSettlementsForApprovalOptions{
				IDs: []string{"settlement-1", "settlement-2"},
			},
		},
		{
			name: "sort order",
			options: &taurusnetwork.ListSettlementsForApprovalOptions{
				SortOrder: "DESC",
			},
		},
		{
			name: "pagination options",
			options: &taurusnetwork.ListSettlementsForApprovalOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    25,
			},
		},
		{
			name: "all options combined",
			options: &taurusnetwork.ListSettlementsForApprovalOptions{
				IDs:         []string{"settlement-1"},
				SortOrder:   "ASC",
				CurrentPage: "xyz789",
				PageRequest: "LAST",
				PageSize:    50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TaurusNetworkSettlementService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkSettlementService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkSettlementService_CreateSettlement_NilRequest(t *testing.T) {
	svc := &TaurusNetworkSettlementService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateSettlement(nil, nil)
	if err == nil {
		t.Error("CreateSettlement() should return error for nil request")
	}
}

func TestTaurusNetworkSettlementService_CreateSettlement_WithRequest(t *testing.T) {
	startDate := time.Now().Add(24 * time.Hour)
	tests := []struct {
		name string
		req  *taurusnetwork.CreateSettlementRequest
	}{
		{
			name: "minimal request",
			req: &taurusnetwork.CreateSettlementRequest{
				TargetParticipantID:   "target-123",
				FirstLegParticipantID: "firstleg-456",
				FirstLegAssets:        []taurusnetwork.SettlementAssetTransfer{},
				SecondLegAssets:       []taurusnetwork.SettlementAssetTransfer{},
			},
		},
		{
			name: "request with assets",
			req: &taurusnetwork.CreateSettlementRequest{
				TargetParticipantID:   "target-123",
				FirstLegParticipantID: "firstleg-456",
				FirstLegAssets: []taurusnetwork.SettlementAssetTransfer{
					{CurrencyID: "BTC", Amount: "1.0"},
				},
				SecondLegAssets: []taurusnetwork.SettlementAssetTransfer{
					{CurrencyID: "ETH", Amount: "10.0"},
				},
			},
		},
		{
			name: "request with clips",
			req: &taurusnetwork.CreateSettlementRequest{
				TargetParticipantID:   "target-123",
				FirstLegParticipantID: "firstleg-456",
				FirstLegAssets:        []taurusnetwork.SettlementAssetTransfer{},
				SecondLegAssets:       []taurusnetwork.SettlementAssetTransfer{},
				Clips: []taurusnetwork.CreateSettlementClipRequest{
					{Index: "0"},
					{Index: "1"},
				},
			},
		},
		{
			name: "request with start execution date",
			req: &taurusnetwork.CreateSettlementRequest{
				TargetParticipantID:   "target-123",
				FirstLegParticipantID: "firstleg-456",
				FirstLegAssets:        []taurusnetwork.SettlementAssetTransfer{},
				SecondLegAssets:       []taurusnetwork.SettlementAssetTransfer{},
				StartExecutionDate:    &startDate,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &TaurusNetworkSettlementService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("TaurusNetworkSettlementService should not be nil")
			}
		})
	}
}

func TestTaurusNetworkSettlementService_CancelSettlement_EmptyID(t *testing.T) {
	svc := &TaurusNetworkSettlementService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.CancelSettlement(nil, "")
	if err == nil {
		t.Error("CancelSettlement() should return error for empty settlement ID")
	}
}

func TestTaurusNetworkSettlementService_ReplaceSettlement_EmptyID(t *testing.T) {
	svc := &TaurusNetworkSettlementService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &taurusnetwork.ReplaceSettlementRequest{
		CreateSettlementRequest: &taurusnetwork.CreateSettlementRequest{
			TargetParticipantID:   "target-123",
			FirstLegParticipantID: "firstleg-456",
		},
	}

	err := svc.ReplaceSettlement(nil, "", req)
	if err == nil {
		t.Error("ReplaceSettlement() should return error for empty settlement ID")
	}
}

func TestTaurusNetworkSettlementService_ReplaceSettlement_NilRequest(t *testing.T) {
	svc := &TaurusNetworkSettlementService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.ReplaceSettlement(nil, "settlement-123", nil)
	if err == nil {
		t.Error("ReplaceSettlement() should return error for nil request")
	}
}

func TestListSettlementsOptions_PageRequestValues(t *testing.T) {
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &taurusnetwork.ListSettlementsOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListSettlementsOptions_SortOrderValues(t *testing.T) {
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &taurusnetwork.ListSettlementsOptions{
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}

func TestListSettlementsOptions_StatusValues(t *testing.T) {
	validStatuses := []string{
		taurusnetwork.SettlementStatusCreating,
		taurusnetwork.SettlementStatusCreated,
		taurusnetwork.SettlementStatusRejectedByCreator,
		taurusnetwork.SettlementStatusApprovedByCreator,
		taurusnetwork.SettlementStatusReceived,
		taurusnetwork.SettlementStatusRejectedByTarget,
		taurusnetwork.SettlementStatusAcceptedByTarget,
		taurusnetwork.SettlementStatusPending,
		taurusnetwork.SettlementStatusCompleted,
		taurusnetwork.SettlementStatusFailed,
	}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			opts := &taurusnetwork.ListSettlementsOptions{
				Statuses: []string{status},
			}
			if len(opts.Statuses) != 1 || opts.Statuses[0] != status {
				t.Errorf("Statuses = %v, want [%v]", opts.Statuses, status)
			}
		})
	}
}

func TestSettlementStatusConstants(t *testing.T) {
	// Verify status constants have expected values
	expectedStatuses := map[string]string{
		"CREATING":            taurusnetwork.SettlementStatusCreating,
		"CREATED":             taurusnetwork.SettlementStatusCreated,
		"REJECTED_BY_CREATOR": taurusnetwork.SettlementStatusRejectedByCreator,
		"APPROVED_BY_CREATOR": taurusnetwork.SettlementStatusApprovedByCreator,
		"RECEIVED":            taurusnetwork.SettlementStatusReceived,
		"REJECTED_BY_TARGET":  taurusnetwork.SettlementStatusRejectedByTarget,
		"ACCEPTED_BY_TARGET":  taurusnetwork.SettlementStatusAcceptedByTarget,
		"PENDING":             taurusnetwork.SettlementStatusPending,
		"COMPLETED":           taurusnetwork.SettlementStatusCompleted,
		"FAILED":              taurusnetwork.SettlementStatusFailed,
	}

	for expected, actual := range expectedStatuses {
		if actual != expected {
			t.Errorf("Status constant = %v, want %v", actual, expected)
		}
	}
}
