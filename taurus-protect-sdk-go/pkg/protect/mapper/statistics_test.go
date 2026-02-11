package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestTagStatisticsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordTagStatistics
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns tag statistics with zero values",
			dto:  &openapi.TgvalidatordTagStatistics{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordTagStatistics {
				tagID := "tag-123"
				totalValuation := "1000000.50"
				tagValue := "Production"
				tagColor := "#FF0000"
				return &openapi.TgvalidatordTagStatistics{
					TagID:          &tagID,
					TotalValuation: &totalValuation,
					Tag: &openapi.TgvalidatordTag{
						Id:    &tagID,
						Value: &tagValue,
						Color: &tagColor,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TagStatisticsFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("TagStatisticsFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("TagStatisticsFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.TagID != nil && got.TagID != *tt.dto.TagID {
				t.Errorf("TagID = %v, want %v", got.TagID, *tt.dto.TagID)
			}
			if tt.dto.TotalValuation != nil && got.TotalValuation != *tt.dto.TotalValuation {
				t.Errorf("TotalValuation = %v, want %v", got.TotalValuation, *tt.dto.TotalValuation)
			}
			// Verify tag is mapped if present
			if tt.dto.Tag != nil {
				if got.Tag == nil {
					t.Error("Tag should not be nil when DTO has tag")
				}
			}
		})
	}
}

func TestTagStatisticsSliceFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordTagStatistics
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordTagStatistics{},
			want: 0,
		},
		{
			name: "converts multiple tag statistics",
			dtos: func() []openapi.TgvalidatordTagStatistics {
				tagID1 := "tag-1"
				tagID2 := "tag-2"
				return []openapi.TgvalidatordTagStatistics{
					{TagID: &tagID1},
					{TagID: &tagID2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TagStatisticsSliceFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("TagStatisticsSliceFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("TagStatisticsSliceFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestPortfolioStatisticsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAggregatedStatsData
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns portfolio statistics with zero values",
			dto:  &openapi.TgvalidatordAggregatedStatsData{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordAggregatedStatsData {
				avgBalance := "500000000000000000"
				addressesCount := "150"
				walletsCount := "25"
				totalBalance := "75000000000000000000"
				totalBalanceBase := "150000.00"
				return &openapi.TgvalidatordAggregatedStatsData{
					AvgBalancePerAddress:     &avgBalance,
					AddressesCount:           &addressesCount,
					WalletsCount:             &walletsCount,
					TotalBalance:             &totalBalance,
					TotalBalanceBaseCurrency: &totalBalanceBase,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PortfolioStatisticsFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("PortfolioStatisticsFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("PortfolioStatisticsFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.AvgBalancePerAddress != nil && got.AvgBalancePerAddress != *tt.dto.AvgBalancePerAddress {
				t.Errorf("AvgBalancePerAddress = %v, want %v", got.AvgBalancePerAddress, *tt.dto.AvgBalancePerAddress)
			}
			if tt.dto.AddressesCount != nil && got.AddressesCount != *tt.dto.AddressesCount {
				t.Errorf("AddressesCount = %v, want %v", got.AddressesCount, *tt.dto.AddressesCount)
			}
			if tt.dto.WalletsCount != nil && got.WalletsCount != *tt.dto.WalletsCount {
				t.Errorf("WalletsCount = %v, want %v", got.WalletsCount, *tt.dto.WalletsCount)
			}
			if tt.dto.TotalBalance != nil && got.TotalBalance != *tt.dto.TotalBalance {
				t.Errorf("TotalBalance = %v, want %v", got.TotalBalance, *tt.dto.TotalBalance)
			}
			if tt.dto.TotalBalanceBaseCurrency != nil && got.TotalBalanceBaseCurrency != *tt.dto.TotalBalanceBaseCurrency {
				t.Errorf("TotalBalanceBaseCurrency = %v, want %v", got.TotalBalanceBaseCurrency, *tt.dto.TotalBalanceBaseCurrency)
			}
		})
	}
}

func TestPortfolioStatisticsHistoryPointFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordAggregatedStatsHistoryPoint
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns history point with zero values",
			dto:  &openapi.TgvalidatordAggregatedStatsHistoryPoint{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordAggregatedStatsHistoryPoint {
				pointDate := time.Now()
				walletsCount := "25"
				return &openapi.TgvalidatordAggregatedStatsHistoryPoint{
					PointDate: &pointDate,
					StatsData: &openapi.TgvalidatordAggregatedStatsData{
						WalletsCount: &walletsCount,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PortfolioStatisticsHistoryPointFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("PortfolioStatisticsHistoryPointFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("PortfolioStatisticsHistoryPointFromDTO() returned nil for non-nil input")
			}
			// Verify timestamp if set
			if tt.dto.PointDate != nil && !got.Timestamp.Equal(*tt.dto.PointDate) {
				t.Errorf("Timestamp = %v, want %v", got.Timestamp, *tt.dto.PointDate)
			}
			// Verify statistics is mapped if present
			if tt.dto.StatsData != nil {
				if got.Statistics == nil {
					t.Error("Statistics should not be nil when DTO has stats data")
				}
			}
		})
	}
}

func TestPortfolioStatisticsHistoryPointsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordAggregatedStatsHistoryPoint
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordAggregatedStatsHistoryPoint{},
			want: 0,
		},
		{
			name: "converts multiple history points",
			dtos: func() []openapi.TgvalidatordAggregatedStatsHistoryPoint {
				t1 := time.Now()
				t2 := time.Now().Add(-time.Hour)
				return []openapi.TgvalidatordAggregatedStatsHistoryPoint{
					{PointDate: &t1},
					{PointDate: &t2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PortfolioStatisticsHistoryPointsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("PortfolioStatisticsHistoryPointsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("PortfolioStatisticsHistoryPointsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestPortfolioStatisticsHistoryPointFromDTO_NilStatsData(t *testing.T) {
	pointDate := time.Now()
	dto := &openapi.TgvalidatordAggregatedStatsHistoryPoint{
		PointDate: &pointDate,
		StatsData: nil,
	}

	got := PortfolioStatisticsHistoryPointFromDTO(dto)
	if got == nil {
		t.Fatal("PortfolioStatisticsHistoryPointFromDTO() returned nil for non-nil input")
	}
	if got.Statistics != nil {
		t.Errorf("Statistics should be nil when DTO stats data is nil, got %v", got.Statistics)
	}
	if !got.Timestamp.Equal(pointDate) {
		t.Errorf("Timestamp = %v, want %v", got.Timestamp, pointDate)
	}
}

func TestPortfolioStatisticsHistoryPointFromDTO_NilPointDate(t *testing.T) {
	walletsCount := "10"
	dto := &openapi.TgvalidatordAggregatedStatsHistoryPoint{
		PointDate: nil,
		StatsData: &openapi.TgvalidatordAggregatedStatsData{
			WalletsCount: &walletsCount,
		},
	}

	got := PortfolioStatisticsHistoryPointFromDTO(dto)
	if got == nil {
		t.Fatal("PortfolioStatisticsHistoryPointFromDTO() returned nil for non-nil input")
	}
	// When point date is nil, it should be the zero time value
	if !got.Timestamp.IsZero() {
		t.Errorf("Timestamp should be zero time when nil, got %v", got.Timestamp)
	}
	if got.Statistics == nil {
		t.Error("Statistics should not be nil when DTO has stats data")
	}
}

func TestTagStatisticsFromDTO_NilTag(t *testing.T) {
	tagID := "tag-123"
	totalValuation := "1000.00"
	dto := &openapi.TgvalidatordTagStatistics{
		TagID:          &tagID,
		TotalValuation: &totalValuation,
		Tag:            nil,
	}

	got := TagStatisticsFromDTO(dto)
	if got == nil {
		t.Fatal("TagStatisticsFromDTO() returned nil for non-nil input")
	}
	if got.Tag != nil {
		t.Errorf("Tag should be nil when DTO tag is nil, got %v", got.Tag)
	}
	if got.TagID != tagID {
		t.Errorf("TagID = %v, want %v", got.TagID, tagID)
	}
	if got.TotalValuation != totalValuation {
		t.Errorf("TotalValuation = %v, want %v", got.TotalValuation, totalValuation)
	}
}
