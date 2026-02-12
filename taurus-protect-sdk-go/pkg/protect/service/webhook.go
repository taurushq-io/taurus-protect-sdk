package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WebhookService provides webhook management operations.
type WebhookService struct {
	api       *openapi.WebhooksAPIService
	errMapper *ErrorMapper
}

// NewWebhookService creates a new WebhookService.
func NewWebhookService(client *openapi.APIClient) *WebhookService {
	return &WebhookService{
		api:       client.WebhooksAPI,
		errMapper: NewErrorMapper(),
	}
}

// CreateWebhook creates a new webhook configuration.
// Returns the created webhook with its secret. The secret is only returned once on creation.
func (s *WebhookService) CreateWebhook(ctx context.Context, req *model.CreateWebhookRequest) (*model.CreateWebhookResult, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Type == "" {
		return nil, fmt.Errorf("type is required")
	}
	if req.URL == "" {
		return nil, fmt.Errorf("url is required")
	}

	createReq := openapi.TgvalidatordCreateWebhookRequest{
		Type: &req.Type,
		Url:  &req.URL,
	}

	resp, httpResp, err := s.api.WebhookServiceCreateWebhook(ctx).
		Body(createReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.CreateWebhookResult{
		Webhook: mapper.WebhookFromDTO(resp.Webhook),
	}

	if resp.Secret != nil {
		result.Secret = *resp.Secret
		// Also set the secret on the webhook for convenience
		if result.Webhook != nil {
			result.Webhook.Secret = *resp.Secret
		}
	}

	return result, nil
}

// DeleteWebhook deletes a webhook configuration by ID.
func (s *WebhookService) DeleteWebhook(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	_, httpResp, err := s.api.WebhookServiceDeleteWebhook(ctx, id).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// ListWebhooks retrieves a list of webhooks with optional filtering and pagination.
func (s *WebhookService) ListWebhooks(ctx context.Context, opts *model.ListWebhooksOptions) (*model.ListWebhooksResult, error) {
	req := s.api.WebhookServiceGetWebhooks(ctx)

	if opts != nil {
		if opts.Type != "" {
			req = req.Type_(opts.Type)
		}
		if opts.URL != "" {
			req = req.Url(opts.URL)
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

	result := &model.ListWebhooksResult{
		Webhooks: mapper.WebhooksFromDTO(resp.Webhooks),
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

// UpdateWebhookStatus updates a webhook's status (enabled/disabled).
func (s *WebhookService) UpdateWebhookStatus(ctx context.Context, id string, enabled bool) (*model.Webhook, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	status := "disabled"
	if enabled {
		status = "enabled"
	}

	updateReq := openapi.WebhookServiceUpdateWebhookStatusBody{
		Status: &status,
	}

	resp, httpResp, err := s.api.WebhookServiceUpdateWebhookStatus(ctx, id).
		Body(updateReq).
		Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.WebhookFromDTO(resp.Webhook), nil
}
