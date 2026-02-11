package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Approvers;
import com.taurushq.sdk.protect.client.model.ApproversGroup;
import com.taurushq.sdk.protect.client.model.Attribute;
import com.taurushq.sdk.protect.client.model.InternalAddress;
import com.taurushq.sdk.protect.client.model.InternalWallet;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddress;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistMetadata;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistTrail;
import com.taurushq.sdk.protect.client.model.WhitelistUserSignature;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApprovers;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordApproversGroup;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTrail;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistUserSignature;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWhitelistedAddressAttribute;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.proto.v1.Whitelist;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.Base64;
import org.mapstruct.factory.Mappers;

/**
 * The interface Address mapper.
 */
@Mapper(uses = {ScoreMapper.class})
public interface WhitelistedAddressMapper {
    /**
     * The constant INSTANCE.
     */
    WhitelistedAddressMapper INSTANCE = Mappers.getMapper(WhitelistedAddressMapper.class);

    /**
     * From dto signed whitelisted address.
     *
     * @param signedWhitelistedAddress the signed whitelisted address
     * @return the whitelisted address
     */
    SignedWhitelistedAddress fromDTO(TgvalidatordSignedWhitelistedAddress signedWhitelistedAddress);

    /**
     * From dto signed whitelisted address envelope.
     *
     * @param envelope the signed whitelisted address envelope
     * @return the signed whitelisted address envelope
     */
    @Mapping(target = "tenantId", expression = "java(envelope.getTenantId() != null ? Long.parseLong(envelope.getTenantId()) : 0L)")
    @Mapping(target = "id", expression = "java(envelope.getId() != null ? Long.parseLong(envelope.getId()) : 0L)")
    SignedWhitelistedAddressEnvelope fromDTO(TgvalidatordSignedWhitelistedAddressEnvelope envelope);


    /**
     * From dto whitelist signature.
     *
     * @param whitelistSignature the whitelist signature
     * @return the whitelist signature
     */
    WhitelistSignature fromDTO(TgvalidatordWhitelistSignature whitelistSignature);


    /**
     * From dto whitelist user signature.
     *
     * @param whitelistUserSignature the whitelist user signature
     * @return the whitelisted user signature
     */
    WhitelistUserSignature fromDTO(TgvalidatordWhitelistUserSignature whitelistUserSignature);


    /**
     * From protobuf whitelisted address.
     *
     * @param whitelistedAddress the protobuf whitelisted address
     * @return the whitelisted address
     */
    @Mapping(source = "linkedWalletsList", target = "linkedWallets")
    @Mapping(source = "linkedInternalAddressesList", target = "linkedInternalAddresses")
    WhitelistedAddress fromDTO(Whitelist.WhitelistedAddress whitelistedAddress);

    /**
     * From protobuf internal wallet.
     *
     * @param internalWallet the protobuf internal wallet
     * @return the internal wallet
     */
    InternalWallet fromDTO(Whitelist.InternalWallet internalWallet);

    /**
     * From protobuf internal address.
     *
     * @param internalAddress the protobuf internal address
     * @return the internal address
     */
    InternalAddress fromDTO(Whitelist.InternalAddress internalAddress);

    /**
     * From dto whitelist trail.
     *
     * @param trail the trail
     * @return the whitelist trail
     */
    WhitelistTrail fromDTO(TgvalidatordTrail trail);

    /**
     * From dto whitelist metadata.
     *
     * @param metadata the metadata
     * @return the whitelist metadata
     */
    WhitelistMetadata fromDTO(TgvalidatordMetadata metadata);

    /**
     * From dto approvers.
     *
     * @param approvers the approvers
     * @return the approvers
     */
    Approvers fromDTO(TgvalidatordApprovers approvers);

    /**
     * From dto approvers group.
     *
     * @param approversGroup the approvers group
     * @return the approvers group
     */
    ApproversGroup fromDTO(TgvalidatordApproversGroup approversGroup);

    /**
     * From dto whitelisted address attribute.
     *
     * @param attribute the attribute
     * @return the attribute
     */
    @Mapping(target = "subType", source = "subtype")
    @Mapping(target = "isFile", source = "isfile")
    Attribute fromDTO(TgvalidatordWhitelistedAddressAttribute attribute);

    /**
     * Decodes a WhitelistedAddress protobuf from raw bytes.
     *
     * @param data the raw protobuf bytes
     * @return the decoded WhitelistedAddress
     * @throws InvalidProtocolBufferException if the data cannot be parsed
     */
    default WhitelistedAddress fromBytes(byte[] data) throws InvalidProtocolBufferException {
        return fromDTO(Whitelist.WhitelistedAddress.parseFrom(data));
    }

    /**
     * Decodes a WhitelistedAddress protobuf from a base64-encoded string.
     *
     * @param base64 the base64-encoded protobuf data
     * @return the decoded WhitelistedAddress
     * @throws InvalidProtocolBufferException if the data cannot be parsed
     */
    default WhitelistedAddress fromBase64String(String base64) throws InvalidProtocolBufferException {
        byte[] data = Base64.getDecoder().decode(base64);
        return fromBytes(data);
    }

}


