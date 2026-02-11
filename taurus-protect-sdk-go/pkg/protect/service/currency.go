package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// CurrencyService provides currency management operations.
type CurrencyService struct {
	api       *openapi.CurrenciesAPIService
	errMapper *ErrorMapper
}

// NewCurrencyService creates a new CurrencyService.
func NewCurrencyService(client *openapi.APIClient) *CurrencyService {
	return &CurrencyService{
		api:       client.CurrenciesAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetCurrencies retrieves a list of all currencies.
func (s *CurrencyService) GetCurrencies(ctx context.Context, opts *model.ListCurrenciesOptions) ([]*model.Currency, error) {
	req := s.api.WalletServiceGetCurrencies(ctx)

	if opts != nil {
		if opts.ShowDisabled {
			req = req.ShowDisabled(true)
		}
		if opts.IncludeLogo {
			req = req.IncludeLogo(true)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.CurrenciesFromDTO(resp.Result), nil
}

// GetCurrency retrieves a single currency by its filter criteria.
// The Blockchain and Network fields in opts are required.
func (s *CurrencyService) GetCurrency(ctx context.Context, opts *model.GetCurrencyOptions) (*model.Currency, error) {
	if opts == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}
	if opts.Blockchain == "" {
		return nil, fmt.Errorf("blockchain is required")
	}
	if opts.Network == "" {
		return nil, fmt.Errorf("network is required")
	}

	req := s.api.WalletServiceGetCurrency(ctx).
		UniqueCurrencyFilterBlockchain(opts.Blockchain).
		UniqueCurrencyFilterNetwork(opts.Network)

	if opts.CurrencyID != "" {
		req = req.CurrencyID(opts.CurrencyID)
	}
	if opts.TokenContractAddress != "" {
		req = req.UniqueCurrencyFilterTokenContractAddress(opts.TokenContractAddress)
	}
	if opts.TokenID != "" {
		req = req.UniqueCurrencyFilterTokenID(opts.TokenID)
	}
	if opts.ShowDisabled {
		req = req.ShowDisabled(true)
	}
	if opts.IncludeLogo {
		req = req.IncludeLogo(true)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return mapper.CurrencyFromDTO(resp.Result), nil
}
