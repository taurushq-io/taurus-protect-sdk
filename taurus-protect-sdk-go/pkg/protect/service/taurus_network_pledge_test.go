package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

func TestNewTaurusNetworkPledgeService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestTaurusNetworkPledgeService_GetPledge_EmptyID(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetPledge(nil, "")
	if err == nil {
		t.Error("GetPledge should return error for empty pledgeID")
	}
	if err.Error() != "pledgeID cannot be empty" {
		t.Errorf("GetPledge error = %v, want 'pledgeID cannot be empty'", err)
	}
}

func TestTaurusNetworkPledgeService_CreatePledge_NilRequest(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreatePledge(nil, nil)
	if err == nil {
		t.Error("CreatePledge should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("CreatePledge error = %v, want 'request cannot be nil'", err)
	}
}

func TestTaurusNetworkPledgeService_CreatePledge_MissingSharedAddressID(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &taurusnetwork.CreatePledgeRequest{
		Amount: "1000000",
	}
	_, err := svc.CreatePledge(nil, req)
	if err == nil {
		t.Error("CreatePledge should return error for missing sharedAddressID")
	}
	if err.Error() != "sharedAddressID is required" {
		t.Errorf("CreatePledge error = %v, want 'sharedAddressID is required'", err)
	}
}

func TestTaurusNetworkPledgeService_CreatePledge_MissingAmount(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &taurusnetwork.CreatePledgeRequest{
		SharedAddressID: "addr-123",
	}
	_, err := svc.CreatePledge(nil, req)
	if err == nil {
		t.Error("CreatePledge should return error for missing amount")
	}
	if err.Error() != "amount is required" {
		t.Errorf("CreatePledge error = %v, want 'amount is required'", err)
	}
}

func TestTaurusNetworkPledgeService_UpdatePledge_EmptyID(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.UpdatePledge(nil, "", &taurusnetwork.UpdatePledgeRequest{})
	if err == nil {
		t.Error("UpdatePledge should return error for empty pledgeID")
	}
	if err.Error() != "pledgeID cannot be empty" {
		t.Errorf("UpdatePledge error = %v, want 'pledgeID cannot be empty'", err)
	}
}

func TestTaurusNetworkPledgeService_UpdatePledge_NilRequest(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.UpdatePledge(nil, "pledge-123", nil)
	if err == nil {
		t.Error("UpdatePledge should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("UpdatePledge error = %v, want 'request cannot be nil'", err)
	}
}

func TestTaurusNetworkPledgeService_AddPledgeCollateral_EmptyID(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.AddPledgeCollateral(nil, "", &taurusnetwork.AddPledgeCollateralRequest{})
	if err == nil {
		t.Error("AddPledgeCollateral should return error for empty pledgeID")
	}
	if err.Error() != "pledgeID cannot be empty" {
		t.Errorf("AddPledgeCollateral error = %v, want 'pledgeID cannot be empty'", err)
	}
}

func TestTaurusNetworkPledgeService_AddPledgeCollateral_NilRequest(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.AddPledgeCollateral(nil, "pledge-123", nil)
	if err == nil {
		t.Error("AddPledgeCollateral should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("AddPledgeCollateral error = %v, want 'request cannot be nil'", err)
	}
}

func TestTaurusNetworkPledgeService_AddPledgeCollateral_MissingAmount(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.AddPledgeCollateral(nil, "pledge-123", &taurusnetwork.AddPledgeCollateralRequest{})
	if err == nil {
		t.Error("AddPledgeCollateral should return error for missing amount")
	}
	if err.Error() != "amount is required" {
		t.Errorf("AddPledgeCollateral error = %v, want 'amount is required'", err)
	}
}

func TestTaurusNetworkPledgeService_WithdrawPledge_EmptyID(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.WithdrawPledge(nil, "", &taurusnetwork.WithdrawPledgeRequest{})
	if err == nil {
		t.Error("WithdrawPledge should return error for empty pledgeID")
	}
	if err.Error() != "pledgeID cannot be empty" {
		t.Errorf("WithdrawPledge error = %v, want 'pledgeID cannot be empty'", err)
	}
}

func TestTaurusNetworkPledgeService_WithdrawPledge_NilRequest(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.WithdrawPledge(nil, "pledge-123", nil)
	if err == nil {
		t.Error("WithdrawPledge should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("WithdrawPledge error = %v, want 'request cannot be nil'", err)
	}
}

func TestTaurusNetworkPledgeService_WithdrawPledge_MissingAmount(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.WithdrawPledge(nil, "pledge-123", &taurusnetwork.WithdrawPledgeRequest{})
	if err == nil {
		t.Error("WithdrawPledge should return error for missing amount")
	}
	if err.Error() != "amount is required" {
		t.Errorf("WithdrawPledge error = %v, want 'amount is required'", err)
	}
}

func TestTaurusNetworkPledgeService_InitiateWithdrawPledge_EmptyID(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.InitiateWithdrawPledge(nil, "", &taurusnetwork.InitiateWithdrawPledgeRequest{})
	if err == nil {
		t.Error("InitiateWithdrawPledge should return error for empty pledgeID")
	}
	if err.Error() != "pledgeID cannot be empty" {
		t.Errorf("InitiateWithdrawPledge error = %v, want 'pledgeID cannot be empty'", err)
	}
}

func TestTaurusNetworkPledgeService_InitiateWithdrawPledge_NilRequest(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.InitiateWithdrawPledge(nil, "pledge-123", nil)
	if err == nil {
		t.Error("InitiateWithdrawPledge should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("InitiateWithdrawPledge error = %v, want 'request cannot be nil'", err)
	}
}

func TestTaurusNetworkPledgeService_InitiateWithdrawPledge_MissingAmount(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.InitiateWithdrawPledge(nil, "pledge-123", &taurusnetwork.InitiateWithdrawPledgeRequest{})
	if err == nil {
		t.Error("InitiateWithdrawPledge should return error for missing amount")
	}
	if err.Error() != "amount is required" {
		t.Errorf("InitiateWithdrawPledge error = %v, want 'amount is required'", err)
	}
}

func TestTaurusNetworkPledgeService_Unpledge_EmptyID(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.Unpledge(nil, "")
	if err == nil {
		t.Error("Unpledge should return error for empty pledgeID")
	}
	if err.Error() != "pledgeID cannot be empty" {
		t.Errorf("Unpledge error = %v, want 'pledgeID cannot be empty'", err)
	}
}

func TestTaurusNetworkPledgeService_RejectPledge_EmptyID(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RejectPledge(nil, "", nil)
	if err == nil {
		t.Error("RejectPledge should return error for empty pledgeID")
	}
	if err.Error() != "pledgeID cannot be empty" {
		t.Errorf("RejectPledge error = %v, want 'pledgeID cannot be empty'", err)
	}
}

func TestTaurusNetworkPledgeService_ApprovePledgeActions_NilRequest(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.ApprovePledgeActions(nil, nil)
	if err == nil {
		t.Error("ApprovePledgeActions should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("ApprovePledgeActions error = %v, want 'request cannot be nil'", err)
	}
}

func TestTaurusNetworkPledgeService_ApprovePledgeActions_EmptyIDs(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.ApprovePledgeActions(nil, &taurusnetwork.ApprovePledgeActionsRequest{})
	if err == nil {
		t.Error("ApprovePledgeActions should return error for empty ids")
	}
	if err.Error() != "ids cannot be empty" {
		t.Errorf("ApprovePledgeActions error = %v, want 'ids cannot be empty'", err)
	}
}

func TestTaurusNetworkPledgeService_ApprovePledgeActions_MissingSignature(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.ApprovePledgeActions(nil, &taurusnetwork.ApprovePledgeActionsRequest{
		Ids: []string{"action-1"},
	})
	if err == nil {
		t.Error("ApprovePledgeActions should return error for missing signature")
	}
	if err.Error() != "signature is required" {
		t.Errorf("ApprovePledgeActions error = %v, want 'signature is required'", err)
	}
}

func TestTaurusNetworkPledgeService_RejectPledgeActions_NilRequest(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RejectPledgeActions(nil, nil)
	if err == nil {
		t.Error("RejectPledgeActions should return error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("RejectPledgeActions error = %v, want 'request cannot be nil'", err)
	}
}

func TestTaurusNetworkPledgeService_RejectPledgeActions_EmptyIDs(t *testing.T) {
	svc := &TaurusNetworkPledgeService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RejectPledgeActions(nil, &taurusnetwork.RejectPledgeActionsRequest{})
	if err == nil {
		t.Error("RejectPledgeActions should return error for empty ids")
	}
	if err.Error() != "ids cannot be empty" {
		t.Errorf("RejectPledgeActions error = %v, want 'ids cannot be empty'", err)
	}
}

func TestListPledgesOptions_AllFields(t *testing.T) {
	opts := &taurusnetwork.ListPledgesOptions{
		OwnerParticipantID:       "owner-123",
		TargetParticipantID:      "target-456",
		SharedAddressIDs:         []string{"addr-1", "addr-2"},
		CurrencyID:               "ETH",
		Statuses:                 []string{"ACCEPTED_BY_TARGET", "UNPLEDGED"},
		SortOrder:                "DESC",
		CurrentPage:              "abc123",
		PageRequest:              "NEXT",
		PageSize:                 50,
		AttributeFiltersJSON:     `[{"key":"customKey","value":"customValue"}]`,
		AttributeFiltersOperator: "AND",
	}

	// Verify options can be created with all fields
	if opts.OwnerParticipantID != "owner-123" {
		t.Errorf("OwnerParticipantID = %v, want owner-123", opts.OwnerParticipantID)
	}
	if len(opts.SharedAddressIDs) != 2 {
		t.Errorf("SharedAddressIDs len = %v, want 2", len(opts.SharedAddressIDs))
	}
	if len(opts.Statuses) != 2 {
		t.Errorf("Statuses len = %v, want 2", len(opts.Statuses))
	}
}

func TestListPledgeActionsOptions_AllFields(t *testing.T) {
	opts := &taurusnetwork.ListPledgeActionsOptions{
		PledgeID:    "pledge-1",
		ActionIDs:   []string{"action-1", "action-2"},
		SortOrder:   "DESC",
		CurrentPage: "xyz789",
		PageRequest: "FIRST",
		PageSize:    100,
	}

	if opts.PledgeID != "pledge-1" {
		t.Errorf("PledgeID = %v, want pledge-1", opts.PledgeID)
	}
	if len(opts.ActionIDs) != 2 {
		t.Errorf("ActionIDs len = %v, want 2", len(opts.ActionIDs))
	}
	if opts.PageSize != 100 {
		t.Errorf("PageSize = %v, want 100", opts.PageSize)
	}
}

func TestListPledgeWithdrawalsOptions_AllFields(t *testing.T) {
	opts := &taurusnetwork.ListPledgeWithdrawalsOptions{
		PledgeID:         "pledge-1",
		WithdrawalStatus: "PENDING",
		SortOrder:        "ASC",
		CurrentPage:      "page123",
		PageRequest:      "LAST",
		PageSize:         25,
	}

	if opts.PledgeID != "pledge-1" {
		t.Errorf("PledgeID = %v, want pledge-1", opts.PledgeID)
	}
	if opts.WithdrawalStatus != "PENDING" {
		t.Errorf("WithdrawalStatus = %v, want PENDING", opts.WithdrawalStatus)
	}
	if opts.PageRequest != "LAST" {
		t.Errorf("PageRequest = %v, want LAST", opts.PageRequest)
	}
}

func TestCreatePledgeRequest_WithDurationSetup(t *testing.T) {
	req := &taurusnetwork.CreatePledgeRequest{
		SharedAddressID:     "addr-123",
		CurrencyID:          "ETH",
		Amount:              "1000000000000000000",
		PledgeType:          "COLLATERAL",
		ExternalReferenceID: "ext-ref-001",
		ReconciliationNote:  "Test pledge",
		PledgeDurationSetup: &taurusnetwork.CreatePledgeDurationSetup{
			MinimumDuration:      "2592000s",
			NoticePeriodDuration: "172800s",
		},
		KeyValueAttributes: []taurusnetwork.KeyValue{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
	}

	if req.SharedAddressID != "addr-123" {
		t.Errorf("SharedAddressID = %v, want addr-123", req.SharedAddressID)
	}
	if req.PledgeDurationSetup == nil {
		t.Error("PledgeDurationSetup should not be nil")
	}
	if len(req.KeyValueAttributes) != 2 {
		t.Errorf("KeyValueAttributes len = %v, want 2", len(req.KeyValueAttributes))
	}
}

func TestWithdrawPledgeRequest_AllFields(t *testing.T) {
	req := &taurusnetwork.WithdrawPledgeRequest{
		DestinationSharedAddressID:   "shared-addr-123",
		DestinationInternalAddressID: "internal-addr-456",
		Amount:                       "500000000000000000",
		ExternalReferenceID:          "ext-ref-002",
	}

	if req.Amount != "500000000000000000" {
		t.Errorf("Amount = %v, want 500000000000000000", req.Amount)
	}
	if req.DestinationSharedAddressID != "shared-addr-123" {
		t.Errorf("DestinationSharedAddressID = %v, want shared-addr-123", req.DestinationSharedAddressID)
	}
}
