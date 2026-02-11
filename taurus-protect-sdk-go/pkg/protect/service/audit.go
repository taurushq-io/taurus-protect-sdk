package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// AuditService provides audit trail management operations.
type AuditService struct {
	api       *openapi.AuditAPIService
	errMapper *ErrorMapper
}

// NewAuditService creates a new AuditService.
func NewAuditService(client *openapi.APIClient) *AuditService {
	return &AuditService{
		api:       client.AuditAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListAuditTrails retrieves a list of audit trails with optional filtering and pagination.
func (s *AuditService) ListAuditTrails(ctx context.Context, opts *model.ListAuditTrailsOptions) (*model.ListAuditTrailsResult, error) {
	req := s.api.AuditServiceGetAuditTrails(ctx)

	if opts != nil {
		if opts.ExternalUserID != "" {
			req = req.ExternalUserId(opts.ExternalUserID)
		}
		if len(opts.Entities) > 0 {
			req = req.Entities(opts.Entities)
		}
		if len(opts.Actions) > 0 {
			req = req.Actions(opts.Actions)
		}
		if opts.CreationDateFrom != nil {
			req = req.CreationDateFrom(*opts.CreationDateFrom)
		}
		if opts.CreationDateTo != nil {
			req = req.CreationDateTo(*opts.CreationDateTo)
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
		if len(opts.SortBy) > 0 {
			req = req.SortingSortBy(opts.SortBy)
		}
		if opts.SortOrder != "" {
			req = req.SortingSortOrder(opts.SortOrder)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListAuditTrailsResult{
		AuditTrails: mapper.AuditTrailsFromDTO(resp.Result),
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

// ExportAuditTrails exports audit trails in the specified format (CSV or JSON).
// Note: A maximum of 10000 trails can be exported at a time.
func (s *AuditService) ExportAuditTrails(ctx context.Context, opts *model.ExportAuditTrailsOptions) (*model.ExportAuditTrailsResult, error) {
	req := s.api.AuditServiceExportAuditTrails(ctx)

	if opts != nil {
		if opts.ExternalUserID != "" {
			req = req.ExternalUserId(opts.ExternalUserID)
		}
		if len(opts.Entities) > 0 {
			req = req.Entities(opts.Entities)
		}
		if len(opts.Actions) > 0 {
			req = req.Actions(opts.Actions)
		}
		if opts.CreationDateFrom != nil {
			req = req.CreationDateFrom(*opts.CreationDateFrom)
		}
		if opts.CreationDateTo != nil {
			req = req.CreationDateTo(*opts.CreationDateTo)
		}
		if opts.Format != "" {
			req = req.Format(opts.Format)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ExportAuditTrailsResult{}

	if resp.Result != nil {
		result.Result = *resp.Result
	}

	// Parse total items
	if resp.TotalItems != nil {
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			result.TotalItems = total
		}
	}

	return result, nil
}
