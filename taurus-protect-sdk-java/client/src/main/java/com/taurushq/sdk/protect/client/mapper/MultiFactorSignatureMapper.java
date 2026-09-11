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
     * <p><b>Hand-written, because MapStruct silently produced an all-null bean here.</b> The
     * source is a bare enum and the target has {@code id} and {@code kind}; with no explicit
     * mapping, MapStruct's default bean mapping from an enum source sets no target properties
     * at all, so {@code getKind()} returned null and a caller could not even tell which KIND
     * of entity the {@code payloadToSign} bytes were supposed to cover. The mapper test only
     * asserted non-null, so nothing caught it.
     *
     * <p>{@code id} stays null on purpose: <b>the reply carries no entity id</b>, singular or
     * plural. That absence is the reason the SDK cannot bind {@code payloadToSign} to a
     * verified entity — see {@code MultiFactorSignatureService.getMultiFactorSignatureInfo}
     * and {@code TODOS.md}. Do not invent one.
     *
     * @param dto the OpenAPI entity type DTO
     * @return the domain model carrying the kind, or {@code null} if the DTO was null
     */
    default MultiFactorSignatureEntityType fromEntityTypeDTO(
            final TgvalidatordMultiFactorSignaturesEntityType dto) {
        if (dto == null) {
            return null;
        }
        MultiFactorSignatureEntityType kind = new MultiFactorSignatureEntityType();
        kind.setKind(dto.getValue());
        return kind;
    }

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
