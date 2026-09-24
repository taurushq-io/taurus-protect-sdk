package service

import (
	"context"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// ExchangeService provides exchange account management operations.
type ExchangeService struct {
	api       *openapi.ExchangeAPIService
	errMapper *ErrorMapper
}

// NewExchangeService creates a new ExchangeService.
func NewExchangeService(client *openapi.APIClient) *ExchangeService {
	return &ExchangeService{
		api:       client.ExchangeAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetExchange retrieves a single exchange account by ID.
func (s *ExchangeService) GetExchange(ctx context.Context, id string) (*model.Exchange, error) {
	req := s.api.ExchangeServiceGetExchange(ctx, id)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.ExchangeFromDTO(resp.Result), nil
}

// ListExchanges retrieves a list of exchange accounts with optional filtering and pagination.
func (s *ExchangeService) ListExchanges(ctx context.Context, opts *model.ListExchangesOptions) (*model.ListExchangesResult, error) {
	if opts == nil {
		opts = &model.ListExchangesOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.ExchangeServiceGetExchanges(ctx), window)
	if opts.CurrencyID != "" {
		req = req.CurrencyID(opts.CurrencyID)
	}
	if opts.IncludeBaseCurrencyValuation {
		req = req.IncludeBaseCurrencyValuation(opts.IncludeBaseCurrencyValuation)
	}
	if opts.ExchangeLabel != "" {
		req = req.ExchangeLabel(opts.ExchangeLabel)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}
	if opts.Status != "" {
		req = req.Status(opts.Status)
	}
	if opts.OnlyPositiveBalance {
		req = req.OnlyPositiveBalance(opts.OnlyPositiveBalance)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListExchangesResult{
		Exchanges: mapper.ExchangesFromDTO(resp.Result),
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	result.Page = page

	return result, nil
}

// ListExchangeCounterparties retrieves a list of exchange counterparties with their exposure and limits.
func (s *ExchangeService) ListExchangeCounterparties(ctx context.Context) (*model.ListExchangeCounterpartiesResult, error) {
	req := s.api.ExchangeServiceGetExchangeCounterparties(ctx)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListExchangeCounterpartiesResult{
		Exchanges: mapper.ExchangeCounterpartiesFromDTO(resp.Exchanges),
	}

	if resp.ExchangesTotalValuation != nil {
		result.ExchangesTotalValuation = *resp.ExchangesTotalValuation
	}

	return result, nil
}

// ExportExchanges exports exchange accounts in the specified format (CSV or JSON).
func (s *ExchangeService) ExportExchanges(ctx context.Context, opts *model.ExportExchangesOptions) (*model.ExportExchangesResult, error) {
	req := s.api.ExchangeServiceExportExchanges(ctx)

	if opts != nil {
		if opts.Format != "" {
			req = req.Format(opts.Format)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ExportExchangesResult{}

	if resp.Result != nil {
		result.Result = *resp.Result
	}

	// Parse total items
	if resp.TotalItems != nil {
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			result.TotalItems = total
		}
	}

	return result, nil
}
