package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/cache"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// PriceService provides price and conversion operations.
type PriceService struct {
	api        *openapi.PricesAPIService
	errMapper  *ErrorMapper
	rulesCache *cache.RulesContainerCache
}

// NewPriceService creates a new PriceService with price signature verification.
//
// Rate and decimals feed amount conversion, so an unverified price is a wrong number a
// caller acts on. Whether prices must be signed is decided by the SuperAdmin-verified
// rules container: see helper.VerifyPrice.
//
// Panics if rulesCache is nil, as NewAddressService does.
func NewPriceService(client *openapi.APIClient, rulesCache *cache.RulesContainerCache) *PriceService {
	if rulesCache == nil {
		panic("rulesCache cannot be nil - price signature verification is mandatory")
	}
	return &PriceService{
		api:        client.PricesAPI,
		errMapper:  NewErrorMapper(),
		rulesCache: rulesCache,
	}
}

// verifyPrices checks every price against the container's PRICEUPDATER keys.
func (s *PriceService) verifyPrices(ctx context.Context, prices []*model.Price) error {
	if len(prices) == 0 {
		return nil
	}
	rulesContainer, err := s.rulesCache.Get(ctx)
	if err != nil {
		return err
	}
	if rulesContainer == nil {
		return &model.IntegrityError{
			Message: "rules container required for price signature verification",
		}
	}
	return helper.VerifyPrices(prices, rulesContainer)
}

// GetPrices retrieves all available currency prices.
func (s *PriceService) GetPrices(ctx context.Context) (*model.GetPricesResult, error) {
	resp, httpResp, err := s.api.PriceServiceGetPrices(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	prices := mapper.PricesFromDTO(resp.Result)
	if err := s.verifyPrices(ctx, prices); err != nil {
		return nil, err
	}

	result := &model.GetPricesResult{
		BaseCurrency: safeString(resp.BaseCurrency),
		Prices:       prices,
	}

	return result, nil
}

// Convert converts an amount from one currency to other currencies.
func (s *PriceService) Convert(ctx context.Context, opts *model.ConvertOptions) (*model.ConversionResult, error) {
	if opts == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}
	if opts.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if opts.Amount == "" {
		return nil, fmt.Errorf("amount is required")
	}

	req := s.api.PriceServiceConvert(ctx, opts.Currency).
		Amount(opts.Amount)

	if len(opts.Symbols) > 0 {
		req = req.Symbols(opts.Symbols)
	}
	if len(opts.TargetCurrencyIds) > 0 {
		req = req.TargetCurrencyIds(opts.TargetCurrencyIds)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.ConversionResultFromDTO(resp), nil
}

// GetPriceHistory retrieves the price history for a currency pair.
func (s *PriceService) GetPriceHistory(ctx context.Context, opts *model.GetPriceHistoryOptions) (*model.GetPriceHistoryResult, error) {
	if opts == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}
	if opts.Base == "" {
		return nil, fmt.Errorf("base currency is required")
	}
	if opts.Quote == "" {
		return nil, fmt.Errorf("quote currency is required")
	}

	req := s.api.PriceServiceGetPricesHistory(ctx, opts.Base, opts.Quote)

	if opts.Limit > 0 {
		req = req.Limit(fmt.Sprintf("%d", opts.Limit))
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.GetPriceHistoryResult{
		History: mapper.PriceHistoryPointsFromDTO(resp.Result),
		Period:  safeString(resp.Period),
	}

	return result, nil
}

// ExportPriceHistory exports the price history in a specified format.
func (s *PriceService) ExportPriceHistory(ctx context.Context, opts *model.ExportPriceHistoryOptions) (*model.ExportPriceHistoryResult, error) {
	if opts == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	req := s.api.PriceServiceExportPricesHistory(ctx)

	if len(opts.CurrencyPairs) > 0 {
		req = req.CurrencyPairs(opts.CurrencyPairs)
	}
	if opts.Limit > 0 {
		req = req.Limit(fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Format != "" {
		req = req.Format(opts.Format)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ExportPriceHistoryResult{
		Data:   safeString(resp.Result),
		Period: safeString(resp.Period),
	}

	return result, nil
}

// safeString helper for dereferencing string pointers safely
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
