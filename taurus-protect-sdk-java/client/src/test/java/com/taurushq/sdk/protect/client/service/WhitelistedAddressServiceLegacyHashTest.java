package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Tests for legacy hash computation in WhitelistedAddressService.
 * Verifies backward compatibility for addresses signed before schema changes.
 *
 * <p>Case 1 (Address 509): contractType field was added to schema after signing
 * <p>Case 2 (Address 391): contractType AND label in linkedInternalAddresses were added after signing
 */
class WhitelistedAddressServiceLegacyHashTest {

    // ==================== CASE 1: Address 509 - contractType added ====================

    // Payload WITHOUT contractType (original schema, 2020)
    private static final String CASE1_LEGACY_PAYLOAD = "{\"currency\":\"ETH\",\"addressType\":\"individual\","
            + "\"address\":\"0x012566A179a935ACF1d81d4D237495DE933D12E6\",\"memo\":\"\","
            + "\"label\":\"CMTA20-KYC - 0x012566A179a935ACF1d81d4D237495DE933D12E6 (request 6826)\","
            + "\"customerId\":\"\",\"exchangeAccountId\":\"\",\"linkedInternalAddresses\":[]}";

    // Payload WITH contractType (current schema)
    private static final String CASE1_CURRENT_PAYLOAD = CASE1_LEGACY_PAYLOAD.replace(
            "\"linkedInternalAddresses\":[]}",
            "\"linkedInternalAddresses\":[],\"contractType\":\"\"}");

    // Hash of CASE1_LEGACY_PAYLOAD (what was signed in 2020)
    private static final String CASE1_LEGACY_HASH =
            "cda66e821ec26f2432a717feaa1ef49be39a7ad9e93b6b8fcdce606659e964df";

    // Hash of CASE1_CURRENT_PAYLOAD (what API now returns)
    private static final String CASE1_CURRENT_HASH =
            "d95ae4359bea509c2542acf410649f1e361233da5e1ac7c7a198b6d6a2bbbe1f";

    @Test
    void testCase1_removesContractType() {
        // Verify that removing contractType from current payload produces legacy hash
        String withoutContractType = CASE1_CURRENT_PAYLOAD.replaceAll(",\"contractType\":\"[^\"]*\"", "");
        String computedHash = CryptoTPV1.calculateHexHash(withoutContractType);

        assertEquals(CASE1_LEGACY_HASH, computedHash);
    }

    @Test
    void testCase1_currentPayloadProducesCurrentHash() {
        // Sanity check: verify current payload produces current hash
        String computedHash = CryptoTPV1.calculateHexHash(CASE1_CURRENT_PAYLOAD);
        assertEquals(CASE1_CURRENT_HASH, computedHash);
    }

    @Test
    void testCase1_legacyPayloadProducesLegacyHash() {
        // Sanity check: verify legacy payload produces legacy hash
        String computedHash = CryptoTPV1.calculateHexHash(CASE1_LEGACY_PAYLOAD);
        assertEquals(CASE1_LEGACY_HASH, computedHash);
    }

    @Test
    void testCase1_noContractType_noTransformation() {
        // If payload doesn't have contractType, no transformation should occur
        String withoutContractType = "{\"currency\":\"ETH\",\"linkedInternalAddresses\":[]}";
        String transformed = withoutContractType.replaceAll(",\"contractType\":\"[^\"]*\"", "");

        // No transformation should occur
        assertEquals(withoutContractType, transformed);
    }

    @Test
    void testCase1_contractTypeWithValue() {
        // Also handle contractType with a non-empty value
        String withContractTypeValue =
                "{\"currency\":\"ETH\",\"linkedInternalAddresses\":[],\"contractType\":\"ERC20\"}";
        String transformed = withContractTypeValue.replaceAll(",\"contractType\":\"[^\"]*\"", "");

        assertEquals("{\"currency\":\"ETH\",\"linkedInternalAddresses\":[]}", transformed);
    }

    @Test
    void testCase1_hashDifference() {
        // Verify that the two hashes are indeed different
        assertNotEquals(CASE1_LEGACY_HASH, CASE1_CURRENT_HASH);
    }

    // ==================== CASE 2: Address 391 - contractType + labels in linkedInternalAddresses ====================

