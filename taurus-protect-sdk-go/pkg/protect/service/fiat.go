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

// ListFiatProviderAccounts retrieves one page of fiat provider accounts. Provider and Label are
// required. Continue with Page.NextCursor until Page.HasMore is false.
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
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.FiatProviderServiceGetFiatProviderAccounts(ctx), window).
		Provider(opts.Provider).
		Label(opts.Label)
	if opts.AccountType != "" {
		req = req.AccountType(opts.AccountType)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}
	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	return &model.ListFiatProviderAccountsResult{
		Accounts: mapper.FiatProviderAccountsFromDTO(resp.Result),
		Page:     page,
	}, nil
}

// ListFiatProviderEntities retrieves one page of the entities registered with fiat providers.
// Continue with Page.NextCursor until Page.HasMore is false.
func (s *FiatService) ListFiatProviderEntities(ctx context.Context, opts *model.ListFiatProviderEntitiesOptions) (*model.ListFiatProviderEntitiesResult, error) {
	if opts == nil {
		opts = &model.ListFiatProviderEntitiesOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, "", "")
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.FiatProviderServiceGetFiatProviderEntities(ctx), window)
	if opts.Provider != "" {
		req = req.Provider(opts.Provider)
	}
	if opts.Label != "" {
		req = req.Label(opts.Label)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}
	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	return &model.ListFiatProviderEntitiesResult{
		Entities: mapper.FiatProviderEntitiesFromDTO(resp.Result),
		Page:     page,
	}, nil
}
