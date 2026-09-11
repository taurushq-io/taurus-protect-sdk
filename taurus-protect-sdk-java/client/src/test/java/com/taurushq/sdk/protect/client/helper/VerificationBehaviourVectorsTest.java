package com.taurushq.sdk.protect.client.helper;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.ContractAddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
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

    // ------------------------------------------------------------------------------
    // legacy_hash and rule_tier_candidates
    // ------------------------------------------------------------------------------
    //
    // Until 2026-09-10 both sections were consumed by the GO suite alone, even though the
    // file's `counts` block declares them and the loop above asserts every count.
    // Asserting a section's LENGTH proves the file is well formed; it does not prove the
    // behaviour is checked. Python shipped the single-tier lookup `rule_tier_candidates`
    // exists to forbid while three SDKs had the fix, and no gate went red.

    @Test
    @DisplayName("step 6 parses the payload the signature COVERED, for every vector")
    void parsesThePayloadTheSignatureCovered() {
        // Asserts the PARSED PAYLOAD, not a hash: the legacy-strip injection moves no hash
        // at all, so crypto-test-vectors.json is structurally blind to it.
        for (JsonElement e : vectors.getAsJsonArray("legacy_hash")) {
            JsonObject c = e.getAsJsonObject();
            String description = c.get("description").getAsString();
            String coveredHash =
                    CryptoTPV1.calculateHexHash(c.get("signed_payload").getAsString());
            String delivered = c.get("delivered_payload").getAsString();

            String matchedPayload = null;
            if (coveredHash.equals(CryptoTPV1.calculateHexHash(delivered))) {
                matchedPayload = delivered;
            } else {
                for (LegacyPayloadVariant variant
                        : WhitelistHashHelper.computeLegacyPayloadVariants(delivered)) {
                    if (coveredHash.equals(variant.getHash())) {
                        matchedPayload = variant.getPayload();
                        break;
                    }
                }
            }

            if ("no_match".equals(c.get("expect").getAsString())) {
                assertNull(matchedPayload,
                        description + ": expected no variant to be covered, but one matched");
                continue;
            }
            assertNotNull(matchedPayload, description + ": expected a covered variant");
            assertEquals(c.get("expect_matched_payload").getAsString(), matchedPayload,
                    description);
        }
    }

    @Test
    @DisplayName("every reachable rule tier is returned, for every vector")
    void returnsEveryReachableRuleTier() {
        // When the signed payload omits `network` there is no authenticated way to learn
        // whether that was legitimate, so every tier the unsigned DTO value could have
        // selected must be enforced. Too few lets the server pick the quorum; too many
        // rejects rows governance would accept.
        for (JsonElement e : vectors.getAsJsonArray("rule_tier_candidates")) {
            JsonObject c = e.getAsJsonObject();
            String description = c.get("description").getAsString();
            String blockchain = c.get("blockchain").getAsString();

            List<AddressWhitelistingRules> addressRules = new ArrayList<>();
            List<ContractAddressWhitelistingRules> contractRules = new ArrayList<>();
            for (JsonElement r : c.getAsJsonArray("rules")) {
                JsonObject rule = r.getAsJsonObject();
                AddressWhitelistingRules a = new AddressWhitelistingRules();
                a.setCurrency(rule.get("blockchain").getAsString());
                a.setNetwork(rule.get("network").getAsString());
                addressRules.add(a);

                ContractAddressWhitelistingRules k = new ContractAddressWhitelistingRules();
                k.setBlockchain(rule.get("blockchain").getAsString());
                k.setNetwork(rule.get("network").getAsString());
                contractRules.add(k);
            }
            DecodedRulesContainer container = new DecodedRulesContainer();
            container.setAddressWhitelistingRules(addressRules);
            container.setContractAddressWhitelistingRules(contractRules);

            List<String> expected = new ArrayList<>();
            for (JsonElement x : c.getAsJsonArray("expect_candidates")) {
                JsonObject want = x.getAsJsonObject();
                expected.add(want.get("blockchain").getAsString() + "|"
                        + want.get("network").getAsString());
            }

            List<String> gotAddress = new ArrayList<>();
            for (AddressWhitelistingRules r
                    : container.findAddressWhitelistingRuleCandidates(blockchain)) {
                gotAddress.add(nullToEmpty(r.getCurrency()) + "|" + nullToEmpty(r.getNetwork()));
            }
            assertEquals(expected, gotAddress, description);

            // The asset peer, from the same vectors: the two families name the chain field
            // differently (currency vs blockchain), so a walk reading one name treats every
            // contract rule as a wildcard global default — fail-OPEN, the broadest tier.
            List<String> gotContract = new ArrayList<>();
            for (ContractAddressWhitelistingRules r
                    : container.findContractAddressWhitelistingRuleCandidates(blockchain)) {
                gotContract.add(nullToEmpty(r.getBlockchain()) + "|" + nullToEmpty(r.getNetwork()));
            }
            assertEquals(expected, gotContract, description + " (contract family)");
        }
    }

    private static String nullToEmpty(final String value) {
        return value == null ? "" : value;
    }
}