    // Current payload with contractType and labels in linkedInternalAddresses
    private static final String CASE2_CURRENT_PAYLOAD = "{\"currency\":\"ETH\",\"addressType\":\"individual\","
            + "\"address\":\"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380\",\"memo\":\"\","
            + "\"label\":\"20200324 test address 2\",\"customerId\":\"1\",\"exchangeAccountId\":\"\","
            + "\"linkedInternalAddresses\":["
            + "{\"id\":\"10\",\"address\":\"0x589ef3d7585f54f0539e24253050887c691c9bd8\",\"label\":\"client 0 ETH \"},"
            + "{\"id\":\"13\",\"address\":\"0x669805f31178faf0dca39c8a5c49ecc531b5156e\","
            + "\"label\":\"ETH internal client 02.02\"},"
            + "{\"id\":\"20\",\"address\":\"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d\",\"label\":\"LBR 07.02\"},"
            + "{\"id\":\"98\",\"address\":\"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e\","
            + "\"label\":\"ETH LBR internal client 26.02\"},"
            + "{\"id\":\"25\",\"address\":\"0x9bc28e6710f5bb2511372987f613a436618e28ad\",\"label\":\"LBR IC 13.02\"}],"
            + "\"contractType\":\"\"}";

    // Original payload without contractType and without labels in linkedInternalAddresses
    private static final String CASE2_LEGACY_PAYLOAD = "{\"currency\":\"ETH\",\"addressType\":\"individual\","
            + "\"address\":\"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380\",\"memo\":\"\","
            + "\"label\":\"20200324 test address 2\",\"customerId\":\"1\",\"exchangeAccountId\":\"\","
            + "\"linkedInternalAddresses\":["
            + "{\"id\":\"10\",\"address\":\"0x589ef3d7585f54f0539e24253050887c691c9bd8\"},"
            + "{\"id\":\"13\",\"address\":\"0x669805f31178faf0dca39c8a5c49ecc531b5156e\"},"
            + "{\"id\":\"20\",\"address\":\"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d\"},"
            + "{\"id\":\"98\",\"address\":\"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e\"},"
            + "{\"id\":\"25\",\"address\":\"0x9bc28e6710f5bb2511372987f613a436618e28ad\"}]}";

    // Hash of CASE2_LEGACY_PAYLOAD (what was signed in March 2020)
    private static final String CASE2_LEGACY_HASH =
            "88e4e456f7ca1fc4ca415c6c571f828c0eb047e9f15f36d547c103b2ea0def9b";

    // Hash of CASE2_CURRENT_PAYLOAD (what API now returns)
    private static final String CASE2_CURRENT_HASH =
            "7d62d7f78ed55c716ea1278473d6cac5b60a31e1df941873118932822df32b03";

    @Test
    void testCase2_removesContractTypeAndLabelsInObjects() {
        // Step 1: Remove labels inside objects (pattern: ,"label":"..."}  ->  })
        String step1 = CASE2_CURRENT_PAYLOAD.replaceAll(",\"label\":\"[^\"]*\"}", "}");

        // Step 2: Remove contractType
        String step2 = step1.replaceAll(",\"contractType\":\"[^\"]*\"", "");

        assertEquals(CASE2_LEGACY_HASH, CryptoTPV1.calculateHexHash(step2));
    }

    @Test
    void testCase2_labelPatternDoesNotAffectMainLabel() {
        // Verify that ,"label":"[^"]*"} pattern does NOT match the main address label
        // Main label: ,"label":"20200324 test address 2","customerId":...  (NOT followed by })
        String transformed = CASE2_CURRENT_PAYLOAD.replaceAll(",\"label\":\"[^\"]*\"}", "}");

        // Main label should still be present
        assertTrue(transformed.contains("\"label\":\"20200324 test address 2\""),
                "Main address label should be preserved");

        // Labels in linkedInternalAddresses should be removed
        assertFalse(transformed.contains("\"label\":\"client 0 ETH \""),
                "Labels inside linkedInternalAddresses should be removed");
        assertFalse(transformed.contains("\"label\":\"ETH internal client 02.02\""),
                "Labels inside linkedInternalAddresses should be removed");
    }

    @Test
    void testCase2_currentPayloadProducesCurrentHash() {
        // Sanity check: verify current payload produces current hash
        String computedHash = CryptoTPV1.calculateHexHash(CASE2_CURRENT_PAYLOAD);
        assertEquals(CASE2_CURRENT_HASH, computedHash);
    }

    @Test
    void testCase2_legacyPayloadProducesLegacyHash() {
        // Sanity check: verify legacy payload produces legacy hash
        String computedHash = CryptoTPV1.calculateHexHash(CASE2_LEGACY_PAYLOAD);
        assertEquals(CASE2_LEGACY_HASH, computedHash);
    }

    @Test
    void testCase2_hashDifference() {
        // Verify that the two hashes are indeed different
        assertNotEquals(CASE2_LEGACY_HASH, CASE2_CURRENT_HASH);
    }

    @Test
    void testCase2_onlyRemovingContractTypeIsNotEnough() {
        // Verify that only removing contractType does NOT produce the legacy hash
        String withoutContractType = CASE2_CURRENT_PAYLOAD.replaceAll(",\"contractType\":\"[^\"]*\"", "");
        String hashWithoutContractType = CryptoTPV1.calculateHexHash(withoutContractType);

        // This should NOT match because labels in linkedInternalAddresses are still present
        assertNotEquals(CASE2_LEGACY_HASH, hashWithoutContractType,
                "Removing only contractType should not produce the legacy hash");
    }

