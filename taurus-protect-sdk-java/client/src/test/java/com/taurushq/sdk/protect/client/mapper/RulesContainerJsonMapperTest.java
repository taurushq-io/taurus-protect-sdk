package com.taurushq.sdk.protect.client.mapper;

import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.proto.v1.RequestReply;
import org.junit.jupiter.api.Test;

import java.util.Base64;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class RulesContainerJsonMapperTest {

    @Test
    void rulesContainerRoundTrip() throws Exception {
        RequestReply.RulesContainer src = RequestReply.RulesContainer.newBuilder()
                .setMinimumDistinctUserSignatures(3)
                .build();
        String encoded = Base64.getEncoder().encodeToString(src.toByteArray());

        String json = RulesContainerJsonMapper.rulesContainerJsonFromBase64(encoded);
        assertTrue(json.contains("minimumDistinctUserSignatures"));

        String reEncoded = RulesContainerJsonMapper.rulesContainerBase64FromJson(json);
        assertEquals(encoded, reEncoded);
    }

    /**
     * The JSON bridge is the path the MCP governance tools drive: decode a container to
     * JSON, let a caller edit it, re-encode, submit — and approvers sign the re-encoded
     * bytes. An edited container arrives with its map keys in arbitrary order, so encoding
     * MUST be order-independent: plain toByteArray() emits map&lt;string, bytes&gt; in
     * unspecified order. The input below carries the five properties keys in
     * reverse-sorted order on purpose; feeding back the canonical (already sorted) JSON
     * would pass even with a non-deterministic encoder and prove nothing. Expected bytes
     * are the same vector asserted by encodingIsDeterministicAndCrossSdkStable for the
     * typed encoder, and the same reversed input is asserted in all four SDKs.
     */
    @Test
    void jsonBridgeIsOrderIndependentAndCrossSdkStable() throws Exception {
        final String expected = "CmgKAnUxEgNQRU0aAQMiEAoGa0FscGhhEgZrQWxwaGEiEAoGa0JyYXZvEgZr"
                + "QnJhdm8iFAoIa0NoYXJsaWUSCGtDaGFybGllIhAKBmtEZWx0YRIGa0RlbHRhIg4KBWtFY2hv"
                + "EgVrRWNob0oQCgZrQWxwaGESBmtBbHBoYUoQCgZrQnJhdm8SBmtCcmF2b0oUCghrQ2hhcmxp"
                + "ZRIIa0NoYXJsaWVKEAoGa0RlbHRhEgZrRGVsdGFKDgoFa0VjaG8SBWtFY2hv";
        final String reversedKeyJson =
                "{\"users\":[{\"id\":\"u1\",\"publicKey\":\"PEM\",\"roles\":[\"SUPERADMIN\"],"
                + "\"properties\":{\"kEcho\":\"a0VjaG8=\",\"kDelta\":\"a0RlbHRh\","
                + "\"kCharlie\":\"a0NoYXJsaWU=\",\"kBravo\":\"a0JyYXZv\","
                + "\"kAlpha\":\"a0FscGhh\"}}],\"properties\":{\"kEcho\":\"a0VjaG8=\","
                + "\"kDelta\":\"a0RlbHRh\",\"kCharlie\":\"a0NoYXJsaWU=\","
                + "\"kBravo\":\"a0JyYXZv\",\"kAlpha\":\"a0FscGhh\"}}";

        for (int i = 0; i < 20; i++) {
            assertEquals(expected,
                    RulesContainerJsonMapper.rulesContainerBase64FromJson(reversedKeyJson),
                    "run " + i + " not stable/cross-SDK-equal");
        }
    }

    @Test
    void rulesContainerRejectsInvalidJson() {
        assertThrows(InvalidProtocolBufferException.class,
                () -> RulesContainerJsonMapper.rulesContainerBase64FromJson("not json"));
    }

    @Test
    void ruleMessageRoundTripStable() throws Exception {
        String encoded = RulesContainerJsonMapper.ruleMessageBase64FromJson(
                "RuleFiatAmountRange", "{\"minAmount\":\"1000\"}");
        String decoded = RulesContainerJsonMapper.ruleMessageJsonFromBase64("RuleFiatAmountRange", encoded);
        assertNotNull(decoded);
        assertTrue(decoded.contains("1000"));
        assertEquals(encoded,
                RulesContainerJsonMapper.ruleMessageBase64FromJson("RuleFiatAmountRange", decoded));
    }

    @Test
    void ruleMessageEmptyReturnsNull() throws Exception {
        assertNull(RulesContainerJsonMapper.ruleMessageJsonFromBase64("RuleSource", ""));
        assertNull(RulesContainerJsonMapper.ruleMessageJsonFromBase64("RuleSource", "   "));
    }

    @Test
    void ruleMessageUnknownTypeThrows() {
        assertThrows(IllegalArgumentException.class,
                () -> RulesContainerJsonMapper.ruleMessageBase64FromJson("NoSuchRuleMessage", "{}"));
    }

    @Test
    void ruleMessageEmptyTypeThrows() {
        assertThrows(IllegalArgumentException.class,
                () -> RulesContainerJsonMapper.ruleMessageBase64FromJson("   ", "{}"));
    }
}
