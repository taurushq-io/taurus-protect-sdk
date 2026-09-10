package com.taurushq.sdk.protect.client.mapper;

import com.google.protobuf.Descriptors.Descriptor;
import com.google.protobuf.DynamicMessage;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.util.JsonFormat;
import com.taurushq.sdk.protect.proto.v1.RequestReply;

import java.util.Base64;

/**
 * Canonical proto-JSON &lt;-&gt; base64 bridge for governance rules.
 *
 * <p>Mirrors the Go SDK's rules_container_json.go: operates on the generated
 * protobuf messages via protobuf-java-util's {@link JsonFormat}, independent of
 * the {@code DecodedRulesContainer} convenience model, so callers can round-trip
 * a rules container or an individual rule-cell message through canonical
 * protobuf JSON.
 */
public final class RulesContainerJsonMapper {

    private RulesContainerJsonMapper() {
    }

    /**
     * Decodes a base64 protobuf RulesContainer into canonical proto JSON.
     *
     * @param base64Data base64-encoded protobuf RulesContainer
     * @return canonical protobuf JSON
     * @throws InvalidProtocolBufferException if the bytes are not a valid RulesContainer
     */
    public static String rulesContainerJsonFromBase64(String base64Data)
            throws InvalidProtocolBufferException {
        RequestReply.RulesContainer container =
                RequestReply.RulesContainer.parseFrom(Base64.getDecoder().decode(base64Data));
        return JsonFormat.printer().print(container);
    }

    /**
     * Encodes canonical proto JSON into a base64 protobuf RulesContainer.
     *
     * @param json canonical protobuf JSON
     * @return base64-encoded protobuf RulesContainer
     * @throws InvalidProtocolBufferException if the JSON is not a valid RulesContainer
     */
    public static String rulesContainerBase64FromJson(String json)
            throws InvalidProtocolBufferException {
        RequestReply.RulesContainer.Builder builder = RequestReply.RulesContainer.newBuilder();
        JsonFormat.parser().merge(json, builder);
        return Base64.getEncoder().encodeToString(deterministicBytes(builder.build()));
    }

    /**
     * Encodes a governance rule protobuf message (resolved by name) into base64.
     *
     * <p>messageType is a concrete rule message name from request_reply.proto such
     * as RuleSource, RuleFiatAmountRange, RulesContainer_Line, or
     * RulesContainer_TransactionRules.
     *
     * @param messageType the rule message type name
     * @param json        canonical protobuf JSON
     * @return base64-encoded protobuf message
     * @throws InvalidProtocolBufferException if the JSON is not valid for the type
     */
    public static String ruleMessageBase64FromJson(String messageType, String json)
            throws InvalidProtocolBufferException {
        Descriptor descriptor = resolveMessage(messageType);
        DynamicMessage.Builder builder = DynamicMessage.newBuilder(descriptor);
        JsonFormat.parser().merge(json, builder);
        return Base64.getEncoder().encodeToString(deterministicBytes(builder.build()));
    }

    /**
     * Decodes a base64 protobuf rule message (by name) into canonical proto JSON.
     * Returns {@code null} for empty input (an empty cell means "match any").
     *
     * @param messageType the rule message type name
     * @param base64Data  base64-encoded protobuf message, or empty
     * @return canonical protobuf JSON, or null for empty input
     * @throws InvalidProtocolBufferException if the bytes are not valid for the type
     */
    public static String ruleMessageJsonFromBase64(String messageType, String base64Data)
            throws InvalidProtocolBufferException {
        if (base64Data == null || base64Data.trim().isEmpty()) {
            return null;
        }
        byte[] data = Base64.getDecoder().decode(base64Data);
        if (data.length == 0) {
            return null;
        }
        Descriptor descriptor = resolveMessage(messageType);
        DynamicMessage message = DynamicMessage.parseFrom(descriptor, data);
        return JsonFormat.printer().print(message);
    }

    /**
     * Serializes with map entries in a stable order, reusing the single deterministic
     * serializer on {@link RulesContainerMapper}. The container carries
     * {@code map<string, bytes> properties} at five levels, so a plain
     * {@code toByteArray()} here would emit different bytes across runs and differ
     * from the other SDKs — and these are the bytes a proposal is signed over.
     */
    private static byte[] deterministicBytes(com.google.protobuf.Message message) {
        return RulesContainerMapper.INSTANCE.deterministicBytes(message);
    }

    /**
     * Resolves a rule message descriptor by its generated name, navigating nested
     * types via the underscore-flattened convention (e.g. RulesContainer_Line).
     */
    private static Descriptor resolveMessage(String messageType) {
        if (messageType == null || messageType.trim().isEmpty()) {
            throw new IllegalArgumentException("messageType cannot be empty");
        }
        String[] parts = messageType.trim().split("_");
        Descriptor descriptor = RequestReply.getDescriptor().findMessageTypeByName(parts[0]);
        for (int i = 1; i < parts.length && descriptor != null; i++) {
            descriptor = descriptor.findNestedTypeByName(parts[i]);
        }
        if (descriptor == null) {
            throw new IllegalArgumentException("unsupported governance rule message type: " + messageType);
        }
        return descriptor;
    }
}
