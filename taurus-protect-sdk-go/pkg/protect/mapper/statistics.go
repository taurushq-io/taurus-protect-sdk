package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// TagStatisticsFromDTO converts an OpenAPI TagStatistics to a domain TagStatistics.
func TagStatisticsFromDTO(dto *openapi.TgvalidatordTagStatistics) *model.TagStatistics {
	if dto == nil {
		return nil
	}

	stats := &model.TagStatistics{
		TagID:          safeString(dto.TagID),
		TotalValuation: safeString(dto.TotalValuation),
	}

	// Convert tag info if present
	if dto.Tag != nil {
		stats.Tag = TagFromDTO(dto.Tag)
	}

	return stats
}

// TagStatisticsSliceFromDTO converts a slice of OpenAPI TagStatistics to domain TagStatistics.
func TagStatisticsSliceFromDTO(dtos []openapi.TgvalidatordTagStatistics) []*model.TagStatistics {
	if dtos == nil {
		return nil
	}
	stats := make([]*model.TagStatistics, len(dtos))
	for i := range dtos {
		stats[i] = TagStatisticsFromDTO(&dtos[i])
	}
	return stats
}

// PortfolioStatisticsFromDTO converts an OpenAPI AggregatedStatsData to a domain PortfolioStatistics.
func PortfolioStatisticsFromDTO(dto *openapi.TgvalidatordAggregatedStatsData) *model.PortfolioStatistics {
	if dto == nil {
		return nil
	}

	return &model.PortfolioStatistics{
		AvgBalancePerAddress:     safeString(dto.AvgBalancePerAddress),
		AddressesCount:           safeString(dto.AddressesCount),
		WalletsCount:             safeString(dto.WalletsCount),
		TotalBalance:             safeString(dto.TotalBalance),
		TotalBalanceBaseCurrency: safeString(dto.TotalBalanceBaseCurrency),
	}
}

// PortfolioStatisticsHistoryPointFromDTO converts an OpenAPI AggregatedStatsHistoryPoint to a domain model.
func PortfolioStatisticsHistoryPointFromDTO(dto *openapi.TgvalidatordAggregatedStatsHistoryPoint) *model.PortfolioStatisticsHistoryPoint {
	if dto == nil {
		return nil
	}

	point := &model.PortfolioStatisticsHistoryPoint{}

	if dto.PointDate != nil {
		point.Timestamp = *dto.PointDate
	}

	if dto.StatsData != nil {
		point.Statistics = PortfolioStatisticsFromDTO(dto.StatsData)
	}

	return point
}

// PortfolioStatisticsHistoryPointsFromDTO converts a slice of OpenAPI AggregatedStatsHistoryPoints to domain models.
func PortfolioStatisticsHistoryPointsFromDTO(dtos []openapi.TgvalidatordAggregatedStatsHistoryPoint) []*model.PortfolioStatisticsHistoryPoint {
	if dtos == nil {
		return nil
	}
	points := make([]*model.PortfolioStatisticsHistoryPoint, len(dtos))
	for i := range dtos {
		points[i] = PortfolioStatisticsHistoryPointFromDTO(&dtos[i])
	}
	return points
}
