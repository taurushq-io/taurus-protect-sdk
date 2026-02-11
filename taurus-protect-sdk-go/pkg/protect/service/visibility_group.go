package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// VisibilityGroupService provides operations for managing restricted visibility groups.
type VisibilityGroupService struct {
	api       *openapi.RestrictedVisibilityGroupsAPIService
	errMapper *ErrorMapper
}

// NewVisibilityGroupService creates a new VisibilityGroupService.
func NewVisibilityGroupService(client *openapi.APIClient) *VisibilityGroupService {
	return &VisibilityGroupService{
		api:       client.RestrictedVisibilityGroupsAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListVisibilityGroups retrieves a list of all restricted visibility groups.
func (s *VisibilityGroupService) ListVisibilityGroups(ctx context.Context) ([]*model.VisibilityGroup, error) {
	resp, httpResp, err := s.api.UserServiceGetVisibilityGroups(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.VisibilityGroupsFromDTO(resp.Result), nil
}

// GetUsersByVisibilityGroupID retrieves users in a visibility group by its ID.
func (s *VisibilityGroupService) GetUsersByVisibilityGroupID(ctx context.Context, visibilityGroupID string) ([]*model.User, error) {
	if visibilityGroupID == "" {
		return nil, fmt.Errorf("visibilityGroupID cannot be empty")
	}

	resp, httpResp, err := s.api.UserServiceGetUsersByVisibilityGroupID(ctx, visibilityGroupID).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.UsersFromDTO(resp.Result), nil
}
