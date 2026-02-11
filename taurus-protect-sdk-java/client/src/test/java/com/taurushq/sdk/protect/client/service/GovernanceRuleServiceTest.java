package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.helper.SignatureVerifier;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.GovernanceRules;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.RuleUserSignature;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.apache.commons.codec.binary.Base64;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
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
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class GovernanceRuleServiceTest {

    private static KeyPair keyPair1;
    private static KeyPair keyPair2;

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private List<PublicKey> superAdminKeys;

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

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        superAdminKeys = Collections.singletonList(keyPair1.getPublic());
    }

    // --- Constructor validation ---

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new GovernanceRuleService(null, apiExceptionMapper, superAdminKeys, 1));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new GovernanceRuleService(apiClient, null, superAdminKeys, 1));
    }

    @Test
    void constructor_throwsOnNullSuperAdminKeys() {
        assertThrows(NullPointerException.class, () ->
                new GovernanceRuleService(apiClient, apiExceptionMapper, null, 1));
    }

    @Test
    void constructor_throwsOnEmptySuperAdminKeys() {
        assertThrows(IllegalArgumentException.class, () ->
                new GovernanceRuleService(apiClient, apiExceptionMapper, Collections.emptyList(), 1));
    }

    @Test
    void constructor_throwsOnZeroMinSignatures() {
        assertThrows(IllegalArgumentException.class, () ->
                new GovernanceRuleService(apiClient, apiExceptionMapper, superAdminKeys, 0));
    }

    @Test
    void constructor_throwsOnNegativeMinSignatures() {
        assertThrows(IllegalArgumentException.class, () ->
                new GovernanceRuleService(apiClient, apiExceptionMapper, superAdminKeys, -1));
    }

    // --- Valid construction ---

    @Test
    void constructor_succeeds_withValidArgs() {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);
        assertNotNull(service);
    }

    @Test
    void getSuperAdminPublicKeys_returnsConfiguredKeys() {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);
        assertEquals(superAdminKeys, service.getSuperAdminPublicKeys());
    }

    @Test
    void getMinValidSignatures_returnsConfiguredValue() {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 2);
        assertEquals(2, service.getMinValidSignatures());
    }

    @Test
    void getSuperAdminPublicKeys_returnsUnmodifiableList() {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);
        assertThrows(UnsupportedOperationException.class, () ->
                service.getSuperAdminPublicKeys().add(keyPair2.getPublic()));
    }

    // --- verifyGovernanceRules ---

    @Test
    void verifyGovernanceRules_withValidSignature_returnsRules() throws Exception {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);

        byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);
        String rulesBase64 = Base64.encodeBase64String(rulesData);
        String sig = CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData);

        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(rulesBase64);
        RuleUserSignature ruleSig = new RuleUserSignature();
        ruleSig.setUserId("admin1");
        ruleSig.setSignature(sig);
        rules.setRulesSignatures(Collections.singletonList(ruleSig));

        GovernanceRules result = service.verifyGovernanceRules(rules, 1);
        assertNotNull(result);
        assertEquals(rulesBase64, result.getRulesContainer());
    }

    @Test
    void verifyGovernanceRules_withInvalidSignature_throws() throws Exception {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);

        byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);
        String rulesBase64 = Base64.encodeBase64String(rulesData);
        // Sign with keyPair2, but service only has keyPair1 as superAdmin key
        String sig = CryptoTPV1.calculateBase64Signature(keyPair2.getPrivate(), rulesData);

        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(rulesBase64);
        RuleUserSignature ruleSig = new RuleUserSignature();
        ruleSig.setUserId("admin1");
        ruleSig.setSignature(sig);
        rules.setRulesSignatures(Collections.singletonList(ruleSig));

        assertThrows(IntegrityException.class, () ->
                service.verifyGovernanceRules(rules, 1));
    }

    @Test
    void verifyGovernanceRules_withInsufficientSignatures_throws() throws Exception {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);

        byte[] rulesData = "governance rules payload".getBytes(StandardCharsets.UTF_8);
        String rulesBase64 = Base64.encodeBase64String(rulesData);
        String sig = CryptoTPV1.calculateBase64Signature(keyPair1.getPrivate(), rulesData);

        GovernanceRules rules = new GovernanceRules();
        rules.setRulesContainer(rulesBase64);
        RuleUserSignature ruleSig = new RuleUserSignature();
        ruleSig.setUserId("admin1");
        ruleSig.setSignature(sig);
        rules.setRulesSignatures(Collections.singletonList(ruleSig));

        // Require 2, only 1 valid
        assertThrows(IntegrityException.class, () ->
                service.verifyGovernanceRules(rules, 2));
    }

    // --- getDecodedRulesContainer ---

    @Test
    void getDecodedRulesContainer_throwsOnNullRules() {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);

        assertThrows(NullPointerException.class, () ->
                service.getDecodedRulesContainer(null));
    }

    // --- getRulesById validation ---

    @Test
    void getRulesById_throwsOnNullId() {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);

        assertThrows(IllegalArgumentException.class, () ->
                service.getRulesById(null));
    }

    @Test
    void getRulesById_throwsOnEmptyId() {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);

        assertThrows(IllegalArgumentException.class, () ->
                service.getRulesById(""));
    }

    // --- getRulesHistory validation ---

    @Test
    void getRulesHistory_throwsOnZeroPageSize() {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);

        assertThrows(IllegalArgumentException.class, () ->
                service.getRulesHistory(0));
    }

    @Test
    void getRulesHistory_throwsOnNegativePageSize() {
        GovernanceRuleService service = new GovernanceRuleService(
                apiClient, apiExceptionMapper, superAdminKeys, 1);

        assertThrows(IllegalArgumentException.class, () ->
                service.getRulesHistory(-1));
    }
}
