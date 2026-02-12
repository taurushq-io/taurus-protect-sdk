package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model/taurusnetwork"
)

// TaurusNetworkSettlementService provides Taurus-NETWORK settlement operations.
type TaurusNetworkSettlementService struct {
	api       *openapi.TaurusNetworkSettlementAPIService
	errMapper *ErrorMapper
}

// NewTaurusNetworkSettlementService creates a new TaurusNetworkSettlementService.
func NewTaurusNetworkSettlementService(client *openapi.APIClient) *TaurusNetworkSettlementService {
	return &TaurusNetworkSettlementService{
		api:       client.TaurusNetworkSettlementAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetSettlement retrieves a settlement by ID.
func (s *TaurusNetworkSettlementService) GetSettlement(ctx context.Context, settlementID string) (*taurusnetwork.Settlement, error) {
	if settlementID == "" {
		return nil, fmt.Errorf("settlement ID is required")
	}

	req := s.api.TaurusNetworkServiceGetSettlement(ctx, settlementID)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.SettlementFromDTO(resp.Result), nil
}

// ListSettlements retrieves a list of settlements with optional filtering and pagination.
func (s *TaurusNetworkSettlementService) ListSettlements(ctx context.Context, opts *taurusnetwork.ListSettlementsOptions) (*taurusnetwork.ListSettlementsResult, error) {
	req := s.api.TaurusNetworkServiceGetSettlements(ctx)

	if opts != nil {
		if opts.CounterParticipantID != "" {
			req = req.CounterParticipantID(opts.CounterParticipantID)
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
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.ListSettlementsResult{
		Settlements: mapper.SettlementsFromDTO(resp.Result),
	}

	// Parse cursor pagination info
	if resp.Cursor != nil {
		if resp.Cursor.CurrentPage != nil {
			result.CurrentPage = *resp.Cursor.CurrentPage
		}
		if resp.Cursor.HasPrevious != nil {
			result.HasPrevious = *resp.Cursor.HasPrevious
		}
		if resp.Cursor.HasNext != nil {
			result.HasNext = *resp.Cursor.HasNext
		}
	}

	return result, nil
}

// ListSettlementsForApproval retrieves a list of settlements pending approval.
// Required role: RequestApprover.
func (s *TaurusNetworkSettlementService) ListSettlementsForApproval(ctx context.Context, opts *taurusnetwork.ListSettlementsForApprovalOptions) (*taurusnetwork.ListSettlementsForApprovalResult, error) {
	req := s.api.TaurusNetworkServiceGetSettlementsForApproval(ctx)

	if opts != nil {
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
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
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.ListSettlementsForApprovalResult{
		Settlements: mapper.SettlementsFromDTO(resp.Result),
	}

	// Parse cursor pagination info
	if resp.Cursor != nil {
		if resp.Cursor.CurrentPage != nil {
			result.CurrentPage = *resp.Cursor.CurrentPage
		}
		if resp.Cursor.HasPrevious != nil {
			result.HasPrevious = *resp.Cursor.HasPrevious
		}
		if resp.Cursor.HasNext != nil {
			result.HasNext = *resp.Cursor.HasNext
		}
	}

	return result, nil
}

// CreateSettlement creates a new settlement.
func (s *TaurusNetworkSettlementService) CreateSettlement(ctx context.Context, req *taurusnetwork.CreateSettlementRequest) (*taurusnetwork.CreateSettlementResult, error) {
	if req == nil {
		return nil, fmt.Errorf("create settlement request is required")
	}

	body := mapper.CreateSettlementRequestToDTO(req)
	apiReq := s.api.TaurusNetworkServiceCreateSettlement(ctx).Body(body)

	resp, httpResp, err := apiReq.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &taurusnetwork.CreateSettlementResult{}
	if resp.SettlementID != nil {
		result.SettlementID = *resp.SettlementID
	}

	return result, nil
}

// CancelSettlement cancels a settlement.
func (s *TaurusNetworkSettlementService) CancelSettlement(ctx context.Context, settlementID string) error {
	if settlementID == "" {
		return fmt.Errorf("settlement ID is required")
	}

	req := s.api.TaurusNetworkServiceCancelSettlement(ctx, settlementID).Body(map[string]interface{}{})

	_, httpResp, err := req.Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// ReplaceSettlement replaces a settlement with new attributes at the target side before approval.
func (s *TaurusNetworkSettlementService) ReplaceSettlement(ctx context.Context, settlementID string, req *taurusnetwork.ReplaceSettlementRequest) error {
	if settlementID == "" {
		return fmt.Errorf("settlement ID is required")
	}
	if req == nil {
		return fmt.Errorf("replace settlement request is required")
	}

	body := mapper.ReplaceSettlementRequestToDTO(req)
	apiReq := s.api.TaurusNetworkServiceReplaceSettlement(ctx, settlementID).Body(body)

	_, httpResp, err := apiReq.Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
