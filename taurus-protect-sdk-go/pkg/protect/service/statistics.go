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
	req := s.api.StatisticsServiceGetAggregatedTagStats(ctx)

	if opts != nil {
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
		if opts.SortBy != "" {
			req = req.SortBy(opts.SortBy)
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
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.ListTagStatisticsResult{
		TagStatistics: mapper.TagStatisticsSliceFromDTO(resp.Result),
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
	req := s.api.StatisticsServiceGetPortfolioStatisticsHistory(ctx)

	if opts != nil {
		if opts.IntervalHours > 0 {
			req = req.IntervalHours(fmt.Sprintf("%d", opts.IntervalHours))
		}
		if opts.From != nil {
			req = req.From(*opts.From)
		}
		if opts.To != nil {
			req = req.To(*opts.To)
		}
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
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
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	result := &model.GetPortfolioStatisticsHistoryResult{
		HistoryPoints: mapper.PortfolioStatisticsHistoryPointsFromDTO(resp.Result),
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
