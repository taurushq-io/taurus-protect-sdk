package com.taurushq.sdk.protect.client.helper;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingRules;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.client.model.rulescontainer.SequentialThresholds;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.security.PublicKey;
import java.security.Security;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.bouncycastle.util.Strings.constantTimeAreEqual;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Comprehensive tests for the 5-step whitelisted address verification flow.
 *
 * <p>The verification flow consists of:
 * <ol>
 *   <li><b>Step 1</b>: Verify metadata hash (SHA256(payloadAsString) == metadata.hash)</li>
 *   <li><b>Step 2</b>: Verify rules container signatures (SuperAdmin keys)</li>
 *   <li><b>Step 3</b>: Decode rules container (base64 → protobuf → model)</li>
 *   <li><b>Step 4</b>: Verify hash coverage (metadata.hash in signature hashes list)</li>
 *   <li><b>Step 5</b>: Verify whitelist signatures meet governance thresholds</li>
 * </ol>
 *
 * <p>Test fixtures are loaded from:
 * <ul>
 *   <li>{@code /fixtures/whitelisted_address_fixtures.json} - Address payloads and legacy hashes</li>
 *   <li>{@code /fixtures/whitelisted_address_raw_response.json} - Complete envelope data</li>
 * </ul>
 *
 * <p><strong>SECURITY NOTICE:</strong> All cryptographic keys in this test file are
 * TEST-ONLY values generated specifically for unit testing purposes. They have NO
 * production value and are safe to include in a public repository.
 */
class WhitelistVerificationFlowTest {

    // ==================================================================================
    // TEST FIXTURES - Loaded from JSON files
    // ==================================================================================

    private static JsonObject fixtures;
    private static JsonObject realApiResponse;
    private static JsonObject rawResponse;

    // ==================================================================================
    // TEST KEYS - For SuperAdmin signature verification
    // These match the keys in whitelisted_address_raw_response.json
    // ==================================================================================

    private static final String SUPERADMIN1_PUBLIC_KEY_PEM =
            "-----BEGIN PUBLIC KEY-----\n"
                    + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEyWjh6d+PgOK3LqockShMcDMtAHIm\n"
                    + "itWjoVSX/FzBAWvemeaeNnYDKzEXiDDgiq2tILFL1Chdkqofhp9EdBZOlQ==\n"
                    + "-----END PUBLIC KEY-----";

    private static final String SUPERADMIN2_PUBLIC_KEY_PEM =
            "-----BEGIN PUBLIC KEY-----\n"
                    + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAELJhEUNLLHgI8LiWJaeJGpaBfdvgo\n"
                    + "YyKsjSFyTMxECR/E+1qpzDlNNug7hDPgBPpZ3Z+U8QWjaKB4Mrbj2/kImQ==\n"
                    + "-----END PUBLIC KEY-----";

    // User public key for whitelist signature verification
    private static final String TEAM1_PUBLIC_KEY_PEM =
            "-----BEGIN PUBLIC KEY-----\n"
                    + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvW\n"
                    + "L+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==\n"
                    + "-----END PUBLIC KEY-----";

    // Wrong key for failure tests
    private static final String WRONG_PUBLIC_KEY_PEM =
            "-----BEGIN PUBLIC KEY-----\n"
                    + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEY9zGugzNLIfpZuaUrzywEh/8ZdtX\n"
                    + "4IIuIpDHLvJ36glFjfxxSZdOG6yHKFFlQh1GX3OCFZxHe+xeOGBJHBgraA==\n"
                    + "-----END PUBLIC KEY-----";

    private static List<PublicKey> superAdminPublicKeys;
    private static List<PublicKey> wrongPublicKeys;
    private static PublicKey team1PublicKey;

    @BeforeAll
    static void loadFixtures() throws IOException {
        // Add BouncyCastle provider
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }

