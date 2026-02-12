package com.taurushq.sdk.protect.client.helper;

import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.IntegrityException;
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
        assertTrue(ex.getMessage().contains("only 1 valid signatures found, minimum 2 required"));
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
        assertTrue(ex.getMessage().contains("only 0 valid signatures found, minimum 1 required"));
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
}
