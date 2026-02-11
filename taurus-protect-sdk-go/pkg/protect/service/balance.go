package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// BalanceService provides balance management operations.
type BalanceService struct {
	api       *openapi.BalancesAPIService
	errMapper *ErrorMapper
}

// NewBalanceService creates a new BalanceService.
func NewBalanceService(client *openapi.APIClient) *BalanceService {
	return &BalanceService{
		api:       client.BalancesAPI,
		errMapper: NewErrorMapper(),
	}
}

// GetBalances retrieves the total balances for the tenant, for each asset.
// An asset is identified by a full triplet of attributes: blockchain, contract number, and token ID.
func (s *BalanceService) GetBalances(ctx context.Context, opts *model.GetBalancesOptions) (*model.GetBalancesResult, error) {
	req := s.api.WalletServiceGetBalances(ctx)

	if opts != nil {
		if opts.Currency != "" {
			req = req.Currency(opts.Currency)
		}
		if opts.TokenID != "" {
			req = req.TokenId(opts.TokenID)
		}
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Cursor != "" {
			req = req.Cursor(opts.Cursor)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.GetBalancesResult{
		Balances: mapper.AssetBalancesFromDTO(resp.Balances),
	}

	// Parse total count
	if resp.Total != nil {
		if total, parseErr := strconv.ParseInt(*resp.Total, 10, 64); parseErr == nil {
			result.Total = total
		}
	}

	// Set next cursor for pagination
	if resp.Next != nil {
		result.NextCursor = *resp.Next
	}

	return result, nil
}
