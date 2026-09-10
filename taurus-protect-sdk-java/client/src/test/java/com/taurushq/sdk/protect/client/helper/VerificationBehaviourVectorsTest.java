package com.taurushq.sdk.protect.client.helper;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
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
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * Cross-SDK behaviour vectors for the verification primitives.
 *
 * <p>These are the invariants that had drifted apart before: which (blockchain,
 * network) pair selects the governance rules, and whether a hash is covered by a
 * signature. Both are pure input to outcome mappings, so all four SDKs assert them
 * from one file rather than from four hand-maintained copies — the arrangement that
 * let this SDK's hash-coverage go non-constant-time while its peers did not.
 *
 * <p>To add a case: append to the shared file, bump the matching count there, and
 * consume it in all four suites.
 */
class VerificationBehaviourVectorsTest {

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
                assertCounts();
                return;
            }
        }
        throw new IOException("Cannot find verification-behaviour-vectors.json. Working directory: "
                + System.getProperty("user.dir"));
    }

    /**
     * Counts are asserted so a case added to the shared file without being consumed
     * here fails loudly rather than being silently ignored by this SDK.
     */
    private static void assertCounts() {
        JsonObject counts = vectors.getAsJsonObject("counts");
        for (String key : counts.keySet()) {
            assertEquals(counts.get(key).getAsInt(), vectors.getAsJsonArray(key).size(),
                    key + ": vector count disagrees with the count the file declares");
        }
    }

    @Test
    void resolvesTheGovernanceRuleKeyForEveryVector() {
        for (JsonElement element : vectors.getAsJsonArray("rule_key")) {
            JsonObject c = element.getAsJsonObject();
            String description = c.get("description").getAsString();
            String payload = c.get("payload_as_string").getAsString();
            String dtoBlockchain = c.get("dto_blockchain").getAsString();
            String dtoNetwork = c.get("dto_network").getAsString();

            if ("error".equals(c.get("expect").getAsString())) {
                assertThrows(WhitelistException.class,
                        () -> WhitelistHashHelper.resolveRuleKey(payload, dtoBlockchain, dtoNetwork),
                        description);
                continue;
            }

            String[] resolved;
            try {
                resolved = WhitelistHashHelper.resolveRuleKey(payload, dtoBlockchain, dtoNetwork);
            } catch (WhitelistException e) {
                throw new AssertionError(description + ": unexpected failure", e);
            }
            assertEquals(c.get("blockchain").getAsString(), resolved[0], description);
            assertEquals(c.get("network").getAsString(), resolved[1], description);
        }
    }

    @Test
    void agreesOnHashCoverageForEveryVector() {
        for (JsonElement element : vectors.getAsJsonArray("hash_coverage")) {
            JsonObject c = element.getAsJsonObject();
            String description = c.get("description").getAsString();

            List<WhitelistSignature> signatures = new ArrayList<>();
            for (JsonElement sigElement : c.getAsJsonArray("signatures")) {
                WhitelistSignature signature = new WhitelistSignature();
                for (JsonElement hash : sigElement.getAsJsonArray()) {
                    signature.getHashes().add(hash.getAsString());
                }
                signatures.add(signature);
            }

            assertEquals(c.get("expect").getAsBoolean(),
                    SignatureVerifier.verifyHashCoverage(c.get("hash").getAsString(), signatures),
                    description);
        }
    }

    @Test
    void agreesOnSingleSignatureHashContainmentForEveryVector() {
        for (JsonElement element : vectors.getAsJsonArray("contains_hash")) {
            JsonObject c = element.getAsJsonObject();
            String description = c.get("description").getAsString();

            List<String> hashes = new ArrayList<>();
            for (JsonElement hash : c.getAsJsonArray("hashes")) {
                hashes.add(hash.getAsString());
            }

            assertEquals(c.get("expect").getAsBoolean(),
                    SignatureVerifier.containsHash(hashes, c.get("hash").getAsString()),
                    description);
        }
    }
}
