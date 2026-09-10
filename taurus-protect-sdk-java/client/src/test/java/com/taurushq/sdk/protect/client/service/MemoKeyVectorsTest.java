package com.taurushq.sdk.protect.client.service;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Cross-SDK vectors for the governance verification memo key.
 *
 * <p>The key decides whether ECDSA runs at all, so it MUST be injective. An unprefixed
 * concatenation leaves the boundary between the container and the signature list
 * uncommitted, and a response-controlling attacker can then shift bytes across it to make
 * a MODIFIED container inherit a genuine one's "already verified" status — skipping
 * signature verification on the document that carries every HSM key the SDK trusts.
 *
 * <p>This is a second reader of the shared vector file, in the {@code service} package
 * rather than beside the other behaviour vectors in {@code helper}, because
 * {@code rulesetVerificationKey} is package-private here — the same reason the signed
 * fixtures need a second consumer for their group-threshold section. The generic
 * count-vs-declared guard stays in {@code VerificationBehaviourVectorsTest}, which loops
 * over every section, so there is no risk of the two drifting on it; this file asserts
 * only its own section's count.
 *
 * <p>To add a case: append to the shared file, bump {@code counts.memo_key} there, and
 * consume it in all four suites.
 */
class MemoKeyVectorsTest {

    private static JsonObject vectors;

    @BeforeAll
    static void loadVectors() throws IOException {
        String[] candidates = {
                "../../scripts/resources/verification-behaviour-vectors.json", // from client/
                "../scripts/resources/verification-behaviour-vectors.json",    // from sdk-java/
                "scripts/resources/verification-behaviour-vectors.json",       // from repo root
        };
        for (String candidate : candidates) {
            Path p = Paths.get(candidate).toAbsolutePath().normalize();
            if (Files.exists(p)) {
                String json = new String(Files.readAllBytes(p), StandardCharsets.UTF_8);
                vectors = JsonParser.parseString(json).getAsJsonObject();
                return;
            }
        }
        throw new IOException("Cannot find verification-behaviour-vectors.json. Working directory: "
                + System.getProperty("user.dir"));
    }

    /** Builds a GovernanceRules from a memo_key vector side. */
    private static GovernanceRules ruleset(final JsonObject spec) {
        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(spec.get("rules_container").getAsString());
        List<RuleUserSignature> signatures = new ArrayList<>();
        for (JsonElement element : spec.getAsJsonArray("signatures")) {
            JsonObject sig = element.getAsJsonObject();
            RuleUserSignature entry = new RuleUserSignature();
            entry.setUserId(sig.get("user_id").getAsString());
            entry.setSignature(sig.get("signature").getAsString());
            signatures.add(entry);
        }
        rules.setRulesSignatures(signatures);
        return rules;
    }

    @Test
    void keepsTheMemoKeyInjectiveForEveryVector() {
        JsonArray cases = vectors.getAsJsonArray("memo_key");
        assertEquals(vectors.getAsJsonObject("counts").get("memo_key").getAsInt(), cases.size(),
                "memo_key: vector count disagrees with the count the file declares");
        assertTrue(cases.size() > 0, "memo_key: no vectors loaded");

        // Non-vacuity: without a "distinct" pair the section would pass against a key
        // function that returns a constant.
        boolean anyDistinct = false;
        for (JsonElement element : cases) {
            if ("distinct".equals(element.getAsJsonObject().get("expect").getAsString())) {
                anyDistinct = true;
                break;
            }
        }
        assertTrue(anyDistinct,
                "memo_key: no 'distinct' vector, so the section cannot catch a colliding key");

        for (JsonElement element : cases) {
            JsonObject c = element.getAsJsonObject();
            String description = c.get("description").getAsString();

            String keyA = GovernanceRuleService.rulesetVerificationKey(
                    ruleset(c.getAsJsonObject("a")));
            String keyB = GovernanceRuleService.rulesetVerificationKey(
                    ruleset(c.getAsJsonObject("b")));
            assertNotNull(keyA, description + ": container a must be decodable");
            assertNotNull(keyB, description + ": container b must be decodable");

            if ("distinct".equals(c.get("expect").getAsString())) {
                assertNotEquals(keyA, keyB,
                        description + ": distinct (container, signatures) pairs must not "
                                + "share a memo key");
            } else {
                assertEquals(keyA, keyB,
                        description + ": these inputs describe the same verified document, "
                                + "so the key must match");
            }
        }
    }
}
