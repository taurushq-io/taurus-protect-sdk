package service

import (
	"context"

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

// GetBalances retrieves one page of the tenant's total balances, one per asset. An asset is
// identified by a full triplet of attributes: blockchain, contract number, and token ID.
// Continue with Page.NextCursor until Page.HasMore is false.
func (s *BalanceService) GetBalances(ctx context.Context, opts *model.GetBalancesOptions) (*model.GetBalancesResult, error) {
	if opts == nil {
		opts = &model.GetBalancesOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, "", "")
	if err != nil {
		return nil, err
	}

	// requestCursor only: the legacy limit + bytes cursor pair is not sent.
	req := applyRequestCursorQuery(s.api.WalletServiceGetBalances(ctx), window)
	if opts.Currency != "" {
		req = req.Currency(opts.Currency)
	}
	if opts.TokenID != "" {
		req = req.TokenId(opts.TokenID)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor, HasTotal: true, Total: resp.Total})
	if err != nil {
		return nil, err
	}
	return &model.GetBalancesResult{
		Balances: mapper.AssetBalancesFromDTO(resp.Balances),
		Page:     page,
	}, nil
}
