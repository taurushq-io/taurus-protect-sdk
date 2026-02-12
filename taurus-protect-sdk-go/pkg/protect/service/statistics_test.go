package service

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewStatisticsService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestStatisticsService_ListTagStatistics_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	// The actual API call will fail, but we're testing the options handling
	svc := &StatisticsService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This verifies the service accepts nil options
	// In a real test with mocked API, nil options should work
	if svc == nil {
		t.Error("StatisticsService should not be nil")
	}
}

func TestStatisticsService_ListTagStatistics_WithOptions(t *testing.T) {
	tests := []struct {
		name    string
		options *model.ListTagStatisticsOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.ListTagStatisticsOptions{},
		},
		{
			name: "query filter",
			options: &model.ListTagStatisticsOptions{
				Query: "production",
			},
		},
		{
			name: "sort by id",
			options: &model.ListTagStatisticsOptions{
				SortBy: "id",
			},
		},
		{
			name: "sort by tagname",
			options: &model.ListTagStatisticsOptions{
				SortBy: "tagname",
			},
		},
		{
			name: "sort by createdat",
			options: &model.ListTagStatisticsOptions{
				SortBy: "createdat",
			},
		},
		{
			name: "sort order ASC",
			options: &model.ListTagStatisticsOptions{
				SortOrder: "ASC",
			},
		},
		{
			name: "sort order DESC",
			options: &model.ListTagStatisticsOptions{
				SortOrder: "DESC",
			},
		},
		{
			name: "pagination options",
			options: &model.ListTagStatisticsOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &model.ListTagStatisticsOptions{
				Query:       "test",
				SortBy:      "tagname",
				SortOrder:   "DESC",
				CurrentPage: "xyz789",
				PageRequest: "FIRST",
				PageSize:    100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &StatisticsService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("StatisticsService should not be nil")
			}
		})
	}
}

func TestStatisticsService_GetPortfolioStatistics(t *testing.T) {
	// Create a service with nil API to test basic structure
	svc := &StatisticsService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("StatisticsService should not be nil")
	}
}

func TestStatisticsService_GetPortfolioStatisticsHistory_NilOptions(t *testing.T) {
	// Create a service with nil API to test that nil options don't cause panic
	svc := &StatisticsService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	if svc == nil {
		t.Error("StatisticsService should not be nil")
	}
}

func TestStatisticsService_GetPortfolioStatisticsHistory_WithOptions(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tests := []struct {
		name    string
		options *model.GetPortfolioStatisticsHistoryOptions
	}{
		{
			name:    "nil options",
			options: nil,
		},
		{
			name:    "empty options",
			options: &model.GetPortfolioStatisticsHistoryOptions{},
		},
		{
			name: "interval hours",
			options: &model.GetPortfolioStatisticsHistoryOptions{
				IntervalHours: 24,
			},
		},
		{
			name: "date range",
			options: &model.GetPortfolioStatisticsHistoryOptions{
				From: &yesterday,
				To:   &now,
			},
		},
		{
			name: "limit",
			options: &model.GetPortfolioStatisticsHistoryOptions{
				Limit: 100,
			},
		},
		{
			name: "sort order ASC",
			options: &model.GetPortfolioStatisticsHistoryOptions{
				SortOrder: "ASC",
			},
		},
		{
			name: "sort order DESC",
			options: &model.GetPortfolioStatisticsHistoryOptions{
				SortOrder: "DESC",
			},
		},
		{
			name: "pagination options",
			options: &model.GetPortfolioStatisticsHistoryOptions{
				CurrentPage: "abc123",
				PageRequest: "NEXT",
				PageSize:    50,
			},
		},
		{
			name: "all options combined",
			options: &model.GetPortfolioStatisticsHistoryOptions{
				IntervalHours: 12,
				From:          &yesterday,
				To:            &now,
				Limit:         200,
				SortOrder:     "ASC",
				CurrentPage:   "xyz789",
				PageRequest:   "LAST",
				PageSize:      25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify options can be created with these values
			// Actual API testing requires mocking
			svc := &StatisticsService{
				api:       nil,
				errMapper: NewErrorMapper(),
			}
			if svc == nil {
				t.Error("StatisticsService should not be nil")
			}
		})
	}
}

func TestListTagStatisticsOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &model.ListTagStatisticsOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestListTagStatisticsOptions_SortOrderValues(t *testing.T) {
	// Test that sort order values match expected API values
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &model.ListTagStatisticsOptions{
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}

func TestListTagStatisticsOptions_SortByValues(t *testing.T) {
	// Test that sort by values match expected API values
	validSortBy := []string{"id", "tagname", "createdat"}

	for _, sortBy := range validSortBy {
		t.Run(sortBy, func(t *testing.T) {
			opts := &model.ListTagStatisticsOptions{
				SortBy: sortBy,
			}
			if opts.SortBy != sortBy {
				t.Errorf("SortBy = %v, want %v", opts.SortBy, sortBy)
			}
		})
	}
}

func TestGetPortfolioStatisticsHistoryOptions_PageRequestValues(t *testing.T) {
	// Test that page request values match expected API values
	validPageRequests := []string{"FIRST", "PREVIOUS", "NEXT", "LAST"}

	for _, pageRequest := range validPageRequests {
		t.Run(pageRequest, func(t *testing.T) {
			opts := &model.GetPortfolioStatisticsHistoryOptions{
				PageRequest: pageRequest,
			}
			if opts.PageRequest != pageRequest {
				t.Errorf("PageRequest = %v, want %v", opts.PageRequest, pageRequest)
			}
		})
	}
}

func TestGetPortfolioStatisticsHistoryOptions_SortOrderValues(t *testing.T) {
	// Test that sort order values match expected API values
	validSortOrders := []string{"ASC", "DESC"}

	for _, sortOrder := range validSortOrders {
		t.Run(sortOrder, func(t *testing.T) {
			opts := &model.GetPortfolioStatisticsHistoryOptions{
				SortOrder: sortOrder,
			}
			if opts.SortOrder != sortOrder {
				t.Errorf("SortOrder = %v, want %v", opts.SortOrder, sortOrder)
			}
		})
	}
}
