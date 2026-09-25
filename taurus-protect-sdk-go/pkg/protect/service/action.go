package service

import (
	"context"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ActionService provides action management operations.
type ActionService struct {
	api       *openapi.ActionsAPIService
	errMapper *ErrorMapper
}

// NewActionService creates a new ActionService.
func NewActionService(client *openapi.APIClient) *ActionService {
	return &ActionService{
		api:       client.ActionsAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetAction retrieves a single action by ID.
func (s *ActionService) GetAction(ctx context.Context, id string) (*model.Action, error) {
	req := s.api.ActionServiceGetAction(ctx, id)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.ActionFromDTO(resp.Action), nil
}

// ListActions retrieves one page of actions. Result.Pagination is never nil; continue with its
// NextOffset until HasMore is false.
func (s *ActionService) ListActions(ctx context.Context, opts *model.ListActionsOptions) (*model.ListActionsResult, error) {
	if opts == nil {
		opts = &model.ListActionsOptions{}
	}
	window, err := resolveOffsetWindow(opts.Limit, opts.Offset)
	if err != nil {
		return nil, err
	}

	req := applyOffsetWindow(s.api.ActionServiceGetActions(ctx), window)
	if len(opts.IDs) > 0 {
		req = req.Ids(opts.IDs)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	pagination, err := offsetPagination(rulePlusRows, window, len(resp.Result), 0,
		offsetReply{TotalItems: resp.TotalItems})
	if err != nil {
		return nil, err
	}
	return &model.ListActionsResult{
		Actions:    mapper.ActionsFromDTO(resp.Result),
		Pagination: pagination,
	}, nil
}
