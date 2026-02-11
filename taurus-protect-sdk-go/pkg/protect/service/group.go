package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// GroupService provides group management operations.
type GroupService struct {
	api       *openapi.GroupsAPIService
	errMapper *ErrorMapper
}

// NewGroupService creates a new GroupService.
func NewGroupService(client *openapi.APIClient) *GroupService {
	return &GroupService{
		api:       client.GroupsAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListGroups retrieves a list of groups with optional filtering and pagination.
func (s *GroupService) ListGroups(ctx context.Context, opts *model.ListGroupsOptions) (*model.ListGroupsResult, error) {
	req := s.api.UserServiceGetGroups(ctx)

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
		if len(opts.ExternalGroupIDs) > 0 {
			req = req.ExternalGroupIds(opts.ExternalGroupIDs)
		}
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListGroupsResult{
		Groups: mapper.GroupsFromDTO(resp.Result),
	}

	// Parse total items
	if resp.TotalItems != nil {
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			result.TotalItems = total
		}
	}

	// Set offset from options if provided
	if opts != nil {
		result.Offset = opts.Offset
	}

	return result, nil
}
