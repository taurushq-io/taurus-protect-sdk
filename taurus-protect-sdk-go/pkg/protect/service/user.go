package service

import (
	"context"
	"fmt"

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

// GetMe retrieves the currently authenticated user, with the enforced-in-rules flags computed.
func (s *UserService) GetMe(ctx context.Context) (*model.User, error) {
	resp, httpResp, err := s.api.UserServiceGetMe(ctx).CheckEnforcedInRules(true).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("user not found")
	}

	user := mapper.UserFromDTO(resp.Result)
	defaultComputedUserFlags(user)
	return user, nil
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

// ListUsers retrieves one page of users. Result.Pagination is never nil; continue with its
// NextOffset until HasMore is false.
func (s *UserService) ListUsers(ctx context.Context, opts *model.ListUsersOptions) (*model.ListUsersResult, error) {
	if opts == nil {
		opts = &model.ListUsersOptions{}
	}
	window, err := resolveOffsetWindow(opts.Limit, opts.Offset)
	if err != nil {
		return nil, err
	}

	req := applyOffsetWindow(s.api.UserServiceGetUsers(ctx), window)
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

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	// A synthetic daemon user can be appended beyond the limit, so the next page starts at
	// offset + min(rows, limit).
	pagination, err := offsetPagination(rulePlusMinRowsLimit, window, len(resp.Result), 0,
		offsetReply{TotalItems: resp.TotalItems})
	if err != nil {
		return nil, err
	}
	users := mapper.UsersFromDTO(resp.Result)
	defaultComputedUserFlags(users...)
	return &model.ListUsersResult{
		Users:      users,
		Pagination: pagination,
	}, nil
}

// defaultComputedUserFlags sets an absent EnforcedInRules to false on users read from an
// endpoint that computes it: validatord leaves a false bool out of its JSON. The public-key flag
// is a BoolValue, sent whenever computed, so it is never defaulted.
func defaultComputedUserFlags(users ...*model.User) {
	for _, user := range users {
		if user == nil {
			continue
		}
		user.EnforcedInRules = falseIfNil(user.EnforcedInRules)
		for i := range user.Groups {
			user.Groups[i].EnforcedInRules = falseIfNil(user.Groups[i].EnforcedInRules)
		}
	}
}

func falseIfNil(b *bool) *bool {
	if b != nil {
		return b
	}
	f := false
	return &f
}

// GetUsersByEmail retrieves the users with these email addresses. The emails are sent in
// batches of at most model.MaxPageSize, and each batch is walked to its last page.
func (s *UserService) GetUsersByEmail(ctx context.Context, emails []string) ([]*model.User, error) {
	if len(emails) == 0 {
		return nil, fmt.Errorf("emails cannot be empty")
	}

	var users []*model.User
	for _, batch := range chunkIDs(emails, model.MaxPageSize) {
		opts := &model.ListUsersOptions{Emails: batch, Limit: model.MaxPageSize}
		for {
			result, err := s.ListUsers(ctx, opts)
			if err != nil {
				return nil, err
			}
			users = append(users, result.Users...)
			if !result.Pagination.HasMore {
				break
			}
			opts.Offset = result.Pagination.NextOffset
		}
	}
	return users, nil
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
