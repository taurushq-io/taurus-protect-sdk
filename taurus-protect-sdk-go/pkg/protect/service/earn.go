package service

import (
	"context"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// EarnService provides the Earn reward reads.
type EarnService struct {
	api       *openapi.EarnAPIService
	errMapper *ErrorMapper
}

// NewEarnService creates a new EarnService.
func NewEarnService(client *openapi.APIClient) *EarnService {
	return &EarnService{
		api:       client.EarnAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListRewards retrieves one page of rewards (EarnService_GetRewards). Continue with
// Page.NextCursor until Page.HasMore is false.
func (s *EarnService) ListRewards(ctx context.Context, opts *model.ListEarnRewardsOptions) (*model.ListEarnRewardsResult, error) {
	if opts == nil {
		opts = &model.ListEarnRewardsOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, "", "")
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.EarnServiceGetRewards(ctx), window)
	if opts.RecipientAddressID != "" {
		req = req.RecipientAddressId(opts.RecipientAddressID)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}
	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	return &model.ListEarnRewardsResult{
		Rewards: mapper.EarnRewardsFromDTO(resp.Rewards),
		Page:    page,
	}, nil
}
