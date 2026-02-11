package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.MultiFactorSignatureApprovalResult;
import com.taurushq.sdk.protect.client.model.MultiFactorSignatureEntityType;
import com.taurushq.sdk.protect.client.model.MultiFactorSignatureInfo;
import com.taurushq.sdk.protect.client.model.MultiFactorSignatureResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproveMultiFactorSignatureReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateMultiFactorSignaturesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetMultiFactorSignatureEntitiesInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMultiFactorSignaturesEntityType;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for converting multi-factor signature DTOs to domain models.
 *
 * @see MultiFactorSignatureInfo
 * @see MultiFactorSignatureResult
 */
@Mapper
public interface MultiFactorSignatureMapper {

    /**
     * Singleton instance of the mapper.
     */
    MultiFactorSignatureMapper INSTANCE = Mappers.getMapper(MultiFactorSignatureMapper.class);

    /**
     * Converts a multi-factor signature info reply to a domain model.
     *
     * @param dto the OpenAPI reply DTO
     * @return the domain model
     */
    MultiFactorSignatureInfo fromDTO(TgvalidatordGetMultiFactorSignatureEntitiesInfoReply dto);

    /**
     * Converts an entity type DTO to a domain model.
     *
     * @param dto the OpenAPI entity type DTO
     * @return the domain model
     */
    MultiFactorSignatureEntityType fromEntityTypeDTO(TgvalidatordMultiFactorSignaturesEntityType dto);

    /**
     * Converts a create reply to a result model.
     *
     * @param dto the OpenAPI create reply DTO
     * @return the result model
     */
    MultiFactorSignatureResult fromCreateDTO(TgvalidatordCreateMultiFactorSignaturesReply dto);

    /**
     * Converts an approval reply to a result model.
     *
     * @param dto the OpenAPI approval reply DTO
     * @return the approval result model
     */
    MultiFactorSignatureApprovalResult fromApprovalDTO(TgvalidatordApproveMultiFactorSignatureReply dto);
}
