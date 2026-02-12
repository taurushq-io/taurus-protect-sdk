package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.TokenMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordERCTokenMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFATokenMetadata;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for converting token metadata OpenAPI DTOs to client model objects.
 */
@Mapper
public interface TokenMetadataMapper {

    /**
     * Singleton instance of the mapper.
     */
    TokenMetadataMapper INSTANCE = Mappers.getMapper(TokenMetadataMapper.class);

    /**
     * Maps ERC token metadata from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    TokenMetadata fromERCDTO(TgvalidatordERCTokenMetadata dto);

    /**
     * Maps FA token metadata from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    TokenMetadata fromFADTO(TgvalidatordFATokenMetadata dto);
}