    // ==================== STRATEGY 2: labels added after contractType already existed ====================

    // Payload WITH contractType but WITHOUT labels in linkedInternalAddresses
    // This represents addresses signed after contractType was added but before labels were added
    private static final String STRATEGY2_LEGACY_PAYLOAD = "{\"currency\":\"ETH\",\"addressType\":\"individual\","
            + "\"address\":\"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380\",\"memo\":\"\","
            + "\"label\":\"20200324 test address 2\",\"customerId\":\"1\",\"exchangeAccountId\":\"\","
            + "\"linkedInternalAddresses\":["
            + "{\"id\":\"10\",\"address\":\"0x589ef3d7585f54f0539e24253050887c691c9bd8\"},"
            + "{\"id\":\"13\",\"address\":\"0x669805f31178faf0dca39c8a5c49ecc531b5156e\"},"
            + "{\"id\":\"20\",\"address\":\"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d\"},"
            + "{\"id\":\"98\",\"address\":\"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e\"},"
            + "{\"id\":\"25\",\"address\":\"0x9bc28e6710f5bb2511372987f613a436618e28ad\"}],"
            + "\"contractType\":\"\"}";

    @Test
    void testStrategy2_removesLabelsOnlyKeepsContractType() {
        // Strategy 2: Remove labels from linkedInternalAddresses objects only (keep contractType)
        // Handles addresses signed after contractType was added but before labels were added
        String withoutLabels = CASE2_CURRENT_PAYLOAD.replaceAll(",\"label\":\"[^\"]*\"}", "}");

        // The transformed payload should have contractType but no labels in linkedInternalAddresses
        assertTrue(withoutLabels.contains("\"contractType\":\"\""),
                "contractType should be preserved");
        assertFalse(withoutLabels.contains("\"label\":\"client 0 ETH \""),
                "Labels inside linkedInternalAddresses should be removed");

        // Verify it matches the STRATEGY2_LEGACY_PAYLOAD
        assertEquals(STRATEGY2_LEGACY_PAYLOAD, withoutLabels);
    }

    @Test
    void testStrategy2_legacyPayloadHash() {
        // Compute the hash of STRATEGY2_LEGACY_PAYLOAD (with contractType, without labels)
        String hash = CryptoTPV1.calculateHexHash(STRATEGY2_LEGACY_PAYLOAD);

        // This hash is different from both CASE2_CURRENT_HASH and CASE2_LEGACY_HASH
        assertNotEquals(CASE2_CURRENT_HASH, hash,
                "Strategy 2 hash should differ from current hash");
        assertNotEquals(CASE2_LEGACY_HASH, hash,
                "Strategy 2 hash should differ from fully legacy hash (no contractType, no labels)");
    }

    @Test
    void testAllStrategiesProduceDifferentHashes() {
        // Current payload (with both contractType and labels)
        String currentHash = CryptoTPV1.calculateHexHash(CASE2_CURRENT_PAYLOAD);

        // Strategy 1: Remove contractType only
        String withoutContractType = CASE2_CURRENT_PAYLOAD.replaceAll(",\"contractType\":\"[^\"]*\"", "");
        String strategy1Hash = CryptoTPV1.calculateHexHash(withoutContractType);

        // Strategy 2: Remove labels only (keep contractType)
        String withoutLabels = CASE2_CURRENT_PAYLOAD.replaceAll(",\"label\":\"[^\"]*\"}", "}");
        String strategy2Hash = CryptoTPV1.calculateHexHash(withoutLabels);

        // Strategy 3: Remove both
        String withoutBoth = CASE2_CURRENT_PAYLOAD.replaceAll(",\"label\":\"[^\"]*\"}", "}");
        withoutBoth = withoutBoth.replaceAll(",\"contractType\":\"[^\"]*\"", "");
        String strategy3Hash = CryptoTPV1.calculateHexHash(withoutBoth);

        // All four hashes should be different
        assertNotEquals(currentHash, strategy1Hash, "Current vs Strategy 1");
        assertNotEquals(currentHash, strategy2Hash, "Current vs Strategy 2");
        assertNotEquals(currentHash, strategy3Hash, "Current vs Strategy 3");
        assertNotEquals(strategy1Hash, strategy2Hash, "Strategy 1 vs Strategy 2");
        assertNotEquals(strategy1Hash, strategy3Hash, "Strategy 1 vs Strategy 3");
        assertNotEquals(strategy2Hash, strategy3Hash, "Strategy 2 vs Strategy 3");

        // Strategy 3 should match CASE2_LEGACY_HASH (the original hash from 2020)
        assertEquals(CASE2_LEGACY_HASH, strategy3Hash, "Strategy 3 should match original legacy hash");
    }
}
