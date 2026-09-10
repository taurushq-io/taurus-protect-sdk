package com.taurushq.sdk.protect.client.cache;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.RulesContainerMapper;
import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.service.GovernanceRuleService;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.Security;
import java.security.spec.ECGenParameterSpec;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;

class RulesContainerCacheTest {

    private static GovernanceRuleService governanceRuleService;

    @BeforeAll
    static void setUpKeys() throws Exception {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC", BouncyCastleProvider.PROVIDER_NAME);
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        KeyPair keyPair = generator.generateKeyPair();
        List<PublicKey> keys = Collections.singletonList(keyPair.getPublic());

        governanceRuleService = new GovernanceRuleService(
                new ApiClient(), new ApiExceptionMapper(), keys, 1);
    }

    // --- Constructor validation ---

    @Test
    void constructor_throwsOnNullGovernanceRuleService() {
        assertThrows(NullPointerException.class, () ->
                new RulesContainerCache(null));
    }

    @Test
    void constructor_throwsOnNullGovernanceRuleServiceWithTtl() {
        assertThrows(NullPointerException.class, () ->
                new RulesContainerCache(null, 5000));
    }

    @Test
    void constructor_throwsOnZeroTtl() {
        assertThrows(IllegalArgumentException.class, () ->
                new RulesContainerCache(governanceRuleService, 0));
    }

    @Test
    void constructor_throwsOnNegativeTtl() {
        assertThrows(IllegalArgumentException.class, () ->
                new RulesContainerCache(governanceRuleService, -1));
    }

    // --- Default TTL ---

    @Test
    void defaultTtl_isFiveMinutes() {
        assertEquals(5 * 60 * 1000L, RulesContainerCache.DEFAULT_CACHE_TTL_MS);
    }

    @Test
    void constructor_defaultTtl_setsFiveMinutes() {
        RulesContainerCache cache = new RulesContainerCache(governanceRuleService);
        assertEquals(RulesContainerCache.DEFAULT_CACHE_TTL_MS, cache.getCacheTtlMs());
    }

    @Test
    void constructor_customTtl_setsCorrectly() {
        RulesContainerCache cache = new RulesContainerCache(governanceRuleService, 10000);
        assertEquals(10000, cache.getCacheTtlMs());
    }

    // --- Cache validity ---

    @Test
    void isCacheValid_returnsFalseWhenEmpty() {
        RulesContainerCache cache = new RulesContainerCache(governanceRuleService);
        assertFalse(cache.isCacheValid());
    }

    // --- getDecodedRulesContainer triggers API call ---
    // The cache calls governanceRuleService.getRules() which makes a network call.
    // Without a mocking framework, we verify the error path:
    // the real API client has no credentials configured, so the call will fail
    // (either ApiException from the server or IllegalArgumentException from HMAC auth).

    @Test
    void getDecodedRulesContainer_throwsWhenApiFails() {
        RulesContainerCache cache = new RulesContainerCache(governanceRuleService);
        assertThrows(Exception.class, () -> cache.getDecodedRulesContainer());
    }

    @Test
    void invalidate_throwsWhenApiFails() {
        RulesContainerCache cache = new RulesContainerCache(governanceRuleService);
        assertThrows(Exception.class, () -> cache.invalidate());
    }

    // --- Cache still invalid after failed fetch ---

    @Test
    void isCacheValid_remainsFalseAfterFailedFetch() {
        RulesContainerCache cache = new RulesContainerCache(governanceRuleService);
        try {
            cache.getDecodedRulesContainer();
        } catch (Exception e) {
            // Expected - API call fails (ApiException or IllegalArgumentException)
        }
        assertFalse(cache.isCacheValid());
    }

    // --- Rules container verification (cross-SDK invariant) ---

    /**
     * The cache supplies the HSM public key that address signature verification
     * trusts, so the container it holds must have its SuperAdmin signatures verified
     * before anything reads it. An unverified container makes that check pass against
     * whatever key the container carries.
     *
     * <p>TypeScript had drifted off this invariant: its cache fetched through the
     * generated API and the raw mapper, skipping verification entirely, and nothing in
     * any SDK tested the path. Asserted here on {@code getDecodedRulesContainer(rules)}
     * because that is the ONLY method the cache fetches through
     * ({@code RulesContainerCache} lines 229 and 237) and this suite has no HTTP stub
     * or mocking library to drive the cache offline.
     *
     * <p>The container is wire-valid on purpose: a malformed one would throw
     * IntegrityException from {@code parseFrom} instead, so the test would pass whether
     * or not verification ran.
     */
    @Test
    void getDecodedRulesContainer_refusesAContainerWithNoSignatures() {
        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(wireValidContainer());
        rules.setRulesSignatures(Collections.<RuleUserSignature>emptyList());

        assertThrows(IntegrityException.class, () ->
                governanceRuleService.getDecodedRulesContainer(rules));
    }

    @Test
    void getDecodedRulesContainer_refusesSignaturesFromAnUnconfiguredKey() {
        RuleUserSignature signature = new RuleUserSignature();
        signature.setUserId("attacker");
        signature.setSignature("YmFkLXNpZy1ub3QtZWNkc2E=");

        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(wireValidContainer());
        rules.setRulesSignatures(Collections.singletonList(signature));

        assertThrows(IntegrityException.class, () ->
                governanceRuleService.getDecodedRulesContainer(rules));
    }

    /** Guards the two tests above from passing on a proto parse error instead. */
    @Test
    void wireValidContainer_decodesCleanly() throws Exception {
        DecodedRulesContainer decoded =
                RulesContainerMapper.INSTANCE.fromBase64String(wireValidContainer());
        assertEquals(2, decoded.getMinimumDistinctUserSignatures());
    }

    private static String wireValidContainer() {
        DecodedRulesContainer container = new DecodedRulesContainer();
        container.setMinimumDistinctUserSignatures(2);
        return RulesContainerMapper.INSTANCE.toBase64String(container);
    }
}
