package com.taurushq.sdk.protect.client.helper;

import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.testutil.SignedFixtures;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;

/**
 * Cross-SDK signed fixtures for verification step 2 — the SuperAdmin threshold.
 * Step 5's per-group threshold is gated from the same file by
 * {@code service.SignedFixturesGroupThresholdTest}, which needs a package-private
 * service method.
 *
 * <p>The behaviour vectors cover everything expressible without key material; this
 * file covers what needs real signatures. Until it existed, the repo's own
 * CLAUDE.md recorded that NO automated cross-SDK gate covered the distinct-key
 * rule, and the five cases below lived as five hand-maintained copies in four
 * suites — this SDK's had been packed into a single test method, so one break hid
 * four.
 *
 * <p>The fixture carries PUBLIC keys and signatures only — the gate verifies, it
 * never signs — so there is no private key material in the repo.
 */
class SignedFixturesTest {

    private static JsonObject fixtures;
    private static byte[] container;

    @BeforeAll
    static void loadFixtures() throws IOException {
        fixtures = SignedFixtures.load();
        String b64 = fixtures.get("rules_container_base64").getAsString();
        assertTrue(!b64.isEmpty(), "fixture carries no rules container");
        container = Base64.getDecoder().decode(b64);
    }

    @Test
    void agreesWithEveryRecordedOutcome() {
        for (JsonElement element : fixtures.getAsJsonArray("superadmin_threshold")) {
            JsonObject c = element.getAsJsonObject();
            String description = c.get("description").getAsString();

            List<PublicKey> keys = new ArrayList<>();
            for (JsonElement pem : c.getAsJsonArray("super_admin_keys_pem")) {
                try {
                    keys.add(CryptoTPV1.decodePublicKey(pem.getAsString()));
                } catch (Exception e) {
                    throw new AssertionError(description + ": cannot decode fixture key", e);
                }
            }

            List<RuleUserSignature> signatures = new ArrayList<>();
            for (JsonElement sigElement : c.getAsJsonArray("signatures")) {
                JsonObject s = sigElement.getAsJsonObject();
                RuleUserSignature signature = new RuleUserSignature();
                signature.setUserId(s.get("user_id").getAsString());
                signature.setSignature(s.get("signature").getAsString());
                signatures.add(signature);
            }

            int minValid = c.get("min_valid_signatures").getAsInt();
            boolean expectError = "error".equals(c.get("expect").getAsString());

            try {
                SignatureVerifier.verifyGovernanceRulesSignatures(
                        container, signatures, keys, minValid);
                if (expectError) {
                    fail(description + ": expected verification to fail, it passed");
                }
            } catch (IntegrityException | IllegalArgumentException e) {
                if (!expectError) {
                    throw new AssertionError(
                            description + ": expected verification to pass", e);
                }
            }
        }
    }

    /**
     * Without a case where the entry COUNT meets the threshold but the distinct-key
     * count does not, the whole file would pass against an implementation that
     * counts entries — the regression these fixtures exist to catch.
     */
    @Test
    void containsACaseThatDistinguishesDistinctKeysFromEntries() {
        boolean discriminating = false;
        for (JsonElement element : fixtures.getAsJsonArray("superadmin_threshold")) {
            JsonObject c = element.getAsJsonObject();
            int minValid = c.get("min_valid_signatures").getAsInt();
            int entries = c.getAsJsonArray("signatures").size();
            if ("error".equals(c.get("expect").getAsString()) && minValid > 0 && entries >= minValid) {
                discriminating = true;
                break;
            }
        }
        assertTrue(discriminating,
                "no vector distinguishes distinct-key counting from entry counting");
    }
}
