package service

import (
	"context"
	"fmt"
	"strconv"

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

// ListActions retrieves a list of actions with optional filtering and pagination.
func (s *ActionService) ListActions(ctx context.Context, opts *model.ListActionsOptions) (*model.ListActionsResult, error) {
	req := s.api.ActionServiceGetActions(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListActionsResult{
		Actions: mapper.ActionsFromDTO(resp.Result),
	}

	// Parse total items
	if resp.TotalItems != nil {
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			result.TotalItems = total
		}
	}

	return result, nil
}
