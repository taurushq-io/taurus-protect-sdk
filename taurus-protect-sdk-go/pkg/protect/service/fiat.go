package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// FiatService provides fiat provider management operations.
type FiatService struct {
	api       *openapi.FiatAPIService
	errMapper *ErrorMapper
}

// NewFiatService creates a new FiatService.
func NewFiatService(client *openapi.APIClient) *FiatService {
	return &FiatService{
		api:       client.FiatAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListFiatProviders retrieves a list of all enabled fiat providers and their valuations.
func (s *FiatService) ListFiatProviders(ctx context.Context) (*model.ListFiatProvidersResult, error) {
	req := s.api.FiatProviderServiceGetFiatProviders(ctx)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListFiatProvidersResult{
		FiatProviders: mapper.FiatProvidersFromDTO(resp.FiatProviders),
	}

	if resp.FiatProvidersTotalValuation != nil {
		result.TotalValuation = *resp.FiatProvidersTotalValuation
	}

	return result, nil
}

// GetFiatProviderAccount retrieves a fiat provider account by ID.
func (s *FiatService) GetFiatProviderAccount(ctx context.Context, id string) (*model.FiatProviderAccount, error) {
	req := s.api.FiatProviderServiceGetFiatProviderAccount(ctx, id)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.FiatProviderAccountFromDTO(resp.Result), nil
}

// ListFiatProviderAccounts retrieves a list of fiat provider accounts with optional filtering and pagination.
func (s *FiatService) ListFiatProviderAccounts(ctx context.Context, opts *model.ListFiatProviderAccountsOptions) (*model.ListFiatProviderAccountsResult, error) {
	if opts == nil {
		return nil, fmt.Errorf("options are required: provider and label must be specified")
	}
	if opts.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}
	if opts.Label == "" {
		return nil, fmt.Errorf("label is required")
	}

	req := s.api.FiatProviderServiceGetFiatProviderAccounts(ctx)
	req = req.Provider(opts.Provider)
	req = req.Label(opts.Label)

	if opts.AccountType != "" {
		req = req.AccountType(opts.AccountType)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
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

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListFiatProviderAccountsResult{
		Accounts: mapper.FiatProviderAccountsFromDTO(resp.Result),
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
