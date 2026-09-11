package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.helper.LegacyPayloadVariant;
import com.taurushq.sdk.protect.client.helper.WhitelistHashHelper;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Tests for legacy payload/hash computation on the whitelisted-address flow.
 * Verifies backward compatibility for addresses signed before schema changes.
 *
 * <p>Case 1 (Address 509): contractType field was added to schema after signing
 * <p>Case 2 (Address 391): contractType AND label in linkedInternalAddresses were added after signing
 *
 * <p><b>Every case below drives the PRODUCTION function.</b> This file used to apply the
 * regexes itself with {@code String.replaceAll}, because the production method was private
 * on {@code WhitelistedAddressService}. That made it a test of a copy: the payload-carrying
 * fix could land in production and every assertion here would still pass, asserting the old
 * semantics. Two sibling files had the same defect. If a future change makes the production
 * function unreachable from here, move the function — do not re-inline the regex.
 */
class WhitelistedAddressServiceLegacyHashTest {

    /** The variant produced by a named strategy, or null when that strategy produced none. */
    private static LegacyPayloadVariant variantWithHash(final String payload, final String hash) {
        for (LegacyPayloadVariant variant : WhitelistHashHelper.computeLegacyPayloadVariants(payload)) {
            if (variant.getHash().equals(hash)) {
                return variant;
            }
        }
        return null;
    }

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
        // Removing contractType from the current payload must reproduce BOTH the legacy
        // hash and the legacy payload — the payload is the half step 6 needs.
        LegacyPayloadVariant variant = variantWithHash(CASE1_CURRENT_PAYLOAD, CASE1_LEGACY_HASH);

