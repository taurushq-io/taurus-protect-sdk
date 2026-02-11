package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.AuditTrail;
import com.taurushq.sdk.protect.client.model.AuditTrailResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAuditTrail;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAuditTrailsReply;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting audit OpenAPI DTOs to client model objects.
 */
@Mapper(uses = ApiResponseCursorMapper.class)
public interface AuditMapper {

    /**
     * Singleton instance of the mapper.
     */
    AuditMapper INSTANCE = Mappers.getMapper(AuditMapper.class);

    /**
     * Maps an audit trail from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    AuditTrail fromDTO(TgvalidatordAuditTrail dto);

    /**
     * Maps a list of audit trails from DTOs.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<AuditTrail> fromDTOList(List<TgvalidatordAuditTrail> dtos);

    /**
     * Maps an audit trails reply to a result with pagination.
     *
     * @param reply the OpenAPI reply
     * @return the domain model result
     */
    @Mapping(target = "auditTrails", source = "result")
    AuditTrailResult fromReply(TgvalidatordGetAuditTrailsReply reply);
}
