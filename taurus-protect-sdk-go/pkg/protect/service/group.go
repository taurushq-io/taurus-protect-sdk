package service

import (
	"context"

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

// ListGroups retrieves one page of groups. Result.Pagination is never nil; continue with its
// NextOffset until HasMore is false.
func (s *GroupService) ListGroups(ctx context.Context, opts *model.ListGroupsOptions) (*model.ListGroupsResult, error) {
	if opts == nil {
		opts = &model.ListGroupsOptions{}
	}
	window, err := resolveOffsetWindow(opts.Limit, opts.Offset)
	if err != nil {
		return nil, err
	}

	req := applyOffsetWindow(s.api.UserServiceGetGroups(ctx), window)
	if len(opts.IDs) > 0 {
		req = req.Ids(opts.IDs)
	}
	if len(opts.ExternalGroupIDs) > 0 {
		req = req.ExternalGroupIds(opts.ExternalGroupIDs)
	}
	if opts.Query != "" {
		req = req.Query(opts.Query)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	// A technical group can be appended beyond the limit, so the next page starts at
	// offset + min(rows, limit).
	pagination, err := offsetPagination(rulePlusMinRowsLimit, window, len(resp.Result), 0,
		offsetReply{TotalItems: resp.TotalItems})
	if err != nil {
		return nil, err
	}
	groups := mapper.GroupsFromDTO(resp.Result)
	for _, group := range groups {
		// GetGroups computes the flag on each group and its users; see defaultComputedUserFlags.
		if group == nil {
			continue
		}
		group.EnforcedInRules = falseIfNil(group.EnforcedInRules)
		for i := range group.Users {
			group.Users[i].EnforcedInRules = falseIfNil(group.Users[i].EnforcedInRules)
		}
	}
	return &model.ListGroupsResult{
		Groups:     groups,
		Pagination: pagination,
	}, nil
}
