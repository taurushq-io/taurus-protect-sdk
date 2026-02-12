package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Tag;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTag;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting between OpenAPI tag DTOs and client model objects.
 */
@Mapper
public interface TagMapper {

    /**
     * Singleton instance of the mapper.
     */
    TagMapper INSTANCE = Mappers.getMapper(TagMapper.class);

    /**
     * Converts an OpenAPI tag DTO to a client model Tag.
     *
     * @param dto the OpenAPI DTO
     * @return the client model Tag
     */
    Tag fromDTO(TgvalidatordTag dto);

    /**
     * Converts a list of OpenAPI tag DTOs to client model Tags.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of client model Tags
     */
    List<Tag> fromDTOList(List<TgvalidatordTag> dtos);
}
