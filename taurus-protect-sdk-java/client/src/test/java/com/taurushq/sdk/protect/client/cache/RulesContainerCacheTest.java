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
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicReference;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class RulesContainerCacheTest {

    /** Bounds the wait for a second caller: a wedged slot must FAIL, never hang the suite. */
    private static final long WEDGE_TIMEOUT_SECONDS = 5;

    private static GovernanceRuleService governanceRuleService;
    private static List<PublicKey> superAdminKeys;

    @BeforeAll
    static void setUpKeys() throws Exception {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC", BouncyCastleProvider.PROVIDER_NAME);
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        KeyPair keyPair = generator.generateKeyPair();
        superAdminKeys = Collections.singletonList(keyPair.getPublic());

        governanceRuleService = new GovernanceRuleService(
                new ApiClient(), new ApiExceptionMapper(), superAdminKeys, 1);
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

    // --- Single-flight flag lifecycle (ported from Python's
    // --- test_failed_fetch_still_releases_the_slot) ---

    /**
     * A governance service whose fetch fails exactly the way a hostile {@code /rules}
     * response makes it fail: the container is wire-valid but carries no SuperAdmin
     * signatures, so {@code getDecodedRulesContainer} — the one method the cache's own
     * {@code doFetch()} verifies through — raises {@link IntegrityException}. That
     * exception extends {@code SecurityException} and is therefore UNCHECKED, which is
     * the entire point: the cache's {@code catch (ApiException e)} never sees it.
     *
     * <p>Overriding {@code getRules()} is the least-hacky injection available. This suite
     * has JUnit only — no mocking library, no HTTP stub — and the cache deliberately
     * exposes no setter and no fetcher hook (Go's {@code Set}/{@code SetFetcher} were
     * REMOVED because either could seat a container nothing had verified; do not add a
     * Java equivalent to make a test easier). Overriding the fetch input rather than the
     * fetch outcome also keeps the SDK's real verification code, not a stub, as the thing
     * that throws.
     */
    private static final class UnverifiableRulesService extends GovernanceRuleService {

        UnverifiableRulesService() {
            super(new ApiClient(), new ApiExceptionMapper(), superAdminKeys, 1);
        }

        @Override
        public GovernanceRules getRules() {
            GovernanceRules rules = new GovernanceRules();
            rules.setRulesContainer(wireValidContainer());
            rules.setRulesSignatures(Collections.<RuleUserSignature>emptyList());
            return rules;
        }
    }

    /** A cache call, so the two tests below can share the second-caller harness. */
    private interface CacheCall {
        void invoke() throws Exception;
    }

    /**
     * Runs {@code call} on another thread and returns what it threw, failing the test if
     * it has not finished within {@link #WEDGE_TIMEOUT_SECONDS}.
     *
     * <p>The bounded wait is what turns the wedge into a red test instead of a suite that
     * never finishes: the cache waits on {@code lock.wait()} with no timeout, so a caller
     * that never gets notified never comes back. The thread is a daemon so that a parked
     * one cannot keep the surefire JVM alive after the failure has been reported.
     *
     * @param call what the second caller should do
     * @return the exception it threw
     * @throws InterruptedException if this thread is interrupted while waiting
     */
    private static Exception failureOnAnotherThread(final CacheCall call) throws InterruptedException {
        final AtomicReference<Exception> outcome = new AtomicReference<>();
        final CountDownLatch finished = new CountDownLatch(1);

        Thread second = new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    call.invoke();
                    outcome.set(new IllegalStateException(
                            "the fetch was expected to fail again, not succeed"));
                } catch (Exception e) {
                    // Deliberately broad: the subject of the test is WHICH exception comes
                    // out, so it has to be captured rather than filtered.
                    outcome.set(e);
                } finally {
                    finished.countDown();
                }
            }
        });
        second.setDaemon(true);
        second.start();

        assertTrue(finished.await(WEDGE_TIMEOUT_SECONDS, TimeUnit.SECONDS),
                "a second caller never returned. The failed fetch left `fetching` true, so it "
                        + "parked on the untimed lock.wait() and never woke: one crafted /rules "
                        + "response permanently wedges every address, asset and price "
                        + "verification in the process");
        return outcome.get();
    }

    /**
     * The single-flight slot must be released even when the fetch fails with an UNCHECKED
     * exception. {@code fetching = true} is set under the lock and {@code doFetch()} then
     * runs outside it; when the reset sat after the try rather than in a finally, the
     * unchecked {@link IntegrityException} from container verification skipped it and left
     * the flag true for the life of the process.
     */
    @Test
    void getDecodedRulesContainer_failedFetchReleasesTheSlot() throws Exception {
        final RulesContainerCache cache = new RulesContainerCache(new UnverifiableRulesService());

        assertThrows(IntegrityException.class, () -> cache.getDecodedRulesContainer(),
                "an unsigned container must be refused");

        Exception second = failureOnAnotherThread(() -> cache.getDecodedRulesContainer());

        assertTrue(second instanceof IntegrityException,
                "the retry must surface the real integrity failure rather than a null container "
                        + "or a stale one; got " + second);
        assertFalse(cache.isCacheValid(),
                "a failed fetch must not seat or freshen a container");
    }

    /** {@code invalidate()} sets the same flag, so it needs the same finally. */
    @Test
    void invalidate_failedFetchReleasesTheSlot() throws Exception {
        final RulesContainerCache cache = new RulesContainerCache(new UnverifiableRulesService());

        assertThrows(IntegrityException.class, () -> cache.invalidate(),
                "an unsigned container must be refused");

        Exception second = failureOnAnotherThread(() -> cache.invalidate());

        assertTrue(second instanceof IntegrityException,
                "the retry must surface the real integrity failure; got " + second);
        assertFalse(cache.isCacheValid(),
                "a failed invalidate must not seat or freshen a container");
    }
}
