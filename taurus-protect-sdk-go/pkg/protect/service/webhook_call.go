package service

import (
	"context"

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
	if opts == nil {
		opts = &model.ListWebhookCallsOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.WebhookServiceGetWebhookCalls(ctx), window)
	if opts.EventID != "" {
		req = req.EventID(opts.EventID)
	}
	if opts.WebhookID != "" {
		req = req.WebhookID(opts.WebhookID)
	}
	if opts.Status != "" {
		req = req.Status(opts.Status)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListWebhookCallsResult{
		WebhookCalls: mapper.WebhookCallsFromDTO(resp.Calls),
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	result.Page = page

	return result, nil
}
