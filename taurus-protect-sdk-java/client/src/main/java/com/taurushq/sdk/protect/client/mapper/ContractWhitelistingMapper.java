package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Approvers;
import com.taurushq.sdk.protect.client.model.ApproversGroup;
import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedContractAddress;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedContractAddressEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistMetadata;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistTrail;
import com.taurushq.sdk.protect.client.model.WhitelistedContractAddressResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApprovers;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproversGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddressEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTrail;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistedContractAddressAttribute;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting contract whitelisting OpenAPI DTOs to client model objects.
 * <p>
 * This mapper handles the conversion of whitelisted contract address data, including
 * envelopes, signatures, trails, and attributes from OpenAPI generated models to the
 * SDK's clean domain models.
 */
@Mapper
public interface ContractWhitelistingMapper {

    /**
     * Singleton instance of the mapper.
     */
    ContractWhitelistingMapper INSTANCE = Mappers.getMapper(ContractWhitelistingMapper.class);

    /**
     * Maps a signed whitelisted contract address envelope from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    SignedWhitelistedContractAddressEnvelope fromEnvelopeDTO(
            TgvalidatordSignedWhitelistedContractAddressEnvelope dto);

    /**
     * Maps a list of signed whitelisted contract address envelopes.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<SignedWhitelistedContractAddressEnvelope> fromEnvelopeDTOList(
            List<TgvalidatordSignedWhitelistedContractAddressEnvelope> dtos);

    /**
     * Maps a signed whitelisted contract address from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    SignedWhitelistedContractAddress fromSignedContractDTO(
            TgvalidatordSignedWhitelistedContractAddress dto);

    /**
     * Maps the envelopes reply to a result with pagination info.
     *
     * @param reply the OpenAPI reply
     * @return the domain model result
     */
    @Mapping(target = "contracts", source = "result")
    @Mapping(target = "totalItems", source = "totalItems", qualifiedByName = "stringToLong")
    WhitelistedContractAddressResult fromEnvelopesReply(
            TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply reply);

    /**
     * Converts a string to long, handling null values.
     *
     * @param value the string value
     * @return the long value, or 0 if null
     */
    @Named("stringToLong")
    default long stringToLong(String value) {
        if (value == null || value.isEmpty()) {
            return 0L;
        }
        return Long.parseLong(value);
    }

    /**
     * Maps a whitelist trail from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    WhitelistTrail fromTrailDTO(TgvalidatordTrail dto);

    /**
     * Maps a list of whitelist trails.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<WhitelistTrail> fromTrailDTOList(List<TgvalidatordTrail> dtos);

    /**
     * Maps whitelist metadata from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    WhitelistMetadata fromMetadataDTO(TgvalidatordMetadata dto);

    /**
     * Maps approvers from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    Approvers fromApproversDTO(TgvalidatordApprovers dto);

    /**
     * Maps an approvers group from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    ApproversGroup fromApproversGroupDTO(TgvalidatordApproversGroup dto);

    /**
     * Maps a whitelist signature from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    WhitelistSignature fromSignatureDTO(TgvalidatordWhitelistSignature dto);

    /**
     * Maps a list of whitelist signatures.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<WhitelistSignature> fromSignatureDTOList(List<TgvalidatordWhitelistSignature> dtos);

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