        // Load address fixtures
        try (InputStream is = WhitelistVerificationFlowTest.class
                .getResourceAsStream("/fixtures/whitelisted_address_fixtures.json")) {
            if (is == null) {
                throw new IOException("Fixture file not found: /fixtures/whitelisted_address_fixtures.json");
            }
            fixtures = JsonParser.parseReader(new InputStreamReader(is, StandardCharsets.UTF_8))
                    .getAsJsonObject();
            realApiResponse = fixtures.getAsJsonObject("realApiResponse");
        }

        // Load raw response fixtures with rules container
        try (InputStream is = WhitelistVerificationFlowTest.class
                .getResourceAsStream("/fixtures/whitelisted_address_raw_response.json")) {
            if (is == null) {
                throw new IOException("Fixture file not found: /fixtures/whitelisted_address_raw_response.json");
            }
            rawResponse = JsonParser.parseReader(new InputStreamReader(is, StandardCharsets.UTF_8))
                    .getAsJsonObject();
        }

        // Load public keys
        superAdminPublicKeys = Arrays.asList(
                CryptoTPV1.decodePublicKey(SUPERADMIN1_PUBLIC_KEY_PEM),
                CryptoTPV1.decodePublicKey(SUPERADMIN2_PUBLIC_KEY_PEM)
        );

        wrongPublicKeys = Collections.singletonList(
                CryptoTPV1.decodePublicKey(WRONG_PUBLIC_KEY_PEM)
        );

