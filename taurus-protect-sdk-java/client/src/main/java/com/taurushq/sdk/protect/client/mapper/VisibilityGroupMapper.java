package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.VisibilityGroup;
import com.taurushq.sdk.protect.client.model.VisibilityGroupUser;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalVisibilityGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordInternalVisibilityGroupUser;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting visibility group DTOs to domain models.
 *
 * @see VisibilityGroup
 * @see VisibilityGroupUser
 */
@Mapper
public interface VisibilityGroupMapper {

    /**
     * Singleton instance of the mapper.
     */
    VisibilityGroupMapper INSTANCE = Mappers.getMapper(VisibilityGroupMapper.class);

    /**
     * Converts a visibility group DTO to a domain model.
     *
     * @param dto the OpenAPI visibility group DTO
     * @return the domain model
     */
    VisibilityGroup fromDTO(TgvalidatordInternalVisibilityGroup dto);

    /**
     * Converts a list of visibility group DTOs to domain models.
     *
     * @param dtos the list of OpenAPI visibility group DTOs
     * @return the list of domain models
     */
    List<VisibilityGroup> fromDTOList(List<TgvalidatordInternalVisibilityGroup> dtos);

    /**
     * Converts a visibility group user DTO to a domain model.
     *
     * @param dto the OpenAPI visibility group user DTO
     * @return the domain model
     */
    VisibilityGroupUser fromUserDTO(TgvalidatordInternalVisibilityGroupUser dto);

    /**
     * Converts a list of visibility group user DTOs to domain models.
     *
     * @param dtos the list of OpenAPI visibility group user DTOs
     * @return the list of domain models
     */
    List<VisibilityGroupUser> fromUserDTOList(List<TgvalidatordInternalVisibilityGroupUser> dtos);
}
