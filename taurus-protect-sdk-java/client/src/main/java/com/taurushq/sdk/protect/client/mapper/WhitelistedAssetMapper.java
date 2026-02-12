package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Approvers;
import com.taurushq.sdk.protect.client.model.ApproversGroup;
import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAsset;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAssetEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistMetadata;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistTrail;
import com.taurushq.sdk.protect.client.model.WhitelistUserSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApprovers;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproversGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedContractAddressEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTrail;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistUserSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistedContractAddressAttribute;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for whitelisted assets (contract addresses).
 */
@Mapper
public interface WhitelistedAssetMapper {

    /**
     * The constant INSTANCE.
     */
    WhitelistedAssetMapper INSTANCE = Mappers.getMapper(WhitelistedAssetMapper.class);

    /**
     * Maps from OpenAPI signed contract address to client model.
     *
     * @param signedContractAddress the OpenAPI signed contract address
     * @return the client model signed whitelisted asset
     */
    SignedWhitelistedAsset fromDTO(TgvalidatordSignedWhitelistedContractAddress signedContractAddress);

    /**
     * Maps from OpenAPI signed contract address envelope to client model.
     *
     * @param envelope the OpenAPI envelope
     * @return the client model envelope
     */
    @Mapping(target = "tenantId",
            expression = "java(envelope.getTenantId() != null ? Long.parseLong(envelope.getTenantId()) : 0L)")
    @Mapping(target = "id",
            expression = "java(envelope.getId() != null ? Long.parseLong(envelope.getId()) : 0L)")
    @Mapping(source = "signedContractAddress", target = "signedAsset")
    SignedWhitelistedAssetEnvelope fromDTO(TgvalidatordSignedWhitelistedContractAddressEnvelope envelope);

    /**
     * Maps from OpenAPI whitelist signature to client model.
     *
     * @param whitelistSignature the OpenAPI whitelist signature
     * @return the client model whitelist signature
     */
    WhitelistSignature fromDTO(TgvalidatordWhitelistSignature whitelistSignature);

    /**
     * Maps from OpenAPI whitelist user signature to client model.
     *
     * @param whitelistUserSignature the OpenAPI whitelist user signature
     * @return the client model whitelist user signature
     */
    WhitelistUserSignature fromDTO(TgvalidatordWhitelistUserSignature whitelistUserSignature);

    /**
     * Maps from OpenAPI trail to client model.
     *
     * @param trail the OpenAPI trail
     * @return the client model whitelist trail
     */
    WhitelistTrail fromDTO(TgvalidatordTrail trail);

    /**
     * Maps from OpenAPI metadata to client model.
     *
     * @param metadata the OpenAPI metadata
     * @return the client model whitelist metadata
     */
    WhitelistMetadata fromDTO(TgvalidatordMetadata metadata);

    /**
     * Maps from OpenAPI approvers to client model.
     *
     * @param approvers the OpenAPI approvers
     * @return the client model approvers
     */
    Approvers fromDTO(TgvalidatordApprovers approvers);

    /**
     * Maps from OpenAPI approvers group to client model.
     *
     * @param approversGroup the OpenAPI approvers group
     * @return the client model approvers group
     */
    ApproversGroup fromDTO(TgvalidatordApproversGroup approversGroup);

    /**
     * Maps from OpenAPI whitelisted contract address attribute to client model.
     *
     * @param attribute the OpenAPI attribute
     * @return the client model attribute
     */
    @Mapping(target = "subType", source = "subtype")
    @Mapping(target = "isFile", source = "isfile")
    Attribute fromDTO(TgvalidatordWhitelistedContractAddressAttribute attribute);
}
