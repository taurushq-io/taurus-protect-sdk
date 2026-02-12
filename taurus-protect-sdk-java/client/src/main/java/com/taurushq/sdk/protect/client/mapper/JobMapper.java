package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Job;
import com.taurushq.sdk.protect.client.model.JobStatistics;
import com.taurushq.sdk.protect.client.model.JobStatus;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordJob;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordJobStatistics;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordJobStatus;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting job-related DTOs to domain models.
 */
@Mapper
public interface JobMapper {

    JobMapper INSTANCE = Mappers.getMapper(JobMapper.class);

    /**
     * Converts a Job DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    Job fromDTO(TgvalidatordJob dto);

    /**
     * Converts a list of Job DTOs to domain models.
     *
     * @param dtos the DTOs to convert
     * @return the domain models
     */
    List<Job> fromDTOList(List<TgvalidatordJob> dtos);

    /**
     * Converts a JobStatistics DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    JobStatistics fromStatisticsDTO(TgvalidatordJobStatistics dto);

    /**
     * Converts a JobStatus DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    JobStatus fromStatusDTO(TgvalidatordJobStatus dto);
}
