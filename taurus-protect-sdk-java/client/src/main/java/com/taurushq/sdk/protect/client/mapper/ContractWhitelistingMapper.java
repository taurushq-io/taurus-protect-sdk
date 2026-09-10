package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistedContractAddressAttribute;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for whitelisted contract attributes.
 * <p>
 * Attribute mappers only. The envelope, signed-contract and result mappers that lived here
 * built domain objects from the DTO with no verification of any kind, which made
 * WhitelistedAssetService's six-step chain avoidable for the same entity. Do not reintroduce
 * a mapper that takes a bare contract envelope.
 */
@Mapper
@SuppressWarnings("PMD.ConstantsInInterface")
public interface ContractWhitelistingMapper {

    /**
     * Singleton instance of the mapper.
     */
    ContractWhitelistingMapper INSTANCE = Mappers.getMapper(ContractWhitelistingMapper.class);

    /**
     * Maps a whitelisted contract address attribute from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    @Mapping(target = "subType", source = "subtype")
    @Mapping(target = "isFile", source = "isfile")
    Attribute fromAttributeDTO(TgvalidatordWhitelistedContractAddressAttribute dto);

    /**
     * Maps a list of whitelisted contract address attributes.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<Attribute> fromAttributeDTOList(List<TgvalidatordWhitelistedContractAddressAttribute> dtos);
}
