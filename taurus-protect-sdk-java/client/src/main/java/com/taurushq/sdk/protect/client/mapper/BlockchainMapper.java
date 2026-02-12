package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.BlockchainInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBlockchainEntity;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting blockchain OpenAPI DTOs to client model objects.
 */
@Mapper(uses = CurrencyMapper.class)
public interface BlockchainMapper {

    /**
     * Singleton instance of the mapper.
     */
    BlockchainMapper INSTANCE = Mappers.getMapper(BlockchainMapper.class);

    /**
     * Maps a blockchain entity from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    BlockchainInfo fromDTO(TgvalidatordBlockchainEntity dto);

    /**
     * Maps a list of blockchain entities from DTOs.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<BlockchainInfo> fromDTOList(List<TgvalidatordBlockchainEntity> dtos);
}
