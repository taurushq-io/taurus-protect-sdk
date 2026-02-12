package com.taurushq.sdk.protect.client.mapper;

import com.google.protobuf.ByteString;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingLine;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.CashSettlement;
import com.taurushq.sdk.protect.client.model.rulescontainer.ContractAddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.EvmCallContract;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleColumn;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleLine;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSource;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceInternalWallet;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceType;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.client.model.rulescontainer.SequentialThresholds;
import com.taurushq.sdk.protect.client.model.rulescontainer.TransactionRuleDetails;
import com.taurushq.sdk.protect.client.model.rulescontainer.TransactionRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.XtzCallContract;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import com.taurushq.sdk.protect.proto.v1.RequestReply;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;
import org.mapstruct.factory.Mappers;

import java.security.PublicKey;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Collections;
import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.stream.Collectors;

/**
 * MapStruct mapper for converting protobuf RulesContainer to client model objects.
 */
@Mapper
public interface RulesContainerMapper {

    /**
     * Logger for this mapper.
     */
    Logger LOGGER = Logger.getLogger(RulesContainerMapper.class.getName());

    /**
     * The constant INSTANCE.
     */
    RulesContainerMapper INSTANCE = Mappers.getMapper(RulesContainerMapper.class);

    /**
     * Converts a protobuf RulesContainer to DecodedRulesContainer.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "usersList", target = "users")
    @Mapping(source = "groupsList", target = "groups")
    @Mapping(source = "transactionRulesList", target = "transactionRules")
    @Mapping(source = "addressWhitelistingRulesList", target = "addressWhitelistingRules")
    @Mapping(source = "contractAddressWhitelistingRulesList", target = "contractAddressWhitelistingRules")
    @Mapping(source = "engineIdentitiesList", target = "engineIdentities")
    DecodedRulesContainer fromProto(RequestReply.RulesContainer proto);

    /**
     * Converts a protobuf User to RuleUser.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "publicKeyPem", source = "publicKey")
    @Mapping(target = "publicKey", source = "publicKey", qualifiedByName = "pemToPublicKey")
    @Mapping(target = "roles", source = "rolesList", qualifiedByName = "rolesToStrings")
    RuleUser fromProto(RequestReply.User proto);

    /**
     * Converts a protobuf Group to RuleGroup.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "userIdsList", target = "userIds")
    RuleGroup fromProto(RequestReply.Group proto);

    /**
     * Converts a protobuf GroupThreshold to GroupThreshold.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    GroupThreshold fromProto(RequestReply.GroupThreshold proto);

    /**
     * Converts a protobuf SequentialThresholds to SequentialThresholds.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "thresholdsList", target = "thresholds")
    SequentialThresholds fromProto(RequestReply.SequentialThresholds proto);

    /**
     * Converts a protobuf TransactionRules to TransactionRules.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "columnsList", target = "columns")
    @Mapping(source = "linesList", target = "lines")
    TransactionRules fromProto(RequestReply.RulesContainer.TransactionRules proto);

    /**
     * Converts a protobuf Column to RuleColumn.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "type", source = "type", qualifiedByName = "columnTypeToString")
    RuleColumn fromProto(RequestReply.RulesContainer.Column proto);

    /**
     * Converts a protobuf Line to RuleLine.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "cellsList", target = "cells")
    @Mapping(source = "parallelThresholdsList", target = "parallelThresholds")
    RuleLine fromProto(RequestReply.RulesContainer.Line proto);

    /**
     * Converts a protobuf TransactionRuleDetails to TransactionRuleDetails.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "domain", source = "domain", qualifiedByName = "ruleDomainToString")
    @Mapping(target = "subDomain", source = "subDomain", qualifiedByName = "ruleSubDomainToString")
    @Mapping(target = "cosmosMethodSignatures", source = "cosmosDetails.methodSignaturesList")
    TransactionRuleDetails fromProto(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails proto);

    /**
     * Converts a protobuf EvmCallContract to EvmCallContract.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    EvmCallContract fromProto(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.EvmCallContract proto);

    /**
     * Converts a protobuf XtzCallContract to XtzCallContract.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    XtzCallContract fromProto(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.XtzCallContract proto);

    /**
     * Converts a protobuf CashSettlement to CashSettlement.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    CashSettlement fromProto(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.CashSettlement proto);

    /**
     * Converts a protobuf AddressWhitelistingRules to AddressWhitelistingRules.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "parallelThresholdsList", target = "parallelThresholds")
    @Mapping(source = "linesList", target = "lines")
    AddressWhitelistingRules fromProto(RequestReply.RulesContainer.AddressWhitelistingRules proto);

    /**
     * Converts a protobuf AddressWhitelistingRules.Line to AddressWhitelistingLine.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "cellsList", target = "cells", qualifiedByName = "cellBytesToRuleSources")
    @Mapping(source = "parallelThresholdsList", target = "parallelThresholds")
    AddressWhitelistingLine fromProto(RequestReply.RulesContainer.AddressWhitelistingRules.Line proto);

    /**
     * Converts a protobuf ContractAddressWhitelistingRules to ContractAddressWhitelistingRules.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "parallelThresholdsList", target = "parallelThresholds")
    @Mapping(target = "blockchain", source = "blockchain", qualifiedByName = "blockchainToString")
    ContractAddressWhitelistingRules fromProto(RequestReply.RulesContainer.ContractAddressWhitelistingRules proto);

    /**
     * Converts a PEM string to a PublicKey object.
     *
     * @param pem the PEM string
     * @return the public key, or null if conversion fails
     */
    @Named("pemToPublicKey")
    default PublicKey pemToPublicKey(String pem) {
        if (pem == null || pem.isEmpty()) {
            return null;
        }
        try {
            return CryptoTPV1.decodePublicKey(pem);
        } catch (Exception e) {
            LOGGER.log(Level.WARNING, "Failed to decode PEM public key: " + e.getMessage(), e);
            return null;
        }
    }