        assertTrue(variant != null, "strategy 1 should produce the 2020 legacy hash");
        assertEquals(CASE1_LEGACY_HASH, variant.getHash());
        assertEquals(CASE1_LEGACY_PAYLOAD, variant.getPayload(),
                "the variant must carry the exact bytes the 2020 signer covered");
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
        // A payload with nothing to strip yields no variants at all, rather than a
        // variant equal to the input: an identical "legacy" hash would make the current
        // hash look covered by a legacy signature.
        assertTrue(WhitelistHashHelper.computeLegacyPayloadVariants(
                        "{\"currency\":\"ETH\",\"linkedInternalAddresses\":[]}").isEmpty(),
                "no strip applies, so there must be no variants");
    }

    @Test
    void testCase1_contractTypeWithValue() {
        // Also handle contractType with a non-empty value
        String withContractTypeValue =
                "{\"currency\":\"ETH\",\"linkedInternalAddresses\":[],\"contractType\":\"ERC20\"}";
        List<LegacyPayloadVariant> variants =
                WhitelistHashHelper.computeLegacyPayloadVariants(withContractTypeValue);

        assertEquals(1, variants.size(), "only strategy 1 applies");
        assertEquals("{\"currency\":\"ETH\",\"linkedInternalAddresses\":[]}",
                variants.get(0).getPayload());
    }

    @Test
    void testCase1_hashDifference() {
        // Verify that the two hashes are indeed different
        assertNotEquals(CASE1_LEGACY_HASH, CASE1_CURRENT_HASH);
    }

    // ============ CASE 2: Address 391 - contractType + labels in linkedInternalAddresses ============

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
        LegacyPayloadVariant variant = variantWithHash(CASE2_CURRENT_PAYLOAD, CASE2_LEGACY_HASH);

        assertTrue(variant != null, "strategy 3 should produce the March 2020 legacy hash");
        assertEquals(CASE2_LEGACY_PAYLOAD, variant.getPayload(),
                "the variant must carry the exact bytes the March 2020 signer covered");
    }

    @Test
    void testCase2_labelPatternDoesNotAffectMainLabel() {
        // The main label here is followed by ,"customerId": — not by a closing brace — so
        // the strip leaves it alone. See the NEXT test for why that is a property of THIS
        // FIXTURE and not of the pattern.
        LegacyPayloadVariant strategy2 = variantWithHash(
                CASE2_CURRENT_PAYLOAD, CryptoTPV1.calculateHexHash(STRATEGY2_LEGACY_PAYLOAD));
        assertTrue(strategy2 != null, "strategy 2 should apply to this fixture");

        assertTrue(strategy2.getPayload().contains("\"label\":\"20200324 test address 2\""),
                "Main address label should be preserved");
        assertFalse(strategy2.getPayload().contains("\"label\":\"client 0 ETH \""),
                "Labels inside linkedInternalAddresses should be removed");
        assertFalse(strategy2.getPayload().contains("\"label\":\"ETH internal client 02.02\""),
                "Labels inside linkedInternalAddresses should be removed");
    }

    /**
     * A payload whose MAIN label is the last member of the top-level object. The response
     * shape is chosen by the server, so this is not hypothetical.
     */
    private static final String MAIN_LABEL_LAST_PAYLOAD = "{\"currency\":\"ETH\","
            + "\"address\":\"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380\",\"memo\":\"\","
            + "\"customerId\":\"1\",\"linkedInternalAddresses\":[],"
            + "\"label\":\"20200324 test address 2\"}";

    @Test
    void testCase2_mainLabelAsLastMember_isAlsoStripped() {
        // The comment this file shipped with claimed the pattern "does NOT match the main
        // address label". That is true only when some other field follows it — which the
        // pre-existing fixture happens to guarantee, so the assertion passed against the
        // vulnerable code and proved nothing about the pattern.
        //
        // Put the main label LAST and the strip removes it. THAT is the injection: a server
        // appends a duplicate ,"label":"..." as the final member of a genuinely signed
        // payload, the strip recovers the signed bytes, every signature check passes, and a
        // step 6 that parsed the DELIVERED text would hand back the appended value. It is
        // why the variant carries its payload.
        List<LegacyPayloadVariant> variants =
                WhitelistHashHelper.computeLegacyPayloadVariants(MAIN_LABEL_LAST_PAYLOAD);

        assertEquals(1, variants.size(), "only the label strip applies to this payload");
        assertFalse(variants.get(0).getPayload().contains("20200324 test address 2"),
                "a trailing main label IS removed by the strip — the pattern cannot tell "
                        + "it apart from an inner one, which is exactly the exposure");
        assertEquals("{\"currency\":\"ETH\","
                        + "\"address\":\"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380\",\"memo\":\"\","
                        + "\"customerId\":\"1\",\"linkedInternalAddresses\":[]}",
                variants.get(0).getPayload());
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
        // Strategy 1 alone does NOT reach the legacy hash here: the inner labels are still
        // present. Driven through production, so the strategy ORDER is pinned too.
        List<LegacyPayloadVariant> variants =
                WhitelistHashHelper.computeLegacyPayloadVariants(CASE2_CURRENT_PAYLOAD);

        assertEquals(3, variants.size(), "all three strategies apply to this fixture");
        assertNotEquals(CASE2_LEGACY_HASH, variants.get(0).getHash(),
                "Removing only contractType should not produce the legacy hash");
        assertEquals(CASE2_LEGACY_HASH, variants.get(2).getHash(),
                "strategy 3 is the one that reaches the original 2020 hash");
    }

    // ============ STRATEGY 2: labels added after contractType already existed ============

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
        List<LegacyPayloadVariant> variants =
                WhitelistHashHelper.computeLegacyPayloadVariants(CASE2_CURRENT_PAYLOAD);

        // Strategy 2 is the SECOND variant: remove labels, keep contractType.
        assertEquals(STRATEGY2_LEGACY_PAYLOAD, variants.get(1).getPayload());
        assertTrue(variants.get(1).getPayload().contains("\"contractType\":\"\""),
                "contractType should be preserved");
    }

    @Test
    void testStrategy2_legacyPayloadHash() {
        String hash = CryptoTPV1.calculateHexHash(STRATEGY2_LEGACY_PAYLOAD);

        assertNotEquals(CASE2_CURRENT_HASH, hash,
                "Strategy 2 hash should differ from current hash");
        assertNotEquals(CASE2_LEGACY_HASH, hash,
                "Strategy 2 hash should differ from fully legacy hash (no contractType, no labels)");
    }

    @Test
    void testAllStrategiesProduceDifferentHashes() {
        String currentHash = CryptoTPV1.calculateHexHash(CASE2_CURRENT_PAYLOAD);
        List<LegacyPayloadVariant> variants =
                WhitelistHashHelper.computeLegacyPayloadVariants(CASE2_CURRENT_PAYLOAD);

        assertEquals(3, variants.size(), "three distinct variants for this fixture");
        String strategy1Hash = variants.get(0).getHash();
        String strategy2Hash = variants.get(1).getHash();
        String strategy3Hash = variants.get(2).getHash();

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
