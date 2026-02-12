package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.PortfolioStatistics;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAggregatedStatsData;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for converting statistics OpenAPI DTOs to client model objects.
 */
@Mapper
public interface StatisticsMapper {

    /**
     * Singleton instance of the mapper.
     */
    StatisticsMapper INSTANCE = Mappers.getMapper(StatisticsMapper.class);

    /**
     * Maps aggregated stats data from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    PortfolioStatistics fromDTO(TgvalidatordAggregatedStatsData dto);
}
