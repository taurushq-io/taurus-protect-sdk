package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// TaurusNetworkPledgeService provides Taurus Network pledge operations.
type TaurusNetworkPledgeService struct {
	api       *openapi.TaurusNetworkPledgeAPIService
	errMapper *ErrorMapper
}

// NewTaurusNetworkPledgeService creates a new TaurusNetworkPledgeService.
func NewTaurusNetworkPledgeService(client *openapi.APIClient) *TaurusNetworkPledgeService {
	return &TaurusNetworkPledgeService{
		api:       client.TaurusNetworkPledgeAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetPledge retrieves a pledge by ID.
func (s *TaurusNetworkPledgeService) GetPledge(ctx context.Context, pledgeID string) (*taurusnetwork.Pledge, error) {
	if pledgeID == "" {
		return nil, fmt.Errorf("pledgeID cannot be empty")
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceGetPledge(ctx, pledgeID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("pledge not found")
	}

	return mapper.PledgeFromDTO(resp.Result), nil
}

// ListPledges retrieves a list of pledges with optional filters.
func (s *TaurusNetworkPledgeService) ListPledges(ctx context.Context, opts *taurusnetwork.ListPledgesOptions) ([]*taurusnetwork.Pledge, *model.CursorPagination, error) {
	req := s.api.TaurusNetworkServiceGetPledges(ctx)

	if opts != nil {
		if opts.OwnerParticipantID != "" {
			req = req.OwnerParticipantID(opts.OwnerParticipantID)
		}
		if opts.TargetParticipantID != "" {
			req = req.TargetParticipantID(opts.TargetParticipantID)
		}
		if len(opts.SharedAddressIDs) > 0 {
			req = req.SharedAddressIDs(opts.SharedAddressIDs)
		}
		if opts.CurrencyID != "" {
			req = req.CurrencyID(opts.CurrencyID)
		}
		if len(opts.Statuses) > 0 {
			req = req.Statuses(opts.Statuses)
		}
		if opts.SortOrder != "" {
			req = req.SortOrder(opts.SortOrder)
		}
		if opts.CurrentPage != "" {
			req = req.CursorCurrentPage(opts.CurrentPage)
		}
		if opts.PageRequest != "" {
			req = req.CursorPageRequest(opts.PageRequest)
		}
		if opts.PageSize > 0 {
			req = req.CursorPageSize(fmt.Sprintf("%d", opts.PageSize))
		}
		if opts.AttributeFiltersJSON != "" {
			req = req.AttributeFiltersJson(opts.AttributeFiltersJSON)
		}
		if opts.AttributeFiltersOperator != "" {
			req = req.AttributeFiltersOperator(opts.AttributeFiltersOperator)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	pledges := mapper.PledgesFromDTO(resp.Pledges)
	cursor := mapper.CursorPaginationFromDTO(resp.Cursor)

	return pledges, cursor, nil
}

// CreatePledge creates a new pledge.
func (s *TaurusNetworkPledgeService) CreatePledge(ctx context.Context, req *taurusnetwork.CreatePledgeRequest) (*taurusnetwork.CreatePledgeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.SharedAddressID == "" {
		return nil, fmt.Errorf("sharedAddressID is required")
	}
	if req.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}

	createReq := openapi.TgvalidatordCreatePledgeRequest{
		SharedAddressID: &req.SharedAddressID,
		Amount:          &req.Amount,
	}

	if req.CurrencyID != "" {
		createReq.CurrencyID = &req.CurrencyID
	}
	if req.PledgeType != "" {
		createReq.PledgeType = &req.PledgeType
	}
	if req.ExternalReferenceID != "" {
		createReq.ExternalReferenceId = &req.ExternalReferenceID
	}
	if req.ReconciliationNote != "" {
		createReq.ReconciliationNote = &req.ReconciliationNote
	}

	if req.PledgeDurationSetup != nil {
		durationSetup := openapi.CreatePledgeRequestPledgeDurationSetupRequest{}
		if req.PledgeDurationSetup.MinimumDuration != "" {
			durationSetup.MinimumDuration = &req.PledgeDurationSetup.MinimumDuration
		}
		if req.PledgeDurationSetup.EndOfMinimumDurationDate != nil {
			durationSetup.EndOfMinimumDurationDate = req.PledgeDurationSetup.EndOfMinimumDurationDate
		}
		if req.PledgeDurationSetup.NoticePeriodDuration != "" {
			durationSetup.NoticePeriodDuration = &req.PledgeDurationSetup.NoticePeriodDuration
		}
		createReq.PledgeDurationSetup = &durationSetup
	}

	if len(req.KeyValueAttributes) > 0 {
		kvAttrs := make([]openapi.TgvalidatordKeyValue, len(req.KeyValueAttributes))
		for i, kv := range req.KeyValueAttributes {
			kvAttrs[i] = openapi.TgvalidatordKeyValue{
				Key:   &kv.Key,
				Value: &kv.Value,
			}
		}
		createReq.KeyValueAttributes = kvAttrs
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceCreatePledge(ctx).Body(createReq).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &taurusnetwork.CreatePledgeResponse{
		Pledge:         mapper.PledgeFromDTO(resp.Result),
		PledgeActionID: safeStringPtr(resp.PledgeActionID),
	}, nil
}

// UpdatePledge updates a pledge's modifiable fields.
func (s *TaurusNetworkPledgeService) UpdatePledge(ctx context.Context, pledgeID string, req *taurusnetwork.UpdatePledgeRequest) error {
	if pledgeID == "" {
		return fmt.Errorf("pledgeID cannot be empty")
	}
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	updateBody := openapi.TaurusNetworkServiceUpdatePledgeBody{}

	if req.DefaultDestinationSharedAddressID != "" {
		updateBody.DefaultDestinationSharedAddressID = &req.DefaultDestinationSharedAddressID
	}
	if req.DefaultDestinationInternalAddressID != "" {
		updateBody.DefaultDestinationInternalAddressID = &req.DefaultDestinationInternalAddressID
	}

	_, httpResp, err := s.api.TaurusNetworkServiceUpdatePledge(ctx, pledgeID).Body(updateBody).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// AddPledgeCollateral adds collateral to an existing pledge.
func (s *TaurusNetworkPledgeService) AddPledgeCollateral(ctx context.Context, pledgeID string, req *taurusnetwork.AddPledgeCollateralRequest) (*taurusnetwork.AddPledgeCollateralResponse, error) {
	if pledgeID == "" {
		return nil, fmt.Errorf("pledgeID cannot be empty")
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}

	body := openapi.TaurusNetworkServiceAddPledgeCollateralBody{
		Amount: &req.Amount,
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceAddPledgeCollateral(ctx, pledgeID).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &taurusnetwork.AddPledgeCollateralResponse{
		PledgeActionID: safeStringPtr(resp.PledgeActionID),
	}, nil
}

// WithdrawPledge creates a withdrawal request from a pledge (pledgee initiates).
func (s *TaurusNetworkPledgeService) WithdrawPledge(ctx context.Context, pledgeID string, req *taurusnetwork.WithdrawPledgeRequest) (*taurusnetwork.WithdrawPledgeResponse, error) {
	if pledgeID == "" {
		return nil, fmt.Errorf("pledgeID cannot be empty")
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}

	body := openapi.TaurusNetworkServiceWithdrawPledgeBody{
		Amount: &req.Amount,
	}

	if req.DestinationSharedAddressID != "" {
		body.DestinationSharedAddressID = &req.DestinationSharedAddressID
	}
	if req.DestinationInternalAddressID != "" {
		body.DestinationInternalAddressID = &req.DestinationInternalAddressID
	}
	if req.ExternalReferenceID != "" {
		body.ExternalReferenceID = &req.ExternalReferenceID
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceWithdrawPledge(ctx, pledgeID).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &taurusnetwork.WithdrawPledgeResponse{
		PledgeWithdrawalID: safeStringPtr(resp.PledgeWithdrawalID),
		PledgeActionID:     safeStringPtr(resp.PledgeActionID),
	}, nil
}

// InitiateWithdrawPledge initiates a withdrawal from a pledge (pledgor initiates).
func (s *TaurusNetworkPledgeService) InitiateWithdrawPledge(ctx context.Context, pledgeID string, req *taurusnetwork.InitiateWithdrawPledgeRequest) (*taurusnetwork.InitiateWithdrawPledgeResponse, error) {
	if pledgeID == "" {
		return nil, fmt.Errorf("pledgeID cannot be empty")
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}

	body := openapi.TaurusNetworkServiceInitiateWithdrawPledgeBody{
		Amount: req.Amount,
	}

	if req.DestinationSharedAddressID != "" {
		body.DestinationSharedAddressID = &req.DestinationSharedAddressID
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceInitiateWithdrawPledge(ctx, pledgeID).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &taurusnetwork.InitiateWithdrawPledgeResponse{
		PledgeWithdrawalID: safeStringPtr(resp.PledgeWithdrawalID),
		PledgeActionID:     safeStringPtr(resp.PledgeActionID),
	}, nil
}

// Unpledge unpledges funds from a pledge.
func (s *TaurusNetworkPledgeService) Unpledge(ctx context.Context, pledgeID string) (*taurusnetwork.UnpledgeResponse, error) {
	if pledgeID == "" {
		return nil, fmt.Errorf("pledgeID cannot be empty")
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceUnpledge(ctx, pledgeID).Body(map[string]interface{}{}).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &taurusnetwork.UnpledgeResponse{
		PledgeActionID: safeStringPtr(resp.PledgeActionID),
	}, nil
}

// RejectPledge rejects a pledge.
func (s *TaurusNetworkPledgeService) RejectPledge(ctx context.Context, pledgeID string, req *taurusnetwork.RejectPledgeRequest) error {
	if pledgeID == "" {
		return fmt.Errorf("pledgeID cannot be empty")
	}

	body := map[string]interface{}{}
	if req != nil && req.Comment != "" {
		body["comment"] = req.Comment
	}

	_, httpResp, err := s.api.TaurusNetworkServiceRejectPledge(ctx, pledgeID).Body(body).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// ListPledgeActions retrieves a list of pledge actions.
func (s *TaurusNetworkPledgeService) ListPledgeActions(ctx context.Context, opts *taurusnetwork.ListPledgeActionsOptions) ([]taurusnetwork.PledgeAction, *model.CursorPagination, error) {
	req := s.api.TaurusNetworkServiceGetPledgeActions(ctx)

	if opts != nil {
		if opts.PledgeID != "" {
			req = req.PledgeID(opts.PledgeID)
		}
		if len(opts.ActionIDs) > 0 {
			req = req.Ids(opts.ActionIDs)
		}
		if opts.SortOrder != "" {
			req = req.SortOrder(opts.SortOrder)
		}
		if opts.CurrentPage != "" {
			req = req.CursorCurrentPage(opts.CurrentPage)
		}
		if opts.PageRequest != "" {
			req = req.CursorPageRequest(opts.PageRequest)
		}
		if opts.PageSize > 0 {
			req = req.CursorPageSize(fmt.Sprintf("%d", opts.PageSize))
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	actions := mapper.PledgeActionsFromDTO(resp.Result)
	if err := verifyPledgeActionMetadata(actions); err != nil {
		return nil, nil, err
	}
	cursor := mapper.CursorPaginationFromDTO(resp.Cursor)

	return actions, cursor, nil
}

// ListPledgeActionsForApproval retrieves a list of pledge actions pending approval.
func (s *TaurusNetworkPledgeService) ListPledgeActionsForApproval(ctx context.Context, opts *taurusnetwork.ListPledgeActionsForApprovalOptions) ([]taurusnetwork.PledgeAction, *model.CursorPagination, error) {
	req := s.api.TaurusNetworkServiceGetPledgeActionsForApproval(ctx)

	if opts != nil {
		if len(opts.ActionIDs) > 0 {
			req = req.Ids(opts.ActionIDs)
		}
		if len(opts.Types) > 0 {
			req = req.Types(opts.Types)
		}
		if opts.SortOrder != "" {
			req = req.SortOrder(opts.SortOrder)
		}
		if opts.CurrentPage != "" {
			req = req.CursorCurrentPage(opts.CurrentPage)
		}
		if opts.PageRequest != "" {
			req = req.CursorPageRequest(opts.PageRequest)
		}
		if opts.PageSize > 0 {
			req = req.CursorPageSize(fmt.Sprintf("%d", opts.PageSize))
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	actions := mapper.PledgeActionsFromDTO(resp.Result)
	if err := verifyPledgeActionMetadata(actions); err != nil {
		return nil, nil, err
	}
	cursor := mapper.CursorPaginationFromDTO(resp.Cursor)

	return actions, cursor, nil
}

// verifyPledgeActionMetadata verifies every action's hash against its payload on the read
// paths, so a caller reviewing an action reads a payload the hash actually commits to.
//
// PledgeAction reuses model.RequestMetadata, whose documented contract is that metadata the
// SDK returned has been materialised by VerifyAndMaterialise — and ParsePayloadEntries
// silently falls back to parsing the unverified string, so without this an integrator
// displaying an action for review saw content nothing had checked. The approval path verifies
// again in the same call rather than trusting HashVerified: the flag is a serialized field, so
// an action decoded from a queue or cache can arrive claiming true.
func verifyPledgeActionMetadata(actions []taurusnetwork.PledgeAction) error {
	for i := range actions {
		if actions[i].Metadata == nil {
			continue
		}
		if err := actions[i].Metadata.VerifyAndMaterialise(); err != nil {
			return fmt.Errorf("pledge action %s: %w", actions[i].ID, err)
		}
	}
	return nil
}

// ApprovePledgeActions approves one or more pledge actions, verifying every metadata hash it
// is about to attest to and signing inside the SDK.
//
//	actions ─▶ presence check ─▶ sha256(payload) == hash ─▶ sort by id ─▶ sign once ─▶ POST
//
// It used to take an opaque, caller-computed `signature` string with only a non-empty check, so
// the approver's key attested to a hash nothing had verified. A malicious or compromised server
// could return a pledge action whose payloadAsString describes a benign top-up while its hash is
// that of a withdrawal of the whole collateral to an address of the attacker's choosing: the
// approver reviews the payload, signs the hash, and the movement executed is not the one they
// read. Python already did this correctly and is listed `verifies` in
// scripts/signing-sites/manifest.json; Go and TypeScript were the two that did not.
//
// docs/SERVICES.md has documented this signing-style signature for some time while the shipped
// method took a pre-computed blob, so this makes the code match the documentation rather than
// the other way round.
//
// One signature covers the whole batch, so any action that fails these checks aborts the call
// and nothing is signed or submitted.
func (s *TaurusNetworkPledgeService) ApprovePledgeActions(
	ctx context.Context,
	actions []taurusnetwork.PledgeAction,
	privateKey *ecdsa.PrivateKey,
	comment string,
) (*taurusnetwork.ApprovePledgeActionsResponse, error) {
	if len(actions) == 0 {
		return nil, fmt.Errorf("actions cannot be empty")
	}
	if privateKey == nil {
		return nil, fmt.Errorf("privateKey cannot be nil")
	}

	// Verify every hash this signature will attest to, then sort by id: the endpoint requires
	// ascending order and the signed array must not depend on the caller's ordering.
	sorted := make([]taurusnetwork.PledgeAction, 0, len(actions))
	for _, action := range actions {
		if action.Metadata == nil {
			return nil, fmt.Errorf("pledge action %s has no metadata", action.ID)
		}
		if action.Metadata.Hash == "" {
			return nil, fmt.Errorf("pledge action %s has no metadata hash", action.ID)
		}
		if action.Metadata.PayloadAsString == "" {
			return nil, &model.IntegrityError{
				Message: fmt.Sprintf("refusing to sign pledge action %s: hash exists but the "+
					"payload it commits to is missing, so there is nothing to verify it against",
					action.ID),
			}
		}
		computed := crypto.CalculateHexHash(action.Metadata.PayloadAsString)
		if !helper.ConstantTimeCompare(computed, action.Metadata.Hash) {
			return nil, &model.IntegrityError{
				Message: fmt.Sprintf("refusing to sign pledge action %s: hash verification "+
					"failed: computed=%s, provided=%s", action.ID, computed, action.Metadata.Hash),
			}
		}
		sorted = append(sorted, action)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })

	ids := make([]string, 0, len(sorted))
	hashes := make([]string, 0, len(sorted))
	for _, action := range sorted {
		ids = append(ids, action.ID)
		hashes = append(hashes, action.Metadata.Hash)
	}

	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize hashes: %w", err)
	}
	signature, err := crypto.SignData(privateKey, hashesJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to sign pledge action hashes: %w", err)
	}

	body := openapi.TgvalidatordApprovePledgeActionsRequest{
		Ids:       ids,
		Signature: signature,
		Comment:   comment,
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceApprovePledgeActions(ctx).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return &taurusnetwork.ApprovePledgeActionsResponse{
		Signatures: safeStringPtr(resp.Signatures),
	}, nil
}

// RejectPledgeActions rejects one or more pledge actions.
func (s *TaurusNetworkPledgeService) RejectPledgeActions(ctx context.Context, req *taurusnetwork.RejectPledgeActionsRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if len(req.Ids) == 0 {
		return fmt.Errorf("ids cannot be empty")
	}

	body := openapi.TgvalidatordRejectPledgeActionsRequest{
		Ids:     req.Ids,
		Comment: req.Comment,
	}

	_, httpResp, err := s.api.TaurusNetworkServiceRejectPledgeActions(ctx).Body(body).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// ListPledgeWithdrawals retrieves a list of pledge withdrawals.
func (s *TaurusNetworkPledgeService) ListPledgeWithdrawals(ctx context.Context, opts *taurusnetwork.ListPledgeWithdrawalsOptions) ([]taurusnetwork.PledgeWithdrawal, *model.CursorPagination, error) {
	req := s.api.TaurusNetworkServiceGetPledgesWithdrawals(ctx)

	if opts != nil {
		if opts.PledgeID != "" {
			req = req.PledgeID(opts.PledgeID)
		}
		if opts.WithdrawalStatus != "" {
			req = req.WithdrawalStatus(opts.WithdrawalStatus)
		}
		if opts.SortOrder != "" {
			req = req.SortOrder(opts.SortOrder)
		}
		if opts.CurrentPage != "" {
			req = req.CursorCurrentPage(opts.CurrentPage)
		}
		if opts.PageRequest != "" {
			req = req.CursorPageRequest(opts.PageRequest)
		}
		if opts.PageSize > 0 {
			req = req.CursorPageSize(fmt.Sprintf("%d", opts.PageSize))
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, nil, s.errMapper.MapError(err, httpResp)
	}

	withdrawals := mapper.PledgeWithdrawalsFromDTO(resp.Withdrawals)
	cursor := mapper.CursorPaginationFromDTO(resp.Cursor)

	return withdrawals, cursor, nil
}

// safeStringPtr safely dereferences a string pointer.
func safeStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