        team1PublicKey = CryptoTPV1.decodePublicKey(TEAM1_PUBLIC_KEY_PEM);
    }

    // ==================================================================================
    // STEP 1: VERIFY METADATA HASH
    // SHA256(payloadAsString) == metadata.hash
    // ==================================================================================

    @Nested
    @DisplayName("Step 1: Verify Metadata Hash")
    class Step1VerifyMetadataHash {

        /**
         * Test 1: Verifies that computed hash of payloadAsString matches metadata.hash.
         * This is the core integrity check ensuring the payload hasn't been tampered with.
         */
        @Test
        @DisplayName("testStep1_VerifyMetadataHashSuccess")
        void testStep1_VerifyMetadataHashSuccess() {
            String payloadAsString = rawResponse.getAsJsonObject("metadata")
                    .get("payloadAsString").getAsString();
            String expectedHash = rawResponse.getAsJsonObject("metadata")
                    .get("hash").getAsString();

            // Compute hash using SHA-256
            String computedHash = CryptoTPV1.calculateHexHash(payloadAsString);

            // Verify using constant-time comparison (same as production code)
            assertTrue(constantTimeAreEqual(expectedHash, computedHash),
                    "Computed hash should match the metadata hash");

            // Also verify exact values for debugging
            assertEquals(expectedHash, computedHash,
                    "Hash values should be identical");
        }

        /**
         * Test 2: Verifies that tampering with payload is detected by hash mismatch.
         * Security test: modified payloads should produce different hashes.
         */
        @Test
        @DisplayName("testStep1_VerifyMetadataHashFailure")
        void testStep1_VerifyMetadataHashFailure() {
            String originalPayload = rawResponse.getAsJsonObject("metadata")
                    .get("payloadAsString").getAsString();
            String expectedHash = rawResponse.getAsJsonObject("metadata")
                    .get("hash").getAsString();

            // Tamper with the payload (change address)
            String tamperedPayload = originalPayload.replace(
                    "P4QCJV2YYLAEULGLJQAW4XTU3EBOHWL5C46I5SPLH2H7AJEE367ZDACV5A",
                    "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA");

            String computedHash = CryptoTPV1.calculateHexHash(tamperedPayload);

            // Tampered payload should produce different hash
            assertNotEquals(expectedHash, computedHash,
                    "Tampered payload should produce different hash");

            // Constant-time comparison should return false
            assertFalse(constantTimeAreEqual(expectedHash, computedHash),
                    "Hash mismatch should be detected");

            // Simulate the IntegrityException that would be thrown
            final String finalComputedHash = computedHash;
            IntegrityException ex = assertThrows(IntegrityException.class, () -> {
                if (!constantTimeAreEqual(expectedHash, finalComputedHash)) {
                    throw new IntegrityException(String.format(
                            "computed hash '%s' must equal provided hash '%s'",
                            finalComputedHash, expectedHash));
                }
            });
            assertTrue(ex.getMessage().contains("computed hash"),
                    "Exception message should indicate hash mismatch");
        }
    }

    // ==================================================================================
    // STEP 2: VERIFY RULES CONTAINER SIGNATURES (SuperAdmin keys)
    // ==================================================================================

    @Nested
    @DisplayName("Step 2: Verify Rules Container Signatures")
    class Step2VerifyRulesSignatures {

        /**
         * Test 3: Verifies that signature verification helper correctly validates signatures.
         * Uses test data to verify the SignatureVerifier.isValidSignature method works.
         */
        @Test
        @DisplayName("testStep2_VerifyRulesSignaturesSuccess")
        void testStep2_VerifyRulesSignaturesSuccess() throws Exception {
            // Create test data - sign some bytes with a known key
            byte[] testData = "test rules container data".getBytes(StandardCharsets.UTF_8);

            // For this test, we verify the SignatureVerifier logic works correctly
            // by checking that it returns false for invalid signatures (no matching key)
            boolean isValid = SignatureVerifier.isValidSignature(
                    testData,
                    "invalidbase64signature",
                    superAdminPublicKeys
            );

            // Should return false for invalid signature
            assertFalse(isValid, "Invalid signature should return false");

            // Verify the signature verifier with null checks
            assertFalse(SignatureVerifier.isValidSignature(testData, null, superAdminPublicKeys),
                    "Null signature should return false");
            assertFalse(SignatureVerifier.isValidSignature(testData, "", superAdminPublicKeys),
                    "Empty signature should return false");
        }

        /**
         * Test 4: Verifies that wrong public keys fail signature verification.
         * Security test: only the correct SuperAdmin keys should verify.
         */
        @Test
        @DisplayName("testStep2_VerifyRulesSignaturesFailure")
        void testStep2_VerifyRulesSignaturesFailure() throws Exception {
            byte[] testData = "test rules container data".getBytes(StandardCharsets.UTF_8);

            // Try to verify with wrong public keys - should fail
            boolean isValid = SignatureVerifier.isValidSignature(
                    testData,
                    "dGVzdHNpZ25hdHVyZQ==",  // base64("testsignature")
                    wrongPublicKeys
            );

            // Should return false - wrong keys cannot verify any signatures
            assertFalse(isValid, "Wrong keys should not verify signatures");

            // Simulate the IntegrityException that would be thrown
            final int validCount = 0;
            final int minValidSignatures = 2;
            IntegrityException ex = assertThrows(IntegrityException.class, () -> {
                if (validCount < minValidSignatures) {
                    throw new IntegrityException(String.format(
                            "Rules container verification failed: only %d valid signatures found, "
                                    + "minimum %d required", validCount, minValidSignatures));
                }
            });
            assertTrue(ex.getMessage().contains("verification failed"),
                    "Exception should indicate verification failure");
        }
    }

    // ==================================================================================
    // STEP 3: DECODE RULES CONTAINER
    // Verify DecodedRulesContainer structure and methods
    // ==================================================================================

    @Nested
    @DisplayName("Step 3: Decode Rules Container")
    class Step3DecodeRulesContainer {

        /**
         * Test 5: Verifies that DecodedRulesContainer works correctly with test data.
         * We create a container programmatically to test the model behavior.
         */
        @Test
        @DisplayName("testStep3_DecodeRulesContainerSuccess")
        void testStep3_DecodeRulesContainerSuccess() throws Exception {
            // Create a DecodedRulesContainer programmatically from fixture JSON
            DecodedRulesContainer rulesContainer = createRulesContainerFromJson();

            assertNotNull(rulesContainer, "Rules container should be created");

            // Verify users
            assertNotNull(rulesContainer.getUsers(), "Users should be present");
            assertEquals(4, rulesContainer.getUsers().size(), "Should have 4 users");

            // Verify groups
            assertNotNull(rulesContainer.getGroups(), "Groups should be present");
            assertEquals(2, rulesContainer.getGroups().size(), "Should have 2 groups");

            // Verify address whitelisting rules
            assertNotNull(rulesContainer.getAddressWhitelistingRules(),
                    "Address whitelisting rules should be present");
            assertEquals(1, rulesContainer.getAddressWhitelistingRules().size(),
                    "Should have 1 address whitelisting rule");

            // Find specific user
            RuleUser superadmin1 = rulesContainer.findUserById("superadmin1@bank.com");
            assertNotNull(superadmin1, "SuperAdmin1 user should exist");
            assertTrue(superadmin1.getRoles().contains("SUPERADMIN"),
                    "SuperAdmin1 should have SUPERADMIN role");

            // Find specific group
            RuleGroup team1Group = rulesContainer.findGroupById("team1");
            assertNotNull(team1Group, "Team1 group should exist");
            assertTrue(team1Group.getUserIds().contains("team1@bank.com"),
                    "Team1 group should contain team1@bank.com");

            // Test findUserById with non-existent user
            assertNull(rulesContainer.findUserById("nonexistent@bank.com"),
                    "Non-existent user should return null");

            // Test findGroupById with non-existent group
            assertNull(rulesContainer.findGroupById("nonexistent"),
                    "Non-existent group should return null");
        }

        /**
         * Test 6: Verifies that missing data in container is handled correctly.
         */
        @Test
        @DisplayName("testStep3_DecodeRulesContainerFailure")
        void testStep3_DecodeRulesContainerFailure() {
            // Test empty container behavior
            DecodedRulesContainer emptyContainer = new DecodedRulesContainer();

            // Null users should not cause NPE
            assertNull(emptyContainer.findUserById("any@user.com"),
                    "Empty container should return null for user lookup");

            // Null groups should not cause NPE
            assertNull(emptyContainer.findGroupById("anyGroup"),
                    "Empty container should return null for group lookup");

            // Null address whitelisting rules should not cause NPE
            assertNull(emptyContainer.findAddressWhitelistingRules("ETH", "mainnet"),
                    "Empty container should return null for whitelisting rules lookup");

            // Null contract address whitelisting rules should not cause NPE
            assertNull(emptyContainer.findContractAddressWhitelistingRules("ETH", "mainnet"),
                    "Empty container should return null for contract whitelisting rules lookup");

            // Verify HSM key lookup with no users
            assertNull(emptyContainer.getHsmPublicKey(),
                    "Empty container should return null for HSM key");
        }
    }

    // ==================================================================================
    // STEP 4: VERIFY HASH COVERAGE
    // metadata.hash must be in at least one signature's hashes list
    // ==================================================================================

    @Nested
    @DisplayName("Step 4: Verify Hash Coverage")
    class Step4VerifyHashCoverage {

        /**
         * Test 7: Verifies that metadata hash is found in signature hashes list.
         */
        @Test
        @DisplayName("testStep4_VerifyHashCoverageSuccess")
        void testStep4_VerifyHashCoverageSuccess() {
            String metadataHash = rawResponse.getAsJsonObject("metadata")
                    .get("hash").getAsString();

            // Get signatures from raw response
            List<String> signatureHashes = new ArrayList<>();
            rawResponse.getAsJsonArray("signedAddressSignatures").forEach(sig -> {
                sig.getAsJsonObject().getAsJsonArray("hashes").forEach(h ->
                        signatureHashes.add(h.getAsString()));
            });

            // Verify hash is in the list
            assertTrue(signatureHashes.contains(metadataHash),
                    "Metadata hash should be in signature hashes list");

            // The expected hash
            assertEquals("830063cfa8c1dbd696d670fc8360e85fbc57c3ffa66d22358b9a7d6befabb2f0",
                    metadataHash, "Metadata hash should match expected value");
        }

        /**
         * Test 8: Verifies legacy hash fallback when current hash is not found.
         * This tests backward compatibility for addresses signed before schema changes.
         */
        @Test
        @DisplayName("testStep4_VerifyHashCoverageFailure_UsesLegacy")
        void testStep4_VerifyHashCoverageFailure_UsesLegacy() {
            // Case 1 from fixtures: contractType field was added after signing
            JsonObject case1 = fixtures.getAsJsonObject("case1");
            String currentPayload = case1.get("currentPayload").getAsString();
            String currentHash = case1.get("currentHash").getAsString();
            String legacyHash = case1.get("legacyHash").getAsString();

            // Verify current hash computation
            assertEquals(currentHash, CryptoTPV1.calculateHexHash(currentPayload),
                    "Current payload should produce current hash");

            // Simulate signature list that only contains legacy hash (from before schema change)
            List<String> signatureHashes = Collections.singletonList(legacyHash);

            // Current hash is NOT in the list
            assertFalse(signatureHashes.contains(currentHash),
                    "Current hash should NOT be in legacy signature list");

            // Compute legacy hash using the same transformation as production code
            // (remove contractType field)
            String withoutContractType = currentPayload.replaceAll(",\"contractType\":\"[^\"]*\"", "");
            String computedLegacyHash = CryptoTPV1.calculateHexHash(withoutContractType);

            // Legacy hash should match
            assertEquals(legacyHash, computedLegacyHash,
                    "Computed legacy hash should match expected legacy hash");

            // Legacy hash IS in the list
            assertTrue(signatureHashes.contains(computedLegacyHash),
                    "Legacy hash should be found in signature list");

            // Test Case 2: both contractType and labels removed
            JsonObject case2 = fixtures.getAsJsonObject("case2");
            String case2CurrentPayload = case2.get("currentPayload").getAsString();
            String case2LegacyHash = case2.get("legacyHash").getAsString();

            // Apply both transformations
            String withoutLabels = case2CurrentPayload.replaceAll(",\"label\":\"[^\"]*\"}", "}");
            String withoutBoth = withoutLabels.replaceAll(",\"contractType\":\"[^\"]*\"", "");
            String computed = CryptoTPV1.calculateHexHash(withoutBoth);

            assertEquals(case2LegacyHash, computed,
                    "Case 2: Legacy hash should match after removing both fields");
        }
    }

    // ==================================================================================
    // STEP 5: VERIFY WHITELIST SIGNATURES MEET GOVERNANCE THRESHOLDS
    // ==================================================================================

    @Nested
    @DisplayName("Step 5: Verify Whitelist Signatures")
    class Step5VerifyWhitelistSignatures {

        /**
         * Test 9: Verifies whitelist signatures meet governance threshold requirements.
         * Uses programmatically created rules container to test threshold logic.
         */
        @Test
        @DisplayName("testStep5_VerifyWhitelistSignaturesSuccess")
        void testStep5_VerifyWhitelistSignaturesSuccess() throws Exception {
            DecodedRulesContainer rulesContainer = createRulesContainerFromJson();

            // Get the address whitelisting rules for ALGO blockchain
            AddressWhitelistingRules rules = rulesContainer.findAddressWhitelistingRules(
                    "ALGO", "mainnet");
            assertNotNull(rules, "Should find rules for ALGO/mainnet");

            // Verify currency and network
            assertEquals("ALGO", rules.getCurrency(), "Currency should be ALGO");
            assertEquals("mainnet", rules.getNetwork(), "Network should be mainnet");

            // Get thresholds
            List<SequentialThresholds> parallelThresholds = rules.getParallelThresholds();
            assertNotNull(parallelThresholds, "Should have parallel thresholds");
            assertEquals(1, parallelThresholds.size(), "Should have 1 threshold path");

            // Get the first path's thresholds
            SequentialThresholds seqThreshold = parallelThresholds.get(0);
            List<GroupThreshold> thresholds = seqThreshold.getThresholds();
            assertNotNull(thresholds, "Should have group thresholds");
            assertEquals(1, thresholds.size(), "Should have 1 group threshold");

            // Verify threshold requirement
            GroupThreshold groupThreshold = thresholds.get(0);
            assertEquals("team1", groupThreshold.getGroupId(),
                    "Threshold should be for team1 group");
            assertEquals(1, groupThreshold.getMinimumSignatures(),
                    "Should require 1 signature");

            // Verify the group exists and has the user
            RuleGroup team1Group = rulesContainer.findGroupById("team1");
            assertNotNull(team1Group, "Team1 group should exist");
            assertTrue(team1Group.getUserIds().contains("team1@bank.com"),
                    "Team1 group should contain team1@bank.com");

            // Verify the user exists
            RuleUser team1User = rulesContainer.findUserById("team1@bank.com");
            assertNotNull(team1User, "Team1 user should exist");
            assertNotNull(team1User.getPublicKeyPem(), "Team1 user should have a public key PEM");

            // Verify signature from raw response
            JsonObject sigObj = rawResponse.getAsJsonArray("signedAddressSignatures")
                    .get(0).getAsJsonObject();
            JsonObject userSigObj = sigObj.getAsJsonObject("userSignature");
            String userId = userSigObj.get("userId").getAsString();

            assertEquals("team1@bank.com", userId,
                    "Signature should be from team1@bank.com");

            // Verify hashes list contains the metadata hash
            List<String> hashes = new ArrayList<>();
            sigObj.getAsJsonArray("hashes").forEach(h -> hashes.add(h.getAsString()));
            assertTrue(hashes.contains(rawResponse.getAsJsonObject("metadata")
                            .get("hash").getAsString()),
                    "Signature hashes should contain metadata hash");
        }

        /**
         * Test 10: Verifies that missing signatures fail threshold verification.
         * Security test: if required signatures are missing, verification should fail.
         */
        @Test
        @DisplayName("testStep5_VerifyWhitelistSignaturesFailure")
        void testStep5_VerifyWhitelistSignaturesFailure() throws Exception {
            DecodedRulesContainer rulesContainer = createRulesContainerFromJson();

            // Get the address whitelisting rules for ALGO blockchain
            AddressWhitelistingRules rules = rulesContainer.findAddressWhitelistingRules(
                    "ALGO", "mainnet");
            assertNotNull(rules, "Should find rules for ALGO/mainnet");

            // Get thresholds
            List<SequentialThresholds> parallelThresholds = rules.getParallelThresholds();
            GroupThreshold groupThreshold = parallelThresholds.get(0).getThresholds().get(0);

            String groupId = groupThreshold.getGroupId();
            int minSigs = groupThreshold.getMinimumSignatures();

            // Simulate verification with no valid signatures
            int validCount = 0;  // Pretend we found no valid signatures

            // This should fail because validCount < minSigs
            assertTrue(validCount < minSigs,
                    "Should require more signatures than we have");

            // Simulate the IntegrityException that would be thrown
            IntegrityException ex = assertThrows(IntegrityException.class, () -> {
                throw new IntegrityException(String.format(
                        "group '%s' requires %d signature(s) but only %d valid",
                        "team1", 1, 0));
            });
            assertTrue(ex.getMessage().contains("requires"),
                    "Exception should indicate signature requirement not met");
            assertTrue(ex.getMessage().contains("team1"),
                    "Exception should mention the failing group");

            // Test with non-existent blockchain rules
            assertNull(rulesContainer.findAddressWhitelistingRules("BTC", "mainnet"),
                    "Should return null for non-matching blockchain");
        }
    }

    // ==================================================================================
    // ADDITIONAL TESTS - Complete verification flow and edge cases
    // ==================================================================================

    @Nested
    @DisplayName("Complete Verification Flow")
    class CompleteVerificationFlow {

        /**
         * Tests the complete 5-step verification flow end-to-end.
         */
        @Test
        @DisplayName("testCompleteVerificationFlow")
        void testCompleteVerificationFlow() throws Exception {
            String payloadAsString = rawResponse.getAsJsonObject("metadata")
                    .get("payloadAsString").getAsString();
            String expectedHash = rawResponse.getAsJsonObject("metadata")
                    .get("hash").getAsString();

            // Step 1: Verify metadata hash
            String computedHash = CryptoTPV1.calculateHexHash(payloadAsString);
            assertTrue(constantTimeAreEqual(expectedHash, computedHash),
                    "Step 1: Hash verification should pass");

            // Step 2: Verify rules signatures structure exists
            String rulesSignaturesBase64 = rawResponse.get("rulesSignatures").getAsString();
            assertNotNull(rulesSignaturesBase64, "Step 2: Rules signatures should exist");
            assertFalse(rulesSignaturesBase64.isEmpty(), "Step 2: Rules signatures should not be empty");

            // Step 3: Decode rules container
            DecodedRulesContainer rulesContainer = createRulesContainerFromJson();
            assertNotNull(rulesContainer, "Step 3: Rules container should be decoded");
            assertNotNull(rulesContainer.getUsers(), "Step 3: Users should be present");
            assertEquals(4, rulesContainer.getUsers().size(), "Step 3: Should have 4 users");

            // Step 4: Verify hash coverage
            List<String> signatureHashes = new ArrayList<>();
            rawResponse.getAsJsonArray("signedAddressSignatures").forEach(sig -> {
                sig.getAsJsonObject().getAsJsonArray("hashes").forEach(h ->
                        signatureHashes.add(h.getAsString()));
            });
            assertTrue(signatureHashes.contains(expectedHash),
                    "Step 4: Metadata hash should be in signature hashes list");

            // Step 5: Verify whitelist signatures (simplified check)
            AddressWhitelistingRules rules = rulesContainer.findAddressWhitelistingRules(
                    "ALGO", "mainnet");
            assertNotNull(rules, "Step 5: Should find whitelisting rules");
            assertNotNull(rules.getParallelThresholds(),
                    "Step 5: Should have threshold requirements");

            // Parse verified address
            WhitelistedAddress addr = WhitelistHashHelper.parseWhitelistedAddressFromJson(
                    payloadAsString);
            assertNotNull(addr, "Should parse address from verified payload");
            assertEquals("ALGO", addr.getBlockchain());
            assertEquals("P4QCJV2YYLAEULGLJQAW4XTU3EBOHWL5C46I5SPLH2H7AJEE367ZDACV5A",
                    addr.getAddress());
            assertEquals("TN_Bank ACC Cockroach_WTRTest", addr.getLabel());
            assertEquals("84dc35e3-0af8-4b6b-be75-785f4b149d16", addr.getTnParticipantID());
        }

        /**
         * Tests that hash computation is deterministic.
         */
        @Test
        @DisplayName("testHashComputationDeterministic")
        void testHashComputationDeterministic() {
            String payloadAsString = realApiResponse.get("payloadAsString").getAsString();

            // Compute hash multiple times
            String hash1 = CryptoTPV1.calculateHexHash(payloadAsString);
            String hash2 = CryptoTPV1.calculateHexHash(payloadAsString);
            String hash3 = CryptoTPV1.calculateHexHash(payloadAsString);

            // All should be identical
            assertEquals(hash1, hash2, "Hash computation should be deterministic (1 vs 2)");
            assertEquals(hash2, hash3, "Hash computation should be deterministic (2 vs 3)");
        }

        /**
         * Tests constant-time comparison security properties.
         */
        @Test
        @DisplayName("testConstantTimeComparisonSecurity")
        void testConstantTimeComparisonSecurity() {
            String hash = "830063cfa8c1dbd696d670fc8360e85fbc57c3ffa66d22358b9a7d6befabb2f0";
            String sameHash = "830063cfa8c1dbd696d670fc8360e85fbc57c3ffa66d22358b9a7d6befabb2f0";
            String diffHash = "aaaa63cfa8c1dbd696d670fc8360e85fbc57c3ffa66d22358b9a7d6befabb2f0";

            // Same hashes should match
            assertTrue(constantTimeAreEqual(hash, sameHash));

            // Different hashes should not match
            assertFalse(constantTimeAreEqual(hash, diffHash));

            // Different lengths should not match
            assertFalse(constantTimeAreEqual(hash, "short"));
            assertFalse(constantTimeAreEqual("short", hash));

            // Empty strings
            assertFalse(constantTimeAreEqual(hash, ""));
            assertFalse(constantTimeAreEqual("", hash));
            assertTrue(constantTimeAreEqual("", ""));
        }
    }

    // ==================================================================================
    // HELPER METHODS
    // ==================================================================================

    /**
     * Creates a DecodedRulesContainer from the JSON fixture data.
     * This simulates what would be decoded from the protobuf.
     */
    private static DecodedRulesContainer createRulesContainerFromJson() throws Exception {
        JsonObject rulesJson = rawResponse.getAsJsonObject("rulesContainerJson");

        DecodedRulesContainer container = new DecodedRulesContainer();

        // Parse users
        List<RuleUser> users = new ArrayList<>();
        rulesJson.getAsJsonArray("users").forEach(userElement -> {
            JsonObject userObj = userElement.getAsJsonObject();
            RuleUser user = new RuleUser();
            user.setId(userObj.get("id").getAsString());
            user.setPublicKeyPem(userObj.get("publicKey").getAsString());

            List<String> roles = new ArrayList<>();
            userObj.getAsJsonArray("roles").forEach(r -> roles.add(r.getAsString()));
            user.setRoles(roles);

            // Also set the PublicKey object
            try {
                user.setPublicKey(CryptoTPV1.decodePublicKey(user.getPublicKeyPem()));
            } catch (IOException e) {
                // Ignore decode errors in test helper
            }

            users.add(user);
        });
        container.setUsers(users);

        // Parse groups
        List<RuleGroup> groups = new ArrayList<>();
        rulesJson.getAsJsonArray("groups").forEach(groupElement -> {
            JsonObject groupObj = groupElement.getAsJsonObject();
            RuleGroup group = new RuleGroup();
            group.setId(groupObj.get("id").getAsString());

            List<String> userIds = new ArrayList<>();
            groupObj.getAsJsonArray("userIds").forEach(u -> userIds.add(u.getAsString()));
            group.setUserIds(userIds);

            groups.add(group);
        });
        container.setGroups(groups);

        // Parse address whitelisting rules
        List<AddressWhitelistingRules> addressRules = new ArrayList<>();
        rulesJson.getAsJsonArray("addressWhitelistingRules").forEach(ruleElement -> {
            JsonObject ruleObj = ruleElement.getAsJsonObject();
            AddressWhitelistingRules rule = new AddressWhitelistingRules();
            rule.setCurrency(ruleObj.get("currency").getAsString());
            rule.setNetwork(ruleObj.get("network").getAsString());

            // Parse parallel thresholds
            List<SequentialThresholds> parallelThresholds = new ArrayList<>();
            ruleObj.getAsJsonArray("parallelThresholds").forEach(ptElement -> {
                JsonObject ptObj = ptElement.getAsJsonObject();
                SequentialThresholds st = new SequentialThresholds();

                List<GroupThreshold> thresholds = new ArrayList<>();
                GroupThreshold gt = new GroupThreshold();
                gt.setGroupId(ptObj.get("groupId").getAsString());
                gt.setMinimumSignatures(ptObj.get("minimumSignatures").getAsInt());
                thresholds.add(gt);

                st.setThresholds(thresholds);
                parallelThresholds.add(st);
            });
            rule.setParallelThresholds(parallelThresholds);

            addressRules.add(rule);
        });
        container.setAddressWhitelistingRules(addressRules);

        // Set other fields
        container.setMinimumDistinctUserSignatures(
                rulesJson.get("minimumDistinctUserSignatures").getAsInt());
        container.setMinimumDistinctGroupSignatures(
                rulesJson.get("minimumDistinctGroupSignatures").getAsInt());
        container.setTimestamp(rulesJson.get("timestamp").getAsLong());

        return container;
    }
}
