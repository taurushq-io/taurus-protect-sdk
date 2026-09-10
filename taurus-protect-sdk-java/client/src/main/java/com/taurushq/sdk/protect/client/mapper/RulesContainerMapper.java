package com.taurushq.sdk.protect.client.mapper;

import com.google.protobuf.ByteString;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingLine;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.CashSettlement;
import com.taurushq.sdk.protect.client.model.rulescontainer.ContractAddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.CosmosDetails;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.EvmCallContract;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleColumn;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleLine;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSource;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceExchange;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceExternalAddress;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceInternalAddress;
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
    @Mapping(target = "properties", expression = "java(proto.getPropertiesMap())")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    DecodedRulesContainer fromProto(RequestReply.RulesContainer proto);

    /**
     * Converts a protobuf User to RuleUser.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "publicKeyPem", source = "publicKey")
    @Mapping(target = "publicKey", source = "publicKey", qualifiedByName = "pemToPublicKey")
    @Mapping(target = "roles", source = "rolesValueList", qualifiedByName = "rolesToStrings")
    @Mapping(target = "properties", expression = "java(proto.getPropertiesMap())")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    RuleUser fromProto(RequestReply.User proto);

    /**
     * Converts a protobuf Group to RuleGroup.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "userIdsList", target = "userIds")
    @Mapping(target = "properties", expression = "java(proto.getPropertiesMap())")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    RuleGroup fromProto(RequestReply.Group proto);

    /**
     * Converts a protobuf GroupThreshold to GroupThreshold.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    GroupThreshold fromProto(RequestReply.GroupThreshold proto);

    /**
     * Converts a protobuf SequentialThresholds to SequentialThresholds.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "thresholdsList", target = "thresholds")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    SequentialThresholds fromProto(RequestReply.SequentialThresholds proto);

    /**
     * Converts a protobuf TransactionRules to TransactionRules.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "columnsList", target = "columns")
    @Mapping(source = "linesList", target = "lines")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    TransactionRules fromProto(RequestReply.RulesContainer.TransactionRules proto);

    /**
     * Converts a protobuf Column to RuleColumn.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "type", source = "typeValue", qualifiedByName = "columnTypeToString")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    RuleColumn fromProto(RequestReply.RulesContainer.Column proto);

    /**
     * Converts a protobuf Line to RuleLine.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "cellsList", target = "cells")
    @Mapping(source = "parallelThresholdsList", target = "parallelThresholds")
    @Mapping(target = "properties", expression = "java(proto.getPropertiesMap())")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    RuleLine fromProto(RequestReply.RulesContainer.Line proto);

    /**
     * Converts a protobuf TransactionRuleDetails to TransactionRuleDetails.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "domain", source = "domainValue", qualifiedByName = "ruleDomainToString")
    @Mapping(target = "subDomain", source = "subDomainValue", qualifiedByName = "ruleSubDomainToString")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    TransactionRuleDetails fromProto(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails proto);

    /**
     * Converts a protobuf EvmCallContract to EvmCallContract.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    EvmCallContract fromProto(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.EvmCallContract proto);

    /**
     * Converts a protobuf XtzCallContract to XtzCallContract.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    XtzCallContract fromProto(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.XtzCallContract proto);

    /**
     * Converts a protobuf CashSettlement to CashSettlement.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    CashSettlement fromProto(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.CashSettlement proto);

    /**
     * Converts a protobuf CosmosDetails to CosmosDetails.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "methodSignaturesList", target = "methodSignatures")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    CosmosDetails fromProto(RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.CosmosDetails proto);

    /**
     * Converts a protobuf AddressWhitelistingRules to AddressWhitelistingRules.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "parallelThresholdsList", target = "parallelThresholds")
    @Mapping(source = "linesList", target = "lines")
    @Mapping(target = "properties", expression = "java(proto.getPropertiesMap())")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    AddressWhitelistingRules fromProto(RequestReply.RulesContainer.AddressWhitelistingRules proto);

    /**
     * Converts a protobuf AddressWhitelistingRules.Line to AddressWhitelistingLine.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "cellsList", target = "cells", qualifiedByName = "cellBytesToRuleSources")
    @Mapping(source = "parallelThresholdsList", target = "parallelThresholds")
    @Mapping(target = "properties", expression = "java(proto.getPropertiesMap())")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
    AddressWhitelistingLine fromProto(RequestReply.RulesContainer.AddressWhitelistingRules.Line proto);

    /**
     * Converts a protobuf ContractAddressWhitelistingRules to ContractAddressWhitelistingRules.
     *
     * @param proto the protobuf object
     * @return the client model
     */
    @Mapping(source = "parallelThresholdsList", target = "parallelThresholds")
    @Mapping(target = "blockchain", source = "blockchainValue", qualifiedByName = "blockchainToString")
    @Mapping(target = "properties", expression = "java(proto.getPropertiesMap())")
    @Mapping(target = "unknownFields", expression = "java(unknownBytes(proto))")
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
            if (LOGGER.isLoggable(Level.WARNING)) {
                LOGGER.log(Level.WARNING, "Failed to decode PEM public key: " + e.getMessage(), e);
            }
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
    default List<String> rolesToStrings(List<Integer> roleValues) {
        if (roleValues == null) {
            return Collections.emptyList();
        }
        return roleValues.stream()
                .map(v -> {
                    RequestReply.Role r = RequestReply.Role.forNumber(v);
                    return r != null ? r.name() : Integer.toString(v);
                })
                .collect(Collectors.toList());
    }

    /**
     * Converts a ColumnType enum to string.
     *
     * @param type the column type
     * @return the column type name
     */
    @Named("columnTypeToString")
    default String columnTypeToString(int typeValue) {
        RequestReply.RulesContainer.ColumnType t =
                RequestReply.RulesContainer.ColumnType.forNumber(typeValue);
        return t != null ? t.name() : Integer.toString(typeValue);
    }

    /**
     * Converts a RuleDomain enum to string.
     *
     * @param domain the rule domain
     * @return the domain name
     */
    @Named("ruleDomainToString")
    default String ruleDomainToString(int domainValue) {
        RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.RuleDomain d =
                RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.RuleDomain.forNumber(domainValue);
        return d != null ? d.name() : Integer.toString(domainValue);
    }

    /**
     * Converts a RuleSubDomain enum to string.
     *
     * @param subDomain the rule sub-domain
     * @return the sub-domain name
     */
    @Named("ruleSubDomainToString")
    default String ruleSubDomainToString(int subDomainValue) {
        RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.RuleSubDomain d =
                RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.RuleSubDomain.forNumber(subDomainValue);
        return d != null ? d.name() : Integer.toString(subDomainValue);
    }

    /**
     * Converts a Blockchain enum to string.
     *
     * @param blockchain the blockchain
     * @return the blockchain name
     */
    @Named("blockchainToString")
    default String blockchainToString(int blockchainValue) {
        RequestReply.Blockchain b = RequestReply.Blockchain.forNumber(blockchainValue);
        return b != null ? b.name() : Integer.toString(blockchainValue);
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
            result.add(ruleSourceFromBytes(cellBytes));
        }
        return result;
    }

    /**
     * Decodes a single address-whitelisting source cell into a typed {@link RuleSource}.
     * Cells that do not round-trip byte-identically (unknown type, unknown sub-fields,
     * malformed) retain their verbatim bytes so they can be re-encoded losslessly.
     *
     * @param cellBytes the serialized RuleSource cell
     * @return the decoded source
     */
    default RuleSource ruleSourceFromBytes(ByteString cellBytes) {
        RuleSource source = new RuleSource();
        try {
            RequestReply.RuleSource p = RequestReply.RuleSource.parseFrom(cellBytes);
            source.setType(RuleSourceType.fromValue(p.getTypeValue()));
            ByteString payload = p.getPayload();
            // An empty payload leaves the typed variant unset, as in the other three
            // SDKs: parsing it would materialize a default sub-message, so a caller
            // checking getInternalWallet() != null would see a wallet with an empty path
            // where Go/Python/TypeScript see nothing.
            switch (payload.isEmpty() ? -1 : p.getTypeValue()) {
                case 1: {
                    RequestReply.RuleSourceInternalWallet i = RequestReply.RuleSourceInternalWallet.parseFrom(payload);
                    RuleSourceInternalWallet w = new RuleSourceInternalWallet();
                    w.setPath(i.getPath());
                    source.setInternalWallet(w);
                    break;
                }
                case 2: {
                    RequestReply.RuleSourceInternalAddress i = RequestReply.RuleSourceInternalAddress.parseFrom(payload);
                    RuleSourceInternalAddress a = new RuleSourceInternalAddress();
                    a.setAddress(i.getAddress());
                    a.setPath(i.getPath());
                    source.setInternalAddress(a);
                    break;
                }
                case 4: {
                    RequestReply.RuleSourceExchange i = RequestReply.RuleSourceExchange.parseFrom(payload);
                    RuleSourceExchange x = new RuleSourceExchange();
                    x.setLabel(i.getLabel());
                    source.setExchange(x);
                    break;
                }
                case 5: {
                    RequestReply.RuleSourceExternalAddress i = RequestReply.RuleSourceExternalAddress.parseFrom(payload);
                    RuleSourceExternalAddress x = new RuleSourceExternalAddress();
                    x.setAddress(i.getAddress());
                    x.setMemo(i.getMemo());
                    source.setExternalAddress(x);
                    break;
                }
                default:
                    // -1 (empty payload), RuleSourceAny (0), RuleSourceAnyExchange (3):
                    // no payload to decode; other values are unknown to this SDK.
                    break;
            }
        } catch (InvalidProtocolBufferException e) {
            source.setType(RuleSourceType.UNRECOGNIZED);
            source.setRaw(cellBytes);
            return source;
        }
        // Lossless guard: keep verbatim bytes if the typed value does not re-encode exactly.
        if (!ruleSourceToBytes(source).equals(cellBytes)) {
            source.setRaw(cellBytes);
        }
        return source;
    }

    // ================= Encode (model -> protobuf) =================

    /**
     * Encodes a decoded rules container to a base64-encoded protobuf string suitable
     * for submission as a rules proposal. Server-controlled fields (enforced rules hash
     * and timestamp) are omitted.
     *
     * @param container the decoded rules container
     * @return the base64-encoded protobuf
     */
    default String toBase64String(DecodedRulesContainer container) {
        return Base64.getEncoder().encodeToString(toBytes(container));
    }

    /**
     * Encodes a decoded rules container to protobuf bytes. Server-controlled fields
     * (enforced rules hash and timestamp) are omitted.
     *
     * @param container the decoded rules container
     * @return the serialized protobuf bytes
     */
    default byte[] toBytes(DecodedRulesContainer container) {
        return deterministicBytes(toProto(container));
    }

    /**
     * Serializes a message with map entries in a stable order. The container carries
     * {@code map<string, bytes> properties} at five levels, and the default serializer
     * emits map entries in unspecified order — so the same reviewed container would
     * encode to different bytes across runs, and differently from the other SDKs.
     *
     * @param message the message to serialize
     * @return the serialized bytes, stable for a given message
     */
    default byte[] deterministicBytes(com.google.protobuf.Message message) {
        byte[] out = new byte[message.getSerializedSize()];
        com.google.protobuf.CodedOutputStream stream =
                com.google.protobuf.CodedOutputStream.newInstance(out);
        stream.useDeterministicSerialization();
        try {
            message.writeTo(stream);
            stream.checkNoSpaceLeft();
        } catch (java.io.IOException e) {
            throw new IllegalStateException("failed to serialize " + message.getClass().getName(), e);
        }
        return out;
    }

    /**
     * Converts a decoded rules container to its protobuf representation. Unknown
     * protobuf fields captured at decode are re-attached at each node so a
     * decode/encode round-trip never silently drops data. The server-controlled
     * enforced rules hash and timestamp are intentionally not set.
     *
     * @param container the decoded rules container
     * @return the protobuf RulesContainer
     */
    default RequestReply.RulesContainer toProto(DecodedRulesContainer container) {
        RequestReply.RulesContainer.Builder b = RequestReply.RulesContainer.newBuilder();
        if (container.getUsers() != null) {
            for (RuleUser u : container.getUsers()) {
                b.addUsers(toProtoUser(u));
            }
        }
        if (container.getGroups() != null) {
            for (RuleGroup g : container.getGroups()) {
                b.addGroups(toProtoGroup(g));
            }
        }
        b.setMinimumDistinctUserSignatures(container.getMinimumDistinctUserSignatures());
        b.setMinimumDistinctGroupSignatures(container.getMinimumDistinctGroupSignatures());
        if (container.getTransactionRules() != null) {
            for (TransactionRules t : container.getTransactionRules()) {
                b.addTransactionRules(toProtoTransactionRules(t));
            }
        }
        if (container.getAddressWhitelistingRules() != null) {
            for (AddressWhitelistingRules a : container.getAddressWhitelistingRules()) {
                b.addAddressWhitelistingRules(toProtoAddressWhitelisting(a));
            }
        }
        if (container.getContractAddressWhitelistingRules() != null) {
            for (ContractAddressWhitelistingRules c : container.getContractAddressWhitelistingRules()) {
                b.addContractAddressWhitelistingRules(toProtoContractWhitelisting(c));
            }
        }
        // enforcedRulesHash is server-controlled: intentionally not set.
        putProps(container.getProperties(), b::putAllProperties);
        // timestamp is server-controlled: intentionally not set.
        b.setMinimumCommitmentSignatures(container.getMinimumCommitmentSignatures());
        if (container.getEngineIdentities() != null) {
            b.addAllEngineIdentities(container.getEngineIdentities());
        }
        b.setHsmSlotId(container.getHsmSlotId());
        reattach(b, container.getUnknownFields());
        return b.build();
    }

    /**
     * Converts a rule user to protobuf.
     *
     * @param u the user
     * @return the protobuf user
     */
    default RequestReply.User toProtoUser(RuleUser u) {
        RequestReply.User.Builder b = RequestReply.User.newBuilder()
                .setId(nz(u.getId()))
                .setPublicKey(nz(u.getPublicKeyPem()))
                .addAllRolesValue(rolesFromStrings(u.getRoles()));
        putProps(u.getProperties(), b::putAllProperties);
        reattach(b, u.getUnknownFields());
        return b.build();
    }

    /**
     * Converts a rule group to protobuf.
     *
     * @param g the group
     * @return the protobuf group
     */
    default RequestReply.Group toProtoGroup(RuleGroup g) {
        RequestReply.Group.Builder b = RequestReply.Group.newBuilder().setId(nz(g.getId()));
        if (g.getUserIds() != null) {
            b.addAllUserIds(g.getUserIds());
        }
        putProps(g.getProperties(), b::putAllProperties);
        reattach(b, g.getUnknownFields());
        return b.build();
    }

    /**
     * Converts a group threshold to protobuf.
     *
     * @param t the group threshold
     * @return the protobuf group threshold
     */
    default RequestReply.GroupThreshold toProtoGroupThreshold(GroupThreshold t) {
        RequestReply.GroupThreshold.Builder b = RequestReply.GroupThreshold.newBuilder()
                .setGroupId(nz(t.getGroupId()))
                .setMinimumSignatures(t.getMinimumSignatures());
        reattach(b, t.getUnknownFields());
        return b.build();
    }

    /**
     * Converts sequential thresholds to protobuf.
     *
     * @param s the sequential thresholds
     * @return the protobuf sequential thresholds
     */
    default RequestReply.SequentialThresholds toProtoSequentialThresholds(SequentialThresholds s) {
        RequestReply.SequentialThresholds.Builder b = RequestReply.SequentialThresholds.newBuilder();
        if (s.getThresholds() != null) {
            for (GroupThreshold t : s.getThresholds()) {
                b.addThresholds(toProtoGroupThreshold(t));
            }
        }
        reattach(b, s.getUnknownFields());
        return b.build();
    }

    /**
     * Converts transaction rules to protobuf.
     *
     * @param t the transaction rules
     * @return the protobuf transaction rules
     */
    default RequestReply.RulesContainer.TransactionRules toProtoTransactionRules(TransactionRules t) {
        RequestReply.RulesContainer.TransactionRules.Builder b =
                RequestReply.RulesContainer.TransactionRules.newBuilder().setKey(nz(t.getKey()));
        if (t.getColumns() != null) {
            for (RuleColumn c : t.getColumns()) {
                b.addColumns(toProtoColumn(c));
            }
        }
        if (t.getLines() != null) {
            for (RuleLine l : t.getLines()) {
                b.addLines(toProtoLine(l));
            }
        }
        if (t.getDetails() != null) {
            b.setDetails(toProtoDetails(t.getDetails()));
        }
        reattach(b, t.getUnknownFields());
        return b.build();
    }

    /**
     * Converts a rule column to protobuf.
     *
     * @param c the column
     * @return the protobuf column
     */
    default RequestReply.RulesContainer.Column toProtoColumn(RuleColumn c) {
        RequestReply.RulesContainer.Column.Builder b = RequestReply.RulesContainer.Column.newBuilder()
                .setTypeValue(enumValue(RequestReply.RulesContainer.ColumnType.class, c.getType()))
                .setName(nz(c.getName()))
                .setMetadataKey(nz(c.getMetadataKey()));
        reattach(b, c.getUnknownFields());
        return b.build();
    }

    /**
     * Converts a rule line to protobuf.
     *
     * @param l the line
     * @return the protobuf line
     */
    default RequestReply.RulesContainer.Line toProtoLine(RuleLine l) {
        RequestReply.RulesContainer.Line.Builder b = RequestReply.RulesContainer.Line.newBuilder();
        if (l.getCells() != null) {
            b.addAllCells(l.getCells());
        }
        if (l.getParallelThresholds() != null) {
            for (SequentialThresholds s : l.getParallelThresholds()) {
                b.addParallelThresholds(toProtoSequentialThresholds(s));
            }
        }
        putProps(l.getProperties(), b::putAllProperties);
        b.setPriority(l.getPriority());
        reattach(b, l.getUnknownFields());
        return b.build();
    }

    /**
     * Converts transaction rule details to protobuf.
     *
     * @param d the details
     * @return the protobuf details
     */
    default RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails toProtoDetails(TransactionRuleDetails d) {
        RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.Builder b =
                RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.newBuilder()
                        .setDomainValue(enumValue(
                                RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.RuleDomain.class, d.getDomain()))
                        .setSubDomainValue(enumValue(
                                RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.RuleSubDomain.class, d.getSubDomain()))
                        .setBlockchain(nz(d.getBlockchain()))
                        .setNetwork(nz(d.getNetwork()));
        if (d.getEvmCallContract() != null) {
            RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.EvmCallContract.Builder evm =
                    RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.EvmCallContract.newBuilder()
                            .setContractType(nz(d.getEvmCallContract().getContractType()))
                            .setMethodSignature(nz(d.getEvmCallContract().getMethodSignature()));
            reattach(evm, d.getEvmCallContract().getUnknownFields());
            b.setEvmCallContract(evm);
        }
        if (d.getXtzCallContract() != null) {
            RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.XtzCallContract.Builder xtz =
                    RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.XtzCallContract.newBuilder()
                            .setContractType(nz(d.getXtzCallContract().getContractType()))
                            .setMethodSignature(nz(d.getXtzCallContract().getMethodSignature()));
            reattach(xtz, d.getXtzCallContract().getUnknownFields());
            b.setXtzCallContract(xtz);
        }
        if (d.getCashSettlement() != null) {
            RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.CashSettlement.Builder cash =
                    RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.CashSettlement.newBuilder()
                            .setProvider(nz(d.getCashSettlement().getProvider()))
                            .setRequestType(nz(d.getCashSettlement().getRequestType()));
            reattach(cash, d.getCashSettlement().getUnknownFields());
            b.setCashSettlement(cash);
        }
        if (d.getCosmosDetails() != null) {
            RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.CosmosDetails.Builder cosmos =
                    RequestReply.RulesContainer.TransactionRules.TransactionRuleDetails.CosmosDetails.newBuilder();
            if (d.getCosmosDetails().getMethodSignatures() != null) {
                cosmos.addAllMethodSignatures(d.getCosmosDetails().getMethodSignatures());
            }
            reattach(cosmos, d.getCosmosDetails().getUnknownFields());
            b.setCosmosDetails(cosmos);
        }
        reattach(b, d.getUnknownFields());
        return b.build();
    }

    /**
     * Converts address-whitelisting rules to protobuf.
     *
     * @param a the address-whitelisting rules
     * @return the protobuf address-whitelisting rules
     */
    default RequestReply.RulesContainer.AddressWhitelistingRules toProtoAddressWhitelisting(AddressWhitelistingRules a) {
        RequestReply.RulesContainer.AddressWhitelistingRules.Builder b =
                RequestReply.RulesContainer.AddressWhitelistingRules.newBuilder()
                        .setCurrency(nz(a.getCurrency()))
                        .setNetwork(nz(a.getNetwork()));
        if (a.getParallelThresholds() != null) {
            for (SequentialThresholds s : a.getParallelThresholds()) {
                b.addParallelThresholds(toProtoSequentialThresholds(s));
            }
        }
        if (a.getLines() != null) {
            for (AddressWhitelistingLine l : a.getLines()) {
                b.addLines(toProtoAddressWhitelistingLine(l));
            }
        }
        putProps(a.getProperties(), b::putAllProperties);
        reattach(b, a.getUnknownFields());
        return b.build();
    }

    /**
     * Converts an address-whitelisting line to protobuf.
     *
     * @param l the line
     * @return the protobuf line
     */
    default RequestReply.RulesContainer.AddressWhitelistingRules.Line toProtoAddressWhitelistingLine(AddressWhitelistingLine l) {
        RequestReply.RulesContainer.AddressWhitelistingRules.Line.Builder b =
                RequestReply.RulesContainer.AddressWhitelistingRules.Line.newBuilder();
        if (l.getCells() != null) {
            for (RuleSource s : l.getCells()) {
                b.addCells(ruleSourceToBytes(s));
            }
        }
        if (l.getParallelThresholds() != null) {
            for (SequentialThresholds s : l.getParallelThresholds()) {
                b.addParallelThresholds(toProtoSequentialThresholds(s));
            }
        }
        putProps(l.getProperties(), b::putAllProperties);
        reattach(b, l.getUnknownFields());
        return b.build();
    }

    /**
     * Converts contract-address-whitelisting rules to protobuf.
     *
     * @param c the contract-address-whitelisting rules
     * @return the protobuf contract-address-whitelisting rules
     */
    default RequestReply.RulesContainer.ContractAddressWhitelistingRules toProtoContractWhitelisting(ContractAddressWhitelistingRules c) {
        RequestReply.RulesContainer.ContractAddressWhitelistingRules.Builder b =
                RequestReply.RulesContainer.ContractAddressWhitelistingRules.newBuilder()
                        .setNetwork(nz(c.getNetwork()))
                        .setBlockchainValue(enumValue(RequestReply.Blockchain.class, c.getBlockchain()));
        if (c.getParallelThresholds() != null) {
            for (SequentialThresholds s : c.getParallelThresholds()) {
                b.addParallelThresholds(toProtoSequentialThresholds(s));
            }
        }
        putProps(c.getProperties(), b::putAllProperties);
        reattach(b, c.getUnknownFields());
        return b.build();
    }

    /**
     * Encodes a typed rule source to its serialized cell bytes.
     *
     * @param s the rule source
     * @return the serialized RuleSource cell
     */
    default ByteString ruleSourceToBytes(RuleSource s) {
        if (s.getRaw() != null) {
            return s.getRaw();
        }
        RuleSourceType type = s.getType() != null ? s.getType() : RuleSourceType.RuleSourceAny;
        RequestReply.RuleSource.Builder b = RequestReply.RuleSource.newBuilder().setTypeValue(type.getValue());
        byte[] payload = null;
        switch (type) {
            case RuleSourceInternalWallet:
                if (s.getInternalWallet() != null) {
                    payload = RequestReply.RuleSourceInternalWallet.newBuilder()
                            .setPath(nz(s.getInternalWallet().getPath())).build().toByteArray();
                }
                break;
            case RuleSourceInternalAddress:
                if (s.getInternalAddress() != null) {
                    payload = RequestReply.RuleSourceInternalAddress.newBuilder()
                            .setAddress(nz(s.getInternalAddress().getAddress()))
                            .setPath(nz(s.getInternalAddress().getPath())).build().toByteArray();
                }
                break;
            case RuleSourceExchange:
                if (s.getExchange() != null) {
                    payload = RequestReply.RuleSourceExchange.newBuilder()
                            .setLabel(nz(s.getExchange().getLabel())).build().toByteArray();
                }
                break;
            case RuleSourceExternalAddress:
                if (s.getExternalAddress() != null) {
                    payload = RequestReply.RuleSourceExternalAddress.newBuilder()
                            .setAddress(nz(s.getExternalAddress().getAddress()))
                            .setMemo(nz(s.getExternalAddress().getMemo())).build().toByteArray();
                }
                break;
            default:
                // RuleSourceAny, RuleSourceAnyExchange: no payload.
                break;
        }
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteString();
    }

    // ================= Encode helpers =================

    /**
     * Captures a protobuf message's unknown fields as bytes for lossless re-encoding.
     *
     * @param message the protobuf message
     * @return the unknown-field bytes (empty if none)
     */
    default ByteString unknownBytes(com.google.protobuf.Message message) {
        return ByteString.copyFrom(message.getUnknownFields().toByteArray());
    }

    /**
     * Re-attaches preserved unknown fields to a builder. Malformed preserved bytes are ignored.
     *
     * @param builder the builder
     * @param unknown the preserved unknown-field bytes
     * @param <B>     the builder type
     */
    default <B extends com.google.protobuf.Message.Builder> void reattach(B builder, ByteString unknown) {
        if (unknown != null && !unknown.isEmpty()) {
            try {
                builder.setUnknownFields(com.google.protobuf.UnknownFieldSet.parseFrom(unknown.toByteArray()));
            } catch (InvalidProtocolBufferException e) {
                if (LOGGER.isLoggable(Level.WARNING)) {
                    LOGGER.log(Level.WARNING, "Ignoring malformed preserved unknown fields: " + e.getMessage());
                }
            }
        }
    }

    /**
     * Applies a properties map to a builder consumer when non-null.
     *
     * @param props    the properties
     * @param consumer the builder's putAll consumer
     */
    default void putProps(java.util.Map<String, ByteString> props, java.util.function.Consumer<java.util.Map<String, ByteString>> consumer) {
        if (props != null) {
            consumer.accept(props);
        }
    }

    /**
     * Returns "" for a null string (protobuf setters reject null).
     *
     * @param s the string
     * @return a non-null string
     */
    default String nz(String s) {
        return s == null ? "" : s;
    }

    /**
     * Converts role name strings to protobuf Role enums, skipping unrecognized names.
     *
     * @param names the role names
     * @return the protobuf roles
     */
    default List<Integer> rolesFromStrings(List<String> names) {
        List<Integer> out = new ArrayList<>();
        if (names == null) {
            return out;
        }
        for (String n : names) {
            // Dropping a role silently would strip a privilege from data SuperAdmins are
            // about to sign, so a role newer than this SDK passes through by number.
            out.add(enumValue(RequestReply.Role.class, n));
        }
        return out;
    }

    /**
     * Resolves a protobuf enum's numeric value from its name, tolerating names unknown
     * to this SDK version (returns 0, the proto default, rather than crashing).
     *
     * @param enumType the protobuf enum class
     * @param name     the enum constant name
     * @param <E>      the enum type
     * @return the enum's numeric value, or 0 if the name is unknown
     */
    @SuppressWarnings("PMD.EmptyCatchBlock")
    default <E extends Enum<E> & com.google.protobuf.ProtocolMessageEnum> int enumValue(Class<E> enumType, String name) {
        if (name == null || name.isEmpty()) {
            return 0;
        }
        try {
            return Enum.valueOf(enumType, name).getNumber();
        } catch (IllegalArgumentException e) {
            // Either a name this SDK version does not know, or the UNRECOGNIZED sentinel
            // whose getNumber() throws. Both fall through to the numeric form below.
        }
        try {
            return Integer.parseInt(name);
        } catch (NumberFormatException e) {
            throw new IllegalArgumentException(
                    "unknown " + enumType.getSimpleName() + " \"" + name + "\"", e);
        }
    }
}
