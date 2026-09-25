package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// TaurusNetworkLendingService provides Taurus Network lending operations.
type TaurusNetworkLendingService struct {
	api       *openapi.TaurusNetworkLendingAPIService
	errMapper *ErrorMapper
}

// NewTaurusNetworkLendingService creates a new TaurusNetworkLendingService.
func NewTaurusNetworkLendingService(client *openapi.APIClient) *TaurusNetworkLendingService {
	return &TaurusNetworkLendingService{
		api:       client.TaurusNetworkLendingAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetLendingAgreement retrieves a lending agreement by ID.
func (s *TaurusNetworkLendingService) GetLendingAgreement(ctx context.Context, lendingAgreementID string) (*taurusnetwork.LendingAgreement, error) {
	if lendingAgreementID == "" {
		return nil, fmt.Errorf("lendingAgreementID cannot be empty")
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceGetLendingAgreement(ctx, lendingAgreementID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.LendingAgreementFromDTO(resp.Result), nil
}

// ListLendingAgreements retrieves a list of lending agreements with optional filtering and pagination.
func (s *TaurusNetworkLendingService) ListLendingAgreements(ctx context.Context, opts *taurusnetwork.ListLendingAgreementsOptions) (*taurusnetwork.ListLendingAgreementsResult, error) {
	if opts == nil {
		opts = &taurusnetwork.ListLendingAgreementsOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.TaurusNetworkServiceGetLendingAgreements(ctx), window)
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.ListLendingAgreementsResult{
		LendingAgreements: mapper.LendingAgreementsFromDTO(resp.LendingAgreements),
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	result.Page = page

	return result, nil
}

// ListLendingAgreementsForApproval retrieves lending agreements pending approval.
func (s *TaurusNetworkLendingService) ListLendingAgreementsForApproval(ctx context.Context, opts *taurusnetwork.ListLendingAgreementsForApprovalOptions) (*taurusnetwork.ListLendingAgreementsForApprovalResult, error) {
	if opts == nil {
		opts = &taurusnetwork.ListLendingAgreementsForApprovalOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.TaurusNetworkServiceGetLendingAgreementsForApproval(ctx), window)
	if len(opts.IDs) > 0 {
		req = req.Ids(opts.IDs)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.ListLendingAgreementsForApprovalResult{
		LendingAgreements: mapper.LendingAgreementsFromDTO(resp.Result),
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	result.Page = page

	return result, nil
}

// CreateLendingAgreement creates a new lending agreement.
func (s *TaurusNetworkLendingService) CreateLendingAgreement(ctx context.Context, req *taurusnetwork.CreateLendingAgreementRequest) (*taurusnetwork.LendingAgreement, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	body := openapi.TgvalidatordCreateLendingAgreementRequest{}

	if req.LendingOfferID != "" {
		body.LendingOfferID = &req.LendingOfferID
	}
	if req.LenderParticipantID != "" {
		body.LenderParticipantID = &req.LenderParticipantID
	}
	if req.CurrencyID != "" {
		body.CurrencyID = &req.CurrencyID
	}
	if req.Amount != "" {
		body.Amount = &req.Amount
	}
	if req.AnnualPercentageYield != "" {
		body.AnnualPercentageYield = &req.AnnualPercentageYield
	}
	if req.Duration != "" {
		body.Duration = &req.Duration
	}
	if req.BorrowerSharedAddressID != "" {
		body.BorrowerSharedAddressID = &req.BorrowerSharedAddressID
	}

	// Convert collaterals
	if len(req.Collaterals) > 0 {
		collaterals := make([]openapi.CreateLendingAgreementRequestLoanCollateralRequest, len(req.Collaterals))
		for i, c := range req.Collaterals {
			collaterals[i] = openapi.CreateLendingAgreementRequestLoanCollateralRequest{
				CurrencyID: stringPtr(c.CurrencyID),
				Amount:     stringPtr(c.Amount),
			}
		}
		body.Collaterals = collaterals
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceCreateLendingAgreement(ctx).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	// CreateLendingAgreement only returns the ID, not the full agreement
	// Return a partial agreement with just the ID
	if resp.LendingAgreementID != nil {
		return &taurusnetwork.LendingAgreement{ID: *resp.LendingAgreementID}, nil
	}
	return &taurusnetwork.LendingAgreement{}, nil
}

// UpdateLendingAgreement updates a lending agreement.
func (s *TaurusNetworkLendingService) UpdateLendingAgreement(ctx context.Context, lendingAgreementID string, req *taurusnetwork.UpdateLendingAgreementRequest) error {
	if lendingAgreementID == "" {
		return fmt.Errorf("lendingAgreementID cannot be empty")
	}
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	body := openapi.TaurusNetworkServiceUpdateLendingAgreementBody{}

	if req.LenderSharedAddressID != "" {
		body.LenderSharedAddressID = &req.LenderSharedAddressID
	}

	_, httpResp, err := s.api.TaurusNetworkServiceUpdateLendingAgreement(ctx, lendingAgreementID).Body(body).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// RepayLendingAgreement repays a lending agreement.
func (s *TaurusNetworkLendingService) RepayLendingAgreement(ctx context.Context, lendingAgreementID string, req *taurusnetwork.RepayLendingAgreementRequest) error {
	if lendingAgreementID == "" {
		return fmt.Errorf("lendingAgreementID cannot be empty")
	}
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	body := openapi.TaurusNetworkServiceRepayLendingAgreementBody{}

	if req.RepayerSharedAddressID != "" {
		body.RepayerSharedAddressID = &req.RepayerSharedAddressID
	}

	_, httpResp, err := s.api.TaurusNetworkServiceRepayLendingAgreement(ctx, lendingAgreementID).Body(body).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// CancelLendingAgreement cancels a lending agreement that is not yet approved by the lender.
func (s *TaurusNetworkLendingService) CancelLendingAgreement(ctx context.Context, lendingAgreementID string) error {
	if lendingAgreementID == "" {
		return fmt.Errorf("lendingAgreementID cannot be empty")
	}

	body := make(map[string]interface{})

	_, httpResp, err := s.api.TaurusNetworkServiceCancelLendingAgreement(ctx, lendingAgreementID).Body(body).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// CreateLendingAgreementAttachment creates an attachment for a lending agreement.
func (s *TaurusNetworkLendingService) CreateLendingAgreementAttachment(ctx context.Context, lendingAgreementID string, req *taurusnetwork.CreateLendingAgreementAttachmentRequest) error {
	if lendingAgreementID == "" {
		return fmt.Errorf("lendingAgreementID cannot be empty")
	}
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	body := openapi.TaurusNetworkServiceCreateLendingAgreementAttachmentBody{}

	if req.Name != "" {
		body.Name = &req.Name
	}
	if req.Type != "" {
		body.Type = &req.Type
	}
	if req.ContentType != "" {
		body.ContentType = &req.ContentType
	}
	if req.Value != "" {
		body.Value = &req.Value
	}

	_, httpResp, err := s.api.TaurusNetworkServiceCreateLendingAgreementAttachment(ctx, lendingAgreementID).Body(body).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// ListLendingAgreementAttachments retrieves attachments for a lending agreement.
func (s *TaurusNetworkLendingService) ListLendingAgreementAttachments(ctx context.Context, lendingAgreementID string) ([]*taurusnetwork.LendingAgreementAttachment, error) {
	if lendingAgreementID == "" {
		return nil, fmt.Errorf("lendingAgreementID cannot be empty")
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceGetLendingAgreementAttachments(ctx, lendingAgreementID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.LendingAgreementAttachmentsFromDTO(resp.Result), nil
}

// GetLendingOffer retrieves a lending offer by ID.
func (s *TaurusNetworkLendingService) GetLendingOffer(ctx context.Context, offerID string) (*taurusnetwork.LendingOffer, error) {
	if offerID == "" {
		return nil, fmt.Errorf("offerID cannot be empty")
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceGetLendingOffer(ctx, offerID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.LendingOfferFromDTO(resp.LendingOffer), nil
}

// ListLendingOffers retrieves a list of lending offers with optional filtering and pagination.
func (s *TaurusNetworkLendingService) ListLendingOffers(ctx context.Context, opts *taurusnetwork.ListLendingOffersOptions) (*taurusnetwork.ListLendingOffersResult, error) {
	if opts == nil {
		opts = &taurusnetwork.ListLendingOffersOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.TaurusNetworkServiceGetLendingOffers(ctx), window)
	if len(opts.CurrencyIDs) > 0 {
		req = req.CurrencyIDsCurrencyIDs(opts.CurrencyIDs)
	}
	if opts.ParticipantID != "" {
		req = req.ParticipantID(opts.ParticipantID)
	}
	if opts.Duration != "" {
		req = req.Duration(opts.Duration)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.ListLendingOffersResult{
		LendingOffers: mapper.LendingOffersFromDTO(resp.LendingOffers),
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	result.Page = page

	return result, nil
}

// CreateLendingOffer creates a new lending offer.
func (s *TaurusNetworkLendingService) CreateLendingOffer(ctx context.Context, req *taurusnetwork.CreateLendingOfferRequest) (*taurusnetwork.LendingOffer, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	body := openapi.TgvalidatordCreateLendingOfferRequest{}

	if req.CurrencyID != "" {
		body.CurrencyID = &req.CurrencyID
	}
	if req.Amount != "" {
		body.Amount = &req.Amount
	}
	if req.AnnualPercentageYield != "" {
		body.AnnualPercentageYield = &req.AnnualPercentageYield
	}
	if req.Duration != "" {
		body.Duration = &req.Duration
	}

	// Convert collateral requirements
	if len(req.CollateralRequirements) > 0 {
		requirements := make([]openapi.CreateLendingOfferRequestCollateralRequirement, len(req.CollateralRequirements))
		for i, r := range req.CollateralRequirements {
			requirements[i] = openapi.CreateLendingOfferRequestCollateralRequirement{
				CurrencyID: stringPtr(r.CurrencyID),
				Ratio:      stringPtr(r.CollateralRatio),
			}
		}
		body.CollateralRequirement = requirements
	}

	resp, httpResp, err := s.api.TaurusNetworkServiceCreateLendingOffer(ctx).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	// CreateLendingOffer only returns the ID, not the full offer
	// Return a partial offer with just the ID
	if resp.OfferID != nil {
		return &taurusnetwork.LendingOffer{ID: *resp.OfferID}, nil
	}
	return &taurusnetwork.LendingOffer{}, nil
}

// DeleteLendingOffer deletes a lending offer by ID.
func (s *TaurusNetworkLendingService) DeleteLendingOffer(ctx context.Context, offerID string) error {
	if offerID == "" {
		return fmt.Errorf("offerID cannot be empty")
	}

	_, httpResp, err := s.api.TaurusNetworkServiceDeleteLendingOffer(ctx, offerID).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// DeleteLendingOffers deletes all lending offers created by the current participant.
func (s *TaurusNetworkLendingService) DeleteLendingOffers(ctx context.Context) error {
	_, httpResp, err := s.api.TaurusNetworkServiceDeleteLendingOffers(ctx).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// stringPtr returns a pointer to the given string, or nil if empty.
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
