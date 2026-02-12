package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// UserService provides user management operations.
type UserService struct {
	api       *openapi.UsersAPIService
	errMapper *ErrorMapper
}

// NewUserService creates a new UserService.
func NewUserService(client *openapi.APIClient) *UserService {
	return &UserService{
		api:       client.UsersAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetMe retrieves the currently authenticated user.
func (s *UserService) GetMe(ctx context.Context) (*model.User, error) {
	resp, httpResp, err := s.api.UserServiceGetMe(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("user not found")
	}

	return mapper.UserFromDTO(resp.Result), nil
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	resp, httpResp, err := s.api.UserServiceGetUser(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("user not found")
	}

	return mapper.UserFromDTO(resp.Result), nil
}

// ListUsers retrieves a list of users with optional filtering.
func (s *UserService) ListUsers(ctx context.Context, opts *model.ListUsersOptions) (*model.ListUsersResult, error) {
	req := s.api.UserServiceGetUsers(ctx)

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
		if len(opts.ExternalUserIDs) > 0 {
			req = req.ExternalUserIds(opts.ExternalUserIDs)
		}
		if len(opts.Emails) > 0 {
			req = req.Emails(opts.Emails)
		}
		if len(opts.Roles) > 0 {
			req = req.Roles(opts.Roles)
		}
		if len(opts.GroupIDs) > 0 {
			req = req.GroupIds(opts.GroupIDs)
		}
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
		if opts.Status != "" {
			req = req.Status(opts.Status)
		}
		if opts.TotpEnabled != nil {
			req = req.TotpEnabled(*opts.TotpEnabled)
		}
		if opts.ExcludeTechnicalUsers {
			req = req.ExcludeTechnicalUsers(true)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListUsersResult{
		Users: mapper.UsersFromDTO(resp.Result),
	}

	// Parse total items
	if resp.TotalItems != nil {
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			result.TotalItems = total
		}
	}

	// Set offset from options
	if opts != nil {
		result.Offset = opts.Offset
	}

	return result, nil
}

// GetUsersByEmail retrieves users by their email addresses.
func (s *UserService) GetUsersByEmail(ctx context.Context, emails []string) ([]*model.User, error) {
	if len(emails) == 0 {
		return nil, fmt.Errorf("emails cannot be empty")
	}

	req := s.api.UserServiceGetUsers(ctx).Emails(emails)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.UsersFromDTO(resp.Result), nil
}

// CreateUserAttribute creates an attribute for a user.
func (s *UserService) CreateUserAttribute(ctx context.Context, userID, key, value string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	body := openapi.UserServiceCreateAttributeBody{}
	body.SetKey(key)
	body.SetValue(value)

	_, httpResp, err := s.api.UserServiceCreateAttribute(ctx, userID).Body(body).Execute()
	if err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}
