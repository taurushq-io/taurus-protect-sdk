package com.taurushq.sdk.protect.client.helper;

import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.apache.commons.codec.binary.Base64;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.Security;
import java.security.spec.ECGenParameterSpec;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class SignatureVerifierTest {

    private static KeyPair keyPair1;
    private static KeyPair keyPair2;

    @BeforeAll
    static void setUpKeys() throws Exception {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC", BouncyCastleProvider.PROVIDER_NAME);
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        keyPair1 = generator.generateKeyPair();
        keyPair2 = generator.generateKeyPair();
    }

    // --- distinct-signer counting ---

    private static RuleUserSignature signatureEntry(String userId, String signature) {
        RuleUserSignature sig = new RuleUserSignature();
        sig.setUserId(userId);
        sig.setSignature(signature);
        return sig;
    }

    private static GovernanceRules rulesWith(String rulesBase64, List<RuleUserSignature> signatures) {
        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(rulesBase64);
        rules.setRulesSignatures(signatures);
        return rules;
    }

    /**
     * The threshold counts distinct signing keys. Counting signature entries would let one
     * compromised key satisfy any threshold, since ECDSA is randomized and a single key can
     * emit unlimited distinct valid signatures over the same data.
     *
     * <p>These five cases are mirrored in all four SDKs and are kept as separate tests on
     * purpose: packed into one @Test, JUnit's fail-fast meant a break in the second case
     * stopped the remaining three from ever running and reported one failure instead of
     * four. Go uses subtests and Python/TypeScript use separate tests for the same reason.
     */
    @Test
    void distinctSigners_twoDistinctKeysMeetThresholdTwo() throws Exception {
        Fixture f = new Fixture();

        assertDoesNotThrow(() -> SignatureVerifier.verifyGovernanceRules(
                rulesWith(f.rulesBase64, Arrays.asList(
                        signatureEntry("admin1", f.key1SigA), signatureEntry("admin2", f.key2Sig))),
                2, f.bothKeys));
    }

    @Test
    void distinctSigners_twoSignaturesFromOneKeyFailAtTwoButPassAtOne() throws Exception {
        Fixture f = new Fixture();
        List<RuleUserSignature> sameKeyTwice = Arrays.asList(
                signatureEntry("admin1", f.key1SigA), signatureEntry("admin1-again", f.key1SigB));

        assertThrows(IntegrityException.class, () -> SignatureVerifier.verifyGovernanceRules(
                rulesWith(f.rulesBase64, sameKeyTwice), 2, f.bothKeys));
        assertDoesNotThrow(() -> SignatureVerifier.verifyGovernanceRules(
                rulesWith(f.rulesBase64, sameKeyTwice), 1, f.bothKeys));
    }

    @Test
    void distinctSigners_replayedIdenticalSignatureCountsOnce() throws Exception {
        Fixture f = new Fixture();

        assertThrows(IntegrityException.class, () -> SignatureVerifier.verifyGovernanceRules(
                rulesWith(f.rulesBase64, Arrays.asList(
                        signatureEntry("admin1", f.key1SigA), signatureEntry("admin1-replay", f.key1SigA))),
                2, f.bothKeys));
    }

    @Test
    void distinctSigners_sameKeyConfiguredTwiceCountsOnce() throws Exception {
        Fixture f = new Fixture();
        List<RuleUserSignature> sameKeyTwice = Arrays.asList(
                signatureEntry("admin1", f.key1SigA), signatureEntry("admin1-again", f.key1SigB));

        assertThrows(IntegrityException.class, () -> SignatureVerifier.verifyGovernanceRules(
                rulesWith(f.rulesBase64, sameKeyTwice), 2, f.keyOneTwice));
    }

    @Test
    void distinctSigners_unconfiguredKeyContributesNothing() throws Exception {
        Fixture f = new Fixture();

        assertThrows(IntegrityException.class, () -> SignatureVerifier.verifyGovernanceRules(
                rulesWith(f.rulesBase64, Arrays.asList(
                        signatureEntry("admin1", f.key1SigA), signatureEntry("stranger", f.key2Sig))),
                2, Collections.singletonList(keyPair1.getPublic())));
    }

    /** Signatures and key lists shared by the five distinct-signer cases. */
    private final class Fixture {
        private final String rulesBase64;
        private final String key1SigA;
        private final String key1SigB;
        private final String key2Sig;
        private final List<PublicKey> bothKeys;
        private final List<PublicKey> keyOneTwice;

        private Fixture() throws Exception {
            byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);
            rulesBase64 = Base64.encodeBase64String(rulesData);
            key1SigA = CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData);
            key1SigB = CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData);
            key2Sig = CryptoTPV1.calculateBase64Signature(keyPair2.getPrivate(), rulesData);
            bothKeys = Arrays.asList(keyPair1.getPublic(), keyPair2.getPublic());
            keyOneTwice = Arrays.asList(keyPair1.getPublic(), keyPair1.getPublic());
        }
    }

    // --- verifySignature() tests ---

    @Test
    void verifySignature_withValidSignature_returnsTrue() throws Exception {
        byte[] data = "test data".getBytes(StandardCharsets.UTF_8);
        java.security.Signature signer = java.security.Signature.getInstance("SHA256withPLAIN-ECDSA");
        signer.initSign(keyPair1.getPrivate());
        signer.update(data);
        byte[] signature = signer.sign();

        assertTrue(SignatureVerifier.verifySignature(data, signature, keyPair1.getPublic()));
    }

    @Test
    void verifySignature_withWrongKey_returnsFalse() throws Exception {
        byte[] data = "test data".getBytes(StandardCharsets.UTF_8);
        java.security.Signature signer = java.security.Signature.getInstance("SHA256withPLAIN-ECDSA");
        signer.initSign(keyPair1.getPrivate());
        signer.update(data);
        byte[] signature = signer.sign();

        // Verify with wrong public key
        assertFalse(SignatureVerifier.verifySignature(data, signature, keyPair2.getPublic()));
    }

    @Test
    void verifySignature_withCorruptedSignature_returnsFalse() throws Exception {
        byte[] data = "test data".getBytes(StandardCharsets.UTF_8);
        java.security.Signature signer = java.security.Signature.getInstance("SHA256withPLAIN-ECDSA");
        signer.initSign(keyPair1.getPrivate());
        signer.update(data);
        byte[] signature = signer.sign();

        // Corrupt the signature
        signature[0] ^= 0xFF;

        assertFalse(SignatureVerifier.verifySignature(data, signature, keyPair1.getPublic()));
    }

    @Test
    void verifySignature_throwsOnNullData() {
        assertThrows(NullPointerException.class, () ->
                SignatureVerifier.verifySignature(null, new byte[]{1}, keyPair1.getPublic()));
    }

    @Test
    void verifySignature_throwsOnNullSignature() {
        assertThrows(NullPointerException.class, () ->
                SignatureVerifier.verifySignature(new byte[]{1}, null, keyPair1.getPublic()));
    }

    @Test
    void verifySignature_throwsOnNullPublicKey() {
        assertThrows(NullPointerException.class, () ->
                SignatureVerifier.verifySignature(new byte[]{1}, new byte[]{1}, null));
    }

    // --- isValidSignature() tests ---

    @Test
    void isValidSignature_withMultipleKeysOneValid_returnsTrue() throws Exception {
        byte[] data = "governance rules data".getBytes(StandardCharsets.UTF_8);
        String base64Sig = CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), data);

        List<PublicKey> keys = Arrays.asList(keyPair2.getPublic(), keyPair1.getPublic());
        assertTrue(SignatureVerifier.isValidSignature(data, base64Sig, keys));
    }

    @Test
    void isValidSignature_withNoValidKey_returnsFalse() throws Exception {
        byte[] data = "governance rules data".getBytes(StandardCharsets.UTF_8);
        String base64Sig = CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), data);

        // Use only keyPair2 for verification (doesn't match)
        List<PublicKey> keys = Collections.singletonList(keyPair2.getPublic());
        assertFalse(SignatureVerifier.isValidSignature(data, base64Sig, keys));
    }

    @Test
    void isValidSignature_withCorruptedBase64Signature_returnsFalse() {
        byte[] data = "governance rules data".getBytes(StandardCharsets.UTF_8);
        String corruptedSig = Base64.encodeBase64String(new byte[64]); // all zeros

        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        assertFalse(SignatureVerifier.isValidSignature(data, corruptedSig, keys));
    }

    // --- verifyGovernanceRules() tests ---

    @Test
    void verifyGovernanceRules_throwsOnNullRules() {
        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        assertThrows(NullPointerException.class, () ->
                SignatureVerifier.verifyGovernanceRules(null, 1, keys));
    }

    @Test
    void verifyGovernanceRules_throwsOnZeroMinSignatures() {
        GovernanceRules rules = new GovernanceRules();
        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        assertThrows(IllegalArgumentException.class, () ->
                SignatureVerifier.verifyGovernanceRules(rules, 0, keys));
    }

    @Test
    void verifyGovernanceRules_throwsOnNegativeMinSignatures() {
        GovernanceRules rules = new GovernanceRules();
        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        assertThrows(IllegalArgumentException.class, () ->
                SignatureVerifier.verifyGovernanceRules(rules, -1, keys));
    }

    @Test
    void verifyGovernanceRules_throwsOnNullKeys() {
        GovernanceRules rules = new GovernanceRules();
        assertThrows(NullPointerException.class, () ->
                SignatureVerifier.verifyGovernanceRules(rules, 1, null));
    }

    @Test
    void verifyGovernanceRules_throwsOnEmptyKeys() {
        GovernanceRules rules = new GovernanceRules();
        List<PublicKey> keys = Collections.emptyList();
        assertThrows(IllegalArgumentException.class, () ->
                SignatureVerifier.verifyGovernanceRules(rules, 1, keys));
    }

    @Test
    void verifyGovernanceRules_throwsOnNullRulesContainer() {
        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(null);
        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        IntegrityException ex = assertThrows(IntegrityException.class, () ->
                SignatureVerifier.verifyGovernanceRules(rules, 1, keys));
        assertTrue(ex.getMessage().contains("rulesContainer is null"));
    }

    @Test
    void verifyGovernanceRules_throwsOnEmptySignatures() {
        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(Base64.encodeBase64String("data".getBytes(StandardCharsets.UTF_8)));
        rules.setRulesSignatures(Collections.emptyList());
        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        IntegrityException ex = assertThrows(IntegrityException.class, () ->
                SignatureVerifier.verifyGovernanceRules(rules, 1, keys));
        assertTrue(ex.getMessage().contains("no signatures present"));
    }

    @Test
    void verifyGovernanceRules_throwsOnNullSignatures() {
        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(Base64.encodeBase64String("data".getBytes(StandardCharsets.UTF_8)));
        rules.setRulesSignatures(null);
        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        IntegrityException ex = assertThrows(IntegrityException.class, () ->
                SignatureVerifier.verifyGovernanceRules(rules, 1, keys));
        assertTrue(ex.getMessage().contains("no signatures present"));
    }

    @Test
    void verifyGovernanceRules_withValidSignature_succeeds() throws Exception {
        byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);
        String rulesBase64 = Base64.encodeBase64String(rulesData);

        String base64Sig = CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData);

        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(rulesBase64);
        RuleUserSignature sig = new RuleUserSignature();
        sig.setUserId("admin1");
        sig.setSignature(base64Sig);
        rules.setRulesSignatures(Collections.singletonList(sig));

        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        assertDoesNotThrow(() -> SignatureVerifier.verifyGovernanceRules(rules, 1, keys));
    }

    @Test
    void verifyGovernanceRules_withInsufficientSignatures_throws() throws Exception {
        byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);
        String rulesBase64 = Base64.encodeBase64String(rulesData);

        String base64Sig = CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData);

        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(rulesBase64);
        RuleUserSignature sig = new RuleUserSignature();
        sig.setUserId("admin1");
        sig.setSignature(base64Sig);
        rules.setRulesSignatures(Collections.singletonList(sig));

        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        // Require 2 signatures but only 1 is valid
        IntegrityException ex = assertThrows(IntegrityException.class, () ->
                SignatureVerifier.verifyGovernanceRules(rules, 2, keys));
        assertTrue(ex.getMessage().contains("only 1 distinct valid signers found, minimum 2 required"));
    }

    @Test
    void verifyGovernanceRules_withInvalidSignature_throws() throws Exception {
        byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);
        String rulesBase64 = Base64.encodeBase64String(rulesData);

        // Sign with keyPair2 but verify with keyPair1
        String base64Sig = CryptoTPV1.calculateBase64Signature(keyPair2.getPrivate(), rulesData);

        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(rulesBase64);
        RuleUserSignature sig = new RuleUserSignature();
        sig.setUserId("admin1");
        sig.setSignature(base64Sig);
        rules.setRulesSignatures(Collections.singletonList(sig));

        List<PublicKey> keys = Collections.singletonList(keyPair1.getPublic());
        IntegrityException ex = assertThrows(IntegrityException.class, () ->
                SignatureVerifier.verifyGovernanceRules(rules, 1, keys));
        assertTrue(ex.getMessage().contains("only 0 distinct valid signers found, minimum 1 required"));
    }

    @Test
    void verifyGovernanceRules_withMultipleSignatures_countsValid() throws Exception {
        byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);
        String rulesBase64 = Base64.encodeBase64String(rulesData);

        String sig1Base64 = CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData);
        String sig2Base64 = CryptoTPV1.calculateBase64Signature(keyPair2.getPrivate(), rulesData);

        RuleUserSignature ruleSig1 = new RuleUserSignature();
        ruleSig1.setUserId("admin1");
        ruleSig1.setSignature(sig1Base64);

        RuleUserSignature ruleSig2 = new RuleUserSignature();
        ruleSig2.setUserId("admin2");
        ruleSig2.setSignature(sig2Base64);

        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(rulesBase64);
        rules.setRulesSignatures(Arrays.asList(ruleSig1, ruleSig2));

        List<PublicKey> keys = Arrays.asList(keyPair1.getPublic(), keyPair2.getPublic());
        // Require 2 valid signatures, both should pass
        assertDoesNotThrow(() -> SignatureVerifier.verifyGovernanceRules(rules, 2, keys));
    }

    @Test
    void verifyGovernanceRules_withNullSignatureEntry_skipsItAndKeepsCounting() throws Exception {
        byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);
        String rulesBase64 = Base64.encodeBase64String(rulesData);

        // A partially-populated approval row must simply not count, matching the explicit
        // skip in the other three SDKs, and must not abort the remaining signers.
        RuleUserSignature nullSig = new RuleUserSignature();
        nullSig.setUserId("pending-approver");
        nullSig.setSignature(null);

        RuleUserSignature emptySig = new RuleUserSignature();
        emptySig.setUserId("also-pending");
        emptySig.setSignature("");

        RuleUserSignature validSig = new RuleUserSignature();
        validSig.setUserId("admin1");
        validSig.setSignature(CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData));

        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(rulesBase64);
        rules.setRulesSignatures(Arrays.asList(nullSig, emptySig, validSig));

        List<PublicKey> keys = Arrays.asList(keyPair1.getPublic(), keyPair2.getPublic());
        assertDoesNotThrow(() -> SignatureVerifier.verifyGovernanceRules(rules, 1, keys));
    }

    @Test
    void verifyGovernanceRulesSignatures_oneKeyCannotMeetThresholdTwo() throws Exception {
        byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);

        RuleUserSignature sig1 = new RuleUserSignature();
        sig1.setUserId("attacker-label-1");
        sig1.setSignature(CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData));

        RuleUserSignature sig2 = new RuleUserSignature();
        sig2.setUserId("attacker-label-2");
        sig2.setSignature(CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData));

        List<PublicKey> keys = Arrays.asList(keyPair1.getPublic(), keyPair2.getPublic());
        IntegrityException ex = assertThrows(IntegrityException.class, () ->
                SignatureVerifier.verifyGovernanceRulesSignatures(
                        rulesData, Arrays.asList(sig1, sig2), keys, 2));
        assertTrue(ex.getMessage().contains("only 1 distinct valid signers found"));
    }

    /**
     * verifyHashCoverage is the peer of Go VerifyHashCoverage, Python
     * verify_hash_coverage and TS verifyHashCoverage; this SDK was the only one without
     * it, leaving callers to compare hashes by hand.
     */
    @Test
    void verifyHashCoverage_findsAHashCoveredByAnySignature() {
        WhitelistSignature first = new WhitelistSignature();
        first.getHashes().addAll(Arrays.asList("aaa", "bbb"));
        WhitelistSignature second = new WhitelistSignature();
        second.getHashes().add("ccc");
        List<WhitelistSignature> signatures = Arrays.asList(first, second);

        assertTrue(SignatureVerifier.verifyHashCoverage("bbb", signatures));
        assertTrue(SignatureVerifier.verifyHashCoverage("ccc", signatures));
        assertFalse(SignatureVerifier.verifyHashCoverage("zzz", signatures));
    }

    @Test
    void verifyHashCoverage_isFalseForEmptyOrMissingInput() {
        WhitelistSignature withHashes = new WhitelistSignature();
        withHashes.getHashes().add("aaa");

        assertFalse(SignatureVerifier.verifyHashCoverage(null, Collections.singletonList(withHashes)));
        assertFalse(SignatureVerifier.verifyHashCoverage("", Collections.singletonList(withHashes)));
        assertFalse(SignatureVerifier.verifyHashCoverage("aaa", null));
        assertFalse(SignatureVerifier.verifyHashCoverage("aaa", Collections.<WhitelistSignature>emptyList()));
        // A signature with no hashes must be skipped, not throw.
        assertFalse(SignatureVerifier.verifyHashCoverage("aaa",
                Collections.singletonList(new WhitelistSignature())));
    }

    /**
     * containsHash is the per-signature half of the pair. Both whitelist services
     * used to carry their own {@code List.contains} copies, which compare with
     * String.equals and return on the first match; those are gone and every call
     * site routes here.
     */
    @Test
    void containsHash_findsAHashInOneSignaturesList() {
        List<String> hashes = Arrays.asList("aaa", "bbb", "ccc");

        assertTrue(SignatureVerifier.containsHash(hashes, "aaa"));
        assertTrue(SignatureVerifier.containsHash(hashes, "ccc"));
        assertFalse(SignatureVerifier.containsHash(hashes, "zzz"));
    }

    @Test
    void containsHash_isFalseForEmptyOrMissingInput() {
        List<String> hashes = Arrays.asList("aaa");

        assertFalse(SignatureVerifier.containsHash(null, "aaa"));
        assertFalse(SignatureVerifier.containsHash(hashes, null));
        assertFalse(SignatureVerifier.containsHash(hashes, ""));
        assertFalse(SignatureVerifier.containsHash(Collections.<String>emptyList(), "aaa"));
    }

    @Test
    void containsHash_skipsNullEntriesRatherThanThrowing() {
        List<String> hashes = Arrays.asList(null, "bbb", null);

        assertTrue(SignatureVerifier.containsHash(hashes, "bbb"));
        assertFalse(SignatureVerifier.containsHash(hashes, "aaa"));
    }

    /**
     * A prefix of a stored hash must not match it. String.equals gets this right
     * too, so the case exists to pin the length handling in MessageDigest.isEqual
     * rather than to catch a difference between the two.
     */
    @Test
    void containsHash_doesNotMatchOnAPrefix() {
        List<String> hashes = Arrays.asList("abcdef");

        assertFalse(SignatureVerifier.containsHash(hashes, "abc"));
        assertFalse(SignatureVerifier.containsHash(hashes, "abcdefgh"));
        assertTrue(SignatureVerifier.containsHash(hashes, "abcdef"));
    }
}
