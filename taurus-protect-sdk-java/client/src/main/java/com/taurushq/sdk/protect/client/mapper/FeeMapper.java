package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Fee;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordKeyValue;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting fee OpenAPI DTOs to client model objects.
 */
@Mapper
public interface FeeMapper {

    /**
     * Singleton instance of the mapper.
     */
    FeeMapper INSTANCE = Mappers.getMapper(FeeMapper.class);

    /**
     * Maps a fee from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    Fee fromDTO(TgvalidatordKeyValue dto);

    /**
     * Maps a list of fees from DTOs.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<Fee> fromDTOList(List<TgvalidatordKeyValue> dtos);
}
