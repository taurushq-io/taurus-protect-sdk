package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.RulesContainerMapper;
import com.taurushq.sdk.protect.client.model.ExcludedRuleset;
import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.time.OffsetDateTime;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Governance reads must verify SuperAdmin signatures. The container carries the HSM
 * public key that address verification trusts, so an unverified one is the worst object
 * this SDK can hand back — and {@code getRulesHistory} returned it unverified with no
 * test covering that.
 */
class GovernanceRuleVerificationTest {

    private static PublicKey superAdminKey;
    private static String wireValidUnsignedContainer;

    @BeforeAll
    static void setUp() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        superAdminKey = generator.generateKeyPair().getPublic();

        // A container that DECODES cleanly but carries no signatures. It has to decode:
        // a malformed blob throws from the parse instead, and the test then passes
        // whether verification ran or not. This trap has bitten this repo twice.
        DecodedRulesContainer container = new DecodedRulesContainer();
        container.setMinimumDistinctUserSignatures(1);
        wireValidUnsignedContainer = RulesContainerMapper.INSTANCE.toBase64String(container);
    }

    private static GovernanceRuleService service() {
        return new GovernanceRuleService(new ApiClient(), new ApiExceptionMapper(),
                Collections.singletonList(superAdminKey), 1);
    }

    private static GovernanceRules entry(final String container, final OffsetDateTime created) {
        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(container);
        rules.setCreationDate(created);
        return rules;
    }

    /** The guard that keeps the tests below non-vacuous. */
    @Test
    void wireValidContainerDecodesCleanly() throws Exception {
        assertTrue(RulesContainerMapper.INSTANCE.fromBase64String(wireValidUnsignedContainer)
                != null, "fixture must decode, or the verification tests pass for the wrong reason");
    }

    /**
     * History is LENIENT where the single-ruleset reads are strict, deliberately: a
     * SuperAdmin key rotation makes every pre-rotation ruleset unverifiable, so a strict
     * page would deny access to the whole audit trail. It must still NAME what it dropped.
     */
    @Test
    void historyExcludesAndNamesUnverifiedEntries() {
        OffsetDateTime created = OffsetDateTime.now();
        List<ExcludedRuleset> excluded = new ArrayList<>();

        List<GovernanceRules> kept = service().verifiedHistoryEntries(Arrays.asList(
                entry(wireValidUnsignedContainer, created),
                entry(wireValidUnsignedContainer, created)), excluded);

        assertTrue(kept.isEmpty(), "no entry verified, so none may be returned");
        assertEquals(2, excluded.size(), "both entries must be named as excluded");
        for (ExcludedRuleset ex : excluded) {
            assertFalse(ex.getReason() == null || ex.getReason().isEmpty(),
                    "an exclusion must carry a reason");
            assertEquals(created, ex.getCreationDate());
        }
    }

    /**
     * The memo must key on the signature set, not just the container bytes: otherwise
     * re-signing the same container would be served from a stale hit.
     */
    @Test
    void memoKeysOnTheSignatureSet() {
        GovernanceRules base = entry(wireValidUnsignedContainer, null);

        GovernanceRules resigned = entry(wireValidUnsignedContainer, null);
        RuleUserSignature sig = new RuleUserSignature();
        sig.setUserId("u1");
        sig.setSignature("sig");
        resigned.setRulesSignatures(Collections.singletonList(sig));

        assertNotEquals(GovernanceRuleService.rulesetVerificationKey(base),
                GovernanceRuleService.rulesetVerificationKey(resigned),
                "a container with a different signature set must not share a memo key");
    }

    /** Signature order is server-controlled, so it must not change the memo key. */
    @Test
    void memoKeyIgnoresSignatureOrder() {
        RuleUserSignature one = new RuleUserSignature();
        one.setUserId("u1");
        one.setSignature("s1");
        RuleUserSignature two = new RuleUserSignature();
        two.setUserId("u2");
        two.setSignature("s2");

        GovernanceRules a = entry(wireValidUnsignedContainer, null);
        a.setRulesSignatures(Arrays.asList(one, two));
        GovernanceRules b = entry(wireValidUnsignedContainer, null);
        b.setRulesSignatures(Arrays.asList(two, one));

        assertEquals(GovernanceRuleService.rulesetVerificationKey(a),
                GovernanceRuleService.rulesetVerificationKey(b));
    }

    /**
     * userId is not read by verification, so it must not participate in the memo key --
     * otherwise a field the server controls freely and verification ignores can influence
     * a verification-SKIP decision.
     */
    @Test
    void memoKeyIgnoresUserId() {
        RuleUserSignature mine = new RuleUserSignature();
        mine.setUserId("u1");
        mine.setSignature("s1");
        RuleUserSignature theirs = new RuleUserSignature();
        theirs.setUserId("attacker-chosen");
        theirs.setSignature("s1");

        GovernanceRules a = entry(wireValidUnsignedContainer, null);
        a.setRulesSignatures(Collections.singletonList(mine));
        GovernanceRules b = entry(wireValidUnsignedContainer, null);
        b.setRulesSignatures(Collections.singletonList(theirs));

        assertEquals(GovernanceRuleService.rulesetVerificationKey(a),
                GovernanceRuleService.rulesetVerificationKey(b),
                "userId is ignored by verification, so it must not change the memo key");
    }

    /**
     * The memo key must be INJECTIVE. A plain concatenation leaves the boundary between
     * the container and the signature list uncommitted, so a response-controlling
     * attacker can shift bytes across it and make a MODIFIED container collide with a
     * genuine one's key -- inheriting its "already verified" status and skipping ECDSA
     * entirely.
     *
     * <p>Three tail shapes, one per plausible separator-free encoding, so this test is red
     * against ALL of them rather than only the one that happens to be in place.
     */
    @Test
    void memoKeyIsInjective() {
        RuleUserSignature sig = new RuleUserSignature();
        sig.setSignature("s1");
        GovernanceRules signed = entry(wireValidUnsignedContainer, null);
        signed.setRulesSignatures(Collections.singletonList(sig));

        byte[] raw = java.util.Base64.getDecoder().decode(wireValidUnsignedContainer);
        byte[][] tails = {
            "s1".getBytes(java.nio.charset.StandardCharsets.UTF_8),
            new byte[] {0, 's', '1'},
            new byte[] {0, 0, 's', '1'},
        };

        for (byte[] tail : tails) {
            byte[] folded = new byte[raw.length + tail.length];
            System.arraycopy(raw, 0, folded, 0, raw.length);
            System.arraycopy(tail, 0, folded, raw.length, tail.length);
            GovernanceRules foldedRules =
                    entry(java.util.Base64.getEncoder().encodeToString(folded), null);

            assertNotEquals(GovernanceRuleService.rulesetVerificationKey(signed),
                    GovernanceRuleService.rulesetVerificationKey(foldedRules),
                    "distinct (container, signatures) pairs must not share a memo key");
        }
    }

    /**
     * An undecodable container has no stable identity, so it must not be memoised -- and
     * verification must still get the final say.
     */
    @Test
    void memoSkipsUndecodableContainer() {
        assertNull(GovernanceRuleService.rulesetVerificationKey(entry("!!!not base64!!!", null)),
                "an undecodable container must not produce a memo key");
    }

    /**
     * The memo key guards the LIST and the signature STRING; it must guard the ELEMENT
     * too. MapStruct preserves nulls, so a response carrying {@code "rulesSignatures":
     * [null]} arrives as a list holding a null entry, and an unguarded
     * {@code sig.getSignature()} threw NullPointerException — from the memo key, which
     * runs BEFORE verification. One crafted response therefore denied every governance
     * read, and did it with an NPE out of a private helper, which reads as an SDK defect
     * rather than a rejected response.
     */
    @Test
    void memoKeyToleratesANullSignatureElement() {
        GovernanceRules rules = entry(wireValidUnsignedContainer, null);
        rules.setRulesSignatures(Collections.<RuleUserSignature>singletonList(null));

        assertNotNull(GovernanceRuleService.rulesetVerificationKey(rules),
                "a null signature element must not stop the memo key being computed");
    }

    /**
     * The whole read path must reject a null element as an integrity failure, not blow up
     * on it. {@code SignatureVerifier} already skips a null entry, so the only thing
     * standing between a crafted response and an NPE was the memo key.
     */
    @Test
    void aNullSignatureElementIsExcluded_notAnNpe() {
        GovernanceRules rules = entry(wireValidUnsignedContainer, null);
        rules.setRulesSignatures(Collections.<RuleUserSignature>singletonList(null));

        List<ExcludedRuleset> excluded = new ArrayList<>();
        List<GovernanceRules> kept = service().verifiedHistoryEntries(
                Collections.singletonList(rules), excluded);

        assertTrue(kept.isEmpty(), "a null entry signs nothing, so the ruleset cannot verify");
        assertEquals(1, excluded.size(),
                "the row must be excluded and named, not escape as an unchecked NPE");
    }

    /**
     * A null element keeps its SLOT in the key rather than being dropped. It still holds a
     * position in the sorted list, still contributes its own 8-byte length prefix (of
     * zero) and is still counted, so {@code []}, {@code [null]} and {@code [null, null]}
     * stay three distinct keys. Skipping the element instead would collapse them onto one
     * — exactly the non-injectivity the length prefixes exist to prevent.
     */
    @Test
    void memoKeyCountsNullSignatureElements() {
        GovernanceRules none = entry(wireValidUnsignedContainer, null);
        none.setRulesSignatures(Collections.<RuleUserSignature>emptyList());
        GovernanceRules one = entry(wireValidUnsignedContainer, null);
        one.setRulesSignatures(Collections.<RuleUserSignature>singletonList(null));
        GovernanceRules two = entry(wireValidUnsignedContainer, null);
        two.setRulesSignatures(Arrays.<RuleUserSignature>asList(null, null));

        assertNotEquals(GovernanceRuleService.rulesetVerificationKey(none),
                GovernanceRuleService.rulesetVerificationKey(one),
                "an empty list and a list holding one null must not share a memo key");
        assertNotEquals(GovernanceRuleService.rulesetVerificationKey(one),
                GovernanceRuleService.rulesetVerificationKey(two),
                "signature-list length must stay committed to the key");
    }

    /** A failure must resurface on every call rather than being remembered as a verdict. */
    @Test
    void memoDoesNotCacheFailures() {
        GovernanceRuleService svc = service();
        List<GovernanceRules> entries =
                Collections.singletonList(entry(wireValidUnsignedContainer, null));

        List<ExcludedRuleset> first = new ArrayList<>();
        svc.verifiedHistoryEntries(entries, first);
        List<ExcludedRuleset> second = new ArrayList<>();
        svc.verifiedHistoryEntries(entries, second);

        assertEquals(1, first.size());
        assertEquals(1, second.size(),
                "a memoised FAILURE would let the second call through; failures must not be cached");
    }
}
