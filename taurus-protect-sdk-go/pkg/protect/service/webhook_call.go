package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WebhookCallService provides webhook call management operations.
type WebhookCallService struct {
	api       *openapi.WebhookCallsAPIService
	errMapper *ErrorMapper
}

// NewWebhookCallService creates a new WebhookCallService.
func NewWebhookCallService(client *openapi.APIClient) *WebhookCallService {
	return &WebhookCallService{
		api:       client.WebhookCallsAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListWebhookCalls retrieves a list of webhook calls with optional filtering and pagination.
func (s *WebhookCallService) ListWebhookCalls(ctx context.Context, opts *model.ListWebhookCallsOptions) (*model.ListWebhookCallsResult, error) {
	req := s.api.WebhookServiceGetWebhookCalls(ctx)

	if opts != nil {
		if opts.EventID != "" {
			req = req.EventID(opts.EventID)
		}
		if opts.WebhookID != "" {
			req = req.WebhookID(opts.WebhookID)
		}
		if opts.Status != "" {
			req = req.Status(opts.Status)
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
		if opts.SortOrder != "" {
			req = req.SortOrder(opts.SortOrder)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListWebhookCallsResult{
		WebhookCalls: mapper.WebhookCallsFromDTO(resp.Calls),
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