    /**
     * Converts a list of Role enums to strings.
     *
     * @param roles the roles
     * @return the role names as strings
     */
    @Named("rolesToStrings")
    default List<String> rolesToStrings(List<RequestReply.Role> roles) {
        if (roles == null) {
            return Collections.emptyList();
        }
        return roles.stream()
                .map(RequestReply.Role::name)
                .collect(Collectors.toList());
    }

    /**
     * Converts a ColumnType enum to string.
     *
     * @param type the column type
     * @return the column type name
     */
    @Named("columnTypeToString")
    default String columnTypeToString(RequestReply.RulesContainer.ColumnType type) {
        return type != null ? type.name() : null;
    }

    /**
     * Converts a RuleDomain enum to string.
     *
     * @param domain the rule domain
     * @return the domain name
     */
    @Named("ruleDomainToString")
    default String ruleDomainToString(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.RuleDomain domain) {
        return domain != null ? domain.name() : null;
    }

    /**
     * Converts a RuleSubDomain enum to string.
     *
     * @param subDomain the rule sub-domain
     * @return the sub-domain name
     */
    @Named("ruleSubDomainToString")
    default String ruleSubDomainToString(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.RuleSubDomain subDomain) {
        return subDomain != null ? subDomain.name() : null;
    }

    /**
     * Converts a Blockchain enum to string.
     *
     * @param blockchain the blockchain
     * @return the blockchain name
     */
    @Named("blockchainToString")
    default String blockchainToString(RequestReply.Blockchain blockchain) {
        return blockchain != null ? blockchain.name() : null;
    }

    /**
     * Converts a list of ByteString to list of ByteString (passthrough).
     *
     * @param bytes the bytes
     * @return the same bytes
     */
    default List<ByteString> mapByteStrings(List<ByteString> bytes) {
        return bytes;
    }

    /**
     * Decodes a protobuf RulesContainer from raw bytes.
     *
     * @param data the raw protobuf bytes
     * @return the decoded DecodedRulesContainer
     * @throws InvalidProtocolBufferException if the data cannot be parsed
     */
    default DecodedRulesContainer fromBytes(byte[] data) throws InvalidProtocolBufferException {
        return fromProto(RequestReply.RulesContainer.parseFrom(data));
    }

    /**
     * Decodes a protobuf RulesContainer from a base64-encoded string.
     *
     * @param base64 the base64-encoded protobuf data
     * @return the decoded DecodedRulesContainer
     * @throws InvalidProtocolBufferException if the data cannot be parsed
     */
    default DecodedRulesContainer fromBase64String(String base64) throws InvalidProtocolBufferException {
        byte[] data = Base64.getDecoder().decode(base64);
        return fromBytes(data);
    }

    /**
     * Decodes UserSignatures protobuf from raw bytes into a list of RuleUserSignature model objects.
     *
     * @param data the raw protobuf bytes
     * @return the list of decoded RuleUserSignature objects
     * @throws InvalidProtocolBufferException if the data cannot be parsed
     */
    default List<RuleUserSignature> userSignaturesFromBytes(byte[] data) throws InvalidProtocolBufferException {
        RequestReply.UserSignatures proto = RequestReply.UserSignatures.parseFrom(data);
        return proto.getSignaturesList().stream()
                .map(sig -> {
                    RuleUserSignature r = new RuleUserSignature();
                    r.setUserId(sig.getUserId());
                    r.setSignature(Base64.getEncoder().encodeToString(sig.getSignature().toByteArray()));
                    return r;
                })
                .collect(Collectors.toList());
    }

    /**
     * Decodes UserSignatures protobuf from a base64-encoded string into a list of RuleUserSignature model objects.
     *
     * @param base64 the base64-encoded protobuf data
     * @return the list of decoded RuleUserSignature objects
     * @throws InvalidProtocolBufferException if the data cannot be parsed
     */
    default List<RuleUserSignature> userSignaturesFromBase64String(String base64) throws InvalidProtocolBufferException {
        byte[] data = Base64.getDecoder().decode(base64);
        return userSignaturesFromBytes(data);
    }

    /**
     * Converts a list of cell ByteStrings to a list of RuleSource objects.
     * Each cell is a serialized RuleSource protobuf message.
     *
     * @param cells the cell bytes
     * @return the list of decoded RuleSource objects
     */
    @Named("cellBytesToRuleSources")
    default List<RuleSource> cellBytesToRuleSources(List<ByteString> cells) {
        if (cells == null) {
            return Collections.emptyList();
        }
        List<RuleSource> result = new ArrayList<>();
        for (ByteString cellBytes : cells) {
            try {
                RequestReply.RuleSource protoSource = RequestReply.RuleSource.parseFrom(cellBytes);
                RuleSource source = new RuleSource();
                source.setType(RuleSourceType.fromValue(protoSource.getTypeValue()));

                // Decode payload based on type
                if (source.getType() == RuleSourceType.RuleSourceInternalWallet
                        && !protoSource.getPayload().isEmpty()) {
                    RequestReply.RuleSourceInternalWallet protoWallet =
                            RequestReply.RuleSourceInternalWallet.parseFrom(protoSource.getPayload());
                    RuleSourceInternalWallet wallet = new RuleSourceInternalWallet();
                    wallet.setPath(protoWallet.getPath());
                    source.setInternalWallet(wallet);
                }

                result.add(source);
            } catch (InvalidProtocolBufferException ignored) {
                // Skip malformed cells - continue processing remaining cells
            }
        }
        return result;
    }
}
