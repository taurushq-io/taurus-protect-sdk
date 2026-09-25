package service

import (
	"context"
	"fmt"
	"strconv"

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

// ListPrices retrieves one page of currency prices (PriceService_QueryPricesV2), each verified
// against the PRICEUPDATER keys of the rules container. Continue with Page.NextCursor until
// Page.HasMore is false.
func (s *PriceService) ListPrices(ctx context.Context, opts *model.ListPricesOptions) (*model.ListPricesResult, error) {
	if opts == nil {
		opts = &model.ListPricesOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, "", "")
	if err != nil {
		return nil, err
	}

	body := openapi.TgvalidatordQueryPricesV2Request{Cursor: window.body()}
	if opts.OnlyPrimary {
		onlyPrimary := true
		body.OnlyPrimary = &onlyPrimary
	}
	if opts.SortOrder != "" {
		body.SortOrder = &opts.SortOrder
	}
	// The currency filter is a oneof: from, fromTo or to.
	switch {
	case opts.FromCurrencyID != "" && len(opts.ToCurrencyIDs) > 0:
		body.FromTo = &openapi.TgvalidatordCurrencyFromToFilter{
			CurrencyFromId: &opts.FromCurrencyID,
			CurrencyToIds:  opts.ToCurrencyIDs,
		}
	case opts.FromCurrencyID != "":
		body.From = &openapi.TgvalidatordCurrencyFromFilter{CurrencyFromId: &opts.FromCurrencyID}
	case len(opts.ToCurrencyIDs) > 0:
		body.To = &openapi.TgvalidatordCurrencyToFilter{CurrencyToIds: opts.ToCurrencyIDs}
	}

	resp, httpResp, err := s.api.PriceServiceQueryPricesV2(ctx).Body(body).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	prices := mapper.PricesFromDTO(resp.Result)
	if err := s.verifyPrices(ctx, prices); err != nil {
		return nil, err
	}

	return &model.ListPricesResult{
		BaseCurrency: safeString(resp.BaseCurrency),
		Prices:       prices,
		Page:         page,
	}, nil
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

	// Price history cannot page: Limit is the number of newest daily points, up to a year.
	limit, err := resolveSize("Limit", opts.Limit, model.MaxPriceHistoryLimit)
	if err != nil {
		return nil, err
	}
	req := s.api.PriceServiceGetPricesHistory(ctx, opts.Base, opts.Quote).Limit(strconv.FormatInt(limit, 10))

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

	// The export cannot page and has no SDK maximum; the server bounds it.
	limit, err := resolveSize("Limit", opts.Limit, 0)
	if err != nil {
		return nil, err
	}
	req := s.api.PriceServiceExportPricesHistory(ctx).Limit(strconv.FormatInt(limit, 10))
	if len(opts.CurrencyPairs) > 0 {
		req = req.CurrencyPairs(opts.CurrencyPairs)
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
