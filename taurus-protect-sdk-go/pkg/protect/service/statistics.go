package service

import (
	"context"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// StatisticsService provides statistics retrieval operations.
type StatisticsService struct {
	api       *openapi.StatisticsAPIService
	errMapper *ErrorMapper
}

// NewStatisticsService creates a new StatisticsService.
func NewStatisticsService(client *openapi.APIClient) *StatisticsService {
	return &StatisticsService{
		api:       client.StatisticsAPI,
		errMapper: NewErrorMapper(),
	}
}

// ListTagStatistics retrieves a list of tag statistics with optional filtering and pagination.
func (s *StatisticsService) ListTagStatistics(ctx context.Context, opts *model.ListTagStatisticsOptions) (*model.ListTagStatisticsResult, error) {
	if opts == nil {
		opts = &model.ListTagStatisticsOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.StatisticsServiceGetAggregatedTagStats(ctx), window)
	if opts.Query != "" {
		req = req.Query(opts.Query)
	}
	if opts.SortBy != "" {
		req = req.SortBy(opts.SortBy)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListTagStatisticsResult{
		TagStatistics: mapper.TagStatisticsSliceFromDTO(resp.Result),
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	result.Page = page

	return result, nil
}

// GetPortfolioStatistics retrieves the global portfolio statistics.
func (s *StatisticsService) GetPortfolioStatistics(ctx context.Context) (*model.PortfolioStatistics, error) {
	resp, httpResp, err := s.api.StatisticsServiceGetPortfolioStatistics(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.PortfolioStatisticsFromDTO(resp.Result), nil
}

// GetPortfolioStatisticsHistory retrieves the portfolio statistics history with optional filtering and pagination.
func (s *StatisticsService) GetPortfolioStatisticsHistory(ctx context.Context, opts *model.GetPortfolioStatisticsHistoryOptions) (*model.GetPortfolioStatisticsHistoryResult, error) {
	if opts == nil {
		opts = &model.GetPortfolioStatisticsHistoryOptions{}
	}
	window, err := resolveCursorWindow(opts.PageSize, opts.Cursor, opts.CurrentPage, opts.PageRequest)
	if err != nil {
		return nil, err
	}

	req := applyCursorQuery(s.api.StatisticsServiceGetPortfolioStatisticsHistory(ctx), window)
	if opts.IntervalHours > 0 {
		req = req.IntervalHours(fmt.Sprintf("%d", opts.IntervalHours))
	}
	if opts.From != nil {
		req = req.From(*opts.From)
	}
	if opts.To != nil {
		req = req.To(*opts.To)
	}
	if opts.SortOrder != "" {
		req = req.SortOrder(opts.SortOrder)
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.GetPortfolioStatisticsHistoryResult{
		HistoryPoints: mapper.PortfolioStatisticsHistoryPointsFromDTO(resp.Result),
	}

	page, err := cursorPage(window.pageSize, cursorReply{Cursor: resp.Cursor})
	if err != nil {
		return nil, err
	}
	result.Page = page

	return result, nil
}
